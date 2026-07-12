package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/kdse/runtime/internal/context"
	"github.com/kdse/runtime/internal/types"
)

type Manager struct {
	repoPath string
}

func NewManager(repoPath string) *Manager {
	return &Manager{repoPath: repoPath}
}

func (m *Manager) SaveState(ctx *context.RuntimeContext) error {
	kdseDir := filepath.Join(m.repoPath, ".kdse")
	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		return err
	}

	sessionState := &types.SessionState{
		SessionID:      ctx.SessionID,
		Repository:     *ctx.Repository,
		Phase:          ctx.Phase,
		State:          ctx.State,
		StartedAt:      ctx.Timestamp,
		UpdatedAt:      time.Now().Format(time.RFC3339),
		Artifacts:      ctx.Repository.Artifacts,
		Dimensions:     ctx.Dimensions,
		Findings:       ctx.Findings,
		Recommendations: m.generateRecommendations(ctx),
	}

	statePath := filepath.Join(kdseDir, "state.json")
	data, err := json.MarshalIndent(sessionState, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statePath, data, 0644)
}

func (m *Manager) LoadState() (*types.SessionState, error) {
	statePath := filepath.Join(m.repoPath, ".kdse", "state.json")
	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, err
	}

	var state types.SessionState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

func (m *Manager) generateRecommendations(ctx *context.RuntimeContext) []types.Recommendation {
	var recs []types.Recommendation

	if ctx.Dimensions == nil {
		return recs
	}

	lowestDim := ""
	lowestScore := 10.0
	for dim, score := range ctx.Dimensions {
		if score < lowestScore {
			lowestScore = score
			lowestDim = dim
		}
	}

	if lowestDim != "" {
		recs = append(recs, types.Recommendation{
			ID:        "KDSE-ACT-001",
			Action:    "Improve " + lowestDim + " dimension",
			Dimension: lowestDim,
			Impact:    1.0,
			Priority:  "High",
			Rationale: "This dimension has the lowest score and represents the highest-value improvement opportunity",
		})
	}

	if ctx.Repository != nil && !ctx.Repository.IsGitRepo {
		recs = append(recs, types.Recommendation{
			ID:        "KDSE-ACT-002",
			Action:    "Initialize version control",
			Dimension: "Governance",
			Impact:    0.5,
			Priority:  "Medium",
			Rationale: "Version control is fundamental to engineering governance",
		})
	}

	return recs
}

func (m *Manager) UpdateState(phase, newState string) error {
	st, err := m.LoadState()
	if err != nil {
		return err
	}

	if phase != "" {
		st.Phase = phase
	}
	if newState != "" {
		st.State = newState
	}
	st.UpdatedAt = time.Now().Format(time.RFC3339)

	kdseDir := filepath.Join(m.repoPath, ".kdse")
	statePath := filepath.Join(kdseDir, "state.json")
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statePath, data, 0644)
}
