#!/bin/bash
# KDSE Workspace Acceptance Tests
# Tests the .kdse/ workspace architecture enforcement

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
    ((TESTS_PASSED++))
}

fail() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    ((TESTS_FAILED++))
}

info() {
    echo -e "${YELLOW}ℹ INFO${NC}: $1"
}

# Setup test repository
setup_test_repo() {
    info "Creating test repository..."
    TEST_DIR=$(mktemp -d)
    cd "$TEST_DIR"
    git init --quiet
    echo "# Test Repository" > README.md
    git add README.md
    git commit --quiet -m "Initial commit"
    echo "Test repository created at: $TEST_DIR"
}

cleanup_test_repo() {
    if [ -n "$TEST_DIR" ] && [ -d "$TEST_DIR" ]; then
        rm -rf "$TEST_DIR"
        info "Cleaned up test repository"
    fi
}

# Trap to cleanup on exit
trap cleanup_test_repo EXIT

echo "================================================"
echo "KDSE Workspace Architecture Acceptance Tests"
echo "================================================"
echo ""

# Build the MCP server
info "Building MCP server..."
cd /workspace/project/KDSE/mcp
go build -o kdse-mcp-test . 2>/dev/null || {
    fail "Failed to build MCP server"
    exit 1
}
MCP_SERVER="/workspace/project/KDSE/mcp/kdse-mcp-test"
pass "MCP server built successfully"

# Test 1: Initialize creates .kdse/ directory
echo ""
info "Test 1: Initialize creates .kdse/ directory"
setup_test_repo

# Send initialize request
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"initialize"},"id":1}' | $MCP_SERVER > /tmp/initialize_output.json

# Check if .kdse/ was created
if [ -d ".kdse" ]; then
    pass ".kdse/ directory was created"
else
    fail ".kdse/ directory was NOT created"
fi

# Check if initialize response contains workspace info
if grep -q '"workspace"' /tmp/initialize_output.json; then
    pass "Initialize response contains workspace info"
else
    fail "Initialize response missing workspace info"
fi

# Test 2: No methodology directories in repo root
echo ""
info "Test 2: No methodology directories in repo root"

# Create a legacy directory
mkdir -p foundation
echo "# Foundation" > foundation/test.md

# Run status and check for warning
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"status"},"id":2}' | $MCP_SERVER > /tmp/status_output.json

if grep -q "legacy_dirs_warning" /tmp/status_output.json; then
    pass "Legacy directory detected and reported"
else
    fail "Legacy directory NOT detected"
fi

# Test 3: Migrate moves legacy directories to .kdse/
echo ""
info "Test 3: Migrate moves legacy directories to .kdse/"

# Run migrate
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"migrate"},"id":3}' | $MCP_SERVER > /tmp/migrate_output.json

# Check if foundation was migrated
if [ -d ".kdse/foundation" ]; then
    pass "Foundation migrated to .kdse/foundation/"
else
    fail "Foundation NOT migrated"
fi

# Check if original foundation is gone
if [ -d "foundation" ]; then
    fail "Original foundation/ still exists"
else
    pass "Original foundation/ removed"
fi

# Test 4: Collect creates artifacts under .kdse/
echo ""
info "Test 4: Collect creates artifacts under .kdse/"

# Run collect
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"collect"},"id":4}' | $MCP_SERVER > /tmp/collect_output.json

# Check response
if grep -q '.kdse/artifacts' /tmp/collect_output.json; then
    pass "Collect response references .kdse/artifacts/"
else
    fail "Collect response missing .kdse/artifacts/ reference"
fi

# Test 5: Foundation tool creates foundation under .kdse/
echo ""
info "Test 5: Foundation tool creates foundation under .kdse/"

# Run foundation
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"foundation"},"id":5}' | $MCP_SERVER > /tmp/foundation_output.json

# Check if .kdse/foundation exists
if [ -d ".kdse/foundation" ]; then
    pass ".kdse/foundation/ created/verified"
else
    fail ".kdse/foundation/ NOT created"
fi

# Test 6: Audit creates reports under .kdse/
echo ""
info "Test 6: Audit creates reports under .kdse/"

# Run audit
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"audit"},"id":6}' | $MCP_SERVER > /tmp/audit_output.json

# Check response
if grep -q '.kdse/reports' /tmp/audit_output.json; then
    pass "Audit response references .kdse/reports/"
else
    fail "Audit response missing .kdse/reports/ reference"
fi

# Test 7: Status shows correct workspace root
echo ""
info "Test 7: Status shows correct workspace root"

# Check if status output contains workspace_root
if grep -q "workspace_root" /tmp/status_output.json && grep -q '.kdse' /tmp/status_output.json; then
    pass "Status shows workspace_root pointing to .kdse/"
else
    fail "Status missing workspace_root info"
fi

# Cleanup test repo for final test
cleanup_test_repo

# Test 8: Help includes workspace architecture info
echo ""
info "Test 8: Help includes workspace architecture info"

cd /workspace/project/KDSE
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"help"},"id":7}' | $MCP_SERVER > /tmp/help_output.json

if grep -q "workspace" /tmp/help_output.json; then
    pass "Help includes workspace information"
else
    fail "Help missing workspace information"
fi

# Test 9: Tools list includes new tools
echo ""
info "Test 9: Tools list includes new tools"

echo '{"jsonrpc":"2.0","method":"tools/list","params":{},"id":8}' | $MCP_SERVER > /tmp/tools_output.json

for tool in collect foundation audit migrate; do
    if grep -q "\"$tool\"" /tmp/tools_output.json; then
        pass "Tool list includes $tool"
    else
        fail "Tool list missing $tool"
    fi
done

# Summary
echo ""
echo "================================================"
echo "Test Summary"
echo "================================================"
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Failed: $TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed.${NC}"
    exit 1
fi
