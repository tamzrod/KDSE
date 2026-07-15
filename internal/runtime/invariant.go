package runtime

import (
	"fmt"
	"os"
	"path/filepath"
)

// InvariantEngine manages runtime phase invariants
type InvariantEngine struct {
	invariants map[string]InvariantCheck
	kdsePath   string
}

// InvariantCheck defines how to verify an invariant
type InvariantCheck struct {
	Description string
	Verify      func(kdsePath string) (bool, string)
}

// NewInvariantEngine creates a new invariant engine
func NewInvariantEngine() *InvariantEngine {
	engine := &InvariantEngine{
		invariants: make(map[string]InvariantCheck),
	}
	engine.registerDefaultInvariants()
	return engine
}

// SetKDSEPath sets the KDSE workspace path for verification
func (e *InvariantEngine) SetKDSEPath(path string) {
	e.kdsePath = path
}

// registerDefaultInvariants registers the standard KDSE invariants
func (e *InvariantEngine) registerDefaultInvariants() {
	e.invariants["Runtime initialized"] = InvariantCheck{
		Description: "KDSE workspace must exist with manifest",
		Verify: func(kdsePath string) (bool, string) {
			manifestPath := filepath.Join(kdsePath, FileManifest)
			if _, err := os.Stat(manifestPath); err != nil {
				return false, fmt.Sprintf("Manifest not found: %v", err)
			}
			return true, "Runtime initialized"
		},
	}

	e.invariants["Foundation exists"] = InvariantCheck{
		Description: "Foundation directory must exist",
		Verify: func(kdsePath string) (bool, string) {
			foundationPath := filepath.Join(kdsePath, DirFoundation)
			info, err := os.Stat(foundationPath)
			if err != nil {
				return false, fmt.Sprintf("Foundation directory not found: %v", err)
			}
			if !info.IsDir() {
				return false, "Foundation path exists but is not a directory"
			}
			return true, "Foundation exists"
		},
	}

	e.invariants["Knowledge collected"] = InvariantCheck{
		Description: "Knowledge directory must have at least one artifact",
		Verify: func(kdsePath string) (bool, string) {
			knowledgePath := filepath.Join(kdsePath, DirKnowledge)
			entries, err := os.ReadDir(knowledgePath)
			if err != nil {
				return false, fmt.Sprintf("Knowledge directory not accessible: %v", err)
			}
			count := 0
			for _, entry := range entries {
				if !entry.IsDir() {
					count++
				}
			}
			if count == 0 {
				return false, "No knowledge artifacts collected"
			}
			return true, fmt.Sprintf("Knowledge collected: %d artifacts", count)
		},
	}

	e.invariants["Architecture approved"] = InvariantCheck{
		Description: "Architecture document must exist in foundation",
		Verify: func(kdsePath string) (bool, string) {
			archPath := filepath.Join(kdsePath, DirFoundation, "architecture.md")
			if _, err := os.Stat(archPath); err != nil {
				return false, "Architecture document not found"
			}
			return true, "Architecture approved"
		},
	}

	e.invariants["Implementation complete"] = InvariantCheck{
		Description: "Implementation artifacts must exist",
		Verify: func(kdsePath string) (bool, string) {
			// Check for any code artifacts
			artifactsPath := filepath.Join(kdsePath, DirArtifacts)
			if _, err := os.Stat(artifactsPath); err != nil {
				return false, "Artifacts directory not found"
			}
			entries, err := os.ReadDir(artifactsPath)
			if err != nil {
				return false, fmt.Sprintf("Artifacts directory not accessible: %v", err)
			}
			count := 0
			for _, entry := range entries {
				if !entry.IsDir() {
					count++
				}
			}
			if count == 0 {
				return false, "No implementation artifacts found"
			}
			return true, fmt.Sprintf("Implementation artifacts: %d", count)
		},
	}

	e.invariants["Verification complete"] = InvariantCheck{
		Description: "Verification report must exist",
		Verify: func(kdsePath string) (bool, string) {
			reportsPath := filepath.Join(kdsePath, DirReports)
			if _, err := os.Stat(reportsPath); err != nil {
				return false, "Reports directory not found"
			}
			entries, err := os.ReadDir(reportsPath)
			if err != nil {
				return false, fmt.Sprintf("Reports directory not accessible: %v", err)
			}
			count := 0
			for _, entry := range entries {
				if !entry.IsDir() {
					count++
				}
			}
			if count == 0 {
				return false, "No verification reports found"
			}
			return true, fmt.Sprintf("Verification reports: %d", count)
		},
	}

	e.invariants["All phases complete"] = InvariantCheck{
		Description: "All required directories must exist",
		Verify: func(kdsePath string) (bool, string) {
			requiredDirs := []string{
				DirFoundation, DirKnowledge, DirLaboratory,
				DirEvidence, DirReferences, DirTraceability,
				DirReports, DirConfig, DirState, DirArtifacts,
			}
			missing := []string{}
			for _, dir := range requiredDirs {
				path := filepath.Join(kdsePath, dir)
				info, err := os.Stat(path)
				if err != nil || !info.IsDir() {
					missing = append(missing, dir)
				}
			}
			if len(missing) > 0 {
				return false, fmt.Sprintf("Missing directories: %v", missing)
			}
			return true, "All phases complete"
		},
	}
}

// Check verifies if a required invariant is satisfied
func (e *InvariantEngine) Check(phase, required string) (bool, string) {
	if e.kdsePath == "" {
		return false, "KDSE path not set"
	}

	check, exists := e.invariants[required]
	if !exists {
		return false, fmt.Sprintf("Unknown invariant: %s", required)
	}

	return check.Verify(e.kdsePath)
}

// CheckAll verifies multiple invariants
func (e *InvariantEngine) CheckAll(requirements []string) (bool, map[string]string) {
	results := make(map[string]string)
	allPassed := true

	for _, req := range requirements {
		passed, msg := e.Check("", req)
		results[req] = msg
		if !passed {
			allPassed = false
		}
	}

	return allPassed, results
}

// CanTransitionTo checks if transition to a phase is allowed
func (e *InvariantEngine) CanTransitionTo(phase string) (bool, []string) {
	invariants := DefaultInvariants()

	for _, inv := range invariants {
		if inv.Phase == phase {
			allPassed := true
			failed := []string{}

			for _, req := range inv.Requires {
				passed, _ := e.Check("", req)
				if !passed {
					allPassed = false
					failed = append(failed, req)
				}
			}

			return allPassed, failed
		}
	}

	// Unknown phase - allow transition
	return true, []string{}
}

// GetInvariantDescription returns the description for an invariant
func (e *InvariantEngine) GetInvariantDescription(name string) string {
	if check, exists := e.invariants[name]; exists {
		return check.Description
	}
	return ""
}
