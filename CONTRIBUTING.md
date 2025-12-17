# Contributing to EcomGo

We welcome contributions! Please follow these guidelines to ensure a smooth contribution process.

## Code of Conduct

Be respectful and professional in all interactions.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-org/ecomgo.git`
3. Create a branch: `git checkout -b feature/your-feature-name`
4. Set up development environment: see [SETUP.md](./SETUP.md)

## Development Process

### 1. Before You Start

- Check existing issues and pull requests
- Create an issue describing your feature or fix
- Wait for approval from maintainers

### 2. Code Standards

- Follow Go style guide: https://golang.org/doc/effective_go
- Run linter: `golangci-lint run`
- Format code: `go fmt ./...`
- Add comments for exported functions

### 3. Testing

- Write tests for new functionality
- Run tests: `go test ./...`
- Maintain >80% code coverage
- Test with both MySQL and PostgreSQL

### 4. Security

- Never commit sensitive information (.env files)
- Don't hardcode credentials or secrets
- Use environment variables for configuration
- Review [Security Considerations](./README.md#security-considerations)

### 5. Documentation

- Update README.md if adding features
- Update API_DOCUMENTATION.md for new endpoints
- Add inline code comments for complex logic
- Update ARCHITECTURE.md if changing system design

### 6. Commit Messages

Use clear, descriptive commit messages:

```
feat: add user password reset endpoint
fix: resolve database connection timeout issue
docs: update API documentation
refactor: simplify user service logic
test: add tests for registration handler
```

## Pull Request Process

1. Update documentation (README.md, API docs, etc.)
2. Add or update tests
3. Ensure all tests pass: `go test ./...`
4. Run linter: `golangci-lint run`
5. Submit PR with clear description
6. Link related issues
7. Request review from maintainers

## PR Template

```
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Related Issue
Fixes #(issue number)

## How Has This Been Tested?
Description of test process

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] No hardcoded secrets
```

## Feature Requests

Create an issue with:
- Clear description of requested feature
- Use case and motivation
- Potential implementation approach
- Any related issues or PRs

## Bug Reports

Create an issue with:
- Clear, descriptive title
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment details (OS, Go version, database)
- Logs or error messages

## Areas for Contribution

### High Priority

- Implement bcrypt for password hashing (SECURITY)
- Add JWT token validation middleware
- Add rate limiting for auth endpoints
- Write comprehensive test suite

### Medium Priority

- Implement Keycloak integration
- Add request/response validation middleware
- Add pagination to endpoints
- Add email verification on registration

### Low Priority

- Add more detailed code comments
- Improve error messages
- Optimize database queries
- Add monitoring/observability

## Code Review

All code is reviewed by maintainers before merging. Reviews focus on:

- Code quality and style
- Test coverage
- Security implications
- Performance impact
- Documentation completeness

## License

By contributing to EcomGo, you agree that your contributions will be licensed under the Apache License 2.0. See the [LICENSE](./LICENSE) file for details.

## Questions?

Open an issue with the "question" label or contact maintainers.

Thank you for contributing!
