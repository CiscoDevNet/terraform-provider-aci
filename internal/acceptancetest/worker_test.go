package acceptancetest

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type acceptanceWorkerInput struct {
	Name  string
	Mode  string
	Trace bool
}

type acceptanceWorkerExpected struct {
	Action string
	Output string
	Trace  bool
}

func TestAcceptanceSubprocessResults(t *testing.T) {
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_pass_with_legacy_output_and_expected_error_text", Input: acceptanceWorkerInput{Mode: "pass"}, Expected: acceptanceWorkerExpected{Action: "pass", Output: `{"imdata":[],"totalCount":"0"}`}},
		{Name: "test_failed_assertion", Input: acceptanceWorkerInput{Mode: "fail"}, Expected: acceptanceWorkerExpected{Action: "fail", Output: "local assertion failed"}},
		{Name: "test_skip_reason", Input: acceptanceWorkerInput{Mode: "skip"}, Expected: acceptanceWorkerExpected{Action: "skip", Output: "APIC 5.2(8h) is outside the supported range: 6.0(2h)-"}},
		{Name: "test_crash", Input: acceptanceWorkerInput{Mode: "panic"}, Expected: acceptanceWorkerExpected{Action: "fail", Output: "local worker panic"}},
		{Name: "test_no_matching_test_is_not_a_pass", Input: acceptanceWorkerInput{Name: "TestAcceptanceWorkerDoesNotExist", Mode: "pass"}, Expected: acceptanceWorkerExpected{Action: "fail", Output: "intended acceptance test did not run"}},
		{Name: "test_invalid_prefix_is_a_failure", Input: acceptanceWorkerInput{Mode: "invalid_prefix"}, Expected: acceptanceWorkerExpected{Action: "fail", Output: "prefix length must be positive"}},
		{Name: "test_existing_trace_file", Input: acceptanceWorkerInput{Mode: "trace", Trace: true}, Expected: acceptanceWorkerExpected{Action: "pass", Output: "local trace written", Trace: true}},
		{Name: "test_missing_trace_file", Input: acceptanceWorkerInput{Mode: "pass", Trace: true}, Expected: acceptanceWorkerExpected{Action: "pass", Output: `{"imdata":[],"totalCount":"0"}`}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(acceptanceWorkerInput)
			expected := testCase.Expected.(acceptanceWorkerExpected)
			if input.Name == "" {
				input.Name = "TestAcceptanceWorkerFixture"
			}
			child := worker{
				Name:        input.Name,
				LogPath:     filepath.Join(t.TempDir(), "detail.log"),
				Timeout:     10 * time.Second,
				Environment: []string{"TF_ACC=", "ACI_TEST_LOCAL_FIXTURE=" + input.Mode},
			}
			if input.Trace {
				child.TraceMask = filepath.Join(filepath.Dir(child.LogPath), "%s.trace.log")
			}
			result := child.run(context.Background())
			assert.Equal(t, expected.Action, result.Action)
			if expected.Trace {
				assert.Equal(t, fmt.Sprintf(child.TraceMask, input.Name), result.TracePath)
			} else {
				assert.Empty(t, result.TracePath)
			}
			contents, err := os.ReadFile(child.LogPath)
			if assert.NoError(t, err) {
				assert.Contains(t, string(contents), expected.Output)
			}
			if expected.Action == "skip" {
				assert.Equal(t, expected.Output, result.SkipReason)
			}
		})
	}
}

func TestAcceptanceSubprocessTimeout(t *testing.T) {
	test.InitializeTest(t)

	child := worker{
		Name:        "TestAcceptanceWorkerFixture",
		LogPath:     filepath.Join(t.TempDir(), "timeout.log"),
		Timeout:     150 * time.Millisecond,
		Environment: []string{"TF_ACC=", "ACI_TEST_LOCAL_FIXTURE=sleep"},
	}
	result := child.run(context.Background())
	assert.Equal(t, "fail", result.Action)
	assert.Less(t, result.Elapsed, 5*time.Second)
}

func TestAcceptanceSubprocessCancellation(t *testing.T) {
	test.InitializeTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	directory := t.TempDir()
	heartbeat := filepath.Join(directory, "heartbeat")
	child := worker{
		Name:        "TestAcceptanceWorkerFixture",
		LogPath:     filepath.Join(directory, "cancel.log"),
		Timeout:     10 * time.Second,
		Environment: []string{"TF_ACC=", "ACI_TEST_LOCAL_FIXTURE=tree", "ACI_TEST_HEARTBEAT=" + heartbeat},
	}
	finished := make(chan workerResult, 1)
	go func() { finished <- child.run(ctx) }()
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(heartbeat); err == nil {
			break
		}
		if time.Now().After(deadline) {
			cancel()
			t.Fatal("local descendant process did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	select {
	case result := <-finished:
		assert.Equal(t, "fail", result.Action)
	case <-time.After(5 * time.Second):
		t.Fatal("cancelled worker did not finish")
	}
	before, err := os.ReadFile(heartbeat)
	if !assert.NoError(t, err) {
		return
	}
	time.Sleep(100 * time.Millisecond)
	after, err := os.ReadFile(heartbeat)
	if assert.NoError(t, err) {
		assert.Equal(t, string(before), string(after), "descendant must stop updating its heartbeat")
	}
}

func TestAcceptanceWorkerFixture(t *testing.T) {
	// This worker is selected explicitly by the local subprocess tests and never calls APIC.
	mode := os.Getenv("ACI_TEST_LOCAL_FIXTURE")
	if mode == "" {
		t.Skip("local subprocess fixture")
	}
	switch mode {
	case "pass":
		fmt.Print(`{"imdata":[],"totalCount":"0"}`)
		t.Log("Read Missing Object (Expect Error): expected error")
	case "fail":
		t.Error("local assertion failed")
	case "skip":
		Skip(t, "APIC 5.2(8h) is outside the supported range: 6.0(2h)-")
	case "panic":
		panic("local worker panic")
	case "invalid_prefix":
		Prefix(t, PrefixOptions{Length: -1})
	case "trace":
		path := fmt.Sprintf(os.Getenv("TF_LOG_PATH_MASK"), strings.ReplaceAll(t.Name(), "/", "__"))
		if err := os.WriteFile(path, []byte("local trace"), 0600); err != nil {
			t.Fatal(err)
		}
		t.Log("local trace written")
	case "sleep":
		time.Sleep(30 * time.Second)
	case "tree":
		executable, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}
		command := exec.Command(executable, "-test.run=^TestAcceptanceWorkerFixture$", "-test.v=test2json")
		command.Env = childEnvironment("ACI_TEST_LOCAL_FIXTURE=heartbeat")
		command.Stdout, command.Stderr = os.Stdout, os.Stderr
		if err := command.Start(); err != nil {
			t.Fatal(err)
		}
		defer command.Wait()
		time.Sleep(30 * time.Second)
	case "heartbeat":
		for {
			if err := os.WriteFile(os.Getenv("ACI_TEST_HEARTBEAT"), []byte(time.Now().String()), 0600); err != nil {
				t.Fatal(err)
			}
			time.Sleep(10 * time.Millisecond)
		}
	default:
		t.Fatalf("unknown local worker mode %q", mode)
	}
}
