// Package knowledge implements minimal Knowledge Promotion capability.
// Responsibility: Notebook Entry → Knowledge Candidate
// No review workflows. No approval workflows. No policy engines.
package knowledge

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Status represents the promotion status
type Status string

const (
	StatusNotebook  Status = "notebook"   // Entry in notebook
	StatusCandidate Status = "candidate"  // Promoted to candidate
	StatusPromoted  Status = "promoted"   // Promoted to knowledge
)

// Entry represents a knowledge entry
type Entry struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Source      string    `json:"source"`      // notebook, external, derived
	Status      Status    `json:"status"`
	Tags        []string  `json:"tags,omitempty"`
	Links       []string  `json:"links,omitempty"`
	CreatedAt   string    `json:"created_at"`
	PromotedAt  string    `json:"promoted_at,omitempty"`
	PromotedBy  string    `json:"promoted_by,omitempty"`
}

// Manager handles knowledge promotion
type Manager struct {
	repoPath string
	entries  map[string]*Entry
}

// NewManager creates a new knowledge manager
func NewManager(repoPath string) *Manager {
	return &Manager{
		repoPath: repoPath,
		entries:  make(map[string]*Entry),
	}
}

// Load loads entries from storage
func (m *Manager) Load() error {
	path := m.entriesPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var entries []*Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}

	for _, e := range entries {
		m.entries[e.ID] = e
	}

	return nil
}

// CreateNotebookEntry creates a notebook entry
func (m *Manager) CreateNotebookEntry(title, content, source string, tags []string) (string, error) {
	id := m.generateID()

	entry := &Entry{
		ID:        id,
		Title:     title,
		Content:   content,
		Source:    source,
		Status:    StatusNotebook,
		Tags:      tags,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	m.entries[id] = entry
	return id, m.save()
}

// PromoteToCandidate promotes a notebook entry to candidate
func (m *Manager) PromoteToCandidate(id string) error {
	entry, ok := m.entries[id]
	if !ok {
		return &NotFoundError{ID: id}
	}

	if entry.Status != StatusNotebook {
		return &InvalidTransitionError{From: string(entry.Status), To: string(StatusCandidate)}
	}

	entry.Status = StatusCandidate
	entry.PromotedAt = time.Now().Format(time.RFC3339)
	entry.PromotedBy = "runtime"

	return m.save()
}

// PromoteToKnowledge promotes a candidate to knowledge (if needed later)
func (m *Manager) PromoteToKnowledge(id string) error {
	entry, ok := m.entries[id]
	if !ok {
		return &NotFoundError{ID: id}
	}

	if entry.Status != StatusCandidate {
		return &InvalidTransitionError{From: string(entry.Status), To: string(StatusPromoted)}
	}

	entry.Status = StatusPromoted
	entry.PromotedAt = time.Now().Format(time.RFC3339)

	return m.save()
}

// Get returns an entry by ID
func (m *Manager) Get(id string) *Entry {
	return m.entries[id]
}

// List returns entries filtered by status
func (m *Manager) List(status Status) []*Entry {
	var result []*Entry
	for _, e := range m.entries {
		if status == "" || e.Status == status {
			result = append(result, e)
		}
	}
	return result
}

// Delete removes an entry
func (m *Manager) Delete(id string) error {
	if _, ok := m.entries[id]; !ok {
		return &NotFoundError{ID: id}
	}
	delete(m.entries, id)
	return m.save()
}

// errors
type NotFoundError struct {
	ID string
}

func (e *NotFoundError) Error() string {
	return "entry not found: " + e.ID
}

type InvalidTransitionError struct {
	From string
	To   string
}

func (e *InvalidTransitionError) Error() string {
	return "invalid transition from " + e.From + " to " + e.To
}

// entriesPath returns the path to entries storage
func (m *Manager) entriesPath() string {
	return filepath.Join(m.repoPath, ".kdse", "knowledge", "entries.json")
}

// save persists entries to storage
func (m *Manager) save() error {
	path := m.entriesPath()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	entries := make([]*Entry, 0, len(m.entries))
	for _, e := range m.entries {
		entries = append(entries, e)
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// generateID creates a unique entry ID
func (m *Manager) generateID() string {
	count := len(m.entries) + 1
	return "KDSE-KNOW-" + time.Now().Format("20060102") + "-" + string(rune('A'+count%26))
}

// Format formats entries for display
func Format(entries []*Entry) string {
	var result string
	for _, e := range entries {
		result += "## " + e.Title + "\n"
		result += "**ID:** " + e.ID + "\n"
		result += "**Status:** " + string(e.Status) + "\n"
		result += "**Created:** " + e.CreatedAt + "\n"
		if e.PromotedAt != "" {
			result += "**Promoted:** " + e.PromotedAt + "\n"
		}
		result += "\n" + e.Content + "\n\n"
	}
	return result
}
