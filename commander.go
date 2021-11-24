package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

// commander keeps track of a directory and runs commands on them providing
// a crude form of parallelism.
type commander struct {
	*sync.Mutex
	path string
	// checkin is a buffered channel. It's length limits the amount of goroutines running commands at once.
	checkin chan struct{}
}

func newCommander() commander {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal("getting current dir:", err)
	}
	goroutines := runtime.NumCPU() / 2
	if goroutines == 0 {
		goroutines = 1
	}
	log.Println("setting max simultaneous commands to", goroutines)
	return commander{
		path:    path,
		checkin: make(chan struct{}, goroutines),
		Mutex:   &sync.Mutex{},
	}
}

func (c *commander) cloneOrUpdateRepo(repo string) {
	if _, err := c.Stat(repo); err != nil {
		// Repo does not exist.
		log.Printf("repo not found. cloning %s", repo)
		d := filepath.Dir(repo)
		if _, err := c.Stat(repo); err != nil {
			log.Printf("creating directory %s", d)
			c.Mkdir(d, dirMode)
		}
		c.Chdir(d)
		c.run(false, "git", "clone", fmt.Sprintf("%s/%s", hostURL, repo))
		return
	}

	c.Chdir(repo)
	log.Printf("repo exists, updating %s", repo)
	c.run(false, "git", "fetch")
	c.run(false, "git", "pull")
}

func (r commander) Stat(path string) (os.FileInfo, error) {
	return os.Stat(r.join(path))
}
func (r commander) Mkdir(d string, mode fs.FileMode) error {
	return os.MkdirAll(r.join(d), mode)
}
func (r *commander) Chdir(path string) {
	r.path = r.join(path)
}
func (r commander) join(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(r.path, path)
}

func (r *commander) run(async bool, name string, arg ...string) {
	r.checkin <- struct{}{} // Check-in for work.
	cmd := exec.Command(name, arg...)
	r.Lock() // protect Chdir and cmd.Start
	defer r.Unlock()
	path := r.path
	err := os.Chdir(path)
	cwd, _ := os.Getwd()
	if err != nil {
		log.Fatal("os.Chdir error: ", err)
	}
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	err = cmd.Start()
	if err != nil {
		log.Fatalf("%s\ncmd %s with err: %q at dir %q", b.String(), cmd.String(), err, cwd)
	}

	err = os.Chdir(r.path)
	if err != nil {
		log.Fatalf("changing directory to %s, from %s: %s", r.path, path, err)
	}
	done := make(chan struct{})
	go func() {
		if async {
			done <- struct{}{}
		}
		err = cmd.Wait()
		if err != nil {
			log.Fatalf("%s\ncmd %s with err: %v at dir %q", b.String(), cmd.String(), err, path)
		}
		log.Printf("cmd %s finished with output:\n%s", cmd, b.String())
		<-r.checkin // Check-out
		if !async {
			done <- struct{}{}
		}
	}()
	<-done
}
