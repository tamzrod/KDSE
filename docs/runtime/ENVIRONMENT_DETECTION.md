# Environment Detection Algorithm

**Document Version:** 1.0  
**Type:** Normative Technical Specification  
**Effective Date:** 2026-07-19  
**Supersedes:** N/A (New Document)

---

## Purpose

Environment Detection classifies the current workspace into exactly one environment type before any planning or engineering activity. This ensures KDSE behaves correctly based on where it is executing:

- **KDSE_RUNTIME**: Avoid treating KDSE itself as a target project
- **KDSE_PROJECT**: Apply project-specific engineering workflows
- **SOFTWARE_PROJECT**: Treat as an existing project requiring initialization
- **BLANK_WORKSPACE**: Prompt for project creation before proceeding

---

## Environment Types

### BLANK_WORKSPACE

**Definition**: No recognizable project exists, no Git repository, no KDSE project.

**Evidence Required**:
- No project indicators (go.mod, package.json, pyproject.toml, etc.)
- No `.kdse/` directory
- No Git repository

**Detection Rationale**:
- Zero project structure indicates a workspace that needs initialization
- KDSE cannot proceed without a valid target project
- Avoids attempting engineering work in empty directories

**Detection Rule**:
```go
func isBlankWorkspace(e Evidence) bool {
    return !e.HasKDSEProject && 
           !e.HasGoMod && 
           !e.HasPackageJSON && 
           !e.HasPyProject && 
           !e.HasCargoToml && 
           !e.HasPomXML && 
           !e.HasMakefile
}
```

**Failure Cases**:
- User accidentally runs KDSE in home directory
- Container starts with empty working directory
- Network drive with stale mounts

---

### SOFTWARE_PROJECT

**Definition**: An existing software project that may contain Git, is not a KDSE runtime, and has no initialized KDSE project.

**Evidence Required**:
- At least one project indicator file:
  - `go.mod` (Go)
  - `package.json` (Node.js)
  - `pyproject.toml` or `setup.py` (Python)
  - `Cargo.toml` (Rust)
  - `pom.xml` (Java)
  - `Makefile` (C/C++)
- No `.kdse/` directory
- Git repository is optional evidence (configurable)

**Detection Rationale**:
- KDSE attaches to existing projects, never creating them
- Project indicators prove this is a legitimate software project
- Absence of `.kdse/` confirms not yet initialized for KDSE
- Git is evidence of project maturity but not required

**Detection Rule**:
```go
func isSoftwareProject(e Evidence) bool {
    if e.HasKDSEProject {
        return false  // Not a software project if it's a KDSE project
    }
    
    hasIndicator := e.HasGoMod || 
                    e.HasPackageJSON || 
                    e.HasPyProject || 
                    e.HasCargoToml || 
                    e.HasPomXML || 
                    e.HasMakefile
    
    return hasIndicator
}
```

**Failure Cases**:
- Project with unusual structure (e.g., monorepo without standard indicators)
- Mixed-language projects with conflicting indicators
- Projects with only documentation files

---

### KDSE_RUNTIME

**Definition**: The current repository is the KDSE framework/runtime itself. Development work targets KDSE, not a normal software project.

**Evidence Required** (at least 2):
1. `go.mod` with module name `github.com/kdse/runtime`
2. `internal/runtime/` directory exists
3. `internal/detection/` directory exists
4. `templates/` directory exists
5. `docs/` directory exists

**Why Multiple Markers?**:
- Single markers could match other projects coincidentally
- KDSE runtime has a distinctive structure with multiple internal packages
- Prevents false positives from similarly-named repositories

**Detection Rationale**:
- KDSE must recognize itself to avoid self-modification
- Prevents treating KDSE as a target for engineering governance
- Enables special commands like `kdse self-assess`

**Repository Markers**:

| Marker | Location | Description |
|--------|----------|-------------|
| Module Name | `go.mod` | Must contain `github.com/kdse/runtime` |
| Runtime Package | `internal/runtime/` | Core runtime implementation |
| Detection Package | `internal/detection/` | Environment detection subsystem |
| Templates | `templates/` | KDSE project templates |
| Documentation | `docs/` | KDSE documentation |

**Detection Rule**:
```go
func isKDSERuntime(e Evidence) bool {
    markerCount := 0
    
    if e.HasGoModKDSE {           markerCount++ }  // github.com/kdse/runtime
    if e.HasRuntimeDir {          markerCount++ }  // internal/runtime/
    if e.HasDetectionDir {        markerCount++ }  // internal/detection/
    if e.HasTemplatesDir {        markerCount++ }  // templates/
    if e.HasDocsDir {             markerCount++ }  // docs/
    
    return markerCount >= 2  // Require at least 2 markers
}
```

**Failure Cases**:
- Forked KDSE repository with different structure
- KDSE inside a monorepo with multiple packages
- Development branches with temporary modifications

---

### KDSE_PROJECT

**Definition**: An existing project that was initialized by KDSE. Contains valid KDSE runtime metadata. Engineering work targets the project, not the KDSE runtime.

**Evidence Required**:
1. `.kdse/` directory MUST exist
2. `manifest.yaml` or `manifest.json` MUST exist
3. Manifest MUST be valid (parseable JSON/YAML)
4. Must NOT be the KDSE runtime itself

**Manifest Structure**:
```json
{
  "schema": "https://kdse.dev/schemas/manifest/v1.0",
  "version": "1.0.0",
  "generated": "2026-07-11T00:00:00Z",
  "status": "ACTIVE",
  "runtime": {
    "version": "2.0"
  }
}
```

**Detection Rationale**:
- Valid manifest proves KDSE was initialized
- Status field indicates initialization state
- Separate from KDSE_RUNTIME to enable self-hosting
- Project type from manifest supports workflow routing

**Detection Rule**:
```go
func isKDSEProject(e Evidence) bool {
    if !e.HasKDSEProject {
        return false  // .kdse/ must exist
    }
    
    if !e.HasManifest {
        return false  // manifest must exist
    }
    
    if e.ManifestValid {
        return false  // manifest must be valid
    }
    
    // Must not be the KDSE runtime itself
    if e.HasGoModKDSE || e.HasRuntimeDir {
        return false
    }
    
    return true
}
```

**Failure Cases**:
- Corrupted manifest file
- Incomplete initialization (partial .kdse/ structure)
- Manifest from incompatible KDSE version

---

## Detection Algorithm

### Order of Evaluation

Detection follows strict priority order to ensure deterministic results:

```
1. KDSE_RUNTIME (Priority 1)
   ↓ Not matched?
2. KDSE_PROJECT (Priority 2)
   ↓ Not matched?
3. SOFTWARE_PROJECT (Priority 3)
   ↓ Not matched?
4. BLANK_WORKSPACE (Priority 4)
```

**Why This Order?**

1. **KDSE_RUNTIME first**: Prevents KDSE from treating itself as a target. This is the most critical distinction.

2. **KDSE_PROJECT second**: Once KDSE itself is ruled out, check if this is a KDSE-enabled project. These have special handling requirements.

3. **SOFTWARE_PROJECT third**: A normal project needs different initialization paths than a KDSE project.

4. **BLANK_WORKSPACE last**: If nothing matched, the workspace is empty and needs guidance.

### Detection Flow

```
START
  ↓
Collect Evidence
  ├── Git Repository (.git/)
  ├── Project Indicators (go.mod, package.json, etc.)
  ├── KDSE Markers (.kdse/, manifest)
  └── Runtime Markers (internal/runtime/, etc.)
  ↓
Evaluate KDSE_RUNTIME
  └── Has ≥2 KDSE markers? → YES → return KDSE_RUNTIME
  ↓ NO
Evaluate KDSE_PROJECT
  └── Has .kdse/ + valid manifest? → YES → return KDSE_PROJECT
  ↓ NO
Evaluate SOFTWARE_PROJECT
  └── Has project indicators? → YES → return SOFTWARE_PROJECT
  ↓ NO
return BLANK_WORKSPACE
```

---

## Evidence Collection

### Evidence Structure

```go
type Evidence struct {
    // Git Evidence
    HasGitRepo     bool   // .git directory exists
    GitRemoteURL   string // Remote URL if available
    GitRootPath    string // Git repository root
    
    // Project Evidence
    HasGoMod       bool   // go.mod exists
    HasPackageJSON bool   // package.json exists
    HasPyProject   bool   // pyproject.toml or setup.py
    HasCargoToml   bool   // Cargo.toml exists
    HasPomXML      bool   // pom.xml exists
    HasMakefile    bool   // Makefile exists
    
    // KDSE Evidence
    HasKDSEProject bool   // .kdse/ directory exists
    HasManifest    bool   // manifest.yaml/json exists
    ManifestValid  bool   // manifest is parseable
    ManifestData   *ManifestInfo
    
    // KDSE Runtime Markers
    HasGoModKDSE    bool   // module github.com/kdse/runtime
    HasRuntimeDir   bool   // internal/runtime/ exists
    HasDetectionDir bool   // internal/detection/ exists
    HasTemplatesDir bool   // templates/ exists
    HasDocsDir      bool   // docs/ exists
    
    // Generic Indicators
    HasReadme  bool   // README.md exists
    HasSrcDir  bool   // src/ exists
    HasTestDir bool   // tests/ exists
    HasLicense bool   // LICENSE exists
    
    // Path Information
    RepoRoot    string // Repository root
    ModuleName  string // Go module name
    ProjectType string // Detected project type
}
```

### Evidence Collection Priority

1. **Fast checks first**: Directory existence checks (O(1))
2. **File parsing second**: go.mod module name extraction
3. **Manifest validation third**: JSON/YAML parsing
4. **Git operations last**: Only if required by options

---

## Confidence Calculation

### Confidence Levels

| Environment | Confidence Basis |
|-------------|------------------|
| KDSE_RUNTIME | Number of KDSE markers / 5 |
| KDSE_PROJECT | 0.95 (valid manifest), 0.7 (manifest only) |
| SOFTWARE_PROJECT | 0.6 + (indicators × 0.05), max 0.95 |
| BLANK_WORKSPACE | 0.8 (default) |

### Confidence Threshold

- **≥ 0.8**: High confidence, proceed with environment-specific workflows
- **0.6 - 0.8**: Medium confidence, log warnings, proceed with care
- **< 0.6**: Low confidence, may need manual verification

---

## Failure Cases and Handling

### Insufficient Evidence

When evidence is contradictory or incomplete:

```
Detection Options:
  AllowAmbiguous: false  // Default
    → Return UNKNOWN environment
    
  AllowAmbiguous: true
    → Return most likely environment with warnings
```

### Contradictory Evidence

Example: `.kdse/` exists but no valid manifest, AND project indicators present.

**Handling**:
1. Log warning about contradictory evidence
2. Attempt recovery (e.g., check for manifest.json if manifest.yaml missing)
3. Return to BLANK_WORKSPACE if unrecoverable

### Detection Errors

| Error | Severity | Handling |
|-------|----------|----------|
| Cannot read directory | WARNING | Continue with partial evidence |
| Cannot parse manifest | WARNING | Mark manifest as invalid |
| Cannot access .git | WARNING | Continue without Git evidence |
| Path not accessible | ERROR | Return UNKNOWN |

---

## API Reference

### Detector

```go
// Create a detector for a repository path
detector := NewEnvironmentDetector("/path/to/repo")

// Create with custom options
opts := &DetectionOptions{
    RequireGitRepo: true,
    StrictManifestValidation: true,
}
detector := NewEnvironmentDetectorWithOptions("/path/to/repo", opts)

// Detect environment
result := detector.Detect()

// Access results
fmt.Println(result.Environment)     // Environment type
fmt.Println(result.Confidence)     // 0.0 - 1.0
fmt.Println(result.Evidence)        // Full evidence
fmt.Println(result.Warnings)        // Non-fatal issues
```

### Convenience Functions

```go
// Quick environment detection
env := DetectEnvironment("/path/to/repo")

// Access specific evidence
result := NewEnvironmentDetector("/path/to/repo").Detect()
if result.Evidence.HasGitRepo {
    fmt.Println("Git repository detected")
}
```

---

## Classification Rules Reference

### Rule: KDSE Runtime Detection

| Property | Value |
|----------|-------|
| Priority | 1 |
| Result | KDSE_RUNTIME |
| Required Evidence | ≥2 of: go.mod(github.com/kdse/runtime), internal/runtime/, internal/detection/, templates/, docs/ |

**Rationale**: The KDSE runtime has a distinctive multi-package structure. No single marker is sufficient to distinguish it from similarly-named repositories.

### Rule: KDSE Project Detection

| Property | Value |
|----------|-------|
| Priority | 2 |
| Result | KDSE_PROJECT |
| Required Evidence | .kdse/ directory, manifest.yaml/json, valid manifest |

**Rationale**: KDSE initialization creates a specific directory structure with a valid manifest. This proves intentional initialization, not accidental file presence.

### Rule: Software Project Detection

| Property | Value |
|----------|-------|
| Priority | 3 |
| Result | SOFTWARE_PROJECT |
| Required Evidence | ≥1 project indicator, no .kdse/ |

**Rationale**: Language-specific files (go.mod, package.json, etc.) indicate a legitimate software project. Git is optional evidence.

### Rule: Blank Workspace Detection

| Property | Value |
|----------|-------|
| Priority | 4 |
| Result | BLANK_WORKSPACE |
| Required Evidence | No project indicators, no .kdse/ |

**Rationale**: Default classification when no other evidence is present. Indicates an empty or unrecognized workspace.

---

## Configuration Options

### DetectionOptions

```go
type DetectionOptions struct {
    // AllowAmbiguous allows returning a classification even with conflicting evidence
    AllowAmbiguous bool
    
    // RequireGitRepo requires a Git repository for software project detection
    RequireGitRepo bool
    
    // StrictManifestValidation enables strict validation of manifest files
    StrictManifestValidation bool
}
```

### Option Behaviors

| Option | Default | Effect |
|--------|---------|--------|
| AllowAmbiguous | false | When true, returns classification even with insufficient evidence |
| RequireGitRepo | false | When true, software projects without Git are classified as BLANK_WORKSPACE |
| StrictManifestValidation | true | When true, invalid manifests cause KDSE_PROJECT detection to fail |

---

## Extensibility

### Adding New Environment Types

1. Add new `EnvironmentType` constant
2. Implement new detection rule in `classify()`
3. Update `calculateConfidence()` for new type
4. Add tests for new type
5. Update documentation

### Adding New Evidence Markers

1. Add new field to `Evidence` struct
2. Implement evidence collection in appropriate `collect*()` function
3. Update detection rules to use new evidence
4. Add tests for new evidence
5. Document new marker in this specification

---

## Testing

### Test Coverage Requirements

| Environment | Test Cases |
|-------------|-----------|
| BLANK_WORKSPACE | Empty dir, only docs, only readme |
| SOFTWARE_PROJECT | Each project type, with/without Git |
| KDSE_RUNTIME | All marker combinations, partial markers |
| KDSE_PROJECT | Valid/invalid manifests, partial structures |

### Integration Testing

```bash
# Test in KDSE repository (should detect KDSE_RUNTIME)
cd /path/to/kdse && go test ./internal/detection/...

# Test in KDSE project (should detect KDSE_PROJECT)
cd /path/to/my-project/.kdse/.. && go test ./internal/detection/...

# Test in software project
cd /path/to/go-project && go test ./internal/detection/...
```

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-07-19 | Initial specification |

---

*This document defines the Environment Detection algorithm for KDSE. It is normative for the detection subsystem implementation.*
