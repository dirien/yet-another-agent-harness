<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# OWASP Top 10 analysis patterns

For each category, investigate actual code paths — not just grep
matches. Read surrounding code to determine if a match is a real
vulnerability, a false positive, or already mitigated.

## A01: Broken access control

**Investigate:**

- API handlers/routes: do they check authorization before acting?
- Direct object references: can users access resources by changing
  IDs in URLs?
- Missing function-level access control: are admin endpoints
  protected?
- CORS misconfiguration: is `Access-Control-Allow-Origin: *` used?

**Evidence to look for:**

- Handlers without middleware/decorator for auth checks
- ID parameters used directly in database queries without ownership
  verification
- Admin routes accessible without role checks
- CORS allowing all origins in production config

## A02: Cryptographic failures

**Investigate:**

- Password storage: bcrypt/argon2/scrypt, or weak MD5/SHA1?
- Data in transit: TLS enforced? HTTP redirects to HTTPS?
- Sensitive data at rest: encrypted? Which algorithm?
- Random number generation: crypto/rand or math/rand?

**Language-specific patterns:**

- Go: `crypto/rand` (secure) vs `math/rand` (insecure for security)
- Node: `crypto.randomBytes` (secure) vs `Math.random` (insecure)
- Python: `secrets` module (secure) vs `random` (insecure)

## A03: Injection

**Investigate:**

- SQL: string concatenation in queries vs parameterized statements
- Command injection: user input in exec/system calls
- NoSQL: `$where`, `$regex`, unsanitized operators in MongoDB queries
- LDAP, XPath, template injection patterns
- GraphQL: unbounded queries, introspection in production

**Key distinction:** A query using `fmt.Sprintf("SELECT * FROM users
WHERE id = %s", userInput)` is vulnerable. The same query using
`db.Query("SELECT * FROM users WHERE id = $1", userInput)` is safe.
Read the actual query construction, not just the presence of `Query`.

## A04: Insecure design

**Investigate:**

- Rate limiting on authentication endpoints
- Account lockout after failed login attempts
- Input validation at trust boundaries (API entry points)
- Business logic flaws: can workflows be bypassed or reordered?
- Missing CAPTCHA on public-facing forms

## A05: Security misconfiguration

**Investigate:**

- Debug mode enabled in production configs
- Default credentials still active
- Unnecessary features enabled (directory listing, stack traces)
- Security headers missing (CSP, HSTS, X-Frame-Options)
- Overly permissive CORS, file permissions, or IAM policies
- Error messages leaking internal details

## A06: Vulnerable and outdated components

**Tools to run (if available):**

```bash
# Multi-language
trivy fs --severity HIGH,CRITICAL .

# Go
govulncheck ./...

# Node
npm audit --audit-level=moderate

# Python
pip-audit

# Rust
cargo audit
```

**Also check:** Are lockfiles committed? Does CI run vulnerability
scanning? When were dependencies last updated?

## A07: Identification and authentication failures

**Investigate:**

- Password complexity requirements enforced?
- Session tokens: secure, httpOnly, sameSite flags on cookies?
- JWT: short expiry? Refresh token rotation? Algorithm pinned
  (not `alg: none`)?
- Credential stuffing protection: rate limiting + account lockout?
- Password reset: token-based with expiry, or insecure
  (security questions, email link without token)?

## A08: Software and data integrity failures

**Investigate:**

- CI/CD pipeline: are secrets injected securely?
- Dependency resolution: lockfile committed and verified?
- Code signing: are releases signed?
- Deserialization: untrusted data deserialized without validation?

**Language-specific deserialization risks:**

- Python: unsafe deserialization (pickle, shelve, yaml.load)
- Node: dynamic code execution on user input
- Go/Java: custom deserialization of untrusted input

## A09: Security logging and monitoring failures

**Investigate:**

- Are login attempts (success and failure) logged?
- Are authorization failures logged?
- Are logs protected from tampering?
- Is sensitive data (passwords, tokens, PII) excluded from logs?
- Is there alerting on suspicious patterns?

## A10: Server-side request forgery (SSRF)

**Investigate:**

- HTTP clients that accept user-provided URLs
- URL validation: allowlist vs blocklist approach?
- Internal network access from user-controlled requests
- Redirect following: disabled or limited?
- Cloud metadata endpoints accessible (169.254.169.254)?
