# Findings: TypeScript → `.sets` → Go (iteration 2)

## Pipeline

```text
input/types.ts ──ts-to-sets/STANDARD.md──▶ sets/types.sets + sets/ts.sets ──sets-to-go/STANDARD.md──▶ go/*.go
                  (VM + TS knowledge)        (the only hand-off)              (VM + .sets only)
```

- Both standards use `VALUE-MODEL.md` 0.2.0 and nothing else. Each had exactly two passes (see each `REVIEW.md`).
- The Go output builds, passes `go vet`, and passes 8 vector tests (`go/model_test.go`, vectors G-01…G-08).

**Caveat on isolation.** The rules enforce it. The Go standard cites only the VM and `.sets` text, and the Go generation read only `sets/`. But the same author wrote both sides in one session, so the Go standard may still be shaped by knowing what the IR would contain. A strict test runs the Go side in a separate agent that has only `VALUE-MODEL.md` and `sets/`.

## 1. Iteration-1 findings, re-checked

| it1 | Problem | Status in it2 | Evidence |
| --- | ------- | ------------- | -------- |
| F1 | Tuple vs. array: one JS value, two encodings | **Fixed.** A JS array is always a sequence, and tuple types are sequence patterns (TS-VAL-6). VM-SEQ-1 allows this: "the source's type decides", and the runtime type is a variable-length array. | `Ambiguous` normalizes to `[]TsNumber` in Go (GO-UNION-0); E-13 |
| F2 | `Set`/`Map`/`Date` collide with other data | **Partly fixed.** `Set` → `SetOf` no longer collides. `Map` vs. object and `Date` vs. number are now *declared* justified (same data, different representation) and reported by TS-WIRE-1. | E-10, E-11, E-12 |
| F3 | `bigint` vs. `number` | **Re-decided.** `1n` and `1` are the same value, justified by VM-AX-2 (a value has no representation). The domains stay distinct (`*big.Int` vs. `float64` in Go). | E-15, T-03 |
| F4 | VM not self-contained | **Open, and worse.** Both standards now depend on assumed `base` names. The Go standard's GO-BASE-1 table is a guess at `base.sets`: `Natural`, `Char`, `Float32`, `Real`, `OnlyKeys` and `OptionalFields` were never seen. | GO-BASE-1, prelude header |
| F5 | Exact decimals for literals | **Works end to end.** The 55-digit `.sets` literal for TS `0.1` becomes a Go constant that rounds back to exactly `0.1`. | G-01 |
| F6 | Top types bounded by `U` | **Fixed in `.sets`** (`TsValue`), but **invisible in Go**: every top type becomes `any`. | T-30…T-33 |
| F7 | Record keys | Still needs the VM-REC-2 override. | `TsRecord` |
| F8 | Lone surrogates | Specified: a projection diagnostic (TS-VAL-1). Not exercised. | P-42 |
| F9 | `{()}` roster | Avoided, because `Set` no longer uses it. | — |
| F10 | Recursion through functions | Works: `Tree[T]` in Go. | G-08 |
| F11 | Default type parameters | Still lost. | `Result` |
| F12 | No `Fractional` form above 2^52 | **Fixed** by writing whole numbers in `Integer` form. `1e21` and the largest date are exact in `.sets` and in Go. The VM text in §5 "Both forms" is still wrong and should be corrected. | G-02, M-07b, M-17 |

## 2. New findings, most severe first

| # | Finding | Where | VM rules |
| - | ------- | ----- | -------- |
| N1 | **Structure survives; constraints don't.** Almost every Go type is *wider*: `TsNumber`, `TsDate`, `Status`, `Level`, `UserKey` and `TsNumericKey` accept any value of their base Go type. Rules such as ranges, rosters, the whole-number restriction and prose extents are carried to Go only as comments. Membership is the core of `.sets`, and it currently stops at the IR. The next step is for the Go standard to generate validators from builder rules (VM-AX-6 operators and literal patterns are mechanical enough). Prose extents (`UserKey`, `TsNumericKey`) can't be validated at all. | GO-EXACT, all wider rows in VECTORS.md §5 | VM-SUB-1, VM-SUB-2, SI-14 |
| N2 | **Other kinds of absence are opaque to a target.** `age?: number` became `Opt[Either[TsNumber, TsUndefined]]`. That's correct, and all three states are kept (G-04). But Go can't know that `undefined` is absence-like: VM-MISS-1 leaves "kinds of absence in between" to the language standard, so in the IR it's just an atom. The natural `.sets` spelling, `undefined in Phi`, is exactly what VM-AX-5 calls undefined. The VM needs a way to declare a second nothing that targets can recognize. | `User.Age`, `Span.F1` | VM-MISS-1, VM-AX-5 |
| N3 | **Targets need set algebra the VM doesn't give.** `EmptyThings` is the same domain as `TsValue` (Go `any`), but Go produced `*Either[TsNonNullish, TsUndefined]`: the Go standard can't prove `TsNonNullish ∪ {Nothing} ∪ TsUndefined = TsValue` through `not in` refinements. GO-UNION-0 recognizes only "visible" containment. The VM defines meaning but no normal form or simplification rules, so each target reinvents an incomplete subset check, and different targets will disagree. | `EmptyThings`, GO-UNION-0 | VM-RULE-1, VM-RULE-4 |
| N4 | **Open records cost every target an `Extra` field.** VM-REC-3 (open by default) together with VM-WIRE-1 (one to one) forces every Go struct to carry `Extra map[TsKey]any`, or else `{name}` and `{name, x}` collapse. The source said nothing about extra keys; it's TS's structural typing. That's correct, but heavy. A source with a real closedness signal would emit `OnlyKeys`; TS has none. | every struct | VM-REC-3, VM-REC-4, VM-WIRE-1 |
| N5 | **Prose carries source semantics across the boundary.** `TsNumericKey`'s extent says "JavaScript's `String(n)`". A target binds described subdomains by prose (VM §5), so it has to understand the source language to bind exactly. The prose channel breaks the "no knowledge of the source" goal. Names leak too (`TsNumber` in Go), though harmlessly. | `TsNumericKey`, `UserKey` | VM-SUB-1, VM-SUB-2 |
| N6 | **Rosters partly solve UP-1.** `Level = {0, 1, 2}` and `One = {1}` became `int8` in Go, because whole numbers are in `Integer` form and the Go standard picks the narrowest integer type. A general `number` stays `float64`. So the integer/float split reaches the target only for literal and enum domains. That's a real improvement over iteration 1, where everything was `Fractional`. | T-05, T-07 | VM-AX-2, UP-1 |
| N7 | **The record/map boundary is a heuristic.** VM says "a record is a map". GO-REC-1 has to decide struct vs. map by whether literal field patterns exist, and decide record constraint vs. index signature by whether the value domain is `U`. `{ r suchthat r in TsRecord }` with no fields becomes a struct holding only `Extra`, though `map` would be the natural choice. | GO-REC-1, `Headers` | VM-REC-1, VM-MAP-1 |
| N8 | **Overlapping unions need a first-branch rule.** `SetOrDict` holds `{}` in both branches (empty set = empty map, VM-REC-1), so `Either` picks branch A. This is a deterministic tie-break, not one to one (VM-WIRE-1 diagnostic). | `SetOrDict`, GO-UNION-6/7 | VM-REC-1, VM-WIRE-1 |
| N9 | **Generic bounds disappear.** `{T sub TsObject}` becomes `[T any]`. Go constraints can't express domain membership. | `Page` | — |
| N10 | **Variant names come from literals, not meaning.** `ResultTrue`/`ResultFalse`, `ShapeCircle`/`ShapeRect`. They're mechanical and stable, but `ok: true` doesn't become `Ok`. | GO-UNION-5 | VM-RULE-3 |

## 3. What a TS distinction became in Go

| TypeScript distinction | Survives to Go as | |
| ---------------------- | ----------------- | - |
| missing / `null` / `undefined` | `Opt.Present` / `*T == nil` / `Either.IsB` | kept |
| `number` / `bigint` | `float64` / `*big.Int` | kept (by domain) |
| `[number, number]` / `number[]` | `[2]TsNumber` / `[]TsNumber` | kept |
| `Set` / `Map` / object | `map[T]struct{}` / `map[K]V` / struct | kept (by domain) |
| discriminated union | sealed interface | kept |
| open object types | `Extra` field | kept |
| literal and enum types | named type + constants | type kept, membership not checked (N1) |
| `Date` | `TsDate` (float64) | as a number only |
| template literal type | `string` + comment | prose only (N1, N5) |
| generic constraint, default type argument | — | lost (N9, F11) |
| `readonly`, methods | — | dropped in the source standard |
| function-typed field | opaque `TsFunction` handle | kept as a placeholder |

## 4. Suggested VM work for iteration 3

1. Fix §5 "Both forms": the `Fractional` form exists only below 2^52 (F12).
2. A recognizable way to declare another kind of absence (N2).
3. A small set of normalization rules (flattening, absorption, `not in` against a superset) so targets agree on union simplification (N3).
4. Make `base.sets` part of what the VM ships, or inline the definitions the VM relies on (F4).
5. Decide whether validators are a target-standard obligation (N1).
