/**
 * Engine-specific types for IncludeKit mock engine
 * These types are NOT in the universal format schema
 */

import type {
  Statement,
  Mutation,
  Dependencies
} from '@includekit/spec';

/**
 * Application schema (engine-specific)
 */
export interface AppSchema {
  version: number;
  models: Array<{
    name: string;
    id: { kind: string };
    relations?: Array<{
      name: string;
      target: string;
      kind: string;
    }>;
  }>;
}

/**
 * Request to add a query and track its dependencies
 */
export interface AddQueryRequest {
  shape: Statement;
  result_hint?: Record<string, any[]>;
}

/**
 * Response from addQuery
 */
export interface AddQueryResponse {
  shape_id: string;
  dependencies: Dependencies;
}

/**
 * Response from computeShapeId
 */
export interface ShapeIdResponse {
  shape_id: string;
}

/**
 * Response from invalidate
 */
export interface InvalidateResponse {
  evict: string[];
}

/**
 * Request to explain invalidation
 */
export interface ExplainRequest {
  mutation: Mutation;
  shape_id: string;
}

/**
 * Response from explainInvalidation
 */
export interface ExplainResponse {
  invalidate: boolean;
  reasons: string[];
}

/**
 * Version information
 */
export interface VersionInfo {
  core: string;
  contract: string;
  abi: string;
}

/**
 * Engine interface matching WASM exports
 */
export interface IIncludeKitEngine {
  setSchema(schema: AppSchema): void;
  computeShapeId(statement: Statement): ShapeIdResponse;
  addQuery(request: AddQueryRequest): AddQueryResponse;
  invalidate(mutation: Mutation): InvalidateResponse;
  explainInvalidation(request: ExplainRequest): ExplainResponse;
  reset(): void;
  getVersion(): VersionInfo;
}
