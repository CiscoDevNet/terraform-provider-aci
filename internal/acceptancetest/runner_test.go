package acceptancetest

import (
	"context"
	"fmt"
	"os"
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
		})
	}
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
