# Retrieval Rules

## Discovery Mechanisms

### Tag-Based Search

Primary discovery mechanism.

```bash
# Find all playwright experiences
search: [playwright]

# Find screenshot experiences
search: [screenshot, playwright]

# Find docker + build experiences  
search: [docker, build]
```

### Situation Matching

Secondary mechanism: AI reads situation fields.

```
When AI encounters a problem:
1. Query experiences by relevant tags
2. Read situation fields
3. Match: "Am I in the same situation?"
4. If match, apply solution
5. Verify with verify field
```

### Confidence Filtering

Use confidence to prioritize:

```yaml
# Prefer proven/high first
# Fall back to medium
# Treat experimental as hints
# Ignore low confidence unless no alternatives
```

## Retrieval Decision Tree

```
AI encounters problem
        │
        ▼
Search experiences by problem tags
        │
        ▼
┌───────────────────────┐
│ Match found?          │
└───────────────────────┘
        │
   YES  │  NO
        ▼        ▼
Apply    Fall back to
solution investigation
        │
        ▼
Did investigation succeed?
        │
   YES  │  NO
        ▼        ▼
Capture  Skip
new DX   (known unsolved)
```

## Quality Signals

| Signal | Meaning |
|--------|---------|
| Proven confidence | Safe to use, no alternatives needed |
| High confidence | Safe to use, well-validated |
| Medium confidence | Reasonable to try, may need adjustment |
| Low confidence | Experimental, try with caution |
| Experimental | Hint only, validate independently |

## Retrieval Best Practices

### Before Applying an Experience

1. **Verify situation match**
   - Read the situation field carefully
   - Ask: "Am I in the same situation?"
   - If uncertain, proceed with caution

2. **Check evidence**
   - Review any attached evidence
   - Understand the original context
   - Look for similar experiences

3. **Apply verification first**
   - Run the verify command first
   - Confirm the problem exists before applying solution
   - Don't apply solution blindly

### When Multiple Experiences Match

1. **Prefer higher confidence**
   - Proven > High > Medium > Low > Experimental

2. **Prefer more specific tags**
   - `[playwright, screenshot, chromium, docker]` > `[playwright]`

3. **Consider recency**
   - Newer experiences may reflect current environment
   - Older experiences may be outdated

### When No Experience Matches

1. **Investigate independently**
   - This is new knowledge to capture
   - Follow the Capture Rules

2. **Document the investigation**
   - Even failed attempts provide learning
   - Capture what didn't work

3. **Share findings**
   - Future sessions may face similar problems
   - Your investigation enables their solution

## Experience Lifecycle in Retrieval

```
┌─────────────┐
│ Discovered  │ ← Tag search finds experience
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Evaluated   │ ← Situation match checked
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Applied     │ ← Solution executed
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Verified    │ ← Verify method confirms
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Validated   │ ← Success confirmed
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Elevated    │ ← Confidence may increase
└─────────────┘
```
