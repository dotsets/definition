# Vector results (iteration 3: TypeScript → `.sets`)

Expected results under [STANDARD.md](STANDARD.md) for every `P-`, `E-` and `M-` vector of iteration 2's [VECTORS.md](../../iteration2-typescript-2-golang/VECTORS.md) §1–§3, and new vectors for what's new.

- **Domains:** [sets/types.sets](../sets/types.sets) and [sets/prelude.sets](../sets/prelude.sets).
- **Notation:** `.sets` value notation (VM §9).
- **Basis** is the rule that fixes the result: a `TS-*` rule of this standard, or a `VM-*` rule alone.
- **vs. it2** is `same`, `changed (why)` or `not applicable`. No vector turned out to be not applicable.
- **IDs** are iteration 2's. New vectors continue the numbering (`P-45`, `E-17`, `M-116`).

Every membership row was computed a second time with a small hand-written model of the two `.sets` files (239 checks, no disagreement). The model mirrors the definitions; it doesn't parse `.sets`.

## 0. Renames

The hand-off may not name the source language (TS-OUT-1), so the prelude names changed. A row that differs from iteration 2 only by one of these domain names is marked `same`.

| Iteration 2 | Iteration 3 |
| ----------- | ----------- |
| `TsNumber` | `Double` |
| `TsUndefined`, atom `undefined` | `Unset`, atom `unset` |
| `TsRecord` | `Object` |
| `TsNumericKey` | `NumericKey` |
| `TsDate` | `EpochMillis` |
| `TsNonNullish` | `Something` |
| `TsObject` | `Composite` |
| `TsValue` | `Anything` |

## 1. Projection: JavaScript value → `.sets` value

| ID | JavaScript value (static type) | `.sets` value | Basis | vs. it2 |
| -- | ------------------------------ | ------------- | ----- | ------- |
| P-01 | `"abc"` | `"abc"` | TS-STR-1 | same |
| P-02 | `""` | `""` | TS-STR-1 | same |
| P-03 | `"a, b"` | `"a\, b"` | TS-STR-1 | same |
| P-04 | `'say "hi"'` | `"say \"hi\""` | TS-STR-1 | same |
| P-05 | `"x\n"` | `"x\n"` | TS-STR-1 | same |
| P-06 | `"{a}"` | `"\{a\}"` | TS-STR-1 | same |
| P-07 | `1` | `1` | TS-NUM-1 | same |
| P-08 | `0` | `0` | TS-NUM-1 | same |
| P-09 | `-0` | `NegativeZero` | TS-NUM-1 | same |
| P-10 | `0.5` | `0.5` | TS-NUM-1 | same |
| P-11 | `0.1` | `0.1000000000000000055511151231257827021181583404541015625` | TS-NUM-1 | same |
| P-12 | `3.14` | `3.140000000000000124344978758017532527446746826171875` | TS-NUM-1 | same |
| P-13 | `1e21` | `1000000000000000000000` (in `Double`) | TS-NUM-1 | same |
| P-14 | `NaN` | `NotANumber` | TS-NUM-1 | same |
| P-15 | `Infinity`, `-Infinity` | `PositiveInfinity`, `NegativeInfinity` | TS-NUM-1 | same |
| P-16 | `1n` | `1` | TS-NUM-2 | same |
| P-17 | `-7n` | `-7` | TS-NUM-2 | same |
| P-18 | `true`, `false` | `true`, `false` | TS-MAP-1 | same |
| P-19 | `null` | `Nothing` | TS-NIL-1 | same |
| P-20 | `undefined` | `unset` | TS-NIL-1 | changed (was `undefined`; the atom has a neutral name, TS-OUT-1) |
| P-21 | `{}` | `{}` | TS-OBJ-1 | same |
| P-22 | `{ name: "a", age: 1 }` | `{("name", "a"), ("age", 1)}` | TS-OBJ-1 | same |
| P-23 | `{ 1: "a" }` | `{("1", "a")}` | TS-OBJ-1 | same |
| P-24 | `{ age: undefined }` | `{("age", unset)}` | TS-NIL-1 | changed (was `{("age", undefined)}`; atom renamed) |
| P-25 | `{ age: null }` | `{("age", Nothing)}` | TS-NIL-1 | same |
| P-26 | `[]` (`string[]`) | `()` | TS-SEQ-1 | same |
| P-27 | `["a", "b"]` (`string[]`) | `("a", ("b", ()))` | TS-SEQ-1 | same |
| P-28 | `[1, 2]` (`[number, number]`) | `(1, (2, ()))` | TS-TUP-1 | same |
| P-29 | `[1, 2]` (`number[]`) | `(1, (2, ()))` | TS-SEQ-1 | same |
| P-30 | `[1]` (`[number, number?]`) | `(1, ())` | TS-TUP-1 | same |
| P-31 | `["a", 1, 2]` (`[string, ...number[]]`) | `("a", (1, (2, ())))` | TS-TUP-1 | same |
| P-32 | `["a", 1, 2]` (`(string \| number)[]`) | `("a", (1, (2, ())))` | TS-SEQ-1 | same |
| P-33 | `[[1]]` (`number[][]`) | `((1, ()), ())` | TS-SEQ-1 | same |
| P-34 | `new Set(["a", "b"])` | `{"a", "b"}` | TS-LIB-1 | same |
| P-35 | `new Set()` | `{}` | TS-LIB-1 | same |
| P-36 | `new Map([["a", 1]])` | `{("a", 1)}` | TS-LIB-1 | same |
| P-37 | `new Map([[1, "a"]])` | `{(1, "a")}` | TS-LIB-1 | same |
| P-38 | `new Date(0)` | `0` | TS-LIB-1 | same |
| P-39 | `new Date(NaN)` | `NotANumber` | TS-LIB-1 | same |
| P-40 | `Color.Red` | `"RED"` | TS-ENUM-1 | same |
| P-41 | `Level.Mid` | `1` | TS-ENUM-1 | same |
| P-42 | `"\uD800"` | no value; diagnostic | TS-STR-1 | same |
| P-43 | `new Set([[1, 2]])` | `{(1, (2, ()))}` | TS-LIB-1, TS-SEQ-1 | same (the value is now known to equal a `Map` value: E-17) |
| P-44 | `4503599627370496` (2^52) | `4503599627370496` | TS-NUM-1 | same |

## 2. Equality

| ID | Left | Right | Equal? | Basis | vs. it2 |
| -- | ---- | ----- | ------ | ----- | ------- |
| E-01 | `1` | `1.0` | no | VM-EQ-1 | same |
| E-02 | `1.5` | `1.50` | yes | VM-LEX-2 | same |
| E-03 | `"1"` | `1` | no | VM-KIND-1 | same |
| E-04 | `true` | `"true"` | no | VM-LEX-3 | same |
| E-05 | `()` | `{}` | yes | VM-AX-5 | same |
| E-06 | `Nothing` | `{}` | no | VM-AX-5 | same |
| E-07 | `Nothing` | `unset` | no | TS-NIL-1 | same (the atom was `undefined`) |
| E-08 | `(1)` | `(1, ())` | no | VM-EQ-1 | same |
| E-09 | `{("a", 1), ("b", 2)}` | `{("b", 2), ("a", 1)}` | yes (order is lost) | VM-EQ-1, TS-WIRE-2 | same |
| E-10 | `new Set(["a"])` → `{"a"}` | `{ a: {} }` → `{("a", {})}` | no | TS-LIB-1 | same |
| E-11 | `new Map([["a", 1]])` | `{ a: 1 }` | yes, justified; diagnostic | TS-WIRE-1 | same (`Scores` and `Dictionary` now carry an `Overlap:` note) |
| E-12 | `new Date(0)` | `0` | yes, justified; diagnostic | TS-WIRE-1 | same |
| E-13 | P-28 | P-29 | yes | TS-TUP-1 | same |
| E-14 | P-31 | P-32 | yes | TS-SEQ-1 | same |
| E-15 | `1n` → `1` | `1` → `1` | yes, justified; diagnostic | TS-NUM-2, TS-WIRE-1 | same (`One` and `OneN` now carry an `Overlap:` note) |
| E-16 | `new Set()` | `new Map()`, `{}`, `[]` | yes (all `{}`) | TS-WIRE-1 | same |

## 3. Membership

### 3.1 Prelude

| ID | Domain | Value | Member? | Basis | vs. it2 |
| -- | ------ | ----- | ------- | ----- | ------- |
| M-01 | `Double` | `1` | yes | TS-NUM-1 | same |
| M-01a | `Double` | `1.0` | no (whole numbers use `Integer` form) | TS-NUM-1 | same |
| M-02 | `Double` | `0.5` | yes | TS-PRE-1 | same |
| M-03 | `Double` | `NotANumber` | yes | VM-AX-3 | same |
| M-04 | `Double` | `NegativeZero` | yes | VM-AX-3 | same |
| M-05 | `Double` | `"1"` | no | VM-KIND-1 | same |
| M-06 | `Double` | `0.1` | no (not binary64) | VM-SUB-0 | same |
| M-07 | `Double` | `9007199254740993` | no | VM-SUB-0 | same |
| M-07a | `Double` | `4503599627370495` | yes | VM-SUB-0 | same |
| M-07b | `Double` | `4503599627370496`, `1000000000000000000000` | yes | TS-NUM-1 | same |
| M-07c | `Double` | `4503599627370496.0` | no (it's a `WholeFractional`, and it isn't in `Float64` either) | TS-PRE-1 | same |
| M-09 | `Unset` | `unset` | yes | TS-NIL-1 | same (the atom was `undefined`) |
| M-10 | `Unset` | `Nothing` | no | VM-AX-5 | same |
| M-11 | `Object` | `{("a", 1)}` | yes | TS-OBJ-1 | same |
| M-12 | `Object` | `{(1, "a")}` | no | TS-OBJ-1 | same |
| M-13 | `Object` | `{("a", 1), ("a", 2)}` | no | VM-MAP-1 | same |
| M-14 | `NumericKey` | `"1.5"`, `"NaN"`, `"1e+21"` | yes (prose) | TS-PRE-1 | same (the prose no longer names the source) |
| M-15 | `NumericKey` | `"-0"`, `"01"`, `"x"` | no (prose) | TS-PRE-1 | same |
| M-16 | `EpochMillis` | `0` | yes | TS-LIB-1 | same |
| M-17 | `EpochMillis` | `8640000000000000` | yes | VM-AX-6 | same |
| M-18 | `EpochMillis` | `8640000000000001` | no | VM-AX-6 | same |
| M-19 | `EpochMillis` | `0.5` | no | TS-LIB-1 | same |
| M-20 | `EpochMillis` | `0.0` | no | TS-NUM-1 | same |
| M-21 | `EpochMillis` | `PositiveInfinity` | no | TS-PRE-1 | same |
| M-22 | `EpochMillis` | `NotANumber` | yes | TS-LIB-1 | same |
| M-23 | `Something` | `0`, `""`, `{}` | yes | TS-MAP-1 | same |
| M-24 | `Something` | `Nothing`, `unset` | no | TS-MAP-1 | same (the atom was `undefined`) |
| M-25 | `Composite` | `{}`, `(1, (2, ()))`, `{"a"}` | yes | TS-MAP-1 | same |
| M-26 | `Composite` | `"a"`, `1`, `true`, `Nothing` | no | TS-MAP-1 | same |
| M-27 | `Composite`, `Anything` | `Circle` | no | TS-PRE-1 | same |
| M-28 | `Composite` | `0` as a `Date` | no: a `Date` is a number (L2) | TS-LIB-1 | same |
| M-29 | `Anything` | `(1)` (a 1-tuple) | no: not a sequence | TS-PRE-1 | same |

### 3.2 Primitives, literals, enums

| ID | Domain | Value | Member? | Basis | vs. it2 |
| -- | ------ | ----- | ------- | ----- | ------- |
| M-30 | `Big` | `1` | yes | TS-MAP-1 | same |
| M-31 | `Big` | `1.0` | no | VM-AX-1 | same |
| M-32 | `One` | `1` | yes | TS-LIT-1 | same |
| M-33 | `One` | `1.0` | no | TS-LIT-1 | same |
| M-34 | `OneN` | `1` | yes | TS-LIT-1 | same |
| M-35 | `Tenth` | `0.1000000000000000055511151231257827021181583404541015625` | yes | TS-LIT-1 | same |
| M-36 | `Tenth` | `0.1` | no | VM-EQ-1 | same |
| M-37 | `Status` | `"active"` | yes | TS-UNION-1 | same |
| M-38 | `Status` | `active`, `"Active"` | no | VM-LEX-3 | same |
| M-39 | `Color` | `"RED"` | yes | TS-ENUM-1 | same |
| M-40 | `Color` | `Red` | no | TS-ENUM-1 | same (the name `Red` is now in a `Names:` note) |
| M-41 | `Level` | `0` | yes | TS-ENUM-1 | same |
| M-42 | `Level` | `0.0`, `3` | no | TS-NUM-1 | same |

### 3.3 Records

`u0` = `{("id", "u1"), ("name", "A"), ("email", Nothing), ("tags", ()), ("createdAt", 0)}`.

| ID | Domain | Value | Member? | Basis | vs. it2 |
| -- | ------ | ----- | ------- | ----- | ------- |
| M-50 | `User` | `u0` | yes | TS-OBJ-1 | same |
| M-51 | `User` | `u0 + ("age", 30)` | yes | TS-OBJ-3 | same |
| M-52 | `User` | `u0 + ("age", unset)` | yes | TS-OBJ-3 | same (the atom was `undefined`) |
| M-53 | `User` | `u0 + ("nick", "a")` | yes | VM-REC-3 | same |
| M-54 | `User` | `u0` with `("tags", ("x", ("y", ())))` | yes | TS-SEQ-1 | same |
| M-55 | `User` | `u0` with `("createdAt", NotANumber)` | yes | TS-LIB-1 | same |
| M-56 | `User` | `u0` without `email` | no | TS-OBJ-2 | same |
| M-57 | `User` | `u0 + ("age", 30.0)` | no (whole number in `Fractional` form) | TS-NUM-1 | same |
| M-58 | `User` | `u0 + ("age", Nothing)` | no | TS-OBJ-3 | same |
| M-59 | `User` | `u0 + ("age", "x")` | no | TS-OBJ-3 | same |
| M-60 | `User` | `u0` with `("tags", ("x"))` | no | VM-SEQ-1 | same |
| M-61 | `User` | `u0 + (1, "x")` | no | TS-OBJ-1 | same |
| M-62 | `Shape` | `{("kind", "circle"), ("radius", 1)}` | yes | TS-OBJ-2 | same (the builder now says `kind in {"circle"}`, not a literal pattern) |
| M-63 | `Shape` | `{("kind", "rect"), ("w", 1), ("h", 2.5)}` | yes | TS-OBJ-2 | same |
| M-64 | `Shape` | `{("kind", "rect"), ("w", 1)}` | no | TS-OBJ-2 | same |
| M-65 | `Shape` | `{("kind", "square"), ("radius", 1)}` | no | TS-OBJ-2 | same |
| M-66 | `Shape` | `{("kind", circle), ("radius", 1)}` | no | VM-LEX-3 | same |
| M-67 | `Headers` | `{("contentType", "a"), ("x-id", "1")}` | yes | TS-OBJ-4 | same |
| M-68 | `Headers` | `{("contentType", "a"), ("x-n", 1)}` | no | TS-OBJ-4 | same |
| M-69 | `Headers` | `{}` | no | TS-OBJ-2 | same |
| M-70 | `Flags` | `{("read", true), ("write", false)}` | yes | TS-SCOPE-1, TS-OBJ-2 | same |
| M-71 | `Flags` | `{("read", true)}` | no | TS-OBJ-2 | same |
| M-72 | `Flags` | `{("read", "true"), ("write", false)}` | no | VM-AX-4 | same |
| M-73 | `Account` | `{("id", "a"), ("balance", 0)}` | no | TS-OBJ-2 | same |
| M-74 | `Account` | `{("id", "a"), ("balance", 0), ("onChange", (callable, P, Unset))}`, with `P` = `{ (n, ()) suchthat n in Double }` | yes by the builder; a target never checks the `onChange` field (VM-PRES-3) | TS-FN-1 | changed (was "can't be written": the function domain was an opaque axiom, and is now a tagged roster with a placeholder) |
| M-75 | `AuditedUser` | `u0 + ("updatedAt", 0)` | yes | TS-INTER-1 | same |
| M-76 | `AuditedUser` | `u0` | no | TS-INTER-1 | same |

### 3.4 Tuples and sequences

| ID | Domain | Value | Member? | Basis | vs. it2 |
| -- | ------ | ----- | ------- | ----- | ------- |
| M-80 | `Point` | `(1, (2.5, ()))` | yes | TS-TUP-1 | same |
| M-81 | `Point` | `(1.0, 2.0)` | no | TS-TUP-1 | same |
| M-82 | `Point` | `(1, (2, (3, ())))` | no | TS-TUP-1 | same |
| M-83 | `Point` | `(1, ())` | no | TS-TUP-1 | same |
| M-84 | `Span` | `(1, ())`, `(1, (2, ()))`, `(1, (unset, ()))` | yes | TS-TUP-1 | same (the atom was `undefined`) |
| M-85 | `Span` | `()` | no | TS-TUP-1 | same |
| M-86 | `Args` | `("a", ())`, `("a", (1, (2, ())))` | yes | TS-TUP-1 | same |
| M-87 | `Args` | `("a")`, `(1, ())` | no | TS-TUP-1 | same |
| M-88 | `Ambiguous` | `(1, (2, ()))` | yes, one encoding | TS-UNION-2 | same (the domain is now written `Seq(Double)`: the tuple branch adds no members) |
| M-89 | `Seq(Double)` | `(1, (2, ()))` (a `Point` value) | yes: `Point` is inside `Seq(Double)` | TS-TUP-1 | same |

### 3.5 Maps, sets, generics, recursion

| ID | Domain | Value | Member? | Basis | vs. it2 |
| -- | ------ | ----- | ------- | ----- | ------- |
| M-90 | `Dictionary` | `{}`, `{("a", 1)}` | yes | TS-OBJ-4 | same |
| M-91 | `Dictionary` | `{("a", "x")}`, `{(1, 1)}` | no | TS-OBJ-4 | same |
| M-92 | `ByIndex` | `{("0", "a")}` | yes | TS-OBJ-4 | same |
| M-93 | `ByIndex` | `{(0, "a")}`, `{("x", "a")}` | `{(0, "a")}` no; `{("x", "a")}` **yes** | TS-OBJ-4 | changed (`{("x", "a")}` was no: a number index signature constrains only keys that are number text, and other keys stay open, VM-REC-3) |
| M-94 | `TagSet` | `{}`, `{"a", "b"}` | yes | TS-LIB-1 | same |
| M-95 | `TagSet` | `{("a", ())}`, `{1}` | no | TS-LIB-1 | same |
| M-96 | `Box(String)` | `{("value", "a")}` | yes | TS-GEN-1 | same |
| M-97 | `Box(String)` | `{("value", 1)}` | no | TS-GEN-1 | same |
| M-98 | `Result(Double, String)` | `{("ok", true), ("value", 1)}`, `{("ok", false), ("error", "e")}` | yes | TS-OBJ-2 | same (the builder now says `ok in {true}`) |
| M-99 | `Result(Double, String)` | `{("ok", true), ("error", "e")}`, `{("ok", "true"), ("value", 1)}` | no | TS-OBJ-2 | same |
| M-100 | `UserPage` | `{("items", (u0, ())), ("next", Nothing)}` | yes | TS-GEN-1 | same |
| M-101 | `UserPage` | `{("items", ({}, ())), ("next", Nothing)}` | no | TS-GEN-1 | same |
| M-102 | `UserPage` | `{("items", ())}` | no | TS-OBJ-2 | same |
| M-103 | `Tree(String)` | `{("value", "a"), ("children", ())}` | yes | TS-RECUR-1 | same |
| M-104 | `Tree(String)` | `{("value", "a"), ("children", ({("value", "b"), ("children", ())}, ()))}` | yes | TS-RECUR-1 | same |
| M-105 | `Tree(String)` | `{("value", "a"), ("children", ("b", ()))}` | no | TS-RECUR-1 | same |
| M-106 | `Json` | `"a"`, `1`, `1.5`, `true`, `Nothing`, `()`, `{("a", (1, ()))}` | yes | TS-RECUR-1 | same |
| M-107 | `Json` | `1.0`, `unset`, `("a")`, `{(1, "a")}` | no | TS-NUM-1 | same (the atom was `undefined`) |
| M-108 | `UserKey` | `"user-1"` | yes (prose) | TS-TPL-1 | same |
| M-109 | `UserKey` | `"User-1"` | no (prose) | TS-TPL-1 | same |
| M-110 | `Nullish` | `Nothing`, `unset` | yes | TS-MAP-1 | same (the atom was `undefined`) |
| M-111 | `Nullish` | `{}` | no | VM-AX-5 | same |
| M-112 | `EmptyThings` | every `Anything` | yes | TS-MAP-1 | same (the domain is now written `Anything`, so a target needs no set algebra; but see M-121) |
| M-112a | `EmptyThings` | `Circle` | no | TS-PRE-1 | same |
| M-113 | `SetOrDict` | `{"a"}` | yes, `FiniteSet` branch only | TS-LIB-1 | same |
| M-114 | `SetOrDict` | `{("a", ())}` | yes, `FiniteMap` branch only | TS-LIB-1 | same |
| M-115 | `SetOrDict` | `{}` | yes, both branches | TS-WIRE-1 | same (now stated in an `Overlap:` note) |

## 4. New vectors

### 4.1 Projection

`P` = `{ (n, ()) suchthat n in Double }`, the argument list of `Account.onChange`.

| ID | JavaScript value (static type) | `.sets` value | Basis | What it shows |
| -- | ------------------------------ | ------------- | ----- | ------------- |
| P-45 | `[ , 1]` (`(number \| undefined)[]`) | `(unset, (1, ()))` | TS-SEQ-1 | A hole is `unset`. |
| P-46 | `(n: number) => {}` (in `Account.onChange`) | `(callable, P, Unset)` | TS-FN-1 | A function is the placeholder of its `Callable` domain. |
| P-47 | `new Account()` | `{("id", ""), ("balance", 0), ("onChange", (callable, P, Unset))}` | TS-OBJ-1, TS-OBJ-5 | A class instance is a record; the private field is data; the method isn't there. |
| P-48 | `new Map([[1, [2]]])` | `{(1, (2, ()))}` | TS-LIB-1 | The same value as P-43. |
| P-49 | `"\u0001"` | `"\u{1}"` | TS-STR-1 | A control character is always an escape. |
| P-50 | `new Date(8.64e15)` | `8640000000000000` | TS-LIB-1 | The last valid instant, in `Integer` form. |
| P-51 | `{ a: {} }` | `{("a", {})}` | TS-OBJ-1 | Not equal to `new Set(["a"])` (E-10). |

### 4.2 Equality

| ID | Left | Right | Equal? | Basis | What it shows |
| -- | ---- | ----- | ------ | ----- | ------------- |
| E-17 | P-43, `new Set([[1, 2]])` | P-48, `new Map([[1, [2]]])` | yes; **not justified** | TS-WIRE-1, L9 | A set of sequences is a map. This is a collision of different data that the standard can't remove. |
| E-18 | `{ age: undefined }` → `{("age", unset)}` | `{}` → `{}` | no | TS-NIL-1 | A property holding `undefined` isn't a missing property. |
| E-19 | `[ , ]` → `(unset, ())` | `[undefined]` → `(unset, ())` | yes | TS-WIRE-1 | A hole and `undefined` are one value. |
| E-20 | P-47, `new Account()` | the plain object with the same three properties | yes, justified | TS-WIRE-1 | The difference is behavior. |
| E-21 | one function in `Account.onChange` | another function of the same type | yes: both are `(callable, P, Unset)` | TS-WIRE-1 | A function isn't data. |
| E-22 | the domain `One` | the domain `OneN` | yes, both are `{1}` | TS-NUM-2 | `1` and `1n` are one value, so the two literal types are one domain. |
| E-23 | the domain `Dictionary` | the domain `Scores` | yes, both are `FiniteMap(String, Double)` | TS-LIB-1 | A `Map` and an index signature hold the same pairs. |

### 4.3 Membership

| ID | Domain | Value | Member? | Basis |
| -- | ------ | ----- | ------- | ----- |
| M-116 | `WholeFractional` | `1.0` | yes | TS-PRE-1 |
| M-116a | `WholeFractional` | `1`, `1.5` | no | TS-PRE-1 |
| M-117 | `Double` | `0`, `-7`, `2.5` | yes | TS-NUM-1 |
| M-117a | `Double` | `0.0`, `true`, `unset` | no | TS-NUM-1 |
| M-118 | `Finite` | `{}`, `{1, 2}`, `{("a", 1)}`, `()` | yes (prose) | TS-PRE-1 |
| M-118a | `Finite` | `Integer`, `1`, `(1, 2)`, `Nothing` | no (prose) | TS-PRE-1 |
| M-119 | `TagSet` | `String` (the domain itself) | no: it's in `SetOf(String)` but isn't finite | VM-SET-1 |
| M-120 | `NumericKey` | `"0"`, `"-7"`, `"1500"`, `"0.001"`, `"1.5e-7"`, `"Infinity"`, `"9007199254740992"` | yes (prose) | TS-PRE-1 |
| M-120a | `NumericKey` | `"1.0"`, `"1e21"`, `"+1"`, `".5"`, `" 1"`, `""`, `1` | no (prose) | TS-PRE-1 |
| M-121 | `Anything` | `(Circle, ())`, `{Circle}` | yes: what is nested isn't checked (L3) | TS-PRE-1 |
| M-122 | `Object` | `{("a", Circle)}` | yes (L3) | VM-REC-3 |
| M-123 | `ByIndex` | `{}`, `{("x", 5)}`, `{("1.5", "a"), ("x", 5)}` | yes | TS-OBJ-4 |
| M-123a | `ByIndex` | `{("0", 5)}`, `{("1.5", 5)}`, `{("NaN", Nothing)}` | no | TS-OBJ-4 |
| M-124 | `Scores` | `{}`, `{("a", 1)}` | yes | TS-LIB-1 |
| M-124a | `Scores` | `{("a", "x")}`, `{"a"}` | no | TS-LIB-1 |
| M-125 | `Timestamped` | `{("updatedAt", 0)}`, `{("updatedAt", NotANumber), ("x", 1)}` | yes | TS-OBJ-1 |
| M-125a | `Timestamped` | `{("updatedAt", 0.5)}`, `{}` | no | TS-OBJ-2 |
| M-126 | `Pi` | `3.140000000000000124344978758017532527446746826171875` | yes | TS-LIT-1 |
| M-126a | `Pi` | `3.14` | no | VM-EQ-1 |
| M-127 | `Callable(P, Unset)` | `(callable, P, Unset)` | yes | TS-FN-1 |
| M-127a | `Callable(P, Unset)` | `callable`, `{}`, `(callable)` | no | TS-FN-1 |
| M-128 | `Anything` | `(callable, P, Unset)` | no: a runtime value isn't data | VM-PRES-3 |
| M-129 | `UserKey` | `"user-"`, `"user-a b"` | yes (prose) | TS-TPL-1 |
| M-129a | `UserKey` | `"user"`, `"xuser-1"`, `""`, `user-1` | no (prose) | TS-TPL-1 |
| M-130 | `Headers` | `{("contentType", "a"), (1, "b")}` | no: a key that isn't a `String` | TS-OBJ-4 |
| M-131 | `Composite` | `u0` | yes, so `User` satisfies the bound of `Page` | TS-GEN-1 |
| M-132 | `Span` | `(1, (Nothing, ()))`, `(1, (2, (3, ())))`, `(unset, ())` | no | TS-TUP-1 |
| M-133 | `SetOrDict` | `{("a", 1)}` | yes, `FiniteMap` branch | TS-OBJ-4 |
| M-133a | `SetOrDict` | `{("a", Nothing)}`, `{"a", ("b", 1)}`, `{1}` | no | TS-MAP-1 |
| M-134 | `Count` | `2.5` | yes | TS-MAP-1 |
| M-134a | `Count` | `"2.5"` | no | VM-KIND-1 |
| M-135 | `Id` | `"u1"` | yes | TS-MAP-1 |
| M-135a | `Id` | `1` | no | VM-KIND-1 |
| M-136 | `Json` | `{}`, `{("a", {})}`, `NotANumber` | yes; `{}` is in the `Seq` branch and the `FiniteMap` branch | TS-WIRE-1 |
| M-136a | `Json` | `{"a"}` | no | TS-RECUR-1 |
| M-137 | `Args` | `("a", (1, (unset, ())))` | no | TS-TUP-1 |
| M-138 | `EpochMillis` | `-8640000000000000`, `-1` | yes | TS-PRE-1 |
| M-138a | `EpochMillis` | `-8640000000000001`, `NegativeZero` | no | TS-PRE-1 |

## 5. Summary

| Group | Vectors | Same | Changed | Not applicable |
| ----- | ------- | ---- | ------- | -------------- |
| Projection (P-01…P-44) | 44 | 42 | 2 (P-20, P-24: atom renamed) | 0 |
| Equality (E-01…E-16) | 16 | 16 | 0 | 0 |
| Membership (M-01…M-115) | 109 | 107 | 2 (M-74: function placeholder; M-93: open number index signature) | 0 |
| New | 7 `P-`, 7 `E-`, 38 `M-` | | | |
