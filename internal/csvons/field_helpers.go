package csvons

// requiredFieldValues validates a generated field expression and ensures
// value extraction channel is available. It aborts via failf on invalid input.
func requiredFieldValues(fieldExpr FieldExpr, fieldName string, fields []string, records [][]string) <-chan string {
	if fieldExpr == nil {
		failf("field expression [%s] is nil", fieldName)
		return nil
	}
	vals := fieldExpr.FieldValue(fields, records)
	if vals == nil {
		failf("field expression [%s] cannot resolve values", fieldName)
		return nil
	}
	return vals
}
