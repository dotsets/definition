# Review of the TypeScript → `.sets` standard (pass 2)

- **Pass 1** wrote [STANDARD.md](STANDARD.md) from `VALUE-MODEL.md` 0.2.0 and `input/types.ts`.
- **Pass 2** checked every `TS-*` rule against every `VM-*` rule, then ran the standard by hand over the 13 blocks of `input/types.ts`. The fixes below are in STANDARD.md.
- There is no third pass. Problems found after pass 2 are in the iteration report, not in the standard.

## 1. Findings and fixes

| # | Rule | Finding | VM rule | Fix |
| - | ---- | ------- | ------- | --- |
| R1 | TS-MAP-1 | Pass 1 named nine prelude domains and defined none. `Double`, `Object` and the rest couldn't be checked against the VM, and the standard didn't say what the prelude must hold. | VM-SUB-0, VM-SUB-2, VM-AX-6 | Added §3 (TS-PRE-1): all 13 definitions, the two prose extents, and the `base` names used. Each definition was then checked (§2 below). |
| R2 | TS-OUT-* | No rule said how an output file gets prelude names. Only a limit mentioned the form. | VM-NAME-1 | Added TS-OUT-4: one `X in prelude` line per name. |
| R3 | TS-OUT-3, TS-UNION-2 | The `Overlap:` note fired when "two branches of a union share values". That's normal for open records (VM-REC-3) and isn't what VM-WIRE-1 flags. It also missed two declared types that become one domain (`One` and `OneN`, `Dictionary` and `Scores`). | VM-WIRE-1, VM-PRES-1 | The note now means exactly a TS-WIRE-1 case: one type admits both source values, or two declared types differ only by such a case. TS-UNION-2 no longer mentions overlaps. |
| R4 | TS-WIRE-1 | The table missed collisions the standard itself creates: all functions of one type (TS-MAP-1), an invalid `Date` and `NaN`, and `0n` with `new Date(0)`. | VM-WIRE-1 | Rows added. The rule now says the adapter reports the case and writes the note. |
| R5 | TS-LIB-1 | The reason for the `SetOf` reading was one-sided: it named the collision of the other reading and not its own. The VM-PRES-2 example uses the other reading. | VM-SET-1, VM-PRES-2 | Rewritten: both readings collide, `SetOf` collides less, and the VM-PRES-2 example value is named as the alternative. |
| R6 | TS-NIL-1 | `unset` was called "a second nothing", which reads like a second member of `Phi`. VM-AX-5 leaves that undefined. | VM-AX-5 | Reworded: an ordinary atom in the roster `Unset`, not a member of `Phi`. |
| R7 | TS-STR-1 | The rule showed three examples and no list of what to escape. `{` and `}` weren't shown at all. | VM-LEX-0, VM-LEX-1, VM-LEX-6 | The five escaped characters and the control-character forms are listed. Keys and enum values use the same rule. |
| R8 | TS-TPL-1 | "Says what may stand in each gap" isn't precise enough to write a check from for gaps such as `${number}`. The standard also didn't say how its described subdomains overlap. | VM-SUB-2 | Only `string` gaps are covered, with the exact wording of the extent. Other gaps moved to TS-SCOPE-4. §3 says the described subdomains overlap as their extents say. |
| R9 | TS-TUP-1 | The rule departs from the wording of two Invariant rules without saying how it reads them. A tuple with both optional and rest elements was unspecified. | VM-SEQ-1, VM-TUP-1 | The reading is stated: a VM tuple would give one array two `.sets` values, so it's not an override (see O1). Added: only the longest shape has the tail. |
| R10 | TS-SCOPE-3 | Pass 1 said "diagnostic" in two rules and "reported" in two others, and defined neither. | VM-WIRE-1 | Defined once in TS-SCOPE-3. TS-WIRE-2 now uses the same word. |
| R11 | TS-OUT-2 | The rule didn't say how a declaration is written. | — | Added: each declaration is `X = …`, a generic one a function. |
| R12 | TS-SCOPE-4 | `Promise` was "not covered" although VM-PRES-3 names its domain. | VM-PRES-3 | The reason is given: `Async(T)` is declared in `RUNTIME.md`, which this iteration didn't read. |
| R13 | §11 | Four limits were missing: a function is an `object` but isn't in `Composite`; the broad domains are exact one level deep; a default type argument is only a note; `Callable` is this standard's own. The unjustified collision had no entry. | VM-PRES-1, VM-PRES-3 | L2 extended, L3, L5, L6 and L9 added, L7 extended. |
| R14 | TS-OBJ-4 | Nothing told a target that a `FiniteMap` call and a builder are different shapes on purpose (iteration 2, N7). | VM-REC-1, VM-MAP-1 | One sentence added: a call is a table, a builder names its fields. |
| R15 | TS-NUM-1 | The rule said a whole number from 2^52 up has no `Fractional` form "in `Binary64Finite`", but the VM prose says every representable whole number has both forms. | VM-SUB-0 | The rule cites the builder, which is the definition (see O3). |
| R16 | TS-MAP-1, TS-FN-1 | "A placeholder" had no spelling, and the standard didn't say a runtime value is never checked. | VM-PRES-3 | The placeholder is `(callable, P, R)`, and TS-FN-1 says it's never serialized or checked as data. |
| R17 | TS-REC-1 | The tag reads as "record". | — | Renamed TS-RECUR-1. |
| R18 | TS-OUT-1 | "Nothing names the source language" could be read as a rule to rename a declared type such as `Json`. | VM-NAME-1 | Added: a declared type keeps its own name. |

## 2. Checked, no change

| VM rule | What was checked | Result |
| ------- | ---------------- | ------ |
| VM-PRES-1 | TS-OBJ-5 drops `readonly` and the access modifiers without a note. | They don't change which values the type holds, so there's no domain distinction to keep. |
| VM-KIND-2, UP-1 | `number` gives one combined domain. | `Double` holds both forms' worth of values, each once. |
| VM-LEX-2, VM-LEX-2a | Number spellings. | Canonical decimal only. No notation is added, so VM-LEX-2a isn't overridden. |
| VM-RULE-1 | `WholeFractional` and `EpochMillis`. | Ranges use `and`; there's no chained comparison and no unary sign on a binder. |
| VM-RULE-2 | Bare symbols in rules. | TS-OBJ-2 keeps atoms in rosters (`ok in {true}`). `Double` negates a named domain (`x not in WholeFractional`), so no binder sits under `not`. |
| VM-RULE-3 | `()` as a literal in a sequence shape. | A literal matches an equal value; `()` and `{}` are treated as equal. |
| VM-LEX-4 | `P` and `R` inside the roster of `Callable`. | They're domain parameters, so they stand for the domains themselves, which is what the tag needs. |
| VM-AX-3, VM-AX-4 | Non-finite numbers and booleans. | Core atoms are used as they are. |
| VM-AX-7, VM-AX-8 | New axioms. | None. `EpochMillis` is a builder over `Integer`; `Callable` is a roster. |
| VM-SUB-2 | `Finite sub U`, `NumericKey sub String`. | Each extent only narrows its parent. |
| VM-REC-2 (override) | `Object` in place of `Record`. | `FiniteMap(String, U)` is inside `Map(U, U)`. |
| VM-REC-3, VM-REC-4 | Openness. | Every record is open. A string index signature names every key, which a builder may do. `OnlyKeys` is never written. |
| VM-MISS-1, VM-MISS-2 | Optional and nullable properties. | Missing, `unset` and `Nothing` are three records. `OptionalField` and `OptionalFields` are used as `base` defines them. |
| VM-SET-1, VM-MAP-1 | Finiteness. | `FiniteSet` and `FiniteMap` narrow with `Finite`. |
| VM-RECUR-1 | Cyclic values. | No reference encoding is defined; a cyclic value is a diagnostic. |
| VM-EQ-1 | One value per source value. | The projection is a function: one number form, one array encoding. |
| SI-8 | `union` and `intersect`. | Binary, nested to the right. |
| Dry run | All 13 blocks of `input/types.ts`. | Every declaration has exactly one rule path. The results are `sets/types.sets`. |

## 3. Left open

These need a decision outside this standard. They're repeated in the iteration report.

| # | Problem | VM rule |
| - | ------- | ------- |
| O1 | VM-TUP-1 says language tuples are tuples. A TypeScript tuple is an array, so one array would get two encodings. TS-TUP-1 reads the two Invariant rules in the only way that keeps the projection a function. The VM wording should say which kind of language tuple it means. | VM-TUP-1, VM-SEQ-1 |
| O2 | A set of non-empty sequences with distinct first elements is equal to a map. No language standard can avoid this with the `SetOf` reading. | VM-SEQ-1, VM-SET-1, VM-MAP-1 |
| O3 | §5 "Both forms" says a float domain holds both forms of every representable whole number. By the builder, the `Fractional` form exists only below 2^52 (iteration 2, F12). | VM-SUB-0 |
| O4 | A bare atom can't be written as a literal in a pattern, because a bare symbol in a rule is a binder. VM-RULE-3 shows strings, numbers and domain names, and no atom. | VM-RULE-2, VM-RULE-3 |
| O5 | The VM points to `base.sets` and `RUNTIME.md` for names it relies on. This iteration couldn't read either (iteration 2, F4). | VM-PRES-3, VM-NAME-1 |
| O6 | The VM mentions the import form `X in file` once and doesn't define it. | VM-NAME-1 |
| O7 | The VM-PRES-2 example gives `new Set(["a", "b"])` the value `{("a", ()), ("b", ())}`. Under this standard it's `{"a", "b"}`, by the choice VM-SET-1 allows. | VM-PRES-2, VM-SET-1 |
