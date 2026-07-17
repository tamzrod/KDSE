#!/usr/bin/env bash
#===============================================================================
# KDSE Runtime Uninstallation Script
#===============================================================================
# Purpose: Remove KDSE Runtime installation
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

Uninstall KDSE Runtime.

OPTIONS:
    -h, --help          Show this help message
    -v, --verbose       Verbose output
    -f, --force         Force uninstall without confirmation
    -k, --keep-reports  Keep reports directory (default: preserved anyway)
    -a, --keep-all      Keep all user data (reports, history, runtime, cache)

PRESERVES BY DEFAULT:
    reports/            - Session reports (unless -a specified)
    history/            - Session history (unless -a specified)

EXAMPLES:
    $(get_script_name)              # Interactive uninstall
    $(get_script_name) -f           # Force uninstall
    $(get_script_name) -a           # Keep all user data

EOF
    exit $EXIT_OK
}

#-------------------------------------------------------------------------------
# Global State
#-------------------------------------------------------------------------------

KEEP_REPORTS=1
KEEP_HISTORY=1
KEEP_RUNTIME=0
KEEP_CACHE=0
KEEP_CONFIGURATION=0

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
        -f|--force)
            FORCE_UNINSTALL=1
            shift
            ;;
        -k|--keep-reports)
            KEEP_REPORTS=1
            shift
            ;;
        -a|--keep-all)
            KEEP_REPORTS=1
            KEEP_HISTORY=1
            KEEP_RUNTIME=1
            KEEP_CACHE=1
            KEEP_CONFIGURATION=1
            shift
            ;;
        *)
            log_error "Unknown option: $1"
            usage
            ;;
    esac
done

#-------------------------------------------------------------------------------
# Confirmation
#-------------------------------------------------------------------------------

confirm_uninstall() {
    local install_path=$(get_install_path)
    
    if [[ "${FORCE_UNINSTALL:-0}" == "1" ]]; then
        return 0
    fi
    
    echo ""
    print_summary_header "KDSE Runtime Uninstallation"
    echo ""
    echo "Installation Path: $install_path"
    echo ""
    echo "This will remove:"
    echo "  - KDSE Runtime installation"
    echo "  - Normative documents"
    echo "  - Informative documents"
    echo "  - Cached data"
    echo ""
    echo "This will preserve:"
    [[ $KEEP_REPORTS -eq 1 ]] && echo "  - reports/"
    [[ $KEEP_HISTORY -eq 1 ]] && echo "  - history/"
    [[ $KEEP_RUNTIME -eq 1 ]] && echo "  - runtime/"
    [[ $KEEP_CACHE -eq 1 ]] && echo "  - cache/"
    [[ $KEEP_CONFIGURATION -eq 1 ]] && echo "  - configuration/"
    echo ""
    echo -n "Continue? [y/N] "
    
    local response
    read response || response="n"
    
    case "$response" in
        [yY]|[yY][eE][sS])
            return 0
            ;;
        *)
            log_info "Uninstallation cancelled"
            exit $EXIT_OK
            ;;
    esac
}

#-------------------------------------------------------------------------------
# Pre-Uninstall Checks
#-------------------------------------------------------------------------------

pre_uninstall_checks() {
    log_info "Running pre-uninstallation checks..."
    
    # Check if installed
    if ! manifest_exists; then
        log_warning "KDSE Runtime is not installed"
        exit $EXIT_NOT_INSTALLED
    fi
    
    # Validate install path
    validate_install_path || exit $EXIT_ERROR
    
    log_info "Pre-uninstallation checks passed"
    return 0
}

#-------------------------------------------------------------------------------
# Backup User Data
#-------------------------------------------------------------------------------

backup_before_uninstall() {
    log_info "Backing up user data..."
    
    local install_path=$(get_install_path)
    
    # Determine what to preserve
    local preserve=()
    [[ $KEEP_REPORTS -eq 1 ]] && preserve+=("reports")
    [[ $KEEP_HISTORY -eq 1 ]] && preserve+=("history")
    [[ $KEEP_RUNTIME -eq 1 ]] && preserve+=("runtime")
    [[ $KEEP_CACHE -eq 1 ]] && preserve+=("cache")
    [[ $KEEP_CONFIGURATION -eq 1 ]] && preserve+=("configuration")
    
    if [[ ${#preserve[@]} -eq 0 ]]; then
        log_info "No data to preserve"
        return 0
    fi
    
    # Create backup in home directory
    local backup_parent="${KDSE_HOME}/.kdse-backup"
    local backup_dir="${backup_parent}/$(date +%Y%m%d-%H%M%S)"
    
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
    
    log_info "User data backed up to: $backup_dir"
    log_info "Backup location: $backup_parent"
    
    return 0
}

#-------------------------------------------------------------------------------
# Remove Installation
#-------------------------------------------------------------------------------

remove_installation() {
    log_info "Removing KDSE Runtime installation..."
    
    local install_path=$(get_install_path)
    local errors=0
    
    # Remove in order of dependency
    local remove_patterns=(
        ".backup"
        ".git"
        "standards"
        "manifest.json"
        "manifest.yaml"
        "manifest.yaml.backup"
        "config.sh"
    )
    
    for pattern in "${remove_patterns[@]}"; do
        local target="${install_path}/${pattern}"
        
        if [[ -e "$target" ]]; then
            if [[ "$pattern" == ".git" ]]; then
                # Git directory needs special handling
                rm -rf "$target" 2>/dev/null || {
                    log_warning "Could not remove: $target"
                    ((errors++))
                }
            elif [[ "$pattern" == ".backup" ]]; then
                # Skip backup during removal (we already backed up)
                continue
            else
                rm -rf "$target" 2>/dev/null || {
                    log_warning "Could not remove: $target"
                    ((errors++))
                }
            fi
        fi
    done
    
    # Remove preserved directories if not keeping
    [[ $KEEP_REPORTS -eq 0 ]] && rm -rf "${install_path}/reports" 2>/dev/null || true
    [[ $KEEP_HISTORY -eq 0 ]] && rm -rf "${install_path}/history" 2>/dev/null || true
    [[ $KEEP_RUNTIME -eq 0 ]] && rm -rf "${install_path}/runtime" 2>/dev/null || true
    [[ $KEEP_CACHE -eq 0 ]] && rm -rf "${install_path}/cache" 2>/dev/null || true
    [[ $KEEP_CONFIGURATION -eq 0 ]] && rm -rf "${install_path}/configuration" 2>/dev/null || true
    
    # Remove empty installation directory
    if [[ -d "$install_path" ]]; then
        if [[ -z "$(ls -A "$install_path" 2>/dev/null)" ]]; then
            rmdir "$install_path" 2>/dev/null || {
                log_warning "Could not remove empty directory: $install_path"
                ((errors++))
            }
        else
            # Check what's left
            local remaining=$(ls -A "$install_path" 2>/dev/null || echo "")
            if [[ -n "$remaining" ]]; then
                log_warning "Non-empty installation directory: $install_path"
                log_warning "Remaining items: $remaining"
            fi
        fi
    fi
    
    if [[ $errors -gt 0 ]]; then
        log_warning "Uninstallation completed with $errors error(s)"
    fi
    
    return 0
}

#-------------------------------------------------------------------------------
# Clean Up Backup Directory
#-------------------------------------------------------------------------------

cleanup_old_backups() {
    log_info "Checking for old backups..."
    
    local backup_parent="${KDSE_HOME}/.kdse-backup"
    
    if [[ ! -d "$backup_parent" ]]; then
        return 0
    fi
    
    # Count backups
    local backup_count=$(find "$backup_parent" -maxdepth 1 -type d 2>/dev/null | wc -l)
    ((backup_count--)) # Subtract the parent directory itself
    
    if [[ $backup_count -gt 3 ]]; then
        log_info "Found $backup_count old backups, keeping 3 most recent"
        
        # Remove oldest backups beyond 3
        find "$backup_parent" -maxdepth 1 -type d -not -name "$(basename "$backup_parent")" \
            -not -newer "$(find "$backup_parent" -maxdepth 1 -type d -not -name "$(basename "$backup_parent")" | tail -3)" \
            2>/dev/null | xargs rm -rf 2>/dev/null || true
    fi
    
    return 0
}

#-------------------------------------------------------------------------------
# Print Summary
#-------------------------------------------------------------------------------

print_summary() {
    local install_path=$(get_install_path)
    
    print_summary_header "Uninstallation Complete"
    
    echo ""
    echo "Installation Path:  $install_path"
    echo "Removed:           $(get_timestamp)"
    echo ""
    
    echo "Preserved Data:"
    [[ $KEEP_REPORTS -eq 1 ]] && echo "  reports/       -> ${KDSE_HOME}/.kdse-backup/"
    [[ $KEEP_HISTORY -eq 1 ]] && echo "  history/       -> ${KDSE_HOME}/.kdse-backup/"
    [[ $KEEP_RUNTIME -eq 1 ]] && echo "  runtime/       -> ${KDSE_HOME}/.kdse-backup/"
    [[ $KEEP_CACHE -eq 1 ]] && echo "  cache/         -> ${KDSE_HOME}/.kdse-backup/"
    [[ $KEEP_CONFIGURATION -eq 1 ]] && echo "  configuration/ -> ${KDSE_HOME}/.kdse-backup/"
    echo ""
    
    echo "Removed:"
    echo "  KDSE Runtime"
    echo "  Normative documents"
    echo "  Informative documents"
    echo "  Git repository"
    echo ""
    
    if [[ -d "${KDSE_HOME}/.kdse-backup" ]]; then
        echo "Backup location: ${KDSE_HOME}/.kdse-backup/"
    fi
    echo ""
    
    print_summary_footer
}

#-------------------------------------------------------------------------------
# Main
#-------------------------------------------------------------------------------

main() {
    # Pre-uninstall checks
    pre_uninstall_checks || exit $?
    
    # Confirm uninstallation
    confirm_uninstall || exit $?
    
    # Backup user data
    backup_before_uninstall || exit $EXIT_ERROR
    
    # Remove installation
    remove_installation || exit $EXIT_ERROR
    
    # Clean up old backups
    cleanup_old_backups
    
    # Print summary
    print_summary
    
    exit $EXIT_OK
}

# Run
main "$@"
