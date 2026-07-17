#!/usr/bin/env bash
#===============================================================================
# KDSE Debug Engine
#===============================================================================
# Purpose: Core engine for evidence-driven debugging
# Design:  Deterministic, evidence-based, confidence-driven root cause analysis
#===============================================================================

set -euo pipefail

#-------------------------------------------------------------------------------
# Debug Engine Configuration
#-------------------------------------------------------------------------------

DEBUG_ENGINE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
KDSE_DEBUG_DIR="${KDSE_DEBUG_DIR:-.kdse/debug}"
KDSE_BOOTSTRAP_DIR="${KDSE_BOOTSTRAP_DIR:-.kdse/bootstrap}"

# Source configuration
CONFIG_FILE="${KDSE_BOOTSTRAP_DIR}/debug-config.yaml"

# Session state
SESSION_ID=""
SESSION_DIR=""
SESSION_STARTED=""
CURRENT_STATE="INITIAL"

# Evidence counter
EVIDENCE_COUNTER=0
HYPOTHESIS_COUNTER=0
EXPERIMENT_COUNTER=0

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

#-------------------------------------------------------------------------------
# Utility Functions
#-------------------------------------------------------------------------------

log_info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
    echo -e "${GREEN}[OK]${NC} $*"
}

log_warning() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*"
}

log_debug() {
    if [[ "${KDSE_DEBUG:-0}" == "1" ]]; then
        echo -e "${CYAN}[DEBUG]${NC} $*"
    fi
}

get_timestamp() {
    date -u +"%Y-%m-%dT%H:%M:%SZ"
}

get_timestamp_short() {
    date -u +"%Y%m%d%H%M%S"
}

generate_id() {
    local prefix="$1"
    local counter="$2"
    printf "%s-%s-%04d" "$prefix" "$(date -u +%Y%m%d)" "$counter"
}

#-------------------------------------------------------------------------------
# Configuration Loading
#-------------------------------------------------------------------------------

load_config() {
    if [[ -f "$CONFIG_FILE" ]]; then
        log_debug "Loading debug configuration from $CONFIG_FILE"
        # For YAML, we use grep/sed for simple values
        CONFIDENCE_THRESHOLD=$(grep "^  threshold:" "$CONFIG_FILE" | sed 's/.*: *//')
        CONFIDENCE_MIN_INITIAL=$(grep "^  minimum_initial:" "$CONFIG_FILE" | sed 's/.*: *//')
        CONFIDENCE_MAX_INITIAL=$(grep "^  maximum_initial:" "$CONFIG_FILE" | sed 's/.*: *//')
        CONFIDENCE_DEFAULT=$(grep "^  default_initial:" "$CONFIG_FILE" | sed 's/.*: *//')
        HYPOTHESIS_MAX_ACTIVE=$(grep "^  max_active:" "$CONFIG_FILE" | sed 's/.*: *//')
        LOOPS_ENABLED=$(grep "^  detection_enabled:" "$CONFIG_FILE" | sed 's/.*: *//')
        LOOPS_MAX_REPETITIONS=$(grep "^  max_repetitions:" "$CONFIG_FILE" | sed 's/.*: *//')
        
        # Defaults if not found
        CONFIDENCE_THRESHOLD=${CONFIDENCE_THRESHOLD:-90}
        CONFIDENCE_MIN_INITIAL=${CONFIDENCE_MIN_INITIAL:-20}
        CONFIDENCE_MAX_INITIAL=${CONFIDENCE_MAX_INITIAL:-60}
        CONFIDENCE_DEFAULT=${CONFIDENCE_DEFAULT:-40}
        HYPOTHESIS_MAX_ACTIVE=${HYPOTHESIS_MAX_ACTIVE:-5}
        LOOPS_ENABLED=${LOOPS_ENABLED:-true}
        LOOPS_MAX_REPETITIONS=${LOOPS_MAX_REPETITIONS:-3}
    else
        log_warning "Configuration file not found, using defaults"
        CONFIDENCE_THRESHOLD=90
        CONFIDENCE_MIN_INITIAL=20
        CONFIDENCE_MAX_INITIAL=60
        CONFIDENCE_DEFAULT=40
        HYPOTHESIS_MAX_ACTIVE=5
        LOOPS_ENABLED=true
        LOOPS_MAX_REPETITIONS=3
    fi
}

#-------------------------------------------------------------------------------
# Session Management
#-------------------------------------------------------------------------------

debug_init() {
    local failure_desc="${1:-}"
    local severity="${2:-medium}"
    
    load_config
    
    # Generate session ID
    EVIDENCE_COUNTER=1
    HYPOTHESIS_COUNTER=1
    EXPERIMENT_COUNTER=1
    SESSION_ID="DEBUG-$(get_timestamp_short)"
    SESSION_STARTED=$(get_timestamp)
    
    # Create session directory
    SESSION_DIR="${KDSE_DEBUG_DIR}/sessions/${SESSION_ID}"
    mkdir -p "$SESSION_DIR"
    mkdir -p "${KDSE_DEBUG_DIR}/evidence"
    mkdir -p "${KDSE_DEBUG_DIR}/hypotheses"
    mkdir -p "${KDSE_DEBUG_DIR}/reports"
    mkdir -p "${KDSE_DEBUG_DIR}/loops"
    
    # Initialize evidence store
    echo '{"evidence": []}' > "${KDSE_DEBUG_DIR}/evidence/store.json"
    
    # Initialize hypothesis registry
    echo '{"hypotheses": []}' > "${KDSE_DEBUG_DIR}/hypotheses/registry.json"
    
    # Initialize loop history
    echo '{"loops": []}' > "${KDSE_DEBUG_DIR}/loops/history.json"
    
    # Create session metadata
    cat > "${SESSION_DIR}/session.json" <<EOF
{
  "session_id": "${SESSION_ID}",
  "started_at": "${SESSION_STARTED}",
  "failure": {
    "description": "${failure_desc}",
    "severity": "${severity}"
  },
  "state": "INITIAL",
  "evidence_count": 0,
  "hypothesis_count": 0,
  "experiment_count": 0,
  "root_cause_selected": null
}
EOF
    
    CURRENT_STATE="EVIDENCE_COLLECTION"
    update_session_state
    
    echo ""
    echo -e "${BOLD}============================================================${NC}"
    echo -e "${BOLD} KDSE Debug Runtime${NC} - Debug Session Started"
    echo -e "${BOLD}============================================================${NC}"
    echo ""
    echo "Session ID: ${SESSION_ID}"
    echo "Started:    ${SESSION_STARTED}"
    echo "State:      ${CURRENT_STATE}"
    echo ""
    if [[ -n "$failure_desc" ]]; then
        echo "Failure:    ${failure_desc}"
        echo ""
    fi
    log_info "Evidence collection phase active"
    echo ""
}

update_session_state() {
    if [[ -z "$SESSION_DIR" ]]; then
        return 1
    fi
    
    # Update state in session.json
    if [[ -f "${SESSION_DIR}/session.json" ]]; then
        local temp_file="${SESSION_DIR}/session.json.tmp"
        sed "s/\"state\": \"[^\"]*\"/\"state\": \"${CURRENT_STATE}\"/" "${SESSION_DIR}/session.json" > "$temp_file"
        mv "$temp_file" "${SESSION_DIR}/session.json"
    fi
}

get_session_state() {
    if [[ -f "${SESSION_DIR}/session.json" ]]; then
        grep '"state"' "${SESSION_DIR}/session.json" | sed 's/.*: *"\([^"]*\)".*/\1/'
    else
        echo "UNKNOWN"
    fi
}

#-------------------------------------------------------------------------------
# Evidence Management
#-------------------------------------------------------------------------------

debug_collect_evidence() {
    local type="$1"
    local content="$2"
    local source="${3:-}"
    local tags="${4:-}"
    
    if [[ "$CURRENT_STATE" != "EVIDENCE_COLLECTION" && "$CURRENT_STATE" != "HYPOTHESIS_GENERATION" && "$CURRENT_STATE" != "INITIAL" ]]; then
        log_warning "Evidence collection not allowed in state: $CURRENT_STATE"
        return 1
    fi
    
    local evidence_id="E-$(printf "%04d" $EVIDENCE_COUNTER)"
    EVIDENCE_COUNTER=$((EVIDENCE_COUNTER + 1))
    
    local timestamp=$(get_timestamp)
    local evidence_json=""
    
    # Build evidence JSON
    evidence_json=$(cat <<EOF
{
  "id": "${evidence_id}",
  "type": "${type}",
  "timestamp": "${timestamp}",
  "source": "${source}",
  "content": $(echo "$content" | jq -Rs .),
  "tags": [$(echo "$tags" | jq -R 'split(",")' 2>/dev/null || echo "[]")]
}
EOF
)
    
    # Add to evidence store
    local store_file="${KDSE_DEBUG_DIR}/evidence/store.json"
    local temp_store="${store_file}.tmp"
    
    # Use jq to add evidence if available, otherwise use sed
    if command -v jq &> /dev/null; then
        echo "$(jq ".evidence += [${evidence_json}]" "$store_file")" > "$temp_store"
    else
        # Manual JSON manipulation
        local content_escaped=$(echo "$content" | sed 's/"/\\"/g' | tr '\n' ' ')
        evidence_json=$(cat <<EOF
{
  "id": "${evidence_id}",
  "type": "${type}",
  "timestamp": "${timestamp}",
  "source": "${source}",
  "content": "${content_escaped}",
  "tags": [${tags}]
}
EOF
)
        sed 's/\(.*"evidence": \[\)\(.*\)\(.*\)/\1\2,'"${evidence_json}"'\3/' "$store_file" > "$temp_store"
    fi
    mv "$temp_store" "$store_file"
    
    # Update session metadata
    local evidence_count=$(grep -c '"id":' "$store_file" 2>/dev/null || echo "0")
    sed -i "s/\"evidence_count\": [0-9]*/\"evidence_count\": ${evidence_count}/" "${SESSION_DIR}/session.json" 2>/dev/null || true
    
    echo ""
    echo -e "${GREEN}✓ Evidence Collected${NC}"
    echo "  ID:       ${evidence_id}"
    echo "  Type:     ${type}"
    echo "  Source:   ${source}"
    echo ""
    
    return 0
}

debug_list_evidence() {
    local store_file="${KDSE_DEBUG_DIR}/evidence/store.json"
    
    if [[ ! -f "$store_file" ]]; then
        log_error "Evidence store not found"
        return 1
    fi
    
    echo ""
    echo -e "${BOLD}Evidence Collected${NC}"
    echo "============================================================"
    
    if command -v jq &> /dev/null; then
        jq -r '.evidence[] | "[\(.id)] \(.type) - \(.source // "unknown")"' "$store_file" 2>/dev/null || {
            log_warning "No evidence collected yet"
        }
    else
        grep '"id":' "$store_file" | sed 's/.*"\(E-[^"]*\)".*/\1/'
    fi
    
    echo ""
}

#-------------------------------------------------------------------------------
# Hypothesis Management
#-------------------------------------------------------------------------------

debug_new_hypothesis() {
    local description="$1"
    local initial_confidence="${2:-${CONFIDENCE_DEFAULT}}"
    local components="${3:-}"
    
    if [[ "$CURRENT_STATE" != "HYPOTHESIS_GENERATION" && "$CURRENT_STATE" != "EVIDENCE_COLLECTION" ]]; then
        log_warning "Hypothesis creation not allowed in state: $CURRENT_STATE"
        return 1
    fi
    
    local hypothesis_id="H-$(printf "%04d" $HYPOTHESIS_COUNTER)"
    HYPOTHESIS_COUNTER=$((HYPOTHESIS_COUNTER + 1))
    
    # Check max active hypotheses
    local active_count=$(jq -r '.hypotheses | map(select(.status == "active")) | length' "${KDSE_DEBUG_DIR}/hypotheses/registry.json" 2>/dev/null || echo "0")
    if [[ "$active_count" -ge "$HYPOTHESIS_MAX_ACTIVE" ]]; then
        log_warning "Maximum active hypotheses (${HYPOTHESIS_MAX_ACTIVE}) reached"
        log_info "Rejecting weakest hypothesis..."
        debug_reject_weakest
    fi
    
    local timestamp=$(get_timestamp)
    local hypothesis_json=""
    
    # Build hypothesis JSON
    hypothesis_json=$(cat <<EOF
{
  "id": "${hypothesis_id}",
  "description": "${description}",
  "status": "active",
  "confidence": {
    "initial": ${initial_confidence},
    "current": ${initial_confidence},
    "threshold": ${CONFIDENCE_THRESHOLD}
  },
  "supporting_evidence": [],
  "contradicting_evidence": [],
  "affected_components": [$(echo "$components" | jq -R 'split(",")' 2>/dev/null || echo "[]")],
  "experiments": [],
  "created_at": "${timestamp}",
  "updated_at": "${timestamp}"
}
EOF
)
    
    # Add to hypothesis registry
    local registry_file="${KDSE_DEBUG_DIR}/hypotheses/registry.json"
    local temp_registry="${registry_file}.tmp"
    
    if command -v jq &> /dev/null; then
        echo "$(jq ".hypotheses += [${hypothesis_json}]" "$registry_file")" > "$temp_registry"
    else
        sed 's/\(.*"hypotheses": \[\)\(.*\)\(.*\)/\1\2,'"${hypothesis_json}"'\3/' "$registry_file" > "$temp_registry"
    fi
    mv "$temp_registry" "$registry_file"
    
    # Update session metadata
    local hypothesis_count=$(grep -c '"id": "H-' "$registry_file" 2>/dev/null || echo "0")
    sed -i "s/\"hypothesis_count\": [0-9]*/\"hypothesis_count\": ${hypothesis_count}/" "${SESSION_DIR}/session.json" 2>/dev/null || true
    
    echo ""
    echo -e "${GREEN}✓ Hypothesis Created${NC}"
    echo "  ID:            ${hypothesis_id}"
    echo "  Description:   ${description}"
    echo "  Confidence:    ${initial_confidence}%"
    echo ""
    
    return 0
}

debug_reject_weakest() {
    local registry_file="${KDSE_DEBUG_DIR}/hypotheses/registry.json"
    
    if [[ ! -f "$registry_file" ]]; then
        return 1
    fi
    
    # Find hypothesis with lowest confidence
    local weakest=$(jq -r '.hypotheses | map(select(.status == "active")) | min_by(.confidence.current) | .id' "$registry_file" 2>/dev/null)
    
    if [[ -n "$weakest" && "$weakest" != "null" ]]; then
        log_info "Rejecting weakest hypothesis: $weakest"
        debug_update_hypothesis_status "$weakest" "rejected"
    fi
}

debug_update_hypothesis_status() {
    local hypothesis_id="$1"
    local new_status="$2"
    local registry_file="${KDSE_DEBUG_DIR}/hypotheses/registry.json"
    
    if command -v jq &> /dev/null; then
        local temp_registry="${registry_file}.tmp"
        echo "$(jq "(.hypotheses[] | select(.id == \"${hypothesis_id}\") | .status) = \"${new_status}\" | (.hypotheses[] | select(.id == \"${hypothesis_id}\") | .updated_at) = \"$(get_timestamp)\"" "$registry_file")" > "$temp_registry"
        mv "$temp_registry" "$registry_file"
    fi
}

debug_list_hypotheses() {
    local registry_file="${KDSE_DEBUG_DIR}/hypotheses/registry.json"
    
    if [[ ! -f "$registry_file" ]]; then
        log_error "Hypothesis registry not found"
        return 1
    fi
    
    echo ""
    echo -e "${BOLD}Hypotheses${NC}"
    echo "============================================================"
    
    if command -v jq &> /dev/null; then
        jq -r '.hypotheses[] | "[\(.id)] \(.description) [\(if .status == "active" then "\(.confidence.current)%" else .status end)]"' "$registry_file" 2>/dev/null || {
            log_warning "No hypotheses created yet"
        }
    else
        grep '"id":' "$registry_file" | head -5
    fi
    
    echo ""
}

#-------------------------------------------------------------------------------
# Evidence Evaluation
#-------------------------------------------------------------------------------

debug_evaluate() {
    local hypothesis_id="$1"
    local evidence_id="$2"
    local impact="$3"  # "supporting" or "contradicting"
    
    if [[ "$CURRENT_STATE" != "EVIDENCE_EVALUATION" && "$CURRENT_STATE" != "HYPOTHESIS_GENERATION" ]]; then
        log_warning "Evidence evaluation not allowed in state: $CURRENT_STATE"
        return 1
    fi
    
    local registry_file="${KDSE_DEBUG_DIR}/hypotheses/registry.json"
    local store_file="${KDSE_DEBUG_DIR}/evidence/store.json"
    
    # Get evidence type for confidence impact
    local evidence_type=$(jq -r ".evidence[] | select(.id == \"${evidence_id}\") | .type" "$store_file" 2>/dev/null)
    
    # Calculate confidence impact
    local impact_value=0
    case "$evidence_type" in
        exception) impact_value=20 ;;
        test_failure) impact_value=15 ;;
        log) impact_value=10 ;;
        source) impact_value=10 ;;
        config) impact_value=5 ;;
        state) impact_value=5 ;;
        dependency) impact_value=5 ;;
        *) impact_value=5 ;;
    esac
    
    if [[ "$impact" == "contradicting" ]]; then
        impact_value=$((-25))
    fi
    
    # Update hypothesis
    if command -v jq &> /dev/null; then
        local temp_registry="${registry_file}.tmp"
        
        # Get current confidence
        local current_conf=$(jq -r ".hypotheses[] | select(.id == \"${hypothesis_id}\") | .confidence.current" "$registry_file")
        local new_conf=$((current_conf + impact_value))
        
        # Clamp to 0-100
        if [[ "$new_conf" -gt 100 ]]; then
            new_conf=100
        elif [[ "$new_conf" -lt 0 ]]; then
            new_conf=0
        fi
        
        if [[ "$impact" == "supporting" ]]; then
            jq "(.hypotheses[] | select(.id == \"${hypothesis_id}\") | .supporting_evidence) += [\"${evidence_id}\"] | (.hypotheses[] | select(.id == \"${hypothesis_id}\") | .confidence.current) = ${new_conf} | (.hypotheses[] | select(.id == \"${hypothesis_id}\") | .updated_at) = \"$(get_timestamp)\"" "$registry_file" > "$temp_registry"
        else
            jq "(.hypotheses[] | select(.id == \"${hypothesis_id}\") | .contradicting_evidence) += [\"${evidence_id}\"] | (.hypotheses[] | select(.id == \"${hypothesis_id}\") | .confidence.current) = ${new_conf} | (.hypotheses[] | select(.id == \"${hypothesis_id}\") | .updated_at) = \"$(get_timestamp)\"" "$registry_file" > "$temp_registry"
        fi
        mv "$temp_registry" "$registry_file"
    fi
    
    # Get updated confidence
    local new_confidence=$(jq -r ".hypotheses[] | select(.id == \"${hypothesis_id}\") | .confidence.current" "$registry_file" 2>/dev/null)
    
    echo ""
    echo -e "${GREEN}✓ Evidence Evaluated${NC}"
    echo "  Hypothesis:  ${hypothesis_id}"
    echo "  Evidence:    ${evidence_id}"
    echo "  Impact:      ${impact} ($impact_value%)"
    echo "  New Confidence: ${new_confidence}%"
    echo ""
    
    return 0
}

#-------------------------------------------------------------------------------
# Confidence Assessment
#-------------------------------------------------------------------------------

debug_check_confidence() {
    local hypothesis_id="${1:-}"
    local registry_file="${KDSE_DEBUG_DIR}/hypotheses/registry.json"
    
    echo ""
    echo -e "${BOLD}Confidence Assessment${NC}"
    echo "============================================================"
    
    if command -v jq &> /dev/null; then
        if [[ -n "$hypothesis_id" ]]; then
            # Check specific hypothesis
            local confidence=$(jq -r ".hypotheses[] | select(.id == \"${hypothesis_id}\") | .confidence.current" "$registry_file" 2>/dev/null)
            local threshold=$(jq -r ".hypotheses[] | select(.id == \"${hypothesis_id}\") | .confidence.threshold" "$registry_file" 2>/dev/null)
            
            echo "Hypothesis:   ${hypothesis_id}"
            echo "Confidence:  ${confidence}%"
            echo "Threshold:   ${threshold}%"
            
            if (( $(echo "$confidence >= $threshold" | bc -l 2>/dev/null || echo "0") )); then
                echo ""
                log_success "Confidence threshold met! Root cause can be selected."
                return 0
            else
                echo ""
                log_warning "Confidence below threshold. Continue collecting evidence."
                return 1
            fi
        else
            # Show all active hypotheses
            jq -r '.hypotheses | map(select(.status == "active")) | sort_by(.confidence.current) | reverse[] | "[\(.id)] \(.description) - \(.confidence.current)%"' "$registry_file" 2>/dev/null
        fi
    else
        grep -o '"id": "H-[^"]*".*"current": [0-9]*' "$registry_file" | sed 's/"id": "\(H-[^"]*\)".*"current": \([0-9]*\)/\1: \2%/'
    fi
    
    echo ""
}

#-------------------------------------------------------------------------------
# Root Cause Selection
#-------------------------------------------------------------------------------

debug_select_root_cause() {
    local hypothesis_id="${1:-}"
    local registry_file="${KDSE_DEBUG_DIR}/hypotheses/registry.json"
    
    if [[ -z "$hypothesis_id" ]]; then
        # Auto-select highest confidence hypothesis
        hypothesis_id=$(jq -r '.hypotheses | map(select(.status == "active")) | max_by(.confidence.current) | .id' "$registry_file" 2>/dev/null)
    fi
    
    local confidence=$(jq -r ".hypotheses[] | select(.id == \"${hypothesis_id}\") | .confidence.current" "$registry_file" 2>/dev/null)
    
    # Check confidence threshold
    if (( $(echo "$confidence < $CONFIDENCE_THRESHOLD" | bc -l 2>/dev/null || echo "1") )); then
        log_error "Confidence ($confidence%) below threshold ($CONFIDENCE_THRESHOLD%)"
        log_info "Cannot select root cause without operator override"
        return 1
    fi
    
    # Update hypothesis status
    debug_update_hypothesis_status "$hypothesis_id" "selected"
    
    # Update session
    CURRENT_STATE="ROOT_CAUSE_SELECTED"
    update_session_state
    sed -i "s/\"root_cause_selected\": null/\"root_cause_selected\": \"${hypothesis_id}\"/" "${SESSION_DIR}/session.json" 2>/dev/null || true
    
    echo ""
    echo -e "${GREEN}✓ Root Cause Selected${NC}"
    echo "============================================================"
    echo "Hypothesis:  ${hypothesis_id}"
    echo "Confidence: ${confidence}%"
    echo ""
    
    # Show recommendation
    local description=$(jq -r ".hypotheses[] | select(.id == \"${hypothesis_id}\") | .description" "$registry_file" 2>/dev/null)
    echo "Description:"
    echo "  $description"
    echo ""
    
    return 0
}

#-------------------------------------------------------------------------------
# Loop Detection
#-------------------------------------------------------------------------------

debug_check_loop() {
    local pattern_type="$1"
    local key="$2"
    
    if [[ "$LOOPS_ENABLED" != "true" ]]; then
        return 0
    fi
    
    local loop_history_file="${KDSE_DEBUG_DIR}/loops/history.json"
    local loop_key="${pattern_type}:${key}"
    
    # Check if this pattern has been seen before
    if [[ -f "$loop_history_file" ]]; then
        local count=$(jq -r ".loops[] | select(.key == \"${loop_key}\") | .count // 0" "$loop_history_file" 2>/dev/null || echo "0")
        
        if [[ "$count" -ge "$LOOPS_MAX_REPETITIONS" ]]; then
            log_warning "Loop detected: ${pattern_type} - ${key}"
            log_info "This investigation has been performed $count times"
            return 1
        fi
        
        # Increment count
        if command -v jq &> /dev/null; then
            if jq -e ".loops[] | select(.key == \"${loop_key}\")" "$loop_history_file" &>/dev/null; then
                # Update existing
                local temp_history="${loop_history_file}.tmp"
                jq "(.loops[] | select(.key == \"${loop_key}\") | .count) += 1 | (.loops[] | select(.key == \"${loop_key}\") | .last_occurrence) = \"$(get_timestamp)\"" "$loop_history_file" > "$temp_history"
                mv "$temp_history" "$loop_history_file"
            else
                # Add new
                local new_loop=$(cat <<EOF
{"key": "${loop_key}", "pattern_type": "${pattern_type}", "count": 1, "first_occurrence": "$(get_timestamp)", "last_occurrence": "$(get_timestamp)"}
EOF
)
                local temp_history="${loop_history_file}.tmp"
                echo "$(jq ".loops += [${new_loop}]" "$loop_history_file")" > "$temp_history"
                mv "$temp_history" "$loop_history_file"
            fi
        fi
    fi
    
    return 0
}

#-------------------------------------------------------------------------------
# Report Generation
#-------------------------------------------------------------------------------

debug_generate_report() {
    local format="${1:-json}"
    local report_file="${KDSE_DEBUG_DIR}/reports/${SESSION_ID}-report.${format}"
    
    local registry_file="${KDSE_DEBUG_DIR}/hypotheses/registry.json"
    local store_file="${KDSE_DEBUG_DIR}/evidence/store.json"
    
    echo ""
    echo -e "${BOLD}Generating Debug Report${NC}"
    echo "============================================================"
    
    # Create report
    if command -v jq &> /dev/null; then
        local root_cause=$(jq -r '.hypotheses[] | select(.status == "selected") // empty' "$registry_file" 2>/dev/null)
        local all_hypotheses=$(jq '.hypotheses' "$registry_file" 2>/dev/null)
        local all_evidence=$(jq '.evidence' "$store_file" 2>/dev/null)
        local evidence_count=$(jq '.evidence | length' "$store_file" 2>/dev/null)
        local hypothesis_count=$(jq '.hypotheses | length' "$registry_file" 2>/dev/null)
        
        local completed_at=$(get_timestamp)
        
        cat > "$report_file" <<EOF
{
  "session_id": "${SESSION_ID}",
  "started_at": "${SESSION_STARTED}",
  "completed_at": "${completed_at}",
  "failure": $(jq '.failure' "${SESSION_DIR}/session.json" 2>/dev/null),
  "root_cause": ${root_cause:-null},
  "evidence": {
    "collected": ${evidence_count:-0},
    "items": ${all_evidence:-[]}
  },
  "hypotheses": {
    "total": ${hypothesis_count:-0},
    "items": ${all_hypotheses:-[]}
  },
  "confidence_threshold": ${CONFIDENCE_THRESHOLD},
  "status": "completed"
}
EOF
    else
        # Simple text report
        cat > "$report_file" <<EOF
Debug Session Report
===================
Session: ${SESSION_ID}
Started: ${SESSION_STARTED}
Completed: $(get_timestamp)

Evidence Collected: $(grep -c '"id":' "$store_file" 2>/dev/null || echo "0")
Hypotheses Generated: $(grep -c '"id": "H-' "$registry_file" 2>/dev/null || echo "0")

Root Cause: $(grep '"status": "selected"' "$registry_file" | sed 's/.*"id": "\(H-[^"]*\)".*/\1/' || echo "None selected")
EOF
    fi
    
    echo "Report: ${report_file}"
    log_success "Debug report generated"
    echo ""
    
    return 0
}

#-------------------------------------------------------------------------------
# Session Completion
#-------------------------------------------------------------------------------

debug_complete() {
    local status="${1:-completed}"
    
    CURRENT_STATE="COMPLETED"
    update_session_state
    
    echo ""
    echo -e "${BOLD}============================================================${NC}"
    echo -e "${BOLD} Debug Session ${status^}${NC}"
    echo -e "${BOLD}============================================================${NC}"
    echo ""
    echo "Session ID: ${SESSION_ID}"
    echo "Duration:   $(get_timestamp) to ${SESSION_STARTED}"
    echo ""
    
    # Generate final report
    debug_generate_report
    
    return 0
}

#-------------------------------------------------------------------------------
# State Transitions
#-------------------------------------------------------------------------------

debug_next_phase() {
    case "$CURRENT_STATE" in
        "INITIAL")
            CURRENT_STATE="EVIDENCE_COLLECTION"
            ;;
        "EVIDENCE_COLLECTION")
            CURRENT_STATE="HYPOTHESIS_GENERATION"
            ;;
        "HYPOTHESIS_GENERATION")
            CURRENT_STATE="EVIDENCE_EVALUATION"
            ;;
        "EVIDENCE_EVALUATION")
            CURRENT_STATE="CONFIDENCE_ASSESSMENT"
            ;;
        "CONFIDENCE_ASSESSMENT")
            CURRENT_STATE="ROOT_CAUSE_SELECTED"
            ;;
        "ROOT_CAUSE_SELECTED")
            CURRENT_STATE="IMPLEMENTING"
            ;;
        "IMPLEMENTING")
            CURRENT_STATE="VERIFICATION"
            ;;
        "VERIFICATION")
            CURRENT_STATE="REGRESSION_TESTS"
            ;;
        "REGRESSION_TESTS")
            CURRENT_STATE="COMPLETED"
            ;;
        *)
            log_warning "Unknown state: $CURRENT_STATE"
            ;;
    esac
    
    update_session_state
    echo ""
    log_info "Phase: $CURRENT_STATE"
    echo ""
}

#-------------------------------------------------------------------------------
# Help
#-------------------------------------------------------------------------------

debug_help() {
    cat <<EOF
KDSE Debug Engine
=================

USAGE:
    source engine.sh && debug_init "failure description"
    debug_collect_evidence <type> <content> [source] [tags]
    debug_new_hypothesis <description> [confidence] [components]
    debug_evaluate <hypothesis_id> <evidence_id> <supporting|contradicting>
    debug_check_confidence [hypothesis_id]
    debug_select_root_cause [hypothesis_id]
    debug_check_loop <pattern_type> <key>
    debug_generate_report [json|markdown]
    debug_complete [status]
    debug_next_phase

STATES:
    INITIAL -> EVIDENCE_COLLECTION -> HYPOTHESIS_GENERATION ->
    EVIDENCE_EVALUATION -> CONFIDENCE_ASSESSMENT ->
    ROOT_CAUSE_SELECTED -> IMPLEMENTING -> VERIFICATION ->
    REGRESSION_TESTS -> COMPLETED

EVIDENCE TYPES:
    exception, test_failure, log, source, state, config, dependency

LOOP PATTERNS:
    file_inspection, module_reload, schema_check,
    database_inspection, cache_clear, import_analysis

EXAMPLES:
    debug_init "Database connection timeout"
    debug_collect_evidence "exception" "Connection refused" "app.py:42" "database,network"
    debug_new_hypothesis "Network firewall blocking port 5432" 40 "Database"
    debug_evaluate "H-0001" "E-0001" "supporting"
    debug_check_confidence "H-0001"
    debug_select_root_cause
    debug_generate_report

EOF
}

#-------------------------------------------------------------------------------
# Main
#-------------------------------------------------------------------------------

# If sourced, don't run main
if [[ "${BASH_SOURCE[0]}" != "${0}" ]]; then
    return 0
fi

# Parse arguments
COMMAND="${1:-help}"
shift 2>/dev/null || true

case "$COMMAND" in
    init)
        debug_init "$@"
        ;;
    collect)
        debug_collect_evidence "$@"
        ;;
    hypothesis|new)
        debug_new_hypothesis "$@"
        ;;
    evaluate)
        debug_evaluate "$@"
        ;;
    confidence|check)
        debug_check_confidence "$@"
        ;;
    select|root-cause)
        debug_select_root_cause "$@"
        ;;
    loop)
        debug_check_loop "$@"
        ;;
    report)
        debug_generate_report "$@"
        ;;
    complete|finish)
        debug_complete "$@"
        ;;
    next|phase)
        debug_next_phase
        ;;
    list-evidence|evidence)
        debug_list_evidence
        ;;
    list-hypotheses|hypotheses)
        debug_list_hypotheses
        ;;
    state)
        echo "Current state: $CURRENT_STATE"
        ;;
    help|--help|-h)
        debug_help
        ;;
    *)
        log_error "Unknown command: $COMMAND"
        debug_help
        exit 1
        ;;
esac
