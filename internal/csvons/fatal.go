package csvons

import "fmt"

// failf aborts validation flow without terminating the whole process.
// Callers can recover panic values and convert them to structured output.
func failf(format string, args ...any) {
	panic(fmt.Errorf(format, args...))
}
