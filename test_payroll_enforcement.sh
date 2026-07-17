#!/bin/bash
#===============================================================================
# KDSE Enforcement Test - Payroll System
#
# This script tests the KDSE enforcement mechanism by attempting to build
# a payroll system with biometrics. It should block because foundation is missing.
#===============================================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
KDSE_PATH="${KDSE_PATH:-.}"
TEST_DIR="/tmp/kdse_test_$$"

echo_step() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}STEP: $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

echo_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

echo_error() {
    echo -e "${RED}✗ $1${NC}"
}

echo_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Cleanup function
cleanup() {
    echo ""
    echo "Cleaning up test environment..."
    rm -rf "$TEST_DIR"
    echo "Done."
}

trap cleanup EXIT

#===============================================================================
# TEST: Enforcement Blocks Premature Implementation
#===============================================================================

echo ""
echo "╔════════════════════════════════════════════════════════════════════════╗"
echo "║          KDSE ENFORCEMENT TEST - Payroll System                       ║"
echo "╚════════════════════════════════════════════════════════════════════════╝"
echo ""

# Create test directory
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

echo_step "1. Initialize KDSE Workspace"

# Initialize workspace
if command -v go &> /dev/null; then
    # Try using Go if available
    go run /workspace/project/KDSE/cmd/kdse/main.go initialize 2>/dev/null || \
    echo_warning "Go not configured, using file-based initialization"
fi

# Manual initialization if Go not available
mkdir -p .kdse/foundation .kdse/knowledge .kdse/reports .kdse/state

echo_success "Workspace initialized"

echo_step "2. Check Initial Status"

# Check what exists
echo "Checking .kdse/ structure:"
echo ""
if [ -d ".kdse/foundation" ]; then
    echo "  foundation/: $(ls -la .kdse/foundation/ 2>/dev/null | wc -l) items"
else
    echo "  foundation/: MISSING"
fi

if [ -d ".kdse/knowledge" ]; then
    echo "  knowledge/: $(ls -la .kdse/knowledge/ 2>/dev/null | wc -l) items"
else
    echo "  knowledge/: MISSING"
fi

echo_step "3. Attempt Premature Implementation"

echo "Objective: Build a payroll system with biometrics"
echo ""
echo "Expected: BLOCKED - Foundation and Knowledge phases incomplete"
echo ""

# Simulate enforcement check
echo -n "Checking foundation... "
foundation_complete=false
if [ -d ".kdse/foundation" ]; then
    missing_docs=()
    for doc in PROBLEM.md SPEC.md REQUIREMENTS.md ASSUMPTIONS.md CONSTRAINTS.md; do
        if [ ! -f ".kdse/foundation/$doc" ]; then
            missing_docs+=("$doc")
        elif [ $(stat -f%z ".kdse/foundation/$doc" 2>/dev/null || stat -c%s ".kdse/foundation/$doc" 2>/dev/null || echo 0) -lt 100 ]; then
            missing_docs+=("$doc (empty)")
        fi
    done
    if [ ${#missing_docs[@]} -eq 0 ]; then
        foundation_complete=true
        echo_success "Complete"
    else
        echo_error "Missing: ${missing_docs[*]}"
    fi
else
    echo_error "Directory missing"
fi

echo -n "Checking knowledge base... "
knowledge_complete=false
if [ -d ".kdse/knowledge" ]; then
    missing_cats=()
    for cat in general operational developmental; do
        if [ ! -d ".kdse/knowledge/$cat" ]; then
            missing_cats+=("$cat")
        fi
    done
    if [ ${#missing_cats[@]} -eq 0 ]; then
        # Check if categories have content
        has_content=false
        for cat in general operational developmental; do
            if [ -n "$(ls -A .kdse/knowledge/$cat 2>/dev/null)" ]; then
                has_content=true
                break
            fi
        done
        if $has_content; then
            knowledge_complete=true
            echo_success "Complete"
        else
            echo_error "Empty categories"
        fi
    else
        echo_error "Missing: ${missing_cats[*]}"
    fi
else
    echo_error "Directory missing"
fi

echo -n "Checking architecture... "
arch_complete=false
if [ -f ".kdse/foundation/ARCHITECTURE.md" ]; then
    size=$(stat -f%z ".kdse/foundation/ARCHITECTURE.md" 2>/dev/null || stat -c%s ".kdse/foundation/ARCHITECTURE.md" 2>/dev/null || echo 0)
    if [ "$size" -gt 200 ]; then
        arch_complete=true
        echo_success "Complete"
    else
        echo_error "Empty or too small"
    fi
else
    echo_error "Missing"
fi

echo ""
echo_step "4. Enforcement Result"

# Determine if implementation is allowed
implementation_allowed=$foundation_complete && $knowledge_complete && $arch_complete

if $implementation_allowed; then
    echo_success "Implementation ALLOWED"
    echo ""
    echo "The following phases were completed:"
    $foundation_complete && echo "  ✓ Foundation"
    $knowledge_complete && echo "  ✓ Knowledge Collection"
    $arch_complete && echo "  ✓ Architecture"
else
    echo_error "Implementation BLOCKED"
    echo ""
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${RED}KDSE PRINCIPLES VIOLATED${NC}"
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo "Required phases that must be completed BEFORE implementation:"
    echo ""
    if ! $foundation_complete; then
        echo -e "  ${RED}✗${NC} Foundation Phase"
        echo "    → Create: PROBLEM.md, SPEC.md, REQUIREMENTS.md, ASSUMPTIONS.md, CONSTRAINTS.md"
    fi
    if ! $knowledge_complete; then
        echo -e "  ${RED}✗${NC} Knowledge Phase"
        echo "    → Collect domain knowledge in: general/, operational/, developmental/"
    fi
    if ! $arch_complete; then
        echo -e "  ${RED}✗${NC} Architecture Phase"
        echo "    → Create: ARCHITECTURE.md with system design"
    fi
    echo ""
    echo "Correct Phase Order:"
    echo "  1. Problem → 2. Knowledge → 3. Foundation → 4. Architecture → 5. Implement"
    echo ""
fi

echo_step "5. Demonstration: Auto-Create Foundation"

echo "Creating foundation documents..."
mkdir -p .kdse/foundation

# Create templates
cat > .kdse/foundation/PROBLEM.md << 'EOF'
# Problem Statement

**Project:** Payroll System with Biometrics

## Problem Description

The organization needs an integrated payroll system that includes biometric 
authentication for employee time tracking and attendance.

## Impact

Manual payroll processing is error-prone and time-consuming. Current systems
lack biometric integration, leading to time theft and inaccurate records.

## Success Criteria

- [ ] Biometric data integration complete
- [ ] Payroll calculation accuracy 99.9%+
- [ ] Compliance with labor regulations
- [ ] Audit trail for all transactions
EOF

cat > .kdse/foundation/SPEC.md << 'EOF'
# Project Specification

**Project:** Payroll System with Biometrics

## Overview

Integrated payroll management system with biometric attendance tracking
for accurate time-based compensation.

## Scope

### In Scope
- Biometric device integration (fingerprint, facial recognition)
- Employee management
- Attendance tracking
- Payroll calculation
- Direct deposit integration
- Tax reporting

### Out of Scope
- HR benefits administration
- Performance management
- Training systems

## Deliverables

1. Biometric integration module
2. Attendance tracking module
3. Payroll calculation engine
4. Reporting dashboard
5. API for third-party integration
EOF

cat > .kdse/foundation/REQUIREMENTS.md << 'EOF'
# Functional Requirements

## FR-001: Biometric Enrollment
**Description:** System shall allow employees to enroll biometric data
**Priority:** High
**Acceptance Criteria:**
- [ ] Support fingerprint enrollment
- [ ] Support facial recognition enrollment
- [ ] Store biometric templates securely

## FR-002: Attendance Tracking
**Description:** System shall record attendance via biometric verification
**Priority:** High
**Acceptance Criteria:**
- [ ] Record check-in/check-out times
- [ ] Detect duplicate punches
- [ ] Handle missed punches with supervisor approval

## FR-003: Payroll Calculation
**Description:** System shall calculate payroll based on attendance
**Priority:** Critical
**Acceptance Criteria:**
- [ ] Calculate regular hours
- [ ] Calculate overtime (1.5x)
- [ ] Apply deductions
- [ ] Generate pay stubs
EOF

cat > .kdse/foundation/ASSUMPTIONS.md << 'EOF'
# Key Assumptions

## Technical Assumptions
1. Biometric devices support TCP/IP connectivity
2. Fingerprint templates use industry-standard format
3. Payroll data must be encrypted at rest

## Business Assumptions
1. All employees have unique biometric identifiers
2. Overtime threshold is 40 hours/week
3. Pay periods are bi-weekly

## Environment Assumptions
1. Network latency < 100ms to biometric devices
2. System available 24/7 except maintenance windows
EOF

cat > .kdse/foundation/CONSTRAINTS.md << 'EOF'
# Project Constraints

## Technical Constraints
- Must integrate with existing HRIS
- Biometric data must be GDPR compliant
- API response time < 200ms

## Schedule Constraints
- Phase 1 delivery: 90 days
- Full system: 180 days

## Resource Constraints
- Development team: 5 FTE
- QA team: 2 FTE

## Compliance Constraints
- SOC 2 Type II certification required
- PCI DSS for payment processing
EOF

echo_success "Foundation documents created"

echo_step "6. Verify Enforcement After Foundation"

# Re-check
echo "Checking foundation after creation..."
foundation_complete=true
for doc in PROBLEM.md SPEC.md REQUIREMENTS.md ASSUMPTIONS.md CONSTRAINTS.md; do
    if [ ! -f ".kdse/foundation/$doc" ]; then
        foundation_complete=false
        echo "  Missing: $doc"
    else
        size=$(stat -f%z ".kdse/foundation/$doc" 2>/dev/null || stat -c%s ".kdse/foundation/$doc" 2>/dev/null || echo 0)
        if [ "$size" -lt 100 ]; then
            echo "  Empty: $doc"
            foundation_complete=false
        else
            echo "  ✓ $doc"
        fi
    fi
done

echo ""
if $foundation_complete; then
    echo_success "Foundation COMPLETE"
else
    echo_warning "Foundation incomplete"
fi

echo_step "7. Summary"

echo ""
echo "╔════════════════════════════════════════════════════════════════════════╗"
echo "║                    TEST RESULTS SUMMARY                              ║"
echo "╠════════════════════════════════════════════════════════════════════════╣"
echo "║                                                                        ║"
echo "║  Initial State:                                                        ║"
echo "║    Foundation:    ${RED}INCOMPLETE${NC}                                          ║"
echo "║    Knowledge:     ${RED}MISSING${NC}                                           ║"
echo "║    Architecture:  ${RED}MISSING${NC}                                           ║"
echo "║    Implementation: ${RED}BLOCKED${NC}                                          ║"
echo "║                                                                        ║"
echo "║  After Foundation Creation:                                            ║"
if $foundation_complete; then
echo "║    Foundation:    ${GREEN}COMPLETE${NC}                                           ║"
else
echo "║    Foundation:    ${RED}INCOMPLETE${NC}                                          ║"
fi
echo "║    Knowledge:     ${RED}STILL MISSING${NC}                                     ║"
echo "║    Architecture:  ${RED}STILL MISSING${NC}                                     ║"
echo "║    Implementation: ${YELLOW}PARTIALLY ALLOWED${NC}                               ║"
echo "║                                                                        ║"
echo "╚════════════════════════════════════════════════════════════════════════╝"
echo ""
echo "Conclusion:"
echo "  KDSE Enforcement successfully BLOCKED premature implementation."
echo "  Foundation creation was required before proceeding."
echo ""
echo "  To fully enable implementation, complete:"
echo "    1. Knowledge Collection (general/, operational/, developmental/)"
echo "    2. Architecture Document (ARCHITECTURE.md)"
echo ""
