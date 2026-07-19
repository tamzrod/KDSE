# KDSE Toolchain Management

## Overview

KDSE Runtime now includes automatic toolchain management to ensure all required development tools are available before entering the implementation phase. This eliminates the common problem of tests being silently skipped due to missing tooling.

## Key Principles

1. **Verification is mandatory** - KDSE never silently skips verification due to missing tooling
2. **Writable installation only** - Toolchains are installed to workspace locations, not system directories
3. **Automatic verification** - All required toolchains are verified before implementation begins
4. **Fail loudly** - If toolchains are missing, KDSE provides actionable instructions and stops

## Supported Toolchains

| Toolchain | Detection Method | Auto-Install |
|-----------|-----------------|---------------|
| Go | `go version` | ✓ (downloads to writable location) |
| Node.js | `node --version` | ✓ (downloads to writable location) |
| Python | `python3 --version` | Manual (requires package manager) |
| Java | `java -version` | Manual (requires JDK) |
| .NET | `dotnet --version` | Manual (installer required) |
| Rust | `rustc --version` | Manual (rustup required) |

## Installation Paths

Toolchains are installed to a configurable writable location:

```
Default: ~/.kdse/tools/
Custom:  $KDSE_TOOLS_PATH
```

Within this directory:
```
~/.kdse/tools/
├── go/           # Go installation
├── node/         # Node.js installation
└── gopath/       # Go workspace
```

## Usage

### Automatic Verification

When KDSE initializes, it automatically:

1. Detects which toolchains are required by the project
2. Checks if each toolchain is available in PATH
3. If missing, attempts to install to the writable location
4. Updates PATH and environment variables
5. Verifies the installation works

### Example: Go Project

```bash
# Create a new Go project
mkdir my-project && cd my-project
go mod init github.com/user/my-project

# Initialize KDSE
kdse initialize

# KDSE automatically:
# 1. Detects go.mod exists
# 2. Checks for Go installation
# 3. If missing, downloads to ~/.kdse/tools/go/
# 4. Updates PATH
# 5. Verifies Go works
# 6. Proceeds with implementation
```

### Example: Node.js Project

```bash
# Create a new Node.js project
mkdir my-app && cd my-app
npm init

# Initialize KDSE
kdse initialize

# KDSE automatically detects package.json and handles Node.js
```

### Example: Python Project

```bash
# Create a new Python project
mkdir my-app && cd my-app
python3 -m venv venv
source venv/bin/activate
pip freeze > requirements.txt

# Initialize KDSE
kdse initialize
```

## Configuration

### Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `KDSE_TOOLS_PATH` | Installation base path | `~/.kdse/tools` |
| `KDSE_TOOLCHAIN_LEVEL` | Verification strictness | `required` |

### Verification Levels

- `required` - Only verify toolchains required by the project (default)
- `all` - Verify all supported toolchains
- `none` - Skip verification (not recommended)

## Error Handling

### Missing Toolchain

When a required toolchain is missing:

```
╔═══════════════════════════════════════════════════════════════╗
║         KDSE Runtime Toolchain Verification                 ║
╠═══════════════════════════════════════════════════════════════╣
║ Project: /workspace/project/my-app
║ Required Toolchains:
║   Go: ✗ MISSING
╠═══════════════════════════════════════════════════════════════╣
║ Status: VERIFICATION FAILED                                ║
╠═══════════════════════════════════════════════════════════════╣
║ Missing Toolchains:
║   - go
╚═══════════════════════════════════════════════════════════════╝

To install Go:
  1. Download from: https://go.dev/dl/
  2. Extract to:   ~/.kdse/tools/go
  3. Add to PATH:  export PATH=~/.kdse/tools/go/bin:$PATH
```

### Installation Log

KDSE maintains a log of toolchain verification:

```
~/.kdse/runtime/toolchain-log.txt
```

This log includes:
- Timestamp of verification
- Required toolchains
- Verification results
- Any errors encountered

## Programmatic Usage

### Go API

```go
import "github.com/kdse/runtime/internal/toolchain"

// Create a toolchain manager
manager := toolchain.NewManager(projectPath)

// Detect all available toolchains
toolchains := manager.DetectAll()

// Verify a specific toolchain
result := manager.Verify(toolchain.ToolchainGo)

// Ensure a toolchain is available (auto-install if missing)
result := manager.Ensure(toolchain.ToolchainGo)

// Verify all project requirements
results := manager.EnsureProjectRequirements()

// Get environment updates for shell
envUpdates := manager.UpdateEnvironment()
```

### CLI Integration

The `kdse initialize` command now includes toolchain verification:

```bash
kdse initialize
```

Output includes:
- Project detected (with type)
- Required toolchains listed
- Verification status for each
- Instructions for any missing toolchains

## Verification Policy

### Correct Behavior

```
Go missing
    ↓
Install Go
    ↓
Verify installation
    ↓
Run tests
```

### Incorrect Behavior (Now Prevented)

```
Go missing
    ↓
Skip tests  ← NEVER HAPPENS
    ↓
Continue implementation  ← NEVER HAPPENS
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    KDSE Initialization                        │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              Project Discovery                               │
│         (Detect language-specific files)                     │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              Toolchain Detection                            │
│                                                              │
│  For each required toolchain:                               │
│    1. Check PATH                                           │
│    2. If missing, check writable installation              │
│    3. If none, install to ~/.kdse/tools/                  │
│    4. Update PATH/GOROOT/etc.                             │
│    5. Verify installation works                            │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              Verification                                   │
│                                                              │
│  All toolchains verified?                                   │
│    ├── YES → Continue with implementation                   │
│    └── NO  → Report missing toolchains → STOP               │
└─────────────────────────────────────────────────────────────┘
```

## Troubleshooting

### Toolchain Not Found After Installation

If you install a toolchain but KDSE still can't find it:

1. Check PATH includes the toolchain bin directory:
   ```bash
   echo $PATH
   ```

2. Manually add to PATH:
   ```bash
   export PATH=~/.kdse/tools/go/bin:$PATH
   ```

3. Verify installation:
   ```bash
   go version
   ```

### Custom Installation Path

To use a custom installation path:

```bash
export KDSE_TOOLS_PATH=/workspace/.tools
kdse initialize
```

### Skip Verification (Not Recommended)

To skip toolchain verification:

```bash
export KDSE_TOOLCHAIN_LEVEL=none
kdse initialize
```

This is not recommended as it may lead to build failures later.

## Future Enhancements

- Automatic installation for all toolchains (not just Go and Node)
- Integration with version managers (nvm, pyenv, etc.)
- Toolchain version pinning per project
- Cache verification results
- Toolchain update notifications
