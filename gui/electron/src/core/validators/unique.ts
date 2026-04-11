import type { UniqueRule, Metadata, ValidationError } from '../types.ts';
import { csvFileName, makeValidationError } from '../types.ts';
import { generateFieldExpr } from '../field-expr.ts';
import '../field-occurrences.ts';
import { requiredFieldOccurrences, isValidationError } from '../field-helpers.ts';
import { readCsvFile } from '../io/csv-reader.ts';

/**
 * Validates that all values in specified columns are unique.
 * Collects all validation errors.
 */
export function uniqueTest(
  stem: string,
  ruler: UniqueRule,
  metadata: Metadata,
  basePath?: string,
): ValidationError[] {
  const errors: ValidationError[] = [];
  const fileName = csvFileName(stem, metadata);

  if (!ruler.fields || ruler.fields.length === 0) {
    errors.push(
      makeValidationError(
        { file: fileName, rule: 'unique' },
        2,
        `ruler is empty`,
      ),
    );
    return errors;
  }

  // Validate metadata indices
  if (metadata.name_index < 0) {
    errors.push(
      makeValidationError(
        { file: fileName, rule: 'unique' },
        2,
        `name_index [${metadata.name_index}] is less than 0`,
      ),
    );
    return errors;
  }

  if (metadata.data_index <= metadata.name_index) {
    errors.push(
      makeValidationError(
        { file: fileName, rule: 'unique' },
        2,
        `data_index [${metadata.data_index}] is less than or equal to name_index [${metadata.name_index}]`,
      ),
    );
    return errors;
  }

  // Read CSV
  const srcRecords = readCsvFile(stem, metadata, basePath);
  if (!srcRecords || srcRecords.length <= metadata.data_index) {
    errors.push(
      makeValidationError(
        { file: fileName, rule: 'unique' },
        2,
        `src_records length [${srcRecords?.length ?? 0}] <= data_index [${metadata.data_index}]`,
      ),
    );
    return errors;
  }

  const srcFields = srcRecords[metadata.name_index]!;

  for (const fieldName of ruler.fields) {
    const fieldExpr = generateFieldExpr(metadata, fieldName);
    const fieldVals = requiredFieldOccurrences(
      fieldExpr,
      fieldName,
      srcFields,
      srcRecords,
      { file: fileName, rule: 'unique', field: fieldName },
    );

    if (isValidationError(fieldVals)) {
      errors.push(fieldVals);
      continue;
    }

    const existingFields = new Map<string, number>();
    for (const occurrence of fieldVals) {
      const fieldVal = occurrence.value;
      const count = (existingFields.get(fieldVal) ?? 0) + 1;
      existingFields.set(fieldVal, count);

      if (count > 1) {
        errors.push(
          makeValidationError(
            {
              file: fileName,
              rule: 'unique',
              field: fieldName,
              row: occurrence.row,
              value: fieldVal,
            },
            1,
            `src_field [${fieldName}] value [${fieldVal}] already exists`,
          ),
        );
      }
    }
  }

  return errors;
}
