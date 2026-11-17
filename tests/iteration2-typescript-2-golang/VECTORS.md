# Test vectors (iteration 2: TypeScript → `.sets` → Go)

These vectors are the baseline for iteration 3. They carry over iteration 1's IDs (`tests/iteration1-TypeScript/VECTORS.md`) and add Go vectors.

- **Domains:** [sets/types.sets](sets/types.sets) and [sets/ts.sets](sets/ts.sets). **Go:** [go/](go/).
- **Notation:** `.sets` value notation (VM §9).
- **IDs** are stable. An ID is never reused for a different input.
- **Basis** says which rules fix the result:
  - **VM**: value-model rules alone.
  - **TS**: a choice of `ts-to-sets/STANDARD.md`.
  - **GO**: a choice of `sets-to-go/STANDARD.md`.
  - **Issue** `F#` (iteration 1) or `N#` (iteration 2): the result shows an open problem in [FINDINGS.md](FINDINGS.md).
- **vs. it1** is `same`, `changed (was …)`, `new`, or `dropped`.

## 1. Projection: JavaScript value → `.sets` value

| ID | JavaScript value (static type) | `.sets` value | Basis | vs. it1 |
| -- | ------------------------------ | ------------- | ----- | ------- |
| P-01 | `"abc"` | `"abc"` | VM-LEX-1 | same |
| P-02 | `""` | `""` | VM-LEX-1 | same |
| P-03 | `"a, b"` | `"a\, b"` | VM-LEX-0 | same |
| P-04 | `'say "hi"'` | `"say \"hi\""` | VM-LEX-1 | same |
| P-05 | `"x\n"` | `"x\n"` | VM-LEX-6 | same |
| P-06 | `"{a}"` | `"\{a\}"` | VM-LEX-1 | same |
| P-07 | `1` | `1` | TS-VAL-2 | changed (was `1.0`) |
| P-08 | `0` | `0` | TS-VAL-2 | changed (was `0.0`) |
| P-09 | `-0` | `NegativeZero` | VM-AX-3 | same |
| P-10 | `0.5` | `0.5` | TS-VAL-2 | same |
| P-11 | `0.1` | `0.1000000000000000055511151231257827021181583404541015625` | VM-AX-2 | same |
| P-12 | `3.14` | `3.140000000000000124344978758017532527446746826171875` | VM-AX-2 | same |
| P-13 | `1e21` | `1000000000000000000000` (in `TsNumber`) | TS-VAL-2 | changed (was `…000.0`, not a member; F12 fixed) |
| P-14 | `NaN` | `NotANumber` | VM-AX-3 | same |
| P-15 | `Infinity`, `-Infinity` | `PositiveInfinity`, `NegativeInfinity` | VM-AX-3 | same |
| P-16 | `1n` | `1` | TS-VAL-3 | same |
| P-17 | `-7n` | `-7` | TS-VAL-3 | same |
| P-18 | `true`, `false` | `true`, `false` | VM-AX-4 | same |
| P-19 | `null` | `Nothing` | VM-AX-5 | same |
| P-20 | `undefined` | `undefined` | TS-VAL-4 | same |
| P-21 | `{}` | `{}` | VM-REC-1 | same |
| P-22 | `{ name: "a", age: 1 }` | `{("name", "a"), ("age", 1)}` | VM-REC-1 | changed (was `1.0`) |
| P-23 | `{ 1: "a" }` | `{("1", "a")}` | TS-VAL-5 | same |
| P-24 | `{ age: undefined }` | `{("age", undefined)}` | VM-MISS-1 | same |
| P-25 | `{ age: null }` | `{("age", Nothing)}` | VM-MISS-1 | same |
| P-26 | `[]` (`string[]`) | `()` | VM-SEQ-1 | same |
| P-27 | `["a", "b"]` (`string[]`) | `("a", ("b", ()))` | VM-SEQ-1 | same |
| P-28 | `[1, 2]` (`[number, number]`) | `(1, (2, ()))` | TS-VAL-6 | changed (was `(1.0, 2.0)`; F1 fixed) |
| P-29 | `[1, 2]` (`number[]`) | `(1, (2, ()))` | VM-SEQ-1 | changed (was `(1.0, (2.0, ()))`) |
| P-30 | `[1]` (`[number, number?]`) | `(1, ())` | TS-TUP-1 | changed (was `(1.0)`) |
| P-31 | `["a", 1, 2]` (`[string, ...number[]]`) | `("a", (1, (2, ())))` | TS-TUP-1 | changed (number form) |
| P-32 | `["a", 1, 2]` (`(string \| number)[]`) | `("a", (1, (2, ())))` | VM-SEQ-1 | changed (number form) |
| P-33 | `[[1]]` (`number[][]`) | `((1, ()), ())` | VM-SEQ-1 | changed (number form) |
| P-34 | `new Set(["a", "b"])` | `{"a", "b"}` | TS-VAL-7 | changed (was a record of keys; F2) |
| P-35 | `new Set()` | `{}` | VM-SET-1 | same |
| P-36 | `new Map([["a", 1]])` | `{("a", 1)}` | TS-VAL-7 | changed (number form) |
| P-37 | `new Map([[1, "a"]])` | `{(1, "a")}` | TS-VAL-7 | changed (number form) |
| P-38 | `new Date(0)` | `0` | TS-VAL-7 | changed (was `0.0`) |
| P-39 | `new Date(NaN)` | `NotANumber` | TS-VAL-7 | same |
| P-40 | `Color.Red` | `"RED"` | TS-ENUM-1 | same |
| P-41 | `Level.Mid` | `1` | TS-ENUM-1 | changed (was `1.0`) |
| P-42 | `"\uD800"` | no value; diagnostic | TS-VAL-1 | changed (was unspecified; F8) |
| P-43 | `new Set([[1, 2]])` | `{(1, (2, ()))}` | TS-VAL-6, TS-VAL-7 | new |
| P-44 | `4503599627370496` (2^52) | `4503599627370496` | TS-VAL-2 | new |

## 2. Equality

| ID | Left | Right | Equal? | Basis | vs. it1 |
| -- | ---- | ----- | ------ | ----- | ------- |
| E-01 | `1` | `1.0` | no | VM-EQ-1 | same |
| E-02 | `1.5` | `1.50` | yes | VM-LEX-2 | same |
| E-03 | `"1"` | `1` | no | VM-KIND-1 | same |
| E-04 | `true` | `"true"` | no | VM-LEX-3 | same |
| E-05 | `()` | `{}` | yes | VM-AX-5 | same |
| E-06 | `Nothing` | `{}` | no | VM-AX-5 | same |
| E-07 | `Nothing` | `undefined` | no | TS-VAL-4 | same |
| E-08 | `(1)` | `(1, ())` | no | VM-EQ-1 | same (number form) |
| E-09 | `{("a", 1), ("b", 2)}` | `{("b", 2), ("a", 1)}` | yes (order is lost) | VM-EQ-1 | same |
| E-10 | `new Set(["a"])` → `{"a"}` | `{ a: {} }` → `{("a", {})}` | no | TS-VAL-7 | changed (was yes; F2 fixed) |
| E-11 | `new Map([["a", 1]])` | `{ a: 1 }` | yes, justified; TS-WIRE-1 diagnostic | TS-VAL-7 | same result, now declared justified |
| E-12 | `new Date(0)` | `0` | yes, justified; TS-WIRE-1 diagnostic | TS-VAL-7 | same result, now declared justified |
| E-13 | P-28 | P-29 | yes | TS-VAL-6 | changed (was no; F1 fixed) |
| E-14 | P-31 | P-32 | yes | VM-SEQ-1 | same |
| E-15 | `1n` → `1` | `1` → `1` | yes, justified; TS-WIRE-1 diagnostic | TS-VAL-3 | changed (was no) |
| E-16 | `new Set()` | `new Map()`, `{}`, `[]` | yes (all `{}`) | VM-REC-1 | new |

## 3. Membership

### 3.1 Prelude

| ID | Domain | Value | Member? | Basis | vs. it1 |
| -- | ------ | ----- | ------- | ----- | ------- |
| M-01 | `TsNumber` | `1` | yes | TS-VAL-2 | changed (was `1.0` yes) |
| M-01a | `TsNumber` | `1.0` | no (whole numbers use `Integer` form) | TS-VAL-2 | new |
| M-02 | `TsNumber` | `0.5` | yes | VM-SUB-0 | same |
| M-03 | `TsNumber` | `NotANumber` | yes | VM-AX-3 | same |
| M-04 | `TsNumber` | `NegativeZero` | yes | VM-AX-3 | same |
| M-05 | `TsNumber` | `"1"` | no | VM-KIND-1 | changed (old M-05 is now M-01a) |
| M-06 | `TsNumber` | `0.1` | no (not binary64) | VM-SUB-0 | same |
| M-07 | `TsNumber` | `9007199254740993` | no | VM-SUB-0 | changed (number form) |
| M-07a | `TsNumber` | `4503599627370495` | yes | VM-SUB-0 | changed (number form) |
| M-07b | `TsNumber` | `4503599627370496`, `1000000000000000000000` | yes | TS-VAL-2 | changed (was no; F12 fixed) |
| M-07c | `TsNumber` | `4503599627370496.0` | no (no `Fractional` form at 2^52) | VM-SUB-0 | new |
| M-09 | `TsUndefined` | `undefined` | yes | TS-VAL-4 | same |
| M-10 | `TsUndefined` | `Nothing` | no | VM-AX-5 | same |
| M-11 | `TsRecord` | `{("a", 1)}` | yes | TS-OBJ-1 | same |
| M-12 | `TsRecord` | `{(1, "a")}` | no | TS-OBJ-1 | same |
| M-13 | `TsRecord` | `{("a", 1), ("a", 2)}` | no | VM-MAP-1 | same |
| M-14 | `TsNumericKey` | `"1.5"`, `"NaN"`, `"1e+21"` | yes (prose) | VM-SUB-2 | same |
| M-15 | `TsNumericKey` | `"-0"`, `"01"`, `"x"` | no (prose) | VM-SUB-2 | same |
| M-16 | `TsDate` | `0` | yes | TS-VAL-7 | changed (was `0.0`) |
| M-17 | `TsDate` | `8640000000000000` | yes | VM-AX-6 | changed (was no; F12 fixed) |
| M-18 | `TsDate` | `8640000000000001` | no | VM-AX-6 | changed (number form) |
| M-19 | `TsDate` | `0.5` | no | TS-VAL-7 | same |
| M-20 | `TsDate` | `0.0` | no | TS-VAL-2 | changed (was `0` no) |
| M-21 | `TsDate` | `PositiveInfinity` | no | VM-AX-6 | same |
| M-22 | `TsDate` | `NotANumber` | yes | TS-VAL-7 | same |
| M-23 | `TsNonNullish` | `0`, `""`, `{}` | yes | TS-PRIM-1 | changed (number form) |
| M-24 | `TsNonNullish` | `Nothing`, `undefined` | no | TS-PRIM-1 | same |
| M-25 | `TsObject` | `{}`, `(1, (2, ()))`, `{"a"}` | yes | TS-PRIM-1 | changed (tuple now a sequence) |
| M-26 | `TsObject` | `"a"`, `1`, `true`, `Nothing` | no | TS-PRIM-1 | same |
| M-27 | `TsObject`, `TsValue` | `Circle` | no | TS-PRIM-1 | changed (was yes; F6 fixed) |
| M-28 | `TsObject` | `0` as a `Date` | no: a `Date` is a number (known limit) | TS-PRIM-1 | new |
| M-29 | `TsValue` | `(1)` (a 1-tuple) | no: JS has no tuples | TS-PRIM-1 | new |

### 3.2 Primitives, literals, enums

| ID | Domain | Value | Member? | Basis | vs. it1 |
| -- | ------ | ----- | ------- | ----- | ------- |
| M-30 | `Big` | `1` | yes | TS-PRIM-1 | same |
| M-31 | `Big` | `1.0` | no | VM-AX-1 | same |
| M-32 | `One` | `1` | yes | TS-LIT-1 | changed (was `1.0`) |
| M-33 | `One` | `1.0` | no | TS-LIT-1 | changed (was `1` no) |
| M-34 | `OneN` | `1` | yes | TS-LIT-1 | same |
| M-35 | `Tenth` | `0.1000000000000000055511151231257827021181583404541015625` | yes | TS-LIT-1 | same |
| M-36 | `Tenth` | `0.1` | no | VM-EQ-1 | same |
| M-37 | `Status` | `"active"` | yes | TS-LIT-1 | same |
| M-38 | `Status` | `active`, `"Active"` | no | VM-LEX-3 | same |
| M-39 | `Color` | `"RED"` | yes | TS-ENUM-1 | same |
| M-40 | `Color` | `Red` | no | TS-ENUM-1 | same |
| M-41 | `Level` | `0` | yes | TS-ENUM-1 | changed (was `0.0`) |
| M-42 | `Level` | `0.0`, `3` | no | TS-VAL-2 | changed |

### 3.3 Records

`u0` = `{("id", "u1"), ("name", "A"), ("email", Nothing), ("tags", ()), ("createdAt", 0)}`.

| ID | Domain | Value | Member? | Basis | vs. it1 |
| -- | ------ | ----- | ------- | ----- | ------- |
| M-50 | `User` | `u0` | yes | VM-REC-2 | changed (`createdAt` form) |
| M-51 | `User` | `u0 + ("age", 30)` | yes | VM-MISS-2 | changed (was `30.0`) |
| M-52 | `User` | `u0 + ("age", undefined)` | yes | TS-OBJ-2 | same |
| M-53 | `User` | `u0 + ("nick", "a")` | yes | VM-REC-3 | same |
| M-54 | `User` | `u0` with `("tags", ("x", ("y", ())))` | yes | VM-SEQ-1 | same |
| M-55 | `User` | `u0` with `("createdAt", NotANumber)` | yes | TS-VAL-7 | same |
| M-56 | `User` | `u0` without `email` | no | VM-RULE-2 | same |
| M-57 | `User` | `u0 + ("age", 30.0)` | no (whole number in `Fractional` form) | TS-VAL-2 | changed (was `30` no) |
| M-58 | `User` | `u0 + ("age", Nothing)` | no | VM-MISS-2 | same |
| M-59 | `User` | `u0 + ("age", "x")` | no | VM-MISS-2 | same |
| M-60 | `User` | `u0` with `("tags", ("x"))` | no | VM-SEQ-1 | same |
| M-61 | `User` | `u0 + (1, "x")` | no | TS-OBJ-1 | same |
| M-62 | `Shape` | `{("kind", "circle"), ("radius", 1)}` | yes | VM-RULE-3 | changed (number form) |
| M-63 | `Shape` | `{("kind", "rect"), ("w", 1), ("h", 2.5)}` | yes | VM-RULE-3 | changed (number form) |
| M-64 | `Shape` | `{("kind", "rect"), ("w", 1)}` | no | VM-RULE-2 | same |
| M-65 | `Shape` | `{("kind", "square"), ("radius", 1)}` | no | VM-RULE-3 | same |
| M-66 | `Shape` | `{("kind", circle), ("radius", 1)}` | no | VM-LEX-3 | same |
| M-67 | `Headers` | `{("contentType", "a"), ("x-id", "1")}` | yes | TS-OBJ-4 | same |
| M-68 | `Headers` | `{("contentType", "a"), ("x-n", 1)}` | no | TS-OBJ-4 | same |
| M-69 | `Headers` | `{}` | no | VM-RULE-2 | same |
| M-70 | `Flags` | `{("read", true), ("write", false)}` | yes | TS-OBJ-4 | same |
| M-71 | `Flags` | `{("read", true)}` | no | TS-OBJ-4 | same |
| M-72 | `Flags` | `{("read", "true"), ("write", false)}` | no | VM-AX-4 | same |
| M-73 | `Account` | `{("id", "a"), ("balance", 0)}` | no | TS-CLASS-1 | same |
| M-74 | `Account` | with an `onChange` value | can't be written (opaque axiom) | VM-AX-7 | same |
| M-75 | `AuditedUser` | `u0 + ("updatedAt", 0)` | yes | TS-UNION-1 | same |
| M-76 | `AuditedUser` | `u0` | no | TS-UNION-1 | same |

### 3.4 Tuples and sequences

| ID | Domain | Value | Member? | Basis | vs. it1 |
| -- | ------ | ----- | ------- | ----- | ------- |
| M-80 | `Point` | `(1, (2.5, ()))` | yes | TS-TUP-1 | changed (was `(1.0, 2.0)`) |
| M-81 | `Point` | `(1.0, 2.0)` | no (old tuple encoding) | TS-VAL-6 | changed |
| M-82 | `Point` | `(1, (2, (3, ())))` | no | TS-TUP-1 | changed |
| M-83 | `Point` | `(1, ())` | no | TS-TUP-1 | changed |
| M-84 | `Span` | `(1, ())`, `(1, (2, ()))`, `(1, (undefined, ()))` | yes | TS-TUP-1 | changed |
| M-85 | `Span` | `()` | no | TS-TUP-1 | same |
| M-86 | `Args` | `("a", ())`, `("a", (1, (2, ())))` | yes | TS-TUP-1 | changed (number form) |
| M-87 | `Args` | `("a")`, `(1, ())` | no | TS-TUP-1 | same |
| M-88 | `Ambiguous` | `(1, (2, ()))` | yes, one encoding | TS-VAL-6 | changed (was two; F1 fixed) |
| M-89 | `Seq(TsNumber)` | `(1, (2, ()))` (a `Point` value) | yes: `Point` ⊂ `Seq(TsNumber)` | VM-SEQ-1 | new |

### 3.5 Maps, sets, generics, recursion

| ID | Domain | Value | Member? | Basis | vs. it1 |
| -- | ------ | ----- | ------- | ----- | ------- |
| M-90 | `Dictionary` | `{}`, `{("a", 1)}` | yes | VM-MAP-1 | changed (number form) |
| M-91 | `Dictionary` | `{("a", "x")}`, `{(1, 1)}` | no | VM-MAP-1 | same |
| M-92 | `ByIndex` | `{("0", "a")}` | yes | TS-OBJ-4 | same |
| M-93 | `ByIndex` | `{(0, "a")}`, `{("x", "a")}` | no | TS-OBJ-4 | same |
| M-94 | `TagSet` | `{}`, `{"a", "b"}` | yes | VM-SET-1 | changed (was a record of keys) |
| M-95 | `TagSet` | `{("a", ())}`, `{1}` | no | VM-SET-1 | changed |
| M-96 | `Box(String)` | `{("value", "a")}` | yes | TS-GEN-1 | same |
| M-97 | `Box(String)` | `{("value", 1)}` | no | TS-GEN-1 | same |
| M-98 | `Result(TsNumber, String)` | `{("ok", true), ("value", 1)}`, `{("ok", false), ("error", "e")}` | yes | VM-RULE-3 | changed (number form) |
| M-99 | `Result(TsNumber, String)` | `{("ok", true), ("error", "e")}`, `{("ok", "true"), ("value", 1)}` | no | VM-RULE-3 | same |
| M-100 | `UserPage` | `{("items", (u0, ())), ("next", Nothing)}` | yes | TS-GEN-1 | same |
| M-101 | `UserPage` | `{("items", ({}, ())), ("next", Nothing)}` | no | TS-GEN-1 | same |
| M-102 | `UserPage` | `{("items", ())}` | no | VM-RULE-2 | same |
| M-103 | `Tree(String)` | `{("value", "a"), ("children", ())}` | yes | VM-RECUR-1 | same |
| M-104 | `Tree(String)` | nested one level | yes | VM-RECUR-1 | same |
| M-105 | `Tree(String)` | `{("value", "a"), ("children", ("b", ()))}` | no | VM-RECUR-1 | same |
| M-106 | `Json` | `"a"`, `1`, `1.5`, `true`, `Nothing`, `()`, `{("a", (1, ()))}` | yes | VM-RECUR-1 | changed (number form) |
| M-107 | `Json` | `1.0`, `undefined`, `("a")`, `{(1, "a")}` | no | TS-VAL-2 | changed (`1` → `1.0`) |
| M-108 | `UserKey` | `"user-1"` | yes (prose) | VM-SUB-2 | same |
| M-109 | `UserKey` | `"User-1"` | no (prose) | VM-SUB-2 | same |
| M-110 | `Nullish` | `Nothing`, `undefined` | yes | TS-PRIM-1 | same |
| M-111 | `Nullish` | `{}` | no | VM-AX-5 | same |
| M-112 | `EmptyThings` | every `TsValue` | yes | TS-PRIM-1 | changed (was any value) |
| M-112a | `EmptyThings` | `Circle` | no | TS-PRIM-1 | new (F6 fixed) |
| M-113 | `SetOrDict` | `{"a"}` | yes, `SetOf` branch only | TS-VAL-7 | changed (F2 fixed) |
| M-114 | `SetOrDict` | `{("a", ())}` | yes, `Map` branch only | TS-VAL-7 | new |
| M-115 | `SetOrDict` | `{}` | yes, both branches | VM-REC-1; N8 | new |

## 4. Go value projection (`.sets` value → Go value)

Each row group is checked by the Go test of the same number in [go/model_test.go](go/model_test.go).

| ID | Domain | `.sets` value | Go value | Basis |
| -- | ------ | ------------- | -------- | ----- |
| G-01a | `Pi` | `3.140000000000000124344978758017532527446746826171875` | `PiValue`, `== 3.14` | GO-ROSTER-1 |
| G-01b | `Tenth` | `0.1000000000000000055511151231257827021181583404541015625` | `TenthValue`, `== 0.1` | GO-ROSTER-1 |
| G-02a | `TsNumber` | `1000000000000000000000` | `TsNumber(1e21)`, exact | GO-VAL-1 |
| G-02b | `TsDate` | `8640000000000000` | `TsDate(8.64e15)`, exact | GO-VAL-1 |
| G-03a | `TsDate` | `NotANumber` | `math.NaN()` | GO-BASE-1 |
| G-03b | `TsNumber` | `NegativeZero` | `math.Copysign(0, -1)` | GO-BASE-1 |
| G-04a | `User` | `u0` | `User{Id: "u1", Name: "A", Email: nil, Tags: []string{}, CreatedAt: 0}`, `Age.Present == false` | GO-REC-1 |
| G-04b | `User` | `u0 + ("age", undefined)` | `Age: Opt{Present: true, Value: Either{IsB: true, B: TsUndefinedValue}}` | GO-UNION-6; N2 |
| G-04c | `User` | `u0 + ("age", 30)` | `Age: Opt{Present: true, Value: Either{A: 30}}` | GO-UNION-6 |
| G-05 | `User` | `u0 + ("nick", "a")` | `Extra: {"nick": "a"}` | GO-REC-1; N4 |
| G-06a | `Shape` | `{("kind", "circle"), ("radius", 1)}` | `ShapeCircle{Radius: 1}` | GO-UNION-5 |
| G-06b | `Result(TsNumber, String)` | `{("ok", false), ("error", "e")}` | `ResultFalse[TsNumber, string]{Error: "e"}` | GO-UNION-5 |
| G-07a | `Point` | `(1, (2, ()))` | `Point{1, 2}` | GO-SEQ-1 |
| G-07b | `Ambiguous` | `(1, (2, ()))` | `Ambiguous{1, 2}` | GO-UNION-0 |
| G-07c | `Span` | `(1, ())` | `Span{F0: 1}`, `F1.Present == false` | GO-UNION-4 |
| G-07d | `Args` | `("a", (1, (2, ())))` | `Args{F0: "a", Rest: []TsNumber{1, 2}}` | GO-SEQ-2 |
| G-08a | `Big` | `1` | `big.NewInt(1)` | GO-BASE-1 |
| G-08b | `TagSet` | `{"a"}` | `TagSet{"a": {}}` | GO-CALL-1 |
| G-08c | `Tree(String)` | one child `"b"` | `Tree[string]{Value: "a", Children: []Tree[string]{{Value: "b"}}}` | GO-RECUR-1 |
| G-08d | `Json` | `{("a", (1, (Nothing, ())))}` | `map[string]Json{"a": []Json{1.0, nil}}` | GO-UNION-1/7 |
| G-09 | `SetOrDict` | `{}` | `SetOrDict{IsB: false}` (held as the `SetOf` branch) | GO-UNION-6; N8 (not tested) |
| G-10 | `User` | `u0` with `email` `"x"` | `Email: &"x"` | GO-UNION-1 (not tested) |

Note: in G-08d, a `Json` number in Go is a `float64`. The `.sets` value is `1` (`Integer` form) and the Go value is `1.0`. Going back from Go to `.sets` must pick the `Integer` form for whole numbers, because `TsNumber` excludes whole `Fractional`s. Nothing tests that yet.

## 5. Go types (domain → Go type)

| ID | Domain | Go type | Exactness | Basis |
| -- | ------ | ------- | --------- | ----- |
| T-01 | `Id` | `= string` | wider | GO-DECL-1 |
| T-02 | `Count` | `= TsNumber` | wider | GO-DECL-1 |
| T-03 | `Big` | `= *big.Int` | wider (`nil`) | GO-BASE-1 |
| T-04 | `Pi`, `Tenth` | `float64` + one constant | wider; N1 | GO-ROSTER-1 |
| T-05 | `One`, `OneN` | `int8` + one constant | wider; N1, N6 | GO-ROSTER-1 |
| T-06 | `Status`, `Color` | `string` + constants | wider; N1 | GO-ROSTER-1 |
| T-07 | `Level` | `int8` + `Level0…2` | wider; N1, N6 | GO-ROSTER-1 |
| T-08 | `User` | struct of 6 fields + `Extra` | wider | GO-REC-1 |
| T-09 | `Shape` | sealed interface, `ShapeCircle`, `ShapeRect` | wider | GO-UNION-5 |
| T-10 | `Point` | `[2]TsNumber` | wider | GO-SEQ-1 |
| T-11 | `Span` | `struct{F0; F1 Opt[Either[TsNumber, TsUndefined]]}` | wider | GO-UNION-4 |
| T-12 | `Args` | `struct{F0 string; Rest []TsNumber}` | wider | GO-SEQ-2 |
| T-13 | `Dictionary` | `map[string]TsNumber` | wider | GO-CALL-1 |
| T-14 | `ByIndex` | `map[TsNumericKey]string` | wider; N5 | GO-CALL-1 |
| T-15 | `Headers` | `struct{ContentType; Rest map[string]string}` | wider | GO-REC-1 |
| T-16 | `Box` | `Box[T any]` | wider | GO-DECL-1 |
| T-17 | `Result` | `Result[T, E any]` sealed, `ResultTrue`, `ResultFalse` | wider; N10 | GO-UNION-5 |
| T-18 | `Page` | `Page[T any]` (bound lost) | wider; N9 | GO-DECL-1 |
| T-19 | `UserPage` | `= Page[User]` | wider | GO-DECL-1 |
| T-20 | `Tree` | `Tree[T any]`, `Children []Tree[T]` | wider | GO-RECUR-1 |
| T-21 | `Json` | `= any` | wider | GO-UNION-7 |
| T-22 | `TagSet` | `map[string]struct{}` | wider | GO-CALL-1 |
| T-23 | `AuditedUser` | one merged struct (7 fields + `Extra`) | wider | GO-INTER-1 |
| T-24 | `UserKey` | `string` | wider; N1 | GO-DECL-1 |
| T-25 | `Account` | struct with `OnChange TsFunction` | wider | GO-REC-1 |
| T-26 | `Ambiguous` | `[]TsNumber` | wider | GO-UNION-0 |
| T-27 | `EmptyThings` | `= *Either[TsNonNullish, TsUndefined]` (should be `any`) | wider; N3 | GO-UNION-1/6 |
| T-28 | `Nullish` | `= *TsUndefined` | exact | GO-UNION-1 |
| T-29 | `SetOrDict` | `Either[map[string]struct{}, map[string]TsNonNullish]` | not one to one at `{}`; N8 | GO-UNION-6 |
| T-30 | `TsValue` | `= any` | wider; F6 | GO-UNION-7 |
| T-31 | `TsObject` | `= any` | wider | GO-REF-1 |
| T-32 | `TsNonNullish` | `= any` | wider | GO-REF-1 |
| T-33 | `TsPrimitive` | `= any` | wider | GO-UNION-7 |
| T-34 | `TsDate` | `TsNumber` | wider | GO-UNION-3 |
| T-35 | `TsNumber` | `float64` | wider (whole `Fractional`s aren't excluded) | GO-REF-1 |
| T-36 | `TsUndefined` | `string` + `TsUndefinedValue` | wider | GO-ROSTER-1 |
| T-37 | `TsSymbol`, `TsFunction` | `struct{ Handle any }` | — | GO-DECL-1 |

## 6. Spelling checks

These are unchanged from iteration 1. No standard produces L-07 (`{()}`) any more, but it stays as an open VM question.

| ID | Text | Expected | Basis |
| -- | ---- | -------- | ----- |
| L-01 | `{"a, b"}` | error (split into `"a` and ` b"`) | VM-LEX-1 |
| L-02 | `{"a\, b"}` | the string `a, b` | VM-LEX-0 |
| L-03 | `{\Number}` | the atom `Number`, not the domain | VM-LEX-4 |
| L-04 | `{\1}` | the atom `1`, not a number | VM-LEX-6 |
| L-05 | `{1e3}` | the atom `1e3` | VM-LEX-2a |
| L-06 | `{+007.50}` | the number `7.5` | VM-LEX-2 |
| L-07 | `{()}` | unclear: the domain holding the empty value, or an atom `()` | Issue F9 |
| L-08 | `{a"b}` | error (unclosed quote) | VM-LEX-1 |
