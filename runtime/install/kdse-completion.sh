#!/usr/bin/env bash
#===============================================================================
# KDSE Runtime Command Completion
#===============================================================================
# Purpose: Bash completion for KDSE Runtime commands
# Usage:   source this file in your .bashrc or .zshrc
#===============================================================================

# Command list
_kdse_commands="status version update verify doctor install uninstall history report run resume audit help"

# Option list
_kdse_options="-h --help -v --verbose -q --quiet -f --force"

#-------------------------------------------------------------------------------
# Completion Function
#-------------------------------------------------------------------------------

_kdse_complete() {
    local cur prev words cword
    
    # Initialize bash completion variables
    _init_completion || return
    
    # Current word being completed
    cur="${words[$cword]}"
    
    # Previous word
    prev="${words[$((cword - 1))]}"
    
    # Check if we're completing a command or option
    if [[ $cword -eq 1 ]]; then
        # Completing the command itself
        COMPREPLY=($(compgen -W "${_kdse_commands}" -- "$cur"))
    elif [[ "$prev" == "-"* ]]; then
        # Completing an option
        COMPREPLY=($(compgen -W "${_kdse_options}" -- "$cur"))
    else
        # Completing arguments for specific commands
        case "$prev" in
            install|uninstall)
                COMPREPLY=($(compgen -W "-f --force -v --verbose -h --help" -- "$cur"))
                ;;
            verify)
                COMPREPLY=($(compgen -W "-v --verbose -q --quiet -j --json -h --help" -- "$cur"))
                ;;
            status|version|doctor|history|report|run|resume|audit)
                COMPREPLY=($(compgen -W "-v --verbose -q --quiet -h --help" -- "$cur"))
                ;;
            update)
                COMPREPLY=($(compgen -W "-v --verbose -f --force -h --help" -- "$cur"))
                ;;
            *)
                COMPREPLY=()
                ;;
        esac
    fi
    
    # Handle file/directory completion for options that take arguments
    if [[ "$prev" == "-r" || "$prev" == "--repo" || "$prev" == "-b" || "$prev" == "--branch" || "$prev" == "-p" || "$prev" == "--path" ]]; then
        COMPREPLY=()
    fi
    
    return 0
}

#-------------------------------------------------------------------------------
# Register Completion
#-------------------------------------------------------------------------------

# Register the completion function
complete -F _kdse_complete kdse

# Also complete for full path
complete -F _kdse_complete "${HOME}/.kdse/kdse"

#-------------------------------------------------------------------------------
# Quick Install Instructions
#-------------------------------------------------------------------------------

# Add to ~/.bashrc or ~/.zshrc:
#
# For bash:
#   source /path/to/kdse-completion.sh
#
# For zsh:
#   autoload bashcompinit
#   bashcompinit
#   source /path/to/kdse-completion.sh
#
# Or copy to /etc/bash_completion.d/:
#   sudo cp kdse-completion.sh /etc/bash_completion.d/
#   source /etc/bash_completion.d/kdse-completion.sh
