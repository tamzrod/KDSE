#!/usr/bin/env bash
#===============================================================================
# KDSE Runtime Installation Framework - Common Functions
#===============================================================================
# Purpose: Shared utilities for KDSE installation, sync, verify, uninstall
# Design:  Deterministic, idempotent, no AI logic, no engineering decisions
#===============================================================================

#-------------------------------------------------------------------------------
# Configuration
#-------------------------------------------------------------------------------

# KDSE home directory (where .kdse will be created)
KDSE_HOME="${KDSE_HOME:-$HOME}"

# Installation directory name
KDSE_DIR=".kdse"

# Manifest file name
MANIFEST_FILE="manifest.json"

# Configuration file name
CONFIG_FILE="config.sh"

# Default KDSE repository
KDSE_REPO="${KDSE_REPO:-https://github.com/tamzrod/KDSE.git}"
KDSE_BRANCH="${KDSE_BRANCH:-main}"

#-------------------------------------------------------------------------------
# Color Output
#-------------------------------------------------------------------------------

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Output functions - silent by default, verbose on failure
VERBOSE="${VERBOSE:-0}"
DEBUG="${DEBUG:-0}"

log_info() {
    if [[ "$VERBOSE" == "1" ]]; then
        echo -e "${BLUE}[INFO]${NC} $*"
    fi
}

log_success() {
    echo -e "${GREEN}[OK]${NC} $*"
}

log_warning() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*" >&2
}

log_debug() {
    if [[ "$DEBUG" == "1" ]]; then
        echo -e "[DEBUG] $*" >&2
    fi
}

#-------------------------------------------------------------------------------
# Path Utilities
#-------------------------------------------------------------------------------

# Get full installation path
get_install_path() {
    echo "${KDSE_HOME}/${KDSE_DIR}"
}

# Get manifest path
get_manifest_path() {
    echo "$(get_install_path)/${MANIFEST_FILE}"
}

# Get configuration path
get_config_path() {
    echo "$(get_install_path)/${CONFIG_FILE}"
}

# Get reports path
get_reports_path() {
    echo "$(get_install_path)/reports"
}

# Get history path
get_history_path() {
    echo "$(get_install_path)/history"
}

# Get runtime path
get_runtime_path() {
    echo "$(get_install_path)/runtime"
}

# Get cache path
get_cache_path() {
    echo "$(get_install_path)/cache"
}

#-------------------------------------------------------------------------------
# Directory Structure
#-------------------------------------------------------------------------------

# Required directories to create
REQUIRED_DIRS=(
    "reports"
    "history"
    "runtime"
    "cache"
    "configuration"
    "standards/normative"
    "standards/informative"
)

#-------------------------------------------------------------------------------
# Normative Documents (must be installed)
#-------------------------------------------------------------------------------

# Core normative documents from KDSE repository
NORMATIVE_DOCS=(
    "docs/foundation/003-core-principles.md"
    "docs/foundation/006-chain-of-authority.md"
    "docs/foundation/004-engineering-model.md"
    "docs/audit/COMPLIANCE_AUDIT.md"
    "docs/audit/FOUNDATION_AUDIT.md"
    "docs/audit/AUDIT_SCORING.md"
)

#-------------------------------------------------------------------------------
# Informative Documents (installed for reference)
#-------------------------------------------------------------------------------

INFORMATIVE_DOCS=(
    "docs/foundation/000-what-is-kdse.md"
    "docs/foundation/001-why-kdse-exists.md"
    "docs/foundation/002-scope.md"
    "docs/foundation/007-glossary.md"
    "runtime/README.md"
    "runtime/ARCHITECTURE.md"
    "runtime/EXECUTION_MODEL.md"
)

#-------------------------------------------------------------------------------
# File Checksums
#-------------------------------------------------------------------------------

# Calculate SHA256 checksum of a file
calculate_checksum() {
    local file="$1"
    if [[ -f "$file" ]]; then
        if command -v sha256sum &>/dev/null; then
            sha256sum "$file" | awk '{print $1}'
        elif command -v shasum &>/dev/null; then
            shasum -a 256 "$file" | awk '{print $1}'
        elif command -v openssl &>/dev/null; then
            openssl dgst -sha256 "$file" | awk '{print $2}'
        else
            log_error "No checksum tool available"
            return 1
        fi
    else
        return 1
    fi
}

#-------------------------------------------------------------------------------
# JSON Utilities
#-------------------------------------------------------------------------------

# Simple JSON string escaping
json_escape() {
    local str="$1"
    str="${str//\\/\\\\}"
    str="${str//\"/\\\"}"
    str="${str//$'\n'/\\n}"
    str="${str//$'\r'/\\r}"
    str="${str//$'\t'/\\t}"
    echo "$str"
}

#-------------------------------------------------------------------------------
# Version Utilities
#-------------------------------------------------------------------------------

# Get current timestamp in ISO format
get_timestamp() {
    date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u +"%Y-%m-%d %H:%M:%S UTC"
}

# Get git commit hash if available
get_git_version() {
    if command -v git &>/dev/null && [[ -d ".git" ]]; then
        git rev-parse HEAD 2>/dev/null | head -c 8
    else
        echo "unknown"
    fi
}

#-------------------------------------------------------------------------------
# Directory Operations
#-------------------------------------------------------------------------------

# Create directory if it doesn't exist
ensure_directory() {
    local dir="$1"
    if [[ ! -d "$dir" ]]; then
        mkdir -p "$dir" || {
            log_error "Failed to create directory: $dir"
            return 1
        }
        log_debug "Created directory: $dir"
    fi
    return 0
}

# Remove directory if empty or force
remove_directory() {
    local dir="$1"
    local force="${2:-0}"
    
    if [[ ! -d "$dir" ]]; then
        return 0
    fi
    
    if [[ "$force" == "1" ]]; then
        rm -rf "$dir" || {
            log_error "Failed to remove directory: $dir"
            return 1
        }
    else
        # Only remove if empty
        if [[ -z "$(ls -A "$dir" 2>/dev/null)" ]]; then
            rmdir "$dir" || {
                log_error "Failed to remove empty directory: $dir"
                return 1
            }
        else
            log_warning "Directory not empty, skipping: $dir"
            return 1
        fi
    fi
    return 0
}

#-------------------------------------------------------------------------------
# File Operations
#-------------------------------------------------------------------------------

# Copy file preserving attributes
copy_file() {
    local src="$1"
    local dest="$2"
    
    if [[ ! -f "$src" ]]; then
        log_error "Source file not found: $src"
        return 1
    fi
    
    ensure_directory "$(dirname "$dest")" || return 1
    
    cp -p "$src" "$dest" || {
        log_error "Failed to copy $src to $dest"
        return 1
    }
    
    log_debug "Copied: $src -> $dest"
    return 0
}

# Remove file if exists
remove_file() {
    local file="$1"
    
    if [[ -f "$file" ]]; then
        rm -f "$file" || {
            log_error "Failed to remove file: $file"
            return 1
        }
        log_debug "Removed: $file"
    fi
    return 0
}

#-------------------------------------------------------------------------------
# Manifest Operations
#-------------------------------------------------------------------------------

# Read manifest and extract value by key
manifest_get() {
    local key="$1"
    local manifest=$(get_manifest_path)
    
    if [[ ! -f "$manifest" ]]; then
        return 1
    fi
    
    # Simple grep-based extraction for basic values
    case "$key" in
        version)
            grep '"version"' "$manifest" | head -1 | sed 's/.*: *"\([^"]*\)".*/\1/'
            ;;
        installed)
            grep '"installed"' "$manifest" | head -1 | sed 's/.*: *"\([^"]*\)".*/\1/'
            ;;
        repo)
            grep '"repo"' "$manifest" | head -1 | sed 's/.*: *"\([^"]*\)".*/\1/'
            ;;
        branch)
            grep '"branch"' "$manifest" | head -1 | sed 's/.*: *"\([^"]*\)".*/\1/'
            ;;
        *)
            return 1
            ;;
    esac
}

# Check if manifest exists and is valid
manifest_exists() {
    local manifest=$(get_manifest_path)
    [[ -f "$manifest" ]] && grep -q '"kdse"' "$manifest"
}

#-------------------------------------------------------------------------------
# Configuration Operations
#-------------------------------------------------------------------------------

# Source configuration file if it exists
load_config() {
    local config=$(get_config_path)
    
    if [[ -f "$config" ]]; then
        # Don't allow execution of arbitrary code
        if grep -qE '^\s*(eval|exec|source|\.)\s+' "$config" 2>/dev/null; then
            log_error "Configuration file contains unsafe commands"
            return 1
        fi
        set -a
        source "$config"
        set +a
        log_debug "Loaded configuration from: $config"
    fi
    return 0
}

#-------------------------------------------------------------------------------
# Platform Detection
#-------------------------------------------------------------------------------

# Detect platform
detect_platform() {
    local platform="unknown"
    
    case "$(uname -s)" in
        Linux*)
            platform="linux"
            ;;
        Darwin*)
            platform="macos"
            ;;
        *)
            platform="unknown"
            ;;
    esac
    
    echo "$platform"
}

# Check if running on supported platform
is_supported_platform() {
    local platform=$(detect_platform)
    [[ "$platform" == "linux" || "$platform" == "macos" ]]
}

#-------------------------------------------------------------------------------
# Dependency Checks
#-------------------------------------------------------------------------------

# Check for required commands
check_dependencies() {
    local missing=()
    local deps=(mkdir cp rm cat grep sed date)
    
    for cmd in "${deps[@]}"; do
        if ! command -v "$cmd" &>/dev/null; then
            missing+=("$cmd")
        fi
    done
    
    if [[ ${#missing[@]} -gt 0 ]]; then
        log_error "Missing required commands: ${missing[*]}"
        return 1
    fi
    
    return 0
}

#-------------------------------------------------------------------------------
# Validation
#-------------------------------------------------------------------------------

# Validate installation path
validate_install_path() {
    local path=$(get_install_path)
    
    # Check if path exists and is a directory
    if [[ -e "$path" && ! -d "$path" ]]; then
        log_error "Installation path exists but is not a directory: $path"
        return 1
    fi
    
    # Check write permissions
    if [[ -e "$path" && ! -w "$path" ]]; then
        log_error "Installation path is not writable: $path"
        return 1
    fi
    
    return 0
}

#-------------------------------------------------------------------------------
# Summary Output
#-------------------------------------------------------------------------------

# Print installation summary header
print_summary_header() {
    echo "============================================================"
    echo " KDSE Runtime - $1"
    echo "============================================================"
}

# Print summary footer
print_summary_footer() {
    echo "============================================================"
}

# Print status line
print_status() {
    local status="$1"
    local message="$2"
    
    case "$status" in
        PASS|OK|SUCCESS)
            log_success "$message"
            ;;
        FAIL|ERROR)
            log_error "$message"
            ;;
        WARN|WARNING)
            log_warning "$message"
            ;;
        *)
            echo "[$status] $message"
            ;;
    esac
}

#-------------------------------------------------------------------------------
# Exit Codes
#-------------------------------------------------------------------------------

EXIT_OK=0
EXIT_ERROR=1
EXIT_ALREADY_INSTALLED=2
EXIT_NOT_INSTALLED=3
EXIT_INVALID_ARGS=4
EXIT_MISSING_DEPS=5

#-------------------------------------------------------------------------------
# Script Information
#-------------------------------------------------------------------------------

# Get script name
get_script_name() {
    basename "$0"
}

# Get script directory
get_script_dir() {
    dirname "$(readlink -f "$0" 2>/dev/null || echo "$0")"
}

#-------------------------------------------------------------------------------
# Export Functions for Use in Other Scripts
#-------------------------------------------------------------------------------

export -f log_info log_success log_warning log_error log_debug
export -f get_install_path get_manifest_path get_config_path
export -f get_reports_path get_history_path get_runtime_path get_cache_path
export -f ensure_directory remove_directory copy_file remove_file
export -f calculate_checksum manifest_exists manifest_get load_config
export -f detect_platform is_supported_platform check_dependencies
export -f validate_install_path print_summary_header print_summary_footer print_status
export -f get_timestamp get_git_version get_script_name get_script_dir
export -f json_escape
