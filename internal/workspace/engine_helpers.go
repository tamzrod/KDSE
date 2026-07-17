package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	kdseVersion = "1.0.0"
)

// Helper functions for the workspace engine

func (e *DefaultEngine) exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (e *DefaultEngine) kdsePath() string {
	return filepath.Join(e.workspacePath, ".kdse")
}

func (e *DefaultEngine) loadRuntimeConfig() (*RuntimeConfig, error) {
	configPath := filepath.Join(e.kdsePath(), "runtime.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, ErrRuntimeInvalid
	}

	// Simple YAML parsing
	var config RuntimeConfig
	if strings.Contains(string(data), "type:") && strings.Contains(string(data), "version:") {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "type:") {
				t := strings.TrimPrefix(line, "type:")
				config.Type = RuntimeType(strings.TrimSpace(t))
			}
			if strings.HasPrefix(line, "version:") {
				v := strings.TrimPrefix(line, "version:")
				config.Version = strings.TrimSpace(v)
			}
		}
	}

	if config.Type == "" || config.Version == "" {
		return nil, ErrRuntimeInvalid
	}

	return &config, nil
}

func (e *DefaultEngine) loadPhase() (Phase, error) {
	phasePath := filepath.Join(e.kdsePath(), "phase.yaml")
	data, err := os.ReadFile(phasePath)
	if err != nil {
		return "", ErrPhaseInvalid
	}

	// Simple YAML parsing
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "current:") {
			p := strings.TrimPrefix(line, "current:")
			phase := Phase(strings.TrimSpace(p))
			
			// Validate phase
			validPhases := map[Phase]bool{
				PhaseInitialization: true,
				PhaseKnowledge:      true,
				PhaseArchitecture:   true,
				PhaseImplementation: true,
				PhaseVerification:   true,
				PhaseReports:        true,
			}
			
			if !validPhases[phase] {
				return "", ErrPhaseInvalid
			}
			
			return phase, nil
		}
	}

	return "", ErrPhaseInvalid
}

func (e *DefaultEngine) persistPhase(phase Phase) error {
	phasePath := filepath.Join(e.kdsePath(), "phase.yaml")
	content := fmt.Sprintf(`current: %s
previous: %s
updated: %s
`, phase, "", time.Now().Format(time.RFC3339))
	return os.WriteFile(phasePath, []byte(content), 0644)
}

func (e *DefaultEngine) initializePhase(dir string) error {
	phasePath := filepath.Join(dir, "phase.yaml")
	content := fmt.Sprintf(`current: %s
previous: none
initialized: %s
`, PhaseInitialization, time.Now().Format(time.RFC3339))
	return os.WriteFile(phasePath, []byte(content), 0644)
}

func (e *DefaultEngine) createRuntimeFiles(dir string, opts InitOptions) error {
	// Create runtime.yaml
	runtimeContent := fmt.Sprintf(`runtime:
  type: %s
  version: %s
  commit: %s

template:
  version: "%s"
  commit: ""

workspace:
  version: %s
  root: %s

session:
  created: %s
`, opts.Type, kdseVersion, "", opts.Template, kdseVersion, e.workspacePath, time.Now().Format(time.RFC3339))
	if err := os.WriteFile(filepath.Join(dir, "runtime.yaml"), []byte(runtimeContent), 0644); err != nil {
		return err
	}

	// Create workspace.yaml
	workspaceContent := fmt.Sprintf(`workspace:
  version: %s
  initialized: %s
  root: %s
`, kdseVersion, time.Now().Format(time.RFC3339), e.workspacePath)
	if err := os.WriteFile(filepath.Join(dir, "workspace.yaml"), []byte(workspaceContent), 0644); err != nil {
		return err
	}

	// Create methodology.yaml
	methodologyContent := fmt.Sprintf(`methodology:
  version: %s
  phases:
    - initialization
    - knowledge
    - architecture
    - implementation
    - verification
    - reports
`, kdseVersion)
	if err := os.WriteFile(filepath.Join(dir, "methodology.yaml"), []byte(methodologyContent), 0644); err != nil {
		return err
	}

	// Create session.yaml
	sessionContent := fmt.Sprintf(`session:
  id: ""
  created: %s
  ended: ""
`, time.Now().Format(time.RFC3339))
	if err := os.WriteFile(filepath.Join(dir, "session.yaml"), []byte(sessionContent), 0644); err != nil {
		return err
	}

	// Create metadata.yaml
	metadataContent := fmt.Sprintf(`metadata:
  created: %s
  version: %s
  engine: workspace-engine
`, time.Now().Format(time.RFC3339), kdseVersion)
	if err := os.WriteFile(filepath.Join(dir, "metadata.yaml"), []byte(metadataContent), 0644); err != nil {
		return err
	}

	return nil
}

func (e *DefaultEngine) createDirectoryReadme(dir string) error {
	readmePath := filepath.Join(dir, "README.md")
	readme := fmt.Sprintf(`# %s

This directory contains %s artifacts.

`, filepath.Base(dir), filepath.Base(dir))
	return os.WriteFile(readmePath, []byte(readme), 0644)
}

func (e *DefaultEngine) loadPhaseHistory() ([]PhaseTransition, error) {
	historyPath := filepath.Join(e.kdsePath(), "phase-history.yaml")
	data, err := os.ReadFile(historyPath)
	if os.IsNotExist(err) {
		return []PhaseTransition{}, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse history
	var history []PhaseTransition
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "-") {
			history = append(history, PhaseTransition{})
		}
	}

	return history, nil
}

func (e *DefaultEngine) recordTransition(t *Transition) error {
	historyPath := filepath.Join(e.kdsePath(), "phase-history.yaml")

	// Load existing history
	history, _ := e.loadPhaseHistory()

	// Add new transition
	history = append(history, PhaseTransition{
		From:      t.From,
		To:        t.To,
		Timestamp: t.Timestamp,
	})

	// Write history
	var content strings.Builder
	content.WriteString("# Phase History\n\n")
	for _, h := range history {
		content.WriteString(fmt.Sprintf("- %s → %s (%s)\n", h.From, h.To, h.Timestamp.Format(time.RFC3339)))
	}

	return os.WriteFile(historyPath, []byte(content.String()), 0644)
}

func (e *DefaultEngine) loadArtifactPaths(phase Phase) ([]string, error) {
	dirMap := map[Phase]string{
		PhaseInitialization: e.kdsePath(),
		PhaseKnowledge:      filepath.Join(e.kdsePath(), "knowledge"),
		PhaseArchitecture:   filepath.Join(e.kdsePath(), "architecture"),
		PhaseImplementation: filepath.Join(e.kdsePath(), "implementation"),
		PhaseVerification:   filepath.Join(e.kdsePath(), "verification"),
		PhaseReports:        filepath.Join(e.kdsePath(), "reports"),
	}

	dir, ok := dirMap[phase]
	if !ok {
		return nil, ErrPhaseInvalid
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, entry := range entries {
		if !entry.IsDir() {
			relPath, _ := filepath.Rel(e.workspacePath, filepath.Join(dir, entry.Name()))
			paths = append(paths, relPath)
		}
	}

	return paths, nil
}

func (e *DefaultEngine) listArtifacts(dir string) ([]Artifact, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var artifacts []Artifact
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		path, _ := filepath.Rel(e.workspacePath, filepath.Join(dir, entry.Name()))
		artifacts = append(artifacts, Artifact{
			Path: path,
			Type: filepath.Ext(entry.Name()),
			Size: info.Size(),
		})
	}

	return artifacts, nil
}

func (e *DefaultEngine) countArtifacts() int {
	count := 0
	dirs := []string{
		filepath.Join(e.kdsePath(), "knowledge"),
		filepath.Join(e.kdsePath(), "architecture"),
		filepath.Join(e.kdsePath(), "implementation"),
		filepath.Join(e.kdsePath(), "verification"),
		filepath.Join(e.kdsePath(), "reports"),
	}

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				count++
			}
		}
	}

	return count
}

func (e *DefaultEngine) countReports() int {
	reportsDir := filepath.Join(e.kdsePath(), "reports")
	entries, err := os.ReadDir(reportsDir)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			count++
		}
	}

	return count
}

func (e *DefaultEngine) loadSession() (*Session, error) {
	sessionPath := filepath.Join(e.kdsePath(), "session.yaml")
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// Simple parsing
	session := &Session{Metadata: make(map[string]string)}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "id:") {
			session.ID = strings.TrimSpace(strings.TrimPrefix(line, "id:"))
		}
		if strings.HasPrefix(line, "created:") {
			created, _ := time.Parse(time.RFC3339, strings.TrimSpace(strings.TrimPrefix(line, "created:")))
			session.Created = created
		}
	}

	if session.ID == "" {
		session.ID = generateSessionID()
	}

	return session, nil
}

func (e *DefaultEngine) persistSession(session *Session) error {
	sessionPath := filepath.Join(e.kdsePath(), "session.yaml")
	content := fmt.Sprintf(`session:
  id: %s
  created: %s
  ended: %s
`, session.ID, session.Created.Format(time.RFC3339), session.Ended.Format(time.RFC3339))
	return os.WriteFile(sessionPath, []byte(content), 0644)
}

// Report formatting functions
func (e *DefaultEngine) formatPhaseReport(ws *Workspace) string {
	var b strings.Builder
	b.WriteString("# Phase Report\n\n")
	b.WriteString(fmt.Sprintf("| Field | Value |\n"))
	b.WriteString(fmt.Sprintf("|-------|-------|\n"))
	b.WriteString(fmt.Sprintf("| Current Phase | %s |\n", ws.State.CurrentPhase))
	b.WriteString(fmt.Sprintf("| Phase History | %d transitions |\n", len(ws.State.PhaseHistory)))
	b.WriteString(fmt.Sprintf("| Generated | %s |\n\n", time.Now().Format(time.RFC3339)))
	return b.String()
}

func (e *DefaultEngine) formatVerificationReport(ws *Workspace) string {
	var b strings.Builder
	b.WriteString("# Verification Report\n\n")
	b.WriteString(fmt.Sprintf("| Field | Value |\n"))
	b.WriteString(fmt.Sprintf("|-------|-------|\n"))
	b.WriteString(fmt.Sprintf("| Workspace | %s |\n", ws.Root))
	b.WriteString(fmt.Sprintf("| Runtime | %s v%s |\n", ws.Runtime.Type, ws.Runtime.Version))
	b.WriteString(fmt.Sprintf("| Verified | %s |\n\n", time.Now().Format(time.RFC3339)))
	return b.String()
}

func (e *DefaultEngine) formatSummaryReport(ws *Workspace) string {
	var b strings.Builder
	b.WriteString("# Summary Report\n\n")
	b.WriteString(fmt.Sprintf("## Workspace Overview\n\n"))
	b.WriteString(fmt.Sprintf("- **Path:** %s\n", ws.Root))
	b.WriteString(fmt.Sprintf("- **Phase:** %s\n", ws.State.CurrentPhase))
	b.WriteString(fmt.Sprintf("- **Artifacts:** %d\n", ws.State.ArtifactCount))
	b.WriteString(fmt.Sprintf("- **Reports:** %d\n\n", ws.State.ReportCount))
	return b.String()
}

func (e *DefaultEngine) formatProgressReport(ws *Workspace) string {
	var b strings.Builder
	b.WriteString("# Progress Report\n\n")
	b.WriteString("## Phase Progress\n\n")

	phases := []Phase{
		PhaseInitialization,
		PhaseKnowledge,
		PhaseArchitecture,
		PhaseImplementation,
		PhaseVerification,
		PhaseReports,
	}

	currentIndex := -1
	for i, p := range phases {
		if p == ws.State.CurrentPhase {
			currentIndex = i
			break
		}
	}

	for i, p := range phases {
		check := "[ ]"
		if i < currentIndex {
			check = "[✓]"
		} else if i == currentIndex {
			check = "[→]"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", check, p))
	}

	b.WriteString("\n")
	return b.String()
}

func generateSessionID() string {
	return fmt.Sprintf("KDSE-SESSION-%s", time.Now().Format("20060102-150405"))
}

func atomicMove(src, dst string) error {
	// Ensure destination parent exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	// Check if destination exists
	if _, err := os.Stat(dst); err == nil {
		// Destination exists, remove it
		if err := os.RemoveAll(dst); err != nil {
			return err
		}
	}

	// Move using rename
	if err := os.Rename(src, dst); err != nil {
		// Fall back to copy+remove on cross-device
		return copyDir(src, dst)
	}

	return nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// Helper for JSON operations (if needed)
func toJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}
