// Package agreement implements the KDSE Agreement subsystem.
// The Agreement establishes project identity, methodology, constraints, and context.
// Subsequent runtime interactions reference the Agreement instead of rebuilding context.
package agreement

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Agreement represents the authoritative project agreement.
// All engineering work references this agreement.
type Agreement struct {
	// Identity
	ProjectName    string `json:"project_name"`
	ProjectPath    string `json:"project_path"`
	InitializedAt  string `json:"initialized_at"`
	LastModifiedAt string `json:"last_modified_at"`

	// Versions
	MethodologyVersion string `json:"methodology_version"`
	RuntimeVersion     string `json:"runtime_version"`

	// Current state
	CurrentPhase       string   `json:"current_phase"`
	ActiveSessions     []string `json:"active_sessions,omitempty"`
	CompletedWorkOrders []string `json:"completed_work_orders,omitempty"`

	// Constraints
	Constraints ConstraintSet `json:"constraints"`

	// Shared assumptions accumulated during engineering
	Assumptions []Assumption `json:"assumptions,omitempty"`

	// Engineering context (derived from artifacts)
	Context EngineeringContext `json:"context"`
}

// ConstraintSet contains all architectural constraints
type ConstraintSet struct {
	SubsystemOwnership    []OwnershipConstraint    `json:"subsystem_ownership,omitempty"`
	DeploymentOwnership   []OwnershipConstraint    `json:"deployment_ownership,omitempty"`
	RuntimeOwnership      []OwnershipConstraint    `json:"runtime_ownership,omitempty"`
	ConfigurationRules   []ConfigurationConstraint `json:"configuration_rules,omitempty"`
	AuthorityBoundaries  []AuthorityBoundary      `json:"authority_boundaries,omitempty"`
}

// OwnershipConstraint defines who owns a subsystem/resource
type OwnershipConstraint struct {
	Path     string `json:"path"`      // Subsystem or resource path
	Owner    string `json:"owner"`     // Owner identifier
	Contact  string `json:"contact,omitempty"`
	Since    string `json:"since"`
	Approved string `json:"approved_by,omitempty"`
}

// ConfigurationConstraint defines configuration rules
type ConfigurationConstraint struct {
	Rule      string `json:"rule"`
	AppliesTo string `json:"applies_to"`
	Enforced  bool   `json:"enforced"`
}

// AuthorityBoundary defines authority limits
type AuthorityBoundary struct {
	Boundary string `json:"boundary"`
	Allowed  []string `json:"allowed_actions,omitempty"`
	Denied   []string `json:"denied_actions,omitempty"`
}

// Assumption represents a shared assumption
type Assumption struct {
	ID        string `json:"id"`
	Statement string `json:"statement"`
	CreatedAt string `json:"created_at"`
	Validated bool   `json:"validated"`
}

// EngineeringContext contains derived engineering context
type EngineeringContext struct {
	Subsystems     []string `json:"subsystems"`
	Languages      []string `json:"languages,omitempty"`
	Frameworks     []string `json:"frameworks,omitempty"`
	Infrastructure []string `json:"infrastructure,omitempty"`
	Dependencies   []string `json:"dependencies,omitempty"`
}

// Manager handles Agreement lifecycle
type Manager struct {
	repoPath  string
	agreement *Agreement
}

// NewManager creates a new Agreement manager
func NewManager(repoPath string) *Manager {
	return &Manager{
		repoPath: repoPath,
	}
}

// Load loads the existing agreement
func (m *Manager) Load() (*Agreement, error) {
	path := m.agreementPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var agreement Agreement
	if err := json.Unmarshal(data, &agreement); err != nil {
		return nil, err
	}

	m.agreement = &agreement
	return &agreement, nil
}

// Create creates a new agreement
func (m *Manager) Create(projectName, projectPath, methodVersion, runtimeVersion string) (*Agreement, error) {
	agreement := &Agreement{
		ProjectName:        projectName,
		ProjectPath:       projectPath,
		InitializedAt:     time.Now().Format(time.RFC3339),
		LastModifiedAt:    time.Now().Format(time.RFC3339),
		MethodologyVersion: methodVersion,
		RuntimeVersion:     runtimeVersion,
		CurrentPhase:      "Problem",
		Constraints:       ConstraintSet{},
		Assumptions:       []Assumption{},
		Context:           EngineeringContext{},
	}

	if err := m.save(agreement); err != nil {
		return nil, err
	}

	m.agreement = agreement
	return agreement, nil
}

// Get returns the current agreement (loads if not in memory)
func (m *Manager) Get() (*Agreement, error) {
	if m.agreement != nil {
		return m.agreement, nil
	}
	return m.Load()
}

// UpdatePhase updates the current phase
func (m *Manager) UpdatePhase(phase string) error {
	agreement, err := m.Get()
	if err != nil {
		return err
	}

	agreement.CurrentPhase = phase
	agreement.LastModifiedAt = time.Now().Format(time.RFC3339)
	return m.save(agreement)
}

// AddConstraint adds an ownership constraint
func (m *Manager) AddSubsystemOwnership(path, owner, contact, approvedBy string) error {
	agreement, err := m.Get()
	if err != nil {
		return err
	}

	constraint := OwnershipConstraint{
		Path:     path,
		Owner:    owner,
		Contact:  contact,
		Since:    time.Now().Format(time.RFC3339),
		Approved: approvedBy,
	}

	agreement.Constraints.SubsystemOwnership = append(
		agreement.Constraints.SubsystemOwnership,
		constraint,
	)
	agreement.LastModifiedAt = time.Now().Format(time.RFC3339)
	return m.save(agreement)
}

// AddAuthorityBoundary adds an authority boundary
func (m *Manager) AddAuthorityBoundary(boundary string, allowed, denied []string) error {
	agreement, err := m.Get()
	if err != nil {
		return err
	}

	ab := AuthorityBoundary{
		Boundary: boundary,
		Allowed:  allowed,
		Denied:   denied,
	}

	agreement.Constraints.AuthorityBoundaries = append(
		agreement.Constraints.AuthorityBoundaries,
		ab,
	)
	agreement.LastModifiedAt = time.Now().Format(time.RFC3339)
	return m.save(agreement)
}

// AddAssumption adds a shared assumption
func (m *Manager) AddAssumption(statement string) (string, error) {
	agreement, err := m.Get()
	if err != nil {
		return "", err
	}

	id := generateAssumptionID(len(agreement.Assumptions) + 1)
	assumption := Assumption{
		ID:        id,
		Statement: statement,
		CreatedAt: time.Now().Format(time.RFC3339),
		Validated: false,
	}

	agreement.Assumptions = append(agreement.Assumptions, assumption)
	agreement.LastModifiedAt = time.Now().Format(time.RFC3339)
	return id, m.save(agreement)
}

// ValidateConstraint checks if an action violates constraints
func (m *Manager) ValidateConstraint(subsystemPath, action string) (bool, string) {
	agreement, err := m.Get()
	if err != nil {
		return true, "" // No agreement = no constraints
	}

	// Check subsystem ownership
	for _, oc := range agreement.Constraints.SubsystemOwnership {
		if oc.Path == subsystemPath {
			// Check authority boundaries
			for _, ab := range agreement.Constraints.AuthorityBoundaries {
				if ab.Boundary == subsystemPath {
					for _, denied := range ab.Denied {
						if denied == action || denied == "*" {
							return false, "Action '" + action + "' denied by authority boundary"
						}
					}
				}
			}
		}
	}

	return true, ""
}

// UpdateContext updates engineering context
func (m *Manager) UpdateContext(ctx EngineeringContext) error {
	agreement, err := m.Get()
	if err != nil {
		return err
	}

	agreement.Context = ctx
	agreement.LastModifiedAt = time.Now().Format(time.RFC3339)
	return m.save(agreement)
}

// agreementPath returns the path to the agreement file
func (m *Manager) agreementPath() string {
	return filepath.Join(m.repoPath, ".kdse", "agreement.json")
}

// save persists the agreement
func (m *Manager) save(agreement *Agreement) error {
	path := m.agreementPath()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(agreement, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func generateAssumptionID(n int) string {
	return "KDSE-ASM-" + time.Now().Format("20060102") + "-" + string(rune('A'+n%26))
}

// Format formats the agreement for display
func Format(a *Agreement) string {
	data, _ := json.MarshalIndent(a, "", "  ")
	return string(data)
}
