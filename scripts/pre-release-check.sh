#!/bin/bash

set -e

echo "🚀 Pre-Release Quality Checks"
echo "============================"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        return 1
    fi
}

print_warning() {
    echo -e "${YELLOW}⚠️ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 1. Code Format Check
echo ""
echo "1. Code Format Check"
echo "-------------------"
UNFORMATTED=$(gofmt -s -l .)
if [ -n "$UNFORMATTED" ]; then
    print_error "Code formatting issues found:"
    echo "$UNFORMATTED"
    echo ""
    print_info "🔧 To fix: gofmt -s -w ."
    print_error "⏸️ Please fix formatting issues before release"
    exit 1
else
    print_status 0 "Code formatting check passed"
fi

# 2. Code Quality Check
echo ""
echo "2. Code Quality Check"
echo "--------------------"
print_info "Running go vet..."
if ! go vet ./...; then
    print_error "go vet found issues"
    print_info "🔧 To fix: go vet ./..."
    print_error "⏸️ Please fix vet issues before release"
    exit 1
fi
print_status 0 "go vet check passed"

print_info "Installing golangci-lint..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

print_info "Running golangci-lint..."
if ! golangci-lint run; then
    print_error "golangci-lint found issues"
    print_info "🔧 To fix: golangci-lint run"
    print_error "⏸️ Please fix linting issues before release"
    exit 1
fi
print_status 0 "golangci-lint check passed"

# 3. Code Review Check
echo ""
echo "3. Code Review Check"
echo "-------------------"

# Check for TODO/FIXME comments
TODO_COUNT=$(grep -r "TODO\|FIXME" --include="*.go" . | wc -l)
if [ "$TODO_COUNT" -gt 0 ]; then
    print_warning "Found $TODO_COUNT TODO/FIXME comments"
    echo "Consider addressing these before release"
fi

# Check for function documentation
UNDOC_FUNCS=$(grep -r "^func [A-Z]" --include="*.go" . | grep -v "//" | wc -l)
if [ "$UNDOC_FUNCS" -gt 0 ]; then
    print_warning "Found $UNDOC_FUNCS undocumented public functions"
    echo "Consider adding documentation for public functions"
fi

# Check for hardcoded credentials
if grep -r "password\|secret\|key" --include="*.go" . | grep -v "_test.go" | grep -v "//.*password\|//.*secret\|//.*key" > /dev/null; then
    print_warning "Potential hardcoded credentials found"
    echo "Please review and ensure no sensitive data is hardcoded"
fi

# Check project structure
if [ ! -d "cmd" ] || [ ! -d "internal" ]; then
    print_error "Invalid project structure - missing cmd or internal directories"
    exit 1
fi
print_status 0 "Project structure is correct"

print_status 0 "Code review checks completed"

# 4. Test Suite
echo ""
echo "4. Test Suite"
echo "------------"
print_info "Running tests with coverage..."
if ! go test -v -race -coverprofile=coverage.out ./...; then
    print_error "Tests failed"
    print_info "🔧 To fix: go test ./..."
    print_error "⏸️ Please fix failing tests before release"
    exit 1
fi

COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
print_info "📊 Test coverage: ${COVERAGE}%"

if (( $(echo "$COVERAGE < 70" | bc -l) )); then
    print_warning "Test coverage is below 70%"
    echo "Consider adding more tests to improve coverage"
else
    print_status 0 "Test coverage is acceptable (${COVERAGE}%)"
fi

print_status 0 "All tests passed"

# 5. Build Test
echo ""
echo "5. Build Test"
echo "------------"
print_info "Testing build for current platform..."
if ! make build; then
    print_error "Build failed"
    print_error "⏸️ Please fix build issues before release"
    exit 1
fi
print_status 0 "Build test passed"

# 6. Version Check
echo ""
echo "6. Version Check"
echo "---------------"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
print_info "Current version: $VERSION"

if [ "$VERSION" == "dev" ]; then
    print_warning "No version tag found - this is a development build"
    echo "For release, create a tag: git tag v1.0.0"
else
    print_status 0 "Version tag found: $VERSION"
fi

# Summary
echo ""
echo "🎉 Pre-Release Quality Checks Complete!"
echo "====================================="
echo -e "${GREEN}✅ Code formatting: PASSED${NC}"
echo -e "${GREEN}✅ Code quality: PASSED${NC}"
echo -e "${GREEN}✅ Code review: PASSED${NC}"
echo -e "${GREEN}✅ Test suite: PASSED (${COVERAGE}% coverage)${NC}"
echo -e "${GREEN}✅ Build test: PASSED${NC}"
echo ""
echo -e "${GREEN}🚀 Ready for release!${NC}"
echo ""
echo "Next steps:"
echo "1. If version is 'dev', create a tag: git tag v1.0.0"
echo "2. Push the tag: git push origin v1.0.0"
echo "3. GitHub Actions will create the release automatically"
