#!/usr/bin/env python3
"""
KDSE MCP Tools - Python wrapper for KDSE HTTP MCP server.

This module provides Python bindings for the KDSE MCP tools that automatically
enforce KDSE principles. The enforcement is built into the runtime, not just
prompt instructions.

Usage:
    from kdse_shttp_tools import shttp_initialize, shttp_execute, shttp_status
    
    # Initialize KDSE workspace
    result = shttp_initialize()
    
    # Execute with automatic enforcement
    result = shttp_execute("Build a payroll system")
    
    # Check status
    status = shttp_status()
"""

import json
import os
import subprocess
import sys
from typing import Any, Dict, List, Optional
from urllib.request import urlopen, Request
from urllib.error import URLError, HTTPError

# Default MCP server configuration
DEFAULT_HOST = "localhost"
DEFAULT_PORT = 8080
KDSE_PORT_ENV = "KDSE_PORT"

# Try to get port from environment or use default
def get_port() -> int:
    return int(os.environ.get(KDSE_PORT_ENV, DEFAULT_PORT))

def get_base_url() -> str:
    return f"http://{DEFAULT_HOST}:{get_port()}"

# ============================================================================
# HTTP Request Helpers
# ============================================================================

def _make_request(endpoint: str, method: str = "GET", data: Optional[Dict] = None) -> Dict[str, Any]:
    """
    Make HTTP request to KDSE MCP server.
    
    Args:
        endpoint: API endpoint (e.g., "/status", "/execute")
        method: HTTP method (GET or POST)
        data: Request body data (for POST requests)
        
    Returns:
        Dict with 'success' and either 'data' or 'error'
    """
    url = f"{get_base_url()}{endpoint}"
    
    try:
        if data:
            json_data = json.dumps(data).encode('utf-8')
            req = Request(url, data=json_data, headers={'Content-Type': 'application/json'})
            req.get_method = lambda: method
        else:
            req = Request(url)
        
        with urlopen(req, timeout=30) as response:
            result = json.loads(response.read().decode('utf-8'))
            return result
            
    except HTTPError as e:
        return {
            "success": False,
            "error": f"HTTP Error {e.code}: {e.reason}",
            "data": None
        }
    except URLError as e:
        return {
            "success": False,
            "error": f"Connection Error: {e.reason}",
            "data": None
        }
    except Exception as e:
        return {
            "success": False,
            "error": str(e),
            "data": None
        }


def _post(endpoint: str, data: Optional[Dict] = None) -> Dict[str, Any]:
    """POST request helper"""
    return _make_request(endpoint, "POST", data)


def _get(endpoint: str) -> Dict[str, Any]:
    """GET request helper"""
    return _make_request(endpoint, "GET")


# ============================================================================
# Core KDSE Tools (MCP Protocol)
# ============================================================================

def shttp_status() -> Dict[str, Any]:
    """
    Get KDSE workspace and session status.
    
    Returns:
        Dict containing workspace readiness, session info, and orchestration state.
        
    Example:
        >>> status = shttp_status()
        >>> print(status['data']['kdse']['workspace_exists'])
        True
    """
    return _get("/status")


def shttp_initialize() -> Dict[str, Any]:
    """
    Initialize KDSE workspace and session.
    
    This creates the .kdse/ directory structure, session state, and runtime files.
    Must be called before any other KDSE operations.
    
    Returns:
        Dict containing session_id, phase, and runtime initialization results.
        
    Example:
        >>> result = shttp_initialize()
        >>> print(result['success'])
        True
    """
    return _post("/initialize")


def shttp_execute(
    objective: str,
    force: bool = False,
    auto_foundation: bool = True,
    auto_knowledge: bool = True,
    enforcement_level: str = "strict"
) -> Dict[str, Any]:
    """
    Execute a KDSE operation with AUTOMATIC enforcement of KDSE principles.
    
    This is the MAIN tool for KDSE operations. It automatically:
    - Checks if foundation exists and is complete
    - Checks if knowledge base has been built
    - Verifies architecture document exists
    - BLOCKS implementation if phases are incomplete
    - Auto-creates missing foundation documents (if enabled)
    - Generates audit warnings for violations
    
    Args:
        objective: The user's request (e.g., "Build a payroll system with biometrics")
        force: Skip enforcement (DANGEROUS - use with caution)
        auto_foundation: Auto-create missing foundation documents
        auto_knowledge: Auto-create missing knowledge categories
        enforcement_level: "off", "warning", "strict", or "hard"
        
    Returns:
        Dict containing:
        - blocked: bool - True if operation was blocked
        - blocked_reason: str - Why operation was blocked
        - violations: list - Any enforcement violations
        - audit_warnings: list - Non-blocking warnings
        - work_order: dict - Current phase work order
        - next_steps: list - Suggested next steps
        - report: str - Formatted report
        
    Example:
        >>> result = shttp_execute("Build a payroll system with biometrics")
        >>> if result['blocked']:
        ...     print("Foundation required:", result['violations'])
        ... else:
        ...     print("Ready to implement!")
    """
    return _post("/execute", {
        "objective": objective,
        "force": force,
        "auto_foundation": auto_foundation,
        "auto_knowledge": auto_knowledge,
        "enforcement_level": enforcement_level
    })


def shttp_foundation(operation: str = "status", documents: Optional[List[str]] = None) -> Dict[str, Any]:
    """
    Manage foundation documents.
    
    Args:
        operation: "status" to get status, "create" to create documents
        documents: List of document names to create (e.g., ["PROBLEM.md", "SPEC.md"])
        
    Returns:
        Dict containing foundation status or creation results.
        
    Example:
        >>> status = shttp_foundation("status")
        >>> print(status['data']['complete'])
        False
        >>> result = shttp_foundation("create", ["PROBLEM.md", "SPEC.md"])
    """
    if operation == "status":
        return _get("/foundation")
    else:
        return _post("/foundation", {"documents": documents or []})


def shttp_knowledge(operation: str = "status") -> Dict[str, Any]:
    """
    Manage knowledge base.
    
    Args:
        operation: "status" to get status, "create" to create categories
        
    Returns:
        Dict containing knowledge base status or creation results.
    """
    if operation == "status":
        return _get("/knowledge")
    else:
        return _post("/knowledge")


def shttp_audit() -> Dict[str, Any]:
    """
    Run compliance audit against KDSE standards.
    
    Returns:
        Dict containing audit report with violations, warnings, and evidence.
    """
    return _get("/audit")


def shttp_compliance() -> Dict[str, Any]:
    """
    Check runtime compliance.
    
    Returns:
        Dict containing compliance report and formatted output.
    """
    return _get("/compliance")


def shttp_enforce(
    operation: str = "implement",
    enforcement_level: str = "strict",
    auto_create: bool = True
) -> Dict[str, Any]:
    """
    Run enforcement check for a specific operation.
    
    Args:
        operation: Operation to check ("implement", "foundation", "knowledge", "architecture")
        enforcement_level: Enforcement strictness level
        auto_create: Auto-create missing artifacts
        
    Returns:
        Dict containing validation result and enforcement report.
    """
    return _post("/enforce", {
        "operation": operation,
        "enforcement_level": enforcement_level,
        "auto_create": auto_create
    })


def shttp_phase(operation: str = "status", new_phase: str = None, evidence: List[str] = None) -> Dict[str, Any]:
    """
    Manage orchestration phases.
    
    Args:
        operation: "status" to get current phase, "transition" to change phase
        new_phase: Target phase ("problem", "knowledge", "foundation", "audit", 
                   "assessment", "architecture", "implementation", "complete")
        evidence: List of evidence for phase transition
        
    Returns:
        Dict containing current phase info or transition result.
    """
    if operation == "status":
        return _get("/phase")
    else:
        return _post("/phase", {
            "phase": new_phase,
            "evidence": evidence or []
        })


def shttp_analyze() -> Dict[str, Any]:
    """
    Analyze repository structure.
    
    Creates a repository-analysis.md document in the knowledge base.
    
    Returns:
        Dict containing analysis results (file counts, types, etc.)
    """
    return _post("/analyze")


def shttp_collect() -> Dict[str, Any]:
    """
    Collect engineering evidence.
    
    Returns:
        Dict containing collection results.
    """
    return _post("/collect")


def shttp_report() -> Dict[str, Any]:
    """
    Generate comprehensive status report.
    
    Returns:
        Dict containing session, foundation, knowledge, and compliance status.
    """
    return _get("/report")


# ============================================================================
# Convenience Functions
# ============================================================================

def is_initialized() -> bool:
    """Check if KDSE workspace is initialized."""
    status = shttp_status()
    return status.get("success", False) and status.get("data", {}).get("kdse", {}).get("workspace_exists", False)


def can_implement() -> Dict[str, Any]:
    """
    Check if implementation is allowed.
    
    This is a quick check that returns whether the current state allows
    implementation or what needs to be done first.
    
    Returns:
        Dict with 'allowed' bool and 'violations' list if not allowed.
    """
    return shttp_enforce("implement")


def check_knowledge_ready() -> bool:
    """Quick check if knowledge base is ready."""
    status = shttp_knowledge()
    if not status.get("success"):
        return False
    data = status.get("data", {})
    return data.get("complete", False)


def check_foundation_ready() -> bool:
    """Quick check if foundation is ready."""
    status = shttp_foundation()
    if not status.get("success"):
        return False
    data = status.get("data", {})
    return data.get("complete", False)


def get_current_phase() -> Optional[str]:
    """Get current orchestration phase."""
    result = shttp_phase()
    if result.get("success"):
        return result.get("data", {}).get("current_phase")
    return None


def transition_to(phase: str) -> Dict[str, Any]:
    """
    Transition to a new orchestration phase.
    
    Args:
        phase: Target phase name
        
    Returns:
        Dict with transition result.
    """
    return shttp_phase("transition", phase)


# ============================================================================
# Demo / Test Functions
# ============================================================================

def demo_enforcement():
    """
    Demonstrate KDSE enforcement with the payroll system example.
    
    This shows how KDSE blocks premature implementation.
    """
    print("=" * 70)
    print("KDSE ENFORCEMENT DEMONSTRATION")
    print("=" * 70)
    print()
    
    # Check status first
    print("1. Checking KDSE status...")
    status = shttp_status()
    if status.get("success"):
        print(f"   Workspace exists: {status['data']['kdse']['workspace_exists']}")
        print(f"   Session active: {status['data']['guard']['session_active']}")
    print()
    
    # Try to execute (should block if no foundation)
    print("2. Attempting to execute 'Build payroll system with biometrics'...")
    result = shttp_execute("Build a payroll system with biometrics")
    
    if result.get("blocked", False):
        print("   ❌ BLOCKED - Foundation/Knowledge phase incomplete")
        print(f"   Reason: {result.get('blocked_reason', 'Unknown')}")
        print()
        if result.get("violations"):
            print("   Violations:")
            for v in result["violations"]:
                print(f"     - [{v.get('code')}] {v.get('message')}")
                if v.get("corrective_action"):
                    print(f"       Fix: {v.get('corrective_action')}")
        print()
        print("   Next steps:")
        for step in result.get("next_steps", [])[:5]:
            print(f"     → {step}")
    else:
        print("   ✓ Allowed to proceed")
    
    print()
    print("=" * 70)


if __name__ == "__main__":
    # If run directly, demonstrate the tools
    demo_enforcement()
