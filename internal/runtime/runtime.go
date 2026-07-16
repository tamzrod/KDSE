// Package runtime implements evidence-driven KDSE runtime management.
// Every operation follows: Execute → Verify → Report
package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// VerificationResult represents the verification status of an artifact
type VerificationResult struct {
	Artifact   string `json:"artifact"`
	Path       string `json:"path"`
	Status     string `json:"status"` // "PASS" or "FAIL"
	Evidence   string `json:"evidence,omitempty"`
	Error      string `json:"error,omitempty"`
	Timestamp  string `json:"timestamp"`
}

// RuntimeManifest defines the complete runtime structure
type RuntimeManifest struct {
	Version     string                 `json:"version"`
	CreatedAt   string                 `json:"created_at"`
	Directories []ManifestDirectory    `json:"directories"`
	Files       []ManifestFile         `json:"files"`
	Invariants  []RuntimeInvariant     `json:"invariants"`
}

// ManifestDirectory describes a required directory
type ManifestDirectory struct {
	Path     string `json:"path"`
	Required bool   `json:"required"`
	Purpose  string `json:"purpose"`
}

// ManifestFile describes a required file
type ManifestFile struct {
	Path     string `json:"path"`
	Required bool   `json:"required"`
	Purpose  string `json:"purpose"`
	Template string `json:"template,omitempty"`
}

// RuntimeInvariant defines a phase transition requirement
type RuntimeInvariant struct {
	Phase       string   `json:"phase"`
	Requires    []string `json:"requires"`
	Description string   `json:"description"`
}

// InitializeResult is the result of runtime initialization
type InitializeResult struct {
	Success       bool                 `json:"success"`
	WorkspacePath string               `json:"workspace_path"`
	Confidence    float64              `json:"confidence"`
	Verification  []VerificationResult `json:"verification"`
	Evidence      []string             `json:"evidence"`
	Errors        []string             `json:"errors,omitempty"`
	Timestamp     string               `json:"timestamp"`
}

// VerificationReport is the result of runtime verification
type VerificationReport struct {
	Success      bool                 `json:"success"`
	Confidence   float64              `json:"confidence"`
	Components   []VerificationResult `json:"components"`
	Missing      []string             `json:"missing,omitempty"`
	Failed       []string             `json:"failed,omitempty"`
	Timestamp    string               `json:"timestamp"`
}

// Standard runtime directories
const (
	DirRuntime        = "runtime"
	DirFoundation     = "foundation"
	DirKnowledge      = "knowledge"
	DirLaboratory     = "laboratory"
	DirEvidence       = "evidence"
	DirReferences     = "references"
	DirTraceability   = "traceability"
	DirReports        = "reports"
	DirConfig         = "config"
	DirState          = "state"
	DirArtifacts      = "artifacts"
	DirSessions       = "sessions"
	DirNormalized     = "normalized"
	DirCache          = "cache"
	DirSomeday       = "someday"
)

// Standard runtime files
const (
	FileManifest        = "manifest.yaml"
	FileSessionState    = "session-state.yaml"
	FileRuntimeConfig   = "runtime.yaml"
	FileKnowledgeIndex  = "knowledge-index.yaml"
	FileArtifactIndex   = "artifact-index.yaml"
)

// Runtime defines the evidence-driven runtime manager
type Runtime struct {
	repoPath  string
	kdsePath  string
	manifest  *RuntimeManifest
	verified  bool
	invariant *InvariantEngine
}

// New creates a new Runtime for the given repository path
func New(repoPath string) *Runtime {
	return &Runtime{
		repoPath:  repoPath,
		kdsePath:  filepath.Join(repoPath, ".kdse"),
		manifest:  DefaultManifest(),
		verified:  false,
		invariant: NewInvariantEngine(),
	}
}

// Initialize creates a full operational KDSE runtime
// Returns evidence of what was created and verified
func (r *Runtime) Initialize() *InitializeResult {
	result := &InitializeResult{
		WorkspacePath: r.kdsePath,
		Timestamp:     time.Now().Format(time.RFC3339),
		Verification:  []VerificationResult{},
		Evidence:      []string{},
	}

	// Execute: Create all directories
	r.executeDirectoryCreation(result)

	// Execute: Create all required files
	r.executeFileCreation(result)

	// Verify: Check every artifact
	r.verifyAllArtifacts(result)

	// Calculate confidence based on verification
	result.Confidence = r.calculateConfidence(result.Verification)

	// Determine success
	result.Success = r.determineSuccess(result)

	return result
}

// executeDirectoryCreation creates all required directories
func (r *Runtime) executeDirectoryCreation(result *InitializeResult) {
	directories := []string{
		DirRuntime,
		DirFoundation,
		DirKnowledge,
		DirLaboratory,
		DirEvidence,
		DirReferences,
		DirTraceability,
		DirReports,
		DirConfig,
		DirState,
		DirArtifacts,
		DirSessions,
		DirNormalized,
		DirCache,
		DirSomeday,
	}

	for _, dir := range directories {
		path := filepath.Join(r.kdsePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			result.Verification = append(result.Verification, VerificationResult{
				Artifact:  dir,
				Path:      path,
				Status:    "FAIL",
				Error:     err.Error(),
				Timestamp: time.Now().Format(time.RFC3339),
			})
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to create directory %s: %v", dir, err))
		}
	}

	// Create knowledge subdirectories
	knowledgeSubdirs := []string{"general", "operational", "developmental"}
	for _, subdir := range knowledgeSubdirs {
		path := filepath.Join(r.kdsePath, DirKnowledge, subdir)
		if err := os.MkdirAll(path, 0755); err != nil {
			result.Verification = append(result.Verification, VerificationResult{
				Artifact:  DirKnowledge + "/" + subdir,
				Path:      path,
				Status:    "FAIL",
				Error:     err.Error(),
				Timestamp: time.Now().Format(time.RFC3339),
			})
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to create directory %s/%s: %v", DirKnowledge, subdir, err))
		}
	}

	// Create laboratory subdirectories
	labSubdirs := []string{"experiments", "reports"}
	for _, subdir := range labSubdirs {
		path := filepath.Join(r.kdsePath, DirLaboratory, subdir)
		if err := os.MkdirAll(path, 0755); err != nil {
			result.Verification = append(result.Verification, VerificationResult{
				Artifact:  DirLaboratory + "/" + subdir,
				Path:      path,
				Status:    "FAIL",
				Error:     err.Error(),
				Timestamp: time.Now().Format(time.RFC3339),
			})
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to create directory %s/%s: %v", DirLaboratory, subdir, err))
		}
	}

	// Create someday subdirectories
	somedaySubdirs := []string{"ideas", "archived", "promoted"}
	for _, subdir := range somedaySubdirs {
		path := filepath.Join(r.kdsePath, DirSomeday, subdir)
		if err := os.MkdirAll(path, 0755); err != nil {
			result.Verification = append(result.Verification, VerificationResult{
				Artifact:  DirSomeday + "/" + subdir,
				Path:      path,
				Status:    "FAIL",
				Error:     err.Error(),
				Timestamp: time.Now().Format(time.RFC3339),
			})
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to create directory %s/%s: %v", DirSomeday, subdir, err))
		}
	}
}

// executeFileCreation creates all required runtime files and templates
func (r *Runtime) executeFileCreation(result *InitializeResult) {
	// Core runtime files
	files := map[string]string{
		FileManifest:       r.generateManifestContent(),
		FileSessionState:   r.generateSessionStateContent(),
		FileRuntimeConfig:  r.generateRuntimeConfigContent(),
		FileKnowledgeIndex: r.generateKnowledgeIndexContent(),
		FileArtifactIndex:  r.generateArtifactIndexContent(),
	}

	for filename, content := range files {
		path := filepath.Join(r.kdsePath, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			result.Verification = append(result.Verification, VerificationResult{
				Artifact:  filename,
				Path:      path,
				Status:    "FAIL",
				Error:     err.Error(),
				Timestamp: time.Now().Format(time.RFC3339),
			})
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to create file %s: %v", filename, err))
		}
	}

	// Create someday manifest
	somedayManifestPath := filepath.Join(r.kdsePath, DirSomeday, "someday.yaml")
	somedayContent := r.generateSomedayManifestContent()
	if err := os.WriteFile(somedayManifestPath, []byte(somedayContent), 0644); err != nil {
		result.Verification = append(result.Verification, VerificationResult{
			Artifact:  DirSomeday + "/someday.yaml",
			Path:      somedayManifestPath,
			Status:    "FAIL",
			Error:     err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		})
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to create file %s: %v", DirSomeday+"/someday.yaml", err))
	}

	// Create all template files
	r.createTemplateFiles(result)
}

// createTemplateFiles creates all engineering templates
func (r *Runtime) createTemplateFiles(result *InitializeResult) {
	// Foundation templates
	foundationTemplates := map[string]string{
		"README.md":         r.generateFoundationReadme(),
		"PROBLEM.md":        r.generateProblemTemplate(),
		"SPEC.md":           r.generateSpecTemplate(),
		"ARCHITECTURE.md":   r.generateArchitectureTemplate(),
		"ASSUMPTIONS.md":    r.generateAssumptionsTemplate(),
		"REQUIREMENTS.md":   r.generateRequirementsTemplate(),
	}

	foundationPath := filepath.Join(r.kdsePath, DirFoundation)
	for filename := range foundationTemplates {
		r.writeTemplateFile(filepath.Join(foundationPath, filename), filename, result)
	}

	// Knowledge templates
	knowledgeTemplates := map[string]string{
		"README.md": r.generateKnowledgeReadme(),
	}
	knowledgePath := filepath.Join(r.kdsePath, DirKnowledge)
	for filename := range knowledgeTemplates {
		r.writeTemplateFile(filepath.Join(knowledgePath, filename), DirKnowledge+"/"+filename, result)
	}

	// Knowledge subdirectories
	knowledgeSubs := map[string]map[string]string{
		"general": {
			"README.md": r.generateKnowledgeSubReadme("general", "General engineering knowledge, patterns, and principles"),
		},
		"operational": {
			"README.md": r.generateKnowledgeSubReadme("operational", "Operational knowledge including runbooks and procedures"),
		},
		"developmental": {
			"README.md": r.generateKnowledgeSubReadme("developmental", "Development-specific knowledge including patterns and practices"),
		},
	}
	for subdir, files := range knowledgeSubs {
		for filename := range files {
			r.writeTemplateFile(filepath.Join(knowledgePath, subdir, filename), DirKnowledge+"/"+subdir+"/"+filename, result)
		}
	}

	// Laboratory templates
	labTemplates := map[string]string{
		"README.md": r.generateLaboratoryReadme(),
	}
	labPath := filepath.Join(r.kdsePath, DirLaboratory)
	for filename := range labTemplates {
		r.writeTemplateFile(filepath.Join(labPath, filename), DirLaboratory+"/"+filename, result)
	}

	// Laboratory subdirectories
	labSubs := map[string]map[string]string{
		"experiments": {
			"README.md": r.generateLabSubReadme("experiments", "Hypothesis-driven experiments with documented outcomes"),
		},
		"reports": {
			"README.md": r.generateLabSubReadme("reports", "Laboratory reports and validation results"),
		},
	}
	for subdir, files := range labSubs {
		for filename := range files {
			r.writeTemplateFile(filepath.Join(labPath, subdir, filename), DirLaboratory+"/"+subdir+"/"+filename, result)
		}
	}

	// Evidence templates
	evidenceTemplates := map[string]string{
		"README.md": r.generateEvidenceReadme(),
	}
	evidencePath := filepath.Join(r.kdsePath, DirEvidence)
	for filename := range evidenceTemplates {
		r.writeTemplateFile(filepath.Join(evidencePath, filename), DirEvidence+"/"+filename, result)
	}

	// Reports templates
	reportsTemplates := map[string]string{
		"README.md": r.generateReportsReadme(),
	}
	reportsPath := filepath.Join(r.kdsePath, DirReports)
	for filename := range reportsTemplates {
		r.writeTemplateFile(filepath.Join(reportsPath, filename), DirReports+"/"+filename, result)
	}

	// References templates
	referencesTemplates := map[string]string{
		"README.md": r.generateReferencesReadme(),
	}
	referencesPath := filepath.Join(r.kdsePath, DirReferences)
	for filename := range referencesTemplates {
		r.writeTemplateFile(filepath.Join(referencesPath, filename), DirReferences+"/"+filename, result)
	}

	// Traceability templates
	traceTemplates := map[string]string{
		"README.md": r.generateTraceabilityReadme(),
	}
	tracePath := filepath.Join(r.kdsePath, DirTraceability)
	for filename := range traceTemplates {
		r.writeTemplateFile(filepath.Join(tracePath, filename), DirTraceability+"/"+filename, result)
	}

	// Artifacts templates
	artifactsTemplates := map[string]string{
		"README.md": r.generateArtifactsReadme(),
	}
	artifactsPath := filepath.Join(r.kdsePath, DirArtifacts)
	for filename := range artifactsTemplates {
		r.writeTemplateFile(filepath.Join(artifactsPath, filename), DirArtifacts+"/"+filename, result)
	}
}

// writeTemplateFile writes a template file and records verification
func (r *Runtime) writeTemplateFile(path, artifactName string, result *InitializeResult) {
	content := r.getTemplateContent(path)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		result.Verification = append(result.Verification, VerificationResult{
			Artifact:  artifactName,
			Path:      path,
			Status:    "FAIL",
			Error:     err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		})
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to create template %s: %v", artifactName, err))
	}
}

// getTemplateContent returns the template content based on filename
func (r *Runtime) getTemplateContent(path string) string {
	filename := filepath.Base(path)
	dir := filepath.Base(filepath.Dir(path))

	switch {
	case filename == "PROBLEM.md":
		return r.generateProblemTemplate()
	case filename == "SPEC.md":
		return r.generateSpecTemplate()
	case filename == "ARCHITECTURE.md":
		return r.generateArchitectureTemplate()
	case filename == "ASSUMPTIONS.md":
		return r.generateAssumptionsTemplate()
	case filename == "REQUIREMENTS.md":
		return r.generateRequirementsTemplate()
	case dir == "knowledge" && filename == "README.md":
		return r.generateKnowledgeReadme()
	case filename == "README.md":
		return r.generateReadmeForDir(dir)
	default:
		return r.generateDefaultReadme()
	}
}

// Foundation templates

func (r *Runtime) generateFoundationReadme() string {
	return `# Foundation Directory

This directory contains the project foundation documents that define the engineering work.

## Contents

| Document | Purpose | Phase |
|----------|---------|-------|
| PROBLEM.md | Problem statement | Problem |
| SPEC.md | Solution specification | Foundation |
| ARCHITECTURE.md | System architecture | Architecture |
| ASSUMPTIONS.md | Project assumptions | Foundation |
| REQUIREMENTS.md | Detailed requirements | Foundation |

## Engineering Phases

1. **Problem Phase** - Update PROBLEM.md
2. **Foundation Phase** - Update SPEC.md, ASSUMPTIONS.md, REQUIREMENTS.md
3. **Architecture Phase** - Update ARCHITECTURE.md

## Rule

DO NOT create new files in this directory.
Only update the existing templates with project-specific content.
`
}

func (r *Runtime) generateProblemTemplate() string {
	return `# Problem Statement

**Created:** ` + time.Now().Format("2006-01-02") + `
**Phase:** Problem

## Problem Summary

[One paragraph describing the core problem being solved]

## Problem Details

### Background
[Context and history that led to this problem]

### Impact
[Who is affected and how]

### Symptoms
[Observable symptoms of the problem]

### Constraints
[Known constraints and limitations]

## Success Criteria

- [Criterion 1]
- [Criterion 2]
- [Criterion 3]

---
*This is a template. Replace with project-specific content.*
`
}

func (r *Runtime) generateSpecTemplate() string {
	return `# Solution Specification

**Created:** ` + time.Now().Format("2006-01-02") + `
**Phase:** Foundation

## Overview

[High-level description of the solution]

## Goals

### Primary Goals
1. [Goal 1]
2. [Goal 2]
3. [Goal 3]

### Secondary Goals
1. [Goal 1]
2. [Goal 2]

## Non-Goals

- [What this solution will NOT address]

## Approach

[High-level approach to solving the problem]

## Scope

### In Scope
- [Item 1]
- [Item 2]

### Out of Scope
- [Item 1]
- [Item 2]

## Dependencies

- [Dependency 1]
- [Dependency 2]

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| [Risk 1] | [Impact] | [Mitigation] |

---
*This is a template. Replace with project-specific content.*
`
}

func (r *Runtime) generateArchitectureTemplate() string {
	return `# System Architecture

**Created:** ` + time.Now().Format("2006-01-02") + `
**Phase:** Architecture

## Architecture Overview

[High-level architectural description]

## System Components

### Component 1
- **Purpose**: [Purpose]
- **Responsibilities**: [Responsibilities]
- **Dependencies**: [Dependencies]

### Component 2
- **Purpose**: [Purpose]
- **Responsibilities**: [Responsibilities]
- **Dependencies**: [Dependencies]

## Data Model

### Entity 1
| Field | Type | Description |
|-------|------|-------------|
| [field] | [type] | [description] |

## API Design

### Endpoint 1
- **Method**: [GET/POST/PUT/DELETE]
- **Path**: [Path]
- **Description**: [Description]

## Technology Stack

| Component | Technology | Version |
|-----------|------------|---------|
| [Component] | [Technology] | [Version] |

## Security Considerations

- [Security consideration 1]
- [Security consideration 2]

## Deployment Architecture

[How the system will be deployed]

---
*This is a template. Replace with project-specific content.*
`
}

func (r *Runtime) generateAssumptionsTemplate() string {
	return `# Project Assumptions

**Created:** ` + time.Now().Format("2006-01-02") + `
**Phase:** Foundation

## Technology Assumptions

| Assumption | Rationale | Status |
|------------|-----------|--------|
| [Assumption 1] | [Why we assume this] | [Valid/Invalid] |
| [Assumption 2] | [Why we assume this] | [Valid/Invalid] |

## Environment Assumptions

| Assumption | Rationale | Status |
|------------|-----------|--------|
| [Assumption 1] | [Why we assume this] | [Valid/Invalid] |

## Business Assumptions

| Assumption | Rationale | Status |
|------------|-----------|--------|
| [Assumption 1] | [Why we assume this] | [Valid/Invalid] |

## Dependency Assumptions

| Assumption | Rationale | Status |
|------------|-----------|--------|
| [Assumption 1] | [Why we assume this] | [Valid/Invalid] |

## Review History

| Date | Reviewer | Changes |
|------|----------|---------|
| [Date] | [Name] | [Changes made] |

---
*This is a template. Replace with project-specific content.*
`
}

func (r *Runtime) generateRequirementsTemplate() string {
	return `# Detailed Requirements

**Created:** ` + time.Now().Format("2006-01-02") + `
**Phase:** Foundation

## Functional Requirements

### FR-001: [Requirement Title]
- **Description**: [Detailed description]
- **Priority**: [High/Medium/Low]
- **Acceptance Criteria**:
  - [Criterion 1]
  - [Criterion 2]

### FR-002: [Requirement Title]
- **Description**: [Detailed description]
- **Priority**: [High/Medium/Low]
- **Acceptance Criteria**:
  - [Criterion 1]
  - [Criterion 2]

## Non-Functional Requirements

### Performance
- [Requirement 1]
- [Requirement 2]

### Security
- [Requirement 1]
- [Requirement 2]

### Reliability
- [Requirement 1]
- [Requirement 2]

### Usability
- [Requirement 1]
- [Requirement 2]

## Constraints

- [Constraint 1]
- [Constraint 2]

## Requirements Traceability

| Requirement | Source | Related Artifacts |
|-------------|--------|-------------------|
| [ID] | [Source] | [Artifacts] |

---
*This is a template. Replace with project-specific content.*
`
}

// Knowledge directory templates

func (r *Runtime) generateKnowledgeReadme() string {
	return `# Knowledge Directory

This directory contains all collected engineering knowledge for the project.

## Structure

~~~
knowledge/
├── general/         # General engineering knowledge
├── operational/     # Operational knowledge (runbooks, procedures)
└── developmental/   # Development knowledge (patterns, practices)
~~~

## Purpose

Knowledge artifacts provide context and guidance for engineering decisions.
They are collected during the Knowledge Phase and referenced throughout the project.

## Engineering Rule

DO NOT create new subdirectories in this directory.
Only add files to the existing subdirectories.

## Categories

### general/
General engineering principles, patterns, and terminology.

### operational/
Runbooks, operational procedures, and incident response knowledge.

### developmental/
Development patterns, coding practices, and technical decisions.

---
*This is a template. Populate with project-specific knowledge.*
`
}

func (r *Runtime) generateKnowledgeSubReadme(subdir, description string) string {
	return `# ` + description + `

**Category:** ` + subdir + `
**Created:** ` + time.Now().Format("2006-01-02") + `

## Purpose

` + description + `

## Contents

[Add knowledge artifacts here]

## Guidelines

- Use clear, descriptive filenames
- Include metadata (author, date, source)
- Link related artifacts

---
*This is a template. Populate with project-specific content.*
`
}

// Laboratory directory templates

func (r *Runtime) generateLaboratoryReadme() string {
	return `# Laboratory Directory

This directory contains the experimental and validation workspace.

## Structure

~~~
laboratory/
├── experiments/    # Hypothesis-driven experiments
└── reports/       # Laboratory reports and validation results
~~~

## Purpose

The laboratory provides a controlled environment for:
- Validating architectural decisions
- Testing implementation approaches
- Proving hypotheses
- Documenting experiments

## Engineering Rule

DO NOT create new subdirectories in this directory.
Only add files to the existing subdirectories.

## Workflow

1. Create experiment with clear hypothesis
2. Document methodology
3. Execute and record results
4. Analyze and conclude
5. Generate report

---
*This is a template. Populate with experimental results.*
`
}

func (r *Runtime) generateLabSubReadme(subdir, description string) string {
	return `# ` + description + `

**Category:** ` + subdir + `
**Created:** ` + time.Now().Format("2006-01-02") + `

## Purpose

` + description + `

## Guidelines

- Name experiments with clear identifiers (EXP-001, LAB-001, etc.)
- Include hypothesis, methodology, and results
- Reference related foundation documents

---
*This is a template. Populate with experimental content.*
`
}

// Evidence directory templates

func (r *Runtime) generateEvidenceReadme() string {
	return `# Evidence Directory

This directory contains all verification evidence for the project.

## Purpose

Evidence artifacts prove that engineering work meets requirements:
- Test results
- Performance metrics
- Validation reports
- Verification artifacts

## Engineering Rule

DO NOT create new files in this directory.
Only add verification evidence during the Verification Phase.

## Evidence Types

| Type | Description |
|------|-------------|
| test-results | Test execution results |
| metrics | Performance and quality metrics |
| validation | Formal validation artifacts |
| verification | Verification evidence |

---
*This is a template. Populate with verification evidence.*
`
}

// Reports directory templates

func (r *Runtime) generateReportsReadme() string {
	return `# Reports Directory

This directory contains all generated reports for the project.

## Purpose

Reports document engineering activities and outcomes:
- Status reports
- Progress reports
- Final reports
- Audit reports

## Engineering Rule

DO NOT create new files in this directory.
Only add reports during the Documentation Phase.

## Report Types

| Type | Frequency | Audience |
|------|-----------|----------|
| Status | Weekly | Team |
| Progress | Phase-based | Stakeholders |
| Final | Project end | All |
| Audit | As needed | Compliance |

---
*This is a template. Populate with project reports.*
`
}

// References directory templates

func (r *Runtime) generateReferencesReadme() string {
	return `# References Directory

This directory contains external reference materials for the project.

## Purpose

Reference materials provide context and guidance:
- External documentation
- Standards and specifications
- External research
- Training materials

## Engineering Rule

DO NOT create new files in this directory.
Only add reference materials when relevant to the project.

## Categories

| Category | Description |
|----------|-------------|
| standards | Industry standards |
| docs | External documentation |
| research | Research papers and articles |
| training | Training and learning materials |

---
*This is a template. Populate with reference materials.*
`
}

// Traceability directory templates

func (r *Runtime) generateTraceabilityReadme() string {
	return `# Traceability Directory

This directory contains requirement traceability artifacts.

## Purpose

Traceability links requirements to implementation:
- Requirement → Design mappings
- Design → Implementation mappings
- Implementation → Test mappings

## Engineering Rule

DO NOT create new files in this directory.
Only update traceability matrices during the project.

## Traceability Matrix

| Requirement | Design | Implementation | Test |
|-------------|--------|----------------|------|
| [Req ID] | [Design ref] | [Impl ref] | [Test ref] |

---
*This is a template. Populate with traceability mappings.*
`
}

// Artifacts directory templates

func (r *Runtime) generateArtifactsReadme() string {
	return `# Artifacts Directory

This directory contains project artifacts and deliverables.

## Purpose

Artifacts are the tangible outputs of engineering work:
- Source code
- Configuration files
- Build artifacts
- Deliverables

## Engineering Rule

DO NOT create new files in this directory.
Only add artifacts during Implementation Phase.

## Categories

| Category | Description |
|----------|-------------|
| source | Source code |
| config | Configuration files |
| builds | Build artifacts |
| docs | Generated documentation |

---
*This is a template. Populate with project artifacts.*
`
}

// Generic README generator

func (r *Runtime) generateDefaultReadme() string {
	return `# Directory

**Created:** ` + time.Now().Format("2006-01-02") + `

## Purpose

[Describe the purpose of this directory]

## Contents

[List contents]

## Guidelines

[Usage guidelines]

---
*This is a template. Replace with specific content.*
`
}

func (r *Runtime) generateReadmeForDir(dir string) string {
	switch dir {
	case "runtime":
		return `# Runtime Directory

This directory contains runtime execution state and logs.
`
	case "config":
		return `# Config Directory

This directory contains runtime configuration files.
`
	case "state":
		return `# State Directory

This directory contains session and runtime state.
`
	case "sessions":
		return `# Sessions Directory

This directory contains session history.
`
	case "normalized":
		return `# Normalized Directory

This directory contains normalized documentation output.
`
	case "cache":
		return `# Cache Directory

This directory contains cached computations.
`
	case "someday":
		return `# Someday/Maybe Directory

This directory contains deferred engineering ideas.

## Structure

~~~
someday/
├── ideas/       # Active someday ideas
├── archived/    # Archived ideas
└── promoted/   # Promoted ideas
~~~
`
	default:
		return r.generateDefaultReadme()
	}
}

// verifyAllArtifacts verifies every created artifact
func (r *Runtime) verifyAllArtifacts(result *InitializeResult) {
	// Verify directories
	directories := []string{
		DirRuntime, DirFoundation, DirKnowledge, DirLaboratory,
		DirEvidence, DirReferences, DirTraceability, DirReports,
		DirConfig, DirState, DirArtifacts, DirSessions,
		DirNormalized, DirCache, DirSomeday,
	}

	for _, dir := range directories {
		path := filepath.Join(r.kdsePath, dir)
		result.Verification = append(result.Verification, r.verifyDirectory(path, dir))
	}

	// Verify knowledge subdirectories
	knowledgeSubdirs := []string{"general", "operational", "developmental"}
	for _, subdir := range knowledgeSubdirs {
		path := filepath.Join(r.kdsePath, DirKnowledge, subdir)
		result.Verification = append(result.Verification, r.verifyDirectory(path, DirKnowledge+"/"+subdir))
	}

	// Verify laboratory subdirectories
	labSubdirs := []string{"experiments", "reports"}
	for _, subdir := range labSubdirs {
		path := filepath.Join(r.kdsePath, DirLaboratory, subdir)
		result.Verification = append(result.Verification, r.verifyDirectory(path, DirLaboratory+"/"+subdir))
	}

	// Verify someday subdirectories
	somedaySubdirs := []string{"ideas", "archived", "promoted"}
	for _, subdir := range somedaySubdirs {
		path := filepath.Join(r.kdsePath, DirSomeday, subdir)
		result.Verification = append(result.Verification, r.verifyDirectory(path, DirSomeday+"/"+subdir))
	}

	// Verify core runtime files
	files := []string{
		FileManifest, FileSessionState, FileRuntimeConfig,
		FileKnowledgeIndex, FileArtifactIndex,
	}

	for _, file := range files {
		path := filepath.Join(r.kdsePath, file)
		result.Verification = append(result.Verification, r.verifyFile(path, file))
	}

	// Verify foundation templates
	foundationTemplates := []string{"README.md", "PROBLEM.md", "SPEC.md", "ARCHITECTURE.md", "ASSUMPTIONS.md", "REQUIREMENTS.md"}
	for _, tpl := range foundationTemplates {
		path := filepath.Join(r.kdsePath, DirFoundation, tpl)
		result.Verification = append(result.Verification, r.verifyFile(path, DirFoundation+"/"+tpl))
	}

	// Verify directory READMEs
	readmeFiles := []string{
		DirKnowledge + "/README.md",
		DirLaboratory + "/README.md",
		DirEvidence + "/README.md",
		DirReports + "/README.md",
		DirReferences + "/README.md",
		DirTraceability + "/README.md",
		DirArtifacts + "/README.md",
	}
	for _, readme := range readmeFiles {
		path := filepath.Join(r.kdsePath, readme)
		result.Verification = append(result.Verification, r.verifyFile(path, readme))
	}

	// Verify knowledge subdirectory READMEs
	for _, subdir := range knowledgeSubdirs {
		path := filepath.Join(r.kdsePath, DirKnowledge, subdir, "README.md")
		result.Verification = append(result.Verification, r.verifyFile(path, DirKnowledge+"/"+subdir+"/README.md"))
	}

	// Verify laboratory subdirectory READMEs
	for _, subdir := range labSubdirs {
		path := filepath.Join(r.kdsePath, DirLaboratory, subdir, "README.md")
		result.Verification = append(result.Verification, r.verifyFile(path, DirLaboratory+"/"+subdir+"/README.md"))
	}

	// Verify someday manifest
	somedayManifestPath := filepath.Join(r.kdsePath, DirSomeday, "someday.yaml")
	result.Verification = append(result.Verification, r.verifyFile(somedayManifestPath, DirSomeday+"/someday.yaml"))

	// Add evidence for successful verifications
	for _, v := range result.Verification {
		if v.Status == "PASS" {
			result.Evidence = append(result.Evidence, fmt.Sprintf("%s: %s", v.Artifact, v.Path))
		}
	}
}

// verifyDirectory checks if a directory exists and is accessible
func (r *Runtime) verifyDirectory(path, name string) VerificationResult {
	info, err := os.Stat(path)
	if err != nil {
		return VerificationResult{
			Artifact:  name,
			Path:      path,
			Status:    "FAIL",
			Error:     fmt.Sprintf("Directory does not exist: %v", err),
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	if !info.IsDir() {
		return VerificationResult{
			Artifact:  name,
			Path:      path,
			Status:    "FAIL",
			Error:     "Path exists but is not a directory",
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	return VerificationResult{
		Artifact:  name,
		Path:      path,
		Status:    "PASS",
		Evidence:  fmt.Sprintf("Directory exists, mode: %o", info.Mode()),
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// verifyFile checks if a file exists and is readable
func (r *Runtime) verifyFile(path, name string) VerificationResult {
	info, err := os.Stat(path)
	if err != nil {
		return VerificationResult{
			Artifact:  name,
			Path:      path,
			Status:    "FAIL",
			Error:     fmt.Sprintf("File does not exist: %v", err),
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	if info.IsDir() {
		return VerificationResult{
			Artifact:  name,
			Path:      path,
			Status:    "FAIL",
			Error:     "Path exists but is a directory, not a file",
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	// Verify file is readable
	content, err := os.ReadFile(path)
	if err != nil {
		return VerificationResult{
			Artifact:  name,
			Path:      path,
			Status:    "FAIL",
			Error:     fmt.Sprintf("File exists but is not readable: %v", err),
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	return VerificationResult{
		Artifact:  name,
		Path:      path,
		Status:    "PASS",
		Evidence:  fmt.Sprintf("File readable, size: %d bytes", len(content)),
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// calculateConfidence calculates verification confidence
func (r *Runtime) calculateConfidence(verifications []VerificationResult) float64 {
	if len(verifications) == 0 {
		return 0.0
	}

	passed := 0
	for _, v := range verifications {
		if v.Status == "PASS" {
			passed++
		}
	}

	return float64(passed) / float64(len(verifications))
}

// determineSuccess checks if all required artifacts passed
func (r *Runtime) determineSuccess(result *InitializeResult) bool {
	for _, v := range result.Verification {
		if v.Status == "FAIL" {
			return false
		}
	}
	return len(result.Errors) == 0
}

// Verify performs a complete runtime verification
func (r *Runtime) Verify() *VerificationReport {
	report := &VerificationReport{
		Components: []VerificationResult{},
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	// Check workspace
	report.Components = append(report.Components, r.verifyWorkspace())

	// Check main directories
	mainDirs := []string{DirRuntime, DirFoundation, DirKnowledge, DirLaboratory, DirEvidence, DirReferences, DirTraceability, DirReports, DirConfig, DirState, DirArtifacts, DirSomeday}
	for _, dir := range mainDirs {
		report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, dir), dir))
	}

	// Check knowledge subdirectories
	knowledgeSubs := []string{"general", "operational", "developmental"}
	for _, sub := range knowledgeSubs {
		report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, DirKnowledge, sub), DirKnowledge+"/"+sub))
	}

	// Check laboratory subdirectories
	labSubs := []string{"experiments", "reports"}
	for _, sub := range labSubs {
		report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, DirLaboratory, sub), DirLaboratory+"/"+sub))
	}

	// Check someday subdirectories
	somedaySubs := []string{"ideas", "archived", "promoted"}
	for _, sub := range somedaySubs {
		report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, DirSomeday, sub), DirSomeday+"/"+sub))
	}

	// Check foundation templates
	foundationTemplates := []string{"PROBLEM.md", "SPEC.md", "ARCHITECTURE.md", "ASSUMPTIONS.md", "REQUIREMENTS.md", "README.md"}
	for _, tpl := range foundationTemplates {
		report.Components = append(report.Components, r.verifyFile(filepath.Join(r.kdsePath, DirFoundation, tpl), DirFoundation+"/"+tpl))
	}

	// Check directory READMEs
	readmeDirs := []string{DirKnowledge, DirLaboratory, DirEvidence, DirReports, DirReferences, DirTraceability, DirArtifacts}
	for _, dir := range readmeDirs {
		report.Components = append(report.Components, r.verifyFile(filepath.Join(r.kdsePath, dir, "README.md"), dir+"/README.md"))
	}

	// Check knowledge subdirectory READMEs
	for _, sub := range knowledgeSubs {
		report.Components = append(report.Components, r.verifyFile(filepath.Join(r.kdsePath, DirKnowledge, sub, "README.md"), DirKnowledge+"/"+sub+"/README.md"))
	}

	// Check laboratory subdirectory READMEs
	for _, sub := range labSubs {
		report.Components = append(report.Components, r.verifyFile(filepath.Join(r.kdsePath, DirLaboratory, sub, "README.md"), DirLaboratory+"/"+sub+"/README.md"))
	}

	// Check core runtime files
	coreFiles := []string{FileManifest, FileSessionState, FileRuntimeConfig, FileKnowledgeIndex, FileArtifactIndex}
	for _, file := range coreFiles {
		report.Components = append(report.Components, r.verifyFile(filepath.Join(r.kdsePath, file), file))
	}

	// Check someday manifest
	report.Components = append(report.Components, r.verifyFile(filepath.Join(r.kdsePath, DirSomeday, "someday.yaml"), DirSomeday+"/someday.yaml"))

	// Calculate confidence
	report.Confidence = r.calculateConfidence(report.Components)

	// Determine overall status
	report.Success = r.determineVerificationSuccess(report)

	return report
}

// verifyWorkspace checks if the .kdse workspace exists
func (r *Runtime) verifyWorkspace() VerificationResult {
	info, err := os.Stat(r.kdsePath)
	if err != nil {
		return VerificationResult{
			Artifact:  "Workspace",
			Path:      r.kdsePath,
			Status:    "FAIL",
			Error:     "Workspace does not exist",
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	if !info.IsDir() {
		return VerificationResult{
			Artifact:  "Workspace",
			Path:      r.kdsePath,
			Status:    "FAIL",
			Error:     "Workspace path exists but is not a directory",
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	return VerificationResult{
		Artifact:  "Workspace",
		Path:      r.kdsePath,
		Status:    "PASS",
		Evidence:  "Workspace directory exists",
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// determineVerificationSuccess checks if all components passed
func (r *Runtime) determineVerificationSuccess(report *VerificationReport) bool {
	for _, c := range report.Components {
		if c.Status == "FAIL" {
			report.Failed = append(report.Failed, c.Artifact)
		}
	}
	return len(report.Failed) == 0
}

// CheckInvariant verifies if a phase transition is allowed
func (r *Runtime) CheckInvariant(phase, required string) (bool, string) {
	return r.invariant.Check(phase, required)
}

// DefaultManifest returns the default runtime manifest
func DefaultManifest() *RuntimeManifest {
	return &RuntimeManifest{
		Version:   "1.0.0",
		CreatedAt: time.Now().Format(time.RFC3339),
		Directories: []ManifestDirectory{
			{Path: DirRuntime, Required: true, Purpose: "Runtime execution state and logs"},
			{Path: DirFoundation, Required: true, Purpose: "Project foundation documents"},
			{Path: DirKnowledge, Required: true, Purpose: "Collected engineering knowledge"},
			{Path: DirLaboratory, Required: true, Purpose: "Testing and experimentation"},
			{Path: DirEvidence, Required: true, Purpose: "Evidence artifacts"},
			{Path: DirReferences, Required: true, Purpose: "Reference materials"},
			{Path: DirTraceability, Required: true, Purpose: "Requirement traceability"},
			{Path: DirReports, Required: true, Purpose: "Generated reports"},
			{Path: DirConfig, Required: true, Purpose: "Runtime configuration"},
			{Path: DirState, Required: true, Purpose: "Session and runtime state"},
			{Path: DirArtifacts, Required: true, Purpose: "Artifact inventory"},
			{Path: DirSessions, Required: true, Purpose: "Session history"},
			{Path: DirNormalized, Required: true, Purpose: "Normalized documentation"},
			{Path: DirCache, Required: true, Purpose: "Cached computations"},
			{Path: DirSomeday, Required: true, Purpose: "Someday/Maybe knowledge repository"},
		},
		Files: []ManifestFile{
			{Path: FileManifest, Required: true, Purpose: "Runtime manifest definition"},
			{Path: FileSessionState, Required: true, Purpose: "Current session state"},
			{Path: FileRuntimeConfig, Required: true, Purpose: "Runtime configuration"},
			{Path: FileKnowledgeIndex, Required: true, Purpose: "Knowledge artifact index"},
			{Path: FileArtifactIndex, Required: true, Purpose: "Artifact inventory"},
		},
		Invariants: DefaultInvariants(),
	}
}

// DefaultInvariants returns the default runtime invariants
func DefaultInvariants() []RuntimeInvariant {
	return []RuntimeInvariant{
		{Phase: "Problem", Requires: []string{"Runtime initialized"}, Description: "Runtime must be initialized before problem phase"},
		{Phase: "Foundation", Requires: []string{"Runtime initialized"}, Description: "Foundation requires initialized runtime"},
		{Phase: "Knowledge", Requires: []string{"Foundation exists"}, Description: "Knowledge collection requires foundation"},
		{Phase: "Architecture", Requires: []string{"Knowledge collected"}, Description: "Architecture requires knowledge"},
		{Phase: "Implementation", Requires: []string{"Architecture approved"}, Description: "Implementation requires approved architecture"},
		{Phase: "Verification", Requires: []string{"Implementation complete"}, Description: "Verification requires implementation"},
		{Phase: "Documentation", Requires: []string{"Verification complete"}, Description: "Documentation requires verification"},
		{Phase: "Audit", Requires: []string{"All phases complete"}, Description: "Audit requires all phases complete"},
	}
}

// File generation templates

func (r *Runtime) generateManifestContent() string {
	m := r.manifest
	m.CreatedAt = time.Now().Format(time.RFC3339)
	data, _ := json.MarshalIndent(m, "", "  ")
	return string(data)
}

func (r *Runtime) generateSessionStateContent() string {
	state := map[string]interface{}{
		"version":       "1.0.0",
		"session_id":    fmt.Sprintf("KDSE-SESSION-%s", time.Now().Format("20060102-150405")),
		"status":        "Initialized",
		"phase":         "Problem",
		"confidence":    0.0,
		"evidence":      []string{},
		"created_at":    time.Now().Format(time.RFC3339),
		"last_verified": nil,
	}
	data, _ := json.MarshalIndent(state, "", "  ")
	return string(data)
}

func (r *Runtime) generateRuntimeConfigContent() string {
	config := map[string]interface{}{
		"version":               "1.0.0",
		"runtime":               "evidence-driven",
		"strict_mode":           true,
		"confidence_threshold":  0.7,
		"evidence_threshold":    0.6,
		"max_cycles":            100,
		"auto_verify":           true,
		"enforce_invariants":    true,
		"created_at":            time.Now().Format(time.RFC3339),
	}
	data, _ := json.MarshalIndent(config, "", "  ")
	return string(data)
}

func (r *Runtime) generateKnowledgeIndexContent() string {
	index := map[string]interface{}{
		"version":       "1.0.0",
		"last_updated":  time.Now().Format(time.RFC3339),
		"artifacts":     []map[string]interface{}{},
		"categories": map[string]int{
			"architecture":   0,
			"design":         0,
			"implementation": 0,
			"testing":        0,
			"documentation":  0,
		},
		"total_count": 0,
	}
	data, _ := json.MarshalIndent(index, "", "  ")
	return string(data)
}

func (r *Runtime) generateArtifactIndexContent() string {
	index := map[string]interface{}{
		"version":      "1.0.0",
		"last_updated": time.Now().Format(time.RFC3339),
		"artifacts":    []map[string]interface{}{},
		"categories": map[string]int{
			"foundation":    0,
			"evidence":      0,
			"reference":     0,
			"traceability":  0,
			"report":        0,
		},
		"total_count": 0,
	}
	data, _ := json.MarshalIndent(index, "", "  ")
	return string(data)
}

// generateSomedayManifestContent generates the initial someday manifest
func (r *Runtime) generateSomedayManifestContent() string {
	manifest := map[string]interface{}{
		"version":      "1.0.0",
		"created_at":   time.Now().Format(time.RFC3339),
		"last_updated": time.Now().Format(time.RFC3339),
		"ideas":        []string{},
		"total_count":  0,
		"by_status": map[string]int{
			"SOMEDAY":     0,
			"LABORATORY":  0,
			"PROMOTED":    0,
			"IMPLEMENTED": 0,
			"ARCHIVED":    0,
			"REJECTED":    0,
		},
		"by_priority": map[string]int{
			"priority-1": 0,
			"priority-2": 0,
			"priority-3": 0,
			"priority-4": 0,
			"priority-5": 0,
		},
		"next_idea_id": 1,
	}
	data, _ := json.MarshalIndent(manifest, "", "  ")
	return string(data)
}

// FormatInitializeResult formats initialization result for display
func FormatInitializeResult(result *InitializeResult) string {
	var output string

	output += "╔═══════════════════════════════════════════════════════════════╗\n"
	output += "║           KDSE Evidence-Driven Initialization                 ║\n"
	output += "╠═══════════════════════════════════════════════════════════════╣\n"
	output += fmt.Sprintf("║ Workspace: %s\n", result.WorkspacePath)
	output += fmt.Sprintf("║ Confidence: %.2f\n", result.Confidence)
	output += fmt.Sprintf("║ Status: %s\n", boolToStatus(result.Success))
	output += "╠═══════════════════════════════════════════════════════════════╣\n"
	output += "║ Verification Results                                         ║\n"

	for _, v := range result.Verification {
		statusIcon := "✓"
		if v.Status == "FAIL" {
			statusIcon = "✗"
		}
		output += fmt.Sprintf("║ %s %-12s %s\n", statusIcon, v.Artifact, v.Status)
	}

	output += "╠═══════════════════════════════════════════════════════════════╣"
	if result.Success {
		output += "\n║ Status: OPERATIONAL                                          ║"
	} else {
		output += "\n║ Status: FAILED                                               ║"
		if len(result.Errors) > 0 {
			output += "\n╠═══════════════════════════════════════════════════════════════╣"
			output += "\n║ Errors                                                       ║"
			for _, err := range result.Errors {
				output += fmt.Sprintf("\n║   • %s\n", err)
			}
		}
	}

	output += "\n╚═══════════════════════════════════════════════════════════════╝\n"
	return output
}

// FormatVerificationReport formats verification report for display
func FormatVerificationReport(report *VerificationReport) string {
	var output string

	output += "╔═══════════════════════════════════════════════════════════════╗\n"
	output += "║              KDSE Runtime Self-Audit                         ║\n"
	output += "╠═══════════════════════════════════════════════════════════════╣"

	for _, c := range report.Components {
		statusIcon := "PASS"
		if c.Status == "FAIL" {
			statusIcon = "FAIL"
		}
		output += fmt.Sprintf("║ %-12s %-8s %s\n", c.Artifact, statusIcon, c.Path)
	}

	output += "╠═══════════════════════════════════════════════════════════════╣"
	output += fmt.Sprintf("║ Confidence: %.2f\n", report.Confidence)

	if report.Success {
		output += "║ Status: OPERATIONAL                                          ║\n"
	} else {
		output += "║ Status: FAILED                                               ║\n"
		if len(report.Failed) > 0 {
			output += "║ Failed Components:                                           ║\n"
			for _, f := range report.Failed {
				output += fmt.Sprintf("║   • %s\n", f)
			}
		}
	}

	output += "╚═══════════════════════════════════════════════════════════════╝\n"
	return output
}

func boolToStatus(b bool) string {
	if b {
		return "SUCCESS"
	}
	return "FAILED"
}
