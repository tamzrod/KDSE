package types

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Repository struct {
	Path      string   `json:"path"`
	Name      string   `json:"name"`
	Phase     string   `json:"phase"`
	Artifacts []string `json:"artifacts"`
	IsGitRepo bool     `json:"is_git_repo"`
}

type SessionState struct {
	SessionID       string           `json:"session_id"`
	Repository      Repository       `json:"repository"`
	Phase           string           `json:"phase"`
	State           string           `json:"state"`
	StartedAt       string           `json:"started_at"`
	UpdatedAt       string           `json:"updated_at"`
	Artifacts       []string         `json:"artifacts"`
	Dimensions      map[string]float64 `json:"dimensions"`
	Findings        []Finding        `json:"findings"`
	Recommendations []Recommendation `json:"recommendations"`
}

type Finding struct {
	ID        string `json:"id"`
	Severity  string `json:"severity"`
	Dimension string `json:"dimension"`
	Title     string `json:"title"`
	Details   string `json:"details"`
}

type Recommendation struct {
	ID        string  `json:"id"`
	Action    string  `json:"action"`
	Dimension string  `json:"dimension"`
	Impact    float64 `json:"impact"`
	Priority  string  `json:"priority"`
	Rationale string  `json:"rationale"`
}

type EngineeringPhase string

const (
	PhaseConcept    EngineeringPhase = "Concept"
	PhaseDefined    EngineeringPhase = "Defined"
	PhaseStructured EngineeringPhase = "Structured"
	PhaseUsable     EngineeringPhase = "Usable"
	PhaseValidated  EngineeringPhase = "Validated"
	PhaseProven     EngineeringPhase = "Proven"
)

type RuntimeState string

const (
	StateIdle            RuntimeState = "Idle"
	StateLoading         RuntimeState = "Loading"
	StateVerifying       RuntimeState = "Verifying"
	StateAssessing       RuntimeState = "Assessing"
	StateReporting       RuntimeState = "Reporting"
	StatePendingApproval RuntimeState = "Pending Approval"
	StateImplementing    RuntimeState = "Implementing"
	StateVerifyingResult RuntimeState = "Verifying Results"
	StateComplete        RuntimeState = "Complete"
	StateClosed          RuntimeState = "Closed"
)

type RuntimeReport struct {
	ReportID       string            `json:"report_id"`
	SessionID      string            `json:"session_id"`
	GeneratedAt    time.Time         `json:"generated_at"`
	Repository     string            `json:"repository"`
	Phase          string            `json:"phase"`
	OverallScore   float64           `json:"overall_score"`
	MaturityLevel  string            `json:"maturity_level"`
	Dimensions     map[string]float64 `json:"dimensions"`
	Findings       []Finding         `json:"findings"`
	Recommendations []Recommendation `json:"recommendations"`
	NextAction     *Recommendation   `json:"next_action,omitempty"`
}

func GenerateSessionID() string {
	return time.Now().Format("KDSE-RT-2006-01-02-150405")
}

func GetMaturityLevel(score float64) string {
	switch {
	case score < 2:
		return "Concept"
	case score < 4:
		return "Defined"
	case score < 6:
		return "Structured"
	case score < 8:
		return "Usable"
	case score < 9:
		return "Validated"
	default:
		return "Proven"
	}
}

func (r *RuntimeReport) Format() string {
	var b strings.Builder

	b.WriteString("# Runtime Report\n\n")
	b.WriteString("| Field | Value |\n")
	b.WriteString("|-------|-------|\n")
	b.WriteString(fmt.Sprintf("| Report ID | %s |\n", r.ReportID))
	b.WriteString(fmt.Sprintf("| Generated | %s |\n", r.GeneratedAt.Format(time.RFC3339)))
	b.WriteString(fmt.Sprintf("| Repository | %s |\n", r.Repository))
	b.WriteString(fmt.Sprintf("| Phase | %s |\n", r.Phase))
	b.WriteString(fmt.Sprintf("| Overall Score | %.1f/10 (%s) |\n", r.OverallScore, r.MaturityLevel))
	b.WriteString(fmt.Sprintf("| Report Version | 1.0 |\n"))

	b.WriteString("\n## Current Status\n\n")
	b.WriteString(fmt.Sprintf("**Overall Score:** %.1f/10 (%s)\n\n", r.OverallScore, r.MaturityLevel))

	if len(r.Findings) > 0 {
		critical := 0
		high := 0
		for _, f := range r.Findings {
			switch f.Severity {
			case "Critical":
				critical++
			case "High":
				high++
			}
		}
		b.WriteString(fmt.Sprintf("Critical findings: %d | High findings: %d\n\n", critical, high))
	}

	b.WriteString("## Compliance Status\n\n")
	b.WriteString("### Dimension Scores\n\n")
	b.WriteString("| Dimension | Score | Status |\n")
	b.WriteString("|-----------|-------|--------|\n")

	scoreStatus := func(score float64) string {
		switch {
		case score >= 7:
			return "✅"
		case score >= 5:
			return "⚠️"
		default:
			return "❌"
		}
	}

	for dim, score := range r.Dimensions {
		b.WriteString(fmt.Sprintf("| %s | %.1f/10 | %s |\n", dim, score, scoreStatus(score)))
	}

	if len(r.Findings) > 0 {
		b.WriteString("\n## Summary of Findings\n\n")
		for _, f := range r.Findings {
			b.WriteString(fmt.Sprintf("### [%s] %s\n", f.Severity, f.Title))
			b.WriteString(fmt.Sprintf("**Dimension:** %s\n\n", f.Dimension))
			b.WriteString(fmt.Sprintf("%s\n\n", f.Details))
		}
	}

	if r.NextAction != nil {
		b.WriteString("## Highest Priority Recommendation\n\n")
		b.WriteString(fmt.Sprintf("**Action:** %s\n\n", r.NextAction.Action))
		b.WriteString(fmt.Sprintf("**Expected Impact:** +%.1f points\n\n", r.NextAction.Impact))
		b.WriteString(fmt.Sprintf("**Rationale:** %s\n\n", r.NextAction.Rationale))
	}

	b.WriteString("## Next Steps\n\n")
	b.WriteString("```\n")
	b.WriteString("[✓] Initialize\n")
	b.WriteString("[✓] Verify Standards\n")
	b.WriteString("[✓] Assess Repository\n")
	b.WriteString("[→] Await Approval\n")
	b.WriteString("[ ] Implement\n")
	b.WriteString("[ ] Verify\n")
	b.WriteString("[ ] Complete or Continue\n")
	b.WriteString("```\n\n")

	b.WriteString("---\n\n")
	b.WriteString("*Generated by KDSE Runtime v1.0.0*\n")

	return b.String()
}

func (r *RuntimeReport) Save(repoPath string) error {
	kdseDir := filepath.Join(repoPath, ".kdse", "reports")
	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		return err
	}

	reportPath := filepath.Join(kdseDir, fmt.Sprintf("%s.md", r.ReportID))
	return os.WriteFile(reportPath, []byte(r.Format()), 0644)
}
