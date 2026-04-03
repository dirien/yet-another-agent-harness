<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Principles guide

Detailed reference for each design principle audited by
`design-principles`. Covers violation patterns, language-specific
manifestations, and scoring guidance.

## SOLID

### S — Single Responsibility Principle

**Statement**: A module, class, or function should have one reason to
change — one primary actor or concern that drives its evolution.

**Violation patterns**:

- A file mixes HTTP handling, business logic, and database access
- A function validates input, transforms data, persists it, and sends
  a notification
- A struct holds both configuration and runtime state
- A module is changed by both the "frontend team" and the "data team"

**Language context**:

- **Go**: Look for large files in `handlers/` or `services/` that also
  contain SQL queries or external API calls. A handler that calls
  `db.Query(...)` directly is a SRP violation.
- **Python/JS**: Classes that inherit from `Model` and also contain
  email-sending logic.
- **Functional**: Functions with side effects mixed into pure
  transformation pipelines.

**Scoring**: Impact HIGH if the violation causes test difficulty or
frequent merge conflicts. MEDIUM if it is just conceptually mixed.

---

### O — Open/Closed Principle

**Statement**: Software entities should be open for extension but
closed for modification. Adding new behavior should not require
editing existing code.

**Violation patterns**:

- `switch`/`if-else` chains on a type discriminator field that grow
  with every new type
- A central "factory" function with a case for every concrete type
- Adding a new payment method requires editing the core payment processor
- No interface or strategy pattern at behavioral extension points

**Language context**:

- **Go**: Missing interfaces at package boundaries. Concrete struct
  types passed across layers. Absence of `io.Reader`/`io.Writer`-style
  abstractions where they would allow substitution.
- **OOP languages**: Prefer abstract classes/interfaces over concrete
  base class extension.

**Scoring**: Impact HIGH when the pattern is in frequently extended
code (e.g., an event system, a plugin registry). LOW if the switch is
stable and unlikely to grow.

---

### L — Liskov Substitution Principle

**Statement**: Subtypes must be substitutable for their base types
without altering the correctness of the program.

**Violation patterns**:

- An override method panics or throws on input the base accepts
- An override ignores a parameter that the base uses (weakening
  behavior)
- A subtype adds preconditions (e.g., "must call Init first") not
  present in the base
- `type assert` or `instanceof` checks before calling a method

**Language context**:

- **Go**: Interfaces are implicit. Look for `interface{}` or `any`
  with type switches that assume a specific concrete type. Also: a
  method on an interface that panics in one implementation but not
  others.
- **OOP languages**: `NotImplementedException` or empty override bodies
  are clear signals.

**Scoring**: Impact CRITICAL if the violation causes runtime panics or
data corruption. HIGH if it causes incorrect behavior. MEDIUM if it is
a design smell without current user impact.

---

### I — Interface Segregation Principle

**Statement**: Clients should not be forced to depend on methods they
do not use. Prefer small, focused interfaces.

**Violation patterns**:

- An interface with > 5 methods where most consumers only use 1–2
- A `Repository` interface requiring both read and write methods when
  read-only consumers exist
- Mock/stub implementations with empty bodies for unused methods
- A single "god interface" exported from a package

**Language context**:

- **Go**: The Go proverb "the bigger the interface, the weaker the
  abstraction" directly encodes ISP. A `io.Reader` (1 method) is
  better than a `ReadWriteCloser` passed to read-only code.
- **OOP languages**: Large abstract base classes with many abstract
  methods.

**Scoring**: Impact MEDIUM when stubs are required. HIGH when it
prevents reuse or forces unwanted dependencies.

---

### D — Dependency Inversion Principle

**Statement**: High-level modules should not depend on low-level
modules. Both should depend on abstractions. Abstractions should not
depend on details.

**Violation patterns**:

- A business logic function directly instantiates a database client
  (`db := sql.Open(...)`)
- No interface between the domain layer and the infrastructure layer
- Functions that are impossible to unit test without real external
  services
- Package-level `var` for concrete service instances

**Language context**:

- **Go**: Constructor injection via interfaces is idiomatic. Direct
  `sql.DB` in a business struct (not injected) is a DIP violation.
- **JS/TS**: Direct `require`/`import` of concrete modules at the call
  site instead of passing dependencies.
- **Python**: `import boto3` inside a function body with no injection
  point.

**Scoring**: Impact HIGH when it blocks testability. CRITICAL when it
creates environment coupling (e.g., must have a real database to run
any test).

---

## DRY — Don't Repeat Yourself

**Statement**: Every piece of knowledge must have a single,
unambiguous, authoritative representation within a system.

**Violation patterns**:

- The same validation logic implemented in multiple handlers
- A data transformation copied across services
- Constants defined in multiple packages
- The same SQL query text in multiple places

**False positives to avoid**:

- Two functions with similar shape but different semantics are NOT DRY
  violations — they will evolve independently
- Test code duplicating production behavior for clarity is acceptable
- Boilerplate that must remain separate for framework reasons (e.g.,
  two HTTP handlers with similar structure) is low-severity

**Scoring**: Impact scales with how often the duplication diverges.
If both copies must always change together, impact is HIGH. If they
are likely to diverge (by design), impact is LOW.

---

## YAGNI — You Aren't Gonna Need It

**Statement**: Do not add functionality until it is actually needed.

**Violation patterns**:

- Interfaces with a single implementation and no planned second
- Configuration options that are never read
- `// TODO: support multi-tenant` abstractions in a single-tenant app
- Generic type parameters where a concrete type would suffice
- Unused exported functions or exported constants

**False positives**:

- Framework-required boilerplate (e.g., interface stubs for test mocks)
- Standard library conventions (e.g., implementing `fmt.Stringer`)
- Deliberate extensibility points documented as architectural decisions

**Scoring**: Impact LOW to MEDIUM (dead code slows readers). Effort S
(delete it). Always P2 or lower unless it is causing active confusion.

---

## KISS — Keep It Simple

**Statement**: Most systems work best if they are kept simple rather
than made complicated.

**Violation patterns**:

- A 3-layer abstraction where a single function would work
- A plugin architecture for a system with 2 known use cases
- Custom serialization for a standard format (JSON, protobuf)
- Metaprogramming or reflection where direct code would be clearer
- Recursive data structures for inherently flat data

**False positives**:

- Complexity justified by measurable performance requirements
- Abstraction required for testability (see DIP)
- Well-known design patterns applied to appropriate problems

**Scoring**: Impact depends on team familiarity. MEDIUM for unexpected
complexity. HIGH if it is in the critical path and slows feature
development.

---

## Law of Demeter

**Statement**: A unit should only talk to its immediate collaborators.
Do not reach through objects to call methods on their internals.

**Violation patterns**:

- `order.GetCustomer().GetAddress().GetCity()`
- Calling `a.B.C.DoThing()` instead of `a.DoThing()` (which delegates)
- Functions that navigate object graphs to extract values

**Scoring**: Impact MEDIUM (tight coupling, brittle to refactoring).
Effort S–M (add a delegation method).

---

## Separation of Concerns

**Statement**: Different concerns of a system should be separated into
distinct sections with minimal overlap.

**Violation patterns**:

- HTTP status code logic in a domain service
- Formatting/rendering logic in a data access layer
- Database transactions spanning multiple bounded contexts in a
  single call
- Business rules scattered across UI event handlers

**Language context**:

- **Go**: `net/http` handler calling `database/sql` directly
- **Node**: Express route handler containing Mongoose queries
- **Python**: Django view containing raw SQL

**Scoring**: Impact HIGH when it prevents independent testing of layers.
CRITICAL when it creates cross-context transactions.

---

## Composition over Inheritance

**Statement**: Favor object composition over class inheritance for
achieving code reuse.

**Violation patterns**:

- Inheritance hierarchies deeper than 2 levels
- Base classes modified to satisfy subclass requirements
- Subclasses overriding many methods without calling `super`
- Behavior shared via inheritance that could be a injected collaborator

**Language context**:

- **Go**: Go has no class inheritance. Violations here mean: embedded
  structs used to share behavior across unrelated types, or missing
  interface usage where composition would enable substitution.
- **Python/JS/Java**: Deep class trees where mixins or strategy
  injection would be cleaner.

**Scoring**: Impact MEDIUM to HIGH depending on how brittle the
hierarchy has become.

---

## 12-Factor App (code-relevant subset)

Only three factors are code-level concerns. The others (build/release/
run, port binding, concurrency, disposability, dev/prod parity, admin
processes) are deployment concerns outside this skill's scope.

### Config (Factor III)

Hardcoded values that should come from environment variables:

- URLs, hostnames, port numbers in source code
- API keys, tokens, passwords as string literals
- Feature flags as constants
- Thresholds and limits that differ between environments

Check with:

```bash
grep -rn '"http://\|"https://' --include='*.go' . | grep -v '_test.go' | grep -v vendor
grep -rn 'localhost\|127\.0\.0\.1' --include='*.go' . | grep -v '_test.go'
```

### Logs (Factor XI)

Ad-hoc logging instead of structured streams:

- `fmt.Println`, `print()`, `console.log` in production paths
- Log lines with string concatenation instead of structured fields
- Log levels inconsistently applied

### Stateless processes (Factor VI)

In-memory state that prevents horizontal scaling:

- Package-level mutable globals holding user or session data
- Local file system caching of request-scoped data
- In-process caches without TTL or distributed backing

**Scoring**: Config violations are HIGH impact when they block
deployment to multiple environments. Logs and stateless are MEDIUM
unless the service is actively scaled.
