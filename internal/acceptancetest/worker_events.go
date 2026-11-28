package acceptancetest

import (
	"encoding/json"
	"io"
	"strings"
)

type workerEvents struct {
	Name       string
	Ran        bool
	Action     string
	SkipReason string
	pending    string
	lastOutput string
}

type testEvent struct {
	Action string
	Test   string
	Output string
}

func (events *workerEvents) read(output io.Reader, detailedLog io.Writer) error {
	decoder := json.NewDecoder(output)
	for {
		var event testEvent
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if _, err := io.WriteString(detailedLog, event.Output); err != nil {
			return err
		}
		events.record(event)
	}
}

func (events *workerEvents) record(event testEvent) {
	if event.Test != events.Name {
		return
	}
	switch event.Action {
	case "run":
		events.Ran = true
	case "pass", "fail", "skip":
		events.Action = event.Action
	case "output":
		events.recordOutput(event.Output)
	}
}

func (events *workerEvents) recordOutput(output string) {
	// test2json can split one output line across events; retain the unfinished part.
	events.pending += output
	for {
		line, remaining, found := strings.Cut(events.pending, "\n")
		if !found {
			return
		}
		events.pending = remaining
		if _, reason, found := strings.Cut(line, skipPrefix); found {
			events.SkipReason = strings.TrimSpace(reason)
		}
		if line = strings.TrimSpace(line); line != "" && !strings.HasPrefix(line, "===") && !strings.HasPrefix(line, "---") {
			events.lastOutput = line
		}
	}
}
