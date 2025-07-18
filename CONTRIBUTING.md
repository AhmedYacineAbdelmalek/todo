# Contributing to Smart Todo CLI

Thank you for considering contributing to Smart Todo CLI! We welcome contributions from the community.

## 🚀 Getting Started

### Prerequisites
- Go 1.21 or later
- Git

### Setting up the Development Environment

1. **Fork the repository**
   ```bash
   # Click the "Fork" button on GitHub, then clone your fork
   git clone https://github.com/yourusername/todo.git
   cd todo
   ```

2. **Set up the Go environment**
   ```bash
   cd todo
   go mod download
   ```

3. **Build and test**
   ```bash
   go build -o todo .
   ./todo version
   ```

4. **Run tests**
   ```bash
   go test ./...
   ```

## 📝 Development Guidelines

### Code Style
- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions small and focused

### Commit Messages
Use conventional commits format:
```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```bash
feat(add): support for recurring tasks
fix(list): resolve date filtering edge case
docs(readme): update installation instructions
```

### Branch Naming
- `feature/description` - for new features
- `fix/description` - for bug fixes
- `docs/description` - for documentation updates

## 🔧 Project Structure

```
todo/
├── cmd/                    # Cobra commands
│   ├── add.go             # Task creation command
│   ├── list.go            # Task listing command
│   ├── mark.go            # Task completion command
│   ├── delete.go          # Task deletion command
│   ├── root.go            # Root command configuration
│   └── version.go         # Version command
├── taskdata/              # Data layer
│   └── task.go            # Task struct and storage operations
├── main.go                # Application entry point
├── go.mod                 # Go module definition
└── go.sum                 # Go module checksums
```

## 🧪 Testing

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Adding Tests
- Place test files alongside the code they test
- Use the `_test.go` suffix
- Follow table-driven test patterns where applicable
- Test both success and error cases

Example test structure:
```go
func TestAddTask(t *testing.T) {
    tests := []struct {
        name        string
        description string
        wantErr     bool
    }{
        {"valid task", "Complete project", false},
        {"empty task", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 🐛 Reporting Issues

### Bug Reports
Include:
- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Go version)
- Relevant command output or error messages

### Feature Requests
Include:
- Clear description of the desired feature
- Use case and motivation
- Proposed implementation approach (if any)
- Examples of similar features in other tools

## 📋 Pull Request Process

1. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```

2. **Make your changes**
   - Follow the coding standards
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**
   ```bash
   go test ./...
   go build -o todo .
   ./todo --help  # Smoke test
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add amazing new feature"
   ```

5. **Push to your fork**
   ```bash
   git push origin feature/amazing-feature
   ```

6. **Create a Pull Request**
   - Use a clear, descriptive title
   - Describe what your PR does and why
   - Reference any related issues
   - Include screenshots for UI changes

### PR Requirements
- [ ] Code follows project style guidelines
- [ ] Tests pass (`go test ./...`)
- [ ] Code builds successfully (`go build`)
- [ ] Documentation updated (if applicable)
- [ ] Self-review completed
- [ ] Ready for code review

## 🎯 Areas for Contribution

### High Priority
- **Performance optimizations** - Improve command execution speed
- **Cross-platform testing** - Ensure compatibility across OS
- **Documentation improvements** - Better examples and guides
- **Error handling** - More user-friendly error messages

### Medium Priority
- **New commands** - Additional task management features
- **Export/Import** - Support for different file formats
- **Plugins** - Extensibility framework
- **Themes** - Customizable output styling

### Low Priority
- **Integration tests** - End-to-end testing scenarios
- **Benchmarks** - Performance measurement tools
- **Localization** - Multi-language support

## 💡 Development Tips

### Local Testing
```bash
# Build and install locally
go build -o todo .
sudo cp todo /usr/local/bin/todo-dev

# Test with sample data
todo-dev add "Test task" --priority high
todo-dev list
todo-dev mark 1
todo-dev delete 1
```

### Debugging
```bash
# Enable verbose output
todo --debug command

# Use Go's built-in tools
go run -race .  # Race condition detection
go vet ./...    # Static analysis
```

### Performance Testing
```bash
# Benchmark specific functions
go test -bench=. ./taskdata

# Profile memory usage
go test -memprofile=mem.prof ./taskdata
go tool pprof mem.prof
```

## 🔒 Security

### Reporting Security Issues
Please report security vulnerabilities privately by emailing [maintainer email]. Do not create public GitHub issues for security problems.

### Security Guidelines
- Validate all user inputs
- Use secure file permissions for task storage
- Avoid storing sensitive information in tasks
- Follow Go security best practices

## 📚 Resources

### Learning Resources
- [Go Documentation](https://golang.org/doc/)
- [Cobra CLI Guide](https://github.com/spf13/cobra)
- [Effective Go](https://golang.org/doc/effective_go.html)

### Project Resources
- [GitHub Repository](https://github.com/AhmedYacineAbdelmalek/todo)
- [Issue Tracker](https://github.com/AhmedYacineAbdelmalek/todo/issues)
- [Releases](https://github.com/AhmedYacineAbdelmalek/todo/releases)

## 🙏 Recognition

Contributors will be recognized in:
- GitHub contributors page
- Release notes for significant contributions
- README acknowledgments for major features

Thank you for helping make Smart Todo CLI better! 🚀
