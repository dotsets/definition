# Core Source Standard (`any language → .sets`)

- **Standard version:** 0.2.0
- **Source language / version:** none. This is the base that every `<language> → .sets` standard extends.
- **Targets `.sets` spec at commit:** `f5a56cc` (`GUIDE.md`)
- **Depends on:** [Core Value Model](../VALUE-MODEL.md), [`base.sets`](../base.sets)
- **Status:** draft

Rule IDs are `CS-*`. **Invariant** rules can't be overridden; **Default** rules can be, through the language standard's Overrides table ([standards/README.md](../../README.md)).

<!-- mtoc-start -->

* [1. Scope](#1-scope)
* [2. Conventions](#2-conventions)
* [3. Axioms](#3-axioms)
* [4. Decisions](#4-decisions)
* [5. Mapping rules](#5-mapping-rules)
  * [5.1 Scalars](#51-scalars)
  * [5.2 Literals and enumerations](#52-literals-and-enumerations)
  * [5.3 Records](#53-records)
  * [5.4 Optional and nullable](#54-optional-and-nullable)
  * [5.5 Positional tuples and fixed arrays](#55-positional-tuples-and-fixed-arrays)
  * [5.6 Untagged unions](#56-untagged-unions)
  * [5.7 Tagged unions](#57-tagged-unions)
  * [5.8 Intersections](#58-intersections)
  * [5.9 Top and bottom](#59-top-and-bottom)
  * [5.10 Collections](#510-collections)
  * [5.11 Aliases, nominal types and brands](#511-aliases-nominal-types-and-brands)
  * [5.12 Generics](#512-generics)
  * [5.13 Recursion](#513-recursion)
  * [5.14 Refinements](#514-refinements)
  * [5.15 Behavior](#515-behavior)
  * [5.16 Modules](#516-modules)
* [6. Construct coverage](#6-construct-coverage)
* [7. Open Gaps](#7-open-gaps)
* [8. Spec Issues](#8-spec-issues)

<!-- mtoc-end -->

## 1. Scope

**CS-SCOPE-1 (Invariant). Every language standard declares its view.** The view is the set of values the standard's domains describe: "TypeScript as JSON over the wire", "Go values as `encoding/json` decodes them", "Python runtime objects". The same type can denote different domains in different views, so a standard without a declared view isn't usable.

**CS-SCOPE-2 (Default). The core view is data.** Unless a language standard says otherwise, the view is serializable data: the values of the [Core Value Model](../VALUE-MODEL.md), carried by the canonical JSON projection (VM-JSON-1). Behavior (functions, channels, promises) and identity (pointers as addresses) are out of scope (§5.15).

## 2. Conventions

**CS-CONV-1 (Invariant). Imports.** Every emitted file begins with imports of exactly the `base` domains it uses, sorted by domain name, then imports from other superdomains, grouped by superdomain name in sorted order:

```sets
Int64 in base
Record in base
String in base
union in base

Address in geo
```

**CS-CONV-2 (Invariant). No unknown domains.** Every domain referenced is declared, constructed, imported or `u` (GUIDE, "Domain Validation"). Adapters don't rely on implementation magic.

**CS-CONV-3 (Default). Domain names.** A named source type becomes a domain of the same name. If the name contains a reserved character (`,`, `{`, `}`), the character is escaped as GUIDE describes (`\,`). If two source types in one module would produce the same domain name, the language standard defines the disambiguation.

**CS-CONV-4 (Invariant). Determinism.** The same source input produces byte-identical `.sets`:
- declarations follow source declaration order, after the imports;
- record fields are listed in key order by Unicode code point (keys are values, so their order carries no meaning);
- roster members follow source order, and duplicate members are dropped;
- `union`/`intersect` chains nest right-associatively, `union(A, union(B, C))`, with operands in source order;
- one declaration per line group, four-space indentation, `//` comments only where §7 requires a diagnostic note.

**CS-CONV-5 (Invariant). Binder names.** The record binder is `r`. Each field's value binder is the key's spelling when that's a valid symbol and doesn't collide with `r` or a domain in scope; otherwise it's `_` followed by the key spelling with reserved characters escaped. Binders in other builders are `x`, `y`, `z`, then `x1`, `x2`, …. Binders are scoped symbols, so this rule only exists to make outputs comparable.

**CS-CONV-6 (Invariant). Rule shape.** Rules evaluate left to right with equal precedence (VM-RULE-1). An adapter emits a rule only when the left-to-right reading is the intended one. When it isn't (for example `(A and B) or (C and D)`), the adapter names the parts and combines them with `union`/`intersect`.

## 3. Axioms

**CS-AX-1 (Invariant). Reuse `base`.** Scalars map to `base` axioms and described subdomains whenever the prose meaning matches (VALUE-MODEL §4–§5). A language standard declares a new axiom only under VM-AX-5.

**CS-AX-2 (Invariant). Disjointness check.** Before declaring an axiom, the language standard shows that no value of the new axiom is a value of an existing axiom. If two of the language's types share values (`int` and `float` both hold `1`), they're described subdomains of one axiom, or unions of axioms, never two axioms.

**CS-AX-3 (Default). Language superdomain.** Language-specific domains go in `<language>.sets`. That file imports from `base` and redefines nothing it provides.

## 4. Decisions

| Question | Decision | Rule |
| -------- | -------- | ---- |
| Record encoding | A set of `(key, value)` tuples: `{ r suchthat r in Record and ("name", name) in r and name in String }`. | VM-REC-*, CS-REC-1 |
| Field identity | Keys are string values inside each tuple. They aren't binder names. | VM-REC-1 |
| Optional / missing / null | A missing field is a missing tuple. Optional fields use `OptionalField`. `null` is a separate atom. | VM-MISS-*, CS-OPT-* |
| Open vs closed | Records are open by default. Closed records add `r in OnlyKeys({…})`. | VM-REC-3, VM-REC-4, CS-REC-4 |
| Nominal vs structural | Structural. A name is a label for a domain. Brands that exist in the value are tags. | CS-NOM-* |
| Generics | Transformers with domain parameters; a bound `T extends X` becomes `{T sub X}`. | CS-GEN-* |
| Recursion | Self-reference by name. | VM-RECUR-1, CS-RECUR-1 |
| Variable-length collections | `Seq`, `SetOf` and `Map` from `base`. | CS-COL-* |
| Modules | One superdomain file per module; imports become `X in module`. | CS-MOD-* |

**Why open records by default.** With open records, README's `Human` is simply `union(Person, Aged)`: `{name: "abc", age: 123}` is a `Person` (it has a `name`) and an `Aged` (it has an `age`), so it's in the union. The visually similar mapping is also the correct one, and no merged rewrite is needed. Languages whose view rejects unknown keys add `OnlyKeys`, and then the same `union` means "exactly one shape".

## 5. Mapping rules

Each rule gives the abstract source construct, the canonical `.sets`, edge cases, and an example in [`../examples/`](../examples/). Those examples hold the canonical domain with member and non-member vectors. Language standards add `<lang>` ↔ `.sets` pairs on top of them.

### 5.1 Scalars

**CS-SCA-1 (Invariant).** A scalar type maps to the narrowest `base` domain whose prose extent equals the type's extent in the view:

| Abstract construct | Canonical `.sets` |
| ------------------ | ----------------- |
| Unicode text | `String` |
| Single character (a Unicode scalar value) | `Char` |
| Arbitrary-precision integer | `Integer` |
| N-bit signed / unsigned integer | `IntN` / `UintN` |
| IEEE binary64 including NaN and infinities | `Float64` |
| IEEE binary64 without them (for example JSON) | `Binary64Finite` |
| Boolean | `Boolean` |
| Null / nil / None | `Null` |
| Unit `()` | `{()}` |
| Byte string | `Seq(Uint8)` |

**CS-SCA-2 (Invariant). Exact or wider, never narrower.** If no `base` domain matches exactly, the language standard declares a described subdomain (VM-SUB-2). It doesn't silently map to a narrower or wider domain. A wider mapping is allowed only with a diagnostic comment on the declaration, for example `Decimal = Real // CS-D-WIDENED: scale and precision not encoded`.

**Edge cases.** Code units versus scalar values (UTF-16 strings that can contain lone surrogates aren't `String` in the core sense; the language standard declares its own described domain). An integer type the view serializes as a JSON string is a `String` subdomain in that view, not an `Integer`.

**Example:** [`examples/01-scalars.sets`](../examples/01-scalars.sets)

### 5.2 Literals and enumerations

**CS-LIT-1 (Invariant).** A literal type is a singleton roster: `"on"` → `{"on"}`, `3` → `{3}`, `true` → `{true}`.

**CS-LIT-2 (Invariant).** A union of literal types is one roster in source order: `"a" | "b" | 3` → `{"a", "b", 3}`.

**CS-ENUM-1 (Default).** An enumeration maps by its value in the view. If members carry string values, they're quoted strings: `{"Succeeded", "Failed"}`. If they carry numbers, they're numeric symbols. If the members are values in their own right with no underlying string or number (Rust unit variants, sealed singletons), they're atoms: `{Succeeded, Failed}`.

**CS-NUM-3 (Default).** A numeric literal is written as its exact decimal value (VM-LEX-2). Every binary64 value has a finite decimal expansion.

**Edge cases.** A literal that spells a domain name in scope becomes `"Number"` as a string or `@Number` as an atom (VM-LEX-4). The empty roster is `phi`, not `{}` (CS-TOP-2).

**Rationale.** Writing enumerations as quoted strings by default would merge atoms and strings, and VM-JSON-2 would then fail to detect the ambiguity.

**Example:** [`examples/02-rosters.sets`](../examples/02-rosters.sets)

### 5.3 Records

**CS-REC-1 (Invariant).** A record, struct, object or class shape maps to the canonical record form (VM-REC-2). `r in Record` comes first; then one clause per field in key order, using binders per CS-CONV-5:

```sets
User = {
    r suchthat
    r in Record and
    ("age", age) in r and age in Int64 and
    ("name", name) in r and name in String
}
```

A required field is the pair `("k", k) in r and k in T`. An optional field is `r in OptionalField({"k"}, T)` in the same key-ordered position.

**CS-REC-2 (Default). Keys are wire keys.** A field's key is its name in the view. In a JSON view that's the serialized key after renames (`json:"user_name"` → `"user_name"`), not the source identifier. Language-specific spellings never appear in `.sets` except as key values.

**CS-REC-3 (Invariant). Nested records stay nested.** A field whose type is another record refers to that record's domain: `("address", address) in r and address in Address`. Embedding, spread and inheritance are resolved by the language standard into the fields the view actually has. If the view flattens embedded fields (Go embedded structs in JSON), the record lists the flattened keys, because that's the value. Otherwise it keeps the nesting.

**CS-REC-4 (Default). Openness follows the view.** Records are open (VM-REC-3). If the view rejects unknown keys (`additionalProperties: false`, `#[serde(deny_unknown_fields)]`, a strict decoder), the record ends with `r in OnlyKeys({k1, …, kn})` listing every key it declares, in key order.

**Edge cases.** An empty open record `{}` → `Record`. An empty closed record → `{ r suchthat r in Record and r in OnlyKeys(phi) }`, whose only member is `{}`. Duplicate keys can't occur (VM-MAP-1). Default values aren't domain information; a field with a default is required or optional according to the view.

**Example:** [`examples/03-record.sets`](../examples/03-record.sets)

### 5.4 Optional and nullable

**CS-OPT-1 (Invariant). Optional field.** A field that can be missing in the view uses `OptionalField`:

```sets
r in OptionalField({"age"}, Integer)
```

**CS-OPT-2 (Invariant). Nullable value.** A value that can be null uses `union(T, Null)`, wherever it appears.

**CS-OPT-3 (Invariant). Keep the states apart.** A field that can be missing, null or a value is `r in OptionalField({"k"}, union(T, Null))`. Only merge missing and `null` when the view itself can't tell them apart, and say so in the language standard. Example: Go `*T` with `omitempty` decodes both a missing key and `null` to `nil`.

**CS-OPT-4 (Default). Language "undefined".** A language with an `undefined` value that differs from absence in the view declares its own atom (`Undefined = {undefined}`) in its superdomain.

**Example:** [`examples/04-optional-nullable.sets`](../examples/04-optional-nullable.sets)

### 5.5 Positional tuples and fixed arrays

**CS-TUP-1 (Invariant).** A positional tuple `[T1, …, Tn]` maps to a tuple builder:

```sets
Pair = { (x, y) suchthat x in String and y in Integer }
```

A fixed-length array `T[n]` is an n-tuple with every component in `T`. A one-element tuple stays a one-element tuple: `(x)` is never `x`.

**Edge cases.** A named tuple whose names appear in the view (on the wire) is a record, not a tuple. A zero-length tuple is `{()}`.

**Example:** [`examples/05-tuple.sets`](../examples/05-tuple.sets)

### 5.6 Untagged unions

**CS-UNI-1 (Invariant).** `A | B` maps to `union(A, B)`.

This is correct for records too, because records are open by default (VM-REC-3). README's `Human`:

```sets
Person = { r suchthat r in Record and ("name", name) in r and name in String }
Aged = { r suchthat r in Record and ("age", age) in r and age in Real }
Human = union(Person, Aged)
```

`{("age", 123), ("name", "abc")}` is in both branches. That overlap is correct, and target adapters must keep it (CT-UNI-3). When the language's records are closed (CS-REC-4), the same `union` of closed branches excludes that value, which is also correct for that view.

**Example:** [`examples/06-union-disjoint.sets`](../examples/06-union-disjoint.sets), [`examples/07-union-overlapping.sets`](../examples/07-union-overlapping.sets)

### 5.7 Tagged unions

**CS-TAG-1 (Invariant).** A sum type whose variants are distinguishable in the value maps to `union` of its variants in source order. The shape of each variant follows the view:
- In the core view, the shape is `Variant({Tag}, Payload)` or the atom `{Tag}` (VM-VAR-1).
- When the view embeds the tag in a record (internally tagged JSON, `{"kind": "circle", "r": 1}`), the variant is a record whose tag field has a singleton domain: `("kind", kind) in r and kind in {"circle"}`.

**Rationale.** The representation follows where the tag lives in the value, not how the source declares it. Both forms are disjoint by construction, which CT-DISJ-1 can detect.

**Example:** [`examples/08-tagged-union.sets`](../examples/08-tagged-union.sets)

### 5.8 Intersections

**CS-INT-1 (Invariant).** `A & B` maps to `intersect(A, B)`. With open records this is the record that has both sets of fields: `intersect(Named, Dated)` admits records with both a `name` and an `at`.

**Edge cases.** Two closed records with different key sets intersect to `phi`, which is the correct domain for that view.

**Example:** [`examples/09-intersection-difference.sets`](../examples/09-intersection-difference.sets)

### 5.9 Top and bottom

**CS-TOP-1 (Invariant).** A type that admits every value in the view (`unknown`, `any`, `interface{}`, `object`) maps to `u` only if the view really admits every value. Otherwise it maps to the union of what the view admits. A JSON view admits `union(String, union(Real, union(Boolean, union(Null, union(Seq(u), Record)))))`; the language standard defines this as a named domain such as `JsonValue`.

**CS-TOP-2 (Invariant).** A type with no values (`never`, `!`, `Nothing`) maps to `phi`.

### 5.10 Collections

**CS-COL-1 (Invariant).** Variable-length list, array, slice → `Seq(T)`. Set → `SetOf(T)`. Map or dictionary → `Map(K, V)`. Index signature `{[k: string]: V}` → `Map(String, V)`.

**CS-COL-2 (Invariant).** A record with fixed fields plus an index signature for the rest is the intersection of the record and a constraint on its other fields. That's expressible only when the fixed fields' values are also in `V`: `intersect(Fixed, Map(String, V))`. Otherwise it's a gap (§7, G-3); the adapter emits the open record `Fixed` with a `// CS-D-GAP index-signature` comment, and reports a diagnostic.

**Example:** [`examples/10-collections.sets`](../examples/10-collections.sets)

### 5.11 Aliases, nominal types and brands

**CS-NOM-1 (Invariant). Names don't create values.** `type A = B` → `A = B`. A nominal newtype (`type Celsius float64`, `struct Meters(f64)` when serialized transparently) whose values are indistinguishable from its base in the view is also `A = B`. The adapter doesn't invent a disjoint axiom to model nominality.

**CS-NOM-2 (Invariant). Brands that exist in the value are part of it.** If the view carries the type's identity (a type tag, a class discriminator), that tag is a field or a variant tag, per CS-TAG-1.

**CS-NOM-3 (Default). Phantom brands are dropped.** A compile-time-only brand (`string & {__brand: "Id"}`) has no runtime value, so it maps to its base domain. The language standard can override this when its toolchain needs brands back. It does so with a described subdomain with no prose extent (`UserId sub String`, which GUIDE reads as the same as `String`), never with an axiom.

### 5.12 Generics

**CS-GEN-1 (Invariant).** A generic type maps to a transformer with one domain parameter per type parameter, in source order. An unbounded parameter is `{T sub u}`; a bound `T extends X` is `{T sub X}`:

```sets
Box({T sub u}): { r suchthat r in Record and ("value", value) in r and value in T }
```

**CS-GEN-2 (Invariant).** An instantiation is a call: `Box<String>` → `Box(String)`. A named instantiation (`type StrBox = Box<String>`) is `StrBox = Box(String)`.

**CS-GEN-3 (Default).** Type-level computation (mapped types, conditional types, `keyof`) is evaluated by the source adapter, and the result is emitted. The type function itself isn't emitted. If the adapter can't evaluate it, it reports a diagnostic and emits nothing for that type.

**Example:** [`examples/11-generic.sets`](../examples/11-generic.sets)

### 5.13 Recursion

**CS-RECUR-1 (Invariant).** A self-referential or mutually recursive type refers to the domain by name (VM-RECUR-1):

```sets
Tree = {
    r suchthat
    r in Record and
    ("children", children) in r and children in Seq(Tree) and
    ("label", label) in r and label in String
}
```

**Example:** [`examples/12-recursion.sets`](../examples/12-recursion.sets)

### 5.14 Refinements

**CS-REF-1 (Invariant).** Value constraints that the language's type system or schema annotations enforce (ranges, lengths, patterns, non-empty) are part of the domain. Use `difference`, `not in`, rosters or predicates when the grammar can express the constraint. Otherwise declare a described subdomain with a prose extent (VM-SUB-2):

```sets
NonEmptyString = difference(String, {""})
```

**CS-REF-2 (Invariant).** A constraint that can't be expressed either way is dropped, and the domain widens. The adapter reports a `CS-D-WIDENED` diagnostic, and adds a matching comment on the declaration.

### 5.15 Behavior

**CS-BEH-1 (Default).** Function types, methods, channels, promises, iterators and references-as-identity have no values in the data view. A field of such a type is rejected with a `CS-D-OUT-OF-VIEW` diagnostic. The field isn't silently dropped, because dropping it would change the record domain.

### 5.16 Modules

**CS-MOD-1 (Invariant).** One source module, package or namespace maps to one `.sets` file. The file name is the superdomain name (GUIDE, "`.sets` files are superdomains"). The language standard defines how module paths become file names.

**CS-MOD-2 (Invariant).** An import of type `X` from module `m` is `X in m`. Re-exports are `X in m` in the re-exporting file. Visibility (exported versus private) isn't domain information; every type referenced by an exported type is emitted.

## 6. Construct coverage

| Construct (abstract) | Status | Rule |
| -------------------- | ------ | ---- |
| Text, char, integers, floats, booleans, null, unit | mapped | CS-SCA-1 |
| Decimals, dates, UUIDs | mapped through language described subdomains | CS-SCA-2, VM-SUB-2 |
| Literal types, literal unions | mapped | CS-LIT-1, CS-LIT-2 |
| Enumerations | mapped | CS-ENUM-1 |
| Pattern / template-literal types | gap (G-4) | CS-REF-2 |
| Records, nested records | mapped | CS-REC-1, CS-REC-3 |
| Embedding / inheritance / spread | mapped (resolved to fields) | CS-REC-3 |
| Optional fields, nullable values | mapped | CS-OPT-* |
| Open and closed records | mapped | CS-REC-4 |
| Positional tuples, fixed arrays | mapped | CS-TUP-1 |
| Untagged unions | mapped | CS-UNI-1 |
| Tagged unions | mapped | CS-TAG-1 |
| Intersections | mapped | CS-INT-1 |
| Top / bottom types | mapped | CS-TOP-1, CS-TOP-2 |
| Lists, sets, maps, index signatures | mapped | CS-COL-1 |
| Fixed fields plus an index signature | mapped when fixed values fit `V`; otherwise gap (G-3) | CS-COL-2 |
| Aliases, newtypes, brands | mapped | CS-NOM-* |
| Generics, bounds | mapped | CS-GEN-1, CS-GEN-2 |
| Type-level computation | mapped by evaluation | CS-GEN-3 |
| Recursive types | mapped | CS-RECUR-1 |
| Range / length refinements | mapped through described subdomains | CS-REF-1 |
| Functions, channels, promises | out of scope | CS-BEH-1 |
| Modules, imports | mapped | CS-MOD-* |

Totals over 23 constructs: 21 mapped (one of them partly a gap), 1 gap, 1 out of scope.

## 7. Open Gaps

**G-1. Ranges and patterns.** The grammar has no ordering or regular expressions. Interim: described subdomains with prose extents (VM-SUB-2). *Non-normative proposal:* a `base` predicate family such as `LessThan(x, y)` with prose semantics, or a `Matches` predicate over a pattern string.

**G-2. Cardinality and finiteness.** `SetOf` and `Map` rely on prose for finiteness and key uniqueness (VM-SET-1, VM-MAP-1). Interim: the prose is normative.

**G-3. "Every other field".** An open record says nothing about the fields it doesn't mention. Constraining only those (an index signature whose value type excludes a fixed field's type) needs a way to name "fields whose key isn't in K". *Non-normative proposal:* a `base` transformer `OtherFields({K sub u}, {V sub u}): { r suchthat r in Record and r not in { s suchthat (k, v) in s and k not in K and v not in V } }`. It uses only existing grammar, so it can be adopted once reviewed.

**G-4. Pattern / template-literal strings.** Same as G-1. Interim: a described subdomain with a prose pattern, or widening to `String` with `CS-D-WIDENED`.

### Diagnostics

| Code | Meaning |
| ---- | ------- |
| `CS-D-WIDENED` | A constraint was dropped and the domain is wider than the source type. |
| `CS-D-GAP` | The construct hits an Open Gap; the output is the documented interim form. |
| `CS-D-OUT-OF-VIEW` | The construct has no values in the declared view. |
| `CS-D-UNEVALUATED` | Type-level computation couldn't be evaluated. |

## 8. Spec Issues

Both directions share one list, kept in [VALUE-MODEL.md §10](../VALUE-MODEL.md#10-spec-issues).
