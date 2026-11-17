# Stress-test findings: value model vs. TypeScript (iteration 1)

Input: [examples/types.ts](examples/types.ts). Output: [examples/types.sets](examples/types.sets). Rules: [STANDARD.md](STANDARD.md) 0.2.0 and `VALUE-MODEL.md` 0.2.0 only.

## Summary

Most of TypeScript mapped directly, using only VM rules: open records (VM-REC-3), required/optional/nullable fields (VM-REC-2, VM-MISS-1/2, VM-AX-5), literal-tagged unions (VM-RULE-3), generics as functions, recursion, and ranged builders (`TsDate`). The value model held up on the data-shape core.

The problems are at the edges: values whose `.sets` encoding depends on the static type, collisions the VM approves that the source can tell apart, and the VM not being self-contained.

## Findings, most severe first

| # | Finding | Where it shows | VM rules |
| - | ------- | -------------- | -------- |
| F1 | **A value's encoding depends on its type.** A JS array `[1, 2]` is `(1.0, 2.0)` under `[number, number]` and `(1.0, (2.0, ()))` under `number[]`. In a union (`Ambiguous`), one runtime value has two `.sets` values. VM-SEQ-1 is Invariant ("the source's type decides which"), so a source standard can't avoid this. It can only report it (TS-WIRE-1). This is the first place where the value model isn't a model of *values*. | `Ambiguous`, `Point` vs. `Seq` | VM-SEQ-1, VM-TUP-1, VM-WIRE-1 |
| F2 | **The VM's blessed wrapper collisions are collisions between different data.** VM-PRES-2 tells TS to read `Set<T>` as a record of keys with empty values. Then `new Set(["a"])` and `{a: {}}` are the same value, `{("a", ())}`, because `()` and `{}` are treated as equal (VM-AX-5). `new Map([["a", 1]])` equals `{a: 1}`, and `new Date(0)` equals `0`. JS tells each pair apart (`instanceof`, `typeof`), so by the VM's own "Collisions" test these are defects, yet VM-PRES-2 recommends them. The two rules contradict each other for any source with runtime type tags. | `SetOrDict`, `TagSet`, `TsDate` | VM-PRES-2, "Collisions", VM-AX-5 |
| F12 | **VM §5 "Both forms" is false for large whole numbers.** The VM says the float domains hold both forms of every representable whole number. A `Fractional` member needs `e < 0` in `m * 2 ^ e`, so `m = x * 2 ^ (0 - e) >= 2x`. Once `abs(x) >= 2 ^ 52`, `m` exceeds `2 ^ 53` and only the `Integer` form exists. With TS-NUM-1's `Fractional`-only numbers, `TsNumber` loses every whole number from 2^52 on, including `1e21` and the upper half of `TsDate`'s range. Found while writing [VECTORS.md](VECTORS.md), so the standard wasn't changed (two-pass limit). Either the VM wording or the `Binary64Finite` builder needs fixing, and TS-NUM-1 must then pick a form for large whole numbers. | `TsNumber`, `TsDate`; vectors P-13, M-07b, M-17 | VM-SUB-0 "Both forms", VM-AX-6 result form |
| F3 | **One number type plus a separate bigint forces a choice of form.** Pass 1 wrote whole `number`s as `Integer`, so `1` and `1n` collided (REVIEW R1). The fix (every `number` in `Fractional` form) works, but it needs `TsNumber = Float64 minus Integer`, because `Binary64Finite` holds both forms of every whole number. Any source with a single number type has to subtract one form. A base domain like `Binary64Finite` restricted to `Fractional` would remove this. It relates to UP-1. | `TsNumber`, `One` vs. `OneN` | VM-AX-1, VM-AX-2, VM-SUB-0 "Both forms", UP-1 |
| F4 | **The VM isn't self-contained.** It uses, but doesn't define: `Float64` (only "adds `IEEESpecial` with `union`"), `Map`, `Seq`, `SetOf`, `OnlyKeys`, `product2`, `singleton`, `intersect`, and the import form `X in file`. It also cites CS-SCOPE-1 and CS-CONV-4 from a standard this test couldn't read. The transpilation had to assume what those names mean, and has no import line. | header of `types.sets`, `ts.sets` | §1 intro, VM-NAME-1, VM-MAP-1, SI-8 |
| F5 | **Literal values aren't the literals the author wrote.** VM-AX-2 makes values exact, so `type Tenth = 0.1` becomes a 55-digit decimal, and `3.14` becomes 52 digits. This is correct, and the VM-LEX-2 note already warns about it, but it means source text like `0.1` can't be round-tripped readably. A target reading `{0.1}` as "the binary64 nearest 0.1" would be wrong by the VM: `0.1` isn't in `Binary64Finite` at all. | `Pi`, `Tenth` | VM-AX-2, VM-LEX-2 note, UP-2 |
| F6 | **Complements are taken against `U`, which is bigger than any language.** `object`, `{}`, `unknown` map to `U` or `U` minus something, so they admit `.sets` values no JS program can produce: arbitrary atoms such as `Circle`, top-level domains such as `Integer`, any tuple. A `TsValue` domain can't be written to bound them, because there is no "tuple of any arity" domain. | `TsObject`, `TsNonNullish`, `Page`'s bound | VM-KIND-1, VM-TUP-1 |
| F7 | **Open records admit keys the language can't have.** `r in Record` accepts `{(1, "a")}`, which isn't a JS object. The standard had to override VM-REC-2 with `TsRecord = Map(String, U)`. This is probably common to every language with string-only keys (JSON, JS, Python kwargs). It might belong in the VM as an option instead of each standard overriding. Also, `{ 1: "a" }` is `{("1", "a")}` in JS, while VM-REC-1's example `{ 1: "abc" } is {(1, "abc")}` reads like a language-neutral fact. | `TsRecord`, `ByIndex` | VM-REC-1, VM-REC-2 |
| F8 | **JS strings can hold lone surrogates; VM `String` can't.** `"\uD800"` is a valid JS string but not a sequence of Unicode scalar values. The standard doesn't address this. It was noticed after the two passes, so it's left as a gap. The options are a String meaning override (VM-KIND-2 allows it) or a projection diagnostic. | not exercised | §4 `String`, VM-KIND-2, VM-PRES-1 |
| F9 | **Roster of the empty tuple.** `Map(T, {()})` needs the domain whose only member is the empty value. VM-RULE-1a defines parentheses inside rules, not inside rosters, and VM-LEX-0's grammar layer only knows `,`, `{` and `}`. So it's unclear whether `{()}` holds the empty tuple or the atom spelled `()`. `{Phi}` is the safer spelling. | `TagSet` | VM-LEX-0, VM-RULE-1a, VM-AX-5 |
| F10 | **Recursion through parameterized functions isn't stated.** VM-RECUR-1 talks about domains referring to themselves. `Tree(T)` refers to `Tree(T)` inside a function body. The least-fixed-point reading clearly extends to it, per argument, but the rule doesn't say so. | `Tree` | VM-RECUR-1 |
| F11 | **No default parameters.** `Result<T, E = string>` can't keep its default, and VM-NAME-1 rules out a second `Result`. Every use site has to fill in the default. That's a small loss of source intent, not of data. | `Result` | VM-NAME-1 |

## What worked without friction

- **Absence.** `age?: number`, `email: string | null` and a present `undefined` stay three distinct states: a missing tuple, `Nothing`, and the atom `undefined` (VM-MISS-1, VM-AX-5).
- **Discriminated unions.** Literal patterns `("kind", "circle") in r` read naturally (VM-RULE-3).
- **Structural typing.** TS's open object types match VM-REC-3 exactly, and `{} | null | undefined` works out to `U`, as it does in TS.
- **Enums.** Dropping member names is a justified collision: a string enum value is a string at runtime.
- **Rest tuples.** `[string, ...number[]]` as `(head, rest)` with `rest in Seq(…)` is a valid sequence, so it agrees with how the same array is encoded as `(string | number)[]`. Only fixed tuples have F1.
- **Ranged builders.** `TsDate` (whole, in range, or `NotANumber`) was written with VM-AX-6 operators and an existential binder, without any prose.

## Test vectors

See [VECTORS.md](VECTORS.md). That file is the baseline for later iterations.
