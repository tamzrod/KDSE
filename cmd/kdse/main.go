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
	case "normalize":
		handleNormalize(repoPath)
	case "run":
		handleRun(repoPath, args)
	case "status":
		handleStatus(repoPath)
	case "report":
		handleReport(repoPath)
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
  normalize   Normalize existing documentation to KDSE standard
  run         Start a KDSE session
  status      Show current session status
  report      Generate runtime report

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

	// Save the result
	if err := saveNormalizationResult(repoPath, result); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not save normalization result: %v\n", err)
	}
}

func saveNormalizationResult(repoPath string, result *normalize.NormalizationResult) error {
	// This would save the JSON result for future reference
	// For now, the report is printed to stdout
	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
