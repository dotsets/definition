# TypeScript → `.sets` Source Standard (iteration 2)

- **Standard version:** 0.2.0 (second and final pass; see [REVIEW.md](REVIEW.md))
- **Based on:** `standards/core/VALUE-MODEL.md` 0.2.0 only. Rebuilt from scratch with iteration 1's findings (`tests/iteration1-TypeScript/FINDINGS.md`, F1–F12) in mind; iteration 1's standard isn't reused.
- **Status:** test draft
- **Prelude:** [`../sets/ts.sets`](../sets/ts.sets). Every output file assumes `base` and the prelude are in scope.

The output of this standard must be readable by a target adapter that knows only the value model and `.sets`. So everything TypeScript-specific that a target needs is written in `.sets` (the prelude), never only in this document.

## 1. View and scope

**TS-VIEW-1. Runtime view.** A TypeScript type maps to the domain of the JavaScript runtime values it admits, read as data. Assumed compiler settings: `strict: true`, `exactOptionalPropertyTypes: false`.

**TS-SCOPE-1. Input.** Type aliases, interfaces, enums, and classes' instance fields. Each named type becomes one named domain, and each generic type becomes a function. Out of scope: values, function bodies, modules, and symbol-keyed properties.

**TS-NAME-1. Names.** A domain keeps its TypeScript name. A name already declared in scope (VM-NAME-1) gets the prefix `Ts`. Anonymous types are written inline.

## 2. Values

This section fixes the projection from JavaScript values to `.sets` values. The domains in §3–§6 follow from it.

**TS-VAL-1. Strings.** A JS string is a `String` (VM-LEX-1). A string that holds a lone surrogate isn't a sequence of Unicode scalar values (VM §4), so it has no `.sets` value, and the projection reports it (VM-WIRE-1). *(F8)*

**TS-VAL-2. Numbers.** A JS `number` is written in its canonical form (VM-LEX-2):
- a finite whole number in `Integer` form: `1` → `1`, `1e21` → `1000000000000000000000`;
- any other finite number in `Fractional` form, as the exact decimal of the binary64 value: `0.5` → `0.5`, `0.1` → `0.1000000000000000055511151231257827021181583404541015625`;
- `NaN`, `Infinity`, `-Infinity`, `-0` → `NotANumber`, `PositiveInfinity`, `NegativeInfinity`, `NegativeZero` (VM-AX-3).

Every binary64 whole number has an `Integer` form in `Binary64Finite`, but large ones have no `Fractional` form. So `Integer` form is the only choice that covers all of them. *(F12, F3)*

**TS-VAL-3. BigInts.** A `bigint` is an `Integer`: `1n` → `1`. `1n` and `1` are therefore the same value. That collision is justified: both are the number one, and VM-AX-2 says the representation that holds a number isn't part of its value. The domains stay distinct (`bigint` → `Integer`, `number` → `TsNumber`). Where one domain admits both, the projection back to JS is ambiguous and reports it (TS-WIRE-1). *(revisits iteration 1's R1)*

**TS-VAL-4. Other primitives.** `true`/`false` → `true`/`false` (VM-AX-4). `null` → `Nothing` (VM-AX-5). `undefined` → the atom `undefined`. A `symbol` is a member of the axiom `TsSymbol`, and a function is a member of `TsFunction` (VM-AX-7). Neither has a written notation.

**TS-VAL-5. Objects.** A plain object is a record of its own string-keyed properties (VM-REC-1). Keys are strings, so `{1: "a"}` → `{("1", "a")}`.

**TS-VAL-6. Arrays.** A JS array is always a sequence (VM-SEQ-1), whatever its static type: `[1, 2]` → `(1, (2, ()))`. At runtime a JS array is variable-length (even a tuple-typed one can be `push`ed), so by VM-SEQ-1 it isn't a fixed-length product. Tuple types constrain positions in the sequence (TS-TUP-1). *(F1)*

**TS-VAL-7. Library objects** narrow to the data they wrap (VM-PRES-2):
- A `Set` is the set of its members (VM-SET-1, as a `SetOf` value): `new Set(["a"])` → `{"a"}`. *(F2, F9)*
- A `Map` is the map of its entries (VM-MAP-1): `new Map([[1, "a"]])` → `{(1, "a")}`.
- A `Date` is its time value, a `number`: `new Date(0)` → `0`.

Insertion order of `Map`, `Set` and object keys isn't kept: domains have no order (VM-EQ-1).

## 3. Primitive types

**TS-PRIM-1.**

| TypeScript | `.sets` |
| ---------- | ------- |
| `string` | `String` |
| `number` | `TsNumber` |
| `bigint` | `Integer` |
| `boolean` | `Boolean` |
| `null` | `{Nothing}` |
| `undefined`, `void` | `TsUndefined` |
| `symbol` | `TsSymbol` |
| `any`, `unknown` | `TsValue` (every JS value) |
| `never` | `Phi` |
| `object` | `TsObject` |
| `{}` | `TsNonNullish` |
| any function type | `TsFunction` |

`TsValue` bounds the top types by what JS can hold, not by `U`. *(F6)*

The prelude is normative. In one line each:
- `TsNumber` is `Float64` without the whole-valued `Fractional`s (TS-VAL-2).
- `TsUndefined` is `{undefined}`.
- `TsRecord` is `Map(String, U)`.
- `TsNumericKey sub String` holds the strings `String(n)` produces.
- `TsDate` is the whole `TsNumber`s within ±8.64e15, plus `NotANumber`.
- `TsValue` is the recursive union of the primitives, `Map(TsValue, TsValue)`, `Seq(TsValue)`, `SetOf(TsValue)` and `TsFunction`.
- `TsObject` is `TsValue` without the primitives, and `TsNonNullish` is `TsValue` without `Nothing` and `undefined`.

**Known limit.** A `Date` value is a number (TS-VAL-7), so it isn't in `TsObject`, though TS lets `Date` stand for `object`. *(R2)*

**TS-LIT-1. Literal types** are singleton rosters of their TS-VAL value, spelled with VM-LEX-1 escapes (`"a, b"` → `{"a\, b"}`): `"a"` → `{"a"}`, `1` → `{1}`, `0.5` → `{0.5}`, `1n` → `{1}`, `true` → `{true}`. A union of literals is one roster.

**TS-ENUM-1.** An enum is the roster of its member values. Member names aren't data and are dropped.

**TS-TMPL-1.** A template literal type is a described subdomain of `String` with a prose extent (VM-SUB-2), written as a comment directly above it.

## 4. Composition

**TS-UNION-1.** `A | B | C` → `union(A, union(B, C))`. `A & B` → `intersect(A, B)`.

**TS-GEN-1.** `type F<T> = …` → `F({T sub U}): …`. `T extends C` → `{T sub C}`. `F<X>` → `F(X)`. A default type argument can't be expressed, so every use writes it out. *(F11)*

**TS-GEN-2.** Mapped, conditional, `keyof` and `typeof` types are resolved by the checker first, or they're out of scope.

**TS-RECUR-1.** Recursive types become recursive domains or functions. A recursive function is read per argument, as the least fixed point (VM-RECUR-1). *(F10)*

## 5. Objects

**TS-OBJ-1. Open records with string keys.** An object type is an open record (VM-REC-3) whose binder is constrained by `r in TsRecord` (`Map(String, U)`) instead of `r in Record`. This is the only override (§8). Each property `k: T` is a required field: `("k", k) in r and k in T`. *(F7)*

**TS-OBJ-2. Optional properties.** `k?: T` → `r in OptionalField("k", union(T, TsUndefined))`.

**TS-OBJ-3.** `readonly` is dropped: it's not a property of the domain.

**TS-OBJ-4. Index signatures.** `[k: string]: V` → the clause `r in Map(String, V)`. `[k: number]: V` → `r in Map(TsNumericKey, V)`. Alone, either one is just the `Map` call. `Record<K, V>` with a literal-union `K` is a record that requires each key.

**TS-CLASS-1.** A class is the open record of its instance fields, `private` included. Methods live on the prototype and are dropped.

## 6. Arrays, tuples, wrappers

**TS-ARR-1.** `T[]`, `Array<T>`, `readonly T[]` → `Seq(T)`.

**TS-TUP-1. Tuples are sequence patterns.** By TS-VAL-6 a tuple value is a sequence, so a tuple type is a builder over a sequence pattern:
- `[A, B]` → `{ (a, (b, ())) suchthat a in A and b in B }`
- `[A, B?]` → `union({ (a, ()) suchthat a in A }, { (a, (b, ())) suchthat a in A and b in union(B, TsUndefined) })`
- `[A, ...B[]]` → `{ (a, rest) suchthat a in A and rest in Seq(B) }`

A rest element anywhere but last (`[...A[], B]`) is out of scope. Every tuple domain is then inside `Seq(U)`, and `[number, number] | number[]` has one encoding per value. *(F1)*

**TS-LIB-1.**

| TypeScript | `.sets` |
| ---------- | ------- |
| `Set<T>` | `SetOf(T)`. A JS set is finite; VM-SET-1 asks for that narrowing, but the grammar can't state it (VM SI-14). |
| `Map<K, V>` | `Map(K, V)` |
| `Date` | `TsDate` |

## 7. Projection back to JavaScript

**TS-WIRE-1.** The projection in §2 is one-to-one except in these places, which report a diagnostic when a domain admits both sides (VM-WIRE-1):
- `number` vs. `bigint`, for whole values (TS-VAL-3);
- `Date` vs. `number` (TS-VAL-7);
- `Map` vs. plain object, for string keys (TS-VAL-7);
- a class instance vs. a plain object with the same fields;
- the empty object, array, `Set` and `Map`, which are all `{}` (VM-REC-1 allows this).

## 8. Overrides

| VM rule | Override |
| ------- | -------- |
| VM-REC-2 (Default) | Record binder `r in TsRecord` instead of `r in Record`. |
