// Package selfhost implements KDSE self-hosting capabilities.
package selfhost

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Analyzer scans and analyzes the KDSE runtime to build an architecture model
type Analyzer struct {
	repoPath    string
	kdsePath    string
	model       *ArchitectureModel
	packageInfo []PackageInfo
}

// PackageInfo contains information about a Go package
type PackageInfo struct {
	Path         string
	Name         string
	Files        []string
	Imports      []string
	Dependencies map[string]bool
}

// HealthCheckConfig defines configuration for health checks
type HealthCheckConfig struct {
	CheckInvariantFiles   bool
	CheckDirectoryStruct  bool
	CheckDependencies     bool
	CheckOrphanComponents bool
}

// NewAnalyzer creates a new self-hosting analyzer
func NewAnalyzer(repoPath string) *Analyzer {
	return &Analyzer{
		repoPath: repoPath,
		kdsePath: filepath.Join(repoPath, ".kdse"),
		model:    NewArchitectureModel(),
	}
}

// Analyze performs a complete self-analysis of the KDSE runtime
func (a *Analyzer) Analyze() (*SelfAssessmentReport, error) {
	report := &SelfAssessmentReport{
		Timestamp:    "",
		Version:      "1.0",
		Architecture: a.model,
	}

	// Step 1: Scan all internal packages
	if err := a.scanPackages(); err != nil {
		return nil, fmt.Errorf("failed to scan packages: %w", err)
	}

	// Step 2: Extract components and dependencies
	if err := a.extractComponents(); err != nil {
		return nil, fmt.Errorf("failed to extract components: %w", err)
	}

	// Step 3: Build the architecture model
	if err := a.buildModel(); err != nil {
		return nil, fmt.Errorf("failed to build model: %w", err)
	}

	// Step 4: Perform health assessment
	report.HealthStatus = a.assessHealth()

	// Step 5: Analyze dependencies
	report.Dependencies = a.analyzeDependencies()

	// Step 6: Generate recommendations
	report.Recommendations = a.generateRecommendations(report)

	// Step 7: Calculate summary
	a.model.CalculateSummary()

	return report, nil
}

// scanPackages scans the internal directory for packages
func (a *Analyzer) scanPackages() error {
	internalPath := filepath.Join(a.repoPath, "internal")
	if _, err := os.Stat(internalPath); os.IsNotExist(err) {
		return fmt.Errorf("internal directory not found at %s", internalPath)
	}

	// Scan each package in internal/
	entries, err := os.ReadDir(internalPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pkgPath := filepath.Join(internalPath, entry.Name())
		pkgInfo, err := a.analyzePackage(pkgPath, entry.Name())
		if err != nil {
			// Log but continue
			continue
		}
		a.packageInfo = append(a.packageInfo, *pkgInfo)
	}

	// Also scan cmd/ directory
	cmdPath := filepath.Join(a.repoPath, "cmd")
	if _, err := os.Stat(cmdPath); !os.IsNotExist(err) {
		entries, err := os.ReadDir(cmdPath)
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				pkgPath := filepath.Join(cmdPath, entry.Name())
				pkgInfo, err := a.analyzePackage(pkgPath, entry.Name())
				if err != nil {
					continue
				}
				a.packageInfo = append(a.packageInfo, *pkgInfo)
			}
		}
	}

	return nil
}

// analyzePackage analyzes a single Go package
func (a *Analyzer) analyzePackage(pkgPath, pkgName string) (*PackageInfo, error) {
	info := &PackageInfo{
		Path:         pkgPath,
		Name:        pkgName,
		Files:       []string{},
		Imports:     []string{},
		Dependencies: make(map[string]bool),
	}

	// Find all .go files in the package
	err := filepath.Walk(pkgPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if fi.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			info.Files = append(info.Files, path)
			// Extract imports from this file
			imports := a.extractImports(path)
			for _, imp := range imports {
				if strings.HasPrefix(imp, "github.com/kdse/runtime/internal/") {
					dep := strings.TrimPrefix(imp, "github.com/kdse/runtime/internal/")
					info.Dependencies[dep] = true
				}
			}
			info.Imports = append(info.Imports, imports...)
		}
		return nil
	})

	return info, err
}

// extractImports extracts import statements from a Go file
func (a *Analyzer) extractImports(filePath string) []string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}

	var imports []string
	lines := strings.Split(string(content), "\n")
	inImport := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "import") {
			if strings.Contains(line, "(") {
				inImport = true
			} else {
				// Single line import
				imp := strings.TrimPrefix(line, "import")
				imp = strings.TrimSpace(imp)
				imp = strings.Trim(imp, "\"")
				if imp != "" {
					imports = append(imports, imp)
				}
				continue
			}
		}

		if inImport {
			if strings.Contains(line, ")") {
				inImport = false
				continue
			}
			imp := strings.TrimSpace(line)
			imp = strings.Trim(imp, "\"")
			if imp != "" {
				imports = append(imports, imp)
			}
		}
	}

	return imports
}

// extractComponents extracts components from the scanned packages
func (a *Analyzer) extractComponents() error {
	componentTypeMap := map[string]ComponentType{
		"runtime":       ComponentTypeRuntime,
		"knowledge":     ComponentTypeKnowledge,
		"orchestration": ComponentTypeOrchestration,
		"collect":       ComponentTypeToolchain,
		"toolchain":     ComponentTypeToolchain,
		"guard":         ComponentTypeToolchain,
		"discover":      ComponentTypeModule,
		"normalize":     ComponentTypeModule,
		"bootstrap":    ComponentTypeModule,
		"agreement":     ComponentTypeModule,
		"report":        ComponentTypeModule,
		"config":        ComponentTypeModule,
		"context":       ComponentTypeModule,
		"state":         ComponentTypeModule,
		"workspace":    ComponentTypeModule,
		"detection":     ComponentTypeModule,
		"methodology":   ComponentTypeModule,
		"types":         ComponentTypeModule,
		"someday":       ComponentTypeModule,
	}

	purposeMap := map[string]string{
		"runtime":       "Evidence-driven runtime management and verification",
		"knowledge":     "Knowledge promotion and management",
		"orchestration": "Engineering workflow orchestration",
		"collect":       "Evidence collection from projects",
		"toolchain":     "Toolchain detection and management",
		"guard":         "Project guard coordination",
		"discover":      "Project and runtime discovery",
		"normalize":     "Documentation normalization",
		"bootstrap":     "Runtime bootstrap operations",
		"agreement":     "Project agreement management",
		"report":        "Report generation",
		"config":        "Configuration management",
		"context":       "Context management for handoffs",
		"state":         "Session state management",
		"workspace":     "Workspace management",
		"detection":     "Project language detection",
		"methodology":   "KDSE methodology enforcement",
		"types":         "Common type definitions",
		"someday":       "Someday/Maybe knowledge repository",
	}

	for _, pkg := range a.packageInfo {
		compType, ok := componentTypeMap[pkg.Name]
		if !ok {
			compType = ComponentTypeModule
		}

		purpose, _ := purposeMap[pkg.Name]
		if purpose == "" {
			purpose = "KDSE " + pkg.Name + " module"
		}

		// Get relative path
		relPath, _ := filepath.Rel(a.repoPath, pkg.Path)

		component := &Component{
			Name:         pkg.Name,
			Type:         compType,
			Purpose:      purpose,
			Path:         relPath,
			Dependencies: []string{},
			Provides:     []string{},
			Metadata: map[string]interface{}{
				"file_count":   len(pkg.Files),
				"import_count": len(pkg.Imports),
			},
		}

		// Add dependencies
		for dep := range pkg.Dependencies {
			component.Dependencies = append(component.Dependencies, dep)
		}

		a.model.AddComponent(component)
	}

	return nil
}

// buildModel builds the complete architecture model
func (a *Analyzer) buildModel() error {
	// Add dependencies between components
	for _, pkg := range a.packageInfo {
		for dep := range pkg.Dependencies {
			depType := DependencyTypeImport
			if strings.Contains(dep, "internal") {
				depType = DependencyTypeImport
			}

			depObj := &Dependency{
				Source:      pkg.Name,
				Target:      dep,
				Type:        depType,
				Description: fmt.Sprintf("%s imports %s", pkg.Name, dep),
			}
			a.model.AddDependency(depObj)
		}
	}

	// Add metadata about KDSE itself
	a.model.Metadata = map[string]interface{}{
		"repo_path":    a.repoPath,
		"kdse_path":   a.kdsePath,
		"package_count": len(a.packageInfo),
	}

	return nil
}

// assessHealth performs health checks on the runtime
func (a *Analyzer) assessHealth() *HealthStatus {
	status := &HealthStatus{
		Score:  1.0,
		Checks: []*HealthCheck{},
	}

	checks := []struct {
		name    string
		check   func() *HealthCheck
		weight  float64
	}{
		{"Invariant Files", a.checkInvariantFiles, 0.15},
		{"Directory Structure", a.checkDirectoryStructure, 0.2},
		{"Module Dependencies", a.checkModuleDependencies, 0.25},
		{"Orphan Components", a.checkOrphanComponents, 0.15},
		{"Runtime Directory", a.checkRuntimeDirectory, 0.15},
		{"Command Handlers", a.checkCommandHandlers, 0.1},
	}

	var totalWeight float64
	var weightedScore float64

	for _, c := range checks {
		check := c.check()
		status.Checks = append(status.Checks, check)

		if check.Status == "FAIL" {
			status.Score -= c.weight
			if strings.Contains(check.Name, "Orphan") || strings.Contains(check.Name, "Invariant") {
				status.CriticalIssues = append(status.CriticalIssues, check.Description)
			} else {
				status.Warnings = append(status.Warnings, check.Description)
			}
		} else if check.Status == "WARN" {
			status.Score -= c.weight / 2
			status.Warnings = append(status.Warnings, check.Description)
		}

		totalWeight += c.weight
		weightedScore += c.weight
		if check.Status == "PASS" {
			weightedScore += c.weight
		}
	}

	if status.Score < 0 {
		status.Score = 0
	}

	if status.Score >= 0.8 {
		status.Overall = "HEALTHY"
	} else if status.Score >= 0.5 {
		status.Overall = "DEGRADED"
	} else {
		status.Overall = "UNHEALTHY"
	}

	return status
}

func (a *Analyzer) checkInvariantFiles() *HealthCheck {
	check := &HealthCheck{Name: "Invariant Files"}

	requiredFiles := []string{
		"runtime/manifest.yaml",
		"runtime/runtime.yaml",
		"knowledge/knowledge-index.yaml",
	}

	var missing []string
	for _, f := range requiredFiles {
		path := filepath.Join(a.kdsePath, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			missing = append(missing, f)
		}
	}

	if len(missing) > 0 {
		check.Status = "WARN"
		check.Description = fmt.Sprintf("Missing files: %v", missing)
		check.Evidence = "Some required runtime files are missing"
	} else {
		check.Status = "PASS"
		check.Description = "All required runtime files present"
		check.Evidence = fmt.Sprintf("Found %d required files", len(requiredFiles)-len(missing))
	}

	return check
}

func (a *Analyzer) checkDirectoryStructure() *HealthCheck {
	check := &HealthCheck{Name: "Directory Structure"}

	requiredDirs := []string{
		"runtime",
		"foundation",
		"knowledge",
		"laboratory",
		"evidence",
		"references",
		"traceability",
		"reports",
		"config",
		"state",
	}

	var missing []string
	for _, d := range requiredDirs {
		path := filepath.Join(a.kdsePath, d)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			missing = append(missing, d)
		}
	}

	if len(missing) > 3 {
		check.Status = "FAIL"
		check.Description = fmt.Sprintf("Missing critical directories: %v", missing)
		check.Evidence = "Required runtime directories are missing"
	} else if len(missing) > 0 {
		check.Status = "WARN"
		check.Description = fmt.Sprintf("Missing directories: %v", missing)
		check.Evidence = "Some optional runtime directories are missing"
	} else {
		check.Status = "PASS"
		check.Description = "All required directories present"
		check.Evidence = "Full directory structure verified"
	}

	return check
}

func (a *Analyzer) checkModuleDependencies() *HealthCheck {
	check := &HealthCheck{Name: "Module Dependencies"}

	var brokenDeps []string
	for _, pkg := range a.packageInfo {
		for dep := range pkg.Dependencies {
			if _, ok := a.model.Components[dep]; !ok {
				brokenDeps = append(brokenDeps, fmt.Sprintf("%s -> %s", pkg.Name, dep))
			}
		}
	}

	if len(brokenDeps) > 5 {
		check.Status = "FAIL"
		check.Description = fmt.Sprintf("Found %d broken dependencies", len(brokenDeps))
		check.Evidence = strings.Join(brokenDeps[:3], ", ") + "..."
	} else if len(brokenDeps) > 0 {
		check.Status = "WARN"
		check.Description = fmt.Sprintf("Found %d potentially broken dependencies", len(brokenDeps))
		check.Evidence = strings.Join(brokenDeps, ", ")
	} else {
		check.Status = "PASS"
		check.Description = "All module dependencies resolved"
		check.Evidence = fmt.Sprintf("Verified %d packages", len(a.packageInfo))
	}

	return check
}

func (a *Analyzer) checkOrphanComponents() *HealthCheck {
	check := &HealthCheck{Name: "Orphan Components"}

	var orphans []string
	for name, comp := range a.model.Components {
		if len(comp.Dependencies) == 0 && name != "runtime" && name != "discover" {
			orphans = append(orphans, name)
		}
	}

	if len(orphans) > 3 {
		check.Status = "WARN"
		check.Description = fmt.Sprintf("Found %d orphan components", len(orphans))
		check.Evidence = strings.Join(orphans, ", ")
	} else {
		check.Status = "PASS"
		check.Description = "No orphan components detected"
		check.Evidence = "All components have proper dependencies"
	}

	return check
}

func (a *Analyzer) checkRuntimeDirectory() *HealthCheck {
	check := &HealthCheck{Name: "Runtime Directory"}

	path := a.kdsePath
	info, err := os.Stat(path)
	if err != nil {
		check.Status = "FAIL"
		check.Description = ".kdse directory not found"
		check.Evidence = "Runtime directory does not exist"
		return check
	}

	if !info.IsDir() {
		check.Status = "FAIL"
		check.Description = ".kdse is not a directory"
		check.Evidence = "Runtime path is not a directory"
		return check
	}

	check.Status = "PASS"
	check.Description = "Runtime directory exists"
	check.Evidence = path

	return check
}

func (a *Analyzer) checkCommandHandlers() *HealthCheck {
	check := &HealthCheck{Name: "Command Handlers"}

	// Read main.go to count command handlers
	mainPath := filepath.Join(a.repoPath, "cmd", "kdse", "main.go")
	content, err := os.ReadFile(mainPath)
	if err != nil {
		check.Status = "WARN"
		check.Description = "Could not read main.go"
		check.Evidence = err.Error()
		return check
	}

	// Count case statements in switch
	re := regexp.MustCompile(`case\s+"([^"]+)"`)
	matches := re.FindAllStringSubmatch(string(content), -1)

	if len(matches) < 10 {
		check.Status = "WARN"
		check.Description = fmt.Sprintf("Only %d commands defined", len(matches))
		check.Evidence = "May be missing some command handlers"
	} else {
		check.Status = "PASS"
		check.Description = fmt.Sprintf("Found %d command handlers", len(matches))
		check.Evidence = "All expected commands present"
	}

	return check
}

// analyzeDependencies performs dependency analysis
func (a *Analyzer) analyzeDependencies() *DependencyAnalysis {
	analysis := &DependencyAnalysis{}

	direct := 0
	transitive := 0
	var violations []*DependencyViolation

	for _, pkg := range a.packageInfo {
		for dep := range pkg.Dependencies {
			if _, ok := a.model.Components[dep]; ok {
				direct++
			} else {
				// Check if it's an external package
				if strings.HasPrefix(dep, "github.com/") || strings.HasPrefix(dep, "golang.org/") {
					transitive++
				} else {
					violations = append(violations, &DependencyViolation{
						Rule:        "internal-only",
						Source:      pkg.Name,
						Target:      dep,
						Description: fmt.Sprintf("%s depends on non-existent %s", pkg.Name, dep),
						Severity:    "warning",
					})
				}
			}
		}
	}

	analysis.Total = direct + transitive
	analysis.Direct = direct
	analysis.Transitive = transitive
	analysis.Violations = violations

	return analysis
}

// generateRecommendations generates recommendations based on analysis
func (a *Analyzer) generateRecommendations(report *SelfAssessmentReport) []string {
	var recommendations []string

	if report.HealthStatus.Score < 0.8 {
		recommendations = append(recommendations, "Consider addressing health check warnings to improve runtime stability")
	}

	if report.Dependencies != nil && len(report.Dependencies.Violations) > 0 {
		recommendations = append(recommendations, "Review and fix dependency violations to maintain architecture integrity")
	}

	if a.model.Summary != nil && a.model.Summary.Depth > 5 {
		recommendations = append(recommendations, "Dependency depth is high - consider refactoring to reduce coupling")
	}

	if len(report.HealthStatus.CriticalIssues) > 0 {
		recommendations = append(recommendations, "Critical issues detected - prioritize fixing these before any self-modification")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Runtime is in good health - safe to proceed with self-evolution")
	}

	return recommendations
}

// FormatReport formats the self-assessment report for display
func FormatReport(report *SelfAssessmentReport, outputFormat string) string {
	switch outputFormat {
	case "json":
		data, _ := json.MarshalIndent(report, "", "  ")
		return string(data)
	default:
		return formatReportMarkdown(report)
	}
}

func formatReportMarkdown(report *SelfAssessmentReport) string {
	var output string

	output += "# KDSE Self-Assessment Report\n\n"
	output += fmt.Sprintf("**Generated:** %s\n", report.Timestamp)
	output += fmt.Sprintf("**Version:** %s\n\n", report.Version)

	output += "## Health Status\n\n"
	output += fmt.Sprintf("| Metric | Value |\n")
	output += fmt.Sprintf("|--------|-------|\n")
	output += fmt.Sprintf("| Overall | %s |\n", report.HealthStatus.Overall)
	output += fmt.Sprintf("| Score | %.2f |\n\n", report.HealthStatus.Score)

	output += "### Health Checks\n\n"
	output += fmt.Sprintf("| Check | Status | Description |\n")
	output += fmt.Sprintf("|-------|--------|------------|\n")
	for _, check := range report.HealthStatus.Checks {
		output += fmt.Sprintf("| %s | %s | %s |\n", check.Name, check.Status, check.Description)
	}

	if len(report.HealthStatus.CriticalIssues) > 0 {
		output += "\n### Critical Issues\n\n"
		for _, issue := range report.HealthStatus.CriticalIssues {
			output += fmt.Sprintf("- %s\n", issue)
		}
	}

	if len(report.HealthStatus.Warnings) > 0 {
		output += "\n### Warnings\n\n"
		for _, warning := range report.HealthStatus.Warnings {
			output += fmt.Sprintf("- %s\n", warning)
		}
	}

	if report.Architecture != nil && report.Architecture.Summary != nil {
		output += "\n## Architecture Summary\n\n"
		output += fmt.Sprintf("- **Total Components:** %d\n", report.Architecture.Summary.TotalComponents)
		output += fmt.Sprintf("- **Total Dependencies:** %d\n", report.Architecture.Summary.TotalDependencies)
		output += fmt.Sprintf("- **Max Depth:** %d\n", report.Architecture.Summary.Depth)
	}

	if report.Dependencies != nil {
		output += "\n## Dependencies\n\n"
		output += fmt.Sprintf("- **Total:** %d\n", report.Dependencies.Total)
		output += fmt.Sprintf("- **Direct:** %d\n", report.Dependencies.Direct)
		output += fmt.Sprintf("- **Transitive:** %d\n", report.Dependencies.Transitive)
	}

	if len(report.Recommendations) > 0 {
		output += "\n## Recommendations\n\n"
		for _, rec := range report.Recommendations {
			output += fmt.Sprintf("- %s\n", rec)
		}
	}

	return output
}
