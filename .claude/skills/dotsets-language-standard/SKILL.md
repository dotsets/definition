---
name: dotsets-language-standard
description: Build a language standard for `.sets` in either direction — a source standard (`<lang> → .sets`, how a language's type system maps onto `.sets` domains) or a target standard (`.sets → <lang>`, how any `.sets` domain is represented in a language such as Go, Rust, TypeScript, Python, Protobuf, JSON Schema). Produces a normative document plus supporting `.sets`/example files. Use when asked to create, draft, extend, or review a "language standard", "source/target adapter spec", "<lang> → .sets mapping", ".sets → <lang> mapping", "codegen rules for .sets", or to decide how a specific type construct or `.sets` domain form should be translated.
---

# dotsets-language-standard

`.sets` is a semantic boundary between languages (README, "Source and Target Independence"):

```text
Source Language → [source standard] → .sets domain → [target standard] → Target Language
```

- A **source standard** is the contract for a `Source → .sets` adapter. For every type construct of one language it answers: *which values does this construct admit, and what is the canonical `.sets` text for that domain?*
- A **target standard** is the contract for a `.sets → Target` adapter. For every `.sets` domain form it answers: *what target-language representation holds exactly the members of this domain, and how does the target check membership where its type system can't?*

The two sides never know about each other. A target standard **must not assume any particular source language**; a source standard must produce output a target adapter can consume with only `.sets` knowledge.

## Step 0 — Pick the direction and re-read the spec

1. Determine the direction from the request: "TypeScript to sets" → source; "generate Go from sets" → target. If the user asks for "a standard for `<lang>`" without a direction, ask whether they want source, target, or both. For both, build the source standard first, then the target standard, and keep them in separate documents.
2. The spec is v0 and changes often. Read the current `README.md` and `GUIDE.md` at the repo root in full before writing anything. Don't rely on memory of earlier versions. If anything below conflicts with `GUIDE.md`, `GUIDE.md` wins, and you note the conflict in the standard's **Spec Issues** section.

## Shared principles (both directions)

1. **Domains are sets of values.** Every decision is justified by membership: which concrete values are in, which are out. Use the README's `Person | Aged` example as the reference case: `{name: "abc", age: 123}` is a valid `Human`, so neither a naive `union(Person, Aged)` (source side) nor a Go constraint interface (target side) is correct.
2. **No primitives.** Atomic kinds (`String`, `Integer`, …) are declared top-level axiomatic domains, not built-ins.
3. **Top-level domains are disjoint.** Separate axioms never overlap; wider kinds are built with `union`.
4. **Parentheses are values.** `"abc" != ("abc")`. Nothing is implicitly flattened, and tuple nesting is significant.
5. **No unknown domains.** Every domain referenced is declared, constructed, imported (`X in file`), or the universal domain `u`/`U`.
6. **Only real grammar is normative.** Use only constructs that `GUIDE.md` defines: top-level axioms, roster `{a, b}`, builder `{ m suchthat rule }`, `in`, `sub`, `and`, `or`, `not in`, transformers/predicates `F({X sub u}): ...`, superdomain imports, tuples `( )`. Don't invent syntax silently. Anything else goes under **Open Gaps** as a clearly marked *non-normative proposal*.
7. **Language-specific representation stays on its side of the boundary.** JSON tags, naming conventions, package paths and codegen hints never go into `.sets`. Conversely, a target standard doesn't read meaning into `.sets` beyond what the domain says.

Don't edit `README.md` or `GUIDE.md`; they are hand-authored. If a standard exposes a spec problem, list it under **Spec Issues** and tell the user.

---

## A. Source standard (`<lang> → .sets`)

### Deliverables

Unless the user says otherwise, produce these under `standards/<language>/source/`:

1. `STANDARD.md`: the normative standard, following [references/source-standard-template.md](references/source-standard-template.md).
2. `<language>.sets`: the language's axioms superdomain (the file name is the domain name, per GUIDE "`.sets` files are superdomains"). It holds the top-level axiomatic domains and any language-specific helper transformers. It imports generic operations from `base` (`union in base`, `u in base`, …) instead of redefining them.
3. `examples/*.<ext>` + `examples/*.sets` pairs: each source snippet next to its canonical `.sets` output. Every mapping rule has at least one pair.

### Workflow

1. **Scope.** Pin down the version and the "view" of values: for example, TypeScript *as JSON over the wire* vs. TypeScript *runtime objects*. The view changes the domain.
2. **Inventory.** Walk [references/source-construct-checklist.md](references/source-construct-checklist.md). Mark each construct **mapped**, **out of scope** (with reason), **gap**, or **N/A**.
3. **Axioms.** Choose the minimal set of disjoint atoms that can build every scalar type via `union`/`difference`/roster. Describe in prose what values each one denotes, and decide edge values (NaN, ±Inf, −0, the empty string) explicitly. Literal types become roster/singleton domains; use quotes when a literal collides with a domain symbol (`{"Number"}`).
4. **Mapping rules.** For each construct, write the source form, the value semantics, **one canonical `.sets` form** (adapters must emit it so outputs are comparable), an example pair, the edge cases (overlap, extra keys, empty, absent, null, recursion) and the rationale for not using the visually similar mapping.
5. **Decisions.** Answer explicitly:
   - **Record encoding:** ordered tuple `(a, b)` or unordered domain `{a, b}`; nested or flattened composition. (GUIDE uses `{name, age}`; README uses `(name)`. Choose one and say why.)
   - **Field identity:** how field names survive (positional only, or a field-name atom paired with the value). If the grammar can't carry them, record a Gap.
   - **Optional / absent / null / undefined:** which are distinct values, and how absence is represented.
   - **Open vs. closed records:** whether extra fields are admitted.
   - **Nominal vs. structural:** whether same-shaped named types are the same domain, and how a brand is represented.
   - **Generics:** map to transformers with domain parameters, e.g. `Box({T sub u}): { (v) suchthat v in T }`; a bound `T extends X` maps to `{T sub X}`.
   - **Recursion**, **variable-length collections**, and **modules → superdomain files**.
6. **Validate** (shared checklist below, plus):
   - [ ] No two axioms that should overlap are declared as separate top-level domains.
   - [ ] Each union/intersection rule is tested against an overlapping value and a value with an extra field.
   - [ ] Someone who knows only `.sets` could write a target adapter from the outputs and the axioms prose.

---

## B. Target standard (`.sets → <lang>`)

### Deliverables

Unless the user says otherwise, produce these under `standards/<language>/target/`:

1. `STANDARD.md`: the normative standard, following [references/target-standard-template.md](references/target-standard-template.md).
2. `bindings.md` (or a section of `STANDARD.md`): the **axiom binding table**, mapping each known axiomatic domain to its target representation.
3. `examples/*.sets` + `examples/*.<ext>` pairs: each `.sets` domain next to its generated target code, including the membership check where one is needed. Every rule has at least one pair.
4. Optionally, a runtime-support sketch: the minimal helper library the generated code depends on (validators, tagged-value wrappers).

### Target-specific principles

1. **Exactness, or a documented approximation.** For each domain form, the target representation is classified as one of:
   - **Exact:** the target type admits exactly the domain's members.
   - **Over-approximate + check:** the target type admits a superset, and the standard prescribes a generated membership check (validator, constructor, `Unmarshal` hook) that rejects non-members. This is the usual way to handle `difference`, `not in`, refinements and overlapping unions.
   - **Under-approximate:** some members are unrepresentable. This is **forbidden by default**. Allow it only if the standard names the lost values and the adapter emits a diagnostic every time it happens.
   - **Unsupported:** the adapter rejects the input with a named diagnostic. That is better than silent wrong code.
2. **Shape follows the domain, not a guessed source.** Don't reverse-engineer what the source language "probably" was. `{(x, y) suchthat (x) in Person or (y) in Aged or ...}` is a domain of 2-tuples with optional components, whatever language produced it. The README's Go `Human` with pointer embeddings is the reference example: it represents the domain, not TypeScript syntax.
3. **The wire is part of the domain.** If generated types serialize, decoding must accept every member and (with checks) reject non-members. Test round-trips: member → encode → decode → the same member.
4. **Deterministic output.** Given the same `.sets` input, the adapter emits byte-identical code. Define naming (domain name → identifier, collision escaping, quoted roster symbols → identifiers), member ordering and file layout.
5. **Axioms are bound by meaning, not by file of origin.** The binding table maps an axiom to a target type based on the prose semantics in its superdomain (e.g. "`Float64`: IEEE-754 binary64"). For an axiom with no binding, define the fallback (an opaque wrapper type plus a diagnostic, or rejection). If axioms only exist in source-specific superdomains (e.g. `typescript.sets`), that couples the target to the source; record it under **Spec Issues** and propose a shared axioms superdomain.

### Workflow

1. **Scope.** The target language version, the runtime/serialization context (e.g. Go + `encoding/json`, Rust + `serde`, Python + `pydantic`) and the minimum toolchain. Say whether generated code may depend on a runtime support library.
2. **Inventory.** Walk [references/target-domain-checklist.md](references/target-domain-checklist.md). It lists every `.sets` domain form, not source constructs. Mark each **exact**, **over + check**, **under (documented)**, **unsupported**, or **gap** (the `.sets` grammar can't express it yet, so there is nothing to map).
3. **Axiom bindings.** Fill the binding table: axiom → target type → exactness → the check needed (e.g. `Integer` → Go `int64` is under-approximate for big integers; say so, or bind to `*big.Int`).
4. **Representation rules.** For each domain form, write the `.sets` form, its membership semantics, the **canonical target representation**, its exactness class, the generated membership check (if any), the wire behavior, an example pair, edge cases (overlapping members, the empty domain `phi`, `u`, nested tuples, quoted roster symbols) and the rationale for the alternatives you rejected.
5. **Decisions.** Answer explicitly:
   - **Tuples:** how n-tuples and nested tuples become target types (struct with positional fields, array, record), and how positional fields are named when the domain gives no names.
   - **Unions:** disjoint (a tagged sum type, sealed interface, enum with payloads) vs. overlapping (embedding with optional parts, as in the README's Go `Human`). Say how the adapter decides disjointness: disjoint axioms, disjoint rosters, otherwise assume overlap.
   - **Intersections and differences:** usually the representation of the wider domain plus a check.
   - **Rosters:** enum or constants; how symbols that aren't valid identifiers are escaped.
   - **Transformers with domain parameters:** generics where the target has them (`{T sub X}` → type parameter bounded by X's representation); otherwise monomorphize per call site. Transformers returning scalars or vectors (not domains) are code, not types: decide whether they're emitted as functions.
   - **Predicates:** emitted as validation functions/methods, or used only to build checks.
   - **`phi` and `u`:** the empty/uninhabited type (`never`, `!`, an uninstantiable type) and the top type (`any`, `interface{}`, `object`), or rejection.
   - **Superdomain imports:** how `X in file` maps to target packages/modules/imports.
6. **Validate** (shared checklist below, plus):
   - [ ] Every rule has an exactness class; every over-approximation has a prescribed check; every under-approximation names the lost values and has a diagnostic.
   - [ ] Nothing in the standard depends on knowing the source language.
   - [ ] Overlapping-union examples (including the README `Human`) produce a representation that holds a value belonging to both branches.
   - [ ] Round-trip and rejection tests are specified for at least one member and one non-member per rule.
   - [ ] The naming and ordering rules make output deterministic.

---

## Shared validation checklist

- [ ] Every `.sets` snippet uses only grammar present in the current `GUIDE.md`; anything else sits inside a labeled non-normative proposal.
- [ ] Every referenced domain is declared, constructed, imported, or `u`.
- [ ] Tuple arity and nesting are consistent across all rules.
- [ ] Open Gaps and Spec Issues are listed, not hidden.
- [ ] Every rule has at least one example pair on disk.

## Report

Tell the user:
- the direction and the files written
- the coverage counts (mapped / exact / approximate / gap / unsupported)
- the key decisions (record encoding and optionality for source; union representation and the axiom bindings for target)
- any Spec Issues that need the spec author's decision

Keep it short; the standard holds the detail.

## Known spec rough edges (verify against current GUIDE)

These existed when this skill was written. Check whether they're resolved before you mention them:

- The universal domain is spelled both `u` and `U` (GUIDE aliases `U = u` in `base`). Pick one.
- Record members appear both as `{name, age}` (GUIDE builder section) and as `(name)` (README, GUIDE tuples section).
- `%`, `=` and `sub or =` appear in a non-normative GUIDE note but are not defined grammar.
- There is no grammar yet for cardinality, variable-length sequences, field names or recursion.
- There is no shared axioms superdomain, so target bindings risk depending on source-specific superdomains.
