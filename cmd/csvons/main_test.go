package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
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
