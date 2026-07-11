#!/usr/bin/env bash
#===============================================================================
# KDSE Runtime Verification Script
#===============================================================================
# Purpose: Verify KDSE Runtime installation integrity
# Design:  Deterministic, idempotent, no AI logic, no engineering decisions
#===============================================================================

set -euo pipefail

# Load common utilities
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/common.sh"

#-------------------------------------------------------------------------------
# Usage
#-------------------------------------------------------------------------------

usage() {
    cat <<EOF
Usage: $(get_script_name) [OPTIONS]

Verify KDSE Runtime installation integrity.

OPTIONS:
    -h, --help          Show this help message
    -v, --verbose       Verbose output
    -q, --quiet         Quiet mode (only show PASS/FAIL)
    -j, --json          Output in JSON format

CHECKS PERFORMED:
    1. Directory layout
    2. Mandatory documents
    3. Manifest integrity
    4. Configuration validity
    5. Version consistency

EXAMPLES:
    $(get_script_name)              # Full verification
    $(get_script_name) -v           # Verbose verification
    $(get_script_name) -q           # Quiet mode
    $(get_script_name) -j           # JSON output

EXIT CODES:
    0 - PASS    All checks passed
    1 - FAIL    One or more checks failed
    2 - ERROR   Cannot verify (not installed)

EOF
    exit $EXIT_OK
}

#-------------------------------------------------------------------------------
# Global State
#-------------------------------------------------------------------------------

CHECKS_PASSED=0
CHECKS_FAILED=0
CHECKS_WARNING=0
JSON_OUTPUT=0
QUIET_MODE=0
RESULTS=()

#-------------------------------------------------------------------------------
# Parse Arguments
#-------------------------------------------------------------------------------

while [[ $# -gt 0 ]]; do
    case "$1" in
        -h|--help)
            usage
            ;;
        -v|--verbose)
            VERBOSE=1
            shift
            ;;
        -q|--quiet)
            QUIET_MODE=1
            shift
            ;;
        -j|--json)
            JSON_OUTPUT=1
            shift
            ;;
        *)
            log_error "Unknown option: $1"
            usage
            ;;
    esac
done

#-------------------------------------------------------------------------------
# Result Recording
#-------------------------------------------------------------------------------

record_check() {
    local check_name="$1"
    local status="$2"
    local message="${3:-}"
    
    RESULTS+=("{\"check\":\"$check_name\",\"status\":\"$status\",\"message\":\"$(json_escape "$message")\"}")
    
    case "$status" in
        PASS)
            CHECKS_PASSED=$((CHECKS_PASSED + 1))
            if [[ "$VERBOSE" == "1" ]]; then
                log_success "$check_name"
                [[ -n "$message" ]] && log_info "  $message"
            fi
            ;;
        FAIL)
            CHECKS_FAILED=$((CHECKS_FAILED + 1))
            log_error "$check_name"
            [[ -n "$message" ]] && log_error "  $message"
            ;;
        WARN)
            CHECKS_WARNING=$((CHECKS_WARNING + 1))
            log_warning "$check_name"
            [[ -n "$message" ]] && log_warning "  $message"
            ;;
    esac
}

#-------------------------------------------------------------------------------
# Check 1: Directory Layout
#-------------------------------------------------------------------------------

check_directory_layout() {
    log_info "Checking directory layout..."
    
    local install_path=$(get_install_path)
    local errors=0
    local missing=()
    
    # Check root directory exists
    if [[ ! -d "$install_path" ]]; then
        record_check "Directory Layout" "FAIL" "Installation directory not found: $install_path"
        return 0
    fi
    
    # Check required directories
    for dir in "${REQUIRED_DIRS[@]}"; do
        if [[ ! -d "${install_path}/${dir}" ]]; then
            missing+=("$dir")
            errors=$((errors + 1))
        fi
    done
    
    if [[ $errors -gt 0 ]]; then
        record_check "Directory Layout" "FAIL" "Missing directories: ${missing[*]}"
        return 0
    fi
    
    record_check "Directory Layout" "PASS" "All required directories present"
    return 0
}

#-------------------------------------------------------------------------------
# Check 2: Mandatory Documents
#-------------------------------------------------------------------------------

check_mandatory_documents() {
    log_info "Checking mandatory documents..."
    
    local install_path=$(get_install_path)
    local errors=0
    local missing=()
    local present=()
    
    for doc in "${NORMATIVE_DOCS[@]}"; do
        local doc_name=$(basename "$doc")
        local doc_path="${install_path}/standards/normative/${doc_name}"
        
        if [[ -f "$doc_path" ]]; then
            present+=("$doc_name")
        else
            missing+=("$doc_name")
            errors=$((errors + 1))
        fi
    done
    
    if [[ $errors -gt 0 ]]; then
        record_check "Mandatory Documents" "FAIL" "Missing: ${missing[*]}"
        return 0
    fi
    
    record_check "Mandatory Documents" "PASS" "${#present[@]} normative documents present"
    return 0
}

#-------------------------------------------------------------------------------
# Check 3: Manifest Integrity
#-------------------------------------------------------------------------------

check_manifest() {
    log_info "Checking manifest..."
    
    local manifest=$(get_manifest_path)
    
    # Check manifest exists
    if [[ ! -f "$manifest" ]]; then
        record_check "Manifest" "FAIL" "Manifest file not found"
        return 0
    fi
    
    # Check manifest is valid JSON
    if ! grep -q '"version"' "$manifest" 2>/dev/null; then
        record_check "Manifest" "FAIL" "Invalid manifest format"
        return 0
    fi
    
    # Extract and verify required fields
    local version=$(manifest_get "version")
    local installed=$(manifest_get "installed")
    local repo=$(manifest_get "repo")
    local branch=$(manifest_get "branch")
    
    local missing_fields=()
    [[ -z "$version" ]] && missing_fields+=("version")
    [[ -z "$installed" ]] && missing_fields+=("installed")
    [[ -z "$repo" ]] && missing_fields+=("repo")
    [[ -z "$branch" ]] && missing_fields+=("branch")
    
    if [[ ${#missing_fields[@]} -gt 0 ]]; then
        record_check "Manifest" "FAIL" "Missing fields: ${missing_fields[*]}"
        return 0
    fi
    
    record_check "Manifest" "PASS" "Version: $version, Branch: $branch"
    return 0
}

#-------------------------------------------------------------------------------
# Check 4: Configuration Validity
#-------------------------------------------------------------------------------

check_configuration() {
    log_info "Checking configuration..."
    
    local config=$(get_config_path)
    
    # Check configuration exists
    if [[ ! -f "$config" ]]; then
        record_check "Configuration" "FAIL" "Configuration file not found"
        return 0
    fi
    
    # Check configuration is readable
    if [[ ! -r "$config" ]]; then
        record_check "Configuration" "FAIL" "Configuration file not readable"
        return 0
    fi
    
    # Check for obvious syntax errors (basic check)
    if grep -qE '^\s*(eval|exec|;)\s+' "$config" 2>/dev/null; then
        record_check "Configuration" "FAIL" "Configuration contains unsafe commands"
        return 0
    fi
    
    # Load and verify key variables
    load_config
    
    if [[ -z "${KDSE_INSTALL_PATH:-}" ]]; then
        record_check "Configuration" "WARN" "KDSE_INSTALL_PATH not set"
        return 0
    fi
    
    record_check "Configuration" "PASS" "Configuration valid"
    return 0
}

#-------------------------------------------------------------------------------
# Check 5: Version Consistency
#-------------------------------------------------------------------------------

check_version_consistency() {
    log_info "Checking version consistency..."
    
    local install_path=$(get_install_path)
    local manifest=$(get_manifest_path)
    
    # Get version from manifest
    local manifest_version=$(manifest_get "version" 2>/dev/null || echo "unknown")
    
    # Get version from git
    local git_version=$(get_git_version)
    
    # Get version from config if available
    load_config
    local config_version="${KDSE_VERSION:-unknown}"
    
    # Check if versions are consistent
    if [[ "$manifest_version" == "$git_version" ]]; then
        record_check "Version Consistency" "PASS" "Manifest and git in sync: $git_version"
        return 0
    fi
    
    # Version mismatch - may need sync
    record_check "Version Consistency" "WARN" "Manifest ($manifest_version) differs from git ($git_version)"
    return 0
}

#-------------------------------------------------------------------------------
# Check 6: Preserved Directories
#-------------------------------------------------------------------------------

check_preserved_directories() {
    log_info "Checking preserved directories..."
    
    local preserved=(
        "reports"
        "history"
        "runtime"
        "cache"
        "configuration"
    )
    
    local install_path=$(get_install_path)
    local all_present=1
    
    for dir in "${preserved[@]}"; do
        if [[ ! -d "${install_path}/${dir}" ]]; then
            log_warning "Preserved directory missing: $dir"
            all_present=0
        fi
    done
    
    if [[ $all_present -eq 1 ]]; then
        record_check "Preserved Directories" "PASS" "All preserved directories present"
        return 0
    fi
    
    record_check "Preserved Directories" "WARN" "Some preserved directories missing"
    return 0
}

#-------------------------------------------------------------------------------
# Check 7: Standards Structure
#-------------------------------------------------------------------------------

check_standards_structure() {
    log_info "Checking standards structure..."
    
    local install_path=$(get_install_path)
    local standards_dir="${install_path}/standards"
    
    # Check standards directory
    if [[ ! -d "$standards_dir" ]]; then
        record_check "Standards Structure" "FAIL" "Standards directory not found"
        return 0
    fi
    
    # Check normative subdirectory
    if [[ ! -d "${standards_dir}/normative" ]]; then
        record_check "Standards Structure" "FAIL" "Normative directory not found"
        return 0
    fi
    
    # Check informative subdirectory
    if [[ ! -d "${standards_dir}/informative" ]]; then
        record_check "Standards Structure" "WARN" "Informative directory not found"
        return 0
    fi
    
    record_check "Standards Structure" "PASS" "Standards structure valid"
    return 0
}

#-------------------------------------------------------------------------------
# Print Results
#-------------------------------------------------------------------------------

print_results() {
    local exit_code=0
    
    if [[ $CHECKS_FAILED -gt 0 ]]; then
        exit_code=1
    fi
    
    if [[ "$JSON_OUTPUT" == "1" ]]; then
        # JSON output
        cat <<EOF
{
    "verification": {
        "status": "$([ $CHECKS_FAILED -eq 0 ] && echo "PASS" || echo "FAIL")",
        "timestamp": "$(get_timestamp)",
        "install_path": "$(get_install_path)",
        "summary": {
            "passed": $CHECKS_PASSED,
            "failed": $CHECKS_FAILED,
            "warnings": $CHECKS_WARNING
        },
        "checks": [$(IFS=,; echo "${RESULTS[*]}")]
    }
}
EOF
    elif [[ "$QUIET_MODE" == "1" ]]; then
        # Quiet mode - just show PASS/FAIL
        if [[ $CHECKS_FAILED -eq 0 ]]; then
            echo "PASS"
        else
            echo "FAIL"
        fi
    else
        # Standard output
        print_summary_header "Verification Summary"
        
        echo ""
        echo "Installation Path: $(get_install_path)"
        echo ""
        echo "Results:"
        echo "  PASSED:  $CHECKS_PASSED"
        echo "  FAILED:  $CHECKS_FAILED"
        echo "  WARNINGS: $CHECKS_WARNING"
        echo ""
        
        if [[ $CHECKS_FAILED -eq 0 ]]; then
            log_success "All verification checks passed"
        else
            log_error "Verification failed with $CHECKS_FAILED error(s)"
        fi
        
        echo ""
        print_summary_footer
    fi
    
    return $exit_code
}

#-------------------------------------------------------------------------------
# Main
#-------------------------------------------------------------------------------

main() {
    # Check if installed - with legacy support
    if ! manifest_exists; then
        if [[ "$JSON_OUTPUT" == "1" ]]; then
            cat <<EOF
{
    "verification": {
        "status": "ERROR",
        "timestamp": "$(get_timestamp)",
        "error": "KDSE Runtime is not installed"
    }
}
EOF
        else
            print_summary_header "Verification"
            log_error "KDSE Runtime is not installed"
            log_info "Run install.sh first"
            print_summary_footer
        fi
        exit $EXIT_NOT_INSTALLED
    fi
    
    # Auto-detect and migrate legacy installations before verification
    auto_detect_and_migrate || {
        log_warning "Failed to migrate legacy installation, continuing with verification"
    }
    
    print_summary_header "KDSE Runtime Verification"
    echo ""
    
    # Show detection info
    local format=$(detect_manifest_format)
    echo "Detected Installation Format: $format"
    if [[ "$MIGRATION_PERFORMED" == "1" ]]; then
        echo "Migration Performed: YES (${MIGRATION_SOURCE_FORMAT} -> JSON)"
        echo ""
    fi
    
    # Run all checks
    check_directory_layout
    check_mandatory_documents
    check_manifest
    check_configuration
    check_version_consistency
    check_preserved_directories
    check_standards_structure
    
    echo ""
    print_results
}

#-------------------------------------------------------------------------------
# Auto-Detect and Migrate Legacy Installations
#-------------------------------------------------------------------------------

auto_detect_and_migrate() {
    local format=$(detect_manifest_format)
    
    case "$format" in
        json)
            log_info "Detected installation format: JSON (current)"
            return 0
            ;;
        yaml)
            log_info "Detected installation format: YAML (legacy)"
            log_info "Automatic migration will be performed..."
            migrate_manifest_yaml_to_json || {
                log_error "Migration failed"
                return 1
            }
            return 0
            ;;
        none)
            log_info "No existing KDSE installation detected"
            return 0
            ;;
    esac
}

# Run
main "$@"
