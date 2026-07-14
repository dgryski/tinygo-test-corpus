package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
	wasip1 := flag.Bool("wasip1", false, "run tests on wasip1")
	wasip2 := flag.Bool("wasip2", false, "run tests on wasip2")
	wizer := flag.Bool("wizer", false, "run wizer on wasi output")
	keepGoing := flag.Bool("k", false, "keep going after a failed test, logging all failures at the end")
	verbose := flag.Bool("v", false, "print verbose output")
	noUpdate := flag.Bool("no-update", false, "don't update repos")

	flag.Parse()

	repos, err := loadRepos(*configYaml)
	if err != nil {
		log.Fatalf("unable to load repo configuration: %v", err)
	}

	var countSubdir, countRepo int
	defer func() {
		log.Printf("Finished!\n%d/%d repos tested\n%d passed subdir tests\n", countRepo, len(repos), countSubdir)
	}()

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
	if *wasip1 {
		target = "wasip1"
		os.Setenv("WASMTIME_BACKTRACE_DETAILS", "1")
	}

	if *wasip2 {
		target = "wasip2"
		os.Setenv("WASMTIME_BACKTRACE_DETAILS", "1")
	}

	// Workspace setup and cleanup.
	goos := newCommander(*parallelism)
	goos.Run(*compiler, "clean")

	goos.Mkdir(corpusFolderName, dirMode) // force directory creation if not exist.
	_, err = goos.Stat(corpusFolderName)
	if err != nil {
		log.Fatal("reading corpus directory: ", err)
	}
	goos.Chdir(corpusFolderName)
	corpusDir := goos.path

	// Commence testing logic.
	for _, repo := range repos {
		if (*wasip1 || *wasip2) && repo.SkipWASI == "true" {
			log.Printf("skipping non-wasi package %v", repo.Repo)
			continue
		}

		goos.Chdir(corpusDir)
		if !*noUpdate {
			goos.cloneOrUpdateRepo(repo.Repo)
		}
		repoBase := filepath.Join(corpusDir, repo.Repo)
		goos.Chdir(repoBase)

		if _, err := goos.Stat("go.mod"); err != nil {
			log.Printf("creating %s/go.mod: running `go mod init`\n", repoBase)
			goos.Run("go", "mod", "init", fmt.Sprintf("%s/%s", host, repo.Repo))
			goos.Run("go", "get", "-t", ".")
		}
		tags := "runtime_asserts"
		if repo.Tags != "" {
			tags += " " + repo.Tags
		}

		var dirs []Subdir
		if repo.Skip != "true" {
			dirs = append(dirs, Subdir{Pkg: repo.Repo, Skip: repo.Skip, SkipWASI: repo.SkipWASI})
		}
		if len(repo.Subdirs) > 0 {
			dirs = append(dirs, repo.Subdirs...)
		}

		for _, subdir := range dirs {
			if (*wasip1 || *wasip2) && subdir.SkipWASI == "true" {
				log.Printf("skipping non-wasi package %v/%v", repo.Repo, subdir.Pkg)
				continue
			}
			if subdir.Pkg != repo.Repo {
				goos.Chdir(subdir.Pkg)
			}
			cmd := []string{"test", "-target=" + target, "-gc=precise", "-tags=" + tags}

			var skips []string
			if subdir.Skip != "" {
				skips = append(skips, subdir.Skip)
			}
			if subdir.SkipWASI != "" && (*wasip1 || *wasip2) {
				skips = append(skips, subdir.SkipWASI)
			}
			if len(skips) > 0 {
				cmd = append(cmd, "-skip="+strings.Join(skips, "|"))
			}
			if *wizer {
				cmd = append(cmd, "-wizer-init")
			}
			if *keepGoing {
				goos.StartNonFatal(*verbose, *compiler, cmd...)
			} else {
				goos.Start(*verbose, *compiler, cmd...)
			}
			countSubdir++
			goos.Chdir(repoBase)
		}
		countRepo++
		log.Printf("finished module %d/%d %s", countRepo, len(repos), repo.Repo)
	}
	goos.Wait()
	if f := goos.fails(); f > 0 {
		log.Fatalf("%d failures running tests.", f)
	}
}

type T struct {
	Repo     string
	Tags     string
	Subdirs  []Subdir
	SkipWASI string
	Skip     string
}

type Subdir struct {
	Pkg      string
	SkipWASI string
	Skip     string
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
