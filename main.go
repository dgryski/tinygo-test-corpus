package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	if strings.HasSuffix(*compiler, "tinygo") {
		mustrun(*compiler, "clean")
	}
	if err != nil {
		log.Fatalf("calling `%v clean`: %s", *compiler, err)
	}
	os.Mkdir(corpusDir, dirMode) // force directory creation if not exist.
	_, err = os.ReadDir(corpusDir)
	if err != nil {
		log.Fatal("reading corpus directory: ", err)
	}

	// Commence testing logic. Start from latest repo additions (end of repos).
	oos := newCommander()
	for i := len(repos) - 1; i >= 0; i-- {
		repo := repos[i]
		oos.Chdir(corpusDir)
		oos.cloneOrUpdateRepo(repo.Repo)
		repoBase := filepath.Join(corpusDir, repo.Repo)
		oos.Chdir(repoBase)

		if _, err := oos.Stat("go.mod"); err != nil {
			log.Printf("creating %s/go.mod: running `go mod init`\n", repoBase)
			oos.run(false, "go", "mod", "init", fmt.Sprintf("%s/%s", host, repo.Repo))
			oos.run(false, "go", "get", "-t", ".")
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
				oos.Chdir(subdir)
			}
			oos.run(true, *compiler, "test", "-v", "-tags="+tags)
			if subdir != "." {
				oos.Chdir(repoBase)
			}
		}
		countRepo++
		log.Printf("finished module %d/%d %s", countRepo, len(repos), repo.Repo)
	}
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
