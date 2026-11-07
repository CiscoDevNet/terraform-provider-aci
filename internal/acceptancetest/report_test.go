package acceptancetest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type reportRenderExpected struct {
	Contains    []string
	NotContains []string
}

func TestAcceptanceReportRendering(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	started := time.Date(2026, 11, 7, 10, 0, 0, 0, time.UTC)
	finished := started.Add(2*time.Minute + 3*time.Second)
	testCases := []test.TestCase{
		{
			Name: "test_mixed_results_and_wall_clock_duration",
			Input: runReport{Started: started, Finished: finished, Scheduled: 4, Results: []reportResult{
				{Name: "TestZ", workerResult: workerResult{Action: "skip", SkipReason: "APIC 5.2(8h) is outside the supported range: 6.0(2h)-", Elapsed: time.Second}},
				{Name: "TestB", workerResult: workerResult{Action: "pass", Elapsed: 90 * time.Second}},
				{Name: "TestC", workerResult: workerResult{Action: "fail", Elapsed: 60 * time.Second}},
				{Name: "TestA", workerResult: workerResult{Action: "fail", Error: "Child test exited unsuccessfully: exit status 1", Elapsed: 50 * time.Second}},
			}},
			Expected: reportRenderExpected{Contains: []string{
				"**4 of 4** test invocations in **2m3s**: **1 passed**, **2 failed**, and **1 skipped**",
				"Overall result: **FAILED**",
				"| Failed | TestA | 50.00s | Child test exited unsuccessfully: exit status 1 | — | — |",
				"| Failed | TestC | 60.00s | See detailed log | — | — |",
				"| Passed | TestB | 90.00s | — | — | — |",
				"| Skipped | TestZ | 1.00s | APIC 5.2(8h) is outside the supported range: 6.0(2h)- |",
			}},
		},
		{
			Name: "test_no_failures",
			Input: runReport{Started: started, Finished: finished, Scheduled: 1, Results: []reportResult{
				{Name: "TestA", workerResult: workerResult{Action: "pass"}},
			}},
			Expected: reportRenderExpected{
				Contains:    []string{"Overall result: **PASSED**", "No failed tests.", "| Passed | TestA |"},
				NotContains: []string{"| Failed |"},
			},
		},
		{
			Name:  "test_unfinished_invocations",
			Input: runReport{Started: started, Finished: finished, Scheduled: 2},
			Expected: reportRenderExpected{Contains: []string{
				"**0 of 2** test invocations", "Overall result: **INCOMPLETE**", "The test run did not complete; some scheduled test invocations have no result.", "No failed tests.", "No successful or skipped tests.",
			}},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			report := testCase.Input.(runReport)
			expected := testCase.Expected.(reportRenderExpected)
			actual := report.render()
			for _, text := range expected.Contains {
				assert.Contains(t, actual, text)
			}
			for _, text := range expected.NotContains {
				assert.NotContains(t, actual, text)
			}
			failed, successful, found := strings.Cut(actual, "<details>")
			assert.True(t, found)
			assert.NotContains(t, failed, "| Passed |")
			assert.NotContains(t, failed, "| Skipped |")
			assert.NotContains(t, successful, "| Failed |")
		})
	}
}

func TestAcceptanceReportRepeatedInvocations(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	started := time.Now().Add(-time.Minute)
	report := runReport{Started: started, Finished: started.Add(time.Minute), Scheduled: 3, Results: []reportResult{
		{Name: "TestRepeated", workerResult: workerResult{Action: "pass", Elapsed: 2 * time.Second, LogPath: "TestRepeated-002.log"}},
		{Name: "TestZ", workerResult: workerResult{Action: "pass", Elapsed: 3 * time.Second}},
		{Name: "TestRepeated", workerResult: workerResult{Action: "pass", Elapsed: time.Second, LogPath: "TestRepeated-001.log"}},
	}}
	actual := report.render()
	assert.Contains(t, actual, "**3 passed**")
	assert.Equal(t, 2, strings.Count(actual, "| Passed | TestRepeated |"))
	assert.Less(t, strings.Index(actual, "| TestRepeated | 1.00s |"), strings.Index(actual, "| TestRepeated | 2.00s |"))
	assert.Less(t, strings.Index(actual, "| TestRepeated | 2.00s |"), strings.Index(actual, "| TestZ |"))
	assert.Equal(t, "TestRepeated-002.log", report.Results[0].LogPath, "rendering must preserve collection order")
}

func TestAcceptanceReportFileLinks(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	directory := t.TempDir()
	logPath := filepath.Join(directory, "test.log")
	tracePath := filepath.Join(directory, "trace #1.log")
	for _, path := range []string{logPath, tracePath} {
		if err := os.WriteFile(path, []byte("local test log"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	testCases := []test.TestCase{
		{Name: "test_log_link", Input: logPath, Expected: "[Log](./test.log)"},
		{Name: "test_encoded_filename", Input: tracePath, Expected: "[Log](./trace%20%231.log)"},
		{Name: "test_no_file", Input: "", Expected: "—"},
		{Name: "test_missing_file", Input: filepath.Join(directory, "missing.log"), Expected: "—"},
		{Name: "test_directory", Input: directory, Expected: "—"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, testCase.Expected, reportFileLink("Log", testCase.Input.(string)))
		})
	}
	started := time.Now()
	report := runReport{Started: started, Finished: started, Scheduled: 1, Results: []reportResult{{
		Name: "TestLocal", workerResult: workerResult{Action: "pass", LogPath: logPath, TracePath: tracePath},
	}}}
	actual := report.render()
	assert.Contains(t, actual, "[Log](./test.log) | [Trace](./trace%20%231.log)")
	assert.NotContains(t, actual, directory, "links must survive copying the report directory")
}

func TestAcceptanceReportMarkdownEscaping(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	started := time.Now()
	report := runReport{Started: started, Finished: started, Scheduled: 1, Results: []reportResult{{
		Name: "TestLocal", workerResult: workerResult{Action: "skip", SkipReason: "APIC <6.0> | unavailable\ntry another version"},
	}}}
	actual := report.render()
	assert.Contains(t, actual, "| Skipped | TestLocal | 0.00s | APIC &lt;6.0&gt; &#124; unavailable try another version | — | — |")
}
