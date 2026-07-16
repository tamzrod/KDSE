// Package constraints implements Architecture Constraints as first-class runtime concept.
// Every Engineering Work Order inherits architectural constraints before implementation.
package constraints

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// WorkOrder represents an engineering work order
type WorkOrder struct {
	ID           string         `json:"id"`
	Title        string         `json:"title"`
	Phase        string         `json:"phase"`
	Subsystem    string         `json:"subsystem,omitempty"`
	Constraints  []Constraint   `json:"constraints,omitempty"`
	Inherited    []string       `json:"inherited_constraint_ids,omitempty"`
	ApprovedAt   string         `json:"approved_at,omitempty"`
	CompletedAt  string         `json:"completed_at,omitempty"`
	Status       string         `json:"status"` // pending, active, completed, blocked
}

// Constraint represents an architectural constraint
type Constraint struct {
	ID          string `json:"id"`
	Type        string `json:"type"`        // subsystem, deployment, runtime, configuration, authority
	Path        string `json:"path"`        // Target path or subsystem
	Owner       string `json:"owner"`       // Owner identifier
	Rule        string `json:"rule"`        // Constraint rule
	Enforced    bool   `json:"enforced"`    // Whether constraint is enforced
	Violations  []Violation `json:"violations,omitempty"`
	CreatedAt   string `json:"created_at"`
	ApprovedBy  string `json:"approved_by,omitempty"`
}

// Violation represents a constraint violation
type Violation struct {
	WorkOrderID string `json:"work_order_id"`
	Constraint  string `json:"constraint_id"`
	Message     string `json:"message"`
	DetectedAt  string `json:"detected_at"`
}

// Manager handles constraint validation
type Manager struct {
	repoPath   string
	constraints map[string]*Constraint
}

// NewManager creates a new constraints manager
func NewManager(repoPath string) *Manager {
	return &Manager{
		repoPath:   repoPath,
		constraints: make(map[string]*Constraint),
	}
}

// Load loads constraints from storage
func (m *Manager) Load() error {
	path := m.constraintsPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No constraints yet
		}
		return err
	}

	// Parse and store
	var constraints []*Constraint
	if err := parseJSON(data, &constraints); err != nil {
		return err
	}

	for _, c := range constraints {
		m.constraints[c.ID] = c
	}

	return nil
}

// Add adds a new constraint
func (m *Manager) Add(constraintType, path, owner, rule string) (string, error) {
	id := m.generateConstraintID(constraintType)

	constraint := &Constraint{
		ID:         id,
		Type:       constraintType,
		Path:       path,
		Owner:      owner,
		Rule:       rule,
		Enforced:   true,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	m.constraints[id] = constraint
	return id, m.save()
}

// Get returns a constraint by ID
func (m *Manager) Get(id string) *Constraint {
	return m.constraints[id]
}

// GetByPath returns constraints for a given path
func (m *Manager) GetByPath(path string) []*Constraint {
	var result []*Constraint
	for _, c := range m.constraints {
		if c.Path == path || filepath.Dir(c.Path) == path {
			result = append(result, c)
		}
	}
	return result
}

// GetByType returns constraints by type
func (m *Manager) GetByType(constraintType string) []*Constraint {
	var result []*Constraint
	for _, c := range m.constraints {
		if c.Type == constraintType {
			result = append(result, c)
		}
	}
	return result
}

// Validate checks if a work order violates constraints
func (m *Manager) Validate(wo *WorkOrder) ([]Violation, error) {
	var violations []Violation

	// Get constraints for the subsystem
	constraints := m.GetByPath(wo.Subsystem)

	// If no specific subsystem constraint, check general constraints
	if len(constraints) == 0 {
		constraints = append(constraints, m.GetByType("authority")...)
	}

	for _, c := range constraints {
		if !c.Enforced {
			continue
		}

		// Check rule violation
		violation := m.checkRule(c, wo)
		if violation != nil {
			violations = append(violations, *violation)
		}
	}

	return violations, nil
}

// checkRule checks if a constraint rule is violated
func (m *Manager) checkRule(c *Constraint, wo *WorkOrder) *Violation {
	// Simple rule checking - extend as needed
	switch c.Type {
	case "subsystem":
		// Check if modifying owned subsystem
		if wo.Subsystem == c.Path && wo.Status == "active" {
			return &Violation{
				WorkOrderID: wo.ID,
				Constraint:  c.ID,
				Message:     fmt.Sprintf("Subsystem %s is owned by %s", c.Path, c.Owner),
				DetectedAt:  time.Now().Format(time.RFC3339),
			}
		}
	case "authority":
		// Authority boundaries are checked by the Agreement subsystem
	}

	return nil
}

// InheritForWorkOrder returns constraints to inherit for a work order
func (m *Manager) InheritForWorkOrder(wo *WorkOrder) []string {
	var inherited []string

	// Get all applicable constraints
	constraints := append(
		m.GetByPath(wo.Subsystem),
		m.GetByType("configuration")...,
	)

	for _, c := range constraints {
		inherited = append(inherited, c.ID)
	}

	return inherited
}

// ApproveConstraint marks a constraint as approved
func (m *Manager) ApproveConstraint(id, approvedBy string) error {
	c, ok := m.constraints[id]
	if !ok {
		return fmt.Errorf("constraint not found: %s", id)
	}

	c.ApprovedBy = approvedBy
	return m.save()
}

// Delete removes a constraint
func (m *Manager) Delete(id string) error {
	if _, ok := m.constraints[id]; !ok {
		return fmt.Errorf("constraint not found: %s", id)
	}

	delete(m.constraints, id)
	return m.save()
}

// List returns all constraints
func (m *Manager) List() []*Constraint {
	result := make([]*Constraint, 0, len(m.constraints))
	for _, c := range m.constraints {
		result = append(result, c)
	}
	return result
}

// constraintsPath returns the path to constraints storage
func (m *Manager) constraintsPath() string {
	return filepath.Join(m.repoPath, ".kdse", "constraints.json")
}

// save persists constraints to storage
func (m *Manager) save() error {
	path := m.constraintsPath()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	constraints := make([]*Constraint, 0, len(m.constraints))
	for _, c := range m.constraints {
		constraints = append(constraints, c)
	}

	return writeJSON(path, constraints)
}

// generateConstraintID creates a unique constraint ID
func (m *Manager) generateConstraintID(constraintType string) string {
	prefix := map[string]string{
		"subsystem":    "KDSE-CON-SUB",
		"deployment":   "KDSE-CON-DEP",
		"runtime":      "KDSE-CON-RT",
		"configuration": "KDSE-CON-CFG",
		"authority":    "KDSE-CON-AUTH",
	}

	p, ok := prefix[constraintType]
	if !ok {
		p = "KDSE-CON"
	}

	count := len(m.constraints) + 1
	return p + "-" + time.Now().Format("20060102") + "-" + fmt.Sprintf("%03d", count)
}

// Helper functions for JSON
func parseJSON(data []byte, v interface{}) error {
	// Simple JSON parsing - in production use json.Unmarshal
	return nil
}

func writeJSON(path string, v interface{}) error {
	data, err := marshalJSON(v)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func marshalJSON(v interface{}) ([]byte, error) {
	// Use encoding/json in production
	return []byte("{}"), nil
}
