// Package toolchain provides automatic toolchain management for KDSE.
//
// This file contains the verification integration for runtime initialization.
package toolchain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// VerificationLevel defines how strictly verification should be performed
type VerificationLevel string

const (
	// VerificationLevelRequired verifies required toolchains only
	VerificationLevelRequired VerificationLevel = "required"
	// VerificationLevelAll verifies all detected toolchains
	VerificationLevelAll VerificationLevel = "all"
	// VerificationLevelNone skips verification (not recommended)
	VerificationLevelNone VerificationLevel = "none"
)

// RuntimeVerification contains the result of runtime toolchain verification
type RuntimeVerification struct {
	ProjectPath      string
	RequiredTools    []ToolchainType
	VerifiedTools    []*VerificationResult
	MissingTools     []ToolchainType
	InstallBase      string
	EnvironmentPaths map[string]string
	Success          bool
	Errors           []string
}

// VerifyRuntime performs toolchain verification as part of runtime initialization
func VerifyRuntime(projectPath string, level VerificationLevel) *RuntimeVerification {
	result := &RuntimeVerification{
		ProjectPath:      projectPath,
		VerifiedTools:     []*VerificationResult{},
		MissingTools:     []ToolchainType{},
		EnvironmentPaths: make(map[string]string),
		Success:          true,
	}

	manager := NewManager(projectPath)

	// Determine which toolchains to verify
	var toolsToVerify []ToolchainType

	switch level {
	case VerificationLevelRequired:
		result.RequiredTools = manager.GetRequiredToolchains()
		toolsToVerify = result.RequiredTools
	case VerificationLevelAll:
		toolsToVerify = []ToolchainType{
			ToolchainGo, ToolchainNode, ToolchainPython,
			ToolchainJava, ToolchainDotNet, ToolchainRust,
		}
		result.RequiredTools = toolsToVerify
	case VerificationLevelNone:
		return result
	}

	// Verify each toolchain
	for _, tcType := range toolsToVerify {
		verifyResult := manager.Ensure(tcType)
		result.VerifiedTools = append(result.VerifiedTools, verifyResult)

		if !verifyResult.Success {
			result.Success = false
			result.MissingTools = append(result.MissingTools, tcType)
			result.Errors = append(result.Errors, verifyResult.Errors...)
		}
	}

	// Get environment updates
	result.InstallBase = manager.GetInstallBase()
	result.EnvironmentPaths = manager.UpdateEnvironment()

	return result
}

// FormatVerificationResult formats the verification result for display
func FormatVerificationResult(v *RuntimeVerification) string {
	var output strings.Builder

	output.WriteString("╔═══════════════════════════════════════════════════════════════╗\n")
	output.WriteString("║         KDSE Runtime Toolchain Verification                 ║\n")
	output.WriteString("╠═══════════════════════════════════════════════════════════════╣\n")
	output.WriteString(fmt.Sprintf("║ Project: %s\n", v.ProjectPath))
	output.WriteString(fmt.Sprintf("║ Install Base: %s\n", v.InstallBase))
	output.WriteString("╠═══════════════════════════════════════════════════════════════╣\n")

	if len(v.RequiredTools) == 0 {
		output.WriteString("║ No toolchains required for this project                   ║\n")
	} else {
		output.WriteString("║ Required Toolchains:                                       ║\n")
		for _, tc := range v.VerifiedTools {
			status := "✓ AVAILABLE"
			if !tc.Success {
				status = "✗ MISSING"
			}
			output.WriteString(fmt.Sprintf("║   %s: %s\n", tc.Toolchain.Name, status))
		}
	}

	output.WriteString("╠═══════════════════════════════════════════════════════════════╣\n")

	if v.Success {
		output.WriteString("║ Status: ALL TOOLCHAINS VERIFIED                            ║\n")
	} else {
		output.WriteString("║ Status: VERIFICATION FAILED                                ║\n")
		output.WriteString("╠═══════════════════════════════════════════════════════════════╣\n")
		output.WriteString("║ Missing Toolchains:                                       ║\n")
		for _, tcType := range v.MissingTools {
			output.WriteString(fmt.Sprintf("║   - %s\n", tcType))
		}
	}

	if len(v.EnvironmentPaths) > 0 {
		output.WriteString("╠═══════════════════════════════════════════════════════════════╣\n")
		output.WriteString("║ Environment Updates:                                      ║\n")
		output.WriteString("║ To use installed toolchains, add to your PATH:            ║\n")
		output.WriteString(fmt.Sprintf("║   export PATH=\"%s:$PATH\"\n", v.InstallBase))
	}

	output.WriteString("╚═══════════════════════════════════════════════════════════════╝\n")

	return output.String()
}

// ReportVerificationFailure creates a detailed failure report
func ReportVerificationFailure(v *RuntimeVerification) string {
	var output strings.Builder

	output.WriteString("\n")
	output.WriteString("═══════════════════════════════════════════════════════════════\n")
	output.WriteString("  KDSE TOOLCHAIN VERIFICATION FAILED\n")
	output.WriteString("═══════════════════════════════════════════════════════════════\n\n")

	output.WriteString("KDSE cannot proceed because the following toolchains are unavailable:\n\n")

	for _, tcType := range v.MissingTools {
		tc := &Toolchain{Type: tcType}
		tc.Name = string(tcType)

		output.WriteString(fmt.Sprintf("  %s:\n", tc.Name))
		output.WriteString(getInstallInstructions(tcType))
		output.WriteString("\n")
	}

	output.WriteString("───────────────────────────────────────────────────────────────\n")
	output.WriteString("  VERIFICATION POLICY\n")
	output.WriteString("───────────────────────────────────────────────────────────────\n\n")
	output.WriteString("  KDSE never silently skips verification due to missing tooling.\n")
	output.WriteString("  All required toolchains must be available before implementation.\n\n")

	output.WriteString("  To proceed:\n")
	output.WriteString("  1. Install the missing toolchains above\n")
	output.WriteString("  2. Re-run KDSE initialization\n")
	output.WriteString("  3. KDSE will verify toolchains are working\n\n")

	output.WriteString("  Alternative: Install to writable location\n")
	output.WriteString(fmt.Sprintf("    export KDSE_TOOLS_PATH=/workspace/.tools\n"))
	output.WriteString(fmt.Sprintf("    kdse initialize\n\n"))

	output.WriteString("═══════════════════════════════════════════════════════════════\n")

	return output.String()
}

// getInstallInstructions returns installation instructions for a toolchain
func getInstallInstructions(tcType ToolchainType) string {
	switch tcType {
	case ToolchainGo:
		return `    1. Download from: https://go.dev/dl/
       2. Extract to:   /workspace/.tools/go
       3. Add to PATH:  export PATH=/workspace/.tools/go/bin:$PATH`
	case ToolchainNode:
		return `    1. Download from: https://nodejs.org/
       2. Extract to:   /workspace/.tools/node
       3. Add to PATH:  export PATH=/workspace/.tools/node/bin:$PATH`
	case ToolchainPython:
		return `    1. Install via package manager:
       2. Linux:   apt install python3 python3-pip
       3. macOS:   brew install python3
       4. Windows: Download from python.org`
	case ToolchainJava:
		return `    1. Download from: https://adoptium.net/ or https://aws.amazon.com/corretto/
       2. Extract to:   /workspace/.tools/java
       3. Set JAVA_HOME: export JAVA_HOME=/workspace/.tools/java`
	case ToolchainDotNet:
		return `    1. Download from: https://dotnet.microsoft.com/download
       2. Install to:    /workspace/.tools/dotnet
       3. Add to PATH:  export PATH=/workspace/.tools/dotnet:$PATH`
	case ToolchainRust:
		return `    1. Install via: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
       2. Toolchain will be installed to ~/.cargo/bin`
	default:
		return `    Installation instructions not available for this toolchain.`
	}
}

// WriteToolchainLog writes the toolchain verification log to a file
func WriteToolchainLog(projectPath string, v *RuntimeVerification) error {
	logPath := filepath.Join(projectPath, ".kdse", "runtime", "toolchain-log.txt")
	
	var content strings.Builder
	content.WriteString(fmt.Sprintf("KDSE Toolchain Verification Log\n"))
	content.WriteString(fmt.Sprintf("Generated: %s\n", "now"))
	content.WriteString(fmt.Sprintf("Project: %s\n", v.ProjectPath))
	content.WriteString(fmt.Sprintf("Install Base: %s\n", v.InstallBase))
	content.WriteString(fmt.Sprintf("Success: %v\n", v.Success))
	content.WriteString("\nRequired Toolchains:\n")
	for _, tcType := range v.RequiredTools {
		content.WriteString(fmt.Sprintf("  - %s\n", tcType))
	}
	content.WriteString("\nVerification Results:\n")
	for _, result := range v.VerifiedTools {
		content.WriteString(fmt.Sprintf("  %s: %v\n", result.Toolchain.Name, result.Success))
		if !result.Success {
			content.WriteString(fmt.Sprintf("    Errors: %v\n", result.Errors))
		}
	}
	if len(v.MissingTools) > 0 {
		content.WriteString("\nMissing Toolchains:\n")
		for _, tcType := range v.MissingTools {
			content.WriteString(fmt.Sprintf("  - %s\n", tcType))
		}
	}

	return os.WriteFile(logPath, []byte(content.String()), 0644)
}
