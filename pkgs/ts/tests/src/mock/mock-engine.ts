/**
 * Mock IncludeKit engine for testing without WASM
 */

import type {
  Statement,
  Mutation,
  Dependencies,
  Filter
} from '@includekit/spec';
import { computeQueryShapeId } from '../shapeId.js';
import type {
  IIncludeKitEngine,
  AppSchema,
  AddQueryRequest,
  AddQueryResponse,
  ShapeIdResponse,
  InvalidateResponse,
  ExplainRequest,
  ExplainResponse,
  VersionInfo
} from './interface.js';

export interface MockEngineConfig {
  /**
   * Custom shape ID generator (default: uses real shapeId computation)
   */
  shapeIdGenerator?: (statement: Statement) => string;
  
  /**
   * Eviction behavior:
   * - 'conservative': Evict all shapes for affected models
   * - 'custom': Use customEvictList
   */
  evictBehavior?: 'conservative' | 'custom';
  
  /**
   * Custom evict list (when evictBehavior = 'custom')
   */
  customEvictList?: string[];
  
  /**
   * Track all method calls for assertions
   */
  trackCalls?: boolean;
}

export interface MockEngineCalls {
  setSchema: Array<{ schema: AppSchema }>;
  computeShapeId: Array<{ statement: Statement }>;
  addQuery: Array<{ request: AddQueryRequest }>;
  invalidate: Array<{ mutation: Mutation }>;
  explainInvalidation: Array<{ request: ExplainRequest }>;
  reset: Array<Record<string, never>>;
  getVersion: Array<Record<string, never>>;
}

export class MockIncludeKitEngine implements IIncludeKitEngine {
  private schema?: AppSchema;
  private shapes = new Map<string, Dependencies>();
  private calls: MockEngineCalls;
  
  constructor(private config: MockEngineConfig = {}) {
    this.calls = this.initCalls();
  }
  
  private initCalls(): MockEngineCalls {
    return {
      setSchema: [],
      computeShapeId: [],
      addQuery: [],
      invalidate: [],
      explainInvalidation: [],
      reset: [],
      getVersion: []
    };
  }
  
  setSchema(schema: AppSchema): void {
    if (this.config.trackCalls) {
      this.calls.setSchema.push({ schema });
    }
    this.schema = schema;
  }
  
  computeShapeId(statement: Statement): ShapeIdResponse {
    if (this.config.trackCalls) {
      this.calls.computeShapeId.push({ statement });
    }
    
    let shapeId: string;
    if (this.config.shapeIdGenerator) {
      shapeId = this.config.shapeIdGenerator(statement);
    } else {
      // Use real shapeId computation from testkit
      shapeId = computeQueryShapeId(statement);
    }
    
    return { shape_id: shapeId };
  }
  
  addQuery(request: AddQueryRequest): AddQueryResponse {
    if (this.config.trackCalls) {
      this.calls.addQuery.push({ request });
    }
    
    const { shape_id } = this.computeShapeId(request.shape);
    
    // Build dependencies
    const dependencies: Dependencies = {
      shape_id,
      records: this.extractRecords(request),
      filters: this.extractFilters(request.shape),
      includes: request.shape.includes || []
    };
    
    // Store for invalidation checks
    this.shapes.set(shape_id, dependencies);
    
    return { shape_id, dependencies };
  }
  
  invalidate(mutation: Mutation): InvalidateResponse {
    if (this.config.trackCalls) {
      this.calls.invalidate.push({ mutation });
    }
    
    // Custom evict list (for testing)
    if (this.config.evictBehavior === 'custom' && this.config.customEvictList) {
      return { evict: this.config.customEvictList };
    }
    
    const evict: string[] = [];
    
    for (const [shapeId, deps] of this.shapes.entries()) {
      for (const change of mutation.changes) {
        const shouldEvict = this.shouldInvalidate(change, deps);
        if (shouldEvict) {
          evict.push(shapeId);
          break;
        }
      }
    }
    
    return { evict };
  }
  
  explainInvalidation(request: ExplainRequest): ExplainResponse {
    if (this.config.trackCalls) {
      this.calls.explainInvalidation.push({ request });
    }
    
    const deps = this.shapes.get(request.shape_id);
    if (!deps) {
      return { invalidate: false, reasons: [] };
    }
    
    const reasons: string[] = [];
    
    for (const change of request.mutation.changes) {
      // Check record membership
      if (deps.records[change.model] && deps.records[change.model].length > 0) {
        reasons.push('record_membership');
      }
      
      // Check filter dependencies
      if (deps.filters.length > 0) {
        for (const filter of deps.filters) {
          if (this.filterReferencesModel(filter, change.model)) {
            reasons.push('filter_dependency');
            break;
          }
        }
      }
      
      // Check relation dependencies
      if (deps.includes.length > 0) {
        for (const include of deps.includes) {
          if (include.query?.model === change.model) {
            reasons.push('relation_dependency');
            break;
          }
        }
      }
    }
    
    // Deduplicate reasons
    const uniqueReasons = [...new Set(reasons)];
    
    return {
      invalidate: uniqueReasons.length > 0,
      reasons: uniqueReasons
    };
  }
  
  reset(): void {
    if (this.config.trackCalls) {
      this.calls.reset.push({});
    }
    
    this.schema = undefined;
    this.shapes.clear();
    
    if (this.config.trackCalls) {
      this.calls = this.initCalls();
    }
  }
  
  getVersion(): VersionInfo {
    if (this.config.trackCalls) {
      this.calls.getVersion.push({});
    }
    
    return {
      core: 'mock-0.1.0',
      contract: '0.1.0',
      abi: '1'
    };
  }
  
  // Helpers
  
  private extractRecords(request: AddQueryRequest): Record<string, string[]> {
    if (!request.result_hint || !request.shape.query) {
      return {};
    }
    
    const records: Record<string, string[]> = {};
    const model = request.shape.query.model;
    
    if (request.result_hint[model]) {
      const rows = request.result_hint[model];
      const ids = rows
        .map(row => row?.id?.toString())
        .filter((id): id is string => id !== undefined);
      
      if (ids.length > 0) {
        records[model] = ids;
      }
    }
    
    return records;
  }
  
  private extractFilters(statement: Statement): Filter[] {
    const filters: Filter[] = [];
    
    if (statement.query?.where) {
      filters.push(statement.query.where);
    }
    
    if (statement.having) {
      filters.push(statement.having);
    }
    
    return filters;
  }
  
  private filterReferencesModel(_filter: Filter, _model: string): boolean {
    // Simplified: just check if any condition exists
    // Real engine would check field paths for relation references
    return !!_filter.conditions && _filter.conditions.length > 0;
  }
  
  private shouldInvalidate(change: any, deps: Dependencies): boolean {
    const behavior = this.config.evictBehavior || 'conservative';
    
    if (behavior === 'conservative') {
      // Conservative: evict if model is tracked
      return !!deps.records[change.model];
    }
    
    return false;
  }
  
  // Public API for test assertions
  
  /**
   * Get all tracked method calls (when trackCalls = true)
   */
  getCalls(): MockEngineCalls {
    return this.calls;
  }
  
  /**
   * Set custom evict list for next invalidate() call
   */
  setEvictList(shapeIds: string[]): void {
    this.config.customEvictList = shapeIds;
    this.config.evictBehavior = 'custom';
  }
  
  /**
   * Get stored dependencies for a shape ID
   */
  getDependencies(shapeId: string): Dependencies | undefined {
    return this.shapes.get(shapeId);
  }
}
