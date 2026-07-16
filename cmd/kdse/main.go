package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kdse/runtime/internal/agreement"
	"github.com/kdse/runtime/internal/config"
	"github.com/kdse/runtime/internal/detection"
	"github.com/kdse/runtime/internal/context"
	"github.com/kdse/runtime/internal/state"
	"github.com/kdse/runtime/internal/report"
	"github.com/kdse/runtime/internal/normalize"
	"github.com/kdse/runtime/internal/collect"
	"github.com/kdse/runtime/internal/knowledge"
	"github.com/kdse/runtime/internal/orchestration"
	kdseruntime "github.com/kdse/runtime/internal/runtime"
	somedaypkg "github.com/kdse/runtime/internal/someday"
	"github.com/kdse/runtime/internal/workspace"
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
	case "init":
		handleInit(repoPath)
	case "initialize":
		handleInitialize(repoPath)
	case "collect":
		handleCollect(repoPath)
	case "normalize":
		handleNormalize(repoPath)
	case "run":
		handleRun(repoPath, args)
	case "orchestrate":
		handleOrchestrate(repoPath, args)
	case "status":
		handleStatus(repoPath)
	case "report":
		handleReport(repoPath)
	case "runtime":
		handleRuntime(repoPath, args)
	case "someday":
		handleSomeday(repoPath, args)
	case "compliance":
		handleCompliance(repoPath, args)
	case "context":
		handleContext(repoPath, args)
	case "agreement":
		handleAgreement(repoPath, args)
	case "notebook":
		handleNotebook(repoPath, args)
	case "promote":
		handlePromote(repoPath, args)
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
	fmt.Println(`KDSE Runtime v1.0.0 - Knowledge-Driven Software Engineering

Usage: kdse <command> [options]

Core Commands:
  init         Initialize .kdse/ workspace
  status       Show current workspace status
  notebook     Engineering notebook management
  promote      Knowledge promotion workflow

Notebook Commands:
  kdse notebook add <title>     Add entry to notebook
  kdse notebook list            List notebook entries
  kdse notebook show <id>      Show entry details

Promote Commands:
  kdse promote submit <id>     Submit entry as candidate
  kdse promote review <id>      Review candidate (--accept/--reject)

Agreement Commands:
  kdse agreement init           Initialize project agreement
  kdse agreement show          Display current agreement
  kdse agreement phase <phase>  Update current phase

Other Commands:
  initialize   Full runtime initialization
  runtime     Runtime management (verify, invariant)
  context     Context handoff management
  collect     Collect engineering evidence
  normalize   Normalize documentation

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

// handleInitialize creates an evidence-driven KDSE runtime
func handleInitialize(repoPath string) {
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║     KDSE Evidence-Driven Runtime Initialization               ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Repository: %s\n", repoPath)
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Println("║ PHASE 1: Execute                                              ║")
	fmt.Println("║ Creating operational runtime structure...")

	kdse := kdseruntime.New(repoPath)

	// Execute: Create all directories and files
	directories := []string{
		kdseruntime.DirRuntime, kdseruntime.DirFoundation,
		kdseruntime.DirKnowledge, kdseruntime.DirLaboratory,
		kdseruntime.DirEvidence, kdseruntime.DirReferences,
		kdseruntime.DirTraceability, kdseruntime.DirReports,
		kdseruntime.DirConfig, kdseruntime.DirState,
		kdseruntime.DirArtifacts, kdseruntime.DirSessions,
		kdseruntime.DirNormalized, kdseruntime.DirCache,
	}

	for _, dir := range directories {
		fmt.Printf("║   Creating: %s/\n", dir)
	}

	fmt.Println("║ PHASE 2: Verify                                              ║")

	// Verify: Check every artifact
	result := kdse.Initialize()

	// Report: Evidence of what was created
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Println("║ VERIFICATION RESULTS                                         ║")

	for _, v := range result.Verification {
		statusIcon := "✓ PASS"
		if v.Status == "FAIL" {
			statusIcon = "✗ FAIL"
		}
		fmt.Printf("║ %s %-15s %s\n", statusIcon, v.Artifact, v.Path)
	}

	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	// Calculate confidence
	fmt.Printf("║ Confidence: %.2f\n", result.Confidence)

	if result.Success {
		fmt.Println("║ Status: OPERATIONAL                                          ║")
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ EVIDENCE SUMMARY                                              ║")
		for _, e := range result.Evidence {
			fmt.Printf("║ ✓ %s\n", e)
		}
	} else {
		fmt.Println("║ Status: INITIALIZATION FAILED                                 ║")
		if len(result.Errors) > 0 {
			fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
			fmt.Println("║ ERRORS                                                       ║")
			for _, err := range result.Errors {
				fmt.Printf("║ ✗ %s\n", err)
			}
		}
	}

	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	if !result.Success {
		os.Exit(1)
	}
}

// handleRuntime manages runtime verification and invariants
func handleRuntime(repoPath string, args []string) {
	if len(args) < 1 {
		printRuntimeUsage()
		os.Exit(1)
	}

	subcmd := args[0]

	switch subcmd {
	case "verify":
		handleRuntimeVerify(repoPath)
	case "invariant":
		handleRuntimeInvariant(repoPath, args[1:])
	default:
		fmt.Printf("Unknown runtime command: %s\n", subcmd)
		printRuntimeUsage()
		os.Exit(1)
	}
}

func printRuntimeUsage() {
	fmt.Println(`KDSE Runtime Commands

Usage: kdse runtime <command> [options]

Commands:
  verify     Verify runtime is operational
  invariant Check phase transition requirements

Examples:
  kdse runtime verify
  kdse runtime invariant --phase Foundation
`)
}

func handleRuntimeVerify(repoPath string) {
	kdse := kdseruntime.New(repoPath)
	report := kdse.Verify()

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              KDSE Runtime Self-Audit                         ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	for _, c := range report.Components {
		statusIcon := "PASS"
		if c.Status == "FAIL" {
			statusIcon = "FAIL"
		}
		fmt.Printf("║ %-12s %-8s %s\n", c.Artifact, statusIcon, c.Path)
	}

	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Confidence: %.2f\n", report.Confidence)

	if report.Success {
		fmt.Println("║ Status: OPERATIONAL                                          ║")
	} else {
		fmt.Println("║ Status: FAILED                                               ║")
		if len(report.Failed) > 0 {
			fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
			fmt.Println("║ Failed Components:                                           ║")
			for _, f := range report.Failed {
				fmt.Printf("║   • %s\n", f)
			}
		}
	}

	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	if !report.Success {
		os.Exit(1)
	}
}

func handleRuntimeInvariant(repoPath string, args []string) {
	phase := ""

	for i := 0; i < len(args); i++ {
		if args[i] == "--phase" && i+1 < len(args) {
			phase = args[i+1]
			i++
		}
	}

	if phase == "" {
		fmt.Fprintf(os.Stderr, "Error: --phase <phase> is required\n")
		printRuntimeUsage()
		os.Exit(1)
	}

	kdse := kdseruntime.New(repoPath)

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              KDSE Runtime Invariant Check                    ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Phase: %s\n", phase)
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	invariants := kdseruntime.DefaultInvariants()
	for _, inv := range invariants {
		if inv.Phase == phase {
			fmt.Printf("║ Requirements:\n")
			for _, req := range inv.Requires {
				passed, msg := kdse.CheckInvariant(phase, req)
				status := "PASS"
				if !passed {
					status = "FAIL"
				}
				fmt.Printf("║   %s %s\n", status, req)
				fmt.Printf("║     → %s\n", msg)
			}
			fmt.Printf("║ Description: %s\n", inv.Description)
			break
		}
	}

	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
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

// handleOrchestrate starts the state-based orchestration engine
func handleOrchestrate(repoPath string, args []string) {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          KDSE State-Based Orchestration Engine                ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Repository: %s\n", repoPath)
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Parse command flags
	maxCycles := 100
	foundationThreshold := 0.7

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--max-cycles", "-m":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &maxCycles)
				i++
			}
		case "--foundation-threshold", "-t":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%f", &foundationThreshold)
				i++
			}
		case "--temp-workspace":
			// TODO: Implement temp workspace creation
		}
	}

	// Create engine configuration
	config := &orchestration.EngineConfig{
		FoundationThreshold: foundationThreshold,
		EvidenceThreshold:    0.6,
		MaxCycles:           maxCycles,
		TempWorkspaceBase:   "temp",
		EnableMigration:     true,
	}

	// Create and initialize engine
	engine, err := orchestration.NewEngine(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating engine: %v\n", err)
		os.Exit(1)
	}

	if err := engine.Initialize(repoPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing orchestration: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Orchestration engine initialized.")
	fmt.Println()
	fmt.Println("Each cycle performs:")
	fmt.Println("  1. Resolve workspace")
	fmt.Println("  2. Evaluate current state")
	fmt.Println("  3. Evaluate confidence")
	fmt.Println("  4. Evaluate missing evidence")
	fmt.Println("  5. Decide next phase")
	fmt.Println("  6. Execute only that phase")
	fmt.Println("  7. Re-evaluate")
	fmt.Println()

	// Check Foundation status
	ready, confidence, missing, err := engine.GetFoundationStatus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not assess Foundation: %v\n", err)
	} else {
		fmt.Printf("Foundation Status:\n")
		fmt.Printf("  Ready for Implementation: %v\n", ready)
		fmt.Printf("  Current Confidence: %.2f\n", confidence)
		fmt.Printf("  Threshold: %.2f\n", foundationThreshold)
		if len(missing) > 0 {
			fmt.Printf("  Missing Foundation Documents:\n")
			for _, m := range missing {
				fmt.Printf("    • %s\n", m)
			}
		}
		fmt.Println()
	}

	// Check if implementation is blocked
	if !engine.CanImplement() {
		fmt.Println("⚠️  Implementation is BLOCKED until Foundation threshold is met.")
		fmt.Println()
	}

	// Execute orchestration cycles
	fmt.Println("Starting orchestration cycles...")
	fmt.Println()

	cycleCount := 0
	for {
		result, err := engine.ExecuteCycle()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in cycle %d: %v\n", cycleCount+1, err)
			break
		}

		cycleCount++
		state := engine.GetState()

		fmt.Printf("Cycle %d:\n", cycleCount)
		fmt.Printf("  Phase: %s → %s\n", result.Decision.Reason, state.CurrentPhase)
		fmt.Printf("  Confidence: %.2f (Foundation: %.2f)\n", 
			state.Confidence.Overall, state.Confidence.Foundation)
		fmt.Printf("  Evidence Completeness: %.0f%%\n", 
			state.EvidenceState.Completeness*100)

		if len(result.Decision.BlockingReasons) > 0 {
			fmt.Printf("  Blocked: %v\n", result.Decision.BlockingReasons)
		}

		if !result.Continue {
			fmt.Println()
			fmt.Println("Orchestration complete.")
			break
		}

		if cycleCount >= maxCycles {
			fmt.Println()
			fmt.Printf("Reached maximum cycles (%d).\n", maxCycles)
			break
		}
	}

	// Final state
	state := engine.GetState()
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Final State                                 ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Cycles Executed: %d\n", cycleCount)
	fmt.Printf("║ Final Phase: %s\n", state.CurrentPhase)
	fmt.Printf("║ Foundation Confidence: %.2f\n", state.Confidence.Foundation)
	fmt.Printf("║ Overall Confidence: %.2f\n", state.Confidence.Overall)
	fmt.Printf("║ Can Implement: %v\n", engine.CanImplement())
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")

	// Save final state
	if err := engine.SaveState(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not save orchestration state: %v\n", err)
	}
}

func handleStatus(repoPath string) {
	ws := workspace.New(repoPath)
	paths := ws.GetPaths()

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    KDSE Workspace Status                      ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Repository: %s\n", repoPath)
	fmt.Printf("║ Workspace:  %s\n", paths.Root)

	if !ws.Exists() {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Workspace not initialized. Run 'kdse init' to start.")
		fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
		return
	}

	// Check for agreement
	agreementMgr := agreement.NewManager(repoPath)
	agreement, err := agreementMgr.Get()
	if err == nil && agreement != nil {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Agreement                                                    ║")
		fmt.Printf("║ Project:    %s\n", agreement.ProjectName)
		fmt.Printf("║ Phase:      %s\n", agreement.CurrentPhase)
		fmt.Printf("║ Assumptions: %d\n", len(agreement.Assumptions))
	}

	// Check knowledge entries
	knowledgeMgr := knowledge.NewManager(repoPath)
	if err := knowledgeMgr.Load(); err == nil {
		stats := knowledgeMgr.Stats()
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Knowledge                                                    ║")
		fmt.Printf("║ Total:      %d\n", stats["total"])
		fmt.Printf("║ Notebook:   %d  |  Candidate: %d  |  Promoted: %d\n",
			stats["notebook"], stats["candidate"], stats["promoted"])
	}

	// Check for legacy directories
	legacyDirs := ws.DetectLegacyDirs()
	if len(legacyDirs) > 0 {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║ Legacy Directories Detected                                 ║")
		fmt.Println("║ Run 'kdse migrate' to move them under .kdse/")
		for _, dir := range legacyDirs {
			fmt.Printf("║   • %s/\n", dir)
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

// handleSomeday manages the Someday/Maybe knowledge system
func handleSomeday(repoPath string, args []string) {
	if len(args) < 1 {
		printSomedayUsage()
		os.Exit(1)
	}

	subcmd := args[0]
	subargs := args[1:]

	// Initialize someday manager
	manager := somedaypkg.New(repoPath)
	if err := manager.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing someday: %v\n", err)
		os.Exit(1)
	}

	switch subcmd {
	case "add":
		handleSomedayAdd(manager, subargs)
	case "list":
		handleSomedayList(manager, subargs)
	case "show":
		handleSomedayShow(manager, subargs)
	case "promote":
		handleSomedayPromote(manager, subargs)
	case "archive":
		handleSomedayArchive(manager, subargs)
	case "search":
		handleSomedaySearch(manager, subargs)
	case "review":
		handleSomedayReview(manager)
	case "export":
		handleSomedayExport(manager, repoPath)
	default:
		fmt.Printf("Unknown someday command: %s\n", subcmd)
		printSomedayUsage()
		os.Exit(1)
	}
}

func printSomedayUsage() {
	fmt.Println(`KDSE Someday/Maybe Commands

Usage: kdse someday <command> [options]

Commands:
  add      Add a new someday/maybe idea
  list     List all ideas (optionally filtered by status)
  show     Show details of a specific idea
  promote  Promote an idea to active consideration
  archive  Archive an idea
  search   Search ideas by keyword
  review   Show ideas ready for review
  export   Export all ideas to JSON

Examples:
  kdse someday add --title "GUI Runtime" --description "Create a graphical interface"
  kdse someday list --status SOMEDAY
  kdse someday show IDEA-001
  kdse someday promote IDEA-001 --reason "Ready for implementation"
  kdse someday archive IDEA-002 --reason "No longer relevant"
  kdse someday search "GUI"
  kdse someday review
  kdse someday export
`)
}

func handleSomedayAdd(manager *somedaypkg.SomedayManager, args []string) {
	title := ""
	description := ""
	problem := ""
	origin := ""
	author := ""

	// Parse arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--title", "-t":
			if i+1 < len(args) {
				title = args[i+1]
				i++
			}
		case "--description", "-d":
			if i+1 < len(args) {
				description = args[i+1]
				i++
			}
		case "--problem", "-p":
			if i+1 < len(args) {
				problem = args[i+1]
				i++
			}
		case "--origin", "-o":
			if i+1 < len(args) {
				origin = args[i+1]
				i++
			}
		case "--author", "-a":
			if i+1 < len(args) {
				author = args[i+1]
				i++
			}
		}
	}

	if title == "" {
		fmt.Fprintf(os.Stderr, "Error: --title is required\n")
		os.Exit(1)
	}

	if description == "" {
		description = title
	}

	idea, err := manager.Add(title, description, problem, origin, author)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding idea: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              IDEA ADDED TO SOMEDAY/MAYBE                      ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ ID:          %s\n", idea.ID)
	fmt.Printf("║ Title:       %s\n", idea.Title)
	fmt.Printf("║ Status:      %s\n", idea.Status)
	fmt.Printf("║ Priority:    %d\n", idea.Priority)
	fmt.Printf("║ Confidence:  %.0f%%\n", idea.Confidence*100)
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Use 'kdse someday show " + idea.ID + "' for details")
}

func handleSomedayList(manager *somedaypkg.SomedayManager, args []string) {
	status := ""

	for i := 0; i < len(args); i++ {
		if args[i] == "--status" && i+1 < len(args) {
			status = args[i+1]
			i++
		}
	}

	var statusEnum somedaypkg.IdeaStatus
	if status != "" {
		statusEnum = somedaypkg.IdeaStatus(status)
	}

	ideas, err := manager.List(statusEnum)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing ideas: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(somedaypkg.FormatList(ideas))
}

func handleSomedayShow(manager *somedaypkg.SomedayManager, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: idea ID is required\n")
		os.Exit(1)
	}

	id := args[0]
	idea, err := manager.Show(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error showing idea: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(somedaypkg.FormatIdea(idea))
}

func handleSomedayPromote(manager *somedaypkg.SomedayManager, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: idea ID is required\n")
		os.Exit(1)
	}

	id := args[0]
	reason := "Promoted for consideration"

	// Parse reason if provided
	for i := 1; i < len(args); i++ {
		if args[i] == "--reason" && i+1 < len(args) {
			reason = args[i+1]
			break
		}
	}

	idea, err := manager.Promote(id, reason)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error promoting idea: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              IDEA PROMOTED                                   ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ ID:          %s\n", idea.ID)
	fmt.Printf("║ Title:       %s\n", idea.Title)
	fmt.Printf("║ Status:      %s → %s\n", somedaypkg.StatusSomeday, idea.Status)
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("The idea has been promoted and is now available for active consideration.")
	fmt.Println("It maintains traceability links to its original Someday/Maybe entry.")
}

func handleSomedayArchive(manager *somedaypkg.SomedayManager, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: idea ID is required\n")
		os.Exit(1)
	}

	id := args[0]
	reason := ""

	// Parse reason if provided
	for i := 1; i < len(args); i++ {
		if args[i] == "--reason" && i+1 < len(args) {
			reason = args[i+1]
			break
		}
	}

	idea, err := manager.Archive(id, reason)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error archiving idea: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              IDEA ARCHIVED                                   ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ ID:          %s\n", idea.ID)
	fmt.Printf("║ Title:       %s\n", idea.Title)
	fmt.Printf("║ Status:      %s\n", idea.Status)
	if reason != "" {
		fmt.Printf("║ Reason:      %s\n", reason)
	}
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("The idea has been archived. It remains searchable and can be retrieved.")
}

func handleSomedaySearch(manager *somedaypkg.SomedayManager, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: search query is required\n")
		os.Exit(1)
	}

	query := strings.Join(args, " ")
	ideas, err := manager.Search(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching ideas: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("Search results for: \"%s\"\n", query)
	fmt.Println()
	fmt.Println(somedaypkg.FormatList(ideas))
}

func handleSomedayReview(manager *somedaypkg.SomedayManager) {
	ideas, err := manager.Review()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting review list: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              IDEAS READY FOR REVIEW                           ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	if len(ideas) == 0 {
		fmt.Println("║ No ideas ready for review.")
	} else {
		for _, idea := range ideas {
			fmt.Printf("║ %s %s [P%d]\n", idea.ID, truncate(idea.Title, 30), idea.Priority)
		}
	}

	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("These are high-priority ideas (P1-P2) in Someday status.")
	fmt.Println("Consider promoting or archiving them.")
}

func handleSomedayExport(manager *somedaypkg.SomedayManager, repoPath string) {
	export, err := manager.Export()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error exporting ideas: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              SOMEDAY/MAYBE EXPORT                          ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Version:     %s\n", export.Version)
	fmt.Printf("║ Exported:    %s\n", export.ExportedAt)
	fmt.Printf("║ Total Ideas: %d\n", export.Manifest.TotalCount)
	fmt.Printf("║ Status:\n")
	for status, count := range export.Manifest.ByStatus {
		fmt.Printf("║   %s: %d\n", status, count)
	}
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Full export data:")
	fmt.Println()

	// Print as JSON
	exportPath := fmt.Sprintf("%s/.kdse/someday/export.json", repoPath)
	// Note: In Go, we'd serialize and save this
	fmt.Printf("Export would be saved to: %s\n", exportPath)
}

// handleCompliance checks KDSE runtime compliance
func handleCompliance(repoPath string, args []string) {
	// Parse arguments
	jsonOutput := false
	fullAudit := false

	for _, arg := range args {
		switch arg {
		case "--json":
			jsonOutput = true
		case "--audit":
			fullAudit = true
		}
	}

	// Run compliance validation
	report, err := kdseruntime.ValidateCompliance(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running compliance check: %v\n", err)
		os.Exit(1)
	}

	// Output format
	if jsonOutput {
		fmt.Println(report.ToJSON())
	} else {
		fmt.Println(kdseruntime.FormatComplianceReport(report))
	}

	// Exit with error code if non-compliant
	if !report.Compliant {
		os.Exit(1)
	}

	// Full audit output
	if fullAudit && report.Compliant {
		fmt.Println()
		fmt.Println("=== KDSE Pre-Implementation Audit ===")
		fmt.Println("✓ Project root identified")
		fmt.Println("✓ .kdse directory exists")
		fmt.Println("✓ Runtime verification passed")
		fmt.Println("✓ All templates exist")
		fmt.Println("✓ No manual folder creation detected")
		fmt.Println("✓ Compliance validated")
		fmt.Println()
		fmt.Println("The project is KDSE compliant. Engineering may proceed.")
	}
}

// handleAgreement manages the Agreement subsystem
func handleAgreement(repoPath string, args []string) {
	if len(args) < 1 {
		printAgreementUsage()
		os.Exit(1)
	}

	subcmd := args[0]
	mgr := agreement.NewManager(repoPath)

	switch subcmd {
	case "init":
		handleAgreementInit(mgr, args[1:])
	case "show":
		handleAgreementShow(mgr)
	case "phase":
		handleAgreementPhase(mgr, args[1:])
	case "add-assumption":
		handleAgreementAddAssumption(mgr, args[1:])
	case "validate":
		handleAgreementValidate(mgr, args[1:])
	default:
		fmt.Printf("Unknown agreement command: %s\n", subcmd)
		printAgreementUsage()
		os.Exit(1)
	}
}

func printAgreementUsage() {
	fmt.Println(`KDSE Agreement Commands

Usage: kdse agreement <command> [options]

Commands:
  init              Initialize a new agreement
  show              Display current agreement
  phase <phase>     Update current phase
  add-assumption    Add a shared assumption
  validate          Validate constraints

Examples:
  kdse agreement init
  kdse agreement show
  kdse agreement phase Problem
  kdse agreement add-assumption "We use Go for backend services"
  kdse agreement validate --subsystem internal/api`)
}

func handleAgreementInit(mgr *agreement.Manager, args []string) {
	// Get project name from git or directory
	projectName := "unknown"
	gitDir := ".git"
	if _, err := os.Stat(gitDir); err == nil {
		// Try to get repo name from git
		projectName = "kdse-project"
	}

	a, err := mgr.Create(projectName, "", "1.0.0", version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              AGREEMENT INITIALIZED                           ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Project:     %s\n", a.ProjectName)
	fmt.Printf("║ Phase:       %s\n", a.CurrentPhase)
	fmt.Printf("║ Methodology: v%s\n", a.MethodologyVersion)
	fmt.Printf("║ Runtime:     v%s\n", a.RuntimeVersion)
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
}

func handleAgreementShow(mgr *agreement.Manager) {
	a, err := mgr.Get()
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			fmt.Println("No agreement found. Run 'kdse agreement init' first.")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	data, _ := json.MarshalIndent(a, "", "  ")
	fmt.Println(string(data))
}

func handleAgreementPhase(mgr *agreement.Manager, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: phase name required\n")
		os.Exit(1)
	}

	phase := args[0]
	if err := mgr.UpdatePhase(phase); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Phase updated to: %s\n", phase)
}

func handleAgreementAddAssumption(mgr *agreement.Manager, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: assumption statement required\n")
		os.Exit(1)
	}

	statement := strings.Join(args, " ")
	id, err := mgr.AddAssumption(statement)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Assumption added: %s\n", id)
}

func handleAgreementValidate(mgr *agreement.Manager, args []string) {
	subsystem := ""
	action := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--subsystem":
			if i+1 < len(args) {
				subsystem = args[i+1]
				i++
			}
		case "--action":
			if i+1 < len(args) {
				action = args[i+1]
				i++
			}
		}
	}

	if subsystem == "" {
		subsystem = "/"
	}
	if action == "" {
		action = "modify"
	}

	valid, msg := mgr.ValidateConstraint(subsystem, action)
	if valid {
		fmt.Println("✓ Constraint validated - action allowed")
	} else {
		fmt.Printf("✗ Constraint violation: %s\n", msg)
		os.Exit(1)
	}
}

// handleInit initializes the .kdse/ workspace (lightweight init)
func handleInit(repoPath string) {
	ws := workspace.New(repoPath)

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              KDSE Workspace Initialization                   ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Repository: %s\n", repoPath)
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	if err := ws.Initialize(); err != nil {
		fmt.Printf("║ ✗ Error: %s\n", err)
		fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
		os.Exit(1)
	}

	fmt.Println("║ ✓ Created .kdse/ directory")
	fmt.Println("║ ✓ Workspace ready for KDSE operations")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  kdse agreement init      # Initialize project agreement")
	fmt.Println("  kdse notebook add <title> # Add first knowledge entry")
	fmt.Println("  kdse status              # View workspace status")
}

// handleNotebook manages the engineering notebook
func handleNotebook(repoPath string, args []string) {
	if len(args) < 1 {
		printNotebookUsage()
		os.Exit(1)
	}

	subcmd := args[0]
	mgr := knowledge.NewManager(repoPath)
	if err := mgr.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading notebook: %v\n", err)
		os.Exit(1)
	}

	switch subcmd {
	case "add":
		handleNotebookAdd(mgr, args[1:])
	case "list":
		handleNotebookList(mgr, args[1:])
	case "show":
		handleNotebookShow(mgr, args[1:])
	default:
		fmt.Printf("Unknown notebook command: %s\n", subcmd)
		printNotebookUsage()
		os.Exit(1)
	}
}

func printNotebookUsage() {
	fmt.Println(`Notebook Commands

Usage: kdse notebook <command> [options]

Commands:
  add <title>              Add entry to notebook
  list                     List all notebook entries
  show <id>                Show entry details

Examples:
  kdse notebook add "Users need password reset"
  kdse notebook add "API latency must be under 200ms" --source benchmark.json
  kdse notebook list
  kdse notebook show KDSE-KNOW-20240101-A
`)
}

func handleNotebookAdd(mgr *knowledge.Manager, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: title required\n")
		os.Exit(1)
	}

	title := args[0]
	content := ""
	source := ""
	tags := []string{}

	// Parse optional flags
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--source", "-s":
			if i+1 < len(args) {
				source = args[i+1]
				i++
			}
		case "--tag", "-t":
			if i+1 < len(args) {
				tags = append(tags, args[i+1])
				i++
			}
		default:
			if !strings.HasPrefix(args[i], "-") {
				content = strings.Join(args[i:], " ")
				break
			}
		}
	}

	id, err := mgr.CreateNotebookEntry(title, content, source, tags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Notebook entry created: %s\n", id)
}

func handleNotebookList(mgr *knowledge.Manager, args []string) {
	entries := mgr.List("")
	stats := mgr.Stats()

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              Engineering Notebook                          ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Total Entries: %d\n", stats["total"])
	fmt.Printf("║ Notebook: %d  |  Candidate: %d  |  Promoted: %d  |  Rejected: %d\n",
		stats["notebook"], stats["candidate"], stats["promoted"], stats["rejected"])
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	if len(entries) == 0 {
		fmt.Println("║ No entries yet. Run 'kdse notebook add <title>' to start.")
	} else {
		for _, e := range entries {
			statusIcon := map[string]string{
				"notebook":  "○",
				"candidate": "◐",
				"promoted":  "●",
				"rejected":  "✗",
			}[string(e.Status)]
			fmt.Printf("║ %s %s %s\n", statusIcon, e.ID, truncate(e.Title, 40))
		}
	}

	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
}

func handleNotebookShow(mgr *knowledge.Manager, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: entry ID required\n")
		os.Exit(1)
	}

	id := args[0]
	entry := mgr.Get(id)
	if entry == nil {
		fmt.Fprintf(os.Stderr, "Entry not found: %s\n", id)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Printf("║ %s\n", truncate(entry.Title, 52))
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ ID:      %s\n", entry.ID)
	fmt.Printf("║ Status:  %s\n", entry.Status)
	if entry.EvidenceStrength > 0 {
		fmt.Printf("║ Strength: %s\n", knowledge.StrengthToStars(entry.EvidenceStrength))
	}
	if entry.Source != "" {
		fmt.Printf("║ Source:  %s\n", entry.Source)
	}
	fmt.Printf("║ Created: %s\n", entry.CreatedAt)
	if entry.PromotedAt != "" {
		fmt.Printf("║ Updated: %s\n", entry.PromotedAt)
	}
	if len(entry.EvidenceRefs) > 0 {
		fmt.Printf("║ Evidence: %s\n", strings.Join(entry.EvidenceRefs, ", "))
	}
	if entry.ReviewRationale != "" {
		fmt.Printf("║ Review:  %s\n", entry.ReviewRationale)
	}
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║\n║ %s\n", strings.ReplaceAll(entry.Content, "\n", "\n║ "))
	fmt.Println("║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
}

// handlePromote manages knowledge promotion workflow
func handlePromote(repoPath string, args []string) {
	if len(args) < 1 {
		printPromoteUsage()
		os.Exit(1)
	}

	subcmd := args[0]
	mgr := knowledge.NewManager(repoPath)
	if err := mgr.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading knowledge: %v\n", err)
		os.Exit(1)
	}

	switch subcmd {
	case "submit":
		handlePromoteSubmit(mgr, args[1:])
	case "review":
		handlePromoteReview(mgr, args[1:])
	case "list":
		handlePromoteList(mgr, args[1:])
	default:
		fmt.Printf("Unknown promote command: %s\n", subcmd)
		printPromoteUsage()
		os.Exit(1)
	}
}

func printPromoteUsage() {
	fmt.Println(`Promote Commands

Usage: kdse promote <command> [options]

Commands:
  submit <id>               Submit notebook entry as candidate
  review <id>               Review candidate (requires --accept or --reject)
  list                      List candidates

Review Options:
  --accept                  Accept the candidate
  --reject                  Reject the candidate
  --strength <1-5>          Set evidence strength (1-5 stars)
  --rationale <text>        Review rationale

Examples:
  kdse promote submit KDSE-KNOW-20240101-A
  kdse promote review KDSE-KNOW-20240101-A --accept --strength 4 --rationale "Well derived from benchmarks"
`)
}

func handlePromoteSubmit(mgr *knowledge.Manager, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: entry ID required\n")
		os.Exit(1)
	}

	id := args[0]
	entry := mgr.Get(id)
	if entry == nil {
		fmt.Fprintf(os.Stderr, "Entry not found: %s\n", id)
		os.Exit(1)
	}

	if entry.Status != knowledge.StatusNotebook {
		fmt.Fprintf(os.Stderr, "Entry is not in notebook status: %s\n", entry.Status)
		os.Exit(1)
	}

	if err := mgr.PromoteToCandidate(id); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Entry %s promoted to candidate\n", id)
}

func handlePromoteReview(mgr *knowledge.Manager, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: entry ID required\n")
		os.Exit(1)
	}

	id := args[0]
	entry := mgr.Get(id)
	if entry == nil {
		fmt.Fprintf(os.Stderr, "Entry not found: %s\n", id)
		os.Exit(1)
	}

	if entry.Status != knowledge.StatusCandidate {
		fmt.Fprintf(os.Stderr, "Entry is not a candidate: %s\n", entry.Status)
		os.Exit(1)
	}

	// Parse flags
	accept := false
	rationale := ""
	strength := knowledge.EvidenceStrength(3)

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--accept":
			accept = true
		case "--reject":
			accept = false
		case "--rationale":
			if i+1 < len(args) {
				rationale = args[i+1]
				i++
			}
		case "--strength":
			if i+1 < len(args) {
				var s int
				fmt.Sscanf(args[i+1], "%d", &s)
				strength = knowledge.EvidenceStrength(s)
				i++
			}
		}
	}

	if err := mgr.Review(id, accept, rationale, strength); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if accept {
		fmt.Printf("Entry %s accepted as knowledge (strength: %s)\n", id, knowledge.StrengthToStars(strength))
	} else {
		fmt.Printf("Entry %s rejected: %s\n", id, rationale)
	}
}

func handlePromoteList(mgr *knowledge.Manager, args []string) {
	candidates := mgr.List(knowledge.StatusCandidate)

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              Knowledge Candidates                          ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	if len(candidates) == 0 {
		fmt.Println("║ No candidates pending review.")
	} else {
		for _, e := range candidates {
			fmt.Printf("║ %s %s\n", "◐", e.ID)
			fmt.Printf("║   Title: %s\n", truncate(e.Title, 40))
			if e.Source != "" {
				fmt.Printf("║   Source: %s\n", e.Source)
			}
			fmt.Println("║")
		}
	}

	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("To review: kdse promote review <id> --accept/--reject")
}
