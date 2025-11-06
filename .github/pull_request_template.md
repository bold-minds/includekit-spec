## Description

<!-- Provide a clear and concise description of your changes -->

## Type of Change

<!-- Mark the relevant option with an "x" -->

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] 📝 Documentation update
- [ ] ♻️ Code refactoring (no functional changes)
- [ ] ⚡ Performance improvement
- [ ] 🧪 Test improvements
- [ ] 🔧 Build/CI changes

## Related Issues

<!-- Link to related issues using #issue_number -->

Closes #
Related to #

## Changes Made

<!-- List the specific changes you made -->

- 
- 
- 

## Testing

<!-- Describe how you tested your changes -->

### Test Commands Run

```bash
# Commands you ran to test
./scripts/test.sh
```

### Test Results

- [ ] All TypeScript tests pass
- [ ] All Go tests pass
- [ ] Conformance tests pass
- [ ] No-runtime constraint verified
- [ ] Manual testing completed (if applicable)

## Schema Changes

<!-- If you modified the schema, answer these questions -->

- [ ] No schema changes
- [ ] Schema changes are backward compatible
- [ ] Schema changes require version bump

**If schema was changed:**
- [ ] Version number updated in `VERSION` file
- [ ] Version sync script run (`go run tools/version/sync.go`)
- [ ] Code regenerated (`cd codegen && go run .`)
- [ ] CHANGELOG.md updated

## Breaking Changes

<!-- If this is a breaking change, describe the impact -->

**Does this PR introduce breaking changes?**
- [ ] No
- [ ] Yes (describe below)

**If yes, describe the breaking changes and migration path:**

<!-- 
- What breaks?
- How should users migrate?
- Why is this necessary?
-->

## Documentation

- [ ] README.md updated (if needed)
- [ ] schema/README.md updated (if needed)
- [ ] CHANGELOG.md updated
- [ ] Code comments added/updated
- [ ] TypeDoc/godoc comments added (if new public APIs)

## Code Quality

- [ ] Code follows the project's style guidelines
- [ ] Self-review of code performed
- [ ] Comments added in hard-to-understand areas
- [ ] No console.log or debug code left in
- [ ] No TODO comments without GitHub issue links

## Commit Convention

- [ ] Commits follow [Conventional Commits](https://www.conventionalcommits.org/) format
- [ ] Commit messages are clear and descriptive

**Example:** 
- `feat: add support for custom operators`
- `fix: correct schema field name validation`
- `docs: update TypeScript usage examples`

## Additional Notes

<!-- Any additional information for reviewers -->

## Checklist

- [ ] I have read the [CONTRIBUTING.md](https://github.com/bold-minds/includekit-spec/blob/main/CONTRIBUTING.md)
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings or errors
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published

---

<!-- Thank you for contributing to IncludeKit Spec! 🎉 -->
