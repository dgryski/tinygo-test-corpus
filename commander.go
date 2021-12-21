package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// commander keeps track of a directory and runs commands on them providing
// a crude form of parallelism. Not safe for concurrent use.
type commander struct {
	// path is guaranteed to be the absolute current path of commander.
	path string
	// checkin is a buffered channel. It's length limits the amount of goroutines running commands at once.
	checkin chan *exec.Cmd
	wg      *sync.WaitGroup
}

func newCommander(goroutines int) commander {
	path, err := os.Getwd()
	if err != nil {
		log.Panic("commander: getting current dir:", err)
	}
	if goroutines < 1 {
		log.Panic("commander: invalid number of goroutines argument")
	}
	log.Println("setting max simultaneous commands to", goroutines)
	return commander{
		path:    path,
		checkin: make(chan *exec.Cmd, goroutines),
		wg:      new(sync.WaitGroup),
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
		c.Run("git", "clone", fmt.Sprintf("%s/%s", hostURL, repo))
		return
	}

	c.Chdir(repo)
	log.Printf("repo exists, updating %s", repo)
	c.Run("git", "fetch")
	c.Run("git", "pull")
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

// Start begins command execution in commander's current directory and returns immediately.
// Prints command output on finish. Calls os.Exit(1) on error.
func (r *commander) Start(name string, arg ...string) {
	r.run(true, name, arg...)
}

// Run executes command in commander's current directory and waits for it to finish.
// Prints command output. If command returns non-zero exit code then result is logged
// and os.Exit(1) is called.
func (r *commander) Run(name string, arg ...string) {
	r.run(false, name, arg...)
}

func (r *commander) run(async bool, name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	r.checkin <- cmd // Check-in for work.

	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	cmd.Dir = r.path
	err := cmd.Start()
	if err != nil {
		log.Panicf("%s\ncmd %s with err: %q at dir %q", b.String(), cmd.String(), err, r.path)
	}

	done := make(chan struct{})
	r.wg.Add(1)
	go func() {
		if async {
			done <- struct{}{}
		}

		err = cmd.Wait()
		if err != nil {
			log.Panicf("%s\ncmd %s with err: %v at dir %q", b.String(), cmd.String(), err, cmd.Dir)
		}
		log.Printf("cmd %s finished with output:\n%s", cmd, b.String())
		<-r.checkin // Check-out
		if !async {
			done <- struct{}{}
		}
		r.wg.Done()
	}()
	<-done
}

// Wait blocks execution until there are no more commands being executed asynchronously
// by commander.
func (r *commander) Wait() {
	r.wg.Wait()
}

func (r *commander) terminate() {
	for len(r.checkin) != 0 {
		cmd := <-r.checkin
		cmd.Process.Kill()
	}
	close(r.checkin)
}
