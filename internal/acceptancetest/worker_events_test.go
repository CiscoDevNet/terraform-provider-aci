package acceptancetest

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type workerEventsInput struct {
	Output  string
	FailLog bool
}

type workerEventsExpected struct {
	Action     string
	SkipReason string
	Error      string
}

type failingLogWriter struct{}

func (failingLogWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestWorkerEventsRead(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_pass_and_unattributed_output",
			Input: workerEventsInput{Output: `{"Action":"run","Test":"TestLocal"}
{"Action":"output","Output":"legacy output\n"}
{"Action":"pass","Test":"TestLocal"}`},
			Expected: workerEventsExpected{Action: "pass"},
		},
		{
			Name: "test_skip_reason_across_events",
			Input: workerEventsInput{Output: `{"Action":"output","Test":"TestLocal","Output":"    Skip rea"}
{"Action":"output","Test":"TestLocal","Output":"son: APIC version unsupported\n"}
{"Action":"skip","Test":"TestLocal"}`},
			Expected: workerEventsExpected{Action: "skip", SkipReason: "APIC version unsupported"},
		},
		{
			Name:     "test_malformed_json",
			Input:    workerEventsInput{Output: "invalid json"},
			Expected: workerEventsExpected{Error: "invalid character"},
		},
		{
			Name:     "test_detail_log_write_error",
			Input:    workerEventsInput{Output: `{"Action":"output","Test":"TestLocal","Output":"detail"}`, FailLog: true},
			Expected: workerEventsExpected{Error: io.ErrClosedPipe.Error()},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(workerEventsInput)
			expected := testCase.Expected.(workerEventsExpected)
			events := workerEvents{Name: "TestLocal"}
			var detailedLog io.Writer = &bytes.Buffer{}
			if input.FailLog {
				detailedLog = failingLogWriter{}
			}
			err := events.read(strings.NewReader(input.Output), detailedLog)
			if expected.Error != "" {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), expected.Error)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, expected.Action, events.Action)
			assert.Equal(t, expected.SkipReason, events.SkipReason)
		})
	}
}
