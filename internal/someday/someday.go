// Package someday implements the Someday/Maybe Knowledge Management system.
// This is NOT a task backlog - it is a structured engineering knowledge repository
// for ideas that are intentionally deferred.
package someday

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// IdeaStatus represents the status of a someday item
type IdeaStatus string

const (
	StatusSomeday    IdeaStatus = "SOMEDAY"
	StatusLaboratory IdeaStatus = "LABORATORY"
	StatusPromoted   IdeaStatus = "PROMOTED"
	StatusImplemented IdeaStatus = "IMPLEMENTED"
	StatusArchived   IdeaStatus = "ARCHIVED"
	StatusRejected   IdeaStatus = "REJECTED"
)

// Idea represents a Someday/Maybe item
type Idea struct {
	ID                  string     `json:"id"`
	Title               string     `json:"title"`
	Description         string     `json:"description"`
	Problem             string     `json:"problem,omitempty"`
	PotentialValue      string     `json:"potential_value,omitempty"`
	Origin              string     `json:"origin,omitempty"`
	DateCreated         string     `json:"date_created"`
	Author              string     `json:"author,omitempty"`
	RelatedObjective    string     `json:"related_objective,omitempty"`
	RelatedKnowledge    []string   `json:"related_knowledge,omitempty"`
	Dependencies        []string   `json:"dependencies,omitempty"`
	EstimatedComplexity string     `json:"estimated_complexity,omitempty"`
	Confidence          float64    `json:"confidence"`
	Priority            int        `json:"priority"` // 1-5, 1 is highest
	Status              IdeaStatus `json:"status"`
	Examples            []string   `json:"examples,omitempty"`
	ReasonDeferred      string     `json:"reason_deferred,omitempty"`
	ReviewDate          string     `json:"review_date,omitempty"`
	Tags                []string   `json:"tags,omitempty"`
	PromotionHistory    []Promotion `json:"promotion_history,omitempty"`
	TraceabilityLinks   []string   `json:"traceability_links,omitempty"`
	UpdatedAt           string     `json:"updated_at"`
}

// Promotion represents a status change in the idea's lifecycle
type Promotion struct {
	From     IdeaStatus `json:"from"`
	To       IdeaStatus `json:"to"`
	Date     string     `json:"date"`
	Reason   string     `json:"reason,omitempty"`
	Evidence string     `json:"evidence,omitempty"`
}

// SomedayManager manages the Someday/Maybe repository
type SomedayManager struct {
	repoPath  string
	somedayPath string
	ideasPath   string
	archivedPath string
	promotedPath string
	manifestPath string
}

// New creates a new SomedayManager
func New(repoPath string) *SomedayManager {
	kdsePath := filepath.Join(repoPath, ".kdse")
	return &SomedayManager{
		repoPath:     repoPath,
		somedayPath:  filepath.Join(kdsePath, "someday"),
		ideasPath:    filepath.Join(kdsePath, "someday", "ideas"),
		archivedPath: filepath.Join(kdsePath, "someday", "archived"),
		promotedPath: filepath.Join(kdsePath, "someday", "promoted"),
		manifestPath: filepath.Join(kdsePath, "someday", "someday.yaml"),
	}
}

// Initialize creates the someday directory structure
func (m *SomedayManager) Initialize() error {
	dirs := []string{
		m.somedayPath,
		m.ideasPath,
		m.archivedPath,
		m.promotedPath,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create manifest if it doesn't exist
	if _, err := os.Stat(m.manifestPath); os.IsNotExist(err) {
		manifest := &SomedayManifest{
			Version:      "1.0.0",
			CreatedAt:    time.Now().Format(time.RFC3339),
			LastUpdated:  time.Now().Format(time.RFC3339),
			Ideas:        []string{},
			TotalCount:   0,
			ByStatus:     map[string]int{},
			ByPriority:   map[string]int{},
			NextIdeaID:   1,
		}
		if err := m.saveManifest(manifest); err != nil {
			return err
		}
	}

	return nil
}

// SomedayManifest represents the someday inventory
type SomedayManifest struct {
	Version     string            `json:"version"`
	CreatedAt  string            `json:"created_at"`
	LastUpdated string            `json:"last_updated"`
	Ideas      []string          `json:"ideas"`
	TotalCount int               `json:"total_count"`
	ByStatus   map[string]int    `json:"by_status"`
	ByPriority map[string]int    `json:"by_priority"`
	NextIdeaID int               `json:"next_idea_id"`
}

// Add creates a new someday idea
func (m *SomedayManager) Add(title, description, problem, origin, author string) (*Idea, error) {
	manifest, err := m.loadManifest()
	if err != nil {
		return nil, err
	}

	idea := &Idea{
		ID:              fmt.Sprintf("IDEA-%03d", manifest.NextIdeaID),
		Title:           title,
		Description:     description,
		Problem:         problem,
		Origin:         origin,
		Author:         author,
		DateCreated:    time.Now().Format(time.RFC3339),
		Status:         StatusSomeday,
		Confidence:     0.5, // Default confidence
		Priority:       3,    // Default priority (medium)
		RelatedKnowledge: []string{},
		Dependencies:    []string{},
		Examples:       []string{},
		Tags:           []string{},
		PromotionHistory: []Promotion{},
		TraceabilityLinks: []string{},
		UpdatedAt:      time.Now().Format(time.RFC3339),
	}

	// Save idea file
	ideaPath := filepath.Join(m.ideasPath, fmt.Sprintf("%s.yaml", idea.ID))
	if err := m.saveIdea(ideaPath, idea); err != nil {
		return nil, err
	}

	// Update manifest
	manifest.Ideas = append(manifest.Ideas, idea.ID)
	manifest.TotalCount++
	manifest.ByStatus[string(StatusSomeday)]++
	manifest.ByPriority[fmt.Sprintf("priority-%d", idea.Priority)]++
	manifest.NextIdeaID++
	manifest.LastUpdated = time.Now().Format(time.RFC3339)

	if err := m.saveManifest(manifest); err != nil {
		return nil, err
	}

	return idea, nil
}

// List returns all ideas, optionally filtered by status
func (m *SomedayManager) List(status IdeaStatus) ([]*Idea, error) {
	var ideas []*Idea

	// Load all ideas from ideas directory
	if status == "" || status == StatusSomeday {
		files, err := filepath.Glob(filepath.Join(m.ideasPath, "*.yaml"))
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			idea, err := m.loadIdea(f)
			if err != nil {
				continue
			}
			ideas = append(ideas, idea)
		}
	}

	// Load promoted ideas
	if status == "" || status == StatusPromoted {
		files, err := filepath.Glob(filepath.Join(m.promotedPath, "*.yaml"))
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			idea, err := m.loadIdea(f)
			if err != nil {
				continue
			}
			ideas = append(ideas, idea)
		}
	}

	// Load archived ideas
	if status == "" || status == StatusArchived {
		files, err := filepath.Glob(filepath.Join(m.archivedPath, "*.yaml"))
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			idea, err := m.loadIdea(f)
			if err != nil {
				continue
			}
			ideas = append(ideas, idea)
		}
	}

	// Sort by priority then by date
	sort.Slice(ideas, func(i, j int) bool {
		if ideas[i].Priority != ideas[j].Priority {
			return ideas[i].Priority < ideas[j].Priority
		}
		return ideas[i].DateCreated > ideas[j].DateCreated
	})

	return ideas, nil
}

// Show returns a specific idea by ID
func (m *SomedayManager) Show(id string) (*Idea, error) {
	// Search all directories
	paths := []string{
		filepath.Join(m.ideasPath, fmt.Sprintf("%s.yaml", id)),
		filepath.Join(m.archivedPath, fmt.Sprintf("%s.yaml", id)),
		filepath.Join(m.promotedPath, fmt.Sprintf("%s.yaml", id)),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return m.loadIdea(path)
		}
	}

	return nil, fmt.Errorf("idea not found: %s", id)
}

// Promote moves an idea to the promoted directory
func (m *SomedayManager) Promote(id, reason string) (*Idea, error) {
	idea, err := m.Show(id)
	if err != nil {
		return nil, err
	}

	if idea.Status == StatusPromoted || idea.Status == StatusImplemented {
		return nil, fmt.Errorf("idea %s is already %s", id, idea.Status)
	}

	oldStatus := idea.Status

	// Record promotion
	promotion := Promotion{
		From:   oldStatus,
		To:     StatusPromoted,
		Date:   time.Now().Format(time.RFC3339),
		Reason: reason,
	}
	idea.PromotionHistory = append(idea.PromotionHistory, promotion)
	idea.Status = StatusPromoted
	idea.UpdatedAt = time.Now().Format(time.RFC3339)

	// Move to promoted directory
	oldPath := filepath.Join(m.ideasPath, fmt.Sprintf("%s.yaml", id))
	newPath := filepath.Join(m.promotedPath, fmt.Sprintf("%s.yaml", id))

	if err := os.Rename(oldPath, newPath); err != nil {
		return nil, err
	}

	// Update manifest
	manifest, err := m.loadManifest()
	if err != nil {
		return nil, err
	}
	manifest.ByStatus[string(oldStatus)]--
	manifest.ByStatus[string(StatusPromoted)]++
	manifest.LastUpdated = time.Now().Format(time.RFC3339)
	if err := m.saveManifest(manifest); err != nil {
		return nil, err
	}

	return idea, nil
}

// Archive moves an idea to the archived directory
func (m *SomedayManager) Archive(id, reason string) (*Idea, error) {
	idea, err := m.Show(id)
	if err != nil {
		return nil, err
	}

	if idea.Status == StatusArchived {
		return nil, fmt.Errorf("idea %s is already archived", id)
	}

	oldStatus := idea.Status

	// Record archival
	promotion := Promotion{
		From:   oldStatus,
		To:     StatusArchived,
		Date:   time.Now().Format(time.RFC3339),
		Reason: reason,
	}
	idea.PromotionHistory = append(idea.PromotionHistory, promotion)
	idea.Status = StatusArchived
	idea.ReasonDeferred = reason
	idea.UpdatedAt = time.Now().Format(time.RFC3339)

	// Determine source directory
	var sourcePath, destPath string
	switch oldStatus {
	case StatusSomeday:
		sourcePath = filepath.Join(m.ideasPath, fmt.Sprintf("%s.yaml", id))
	case StatusPromoted:
		sourcePath = filepath.Join(m.promotedPath, fmt.Sprintf("%s.yaml", id))
	default:
		sourcePath = filepath.Join(m.ideasPath, fmt.Sprintf("%s.yaml", id))
	}
	destPath = filepath.Join(m.archivedPath, fmt.Sprintf("%s.yaml", id))

	if err := os.Rename(sourcePath, destPath); err != nil {
		return nil, err
	}

	// Update manifest
	manifest, err := m.loadManifest()
	if err != nil {
		return nil, err
	}
	manifest.ByStatus[string(oldStatus)]--
	manifest.ByStatus[string(StatusArchived)]++
	manifest.LastUpdated = time.Now().Format(time.RFC3339)
	if err := m.saveManifest(manifest); err != nil {
		return nil, err
	}

	return idea, nil
}

// Search finds ideas matching a query
func (m *SomedayManager) Search(query string) ([]*Idea, error) {
	allIdeas, err := m.List("")
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	var results []*Idea

	for _, idea := range allIdeas {
		// Search in title, description, problem, tags
		if strings.Contains(strings.ToLower(idea.Title), query) ||
			strings.Contains(strings.ToLower(idea.Description), query) ||
			strings.Contains(strings.ToLower(idea.Problem), query) {
			results = append(results, idea)
			continue
		}

		// Search in tags
		for _, tag := range idea.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				results = append(results, idea)
				break
			}
		}
	}

	return results, nil
}

// Review returns ideas ready for review
func (m *SomedayManager) Review() ([]*Idea, error) {
	ideas, err := m.List("")
	if err != nil {
		return nil, err
	}

	var review []*Idea
	for _, idea := range ideas {
		// Include ideas with past review dates or high priority
		if idea.Status == StatusSomeday && idea.Priority <= 2 {
			review = append(review, idea)
		}
	}

	return review, nil
}

// Export generates an export of all ideas
func (m *SomedayManager) Export() (*SomedayExport, error) {
	ideas, err := m.List("")
	if err != nil {
		return nil, err
	}

	manifest, err := m.loadManifest()
	if err != nil {
		return nil, err
	}

	return &SomedayExport{
		Version:     "1.0.0",
		ExportedAt: time.Now().Format(time.RFC3339),
		Manifest:   manifest,
		Ideas:      ideas,
	}, nil
}

// SomedayExport represents an export of all someday items
type SomedayExport struct {
	Version     string    `json:"version"`
	ExportedAt  string    `json:"exported_at"`
	Manifest    *SomedayManifest `json:"manifest"`
	Ideas       []*Idea   `json:"ideas"`
}

// DetectIdeaPatterns searches text for someday/maybe patterns
func DetectIdeaPatterns(text string) []string {
	patterns := []string{
		`(?i)\bmaybe\b`,
		`(?i)\bsomeday\b`,
		`(?i)\bfuture\b`,
		`(?i)\binteresting\b`,
		`(?i)\bnot now\b`,
		`(?i)\blater\b`,
		`(?i)\bperhaps\b`,
		`(?i)\bpotentially\b`,
		`(?i)\bif we had time\b`,
		`(?i)\bnice to have\b`,
		`(?i)\bdown the road\b`,
	}

	var matches []string
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(text) {
			matches = append(matches, pattern)
		}
	}

	return matches
}

// File helpers

func (m *SomedayManager) saveIdea(path string, idea *Idea) error {
	data, err := json.MarshalIndent(idea, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (m *SomedayManager) loadIdea(path string) (*Idea, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var idea Idea
	if err := json.Unmarshal(data, &idea); err != nil {
		return nil, err
	}

	return &idea, nil
}

func (m *SomedayManager) loadManifest() (*SomedayManifest, error) {
	data, err := os.ReadFile(m.manifestPath)
	if err != nil {
		return nil, err
	}

	var manifest SomedayManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (m *SomedayManager) saveManifest(manifest *SomedayManifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.manifestPath, data, 0644)
}

// FormatIdea formats an idea for display
func FormatIdea(idea *Idea) string {
	statusIcon := "○"
	switch idea.Status {
	case StatusSomeday:
		statusIcon = "◇"
	case StatusLaboratory:
		statusIcon = "◐"
	case StatusPromoted:
		statusIcon = "◉"
	case StatusImplemented:
		statusIcon = "●"
	case StatusArchived:
		statusIcon = "○"
	case StatusRejected:
		statusIcon = "✗"
	}

	priorityStr := ""
	switch idea.Priority {
	case 1:
		priorityStr = "🔴 CRITICAL"
	case 2:
		priorityStr = "🟠 HIGH"
	case 3:
		priorityStr = "🟡 MEDIUM"
	case 4:
		priorityStr = "🟢 LOW"
	case 5:
		priorityStr = "⚪ TRIVIAL"
	}

	return fmt.Sprintf(`%s %s [%s]

Title:       %s
Description: %s
Problem:     %s
Priority:    %s
Confidence:  %.0f%%
Tags:        %s
Status:      %s

Created:     %s
Author:      %s
Origin:      %s

Related Objective: %s
Dependencies:     %s

%s
`,
		statusIcon, idea.ID, idea.Status,
		idea.Title,
		idea.Description,
		idea.Problem,
		priorityStr,
		idea.Confidence*100,
		strings.Join(idea.Tags, ", "),
		idea.Status,
		idea.DateCreated,
		idea.Author,
		idea.Origin,
		idea.RelatedObjective,
		strings.Join(idea.Dependencies, ", "),
		formatPromotionHistory(idea.PromotionHistory),
	)
}

func formatPromotionHistory(history []Promotion) string {
	if len(history) == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, "Promotion History:")
	for _, p := range history {
		lines = append(lines, fmt.Sprintf("  %s → %s (%s)", p.From, p.To, p.Date))
		if p.Reason != "" {
			lines = append(lines, fmt.Sprintf("    Reason: %s", p.Reason))
		}
	}
	return strings.Join(lines, "\n")
}

// FormatList formats a list of ideas for display
func FormatList(ideas []*Idea) string {
	if len(ideas) == 0 {
		return "No ideas found."
	}

	var lines []string
	lines = append(lines, "")
	lines = append(lines, "╔═══════════════════════════════════════════════════════════════╗")
	lines = append(lines, "║                    SOMEDAY/MAYBE IDEAS                     ║")
	lines = append(lines, "╠═══════════════════════════════════════════════════════════════╣")
	lines = append(lines, fmt.Sprintf("║ Total Ideas: %d                                            ║", len(ideas)))
	lines = append(lines, "╠═══════════════════════════════════════════════════════════════╣")

	// Group by status
	byStatus := make(map[IdeaStatus][]*Idea)
	for _, idea := range ideas {
		byStatus[idea.Status] = append(byStatus[idea.Status], idea)
	}

	for status, statusIdeas := range byStatus {
		lines = append(lines, fmt.Sprintf("║ %s (%d)                                                  ║", status, len(statusIdeas)))
		for _, idea := range statusIdeas {
			priorityStr := fmt.Sprintf("P%d", idea.Priority)
			lines = append(lines, fmt.Sprintf("║   %s %-10s %-20s %s║", priorityStr, idea.ID, truncate(idea.Title, 20), idea.Status))
		}
	}

	lines = append(lines, "╚═══════════════════════════════════════════════════════════════╝")
	lines = append(lines, "")
	lines = append(lines, "Use 'kdse someday show <ID>' for details")
	lines = append(lines, "")

	return strings.Join(lines, "\n")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
