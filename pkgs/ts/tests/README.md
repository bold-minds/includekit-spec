# @includekit/spec-testkit

**Testkit package** for IncludeKit Universal Format - runtime validators, JSON Canonicalization, and shapeId computation.

⚠️ **This package is for testing and development only.** Do not import in production code.

## Installation

```bash
npm install --save-dev @includekit/spec-testkit
```

## Usage

```typescript
import {
  validateQueryShape,
  validateMutationEvent,
  validateDependencies,
  canonicalize,
  computeShapeId,
  computeQueryShapeId,
} from '@includekit/spec-testkit';

const shape = {
  model: 'Post',
  where: { atoms: [{ field: 'published', op: 'eq', value: true }] },
  orderBy: [{ field: 'createdAt', direction: 'desc' }],
  take: 10,
};

// Validate
validateQueryShape(shape); // throws ValidationError if invalid

// Canonicalize and compute shapeId
const canonical = canonicalize(shape);
const shapeId = computeShapeId(canonical);
// or directly:
const id = computeQueryShapeId(shape);

console.log(id); // s_<64-char-hex>
```

## API

### Validators

- `validateQueryShape(shape: any): asserts shape is QueryShape`
- `validateMutationEvent(event: any): asserts event is MutationEvent`
- `validateDependencies(deps: any): asserts deps is Dependencies`

### Canonicalization

- `canonicalize(obj: any): string` - JCS canonicalization
- `canonicalizeQueryShape(shape: any): string` - Removes diagnostic fields first

### ShapeId

- `computeShapeId(canonicalJson: string): string` - Compute from canonical JSON
- `computeQueryShapeId(shape: any): string` - Convenience for QueryShape

## License

Apache-2.0
