# Core Target Standard (`.sets → any language`)

- **Standard version:** 0.2.0
- **Target language / version:** none. This is the base that every `.sets → <language>` standard extends.
- **Runtime / serialization context:** the canonical JSON projection (VM-JSON-1), unless the language standard overrides it
- **Targets `.sets` spec at commit:** `f5a56cc` (`GUIDE.md`)
- **Depends on:** [Core Value Model](../VALUE-MODEL.md), [`base.sets`](../base.sets)
- **Status:** draft

Rule IDs are `CT-*`. **Invariant** rules can't be overridden; **Default** rules can be, through the language standard's Overrides table ([standards/README.md](../../README.md)). Nothing in this standard depends on which language produced the `.sets` input.

<!-- mtoc-start -->

* [1. Scope](#1-scope)
* [2. Conventions](#2-conventions)
  * [2.1 Representation vocabulary](#21-representation-vocabulary)
  * [2.2 Exactness classes](#22-exactness-classes)
  * [2.3 Naming](#23-naming)
  * [2.4 Ordering and layout](#24-ordering-and-layout)
  * [2.5 Diagnostics](#25-diagnostics)
* [3. Axiom bindings](#3-axiom-bindings)
* [4. Decisions](#4-decisions)
* [5. Representation rules](#5-representation-rules)
  * [5.1 Declarations and imports](#51-declarations-and-imports)
  * [5.2 Rosters](#52-rosters)
  * [5.3 Scalar builders and refinements](#53-scalar-builders-and-refinements)
  * [5.4 Tuples](#54-tuples)
  * [5.5 Records](#55-records)
  * [5.6 Optional fields and null](#56-optional-fields-and-null)
  * [5.7 Unions](#57-unions)
  * [5.8 Intersections and differences](#58-intersections-and-differences)
  * [5.9 Collections](#59-collections)
  * [5.10 Transformers and generics](#510-transformers-and-generics)
  * [5.11 Predicates and value transformers](#511-predicates-and-value-transformers)
  * [5.12 `u` and `phi`](#512-u-and-phi)
  * [5.13 Recursion](#513-recursion)
* [6. Domain-form coverage](#6-domain-form-coverage)
* [7. Runtime support](#7-runtime-support)
* [8. Tests](#8-tests)
* [9. Open Gaps](#9-open-gaps)
* [10. Spec Issues](#10-spec-issues)

<!-- mtoc-end -->

## 1. Scope

**CT-SCOPE-1 (Invariant). Input.** Any `.sets` file that is valid under GUIDE and VALUE-MODEL, and whose domains all resolve to `base`, a bound language superdomain (§3), or declarations in the input.

**CT-SCOPE-2 (Invariant). Output.** For each named domain, the adapter generates:
1. a **type** (the representation, §2.1);
2. a **membership check** whenever the representation is over-approximate;
3. a **codec** for the wire (VM-JSON-1 by default), where decoding accepts every member and rejects every non-member.

**CT-SCOPE-3 (Default). Runtime support.** Generated code may depend on a small support library (§7). A language standard that forbids that inlines the helpers.

## 2. Conventions

### 2.1 Representation vocabulary

Language standards instantiate these abstract representations. Rules in §5 are stated in terms of them.

| Representation | Meaning | Typical instantiation |
| -------------- | ------- | --------------------- |
| `Scalar(b)` | A bound native type `b` (§3). | `string`, `int32`, `f64` |
| `Enum(m1…mn)` | Closed set of constants of one kind. | enum, string-literal union, consts |
| `Record(k1: R1 … kn: Rn; rest)` | Named fields with wire keys. An open record also has `rest`, which keeps fields the domain doesn't mention (CT-REC-2). | struct, class, interface, dataclass |
| `Tuple(R1 … Rn)` | Positional fields. | tuple, struct with `Item0…`, fixed array |
| `List(R)`, `Set(R)`, `Dict(RK, RV)` | Native collections. | slice, `Vec`, `list`, `Map` |
| `Presence(R)` | Field present or missing, distinguishable. | `Option<T>`, pointer plus `omitempty`, `?:` |
| `Nullable(R)` | Value or `null`, distinguishable. | `T \| null`, `Option<T>`, pointer |
| `Sum(B1 … Bn)` | Exactly one of disjoint branches. | tagged enum, sealed interface, discriminated union |
| `Overlap(B1 … Bn)` | One or more of overlapping branches at once. | struct of optional branch views (README Go `Human`) |
| `Checked(R, c)` | `R` plus a generated check `c`. | newtype with validating constructor and decoder |
| `Generic(P1 … Pn ⇒ R)` | Type parameters bounded by representations. | generics; or monomorphized copies |
| `Top`, `Bottom` | Any value; no value. | `any`/`interface{}`; `never`/`!` |
| `Opaque(axiom)` | Wrapper for an unbound axiom (CT-AXB-3). | struct around raw wire data |

### 2.2 Exactness classes

**CT-EX-1 (Invariant).** Every rule, binding and generated declaration has one class:

| Class | Meaning | Required |
| ----- | ------- | -------- |
| **exact** | The representation holds exactly the domain's members. | nothing |
| **over + check** | It holds a superset. | a generated check, run by constructors and decoders, that rejects non-members |
| **under** | Some members can't be represented. | forbidden by default; allowed only if the language standard names the lost values, and the adapter emits `CT-D-UNDER` every time |
| **unsupported** | No acceptable representation. | the adapter rejects the input with a named diagnostic |

**CT-EX-2 (Invariant). Silent wrong code is never an option.** If a domain form has no rule, the adapter reports `CT-D-UNSUPPORTED`. It doesn't fall back to a "close enough" type.

### 2.3 Naming

**CT-NAME-1 (Default).** A domain name becomes an identifier by the language's casing convention for types. Non-identifier characters are replaced by `_uXXXX` (the code point in uppercase hex). If the result is a reserved word, `_` is appended.

**CT-NAME-2 (Default). Synthesized names.** Anonymous domains that need a named type are named from their position. The names nest in this order:

| Position | Name |
| -------- | ---- |
| Field value of record `P` with key `k` | `P` + `Pascal(k)` |
| Component `i` (0-based) of tuple `P` | `P` + `Item` + `i` |
| Element of `Seq`/`SetOf` in `P` | `P` + `Item` |
| Key / value of `Map` in `P` | `P` + `Key` / `P` + `Value` |
| Branch `i` of a union `P` | `P` + the branch domain's name, or `P` + `Branch` + `i` if anonymous |
| Variant `Tag` of `P` | `P` + `Pascal(Tag)` |

**CT-NAME-3 (Invariant). Roster members and keys keep their spelling on the wire.** An identifier derived from a roster member or a record key is only a name. The codec always uses the exact value: the quoted string, the atom spelling, or the number.

**CT-NAME-4 (Invariant). Collisions.** If two names collide after CT-NAME-1 and CT-NAME-2, the later one in emission order (§2.4) gets the suffix `2`, then `3`, and so on.

### 2.4 Ordering and layout

**CT-ORD-1 (Invariant). Determinism.** The same `.sets` input produces byte-identical output.
- One `.sets` file (superdomain) maps to one target module. The language standard defines the path mapping.
- Declarations follow `.sets` declaration order. Synthesized types come right after the declaration that introduced them, in depth-first order.
- Record fields follow key order by code point. Tuple components keep their positions. Union branches keep their operand order with nested `union` calls flattened left to right. Enum members follow roster order.
- Imports are sorted by module path.

### 2.5 Diagnostics

| Code | Severity | Meaning |
| ---- | -------- | ------- |
| `CT-D-UNSUPPORTED` | error | The domain form has no representation rule. |
| `CT-D-UNDER` | warning | The representation loses members; the message names them. |
| `CT-D-UNBOUND` | warning | An axiom has no binding; `Opaque` is used (CT-AXB-3). |
| `CT-D-WIRE-AMBIGUOUS` | error | The wire projection isn't one-to-one on this domain (VM-JSON-2). |
| `CT-D-RECURSION` | error | The target can't represent this recursive domain. |

## 3. Axiom bindings

**CT-AXB-1 (Invariant). Bind by meaning.** Each language standard fills this table for its target. A domain is bound by its prose in VALUE-MODEL (or in the superdomain that declares it), never by the file it came from. The "Required exactness" column is the minimum a binding must reach without a `CT-D-UNDER` diagnostic.

| Domain (`base`) | Denotes | Required exactness | Check when over |
| --------------- | ------- | ------------------ | --------------- |
| `String` | Unicode scalar sequences | exact, or over + check where native strings allow lone surrogates or invalid bytes | valid scalar sequence |
| `Char` | One-scalar strings | exact, or over + check on `String` | length is one scalar |
| `Integer` | All integers | exact (big integer); a bounded native integer is **under** | — |
| `Natural` | Integers ≥ 0 | exact or over + check | `n ≥ 0` |
| `IntN`, `UintN` | Bounded integers | exact (native width) or over + check (wider native) | range |
| `Fractional` | Non-integer finite reals | exact (arbitrary-precision rational/decimal) or **under** | not an integer |
| `Real` | Finite reals | as `Fractional` | — |
| `Binary64Finite` | Finite binary64 values | over + check on native `f64` | finite |
| `Binary32Finite` | Finite binary32 values | exact or over + check | finite, representable in binary32 |
| `Float64`, `Float32` | Including `IEEESpecial` | exact on native floats | — |
| `Boolean` | `{true, false}` | exact | — |
| `Null` | `{null}` | exact (unit / null type) | — |
| `Record` | All records | `Record(; rest)`: a string-keyed `Dict` of raw values | — |

**CT-AXB-2 (Default). Native numeric choice.** A described integer subdomain binds to the narrowest native type that contains it. If none matches exactly, the binding uses a wider native type with a range check. A union of numeric domains is bound as a single scalar under CT-UNI-5.

**CT-AXB-3 (Invariant). Unbound axioms.** An axiom or described subdomain with no row in the language standard's table becomes `Opaque(name)`. That type keeps the raw wire value so it round-trips, and the adapter emits `CT-D-UNBOUND`. Guessing a native type from the axiom's name isn't allowed.

## 4. Decisions

| Question | Decision | Rule |
| -------- | -------- | ---- |
| Tuples | `Tuple`, positional components named `Item0…ItemN` where the language needs names. | CT-TUP-1 |
| Records | Recognized by shape (`r in Record` plus field clauses), represented as `Record`, with `rest` when open. | CT-REC-1, CT-REC-2 |
| Disjoint unions | `Sum`, decoded by trying branches in order; disjointness guarantees at most one accepts. | CT-UNI-1, CT-UNI-2 |
| Overlapping unions | Merged `Record` + check when all branches are records; otherwise `Overlap`. | CT-UNI-3, CT-UNI-4 |
| How disjointness is decided | A conservative shape algorithm; anything it can't prove disjoint is assumed to overlap. | CT-DISJ-1 |
| Intersections, differences | Representation of the narrower operand, plus a check. | CT-INT-1, CT-DIF-1 |
| Rosters | `Enum` per kind; mixed kinds become a `Sum` of enums. | CT-ROS-* |
| Transformers returning domains | Generics, or monomorphization per call site. | CT-GEN-* |
| Predicates, value transformers | Functions, used inside checks. | CT-FN-* |
| `u` / `phi` | `Top` / `Bottom`, or rejection. | CT-TOP-* |
| Imports | Module imports; `base` is bound, not emitted. | CT-DECL-3 |

## 5. Representation rules

Each rule states its exactness. Examples live in [`../examples/`](../examples/): each file holds the canonical domain and member / non-member vectors. Language standards add `.sets` ↔ `<lang>` pairs.

### 5.1 Declarations and imports

**CT-DECL-1 (Invariant). Top-level axiom.** Bound per §3. It is exact or over + check, or `Opaque` with a diagnostic.

**CT-DECL-2 (Default). Alias.** `A = B` becomes a type alias when the language has transparent aliases (exact). Otherwise it reuses `B`'s type under both names.

**CT-DECL-3 (Invariant). Imports.** `X in m` becomes an import of `X`'s type from the module generated for `m`. Imports from `base` generate nothing; their domains are bound (§3) or built in (§5.9, §5.12).

**CT-DECL-4 (Invariant). Described subdomains.** `X sub Y` with a prose extent is bound like an axiom (CT-AXB-1). Without a binding, it becomes `Checked(Y's representation, c)` if the prose supplies a check. With no prose extent, GUIDE reads `X` as the same as `Y` ("essentially becoming an alias"), so `X` gets `Y`'s representation (exact).

### 5.2 Rosters

**CT-ROS-1 (Invariant). One-kind rosters.** A roster whose members are all strings, all atoms, or all numbers becomes `Enum` (exact). The codec writes each member's exact value (CT-NAME-3).

**CT-ROS-2 (Invariant). Mixed rosters.** A roster with members of several kinds becomes `Sum` of one `Enum` per kind, in the order each kind first appears (exact). If two members project to the same wire value (`Circle` and `"Circle"`), the adapter reports `CT-D-WIRE-AMBIGUOUS`.

**CT-ROS-3 (Invariant). Singleton.** A singleton roster is a one-member `Enum`, or a constant where the language has literal types (exact). It's never widened to the member's kind without a check.

**CT-ROS-4 (Invariant). Domain-valued members.** A roster whose member is a domain symbol (`{Number}`) is a domain of domains. It is **unsupported** (`CT-D-UNSUPPORTED`) until a language standard defines a representation for domain values.

**CT-ROS-5 (Invariant).** `{Boolean members}` and `{null}` use the bound `Boolean` and `Null` types, never a generated enum.

**Example:** [`examples/02-rosters.sets`](../examples/02-rosters.sets)

### 5.3 Scalar builders and refinements

**CT-BLD-1 (Invariant).** `{ x suchthat x in D }` is `D`'s representation (exact).

**CT-BLD-2 (Invariant). `not in` and `difference`.** `{ x suchthat x in A and x not in B }` and `difference(A, B)` become `Checked(A's representation, not member of B)` (over + check).

**CT-BLD-3 (Invariant). Equality.** `x in singleton(v)` is equivalent to the roster `{v}` (CT-ROS-3).

**CT-BLD-4 (Invariant). Predicate rules.** A rule that calls a predicate becomes `Checked`, with the predicate emitted per CT-FN-1.

### 5.4 Tuples

**CT-TUP-1 (Default).** `{ (x1 … xn) suchthat x1 in D1 and … }` becomes `Tuple(R1 … Rn)` (exact if each `Ri` is exact). Where the language needs field names, components are named `Item0 … ItemN-1`. The wire form is a JSON array of length `n`.

**CT-TUP-2 (Invariant). No flattening.** A nested tuple `((x), (y))` is a `Tuple` whose components are 1-tuples. A 1-tuple `(x)` is `Tuple(R)` and never `R`.

**CT-TUP-3 (Invariant).** The empty tuple `()` is a unit-like `Tuple()` that projects to `[]`.

**Example:** [`examples/05-tuple.sets`](../examples/05-tuple.sets)

### 5.5 Records

**CT-REC-1 (Invariant). Recognition.** A builder `{ r suchthat r in Record and C1 and … and Cn }` is a record when every clause `Ci` is one of:
- a required field: `("k", b) in r and b in V`, where `"k"` is a string literal and `b` is a binder used nowhere else;
- an optional field: `r in OptionalField({"k"}, V)`;
- a closing clause: `r in OnlyKeys(K)`, where `K` is a roster of string literals.

Keys are pairwise distinct. The representation is `Record(k: R(V) … )`, with optional fields as `Presence(R(V))` (CT-OPT-1). A required or optional field is exact when `R(V)` is exact.

**CT-REC-2 (Invariant). Open records keep unknown fields.** Without `OnlyKeys`, the record is open (VM-REC-3), and the representation gets a `rest` part: a string-keyed `Dict` of raw values. Decoding puts every unmentioned field in `rest`, and encoding writes it back, so open records round-trip exactly. Dropping `rest` is **under** (the lost values are every record with an extra field) and needs `CT-D-UNDER`.

**CT-REC-3 (Invariant). Closed decoding.** With `OnlyKeys`, there's no `rest`, and decoding rejects any key outside `K`.

**CT-REC-4 (Invariant). Other record rules.** A record rule with `or`, `not in` or predicate clauses becomes `Checked(Record(…), rule)` over the fields the rule mentions (over + check). If no field structure can be found, it's `Checked(Record(; rest), rule)`.

**Example:** [`examples/03-record.sets`](../examples/03-record.sets)

### 5.6 Optional fields and null

**CT-OPT-1 (Invariant).** A field `r in OptionalField({"k"}, V)` is `Presence(R(V))`. On the wire, a missing field is an omitted key.

**CT-OPT-2 (Invariant).** A domain `union(V, Null)` is `Nullable(R(V))`.

**CT-OPT-3 (Invariant). Three states.** `OptionalField({"k"}, union(V, Null))` is `Presence(Nullable(R(V)))`. All three states (missing, `null`, value) must survive a round trip. If the language's idiomatic optional type merges them (Go `*T` with `omitempty`), the language standard prescribes a three-state wrapper. Merging them is **under** (the lost value is the missing key or `null`) and needs `CT-D-UNDER`.

**Example:** [`examples/04-optional-nullable.sets`](../examples/04-optional-nullable.sets)

### 5.7 Unions

**CT-DISJ-1 (Invariant). Disjointness.** The adapter proves two domains `A` and `B` disjoint when any of these holds. Otherwise it assumes they overlap.
1. Their members have different kinds (VM-KIND-1): string, number, atom, tuple, finite domain.
2. Both are tuples of different arity.
3. Both are records, one requires key `k`, and the other is closed with `OnlyKeys` that excludes `k`.
4. Both are records that require the same key `k` with provably disjoint value domains (this is how tag fields discriminate).
5. Both are rooted in different top-level axioms, or are rosters with no common member, or one is a roster none of whose members is in the other.
6. Both are `Variant`s with different tags.
7. Both are described subdomains whose prose extents are disjoint (for example `Int8` and a roster of strings; not `Int8` and `Uint8`).

Rule 1 treats `Seq` (whose members are `()` and 2-tuples) as overlapping any 0- or 2-tuple domain, unless the components prove otherwise.

**CT-UNI-1 (Invariant). Disjoint unions.** If all branches are pairwise disjoint, the union becomes `Sum(B1 … Bn)` (exact).

**CT-UNI-2 (Invariant). Decoding a `Sum`.** The decoder tries branches in operand order and takes the first that accepts. Disjointness means at most one accepts. When a cheap discriminator exists (JSON kind, required key, tag value), the language standard may use it, but the result must be the same as trying in order. Encoding writes the active branch's value.

**CT-UNI-3 (Invariant). Overlapping unions keep the overlap.** If any two branches may overlap, the representation must be able to hold a value that belongs to both, and must report membership in each branch independently. A representation in which "is `A`" implies "is not `B`" is wrong. This is the README `Human` case.

**CT-UNI-4 (Invariant). Overlapping record unions merge.** If every branch is a record, the union becomes one `Record` with:
- one field per key mentioned by any branch, typed by the union of the value domains the branches give that key, and wrapped in `Presence` unless every branch requires it;
- `rest` if any branch is open. A value for a mentioned key that falls outside that field's type (an open branch admits anything for keys it doesn't mention) is kept in `rest` under its key, so nothing is lost;
- a check that the value is in at least one branch, and a generated `Is<Branch>()` per branch that tests that branch alone (over + check).

```text
Human = union(Person, Aged)          (both open)
  ⇒ Record(age: Presence(Real), name: Presence(String); rest)
    check:   IsPerson() or IsAged()
    IsPerson():  name present            IsAged():  age present
    {"name":"abc","age":123}  ⇒  both true
    {"name":"abc","age":"x"}  ⇒  name set, "age" kept in rest, IsPerson() true, IsAged() false
```

This is the README Go `Human`, with `Person` and `Aged` as views over one record instead of two embedded structs.

**CT-UNI-5 (Default). Same-kind scalar unions.** A union of numeric domains, or of string domains, becomes the narrowest bound scalar that contains every branch, plus a check that the value is in some branch (over + check). For example, `union(Int32, Uint32)` becomes a 64-bit integer with a check.

**CT-UNI-6 (Invariant). Other overlapping unions.** These become `Overlap(B1 … Bn)`: one optional slot per branch. Decoding sets **every** slot whose branch accepts the value, and at least one must accept. Encoding writes the value from the first set slot. All set slots hold the same value, so they project the same way (over + check: the check is that at least one slot is set and all set slots agree).

**Example:** [`examples/06-union-disjoint.sets`](../examples/06-union-disjoint.sets), [`examples/07-union-overlapping.sets`](../examples/07-union-overlapping.sets), [`examples/08-tagged-union.sets`](../examples/08-tagged-union.sets)

### 5.8 Intersections and differences

**CT-INT-1 (Invariant).** `intersect(A, B)` becomes `Checked(R(A), member of B)`, where `A` is the operand with the narrower representation (over + check). Two records intersect fieldwise into one `Record`: every key of either side, required if either side requires it, typed by the intersection of both sides' value domains, with `rest` only if both sides are open (exact when each field intersection is exact). If a closed side excludes a key the other side requires, the result is `phi` (CT-TOP-2).

**CT-DIF-1 (Invariant).** `difference(A, B)` becomes `Checked(R(A), not member of B)` (over + check). If `B` is provably disjoint from `A`, it's just `R(A)` (exact).

**Example:** [`examples/09-intersection-difference.sets`](../examples/09-intersection-difference.sets)

### 5.9 Collections

**CT-COL-1 (Invariant).** `Seq(T)` becomes `List(R(T))` (exact). The adapter recognizes `Seq` by its `base` binding and never generates nested pairs.

**CT-COL-2 (Default).** `SetOf(T)` becomes `Set(R(T))` when `R(T)` is hashable or comparable in the target (exact). Otherwise it's `Checked(List(R(T)), no duplicates)` (over + check). The wire form follows VM-JSON-1: a sorted array.

**CT-COL-3 (Default).** `Map(K, V)` becomes `Dict(R(K), R(V))` when `R(K)` can be a native key (exact). Otherwise it's `Checked(List(Tuple(R(K), R(V))), unique keys)` (over + check).

**Example:** [`examples/10-collections.sets`](../examples/10-collections.sets)

### 5.10 Transformers and generics

**CT-GEN-1 (Default).** A transformer that returns a domain and has only domain parameters `{T sub X}` becomes `Generic`. Each parameter becomes a type parameter bounded by `R(X)`, or is unbounded when `X` is `u`. If the target has no generics, or can't express the bound, each call site is monomorphized under CT-NAME-2, using the call's argument names (`Box(String)` → `BoxString`).

**CT-GEN-2 (Invariant). Call sites.** `N = F(A1 … An)` becomes an instantiation alias, or a monomorphized type named `N`.

**CT-GEN-3 (Invariant).** `base` operations and structural transformers (`union`, `intersect`, `difference`, `singleton`, `product2`, `OptionalField`, `OnlyKeys`, `Variant`, `Seq`, `SetOf`, `Map`) are recognized by binding. They aren't emitted as generics.

**Example:** [`examples/11-generic.sets`](../examples/11-generic.sets)

### 5.11 Predicates and value transformers

**CT-FN-1 (Default).** A predicate becomes a validation function that returns a boolean. It is emitted only if a check uses it.

**CT-FN-2 (Default).** A transformer that returns a scalar or a tuple is code, not a type. It's emitted as a function only when a check uses it. Otherwise it's skipped, without a diagnostic.

### 5.12 `u` and `phi`

**CT-TOP-1 (Invariant).** `u` becomes `Top`. Its codec accepts any wire value and keeps it as raw data, so it round-trips exactly. A target with no top type reports `CT-D-UNSUPPORTED`.

**CT-TOP-2 (Invariant).** `phi` becomes `Bottom`. If the target has no uninhabited type, it uses a type with no public constructor and a decoder that always fails (exact). A record field of domain `phi` makes the whole record `phi`; the adapter reports this as a warning, because it's usually a source mistake.

### 5.13 Recursion

**CT-RECUR-1 (Invariant).** A self-referential domain becomes a recursive type. The language standard adds whatever indirection the target needs (a box, pointer or reference). It keeps the representation exact; the indirection isn't a value. A target that can't do this reports `CT-D-RECURSION`.

**Example:** [`examples/12-recursion.sets`](../examples/12-recursion.sets)

## 6. Domain-form coverage

| Domain form | Exactness / status | Rule |
| ----------- | ------------------ | ---- |
| Top-level axiom | exact or over + check (per binding) | CT-DECL-1, §3 |
| Alias | exact | CT-DECL-2 |
| Superdomain import / file | exact | CT-DECL-3 |
| Described subdomain | per binding, or over + check | CT-DECL-4 |
| Roster: strings / atoms / numbers | exact | CT-ROS-1 |
| Roster: mixed kinds | exact (or ambiguity error) | CT-ROS-2 |
| Singleton | exact | CT-ROS-3 |
| Roster of domain symbols | unsupported | CT-ROS-4 |
| Scalar builder `x in D` | exact | CT-BLD-1 |
| `not in`, `difference` | over + check | CT-BLD-2, CT-DIF-1 |
| Tuple, nested tuple, empty tuple | exact | CT-TUP-* |
| Record, open (with `rest`) or closed | exact | CT-REC-1 to CT-REC-3 |
| Record with other rule clauses | over + check | CT-REC-4 |
| `OptionalField`, nullable, three-state | exact | CT-OPT-* |
| Disjoint union | exact | CT-UNI-1 |
| Overlapping record union | over + check | CT-UNI-4 |
| Same-kind scalar union | over + check | CT-UNI-5 |
| Other overlapping union | over + check | CT-UNI-6 |
| `intersect` | over + check | CT-INT-1 |
| `Seq` | exact | CT-COL-1 |
| `SetOf`, `Map` | exact or over + check | CT-COL-2, CT-COL-3 |
| Transformer returning a domain (generics) | exact | CT-GEN-1 |
| Predicate, value transformer | code, not type | CT-FN-* |
| `u`, `phi` | exact | CT-TOP-* |
| Recursive domain | exact | CT-RECUR-1 |
| Cardinality constraints | gap | §9 |

Totals over 26 forms: 14 exact, 6 over + check, 3 exact or over + check depending on the binding, 1 unsupported, 1 gap, and 1 that is code rather than a type (predicates and value transformers).

## 7. Runtime support

The smallest helper API that the rules assume. Language standards instantiate it or inline it.

```text
Presence<T>       missing | present(T)
Nullable<T>       null | value(T)
Raw               an unparsed wire value (for Top, Opaque and rest)
check(value) → ok | error(path, rule)      generated per Checked type
decode(wire) → T | error(path, rule)       runs check before returning
encode(T) → wire                            never fails for a checked value
```

Errors carry a path (`$.items[2].age`) and the domain that rejected the value.

## 8. Tests

**CT-TEST-1 (Invariant).** For every rule a language standard instantiates, its test suite includes:
- at least one **member** that round-trips: value → encode → decode → an equal value (VM-EQ-1);
- at least one **non-member** that decode rejects;
- for overlapping unions, one value in **both** branches whose membership reports are both true after decoding.

The member and non-member vectors in [`../examples/`](../examples/) are the shared baseline.

## 9. Open Gaps

**G-T1. Cardinality.** A domain can't state "exactly n members" or "at most n", so collections have no length constraints beyond tuples. Interim: none needed. The source side reports `CS-D-WIDENED`.

**G-T2. Domain values.** Rosters of domain symbols (CT-ROS-4) need a value-level representation of domains. Interim: unsupported.

**G-T3. Disjointness beyond shape.** CT-DISJ-1 is conservative. For example, it can't prove `{ x suchthat x in Integer and IsEven(x) }` and `{ x suchthat x in Integer and IsOdd(x) }` disjoint, so their union becomes `Overlap` (correct, but heavier than needed). *Non-normative proposal:* a `disjoint` declaration that source adapters may emit and target adapters trust only after checking it against vectors.

## 10. Spec Issues

Both directions share one list, kept in [VALUE-MODEL.md §10](../VALUE-MODEL.md#10-spec-issues).
