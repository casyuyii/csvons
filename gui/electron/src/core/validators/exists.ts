import type { ExistsRule, Metadata, ValidationError } from '../types.ts';
import { csvFileName, makeValidationError } from '../types.ts';
import { generateFieldExpr } from '../field-expr.ts';
import '../field-occurrences.ts';
import { requiredFieldOccurrences, isValidationError } from '../field-helpers.ts';
import { readCsvFile } from '../io/csv-reader.ts';

/**
 * Validates that values in source columns exist in destination columns.
 * Collects all validation errors instead of stopping on first failure.
 */
export function existsTest(
  stem: string,
  rules: ExistsRule[],
  metadata: Metadata,
  basePath?: string,
): ValidationError[] {
  const errors: ValidationError[] = [];
  const fileName = csvFileName(stem, metadata);

  if (rules.length === 0) {
    errors.push(
      makeValidationError(
        { file: fileName, rule: 'exists' },
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
        { file: fileName, rule: 'exists' },
        2,
        `name_index [${metadata.name_index}] is less than 0`,
      ),
    );
    return errors;
  }

  if (metadata.data_index <= metadata.name_index) {
    errors.push(
      makeValidationError(
        { file: fileName, rule: 'exists' },
        2,
        `data_index [${metadata.data_index}] is less than or equal to name_index [${metadata.name_index}]`,
      ),
    );
    return errors;
  }

  // Read source CSV
  const srcRecords = readCsvFile(stem, metadata, basePath);
  if (!srcRecords || srcRecords.length <= metadata.data_index) {
    errors.push(
      makeValidationError(
        { file: fileName, rule: 'exists' },
        2,
        `src_records length [${srcRecords?.length ?? 0}] <= data_index [${metadata.data_index}]`,
      ),
    );
    return errors;
  }

  const srcFields = srcRecords[metadata.name_index]!;

  for (const exist of rules) {
    const dstFileName = csvFileName(exist.dst_file_stem, metadata);

    // Read destination CSV
    const dstRecords = readCsvFile(exist.dst_file_stem, metadata, basePath);
    if (!dstRecords || dstRecords.length <= metadata.data_index) {
      errors.push(
        makeValidationError(
          { file: dstFileName, rule: 'exists' },
          2,
          `dst_records length [${dstRecords?.length ?? 0}] <= data_index [${metadata.data_index}]`,
        ),
      );
      return errors;
    }

    const dstFields = dstRecords[metadata.name_index]!;

    for (const field of exist.fields) {
      // Create field expressions for source and destination
      const srcFieldExpr = generateFieldExpr(metadata, field.src);
      const srcFieldVals = requiredFieldOccurrences(
        srcFieldExpr,
        field.src,
        srcFields,
        srcRecords,
        { file: fileName, rule: 'exists', field: field.src },
      );

      if (isValidationError(srcFieldVals)) {
        errors.push(srcFieldVals);
        continue;
      }

      const dstFieldExpr = generateFieldExpr(metadata, field.dst);
      const dstFieldVals = requiredFieldOccurrences(
        dstFieldExpr,
        field.dst,
        dstFields,
        dstRecords,
        { file: dstFileName, rule: 'exists', field: field.dst },
      );

      if (isValidationError(dstFieldVals)) {
        errors.push(dstFieldVals);
        continue;
      }

      // Build set of all destination values for lookup
      const dstValueSet = new Set<string>();
      for (const dstOcc of dstFieldVals) {
        dstValueSet.add(dstOcc.value);
      }

      // Check each source value exists in destination
      const searchedFields = new Set<string>();
      for (const srcOccurrence of srcFieldVals) {
        const fieldVal = srcOccurrence.value;

        if (searchedFields.has(fieldVal)) continue;

        if (!dstValueSet.has(fieldVal)) {
          errors.push(
            makeValidationError(
              {
                file: fileName,
                rule: 'exists',
                field: field.src,
                row: srcOccurrence.row,
                value: fieldVal,
              },
              1,
              `src_field [${field.src}] value [${fieldVal}] not found in dst_records`,
            ),
          );
        }

        searchedFields.add(fieldVal);
      }
    }
  }

  return errors;
}
