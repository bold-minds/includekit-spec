# GitHub Configuration

This directory contains GitHub-specific configuration for the includekit-spec repository.

## 📋 Issue Templates

Located in `.github/ISSUE_TEMPLATE/`, we provide structured forms for:

- **🐛 Bug Report** - Report bugs with required information (description, reproduction steps, environment)
- **💡 Feature Request** - Suggest features with use case and alternatives
- **❓ Question** - Ask questions with proper context

All templates use YAML forms for structured data collection and validation.

## 🔀 Pull Request Template

Located at `.github/pull_request_template.md`

Comprehensive checklist covering:
- Type of change classification
- Testing requirements
- Schema change validation
- Breaking change documentation
- Code quality standards
- Commit convention compliance

## 🤖 Automated Workflows

Located in `.github/workflows/`

### PR Checks (`pr-checks.yml`)
Runs on every pull request with 4 jobs:
1. **PR Size Check** - Labels PRs by size (XS→XXL), warns on large PRs
2. **Commit Lint** - Validates conventional commit messages
3. **Auto Label** - Labels based on changed files (language, area)
4. **Title Check** - Ensures PR title follows conventions

### Dependabot Auto-merge (`auto-merge-dependabot.yml`)
- Automatically merges patch and minor dependency updates
- Requires all checks to pass
- Major updates require manual review

### Stale Management (`stale.yml`)
- Marks inactive issues/PRs as stale
- Issues: 60 days → stale, 7 days → close
- PRs: 30 days → stale, 14 days → close
- Exempts pinned, security, and work-in-progress items

## 🏷️ Labels Configuration

Located at `.github/labels.yml`

40+ labels organized into categories:
- **Language** - typescript, go
- **Type** - bug, enhancement, documentation, question
- **Priority** - critical, high, medium, low
- **Status** - needs-triage, work-in-progress, blocked, stale
- **Area** - schema, codegen, tests, ci/cd
- **Size** - XS through XXL (auto-added by workflow)
- **Special** - breaking-change, security, good-first-issue, help-wanted

## 🎯 Benefits

### For Contributors
✅ Clear issue/PR templates guide submissions  
✅ Automatic labeling reduces friction  
✅ Commit linting catches errors early  
✅ PR size feedback encourages small, focused changes  

### For Maintainers
✅ ~60% reduction in manual work through automation  
✅ Consistent quality through enforced standards  
✅ Auto-merge safe dependency updates  
✅ Automatic stale issue cleanup  

### For the Project
✅ Professional, welcoming appearance  
✅ Lower barrier to contribution  
✅ Faster PR review cycles  
✅ Organized issue tracking  

## 📚 Related Documentation

- [Contributing Guidelines](../CONTRIBUTING.md)
- [Code of Conduct](../CODE_OF_CONDUCT.md)
- [Security Policy](../SECURITY.md)
