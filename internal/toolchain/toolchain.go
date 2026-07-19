// Package toolchain provides automatic toolchain management for KDSE.
//
// KDSE Runtime owns the responsibility of ensuring all required development
// toolchains are available before entering the implementation phase.
//
// Key Principles:
// - KDSE never silently skips verification due to missing tooling
// - Toolchains are installed into writable workspace locations only
// - PATH and environment variables are updated automatically
// - Verification fails loudly instead of silently skipping
package toolchain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Toolchain represents a development toolchain (Go, Node.js, Python, etc.)
type Toolchain struct {
	Name            string
	Type            ToolchainType
	Version         string
	InstallPath     string
	IsAvailable     bool
	IsInstalled     bool
	DetectionMethod string
	VersionCommand  string
	InstallCommands []string
}

// ToolchainType identifies the type of toolchain
type ToolchainType string

const (
	ToolchainGo     ToolchainType = "go"
	ToolchainNode   ToolchainType = "node"
	ToolchainPython ToolchainType = "python"
	ToolchainJava   ToolchainType = "java"
	ToolchainDotNet ToolchainType = "dotnet"
	ToolchainRust   ToolchainType = "rust"
	ToolchainRuby   ToolchainType = "ruby"
	ToolchainPHP    ToolchainType = "php"
)

// VerificationResult represents the result of toolchain verification
type VerificationResult struct {
	Success      bool
	Toolchain    *Toolchain
	Errors       []string
	Evidence     []string
	InstalledAt  string
	VerifiedAt   string
}

// InstallResult represents the result of toolchain installation
type InstallResult struct {
	Success     bool
	Toolchain   *Toolchain
	InstallPath string
	Errors      []string
	Output      string
}

// Manager handles toolchain detection, installation, and verification
type Manager struct {
	installBase string
	projectPath string
	tools      map[ToolchainType]*Toolchain
}

// DefaultInstallBase returns the default installation base path
func DefaultInstallBase() string {
	return filepath.Join(os.Getenv("HOME"), ".kdse", "tools")
}

// NewManager creates a new toolchain manager
func NewManager(projectPath string) *Manager {
	installBase := os.Getenv("KDSE_TOOLS_PATH")
	if installBase == "" {
		installBase = DefaultInstallBase()
	}

	return &Manager{
		installBase: installBase,
		projectPath: projectPath,
		tools:      make(map[ToolchainType]*Toolchain),
	}
}

// SetInstallBase sets the base path for toolchain installations
func (m *Manager) SetInstallBase(path string) {
	m.installBase = path
}

// DetectAll detects all available toolchains in the system
func (m *Manager) DetectAll() []*Toolchain {
	var result []*Toolchain

	// Pre-defined toolchains to detect
	toolchains := []Toolchain{
		{Name: "Go", Type: ToolchainGo, VersionCommand: "version", DetectionMethod: "go version"},
		{Name: "Node.js", Type: ToolchainNode, VersionCommand: "--version", DetectionMethod: "node --version"},
		{Name: "Python", Type: ToolchainPython, VersionCommand: "--version", DetectionMethod: "python3 --version"},
		{Name: "Java", Type: ToolchainJava, VersionCommand: "-version", DetectionMethod: "java -version"},
		{Name: ".NET", Type: ToolchainDotNet, VersionCommand: "--version", DetectionMethod: "dotnet --version"},
		{Name: "Rust", Type: ToolchainRust, VersionCommand: "--version", DetectionMethod: "rustc --version"},
	}

	for _, tc := range toolchains {
		detected := m.detect(&tc)
		m.tools[tc.Type] = detected
		result = append(result, detected)
	}

	return result
}

// Detect finds a specific toolchain
func (m *Manager) Detect(toolchainType ToolchainType) *Toolchain {
	if tc, ok := m.tools[toolchainType]; ok {
		return tc
	}

	tc := m.createToolchain(toolchainType)
	detected := m.detect(tc)
	m.tools[toolchainType] = detected
	return detected
}

// Verify checks if a toolchain is available and working
func (m *Manager) Verify(toolchainType ToolchainType) *VerificationResult {
	tc := m.Detect(toolchainType)

	result := &VerificationResult{
		Toolchain: tc,
		Success:   false,
	}

	if !tc.IsAvailable {
		result.Errors = append(result.Errors, fmt.Sprintf("%s is not available in PATH", tc.Name))
		return result
	}

	// Verify by running version command
	cmd := exec.Command(m.getCommand(tc.Type), strings.Fields(tc.VersionCommand)...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to verify %s: %v", tc.Name, err))
		result.Errors = append(result.Errors, string(output))
		return result
	}

	result.Success = true
	result.VerifiedAt = tc.Version
	result.Evidence = append(result.Evidence, fmt.Sprintf("%s is available: %s", tc.Name, strings.TrimSpace(string(output))))
	return result
}

// Ensure verifies a toolchain is available, installing if necessary
func (m *Manager) Ensure(toolchainType ToolchainType) *VerificationResult {
	// First check if available
	result := m.Verify(toolchainType)
	if result.Success {
		return result
	}

	// Try to install
	installResult := m.Install(toolchainType)
	if !installResult.Success {
		result.Errors = append(result.Errors, installResult.Errors...)
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to install %s", toolchainType))
		return result
	}

	// Re-verify after installation
	return m.Verify(toolchainType)
}

// Install installs a toolchain to the writable location
func (m *Manager) Install(toolchainType ToolchainType) *InstallResult {
	tc := m.Detect(toolchainType)

	result := &InstallResult{
		Toolchain:   tc,
		InstallPath:  filepath.Join(m.installBase, tc.Name),
	}

	if tc.IsAvailable {
		result.Success = true
		result.Output = fmt.Sprintf("%s is already available", tc.Name)
		return result
	}

	// Install based on toolchain type
	var err error
	var output string

	switch toolchainType {
	case ToolchainGo:
		output, err = m.installGo(tc, result.InstallPath)
	case ToolchainNode:
		output, err = m.installNode(tc, result.InstallPath)
	case ToolchainPython:
		output, err = m.installPython(tc, result.InstallPath)
	default:
		err = fmt.Errorf("automatic installation for %s is not yet supported", tc.Name)
	}

	result.Output = output
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
		return result
	}

	// Update PATH
	if err := m.updatePATH(result.InstallPath); err != nil {
		result.Errors = append(result.Errors, err.Error())
		return result
	}

	result.Success = true
	tc.IsInstalled = true
	tc.InstallPath = result.InstallPath

	return result
}

// GetRequiredToolchains determines which toolchains are required for a project
func (m *Manager) GetRequiredToolchains() []ToolchainType {
	var required []ToolchainType

	// Check project files to determine required toolchains
	projectFiles := map[ToolchainType][]string{
		ToolchainGo:     {"go.mod", "go.sum"},
		ToolchainNode:   {"package.json", "package-lock.json", "yarn.lock"},
		ToolchainPython: {"pyproject.toml", "setup.py", "requirements.txt", "Pipfile"},
		ToolchainJava:   {"pom.xml", "build.gradle", "build.gradle.kts"},
		ToolchainDotNet: {".sln", ".csproj"},
		ToolchainRust:   {"Cargo.toml", "Cargo.lock"},
	}

	for toolchainType, files := range projectFiles {
		for _, file := range files {
			path := filepath.Join(m.projectPath, file)
			if _, err := os.Stat(path); err == nil {
				required = append(required, toolchainType)
				break
			}
		}
	}

	return required
}

// EnsureProjectRequirements verifies all toolchains required by the project
func (m *Manager) EnsureProjectRequirements() []*VerificationResult {
	var results []*VerificationResult

	required := m.GetRequiredToolchains()
	if len(required) == 0 {
		return results
	}

	for _, tcType := range required {
		result := m.Ensure(tcType)
		results = append(results, result)
	}

	return results
}

// UpdateEnvironmentUpdates exports the environment variables needed for toolchains
func (m *Manager) UpdateEnvironment() map[string]string {
	env := make(map[string]string)

	// Add toolchain bin directories to PATH
	pathEntries := []string{
		filepath.Join(m.installBase, "go", "bin"),
		filepath.Join(m.installBase, "node", "bin"),
	}

	currentPath := os.Getenv("PATH")
	for _, entry := range pathEntries {
		if _, err := os.Stat(entry); err == nil {
			env["PATH"] = entry + string(os.PathListSeparator) + currentPath
		}
	}

	// Set GOROOT if Go is installed
	goRoot := filepath.Join(m.installBase, "go")
	if _, err := os.Stat(goRoot); err == nil {
		env["GOROOT"] = goRoot
		env["GOPATH"] = filepath.Join(m.installBase, "gopath")
	}

	return env
}

// GetInstallBase returns the current installation base path
func (m *Manager) GetInstallBase() string {
	return m.installBase
}

// createToolchain creates a new toolchain definition
func (m *Manager) createToolchain(tcType ToolchainType) *Toolchain {
	tc := &Toolchain{Type: tcType}

	switch tcType {
	case ToolchainGo:
		tc.Name = "Go"
		tc.VersionCommand = "version"
		tc.InstallCommands = []string{
			"https://go.dev/dl/go{VERSION}.{OS}-{ARCH}.tar.gz",
		}
	case ToolchainNode:
		tc.Name = "Node.js"
		tc.VersionCommand = "--version"
		tc.InstallCommands = []string{
			"https://nodejs.org/dist/v{VERSION}/node-v{VERSION}-{OS}-{ARCH}.tar.gz",
		}
	case ToolchainPython:
		tc.Name = "Python"
		tc.VersionCommand = "--version"
		tc.InstallCommands = []string{
			"https://www.python.org/ftp/python/{VERSION}/Python-{VERSION}.tgz",
		}
	case ToolchainJava:
		tc.Name = "Java"
		tc.VersionCommand = "-version"
	case ToolchainDotNet:
		tc.Name = ".NET"
		tc.VersionCommand = "--version"
	case ToolchainRust:
		tc.Name = "Rust"
		tc.VersionCommand = "--version"
	}

	return tc
}

// detect checks if a toolchain is available
func (m *Manager) detect(tc *Toolchain) *Toolchain {
	cmd := m.getCommand(tc.Type)
	if cmd == "" {
		tc.IsAvailable = false
		return tc
	}

	// Check if command exists
	execCmd := exec.Command(cmd, strings.Fields(tc.VersionCommand)...)
	output, err := execCmd.CombinedOutput()

	if err != nil {
		tc.IsAvailable = false
		return tc
	}

	tc.IsAvailable = true
	tc.Version = strings.TrimSpace(string(output))

	// Find the actual path
	whichCmd := exec.Command("which", cmd)
	whichOutput, _ := whichCmd.Output()
	tc.InstallPath = strings.TrimSpace(string(whichOutput))

	return tc
}

// getCommand returns the command name for a toolchain type
func (m *Manager) getCommand(tcType ToolchainType) string {
	commands := map[ToolchainType]string{
		ToolchainGo:     "go",
		ToolchainNode:   "node",
		ToolchainPython: "python3",
		ToolchainJava:   "java",
		ToolchainDotNet: "dotnet",
		ToolchainRust:   "rustc",
		ToolchainRuby:   "ruby",
		ToolchainPHP:    "php",
	}

	if cmd, ok := commands[tcType]; ok {
		return cmd
	}
	return ""
}

// updatePATH adds a directory to PATH
func (m *Manager) updatePATH(dir string) error {
	// Ensure the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create toolchain directory: %w", err)
	}

	// Prepend to PATH in a persistent way (for the current session)
	currentPath := os.Getenv("PATH")
	newPath := dir + string(os.PathListSeparator) + currentPath

	// Update the PATH for the current process
	// Note: This only affects the current process
	// For persistence, the caller should export this in their shell
	_ = os.Setenv("PATH", newPath)
	return nil
}

// installGo downloads and installs Go
func (m *Manager) installGo(tc *Toolchain, installPath string) (string, error) {
	// Determine OS and architecture
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Map Go architectures
	archMap := map[string]string{
		"amd64": "amd64",
		"386":   "386",
		"arm64": "arm64",
	}
	arch, ok := archMap[goarch]
	if !ok {
		arch = "amd64"
	}

	// Map Go OS names
	osMap := map[string]string{
		"linux":   "linux",
		"darwin":  "darwin",
		"windows": "windows",
	}
	os, ok := osMap[goos]
	if !ok {
		os = "linux"
	}

	// Default Go version
	version := "1.21.0"
	url := fmt.Sprintf("https://go.dev/dl/go%s.%s-%s.tar.gz", version, os, arch)

	return m.downloadAndExtract(url, filepath.Join(installPath, "go"))
}

// installNode downloads and installs Node.js
func (m *Manager) installNode(tc *Toolchain, installPath string) (string, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Determine Node.js version and URL based on platform
	var url string
	switch goos {
	case "linux":
		if goarch == "amd64" || goarch == "arm64" {
			url = fmt.Sprintf("https://nodejs.org/dist/v20.10.0/node-v20.10.0-linux-%s.tar.xz", goarch)
		}
	case "darwin":
		url = "https://nodejs.org/dist/v20.10.0/node-v20.10.0-darwin-x64.tar.gz"
	}

	if url == "" {
		return "", fmt.Errorf("unsupported platform for Node.js: %s/%s", goos, goarch)
	}

	return m.downloadAndExtract(url, filepath.Join(installPath, "node"))
}

// installPython downloads and installs Python
func (m *Manager) installPython(tc *Toolchain, installPath string) (string, error) {
	goos := runtime.GOOS

	// Python installation typically requires compilation
	// For now, provide instructions
	if goos == "linux" {
		return "", fmt.Errorf("please install Python via your package manager: apt install python3 python3-pip")
	}

	return "", fmt.Errorf("automatic Python installation is not yet supported for %s", goos)
}

// downloadAndExtract downloads and extracts a tarball
func (m *Manager) downloadAndExtract(url, destPath string) (string, error) {
	// Create destination directory
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Download using curl or wget
	var cmd *exec.Cmd
	if _, err := exec.LookPath("curl"); err == nil {
		cmd = exec.Command("curl", "-fsSL", "-o", "/tmp/toolchain.tar.gz", url)
	} else if _, err := exec.LookPath("wget"); err == nil {
		cmd = exec.Command("wget", "-O", "/tmp/toolchain.tar.gz", url)
	} else {
		return "", fmt.Errorf("neither curl nor wget is available")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("download failed: %w", err)
	}

	// Extract
	extractCmd := exec.Command("tar", "-xzf", "/tmp/toolchain.tar.gz", "-C", destPath, "--strip-components=1")
	_, err = extractCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("extraction failed: %w", err)
	}

	// Clean up
	os.Remove("/tmp/toolchain.tar.gz")

	return fmt.Sprintf("Downloaded and extracted to %s", destPath), nil
}

// String returns a string representation of a toolchain
func (tc *Toolchain) String() string {
	status := "unavailable"
	if tc.IsAvailable {
		status = fmt.Sprintf("available (%s)", tc.Version)
	}
	return fmt.Sprintf("%s [%s]", tc.Name, status)
}

// IsRequired returns true if a toolchain is required for the given project
func IsRequired(projectPath string, tcType ToolchainType) bool {
	projectIndicators := map[ToolchainType][]string{
		ToolchainGo:     {"go.mod", "go.sum"},
		ToolchainNode:   {"package.json", "package-lock.json"},
		ToolchainPython: {"pyproject.toml", "setup.py", "requirements.txt"},
		ToolchainJava:   {"pom.xml", "build.gradle"},
		ToolchainDotNet: {".sln", ".csproj"},
		ToolchainRust:   {"Cargo.toml"},
	}

	files, ok := projectIndicators[tcType]
	if !ok {
		return false
	}

	for _, file := range files {
		path := filepath.Join(projectPath, file)
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}
