package lifecycle

import (
	"context"
)

// Lifecycle defines the methodology's phase lifecycle management
// This interface is implemented by the methodology package and consumed by the workspace engine
type Lifecycle interface {
	// GetPhases returns all phases in order
	GetPhases() []Phase

	// GetPhaseName returns the display name for a phase
	GetPhaseName(phase Phase) string

	// GetPhaseDescription returns the description for a phase
	GetPhaseDescription(phase Phase) string

	// GetRequiredArtifacts returns required artifacts for a phase
	GetRequiredArtifacts(phase Phase) []ArtifactSpec

	// IsValidTransition checks if a phase transition is valid
	IsValidTransition(from, to Phase) bool

	// GetNextPhase returns the valid next phase
	GetNextPhase(current Phase) (Phase, error)

	// GetPreviousPhase returns the previous phase
	GetPreviousPhase(current Phase) (Phase, error)

	// GetCompletionCriteria returns criteria for phase completion
	GetCompletionCriteria(phase Phase) []CompletionCriterion

	// IsTerminalPhase checks if a phase is terminal
	IsTerminalPhase(phase Phase) bool

	// ValidatePhase validates a phase and its artifacts
	ValidatePhase(ctx context.Context, phase Phase, artifactPaths []string) (*PhaseValidation, error)
}

// PhaseValidation contains the result of phase validation
type PhaseValidation struct {
	Phase         Phase
	Valid         bool
	Verified      []ArtifactSpec
	Missing       []ArtifactSpec
	Invalid       []ArtifactSpec
	CriteriaMet   []CompletionCriterion
	CriteriaUnmet []CompletionCriterion
	Errors        []ValidationError
}

// ValidationError contains error details
type ValidationError struct {
	Type    string
	Message string
	Path    string
}

// DefaultLifecycle implements Lifecycle interface
type DefaultLifecycle struct{}

// NewLifecycle creates a new default lifecycle
func NewLifecycle() *DefaultLifecycle {
	return &DefaultLifecycle{}
}

// GetPhases returns all phases in order
func (l *DefaultLifecycle) GetPhases() []Phase {
	return AllPhases()
}

// GetPhaseName returns the display name for a phase
func (l *DefaultLifecycle) GetPhaseName(phase Phase) string {
	return GetPhaseName(phase)
}

// GetPhaseDescription returns the description for a phase
func (l *DefaultLifecycle) GetPhaseDescription(phase Phase) string {
	return GetPhaseDescription(phase)
}

// GetRequiredArtifacts returns required artifacts for a phase
func (l *DefaultLifecycle) GetRequiredArtifacts(phase Phase) []ArtifactSpec {
	return GetRequiredArtifacts(phase)
}

// IsValidTransition checks if a phase transition is valid
func (l *DefaultLifecycle) IsValidTransition(from, to Phase) bool {
	return IsValidTransition(from, to)
}

// GetNextPhase returns the valid next phase
func (l *DefaultLifecycle) GetNextPhase(current Phase) (Phase, error) {
	return GetNextPhase(current)
}

// GetPreviousPhase returns the previous phase
func (l *DefaultLifecycle) GetPreviousPhase(current Phase) (Phase, error) {
	return GetPreviousPhase(current)
}

// GetCompletionCriteria returns criteria for phase completion
func (l *DefaultLifecycle) GetCompletionCriteria(phase Phase) []CompletionCriterion {
	spec, err := GetPhaseSpec(phase)
	if err != nil {
		return nil
	}
	return spec.CompletionCriteria
}

// IsTerminalPhase checks if a phase is terminal
func (l *DefaultLifecycle) IsTerminalPhase(phase Phase) bool {
	return IsTerminalPhase(phase)
}

// ValidatePhase validates a phase
func (l *DefaultLifecycle) ValidatePhase(ctx context.Context, phase Phase, artifactPaths []string) (*PhaseValidation, error) {
	spec, err := GetPhaseSpec(phase)
	if err != nil {
		return nil, err
	}

	validation := &PhaseValidation{
		Phase:         phase,
		Valid:         true,
		Verified:      []ArtifactSpec{},
		Missing:       []ArtifactSpec{},
		Invalid:       []ArtifactSpec{},
		CriteriaMet:   []CompletionCriterion{},
		CriteriaUnmet: []CompletionCriterion{},
	}

	// Build path map for quick lookup
	pathMap := make(map[string]bool)
	for _, p := range artifactPaths {
		pathMap[p] = true
	}

	// Check required artifacts
	for _, artifact := range spec.RequiredArtifacts {
		if !artifact.Required {
			continue
		}

		if pathMap[artifact.Path] {
			validation.Verified = append(validation.Verified, artifact)
		} else {
			validation.Missing = append(validation.Missing, artifact)
			validation.Valid = false
			validation.Errors = append(validation.Errors, ValidationError{
				Type:    "MISSING_ARTIFACT",
				Message: "Required artifact not found",
				Path:    artifact.Path,
			})
		}
	}

	// Check optional artifacts that exist
	for _, artifact := range spec.RequiredArtifacts {
		if artifact.Required {
			continue
		}
		if pathMap[artifact.Path] {
			validation.Verified = append(validation.Verified, artifact)
		}
	}

	// Evaluate completion criteria
	for _, criterion := range spec.CompletionCriteria {
		if criterion.Validate(phase) {
			validation.CriteriaMet = append(validation.CriteriaMet, criterion)
		} else {
			validation.CriteriaUnmet = append(validation.CriteriaUnmet, criterion)
			validation.Valid = false
			validation.Errors = append(validation.Errors, ValidationError{
				Type:    "UNMET_CRITERION",
				Message: criterion.Description,
			})
		}
	}

	return validation, nil
}
