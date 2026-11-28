package acceptancetest

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type worker struct {
	Name        string
	LogPath     string
	TraceMask   string
	Timeout     time.Duration
	Environment []string
}

type workerResult struct {
	Action     string
	SkipReason string
	LogPath    string
	TracePath  string
	Error      string
	Elapsed    time.Duration
}

func (child worker) run(ctx context.Context) (result workerResult) {
	started := time.Now()
	result = workerResult{Action: "fail", LogPath: child.LogPath}
	defer func() {
		result.Elapsed = time.Since(started)
		if child.TraceMask != "" {
			// Match terraform-plugin-testing's substitution, and link only traces that were actually written.
			tracePath := fmt.Sprintf(child.TraceMask, strings.ReplaceAll(child.Name, "/", "__"))
			if info, err := os.Stat(tracePath); err == nil && !info.IsDir() {
				result.TracePath = tracePath
			}
		}
	}()
	logFile, err := os.OpenFile(child.LogPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		result.Error = fmt.Sprintf("Open detailed test log: %v", err)
		return result
	}
	defer logFile.Close()
	fail := func(message string) workerResult {
		result.Error = message
		fmt.Fprintln(logFile, message)
		return result
	}
	// Give the child its own alarm and leave a short interval for its final output to reach the parent.
	ctx, cancel := context.WithTimeout(ctx, child.Timeout+2*time.Second)
	defer cancel()
	command, err := child.buildCommand(ctx)
	if err != nil {
		return fail(err.Error())
	}
	command.Stderr = logFile
	stdout, err := command.StdoutPipe()
	if err != nil {
		return fail(fmt.Sprintf("Capture child test output: %v", err))
	}
	if err := command.Start(); err != nil {
		return fail(fmt.Sprintf("Start child test: %v", err))
	}
	defer stopAcceptanceProcess(command)
	events := workerEvents{Name: child.Name}
	decodeError := events.read(stdout, logFile)
	if decodeError != nil {
		stopAcceptanceProcess(command)
	}
	// Reap the worker before reporting errors, preserving the original read/write failure.
	waitError := command.Wait()
	if decodeError != nil {
		return fail(fmt.Sprintf("Read child test results: %v", decodeError))
	}
	if waitError != nil {
		return fail(fmt.Sprintf("Child test exited unsuccessfully: %v", waitError))
	}
	if !events.Ran || events.Action == "" {
		return fail("The intended acceptance test did not run or did not produce a terminal result")
	}
	result.Action = events.Action
	if events.Action == "skip" {
		result.SkipReason = events.SkipReason
		if result.SkipReason == "" {
			// Framework skips outside PreCheck still retain their emitted diagnostic.
			result.SkipReason = events.lastOutput
		}
	}
	return result
}

func (child worker) buildCommand(ctx context.Context) (*exec.Cmd, error) {
	executable, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("Find test executable: %w", err)
	}
	if child.Timeout <= 0 {
		return nil, fmt.Errorf("No time remains to execute the acceptance test")
	}

	// Reuse the compiled binary for exactly one test and convert its output into JSON events.
	command := exec.CommandContext(ctx, "go", "tool", "test2json", "-t", executable,
		"-test.run=^"+regexp.QuoteMeta(child.Name)+"$", "-test.v=test2json", "-test.count=1", "-test.parallel=1",
		"-test.timeout="+child.Timeout.String())
	overrides := append([]string(nil), child.Environment...)
	// The marker makes the selected test execute its body instead of starting another worker.
	overrides = append(overrides, "ACI_TEST_WORKER="+child.Name)
	if child.TraceMask != "" {
		// Each process owns its logger, so legacy log.Print and framework traces share this test's trace file.
		overrides = append(overrides, "TF_LOG_PATH=", "TF_ACC_LOG_PATH=", "TF_LOG_PATH_MASK="+child.TraceMask)
	}
	command.Env = childEnvironment(overrides...)
	configureAcceptanceProcess(command)
	// Cancel the whole process tree when the shared context or worker deadline expires.
	command.Cancel = func() error { return stopAcceptanceProcess(command) }
	command.WaitDelay = 2 * time.Second
	return command, nil
}

func childEnvironment(overrides ...string) []string {
	// Modify only the child environment; parent tests can be scheduled concurrently.
	environment := os.Environ()
	for _, override := range overrides {
		name, _, _ := strings.Cut(override, "=")
		filtered := environment[:0]
		for _, existing := range environment {
			if !strings.HasPrefix(existing, name+"=") {
				filtered = append(filtered, existing)
			}
		}
		environment = append(filtered, override)
	}
	return environment
}
