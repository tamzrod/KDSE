# Capture Rules

## The Capture Imperative

Capture when ALL of the following are true:

```
□ Failure was OBSERVED during development
□ Solution was FOUND through investigation  
□ Fix was VERIFIED to work
□ Knowledge would otherwise be LOST
□ Experience is from REAL event, not hypothetical
```

## The Rejection Imperative

Reject when ANY of the following are true:

```
✗ Hypothetical failure
✗ Generic programming advice
✗ Imagined workflow
✗ Solution never tested
✗ Already documented elsewhere
✗ Obvious to any competent developer
```

## The Provenance Chain

Every experience must satisfy:

```
Observed → Solved → Verified → Captured

Each step must be real.
Each transition must be documented.
Each record must be from actual development event.
```

## The Evidence Rule

**Evidence is optional but encouraged when:**
- Failure is visual (use screenshots)
- Solution is complex (use logs)
- Proof is required for trust (use test output)

**Evidence is unnecessary when:**
- Situation/solution is self-explanatory
- Evidence is already in Git history
- Words alone convey the knowledge

## The Sixty-Second Rule

Capture must complete in under 60 seconds.

**Time budget:**
| Field | Maximum |
|-------|---------|
| situation | 15 seconds |
| solution | 15 seconds |
| verify | 10 seconds |
| tags | 5 seconds |
| evidence | 10 seconds (optional) |
| **Total** | **55 seconds** |

**If capture takes longer:**
- The experience may not be worth capturing
- The schema may be wrong
- Simplify before recording

## Situation Rule

Situation describes the development conditions under which the experience should be reused.

**DO NOT describe:**
- Commands
- Symptoms
- Error messages alone

**The future AI must determine:**
"Am I in the same situation?"

**Structure:**
```
When [context], [unusual condition] occurs
```

**Examples:**
```yaml
# ✓ Good
situation: "building Go dependencies for the first time in a sandbox with restricted network"

# ✗ Bad - symptom only
situation: "go build fails"

# ✗ Bad - command only  
situation: "running go build ./..."
```

## Confidence Assignment

| When | Assign |
|------|--------|
| First capture | Experimental |
| Worked 2-3 times | Low |
| Worked 3+ times across sessions | Medium |
| Worked consistently with no failures | High |
| Standard practice across repos | Proven |

## Review Checklist

Before committing a Development Experience:

- [ ] Is this from a real development event?
- [ ] Was the failure actually observed?
- [ ] Was the solution actually tested?
- [ ] Was the fix actually verified?
- [ ] Does situation describe context, not commands?
- [ ] Can someone answer "Am I in the same situation?" from this?
- [ ] Are tags specific enough for retrieval?
- [ ] Is evidence referenced when helpful?
