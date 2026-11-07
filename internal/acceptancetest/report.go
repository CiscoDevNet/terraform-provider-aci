package acceptancetest

import (
	"cmp"
	"fmt"
	"html"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

type reportResult struct {
	Name string
	workerResult
}

type runReport struct {
	Started          time.Time
	Finished         time.Time
	Scheduled        int
	IncompleteReason string
	Results          []reportResult
}

func (report runReport) write(directory string) error {
	return os.WriteFile(filepath.Join(directory, "summary.md"), []byte(report.render()), 0600)
}

func (report runReport) render() string {
	// Sort a copy so report rendering cannot reorder the collected worker results.
	results := slices.Clone(report.Results)
	slices.SortFunc(results, func(a, b reportResult) int {
		if order := cmp.Compare(a.Name, b.Name); order != 0 {
			return order
		}
		return cmp.Compare(a.LogPath, b.LogPath)
	})
	var failures, successful []reportResult
	passed, skipped := 0, 0
	for _, result := range results {
		switch result.Action {
		case "fail":
			failures = append(failures, result)
		case "pass":
			passed++
			successful = append(successful, result)
		case "skip":
			skipped++
			successful = append(successful, result)
		}
	}

	status := "PASSED"
	if report.IncompleteReason != "" || len(results) < report.Scheduled {
		status = "INCOMPLETE"
	} else if len(failures) > 0 {
		status = "FAILED"
	}
	// Wall-clock time includes scheduling and overlapping workers, rather than adding their durations.
	elapsed := report.Finished.Sub(report.Started).Round(time.Second)
	var contents strings.Builder
	contents.WriteString("# ACI generated acceptance test report\n\n")
	fmt.Fprintf(&contents, "The run has completed **%d of %d** test invocations in **%s**: **%d passed**, **%d failed**, and **%d skipped**. Overall result: **%s**.\n\n", len(results), report.Scheduled, elapsed, passed, len(failures), skipped, status)
	if status == "INCOMPLETE" {
		reason := report.IncompleteReason
		if reason == "" {
			reason = "The test run did not complete; some scheduled test invocations have no result."
		}
		fmt.Fprintf(&contents, "%s\n\n", reportText(reason))
	}
	contents.WriteString("## Failed tests\n\n")
	if len(failures) == 0 {
		contents.WriteString("No failed tests.\n\n")
	} else {
		report.writeTable(&contents, failures)
	}
	contents.WriteString("<details>\n<summary>Successful and skipped tests</summary>\n\n")
	if len(successful) == 0 {
		contents.WriteString("No successful or skipped tests.\n\n")
	} else {
		report.writeTable(&contents, successful)
	}
	contents.WriteString("</details>\n")
	return contents.String()
}

func (report runReport) writeTable(contents *strings.Builder, results []reportResult) {
	contents.WriteString("| Result | Test | Duration | Reason | Log | Trace |\n")
	contents.WriteString("|---|---|---:|---|---|---|\n")
	for _, result := range results {
		status, reason := "Passed", "—"
		switch result.Action {
		case "fail":
			status, reason = "Failed", result.Error
			if reason == "" {
				reason = "See detailed log"
			}
		case "skip":
			status, reason = "Skipped", result.SkipReason
		}
		fmt.Fprintf(contents, "| %s | %s | %.2fs | %s | %s | %s |\n", status, reportText(result.Name), result.Elapsed.Seconds(), reportText(reason), reportFileLink("Log", result.LogPath), reportFileLink("Trace", result.TracePath))
	}
	contents.WriteString("\n")
}

func reportText(value string) string {
	// Keep diagnostics on one table row and preserve characters with Markdown meaning as literal text.
	value = html.EscapeString(strings.Join(strings.Fields(value), " "))
	return strings.NewReplacer("\\", "\\\\", "|", "&#124;", "`", "\\`", "*", "\\*", "_", "\\_", "[", "\\[", "]", "\\]").Replace(value)
}

func reportFileLink(label, path string) string {
	if path == "" {
		return "—"
	}
	if info, err := os.Stat(path); err != nil || info.IsDir() {
		return "—"
	}
	// Worker logs live beside summary.md, so the report remains portable when the directory is copied.
	return fmt.Sprintf("[%s](./%s)", label, url.PathEscape(filepath.Base(path)))
}
