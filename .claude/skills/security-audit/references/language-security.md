<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Language-specific security patterns

Vulnerability patterns unique to specific languages and ecosystems.
Check only the sections relevant to the detected project type.

## Go

**Error swallowing:**

```go
// VULNERABLE: error ignored, operation continues silently
result, _ := dangerousOperation()

// SAFE: error checked and handled
result, err := dangerousOperation()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

Search for `_ =` assignments where the second value is an error.
In security-critical code paths, every error must be handled.

**Goroutine leaks:** Goroutines that block forever on channels or
network calls without timeouts. Check for `go func()` without
context cancellation.

**Unsafe package:** Any use of `unsafe.Pointer` bypasses Go's type
safety. Flag and verify necessity.

**Integer overflow:** Go does not panic on integer overflow. Check
arithmetic on user-provided sizes/counts used for allocation.

## JavaScript / TypeScript

**Prototype pollution:** Merging user-controlled objects into
application objects without sanitization. Check for deep merge,
`Object.assign`, spread operators on untrusted input.

**Dynamic code execution:** Functions that interpret strings as
code (such as `Function()` constructor, `setTimeout` with string
argument, or `vm` module context execution) with user input are
remote code execution vulnerabilities.

**RegExp denial of service (ReDoS):** Regular expressions with
nested quantifiers on untrusted input can cause catastrophic
backtracking. Check patterns like `(a+)+`, `(a|a)+`, `(a|b)*c`.

**Dependency surface area:** Node projects can pull hundreds of
transitive dependencies. Each is an attack surface. Check total
package count and depth.

**XSS in frameworks:** React's `dangerouslySetInnerHTML`, Vue's
`v-html` directive, and Angular's `[innerHTML]` binding with
unsanitized user data bypass framework XSS protections.

## Python

**Unsafe deserialization:** Python's `pickle.load()` on untrusted
data enables arbitrary code execution. The same applies to `shelve`,
`marshal`, and PyYAML's `yaml.load()` — use `yaml.safe_load()`
instead. Flag any deserialization of untrusted input.

**Format string injection:** `str.format()` or f-strings with
user-controlled format strings can leak variables via attribute
traversal chains.

**Subprocess with shell=True:** `subprocess.call(cmd, shell=True)`
with user input enables command injection. Use `shell=False` with
argument lists instead.

**ORM raw queries:** Raw SQL execution that bypasses ORM protections.
Use parameterized queries with bound parameters.

**Django/Flask specific:** Debug mode in production, SECRET_KEY
exposure, missing CSRF middleware, missing authentication decorators
on views.

## Rust

**Unsafe blocks:** Any `unsafe {}` block bypasses Rust's ownership
and borrowing guarantees. Each unsafe block should have a safety
comment explaining why it's correct.

**Unchecked array indexing:** `array[i]` panics on out-of-bounds.
Use `.get(i)` for checked access on untrusted indices.

**FFI boundaries:** `extern "C"` functions lose Rust's safety
guarantees. Validate all data crossing FFI boundaries.

## .NET / C\#

**SQL injection via string interpolation:**
`$"SELECT * FROM users WHERE id = {userId}"` is vulnerable.
Use parameterized queries with `SqlParameter`.

**XML external entity (XXE):** `XmlDocument` with default settings
processes external entities. Set `XmlResolver = null`.

**Unsafe deserialization:** Binary formatters and similar
serialization mechanisms on untrusted data enable remote code
execution. Use `System.Text.Json` instead.

**Path traversal:** `Path.Combine(basePath, userInput)` does not
prevent `../` traversal. Validate the resolved path starts with
the base path.
