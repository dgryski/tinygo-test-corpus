package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// job is a command queued for execution by a commander worker.
type job struct {
	// dir is the directory the command runs in, captured when the job is queued
	// since the commander's current directory changes while jobs are pending.
	dir     string
	name    string
	args    []string
	fatal   bool
	verbose bool
	// tag identifies the module the job belongs to.
	tag string
}

// result is the outcome of a finished job.
type result struct {
	job    job
	output string
	err    error
}

func (j job) String() string {
	return exec.Command(j.name, j.args...).String()
}

// commander keeps track of a directory and runs commands on them providing
// a crude form of parallelism. Not safe for concurrent use.
type commander struct {
	// path is guaranteed to be the absolute current path of commander.
	path string
	// checkin carries jobs to the worker goroutines.
	checkin chan job
	// checkout contains the "finished" results for a run
	checkout chan result
	// onStart, if set, is called from the worker goroutine as it picks up a job.
	// It must be set before the first job is submitted and be safe for concurrent use.
	onStart func(job)
}

func newCommander(goroutines int) *commander {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal("commander: getting current dir:", err)
	}
	if goroutines < 1 {
		log.Fatal("commander: invalid number of goroutines argument")
	}
	log.Println("setting max simultaneous commands to", goroutines)
	c := &commander{
		path:     path,
		checkin:  make(chan job),
		checkout: make(chan result, goroutines),
	}
	for i := 0; i < goroutines; i++ {
		go c.worker()
	}
	return c
}

func (c *commander) worker() {
	for j := range c.checkin {
		if c.onStart != nil {
			c.onStart(j)
		}
		c.checkout <- execute(j)
	}
}

func execute(j job) result {
	cmd := exec.Command(j.name, j.args...)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	cmd.Dir = j.dir
	err := cmd.Run()
	return result{job: j, output: b.String(), err: err}
}

// Job returns a job to run name with args in commander's current directory.
func (c *commander) Job(tag string, fatal, verbose bool, name string, arg ...string) job {
	return job{dir: c.path, name: name, args: arg, fatal: fatal, verbose: verbose, tag: tag}
}

// Submit queues j for a worker. While waiting for a free worker it calls
// handle on any results that finish in the meantime.
func (c *commander) Submit(j job, handle func(result)) {
	for {
		select {
		case c.checkin <- j:
			return
		case r := <-c.checkout:
			handle(r)
		}
	}
}

// Drain calls handle on pending results from checkout and then stops the
// workers. pending is the number of submitted jobs not yet handled.
func (c *commander) Drain(pending int, handle func(result)) {
	for ; pending > 0; pending-- {
		handle(<-c.checkout)
	}
	close(c.checkin)
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

func (r *commander) Stat(path string) (os.FileInfo, error) {
	return os.Stat(r.join(path))
}
func (r *commander) Mkdir(d string, mode fs.FileMode) error {
	return os.MkdirAll(r.join(d), mode)
}
func (r *commander) Chdir(path string) {
	r.path = r.join(path)
}
func (r *commander) join(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(r.path, path)
}

// Run executes command in commander's current directory and waits for it to finish.
// Prints command output. If command returns non-zero exit code then result is logged
// and os.Exit(1) is called. It does not use the worker pool.
func (r *commander) Run(name string, arg ...string) {
	res := execute(r.Job("", true, false, name, arg...))
	if res.err != nil {
		log.Fatalln(res.failure())
	}
	log.Printf("%s\ncmd %s at dir %q", res.output, res.job, res.job.dir)
}

func (r result) failure() string {
	return fmt.Sprintf("%s\ncmd %s with err: %v at dir %q", r.output, r.job, r.err, r.job.dir)
}
