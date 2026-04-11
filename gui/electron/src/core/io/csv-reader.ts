import * as fs from 'node:fs';
import * as path from 'node:path';
import Papa from 'papaparse';
import type { Metadata } from '../types.ts';

/**
 * Reads a CSV file identified by its stem (base name) and metadata.
 * Returns 2D array of strings, or null on error.
 */
export function readCsvFile(
  stem: string,
  metadata: Metadata,
  basePath?: string,
): string[][] | null {
  const folder = basePath
    ? path.resolve(basePath, metadata.csv_file_folder)
    : metadata.csv_file_folder;
  const fullPath = path.join(folder, stem + metadata.extension);

  let content: string;
  try {
    content = fs.readFileSync(fullPath, 'utf-8');
  } catch {
    return null;
  }

  const result = Papa.parse<string[]>(content, {
    header: false,
    skipEmptyLines: false,
  });

  if (result.errors.length > 0) {
    return null;
  }

  return result.data;
}
