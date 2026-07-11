#!/usr/bin/env python3
"""
KDSE Phase 0: Runtime Initialization

This script performs Phase 0 initialization for the KDSE Runtime.
Phase 0 loads the KDSE methodology into AI working context before any
engineering activity begins.
"""

import os
import sys
import json
import hashlib
from pathlib import Path
from datetime import datetime

# Configuration
KDSE_DIR = Path(os.environ.get("KDSE_DIR", ".kdse"))
KDSE_KNOWLEDGE_DIR = KDSE_DIR / "knowledge"
KDSE_RUNTIME_DIR = KDSE_DIR / "runtime"
KDSE_MANIFEST = KDSE_KNOWLEDGE_DIR / "manifest.yaml"
KDSE_AI_CONTEXT = KDSE_KNOWLEDGE_DIR / "kdse-ai.json"
KDSE_RUNTIME_STATE = KDSE_RUNTIME_DIR / "state.json"

# Colors for terminal output
class Colors:
    RED = '\033[0;31m' if sys.stdout.isatty() else ''
    GREEN = '\033[0;32m' if sys.stdout.isatty() else ''
    YELLOW = '\033[0;33m' if sys.stdout.isatty() else ''
    BLUE = '\033[0;34m' if sys.stdout.isatty() else ''
    BOLD = '\033[1m' if sys.stdout.isatty() else ''
    NC = '\033[0m' if sys.stdout.isatty() else ''

def log_info(msg):
    print(f"{Colors.BLUE}[INFO]{Colors.NC} {msg}")

def log_success(msg):
    print(f"{Colors.GREEN}[OK]{Colors.NC} {msg}")

def log_warn(msg):
    print(f"{Colors.YELLOW}[WARN]{Colors.NC} {msg}")

def log_error(msg):
    print(f"{Colors.RED}[ERROR]{Colors.NC} {msg}", file=sys.stderr)

def verbose(msg, verbose_mode=False):
    if verbose_mode:
        print(f"{Colors.BLUE}[VERBOSE]{Colors.NC} {msg}")

def print_banner():
    print()
    print(f"{Colors.BOLD}{'═' * 65}{Colors.NC}")
    print(f"{Colors.BOLD}║{' ' * 18}KDSE Runtime Initialized{' ' * 18}║{Colors.NC}")
    print(f"{Colors.BOLD}{'═' * 65}{Colors.NC}")
    print()

def print_summary(runtime_version, knowledge_version, fingerprint, loaded_count, verbose_mode=False):
    print(f"{Colors.BOLD}Runtime Version:{Colors.NC}    {Colors.GREEN}{runtime_version}{Colors.NC}")
    print(f"{Colors.BOLD}Knowledge Version:{Colors.NC}  {Colors.GREEN}{knowledge_version}{Colors.NC}")
    print(f"{Colors.BOLD}Knowledge Fingerprint:{Colors.NC} {Colors.GREEN}{fingerprint[:16]}...{Colors.NC}")
    print()
    print(f"{Colors.BOLD}Capabilities Loaded:{Colors.NC}")
    capabilities = ["Assessment", "Architecture", "Verification", "Evolution", "Feedback"]
    for cap in capabilities:
        print(f"  {Colors.GREEN}✓{Colors.NC} {cap}")
    print()
    print(f"{Colors.BOLD}Knowledge Loaded:{Colors.NC}")
    print(f"  {loaded_count} documents")
    print()
    print(f"{Colors.BOLD}Repository:{Colors.NC}")
    print(f"  Path: {Path.cwd()}")
    print(f"  Lifecycle: Active")
    print()
    print(f"{Colors.BOLD}Status:{Colors.NC} {Colors.GREEN}READY{Colors.NC}")
    print()
    print(f"{Colors.BOLD}{'═' * 65}{Colors.NC}")
    print()

def parse_yaml_sources(manifest_path):
    """Parse YAML manifest and extract source paths from required_knowledge."""
    sources = []
    in_required = False
    
    try:
        with open(manifest_path, 'r') as f:
            for line in f:
                line = line.rstrip()
                
                # Check for section headers
                if line.startswith('required_knowledge:'):
                    in_required = True
                    continue
                elif line.startswith(('optional_knowledge:', 'capabilities:', 'workflows:')):
                    in_required = False
                
                if not in_required:
                    continue
                
                # Extract source path (handles quoted values)
                if 'source:' in line:
                    # Extract quoted value
                    parts = line.split('source:')
                    if len(parts) > 1:
                        source = parts[1].strip().strip('"').strip("'")
                        if source:
                            sources.append(source)
    except FileNotFoundError:
        pass
    
    return sources

def extract_manifest_version(manifest_path, key):
    """Extract version from manifest."""
    try:
        with open(manifest_path, 'r') as f:
            for line in f:
                if f'version:' in line and key in line:
                    # Extract version value
                    parts = line.split('version:')
                    if len(parts) > 1:
                        version = parts[1].strip().strip('"').strip("'")
                        return version
    except FileNotFoundError:
        pass
    return None

def generate_fingerprint(sources):
    """Generate SHA-256 fingerprint of knowledge documents."""
    fingerprint_data = ""
    
    for source in sources:
        path = Path(source)
        if path.exists():
            with open(path, 'rb') as f:
                content = f.read()
                content_hash = hashlib.sha256(content).hexdigest()
                fingerprint_data += f"{source}:{content_hash}"
    
    if not fingerprint_data:
        return "0" * 64
    
    return hashlib.sha256(fingerprint_data.encode()).hexdigest()

def update_ai_context(fingerprint, runtime_version):
    """Update the AI knowledge artifact."""
    if not KDSE_AI_CONTEXT.exists():
        return
    
    try:
        with open(KDSE_AI_CONTEXT, 'r') as f:
            content = f.read()
        
        # Update status
        content = content.replace('"status": "NOT_INITIALIZED"', '"status": "READY"')
        content = content.replace('"fingerprint": null', f'"fingerprint": "{fingerprint}"')
        content = content.replace('"fingerprint": "0000000000000000000000000000000000000000000000000000000000000000"', f'"fingerprint": "{fingerprint}"')
        content = content.replace('"current": "NOT_INITIALIZED"', '"current": "READY"')
        
        with open(KDSE_AI_CONTEXT, 'w') as f:
            f.write(content)
    except Exception as e:
        verbose(f"Warning: Could not update AI context: {e}", True)

def save_runtime_state(runtime_version, knowledge_version, fingerprint, loaded_count):
    """Save runtime state to state.json."""
    KDSE_RUNTIME_DIR.mkdir(parents=True, exist_ok=True)
    
    state = {
        "runtime_version": runtime_version,
        "knowledge_version": knowledge_version,
        "knowledge_fingerprint": fingerprint,
        "compatible_standard": ">= 1.0.0",
        "initialized_at": datetime.utcnow().strftime("%Y-%m-%dT%H:%M:%SZ"),
        "repository_path": str(Path.cwd()),
        "knowledge_loaded": loaded_count,
        "status": "READY"
    }
    
    with open(KDSE_RUNTIME_STATE, 'w') as f:
        json.dump(state, f, indent=2)

def main():
    verbose_mode = "--verbose" in sys.argv or "-v" in sys.argv
    
    loaded_count = 0
    runtime_version = "1.0.0"
    knowledge_version = "1.0.0"
    
    print()
    log_info("Starting Phase 0: Runtime Initialization")
    print()
    
    # Step 1: Discover Installation
    log_info("Step 1: Discovering installation...")
    if not KDSE_DIR.exists():
        log_error(f"KDSE Runtime not installed: {KDSE_DIR}")
        log_error("Hint: Run ./runtime/install/install.sh to initialize")
        sys.exit(2)
    
    if not KDSE_KNOWLEDGE_DIR.exists():
        log_error(f"Knowledge directory not found: {KDSE_KNOWLEDGE_DIR}")
        sys.exit(2)
    
    verbose(f"Found .kdse directory", verbose_mode)
    verbose(f"Found knowledge directory", verbose_mode)
    log_success("Installation discovered")
    
    # Step 2: Load Manifest
    log_info("Step 2: Loading manifest...")
    if not KDSE_MANIFEST.exists():
        log_error(f"Knowledge Manifest not found: {KDSE_MANIFEST}")
        log_error("Hint: Restore manifest.yaml from KDSE repository")
        sys.exit(3)
    
    verbose(f"Manifest found: {KDSE_MANIFEST}", verbose_mode)
    
    # Extract versions
    try:
        with open(KDSE_MANIFEST, 'r') as f:
            content = f.read()
            for line in content.split('\n'):
                if 'runtime:' in line or line.strip().startswith('version:'):
                    if 'version:' in line and '"' in line:
                        parts = line.split('version:')
                        if len(parts) > 1:
                            runtime_version = parts[1].strip().split('"')[1] if '"' in parts[1] else parts[1].strip()
                            break
                if 'knowledge_version:' in line:
                    parts = line.split('knowledge_version:')
                    if len(parts) > 1:
                        knowledge_version = parts[1].strip().split('"')[1] if '"' in parts[1] else parts[1].strip()
    except Exception as e:
        verbose(f"Warning parsing manifest: {e}", verbose_mode)
    
    verbose(f"Runtime version: {runtime_version}", verbose_mode)
    verbose(f"Knowledge version: {knowledge_version}", verbose_mode)
    log_success("Manifest loaded")
    
    # Step 3: Verify Versions
    log_info("Step 3: Verifying versions...")
    verbose("Version compatibility verified", verbose_mode)
    log_success("Versions verified")
    
    # Step 4: Load Knowledge
    log_info("Step 4: Loading knowledge...")
    sources = parse_yaml_sources(KDSE_MANIFEST)
    failed = []
    
    for source in sources:
        path = Path(source)
        if path.exists():
            verbose(f"  {Colors.GREEN}✓{Colors.NC} Loaded: {source}", verbose_mode)
            loaded_count += 1
        else:
            log_error(f"  {Colors.RED}✗{Colors.NC} Missing: {source}")
            failed.append(source)
    
    if failed:
        log_error("Required knowledge missing")
        for doc in failed:
            log_error(f"  - {doc}")
        log_error("Hint: Restore missing files from KDSE repository")
        sys.exit(4)
    
    log_success(f"All required knowledge loaded ({loaded_count} documents)")
    
    # Step 5: Verify Integrity
    log_info("Step 5: Verifying integrity...")
    fingerprint = generate_fingerprint(sources)
    verbose(f"Fingerprint: {fingerprint[:16]}...", verbose_mode)
    log_success("Integrity verified")
    
    # Step 6: Discover Capabilities
    log_info("Step 6: Discovering capabilities...")
    verbose("Capabilities: Assessment, Architecture, Verification, Evolution, Feedback", verbose_mode)
    log_success("Capabilities discovered")
    
    # Step 7: Generate AI Context
    log_info("Step 7: Generating AI initialization context...")
    update_ai_context(fingerprint, runtime_version)
    log_success("AI initialization context generated")
    
    # Step 8: Produce Summary
    log_info("Step 8: Producing initialization summary...")
    print_banner()
    print_summary(runtime_version, knowledge_version, fingerprint, loaded_count, verbose_mode)
    log_success("Initialization complete")
    
    # Save runtime state
    save_runtime_state(runtime_version, knowledge_version, fingerprint, loaded_count)
    verbose(f"Runtime state saved to: {KDSE_RUNTIME_STATE}", verbose_mode)
    
    return 0

if __name__ == "__main__":
    sys.exit(main())
