# Core Value Model

- **Standard version:** 0.2.0
- **Targets `.sets` spec at commit:** `f5a56cc` (`GUIDE.md`), `b52f182` (`README.md`)
- **Status:** draft
- **Used by:** [core source standard](source/STANDARD.md), [core target standard](target/STANDARD.md), and every language standard

**What the value model is.** The value model extends the `.sets` grammar with the meaning of the language itself. `GUIDE.md` gives a complete grammar: how domains are declared and how rosters, builders and rules are written. The value model applies that grammar to define the core ideas of `.sets`: what a value is, when two values are equal, and what truth and nothing are. It describes `.sets` and how it behaves, without regard to transpilation. It isn't concerned with other languages; the core source and target standards build on it for that.

Truth and nothing show the purpose best. The grammar can already write both outcomes of a rule: `\\ in U` always holds and `\\ in Phi` never does. The value model names them `Tautology` and `Contradiction` (VM-AX-4), and declares the one `Nothing` with `Nothing in Phi` (VM-AX-5). A language may still need its own atoms in `.sets` to control its booleans, such as `{true, false}`, but those are data. In `.sets`, real truth is a rule that holds: the state of `Tautology`.

**Every rule is an invariant.** Nothing in this document can be overridden. It defines what `.sets` text means, so every reader, adapter and standard reads the same text the same way. The core standards build on it with their own Invariant and Default rules ([standards/README.md](../README.md)); a language that seems to need a different meaning raises a Spec Issue (§5).

## Layer 2

`.sets` is built in layers, and the value model is layer 2:

1. **Grammar** (`GUIDE.md`): how `.sets` text is written.
2. **Value model** (this document): what `.sets` means as a language.
3. **Core standards** ([source](source/STANDARD.md), [target](target/STANDARD.md)): the invariants and defaults every language standard shares.
4. **Language standards**: how one language maps to and from `.sets`.

Something new starts in a language standard. If it's justified for every language, it moves into the core standards. If it's justified as part of `.sets` itself, independent of any language, it moves into the value model. Beyond that, it can become grammar, but only with a justified reason to replace an existing way of writing it with new grammar.

Numbers show the path. The grammar already allows symbols spelled like numbers; the value model says which spellings are numbers and what their operators do (§2). That may become grammar later, and then it leaves this document.

<!-- mtoc-start -->

* [Layer 2](#layer-2)
* [0. Principle: preserve the domain](#0-principle-preserve-the-domain)
* [1. Top-level domains](#1-top-level-domains)
* [2. Symbols and numbers](#2-symbols-and-numbers)
* [3. Truth and nothing](#3-truth-and-nothing)
* [4. Ranged and described subdomains](#4-ranged-and-described-subdomains)
* [5. Spec Issues](#5-spec-issues)

<!-- mtoc-end -->

## 0. Principle: preserve the domain

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

## 1. Top-level domains

**VM-KIND-2. Top-level domains are undefined domains.** A top-level axiomatic domain (`Integer`, `String`) has no members of its own in `.sets`; it's an *undefined* domain, not an infinite one. Its prose meaning (number axioms, pending; CS-STR-1 for `String`) fixes what its members are, the same in every context. A view may observe only some of them; it never changes what they are.

Inside `.sets`, top-level domains are always disjoint (GUIDE, "Top domain axiom"). A representation may merge several into one type, as JavaScript's `number` does for `Integer` and `Fractional`; the domains stay disjoint.

**Why keep the distinction in `.sets` at all.** A representation can always merge two domains, but nothing can split one that the IR doesn't separate. Keeping `Integer` and `Fractional` apart preserves the domain for every target that can use the difference ([VM-PRES-1](source/STANDARD.md#0-principles)).

This only works when the source makes the distinction. A source without it (TypeScript, whose only number type is `number`) produces one combined domain, and nothing downstream can tell which form a value was meant to be ([UP-1](#unsolved-problems)).

## 2. Symbols and numbers

`GUIDE.md` defines a symbol as a run of UTF-8 characters. This section fixes which symbols are numbers and what numeric operators do. How a symbol becomes an atom, and how escapes work, is grammar (pending in GUIDE). Quoting is a convention GUIDE leaves to "the transpiler ecosystem"; it's in the core source standard (CS-STR-2), not here.

**VM-LEX-0. Two parsing layers.** A `.sets` parser has two layers:
1. **The grammar layer** implements GUIDE only. It splits rosters at unescaped `,`, `{` and `}`, applies GUIDE's whitespace rule, and produces symbols. It gives quotes no special meaning. It keeps escape sequences in the symbol text as written; it only uses them to decide where a symbol ends.
2. **The value-model layer** implements this document and the `base` rules. It reads each symbol as a number or an atom and decodes escape sequences.

Both readings of the same text must agree. Writing `{a\, b}` satisfies both: the grammar layer sees one symbol because the comma is escaped, and the value-model layer reads it as the atom `a, b`.

An escaped separator isn't a separator. GUIDE's rule that ignores one whitespace byte around a reserved symbol applies only to the unescaped `,`, `{` and `}` of the grammar, never to `\,`, `\{` or `\}`. So in `{a\, b}` the space is part of the atom.

**VM-LEX-2. Numeric symbols are numbers.** A symbol is a number when it starts with a digit (after an optional sign), its spelling is numeric in one of the notations below, and it satisfies `x + 0 = x` (VM-AX-6). Anything else is an atom (GUIDE, pending). An implementation may check the spelling with a pattern and confirm with the identity.

| Notation | Form | Examples |
| -------- | ---- | -------- |
| Decimal | digits, optionally `.` and digits | `42`, `007`, `2.5` |
| Exponent | a decimal, `e` or `E`, then an optionally signed integer exponent; `m e k` is m × 10^k | `1e3`, `6.02e23`, `15e-1` |
| Hexadecimal | `0x` or `0X`, then hex digits (any case) | `0x1F` |
| Hexadecimal float | `0x` or `0X`, hex digits, optionally `.` and hex digits, `p` or `P`, then an optionally signed decimal exponent; `m p k` is m × 2^k | `0x1.8p1`, `0x1p-1074` |
| Octal | `0o` or `0O`, then octal digits | `0o17` |
| Binary | `0b` or `0B`, then binary digits | `0b101` |

- **Sign.** Any notation may start with `+` or `-`. `+3` and `3` are the same number; `-3` is its negative.
- **Separators.** Inside a number, `_` between two digits, or `\` before a digit, only separates: `1_000`, `1\000` and `1000` are the same number. `\` is the suggested separator, since an escaped digit is still that digit (GUIDE, pending). `_` anywhere else (`_1`, `1_`, `1__0`, `0x_1`), or an escape on a non-digit (`1\e3`), makes the symbol an atom. A symbol that starts with an escape (`\1`) is an atom.
- **The form follows the spelling.** A number is `Fractional` if it's written with a point or a negative exponent, and `Integer` otherwise, even when its fractional part is zero. `3`, `1e3`, `0x1F` and `0x1p3` are in `Integer`; `3.0`, `1.0e3`, `15e-1` and `0x1.8p1` are in `Fractional`.
- **Leading zeros** don't change a decimal: `007` is `7`, never octal; octal is always written `0o`.
- **Zero.** `-0` and `+0` are the integer `0`; `-0.0` is the fractional `0.0`. IEEE-754 negative zero is neither; it's the atom `NegativeZero` (number axioms, pending).

**Note: not every number can be written.** Every notation is finite, so a rational such as one third has no literal (UP-2). Exponents keep magnitudes short, and hexadecimal floats write every binary64 value exactly (the value nearest 0.1 is `0x1.999999999999Ap-4`). Operators appear in membership rules (VM-AX-6), not as a way of writing values.

**VM-AX-6. Numeric operators are mathematical.** Numbers are exact mathematical values with no width (number axioms, pending), so operators on them have their mathematical meaning:
- `+`, `-`, `*`, `/` give exact results, with no overflow, rounding or wraparound. `1 / 3` is the exact rational one third. Division by zero has no result.
- `a ^ n` is exponentiation with an `Integer` exponent `n`. For `n >= 0` it's defined for every number `a` (`0 ^ 0` is `1`); for `n < 0` it's `1 / a ^ (0 - n)` and needs `a` to be non-zero. An exponent that isn't an `Integer` has no result, including `2 ^ 2.0`, so results stay exact rationals: `2 ^ -2` is `0.25`.
- `%` is defined on `Integer` operands only. The result has the sign of the divisor (mathematical modulo): `-7 % 3` is `2`. A `Fractional` operand (`7.0 % 3`) or a zero divisor gives no result.
- `<`, `<=`, `>`, `>=` are the ordering of the reals. They compare mathematical values and ignore the form: `1 < 1.5` and `1 <= 1.0` both hold.
- An operator applied to a value that isn't a number (atom, tuple, domain) has no result. Nothing is coerced: `"1" + 1` doesn't exist.
- A comparison or membership test whose operand has no result is false, so the rule doesn't hold for that value. Under `not` it becomes true: `not "a" < 1` holds. That's valid, if unintuitive.

**Result form.** Every result is an `Integer` or a `Fractional` (number axioms, pending):
- `+`, `-`, `*`, `%`, and `^` with `n >= 0`, on `Integer` operands only, give an `Integer`: `2 + 3` is `5`.
- `/`, `^` with `n < 0`, and any operation with a `Fractional` operand give a `Fractional`: `4 / 2` is `2.0`, and `1.5 + 1.5` is `3.0`.

Results compare by mathematical value but are equal only within one form, so `2 + 3 <= 5.0` holds while `2 + 3 in singleton(5.0)` doesn't.

Equality stays `x in singleton(y)` (GUIDE, "Singletons"); `=` and `!=` aren't operators.

## 3. Truth and nothing

**VM-AX-4. Truth isn't a value.** Truthiness in `.sets` is whether a rule holds. Only rules have it; no value does, and no value converts to it. Two parameterless predicates name the two outcomes:

```sets
Tautology: \\ in U
Contradiction: \\ in Phi
```

`Tautology` always holds, because every symbol is in `U`; `Contradiction` never holds, because `x in Phi` never holds as a predicate (VM-AX-5). That's true in every context.

This is narrower than truthiness in many languages, where a value itself counts as true or false: C treats any non-zero integer as true, and JavaScript treats `0`, `""` and `null` as false. `.sets` can't express that kind of truthiness. It's behavior of the language, not data, so a value stays what it is: `0` is the number `0`, and an atom spelled `true` is just an atom.

**VM-AX-5. One nothing.** The core names exactly one nothing, `Nothing`, declared by a deliberate paradox:

```sets
Nothing in Phi
```

`Phi` isn't nothing: it's a domain, something that contains nothing. No value can be a member of `Phi`, so the one declared member is unlike every value, and that makes it nothing. A top-level `Nothing in Phi` is a declaration: it states that `Nothing` is in `Phi`, which makes `Nothing` what it is. It adds no member to `Phi`. Inside a rule, `x in Phi` is a predicate that never holds, as with `Contradiction` in VM-AX-4: `Phi` stays `{}`, and `{ x suchthat x in Phi }` is `Phi`.

`Nothing` is in `U`: the moment it enters `U`'s discourse, it becomes part of it.

Another top-level declaration into `Phi`, such as `Abcd in Phi`, is undefined: `Abcd` may or may not equal `Nothing`. That isn't a special case the standard adds; it follows from `Phi`. Declaring a member into a domain with no members is itself a paradox, and whether non-existence is one state or several is unknown. `Phi` may hold several nothings, each the same and yet different, because the domain of discourse keeps changing as superdomains grow. VM-LEX-3 still holds for `Abcd`, read in this context. The core doesn't decide the question, so no standard can rely on it.

There's no other core null atom. A value that may be absent-as-a-value (a nullable value) is `union(T, {Nothing})`: either a member of `T`, or `Nothing`. Other null-ish atoms (`null`, `undefined`, `nil`) are ordinary atoms; `.sets` gives them no meaning.

`()` and `{}` are conceptually equal, not literally. A tuple states the order of its components; with no components there is no order to state. If `Phi` held more than one nothing, order among them would matter, and `()` and `{}` would differ. The core deals with one nothing, so within its context the two are treated as equal and either may be written. Neither is `Nothing`.

## 4. Ranged and described subdomains

**VM-SUB-0. Ranged domains are builders.** A domain bounded by numeric comparisons is written as a builder with VM-AX-6 operators. `base.sets` defines `Natural`, `Int8`…`Int64` and `Uint8`…`Uint64` this way, for example:

```sets
Int32 = { x suchthat x in Integer and x >= -2147483648 and x <= 2147483647 }
```

Their overlaps follow from the rules: `Uint32` is inside `Int64` because every member satisfies `Int64`'s rule.

The binary floating-point domains are builders too. A finite IEEE-754 value is `m * 2 ^ e` for an integer significand `m` and exponent `e` within the format's bounds. `m` and `e` are rule-local binders, so they're existential (GUIDE, pending):

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

`0` and `0.0` are members (`m = 0`); `−0` isn't, because it's the atom `NegativeZero` (number axioms, pending). `Float64` and `Float32` add `IEEESpecial` with `union`. `Binary32Finite` is inside `Binary64Finite`, and every `Int32` is in `Binary64Finite`, both by the rules.

A top-level `X sub Y` declares a domain `X` whose members are drawn from `Y`, at most all of `Y` (GUIDE, "Difference between `sub` and `in` in this context"). Domains that the grammar still can't state are declared that way, and the prose here fixes each one's exact extent. A described subdomain with no prose extent is read as GUIDE reads it: the same as `Y`.

`base` declares one: `Char`, whose extent is in the core source standard (CS-STR-3).

**VM-SUB-1.** This rule covers the described subdomains of `base.sets` only, so it stays language-agnostic. Those subdomains overlap the way their prose says, regardless of the single `sub` parent written in `base.sets`. Prose is a stopgap for what the grammar can't state yet: once the grammar can express a described subdomain's extent, the subdomain is replaced by a builder, the way the ranged domains were (VM-SUB-0). `Char`, for example, becomes a builder when a string-length operator exists.

**When prose contradicts the grammar.** `X sub Y` is a grammar fact: every member of `X` is a member of `Y`. Prose can only narrow that, by saying which members of `Y` belong to `X`. If the prose instead names a value outside `Y` (for example, prose for `Char sub String` that admits the number `1`), the declaration says two incompatible things:
- Read by its declaration, `X` can't contain the value, because the value isn't in `Y`.
- A target that binds `X` by its prose would accept it.

The two readings then disagree about the same domain, and no representation can be correct for both, so the declaration is a defect in `base.sets`. It's fixed by correcting the prose or the `sub` parent, never by choosing one reading at bind time.

## 5. Spec Issues

Ambiguities and contradictions in `README.md` (`b52f182`) and `GUIDE.md` (`f5a56cc`) that the core standards had to resolve. Each lists the interpretation the core standards assume.

### Open

| # | Issue | Core assumption |
| - | ----- | --------------- |
| SI-18 | GUIDE's `base.sets` example lists only `u`, `phi` and the operations, and defines `U = u` and `Phi = phi` as aliases. `base` now also holds the axioms, described subdomains and structural transformers. | [`base.sets`](base.sets) is the full `base`, with `U` and `Phi` and no lowercase aliases. GUIDE's example could point to it. |
| SI-6 | GUIDE says function parameters must be domains, but `IsOne(x)` and `ToUser(x, y)` take bare parameters. | Core only uses `{X sub D}` and `{x suchthat …}` parameters. |
| SI-7 | `=`, `sub or =`, `!=` (GUIDE) and `subsetof` (README `Human`) appear but aren't defined. | Not used. `%` is now defined by VM-AX-6. |
| SI-8 | `base` shows a variadic `union({A sub u}, {B sub u}, ...)`; variadic parameters aren't defined. | Binary `union`/`intersect`; a chain nests them, and any nesting means the same (GUIDE, pending). |
| SI-12 | "`{x} in U` *can* be written as `{x} sub U`" mixes membership and subset. | In a rule, `s sub T` means every member of `s` is in `T` (used by `SetOf`, `Map`, `OnlyKeys`). Pending in GUIDE; see the notes below. |
| SI-14 | There's no grammar for cardinality, finiteness, string length or patterns. | Prose rules for described subdomains (§4) and the gaps in each standard. Ranges and floating-point domains are covered by VM-AX-6; key uniqueness is in the `Map` builder (VM-MAP-1); finiteness of a language's sets and maps is left to its standard (VM-SET-1, VM-MAP-1). |
| SI-15 | `//` comments are used throughout but aren't part of the grammar. | Allowed; they carry no meaning. |
| SI-16 | README says TypeScript's `Person \| Aged` is an "intersection". As a domain it's a union of open records whose members can satisfy both branches. | Terminology only. `union(Person, Aged)` with open records (VM-REC-3). |
| SI-19 | Is arithmetic needed in `.sets` at all? VM-AX-6 operators exist for ranges (`Int32`) and float formats (`m * 2 ^ e`), but they bring evaluation into a language of domains, and the VM-LEX-2 note already avoids literal math. | Open question for the spec author. Arithmetic stays allowed in comparisons and call arguments, and is forbidden inside tuples (GUIDE, pending). |
| SI-17 | GUIDE doesn't say how symbols that appear only in a rule (`n` in `("name", n) in r and n in String`) are bound. | Existentially over the whole rule. Pending in GUIDE; see the notes below. |

### Unsolved problems

Problems the core standards know about but don't solve yet. These are concerns of the `base` standard, not ambiguities in GUIDE.

| # | Problem | Current state |
| - | ------- | ------------- |
| UP-1 | A source language with one number type (TypeScript `number`) can't say whether a value is an integer or a fraction, so its IR has one combined number domain (for example `Binary64Finite`) instead of separate `Integer` and `Fractional` domains. A target with distinct integer and floating-point types can't choose between them correctly from the IR (VM-KIND-2, VM-PRES-1). | Unsolved. Targets use the combined domain's representation (a floating-point type). Possible directions, none adopted: source-side inference from usage, or a source standard that declares integer-only described subdomains from annotations. |
| UP-2 | Every numeric notation is finite (VM-LEX-2), so a number with no finite form in any of them (one third) can't be written. Exponents and hexadecimal floats keep other literals short. | Unsolved. Writing values through arithmetic is out of scope for now: transpilation in this iteration is non-deterministic and math must stay deterministic, so literal arithmetic is avoided wherever possible. |

### Resolved

| # | Issue | Resolution |
| - | ----- | ---------- |
| SI-1 | The universal domain was spelled both `u` and `U`, and the empty domain both `phi` and `Phi`. | `U` and `Phi` only; `base.sets` drops the lowercase aliases. |
| SI-2 | Record members were `{name, age}` in GUIDE and `(name)` in README. | Records are sets of `(key, value)` tuples: `{r suchthat ("name", n) in r and n in String}` (VM-REC-*). |
| SI-3 | There was no shared axioms superdomain. | `core` is merged into `base`. |
| SI-4 | Quoting was only a convention. | VM-LEX-2 is accepted for now; atoms and escapes are pending in GUIDE; quoted strings are CS-STR-2. |
| SI-5 | `and`/`or` had no precedence. | Boolean algebra with `()` grouping. Pending in GUIDE; see the notes below. |
| SI-9 | The empty tuple wasn't defined. | `()` is the 0-tuple. It's conceptually equal to `{}`, not literally: with no components there is no order to state, and the core treats them as equal (VM-AX-5). |
| SI-10 | Recursion wasn't defined. | Allowed. Pending in GUIDE; see the notes below. |
| SI-11 | Top-level `X sub Y` seemed undefined. | Defined in GUIDE, "Difference between `sub` and `in` in this context" (§4). |
| SI-13 | Literals inside patterns weren't defined. | Allowed. Pending in GUIDE; see the notes below. |

### Notes for the spec author

Decisions made in review that belong in `GUIDE.md`. The standards already rely on them.

1. **Rule precedence and grouping.** Rules follow boolean algebra: arithmetic, then comparisons and membership, then `not`, then `and`, then `or`, each level grouping left to right (except `^`). Parentheses after a name are a call, parentheses around a rule or arithmetic expression group it, and parentheses around values and binders are tuples, which can't contain arithmetic. GUIDE's "Tuples and Vectors" section should say that grouping parentheses aren't values.
2. **Numeric operators** (VM-AX-6). `^ * / % + -` and `< <= > >=` on numbers, with exact mathematical meaning and no coercion; `^` takes integer exponents and groups right to left. This makes ranged and floating-point domains expressible (`Int32`, `Binary64Finite` as builders), which GUIDE's number discussion could show.
3. **Rule-local binders** (SI-17). A symbol that first appears inside a rule is existential over the rule.
4. **Names.** `U` and `Phi` are the only spellings of the universal and empty domains; GUIDE's `u` and `phi` aliases go away. A declared domain or function can't be redeclared, which GUIDE's "Domain Validation" section could state.
5. **Nothing** (VM-AX-5). `Nothing in Phi` declares the one nothing, as a deliberate paradox. `Phi` is the empty domain, not nothing. `()` and `{}` are conceptually equal, not literally. GUIDE's "Tuples and Vectors" section could say so.
6. **Escapes** (VM-LEX-0). An escaped character makes a new symbol and isn't a separator. GUIDE's symbol definition could state the two parsing layers.
7. **Value kinds and equality.** Numbers, atoms, tuples and domains, and when two values or two texts are equal. Redeclaration (item 4) depends on equality: a name is redeclared when it equals one in scope.
8. **Atoms** (SI-4). A symbol that isn't a number or a domain name is an atom, quotes included; item 6 covers escapes.
9. **Recursion** (SI-10). A domain may refer to itself, and holds the values its rule builds in finitely many steps.
10. **Literals in patterns** (SI-13). A literal in a pattern matches only an equal value.
