# KDSE Runtime Discovery Fix - Project-First Architecture

## Overview

KDSE has been redesigned to use a **project-first** architecture. KDSE is an engineering runtime that **attaches** to existing software projects instead of creating them.

## The Problem

The previous architecture assumed:
1. KDSE should initialize first
2. A Git repository must exist
3. KDSE might need to create project structures

This led to failure scenarios:
```
Blank Workspace → KDSE Initialize → Fail → Create Git → Fail → Create .kdse → Fail → Abandon KDSE
```

## New Architecture

KDSE now follows this lifecycle:

```
Workspace
    ↓
Project Discovery
    ↓
Project Exists?
    ├── YES → KDSE Initialization → Continue
    └── NO  → Request Project Initialization → Wait → Retry Discovery
```

### Key Principles

1. **KDSE never creates the project** - The project must exist first
2. **Git is only optional evidence** - Not a requirement
3. **KDSE detects via language-specific files** - go.mod, package.json, pyproject.toml, etc.

## Project Discovery

### Detection Order
1. Scan for language-specific project files
2. Optionally detect Git repository (as evidence)
3. Resolve to project root

### Supported Project Types

| Language/Framework | Indicators |
|---------------------|------------|
| Go | `go.mod`, `go.sum`, `main.go`, `cmd/`, `internal/` |
| Node.js | `package.json`, `package-lock.json`, `node_modules/` |
| Python | `pyproject.toml`, `setup.py`, `requirements.txt`, `venv/` |
| Rust | `Cargo.toml`, `Cargo.lock`, `src/` |
| Java | `pom.xml`, `build.gradle`, `src/main/java/` |
| .NET | `.sln`, `.csproj`, `Program.cs` |
| PHP | `composer.json`, `artisan`, `public/index.php` |
| C/C++ | `Makefile`, `CMakeLists.txt`, `*.c`, `*.h` |

## Blank Workspace Behavior

When no project exists:

**DO NOT:**
- Create `.git`
- Initialize Git
- Create `.kdse`
- Fabricate project files
- Fabricate runtime metadata

**INSTEAD:**
Return a clear message:
```
No software project detected.

KDSE requires a software project before initialization.

Please initialize your project first:
  - Go:       go mod init github.com/user/project
  - Node.js:  npm init
  - Python:   create pyproject.toml or requirements.txt
  - Rust:     cargo init
  - Java:     create pom.xml or build.gradle
  - .NET:     dotnet new

Then run kdse initialize.
```

## Files Modified

### Core Changes

1. **`internal/discover/discover.go`** - Project-first detection
   - `Resolve()` now uses `detectProject()` first
   - Git is detected as optional evidence via `hasGitRepository()`
   - Returns `RuntimePaths` with `ProjectType` and `ProjectIndicators`

2. **`internal/guard/project_guard.go`** - Uses discover package
   - `Validate()` now uses `discover.Resolve()` for project detection
   - Removed Git as requirement

3. **`internal/guard/coordinator.go`** - Updated initialization
   - `EnsureProject()` returns `ProjectInfo` with project details
   - `NoProjectError` provides actionable message

4. **`internal/guard/types.go`** - Updated error messages
   - `ErrNoProjectDetected` with clear hint

5. **`cmd/kdse/main.go`** - Updated CLI
   - Shows project type and Git status on initialization
   - Displays clear message when no project detected

## Usage Examples

### Existing Project (Go)
```bash
cd my-go-project
kdse initialize
# Output shows:
#   Project: my-go-project
#   Type: go
#   Git: detected
```

### Existing Project (Node.js)
```bash
cd my-react-app
kdse initialize
# Output shows:
#   Project: my-react-app
#   Type: node
#   Git: optional
```

### Blank Workspace
```bash
cd empty-workspace
kdse initialize
# Output shows clear message to initialize project first
```

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Client Request                           │
│           (project path or default to cwd)                  │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              detectProject(projectPath)                      │
│                                                              │
│  1. Scan for language-specific files                       │
│  2. Match against ProjectIndicators                       │
│  3. Require at least 2 indicators for valid project        │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              hasGitRepository(projectRoot)                   │
│                                                              │
│  OPTIONAL - Git is only evidence, not required              │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              RuntimePaths                                   │
│                                                              │
│  - ProjectPath: /path/to/project                           │
│  - ProjectRoot: /path/to/project                           │
│  - RuntimePath: /path/to/project/.kdse                      │
│  - ProjectType: go                                          │
│  - ProjectIndicators: [go.mod, main.go, cmd/]              │
│  - IsGitRepo: true/false (optional)                        │
└─────────────────────────────────────────────────────────────┘
```

## Verification

Run tests:
```bash
go test ./internal/discover/... -v
go test ./internal/guard/... -v
```

## Acceptance Criteria

- ✅ KDSE no longer attempts to initialize inside an empty workspace
- ✅ KDSE never creates Git repositories
- ✅ KDSE treats Git only as optional evidence
- ✅ KDSE initializes only after a valid software project exists
- ✅ Blank workspaces produce a clear "project required" response
- ✅ Existing projects initialize normally
- ✅ Recovery logic no longer loops creating Git or .kdse
