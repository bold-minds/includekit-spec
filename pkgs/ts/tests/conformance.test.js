import { test } from 'node:test';
import { strict as assert } from 'node:assert';
import { readFile } from 'fs/promises';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import {
  canonicalizeQueryShape,
  computeShapeId,
  validateStatement,
} from './dist/index.js';

const __dirname = dirname(fileURLToPath(import.meta.url));

test('conformance: query shapes produce expected canonical JSON and shapeId', async () => {
  const vectorsPath = join(__dirname, '..', '..', '..', 'tools', 'tests', 'vectors', 'query-shapes.json');
  const vectorsData = await readFile(vectorsPath, 'utf-8');
  const vectors = JSON.parse(vectorsData);

  for (const vector of vectors) {
    await test(`vector: ${vector.name}`, () => {
      // Validate the shape
      validateStatement(vector.shape);

      // Canonicalize
      const canonical = canonicalizeQueryShape(vector.shape);
      
      // CRITICAL: Compare against expected canonical JSON
      assert.equal(canonical, vector.expectedCanonical,
        `Canonical JSON must match expected for ${vector.name}`);
      
      // Compute shapeId
      const shapeId = computeShapeId(canonical);
      
      // CRITICAL: Compare against expected shapeId
      assert.equal(shapeId, vector.expectedShapeId,
        `ShapeId must match expected for ${vector.name}`);
      
      // Basic format checks
      assert.ok(shapeId.startsWith('s_'), 'shapeId should start with s_');
      assert.equal(shapeId.length, 66, 'shapeId should be 66 characters (s_ + 64 hex)');

      // Verify determinism: same shape should produce same canonical JSON
      const canonical2 = canonicalizeQueryShape(vector.shape);
      assert.equal(canonical, canonical2, 'canonicalization should be deterministic');

      const shapeId2 = computeShapeId(canonical2);
      assert.equal(shapeId, shapeId2, 'shapeId should be deterministic');
    });
  }
});

test('conformance: validation catches invalid shapes', () => {
  assert.throws(() => {
    validateStatement({ query: { model: '' } }); // empty model
  }, /model must be a non-empty string/);

  assert.throws(() => {
    validateStatement({ query: { model: 'Post', order_by: [{ field: '' }] } }); // empty field
  }, /field must be a non-empty string/);
});
