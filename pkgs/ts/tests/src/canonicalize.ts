/**
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
