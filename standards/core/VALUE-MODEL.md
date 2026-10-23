# Core Value Model

- **Standard version:** 0.2.0
- **Targets `.sets` spec at commit:** `f5a56cc` (`GUIDE.md`), `b52f182` (`README.md`)
- **Status:** draft
- **Used by:** [core source standard](source/STANDARD.md), [core target standard](target/STANDARD.md), and every language standard

`.sets` describes domains of values, but `GUIDE.md` doesn't yet say what a string, a record or a list *is* as a value. Without that agreement, a source adapter and a target adapter can read the same `.sets` text differently. This document fixes the values that both directions share. It is the only document a target adapter needs to understand a source adapter's output.

Rules are tagged **Invariant** (a language standard can't override them) or **Default** (a language standard can override them in its Overrides table; see [standards/README.md](../README.md)).

<!-- mtoc-start -->

* [1. Value kinds](#1-value-kinds)
* [2. Lexical conventions for roster symbols](#2-lexical-conventions-for-roster-symbols)
* [3. Rules and binders](#3-rules-and-binders)
* [4. Axioms](#4-axioms)
* [5. Ranged and described subdomains](#5-ranged-and-described-subdomains)
* [6. Structural encodings](#6-structural-encodings)
  * [6.1 Records](#61-records)
  * [6.2 Missing fields and null](#62-missing-fields-and-null)
  * [6.3 Tagged variants](#63-tagged-variants)
  * [6.4 Tuples and the empty tuple](#64-tuples-and-the-empty-tuple)
  * [6.5 Sequences, sets and maps](#65-sequences-sets-and-maps)
  * [6.6 Recursion](#66-recursion)
* [7. Equality](#7-equality)
* [8. Canonical JSON projection](#8-canonical-json-projection)
* [9. Value notation for test vectors](#9-value-notation-for-test-vectors)
* [10. Spec Issues](#10-spec-issues)

<!-- mtoc-end -->

## 1. Value kinds

**VM-KIND-1 (Invariant).** Every value in the core model has exactly one of these kinds:

| Kind | Examples | Written in `.sets` as |
| ---- | -------- | --------------------- |
| String | `"abc"`, `""` | quoted symbol |
| Number | `0`, `-7`, `2.5` | numeric symbol |
| Atom | `true`, `null`, `Circle` | bare symbol |
| Tuple | `("abc")`, `(1, 2)`, `()` | parentheses |
| Finite domain | `{("name", "abc")}`, `{1, 2}`, `{}` | braces |

Records, sets and maps aren't separate kinds. They are finite domains whose members are values (§6). Values of different kinds are never equal: `"1" ≠ 1 ≠ one ≠ (1) ≠ {1}`.

## 2. Lexical conventions for roster symbols

`GUIDE.md` defines a symbol as a run of UTF-8 characters, and says quoting is a convention that "the transpiler ecosystem" may agree on. This section is that agreement, accepted for now.

**VM-LEX-1 (Invariant). Quoted symbols are strings.** `"Succeeded"` is the String value `Succeeded`. Inside quotes, `\"` is a literal quote and `\\` is a literal backslash, in addition to GUIDE's `\,`.

**VM-LEX-2 (Invariant). Numeric symbols are numbers.** A bare symbol that matches `-?(0|[1-9][0-9]*)(\.[0-9]*[1-9])?` is the number it spells in base 10. `3` is in `Integer`; `2.5` is in `Fractional`. This is the only canonical spelling: no `+`, no leading zeros, no trailing fractional zeros, no exponent, and `-0` is written as the atom `NegativeZero` (VM-AX-3).

**VM-LEX-3 (Invariant). Other bare symbols are atoms.** An atom is an opaque value equal only to itself. It isn't a member of `String`, so `Succeeded ≠ "Succeeded"`.

**VM-LEX-4 (Invariant). The `@` escape.** When an atom's spelling would otherwise be read as something else (a domain name in scope, a numeric symbol, or a symbol starting with `"` or `@`), it's written with a leading `@`. `{@Number}` is the atom `Number`; `{Number}` is the domain `Number`. Atoms are written bare whenever that's unambiguous.

**VM-LEX-5 (Invariant). Universal domain spelling.** Standards use `U`. `u` is an alias in `base` for convenience and is never emitted.

## 3. Rules and binders

**VM-RULE-1 (Invariant). Rules follow boolean algebra.** Rules combine with the usual precedence from math, highest first:

| Level | Operators |
| ----- | --------- |
| 1 | arithmetic (VM-AX-6): `^`, then `*`, `/`, `%`, then `+`, `-` |
| 2 | comparisons and membership: `<`, `<=`, `>`, `>=`, `in`, `not in`, `sub` |
| 3 | `and` |
| 4 | `or` |

Operators at the same level group left to right, except `^`, which groups right to left as in math (`2 ^ 3 ^ 2` is `2 ^ 9`). `A or B and C` is `A or (B and C)`; `A and B or C and D` is `(A and B) or (C and D)`.

**VM-RULE-1a (Invariant). Grouping parentheses.** Inside a rule, parentheses that enclose a rule expression (anything containing `in`, `not in`, `sub`, `and`, `or` or a comparison) or an arithmetic expression group it, as in math: `A and (B or C)`, `(x + 1) % 2`. Parentheses that enclose only values and binders are tuples: `("name", name) in r`. Grouping parentheses never create a value, so this doesn't conflict with "parentheses are values" (GUIDE, "Tuples and Vectors"), which is about values.

**VM-RULE-2 (Invariant). Rule-local binders are existential.** A symbol that first appears inside a rule, and isn't the member binder, a function parameter or a domain in scope, is bound existentially over the whole rule. In

```sets
Person = { r suchthat r in Record and ("name", name) in r and name in String }
```

`name` ranges over the values for which the rule holds: `r` is a member if *some* `name` makes the rule true. Because record keys are unique (VM-MAP-1), there's at most one candidate, so "some `name`" means "the value of `r`'s `name` field".

| `r` | Why | Member? |
| --- | --- | ------- |
| `{("age", 1), ("name", "abc")}` | `name = "abc"` puts `("name", name)` in `r`, and `"abc"` is a String. | yes |
| `{("name", 5)}` | The only candidate is `name = 5`, and `5` isn't a String. | no |
| `{("age", 1)}` | No value of `name` puts `("name", name)` in `r`. This is how a required field is expressed. | no |

The binder is existential, not universal: "for all `name`, `("name", name) in r`" would need every value to be paired with `"name"`, which no finite record satisfies.

**VM-RULE-3 (Invariant). Literals in patterns.** A pattern such as `("name", name)` may mix literal values (`"name"`) and binders (`name`). A literal matches only an equal value (VM-EQ-1).

**VM-RULE-4 (Invariant). `sub` in a rule.** `s sub T` holds when every member of `s` is a member of `T`.

## 4. Axioms

`base.sets` declares three top-level axioms. By GUIDE, top-level domains are pairwise disjoint.

| Axiom | Denotes |
| ----- | ------- |
| `String` | Every finite sequence of Unicode scalar values (U+0000–U+D7FF, U+E000–U+10FFFF), including the empty string. Surrogate code points aren't string content. |
| `Integer` | Every mathematical integer, unbounded. |
| `Fractional` | Every finite real number that isn't an integer. |

**VM-AX-1 (Invariant).** `String`, `Integer` and `Fractional` are disjoint. `Real = union(Integer, Fractional)` is every finite real number.

**VM-AX-2 (Invariant). Numbers are mathematical.** A number value has no width, signedness or representation. The binary64 value `1.0` and the integer `1` are the same value and belong to `Integer`. Width is a constraint on the domain (§5), not on the value.

**VM-AX-3 (Invariant). Non-finite values are atoms.** IEEE-754 `NaN`, `+∞`, `−∞` and `−0` are the atoms `NaN`, `PositiveInfinity`, `NegativeInfinity` and `NegativeZero` (the roster `IEEESpecial`). All NaN payloads are the single atom `NaN`. A source standard whose view doesn't distinguish `−0` from `0` (JSON, for example) uses `0` and doesn't emit `NegativeZero`.

**VM-AX-4 (Invariant). Booleans and null are atoms.** `Boolean = {true, false}`, `Null = {null}`. A language whose booleans are integers maps them by its own source standard; the core values stay atoms.

**VM-AX-6 (Invariant). Numeric operators are mathematical.** Because numbers are mathematical values (VM-AX-2), operators on them have their mathematical meaning:
- `+`, `-`, `*`, `/` give exact results, with no overflow, rounding or wraparound. `1 / 3` is the exact rational one third. Division by zero has no result.
- `a ^ n` is exponentiation with an integer exponent `n`. For `n >= 0` it's defined for every number `a` (`0 ^ 0` is `1`); for `n < 0` it's `1 / a ^ (0 - n)` and needs `a` to be non-zero. A non-integer exponent has no result, so results stay exact rationals: `2 ^ -2` is `0.25`.
- `%` is defined on integers. The result has the sign of the divisor (mathematical modulo): `-7 % 3` is `2`. A non-integer operand or a zero divisor gives no result.
- `<`, `<=`, `>`, `>=` are the ordering of the reals.
- An operator applied to a value that isn't a number (string, atom, tuple, domain) has no result. Nothing is coerced: `"1" + 1` doesn't exist.
- A comparison or membership test whose operand has no result is false, so the rule doesn't hold for that value.

Equality stays `x in singleton(y)` (GUIDE, "Singletons"); `=` and `!=` aren't operators.

**VM-AX-5 (Default). New axioms.** A language standard adds a top-level axiom only for values that no `base` domain contains (for example, opaque handles in a runtime view). Values that exist in `base` (dates, UUIDs, decimals) are described subdomains of `base` domains instead, because a new axiom would be disjoint from them.

## 5. Ranged and described subdomains

**VM-SUB-0 (Invariant). Ranged domains are builders.** A domain bounded by numeric comparisons is written as a builder with VM-AX-6 operators. `base.sets` defines `Natural`, `Int8`…`Int64` and `Uint8`…`Uint64` this way, for example:

```sets
Int32 = { x suchthat x in Integer and x >= -2147483648 and x <= 2147483647 }
```

Their overlaps follow from the rules: `Uint32` is inside `Int64` because every member satisfies `Int64`'s rule.

The binary floating-point domains are builders too. A finite IEEE-754 value is `m * 2 ^ e` for an integer significand `m` and exponent `e` within the format's bounds. `m` and `e` are rule-local binders, so they're existential (VM-RULE-2):

```sets
Binary64Finite = {
    x suchthat
    m in Integer and m > -9007199254740992 and m < 9007199254740992 and
    e in Integer and e >= -1074 and e <= 971 and
    x in singleton(m * 2 ^ e)
}
```

| Domain | `abs(m) <` | `e` range |
| ------ | ---------- | --------- |
| `Binary64Finite` | `2 ^ 53` | `-1074 … 971` |
| `Binary32Finite` | `2 ^ 24` | `-149 … 104` |

`0` is a member (`m = 0`); `−0` isn't, because it's the atom `NegativeZero` (VM-AX-3). `Float64` and `Float32` add `IEEESpecial` with `union`. `Binary32Finite` is inside `Binary64Finite`, and every `Int32` is in `Binary64Finite`, both by the rules.

A top-level `X sub Y` declares a domain `X` whose members are drawn from `Y`, at most all of `Y` (GUIDE, "Difference between `sub` and `in` in this context"). Domains that the grammar still can't state are declared that way, and the prose here fixes each one's exact extent. Target adapters bind these domains by this prose. A described subdomain with no prose extent is read as GUIDE reads it: the same as `Y`.

| Domain | Exact members |
| ------ | ------------- |
| `Char` | Strings of exactly one Unicode scalar value. (Needs a string-length operator.) |

**VM-SUB-1 (Invariant).** Described subdomains overlap the way their prose says, regardless of the single `sub` parent written in `base.sets`.

**VM-SUB-2 (Default).** A language standard writes further ranged domains as builders (VM-SUB-0), and declares the remaining described subdomains in its own superdomain using the `sub` form (`Decimal128 sub Real`, `Uuid sub String`, `Rfc3339Timestamp sub String`). Each one needs a prose extent precise enough to write a membership check from.

## 6. Structural encodings

### 6.1 Records

**VM-REC-1 (Invariant). A record value is a set of tuples.** Each member is a field `(key, value)`; keys are strings and unique. `Record = Map(String, U)` in `base` is the domain of all records.

```text
{ name: "abc", age: 123 }   is the value   {("name", "abc"), ("age", 123)}
```

Field names are part of the value, because `{"name": 1}` and `{"nom": 1}` are different data on any wire. A positional tuple `(name, age)` can't keep that difference without relying on binder names, and binder names are scoped symbols that renaming must not affect.

**VM-REC-2 (Invariant). Canonical record form.** A record domain is a builder over a record binder `r`. Each required field is a pattern on `r`:

```sets
Person = { r suchthat r in Record and ("name", name) in r and name in String }
```

`r in Record` comes first, so `r` ranges over records and not over any domain that happens to contain a `("name", …)` tuple.

**VM-REC-3 (Default). Records are open.** A record domain admits fields it doesn't mention, with any value. `{("age", 123), ("name", "abc")}` is a `Person`. This is the README `Human` reading: that value is both a `Person` and an `Aged`, so `union(Person, Aged)` is exactly `Human`.

**VM-REC-4 (Invariant). Closed records say so.** A record domain that rejects unmentioned keys adds `r in OnlyKeys({k1, …, kn})`, listing every key it admits.

### 6.2 Missing fields and null

**VM-MISS-1 (Invariant). A missing field is a missing tuple.** Absence isn't a value. `{name: "abc"}` is `{("name", "abc")}`; it has no `age` tuple at all.

**VM-MISS-2 (Invariant). Optional fields.** `r in OptionalField({"age"}, T)` says: if `r` has an `age` field, its value is in `T`. A field that's neither required nor optional in an open record can hold anything.

**VM-MISS-3 (Invariant). Null is a value.** `null` is an atom. A field that may be missing, null or a value is `r in OptionalField({"k"}, union(T, Null))`, and those are three different record values.

### 6.3 Tagged variants

**VM-VAR-1 (Invariant).** A value of a variant with a payload is the 2-tuple `(Tag, payload)`, where `Tag` is an atom. `Variant({Tag}, P)` is `product2({Tag}, P)`. A variant without a payload is the bare atom `Tag`.

```sets
Shape = union(Variant({Circle}, Real), union(Variant({Square}, Real), {Origin}))
```

### 6.4 Tuples and the empty tuple

**VM-TUP-1 (Invariant).** An n-tuple `(v1, …, vn)` is ordered, nests without flattening (GUIDE, "Tuples and Vectors"), and is never equal to a record, a sequence or a domain. Fixed-length positional products (language tuples, fixed-size arrays) are tuples.

**VM-TUP-2 (Invariant). The empty tuple.** `()` is the 0-tuple: one value, equal only to itself. It is distinct from `{}`, `null` and every other value, so transpilers may use it as a unique marker. `base` uses it to end sequences (`EmptySeq = {()}`); a language standard may also use it for its unit value.

### 6.5 Sequences, sets and maps

**VM-SEQ-1 (Invariant).** A finite sequence is a nested pair ending in the empty tuple: `[]` is `()`, and `[a, b]` is `(a, (b, ()))`. `Seq(T)` is the domain of finite sequences over `T`. Target adapters bind `Seq` to native lists by meaning; they don't generate nested pairs.

**VM-SET-1 (Invariant).** A finite set of values is a finite domain. `SetOf(T)` is the domain of **finite** subdomains of `T`. The `base.sets` text says `s sub T`; the finiteness comes from this rule, because values are finite.

**VM-MAP-1 (Invariant).** A finite map is a finite domain of `(key, value)` pairs with unique keys. `Map(K, V)` is every such domain whose pairs are in `product2(K, V)`. Key uniqueness comes from this rule; the grammar can't state it yet. A record is a map with string keys.

### 6.6 Recursion

**VM-RECUR-1 (Invariant).** A domain may refer to itself, directly or through other named domains. Members are finite values. `Seq` in `base.sets` is recursive.

## 7. Equality

**VM-EQ-1 (Invariant).** Two values are equal exactly when they have the same kind and:
- strings: the same sequence of scalar values (no normalization);
- numbers: the same mathematical value;
- atoms: the same spelling;
- tuples: the same arity and pairwise equal components;
- finite domains: the same members.

`x in singleton(y)` holds exactly when `x` and `y` are equal (GUIDE, "Singletons").

## 8. Canonical JSON projection

JSON is the most common wire. This projection is the default wire for both directions. A language standard that uses a different wire format defines its own projection and keeps to the same rules: it's one-to-one on each domain, and it rejects non-members.

**VM-JSON-1 (Default).**

| Value | JSON |
| ----- | ---- |
| String | string |
| Number | number, written per VM-LEX-2 |
| `true`, `false`, `null` | `true`, `false`, `null` |
| Other atom `A` | the string `"A"` |
| `NaN`, `PositiveInfinity`, `NegativeInfinity` | not representable; see VM-JSON-2 |
| `NegativeZero` | `-0` |
| Record | object with one member per field, including fields the domain doesn't mention |
| `(Tag, payload)` variant | `{"Tag": payload}` |
| Sequence | array |
| Set | array, sorted by the canonical JSON text of its elements |
| Map with string keys | object |
| Map with other keys | array of `[key, value]` arrays, sorted by key text |
| Empty tuple `()` | `[]` |
| Other n-tuple | array of length n |

**VM-JSON-2 (Invariant). Ambiguity is a diagnostic, not a guess.** If a domain admits two different values with the same projection (for example the atom `Circle` and the string `"Circle"`, a 2-tuple and a 2-element sequence, or a variant `(Tag, p)` and a record with the single key `"Tag"`), the projection isn't one-to-one for that domain. The adapter then emits a diagnostic and either rejects the domain or uses a projection override from its language standard. The same applies to non-finite floats on plain JSON.

## 9. Value notation for test vectors

Examples and test vectors write values in `.sets` value notation (§1–§6), and add the JSON projection when it helps. For example:

```text
member:      {("name", "abc")}                   JSON {"name":"abc"}
member:      {("name", "abc"), ("nick", "a")}    JSON {"name":"abc","nick":"a"}   (open record: extra key)
non-member:  {("name", 1)}                       JSON {"name":1}
```

## 10. Spec Issues

Ambiguities and contradictions in `README.md` (`b52f182`) and `GUIDE.md` (`f5a56cc`) that the core standards had to resolve. Each lists the interpretation the core standards assume.

### Open

| # | Issue | Core assumption |
| - | ----- | --------------- |
| SI-18 | GUIDE's `base.sets` example lists only `u`, `phi` and the operations. `base` now also holds the axioms, described subdomains and structural transformers. | [`base.sets`](base.sets) is the full `base`. GUIDE's example could point to it. |
| SI-6 | GUIDE says function parameters must be domains, but `IsOne(x)` and `ToUser(x, y)` take bare parameters. | Core only uses `{X sub D}` and `{x suchthat …}` parameters. |
| SI-7 | `=`, `sub or =`, `!=` (GUIDE) and `subsetof` (README `Human`) appear but aren't defined. | Not used. `%` is now defined by VM-AX-6. |
| SI-8 | `base` shows a variadic `union({A sub u}, {B sub u}, ...)`; variadic parameters aren't defined. | Binary `union`/`intersect`, nested right-associatively (CS-CONV-4). |
| SI-12 | "`{x} in U` *can* be written as `{x} sub U`" mixes membership and subset. | In a rule, `s sub T` means every member of `s` is in `T` (VM-RULE-4; used by `SetOf`, `Map`, `OnlyKeys`). |
| SI-14 | There's no grammar for cardinality, finiteness, string length or patterns. | Prose rules (VM-SET-1, VM-MAP-1, §5) and the gaps in each standard. Ranges and floating-point domains are covered by VM-AX-6. |
| SI-15 | `//` comments are used throughout but aren't part of the grammar. | Allowed; they carry no meaning. |
| SI-16 | README says TypeScript's `Person \| Aged` is an "intersection". As a domain it's a union of open records whose members can satisfy both branches. | Terminology only. `union(Person, Aged)` with open records (VM-REC-3). |
| SI-17 | GUIDE doesn't say how symbols that appear only in a rule (`n` in `("name", n) in r and n in String`) are bound. | Existentially over the whole rule (VM-RULE-2). |

### Resolved

| # | Issue | Resolution |
| - | ----- | ---------- |
| SI-1 | The universal domain was spelled both `u` and `U`. | `U` is canonical; `u` is an alias in `base` (VM-LEX-5). |
| SI-2 | Record members were `{name, age}` in GUIDE and `(name)` in README. | Records are sets of `(key, value)` tuples: `{r suchthat ("name", n) in r and n in String}` (VM-REC-*). |
| SI-3 | There was no shared axioms superdomain. | `core` is merged into `base`. |
| SI-4 | Quoting was only a convention. | VM-LEX-1 to VM-LEX-4 are accepted for now. |
| SI-5 | `and`/`or` had no precedence. | Boolean algebra with `()` grouping (VM-RULE-1, VM-RULE-1a). Pending in GUIDE; see the notes below. |
| SI-9 | The empty tuple wasn't defined. | `()` is the 0-tuple; transpilers may use it as a unique value (VM-TUP-2). |
| SI-10 | Recursion wasn't defined. | Allowed (VM-RECUR-1). |
| SI-11 | Top-level `X sub Y` seemed undefined. | Defined in GUIDE, "Difference between `sub` and `in` in this context" (§5). |
| SI-13 | Literals inside patterns weren't defined. | Allowed (VM-RULE-3). |

### Notes for the spec author

Decisions made in review that belong in `GUIDE.md`. The standards already rely on them.

1. **Rule precedence and grouping** (VM-RULE-1, VM-RULE-1a). Rules follow boolean algebra: arithmetic, then comparisons and membership, then `and`, then `or`, each level grouping left to right. Parentheses around a rule or arithmetic expression group it; parentheses around values and binders are tuples. GUIDE's "Tuples and Vectors" section should say that grouping parentheses aren't values.
2. **Numeric operators** (VM-AX-6). `^ * / % + -` and `< <= > >=` on numbers, with exact mathematical meaning and no coercion; `^` takes integer exponents and groups right to left. This makes ranged and floating-point domains expressible (`Int32`, `Binary64Finite` as builders), which GUIDE's number discussion could show.
3. **Rule-local binders** (VM-RULE-2, SI-17). A symbol that first appears inside a rule is existential over the rule.
