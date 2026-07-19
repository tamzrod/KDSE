// Package selfhost implements KDSE self-hosting capabilities.
package selfhost

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SelfAssessOptions contains options for self-assessment
type SelfAssessOptions struct {
	Component   string
	Verbose     bool
	Output      string // "json" or "markdown"
	SaveReport  bool
	ReportPath  string
}

// EvolveOptions contains options for evolution
type SelfEvolveOptions struct {
	DryRun    bool
	Approve   bool
	Force     bool
	Promote   bool
	Stage     PromotionStage
}

// HandleSelfAssess handles the self-assess command
func HandleSelfAssess(repoPath string, opts *SelfAssessOptions) error {
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              KDSE Self-Assessment                            ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	// Create analyzer
	analyzer := NewAnalyzer(repoPath)

	// Run analysis
	fmt.Println("║ Collecting runtime information...")
	report, err := analyzer.Analyze()
	if err != nil {
		fmt.Fprintf(os.Stderr, "║ Error: %v\n", err)
		fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
		return err
	}

	report.Timestamp = time.Now().Format(time.RFC3339)

	// Save architecture model
	modelPath := filepath.Join(repoPath, ".kdse", "runtime", "architecture-model.json")
	if err := report.Architecture.Save(modelPath); err != nil {
		fmt.Fprintf(os.Stderr, "║ Warning: Failed to save architecture model: %v\n", err)
	} else {
		fmt.Printf("║ Architecture model saved to: %s\n", modelPath)
	}

	// Format output
	outputFormat := "markdown"
	if opts != nil && opts.Output != "" {
		outputFormat = opts.Output
	}

	output := FormatReport(report, outputFormat)

	// Save report if requested
	if opts != nil && opts.SaveReport {
		reportPath := opts.ReportPath
		if reportPath == "" {
			reportPath = filepath.Join(repoPath, ".kdse", "reports",
				fmt.Sprintf("self-assessment-%s.md", time.Now().Format("20060102-150405")))
		}

		dir := filepath.Dir(reportPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "║ Warning: Failed to create reports directory: %v\n", err)
		} else {
			if err := os.WriteFile(reportPath, []byte(output), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "║ Warning: Failed to save report: %v\n", err)
			} else {
				fmt.Printf("║ Report saved to: %s\n", reportPath)
			}
		}
	}

	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Println("║ Health Status                                               ║")
	fmt.Printf("║ Overall: %s (%.2f)\n", report.HealthStatus.Overall, report.HealthStatus.Score)
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Println("║ Health Checks                                               ║")

	for _, check := range report.HealthStatus.Checks {
		statusIcon := map[string]string{
			"PASS": "✓",
			"FAIL": "✗",
			"WARN": "⚠",
		}[check.Status]
		fmt.Printf("║ %s %-15s %s\n", statusIcon, check.Name, check.Status)
	}

	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Println("║ Architecture Summary                                        ║")
	if report.Architecture != nil && report.Architecture.Summary != nil {
		fmt.Printf("║ Total Components: %d\n", report.Architecture.Summary.TotalComponents)
		fmt.Printf("║ Total Dependencies: %d\n", report.Architecture.Summary.TotalDependencies)
		fmt.Printf("║ Max Depth: %d\n", report.Architecture.Summary.Depth)
	}

	if len(report.HealthStatus.CriticalIssues) > 0 {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Critical Issues                                              ║")
		for _, issue := range report.HealthStatus.CriticalIssues {
			fmt.Printf("║ • %s\n", issue)
		}
	}

	if len(report.Recommendations) > 0 {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Recommendations                                             ║")
		for _, rec := range report.Recommendations {
			fmt.Printf("║ → %s\n", rec)
		}
	}

	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")

	// Print full report in selected format if verbose or JSON output
	if opts != nil && (opts.Verbose || opts.Output == "json") {
		fmt.Println()
		fmt.Println("=== Full Report ===")
		fmt.Println(output)
	}

	return nil
}

// HandleEvolve handles the evolve command
func HandleEvolve(repoPath string, opts *SelfEvolveOptions) error {
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              KDSE Self-Evolution                             ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	// Create workflow
	workflow := NewEvolutionWorkflow(repoPath)

	// Start workflow
	if err := workflow.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "║ Error starting workflow: %v\n", err)
		fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
		return err
	}

	fmt.Printf("║ Workflow started: %s\n", workflow.state.ID)
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Current Phase: %s\n", workflow.state.CurrentPhase)

	// Execute phases
	phases := []WorkflowPhase{
		PhaseCollect,
		PhaseKnowledge,
		PhaseArchitecture,
		PhaseImpactAnalysis,
	}

	for _, phase := range phases {
		if workflow.state.CurrentPhase != phase {
			continue
		}

		fmt.Printf("║ Executing phase: %s\n", phase)

		if err := workflow.ExecutePhase(); err != nil {
			fmt.Fprintf(os.Stderr, "║ Error in phase %s: %v\n", phase, err)
			// Continue anyway for dry-run
		}

		// Show phase status
		phaseState := workflow.state.Phases[string(phase)]
		if phaseState != nil {
			fmt.Printf("║   Status: %s\n", phaseState.Status)
			if phaseState.CompletedAt != "" {
				fmt.Printf("║   Completed: %s\n", phaseState.CompletedAt)
			}
		}
	}

	// Handle approval
	if workflow.state.CurrentPhase == PhaseApproval {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Approval Required                                           ║")
		fmt.Println("║                                                            ║")
		fmt.Println("║ The following evidence is required before proceeding:        ║")

		if workflow.state.Assessment != nil {
			fmt.Printf("║ • Health Score: %.2f\n", workflow.state.Assessment.HealthStatus.Score)
			fmt.Printf("║ • Components: %d\n", len(workflow.state.Assessment.Architecture.Components))
		}

		if opts != nil && opts.Approve {
			workflow.Approve("cli")
			fmt.Println("║                                                            ║")
			fmt.Println("║ ✓ Approved by user                                         ║")
		} else if opts != nil && opts.DryRun {
			fmt.Println("║                                                            ║")
			fmt.Println("║ Dry-run mode - approval not granted                         ║")
		} else {
			fmt.Println("║                                                            ║")
			fmt.Println("║ To approve, run: kdse evolve --approve                      ║")
		}
	}

	// Handle implementation phase
	if workflow.state.CurrentPhase == PhaseImplementation {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Implementation Phase                                        ║")

		if opts != nil && opts.DryRun {
			fmt.Println("║ Dry-run mode - skipping implementation                      ║")
		} else {
			fmt.Println("║ Implementation would occur here                           ║")
			fmt.Println("║ (This is a placeholder for actual code changes)             ║")
		}
	}

	// Handle promotion
	if opts != nil && opts.Promote {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Promotion                                                  ║")

		manager := NewPromotionManager(repoPath)
		if err := manager.Initialize(); err != nil {
			fmt.Printf("║ Warning: Could not initialize promotion manager: %v\n", err)
		} else {
			stage := StageStaging
			if opts.Stage != "" {
				stage = opts.Stage
			}

			result := manager.Promote(stage, "Self-evolution", "cli")
			if result.Success {
				fmt.Printf("║ ✓ Promoted to: %s\n", stage)
			} else {
				fmt.Printf("║ ✗ Promotion failed: %v\n", result.Errors)
			}
		}
	}

	// Show final state
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Println("║ Workflow State                                             ║")
	fmt.Printf("║ ID: %s\n", workflow.state.ID)
	fmt.Printf("║ Current Phase: %s\n", workflow.state.CurrentPhase)
	fmt.Printf("║ Updated: %s\n", workflow.state.UpdatedAt)

	// Save workflow state
	if err := workflow.Save(); err != nil {
		fmt.Printf("║ Warning: Failed to save workflow state: %v\n", err)
	} else {
		fmt.Printf("║ State saved to: .kdse/runtime/evolution-state.json\n")
	}

	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")

	return nil
}

// HandleSelfAssessJSON handles self-assess returning JSON output
func HandleSelfAssessJSON(repoPath string) (string, error) {
	analyzer := NewAnalyzer(repoPath)
	report, err := analyzer.Analyze()
	if err != nil {
		return "", err
	}

	report.Timestamp = time.Now().Format(time.RFC3339)

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// HandleImpactAnalysis handles impact analysis for proposed changes
func HandleImpactAnalysis(repoPath string, changes []*Change) (*ImpactAnalysisResult, error) {
	// Load current architecture model
	modelPath := filepath.Join(repoPath, ".kdse", "runtime", "architecture-model.json")
	var model *ArchitectureModel

	if data, err := os.ReadFile(modelPath); err == nil {
		json.Unmarshal(data, &model)
	}

	if model == nil {
		// Generate a new model
		analyzer := NewAnalyzer(repoPath)
		report, err := analyzer.Analyze()
		if err != nil {
			return nil, err
		}
		model = report.Architecture
	}

	// Analyze changes
	impactAnalyzer := NewImpactAnalyzer(model)
	return impactAnalyzer.AnalyzeMultiple(changes), nil
}
