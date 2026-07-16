# KDSE End-to-End Example

This directory demonstrates a complete KDSE workflow from initialization to knowledge promotion.

## Scenario

We're building a user authentication system. This example shows how to:

1. Initialize the KDSE workspace
2. Create the project agreement
3. Add knowledge to the notebook
4. Promote knowledge through the pipeline
5. View workspace status

## Step 1: Initialize Workspace

```bash
# Navigate to your project directory
cd my-auth-service

# Initialize the KDSE workspace
kdse init
```

Output:
```
╔═══════════════════════════════════════════════════════════════╗
║              KDSE Workspace Initialization                   ║
╠═══════════════════════════════════════════════════════════════╣
║ Repository: /path/to/my-auth-service
╠═══════════════════════════════════════════════════════════════╣
║ ✓ Created .kdse/ directory
║ ✓ Workspace ready for KDSE operations
╚═══════════════════════════════════════════════════════════════╝

Next steps:
  kdse agreement init      # Initialize project agreement
  kdse notebook add <title> # Add first knowledge entry
  kdse status              # View workspace status
```

## Step 2: Initialize Project Agreement

```bash
kdse agreement init
```

Output:
```
╔═══════════════════════════════════════════════════════════════╗
║              AGREEMENT INITIALIZED                            ║
╠═══════════════════════════════════════════════════════════════╣
║ Project:     my-auth-service
║ Phase:       Problem
║ Methodology: v1.0
║ Runtime:     v1.0.0
╚═══════════════════════════════════════════════════════════════╝
```

## Step 3: Add Knowledge to Notebook

### Entry 1: User Requirements

```bash
kdse notebook add "Users need secure authentication" \
  --source "customer-interviews.md" \
  --tag "requirements"
```

Output:
```
Notebook entry created: KDSE-KNOW-20240115-A
```

### Entry 2: Performance Requirements

```bash
kdse notebook add "Authentication must complete within 200ms" \
  --source "benchmark-results.json" \
  --tag "performance"
```

Output:
```
Notebook entry created: KDSE-KNOW-20240115-B
```

### Entry 3: Security Constraints

```bash
kdse notebook add "Passwords must use bcrypt with cost factor 12" \
  --source "security-audit-2024.md" \
  --tag "security"
```

Output:
```
Notebook entry created: KDSE-KNOW-20240115-C
```

## Step 4: View Notebook

```bash
kdse notebook list
```

Output:
```
╔═══════════════════════════════════════════════════════════════╗
║              Engineering Notebook                              ║
╠═══════════════════════════════════════════════════════════════╣
║ Total Entries: 3
║ Notebook: 3  |  Candidate: 0  |  Promoted: 0  |  Rejected: 0
╠═══════════════════════════════════════════════════════════════╣
║ ○ KDSE-KNOW-20240115-A Users need secure authentication
║ ○ KDSE-KNOW-20240115-B Authentication must complete within
║ ○ KDSE-KNOW-20240115-C Passwords must use bcrypt with co
╚═══════════════════════════════════════════════════════════════╝
```

## Step 5: Submit Candidates for Review

```bash
# Submit the first entry as a candidate
kdse promote submit KDSE-KNOW-20240115-A

# Submit the third entry (security is strong evidence)
kdse promote submit KDSE-KNOW-20240115-C
```

Output:
```
Entry KDSE-KNOW-20240115-A promoted to candidate
Entry KDSE-KNOW-20240115-B promoted to candidate
```

## Step 6: Review Candidates

```bash
# Review and accept the security entry (strong evidence)
kdse promote review KDSE-KNOW-20240115-C \
  --accept \
  --strength 5 \
  --rationale "Confirmed by security audit and OWASP guidelines"
```

Output:
```
Entry KDSE-KNOW-20240115-C accepted as knowledge (strength: ★★★★★)
```

```bash
# Review the user requirements entry
kdse promote review KDSE-KNOW-20240115-A \
  --accept \
  --strength 4 \
  --rationale "Based on customer interviews, confirms existing assumptions"
```

Output:
```
Entry KDSE-KNOW-20240115-A accepted as knowledge (strength: ★★★★☆)
```

## Step 7: View Final Status

```bash
kdse status
```

Output:
```
╔═══════════════════════════════════════════════════════════════╗
║                    KDSE Workspace Status                      ║
╠═══════════════════════════════════════════════════════════════╣
║ Repository: /path/to/my-auth-service
║ Workspace:  /path/to/my-auth-service/.kdse
╠═══════════════════════════════════════════════════════════════╣
║ Agreement                                                    ║
║ Project:    my-auth-service
║ Phase:      Problem
║ Assumptions: 0
╠═══════════════════════════════════════════════════════════════╣
║ Knowledge                                                    ║
║ Total:      3
║ Notebook:   1  |  Candidate: 0  |  Promoted: 2
╚═══════════════════════════════════════════════════════════════╝
```

## Knowledge Artifacts Created

| ID | Title | Status | Evidence Strength |
|----|-------|--------|-------------------|
| KDSE-KNOW-20240115-A | Users need secure authentication | Promoted | ★★★★☆ |
| KDSE-KNOW-20240115-B | Authentication must complete within 200ms | Notebook | - |
| KDSE-KNOW-20240115-C | Passwords must use bcrypt with cost factor 12 | Promoted | ★★★★★ |

## Next Steps

1. **Derive Architecture**: Use the promoted knowledge to create architectural decisions
2. **Update Phase**: Move from "Problem" to "Architecture" phase:
   ```bash
   kdse agreement phase Architecture
   ```
3. **Add More Evidence**: Continue adding evidence for remaining entries
4. **Review Performance Entry**: Submit and review KDSE-KNOW-20240115-B

## Files Created

The `.kdse/` workspace now contains:

```
.kdse/
├── agreement.json          # Project agreement
└── knowledge/
    └── entries.json        # Knowledge entries with evidence strength
```

## Derivation to Architecture

With knowledge promoted, you can now trace architectural decisions:

```
Knowledge: "Passwords must use bcrypt with cost factor 12"
    │
    │ Authorizes
    ▼
ADR: "Use bcrypt-12 for password hashing"
    │
    │ Guides
    ▼
Implementation: User model with bcrypt hash field
```

This traceability ensures every architectural decision traces back to authoritative knowledge.
