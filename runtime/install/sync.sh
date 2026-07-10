#!/usr/bin/env bash
#===============================================================================
# KDSE Runtime Synchronization Script
#===============================================================================
# Purpose: Synchronize installed KDSE standards with repository
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

Synchronize KDSE Runtime with latest repository standards.

OPTIONS:
    -h, --help          Show this help message
    -v, --verbose       Verbose output
    -f, --force         Force sync even if up-to-date
    -r, --repo URL      KDSE repository URL (default: from config)
    -b, --branch NAME   KDSE branch (default: from config)

PRESERVES:
    reports/            - Session reports
    history/            - Session history
    runtime/            - Runtime state
    cache/              - Cached data
    configuration/      - User configuration

EXAMPLES:
    $(get_script_name)              # Sync with latest
    $(get_script_name) -v          # Verbose sync
    $(get_script_name) -f          # Force sync

EOF
    exit $EXIT_OK
}

#-------------------------------------------------------------------------------
# Parse Arguments
#-------------------------------------------------------------------------------

FORCE_SYNC=0

while [[ $# -gt 0 ]]; do
    case "$1" in
        -h|--help)
            usage
            ;;
        -v|--verbose)
            VERBOSE=1
            shift
            ;;
        -f|--force)
            FORCE_SYNC=1
            shift
            ;;
        -r|--repo)
            KDSE_REPO="$2"
            shift 2
            ;;
        -b|--branch)
            KDSE_BRANCH="$2"
            shift 2
            ;;
        *)
            log_error "Unknown option: $1"
            usage
            ;;
    esac
done

#-------------------------------------------------------------------------------
# Pre-Sync Checks
#-------------------------------------------------------------------------------

pre_sync_checks() {
    log_info "Running pre-synchronization checks..."
    
    # Check if installed
    if ! manifest_exists; then
        log_error "KDSE Runtime is not installed"
        log_info "Run install.sh first"
        exit $EXIT_NOT_INSTALLED
    fi
    
    # Load existing configuration
    load_config
    
    # Check dependencies
    check_dependencies || {
        log_error "Missing required dependencies"
        exit $EXIT_MISSING_DEPS
    }
    
    log_info "Pre-synchronization checks passed"
    return 0
}

#-------------------------------------------------------------------------------
# Fetch Latest Changes
#-------------------------------------------------------------------------------

fetch_latest() {
    log_info "Fetching latest changes..."
    
    local install_path=$(get_install_path)
    
    if [[ ! -d "${install_path}/.git" ]]; then
        log_error "Not a git repository: $install_path"
        return 1
    fi
    
    (cd "$install_path" && git fetch origin "$KDSE_BRANCH" 2>/dev/null) || {
        log_warning "Could not fetch from origin"
    }
    
    log_info "Latest changes fetched"
    return 0
}

#-------------------------------------------------------------------------------
# Check for Updates
#-------------------------------------------------------------------------------

check_for_updates() {
    log_info "Checking for updates..."
    
    local install_path=$(get_install_path)
    local current_branch=$(manifest_get "branch" 2>/dev/null || echo "$KDSE_BRANCH")
    
    # Get current and latest commit
    local current_commit=$(cd "$install_path" && git rev-parse HEAD 2>/dev/null | head -c 8 || echo "unknown")
    local latest_commit=$(cd "$install_path" && git rev-parse "origin/${current_branch}" 2>/dev/null | head -c 8 || echo "unknown")
    
    log_info "Current version: $current_commit"
    log_info "Latest version:  $latest_commit"
    
    if [[ "$current_commit" == "$latest_commit" ]]; then
        if [[ "$FORCE_SYNC" == "0" ]]; then
            log_info "Already up-to-date"
            return 1
        else
            log_info "Force sync requested"
        fi
    fi
    
    return 0
}

#-------------------------------------------------------------------------------
# Backup User Data
#-------------------------------------------------------------------------------

backup_user_data() {
    log_info "Backing up user data..."
    
    local install_path=$(get_install_path)
    
    # Directories to preserve
    local preserve=(
        "reports"
        "history"
        "runtime"
        "cache"
        "configuration"
    )
    
    # Create backup directory
    local backup_dir="${install_path}/.backup"
    ensure_directory "$backup_dir" || return 1
    
    for dir in "${preserve[@]}"; do
        local src="${install_path}/${dir}"
        local dest="${backup_dir}/${dir}"
        
        if [[ -d "$src" ]]; then
            cp -r "$src" "$dest" 2>/dev/null || {
                log_warning "Could not backup: $dir"
            }
        fi
    done
    
    # Also backup manifest and config
    [[ -f "$(get_manifest_path)" ]] && cp "$(get_manifest_path)" "${backup_dir}/" 2>/dev/null || true
    [[ -f "$(get_config_path)" ]] && cp "$(get_config_path)" "${backup_dir}/" 2>/dev/null || true
    
    log_info "User data backed up to: $backup_dir"
    return 0
}

#-------------------------------------------------------------------------------
# Synchronize Documents
#-------------------------------------------------------------------------------

synchronize_documents() {
    log_info "Synchronizing normative documents..."
    
    local install_path=$(get_install_path)
    local standards_dir="${install_path}/standards/normative"
    
    local updated=0
    local added=0
    local failed=0
    
    # Update normative documents
    for doc in "${NORMATIVE_DOCS[@]}"; do
        local src="${install_path}/${doc}"
        local dest="${standards_dir}/$(basename "$doc")"
        
        if [[ -f "$src" ]]; then
            if [[ -f "$dest" ]]; then
                # Compare checksums
                local src_sum=$(calculate_checksum "$src")
                local dest_sum=$(calculate_checksum "$dest")
                
                if [[ "$src_sum" != "$dest_sum" ]]; then
                    cp -p "$src" "$dest" && ((updated++)) || ((failed++))
                fi
            else
                cp -p "$src" "$dest" && ((added++)) || ((failed++))
            fi
        fi
    done
    
    # Sync informative documents
    log_info "Synchronizing informative documents..."
    local info_dir="${install_path}/standards/informative"
    ensure_directory "$info_dir" || return 1
    
    for doc in "${INFORMATIVE_DOCS[@]}"; do
        local src="${install_path}/${doc}"
        local dest="${info_dir}/$(basename "$doc")"
        
        if [[ -f "$src" ]]; then
            if [[ -f "$dest" ]]; then
                local src_sum=$(calculate_checksum "$src")
                local dest_sum=$(calculate_checksum "$dest")
                
                if [[ "$src_sum" != "$dest_sum" ]]; then
                    cp -p "$src" "$dest" && ((updated++)) || ((failed++))
                fi
            else
                cp -p "$src" "$dest" && ((added++)) || ((failed++))
            fi
        fi
    done
    
    log_info "Documents: $updated updated, $added added, $failed failed"
    return 0
}

#-------------------------------------------------------------------------------
# Restore User Data
#-------------------------------------------------------------------------------

restore_user_data() {
    log_info "Restoring user data..."
    
    local install_path=$(get_install_path)
    local backup_dir="${install_path}/.backup"
    
    if [[ ! -d "$backup_dir" ]]; then
        log_info "No backup found, skipping restore"
        return 0
    fi
    
    # Restore preserved directories
    local preserve=(
        "reports"
        "history"
        "runtime"
        "cache"
        "configuration"
    )
    
    for dir in "${preserve[@]}"; do
        local src="${backup_dir}/${dir}"
        local dest="${install_path}/${dir}"
        
        if [[ -d "$src" ]]; then
            # Create destination if needed
            ensure_directory "$dest" || continue
            
            # Restore contents
            cp -r "$src/"* "$dest/" 2>/dev/null || {
                log_warning "Could not restore: $dir"
            }
        fi
    done
    
    # Clean up backup
    rm -rf "$backup_dir"
    
    log_info "User data restored"
    return 0
}

#-------------------------------------------------------------------------------
# Update Manifest
#-------------------------------------------------------------------------------

update_manifest() {
    log_info "Updating manifest..."
    
    local manifest=$(get_manifest_path)
    local version=$(get_git_version)
    local timestamp=$(get_timestamp)
    local current_repo=$(manifest_get "repo" 2>/dev/null || echo "$KDSE_REPO")
    local current_branch=$(manifest_get "branch" 2>/dev/null || echo "$KDSE_BRANCH")
    
    # Preserve original installation date
    local installed=$(manifest_get "installed" 2>/dev/null || echo "$timestamp")
    
    cat > "$manifest" <<EOF
{
    "kdse": "runtime-manifest",
    "version": "$version",
    "installed": "$installed",
    "last_sync": "$timestamp",
    "repo": "${KDSE_REPO:-$current_repo}",
    "branch": "${KDSE_BRANCH:-$current_branch}",
    "platform": "$(detect_platform)",
    "directories": {
        "reports": "$(get_reports_path)",
        "history": "$(get_history_path)",
        "runtime": "$(get_runtime_path)",
        "cache": "$(get_cache_path)",
        "configuration": "$(get_install_path)/configuration"
    },
    "normative_documents": [
$(for i in "${!NORMATIVE_DOCS[@]}"; do
    local doc="${NORMATIVE_DOCS[$i]}"
    local comma=","
    [[ $i -eq $((${#NORMATIVE_DOCS[@]} - 1)) ]] && comma=""
    echo "        \"$doc\"$comma"
done)
    ],
    "informative_documents": [
$(for i in "${!INFORMATIVE_DOCS[@]}"; do
    local doc="${INFORMATIVE_DOCS[$i]}"
    local comma=","
    [[ $i -eq $((${#INFORMATIVE_DOCS[@]} - 1)) ]] && comma=""
    echo "        \"$doc\"$comma"
done)
    ]
}
EOF
    
    log_info "Manifest updated"
    return 0
}

#-------------------------------------------------------------------------------
# Verify Integrity
#-------------------------------------------------------------------------------

verify_integrity() {
    log_info "Verifying integrity..."
    
    local install_path=$(get_install_path)
    local errors=0
    
    # Check manifest
    if [[ ! -f "$(get_manifest_path)" ]]; then
        log_error "Manifest not found"
        ((errors++))
    fi
    
    # Check required directories
    for dir in "${REQUIRED_DIRS[@]}"; do
        if [[ ! -d "${install_path}/${dir}" ]]; then
            log_error "Directory missing: $dir"
            ((errors++))
        fi
    done
    
    # Verify preserved directories have content
    if [[ -d "$(get_reports_path)" ]] && [[ -z "$(ls -A "$(get_reports_path)" 2>/dev/null)" ]]; then
        log_warning "Reports directory is empty (may be expected)"
    fi
    
    if [[ $errors -gt 0 ]]; then
        log_error "Integrity verification failed with $errors errors"
        return 1
    fi
    
    log_info "Integrity verified"
    return 0
}

#-------------------------------------------------------------------------------
# Print Sync Summary
#-------------------------------------------------------------------------------

print_summary() {
    local install_path=$(get_install_path)
    local old_version=$(cd "$install_path" && git rev-parse HEAD^ --short 2>/dev/null | head -c 8 || echo "unknown")
    local new_version=$(get_git_version)
    
    print_summary_header "Synchronization Complete"
    
    echo ""
    echo "Installation Path: $install_path"
    echo "Previous Version:  $old_version"
    echo "Current Version:   $new_version"
    echo "Synchronized:      $(get_timestamp)"
    echo ""
    echo "Preserved Data:"
    echo "  reports/       - Session reports"
    echo "  history/       - Session history"
    echo "  runtime/       - Runtime state"
    echo "  cache/         - Cached data"
    echo "  configuration/ - User configuration"
    echo ""
    
    print_summary_footer
}

#-------------------------------------------------------------------------------
# Main
#-------------------------------------------------------------------------------

main() {
    print_summary_header "KDSE Runtime Synchronization"
    
    # Run sync steps
    pre_sync_checks || exit $?
    fetch_latest || exit $EXIT_ERROR
    
    # Check for updates
    if ! check_for_updates; then
        if [[ "$FORCE_SYNC" == "0" ]]; then
            print_summary_header "Already Up-to-Date"
            echo ""
            echo "KDSE Runtime is already synchronized."
            echo ""
            print_summary_footer
            exit $EXIT_OK
        fi
    fi
    
    echo ""
    backup_user_data || exit $EXIT_ERROR
    synchronize_documents || exit $EXIT_ERROR
    restore_user_data || exit $EXIT_ERROR
    update_manifest || exit $EXIT_ERROR
    
    # Verify and report
    if verify_integrity; then
        print_summary
        exit $EXIT_OK
    else
        log_error "Synchronization verification failed"
        exit $EXIT_ERROR
    fi
}

# Run
main "$@"
