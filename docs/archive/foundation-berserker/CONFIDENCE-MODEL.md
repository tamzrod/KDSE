# Confidence Model

## Definition

Confidence represents **empirical trust** in an experience's reusability.

**Important:** Confidence does NOT represent:
- Certainty of root cause
- Quality of solution
- Severity of problem

Confidence represents ONLY:
- How often has this worked?
- How consistent is the outcome?
- How validated is the approach?

## Confidence Levels

| Level | Definition | Criteria | When to Use |
|-------|------------|----------|-------------|
| **Experimental** | Verified once | Observed once, solution worked once | New discoveries |
| **Low** | Verified multiple times in similar situations | Worked 2-3 times in similar contexts | Initial learning |
| **Medium** | Repeated successfully across multiple sessions | Worked in 3+ different sessions | Confirmed approach |
| **High** | Repeated successfully over time with no contradictions | Worked consistently, no failures reported | Standard approach |
| **Proven** | Standard practice across repositories | Accepted as norm, alternatives rare | Established practice |

## Confidence Progression

```
Experimental → Low → Medium → High → Proven
    ↑
    │ Each step requires evidence of successful reuse
    │
[Capture at Experimental]
[Elevate when repeated]
```

## Confidence Examples

### Experimental

```yaml
confidence: Experimental
# First time discovering that GOPROXY helps in restricted networks
# Solution worked once, needs validation
```

### Low

```yaml
confidence: Low
# Worked on 2-3 sandbox builds
# Similar conditions but not identical
```

### Medium

```yaml
confidence: Medium
# Worked across 5 different Go projects
# Consistent in similar environments
```

### High

```yaml
confidence: High
# Used in 20+ sessions over 6 months
# Never failed when conditions matched
```

### Proven

```yaml
confidence: Proven
# Team-wide standard practice
# New team members learn this immediately
# Alternative approaches considered inferior
```

## Confidence Guidelines

**When to elevate confidence:**
- Same solution worked 3+ times in similar situations
- No contradictory evidence after multiple applications
- Solution becomes standard team practice

**When to downgrade confidence:**
- Solution failed when conditions appeared same
- New evidence suggests alternative is better
- Environment changed (new versions, new tools)

**When to retire experience:**
- Solution no longer applies (deprecated tool/approach)
- Superseded by better solution
- Confidence dropped to unreliable

## Quality Signals

| Signal | Meaning |
|--------|---------|
| Proven confidence | Safe to use, no alternatives needed |
| High confidence | Safe to use, well-validated |
| Medium confidence | Reasonable to try, may need adjustment |
| Low confidence | Experimental, try with caution |
| Experimental | Hint only, validate independently |
