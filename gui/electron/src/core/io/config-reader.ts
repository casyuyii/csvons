import * as fs from 'node:fs';
import type { Metadata } from '../types.ts';

const METADATA_KEY = 'csvons_metadata';

export interface ConfigResult {
  rules: Record<string, Record<string, unknown>>;
  metadata: Metadata;
}

/**
 * Reads and parses a ruler JSON configuration file.
 * Extracts the metadata section and returns remaining keys as rule definitions.
 * Returns null if the file cannot be read, parsed, or lacks valid metadata.
 */
export function readConfigFile(configPath: string): ConfigResult | null {
  let data: string;
  try {
    data = fs.readFileSync(configPath, 'utf-8');
  } catch {
    return null;
  }

  let cfg: Record<string, unknown>;
  try {
    cfg = JSON.parse(data) as Record<string, unknown>;
  } catch {
    return null;
  }

  const metadataRaw = cfg[METADATA_KEY];
  if (!metadataRaw || typeof metadataRaw !== 'object') {
    return null;
  }

  const metadata = metadataRaw as Metadata;

  // Remove metadata key, leaving only file stem rules
  delete cfg[METADATA_KEY];

  const rules: Record<string, Record<string, unknown>> = {};
  for (const [key, value] of Object.entries(cfg)) {
    if (value && typeof value === 'object') {
      rules[key] = value as Record<string, unknown>;
    }
  }

  return { rules, metadata };
}
