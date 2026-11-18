# `.sets` → `<Language>` Standard

- **Standard version:** 0.1.0
- **Target language / version:** `<e.g. Go 1.23>`
- **Runtime / serialization context:** `<e.g. encoding/json>`
- **Targets `.sets` spec at commit:** `<git short hash of GUIDE.md>`
- **Status:** draft | proposed | accepted

## 1. Scope

The input this adapter accepts (any valid `.sets`? only files that import specific superdomains?). The generated artifacts (types, validators, (de)serializers). Whether a runtime support library is allowed. What is explicitly out of scope.

## 2. Conventions

- The universal domain spelling assumed (`u` or `U`).
- **Naming:** domain name → identifier, how quoted roster symbols and non-identifier symbols are escaped, collision resolution.
- **Ordering and layout:** the order declarations are emitted in, and how files/packages are laid out, so output is deterministic.
- **Diagnostics:** the list of diagnostic codes the adapter emits (unsupported form, under-approximation, unbound axiom, …).

## 3. Axiom bindings

| Axiom (superdomain) | Denotes (from axioms prose) | Target type | Exactness | Check required | Notes |
| ------------------- | --------------------------- | ----------- | --------- | -------------- | ----- |

**Unbound axiom fallback:** what happens when an axiom has no row.

## 4. Decisions

One subsection per cross-cutting question (tuples, disjoint vs overlapping unions and how disjointness is decided, intersections/differences, rosters, transformers/generics, predicates, `phi`/`u`, imports). Each has: **Decision**, **Alternatives considered**, **Why**, **Exactness consequence**.

## 5. Representation rules

Repeat per domain form:

### 5.N `<Domain form>`

**`.sets` form**

```sets
...
```

**Membership semantics:** which values are members.

**Canonical target representation**

```<lang>
...
```

**Exactness:** exact | over + check | under (lost values: …) | unsupported (diagnostic: …)

**Membership check:** the generated validator/constructor/decode hook, if any.

**Wire behavior:** how members encode and decode; how non-members are rejected.

**Edge cases:** overlap, `phi`, `u`, nested tuples, quoted symbols, empty roster.

**Rationale:** why not the alternatives.

**Example:** `examples/<file>.sets` ↔ `examples/<file>.<ext>`

## 6. Domain-form coverage

| Domain form | Exactness / status | Rule § | Notes |
| ----------- | ------------------ | ------ | ----- |

## 7. Runtime support (optional)

The minimal helper API the generated code relies on.

## 8. Tests

For each rule: at least one member that must round-trip and one non-member that must be rejected.

## 9. Open Gaps

`.sets` forms that the grammar can't express yet, or target limitations with no acceptable representation. For each: interim adapter behavior and an optional **non-normative** proposal.

## 10. Spec Issues

Ambiguities or contradictions found in `README.md` / `GUIDE.md`, with the interpretation this standard assumed.
