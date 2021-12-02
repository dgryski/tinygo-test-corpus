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
		if !re.MatchString(repo.Repo) {
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
		dirs := []string{"."}
		if len(repo.Subdirs) > 0 {
			dirs = repo.Subdirs
		}

		for _, subdir := range dirs {
			if subdir != "." {
				goos.Chdir(subdir)
			}
			goos.Start(*compiler, "test", "-v", "-tags="+tags)
			countSubdir++
			if subdir != "." {
				goos.Chdir(repoBase)
			}
		}
		countRepo++
		log.Printf("finished module %d/%d %s", countRepo, len(repos), repo.Repo)
	}
}

type T struct {
	Repo    string
	Tags    string
	Subdirs []string
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
