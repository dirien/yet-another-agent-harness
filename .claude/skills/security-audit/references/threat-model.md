<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Threat modeling

Apply STRIDE to the actual codebase — not a generic checklist. Map
real entry points, data flows, and trust boundaries.

## Step 1: Identify entry points

Find every way data enters the application:

- **HTTP/API endpoints** — routes, handlers, controllers
- **CLI arguments** — command-line parsers, flag definitions
- **Environment variables** — configuration loaded at startup
- **File system reads** — config files, user uploads, data imports
- **Database queries** — data from external databases
- **Message queues** — event consumers, webhook receivers
- **Inter-service calls** — gRPC, REST, GraphQL clients

For each entry point, note: what data comes in, who can send it,
and what validation exists.

## Step 2: Map trust boundaries

A trust boundary is where data crosses from a less-trusted zone to
a more-trusted zone:

- User input → application logic (most critical)
- External API response → application logic
- Database result → application logic (usually trusted, but verify)
- Application → filesystem (write operations)
- Application → external service (outbound, data leakage risk)

Draw the boundary: everything outside is untrusted, everything
inside is trusted. Validate at every crossing.

## Step 3: Apply STRIDE

For each entry point and trust boundary crossing, evaluate:

**Spoofing** — can an attacker impersonate a legitimate user or
service?

- Check: authentication at this entry point
- Check: certificate validation for service-to-service calls
- Check: token verification (JWT signature, expiry, issuer)

**Tampering** — can an attacker modify data in transit or at rest?

- Check: TLS for all network communication
- Check: integrity checks on file uploads or downloads
- Check: CSRF protection on state-changing operations
- Check: signed cookies/tokens vs unsigned

**Repudiation** — can an attacker perform an action and deny it?

- Check: are security-relevant actions logged?
- Check: are logs tamper-proof (append-only, centralized)?
- Check: is there enough context in logs to reconstruct events?

**Information disclosure** — can an attacker access unauthorized
data?

- Check: error messages — do they leak stack traces, SQL, or
  internal paths?
- Check: API responses — do they return more data than needed?
- Check: logging — is sensitive data (passwords, tokens, PII)
  written to logs?
- Check: directory listing, source maps, debug endpoints exposed?

**Denial of service** — can an attacker make the system unavailable?

- Check: rate limiting on public endpoints
- Check: request size limits (body, file uploads)
- Check: query complexity limits (GraphQL depth, SQL joins)
- Check: timeout configuration for external calls
- Check: resource exhaustion (memory, connections, file handles)

**Elevation of privilege** — can an attacker gain higher access?

- Check: role-based access control enforced consistently
- Check: privilege escalation via parameter manipulation
- Check: admin functionality accessible without admin role
- Check: insecure direct object references (IDOR)

## Step 4: Prioritize threats

Score each identified threat:

| Factor | Scale |
| ------ | ----- |
| Damage potential | 1-3 (what's the worst case?) |
| Reproducibility | 1-3 (how easy to exploit?) |
| Exploitability | 1-3 (what skill level needed?) |
| Affected users | 1-3 (how many impacted?) |
| Discoverability | 1-3 (how easy to find?) |

Sum = DREAD score (5-15). Map to severity:

- 12-15: CRITICAL
- 9-11: HIGH
- 6-8: MEDIUM
- 5: LOW
