package csvons

import "fmt"

// ValidationError carries structured context about a validation or runtime
// failure so callers can turn recovered panics into machine-readable reports.
type ValidationError struct {
	File     string
	Rule     string
	Field    string
	Row      *int
	Value    string
	Message  string
	Severity string
	Code     int
}

func (e ValidationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "csvons validation failed"
}

// ExitCode returns the CLI exit code that should be used when this error is
// surfaced to callers.
func (e ValidationError) ExitCode() int {
	if e.Code == 2 {
		return 2
	}
	return 1
}

// ValidationContext carries optional issue metadata for structured failures.
type ValidationContext struct {
	File     string
	Rule     string
	Field    string
	Row      *int
	Value    string
	Severity string
}

func (c ValidationContext) validationError(code int, format string, args ...any) ValidationError {
	severity := c.Severity
	if severity == "" {
		severity = "error"
	}

	return ValidationError{
		File:     c.File,
		Rule:     c.Rule,
		Field:    c.Field,
		Row:      c.Row,
		Value:    c.Value,
		Message:  fmt.Sprintf(format, args...),
		Severity: severity,
		Code:     code,
	}
}

func failValidation(ctx ValidationContext, format string, args ...any) {
	panic(ctx.validationError(1, format, args...))
}

func failRuntime(ctx ValidationContext, format string, args ...any) {
	panic(ctx.validationError(2, format, args...))
}

// failf aborts validation flow without terminating the whole process.
// Callers can recover panic values and convert them to structured output.
func failf(format string, args ...any) {
	failRuntime(ValidationContext{}, format, args...)
}

func csvFileName(stem string, metadata *Metadata) string {
	if metadata == nil {
		return stem
	}
	return stem + metadata.Extension
}

func rowPointer(row int) *int {
	rowCopy := row
	return &rowCopy
}
