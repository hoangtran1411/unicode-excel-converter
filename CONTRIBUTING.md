# Contributing to VniConverter

First off, thanks for taking the time to contribute! 🎉

The following is a set of guidelines for contributing to VniConverter. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

## 🚀 Getting Started

1. **Fork the repository** on GitHub.
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/hoangtran1411/convert-vni-to-unicode.git
   cd convert-vni-to-unicode
   ```
3. **Install Dependencies**:
   - Go 1.21+
   - Node.js 20+
   - Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
4. **Run the app**:
   ```bash
   wails dev
   ```

## 🛠️ Development Workflow

1. Create a new branch for your feature or fix:
   ```bash
   git checkout -b feature/amazing-feature
   ```
2. Make your changes.
3. **Run Tests**: Ensure all tests pass before pushing.
   ```bash
   go test ./... -v
   ```
4. **Linting**: We use `golangci-lint`.
   ```bash
   golangci-lint run
   ```
5. Commit your changes with a descriptive message.

## 📝 Pull Requests

1. Push to your fork and submit a Pull Request to the `main` branch.
2. Provide a clear description of the problem and solution.
3. Link to any related Issues.
4. Ensure the **CI** checks pass (Tests + Lint).

## 🐛 Reporting Issues

If you find a bug or have a feature request, please open an Issue on GitHub.
- **Bugs**: Include steps to reproduce, expected behavior, and screenshots if possible.
- **Features**: Describe the feature and why it would be useful.

## 🎨 Coding Standards

- **Go**: Follow standard Go conventions (Effective Go). Use `go fmt`.
- **Frontend**: Clean HTML/CSS/JS. Keep the UI "Premium" (Dark theme, consistent spacing).
- **Comments**: Write comments in English explaining the *why*, not just the *what*.

## 📄 License

By contributing, you agree that your contributions will be licensed under its MIT License.
