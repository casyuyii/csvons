// Command csvons validates CSV files against constraint rules defined in a JSON configuration file.
//
// Usage:
//
//	csvons <ruler.json>
//
// The program reads the specified ruler JSON file, parses the metadata
// and constraint rules, then validates each referenced CSV file against its rules.
//
// Supported constraints:
//   - exists: values in a column must exist in another CSV file's column
//   - unique: values in a column must be unique across all rows
//   - vtype: values must conform to a specified type and optional range
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	csvons "csvons/internal/csvons"
)

type validationSummary struct {
	FilesChecked int   `json:"files_checked"`
	Passed       int   `json:"passed"`
	Failed       int   `json:"failed"`
	DurationMS   int64 `json:"duration_ms"`
}

type validationIssue struct {
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type validationReport struct {
	Summary validationSummary `json:"summary"`
	Issues  []validationIssue `json:"issues"`
}

func main() {
	os.Exit(run())
}

func run() (code int) {
	var format string
	var outputPath string
	var rules map[string]json.RawMessage

	flag.StringVar(&format, "format", "text", "output format: text or json")
	flag.StringVar(&outputPath, "output", "", "optional output file path")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--format text|json] [--output <path>] <ruler.json>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nValidate CSV files against constraint rules defined in a JSON configuration file.\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return 2
	}
	if format != "text" && format != "json" {
		fmt.Fprintf(os.Stderr, "invalid --format value %q; expected text or json\n", format)
		return 2
	}

	configFileName := flag.Arg(0)
	startAt := time.Now()
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}
		msg := fmt.Sprint(recovered)
		_ = emitOutput(format, outputPath, validationReport{
			Summary: validationSummary{
				FilesChecked: len(rules),
				Passed:       0,
				Failed:       1,
				DurationMS:   time.Since(startAt).Milliseconds(),
			},
			Issues: []validationIssue{{Message: msg, Severity: "error"}},
		})
		code = 1
	}()

	var metadata *csvons.Metadata
	rules, metadata = csvons.ReadConfigFile(configFileName)
	if rules == nil || metadata == nil {
		_ = emitOutput(format, outputPath, validationReport{
			Summary: validationSummary{},
			Issues: []validationIssue{{
				Message:  fmt.Sprintf("read config file error: file_name=%s", configFileName),
				Severity: "error",
			}},
		})
		return 2
	}

	for stem, v := range rules {
		rulers := map[string]json.RawMessage{}
		err := json.Unmarshal(v, &rulers)
		if err != nil {
			_ = emitOutput(format, outputPath, validationReport{Summary: validationSummary{}, Issues: []validationIssue{{Message: fmt.Sprintf("error unmarshalling rulers: error=%v", err), Severity: "error"}}})
			return 2
		}

		for k, v := range rulers {
			switch k {
			case "exists":
				var exists []csvons.Exists
				err := json.Unmarshal(v, &exists)
				if err != nil {
					_ = emitOutput(format, outputPath, validationReport{Summary: validationSummary{}, Issues: []validationIssue{{Message: fmt.Sprintf("error unmarshalling exists: error=%v", err), Severity: "error"}}})
					return 2
				}
				csvons.ExistsTest(stem, exists, metadata)

			case "unique":
				var unique csvons.Unique
				err := json.Unmarshal(v, &unique)
				if err != nil {
					_ = emitOutput(format, outputPath, validationReport{Summary: validationSummary{}, Issues: []validationIssue{{Message: fmt.Sprintf("error unmarshalling unique: error=%v", err), Severity: "error"}}})
					return 2
				}
				csvons.UniqueTest(stem, &unique, metadata)

			case "vtype":
				var vtype []csvons.VType
				err := json.Unmarshal(v, &vtype)
				if err != nil {
					_ = emitOutput(format, outputPath, validationReport{Summary: validationSummary{}, Issues: []validationIssue{{Message: fmt.Sprintf("error unmarshalling vtype: error=%v", err), Severity: "error"}}})
					return 2
				}
				csvons.VTypeTest(stem, vtype, metadata)
			default:
				_ = emitOutput(format, outputPath, validationReport{Summary: validationSummary{}, Issues: []validationIssue{{Message: fmt.Sprintf("unknown key %s", k), Severity: "error"}}})
				return 2
			}
		}
	}

	durationMs := time.Since(startAt).Milliseconds()
	report := validationReport{
		Summary: validationSummary{
			FilesChecked: len(rules),
			Passed:       len(rules),
			Failed:       0,
			DurationMS:   durationMs,
		},
		Issues: []validationIssue{},
	}

	if err := emitOutput(format, outputPath, report); err != nil {
		log.Printf("error writing output: %v", err)
		return 2
	}
	return 0
}

func emitOutput(format, outputPath string, report validationReport) error {
	var out []byte
	switch format {
	case "json":
		data, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return err
		}
		out = append(data, '\n')
	default:
		var b bytes.Buffer
		if len(report.Issues) > 0 {
			for _, issue := range report.Issues {
				fmt.Fprintf(&b, "[%s] %s\n", issue.Severity, issue.Message)
			}
		} else {
			fmt.Fprintf(&b, "Validation succeeded: files_checked=%d passed=%d failed=%d duration_ms=%d\n",
				report.Summary.FilesChecked,
				report.Summary.Passed,
				report.Summary.Failed,
				report.Summary.DurationMS,
			)
		}
		out = b.Bytes()
	}

	if outputPath != "" {
		return os.WriteFile(outputPath, out, 0o644)
	}
	_, err := os.Stdout.Write(out)
	return err
}
