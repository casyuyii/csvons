package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEmitOutputJSONToStdout(t *testing.T) {
	report := validationReport{
		Summary: validationSummary{FilesChecked: 2, Passed: 2, Failed: 0, DurationMS: 15},
		Issues:  []validationIssue{},
	}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe create failed: %v", err)
	}
	defer r.Close()

	origStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = origStdout }()

	if err := emitOutput("json", "", report); err != nil {
		t.Fatalf("emitOutput(json) error: %v", err)
	}
	_ = w.Close()

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read stdout pipe failed: %v", err)
	}

	var got validationReport
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("stdout not valid json: %v, out=%q", err, string(out))
	}
	if got.Summary.FilesChecked != 2 || got.Summary.Passed != 2 || got.Summary.Failed != 0 {
		t.Fatalf("unexpected summary: %+v", got.Summary)
	}
	if got.SchemaVersion != reportSchemaVersion {
		t.Fatalf("unexpected schema version: %q", got.SchemaVersion)
	}
}

func TestEmitOutputTextToFile(t *testing.T) {
	report := validationReport{
		Summary: validationSummary{FilesChecked: 1, Passed: 1, Failed: 0, DurationMS: 8},
		Issues:  []validationIssue{},
	}

	dir := t.TempDir()
	outPath := filepath.Join(dir, "report.txt")
	if err := emitOutput("text", outPath, report); err != nil {
		t.Fatalf("emitOutput(text,file) error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file failed: %v", err)
	}

	want := "Validation succeeded: files_checked=1 passed=1 failed=0 duration_ms=8\n"
	if string(data) != want {
		t.Fatalf("unexpected text output\nwant: %q\n got: %q", want, string(data))
	}
}

func TestEmitOutputTextWithIssues(t *testing.T) {
	report := validationReport{
		Summary: validationSummary{},
		Issues: []validationIssue{
			{Severity: "error", Message: "read config failed"},
		},
	}

	dir := t.TempDir()
	outPath := filepath.Join(dir, "report.txt")
	if err := emitOutput("text", outPath, report); err != nil {
		t.Fatalf("emitOutput(text,file) error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file failed: %v", err)
	}

	want := "[error] read config failed\n"
	if string(data) != want {
		t.Fatalf("unexpected text output\\nwant: %q\\n got: %q", want, string(data))
	}
}

func TestRunWithArgsJSONValidationFailureIncludesStructuredIssue(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "users.csv")
	if err := os.WriteFile(csvPath, []byte("Username\nalpha\nalpha\n"), 0o644); err != nil {
		t.Fatalf("write csv failed: %v", err)
	}

	configPath := filepath.Join(dir, "ruler.json")
	config := map[string]any{
		"users": map[string]any{
			"unique": map[string]any{
				"fields": []string{"Username"},
			},
		},
		"csvons_metadata": map[string]any{
			"csv_file_folder": dir,
			"name_index":      0,
			"data_index":      1,
			"extension":       ".csv",
		},
	}
	writeJSONFile(t, configPath, config)

	reportPath := filepath.Join(dir, "report.json")
	code := runWithArgs([]string{"--format", "json", "--output", reportPath, configPath})
	if code != 1 {
		t.Fatalf("unexpected exit code: got %d want 1", code)
	}

	report := readReportFile(t, reportPath)
	if report.SchemaVersion != reportSchemaVersion {
		t.Fatalf("unexpected schema version: %q", report.SchemaVersion)
	}
	if len(report.Issues) != 1 {
		t.Fatalf("unexpected issue count: %d", len(report.Issues))
	}

	issue := report.Issues[0]
	if issue.File != "users.csv" {
		t.Fatalf("unexpected issue file: %q", issue.File)
	}
	if issue.Rule != "unique" {
		t.Fatalf("unexpected issue rule: %q", issue.Rule)
	}
	if issue.Field != "Username" {
		t.Fatalf("unexpected issue field: %q", issue.Field)
	}
	if issue.Row == nil || *issue.Row != 3 {
		t.Fatalf("unexpected issue row: %#v", issue.Row)
	}
	if issue.Value != "alpha" {
		t.Fatalf("unexpected issue value: %q", issue.Value)
	}
	if !strings.Contains(issue.Message, "already exists") {
		t.Fatalf("unexpected issue message: %q", issue.Message)
	}
}

func TestRunWithArgsJSONRuntimeFailureUsesExitCodeTwo(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "users.csv")
	if err := os.WriteFile(csvPath, []byte("Username\nalpha\n"), 0o644); err != nil {
		t.Fatalf("write csv failed: %v", err)
	}

	configPath := filepath.Join(dir, "ruler.json")
	config := map[string]any{
		"users": map[string]any{
			"unique": map[string]any{
				"fields": []string{"MissingField"},
			},
		},
		"csvons_metadata": map[string]any{
			"csv_file_folder": dir,
			"name_index":      0,
			"data_index":      1,
			"extension":       ".csv",
		},
	}
	writeJSONFile(t, configPath, config)

	reportPath := filepath.Join(dir, "report.json")
	code := runWithArgs([]string{"--format", "json", "--output", reportPath, configPath})
	if code != 2 {
		t.Fatalf("unexpected exit code: got %d want 2", code)
	}

	report := readReportFile(t, reportPath)
	if len(report.Issues) != 1 {
		t.Fatalf("unexpected issue count: %d", len(report.Issues))
	}

	issue := report.Issues[0]
	if issue.File != "users.csv" {
		t.Fatalf("unexpected issue file: %q", issue.File)
	}
	if issue.Rule != "unique" {
		t.Fatalf("unexpected issue rule: %q", issue.Rule)
	}
	if issue.Field != "MissingField" {
		t.Fatalf("unexpected issue field: %q", issue.Field)
	}
	if issue.Row != nil {
		t.Fatalf("expected nil row, got %#v", issue.Row)
	}
	if !strings.Contains(issue.Message, "cannot resolve values") {
		t.Fatalf("unexpected issue message: %q", issue.Message)
	}
}

func writeJSONFile(t *testing.T, path string, value any) {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal json failed: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write json file failed: %v", err)
	}
}

func readReportFile(t *testing.T, path string) validationReport {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read report failed: %v", err)
	}

	var report validationReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("unmarshal report failed: %v", err)
	}
	return report
}
