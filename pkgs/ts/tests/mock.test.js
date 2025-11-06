import { test } from 'node:test';
import { strict as assert } from 'node:assert';
import { MockIncludeKitEngine } from './dist/mock/index.js';

test('MockIncludeKitEngine: setSchema stores schema', () => {
  const engine = new MockIncludeKitEngine({ trackCalls: true });
  const schema = {
    version: 1,
    models: [
      { name: 'users', id: { kind: 'string' } },
      { name: 'posts', id: { kind: 'string' } }
    ]
  };
  
  engine.setSchema(schema);
  
  const calls = engine.getCalls();
  assert.equal(calls.setSchema.length, 1);
  assert.deepEqual(calls.setSchema[0].schema, schema);
});

test('MockIncludeKitEngine: computeShapeId returns deterministic shape_id', () => {
  const engine = new MockIncludeKitEngine();
  const statement = {
    query: {
      model: 'users',
      where: {
        conditions: [{ field: 'id', op: 'eq', value: '1' }]
      }
    }
  };
  
  const result1 = engine.computeShapeId(statement);
  const result2 = engine.computeShapeId(statement);
  
  assert.equal(result1.shape_id, result2.shape_id);
  assert.ok(result1.shape_id.startsWith('s_'));
  assert.equal(result1.shape_id.length, 66); // s_ + 64 hex chars
});

test('MockIncludeKitEngine: computeShapeId uses custom generator', () => {
  const engine = new MockIncludeKitEngine({
    shapeIdGenerator: () => 's_custom_test'
  });
  
  const result = engine.computeShapeId({ query: { model: 'users' } });
  assert.equal(result.shape_id, 's_custom_test');
});

test('MockIncludeKitEngine: addQuery returns shape_id and dependencies', () => {
  const engine = new MockIncludeKitEngine();
  const statement = {
    query: {
      model: 'users',
      where: {
        conditions: [{ field: 'id', op: 'eq', value: '1' }]
      }
    }
  };
  
  const result = engine.addQuery({ shape: statement });
  
  assert.ok(result.shape_id);
  assert.equal(result.dependencies.shape_id, result.shape_id);
  assert.deepEqual(result.dependencies.records, {});
  assert.equal(result.dependencies.filters.length, 1);
});

test('MockIncludeKitEngine: addQuery extracts records from result_hint', () => {
  const engine = new MockIncludeKitEngine();
  const statement = {
    query: { model: 'users' }
  };
  
  const result = engine.addQuery({
    shape: statement,
    result_hint: {
      users: [
        { id: '1', name: 'Alice' },
        { id: '2', name: 'Bob' }
      ]
    }
  });
  
  assert.deepEqual(result.dependencies.records.users, ['1', '2']);
});

test('MockIncludeKitEngine: invalidate evicts shapes for affected models', () => {
  const engine = new MockIncludeKitEngine();
  const statement = {
    query: { model: 'users' }
  };
  
  const { shape_id } = engine.addQuery({
    shape: statement,
    result_hint: {
      users: [{ id: '1', name: 'Alice' }]
    }
  });
  
  const mutation = {
    changes: [{
      model: 'users',
      action: 'update',
      set: [{ field: 'name', value: 'Alice Updated' }],
      where: { conditions: [{ field: 'id', op: 'eq', value: '1' }] }
    }]
  };
  
  const result = engine.invalidate(mutation);
  assert.ok(result.evict.includes(shape_id));
});

test('MockIncludeKitEngine: invalidate does not evict unrelated models', () => {
  const engine = new MockIncludeKitEngine();
  const statement = {
    query: { model: 'posts' }
  };
  
  const { shape_id } = engine.addQuery({ shape: statement });
  
  const mutation = {
    changes: [{
      model: 'users',
      action: 'update',
      set: [{ field: 'name', value: 'Test' }]
    }]
  };
  
  const result = engine.invalidate(mutation);
  assert.ok(!result.evict.includes(shape_id));
});

test('MockIncludeKitEngine: invalidate uses custom evict list', () => {
  const engine = new MockIncludeKitEngine();
  engine.setEvictList(['s_custom_1', 's_custom_2']);
  
  const mutation = {
    changes: [{ model: 'users', action: 'update', set: [] }]
  };
  
  const result = engine.invalidate(mutation);
  assert.deepEqual(result.evict, ['s_custom_1', 's_custom_2']);
});

test('MockIncludeKitEngine: explainInvalidation explains record membership', () => {
  const engine = new MockIncludeKitEngine();
  const statement = {
    query: { model: 'users' }
  };
  
  const { shape_id } = engine.addQuery({
    shape: statement,
    result_hint: { users: [{ id: '1' }] }
  });
  
  const mutation = {
    changes: [{ model: 'users', action: 'update', set: [] }]
  };
  
  const result = engine.explainInvalidation({ mutation, shape_id });
  
  assert.equal(result.invalidate, true);
  assert.ok(result.reasons.includes('record_membership'));
});

test('MockIncludeKitEngine: explainInvalidation returns false for unknown shape_id', () => {
  const engine = new MockIncludeKitEngine();
  
  const mutation = {
    changes: [{ model: 'users', action: 'update', set: [] }]
  };
  
  const result = engine.explainInvalidation({
    mutation,
    shape_id: 's_unknown'
  });
  
  assert.equal(result.invalidate, false);
  assert.deepEqual(result.reasons, []);
});

test('MockIncludeKitEngine: reset clears all state', () => {
  const engine = new MockIncludeKitEngine({ trackCalls: true });
  
  engine.addQuery({ shape: { query: { model: 'users' } } });
  assert.equal(engine.getCalls().addQuery.length, 1);
  
  engine.reset();
  
  const calls = engine.getCalls();
  assert.equal(calls.addQuery.length, 0);
});

test('MockIncludeKitEngine: getVersion returns version info', () => {
  const engine = new MockIncludeKitEngine();
  
  const version = engine.getVersion();
  
  assert.equal(version.core, 'mock-0.1.0');
  assert.equal(version.contract, '0.1.0');
  assert.equal(version.abi, '1');
});

test('MockIncludeKitEngine: trackCalls records all method calls', () => {
  const engine = new MockIncludeKitEngine({ trackCalls: true });
  
  engine.setSchema({ version: 1, models: [] });
  engine.computeShapeId({ query: { model: 'users' } });
  engine.addQuery({ shape: { query: { model: 'posts' } } });
  engine.getVersion();
  
  const calls = engine.getCalls();
  assert.equal(calls.setSchema.length, 1);
  assert.equal(calls.computeShapeId.length, 2); // Called by addQuery too
  assert.equal(calls.addQuery.length, 1);
  assert.equal(calls.getVersion.length, 1);
});
