# KDSE Runtime Defect Correction

**Defect ID:** KDSE-RUNTIME-DEFECT-002  
**Title:** Legacy Runtime Installations Are Not Automatically Recognized  
**Document Version:** 1.0  
**Type:** Defect Correction Report  
**Effective Date:** 2026-07-10  
**Status:** Corrected

---

## Executive Summary

This document describes the correction of a runtime compatibility defect where legacy KDSE Runtime installations using `manifest.yaml` were not recognized during synchronization, causing the runtime to incorrectly report "KDSE Runtime is not installed."

**Solution:** Implemented automatic manifest format detection and seamless migration from YAML to JSON format.

---

## Engineering Evidence

### Issue Description

During a Runtime Synchronization of an existing KDSE installation, the runtime incorrectly reported:

```
KDSE Runtime is not installed.
```

The installation actually existed. The runtime expected `.kdse/manifest.json` while the existing installation contained `.kdse/manifest.yaml`.

### Root Cause

The Runtime assumed a single installation format (JSON). The `manifest_exists()` function in `common.sh` only checked for `manifest.json`:

```bash
# Original implementation
manifest_exists() {
    local manifest=$(get_manifest_path)
    [[ -f "$manifest" ]] && grep -q '"kdse"' "$manifest"
}
```

Where `get_manifest_path()` returned:
```
${KDSE_HOME}/${KDSE_DIR}/manifest.json
```

No mechanism existed to detect or handle legacy YAML manifests.

### Impact Assessment

- **User Impact:** High - Existing installations appeared broken
- **Operator Experience:** Poor - Required manual inspection and troubleshooting
- **System Reliability:** Compromised - False negative detection

---

## Design Decision

### Approach: Transparent Auto-Migration

After evaluating multiple approaches, we selected **transparent auto-migration** with the following rationale:

| Approach | Pros | Cons |
|----------|------|------|
| Require manual migration | Simple implementation | Poor UX, operator burden |
| Support both formats permanently | No migration complexity | Code complexity, dual format maintenance |
| Auto-migrate on first access | Seamless UX | One-time migration overhead |
| **Auto-migrate transparently** | **Best UX, single format long-term** | **Migration complexity (acceptable)** |

### Key Design Principles

1. **Zero Operator Intervention**: Migration happens automatically without prompting
2. **Idempotency**: Running migration multiple times produces the same result
3. **Data Preservation**: All user data (reports, history, runtime state, cache) preserved
4. **Auditability**: Migration events are logged and recorded in the manifest
5. **Backward Compatibility**: Both formats supported during detection phase

---

## Migration Strategy

### Detection Algorithm

```
1. Check if .kdse directory exists
2. If manifest.json exists with valid "kdse" marker:
   → Format = JSON (current)
3. Else if manifest.yaml exists with valid "kdse" marker:
   → Format = YAML (legacy)
4. Else:
   → No installation detected
```

### Migration Process

1. **Detect**: Identify legacy YAML manifest
2. **Backup**: Create `manifest.yaml.backup` for rollback
3. **Extract**: Parse version, installed date, repo, branch from YAML
4. **Transform**: Convert to JSON format with migration metadata
5. **Write**: Create new `manifest.json`
6. **Verify**: Confirm JSON manifest is valid
7. **Report**: Log migration completion

### Migration Metadata

The migrated manifest includes a `migration` section:

```json
{
  "kdse": "runtime-manifest",
  "version": "...",
  "installed": "...",
  "migration": {
    "from": "yaml",
    "timestamp": "2026-07-10T12:00:00Z"
  },
  ...
}
```

### Idempotency Guarantees

- Migration only proceeds if `manifest.json` does not exist
- YAML manifest is backed up, not deleted
- Running sync/verify multiple times has no additional effect
- Pre-existing JSON manifests are never modified

---

## Compatibility Policy

### Supported Installation Formats

| Format | Status | Detection | Migration |
|--------|--------|-----------|-----------|
| `manifest.json` | Current | ✅ Full | N/A |
| `manifest.yaml` | Legacy | ✅ Full | ✅ Automatic |

### Future Format Extensibility

The detection system is designed for extensibility:

```bash
detect_manifest_format() {
    # Priority order: newest to oldest
    if check_format "manifest.json"; then echo "json"; return; fi
    if check_format "manifest.yaml"; then echo "yaml"; return; fi
    echo "none"
}
```

When new formats are introduced:
1. Add new format check in priority order
2. Implement format-specific migration function
3. Update documentation

### Extent of Compatibility Support

- **Kept:** All user data directories (reports, history, runtime, cache, configuration)
- **Migrated:** Manifest format only
- **Removed:** Legacy `manifest.yaml` backup after successful migration (optional cleanup)

---

## Implementation Changes

### Files Modified

| File | Changes |
|------|---------|
| `common.sh` | Added detection, migration, and format functions |
| `install.sh` | Auto-migration before install checks |
| `sync.sh` | Auto-migration before sync operations |
| `verify.sh` | Auto-migration before verification |
| `uninstall.sh` | Cleanup of legacy manifest files |

### New Functions Added

| Function | Purpose |
|----------|---------|
| `detect_manifest_format()` | Returns "json", "yaml", or "none" |
| `get_legacy_manifest_path()` | Returns path to YAML manifest |
| `get_manifest_path_by_format()` | Returns path based on format |
| `migrate_manifest_yaml_to_json()` | Performs YAML to JSON migration |
| `auto_migrate_if_needed()` | Wrapper for automatic migration |

### Runtime Output Enhancement

All runtime commands now report:

```
Detected Installation Format: json
Installation Format: json
```

When migration is performed:

```
Detected Installation Format: yaml
Automatic migration will be performed...
Migration complete: manifest.yaml -> manifest.json
Migration Performed: YES (yaml -> JSON)
```

---

## Test Results

### Scenario 1: Fresh Installation

**Setup:** Clean environment, no existing `.kdse` directory  
**Action:** Run `./install.sh`  
**Expected:** Installation succeeds  
**Result:** ✅ PASS

```
$ ./install.sh
============================================================
 KDSE Runtime Installation
============================================================

Repository: https://github.com/tamzrod/KDSE.git
Branch:     main
Path:       /home/user/.kdse

[INFO] No existing KDSE installation detected
[INFO] Pre-installation checks passed
...
[OK] Installation verified successfully
```

### Scenario 2: Existing manifest.json

**Setup:** Existing installation with `manifest.json`  
**Action:** Run `./sync.sh`  
**Expected:** Sync succeeds without migration  
**Result:** ✅ PASS

```
$ ./sync.sh
============================================================
 KDSE Runtime Synchronization
============================================================

[INFO] Detected installation format: JSON (current)
[INFO] Pre-synchronization checks passed
...
[OK] Synchronization complete
```

### Scenario 3: Existing manifest.yaml (Legacy)

**Setup:** Existing installation with `manifest.yaml` (legacy)  
**Action:** Run `./sync.sh`  
**Expected:** Automatic migration, then sync succeeds  
**Result:** ✅ PASS

```
$ ./sync.sh
============================================================
 KDSE Runtime Synchronization
============================================================

[INFO] Detected installation format: YAML (legacy)
[INFO] Automatic migration will be performed...
[INFO] Migrating manifest.yaml to manifest.json...
[OK] Migration complete: manifest.yaml -> manifest.json
[INFO] YAML manifest backed up to: manifest.yaml.backup
[INFO] Pre-synchronization checks passed
...
[OK] Synchronization complete

Installation Format: json
Migration Performed: YES (yaml -> JSON)
```

### Scenario 4: Already Migrated Installation

**Setup:** Previously migrated installation (both `manifest.yaml.backup` and `manifest.json` exist)  
**Action:** Run `./sync.sh` twice  
**Expected:** No changes, no errors  
**Result:** ✅ PASS

```
$ ./sync.sh
============================================================
 KDSE Runtime Synchronization
============================================================

[INFO] Detected installation format: JSON (current)
[INFO] Pre-synchronization checks passed
...
[OK] Synchronization complete

$ ./sync.sh  # Second run
[INFO] Detected installation format: JSON (current)
[INFO] Pre-synchronization checks passed
...
[OK] Already up-to-date
```

### Test Summary

| Scenario | Status | Notes |
|----------|--------|-------|
| Fresh Installation | ✅ PASS | Clean install works correctly |
| Existing JSON | ✅ PASS | No migration needed |
| Legacy YAML | ✅ PASS | Auto-migration works |
| Already Migrated | ✅ PASS | Idempotent operation |

---

## Risks

### Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Migration data loss | Low | High | Backup before migration |
| Migration corruption | Low | Medium | JSON validation after write |
| Idempotency violation | Low | Medium | Check for existing JSON before migration |
| Partial migration state | Low | Low | Atomic file operations |

### Safety Mechanisms

1. **Backup First**: YAML manifest backed up before any changes
2. **Existence Check**: Migration only if JSON doesn't exist
3. **Validation**: New manifest validated before completing
4. **Rollback Path**: Backup file retained for manual recovery if needed

### Known Limitations

- Migration only handles basic manifest fields
- Complex YAML structures with nested data may lose precision
- Timezone information in YAML dates may differ slightly in JSON

---

## Lessons Learned

### Technical Insights

1. **Format Detection**: Always check for multiple formats before assuming
2. **Migration Safety**: Backup before migration, validate after
3. **Idempotency**: Design for repeated execution without side effects
4. **Logging**: Clear migration status helps debugging

### Process Insights

1. **Test Coverage**: Include legacy format in all test scenarios
2. **Documentation**: Document supported formats explicitly
3. **User Communication**: Show migration status in output

### Design Principles Applied

| Principle | Application |
|-----------|-------------|
| Fail gracefully | Continue if migration partially fails |
| Be transparent | Show migration status to operator |
| Preserve data | Never delete without backup |
| Be extensible | Design for future format support |

---

## References

- [Runtime Installation Framework README](README.md)
- [Runtime Architecture](../ARCHITECTURE.md)
- [Common Utilities Source](install/common.sh)

---

## Document History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-07-10 | Initial defect correction report |

---

*This document was created as part of KDSE-RUNTIME-DEFECT-002 correction. The fix ensures legacy installations are automatically recognized and migrated, providing a seamless operator experience.*
