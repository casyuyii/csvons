import * as path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

/** Root of the csvons repo (contains testdata/, ruler/, etc.) */
export const PROJECT_ROOT = path.resolve(__dirname, '../../..');

/** Path to testdata directory */
export const TESTDATA_DIR = path.join(PROJECT_ROOT, 'testdata');

/** Path to ruler directory */
export const RULER_DIR = path.join(PROJECT_ROOT, 'ruler');
