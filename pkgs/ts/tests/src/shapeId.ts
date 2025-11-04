/**
 * Compute shapeId from canonical JSON
 */

import { createHash } from 'crypto';
import { canonicalizeQueryShape } from './canonicalize.js';

export function computeShapeId(canonicalJson: string): string {
  const hash = createHash('sha256').update(canonicalJson, 'utf8').digest('hex');
  return 's_' + hash;
}

export function computeQueryShapeId(shape: any): string {
  const canonical = canonicalizeQueryShape(shape);
  return computeShapeId(canonical);
}
