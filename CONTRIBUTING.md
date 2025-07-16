# Contributing to LazyTodo

Thank you for considering contributing to LazyTodo! We appreciate your help in making this terminal todo manager even better.

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- SQLite (for development)

### Setting up the Development Environment

1. **Fork and Clone**
   ```bash
   git clone https://github.com/yourusername/lazytodo.git
   cd lazytodo
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Build and Run**
   ```bash
   go build -o lazytodo cmd/main.go
   ./lazytodo
   ```

4. **Run Tests**
   ```bash
   go test -v ./...
   ```

## 📝 Development Guidelines

### Code Style

- Follow Go conventions and best practices
- Use `gofmt` to format your code
- Run `go vet` to check for issues
- Ensure all tests pass before submitting

### Project Structure

```
lazytodo/
├── cmd/                    # Application entry points
│   └── main.go
├── internal/              # Private application code
│   ├── ui/               # User interface components
│   ├── storage/          # Data storage layer
│   └── models/           # Data models
├── migrations/           # Database migrations
├── scripts/             # Installation scripts
├── .github/             # GitHub Actions workflows
└── docs/               # Documentation
```

### Commit Messages

We use conventional commits. Please format your commit messages as:

```
type(scope): description

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(ui): add dark mode support
fix(storage): resolve database connection issue
docs(readme): update installation instructions
```

## 🐛 Reporting Issues

### Bug Reports

Please use the GitHub issue template and include:

- **Description**: Clear description of the issue
- **Steps to Reproduce**: Detailed steps to reproduce the bug
- **Expected Behavior**: What you expected to happen
- **Actual Behavior**: What actually happened
- **Environment**: OS, terminal, Go version
- **Screenshots**: If applicable

### Feature Requests

For feature requests, please include:

- **Use Case**: Why is this feature needed?
- **Description**: Detailed description of the feature
- **Mockups**: UI mockups if applicable
- **Alternatives**: Any alternative solutions considered

## 🔧 Making Changes

### Development Workflow

1. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

2. **Make Changes**
   - Write code following our guidelines
   - Add tests for new functionality
   - Update documentation if needed

3. **Test Your Changes**
   ```bash
   # Run tests
   go test -v ./...
   
   # Test builds
   go build -o lazytodo cmd/main.go
   
   # Test installation scripts (if modified)
   bash -n scripts/install-linux.sh
   bash -n scripts/install-mac.sh
   ```

4. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat(ui): add your feature description"
   ```

5. **Push and Create PR**
   ```bash
   git push origin feature/your-feature-name
   ```

### Pull Request Guidelines

- **Title**: Clear, descriptive title
- **Description**: Detailed description of changes
- **Testing**: Describe how you tested the changes
- **Screenshots**: Include screenshots for UI changes
- **Breaking Changes**: Clearly mark any breaking changes

### Code Review Process

1. All PRs require at least one review
2. CI checks must pass
3. Address reviewer feedback promptly
4. Squash commits before merge if requested

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./internal/ui/
```

### Writing Tests

- Write unit tests for new functions
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for good test coverage (>80%)

### Test Categories

1. **Unit Tests**: Test individual functions/methods
2. **Integration Tests**: Test component interactions
3. **UI Tests**: Test user interface components
4. **End-to-End Tests**: Test complete workflows

## 📦 Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Creating Releases

1. **Update Version**
   - Update version in code if applicable
   - Update CHANGELOG.md

2. **Create Tag**
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

3. **GitHub Actions** will automatically:
   - Build binaries for all platforms
   - Create GitHub release
   - Upload release assets
   - Generate checksums

## 🎨 UI/UX Guidelines

### Design Principles

- **Simplicity**: Keep the interface clean and intuitive
- **Consistency**: Use consistent styling throughout
- **Accessibility**: Support different terminal capabilities
- **Performance**: Maintain smooth, responsive interactions

### Color Scheme

- **Primary**: Purple (#7C3AED)
- **Accent**: Green (#10B981)
- **Text**: Follow terminal color capabilities
- **Borders**: Elegant Unicode characters

### Key Bindings

- Follow vim-like conventions where possible
- Provide clear help documentation
- Support common shortcuts (Ctrl+C, Ctrl+Q)

## 🤝 Community

### Communication

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and ideas
- **Pull Requests**: Code contributions

### Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help newcomers feel welcome
- Follow the GitHub Community Guidelines

## 📚 Resources

### Documentation

- [Go Documentation](https://golang.org/doc/)
- [Bubble Tea Framework](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Styling](https://github.com/charmbracelet/lipgloss)
- [SQLite Documentation](https://sqlite.org/docs.html)

### Tools

- [Go Tools](https://golang.org/cmd/)
- [golangci-lint](https://golangci-lint.run/)
- [Delve Debugger](https://github.com/go-delve/delve)

## ❓ Getting Help

If you need help with contributing:

1. Check existing [GitHub Issues](https://github.com/yourusername/lazytodo/issues)
2. Create a new issue with the `question` label
3. Ask in [GitHub Discussions](https://github.com/yourusername/lazytodo/discussions)

Thank you for contributing to LazyTodo! 🎉 