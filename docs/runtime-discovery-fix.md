# KDSE Runtime Discovery Bug Fix

## Root Cause Analysis

### The Problem
KDSE was resolving the runtime location from the server's execution environment (e.g., `/app`) instead of the user's Git repository. This caused `.kdse` to be created in the wrong location.

### Root Cause
Multiple locations in the codebase used `os.Getwd()` to determine the project path:

1. **CLI (`cmd/kdse/main.go:40`)**: Used `os.Getwd()` directly
   ```go
   repoPath, _ := os.Getwd()
   ```

2. **MCP Tools (`cmd/mcp/tools/tools.go:33`)**: Used `os.Getwd()` before resolving to Git root
   ```go
   cwd, _ := os.Getwd()
   gitRoot := cwd
   resolver := detection.NewGitResolver(cwd)
   ```

3. **Orchestration Resolver (`internal/orchestration/resolver.go:17`)**: Used `os.Getwd()` to initialize
   ```go
   wd, err := os.Getwd()
   ```

When the KDSE server runs in a Docker container or MCP environment, `os.Getwd()` returns the container's working directory (like `/app`), not the user's project path. Even though the code tried to resolve to Git root, it was starting from the wrong location.

### The Fundamental Issue
The architecture violated the principle that **KDSE is repository-centric**. The runtime MUST always live at `<git root>/.kdse`, exactly beside `.git` and `.github`.

---

## Solution

### Created Shared Runtime Discovery Package

**File**: `internal/discover/discover.go`

```go
// Resolve discovers the KDSE runtime paths from a project path.
// This function implements the core KDSE runtime discovery rules:
// 1. If projectPath is provided, start from there
// 2. If projectPath is empty, use current working directory
// 3. Resolve to Git repository root
// 4. Runtime path = <git root>/.kdse
func Resolve(projectPath string) (*RuntimePaths, error)

// ResolveRuntime is a convenience function that returns only the runtime path.
func ResolveRuntime(projectPath string) (string, error)

// ResolveRepository is a convenience function that returns only the repository path.
func ResolveRepository(projectPath string) (string, error)
```

### Updated All Components to Use Shared Discovery

| Component | File | Changes |
|-----------|------|---------|
| CLI | `cmd/kdse/main.go` | Uses `discover.Resolve()` with optional project path argument |
| MCP Tools | `cmd/mcp/tools/tools.go` | `NewToolHandler(projectPath)` accepts project path |
| MCP Server | `cmd/mcp/main.go` | Supports `KDSE_PROJECT_PATH` env var |
| Orchestration | `internal/orchestration/resolver.go` | Uses `discover.Resolve()` in constructors |
| Guard | `internal/guard/coordinator.go` | Uses `discover.Resolve()` in `NewCoordinator()` |

---

## Files Modified

### New Files Created
1. `internal/discover/discover.go` - Shared runtime discovery package
2. `internal/discover/discover_test.go` - Comprehensive test suite

### Modified Files
1. `cmd/kdse/main.go` - Updated to use `discover.Resolve()`
2. `cmd/mcp/main.go` - Added `KDSE_PROJECT_PATH` support
3. `cmd/mcp/tools/tools.go` - Updated `NewToolHandler()` signature
4. `internal/orchestration/resolver.go` - Updated constructors
5. `internal/guard/coordinator.go` - Updated to use discover package

---

## Test Coverage

### Test Cases Added

| Test | Description | Status |
|------|-------------|--------|
| `TestResolveFromRepoRoot` | Resolution from repository root | ✓ |
| `TestResolveFromNestedDirectory` | Resolution from nested dir (`repo/cmd/server/internal`) | ✓ |
| `TestResolveFromEmptyPath` | Empty path uses cwd | ✓ |
| `TestResolveNoGitRepository` | Error when no git repo | ✓ |
| `TestResolveRuntime` | Convenience function | ✓ |
| `TestResolveRepository` | Convenience function | ✓ |
| `TestMustResolvePanics` | Panic on error | ✓ |
| `TestEnsureRuntime` | Runtime directory creation | ✓ |
| `TestRuntimePathsString` | String method | ✓ |
| `TestInvalidProjectPath` | Invalid path handling | ✓ |
| `TestHasGitRepository` | Helper function | ✓ |
| `TestSubmoduleGitRepo` | Git submodule support | ✓ |
| `TestDifferentCwdThanProject` | CWD independence | ✓ |
| `TestRepositoryWithGitHub` | Repository with `.github` | ✓ |
| `TestRepositoryAlreadyHasKDSE` | Pre-existing `.kdse` | ✓ |
| `TestRuntimeIndependentOfServerCwd` | **CRITICAL: Server cwd independence** | ✓ |

### Critical Test: `TestRuntimeIndependentOfServerCwd`
This test proves the bug is fixed by:
1. Creating a Git repo at `/tmp/kdse-test-xxx`
2. Changing working directory to `/tmp` (simulating server at `/app`)
3. Resolving runtime with explicit project path
4. Asserting that `RuntimePath` resolves to `<repo>/.kdse`, NOT `/tmp/.kdse`

---

## Initialization Rules

When `kdse initialize` runs:

1. ✓ Resolve Git repository root using `discover.Resolve()`
2. ✓ Runtime path = `<repo>/.kdse`
3. ✓ If missing, create `.kdse`
4. ✓ Never create runtime outside the repository

### Failure Conditions

If no Git repository exists:
- ✓ Returns explicit error: `ErrNoGitRepository`
- ✓ **DO NOT** create `/app/.kdse`, `/repo/.kdse`, `~/.kdse`, `/workspace/.kdse`
- ✓ Runtime must never exist outside a Git repository

---

## Success Criteria Met

### All Commands Now Resolve the Same Runtime

| Command | Resolution Method |
|---------|------------------|
| `kdse initialize` | `discover.Resolve()` → `<git root>/.kdse` |
| `kdse status` | `discover.Resolve()` → `<git root>/.kdse` |
| `kdse collect` | `discover.Resolve()` → `<git root>/.kdse` |
| `kdse assess` | `discover.Resolve()` → `<git root>/.kdse` |
| `kdse report` | `discover.Resolve()` → `<git root>/.kdse` |
| `kdse run` | `discover.Resolve()` → `<git root>/.kdse` |

### Independent of Execution Environment

| Environment | Resolution |
|-------------|------------|
| Server cwd | ✓ Uses project path instead |
| Docker | ✓ Uses project path instead |
| MCP | ✓ Uses project path instead |
| HTTP server | ✓ Uses project path instead |
| OpenHands | ✓ Uses project path instead |
| VS Code | ✓ Uses project path instead |
| Claude Code | ✓ Uses project path instead |
| Terminal location | ✓ Uses project path or cwd |

---

## Usage Examples

### CLI
```bash
# Use current directory (resolves via Git)
kdse status

# Use explicit project path (ignores cwd)
kdse /path/to/project status
kdse initialize /path/to/project
```

### MCP Server
```bash
# Set project path via environment variable
export KDSE_PROJECT_PATH=/path/to/project
kdse-mcp

# For HTTP transport, use header
curl -H "X-KDSE-Project-Path: /path/to/project" http://localhost:8080/status
```

---

## Backwards Compatibility

- CLI: Falls back to cwd if no project path provided
- MCP: Falls back to cwd if `KDSE_PROJECT_PATH` not set
- All internal packages handle resolution failures gracefully

---

## Architecture Summary

```
┌─────────────────────────────────────────────────────────────┐
│                    Client Request                           │
│           (project path or default to cwd)                  │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              discover.Resolve(projectPath)                   │
│                                                              │
│  1. If projectPath is empty, use cwd                       │
│  2. Resolve to Git repository root                         │
│  3. Runtime path = <git root>/.kdse                        │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              RuntimePaths                                   │
│                                                              │
│  - RepositoryPath: /path/to/project                        │
│  - RuntimePath: /path/to/project/.kdse                      │
│  - GitRoot: /path/to/project                               │
│  - IsGitRepo: true                                         │
└─────────────────────────────────────────────────────────────┘
```

---

## Verification

To verify the fix works, run the tests:
```bash
go test ./internal/discover/... -v
```

The critical test `TestRuntimeIndependentOfServerCwd` will pass, proving that:
1. Runtime resolution is independent of server's working directory
2. Project path provided by client takes precedence
3. `.kdse` is always created at `<git root>/.kdse`
