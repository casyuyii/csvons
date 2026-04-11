import type { VTypeRule, Metadata, ValidationError } from '../types.ts';
import { csvFileName, makeValidationError } from '../types.ts';
import { generateFieldExpr } from '../field-expr.ts';
import '../field-occurrences.ts';
import { requiredFieldOccurrences, isValidationError } from '../field-helpers.ts';
import { readCsvFile } from '../io/csv-reader.ts';

/**
 * Valid boolean values matching Go's strconv.ParseBool behavior.
 */
const VALID_BOOLS = new Set([
  '1',
  't',
  'T',
  'TRUE',
  'true',
  'True',
  '0',
  'f',
  'F',
  'FALSE',
  'false',
  'False',
]);

/**
 * Validates that values conform to expected types and optional ranges.
 * Collects all validation errors.
 */
export function vtypeTest(
  stem: string,
  rules: VTypeRule[],
  metadata: Metadata,
  basePath?: string,
): ValidationError[] {
  const errors: ValidationError[] = [];
  const fileName = csvFileName(stem, metadata);

  if (rules.length === 0) {
    errors.push(
      makeValidationError(
        { file: fileName, rule: 'vtype' },
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
        { file: fileName, rule: 'vtype' },
        2,
        `name_index [${metadata.name_index}] is less than 0`,
      ),
    );
    return errors;
  }

  if (metadata.data_index <= metadata.name_index) {
    errors.push(
      makeValidationError(
        { file: fileName, rule: 'vtype' },
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
        { file: fileName, rule: 'vtype' },
        2,
        `src_records length [${srcRecords?.length ?? 0}] <= data_index [${metadata.data_index}]`,
      ),
    );
    return errors;
  }

  const srcFields = srcRecords[metadata.name_index]!;

  for (const vtype of rules) {
    const fieldExpr = generateFieldExpr(metadata, vtype.field);
    const fieldVals = requiredFieldOccurrences(
      fieldExpr,
      vtype.field,
      srcFields,
      srcRecords,
      { file: fileName, rule: 'vtype', field: vtype.field },
    );

    if (isValidationError(fieldVals)) {
      errors.push(fieldVals);
      continue;
    }

    const checkedCache = new Set<string>();

    for (const occurrence of fieldVals) {
      const fieldVal = occurrence.value;

      if (checkedCache.has(fieldVal)) continue;

      switch (vtype.type) {
        case 'int': {
          // Match Go's strconv.ParseInt behavior: accepts leading zeros, optional sign
          const parsed = parseInt(fieldVal, 10);
          if (isNaN(parsed) || !/^[+-]?\d+$/.test(fieldVal)) {
            errors.push(
              makeValidationError(
                {
                  file: fileName,
                  rule: 'vtype',
                  field: vtype.field,
                  row: occurrence.row,
                  value: fieldVal,
                },
                1,
                `src_field [${vtype.field}] value [${fieldVal}] is not an int`,
              ),
            );
            continue;
          }
          if (vtype.range) {
            if (parsed > vtype.range.max || parsed < vtype.range.min) {
              errors.push(
                makeValidationError(
                  {
                    file: fileName,
                    rule: 'vtype',
                    field: vtype.field,
                    row: occurrence.row,
                    value: fieldVal,
                  },
                  1,
                  `src_field [${vtype.field}] value [${fieldVal}] is not in the range [${vtype.range.min}, ${vtype.range.max}]`,
                ),
              );
              continue;
            }
          }
          break;
        }

        case 'float64': {
          const parsed = parseFloat(fieldVal);
          if (isNaN(parsed)) {
            errors.push(
              makeValidationError(
                {
                  file: fileName,
                  rule: 'vtype',
                  field: vtype.field,
                  row: occurrence.row,
                  value: fieldVal,
                },
                1,
                `src_field [${vtype.field}] value [${fieldVal}] is not a float64`,
              ),
            );
            continue;
          }
          if (vtype.range) {
            if (parsed > vtype.range.max || parsed < vtype.range.min) {
              errors.push(
                makeValidationError(
                  {
                    file: fileName,
                    rule: 'vtype',
                    field: vtype.field,
                    row: occurrence.row,
                    value: fieldVal,
                  },
                  1,
                  `src_field [${vtype.field}] value [${fieldVal}] is not in the range [${vtype.range.min}, ${vtype.range.max}]`,
                ),
              );
              continue;
            }
          }
          break;
        }

        case 'bool': {
          if (!VALID_BOOLS.has(fieldVal)) {
            errors.push(
              makeValidationError(
                {
                  file: fileName,
                  rule: 'vtype',
                  field: vtype.field,
                  row: occurrence.row,
                  value: fieldVal,
                },
                1,
                `src_field [${vtype.field}] value [${fieldVal}] is not a bool`,
              ),
            );
            continue;
          }
          break;
        }

        default: {
          errors.push(
            makeValidationError(
              {
                file: fileName,
                rule: 'vtype',
                field: vtype.field,
                row: occurrence.row,
                value: fieldVal,
              },
              2,
              `src_field [${vtype.field}] value [${fieldVal}] is not a valid type`,
            ),
          );
          continue;
        }
      }

      checkedCache.add(fieldVal);
    }
  }

  return errors;
}
