# CLI Runtime

**Document Version:** 1.0  
**Type:** Normative  
**Effective Date:** 2026-07-17

---

## Purpose

This document defines the CLI Runtime for KDSE. The CLI Runtime is a **thin adapter** that translates command-line operations into Workspace Engine calls.

---

## Core Principle: Thin Adapter

The CLI Runtime MUST be a thin adapter. It contains NO business logic.

### Responsibilities

| Responsibility | Description | Required |
|----------------|-------------|----------|
| Parse commands | Parse CLI arguments | YES |
| Call Workspace Engine | Invoke engine methods | YES |
| Display output | Format and display results | YES |
| Handle errors | Convert errors to user messages | YES |

### Non-Responsibilities

The CLI Runtime does NOT:
- Implement engineering rules
- Validate artifacts
- Manage state
- Enforce phases
- Make decisions

---

## Architecture

### Position in Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         User (Human/AI)                              │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ CLI Commands
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         CLI Runtime                                  │
│                       (Thin Adapter)                                 │
│                                                                     │
│  • Parses: kdse init, kdse verify, kdse phase, etc.                   │
│  • Calls: Workspace Engine methods                                   │
│  • Formats: Output for terminal                                       │
│  • NO business logic                                                 │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Engine Interface
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Workspace Engine                               │
│                        (State Owner)                                  │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        KDSE Methodology                               │
│                        (Engineering Rules)                            │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Command Interface

### Command Structure

```go
// Command represents a CLI command
type Command struct {
    Name        string
    Description string
    Args        []Arg
    Flags       []Flag
    Handler     CommandHandler
}

// CommandHandler is the function signature for command handlers
type CommandHandler func(ctx context.Context, args []string, flags map[string]string) error
```

### Command Registry

```go
var commands = []Command{
    {
        Name:        "init",
        Description: "Initialize KDSE runtime in current directory",
        Args:        []Arg{},
        Flags: []Flag{
            {Name: "template", Default: "default", Description: "Template to use"},
            {Name: "type", Default: "cli", Description: "Runtime type"},
        },
        Handler: handleInit,
    },
    {
        Name:        "verify",
        Description: "Verify KDSE runtime state",
        Args:        []Arg{},
        Flags:       []Flag{},
        Handler:     handleVerify,
    },
    {
        Name:        "phase",
        Description: "Show or advance current phase",
        Args: []Arg{
            {Name: "action", Options: []string{"show", "advance"}},
        },
        Flags:       []Flag{},
        Handler:     handlePhase,
    },
    {
        Name:        "report",
        Description: "Generate a report",
        Args: []Arg{
            {Name: "type", Options: []string{"phase", "verification", "summary", "progress"}},
        },
        Flags: []Flag{
            {Name: "output", Default: "", Description: "Output file path"},
        },
        Handler: handleReport,
    },
    {
        Name:        "artifacts",
        Description: "List artifacts for a phase",
        Args: []Arg{
            {Name: "phase", Options: []string{"knowledge", "architecture", "implementation", "verification", "reports"}},
        },
        Flags: []Flag{},
        Handler: handleArtifacts,
    },
}
```

---

## Command Handlers

### Init Command

```go
func handleInit(ctx context.Context, args []string, flags map[string]string) error {
    // 1. Parse arguments
    opts := workspace.InitOptions{
        Path:     getWorkingDir(),
        Type:     parseRuntimeType(flags["type"]),
        Template: flags["template"],
    }
    
    // 2. Call Workspace Engine (NO local logic)
    engine := workspace.NewEngine()
    ws, err := engine.InitializeWorkspace(ctx, opts)
    if err != nil {
        return formatError(err)
    }
    
    // 3. Display result
    fmt.Printf("KDSE initialized successfully at %s\n", ws.Path)
    fmt.Printf("Runtime type: %s\n", ws.Runtime.Type)
    fmt.Printf("Current phase: %s\n", ws.State.CurrentPhase)
    
    return nil
}
```

### Verify Command

```go
func handleVerify(ctx context.Context, args []string, flags map[string]string) error {
    // 1. Call Workspace Engine
    engine := workspace.NewEngine()
    result, err := engine.VerifyWorkspace(ctx)
    if err != nil {
        return formatError(err)
    }
    
    // 2. Display result
    if result.Valid {
        fmt.Printf("✓ Workspace verified\n")
        fmt.Printf("  Phase: %s\n", result.Phase)
        fmt.Printf("  Runtime: %s v%s\n", result.RuntimeInfo.Type, result.RuntimeInfo.Version)
    } else {
        fmt.Printf("✗ Verification failed\n")
        for _, err := range result.Errors {
            fmt.Printf("  Error: %s\n", err.Message)
        }
    }
    
    return nil
}
```

### Phase Command

```go
func handlePhase(ctx context.Context, args []string, flags map[string]string) error {
    engine := workspace.NewEngine()
    action := args[0]
    
    switch action {
    case "show":
        // Get current phase
        phase, err := engine.GetPhase(ctx)
        if err != nil {
            return formatError(err)
        }
        fmt.Printf("Current phase: %s\n", phase.Name)
        fmt.Printf("Description: %s\n", phase.Description)
        
    case "advance":
        // Advance to next phase
        nextPhase := phase.Name // Would come from engine
        transition, err := engine.AdvancePhase(ctx, nextPhase)
        if err != nil {
            return formatError(err)
        }
        fmt.Printf("Phase advanced: %s → %s\n", transition.From, transition.To)
    }
    
    return nil
}
```

### Report Command

```go
func handleReport(ctx context.Context, args []string, flags map[string]string) error {
    // Parse report type
    reportType := parseReportType(args[0])
    
    // Call Workspace Engine
    engine := workspace.NewEngine()
    opts := workspace.ReportOptions{
        Type:   reportType,
        Output: flags["output"],
    }
    
    report, err := engine.GenerateReport(ctx, opts)
    if err != nil {
        return formatError(err)
    }
    
    // Display report
    if opts.Output != "" {
        if err := writeReport(report, opts.Output); err != nil {
            return formatError(err)
        }
        fmt.Printf("Report written to %s\n", opts.Output)
    } else {
        printReport(report)
    }
    
    return nil
}
```

---

## Output Formatting

### Error Formatting

```go
func formatError(err error) error {
    var engineErr *workspace.EngineError
    if errors.As(err, &engineErr) {
        return fmt.Errorf("kdse: %s\n\n%s\n\nHint: %s",
            engineErr.Message,
            formatDetails(engineErr.Details),
            engineErr.Remediation)
    }
    return fmt.Errorf("kdse: %v", err)
}

func formatDetails(details map[string]interface{}) string {
    if len(details) == 0 {
        return ""
    }
    var lines []string
    for k, v := range details {
        lines = append(lines, fmt.Sprintf("  %s: %v", k, v))
    }
    return strings.Join(lines, "\n")
}
```

### Success Formatting

```go
func formatSuccess(msg string) {
    fmt.Printf("✓ %s\n", msg)
}

func formatProgress(current, total int, label string) {
    pct := float64(current) / float64(total) * 100
    fmt.Printf("[%d/%d] %s (%.0f%%)\n", current, total, label, pct)
}
```

---

## Main Entry Point

### Main Function

```go
// cmd/kdse/main.go
package main

import (
    "context"
    "os"
    
    "kdse/cmd/kdse/commands"
    "kdse/internal/workspace"
)

func main() {
    // Create context
    ctx := context.Background()
    
    // Create engine (shared across commands)
    engine := workspace.NewEngine()
    
    // Parse and execute command
    cmd := commands.Parse(os.Args[1:])
    if err := cmd.Handler(ctx, cmd.Args, cmd.Flags); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

### Package Structure

```
cmd/
└── kdse/
    ├── main.go              # Entry point
    ├── commands/
    │   ├── commands.go      # Command registry
    │   ├── init.go          # init command
    │   ├── verify.go        # verify command
    │   ├── phase.go         # phase command
    │   ├── report.go        # report command
    │   └── artifacts.go    # artifacts command
    ├── formatter/
    │   ├── formatter.go     # Output formatting
    │   ├── error.go        # Error formatting
    │   └── table.go        # Table formatting
    └── parser/
        ├── parser.go       # Argument parsing
        └── flags.go        # Flag parsing
```

---

## Dependency Rules

### Allowed Dependencies

```go
// cmd/kdse/main.go
package main

import (
    // OK: Standard library
    "context"
    "fmt"
    "os"
    
    // OK: Internal workspace engine
    "kdse/internal/workspace"
    
    // OK: Internal types
    "kdse/internal/types"
    
    // OK: Internal methodology (read-only interfaces)
    "kdse/internal/methodology/lifecycle"
)
```

### Forbidden Dependencies

```go
// cmd/kdse/main.go - DO NOT DO THIS
import (
    // FORBIDDEN: Internal runtime (except as engine)
    "kdse/internal/runtime"  // NO
    
    // FORBIDDEN: Direct state management
    "kdse/internal/state"    // NO
    
    // FORBIDDEN: Direct bootstrap logic
    "kdse/internal/bootstrap" // NO (use workspace engine)
    
    // FORBIDDEN: Direct knowledge logic
    "kdse/internal/knowledge" // NO (use workspace engine)
)
```

---

## Verification Gate

### Pre-Command Verification

```go
func withVerification(ctx context.Context, cmd Command, handler CommandHandler) CommandHandler {
    return func(ctx context.Context, args []string, flags map[string]string) error {
        // Skip verification for init command
        if cmd.Name == "init" {
            return handler(ctx, args, flags)
        }
        
        // Verify workspace before command
        engine := workspace.NewEngine()
        result, err := engine.VerifyWorkspace(ctx)
        if err != nil {
            return &VerificationError{
                Command: cmd.Name,
                Error:   err,
            }
        }
        
        if !result.Valid {
            return &VerificationError{
                Command: cmd.Name,
                Error:   ErrWorkspaceInvalid,
                Details: result.Errors,
            }
        }
        
        // Add workspace info to context
        ctx = context.WithValue(ctx, WorkspaceKey, result)
        
        return handler(ctx, args, flags)
    }
}
```

---

## Testing

### CLI Test Strategy

```go
// cmd/kdse/commands/commands_test.go
package commands

func TestInitCommand(t *testing.T) {
    // Create temp directory
    dir := t.TempDir()
    
    // Execute init command
    cmd := Command{Name: "init"}
    handler := withVerification(context.Background(), cmd, handleInit)
    
    err := handler(context.Background(), []string{dir}, map[string]string{
        "type": "cli",
    })
    
    // Verify results
    assert.NoError(t, err)
    
    // Verify .kdse exists
    assert.DirExists(t, filepath.Join(dir, ".kdse"))
    
    // Verify required files
    assert.FileExists(t, filepath.Join(dir, ".kdse", "runtime.yaml"))
    assert.FileExists(t, filepath.Join(dir, ".kdse", "phase.yaml"))
}
```

---

## Error Codes

### CLI Error Codes

| Code | Description | Remediation |
|------|-------------|-------------|
| E001 | Workspace not found | Run `kdse init` |
| E002 | Verification failed | Check .kdse directory |
| E003 | Invalid phase | Use `kdse phase show` |
| E004 | Invalid transition | Complete current phase first |
| E005 | Report generation failed | Check workspace state |
| E006 | Unknown command | Use `kdse help` |

---

## Document Relationships

```
CLI_RUNTIME.md
    │
    ├── Defines: CLI adapter responsibilities, commands, formatting
    │
    ├── Referenced By:
    │   ├── All CLI implementations
    │   └── Documentation
    │
    ├── References:
    │   ├── RUNTIME_ARCHITECTURE.md
    │   ├── WORKSPACE_ENGINE.md
    │   └── PRINCIPLES.md (I-001: Thin Adapters)
    │
    └── Related Documents:
        └── MCP_RUNTIME.md
```

---

*This document is normative. The CLI Runtime MUST be a thin adapter with no business logic.*
