# Test vectors (iteration 1, TypeScript)

These vectors are the baseline for later iterations. Every iteration should keep each vector, or record why the expected result changed.

- **Domains:** [examples/types.sets](examples/types.sets) and [ts.sets](ts.sets), from [STANDARD.md](STANDARD.md) 0.2.0 and `VALUE-MODEL.md` 0.2.0.
- **Notation:** `.sets` value notation (VM §9).
- **IDs** are stable. A later iteration may change an expected result, but never reuse an ID for a different input.
- **Basis** says which rules fix the result:
  - **VM**: value-model rules alone. These should survive any change to the TS standard.
  - **TS**: a choice made by the TS standard. A new TS standard may change it.
  - **Issue Fn**: the current expected result shows a known problem from [FINDINGS.md](FINDINGS.md). A fix is expected to change this vector.

## 1. Projection: JavaScript value → `.sets` value

| ID | JavaScript value (static type) | `.sets` value | Basis |
| -- | ------------------------------ | ------------- | ----- |
| P-01 | `"abc"` | `"abc"` | VM-LEX-1 |
| P-02 | `""` | `""` | VM-LEX-1 |
| P-03 | `"a, b"` | `"a\, b"` | VM-LEX-0, VM-LEX-1 |
| P-04 | `'say "hi"'` | `"say \"hi\""` | VM-LEX-1 |
| P-05 | `"x\n"` | `"x\n"` | VM-LEX-6 |
| P-06 | `"{a}"` | `"\{a\}"` | VM-LEX-1 |
| P-07 | `1` | `1.0` | TS-NUM-1 |
| P-08 | `0` | `0.0` | TS-NUM-1 |
| P-09 | `-0` | `NegativeZero` | VM-AX-3 |
| P-10 | `0.5` | `0.5` | TS-NUM-1 |
| P-11 | `0.1` | `0.1000000000000000055511151231257827021181583404541015625` | VM-AX-2, TS-NUM-1; Issue F5 |
| P-12 | `3.14` | `3.140000000000000124344978758017532527446746826171875` | VM-AX-2, TS-NUM-1; Issue F5 |
| P-13 | `1e21` | `1000000000000000000000.0`, which isn't in `TsNumber` (M-07b) | TS-NUM-1; Issue F12 |
| P-14 | `NaN` | `NotANumber` | VM-AX-3 |
| P-15 | `Infinity`, `-Infinity` | `PositiveInfinity`, `NegativeInfinity` | VM-AX-3 |
| P-16 | `1n` | `1` | TS-PRIM-1 |
| P-17 | `-7n` | `-7` | TS-PRIM-1, VM-LEX-2 |
| P-18 | `true`, `false` | `true`, `false` | VM-AX-4 |
| P-19 | `null` | `Nothing` | VM-AX-5 |
| P-20 | `undefined` | `undefined` | TS-PRIM-1 |
| P-21 | `{}` | `{}` | VM-REC-1 |
| P-22 | `{ name: "a", age: 1 }` | `{("name", "a"), ("age", 1.0)}` | VM-REC-1 |
| P-23 | `{ 1: "a" }` | `{("1", "a")}` | TS-OBJ-1; Issue F7 |
| P-24 | `{ age: undefined }` | `{("age", undefined)}` | VM-MISS-1, TS-OBJ-2 |
| P-25 | `{ age: null }` | `{("age", Nothing)}` | VM-MISS-1 |
| P-26 | `[]` (`string[]`) | `()` | VM-SEQ-1 |
| P-27 | `["a", "b"]` (`string[]`) | `("a", ("b", ()))` | VM-SEQ-1 |
| P-28 | `[1, 2]` (`[number, number]`) | `(1.0, 2.0)` | VM-TUP-1; Issue F1 |
| P-29 | `[1, 2]` (`number[]`) | `(1.0, (2.0, ()))` | VM-SEQ-1; Issue F1 |
| P-30 | `[1]` (`[number, number?]`) | `(1.0)` | TS-TUP-1 |
| P-31 | `["a", 1, 2]` (`[string, ...number[]]`) | `("a", (1.0, (2.0, ())))` | TS-TUP-1 |
| P-32 | `["a", 1, 2]` (`(string \| number)[]`) | `("a", (1.0, (2.0, ())))` | VM-SEQ-1 (same as P-31) |
| P-33 | `[[1]]` (`number[][]`) | `((1.0, ()), ())` | VM-SEQ-1, VM-TUP-1 |
| P-34 | `new Set(["a", "b"])` | `{("a", ()), ("b", ())}` | VM-PRES-2, TS-LIB-1 |
| P-35 | `new Set()` | `{}` | VM-PRES-2, VM-REC-1 |
| P-36 | `new Map([["a", 1]])` | `{("a", 1.0)}` | TS-LIB-1; Issue F2 |
| P-37 | `new Map([[1, "a"]])` | `{(1.0, "a")}` | TS-LIB-1 (key stays a number, unlike P-23) |
| P-38 | `new Date(0)` | `0.0` | TS-LIB-1; Issue F2 |
| P-39 | `new Date(NaN)` | `NotANumber` | TS-LIB-1 |
| P-40 | `Color.Red` | `"RED"` | TS-ENUM-1 |
| P-41 | `Level.Mid` | `1.0` | TS-ENUM-1, TS-NUM-1 |
| P-42 | `"\uD800"` (lone surrogate) | no `.sets` value; the projection reports it | Issue F8 |

## 2. Equality

| ID | Left | Right | Equal? | Basis |
| -- | ---- | ----- | ------ | ----- |
| E-01 | `1` | `1.0` | no | VM-LEX-2, VM-EQ-1 |
| E-02 | `1.5` | `1.50` | yes | VM-LEX-2 |
| E-03 | `"1"` | `1` | no | VM-KIND-1 |
| E-04 | `true` | `"true"` | no | VM-LEX-3 |
| E-05 | `()` | `{}` | yes (core context) | VM-AX-5 |
| E-06 | `Nothing` | `{}` | no | VM-AX-5 |
| E-07 | `Nothing` | `undefined` | no | VM-LEX-3, TS-PRIM-1 |
| E-08 | `(1.0)` | `(1.0, ())` | no (different arity) | VM-EQ-1 |
| E-09 | `{("a", 1.0), ("b", 2.0)}` | `{("b", 2.0), ("a", 1.0)}` | yes (key order is lost) | VM-EQ-1 |
| E-10 | P-34 `Set` value `{("a", ())}` | `{ a: {} }` → `{("a", {})}` | yes | VM-AX-5; Issue F2 |
| E-11 | P-36 `Map` value | `{ a: 1 }` → `{("a", 1.0)}` | yes | Issue F2 |
| E-12 | P-38 `Date` value | `0` → `0.0` | yes | Issue F2 |
| E-13 | P-28 | P-29 | no, though the JS value is the same | Issue F1 |
| E-14 | P-31 | P-32 | yes | VM-SEQ-1 |
| E-15 | `1n` → `1` | `1` → `1.0` | no | TS-NUM-1 (REVIEW R1) |

## 3. Membership

Each row checks one value against one domain. "yes" means member, "no" means non-member.

### 3.1 Prelude (`ts.sets`)

| ID | Domain | Value | Member? | Basis |
| -- | ------ | ----- | ------- | ----- |
| M-01 | `TsNumber` | `1.0` | yes | TS-NUM-1 |
| M-02 | `TsNumber` | `0.5` | yes | VM-SUB-0 |
| M-03 | `TsNumber` | `NotANumber` | yes | VM-AX-3 |
| M-04 | `TsNumber` | `NegativeZero` | yes | VM-AX-3 |
| M-05 | `TsNumber` | `1` | no (`Integer` form) | TS-NUM-1; Issue F3 |
| M-06 | `TsNumber` | `0.1` | no (not a binary64 value) | VM-SUB-0; Issue F5 |
| M-07 | `TsNumber` | `9007199254740993.0` | no (2^53 + 1 isn't representable) | VM-SUB-0 |
| M-07a | `TsNumber` | `4503599627370495.0` (2^52 − 1) | yes | VM-SUB-0 |
| M-07b | `TsNumber` | `4503599627370496.0` (2^52), `1000000000000000000000.0` | no, though JS has these numbers | VM-SUB-0; Issue F12 |
| M-07c | `Binary64Finite` | `4503599627370496` (2^52, `Integer` form) | yes | VM-SUB-0 |
| M-08 | `TsNumber` | `"1"` | no | VM-KIND-1 |
| M-09 | `TsUndefined` | `undefined` | yes | TS-PRIM-1 |
| M-10 | `TsUndefined` | `Nothing` | no | VM-AX-5 |
| M-11 | `TsRecord` | `{("a", 1.0)}` | yes | TS-OBJ-1 |
| M-12 | `TsRecord` | `{(1, "a")}` | no (key not a string) | TS-OBJ-1; Issue F7 |
| M-13 | `TsRecord` | `{("a", 1.0), ("a", 2.0)}` | no (repeated key) | VM-REC-1, VM-MAP-1 |
| M-14 | `TsNumericKey` | `"1.5"`, `"NaN"`, `"1e+21"` | yes | TS-OBJ-4 |
| M-15 | `TsNumericKey` | `"-0"`, `"01"`, `"x"` | no | TS-OBJ-4 |
| M-16 | `TsDate` | `0.0` | yes | TS-LIB-1 |
| M-17 | `TsDate` | `8640000000000000.0` (the largest JS time value) | no: above 2^52, so not in `TsNumber` | Issue F12 |
| M-17a | `TsDate` | `4503599627370495.0` | yes | VM-AX-6 |
| M-18 | `TsDate` | `8640000000000001.0` | no (out of range) | VM-AX-6 |
| M-19 | `TsDate` | `0.5` | no (not whole) | TS-LIB-1 |
| M-20 | `TsDate` | `0` | no (`Integer` form) | TS-NUM-1 |
| M-21 | `TsDate` | `PositiveInfinity` | no (comparisons on atoms fail) | VM-AX-6 |
| M-22 | `TsDate` | `NotANumber` | yes | TS-LIB-1 |
| M-23 | `TsNonNullish` | `0.0`, `""`, `{}` | yes | TS-PRIM-1 |
| M-24 | `TsNonNullish` | `Nothing`, `undefined` | no | TS-PRIM-1 |
| M-25 | `TsObject` | `{}`, `(1.0, 2.0)` | yes | TS-PRIM-1 |
| M-26 | `TsObject` | `"a"`, `1.0`, `1`, `true`, `Nothing` | no | TS-PRIM-1 |
| M-27 | `TsObject` | `Circle` (a non-JS atom) | yes | Issue F6 |

### 3.2 Primitives, literals, enums

| ID | Domain | Value | Member? | Basis |
| -- | ------ | ----- | ------- | ----- |
| M-30 | `Big` | `1` | yes | TS-PRIM-1 |
| M-31 | `Big` | `1.0` | no | VM-AX-1 |
| M-32 | `One` | `1.0` | yes | TS-NUM-2 |
| M-33 | `One` | `1` | no | TS-NUM-2 |
| M-34 | `OneN` | `1` | yes | TS-LIT-1 |
| M-35 | `Tenth` | `0.1000000000000000055511151231257827021181583404541015625` | yes | TS-NUM-2 |
| M-36 | `Tenth` | `0.1` | no | VM-EQ-1; Issue F5 |
| M-37 | `Status` | `"active"` | yes | TS-LIT-1 |
| M-38 | `Status` | `active`, `"Active"` | no | VM-LEX-3, VM-EQ-1 |
| M-39 | `Color` | `"RED"` | yes | TS-ENUM-1 |
| M-40 | `Color` | `Red` | no (atom; the member name isn't data) | TS-ENUM-1 |
| M-41 | `Level` | `0.0` | yes | TS-ENUM-1 |
| M-42 | `Level` | `0`, `3.0` | no | TS-NUM-1 |

### 3.3 Records (`User`, `Shape`, `Headers`, `Flags`, `Account`)

Shorthand: `u0` = `{("id", "u1"), ("name", "A"), ("email", Nothing), ("tags", ()), ("createdAt", 0.0)}`. Rows written `u0 +` add fields, and rows written `u0 −` remove them.

| ID | Domain | Value | Member? | Basis |
| -- | ------ | ----- | ------- | ----- |
| M-50 | `User` | `u0` | yes | VM-REC-2 |
| M-51 | `User` | `u0 + ("age", 30.0)` | yes | VM-MISS-2 |
| M-52 | `User` | `u0 + ("age", undefined)` | yes | TS-OBJ-2 |
| M-53 | `User` | `u0 + ("nick", "a")` | yes (open record) | VM-REC-3 |
| M-54 | `User` | `u0` with `("tags", ("x", ("y", ())))` | yes | VM-SEQ-1 |
| M-55 | `User` | `u0` with `("createdAt", NotANumber)` | yes | TS-LIB-1 |
| M-56 | `User` | `u0 − ("email", …)` | no (required field missing) | VM-RULE-2 |
| M-57 | `User` | `u0 + ("age", 30)` | no (`30` is a bigint) | TS-NUM-1 |
| M-58 | `User` | `u0 + ("age", Nothing)` | no | VM-MISS-2 |
| M-59 | `User` | `u0 + ("age", "x")` | no | VM-MISS-2 |
| M-60 | `User` | `u0` with `("tags", ("x"))` | no (1-tuple isn't a sequence) | VM-SEQ-1 |
| M-61 | `User` | `u0 + (1, "x")` | no (non-string key) | TS-OBJ-1 |
| M-62 | `Shape` | `{("kind", "circle"), ("radius", 1.0)}` | yes | VM-RULE-3 |
| M-63 | `Shape` | `{("kind", "rect"), ("w", 1.0), ("h", 2.0)}` | yes | VM-RULE-3 |
| M-64 | `Shape` | `{("kind", "rect"), ("w", 1.0)}` | no | VM-RULE-2 |
| M-65 | `Shape` | `{("kind", "square"), ("radius", 1.0)}` | no | VM-RULE-3 |
| M-66 | `Shape` | `{("kind", circle), ("radius", 1.0)}` | no (atom, not string) | VM-LEX-3 |
| M-67 | `Headers` | `{("contentType", "a"), ("x-id", "1")}` | yes | TS-OBJ-4 |
| M-68 | `Headers` | `{("contentType", "a"), ("x-n", 1.0)}` | no (index signature) | TS-OBJ-4 |
| M-69 | `Headers` | `{}` | no | VM-RULE-2 |
| M-70 | `Flags` | `{("read", true), ("write", false)}` | yes | TS-OBJ-4 |
| M-71 | `Flags` | `{("read", true)}` | no | TS-OBJ-4 |
| M-72 | `Flags` | `{("read", "true"), ("write", false)}` | no | VM-AX-4 |
| M-73 | `Account` | `{("id", "a"), ("balance", 0.0)}` | no (`onChange` missing) | TS-CLASS-1 |
| M-74 | `Account` | a member with an `onChange` value | can't be written: `TsFunction` members have no notation | VM-KIND-2, VM-AX-7 |
| M-75 | `AuditedUser` | `u0 + ("updatedAt", 0.0)` | yes | TS-UNION-1 |
| M-76 | `AuditedUser` | `u0` | no | TS-UNION-1 |

### 3.4 Tuples and sequences

| ID | Domain | Value | Member? | Basis |
| -- | ------ | ----- | ------- | ----- |
| M-80 | `Point` | `(1.0, 2.0)` | yes | VM-TUP-1 |
| M-81 | `Point` | `(1, 2)` | no | TS-NUM-1 |
| M-82 | `Point` | `(1.0, (2.0, ()))` | no | VM-SEQ-1; Issue F1 |
| M-83 | `Point` | `(1.0)`, `(1.0, 2.0, 3.0)` | no | VM-TUP-1 |
| M-84 | `Span` | `(1.0)`, `(1.0, 2.0)`, `(1.0, undefined)` | yes | TS-TUP-1 |
| M-85 | `Span` | `()` | no | TS-TUP-1 |
| M-86 | `Args` | `("a", ())`, `("a", (1.0, (2.0, ())))` | yes | TS-TUP-1 |
| M-87 | `Args` | `("a")`, `(1.0, ())` | no | TS-TUP-1 |
| M-88 | `Ambiguous` | `(1.0, 2.0)` and `(1.0, (2.0, ()))` | both yes; one JS value, so a TS-WIRE-1 diagnostic | Issue F1 |

### 3.5 Maps, sets, generics, recursion

| ID | Domain | Value | Member? | Basis |
| -- | ------ | ----- | ------- | ----- |
| M-90 | `Dictionary` | `{}`, `{("a", 1.0)}` | yes | VM-MAP-1 |
| M-91 | `Dictionary` | `{("a", "x")}`, `{(1, 1.0)}` | no | VM-MAP-1 |
| M-92 | `ByIndex` | `{("0", "a")}`, `{("1.5", "a")}` | yes | TS-OBJ-4 |
| M-93 | `ByIndex` | `{(0.0, "a")}`, `{("x", "a")}` | no | TS-OBJ-4 |
| M-94 | `TagSet` | `{}`, `{("a", ())}` | yes | TS-LIB-1 |
| M-95 | `TagSet` | `{("a", 1.0)}`, `{(1.0, ())}` | no | TS-LIB-1 |
| M-96 | `Box(String)` | `{("value", "a")}` | yes | TS-GEN-1 |
| M-97 | `Box(String)` | `{("value", 1.0)}` | no | TS-GEN-1 |
| M-98 | `Result(TsNumber, String)` | `{("ok", true), ("value", 1.0)}`, `{("ok", false), ("error", "e")}` | yes | VM-RULE-3 |
| M-99 | `Result(TsNumber, String)` | `{("ok", true), ("error", "e")}`, `{("ok", "true"), ("value", 1.0)}` | no | VM-RULE-3, VM-AX-4 |
| M-100 | `UserPage` | `{("items", (u0, ())), ("next", Nothing)}` | yes | TS-GEN-1 |
| M-101 | `UserPage` | `{("items", ({}, ())), ("next", Nothing)}` | no (`{}` isn't a `User`) | TS-GEN-1 |
| M-102 | `UserPage` | `{("items", ())}` | no (`next` missing) | VM-RULE-2 |
| M-103 | `Tree(String)` | `{("value", "a"), ("children", ())}` | yes | VM-RECUR-1 |
| M-104 | `Tree(String)` | `{("value", "a"), ("children", ({("value", "b"), ("children", ())}, ()))}` | yes | VM-RECUR-1; Issue F10 |
| M-105 | `Tree(String)` | `{("value", "a"), ("children", ("b", ()))}` | no | VM-RECUR-1 |
| M-106 | `Json` | `"a"`, `1.0`, `true`, `Nothing`, `()`, `{("a", (1.0, ()))}` | yes | VM-RECUR-1 |
| M-107 | `Json` | `1`, `undefined`, `("a")`, `{(1, "a")}` | no | TS-NUM-1, VM-SEQ-1 |
| M-108 | `UserKey` | `"user-1"`, `"user-"` | yes (by prose) | VM-SUB-2 |
| M-109 | `UserKey` | `"User-1"`, `"use"` | no (by prose) | VM-SUB-2 |
| M-110 | `Nullish` | `Nothing`, `undefined` | yes | TS-PRIM-1 |
| M-111 | `Nullish` | `{}` | no | VM-AX-5 |
| M-112 | `EmptyThings` | any value, e.g. `Circle` | yes (works out to `U`) | TS-PRIM-1 |
| M-113 | `SetOrDict` | `{("a", ())}` | yes, as both the `Set` and the `Record` branch | Issue F2 |

## 4. Spelling checks (value-model layer)

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
