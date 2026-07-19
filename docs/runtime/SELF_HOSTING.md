# Self-Hosting: KDSE Applied to KDSE

## Overview

KDSE now supports self-hosting - the ability to analyze and evolve its own architecture using the same evidence-driven methodology it applies to software projects. This capability ensures that KDSE never modifies itself without first understanding itself.

## Core Principle

> **KDSE never modifies itself without first understanding itself.**

This is enforced through a structured workflow that requires evidence collection, architecture modeling, impact analysis, and explicit approval before any self-modification.

## Architecture

The self-hosting module (`internal/selfhost/`) consists of several components:

### Module Structure

```
internal/selfhost/
├── selfhost.go      # Main entry point and Manager
├── model.go         # Architecture model types
├── analyzer.go      # Runtime self-analysis
├── impact.go        # Impact analysis engine
├── workflow.go      # Evolution workflow engine
├── promotion.go     # Staged promotion with rollback
└── commands.go      # CLI command handlers
```

### Key Concepts

#### Architecture Model

A machine-readable model of KDSE's architecture including:
- **Components**: Modules, commands, toolchains
- **Dependencies**: Import, invoke, config, and data relationships
- **Data Flows**: How data moves between components
- **Summary Statistics**: Component counts, dependency depth, cycles

#### Self-Assessment Report

A comprehensive report generated through runtime analysis:
- Health status with score (0.0 - 1.0)
- Component inventory
- Dependency analysis
- Critical issues and warnings
- Recommendations

#### Impact Analysis

Evidence-based assessment of proposed changes:
- Affected components identification
- Risk scoring (low/medium/high/critical)
- Breaking change detection
- Dependency impact analysis

## Commands

### `kdse self-assess`

Analyzes KDSE's own architecture and produces a self-assessment report.

```bash
# Basic self-assessment
kdse self-assess

# Verbose output with full report
kdse self-assess --verbose

# JSON output for integration
kdse self-assess --output json

# Save report to specific path
kdse self-assess --save --path /path/to/report.md
```

**Options:**
- `--verbose, -v`: Show full report
- `--output, -o`: Output format (markdown, json)
- `--save, -s`: Save report to file

### `kdse evolve`

Executes the self-evolution workflow with phases:
1. **Collect** - Collect runtime information
2. **Knowledge** - Build knowledge about the runtime
3. **Architecture** - Generate architecture model
4. **Impact Analysis** - Assess change impact
5. **Approval** - Require explicit approval
6. **Implementation** - Apply changes
7. **Verification** - Verify changes

```bash
# Preview evolution workflow
kdse evolve --dry-run

# Approve and proceed
kdse evolve --approve

# Approve and promote to staging
kdse evolve --approve --promote --stage staging

# Force (skip some checks)
kdse evolve --force
```

**Options:**
- `--dry-run, -n`: Preview without making changes
- `--approve, -a`: Grant approval to proceed
- `--force, -f`: Force execution
- `--promote, -p`: Promote after successful evolution
- `--stage`: Target promotion stage (staging, canary, stable)

## Evidence-Driven Workflow

The self-hosting workflow follows KDSE's evidence-driven methodology:

```
┌─────────────────────────────────────────────────────────────┐
│                    Evidence Collection                       │
├─────────────────────────────────────────────────────────────┤
│  • Package scanning (internal/*, cmd/*)                     │
│  • Dependency extraction                                    │
│  • Health check execution                                  │
│  • Architecture model generation                            │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    Knowledge Building                        │
├─────────────────────────────────────────────────────────────┤
│  • Component metadata aggregation                           │
│  • Dependency graph construction                            │
│  • Health assessment summary                               │
│  • Recommendation generation                               │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    Impact Analysis                          │
├─────────────────────────────────────────────────────────────┤
│  • Change classification                                   │
│  • Affected component identification                       │
│  • Risk scoring                                           │
│  • Breaking change detection                               │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    Approval Gate                            │
├─────────────────────────────────────────────────────────────┤
│  • Explicit human approval required                         │
│  • Evidence review                                         │
│  • Risk acknowledgment                                     │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    Implementation                           │
├─────────────────────────────────────────────────────────────┤
│  • Staged rollout support                                  │
│  • Rollback capability                                     │
│  • Promotion management                                    │
└─────────────────────────────────────────────────────────────┘
```

## Promotion Process

Self-modifications go through a staged promotion process:

| Stage | Purpose | Rollback |
|-------|---------|----------|
| **Staging** | Initial testing | Instant |
| **Canary** | Limited rollout | Instant |
| **Stable** | Full deployment | Supported |

### Promotion Commands

```bash
# The evolve command handles promotion automatically with --promote
kdse evolve --approve --promote --stage staging
```

### Rollback

Rollback restores the previous snapshot:

```bash
# Rollback is handled through the promotion system
# When promoting, a snapshot is automatically created
```

## Architecture Model

The architecture model is persisted to:
```
.kkdse/runtime/architecture-model.json
```

### Model Structure

```json
{
  "version": "1.0",
  "generated_at": "2024-01-01T00:00:00Z",
  "components": {
    "runtime": {
      "name": "runtime",
      "type": "runtime",
      "purpose": "Evidence-driven runtime management",
      "path": "internal/runtime",
      "dependencies": ["knowledge", "orchestration"]
    }
  },
  "dependencies": [
    {
      "source": "runtime",
      "target": "knowledge",
      "type": "import",
      "description": "runtime imports knowledge"
    }
  ],
  "summary": {
    "total_components": 20,
    "total_dependencies": 45,
    "max_depth": 4
  }
}
```

## Health Checks

Self-assessment includes comprehensive health checks:

| Check | Weight | Description |
|-------|--------|-------------|
| Invariant Files | 15% | Required runtime files present |
| Directory Structure | 20% | Required directories exist |
| Module Dependencies | 25% | All dependencies resolved |
| Orphan Components | 15% | Components without dependencies |
| Runtime Directory | 15% | .kdse directory accessible |
| Command Handlers | 10% | All commands implemented |

### Health Status Values

- **HEALTHY** (≥0.8): Runtime is in good condition
- **DEGRADED** (0.5-0.8): Some issues detected
- **UNHEALTHY** (<0.5): Critical issues present

## Impact Analysis

Impact analysis classifies changes and assesses risk:

### Change Types
- **add**: New component addition
- **modify**: Existing component modification
- **delete**: Component removal
- **refactor**: Structural changes

### Risk Levels
- **low**: Minimal impact, safe to proceed
- **medium**: Moderate impact, review recommended
- **high**: Significant impact, careful review required
- **critical**: Breaking changes detected, must not proceed without approval

## Output Files

The self-hosting module generates the following outputs:

| File | Location | Description |
|------|----------|-------------|
| Architecture Model | `.kdse/runtime/architecture-model.json` | Machine-readable architecture |
| Self-Assessment Report | `.kdse/reports/self-assessment-*.md` | Human-readable report |
| Evolution State | `.kdse/runtime/evolution-state.json` | Workflow state |
| Promotion History | `.kdse/runtime/promotion/history.json` | Promotion records |
| Snapshots | `.kdse/runtime/promotion/snapshots/*.json` | Point-in-time backups |

## Best Practices

### Before Self-Modification

1. **Run self-assessment first**
   ```bash
   kdse self-assess --verbose
   ```

2. **Review health status**
   - Ensure score ≥ 0.8
   - Address critical issues before proceeding

3. **Check architecture model**
   ```bash
   cat .kdse/runtime/architecture-model.json | jq
   ```

### During Self-Modification

1. **Always use --dry-run first**
   ```bash
   kdse evolve --dry-run
   ```

2. **Review impact analysis**
   - Check affected components
   - Verify risk level is acceptable

3. **Get explicit approval**
   ```bash
   kdse evolve --approve
   ```

### After Self-Modification

1. **Verify changes**
   ```bash
   kdse self-assess
   ```

2. **Check health status**
   - Ensure no regression

3. **Promote incrementally**
   ```bash
   kdse evolve --promote --stage staging
   # Wait, observe, then:
   kdse evolve --promote --stage stable
   ```

## Integration

### Programmatic Access

```go
import "github.com/kdse/runtime/internal/selfhost"

func analyzeMyRuntime() {
    mgr := selfhost.NewManager("/path/to/project")
    
    // Run self-assessment
    report, err := mgr.RunSelfAssessment()
    
    // Analyze impact of changes
    changes := []*selfhost.Change{
        {Type: selfhost.ChangeTypeModify, Component: "runtime"},
    }
    result, err := mgr.AnalyzeImpact(changes)
}
```

### CI/CD Integration

```yaml
# Example: Check health before deployment
- name: Self-assess KDSE
  run: kdse self-assess --output json > health.json
  
- name: Check health score
  run: |
    SCORE=$(cat health.json | jq -r '.health_status.score')
    if (( $(echo "$SCORE < 0.7" | bc -l) )); then
      echo "Health score too low: $SCORE"
      exit 1
    fi
```

## Troubleshooting

### "Architecture model not found"

Run self-assessment first:
```bash
kdse self-assess
```

### "Approval required"

Use the --approve flag:
```bash
kdse evolve --approve
```

### "Cannot promote from current stage"

Ensure valid promotion path:
- Stable → Staging or Canary
- Staging → Canary or Stable
- Canary → Stable

### "Rollback failed"

Check snapshot availability:
```bash
ls .kdse/runtime/promotion/snapshots/
```

## Future Enhancements

Potential improvements for self-hosting:

1. **Automated regression testing** after self-modification
2. **Dependency visualization** with graph output
3. **Change提案提案系统** for proposing and reviewing changes
4. **Scheduled health checks** with alerting
5. **Multi-runtime support** for analyzing external runtimes

## References

- [Architecture Model](./ARCHITECTURE.md)
- [Evidence-Driven Architecture](./EVIDENCE_DRIVEN_ARCHITECTURE.md)
- [Workflow](./WORKFLOW.md)
- [Knowledge Promotion](./KDSE_PROJECT_RUNTIME.md)
