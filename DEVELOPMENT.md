# Development Guide

This document provides information for developers working on the td project.

## Architecture

td follows a clean architecture with clear separation of concerns:

```
cmd/td/           # Main application entry point
internal/
├── config/       # Configuration management
├── storage/      # Data persistence layer
├── task/         # Business logic and domain models
└── ui/           # Terminal user interface
integration/      # Integration tests
```

### Package Overview

- **config**: Handles application configuration loading and defaults
- **storage**: Manages data persistence with validation and integrity checks
- **task**: Contains the core business logic, task management, and undo/redo operations
- **ui**: Implements the terminal interface using Bubble Tea

## Building

```bash
# Build the application
make build

# Or manually
go build ./cmd/td

# Run tests
go test ./...

# Run integration tests
go test ./integration
```

## Code Quality

### Testing

- Unit tests are located alongside the code they test (`*_test.go`)
- Integration tests are in the `integration/` directory
- All tests can be run with `go test ./...`

### Linting

The codebase uses Go's standard formatting and vetting tools:

```bash
go fmt ./...
go vet ./...
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Add documentation comments for exported functions/types
- Keep functions focused on single responsibilities

## Adding Features

### 1. New Commands/Features

1. Add business logic to the appropriate package (likely `task/`)
2. Add UI handling in `ui/update.go`
3. Add key binding in `ui/keymap.go` and config
4. Update help text in `ui/view.go`
5. Add tests for new functionality

### 2. Configuration Options

1. Add fields to `config.Config` struct
2. Update `DefaultConfig()` function
3. Update config loading/merging logic
4. Use the config in appropriate places

### 3. Data Persistence Changes

1. Update `storage.Repository` methods
2. Maintain backward compatibility with existing data
3. Add validation for new data formats
4. Update tests

## Key Design Decisions

### Undo/Redo System

The undo/redo system uses the Command pattern with action types and state snapshots. This allows for:

- Complex operations to be undone atomically
- State to be restored precisely
- Extensibility for new operation types

### Task Management

Tasks are managed through a `TaskManager` that provides:

- Immutable operations (returns new task references)
- Automatic sorting by priority and creation time
- Centralized business logic validation

### Data Validation

Multiple layers of validation ensure data integrity:

1. **Input validation** in UI layer (prevents invalid user input)
2. **Business logic validation** in task layer (ensures consistency)
3. **Data validation** in storage layer (prevents corruption)

## Performance Considerations

- Tasks are cached in the UI layer to avoid repeated filtering
- JSON serialization is used for simplicity and human readability
- File I/O is minimized through batching operations

## Future Enhancements

Areas identified for future improvement:

- Plugin system for custom storage backends
- Advanced filtering and search capabilities
- Cloud synchronization
- Mobile/web companion apps
- Advanced theming and customization
