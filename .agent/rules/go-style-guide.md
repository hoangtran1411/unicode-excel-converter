---
trigger: always_on
---

# Go Style Guide - convert-vni-to-unicode

> **Core Rules** - For full idioms reference, see `go-idioms-reference.md`

This project is a **Desktop application to convert VNI font encodings to Unicode in Excel files** built with:
- **Wails v2** for the Desktop GUI (Windows)
- **Internal packages** for modular business logic

---

## Code Style

- Format with `gofmt`/`goimports`. Run `golangci-lint run ./...` before commit.
- Adhere to [Effective Go](https://go.dev/doc/effective_go).
- Core logic in `internal/` packages. Wails code in root package.

## Project Structure

```
convert-vni-to-unicode/
  main.go              - Wails entry point
  app.go               - Wails app bindings
  internal/
    converter/         - Domain-specific logic
  frontend/            - Wails frontend (HTML/CSS/JS)
  build/               - Build output
```

## Error Handling

- Wrap errors: `fmt.Errorf("context: %w", err)`
- Guard clauses for fail-fast
- Do not log and return the same error
- One responsibility per layer

## Wails Integration

- Bound methods on `*App` struct (PascalCase)
- `runtime.EventsEmit()` for frontend updates
- Return structs with `json` tags
- Use `runtime.*Dialog()` for file selection

## Testing & Linting

- Table-driven tests with `t.Run`
- Target 70% coverage for `internal/`
- Use `make test` and `make lint`

---

## AI Agent Rules (Critical)

### Enforcement

- Prefer clarity over cleverness
- Prefer idiomatic Go over Java/C#/JS patterns
- If unsure, follow Effective Go first

### Context Accuracy

- Documentation links ≠ guarantees of correctness
- For external APIs: prefer explicit function signatures in context
- State assumptions when context is missing

### Library Version Awareness

- Check `go.mod` for actual versions before suggesting APIs
- LLMs hallucinate APIs for newer features not in training data
- Prefer stable APIs over experimental features

### Context Engineering

- Right context at right time, not all docs at once
- Reference existing codebase patterns first
- State missing context rather than guessing

---

## Quick Reference Links

- [Effective Go](https://go.dev/doc/effective_go)
- [Wails v2](https://github.com/wailsapp/wails)
- [golangci-lint](https://github.com/golangci/golangci-lint)
- [Excelize v2](https://github.com/xuri/excelize)

> **Full Reference:** See `.agent/rules/go-idioms-reference.md` for detailed idioms, code examples, and best practices.
