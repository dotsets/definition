# TypeScript → `.sets` Source Standard (iteration 1)

- **Standard version:** 0.2.0 (second and final pass; see [REVIEW.md](REVIEW.md))
- **Based on:** `standards/core/VALUE-MODEL.md` 0.2.0 only. No other standard is used.
- **Status:** test draft
- **Prelude:** [`ts.sets`](ts.sets), which every output file assumes is in scope together with `base`

This standard says how TypeScript type declarations become `.sets` domains. Every rule is subordinate to the value model (VM). Where this standard is silent, the VM decides.

## 1. Scope and view

**TS-VIEW-1. Runtime view.** A TypeScript type maps to the domain of JavaScript runtime values it admits, read as data. The compiler settings assumed are `strict: true` and `exactOptionalPropertyTypes: false`.

**TS-SCOPE-1. Input.** Type aliases, interfaces, enums, and the instance fields of classes. Each named type becomes one named domain; generic types become functions. Out of scope: function bodies, values, modules, and symbol-keyed properties.

**TS-NAME-1. Names.** A domain keeps its TypeScript name. If the name is already declared in scope (VM-NAME-1, for example `Record`, `String`, `Map`), it's prefixed with `Ts`. Anonymous object types are written inline as builders.

## 2. Primitives

**TS-PRIM-1.**

| TypeScript | `.sets` |
| ---------- | ------- |
| `string` | `String` |
| `number` | `TsNumber` (§3) |
| `bigint` | `Integer` |
| `boolean` | `Boolean` |
| `null` | `{Nothing}` |
| `undefined`, `void` | `TsUndefined` = `{undefined}` |
| `symbol` | `TsSymbol` (new axiom, VM-AX-7) |
| `any`, `unknown` | `U` |
| `never` | `Phi` |
| `object` | `TsObject` (everything that isn't a primitive) |
| `{}` | `TsNonNullish` (everything except `null` and `undefined`) |

`null` is the core `Nothing` (VM-AX-5). JavaScript has a second kind of absence, `undefined`. VM-PRES-1 requires it be kept, so it's the atom `undefined` (VM-MISS-1 leaves it to this standard).

## 3. Numbers

**TS-NUM-1. Number values.** A JavaScript `number` is an IEEE-754 binary64 value, and it's always written in `Fractional` form, as the exact decimal expansion of the binary value (VM-AX-2):
- `3` is `3.0`; `0.5` is `0.5`; `0.1` is `0.1000000000000000055511151231257827021181583404541015625`.
- `NaN`, `Infinity`, `-Infinity` and `-0` are `NotANumber`, `PositiveInfinity`, `NegativeInfinity` and `NegativeZero` (VM-AX-3).

`Integer` is reserved for `bigint`, so `1` and `1n` stay different values (VM-PRES-1). `TsNumber` is `Float64` without its `Integer` members.

**TS-NUM-2. Literal types.** A number literal type is the singleton roster of its TS-NUM-1 value: `1` → `{1.0}`.

## 4. Literals, enums, unions, intersections

**TS-LIT-1.** `"a"` → `{"a"}`, `true` → `{true}`, `1n` → `{1}`. A union of literals is one roster: `"a" | "b"` → `{"a", "b"}`.

**TS-ENUM-1.** An enum is the roster of its runtime member values: `enum L { Low, High }` → `{0.0, 1.0}`. Member names aren't data and are dropped.

**TS-UNION-1.** `A | B | C` → `union(A, union(B, C))` (binary, right-nested). `A & B` → `intersect(A, B)`.

**TS-TMPL-1.** A template literal type is a described subdomain `X sub String` with a prose extent (VM-SUB-2).

## 5. Objects

**TS-OBJ-1. Object types are open records with property keys.** TypeScript object types are structural and admit extra properties, which matches VM-REC-3. JavaScript object keys are strings (symbols are out of scope), so the record binder is constrained by `r in TsRecord` instead of `r in Record` (override of VM-REC-2). Each property `k: T` is a required field:

```sets
User = { r suchthat r in TsRecord and ("id", id) in r and id in String }
```

A numeric property name is a string key: `{ 1: "a" }` is `{("1", "a")}`. TypeScript object types never emit `OnlyKeys`.

**TS-OBJ-2. Optional properties.** `k?: T` → `r in OptionalField("k", union(T, TsUndefined))`, because without `exactOptionalPropertyTypes` a present `undefined` is allowed. A missing property has no tuple (VM-MISS-1).

**TS-OBJ-3. `readonly`** is a mutability marker, not a domain property, and is dropped.

**TS-OBJ-4. Index signatures.** `{ [k: string]: V }` and `Record<string, V>` → `Map(String, V)`. `{ [k: number]: V }` → `Map(TsNumericKey, V)`, because the runtime keys are strings. `Record<K, V>` with a literal union `K` is a record with each key required. Combined with named properties, the map is a clause on `r` next to the property clauses.

**TS-CLASS-1. Classes.** A class is the open record of its declared instance fields, `private` ones included. Methods live on the prototype, not the instance, and are dropped. A field whose type is a function type maps to `TsFunction` (a new axiom).

## 6. Arrays and tuples

**TS-ARR-1.** `T[]`, `Array<T>` and `readonly T[]` → `Seq(T)` (VM-SEQ-1).

**TS-TUP-1.** A fixed tuple `[A, B]` → `{ (a, b) suchthat a in A and b in B }` (VM-TUP-1). An optional element `[A, B?]` gives a union of each length, and the element takes `union(B, TsUndefined)`. A rest element `[A, ...B[]]` makes the tuple variable-length, so it's a sequence: `{ (a, rest) suchthat a in A and rest in Seq(B) }`.

## 7. Generics and recursion

**TS-GEN-1.** `type Box<T> = …` → `Box({T sub U}): …`. A constraint `T extends C` → `{T sub C}`. An instantiation `Box<string>` → `Box(String)`.

**TS-GEN-2.** Mapped, conditional, `keyof` and `typeof` types are resolved by the TypeScript checker first. If they can't be resolved to the forms above, they're out of scope.

**TS-REC-1.** Recursive types map to recursive domains and functions (VM-RECUR-1).

## 8. Library wrappers

**TS-LIB-1.** Following VM-PRES-2, wrappers narrow to what they wrap:

| TypeScript | `.sets` |
| ---------- | ------- |
| `Set<T>` | `Map(T, {()})`: keys are the members, values are empty |
| `Map<K, V>` | `Map(K, V)` |
| `Date` | `TsDate`: its time value, a whole-valued `TsNumber` within ±8.64e15, or `NotANumber` |

## 9. Projection

**TS-WIRE-1.** The projection from JavaScript values to `.sets` values follows §2–§8. When a domain admits two `.sets` values for one JavaScript value, it reports a diagnostic (VM-WIRE-1). The known case is a union of a fixed tuple and an array: `[1, 2]` is both `(1.0, 2.0)` and `(1.0, (2.0, ()))`.

## 10. Overrides

| VM rule | Override |
| ------- | -------- |
| VM-REC-2 (Default) | The record binder uses `r in TsRecord` (`Map(TsKey, U)`) instead of `r in Record`, so keys are strings. |
