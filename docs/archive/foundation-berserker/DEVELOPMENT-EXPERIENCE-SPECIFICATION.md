# Development Experience Specification

## Definition

A **Development Experience** is a verified record of development knowledge gained through direct observation, investigation, and resolution during a development session.

## Schema

```yaml
experience:
  # REQUIRED FIELDS
  
  situation: string       # Development context that triggers reuse
  solution: string       # Verified procedure or learning
  verify: string        # Confirmation method
  tags: string[]         # Discovery keywords (3-5)
  
  # OPTIONAL FIELDS
  
  confidence: enum       # Experimental|Low|Medium|High|Proven
  evidence: string[]     # Supporting artifacts
  
  # INHERITED (from Git/environment)
  
  id: auto               # Git commit hash
  created: auto          # Commit timestamp
  author: auto           # Git author
```

## Field Specifications

### situation

**Purpose:** Describe when this experience should be reused.

**Rule:** Answer "Am I in the same situation?"

**Structure:**
```
When [context], [unusual condition] occurs
```

**Examples:**
```yaml
situation: "building Go dependencies for the first time in a sandbox with restricted network"
situation: "playwright screenshot captures page before WebSocket updates complete"
situation: "integration test fails when run before database container starts"
```

**Anti-patterns:**
- Error messages alone
- Commands alone
- Vague descriptions

### solution

**Purpose:** Describe the verified procedure or learning.

**Rule:** One sentence describing what consistently works.

**Examples:**
```yaml
solution: "export GOPROXY=https://proxy.golang.org,direct before building"
solution: "add page.wait_for_load_state('networkidle') before screenshot"
solution: "docker-compose up -d && sleep 5 before running tests"
```

### verify

**Purpose:** Describe how to confirm the solution works.

**Rule:** One sentence describing confirmation method.

**Examples:**
```yaml
verify: "go build ./... completes in under 60 seconds"
verify: "screenshot shows full rendered page with all dynamic content"
verify: "pytest tests/integration/ passes with no connection errors"
```

### tags

**Purpose:** Enable discovery through search.

**Rule:** 3-5 keywords relevant to the problem and solution.

**Examples:**
```yaml
tags: [go, build, timeout, proxy, network]
tags: [playwright, screenshot, timing, networkidle]
tags: [docker, postgres, pytest, integration-test]
```

### confidence (optional)

**Purpose:** Indicate empirical trust level.

**Default:** Experimental (if not specified)

**See:** [Confidence Model](CONFIDENCE-MODEL.md)

### evidence (optional)

**Purpose:** Reference supporting artifacts.

**Rule:** Only when words are insufficient.

**Examples:**
```yaml
evidence: ["screenshot-blank.png", "screenshot-success.png"]
evidence: ["commit-abc123", "test-output.log"]
```

## Sixty-Second Rule

Capture must complete in under 60 seconds.

| Field | Maximum |
|-------|---------|
| situation | 15 seconds |
| solution | 15 seconds |
| verify | 10 seconds |
| tags | 5 seconds |
| evidence | 10 seconds (optional) |
| **Total** | **55 seconds** |

If capture takes longer, the experience may not be worth capturing.
