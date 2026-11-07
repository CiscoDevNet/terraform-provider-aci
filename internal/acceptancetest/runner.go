package acceptancetest

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

const skipPrefix = "Skip reason: "

// Skip reports a reason that can be preserved in the parent run summary.
func Skip(t *testing.T, reason string) {
	t.Helper()
	// Preserve a recognizable reason in the structured child output without parsing arbitrary warnings.
	t.Skip(skipPrefix + reason)
}

// Runner owns worker processes and logs for one acceptance-test package.
// Its zero value is ready for use and must not be copied after first use.
type Runner struct {
	sync.Mutex
	directory string
	attempts  map[string]int
	context   context.Context
	stop      context.CancelFunc
	report    runReport
}

// RunInSubprocess handles a test in an isolated worker when parallel mode is
// enabled. A true return means the caller must return without running its body.
func (run *Runner) RunInSubprocess(t *testing.T) bool {
	t.Helper()
	if os.Getenv("TF_ACC") == "" || os.Getenv("ACI_TEST_PARALLEL") != "1" || os.Getenv("ACI_TEST_WORKER") != "" {
		return false
	}
	// Register before queueing so an interrupted report includes workers that have not started yet.
	deadline, _ := t.Deadline()
	child, ctx, err := run.prepareWorker(t.Name(), deadline)
	if err != nil {
		t.Fatalf("Prepare acceptance-test logs: %v", err)
	}
	// Go's scheduler limits active workers using the parent test command's -parallel setting.
	t.Parallel()
	result := workerResult{Action: "fail", LogPath: child.LogPath}
	if err := ctx.Err(); err != nil {
		result.Error = fmt.Sprintf("Not started: test run stopped (%v)", err)
	} else {
		t.Logf("Detailed acceptance-test log: %s", child.LogPath)
		result = child.run(ctx)
	}
	run.recordResult(t.Name(), result)
	switch result.Action {
	case "skip":
		t.Skip(result.SkipReason)
	case "fail":
		t.Errorf("Acceptance test failed: %s; see %s", result.Error, result.LogPath)
	}
	return true
}

// Finish finalizes the report after the package's tests return, including failed runs.
func (run *Runner) Finish() error {
	run.Lock()
	defer run.Unlock()
	if run.directory == "" || !run.report.Finished.IsZero() {
		return nil
	}
	if run.stop != nil {
		defer run.stop()
	}
	run.report.Finished = time.Now()
	if run.context != nil {
		switch run.context.Err() {
		case context.DeadlineExceeded:
			run.report.IncompleteReason = "The test run did not complete because its time limit was reached."
		case context.Canceled:
			run.report.IncompleteReason = "The test run did not complete because it was interrupted."
		}
	}
	return run.report.write(run.directory)
}

func (run *Runner) prepareWorker(name string, deadline time.Time) (worker, context.Context, error) {
	run.Lock()
	defer run.Unlock()
	if run.directory == "" {
		base, err := logDirectory()
		if err != nil {
			return worker{}, nil, err
		}
		if err := os.MkdirAll(base, 0700); err != nil {
			return worker{}, nil, err
		}
		run.directory, err = os.MkdirTemp(base, time.Now().UTC().Format("20060102T150405-"))
		if err != nil {
			return worker{}, nil, err
		}
		run.attempts = make(map[string]int)
		run.report.Started = time.Now()
		// Cancellation is shared by the parent only; newly scheduled workers also see an interrupt.
		run.context, run.stop = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		if !deadline.IsZero() {
			// Leave time for workers to stop and Finish to write the report before Go's hard timeout.
			// Short runs reserve 10% of their remaining time instead of the full 30 seconds.
			stopSignals := run.stop
			grace := min(30*time.Second, time.Until(deadline)/10)
			var stopDeadline context.CancelFunc
			run.context, stopDeadline = context.WithDeadline(run.context, deadline.Add(-grace))
			run.stop = func() {
				stopDeadline()
				stopSignals()
			}
		}
	}
	run.attempts[name]++
	suffix := fmt.Sprintf("-%03d", run.attempts[name])
	filename := strings.ReplaceAll(name, "/", "__")
	child := worker{
		Name:    name,
		LogPath: filepath.Join(run.directory, filename+suffix+".log"),
		Timeout: 30 * time.Minute,
	}
	if os.Getenv("TF_LOG") != "" || os.Getenv("TF_ACC_LOG") != "" || os.Getenv("TF_LOG_PATH") != "" || os.Getenv("TF_ACC_LOG_PATH") != "" || os.Getenv("TF_LOG_PATH_MASK") != "" {
		child.TraceMask = filepath.Join(run.directory, "%s"+suffix+".trace.log")
	}
	run.report.Scheduled++
	return child, run.context, nil
}

func (run *Runner) recordResult(name string, result workerResult) {
	run.Lock()
	defer run.Unlock()
	run.report.Results = append(run.report.Results, reportResult{Name: name, workerResult: result})
}

func logDirectory() (string, error) {
	if directory := os.Getenv("ACI_TEST_LOG_DIR"); directory != "" {
		return filepath.Abs(directory)
	}
	directory, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(directory, "go.mod")); err == nil {
			return filepath.Join(directory, "test-logs"), nil
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			return "", fmt.Errorf("cannot find repository root for default test logs")
		}
		directory = parent
	}
}
