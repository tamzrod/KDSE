# KDSE Someday/Maybe Knowledge Management

## Overview

Someday/Maybe is a structured engineering knowledge repository for ideas that are intentionally deferred. It is NOT a task backlog or TODO list.

The purpose is to preserve valuable ideas without allowing them to distract the current engineering objective.

## Architectural Principle

```
Current Objective
       ↓
Future Idea
       ↓
Someday/Maybe
       ↓
Revisit when appropriate
```

Ideas remain available but do not interfere with active engineering work.

## Directory Structure

```
.kdse/
    someday/
        someday.yaml       # Inventory manifest
        ideas/             # Active someday ideas
        archived/          # Archived ideas
        promoted/         # Promoted ideas (entering workflow)
```

## Idea Schema

Every Someday item contains:

| Field | Type | Description |
|-------|------|-------------|
| ID | string | Unique identifier (IDEA-001, IDEA-002, ...) |
| Title | string | Brief title |
| Description | string | Detailed description |
| Problem | string | Problem this idea solves |
| PotentialValue | string | Expected value of implementation |
| Origin | string | Where the idea originated |
| DateCreated | string | ISO timestamp |
| Author | string | Creator |
| RelatedObjective | string | Linked engineering objective |
| RelatedKnowledge | array | References to knowledge artifacts |
| Dependencies | array | Required prerequisites |
| EstimatedComplexity | string | LOW/MEDIUM/HIGH/UNKNOWN |
| Confidence | float | Initial confidence (0.0-1.0) |
| Priority | int | 1-5 (1 is highest) |
| Status | enum | SOMEDAY/LABORATORY/PROMOTED/IMPLEMENTED/ARCHIVED/REJECTED |
| Examples | array | Usage examples |
| ReasonDeferred | string | Why deferred |
| ReviewDate | string | When to revisit (optional) |
| Tags | array | Categorization tags |
| PromotionHistory | array | Status transition history |
| TraceabilityLinks | array | Links to derived work |

## Status Lifecycle

```
SOMEDAY
   ↓
LABORATORY (optional experimentation)
   ↓
PROMOTED (entering active workflow)
   ↓
IMPLEMENTED
```

Other states:
- **ARCHIVED** - No longer relevant
- **REJECTED** - Deliberately rejected

## Commands

### `kdse someday add`

Add a new someday/maybe idea.

```bash
kdse someday add \
  --title "GUI Runtime" \
  --description "Create a graphical interface for KDSE" \
  --problem "Users need visual interaction" \
  --origin "User feedback" \
  --author "Team Lead"
```

### `kdse someday list`

List all ideas, optionally filtered by status.

```bash
# List all ideas
kdse someday list

# Filter by status
kdse someday list --status SOMEDAY
kdse someday list --status PROMOTED
```

### `kdse someday show`

Show details of a specific idea.

```bash
kdse someday show IDEA-001
```

### `kdse someday promote`

Promote an idea to active consideration.

```bash
kdse someday promote IDEA-001 --reason "Ready for implementation"
```

Promotion creates traceability. The original idea remains linked.

### `kdse someday archive`

Archive an idea (no longer relevant).

```bash
kdse someday archive IDEA-002 --reason "No longer applicable"
```

### `kdse someday search`

Search ideas by keyword.

```bash
kdse someday search "GUI"
kdse someday search "runtime"
```

### `kdse someday review`

Show high-priority ideas ready for review.

```bash
kdse someday review
```

### `kdse someday export`

Export all ideas.

```bash
kdse someday export
```

## Detection Patterns

During conversations, the runtime detects statements like:

- "Maybe..."
- "Someday..."
- "Future..."
- "Interesting..."
- "Not now..."
- "Later..."
- "Perhaps..."
- "Potentially..."
- "If we had time..."
- "Nice to have..."
- "Down the road..."

When detected, the AI asks:

> "This sounds like a Someday/Maybe idea. Would you like me to capture it?"

Only capture with user approval.

## Promotion Workflow

An idea can be promoted into active engineering:

```
Someday
    ↓
Knowledge Collection
    ↓
Assessment
    ↓
Foundation
    ↓
Architecture
    ↓
Implementation
    ↓
Verification
    ↓
Complete
```

Promotion creates traceability. The original idea remains linked.

## Knowledge Integration

Someday items participate in Knowledge Discovery. When a future objective matches an existing idea:

```
Objective: Create KDSE GUI
    ↓
Knowledge Discovery finds: IDEA-004 (Knowledge Graph Runtime)
    ↓
Runtime recommends reusing it
```

## Laboratory Integration

Ideas are hypotheses. Before entering the Global Runtime they may pass through the Laboratory:

```
Observation
    ↓
Someday
    ↓
Experiment
    ↓
Validated
    ↓
Global Runtime
```

## Runtime Verification

`kdse runtime verify` must verify:

- `.kdse/someday/` exists
- `.kdse/someday/ideas/` exists
- `.kdse/someday/archived/` exists
- `.kdse/someday/promoted/` exists
- `.kdse/someday/someday.yaml` exists
- Manifest contains someday inventory
- Promotion history is valid
- Traceability remains intact

## Design Goals

1. **Never lose engineering ideas** - All ideas are preserved indefinitely
2. **Never interrupt current work** - Ideas do not interfere with active objectives
3. **Searchable** - Ideas can be found when relevant
4. **Traceable** - Promotions maintain links to original ideas
5. **Promotable** - Ideas can enter the active workflow when appropriate

## Implementation

- **Package**: `internal/someday`
- **Template**: `runtime/template/someday/`
- **Commands**: Added to `cmd/kdse/main.go`
- **Initialization**: Updated `internal/runtime/runtime.go`

## Future Enhancements

- [ ] Automatic review reminders based on ReviewDate
- [ ] Integration with project planning tools
- [ ] Confidence decay over time
- [ ] Linked dependency tracking
- [ ] Impact estimation calculations
