package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEmitOutputJSONToWriter(t *testing.T) {
	report := validationReport{
		Summary: validationSummary{FilesChecked: 2, Passed: 2, Failed: 0, DurationMS: 15},
		Issues:  []validationIssue{},
	}

	var stdout bytes.Buffer
	if err := emitOutput("json", "", report, &stdout); err != nil {
		t.Fatalf("emitOutput(json) error: %v", err)
	}

	var got validationReport
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("stdout not valid json: %v, out=%q", err, stdout.String())
	}
	if got.Summary.FilesChecked != 2 || got.Summary.Passed != 2 || got.Summary.Failed != 0 {
		t.Fatalf("unexpected summary: %+v", got.Summary)
	}
}

func TestEmitOutputTextToFile(t *testing.T) {
	report := validationReport{
		Summary: validationSummary{FilesChecked: 1, Passed: 1, Failed: 0, DurationMS: 8},
		Issues:  []validationIssue{},
	}

	dir := t.TempDir()
	outPath := filepath.Join(dir, "report.txt")
	if err := emitOutput("text", outPath, report, &bytes.Buffer{}); err != nil {
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
	if err := emitOutput("text", outPath, report, &bytes.Buffer{}); err != nil {
		t.Fatalf("emitOutput(text,file) error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file failed: %v", err)
	}

	want := "[error] read config failed\n"
	if string(data) != want {
		t.Fatalf("unexpected text output\nwant: %q\n got: %q", want, string(data))
	}
}

func TestRunWithArgsInvalidFormat(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runWithArgs([]string{"--format", "xml", "ruler/ruler_employees.json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected code 2, got %d", code)
	}
	if !bytes.Contains(stderr.Bytes(), []byte("invalid --format")) {
		t.Fatalf("expected invalid --format error in stderr, got: %q", stderr.String())
	}
}

func TestRunWithArgsMissingRequiredArg(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runWithArgs([]string{"--format", "json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected code 2, got %d", code)
	}
	if !bytes.Contains(stderr.Bytes(), []byte("Usage:")) {
		t.Fatalf("expected usage in stderr, got: %q", stderr.String())
	}
}

func TestRunWithArgsMissingConfigProducesJSONIssue(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runWithArgs([]string{"--format", "json", "--quiet", "/tmp/not-found-ruler.json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected code 2, got %d", code)
	}

	var report validationReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("expected json report in stdout, got err=%v out=%q", err, stdout.String())
	}
	if len(report.Issues) == 0 {
		t.Fatalf("expected at least one issue, got none")
	}
	if report.Issues[0].File == "" {
		t.Fatalf("expected issue file context, got %+v", report.Issues[0])
	}
}

func TestRunWithArgsSuccessJSON(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	repoRoot := filepath.Join(cwd, "..", "..")
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	defer func() { _ = os.Chdir(cwd) }()

	var stdout, stderr bytes.Buffer
	rulerPath := filepath.Join("ruler", "ruler_employees.json")
	code := runWithArgs([]string{"--format", "json", "--quiet", rulerPath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected code 0, got %d; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}

	var report validationReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("expected json report in stdout, got err=%v out=%q", err, stdout.String())
	}
	if report.Summary.FilesChecked == 0 {
		t.Fatalf("expected non-zero files_checked, got %+v", report.Summary)
	}
}

func TestRunWithArgsValidationFailureProducesCode1(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	repoRoot := filepath.Join(cwd, "..", "..")
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	defer func() { _ = os.Chdir(cwd) }()

	cfg := `{
  "csvons_metadata": {
    "csv_file_folder": "testdata",
    "name_index": 0,
    "data_index": 1,
    "extension": ".csv",
    "lev1_separator": ";",
    "lev2_separator": ":",
    "field_connector": "|"
  },
  "employees": {
    "vtype": [{"field": "Salary", "type": "bool"}]
  }
}`
	tmpCfg := filepath.Join(t.TempDir(), "bad_rule.json")
	if err := os.WriteFile(tmpCfg, []byte(cfg), 0o644); err != nil {
		t.Fatalf("write temp config failed: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := runWithArgs([]string{"--format", "json", "--quiet", tmpCfg}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected code 1 for validation failure, got %d; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}

	var report validationReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("expected json report in stdout, got err=%v out=%q", err, stdout.String())
	}
	if len(report.Issues) == 0 {
		t.Fatalf("expected validation issue, got none")
	}
	if report.Issues[0].File != "employees" || report.Issues[0].Rule != "vtype" {
		t.Fatalf("expected file/rule context in issue, got %+v", report.Issues[0])
	}
	if report.Issues[0].Field != "Salary" {
		t.Fatalf("expected extracted field name 'Salary', got %+v", report.Issues[0])
	}
	if report.Summary.Failed == 0 {
		t.Fatalf("expected failed count > 0, got %+v", report.Summary)
	}
}

func TestExtractFieldName(t *testing.T) {
	if got := extractFieldName("src_field [Salary] value [85000] is not a bool"); got != "Salary" {
		t.Fatalf("expected Salary, got %q", got)
	}
	if got := extractFieldName("unknown error"); got != "" {
		t.Fatalf("expected empty field, got %q", got)
	}
}
