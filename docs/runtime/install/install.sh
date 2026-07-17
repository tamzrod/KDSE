#!/usr/bin/env bash
#===============================================================================
# KDSE Runtime Installation Script
#===============================================================================
# Purpose: Initialize KDSE Runtime in a repository
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

Initialize KDSE Runtime installation.

OPTIONS:
    -h, --help          Show this help message
    -v, --verbose       Verbose output
    -f, --force         Force reinstall (overwrites existing)
    -r, --repo URL      KDSE repository URL (default: ${KDSE_REPO})
    -b, --branch NAME   KDSE branch (default: ${KDSE_BRANCH})
    -p, --path PATH     Installation path (default: ${KDSE_HOME}/${KDSE_DIR})

EXAMPLES:
    $(get_script_name)                    # Interactive install
    $(get_script_name) -v                 # Verbose install
    $(get_script_name) -f                 # Force reinstall
    $(get_script_name) -r https://github.com/user/kdse.git

EOF
    exit $EXIT_OK
}

#-------------------------------------------------------------------------------
# Parse Arguments
#-------------------------------------------------------------------------------

FORCE_INSTALL=0

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
            FORCE_INSTALL=1
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
        -p|--path)
            KDSE_HOME="$2"
            shift 2
            ;;
        *)
            log_error "Unknown option: $1"
            usage
            ;;
    esac
done

#-------------------------------------------------------------------------------
# Pre-Installation Checks
#-------------------------------------------------------------------------------

pre_install_checks() {
    log_info "Running pre-installation checks..."
    
    # Check dependencies
    check_dependencies || {
        log_error "Missing required dependencies"
        exit $EXIT_MISSING_DEPS
    }
    
    # Check platform
    if ! is_supported_platform; then
        log_warning "Platform may not be fully supported: $(detect_platform)"
    fi
    
    # Validate install path
    validate_install_path || exit $EXIT_ERROR
    
    # Auto-detect and migrate legacy installations
    auto_detect_and_migrate || {
        log_error "Failed to process existing installation"
        exit $EXIT_ERROR
    }
    
    # Check if already installed
    if manifest_exists && [[ "$FORCE_INSTALL" == "0" ]]; then
        log_warning "KDSE Runtime is already installed"
        log_info "Use -f/--force to reinstall"
        exit $EXIT_ALREADY_INSTALLED
    fi
    
    log_info "Pre-installation checks passed"
    return 0
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

#-------------------------------------------------------------------------------
# Create Directory Structure
#-------------------------------------------------------------------------------

create_directories() {
    log_info "Creating directory structure..."
    
    local install_path=$(get_install_path)
    
    for dir in "${REQUIRED_DIRS[@]}"; do
        ensure_directory "${install_path}/${dir}" || {
            log_error "Failed to create directory: $dir"
            return 1
        }
    done
    
    log_info "Directory structure created"
    return 0
}

#-------------------------------------------------------------------------------
# Clone/Sync KDSE Repository
#-------------------------------------------------------------------------------

install_kdse_repository() {
    log_info "Installing KDSE repository..."
    
    local install_path=$(get_install_path)
    local git_dir="${install_path}/.git"
    
    if [[ -d "$git_dir" ]]; then
        log_info "KDSE repository already present, using existing"
        (cd "$install_path" && git fetch origin "$KDSE_BRANCH" 2>/dev/null || true)
    elif [[ -d "$install_path" ]]; then
        # Directory exists but no git - remove and re-clone
        log_info "Cleaning existing directory and re-cloning..."
        rm -rf "$install_path" || {
            log_error "Failed to clean existing directory"
            return 1
        }
        git clone --depth 1 --branch "$KDSE_BRANCH" "$KDSE_REPO" "$install_path" || {
            log_error "Failed to clone repository"
            return 1
        }
    else
        log_info "Cloning KDSE repository: $KDSE_REPO"
        git clone --depth 1 --branch "$KDSE_BRANCH" "$KDSE_REPO" "$install_path" || {
            log_error "Failed to clone repository"
            return 1
        }
    fi
    
    log_info "KDSE repository installed"
    return 0
}

#-------------------------------------------------------------------------------
# Install Normative Documents
#-------------------------------------------------------------------------------

install_normative_documents() {
    log_info "Installing normative documents..."
    
    local install_path=$(get_install_path)
    local standards_dir="${install_path}/standards/normative"
    
    ensure_directory "$standards_dir" || return 1
    
    local installed=0
    local failed=0
    
    for doc in "${NORMATIVE_DOCS[@]}"; do
        local src="${install_path}/${doc}"
        local dest="${standards_dir}/$(basename "$doc")"
        
        if [[ -f "$src" ]]; then
            copy_file "$src" "$dest" && ((installed++)) || ((failed++))
        else
            log_warning "Normative document not found: $doc"
            ((failed++))
        fi
    done
    
    log_info "Installed $installed normative documents, $failed failed"
    
    [[ $failed -eq 0 ]] || log_warning "Some normative documents could not be installed"
    
    return 0
}

#-------------------------------------------------------------------------------
# Install Informative Documents
#-------------------------------------------------------------------------------

install_informative_documents() {
    log_info "Installing informative documents..."
    
    local install_path=$(get_install_path)
    local info_dir="${install_path}/standards/informative"
    
    ensure_directory "$info_dir" || return 1
    
    local installed=0
    local failed=0
    
    for doc in "${INFORMATIVE_DOCS[@]}"; do
        local src="${install_path}/${doc}"
        local dest="${info_dir}/$(basename "$doc")"
        
        if [[ -f "$src" ]]; then
            copy_file "$src" "$dest" && ((installed++)) || ((failed++))
        else
            log_debug "Informative document not found: $doc"
        fi
    done
    
    log_info "Installed $installed informative documents"
    
    return 0
}

#-------------------------------------------------------------------------------
# Generate Manifest
#-------------------------------------------------------------------------------

generate_manifest() {
    log_info "Generating manifest..."
    
    local manifest=$(get_manifest_path)
    local version=$(get_git_version)
    local timestamp=$(get_timestamp)
    
    # Create manifest
    cat > "$manifest" <<EOF
{
    "kdse": "runtime-manifest",
    "version": "$version",
    "installed": "$timestamp",
    "repo": "$KDSE_REPO",
    "branch": "$KDSE_BRANCH",
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
    
    log_info "Manifest generated"
    return 0
}

#-------------------------------------------------------------------------------
# Generate Configuration
#-------------------------------------------------------------------------------

generate_configuration() {
    log_info "Generating configuration..."
    
    local config=$(get_config_path)
    local install_path=$(get_install_path)
    
    cat > "$config" <<EOF
#!/usr/bin/env bash
#===============================================================================
# KDSE Runtime Configuration
#===============================================================================
# This file is auto-generated. Do not edit manually.
#===============================================================================

# Installation
export KDSE_HOME="${KDSE_HOME}"
export KDSE_DIR="${KDSE_DIR}"
export KDSE_INSTALL_PATH="${install_path}"

# Repository
export KDSE_REPO="${KDSE_REPO}"
export KDSE_BRANCH="${KDSE_BRANCH}"

# Paths
export KDSE_REPORTS_DIR="$(get_reports_path)"
export KDSE_HISTORY_DIR="$(get_history_path)"
export KDSE_RUNTIME_DIR="$(get_runtime_path)"
export KDSE_CACHE_DIR="$(get_cache_path)"
export KDSE_CONFIG_DIR="${install_path}/configuration"

# Version
export KDSE_VERSION="$(get_git_version)"
export KDSE_INSTALLED="$(get_timestamp)"

# Platform
export KDSE_PLATFORM="$(detect_platform)"

# Command Path - Add KDSE to PATH
export PATH="${install_path}:\$PATH"

# Aliases for common commands
alias kdse-status='${install_path}/kdse status'
alias kdse-update='${install_path}/kdse update'
alias kdse-verify='${install_path}/kdse verify'
alias kdse-doctor='${install_path}/kdse doctor'
alias kdse-run='${install_path}/kdse run'
alias kdse-resume='${install_path}/kdse resume'

EOF
    
    chmod +x "$config"
    log_info "Configuration generated"
    return 0
}

#-------------------------------------------------------------------------------
# Install KDSE Command
#-------------------------------------------------------------------------------

install_kdse_command() {
    log_info "Installing KDSE command interface..."
    
    local install_path=$(get_install_path)
    local kdse_script="${install_path}/kdse"
    local source_script="${SCRIPT_DIR}/kdse"
    
    # Copy the kdse script
    if [[ -f "$source_script" ]]; then
        cp "$source_script" "$kdse_script" || {
            log_error "Failed to copy kdse script"
            return 1
        }
        chmod +x "$kdse_script"
        log_info "KDSE command installed to: $kdse_script"
    else
        log_warning "KDSE command script not found in repository"
        return 1
    fi
    
    # Create convenience symlink in PATH if possible
    local bin_dir="${KDSE_HOME}/bin"
    if [[ -d "$bin_dir" || -w "$(dirname "$bin_dir")" ]]; then
        ensure_directory "$bin_dir" || true
        local symlink="${bin_dir}/kdse"
        if [[ ! -L "$symlink" && ! -f "$symlink" ]]; then
            ln -sf "$kdse_script" "$symlink" 2>/dev/null && {
                log_info "KDSE command available as: $symlink"
            } || true
        fi
    fi
    
    return 0
}

#-------------------------------------------------------------------------------
# Install Command Registry
#-------------------------------------------------------------------------------

install_command_registry() {
    log_info "Installing command registry..."
    
    local install_path=$(get_install_path)
    local runtime_dir="${install_path}/runtime"
    local registry_dest="${runtime_dir}/commands.yaml"
    local registry_src="${SCRIPT_DIR}/commands.yaml"
    
    # Ensure runtime directory exists
    ensure_directory "$runtime_dir" || {
        log_error "Failed to create runtime directory"
        return 1
    }
    
    # Copy the command registry
    if [[ -f "$registry_src" ]]; then
        cp "$registry_src" "$registry_dest" || {
            log_error "Failed to copy command registry"
            return 1
        }
        log_info "Command registry installed to: $registry_dest"
    else
        log_warning "Command registry not found in repository"
        # Create a minimal registry
        cat > "$registry_dest" <<'REGISTRY_EOF'
# KDSE Runtime Command Registry
# Minimal registry - full registry should be in repository
version: "1.0"
interface_version: "1.0"
commands: []
REGISTRY_EOF
        log_info "Created minimal command registry"
    fi
    
    return 0
}

#-------------------------------------------------------------------------------
# Verify Installation
#-------------------------------------------------------------------------------

verify_installation() {
    log_info "Verifying installation..."
    
    local install_path=$(get_install_path)
    local errors=0
    
    # Check manifest
    if [[ ! -f "$(get_manifest_path)" ]]; then
        log_error "Manifest not found"
        ((errors++))
    fi
    
    # Check configuration
    if [[ ! -f "$(get_config_path)" ]]; then
        log_error "Configuration not found"
        ((errors++))
    fi
    
    # Check required directories
    for dir in "${REQUIRED_DIRS[@]}"; do
        if [[ ! -d "${install_path}/${dir}" ]]; then
            log_error "Directory missing: $dir"
            ((errors++))
        fi
    done
    
    # Check normative documents
    for doc in "${NORMATIVE_DOCS[@]}"; do
        if [[ ! -f "${install_path}/standards/normative/$(basename "$doc")" ]]; then
            log_warning "Normative document missing: $(basename "$doc")"
        fi
    done
    
    if [[ $errors -gt 0 ]]; then
        log_error "Verification failed with $errors errors"
        return 1
    fi
    
    log_info "Installation verified successfully"
    return 0
}

#-------------------------------------------------------------------------------
# Print Installation Summary
#-------------------------------------------------------------------------------

print_summary() {
    local install_path=$(get_install_path)
    local manifest_format=$(detect_manifest_format)
    
    print_summary_header "Installation Complete"
    
    echo ""
    echo "Installation Path: $install_path"
    echo "KDSE Repository:   $KDSE_REPO"
    echo "Branch:           $KDSE_BRANCH"
    echo "Version:          $(get_git_version)"
    echo "Installed:        $(get_timestamp)"
    echo ""
    echo "Installation Format: $manifest_format"
    
    # Show migration info if applicable
    if [[ "$manifest_format" == "yaml" ]]; then
        echo ""
        echo "Migration:        YES (YAML -> JSON)"
        echo "Legacy manifest:  manifest.yaml (backed up)"
    fi
    
    echo ""
    echo "Directory Structure:"
    echo "  .kdse/"
    for dir in "${REQUIRED_DIRS[@]}"; do
        echo "    $dir/"
    done
    echo "  manifest.json"
    echo "  config.sh"
    echo "  kdse"
    echo "  runtime/commands.yaml  # Command registry (AI agent integration)"
    echo ""
    echo "Quick Start:"
    echo "  ${install_path}/kdse status      # Check runtime health"
    echo "  ${install_path}/kdse commands    # List all commands"
    echo "  ${install_path}/kdse update      # Update runtime"
    echo "  ${install_path}/kdse run         # Start session"
    echo ""
    echo "  Add to PATH for convenience:"
    echo "  export PATH=\"${install_path}:\$PATH\""
    echo ""
    
    print_summary_footer
}

#-------------------------------------------------------------------------------
# Main
#-------------------------------------------------------------------------------

main() {
    print_summary_header "KDSE Runtime Installation"
    
    echo ""
    echo "Repository: $KDSE_REPO"
    echo "Branch:     $KDSE_BRANCH"
    echo "Path:       $(get_install_path)"
    echo ""
    
    # Run installation steps
    pre_install_checks || exit $?
    install_kdse_repository || exit $EXIT_ERROR
    create_directories || exit $EXIT_ERROR
    install_normative_documents || exit $EXIT_ERROR
    install_informative_documents || exit $EXIT_ERROR
    generate_manifest || exit $EXIT_ERROR
    generate_configuration || exit $EXIT_ERROR
    install_kdse_command || exit $EXIT_ERROR
    install_command_registry || exit $EXIT_ERROR
    
    # Verify and report
    if verify_installation; then
        print_summary
        exit $EXIT_OK
    else
        log_error "Installation verification failed"
        exit $EXIT_ERROR
    fi
}

# Run
main "$@"
