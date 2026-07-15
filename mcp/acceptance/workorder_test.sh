#!/bin/bash
# KDSE WorkOrder Acceptance Tests
# Tests that the Runtime Owns Methodology - WorkOrders are explicit

# Don't use set -e - we want to continue on individual test failures

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0
TEST_DIR=""

# Helper functions
pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
    ((TESTS_PASSED++)) || true
}

fail() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    ((TESTS_FAILED++)) || true
}

info() {
    echo -e "${YELLOW}ℹ INFO${NC}: $1"
}

# Setup test repository
setup_test_repo() {
    info "Creating test repository..."
    TEST_DIR=$(mktemp -d)
    cd "$TEST_DIR" || return 1
    git init --quiet 2>/dev/null || true
    echo "# Test Repository" > README.md
    git add README.md 2>/dev/null || true
    git commit --quiet -m "Initial commit" 2>/dev/null || true
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
echo "KDSE WorkOrder Acceptance Tests"
echo "Runtime Owns Methodology"
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

# Configure git identity for tests
git config --global user.email "test@test.com" 2>/dev/null || true
git config --global user.name "Test User" 2>/dev/null || true

# Test 1: Initialize creates session
echo ""
info "Test 1: Initialize creates orchestration session"
setup_test_repo

# Send initialize request
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"initialize"},"id":1}' | $MCP_SERVER > /tmp/init_output.json

if grep -q 'session_id' /tmp/init_output.json && grep -q 'strict_mode' /tmp/init_output.json; then
    pass "Session initialized with STRICT mode"
else
    fail "Session not properly initialized"
fi

# Test 2: First execute returns WorkOrder with Problem phase (first phase)
echo ""
info "Test 2: First execute() returns Problem WorkOrder (first phase)"

# Execute with an objective
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"execute","arguments":{"objective":"Create a supermarket inventory system"}},"id":2}' | $MCP_SERVER > /tmp/execute1_output.json

# Check that WorkOrder exists
if grep -q 'work_order' /tmp/execute1_output.json; then
    pass "WorkOrder is present in execute response"
else
    fail "WorkOrder missing from execute response"
fi

# Check that WorkOrder contains Problem phase (first phase)
if grep -q 'Problem' /tmp/execute1_output.json; then
    pass "WorkOrder phase is Problem (first phase)"
else
    fail "WorkOrder phase is not Problem"
fi

# Test 3: WorkOrder contains blocked actions
echo ""
info "Test 3: WorkOrder contains blocked actions"

if grep -q 'blocked_actions' /tmp/execute1_output.json; then
    pass "Blocked actions field exists"
else
    fail "Blocked actions field missing"
fi

if grep -q 'DO NOT generate any code' /tmp/execute1_output.json; then
    pass "Source code generation is blocked"
else
    fail "Source code generation not explicitly blocked"
fi

if grep -q 'DO NOT create project' /tmp/execute1_output.json; then
    pass "Project structure creation is blocked"
else
    fail "Project structure creation not explicitly blocked"
fi

# Test 4: WorkOrder contains completion criteria
echo ""
info "Test 4: WorkOrder contains completion criteria"

if grep -q 'completion_criteria' /tmp/execute1_output.json; then
    pass "Completion criteria field exists"
else
    fail "Completion criteria field missing"
fi

# Test 5: WorkOrder contains expected deliverables
echo ""
info "Test 5: WorkOrder contains expected deliverables"

if grep -q 'expected_deliverables' /tmp/execute1_output.json; then
    pass "Expected deliverables field exists"
else
    fail "Expected deliverables field missing"
fi

# Test 7: No source code created after Foundation WorkOrder
echo ""
info "Test 7: No source code should be created in repository (Foundation phase only)"

# Count source code files
SOURCE_FILES=$(find . -type f \( -name "*.go" -o -name "*.js" -o -name "*.py" -o -name "*.java" -o -name "*.ts" \) 2>/dev/null | wc -l)
if [ "$SOURCE_FILES" -eq 0 ]; then
    pass "No source code files in repository (correct behavior)"
else
    fail "Found $SOURCE_FILES source code files - Foundation WorkOrder violated"
fi

# Test 8: No folders outside .kdse/
echo ""
info "Test 8: No folders created outside .kdse/"

# Check for directories outside .kdse
OUTSIDE_KDSE=$(find . -maxdepth 1 -type d ! -name ".git" ! -name ".kdse" ! -name "." 2>/dev/null | wc -l)
if [ "$OUTSIDE_KDSE" -eq 0 ]; then
    pass "No directories created outside .kdse/"
else
    fail "Directories created outside .kdse/ - Foundation WorkOrder violated"
fi

# Test 9: Execute again advances to Knowledge phase
echo ""
info "Test 9: Second execute() advances to Knowledge phase"

echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"execute","arguments":{"objective":"Create a supermarket inventory system"}},"id":3}' | $MCP_SERVER > /tmp/execute2_output.json

if grep -q 'Knowledge' /tmp/execute2_output.json; then
    pass "Phase advanced to Knowledge Collection"
else
    fail "Phase did not advance correctly"
fi

# Test 10: Knowledge WorkOrder has its own blocked actions
echo ""
info "Test 10: Knowledge WorkOrder has appropriate blocked actions"

if grep -q 'DO NOT modify' /tmp/execute2_output.json; then
    pass "Knowledge WorkOrder blocks code modification"
else
    fail "Knowledge WorkOrder should block code modification"
fi

# Test 11: Third execute advances to Foundation
echo ""
info "Test 11: Third execute() advances to Foundation"

echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"execute","arguments":{"objective":"Create a supermarket inventory system"}},"id":4}' | $MCP_SERVER > /tmp/execute3_output.json

if grep -q 'Foundation' /tmp/execute3_output.json; then
    pass "Phase advanced to Foundation"
else
    fail "Phase did not advance to Foundation"
fi

# Test 12: Foundation WorkOrder is explicit about documents
echo ""
info "Test 12: Foundation WorkOrder explicitly lists all 5 documents"

FOUNDATIONS_DOCS=("SPEC.md" "REQUIREMENTS.md" "ASSUMPTIONS.md" "CONSTRAINTS.md" "GLOSSARY.md")
for doc in "${FOUNDATIONS_DOCS[@]}"; do
    if grep -q "$doc" /tmp/execute3_output.json; then
        pass "Foundation WorkOrder lists $doc"
    else
        fail "Foundation WorkOrder missing $doc"
    fi
done

# Test 13: Strict mode enforcement message
echo ""
info "Test 13: WorkOrder indicates strict mode enforcement"

if grep -q 'strict_mode_enforced' /tmp/execute3_output.json; then
    pass "Strict mode is enforced in WorkOrder"
else
    fail "Strict mode not enforced in WorkOrder"
fi

# Cleanup test repo
cleanup_test_repo

# Test 14: Help text mentions WorkOrders
echo ""
info "Test 14: Help mentions WorkOrders and methodology"

cd /workspace/project/KDSE
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"help"},"id":5}' | $MCP_SERVER > /tmp/help_output.json

if grep -q "WorkOrder" /tmp/help_output.json || grep -q "work_order" /tmp/help_output.json; then
    pass "Help mentions WorkOrders"
else
    info "Help does not mention WorkOrders (optional - can be added)"
fi

# Summary
echo ""
echo "================================================"
echo "Test Summary"
echo "================================================"
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Failed: $TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed! Runtime Owns Methodology.${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed.${NC}"
    exit 1
fi
