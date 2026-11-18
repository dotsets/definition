# Core Source Standard (`any language → .sets`)

- **Standard version:** 0.2.0
- **Source language / version:** none. This is the base that every `<language> → .sets` standard extends.
- **Targets `.sets` spec at commit:** `f5a56cc` (`GUIDE.md`)
- **Depends on:** [Core Value Model](../VALUE-MODEL.md), [`base.sets`](../base.sets)
- **Status:** draft

Rule IDs are `CS-*`. Rules moved here from the value model keep their `VM-*` IDs, and like every value-model rule they're invariants. **Invariant** rules can't be overridden; **Default** rules can be, through the language standard's Overrides table ([standards/README.md](../../README.md)).

<!-- mtoc-start -->

* [0. Principles](#0-principles)
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
* [6. Structural encodings](#6-structural-encodings)
  * [6.1 Records](#61-records)
  * [6.2 Missing fields and null](#62-missing-fields-and-null)
  * [6.3 Tuples](#63-tuples)
  * [6.4 Sequences, sets and maps](#64-sequences-sets-and-maps)
* [7. Value notation for test vectors](#7-value-notation-for-test-vectors)
* [8. Construct coverage](#8-construct-coverage)
* [9. Open Gaps](#9-open-gaps)
* [10. Spec Issues](#10-spec-issues)

<!-- mtoc-end -->

## 0. Principles

**VM-PRES-1. Preserve as much of the domain as the source states.** `.sets` exists to preserve how data is structured, how it's spread across values, and how it's viewed. A `.sets` description keeps every distinction its source makes, even one a likely target can't use: a distinction dropped on the way in can't be recovered on the way out.

How much of the domain survives the IR is the measure of a translation's clarity and correctness. When two encodings are both valid, the one that keeps more of the domain is the better one.

**VM-PRES-2. Preserve the structure of the data.** How data is structured doesn't change between languages. A record is fields and their values, a sequence is ordered elements, and a set is its members, whatever a language calls the type that holds them. A source type maps to the structure of its values, not to the construct the source uses. A struct, an object type and a class instance all hold fields and values, so each maps to a record (VM-REC-1). A class or library type that wraps another structure holds that structure's values and adds none of its own. A JavaScript `Set` is an object whose keys are its members and whose values are empty, so it's a record (VM-SET-1):

```text
new Set(["a", "b"])   is the value   {("a", ()), ("b", ())}
```

The empty value may be written `()` or `Phi`; in the core context they're treated as the same (VM-AX-5).

**Collisions.** Two source values may become one `.sets` value only when they have the same structure and the same contents: the same fields, elements or members, and scalars the view doesn't tell apart (VM-PRES-1). Whether two values pass is a matter of their structure, never of what the source calls them. When they pass, they are one `.sets` value, even if the source gives them different types. When they don't, a mapping that makes them one value is a defect.

**Names.** Two source types with the same structure describe the same domain, but each keeps its own name as a declaration. The name records that the source treated them as different types:

```sets
UserId = String
Email = String
```

TypeScript's `type UserId = string` and `type Email = string` are both strings. `UserId` and `Email` are equal domains, and a value of one is a value of the other. The names are still two declarations, so a target can give each its own type. The name never enters the value: tagging the values (`Tagged`, VM-PRES-3) would make them different data, which the source doesn't have.

## 1. Scope

**CS-SCOPE-1 (Invariant). Every language standard declares its view.** The view is the set of values the standard's domains describe. The same type can denote different domains in different views: a TypeScript `number` is `Binary64Finite` as JSON on the wire, and `Float64` in the running program. The standard declares its view in two parts:
1. A **View** field in its header, stating what is observed (a wire format and its codec, such as `encoding/json` or `JSON.parse`, or the runtime) and which distinctions the view keeps or loses where they matter to the mapping (for example, whether it tells `1` from `1.0`, or a missing field from `null`).
2. A **view domain** that holds every value of the view. When the view is a wire format, the view domain belongs to that wire's superdomain, so every language with the same wire shares it (`JsonValue in json`). Otherwise, the standard declares it in `<language>.sets`.

Every domain the standard emits is a subdomain of the view domain, and a type that admits every value in the view maps to it (CS-TOP-1). A test vector outside the view domain is a defect in the standard.

Every emitted file imports its view domain. That's how a target standard learns the wire to generate codecs for, without learning the source language.

For example:

```text
View: TypeScript values as JSON.stringify writes them and JSON.parse reads them.
  Domain: JsonValue (json.sets, the json superdomain)
  Keeps:  a missing key vs "k": null.
  Loses:  1 vs 1.0 (one number type); undefined vs missing (undefined fields
          aren't written); NaN and infinities (written as null).
```

Under this view, `number` maps to `Binary64Finite`, `x?: T` to `OptionalField("x", T)`, `x: T | null` to `union(T, {Nothing})`, and `unknown` to `JsonValue`.

Every mapping in the standard is relative to this view. A view is always a view of values: identity and cycles aren't values (GUIDE, pending), so even a runtime view narrows to values.

**CS-SCOPE-2 (Invariant). The core view.** The core view is the [Core Value Model](../VALUE-MODEL.md) itself: every core value, observed through no particular wire. Its view domain is `U`. The core examples use this view. A language standard may declare it as its own view (`View: core`), or declare another view under CS-SCOPE-1.

Runtime-only types aren't values of any view. They map to their runtime concept (VM-PRES-3, [RUNTIME.md](../RUNTIME.md)) whatever the view, and the view-domain requirement of CS-SCOPE-1 doesn't apply to them. Identity (pointers as addresses) and cycles aren't values either (GUIDE, pending).

## 2. Conventions

**CS-CONV-1 (Invariant). Imports.** An emitted file imports every symbol it uses and its view domain (CS-SCOPE-1), and nothing else. A core-view file imports `U in base`; a JSON-view file imports `JsonValue in json`, whether or not a declaration uses it. Every domain a file references is declared in it, written inline, or imported; no symbol is built in, `U` included (GUIDE, "Domain Validation").

**CS-CONV-3 (Default). Domain names.** A named source type becomes a domain of the same name. A character that can't appear bare in a symbol, such as a separator (`,`, `{`, `}`), is escaped (GUIDE, pending); the escaped symbol is still that domain's name (GUIDE, pending). When the name is already taken:
- If the source type has the same structure as the domain in scope, it is that domain and nothing new is declared (VM-PRES-2). JavaScript's `String` is `String` from `base`.
- Otherwise the language standard renames it. That covers a name already in scope (a `base` name, an imported name such as `Async` from `runtime`, an earlier declaration) and two source types in one module that would get the same name.

**CS-CONV-5 (Invariant). Binder names.** The record binder is `r`. Each field's value binder is the key's spelling when that's a valid symbol and doesn't collide with `r` or a domain in scope; otherwise it's `_` followed by the key spelling with reserved characters escaped. Binders in other builders are `x`, `y`, `z`, then `x1`, `x2`, …. Binders are scoped symbols, so this rule only exists to make outputs comparable.

**CS-CONV-6 (Invariant). Rule shape.** Rules follow boolean algebra (GUIDE). Adapters write rules inline and add grouping parentheses only where precedence alone would read differently: `(A or B) and C` needs them; `A and B or C and D` doesn't. Adapters don't add redundant parentheses, so output stays comparable.

## 3. Axioms

**CS-AX-1 (Invariant). Reuse `base`; new axioms only for new values.** Scalars map to `base` axioms and described subdomains whenever the prose meaning matches (VALUE-MODEL §4; CS-STR-1; number axioms, pending). Values that exist in `base` (dates, UUIDs, decimals) are described subdomains of `base` domains, because a new axiom would be disjoint from them (GUIDE, "Top domain axiom"). A language standard declares a new top-level axiom only for values no `base` domain contains. Runtime-only types aren't values and never become axioms (VM-PRES-3).

**CS-AX-2 (Invariant). Disjointness check.** Before declaring an axiom, the language standard shows that no value of the new axiom is a value of an existing axiom. If two of the language's types share values (`int` and `float` both hold `1`), they're described subdomains of one axiom, or unions of axioms, never two axioms.

**CS-AX-3 (Default). Language superdomain.** Language-specific domains go in `<language>.sets`. That file imports from `base` and redefines nothing it provides.

**CS-STR-1 (Invariant). String.** The `base` axiom `String` is every finite sequence of Unicode scalar values (U+0000–U+D7FF, U+E000–U+10FFFF), including the empty string. Surrogate code points aren't string content.

**CS-STR-2 (Invariant). Quoted atoms are strings.** An atom (GUIDE, pending) that starts and ends with an unescaped `"` encodes a member of `String`: the characters between the quotes, with escapes decoded (GUIDE, pending). `"Succeeded"` encodes the string `Succeeded`, and `""` the empty string. No other atom is a `String`, so `Succeeded` isn't one.

Inside a symbol, every unescaped `"` must be closed by another one:
- `"a\"b"` is the string `a"b`.
- `a\"b` is the atom `a"b`, not a string.
- `a"b` is an error: its `"` is never closed.

Separators are escaped the same way, inside quotes or not, because the grammar layer doesn't know about quotes (VM-LEX-0):
- `{"a\, b"}` is one string, `a, b`.
- `{"a, b"}` is an error: the grammar layer splits it into `"a` and ` b"`, and neither has a closing quote.

**CS-STR-3 (Invariant). Char.** The `base` described subdomain `Char sub String` holds the strings of exactly one Unicode scalar value. (It becomes a builder once a string-length operator exists, VM-SUB-1.)

**CS-AX-4 (Default). Described subdomains.** A language standard writes further ranged domains as builders (VM-SUB-0), and declares the remaining described subdomains in its superdomain (CS-AX-3) with the `sub` form (`Uuid sub String`, `Rfc3339Timestamp sub String`). Each one needs a prose extent precise enough to write a membership check from. The language standard owns them the way `base` owns its own (VM-SUB-1): it states how they overlap, and prose that contradicts a subdomain's `sub` parent is a defect in that standard. How it describes a domain (prose, builder or structured value, VM-SUB-3) is its own choice.

**VM-SUB-3. Structured values aren't subdomains.** A value whose representation carries information the number alone doesn't is a structured value, not a subdomain of `Real`. A decimal that keeps its scale is the usual case: decimal128 stores `1.50` as coefficient `150` and exponent `-2`, and `1.5` as `15` and `-1`. They're equal numbers but different values (`1.50 * 2` is `3.00`). Since `1.50` and `1.5` are the same `Fractional` (VM-LEX-2), `Decimal128 sub Real` would lose the scale (VM-PRES-1). Keeping it takes a tuple:

```sets
Decimal128 = {
    (c, q) suchthat
    c in Integer and c > 0 - 10 ^ 34 and c < 10 ^ 34 and
    q in Integer and q >= -6176 and q <= 6111
}
```

Here `1.50` is `(150, -2)` and `1.5` is `(15, -1)`, two members. Special values (NaN, infinities, negative zero) are added with `union`, as `Float64` does.

## 4. Decisions

| Question | Decision | Rule |
| -------- | -------- | ---- |
| Record encoding | A set of `(key, value)` tuples: `{ r suchthat r in Record and ("name", name) in r and name in String }`. | VM-REC-*, CS-REC-1 |
| Field identity | Keys are values inside each tuple: any value, usually strings in a wire view. They aren't binder names. | VM-REC-1 |
| Optional / missing / null | A missing field is a missing tuple. Optional fields use `OptionalField`. Null is the value `Nothing`; a nullable value is `union(T, {Nothing})`. | VM-MISS-*, VM-AX-5, CS-OPT-* |
| Open vs closed | Records are open by default. Closed records add `r in OnlyKeys({…})`. | VM-REC-3, VM-REC-4, CS-REC-4 |
| Nominal vs structural | Structural. A name is a label for a domain. Brands that exist in the value are tags. | CS-NOM-* |
| Generics | Transformers with domain parameters; a bound `T extends X` becomes `{T sub X}`. | CS-GEN-* |
| Recursion | Self-reference by name. | GUIDE (pending), CS-RECUR-1 |
| Variable-length collections | `Seq`, `SetOf` and `Map` from `base`. | CS-COL-* |
| Modules | One superdomain file per module; imports become `X in module`. | CS-MOD-* |
| Wrapper types | Mapped to the structure they wrap. | VM-PRES-2, CS-WRAP-1 |

**Why open records by default.** With open records, README's `Human` is simply `union(Person, Aged)`: `{name: "abc", age: 123}` is a `Person` (it has a `name`) and an `Aged` (it has an `age`), so it's in the union. The visually similar mapping is also the correct one, and no merged rewrite is needed. Languages whose view rejects unknown keys add `OnlyKeys`, and then the same `union` means "exactly one shape".

## 5. Mapping rules

Each rule gives the abstract source construct, the canonical `.sets`, edge cases, and an example in [`../examples/`](../examples/). Those examples hold the canonical domain with member and non-member vectors. Language standards add `<lang>` ↔ `.sets` pairs on top of them.

### 5.1 Scalars

**CS-SCA-1 (Invariant).** A scalar type maps to the narrowest `base` domain whose prose extent equals the type's extent in the view. The mapping keeps every distinction the view makes (VM-PRES-1): a view that tells integers from fractions keeps `Integer` and `Fractional` apart (number axioms, pending). Which form a number takes is fixed by the view: a binary64 `1.0` is the `Fractional` `1.0` where the view keeps the distinction, and the `Integer` `1` where it doesn't (UP-1).

| Abstract construct | Canonical `.sets` |
| ------------------ | ----------------- |
| Unicode text | `String` |
| Single character (a Unicode scalar value) | `Char` |
| Arbitrary-precision integer | `Integer` |
| N-bit signed / unsigned integer | `IntN` / `UintN` |
| IEEE binary64 including NaN and infinities | `Float64` |
| IEEE binary64 without them (for example JSON) | `Binary64Finite` |
| Boolean | `Boolean` |
| Null / nil / None | `{Nothing}` (VM-AX-5) |
| Unit `()` | `{()}` |
| Byte string | `Seq(Uint8)` |

`Boolean` is the roster `{true, false}` in `base`: two atoms, like any enum. They're data and hold no truth (VM-AX-4). A boolean type maps by its values in the view: two atom-like values are `Boolean`; C's `_Bool` holds the integers `0` and `1`, so it's `{0, 1}`.

**CS-SCA-2 (Invariant). Exact or wider, never narrower.** If no `base` domain matches exactly, the language standard declares a described subdomain (CS-AX-4). It doesn't silently map to a narrower or wider domain. Widening loses a distinction the source makes (VM-PRES-1), so a wider mapping is allowed only with a diagnostic comment on the declaration, for example `Decimal = Real // CS-D-WIDENED: scale and precision not encoded`.

**CS-SCA-3 (Default). IEEE special values.** A view's NaN maps to `NotANumber`, folding every payload into it. A language whose view can observe payloads keeps them with its own domain in `<language>.sets` (VM-PRES-1). A view with a negative zero maps it to `NegativeZero`; a view without one never produces it.

**Edge cases.** Code units versus scalar values (UTF-16 strings that can contain lone surrogates aren't `String` in the core sense; the language standard declares its own described domain). An integer type the view serializes as a JSON string is a `String` subdomain in that view, not an `Integer`.

**Example:** [`examples/01-scalars.sets`](../examples/01-scalars.sets)

### 5.2 Literals and enumerations

**CS-LIT-1 (Invariant).** A literal type is a singleton roster: `"on"` → `{"on"}`, `3` → `{3}`, `true` → `{true}`.

**CS-LIT-2 (Invariant).** A union of literal types is one roster in source order: `"a" | "b" | 3` → `{"a", "b", 3}`.

**CS-ENUM-1 (Default).** An enumeration maps by its value in the view. If members carry string values, they're quoted strings: `{"Succeeded", "Failed"}`. If they carry numbers, they're numeric symbols. If the members are values in their own right with no underlying string or number (Rust unit variants, sealed singletons), they're unquoted atoms: `{Succeeded, Failed}`.

**CS-NUM-3 (Default).** A numeric literal is written in a notation that spells its value exactly, in the same form (VM-LEX-2). Every binary64 value has an exact hexadecimal float, and an exact decimal.

**Edge cases.** A literal that spells a domain name in scope becomes `"Number"` as a string or `\Number` as an atom (GUIDE, pending). The empty roster is `Phi`, not `{}` (CS-TOP-2).

**Rationale.** Writing enumerations as quoted strings by default would turn values that aren't strings into strings, a distinction the language keeps (VM-PRES-1).

**Example:** [`examples/02-rosters.sets`](../examples/02-rosters.sets)

### 5.3 Records

**CS-REC-1 (Invariant).** A record, struct, object or class shape maps to a record builder (VM-REC-2). Other text can describe the same domain, but adapters emit this form so every reader recognizes it. `r in Record` is required, and adapters write it first; then one clause per field in key order, using binders per CS-CONV-5. The record form keeps every field the view has (VM-PRES-1):

```sets
User = {
    r suchthat
    r in Record and
    ("age", age) in r and age in Int64 and
    ("name", name) in r and name in String
}
```

A required field is the pair `("k", k) in r and k in T`. An optional field is `r in OptionalField("k", T)` in the same key-ordered position.

**CS-REC-2 (Default). Keys are wire keys.** A field's key is its name in the view. In a JSON view that's the serialized key after renames (`json:"user_name"` → `"user_name"`), not the source identifier. Language-specific spellings never appear in `.sets` except as key values.

**CS-REC-3 (Invariant). Nested records stay nested.** A field whose type is another record refers to that record's domain: `("address", address) in r and address in Address`. Embedding, spread and inheritance are resolved by the language standard into the fields the view actually has. If the view flattens embedded fields (Go embedded structs in JSON), the record lists the flattened keys, because that's the value. Otherwise it keeps the nesting.

**CS-REC-4 (Default). Openness follows the view.** Records are open (VM-REC-3). If the view rejects unknown keys (`additionalProperties: false`, `#[serde(deny_unknown_fields)]`, a strict decoder), the record ends with `r in OnlyKeys({k1, …, kn})` listing every key it declares, in key order.

**Edge cases.** An empty open record `{}` → `Record`. An empty closed record → `{ r suchthat r in Record and r in OnlyKeys(Phi) }`, whose only member is `{}`. Duplicate keys can't occur (VM-MAP-1). Default values aren't domain information; a field with a default is required or optional according to the view.

**Example:** [`examples/03-record.sets`](../examples/03-record.sets)

### 5.4 Optional and nullable

**CS-OPT-1 (Invariant). Optional field.** A field that can be missing in the view uses `OptionalField`:

```sets
r in OptionalField("age", Integer)
```

**CS-OPT-2 (Invariant). Nullable value.** A value that can be null uses `union(T, {Nothing})`, wherever it appears (VM-AX-5).

**CS-OPT-3 (Invariant). Keep the states apart.** A field that can be missing, null or a value is `r in OptionalField("k", union(T, {Nothing}))`: a missing tuple, the tuple `("k", Nothing)`, or a value (VM-MISS-1). Only merge missing and null when the view itself can't tell them apart, and say so in the language standard. Example: Go `*T` with `omitempty` decodes both a missing key and `null` to `nil`.

**CS-OPT-4 (Default). Language "undefined".** A language with an `undefined` value that differs from absence in the view declares its own atom (`Undefined = {undefined}`) in its superdomain (VM-AX-5, VM-MISS-1).

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

**CS-TAG-1 (Invariant).** A sum type whose variants are distinguishable in the value maps to `union` of its variants in source order. Core has no variant encoding; each variant is written as the structure the view gives it: a record, a tuple or an atom (VM-PRES-2):
- When the view embeds the tag in a record (internally tagged JSON, `{"kind": "circle", "r": 1}`), the variant is a record whose tag field has a singleton domain: `("kind", kind) in r and kind in {"circle"}`.
- For other shapes (externally or adjacently tagged, a bare tag without payload), the language standard defines the encoding.

**Rationale.** The representation follows where the tag lives in the value, not how the source declares it. Both forms are disjoint by construction, which CT-DISJ-1 can detect.

**Example:** [`examples/08-tagged-union.sets`](../examples/08-tagged-union.sets)

### 5.8 Intersections

**CS-INT-1 (Invariant).** `A & B` maps to `intersect(A, B)`. With open records this is the record that has both sets of fields: `intersect(Named, Dated)` admits records with both a `name` and an `at`.

**Edge cases.** Two closed records with different key sets intersect to `Phi`, which is the correct domain for that view.

**Example:** [`examples/09-intersection-difference.sets`](../examples/09-intersection-difference.sets)

### 5.9 Top and bottom

**CS-TOP-1 (Invariant).** A type that admits every value in the view (`unknown`, `any`, `interface{}`, `object`) maps to `U` only if the view really admits every value. Otherwise it maps to the union of what the view admits. A JSON view admits `JsonValue` ([`json.sets`](../json.sets)).

**CS-TOP-2 (Invariant).** A type with no values (`never`, `!`, Scala's `Nothing`) maps to `Phi`. It isn't the core value `Nothing` (VM-AX-5).

### 5.10 Collections

**CS-COL-1 (Invariant).** Variable-length list, array, slice → `Seq(T)`. Set → `SetOf(T)`, or a record of keys with `()` values when the language standard reads its set type that way (VM-SET-1, VM-PRES-2). Map or dictionary → `Map(K, V)`. Index signature `{[k: string]: V}` → `Map(String, V)`. A language's set and map types are finite, and the language standard narrows `SetOf(T)` and `Map(K, V)` to that.

**CS-COL-2 (Invariant).** A record with fixed fields plus an index signature for the rest is the intersection of the record and a constraint on its other fields. That's expressible only when the fixed fields' values are also in `V`: `intersect(Fixed, Map(String, V))`. Otherwise it's a gap (§9, G-3); the adapter emits the open record `Fixed` with a `// CS-D-GAP index-signature` comment, and reports a diagnostic.

**Example:** [`examples/10-collections.sets`](../examples/10-collections.sets)

### 5.11 Aliases, nominal types and brands

**CS-NOM-1 (Invariant). Names don't create values.** `type A = B` → `A = B`. A nominal newtype (`type Celsius float64`, `struct Meters(f64)` when serialized transparently) whose values are indistinguishable from its base in the view is also `A = B`. The adapter doesn't invent a disjoint axiom to model nominality. Two source types with the same structure keep their own names, never merged into one declaration: `type UserId = string` and `type Email = string` → `UserId = String` and `Email = String` (VM-PRES-2).

**CS-NOM-2 (Invariant). Brands that exist in the value are part of it.** If the view carries the type's identity (a type tag, a class discriminator), that tag is a field or a variant tag, per CS-TAG-1.

**CS-WRAP-1 (Invariant). Wrappers map to the structure they wrap.** A class or library type that wraps another structure (a set over an object, a date over a number or string, a boxed big integer) maps to that structure in the view (VM-PRES-2). Core has no rule for any particular wrapper; the language standard names its wrappers and the structure each one holds.

**CS-NOM-3 (Default). Phantom brands are dropped.** A compile-time-only brand (`string & {__brand: "Id"}`) has no runtime value, so it maps to its base domain. The language standard can override this when its toolchain needs brands back. It does so with a described subdomain with no prose extent (`UserId sub String`, which GUIDE reads as the same as `String`), never with an axiom.

### 5.12 Generics

**CS-GEN-1 (Invariant).** A generic type maps to a transformer with one domain parameter per type parameter, in source order. An unbounded parameter is `{T sub U}`; a bound `T extends X` is `{T sub X}`:

```sets
Box({T sub U}): { r suchthat r in Record and ("value", value) in r and value in T }
```

**CS-GEN-2 (Invariant).** An instantiation is a call: `Box<String>` → `Box(String)`. A named instantiation (`type StrBox = Box<String>`) is `StrBox = Box(String)`.

**CS-GEN-3 (Default).** Type-level computation (mapped types, conditional types, `keyof`) is evaluated by the source adapter, and the result is emitted. The type function itself isn't emitted. If the adapter can't evaluate it, it reports a diagnostic and emits nothing for that type.

**Example:** [`examples/11-generic.sets`](../examples/11-generic.sets)

### 5.13 Recursion

**CS-RECUR-1 (Invariant).** A self-referential or mutually recursive type refers to the domain by name (GUIDE, pending):

```sets
Tree = {
    r suchthat
    r in Record and
    ("children", children) in r and children in Seq(Tree) and
    ("label", label) in r and label in String
}
```

Cyclic data isn't a value (GUIDE, pending); a language standard that keeps it encodes the references as data, such as ids or paths.

**Example:** [`examples/12-recursion.sets`](../examples/12-recursion.sets)

### 5.14 Refinements

**CS-REF-1 (Invariant).** Value constraints that the language's type system or schema annotations enforce (ranges, lengths, patterns, non-empty) are part of the domain. Use numeric operators (VM-LEX-2), `difference`, `not in`, rosters or predicates when the grammar can express the constraint. Otherwise declare a described subdomain with a prose extent (CS-AX-4):

```sets
Percent = { x suchthat x in Real and x >= 0 and x <= 100 }
EvenInt = { x suchthat x in Integer and x % 2 in {0} }
NonEmptyString = difference(String, {""})
```

**CS-REF-2 (Invariant).** A constraint that can't be expressed either way is dropped, and the domain widens. The adapter reports a `CS-D-WIDENED` diagnostic, and adds a matching comment on the declaration.

### 5.15 Behavior

**CS-BEH-1 (Default).** Function types, methods, channels, promises, iterators and references-as-identity have no values in the data view. A field of such a type is rejected with a `CS-D-OUT-OF-VIEW` diagnostic. The field isn't silently dropped, because dropping it would change the record domain.

### 5.16 Modules

**CS-MOD-1 (Invariant).** One source module, package or namespace maps to one `.sets` file. The file name is the superdomain name (GUIDE, "`.sets` files are superdomains"). The language standard defines how module paths become file names.

**CS-MOD-2 (Invariant).** An import of type `X` from module `m` is `X in m`. Re-exports are `X in m` in the re-exporting file. Visibility (exported versus private) isn't domain information; every type referenced by an exported type is emitted.

## 6. Structural encodings

### 6.1 Records

**VM-REC-1. A record value is a set of tuples.** Each member is a field `(key, value)`. A key can be any value, and keys are unique within a record. The record with no fields is `{}`, a set with no tuples. It's the empty domain `Phi`, not a missing value, which is `Nothing` (VM-AX-5). The domain of all records in `base` is:

```sets
Record = Map(U, U)
```

```text
{ name: "abc", age: 123 }   is the value   {("name", "abc"), ("age", 123)}
{ 1: "abc" }                is the value   {(1, "abc")}
{ }                         is the value   {}
```

`{}` is also the empty set, the empty map and, in context, `[]` (VM-AX-5). With nothing inside, there's no field, element or member to tell them apart, so the collision is allowed (VM-PRES-2).

Field names are part of the value, because `{"name": 1}` and `{"nom": 1}` are different data on wires that name fields. A positional tuple `(name, age)` can't keep that difference without relying on binder names, and binder names are scoped symbols that renaming must not affect.

**VM-REC-2. Record builders.** A record domain is written as a builder over a record binder `r`, with each required field as a pattern on `r`:

```sets
Person = { r suchthat r in Record and ("name", name) in r and name in String }
```

`r in Record` is what makes the members records. Without it, `{("name", "a"), ("name", 1)}` would be a `Person` even though its keys repeat, so it isn't a record (VM-REC-1). So would `{("name", "a"), 5}`.

Other text can describe the same domain, and it's still the same domain: a domain is its members, not its spelling.

**VM-REC-3. Records are open.** A record domain admits fields it doesn't mention, with any value. This follows from the builder: a rule constrains only the fields it names, so the same `.sets` text means the same domain to every adapter. A closed source type (a Go struct, a sealed class) says so with `OnlyKeys` (VM-REC-4).

```sets
Person = { r suchthat r in Record and ("name", name) in r and name in String }
Aged = { r suchthat r in Record and ("age", age) in r and age in Integer }
Human = union(Person, Aged)
```

`{("age", 123), ("name", "abc")}` is a `Person`. This is the README `Human` reading: that value is both a `Person` and an `Aged`, so it's a `Human`.

**VM-REC-4. Closed records say so.** A record domain that rejects unmentioned keys adds `r in OnlyKeys({k1, …, kn})`, listing every key it admits: its required keys and its optional keys (VM-MISS-2).

```sets
ClosedPerson = { r suchthat r in Record and ("name", name) in r and name in String and r in OptionalField("age", Integer) and r in OnlyKeys({"name", "age"}) }
```

`{("name", "abc")}` and `{("name", "abc"), ("age", 1)}` are `ClosedPerson`s; `{("name", "abc"), ("id", 1)}` isn't, because `"id"` isn't listed.

### 6.2 Missing fields and null

**VM-MISS-1. A missing field is a missing tuple.** Absence isn't a value: a missing field has no tuple at all. A field that's present with nothing in it is a tuple whose value is `Nothing` (VM-AX-5). These are different records:

```text
{name: "abc"}               is   {("name", "abc")}
{name: "abc", age: null}    is   {("name", "abc"), ("age", Nothing)}
```

Kinds of absence in between, such as JavaScript's `undefined`, are ordinary values; `.sets` gives them no special meaning.

**VM-MISS-2. Optional fields.** `r in OptionalFields(K, V)` says: every field of `r` whose key is in `K`, if present, has a value in `V`. `K` is a domain of keys, so one call can make several fields optional with the same value domain. `OptionalField(k, V)` is the single-key helper, `OptionalFields({k}, V)`. In `base`:

```sets
OptionalFields({K sub U}, {V sub U}): {
    r suchthat
    r in Record and
    r not in { s suchthat (k, v) in s and k in K and v not in V }
}

OptionalField({k suchthat k in U}, {V sub U}): OptionalFields({k}, V)
```

```sets
Person = { r suchthat r in Record and ("name", name) in r and name in String and r in OptionalField("age", Integer) }
```

`{("name", "a")}` and `{("name", "a"), ("age", 1)}` are `Person`s. `{("name", "a"), ("age", "x")}` isn't, and nor is `{("name", "a"), ("age", Nothing)}`, because `Nothing` isn't in `Integer`. A field that may also be null uses `union(Integer, {Nothing})` (VM-AX-5).

A field that's neither required nor optional in an open record can hold anything (VM-REC-3).

### 6.3 Tuples

**VM-TUP-1. Tuples are ordered.** An n-tuple `(v1, …, vn)` is ordered and nests without flattening (GUIDE, "Tuples and Vectors"). A tuple with at least one component is its own kind of value (VM-KIND-1), never a record or another domain; the empty tuple `()` is treated as `{}` in the core context (VM-AX-5). Sequences are built from tuples (VM-SEQ-1), so a nested pair ending in `()`, such as `(a, (b, ()))`, is also a sequence.

### 6.4 Sequences, sets and maps

**VM-SEQ-1. Sequences are nested pairs.** A finite sequence is a nested pair ending in the empty tuple: `[]` is `()`, and `[a, b]` is `(a, (b, ()))`. `Seq(T)` is the domain of finite sequences over `T`. Because `()` is treated as `{}` (VM-AX-5), `[]` collides with the empty record, set and map; that collision is allowed (VM-REC-1).

**VM-SET-1. Sets are domains.** A set of values is a domain, so `{"a", "b"}` is the set of `"a"` and `"b"`. `SetOf(T)` is the domain of all subdomains of `T` (`s sub T`).

**VM-MAP-1. Maps are domains of pairs.** A map is a domain of `(key, value)` pairs with unique keys. `Map(K, V)` is every such domain whose pairs are in `product2(K, V)`; `base` states uniqueness in its builder. A record is a map (VM-REC-1).

## 7. Value notation for test vectors

Examples and test vectors write values in `.sets` value notation (§1–§6). For example:

```text
member:      {("name", "abc")}
member:      {("name", "abc"), ("nick", "a")}    (open record: extra key)
non-member:  {("name", 1)}
```

## 8. Construct coverage

| Construct (abstract) | Status | Rule |
| -------------------- | ------ | ---- |
| Text, char, integers, floats, booleans, null, unit | mapped | CS-SCA-1 |
| Decimals, dates, UUIDs | mapped through language described subdomains | CS-SCA-2, CS-AX-4 |
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
| Wrapper types | mapped to the wrapped structure | CS-WRAP-1 |
| Generics, bounds | mapped | CS-GEN-1, CS-GEN-2 |
| Type-level computation | mapped by evaluation | CS-GEN-3 |
| Recursive types | mapped | CS-RECUR-1 |
| Range refinements | mapped through numeric operators | CS-REF-1 |
| Length refinements | mapped through described subdomains | CS-REF-1 |
| Functions, channels, promises | out of scope | CS-BEH-1 |
| Modules, imports | mapped | CS-MOD-* |

Totals over 25 constructs: 23 mapped (one of them partly a gap), 1 gap, 1 out of scope.

## 9. Open Gaps

**G-1. Lengths and patterns.** Numeric ranges are covered by VM-LEX-2, but there are no string or collection length operators and no regular expressions. Interim: described subdomains with prose extents (CS-AX-4). *Non-normative proposal:* a `Length` transformer on strings and sequences, and a `Matches` predicate over a pattern string.

**G-2. Cardinality and finiteness.** The grammar can't say that a domain is finite. `SetOf` and `Map` admit infinite domains (VM-SET-1, VM-MAP-1), and `Map` states key uniqueness in its builder. Interim: a language standard narrows its finite collections in prose.

**G-3. "Every other field".** An open record says nothing about the fields it doesn't mention. Constraining only those (an index signature whose value type excludes a fixed field's type) needs a way to name "fields whose key isn't in K". *Non-normative proposal:* a `base` transformer `OtherFields({K sub U}, {V sub U}): { r suchthat r in Record and r not in { s suchthat (k, v) in s and k not in K and v not in V } }`. It uses only existing grammar, so it can be adopted once reviewed.

**G-4. Pattern / template-literal strings.** Same as G-1. Interim: a described subdomain with a prose pattern, or widening to `String` with `CS-D-WIDENED`.

### Diagnostics

| Code | Meaning |
| ---- | ------- |
| `CS-D-WIDENED` | A constraint was dropped and the domain is wider than the source type. |
| `CS-D-GAP` | The construct hits an Open Gap; the output is the documented interim form. |
| `CS-D-OUT-OF-VIEW` | The construct has no values in the declared view. |
| `CS-D-UNEVALUATED` | Type-level computation couldn't be evaluated. |

## 10. Spec Issues

Both directions share one list, kept in [VALUE-MODEL.md §5](../VALUE-MODEL.md#5-spec-issues).
