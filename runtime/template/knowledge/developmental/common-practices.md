# Common Engineering Practices

**Type:** Developmental Knowledge  
**Template:** Baseline  
**Version:** 1.0

---

## Purpose

This document catalogs common engineering practices applicable across software development. These practices improve code quality, team collaboration, and project maintainability.

---

## Code Quality Practices

### Write Clean Code

| Practice | Description |
|----------|-------------|
| Meaningful names | Variables, functions, classes named for their purpose |
| Single responsibility | Each function does one thing |
| Small functions | Functions under 20 lines preferred |
| Avoid repetition | DRY (Don't Repeat Yourself) |
| Comments explain why | Code explains how; comments explain why |

### Test Practices

| Practice | Description |
|----------|-------------|
| Test behavior | Test what, not how |
| Arrange-Act-Assert | Structure test clearly |
| One assertion | Prefer single assertion per test |
| Meaningful names | Test names describe expected behavior |
| Fast tests | Keep unit tests fast |

### Example

```python
# BAD
def test():
    x = calculate(1, 2)
    assert x == 3

# GOOD
def test_calculate_returns_sum_of_two_positive_integers():
    """
    When adding two positive integers,
    the result should be their sum.
    """
    # Arrange
    operand_a = 1
    operand_b = 2
    expected_sum = 3
    
    # Act
    result = calculate(operand_a, operand_b)
    
    # Assert
    assert result == expected_sum
```

---

## Version Control Practices

### Commit Messages

| Practice | Description |
|----------|-------------|
| Subject line | 50 chars max, imperative mood |
| Body | Explain what and why, not how |
| Reference issues | Include issue numbers |
| Atomic commits | One logical change per commit |

### Example

```
feat: add user authentication

Implement JWT-based authentication for API endpoints.
Users can now login and receive tokens for subsequent requests.

Closes #123
```

### Branching Strategy

| Branch | Purpose | Naming |
|--------|---------|--------|
| main/master | Production-ready | - |
| develop | Integration | develop |
| feature | New features | feature/description |
| bugfix | Bug fixes | bugfix/description |
| hotfix | Production fixes | hotfix/description |

---

## Review Practices

### Code Review

| Practice | For Author | For Reviewer |
|----------|------------|-------------|
| Small changes | Keep PRs under 400 lines | Review within 24 hours |
| Self-review | Review before requesting | Be constructive |
| Explanations | Explain context in PR | Ask questions, don't command |
| Approvals | Address feedback | Approve when acceptable |

### Review Checklist

- [ ] Code is correct
- [ ] Tests are adequate
- [ ] Code is readable
- [ ] Documentation updated
- [ ] No security issues
- [ ] Performance acceptable
- [ ] Error handling complete

---

## Error Handling Practices

### Principle: Fail Fast and Loud

```python
# BAD: Silent failure
def process(data):
    result = risky_operation(data)
    return result  # Might be None

# GOOD: Explicit failure
def process(data):
    if data is None:
        raise ValueError("data cannot be None")
    result = risky_operation(data)
    if result is None:
        raise ProcessingError("Operation failed")
    return result
```

### Error Response Pattern

```python
class APIError(Exception):
    def __init__(self, message, code, details=None):
        self.message = message
        self.code = code
        self.details = details or {}
        super().__init__(self.message)

# Usage
raise APIError(
    message="User not found",
    code="USER_NOT_FOUND",
    details={"user_id": user_id}
)
```

---

## Logging Practices

### Log Levels

| Level | Use | Example |
|-------|-----|---------|
| DEBUG | Detailed debugging info | "Entering function X" |
| INFO | General events | "User logged in" |
| WARNING | Potential issues | "Retry attempt 2" |
| ERROR | Errors that need attention | "Connection failed" |
| CRITICAL | System failures | "Database unreachable" |

### Structured Logging

```python
# BAD
logger.info(f"User {user_id} performed {action}")

# GOOD
logger.info(
    "User action performed",
    extra={
        "user_id": user_id,
        "action": action,
        "timestamp": datetime.utcnow().isoformat()
    }
)
```

---

## Security Practices

### Input Validation

```python
# Always validate input
def create_user(username, email, password):
    # Validate
    if len(username) < 3:
        raise ValidationError("Username too short")
    if "@" not in email:
        raise ValidationError("Invalid email")
    if len(password) < 8:
        raise ValidationError("Password too short")
    
    # Process
    ...
```

### Principle of Least Privilege

| Practice | Description |
|----------|-------------|
| Minimal permissions | Request only necessary access |
| Separate concerns | Different credentials for different purposes |
| Rotate secrets | Change credentials periodically |
| Audit access | Log who accessed what |

---

## Performance Practices

### Measure Before Optimizing

```python
# BAD: Premature optimization
# Optimization for speed before knowing it's slow
def find_user_fast(users, name):
    return [u for u in users if u.name == name]

# GOOD: Measure first
# Use profiling to identify bottlenecks
def find_user(users, name):
    return [u for u in users if u.name == name]

# If profiling shows this is slow:
# - Add index
# - Use database query
# - Cache results
```

### Caching Strategy

| Strategy | When to Use | Consideration |
|----------|-------------|---------------|
| No cache | Simple, fast operations | Simplicity |
| In-memory | Single instance, small data | Memory usage |
| Distributed | Multi-instance, large data | Consistency |
| CDN | Static assets | TTL management |

---

## Communication Practices

### Async Communication

| Channel | Use For | Not For |
|---------|---------|---------|
| Email | Formal, external | Urgent, internal |
| Chat | Quick questions | Complex discussions |
| Document | Decisions, specs | Real-time collaboration |
| Meeting | Brainstorming, alignment | Status updates |

### Documentation Updates

| Change | Documentation |
|--------|--------------|
| New feature | README, API docs |
| Breaking change | Migration guide |
| Bug fix | Release notes |
| Architecture | ADR |

---

*This document is developmental knowledge. Project-specific practices may be added in .kdse/knowledge/developmental/practices/.*
