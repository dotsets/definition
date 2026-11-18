# Self-review of STANDARD.md (pass 1 → pass 2)

Pass 1 was 0.1.0. Each finding below was checked against the value model only, and fixed in 0.2.0. No third pass was made.

| # | Finding in 0.1.0 | VM rule | Fix in 0.2.0 |
| - | ---------------- | ------- | ------------ |
| R1 | `bigint` → `Integer` and whole `number`s → `Integer` form make `1n` and `1` the same `.sets` value. JavaScript tells them apart (`typeof`, `1n + 1` throws), so this is a collision between different data: a defect. | VM-PRES-1, "Collisions" | Every finite `number` is written in `Fractional` form (`1` → `1.0`). `bigint` keeps `Integer`. `TsNumber` becomes `Float64` minus `Integer`. |
| R2 | TS-OBJ-1 says keys are strings, but `r in Record` admits any key (`Record = Map(U, U)`), so `{(1, "a")}` was a `User`, though no JS object has a number key. Numeric keys (`{1: "a"}`, `{[k: number]: V}`) weren't covered at all. | VM-REC-1, VM-PRES-1 | `r in TsRecord` with `TsRecord = Map(TsKey, U)`, as an override of VM-REC-2 (Default). Numeric keys project to their string form; `[k: number]` uses `TsNumericKey sub String`. |
| R3 | `TsObject` and the `{}` type were named but never defined. | — | Both defined in `ts.sets` as complements. |
| R4 | TS-LIT-1 example `1` → `{1}` contradicts R1. | TS-NUM-1 | `1` → `{1.0}`. |
| R5 | `[A, B?]` didn't allow a present `undefined` in position 2, unlike TS-OBJ-2. | TS-VIEW-1 | Optional elements take `union(B, TsUndefined)`. |
| R6 | A union of a fixed tuple and an array (`[number, number] \| number[]`) gives one JS value two `.sets` values (`(1.0, 2.0)` and `(1.0, (2.0, ()))`). VM-SEQ-1 makes the encoding depend on the type, and pass 1 didn't say what to do. | VM-SEQ-1, VM-WIRE-1 | TS-WIRE-1: the projection reports a diagnostic when a value has two encodings in the domain. |
| R7 | `TsDate` used `Integer` milliseconds, but `Date` wraps a `number`, which is now `Fractional`. | VM-PRES-2, R1 | `TsDate` is whole-valued `TsNumber`s in range, plus `NotANumber`. |
| R8 | Symbol-keyed properties weren't mentioned. | — | Out of scope (TS-SCOPE-1). |
| R9 | Overrides table said "None" though R2 overrides a Default. | VM header | Overrides table lists VM-REC-2. |
