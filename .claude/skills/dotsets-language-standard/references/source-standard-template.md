# `<Language>` → `.sets` Standard

- **Standard version:** 0.1.0
- **Source language / version:** `<e.g. TypeScript 5.x>`
- **Targets `.sets` spec at commit:** `<git short hash of GUIDE.md>`
- **Status:** draft | proposed | accepted

## 1. Scope

What "view" of values this standard maps (e.g. JSON-serializable values only; runtime values including functions; wire format of a specific encoder). What is explicitly out of scope.

## 2. Conventions

- Universal domain spelling used (`u` or `U`).
- Base imports every output file starts with, e.g.

  ```sets
  u in base
  phi in base
  union in base
  intersect in base
  difference in base
  singleton in base
  ```

- Language axioms import, e.g. `String in <language>`.
- Naming: how source type names become `.sets` domain names (and how collisions/escapes are handled).

## 3. Axioms superdomain (`<language>.sets`)

Table of every top-level axiomatic domain:

| Domain | Denotes (values) | Disjoint from | Source types built from it |
| ------ | ---------------- | ------------- | -------------------------- |

Then derived scalar domains (`Number = union(Integer, Decimal)`, etc.) with one-line justification each.

## 4. Decisions

One subsection per cross-cutting question (record encoding, field identity, optional/absent/null, open vs closed records, nominal vs structural, generics, recursion, collections, modules). Each: **Decision**, **Alternatives considered**, **Why**, **Consequence for target adapters**.

## 5. Mapping rules

Repeat per construct:

### 5.N `<Construct name>`

**Source form**

```<lang>
...
```

**Value semantics** — which values are members (prose, with spec citation if subtle).

**Canonical `.sets`**

```sets
...
```

**Edge cases** — overlapping values, extra fields, empty, absent, null, recursion; state membership of each.

**Rationale** — why not the visually similar mapping.

**Example:** `examples/<file>.<ext>` ↔ `examples/<file>.sets`

## 6. Construct coverage

| Construct | Status (mapped / out of scope / gap / N/A) | Rule § | Notes |
| --------- | ------------------------------------------ | ------ | ----- |

## 7. Open Gaps

Constructs that can't be expressed in current `.sets` grammar. Each: what's missing, interim adapter behavior (reject / approximate / emit diagnostic), and an optional **non-normative** grammar proposal.

## 8. Spec Issues

Ambiguities or contradictions found in `README.md` / `GUIDE.md` while writing this standard, with the interpretation this standard assumed.
