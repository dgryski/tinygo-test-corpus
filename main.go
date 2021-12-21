package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"
)

const (
	// underscore prefix so go tool excludes corpus directory.
	corpusFolderName = "_corpus"
	dirMode          = 0777
	host             = "github.com"
	hostURL          = "https://" + host
)

func main() {
	configYaml := flag.String("config", "repos.yaml", "yaml of repositories to run")
	compiler := flag.String("compiler", "tinygo", "use this go compiler")
	runPattern := flag.String("run", "", "compiler will run on all repo names matching this pattern (regexp)")
	parallelism := flag.Int("parallel", 2, "max number of goroutines running compiler at any time")
	wasi := flag.Bool("wasi", false, "run tests on wasi")

	flag.Parse()

	repos, err := loadRepos(*configYaml)
	if err != nil {
		log.Fatalf("unable to load repo configuration: %v", err)
	}

	// Which repos to run.
	re, err := regexp.Compile(*runPattern)
	if err != nil {
		log.Fatal("compiling run regexp:", err)
	}
	if *runPattern != "" {
		// pre-filter the repos list so the progress indicator is correct
		var filtered []T
		for _, repo := range repos {
			if !re.MatchString(repo.Repo) {
				continue
			}
			filtered = append(filtered, repo)
		}
		repos = filtered
	}

	var target string
	if *wasi {
		target = "wasi"
		os.Setenv("WASMTIME_BACKTRACE_DETAILS", "1")
	}

	var countSubdir, countRepo int
	defer func() {
		log.Printf("%d/%d repos tested\n%d passed subdir tests\n", countRepo, len(repos), countSubdir)
		err := recover()
		if err != nil {
			log.Fatalf("Fatal error encountered: %v", err)
		} else {
			log.Println("Finished succesfully")
		}
	}()

	// Workspace setup and cleanup.
	goos := newCommander(*parallelism)
	defer goos.terminate()
	goos.Run(*compiler, "clean")

	goos.Mkdir(corpusFolderName, dirMode) // force directory creation if not exist.
	_, err = goos.Stat(corpusFolderName)
	if err != nil {
		log.Panic("reading corpus directory: ", err)
	}
	goos.Chdir(corpusFolderName)
	corpusDir := goos.path

	// Commence testing logic.
	for _, repo := range repos {
		if *wasi && repo.SkipWASI {
			log.Printf("skipping non-wasi package %v", repo.Repo)
			continue
		}

		goos.Chdir(corpusDir)
		goos.cloneOrUpdateRepo(repo.Repo)
		repoBase := filepath.Join(corpusDir, repo.Repo)
		goos.Chdir(repoBase)

		if _, err := goos.Stat("go.mod"); err != nil {
			log.Printf("creating %s/go.mod: running `go mod init`\n", repoBase)
			goos.Run("go", "mod", "init", fmt.Sprintf("%s/%s", host, repo.Repo))
			goos.Run("go", "get", "-t", ".")
		}
		tags := ""
		if repo.Tags != "" {
			tags = repo.Tags
		}

		dirs := []Subdir{Subdir{Pkg: repo.Repo}}
		if len(repo.Subdirs) > 0 {
			dirs = repo.Subdirs
		}

		for _, subdir := range dirs {
			if *wasi && subdir.SkipWASI {
				log.Printf("skipping non-wasi package %v/%v", repo.Repo, subdir.Pkg)
				continue
			}
			if subdir.Pkg != repo.Repo {
				goos.Chdir(subdir.Pkg)
			}
			goos.Start(*compiler, "test", "-v", "-target="+target, "-tags="+tags)
			countSubdir++
			goos.Chdir(repoBase)
		}
		countRepo++
		log.Printf("finished module %d/%d %s", countRepo, len(repos), repo.Repo)
	}
	goos.Wait()

}

type T struct {
	Repo     string
	Tags     string
	Subdirs  []Subdir
	SkipWASI bool
}

type Subdir struct {
	Pkg      string
	SkipWASI bool
}

func loadRepos(f string) ([]T, error) {

	yf, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	var repos []T
	err = yaml.Unmarshal(yf, &repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}
