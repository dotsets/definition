# Core Value Model

- **Standard version:** 0.2.0
- **Targets `.sets` spec at commit:** `f5a56cc` (`GUIDE.md`), `b52f182` (`README.md`)
- **Status:** draft
- **Used by:** [core source standard](source/STANDARD.md), [core target standard](target/STANDARD.md), and every language standard

`.sets` describes domains of values, but `GUIDE.md` doesn't yet say what a string, a record or a list *is* as a value. Without that agreement, a source adapter and a target adapter can read the same `.sets` text differently. This document fixes the values that both directions share. With [JSON-MODEL.md](JSON-MODEL.md) for the default wire, it is all a target adapter needs to understand a source adapter's output.

Rules are tagged **Invariant** (a language standard can't override them) or **Default** (a language standard can override them in its Overrides table; see [standards/README.md](../README.md)).

<!-- mtoc-start -->

* [0. Principle: preserve the domain](#0-principle-preserve-the-domain)
* [1. Value kinds](#1-value-kinds)
* [2. Lexical conventions for roster symbols](#2-lexical-conventions-for-roster-symbols)
* [3. Rules and binders](#3-rules-and-binders)
* [4. Axioms](#4-axioms)
* [5. Ranged and described subdomains](#5-ranged-and-described-subdomains)
* [6. Structural encodings](#6-structural-encodings)
  * [6.1 Records](#61-records)
  * [6.2 Missing fields and null](#62-missing-fields-and-null)
  * [6.3 Tuples](#63-tuples)
  * [6.4 Sequences, sets and maps](#64-sequences-sets-and-maps)
  * [6.5 Recursion](#65-recursion)
* [7. Equality](#7-equality)
* [8. Wire projections](#8-wire-projections)
* [9. Value notation for test vectors](#9-value-notation-for-test-vectors)
* [10. Spec Issues](#10-spec-issues)

<!-- mtoc-end -->

## 0. Principle: preserve the domain

**VM-PRES-1 (Invariant). Preserve as much of the domain as the source states.** `.sets` exists to preserve how data is structured, how it's spread across values, and how it's viewed. A source adapter keeps every distinction its source language makes, even when a likely target can't use it, so that target standards have as much domain information as possible to generate their code from. A target standard may merge or drop distinctions it can't represent (with the exactness rules of its standard), but it can't recover a distinction the source never wrote down.

How much of the domain survives the IR is the measure of a translation's clarity and correctness. When two encodings are both valid, the standards choose the one that keeps more of the domain.

**VM-PRES-2 (Invariant). `.sets` works with fundamental domains.** The goal is the data structures themselves, which don't change between languages. A class or library type that wraps a fundamental (a JavaScript `Set` over an object, a `Date` over a number) isn't part of the core standard. A language standard defines its wrappers separately, and even common ones should narrow to the fundamentals they wrap. A TypeScript `Set`, for example, is an object whose keys are its members and whose values are empty, so the TypeScript standard can read it as a record with no special treatment (VM-SET-1):

```text
new Set(["a", "b"])   is the value   {("a", ()), ("b", ())}
```

The empty value may be written `()` or `Phi`; in the core context they're treated as the same (VM-AX-5). A language standard binds it to its own nothing-like unit, such as Go's `struct{}{}`.

**Collisions.** When two source values become one `.sets` value, first ask why. If they describe the same domain under different source semantics, the collision is justified: there is no difference in the data to preserve, so VM-PRES-1 holds. Only a collision between different data is a defect to fix.

## 1. Value kinds

**VM-KIND-1 (Invariant).** Every value in the core model has exactly one of these kinds:

| Kind | Examples | Written in `.sets` as |
| ---- | -------- | --------------------- |
| String | `"abc"`, `""` | quoted symbol |
| Number | `0`, `-7`, `2.5` | numeric symbol |
| Atom | `true`, `Circle` | bare symbol |
| Tuple | `("abc")`, `(1, 2)`, `()` | parentheses |
| Domain | `{("name", "abc")}`, `{1, 2}`, `{}`, `Integer` | braces, or a domain name |

Any domain is a value, whatever its size. A domain can therefore be a member of another domain: `{Integer}` is the domain whose only member is the domain `Integer`. Records, sets and maps aren't separate kinds; they are domains whose members are values (§6).

A tuple can have any number of components, and each component can be a value of any kind.

Values of different kinds are never equal: `"1" ≠ 1 ≠ one ≠ (1) ≠ {1}`.

**VM-KIND-2 (Invariant). Top-level domains are undefined domains.** A top-level axiomatic domain (`Integer`, `String`) has no members of its own in `.sets`; it's an *undefined* domain, not an infinite one. Its members are fixed by the standard of the context that uses it: the language standard and its view.

Inside `.sets`, top-level domains are always disjoint (GUIDE, "Top domain axiom"), and `.sets` code must keep to that. A language that doesn't have a separate kind for each domain doesn't have to preserve the distinction: its standard *merges* the domains into one representation. For example:

```sets
Integer
Fractional
Number = union(Integer, Fractional)
```

In `.sets`, `Integer` and `Fractional` are disjoint and `Number` is their union. A JavaScript standard has only one number kind, so it represents all three with the JavaScript `number` type, and in JavaScript `Integer = Fractional = Number`. A C standard keeps them apart, with `Number` strictly larger than either.

A standard may merge in one pass, or in several (for example, first map `Integer` and `Fractional` each to a number type, then merge the identical results). The `.sets` domains don't change either way; merging is a representation decision of the language standard.

**Why keep the distinction in `.sets` at all.** A target can always merge two domains, but it can never split one that the IR doesn't separate. Keeping `Integer` and `Fractional` apart preserves the domain for every target that can use the difference ([VM-PRES-1](#0-principle-preserve-the-domain)).

This only works when the source makes the distinction. A source language without it (TypeScript, whose only number type is `number`) always produces one combined domain, so the IR has no separate `Integer` and `Fractional`. A target with separate integer and floating-point types then can't tell which one a value was meant to be, and choosing correctly is impossible from the IR alone. This is an unsolved problem for the `base` standard; see [UP-1](#unsolved-problems).

## 2. Lexical conventions for roster symbols

`GUIDE.md` defines a symbol as a run of UTF-8 characters, and says quoting is a convention that "the transpiler ecosystem" may agree on. This section is that agreement, accepted for now.

**VM-LEX-0 (Invariant). Two parsing layers.** A `.sets` parser has two layers:
1. **The grammar layer** implements GUIDE only. It splits rosters at unescaped `,`, `{` and `}`, applies GUIDE's whitespace rule, and produces symbols. It gives quotes no special meaning. It keeps escape sequences in the symbol text as written; it only uses them to decide where a symbol ends.
2. **The value-model layer** implements this document and the `base` rules. It reads each symbol as a string, number or atom, decodes escape sequences, and checks quotes.

Both readings of the same text must agree. Writing `{"a\, b"}` satisfies both: the grammar layer sees one symbol because the comma is escaped, and the value-model layer reads it as the string `a, b`.

An escaped separator isn't a separator. GUIDE's rule that ignores one whitespace byte around a reserved symbol applies only to the unescaped `,`, `{` and `}` of the grammar, never to `\,`, `\{` or `\}`. So in `{"a\, b"}` the space is part of the string.

**VM-LEX-1 (Invariant). Quoted symbols are strings.** A symbol that starts and ends with an unescaped `"` is a String; the surrounding quotes aren't part of the value. `"Succeeded"` is the String value `Succeeded`.

Inside a symbol, the value-model layer expects every unescaped `"` to be closed by another one:
- `"a\"b"` is the string `a"b`.
- `a\"b` is the atom `a"b`.
- `a"b` is an error: its `"` is never closed.

Separators are escaped the same way, inside quotes or not, because the grammar layer doesn't know about quotes (VM-LEX-0):
- `{"a\, b"}` is one string, `a, b`.
- `{"a, b"}` is an error: the grammar layer splits it into `"a` and ` b"`, and neither has a closing quote.
- `\{` and `\}` are a literal brace, and `\\` is a literal backslash.

**VM-LEX-2 (Invariant). Numeric symbols are numbers.** A bare symbol that matches `[+-]?[0-9]+(\.[0-9]+)?` is the number it spells in base 10. A leading `+` or `-` gives the number's direction: `+3` and `3` are the same number, and `-3` is its negative.

- **The decimal point decides the axiom.** A symbol without a decimal point is in `Integer`; a symbol with one is in `Fractional`, even when its fractional part is zero. `3` is in `Integer`; `3.0`, `2.5` and `-2.5` are in `Fractional`. So `1` and `1.0` are different values.
- Within `Fractional`, trailing fractional zeros don't change the value: `1.5` and `1.50` are the same number.
- Leading zeros are allowed and don't change the value: `007` is `7`, always read in base 10, never as octal.
- `-0` and `+0` are the integer `0`, and `-0.0` is the fractional `0.0`. IEEE-754 negative zero is neither; it's the atom `NegativeZero` (VM-AX-3).
- The canonical spelling omits `+`, strips leading zeros (keeping a single `0` before the point), and strips trailing fractional zeros while keeping at least one digit after the point: `+007.50` is written `7.5`, and `1.00` is written `1.0`.

**VM-LEX-2a (Default). Other numeric notations.** The core standard defines no exponent or other scientific notation (`1e3`, `0x1F`, `1_000`), because each one is an extra parsing requirement for every adapter. In core, such symbols are atoms. A language standard may define its own notations: either parse them as numbers, or keep them verbatim as atoms. It says which in its Overrides table.

**Note: not every number can be written.** A literal is a finite decimal, so a rational such as one third has no literal, and some exact binary values need very long ones (the binary64 value nearest 0.1 has 55 significant digits). This is unsolved (UP-2). The core standard deliberately doesn't compute values with operators to fill the gap: this iteration of the standards is about non-deterministic transpilation, while math must stay deterministic, so literal arithmetic is avoided wherever possible. Operators appear in membership rules (VM-AX-6), not as a way of writing values.

**VM-LEX-3 (Invariant). Other bare symbols are atoms.** An atom is an opaque value equal only to itself. It isn't a member of `String`, so `Succeeded ≠ "Succeeded"`.

**VM-LEX-4 (Invariant). Escaped symbols aren't domain names.** A symbol that contains an escape (VM-LEX-6) is never read as a domain name. To write an atom spelled like a domain in scope, escape one of its characters; the canonical form escapes the first one. `{\Number}` is the atom `Number`; `{Number}` is the domain `Number`.

`@` has no special meaning: `@Number` is an ordinary atom spelled `@Number`. Atoms are written without escapes whenever that's unambiguous.

**VM-LEX-5 (Invariant). Universal and empty domain spelling.** The universal domain is `U` and the empty domain is `Phi`, both declared in `base`. There are no aliases: `u` and `phi` aren't defined, and a reference to either is an unknown domain unless the file declares it itself.

**VM-NAME-1 (Invariant). Declared names can't be redeclared.** A domain or function that's already declared in scope can't be declared again. That covers names from `base`, names imported with `X in file`, and names declared earlier in the same file. Redeclaring `U`, `Phi`, `union` or `String`, for example, is an error. A file that wants a different meaning uses a different name.

**VM-LEX-6 (Invariant). Escapes make new symbols.** A backslash can escape any character, in any symbol: string, atom or otherwise. The escaped character is taken literally, outside every rule it would otherwise follow, and the result is a new symbol in `U`:
- `\,`, `\{`, `\}` aren't separators (VM-LEX-0).
- `\"` doesn't open or close a string (VM-LEX-1).
- `\1` is an atom spelled `1`, not the number `1` (VM-LEX-2), because the escaped digit isn't read as numeric.
- `\Number` is the atom `Number`, not the domain `Number` (VM-LEX-4).
- `\\` is a literal backslash.

Escaping a character that triggers no rule changes nothing: `a\qb` and `aqb` are the same atom. The canonical form escapes only where a rule would otherwise apply.

**Control characters.** A few escapes stand for characters that can't be written literally. They work in any symbol:

| Escape | Character |
| ------ | --------- |
| `\n` | line feed, U+000A |
| `\r` | carriage return, U+000D |
| `\t` | tab, U+0009 |
| `\u{H}` | the Unicode scalar value with hex code `H` (1 to 6 hex digits), written only as `\u` followed by `{`, for example `\u{0}` or `\u{1F600}` |

Control characters (U+0000–U+001F and U+007F) aren't allowed raw in a symbol; they're always written as escapes. The canonical form uses `\n`, `\r` and `\t` where they apply and `\u{H}` with uppercase hex and no leading zeros for the rest. A `\u{H}` for a surrogate code point (U+D800–U+DFFF) or a value above U+10FFFF is an error.

## 3. Rules and binders

**VM-RULE-1 (Invariant). Rules follow boolean algebra.** Rules combine with the usual precedence from math, highest first:

| Level | Operators |
| ----- | --------- |
| 1 | arithmetic (VM-AX-6): `^`, then `*`, `/`, `%`, then `+`, `-` |
| 2 | comparisons and membership: `<`, `<=`, `>`, `>=`, `in`, `not in`, `sub` |
| 3 | `not` |
| 4 | `and` |
| 5 | `or` |

Operators at the same level group left to right, except `^`, which groups right to left as in math (`2 ^ 3 ^ 2` is `2 ^ 9`). `A or B and C` is `A or (B and C)`; `A and B or C and D` is `(A and B) or (C and D)`.

`not` is boolean negation. It applies to the rule on its right and binds tighter than `and`: `not A and B` is `(not A) and B`, and `not x in D` is `not (x in D)`. `x not in D` is the same as `not x in D`; both spellings are allowed. `not not A` is `A`.

Two things that look like operators aren't:
- **Chained comparisons.** `0 <= x <= 10` isn't a chain. Ranges are written with `and`: `0 <= x and x <= 10`. A comparison whose operand is itself a comparison is an error, not a silently false rule.
- **Unary signs.** A sign is part of a numeric literal (VM-LEX-2), not an operator. Negating a binder is written `0 - x`; `-x` isn't defined.

**VM-RULE-1a (Invariant). Parentheses.** Inside a rule, parentheses mean one of three things:
1. **Call.** Directly after a domain or function name, they call it: `Seq(T)`, `singleton(m * 2 ^ e)`.
2. **Grouping.** Around a rule expression (anything containing `in`, `not in`, `sub`, `not`, `and`, `or` or a comparison) or an arithmetic expression, they group it, as in math: `A and (B or C)`, `(x + 1) % 2`. Grouping parentheses never create a value.
3. **Tuple.** Otherwise they're a tuple. A tuple's components are values and binders only: `("name", name)`, `(x)`, `()`. A single value in parentheses is always a 1-tuple, so `(x) + 1` has no result. Arithmetic isn't allowed inside a tuple; a computed component is bound first: `(z, y)` with `z in singleton(x + 1)`.

Only tuple parentheses make values, so this doesn't conflict with "parentheses are values" (GUIDE, "Tuples and Vectors").

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

The binder is existential over the whole rule even under `not`. `not ("name", n) in r` is valid and means "there is some `n` for which `("name", n)` isn't in `r`", which almost every record satisfies. That reading may not be the intuitive one, but it's well defined. The standard defines what a rule means, not what an adapter meant: an adapter that wants "`r` has no `name` field" writes a rule that says so, for example `r not in { s suchthat ("name", n) in s }`. An adapter that emits a valid but unintended rule has a bug in the adapter, not in the rule.

**VM-RULE-3 (Invariant). Literals in patterns.** A pattern such as `("name", name)` may mix literal values (`"name"`) and binders (`name`). A literal matches only an equal value (VM-EQ-1).

| Pattern | Matches | Doesn't match | Why |
| ------- | ------- | ------------- | --- |
| `(Integer, x)` | `(Integer, 5)` | `(3, 5)` | `Integer` is a domain name, so it matches the domain `Integer` itself (VM-KIND-1), not a number in it. To test membership, bind the component and add a clause: `(k, x) in r and k in Integer`. |
| `(x, x)` | `(1, 1)`, `("a", "a")` | `(1, 2)`, `(1, 1.0)` | A binder used twice must match equal values. |
| `(1.0, x)` | `(1.0, 5)`, `(1.00, 5)` | `(1, 5)` | `1.0` and `1.00` are the same `Fractional` value; `1` is an `Integer` (VM-LEX-2). |

**VM-RULE-4 (Invariant). `sub` in a rule.** `s sub T` holds when every member of `s` is a member of `T`. It means "subset or equal": `s sub s` holds, and a proper subset needs no separate operator because it's already covered.

Both sides must be domains. This is where the two parsing layers (VM-LEX-0) differ:
- In plain `.sets` (the grammar layer), `5 sub U` is valid, because `5` is just a symbol and could name a domain.
- Under this value model, `5` is a number (VM-LEX-2), not a domain, so `5 sub U` is an error. The value-model layer reports it rather than treating it as vacuously true.

## 4. Axioms

`base.sets` declares three top-level axioms. By GUIDE, top-level domains are pairwise disjoint.

Top-level domains are undefined in `.sets` (VM-KIND-2). The table gives each one's meaning in the core context. A language standard may override that meaning for its own context, or merge domains it can't tell apart.

| Axiom | Core meaning |
| ----- | ------- |
| `String` | Every finite sequence of Unicode scalar values (U+0000–U+D7FF, U+E000–U+10FFFF), including the empty string. Surrogate code points aren't string content. |
| `Integer` | Every mathematical integer, unbounded. |
| `Fractional` | Every finite real number written with a fractional part (VM-LEX-2), including those whose fractional part is zero: `1.0`, `2.5`. |

**VM-AX-1 (Invariant).** `String`, `Integer` and `Fractional` are disjoint. `Real = union(Integer, Fractional)` is every finite real number, in both forms.

So `1` and `1.0` are both members of `Real`, as two different values. They're different symbols with the same mathematical meaning; treating them as the same is a semantic choice, not something the value model makes. An adapter may make that choice where it fits. For example, a literal `1.0` that's only serialized, never stored in a typed variable, may be emitted as the integer `1`.

**VM-AX-2 (Invariant). Numbers have no width.** A number value has no width, signedness or machine representation; it's an exact mathematical number in one of two forms, `Integer` or `Fractional` (VM-AX-1). Width is a constraint on the domain (§5), not on the value: `300` is the same value whether a source declares it `int16` or `int64`. Which form a source value takes is fixed by the source standard's view. A binary64 `1.0` may be the `Fractional` `1.0` or the `Integer` `1`, depending on whether that view keeps the distinction (VM-PRES-1, UP-1).

**VM-AX-3 (Default). Non-finite values are atoms.** IEEE-754 NaN, `+∞`, `−∞` and `−0` are the atoms `NotANumber`, `PositiveInfinity`, `NegativeInfinity` and `NegativeZero` (the roster `IEEESpecial`). All NaN payloads are the single atom `NotANumber`. Folding them loses the payload, which some languages can observe; a language standard that needs to keep it overrides this rule with its own atoms (VM-PRES-1). A source that has a negative zero emits `NegativeZero`. A view whose numbers have no signed zero never produces one, so there's nothing to fold.

**VM-AX-4 (Default). Booleans are enum-like atoms.** `Boolean = {true, false}` is a roster of two atoms, like any enum. The atoms are descriptors for data: they have no truthiness, and `.sets` doesn't treat `true` as holding or `false` as failing. A language whose booleans are something else (integers `0` and `1` in C) maps them by its own source standard.

Truth itself isn't a value. A rule holds or it doesn't, and two parameterless predicates name the two outcomes:

```sets
Truef: \\ in U
Falsef: \\ in Phi
```

`Truef` always holds, because every symbol is in `U`; `Falsef` never holds, because `x in Phi` never holds as a predicate (VM-AX-5). That's true in every context. Adapters that need data use the `Boolean` atoms; truthiness rules of a language (what counts as true in an `if`) belong to that language's standard.

**VM-AX-5 (Invariant). One nothing.** The core names exactly one nothing, `Nothing`, declared by a deliberate paradox:

```sets
Nothing in Phi
```

`Phi` isn't nothing: it's a domain, something that contains nothing. No value can be a member of `Phi`, so the one declared member is unlike every value, and that makes it nothing. A top-level `Nothing in Phi` is a declaration: it states that `Nothing` is in `Phi`, which makes `Nothing` what it is. It adds no member to `Phi`. Inside a rule, `x in Phi` is a predicate that never holds, as with `Falsef` in VM-AX-4: `Phi` stays `{}`, and `{ x suchthat x in Phi }` is `Phi`.

`Nothing` is in `U`: the moment it enters `U`'s discourse, it becomes part of it.

Another top-level declaration into `Phi`, such as `Abcd in Phi`, is undefined: `Abcd` may or may not equal `Nothing`. That isn't a special case the standard adds; it follows from `Phi`. Declaring a member into a domain with no members is itself a paradox, and whether non-existence is one state or several is unknown. `Phi` may hold several nothings, each the same and yet different, because the domain of discourse keeps changing as superdomains grow. VM-LEX-3 still holds for `Abcd`, read in this context. The core doesn't decide the question, so no standard can rely on it.

There's no other core null atom. A value that may be absent-as-a-value (a nullable value) is `union(T, {Nothing})`: either a member of `T`, or `Nothing`. A language standard may still declare its own null-ish atoms (`null`, `undefined`, `nil`) when its view distinguishes several kinds of nothing; core gives none of them meaning.

`()` and `{}` are conceptually equal, not literally. A tuple states the order of its components; with no components there is no order to state. If `Phi` held more than one nothing, order among them would matter, and `()` and `{}` would differ. The core deals with one nothing, so within its context the two are treated as equal and either may be written. Neither is `Nothing`.

**VM-AX-6 (Invariant). Numeric operators are mathematical.** Numbers are exact mathematical values with no width (VM-AX-2), so operators on them have their mathematical meaning:
- `+`, `-`, `*`, `/` give exact results, with no overflow, rounding or wraparound. `1 / 3` is the exact rational one third. Division by zero has no result.
- `a ^ n` is exponentiation with an `Integer` exponent `n`. For `n >= 0` it's defined for every number `a` (`0 ^ 0` is `1`); for `n < 0` it's `1 / a ^ (0 - n)` and needs `a` to be non-zero. An exponent that isn't an `Integer` has no result, including `2 ^ 2.0`, so results stay exact rationals: `2 ^ -2` is `0.25`.
- `%` is defined on `Integer` operands only. The result has the sign of the divisor (mathematical modulo): `-7 % 3` is `2`. A `Fractional` operand (`7.0 % 3`) or a zero divisor gives no result.
- `<`, `<=`, `>`, `>=` are the ordering of the reals. They compare mathematical values and ignore the form: `1 < 1.5` and `1 <= 1.0` both hold.
- An operator applied to a value that isn't a number (string, atom, tuple, domain) has no result. Nothing is coerced: `"1" + 1` doesn't exist.
- A comparison or membership test whose operand has no result is false, so the rule doesn't hold for that value. Under `not` it becomes true: `not "a" < 1` holds. That's valid, if unintuitive (VM-RULE-2).

**Result form.** Every result is an `Integer` or a `Fractional` (VM-AX-1):
- `+`, `-`, `*`, `%`, and `^` with `n >= 0`, on `Integer` operands only, give an `Integer`: `2 + 3` is `5`.
- `/`, `^` with `n < 0`, and any operation with a `Fractional` operand give a `Fractional`: `4 / 2` is `2.0`, and `1.5 + 1.5` is `3.0`.

Results compare by mathematical value but are equal only within one form, so `2 + 3 <= 5.0` holds while `2 + 3 in singleton(5.0)` doesn't.

Equality stays `x in singleton(y)` (GUIDE, "Singletons"); `=` and `!=` aren't operators.

**VM-AX-7 (Default). New axioms.** A language standard adds a top-level axiom only for values that no `base` domain contains (for example, opaque handles in a runtime view). Values that exist in `base` (dates, UUIDs, decimals) are described subdomains of `base` domains instead, because a new axiom would be disjoint from them.

**VM-AX-8 (Invariant). Axioms stay disjoint.** Every top-level axiom is disjoint from every other, in `base` and in every language standard (GUIDE, "Top domain axiom"). No override can relax this: a language standard that overrides VM-AX-7 still can't declare an axiom that shares a value with an existing one. Languages that can't tell two axioms apart merge their representations instead (VM-KIND-2); the axioms themselves stay disjoint.

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

**Both forms.** By VM-AX-6's result form, `m * 2 ^ e` is an `Integer` when `e >= 0` and a `Fractional` when `e < 0`. A whole number can be built both ways (`2 = 1 * 2 ^ 1 = 4 * 2 ^ -1`), so the float domains contain both forms of every representable whole number: `2` and `2.0` are both members of `Binary64Finite`. Values with a fractional part, such as `0.5`, are members only as `Fractional`. The integer ranges hold one form only: `Int32` requires `x in Integer`, so `1.0` isn't an `Int32`.

`0` and `0.0` are members (`m = 0`); `−0` isn't, because it's the atom `NegativeZero` (VM-AX-3). `Float64` and `Float32` add `IEEESpecial` with `union`. `Binary32Finite` is inside `Binary64Finite`, and every `Int32` is in `Binary64Finite`, both by the rules.

A top-level `X sub Y` declares a domain `X` whose members are drawn from `Y`, at most all of `Y` (GUIDE, "Difference between `sub` and `in` in this context"). Domains that the grammar still can't state are declared that way, and the prose here fixes each one's exact extent. Target adapters bind these domains by this prose. A described subdomain with no prose extent is read as GUIDE reads it: the same as `Y`.

| Domain | Exact members |
| ------ | ------------- |
| `Char` | Strings of exactly one Unicode scalar value. (Needs a string-length operator.) |

**VM-SUB-1 (Invariant).** This rule covers the described subdomains of `base.sets` only, so it stays language-agnostic. Those subdomains overlap the way their prose says, regardless of the single `sub` parent written in `base.sets`. Prose is a stopgap for what the grammar can't state yet: once the grammar can express a described subdomain's extent, the subdomain is replaced by a builder, the way the ranged domains were (VM-SUB-0). `Char`, for example, becomes a builder when a string-length operator exists.

**When prose contradicts the grammar.** `X sub Y` is a grammar fact: every member of `X` is a member of `Y`. Prose can only narrow that, by saying which members of `Y` belong to `X`. If the prose instead names a value outside `Y` (for example, prose for `Char sub String` that admits the number `1`), the declaration says two incompatible things:
- Read by its declaration, `X` can't contain the value, because the value isn't in `Y`.
- A target that binds `X` by its prose would accept it.

The two readings then disagree about the same domain, and no representation can be correct for both, so the declaration is a defect in `base.sets`. It's fixed by correcting the prose or the `sub` parent, never by choosing one reading at bind time.

**VM-SUB-2 (Default).** A language standard writes further ranged domains as builders (VM-SUB-0), and declares the remaining described subdomains in its own superdomain using the `sub` form (`Uuid sub String`, `Rfc3339Timestamp sub String`). Each one needs a prose extent precise enough to write a membership check from.

The language standard owns its described subdomains the way `base` owns its own (VM-SUB-1): it states how they overlap, and prose that contradicts a subdomain's `sub` parent is a defect in that standard.

This rule is Default on purpose: how a language standard describes its domains (prose, builder or structured value) is its own choice.

A value whose representation carries information the number alone doesn't is a structured value, not a subdomain of `Real`. A decimal that keeps its scale is the usual case: decimal128 stores `1.50` as coefficient `150` and exponent `-2`, and `1.5` as `15` and `-1`. They're equal numbers but different values (`1.50 * 2` is `3.00`). Since `1.50` and `1.5` are the same `Fractional` (VM-LEX-2), `Decimal128 sub Real` would lose the scale (VM-PRES-1). A language standard that keeps it writes the value as a tuple instead:

```sets
Decimal128 = {
    (c, q) suchthat
    c in Integer and c > 0 - 10 ^ 34 and c < 10 ^ 34 and
    q in Integer and q >= -6176 and q <= 6111
}
```

Here `1.50` is `(150, -2)` and `1.5` is `(15, -1)`, two members. Special values (NaN, infinities, negative zero) are added with `union`, as `Float64` does.

## 6. Structural encodings

### 6.1 Records

**VM-REC-1 (Invariant). A record value is a set of tuples.** Each member is a field `(key, value)`. A key can be any value, and keys are unique within a record. The record with no fields is `{}`, a set with no tuples. It's the empty domain `Phi`, not a missing value, which is `Nothing` (VM-AX-5). The domain of all records in `base` is:

```sets
Record = Map(U, U)
```

```text
{ name: "abc", age: 123 }   is the value   {("name", "abc"), ("age", 123)}
{ 1: "abc" }                is the value   {(1, "abc")}
{ }                         is the value   {}
```

`{}` is also the empty set, the empty map and, in context, `[]` (VM-AX-5). That's the same data, nothing inside, under different source semantics, so the collision is allowed (VM-PRES-2).

Field names are part of the value, because `{"name": 1}` and `{"nom": 1}` are different data on wires that name fields. A positional tuple `(name, age)` can't keep that difference without relying on binder names, and binder names are scoped symbols that renaming must not affect.

**VM-REC-2 (Default). Preferred record form.** A record domain is written as a builder over a record binder `r`, with each required field as a pattern on `r`:

```sets
Person = { r suchthat r in Record and ("name", name) in r and name in String }
```

`r in Record` is required (by convention written first). Without it, `{("name", "a"), ("name", 1)}` would be a `Person` even though its keys repeat, so it isn't a record (VM-REC-1). So would `{("name", "a"), 5}`.

This form is preferred, not the only one. How a domain is written can change while the data stays the same: in memory a record or a map is a list of bytes, contiguous or fragmented, that could also be read as an array or a set. An adapter may write a record domain another way when it describes the same domain. The standard gives this form so that adapters can recognize the domain easily.

**VM-REC-3 (Invariant). Records are open.** A record domain admits fields it doesn't mention, with any value. This follows from the builder: a rule constrains only the fields it names, so the same `.sets` text means the same domain to every adapter. A closed source type (a Go struct, a sealed class) says so with `OnlyKeys` (VM-REC-4).

```sets
Person = { r suchthat r in Record and ("name", name) in r and name in String }
Aged = { r suchthat r in Record and ("age", age) in r and age in Integer }
Human = union(Person, Aged)
```

`{("age", 123), ("name", "abc")}` is a `Person`. This is the README `Human` reading: that value is both a `Person` and an `Aged`, so it's a `Human`.

**VM-REC-4 (Invariant). Closed records say so.** A record domain that rejects unmentioned keys adds `r in OnlyKeys({k1, …, kn})`, listing every key it admits: its required keys and its optional keys (VM-MISS-2).

```sets
ClosedPerson = { r suchthat r in Record and ("name", name) in r and name in String and r in OptionalField("age", Integer) and r in OnlyKeys({"name", "age"}) }
```

`{("name", "abc")}` and `{("name", "abc"), ("age", 1)}` are `ClosedPerson`s; `{("name", "abc"), ("id", 1)}` isn't, because `"id"` isn't listed.

### 6.2 Missing fields and null

**VM-MISS-1 (Invariant). A missing field is a missing tuple.** Absence isn't a value: a missing field has no tuple at all. A field that's present with nothing in it is a tuple whose value is `Nothing` (VM-AX-5). These are different records:

```text
{name: "abc"}               is   {("name", "abc")}
{name: "abc", age: null}    is   {("name", "abc"), ("age", Nothing)}
```

Kinds of absence in between, such as JavaScript's `undefined`, belong to the language standard.

**VM-MISS-2 (Invariant). Optional fields.** `r in OptionalFields(K, V)` says: every field of `r` whose key is in `K`, if present, has a value in `V`. `K` is a domain of keys, so one call can make several fields optional with the same value domain. `OptionalField(k, V)` is the single-key helper, `OptionalFields({k}, V)`. In `base`:

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

**VM-TUP-1 (Invariant). Tuples are ordered.** An n-tuple `(v1, …, vn)` is ordered and nests without flattening (GUIDE, "Tuples and Vectors"). A tuple with at least one component is its own kind of value (VM-KIND-1), never a record or another domain; the empty tuple `()` is treated as `{}` in the core context (VM-AX-5). Sequences are built from tuples (VM-SEQ-1), so a nested pair ending in `()`, such as `(a, (b, ()))`, is also a sequence. Fixed-length positional products (language tuples, fixed-size arrays) are tuples.

### 6.4 Sequences, sets and maps

**VM-SEQ-1 (Invariant). Sequences are nested pairs.** A finite sequence is a nested pair ending in the empty tuple: `[]` is `()`, and `[a, b]` is `(a, (b, ()))`. `Seq(T)` is the domain of finite sequences over `T`. Variable-length lists (JS arrays, Go slices, Python lists) are sequences; fixed-length ones are tuples (VM-TUP-1), and the source's type decides which. Because `()` is treated as `{}` (VM-AX-5), `[]` collides with the empty record, set and map; that collision is allowed (VM-REC-1). Target adapters bind `Seq` to native lists by meaning; they don't generate nested pairs.

**VM-SET-1 (Invariant). Sets are domains.** A set of values is a domain, so `{"a", "b"}` is the set of `"a"` and `"b"`. `SetOf(T)` is the domain of all subdomains of `T` (`s sub T`). A language's set type is finite, and its standard narrows to that; it also decides whether the set reads as a `SetOf` value or as a record of keys (VM-PRES-2).

**VM-MAP-1 (Invariant). Maps are domains of pairs.** A map is a domain of `(key, value)` pairs with unique keys. `Map(K, V)` is every such domain whose pairs are in `product2(K, V)`; `base` states uniqueness in its builder. A record is a map (VM-REC-1). A language's map type is finite, and its standard narrows to that.

### 6.5 Recursion

**VM-RECUR-1 (Invariant). Domains may be recursive.** A domain may refer to itself, directly or through other named domains. A recursive domain holds exactly the values built from its rule in finitely many steps (its least fixed point): `Seq(T)` holds every finite sequence, but not an endless one. So `X = X` is `{}`, and `X = union(X, {1})` is `{1}`. A member may still be an infinite domain (VM-KIND-1); only the building is finite. Cyclic data, such as an object that refers to itself, isn't a value; its language standard encodes references (VM-PRES-2).

## 7. Equality

**VM-EQ-1 (Invariant). Equality.** A value that isn't built from other values (a string, a number or an atom) is compared as a symbol. Two such values are equal exactly when their canonical spellings are the same (VM-LEX-*), with no other normalization. So `1.5` is `1.50`, `a\qb` is `aqb`, `1` isn't `1.0`, and `"1"` isn't `1`.

A number computed by an operator (VM-AX-6) may have no spelling, such as `1 / 3` (UP-2). It's equal to a number of the same form (`Integer` or `Fractional`) with the same mathematical value. For numbers that can be written, that's the same as comparing canonical spellings.

Built values compare by structure:
- tuples: the same arity and pairwise equal components;
- domains: the same members. Undefined domains (top-level axioms, VM-KIND-2) are disjoint from one another, and each is equal only to itself. `()` and `{}` are treated as equal in the core context (VM-AX-5).

`Nothing` is equal only to itself; whether another declaration into `Phi` equals it is undefined (VM-AX-5). `x in singleton(y)` holds exactly when `x` and `y` are equal (GUIDE, "Singletons").

## 8. Wire projections

**VM-WIRE-1 (Invariant). A projection is one-to-one or it says so.** A wire projection maps values to a wire format and back. It is one-to-one on each domain and rejects non-members. Where a domain admits two values with the same projection, that is a diagnostic, not a guess. The default projection is JSON, specified in [JSON-MODEL.md](JSON-MODEL.md) (VM-JSON-1, VM-JSON-2). A language standard that uses another wire defines its own projection under this rule.

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
| SI-18 | GUIDE's `base.sets` example lists only `u`, `phi` and the operations, and defines `U = u` and `Phi = phi` as aliases. `base` now also holds the axioms, described subdomains and structural transformers. | [`base.sets`](base.sets) is the full `base`, with `U` and `Phi` and no lowercase aliases (VM-LEX-5). GUIDE's example could point to it. |
| SI-6 | GUIDE says function parameters must be domains, but `IsOne(x)` and `ToUser(x, y)` take bare parameters. | Core only uses `{X sub D}` and `{x suchthat …}` parameters. |
| SI-7 | `=`, `sub or =`, `!=` (GUIDE) and `subsetof` (README `Human`) appear but aren't defined. | Not used. `%` is now defined by VM-AX-6. |
| SI-8 | `base` shows a variadic `union({A sub u}, {B sub u}, ...)`; variadic parameters aren't defined. | Binary `union`/`intersect`, nested right-associatively (CS-CONV-4). |
| SI-12 | "`{x} in U` *can* be written as `{x} sub U`" mixes membership and subset. | In a rule, `s sub T` means every member of `s` is in `T` (VM-RULE-4; used by `SetOf`, `Map`, `OnlyKeys`). |
| SI-14 | There's no grammar for cardinality, finiteness, string length or patterns. | Prose rules for described subdomains (§5) and the gaps in each standard. Ranges and floating-point domains are covered by VM-AX-6; key uniqueness is in the `Map` builder (VM-MAP-1); finiteness of a language's sets and maps is left to its standard (VM-SET-1, VM-MAP-1). |
| SI-15 | `//` comments are used throughout but aren't part of the grammar. | Allowed; they carry no meaning. |
| SI-16 | README says TypeScript's `Person \| Aged` is an "intersection". As a domain it's a union of open records whose members can satisfy both branches. | Terminology only. `union(Person, Aged)` with open records (VM-REC-3). |
| SI-19 | Is arithmetic needed in `.sets` at all? VM-AX-6 operators exist for ranges (`Int32`) and float formats (`m * 2 ^ e`), but they bring evaluation into a language of domains, and the VM-LEX-2 note already avoids literal math. | Open question for the spec author. Arithmetic stays allowed in comparisons and call arguments, and is forbidden inside tuples (VM-RULE-1a). |
| SI-17 | GUIDE doesn't say how symbols that appear only in a rule (`n` in `("name", n) in r and n in String`) are bound. | Existentially over the whole rule (VM-RULE-2). |

### Unsolved problems

Problems the core standards know about but don't solve yet. These are concerns of the `base` standard, not ambiguities in GUIDE.

| # | Problem | Current state |
| - | ------- | ------------- |
| UP-1 | A source language with one number type (TypeScript `number`) can't say whether a value is an integer or a fraction, so its IR has one combined number domain (for example `Binary64Finite`) instead of separate `Integer` and `Fractional` domains. A target with distinct integer and floating-point types can't choose between them correctly from the IR (VM-KIND-2, VM-PRES-1). | Unsolved. Targets use the combined domain's representation (a floating-point type). Possible directions, none adopted: source-side inference from usage, or a source standard that declares integer-only described subdomains from annotations. |
| UP-2 | A literal is a finite decimal (VM-LEX-2), so a number with no finite decimal form (one third) can't be written, and some exact binary values need very long literals. | Unsolved. Writing values through arithmetic is out of scope for now: transpilation in this iteration is non-deterministic and math must stay deterministic, so literal arithmetic is avoided wherever possible. |

### Resolved

| # | Issue | Resolution |
| - | ----- | ---------- |
| SI-1 | The universal domain was spelled both `u` and `U`, and the empty domain both `phi` and `Phi`. | `U` and `Phi` only; the lowercase aliases are dropped (VM-LEX-5). |
| SI-2 | Record members were `{name, age}` in GUIDE and `(name)` in README. | Records are sets of `(key, value)` tuples: `{r suchthat ("name", n) in r and n in String}` (VM-REC-*). |
| SI-3 | There was no shared axioms superdomain. | `core` is merged into `base`. |
| SI-4 | Quoting was only a convention. | VM-LEX-1 to VM-LEX-4 are accepted for now. |
| SI-5 | `and`/`or` had no precedence. | Boolean algebra with `()` grouping (VM-RULE-1, VM-RULE-1a). Pending in GUIDE; see the notes below. |
| SI-9 | The empty tuple wasn't defined. | `()` is the 0-tuple. It's conceptually equal to `{}`, not literally: with no components there is no order to state, and the core treats them as equal (VM-AX-5). |
| SI-10 | Recursion wasn't defined. | Allowed (VM-RECUR-1). |
| SI-11 | Top-level `X sub Y` seemed undefined. | Defined in GUIDE, "Difference between `sub` and `in` in this context" (§5). |
| SI-13 | Literals inside patterns weren't defined. | Allowed (VM-RULE-3). |

### Notes for the spec author

Decisions made in review that belong in `GUIDE.md`. The standards already rely on them.

1. **Rule precedence and grouping** (VM-RULE-1, VM-RULE-1a). Rules follow boolean algebra: arithmetic, then comparisons and membership, then `not`, then `and`, then `or`, each level grouping left to right (except `^`). Parentheses after a name are a call, parentheses around a rule or arithmetic expression group it, and parentheses around values and binders are tuples, which can't contain arithmetic. GUIDE's "Tuples and Vectors" section should say that grouping parentheses aren't values.
2. **Numeric operators** (VM-AX-6). `^ * / % + -` and `< <= > >=` on numbers, with exact mathematical meaning and no coercion; `^` takes integer exponents and groups right to left. This makes ranged and floating-point domains expressible (`Int32`, `Binary64Finite` as builders), which GUIDE's number discussion could show.
3. **Rule-local binders** (VM-RULE-2, SI-17). A symbol that first appears inside a rule is existential over the rule.
4. **Names** (VM-LEX-5, VM-NAME-1). `U` and `Phi` are the only spellings of the universal and empty domains; GUIDE's `u` and `phi` aliases go away. A declared domain or function can't be redeclared, which GUIDE's "Domain Validation" section could state.
5. **Nothing** (VM-AX-5). `Nothing in Phi` declares the one nothing, as a deliberate paradox. `Phi` is the empty domain, not nothing. `()` and `{}` are conceptually equal, not literally. GUIDE's "Tuples and Vectors" section could say so.
6. **Escapes** (VM-LEX-0, VM-LEX-6). An escaped character makes a new symbol and isn't a separator. GUIDE's symbol definition could state the two parsing layers.
