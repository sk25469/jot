# Testing Documentation

This document outlines the testing strategy and test coverage for the jot CLI application.

## Test Structure

The project includes comprehensive test coverage across multiple layers:

### 📦 **Package Tests**

#### `models/` - **100% Coverage**
- Tests all data structures and their methods
- Validates default values and field assignments
- Tests model relationships and embedded types

#### `styles/` - **100% Coverage**
- Tests lipgloss styling components
- Validates color definitions and style applications
- Tests rendering functions (headers, tags, progress bars, etc.)
- Verifies deterministic tag coloring

#### `service/` - **14.4% Coverage**
- Tests core business logic functions
- Validates note content generation and parsing
- Tests hash generation and content preview creation
- Tests word counting and metadata handling

#### `config/` - **23.5% Coverage**
- Tests configuration loading and path handling
- Validates environment variable processing
- Tests default value generation

#### `integration/` - **Integration Tests**
- Tests cross-package interactions
- Validates environment setup requirements
- Tests time handling and error scenarios

## Running Tests

### Basic Test Commands

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Generate detailed HTML coverage report
make test-coverage-detailed
```

### Individual Package Testing

```bash
# Test specific packages
go test ./models -v
go test ./styles -v
go test ./service -v
go test ./config -v
```

## Test Categories

### 🧪 **Unit Tests**
- **Models**: Data structure validation
- **Styles**: UI component rendering
- **Service**: Business logic functions
- **Config**: Configuration handling

### 🔗 **Integration Tests**
- Cross-package functionality
- Environment requirement validation
- Error handling scenarios
- Time and file system operations

### 📊 **Coverage Goals**

| Package | Current Coverage | Target |
|---------|------------------|--------|
| `models` | 100% | ✅ 100% |
| `styles` | 100% | ✅ 100% |
| `service` | 14.4% | 🎯 70%+ |
| `config` | 23.5% | 🎯 70%+ |

## Test Examples

### Model Testing
```go
func TestDefaultListFilter(t *testing.T) {
    filter := DefaultListFilter()
    if filter.Limit != 100 {
        t.Errorf("Expected Limit to be 100, got %d", filter.Limit)
    }
}
```

### Style Testing
```go
func TestGetTagStyle(t *testing.T) {
    style1 := GetTagStyle("golang")
    style2 := GetTagStyle("golang")
    
    rendered1 := style1.Render("golang")
    rendered2 := style2.Render("golang")
    
    if rendered1 != rendered2 {
        t.Errorf("Same tag should produce same styled output")
    }
}
```

### Service Testing
```go
func TestGenerateContentHash(t *testing.T) {
    service := &NoteService{}
    content := "test content"
    
    hash1 := service.generateContentHash(content)
    hash2 := service.generateContentHash(content)
    
    if hash1 != hash2 {
        t.Errorf("Hash should be deterministic")
    }
}
```

## Continuous Integration

The test suite is designed to be run in CI/CD pipelines:

```bash
# Quality check command (includes tests)
make check
```

This runs:
- All unit and integration tests
- Code coverage analysis
- Go vet for static analysis
- Go fmt for code formatting

## Testing Best Practices

### ✅ **What We Test**
- Data structure integrity
- Business logic correctness
- UI component rendering
- Configuration handling
- Error scenarios
- Edge cases and boundary conditions

### 🔍 **Test Patterns**
- **Table-driven tests** for multiple scenarios
- **Deterministic testing** for hash functions and styling
- **Environment isolation** for config tests
- **Mock-friendly design** for future database testing

### 📈 **Future Improvements**
- Database layer testing with test database
- Command-line interface testing
- Performance benchmarks
- End-to-end CLI testing

## Test Maintenance

### Adding New Tests
1. Create `*_test.go` files alongside source code
2. Follow naming convention: `TestFunctionName`
3. Use table-driven tests for multiple scenarios
4. Include edge cases and error conditions

### Running Specific Tests
```bash
# Run specific test function
go test ./models -run TestDefaultListFilter -v

# Run tests matching pattern
go test ./... -run TestRender -v
```

### Debugging Tests
```bash
# Run with race detection
go test ./... -race

# Run with memory profiling
go test ./... -memprofile=mem.prof

# Run with CPU profiling
go test ./... -cpuprofile=cpu.prof
```

---

The comprehensive test suite ensures code quality, prevents regressions, and enables confident refactoring and feature development.