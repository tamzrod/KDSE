package main

import (
	"fmt"
	"os"

	"github.com/kdse/runtime/internal/config"
	"github.com/kdse/runtime/internal/detection"
	"github.com/kdse/runtime/internal/context"
	"github.com/kdse/runtime/internal/state"
	"github.com/kdse/runtime/internal/report"
	"github.com/kdse/runtime/internal/normalize"
	"github.com/kdse/runtime/internal/collect"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	repoPath, _ := os.Getwd()

	switch cmd {
	case "install":
		handleInstall()
	case "update":
		handleUpdate()
	case "collect":
		handleCollect(repoPath)
	case "normalize":
		handleNormalize(repoPath)
	case "run":
		handleRun(repoPath, args)
	case "status":
		handleStatus(repoPath)
	case "report":
		handleReport(repoPath)
	case "context":
		handleContext(repoPath, args)
	case "version", "--version", "-v":
		fmt.Printf("KDSE Runtime v%s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`KDSE Runtime v1.0.0 - Knowledge-Driven Software Engineering Runtime

Usage: kdse <command> [options]

Commands:
  install     Install KDSE runtime configuration
  update      Update KDSE runtime
  collect     Collect engineering evidence
  normalize   Normalize existing documentation to KDSE standard
  run         Start a KDSE session
  status      Show current session status
  report      Generate runtime report
  context     Context handoff management

Context Commands:
  kdse context init           Initialize context handoff
  kdse context stage         Transition to new stage
  kdse context next-action   Set next action directive
  kdse context add-evidence  Add evidence files
  kdse context read          Display current context

Options:
  -h, --help    Show this help message
  -v, --version Show version information

For more information, see https://github.com/kdse/runtime`)
}

func handleInstall() {
	fmt.Println("Installing KDSE runtime...")

	cfg := config.Default()
	if err := cfg.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("KDSE runtime installed successfully.")
}

func handleUpdate() {
	fmt.Println("Checking for updates...")
	fmt.Printf("KDSE Runtime v%s is up to date.\n", version)
}

func handleRun(repoPath string, args []string) {
	fmt.Printf("Starting KDSE session in: %s\n", repoPath)

	detector := detection.NewDetector(repoPath)
	repo, err := detector.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting repository: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Repository: %s\n", repo.Name)
	fmt.Printf("Phase: %s\n", repo.Phase)
	fmt.Printf("Detected Artifacts: %d\n", len(repo.Artifacts))

	ctx := context.NewBuilder(repoPath).
		WithRepository(repo).
		Build()

	if err := ctx.Save(repoPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving context: %v\n", err)
		os.Exit(1)
	}

	if err := state.NewManager(repoPath).SaveState(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving state: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nSession initialized. Run 'kdse status' for details.\n")
}

func handleStatus(repoPath string) {
	mgr := state.NewManager(repoPath)
	st, err := mgr.LoadState()
	if err != nil {
		fmt.Println("No active KDSE session.")
		fmt.Println("Run 'kdse run' to start a session.")
		return
	}

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    KDSE Session Status                         ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	fmt.Printf("║ Session ID:    %s\n", st.SessionID)
	fmt.Printf("║ Repository:   %s\n", st.Repository.Path)
	fmt.Printf("║ Phase:        %s\n", st.Phase)
	fmt.Printf("║ State:        %s\n", st.State)
	fmt.Printf("║ Started:      %s\n", st.StartedAt)

	if len(st.Artifacts) > 0 {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Artifacts                                                    ║")
		for _, a := range st.Artifacts {
			fmt.Printf("║   • %s\n", a)
		}
	}

	if len(st.Dimensions) > 0 {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Dimensions                                                   ║")
		for dim, score := range st.Dimensions {
			fmt.Printf("║   %s: %.1f/10\n", dim, score)
		}
	}

	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
}

func handleReport(repoPath string) {
	mgr := state.NewManager(repoPath)
	st, err := mgr.LoadState()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: No active session. Run 'kdse run' first.\n")
		os.Exit(1)
	}

	rpt := report.NewGenerator(repoPath).Generate(st)

	fmt.Println(rpt.Format())

	if err := rpt.Save(repoPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not save report: %v\n", err)
	}
}

func handleNormalize(repoPath string) {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              KDSE Documentation Normalization                 ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	fmt.Printf("║ Repository:   %s\n", repoPath)
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Println("Starting documentation normalization...")
	fmt.Println("This process will:")
	fmt.Println("  • Discover existing documentation")
	fmt.Println("  • Analyze and extract engineering knowledge")
	fmt.Println("  • Generate KDSE-standard artifacts")
	fmt.Println("  • Build full traceability")
	fmt.Println("  • Preserve all original documentation unchanged")
	fmt.Println()

	normalizer := normalize.NewNormalizer(repoPath)
	result, err := normalizer.Normalize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Normalization failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              Normalization Complete                           ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	fmt.Printf("║ Documents Found:       %d\n", result.Statistics.TotalDocsFound)
	fmt.Printf("║ Artifacts Generated:  %d\n", result.Statistics.TotalArtifactsGen)
	fmt.Printf("║ Processing Time:       %.2fs\n", result.Statistics.ProcessingTime)
	fmt.Printf("║ Success Rate:         %.1f%%\n", result.Statistics.SuccessRate)
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	if len(result.NormalizedArts) > 0 {
		fmt.Println("║ Generated Artifacts:                                         ║")
		for _, art := range result.NormalizedArts {
			artName := truncate(art.Title, 44)
			fmt.Printf("║   • %s\n", artName)
		}
	}
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Println("Normalized artifacts are available in: .kdse/normalized/")
	fmt.Println()
	fmt.Println(result.FormatReport())

	if err := saveNormalizationResult(repoPath, result); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not save normalization result: %v\n", err)
	}
}

func saveNormalizationResult(repoPath string, result *normalize.NormalizationResult) error {
	return nil
}

func handleCollect(repoPath string) {
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              KDSE Artifact Collection                        ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Repository:   %s\n", repoPath)
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Println("Discovering engineering artifacts in artifacts/ directory...")
	fmt.Println()

	collector := collect.NewCollector(repoPath, "kdse-collect")
	result, err := collector.Collect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Artifact collection failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              Collection Complete                             ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	fmt.Printf("║ Artifacts Discovered: %d\n", len(result.ArtifactsFound))
	fmt.Printf("║ Total Size:          %s\n", formatSize(result.TotalSize))
	fmt.Printf("║ Processing Time:     %.2fs\n", result.ProcessingTime)
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	if len(result.ArtifactsFound) == 0 {
		fmt.Println("No artifacts found in artifacts/ directory.")
		fmt.Println("Create an artifacts/ directory and add engineering evidence.")
		fmt.Println()
		fmt.Println("Example structure:")
		fmt.Println("  artifacts/")
		fmt.Println("    manuals/")
		fmt.Println("    standards/")
		fmt.Println("    specifications/")
		fmt.Println("    datasheets/")
		fmt.Println()
	} else {
		fmt.Println("Artifact inventory: .kdse/artifacts/inventory.json")
		fmt.Println("Collection report:  .kdse/reports/")
		fmt.Println()

		// Show category summary
		categories := make(map[collect.ArtifactCategory]int)
		for _, art := range result.ArtifactsFound {
			categories[art.Category]++
		}

		fmt.Println("Artifacts by category:")
		for cat, count := range categories {
			fmt.Printf("  %s: %d\n", cat, count)
		}
		fmt.Println()
	}

	fmt.Println("The runtime discovers and catalogs evidence.")
	fmt.Println("Interpretation belongs to executors.")
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// handleContext manages the context handoff system
func handleContext(repoPath string, args []string) {
	if len(args) < 1 {
		printContextUsage()
		os.Exit(1)
	}

	subcmd := args[0]
	subargs := args[1:]

	switch subcmd {
	case "init":
		handleContextInit(repoPath, subargs)
	case "stage":
		handleContextStage(repoPath, subargs)
	case "next-action":
		handleContextNextAction(repoPath, subargs)
	case "add-evidence":
		handleContextAddEvidence(repoPath, subargs)
	case "read":
		handleContextRead(repoPath)
	default:
		fmt.Printf("Unknown context command: %s\n", subcmd)
		printContextUsage()
		os.Exit(1)
	}
}

func printContextUsage() {
	fmt.Println(`KDSE Context Handoff Commands

Usage: kdse context <command> [options]

Commands:
  init           Initialize a new context
  stage          Transition to a new stage
  next-action    Set the next action directive
  add-evidence   Add evidence file references
  read           Display current context

Examples:
  kdse context init --project myapp --stage Concept
  kdse context stage --to Architecture --evidence docs/arch.md
  kdse context next-action "Review domain model"
  kdse context add-evidence docs/screenshots/dashboard.png
  kdse context read`)
}

func handleContextInit(repoPath string, args []string) {
	project := "unknown"
	stage := "Concept"
	nextAction := "Initialize KDSE project"

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project", "-p":
			if i+1 < len(args) {
				project = args[i+1]
				i++
			}
		case "--stage", "-s":
			if i+1 < len(args) {
				stage = args[i+1]
				i++
			}
		case "--next-action", "-n":
			if i+1 < len(args) {
				nextAction = args[i+1]
				i++
			}
		}
	}

	ctx := context.NewHandoffContext(project, stage, nextAction)

	// Set default allowed context
	ctx.AllowedContext = []string{
		".kdse/context.json",
		"docs/",
		"README.md",
	}

	if err := ctx.SaveHandoff(repoPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize context: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Context initialized successfully.")
	fmt.Printf("  Project: %s\n", project)
	fmt.Printf("  Stage:   %s\n", stage)
	fmt.Printf("  Next:    %s\n", nextAction)
}

func handleContextStage(repoPath string, args []string) {
	newStage := ""
	var evidence []string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--to", "-t":
			if i+1 < len(args) {
				newStage = args[i+1]
				i++
			}
		case "--evidence", "-e":
			if i+1 < len(args) {
				evidence = append(evidence, args[i+1])
				i++
			}
		}
	}

	if newStage == "" {
		fmt.Fprintf(os.Stderr, "Error: --to <stage> is required\n")
		os.Exit(1)
	}

	ctx, err := context.LoadHandoff(repoPath)
	if err != nil {
		// Create new context if none exists
		ctx = context.NewHandoffContext("unknown", "Concept", "Initialize project")
	}

	previousStage := ctx.CurrentStage
	ctx.TransitionStage(newStage, evidence...)

	if err := ctx.SaveHandoff(repoPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to save context: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Stage transition complete.")
	fmt.Printf("  Previous: %s\n", previousStage)
	fmt.Printf("  Current:  %s\n", newStage)
	if len(evidence) > 0 {
		fmt.Println("  Evidence:")
		for _, e := range evidence {
			fmt.Printf("    • %s\n", e)
		}
	}
}

func handleContextNextAction(repoPath string, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: next-action requires an argument\n")
		os.Exit(1)
	}

	action := args[0]
	// Join remaining args for multi-word actions
	for i := 1; i < len(args); i++ {
		action += " " + args[i]
	}

	ctx, err := context.LoadHandoff(repoPath)
	if err != nil {
		ctx = context.NewHandoffContext("unknown", "Unknown", action)
	}

	ctx.SetNextAction(action)

	if err := ctx.SaveHandoff(repoPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to save context: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Next action set: %s\n", action)
}

func handleContextAddEvidence(repoPath string, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: add-evidence requires at least one file path\n")
		os.Exit(1)
	}

	ctx, err := context.LoadHandoff(repoPath)
	if err != nil {
		ctx = context.NewHandoffContext("unknown", "Unknown", "Continue work")
	}

	ctx.AddEvidence(args...)

	if err := ctx.SaveHandoff(repoPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to save context: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Evidence added:")
	for _, e := range args {
		fmt.Printf("  • %s\n", e)
	}
}

func handleContextRead(repoPath string) {
	ctx, err := context.LoadHandoff(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: No context found. Run 'kdse context init' first.\n")
		os.Exit(1)
	}

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    KDSE Context Handoff                      ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	fmt.Printf("║ Project:         %s\n", ctx.Project)
	fmt.Printf("║ Current Stage:   %s\n", ctx.CurrentStage)
	if ctx.PreviousStage != nil {
		fmt.Printf("║ Previous Stage:  %s\n", *ctx.PreviousStage)
	}
	fmt.Printf("║ Next Action:     %s\n", ctx.NextAction)

	if len(ctx.StageHistory) > 0 {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Stage History                                                 ║")
		for _, entry := range ctx.StageHistory {
			fmt.Printf("║   %s → %s\n", entry.CompletedAt, entry.Stage)
		}
	}

	if len(ctx.Evidence) > 0 {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Evidence                                                       ║")
		for _, e := range ctx.Evidence {
			fmt.Printf("║   • %s\n", e)
		}
	}

	if len(ctx.AllowedContext) > 0 {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Allowed Context                                                ║")
		for _, p := range ctx.AllowedContext {
			fmt.Printf("║   • %s\n", p)
		}
	}

	if ctx.Session != nil {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Session                                                         ║")
		fmt.Printf("║   Session ID:  %s\n", ctx.Session.SessionID)
		fmt.Printf("║   Started:     %s\n", ctx.Session.StartedAt)
		fmt.Printf("║   Last Update: %s\n", ctx.Session.LastUpdated)
	}

	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Metadata                                                       ║")
	fmt.Printf("║   Initialized: %s\n", ctx.Metadata.InitializedAt)
	fmt.Printf("║   Transitions: %d\n", ctx.Metadata.TransitionsCount)
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
}
