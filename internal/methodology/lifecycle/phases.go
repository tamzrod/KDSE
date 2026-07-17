package lifecycle

import (
	"errors"
	"fmt"
)

// Phase represents a KDSE engineering phase
type Phase string

const (
	// PhaseInitialization is the first phase where runtime is created
	PhaseInitialization Phase = "initialization"
	// PhaseKnowledge is where requirements are gathered
	PhaseKnowledge Phase = "knowledge"
	// PhaseArchitecture is where system architecture is designed
	PhaseArchitecture Phase = "architecture"
	// PhaseImplementation is where the system is built
	PhaseImplementation Phase = "implementation"
	// PhaseVerification is where the system is tested
	PhaseVerification Phase = "verification"
	// PhaseReports is where results are documented
	PhaseReports Phase = "reports"
)

// AllPhases returns all phases in order
func AllPhases() []Phase {
	return []Phase{
		PhaseInitialization,
		PhaseKnowledge,
		PhaseArchitecture,
		PhaseImplementation,
		PhaseVerification,
		PhaseReports,
	}
}

// PhaseSpec defines the specification for a phase
type PhaseSpec struct {
	Name                string
	Description         string
	AllowedTransitions  []Phase
	RequiredArtifacts   []ArtifactSpec
	CompletionCriteria  []CompletionCriterion
}

// ArtifactSpec defines a required artifact for a phase
type ArtifactSpec struct {
	Path        string
	Type        string
	Required    bool
	Validators  []ArtifactValidator
}

// ArtifactValidator validates an artifact
type ArtifactValidator func(path string) error

// CompletionCriterion defines a criterion for phase completion
type CompletionCriterion struct {
	Description string
	Validate    func(phase Phase) bool
}

// Errors
var (
	ErrNoNextPhase       = errors.New("no next phase available")
	ErrNoPreviousPhase   = errors.New("no previous phase available")
	ErrInvalidTransition = errors.New("invalid phase transition")
	ErrPhaseNotFound     = errors.New("phase not found")
)

// PhaseRegistry contains all phase specifications
var PhaseRegistry = map[Phase]PhaseSpec{
	PhaseInitialization: {
		Name:       "Initialization",
		Description: "Set up runtime and workspace",
		AllowedTransitions: []Phase{PhaseKnowledge},
		RequiredArtifacts: []ArtifactSpec{
			{Path: ".kdse/runtime.yaml", Type: "yaml", Required: true},
			{Path: ".kdse/workspace.yaml", Type: "yaml", Required: true},
			{Path: ".kdse/methodology.yaml", Type: "yaml", Required: true},
			{Path: ".kdse/phase.yaml", Type: "yaml", Required: true},
			{Path: ".kdse/session.yaml", Type: "yaml", Required: true},
			{Path: ".kdse/metadata.yaml", Type: "yaml", Required: true},
		},
		CompletionCriteria: []CompletionCriterion{
			{Description: "All required files exist", Validate: func(p Phase) bool { return true }},
			{Description: "Runtime version is supported", Validate: func(p Phase) bool { return true }},
		},
	},
	PhaseKnowledge: {
		Name:       "Knowledge",
		Description: "Gather and document requirements",
		AllowedTransitions: []Phase{PhaseArchitecture},
		RequiredArtifacts: []ArtifactSpec{
			{Path: "knowledge/requirements.md", Type: "markdown", Required: true},
			{Path: "knowledge/stakeholders.md", Type: "markdown", Required: true},
			{Path: "knowledge/constraints.md", Type: "markdown", Required: true},
			{Path: "knowledge/glossary.md", Type: "markdown", Required: false},
		},
		CompletionCriteria: []CompletionCriterion{
			{Description: "All required files exist", Validate: func(p Phase) bool { return true }},
			{Description: "Content meets quality standards", Validate: func(p Phase) bool { return true }},
		},
	},
	PhaseArchitecture: {
		Name:       "Architecture",
		Description: "Design system architecture",
		AllowedTransitions: []Phase{PhaseImplementation},
		RequiredArtifacts: []ArtifactSpec{
			{Path: "architecture/architecture.md", Type: "markdown", Required: true},
			{Path: "architecture/decisions.md", Type: "markdown", Required: true},
			{Path: "architecture/components.md", Type: "markdown", Required: true},
			{Path: "architecture/interfaces.md", Type: "markdown", Required: false},
		},
		CompletionCriteria: []CompletionCriterion{
			{Description: "Architecture addresses all requirements", Validate: func(p Phase) bool { return true }},
			{Description: "Decisions are documented", Validate: func(p Phase) bool { return true }},
		},
	},
	PhaseImplementation: {
		Name:       "Implementation",
		Description: "Build the system",
		AllowedTransitions: []Phase{PhaseVerification},
		RequiredArtifacts: []ArtifactSpec{
			{Path: "implementation/implementation.md", Type: "markdown", Required: true},
			{Path: "implementation/testing.md", Type: "markdown", Required: true},
			{Path: "implementation/deployment.md", Type: "markdown", Required: false},
		},
		CompletionCriteria: []CompletionCriterion{
			{Description: "Implementation matches architecture", Validate: func(p Phase) bool { return true }},
			{Description: "Testing approach is defined", Validate: func(p Phase) bool { return true }},
		},
	},
	PhaseVerification: {
		Name:       "Verification",
		Description: "Test and validate the implementation",
		AllowedTransitions: []Phase{PhaseReports},
		RequiredArtifacts: []ArtifactSpec{
			{Path: "verification/verification.md", Type: "markdown", Required: true},
			{Path: "verification/test-results.md", Type: "markdown", Required: true},
			{Path: "verification/coverage.md", Type: "markdown", Required: true},
		},
		CompletionCriteria: []CompletionCriterion{
			{Description: "All tests pass", Validate: func(p Phase) bool { return true }},
			{Description: "Coverage meets target", Validate: func(p Phase) bool { return true }},
		},
	},
	PhaseReports: {
		Name:       "Reports",
		Description: "Document results and recommendations",
		AllowedTransitions: []Phase{}, // Terminal phase
		RequiredArtifacts: []ArtifactSpec{
			{Path: "reports/summary.md", Type: "markdown", Required: true},
			{Path: "reports/findings.md", Type: "markdown", Required: true},
			{Path: "reports/recommendations.md", Type: "markdown", Required: true},
		},
		CompletionCriteria: []CompletionCriterion{
			{Description: "Summary is clear", Validate: func(p Phase) bool { return true }},
			{Description: "Recommendations are actionable", Validate: func(p Phase) bool { return true }},
		},
	},
}

// GetPhaseSpec returns the specification for a phase
func GetPhaseSpec(phase Phase) (PhaseSpec, error) {
	spec, ok := PhaseRegistry[phase]
	if !ok {
		return PhaseSpec{}, fmt.Errorf("%w: %s", ErrPhaseNotFound, phase)
	}
	return spec, nil
}

// GetPhaseName returns the display name for a phase
func GetPhaseName(phase Phase) string {
	spec, ok := PhaseRegistry[phase]
	if !ok {
		return string(phase)
	}
	return spec.Name
}

// GetPhaseDescription returns the description for a phase
func GetPhaseDescription(phase Phase) string {
	spec, ok := PhaseRegistry[phase]
	if !ok {
		return ""
	}
	return spec.Description
}

// GetRequiredArtifacts returns required artifacts for a phase
func GetRequiredArtifacts(phase Phase) []ArtifactSpec {
	spec, ok := PhaseRegistry[phase]
	if !ok {
		return nil
	}
	return spec.RequiredArtifacts
}

// IsValidTransition checks if a phase transition is valid
func IsValidTransition(from, to Phase) bool {
	spec, ok := PhaseRegistry[from]
	if !ok {
		return false
	}
	for _, allowed := range spec.AllowedTransitions {
		if allowed == to {
			return true
		}
	}
	return false
}

// GetNextPhase returns the valid next phase
func GetNextPhase(current Phase) (Phase, error) {
	spec, ok := PhaseRegistry[current]
	if !ok {
		return "", ErrPhaseNotFound
	}
	if len(spec.AllowedTransitions) == 0 {
		return "", ErrNoNextPhase
	}
	return spec.AllowedTransitions[0], nil
}

// GetPreviousPhase returns the previous phase
func GetPreviousPhase(current Phase) (Phase, error) {
	phases := AllPhases()
	for i, phase := range phases {
		if phase == current && i > 0 {
			return phases[i-1], nil
		}
	}
	return "", ErrNoPreviousPhase
}

// IsTerminalPhase checks if a phase is terminal
func IsTerminalPhase(phase Phase) bool {
	spec, ok := PhaseRegistry[phase]
	if !ok {
		return false
	}
	return len(spec.AllowedTransitions) == 0
}

// Index returns the index of a phase in the sequence
func Index(phase Phase) int {
	phases := AllPhases()
	for i, p := range phases {
		if p == phase {
			return i
		}
	}
	return -1
}
