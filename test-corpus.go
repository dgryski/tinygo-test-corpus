package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

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
	flag.Parse()

	repos, err := loadRepos(*configYaml)
	if err != nil {
		log.Fatalf("unable to load repo configuration: %v", err)
	}

	var countSubdir, countRepo int
	defer func() {
		log.Printf("Finished!\n%d/%d repos tested\n%d passed subdir tests\n", countRepo, len(repos), countSubdir)
	}()

	// Workspace setup and cleanup.
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal("getting current dir:", err)
	}
	corpusDir := filepath.Join(baseDir, corpusFolderName)
	mustrun("tinygo", "clean")
	if err != nil {
		log.Fatal("calling `tinygo clean`:", err)
	}
	os.Mkdir(corpusDir, dirMode) // force directory creation if not exist.
	_, err = os.ReadDir(corpusDir)
	if err != nil {
		log.Fatal("reading corpus directory: ", err)
	}

	// Commence testing logic.
	for _, repo := range repos {
		os.Chdir(corpusDir)
		cloneOrUpdateRepo(repo.Repo)
		repoBase := filepath.Join(corpusDir, repo.Repo)
		os.Chdir(repoBase)

		if _, err := os.Stat("go.mod"); err != nil {
			log.Printf("creating %s/go.mod: running `go mod init`\n", repoBase)
			mustrun("go", "mod", "init", fmt.Sprintf("%s/%s", host, repo.Repo))
			mustrun("go", "get", "-t", ".")
		}
		tags := ""
		if repo.Tags != "" {
			tags = fmt.Sprintf("%s", repo.Tags)
		}
		dirs := []string{"."}
		if len(repo.Subdirs) > 0 {
			dirs = repo.Subdirs
		}

		for _, subdir := range dirs {
			if subdir != "." {
				os.Chdir(subdir)
			}
			tinyout := make(chan string)
			// Run TinyGo and Go in parallel.
			go func() {
				tinyout <- mustrun("tinygo", "test", "-v", "-tags", tags)
			}()
			out1 := mustrun("go", "test", "-v")
			countSubdir++
			log.Printf("package %s:\n%s\n%s\n", filepath.Join(repo.Repo, subdir), out1, <-tinyout)
			if subdir != "." {
				os.Chdir(repoBase)
			}
		}
		countRepo++
		log.Printf("finished module %d/%d %s", countRepo, len(repos), repo.Repo)
	}
}

func cloneOrUpdateRepo(repo string) {
	if _, err := os.Stat(repo); err != nil {
		// Repo does not exist.
		log.Printf("repo not found. cloning %s", repo)
		d := filepath.Dir(repo)
		if _, err := os.Stat(repo); err != nil {
			log.Printf("creating directory %s", d)
			os.Mkdir(d, dirMode)
		}
		os.Chdir(d)
		mustrun("git", "clone", fmt.Sprintf("%s/%s", hostURL, repo))
		return
	}

	os.Chdir(repo)
	log.Printf("repo exists, updating %s", repo)
	mustrun("git", "fetch")
	mustrun("git", "pull")
}

func mustrun(name string, arg ...string) (stdout string) {
	cmd := exec.Command(name, arg...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		cwd, _ := os.Getwd()
		log.Fatalf("%s\ncmd %s with err: %q at dir %q", string(b), cmd.String(), err, cwd)
	}
	return string(b)
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
