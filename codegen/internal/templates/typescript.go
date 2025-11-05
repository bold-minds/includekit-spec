package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// These are the same templates from the previous TypeScript codegen
// but now embedded in the Go tool

func WriteTypeScriptValidators(dir, schemaPath string) error {
	schemaFile := schemaPath
	if strings.Contains(schemaPath, "/") {
		schemaFile = schemaPath[strings.LastIndex(schemaPath, "/")+1:]
	}

	content := fmt.Sprintf(`/**
 * Runtime validators for IncludeKit Universal Format
 * Auto-generated from schema/%s
 */

import type {`, schemaFile) + `
  QueryShape,
  MutationEvent,
  Dependencies,
  FilterSpec,
  FilterAtom,
  OrderBySpec,
  Scalar,
} from '@includekit/spec';

export class ValidationError extends Error {
  constructor(message: string, public path: string = '') {
    super(message);
    this.name = 'ValidationError';
  }
}

function isScalar(value: any): value is Scalar {
  if (value === null) return true;
  const type = typeof value;
  if (type === 'string' || type === 'number' || type === 'boolean') return true;
  if (type === 'object' && value.json !== undefined && Object.keys(value).length === 1) return true;
  return false;
}

function validateFilterAtom(atom: any, path: string = 'filterAtom'): asserts atom is FilterAtom {
  if (typeof atom !== 'object' || atom === null) {
    throw new ValidationError('FilterAtom must be an object', path);
  }
  if (typeof atom.field !== 'string' || atom.field.length === 0) {
    throw new ValidationError('FilterAtom.field must be a non-empty string', ` + "`${path}.field`" + `);
  }
  if (typeof atom.op !== 'string') {
    throw new ValidationError('FilterAtom.op must be a string', ` + "`${path}.op`" + `);
  }

  const validOps = [
    'eq', 'ne', 'in', 'notIn', 'isNull',
    'gt', 'gte', 'lt', 'lte', 'between',
    'contains', 'startsWith', 'endsWith',
    'like', 'ilike', 'regex',
    'has', 'hasSome', 'hasEvery', 'jsonContains',
    'lenEq', 'lenGt', 'lenLt', 'exists'
  ];

  if (!validOps.includes(atom.op) && !atom.op.startsWith('custom:')) {
    throw new ValidationError(` + "`Invalid operator: ${atom.op}`" + `, ` + "`${path}.op`" + `);
  }

  if (atom.value !== undefined && !isScalar(atom.value) && !Array.isArray(atom.value) && typeof atom.value !== 'object') {
    throw new ValidationError('FilterAtom.value must be a Scalar, array, or object', ` + "`${path}.value`" + `);
  }
}

function validateFilterSpec(spec: any, path: string = 'filterSpec'): asserts spec is FilterSpec {
  if (typeof spec !== 'object' || spec === null) {
    throw new ValidationError('FilterSpec must be an object', path);
  }

  if (spec.and && Array.isArray(spec.and)) {
    spec.and.forEach((s: any, i: number) => validateFilterSpec(s, ` + "`${path}.and[${i}]`" + `));
  }
  if (spec.or && Array.isArray(spec.or)) {
    spec.or.forEach((s: any, i: number) => validateFilterSpec(s, ` + "`${path}.or[${i}]`" + `));
  }
  if (spec.not) {
    validateFilterSpec(spec.not, ` + "`${path}.not`" + `);
  }
  if (spec.atoms && Array.isArray(spec.atoms)) {
    spec.atoms.forEach((a: any, i: number) => validateFilterAtom(a, ` + "`${path}.atoms[${i}]`" + `));
  }
}

function validateOrderBy(orderBy: any, path: string = 'orderBy'): asserts orderBy is OrderBySpec {
  if (typeof orderBy !== 'object' || orderBy === null) {
    throw new ValidationError('OrderBySpec must be an object', path);
  }
  if (typeof orderBy.field !== 'string' || orderBy.field.length === 0) {
    throw new ValidationError('OrderBySpec.field must be a non-empty string', ` + "`${path}.field`" + `);
  }
  if (orderBy.direction !== 'asc' && orderBy.direction !== 'desc') {
    throw new ValidationError('OrderBySpec.direction must be "asc" or "desc"', ` + "`${path}.direction`" + `);
  }
}

export function validateQueryShape(shape: any): asserts shape is QueryShape {
  if (typeof shape !== 'object' || shape === null) {
    throw new ValidationError('QueryShape must be an object', 'queryShape');
  }
  if (typeof shape.model !== 'string' || shape.model.length === 0) {
    throw new ValidationError('QueryShape.model must be a non-empty string', 'queryShape.model');
  }

  if (shape.where) {
    validateFilterSpec(shape.where, 'queryShape.where');
  }
  if (shape.orderBy && Array.isArray(shape.orderBy)) {
    shape.orderBy.forEach((o: any, i: number) => validateOrderBy(o, ` + "`queryShape.orderBy[${i}]`" + `));
  }
  if (shape.take !== undefined && (typeof shape.take !== 'number' || !Number.isInteger(shape.take))) {
    throw new ValidationError('QueryShape.take must be an integer', 'queryShape.take');
  }
  if (shape.skip !== undefined && (typeof shape.skip !== 'number' || !Number.isInteger(shape.skip))) {
    throw new ValidationError('QueryShape.skip must be an integer', 'queryShape.skip');
  }
}

export function validateMutationEvent(event: any): asserts event is MutationEvent {
  if (typeof event !== 'object' || event === null) {
    throw new ValidationError('MutationEvent must be an object', 'mutationEvent');
  }
  if (!Array.isArray(event.changes)) {
    throw new ValidationError('MutationEvent.changes must be an array', 'mutationEvent.changes');
  }

  event.changes.forEach((change: any, i: number) => {
    if (typeof change !== 'object' || change === null) {
      throw new ValidationError(` + "`Change must be an object`" + `, ` + "`mutationEvent.changes[${i}]`" + `);
    }
    if (!['create', 'update', 'delete', 'link', 'unlink'].includes(change.op || change.kind)) {
      throw new ValidationError(` + "`Invalid change operation`" + `, ` + "`mutationEvent.changes[${i}].op`" + `);
    }
  });
}

export function validateDependencies(deps: any): asserts deps is Dependencies {
  if (typeof deps !== 'object' || deps === null) {
    throw new ValidationError('Dependencies must be an object', 'dependencies');
  }
  if (typeof deps.shapeId !== 'string' || !/^s_[0-9a-f]{64}$/.test(deps.shapeId)) {
    throw new ValidationError('Dependencies.shapeId must match pattern ^s_[0-9a-f]{64}$', 'dependencies.shapeId');
  }
  if (typeof deps.records !== 'object' || deps.records === null) {
    throw new ValidationError('Dependencies.records must be an object', 'dependencies.records');
  }
  if (!Array.isArray(deps.filterBounds)) {
    throw new ValidationError('Dependencies.filterBounds must be an array', 'dependencies.filterBounds');
  }
  if (!Array.isArray(deps.relationBounds)) {
    throw new ValidationError('Dependencies.relationBounds must be an array', 'dependencies.relationBounds');
  }
}
`

	return os.WriteFile(filepath.Join(dir, "validators.ts"), []byte(content), 0644)
}

func WriteTypeScriptCanonicalize(dir string) error {
	content := `/**
 * JSON Canonicalization Scheme (JCS) implementation
 * RFC 8785: https://tools.ietf.org/html/rfc8785
 */

export function canonicalize(obj: any): string {
  return JSON.stringify(obj, canonicalReplacer);
}

function canonicalReplacer(_key: string, value: any): any {
  if (value && typeof value === 'object' && !Array.isArray(value)) {
    // Sort object keys
    const sorted: Record<string, any> = {};
    Object.keys(value)
      .sort()
      .forEach((k) => {
        sorted[k] = value[k];
      });
    return sorted;
  }
  return value;
}

export function canonicalizeQueryShape(shape: any): string {
  // Remove diagnostic fields before canonicalization
  const cleaned = JSON.parse(JSON.stringify(shape));
  delete cleaned.orm;
  delete cleaned.adapterVersion;
  return canonicalize(cleaned);
}
`

	return os.WriteFile(filepath.Join(dir, "canonicalize.ts"), []byte(content), 0644)
}

func WriteTypeScriptShapeId(dir string) error {
	content := `/**
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
`

	return os.WriteFile(filepath.Join(dir, "shapeId.ts"), []byte(content), 0644)
}

func WriteTypeScriptIndex(dir string) error {
	content := `export * from './validators.js';
export * from './canonicalize.js';
export * from './shapeId.js';
`

	return os.WriteFile(filepath.Join(dir, "index.ts"), []byte(content), 0644)
}
