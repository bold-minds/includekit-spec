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
 * Manually maintained - must be updated when schema changes
 */

import type {`, schemaFile) + `
  Statement,
  Mutation,
  Dependencies,
  Filter,
  Condition,
  OrderBy,
} from '@includekit/spec';

export class ValidationError extends Error {
  constructor(message: string, public path: string = '') {
    super(message);
    this.name = 'ValidationError';
  }
}

function validateCondition(condition: any, path: string = 'condition'): asserts condition is Condition {
  if (typeof condition !== 'object' || condition === null) {
    throw new ValidationError('Condition must be an object', path);
  }
  if (typeof condition.field !== 'string' || condition.field.length === 0) {
    throw new ValidationError('Condition.field must be a non-empty string', ` + "`${path}.field`" + `);
  }
  if (typeof condition.op !== 'string') {
    throw new ValidationError('Condition.op must be a string', ` + "`${path}.op`" + `);
  }

  const validOps = [
    'eq', 'ne', 'in', 'notIn', 'isNull',
    'gt', 'gte', 'lt', 'lte', 'between',
    'contains', 'startsWith', 'endsWith',
    'like', 'ilike', 'regex',
    'has', 'hasSome', 'hasEvery', 'jsonContains',
    'lenEq', 'lenGt', 'lenLt', 'exists'
  ];

  const isCustomOp = condition.op.startsWith('custom:');
  if (!validOps.includes(condition.op) && !isCustomOp) {
    throw new ValidationError(` + "`Invalid operator: ${condition.op}`" + `, ` + "`${path}.op`" + `);
  }
  
  // value can be any JSON value - no type validation needed
}

function validateFilter(filter: any, path: string = 'filter'): asserts filter is Filter {
  if (typeof filter !== 'object' || filter === null) {
    throw new ValidationError('Filter must be an object', path);
  }

  if (filter.and && Array.isArray(filter.and)) {
    filter.and.forEach((f: any, i: number) => validateFilter(f, ` + "`${path}.and[${i}]`" + `));
  }
  if (filter.or && Array.isArray(filter.or)) {
    filter.or.forEach((f: any, i: number) => validateFilter(f, ` + "`${path}.or[${i}]`" + `));
  }
  if (filter.not) {
    validateFilter(filter.not, ` + "`${path}.not`" + `);
  }
  if (filter.conditions && Array.isArray(filter.conditions)) {
    filter.conditions.forEach((c: any, i: number) => validateCondition(c, ` + "`${path}.conditions[${i}]`" + `));
  }
}

function validateOrderBy(orderBy: any, path: string = 'orderBy'): asserts orderBy is OrderBy {
  if (typeof orderBy !== 'object' || orderBy === null) {
    throw new ValidationError('OrderBy must be an object', path);
  }
  if (typeof orderBy.field !== 'string' || orderBy.field.length === 0) {
    throw new ValidationError('OrderBy.field must be a non-empty string', ` + "`${path}.field`" + `);
  }
  // descending, nulls_first, case_sensitive are all booleans - no validation needed beyond type
}

export function validateStatement(statement: any): asserts statement is Statement {
  if (typeof statement !== 'object' || statement === null) {
    throw new ValidationError('Statement must be an object', 'statement');
  }

  if (statement.query) {
    if (typeof statement.query !== 'object' || statement.query === null) {
      throw new ValidationError('Statement.query must be an object', 'statement.query');
    }
    if (typeof statement.query.model !== 'string' || statement.query.model.length === 0) {
      throw new ValidationError('Statement.query.model must be a non-empty string', 'statement.query.model');
    }
    if (statement.query.where) {
      validateFilter(statement.query.where, 'statement.query.where');
    }
    if (statement.query.order_by && Array.isArray(statement.query.order_by)) {
      statement.query.order_by.forEach((o: any, i: number) => validateOrderBy(o, ` + "`statement.query.order_by[${i}]`" + `));
    }
    if (statement.query.limit !== undefined && (typeof statement.query.limit !== 'number' || !Number.isInteger(statement.query.limit))) {
      throw new ValidationError('Statement.query.limit must be an integer', 'statement.query.limit');
    }
    if (statement.query.offset !== undefined && (typeof statement.query.offset !== 'number' || !Number.isInteger(statement.query.offset))) {
      throw new ValidationError('Statement.query.offset must be an integer', 'statement.query.offset');
    }
  }

  if (statement.pagination) {
    if (typeof statement.pagination !== 'object' || statement.pagination === null) {
      throw new ValidationError('Statement.pagination must be an object', 'statement.pagination');
    }
    const hasForward = statement.pagination.first !== undefined || statement.pagination.after !== undefined;
    const hasBackward = statement.pagination.last !== undefined || statement.pagination.before !== undefined;
    if (hasForward && hasBackward) {
      throw new ValidationError('Cannot mix forward and backward pagination', 'statement.pagination');
    }
  }
}

export function validateMutation(mutation: any): asserts mutation is Mutation {
  if (typeof mutation !== 'object' || mutation === null) {
    throw new ValidationError('Mutation must be an object', 'mutation');
  }
  if (!Array.isArray(mutation.changes)) {
    throw new ValidationError('Mutation.changes must be an array', 'mutation.changes');
  }

  mutation.changes.forEach((change: any, i: number) => {
    if (typeof change !== 'object' || change === null) {
      throw new ValidationError(` + "`Change must be an object`" + `, ` + "`mutation.changes[${i}]`" + `);
    }
    if (!['insert', 'update', 'delete'].includes(change.action)) {
      throw new ValidationError(` + "`Invalid change action: must be insert, update, or delete`" + `, ` + "`mutation.changes[${i}].action`" + `);
    }
    if (typeof change.model !== 'string' || change.model.length === 0) {
      throw new ValidationError('Change.model must be a non-empty string', ` + "`mutation.changes[${i}].model`" + `);
    }
    
    // Validate based on action
    if (change.action === 'insert' && (!Array.isArray(change.set) || change.set.length === 0)) {
      throw new ValidationError('Insert requires non-empty set', ` + "`mutation.changes[${i}].set`" + `);
    }
    if (change.action === 'update' && (!Array.isArray(change.set) || change.set.length === 0)) {
      throw new ValidationError('Update requires non-empty set', ` + "`mutation.changes[${i}].set`" + `);
    }
    if (change.action === 'update' && !change.where) {
      throw new ValidationError('Update requires where clause', ` + "`mutation.changes[${i}].where`" + `);
    }
    if (change.action === 'delete' && !change.where) {
      throw new ValidationError('Delete requires where clause', ` + "`mutation.changes[${i}].where`" + `);
    }
  });
}

export function validateDependencies(deps: any): asserts deps is Dependencies {
  if (typeof deps !== 'object' || deps === null) {
    throw new ValidationError('Dependencies must be an object', 'dependencies');
  }
  if (typeof deps.shape_id !== 'string' || !/^s_[0-9a-f]{64}$/.test(deps.shape_id)) {
    throw new ValidationError('Dependencies.shape_id must match pattern ^s_[0-9a-f]{64}$', 'dependencies.shape_id');
  }
  if (typeof deps.records !== 'object' || deps.records === null) {
    throw new ValidationError('Dependencies.records must be an object', 'dependencies.records');
  }
  if (!Array.isArray(deps.filters)) {
    throw new ValidationError('Dependencies.filters must be an array', 'dependencies.filters');
  }
  if (!Array.isArray(deps.includes)) {
    throw new ValidationError('Dependencies.includes must be an array', 'dependencies.includes');
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
  delete cleaned.orm_version;
  delete cleaned.sdk_version;
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
