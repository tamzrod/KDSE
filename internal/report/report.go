package report

import (
	"fmt"
	"time"

	"github.com/kdse/runtime/internal/types"
)

type Generator struct {
	repoPath string
}

func NewGenerator(repoPath string) *Generator {
	return &Generator{repoPath: repoPath}
}

func (g *Generator) Generate(state *types.SessionState) *types.RuntimeReport {
	overallScore := g.calculateOverallScore(state.Dimensions)

	report := &types.RuntimeReport{
		ReportID:       fmt.Sprintf("KDSE-RT-%s", time.Now().Format("20060102-150405")),
		SessionID:      state.SessionID,
		GeneratedAt:    time.Now(),
		Repository:     state.Repository.Path,
		Phase:          state.Phase,
		OverallScore:   overallScore,
		MaturityLevel:  types.GetMaturityLevel(overallScore),
		Dimensions:     state.Dimensions,
		Findings:       state.Findings,
		Recommendations: state.Recommendations,
	}

	if len(state.Recommendations) > 0 {
		report.NextAction = &state.Recommendations[0]
	}

	return report
}

func (g *Generator) calculateOverallScore(dimensions map[string]float64) float64 {
	if dimensions == nil || len(dimensions) == 0 {
		return 0
	}

	var total float64
	for _, score := range dimensions {
		total += score
	}
	return total / float64(len(dimensions))
}
