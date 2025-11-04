# @includekit/types

**Production types-only package** for IncludeKit Universal Specification.

This package provides **TypeScript type definitions only** with no runtime code.

## Installation

```bash
npm install @includekit/types
```

## Usage

```typescript
import type { QueryShape, MutationEvent, Dependencies } from '@includekit/types';

const shape: QueryShape = {
  model: 'Post',
  where: { atoms: [{ field: 'published', op: 'eq', value: true }] },
  orderBy: [{ field: 'createdAt', direction: 'desc' }],
  take: 10,
};
```

## For Testing & Validation

If you need runtime validators, JSON canonicalization, or shapeId computation, use:

```bash
npm install --save-dev @includekit/types-testkit
```

See [includekit-types-testkit](https://github.com/bold-minds/includekit-types-testkit) for details.

## License

Apache-2.0
