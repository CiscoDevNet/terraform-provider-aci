package acceptancetest

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type acceptanceGuardInput struct {
	Acceptance string
	Parallel   string
	Worker     string
}

func TestDelayAfterTest(t *testing.T) {
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_zero_delay", Input: time.Duration(0)},
		{Name: "test_delay_follows_other_cleanup", Input: 20 * time.Millisecond},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			delay := testCase.Input.(time.Duration)
			var cleanedUp time.Time
			t.Run("fixture", func(t *testing.T) {
				DelayAfterTest(t, delay)
				t.Cleanup(func() { cleanedUp = time.Now() })
			})
			if assert.False(t, cleanedUp.IsZero()) {
				assert.GreaterOrEqual(t, time.Since(cleanedUp), delay)
			}
		})
	}
}

func TestAcceptanceSubprocessGuards(t *testing.T) {
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_acceptance_disabled", Input: acceptanceGuardInput{Parallel: "1"}},
		{Name: "test_parallel_mode_disabled", Input: acceptanceGuardInput{Acceptance: "1"}},
		{Name: "test_worker_does_not_recurse", Input: acceptanceGuardInput{Acceptance: "1", Parallel: "1", Worker: "worker"}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(acceptanceGuardInput)
			t.Setenv("TF_ACC", input.Acceptance)
			t.Setenv("ACI_TEST_PARALLEL", input.Parallel)
			t.Setenv("ACI_TEST_WORKER", input.Worker)
			var runner Runner
			assert.False(t, runner.RunInSubprocess(t))
			assert.False(t, runner.RunSeriallyInSubprocess(t))
		})
	}
}

func TestAcceptanceRunnerScheduling(t *testing.T) {
	test.InitializeTest(t)

	directory := t.TempDir()
	executable, err := os.Executable()
	if !assert.NoError(t, err) {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// A separate test process can opt into parallel workers without changing the unit suite's environment.
	command := exec.CommandContext(ctx, executable, "-test.run=^TestAcceptanceRunnerSchedulingFixture$", "-test.v", "-test.parallel=2", "-test.timeout=15s")
	command.Env = childEnvironment("TF_ACC=1", "ACI_TEST_PARALLEL=1", "ACI_TEST_WORKER=", "ACI_TEST_SCHEDULING_DIRECTORY="+directory, "ACI_TEST_LOG_DIR="+filepath.Join(directory, "logs"))
	configureAcceptanceProcess(command)
	command.Cancel = func() error { return stopAcceptanceProcess(command) }
	command.WaitDelay = 2 * time.Second
	output, err := command.CombinedOutput()
	if !assert.NoError(t, err, string(output)) {
		return
	}
	reports, err := filepath.Glob(filepath.Join(directory, "logs", "*", "summary.md"))
	if assert.NoError(t, err) && assert.Len(t, reports, 1) {
		contents, err := os.ReadFile(reports[0])
		if assert.NoError(t, err) {
			assert.Contains(t, string(contents), "**4 of 4** test invocations")
			assert.Contains(t, string(contents), "**3 passed**, **0 failed**, and **1 skipped**")
		}
	}
}

func TestAcceptanceRunnerSchedulingFixture(t *testing.T) {
	// These workers exercise scheduling and reports locally; they never contact APIC.
	directory := os.Getenv("ACI_TEST_SCHEDULING_DIRECTORY")
	if directory == "" {
		t.Skip("local scheduling fixture")
	}
	runner := Runner{}
	serialFinished := filepath.Join(directory, "serial_finished")
	testCases := []test.TestCase{
		{Name: "parallel_before", Input: false, Expected: "parallel_after"},
		{Name: "serial", Input: true},
		{Name: "serial_skipped", Input: true},
		{Name: "parallel_after", Input: false, Expected: "parallel_before"},
	}
	t.Run("cases", func(t *testing.T) {
		for _, testCase := range testCases {
			runSerially := testCase.Input.(bool)
			t.Run(testCase.Name, func(t *testing.T) {
				if runSerially {
					if runner.RunSeriallyInSubprocess(t) {
						return
					}
					if testCase.Name == "serial_skipped" {
						// A missing skip guard would exceed the outer process timeout.
						DelayAfterTest(t, time.Hour)
						Skip(t, "local skipped fixture")
					}
					DelayAfterTest(t, 20*time.Millisecond)
					t.Cleanup(func() {
						assert.NoError(t, os.WriteFile(serialFinished, nil, 0600))
					})
					return
				}
				if runner.RunInSubprocess(t) {
					return
				}
				assert.FileExists(t, serialFinished)
				if !assert.NoError(t, os.WriteFile(filepath.Join(directory, testCase.Name), nil, 0600)) {
					return
				}
				// Both parallel workers must be active together; sequential execution cannot satisfy this barrier.
				assert.Eventually(t, func() bool {
					_, err := os.Stat(filepath.Join(directory, testCase.Expected.(string)))
					return err == nil
				}, 5*time.Second, 10*time.Millisecond)
			})
			if runSerially && os.Getenv("ACI_TEST_WORKER") == "" {
				// Unlike t.Parallel, a serial invocation must finish before t.Run returns to its parent.
				assert.FileExists(t, serialFinished)
			}
		}
	})
	assert.NoError(t, runner.Finish())
}

func TestAcceptanceSummary(t *testing.T) {
	test.InitializeTest(t)

	t.Setenv("ACI_TEST_LOG_DIR", t.TempDir())
	run := Runner{}
	children := make([]worker, 0, 12)
	for index := range 12 {
		child, _, err := run.prepareWorker(fmt.Sprintf("TestLocal%02d", index), time.Time{})
		if !assert.NoError(t, err) {
			return
		}
		children = append(children, child)
	}
	defer run.stop()
	var workers sync.WaitGroup
	for index, child := range children {
		workers.Add(1)
		go func() {
			defer workers.Done()
			run.recordResult(child.Name, workerResult{
				Action: "skip", SkipReason: fmt.Sprintf("reason %d", index),
				LogPath: child.LogPath, Elapsed: time.Second,
			})
		}()
	}
	workers.Wait()
	assert.NoFileExists(t, filepath.Join(run.directory, "summary.md"), "collecting results must not write the report")
	assert.NoError(t, run.Finish())
	contents, err := os.ReadFile(filepath.Join(run.directory, "summary.md"))
	if assert.NoError(t, err) {
		assert.Contains(t, string(contents), "**12 of 12** test invocations")
		assert.Contains(t, string(contents), "**0 passed**, **0 failed**, and **12 skipped**")
		assert.Contains(t, string(contents), "Overall result: **PASSED**")
		for index, child := range children {
			assert.Contains(t, string(contents), fmt.Sprintf("| Skipped | %s | 1.00s | reason %d |", child.Name, index))
		}
	}
}

type reportLifecycleInput struct {
	Pending     bool
	Interrupted bool
	Result      bool
}

func TestAcceptanceReportLifecycle(t *testing.T) {
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_normal_completion", Input: reportLifecycleInput{Result: true}, Expected: "PASSED"},
		{Name: "test_queued_worker_is_incomplete", Input: reportLifecycleInput{Pending: true, Result: true}, Expected: "INCOMPLETE"},
		{Name: "test_no_worker_results_is_incomplete", Input: reportLifecycleInput{}, Expected: "INCOMPLETE"},
		{Name: "test_interrupted_with_all_results_is_incomplete", Input: reportLifecycleInput{Interrupted: true, Result: true}, Expected: "INCOMPLETE"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Setenv("ACI_TEST_LOG_DIR", t.TempDir())
			input := testCase.Input.(reportLifecycleInput)
			run := Runner{}
			child, _, err := run.prepareWorker("TestLocal", time.Time{})
			if !assert.NoError(t, err) {
				return
			}
			defer run.stop()
			path := filepath.Join(run.directory, "summary.md")
			assert.NoFileExists(t, path, "registering invocations must not write the report")
			if input.Pending {
				_, _, err := run.prepareWorker("TestPending", time.Time{})
				assert.NoError(t, err)
			}
			if input.Interrupted {
				run.stop()
			}
			if input.Result {
				run.recordResult(child.Name, workerResult{Action: "pass", LogPath: child.LogPath})
			}
			assert.NoFileExists(t, path, "only Finish may write the report")
			assert.NoError(t, run.Finish())
			contents, err := os.ReadFile(path)
			if assert.NoError(t, err) {
				assert.Contains(t, string(contents), "Overall result: **"+testCase.Expected.(string)+"**")
				if input.Interrupted {
					assert.Contains(t, string(contents), "The test run did not complete because it was interrupted.")
				}
			}
			fixedTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
			assert.NoError(t, os.Chtimes(path, fixedTime, fixedTime))
			assert.NoError(t, run.Finish())
			again, err := os.ReadFile(path)
			assert.NoError(t, err)
			assert.Equal(t, contents, again, "finalizing twice must preserve the completed report")
			info, err := os.Stat(path)
			if assert.NoError(t, err) {
				assert.True(t, info.ModTime().Equal(fixedTime), "finalizing twice must not write the file again")
			}
		})
	}
}

func TestAcceptanceRunnerDeadline(t *testing.T) {
	test.InitializeTest(t)
	testCases := []test.TestCase{
		{Name: "test_without_go_deadline", Input: time.Duration(0)},
		{Name: "test_long_run_reserves_thirty_seconds", Input: 3 * time.Hour},
		{Name: "test_short_run_has_time_to_execute", Input: 10 * time.Second},
		{Name: "test_expired_deadline", Input: -time.Second},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Setenv("ACI_TEST_LOG_DIR", t.TempDir())
			timeout := testCase.Input.(time.Duration)
			deadline := time.Time{}
			if timeout != 0 {
				deadline = time.Now().Add(timeout)
			}
			run := Runner{}
			child, ctx, err := run.prepareWorker("TestLocal", deadline)
			if !assert.NoError(t, err) {
				return
			}
			defer run.stop()
			assert.Equal(t, 30*time.Minute, child.Timeout)
			softDeadline, hasDeadline := ctx.Deadline()
			assert.Equal(t, timeout != 0, hasDeadline)
			switch {
			case timeout >= time.Hour:
				assert.Equal(t, deadline.Add(-30*time.Second), softDeadline)
			case timeout > 0:
				assert.True(t, softDeadline.Before(deadline))
				assert.True(t, deadline.Sub(softDeadline) <= time.Second)
				assert.NoError(t, ctx.Err(), "a short timeout must not cancel the runner immediately")
			case timeout < 0:
				assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
			}
		})
	}
}

func TestAcceptanceRunnerDeadlineReport(t *testing.T) {
	test.InitializeTest(t)
	t.Setenv("ACI_TEST_LOG_DIR", t.TempDir())
	run := Runner{}
	_, ctx, err := run.prepareWorker("TestLocal", time.Now().Add(50*time.Millisecond))
	if !assert.NoError(t, err) {
		return
	}
	defer run.stop()
	select {
	case <-ctx.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("the shared runner deadline did not expire")
	}
	assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
	assert.NoFileExists(t, filepath.Join(run.directory, "summary.md"))
	assert.NoError(t, run.Finish())
	contents, err := os.ReadFile(filepath.Join(run.directory, "summary.md"))
	if assert.NoError(t, err) {
		assert.Contains(t, string(contents), "Overall result: **INCOMPLETE**")
		assert.Contains(t, string(contents), "The test run did not complete because its time limit was reached.")
	}
}

func TestAcceptanceReportUnusedRunner(t *testing.T) {
	test.InitializeTest(t)
	t.Setenv("ACI_TEST_LOG_DIR", t.TempDir())
	run := Runner{}
	assert.NoError(t, run.Finish())
	assert.Empty(t, run.directory, "unit runs and worker processes must not create another report")
}

func TestAcceptanceRepeatedLogPaths(t *testing.T) {
	test.InitializeTest(t)

	t.Setenv("ACI_TEST_LOG_DIR", t.TempDir())
	run := Runner{}
	first, _, err := run.prepareWorker("TestLocal", time.Time{})
	if !assert.NoError(t, err) {
		return
	}
	defer run.stop()
	second, _, err := run.prepareWorker("TestLocal", time.Time{})
	if assert.NoError(t, err) {
		assert.NotEqual(t, first.LogPath, second.LogPath)
		assert.Equal(t, filepath.Dir(first.LogPath), filepath.Dir(second.LogPath))
		assert.Equal(t, "TestLocal-001.log", filepath.Base(first.LogPath))
		assert.Equal(t, "TestLocal-002.log", filepath.Base(second.LogPath))
	}
}
