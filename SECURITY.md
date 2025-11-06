# Security Policy

## Supported Versions

We release security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability, please follow these steps:

### Private Disclosure Process

1. **Do NOT** open a public GitHub issue for security vulnerabilities
2. **Email** the maintainers directly at support@boldminds.tech
3. **Include** as much information as possible:
   - Type of vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if available)

### What to Expect

- **Initial Response**: Within 48 hours acknowledging receipt
- **Status Updates**: Every 7 days until resolution
- **Fix Timeline**: We aim to release fixes within 30 days for critical issues
- **Disclosure**: We'll coordinate public disclosure with you

## Security Best Practices

### For Users

- **Keep dependencies updated**: Run `npm update` and `go get -u` regularly
- **Verify integrity**: Check package checksums when installing
- **Use official sources**: Only install from npm and official Go module proxy

### For Contributors

- **No secrets in code**: Never commit API keys, passwords, or credentials
- **Validate input**: All public APIs should validate and sanitize inputs
- **Follow least privilege**: Grant minimal necessary permissions
- **Keep dependencies minimal**: Fewer dependencies = smaller attack surface

## Known Security Considerations

### Canonicalization

- The canonicalization implementation is simplified (key ordering only)
- Full RFC 8785 compliance is not required for shapeId determinism
- Prototype pollution protection is implemented

### Type Safety

- Production packages are types-only (no runtime code)
- Testkit packages include validation but are for dev/test only
- Never import testkit packages in production code

## Security Features

### Built-in Protections

- ✅ Prototype pollution prevention in canonicalization
- ✅ Input validation in all testkit validators
- ✅ Type-safe schema definitions
- ✅ No eval() or dynamic code execution
- ✅ Minimal dependencies

### Regular Audits

We perform:
- Manual code reviews for all PRs
- Automated dependency scanning (coming soon)
- Regular security audits of critical paths

## Acknowledgments

We appreciate responsible disclosure. Security researchers who report valid vulnerabilities will be:

- Credited in CHANGELOG (unless you prefer to remain anonymous)
- Acknowledged in security advisories
- Thanked publicly (with your permission)

Thank you for helping keep IncludeKit secure! 🔒
