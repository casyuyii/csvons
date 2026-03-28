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
	File     string `json:"file,omitempty"`
	Rule     string `json:"rule,omitempty"`
	Field    string `json:"field,omitempty"`
	Row      *int   `json:"row,omitempty"`
	Value    string `json:"value,omitempty"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type validationReport struct {
	SchemaVersion string            `json:"schema_version,omitempty"`
	Summary       validationSummary `json:"summary"`
	Issues        []validationIssue `json:"issues"`
}

const reportSchemaVersion = "csvons.validation_report.v1"

func main() {
	os.Exit(run())
}

func run() int {
	return runWithArgs(os.Args[1:])
}

func runWithArgs(args []string) (code int) {
	var format string
	var outputPath string
	var rules map[string]json.RawMessage

	flags := flag.NewFlagSet("csvons", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	flags.StringVar(&format, "format", "text", "output format: text or json")
	flags.StringVar(&outputPath, "output", "", "optional output file path")
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--format text|json] [--output <path>] <ruler.json>\n", flags.Name())
		fmt.Fprintf(os.Stderr, "\nValidate CSV files against constraint rules defined in a JSON configuration file.\n")
		flags.PrintDefaults()
	}

	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() < 1 {
		flags.Usage()
		return 2
	}
	if format != "text" && format != "json" {
		fmt.Fprintf(os.Stderr, "invalid --format value %q; expected text or json\n", format)
		return 2
	}

	configFileName := flags.Arg(0)
	startAt := time.Now()
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}

		issue, issueCode := validationIssueFromRecovered(recovered)
		_ = emitOutput(format, outputPath, validationReport{
			Summary: validationSummary{
				FilesChecked: len(rules),
				Passed:       0,
				Failed:       1,
				DurationMS:   time.Since(startAt).Milliseconds(),
			},
			Issues: []validationIssue{issue},
		})
		code = issueCode
	}()

	metadataFile := func(stem string, metadata *csvons.Metadata) string {
		if metadata == nil {
			return stem
		}
		return stem + metadata.Extension
	}

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

	for stem, rawRules := range rules {
		rulers := map[string]json.RawMessage{}
		if err := json.Unmarshal(rawRules, &rulers); err != nil {
			_ = emitOutput(format, outputPath, validationReport{
				Summary: validationSummary{},
				Issues: []validationIssue{{
					File:     metadataFile(stem, metadata),
					Message:  fmt.Sprintf("error unmarshalling rulers: error=%v", err),
					Severity: "error",
				}},
			})
			return 2
		}

		for ruleName, rawRule := range rulers {
			switch ruleName {
			case "exists":
				var exists []csvons.Exists
				if err := json.Unmarshal(rawRule, &exists); err != nil {
					_ = emitOutput(format, outputPath, validationReport{
						Summary: validationSummary{},
						Issues: []validationIssue{{
							File:     metadataFile(stem, metadata),
							Rule:     ruleName,
							Message:  fmt.Sprintf("error unmarshalling exists: error=%v", err),
							Severity: "error",
						}},
					})
					return 2
				}
				csvons.ExistsTest(stem, exists, metadata)

			case "unique":
				var unique csvons.Unique
				if err := json.Unmarshal(rawRule, &unique); err != nil {
					_ = emitOutput(format, outputPath, validationReport{
						Summary: validationSummary{},
						Issues: []validationIssue{{
							File:     metadataFile(stem, metadata),
							Rule:     ruleName,
							Message:  fmt.Sprintf("error unmarshalling unique: error=%v", err),
							Severity: "error",
						}},
					})
					return 2
				}
				csvons.UniqueTest(stem, &unique, metadata)

			case "vtype":
				var vtype []csvons.VType
				if err := json.Unmarshal(rawRule, &vtype); err != nil {
					_ = emitOutput(format, outputPath, validationReport{
						Summary: validationSummary{},
						Issues: []validationIssue{{
							File:     metadataFile(stem, metadata),
							Rule:     ruleName,
							Message:  fmt.Sprintf("error unmarshalling vtype: error=%v", err),
							Severity: "error",
						}},
					})
					return 2
				}
				csvons.VTypeTest(stem, vtype, metadata)

			default:
				_ = emitOutput(format, outputPath, validationReport{
					Summary: validationSummary{},
					Issues: []validationIssue{{
						File:     metadataFile(stem, metadata),
						Rule:     ruleName,
						Message:  fmt.Sprintf("unknown key %s", ruleName),
						Severity: "error",
					}},
				})
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

func validationIssueFromRecovered(recovered any) (validationIssue, int) {
	switch v := recovered.(type) {
	case csvons.ValidationError:
		return validationIssue{
			File:     v.File,
			Rule:     v.Rule,
			Field:    v.Field,
			Row:      v.Row,
			Value:    v.Value,
			Message:  v.Error(),
			Severity: v.Severity,
		}, v.ExitCode()
	case error:
		return validationIssue{
			Message:  v.Error(),
			Severity: "error",
		}, 2
	default:
		return validationIssue{
			Message:  fmt.Sprint(recovered),
			Severity: "error",
		}, 2
	}
}

func emitOutput(format, outputPath string, report validationReport) error {
	var out []byte
	switch format {
	case "json":
		if report.SchemaVersion == "" {
			report.SchemaVersion = reportSchemaVersion
		}
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
