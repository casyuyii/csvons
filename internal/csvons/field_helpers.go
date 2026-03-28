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

func requiredFieldOccurrences(fieldExpr FieldExpr, fieldName string, fields []string, records [][]string, ctx ValidationContext) <-chan FieldOccurrence {
	ctx.Field = fieldName
	if fieldExpr == nil {
		failRuntime(ctx, "field expression [%s] is nil", fieldName)
		return nil
	}

	provider, ok := fieldExpr.(fieldOccurrenceProvider)
	if !ok {
		failRuntime(ctx, "field expression [%s] cannot resolve values", fieldName)
		return nil
	}

	occurrences := provider.FieldOccurrences(fields, records)
	if occurrences == nil {
		failRuntime(ctx, "field expression [%s] cannot resolve values", fieldName)
		return nil
	}
	return occurrences
}
