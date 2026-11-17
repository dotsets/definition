# TypeScript → `.sets` standard (iteration 3)

- **Direction:** source (TypeScript → `.sets`)
- **Built on:** [`VALUE-MODEL.md`](../../../standards/core/VALUE-MODEL.md) 0.2.0 and nothing else. `VM-*`, `SI-*` and `UP-*` tags refer to it.
- **Output:** [`sets/prelude.sets`](../sets/prelude.sets), the domains this standard adds to `base`, and one `.sets` file per input file.
- **Status:** test iteration, written in two passes. This is the text after pass 2 (see [REVIEW.md](REVIEW.md)).

Rules are tagged `TS-*`. The types in the examples are from [`input/types.ts`](../input/types.ts).

## 1. Scope and output

**TS-SCOPE-1. What is transpiled.** Type declarations: `type`, `interface`, `enum` and `class`. The adapter works on types as the checker resolves them under `strict`. A reference to a declared type stays a name. A utility, mapped, conditional or indexed-access type is replaced by what it resolves to. Statements, expressions and function bodies aren't transpiled.

**TS-SCOPE-2. The view is the value in memory.** A domain holds the values a variable of the type can hold at run time, read as data at rest. It isn't the JSON form of the value.

**TS-SCOPE-3. No rule, no guess.** A construct with no rule here isn't transpiled, and a value with no `.sets` value isn't projected. Either one is a *diagnostic*: a report to whoever runs the adapter. §10 lists the common cases.

**TS-OUT-1. The output stands alone.** A target reads only `sets/` and the value model. So:
- nothing in `sets/` names the source language: no domain, atom, comment or prose extent (a declared type keeps its own name, whatever it is);
- prelude names are concept names (`Double`, `Unset`), never source spellings;
- prose is written with `.sets` values and domains only;
- whatever a target needs is in the `.sets` text: as a definition where the grammar can say it, as a comment where it can't.

**TS-OUT-2. Names.** `type X = T`, `interface X`, `enum X` and `class X` are each written `X = …`, and a generic one as a function `X(…): …` (TS-GEN-1). If the name is already declared in `base` or the prelude, the adapter appends `_` and says so in a comment (VM-NAME-1). A binder is named after its property or tuple label; unlabelled tuple positions are `e0`, `e1`, …, and a rest element is `rest`. Where that spelling is a keyword or a name in scope, any fresh name is used: binder names carry no meaning (VM-REC-1).

**TS-OUT-3. Notes.** The adapter writes these comments in the output. Comments carry no meaning in `.sets` (SI-15); they're for the target's author.

| Note | Says |
| ---- | ---- |
| `Overlap:` | A case of TS-WIRE-1 applies here: the domain admits two source values that are one `.sets` value, or it's the same domain as another declared type that the source keeps apart. |
| `Names:` | An enum's member names, in roster order. |
| `Default:` | A default type argument (TS-GEN-1). |
| `Not data:` | A member that was dropped because it's behavior (TS-OBJ-5). |
| `Extent:` | The exact members of a described subdomain. |

**TS-OUT-4. Imports.** An output file names every prelude domain it uses on a line `X in prelude`, the `X in file` form of VM-NAME-1. `base` names need no line.

## 2. Overrides

| VM rule | VM default | This standard | Why |
| ------- | ---------- | ------------- | --- |
| VM-REC-2 | A record builder starts with `r in Record`. | It starts with `r in Object`, or with `r in FiniteMap(String, V)` under TS-OBJ-4. | Property keys are strings and an object is finite. `Object` is inside `Record`, so the VM's requirement still holds. |
| VM-SUB-2 | Described subdomains are declared in the standard's own superdomain. | A template literal type adds a described subdomain to the output file (TS-TPL-1). | Its extent comes from the input, not from this standard. |

No other Default rule is overridden, and no axiom is added (VM-AX-7).

## 3. Prelude

**TS-PRE-1.** `sets/prelude.sets` declares these domains and nothing else. In the file, each one has a comment that tells a target in neutral words what it means and how it's used (TS-OUT-1).

| Domain | Definition | Holds |
| ------ | ---------- | ----- |
| `WholeFractional` | `{ x suchthat x in Fractional and n in Integer and x >= n and x <= n }` | A `Fractional` with no fractional part: `1.0`. |
| `Double` | `{ x suchthat x in Float64 and x not in WholeFractional }` | Each binary64 value, in one form (TS-NUM-1). |
| `Unset` | `{unset}` | The atom for `undefined` (TS-NIL-1). |
| `EpochMillis` | `union({ x suchthat x in Integer and x >= -8640000000000000 and x <= 8640000000000000 }, {NotANumber})` | The time value of a `Date` (TS-LIB-1). |
| `Finite` | `Finite sub U`, described | Every domain with finitely many members. |
| `FiniteSet(T)` | `intersect(SetOf(T), Finite)` | A finite set (VM-SET-1). |
| `FiniteMap(K, V)` | `intersect(Map(K, V), Finite)` | A finite map (VM-MAP-1). |
| `Object` | `FiniteMap(String, U)` | A finite record with `String` keys (TS-OBJ-1). |
| `NumericKey` | `NumericKey sub String`, described | The text of a `Double` as a property key (TS-OBJ-4). |
| `Composite` | `union(Seq(U), Finite)` | Every sequence, set, map and record. |
| `Something` | `union(String, union(Double, union(Integer, union(Boolean, Composite))))` | Every value except the two nothings. |
| `Anything` | `union(Something, union({Nothing}, Unset))` | Every value. |
| `Callable(P, R)` | `{(callable, P, R)}` | A runtime domain for functions (TS-FN-1). |

The two described subdomains have these extents. The prelude states them in `.sets` terms, precisely enough to write a membership check from (VM-SUB-2). They overlap other domains as their extents say.
- `Finite`: a value is a member exactly when it's a domain with finitely many members.
- `NumericKey`: a `String` is a member exactly when it's the text that `Number::toString` of ECMA-262 gives for some `number`. `"1.5"`, `"NaN"` and `"1e+21"` are members; `"-0"`, `"01"` and `"x"` aren't.

The prelude uses these `base` names as the value model describes them: `U`, `union`, `intersect`, `String`, `Integer`, `Fractional`, `Boolean`, `Nothing`, `Float64`, `Seq`, `SetOf`, `Map`, and the atom `NotANumber`. Output files also use `Phi`, `OptionalField` and `OptionalFields`.

## 4. The mapping

**TS-MAP-1.** Each TypeScript type has one domain, and each run-time value has one `.sets` value. A type variable in the Domain column stands for the domain of that type.

| TypeScript type | Domain | A value | Its `.sets` value |
| --------------- | ------ | ------- | ----------------- |
| `string` | `String` | `"a, b"` | `"a\, b"` |
| `number` | `Double` | `1`, `0.5`, `NaN` | `1`, `0.5`, `NotANumber` |
| `bigint` | `Integer` | `1n` | `1` |
| `boolean` | `Boolean` | `true` | `true` |
| `null` | `{Nothing}` | `null` | `Nothing` |
| `undefined`, `void` | `Unset` | `undefined` | `unset` |
| a literal type | a roster (TS-LIT-1) | `"active"` | `"active"` |
| `T[]` | `Seq(T)` | `["a", "b"]` | `("a", ("b", ()))` |
| a tuple type | a sequence shape (TS-TUP-1) | `[1, 2]` | `(1, (2, ()))` |
| an object type | a record builder (TS-OBJ-1) | `{ name: "a" }` | `{("name", "a")}` |
| `Set<T>` | `FiniteSet(T)` | `new Set(["a", "b"])` | `{"a", "b"}` |
| `Map<K, V>` | `FiniteMap(K, V)` | `new Map([["a", 1]])` | `{("a", 1)}` |
| `Date` | `EpochMillis` | `new Date(0)` | `0` |
| a function type | `Callable(P, R)` (TS-FN-1) | a function | the placeholder `(callable, P, R)` |
| `unknown`, `any`, `{} \| null \| undefined` | `Anything` | | |
| `{}` | `Something` | | |
| `object` | `Composite` | | |
| `never` | `Phi` | | |

## 5. Scalars

**TS-STR-1. Strings.** A string is the `String` with the same characters. Inside the quotes, `,`, `{`, `}`, `"` and `\` are escaped with a backslash, and a control character is written `\n`, `\r`, `\t` or `\u{H}` (VM-LEX-0, VM-LEX-1, VM-LEX-6): `"a\, b"`, `"say \"hi\""`, `"\{a\}"`, `"x\n"`. Property keys and enum values are written the same way. A string that holds a lone surrogate has no value, because `String` holds Unicode scalar values only (VM §4): diagnostic.

**TS-NUM-1. Numbers.** A `number` is a binary64 value, written exactly (VM-AX-2):
- a whole number in `Integer` form: `1`, `0`, and `1000000000000000000000` for `1e21`;
- any other finite number in `Fractional` form, as its full decimal expansion: `0.5`, and `0.1000000000000000055511151231257827021181583404541015625` for `0.1`;
- `NaN`, `Infinity`, `-Infinity` and `-0` as the atoms `NotANumber`, `PositiveInfinity`, `NegativeInfinity` and `NegativeZero` (VM-AX-3).

So one number has one spelling, and `Double` holds exactly these values: `1` is a member and `1.0` isn't. Whole numbers take the `Integer` form because, by the `Binary64Finite` builder, a whole number from 2^52 up has no `Fractional` form. No exponent notation is used (VM-LEX-2a).

`number` has no integer type, so `Double` is one combined domain (VM-KIND-2, UP-1). The `Integer` form is a spelling, not a type.

**TS-NUM-2. BigInts.** A `bigint` is the same number in `Integer`. `1n` and `1` are one value, `1`: a number has no representation (VM-AX-2). The types stay apart: `bigint` is `Integer` and `number` is `Double`.

**TS-NIL-1. Null and undefined.** `null` is `Nothing` (VM-AX-5). `undefined` is the atom `unset`. It's an ordinary atom, the one member of the prelude's `Unset`; VM-AX-5 lets a language standard declare its own null-ish atoms. It isn't a second member of `Phi`. `void` is read as `undefined`. The three ways a property can lack a value stay apart (VM-MISS-1):

```text
{ name: "a" }                    is   {("name", "a")}
{ name: "a", age: undefined }    is   {("name", "a"), ("age", unset)}
{ name: "a", age: null }         is   {("name", "a"), ("age", Nothing)}
```

**TS-LIT-1. Literal types.** A literal type is the roster of its value: `"active"` is `{"active"}`, `1` is `{1}`, `1n` is `{1}`, `true` is `{true}`, and `0.1` is `{0.1000000000000000055511151231257827021181583404541015625}`.

**TS-ENUM-1. Enums.** An enum is the roster of its member values, with a `Names:` note: `enum Color { Red = "RED", Green = "GREEN" }` is `{"RED", "GREEN"}`, and `enum Level { Low, Mid, High }` is `{0, 1, 2}`.

## 6. Records

**TS-OBJ-1. An object type is an open record.** An object type, an interface or a class is a builder over `r in Object`, with one clause per member:

```sets
Timestamped = { r suchthat r in Object and ("updatedAt", updatedAt) in r and updatedAt in EpochMillis }
```

An object value is the record of its own string-keyed properties (VM-REC-1). Keys are `String`s, and a numeric key is its text: `{ 1: "a" }` is `{("1", "a")}`. Symbol-keyed and inherited properties aren't in the record. A class instance is read the same way.

The record is open (VM-REC-3) and no builder uses `OnlyKeys`: TypeScript types are structural, so a value may carry properties its type doesn't mention.

**TS-OBJ-2. Required property.** `k: T` is `("k", k) in r and k in T`. A property of literal type is written the same way, `("kind", kind) in r and kind in {"circle"}`, not as a literal pattern: a bare symbol such as `true` inside a rule would be a binder (VM-RULE-2).

**TS-OBJ-3. Optional property.** `k?: T` is `r in OptionalField("k", union(T, Unset))`: the property may be missing, or present and `undefined`. Under `exactOptionalPropertyTypes` it's `r in OptionalField("k", T)`.

**TS-OBJ-4. Index signatures.**
- `[k: string]: V` covers every key, so the first clause becomes `r in FiniteMap(String, V)`. A type with nothing else is `FiniteMap(String, V)` itself: `{ [k: string]: number }` is `FiniteMap(String, Double)`.
- `[k: number]: V` covers the keys that are the text of a number: `r in OptionalFields(NumericKey, V)` (VM-MISS-2). Other keys stay open.

The form tells a target which keys are data and which are structure: a `FiniteMap` call is a table, and a builder names its fields.

**TS-OBJ-5. What isn't a field.** `readonly`, `public`, `private` and `protected` change nothing: the field is data either way. Methods, accessors, constructors and static members are behavior. They're dropped (VM-PRES-3) with a `Not data:` note. A property of function type is a field (TS-FN-1).

## 7. Sequences

**TS-SEQ-1. Arrays.** `T[]`, `Array<T>` and `readonly T[]` are `Seq(T)`. Every array value is a sequence (VM-SEQ-1), whatever its static type, and a hole is `unset`.

**TS-TUP-1. A tuple type is a sequence shape.** A tuple value is an array at run time, and TypeScript accepts a `[number, number]` wherever a `number[]` is expected. Written as a VM tuple `(1, 2)`, the one array `[1, 2]` would have two `.sets` values, depending on the type it's read through. So a tuple type is the domain of sequences of one shape. This reads "the source's type decides" in VM-SEQ-1 as: every array value has an array type. It isn't an override. The "language tuples" of VM-TUP-1 are tuples that are their own kind of value, which TypeScript's aren't.

- `[A, B]` is `{ (e0, (e1, ())) suchthat e0 in A and e1 in B }`.
- An optional element makes a union of two shapes: one that stops before it, and one that includes it and also admits `unset`. `[A, B?]` is `union({ (e0, ()) suchthat e0 in A }, { (e0, (e1, ())) suchthat e0 in A and e1 in union(B, Unset) })`.
- A rest element is the tail: `[A, ...B[]]` is `{ (e0, rest) suchthat e0 in A and rest in Seq(B) }`. With optional elements, only the shape that includes them all has the tail.

## 8. Combining types

**TS-UNION-1. Unions.** `A | B | C` is `union(A, union(B, C))`: binary and nested to the right (SI-8). Literal branches are written as one roster: `"active" | "suspended"` is `{"active", "suspended"}`. `true | false` is `Boolean`.

**TS-UNION-2. Simplify.** A branch that adds no members is left out, with a comment that says so:
- a branch equal to another branch (`undefined | void`);
- a literal branch whose basic type is also a branch (`1 | number` is `Double`);
- a tuple type whose every position is inside the element domain of an array type that is also a branch (`[number, number] | number[]` is `Seq(Double)`).

Nothing else is simplified.

**TS-INTER-1. Intersections.** `A & B` of object types is `intersect(A, B)`.

**TS-GEN-1. Generics.** A generic type `X<T>` is a function `X({T sub B}): …`, where `B` is the constraint's domain, or `U` when there's no constraint. `X<A>` is the call `X(A)`. A default type argument has no `.sets` form: it's written in a `Default:` note and filled in wherever a use omits it.

**TS-RECUR-1. Recursion.** A type may refer to itself (VM-RECUR-1). A cyclic value isn't a value: diagnostic.

## 9. Library types, template literals and functions

**TS-LIB-1. Wrappers narrow to what they wrap (VM-PRES-2).**

| TypeScript | Domain | Value |
| ---------- | ------ | ----- |
| `Set<T>`, `ReadonlySet<T>` | `FiniteSet(T)` | The domain of its members (VM-SET-1). |
| `Map<K, V>`, `ReadonlyMap<K, V>` | `FiniteMap(K, V)` | The domain of its `(key, value)` pairs (VM-MAP-1). |
| `Date` | `EpochMillis` | Its time value: a whole number of milliseconds since 1970-01-01T00:00:00Z, or `NotANumber` for an invalid date. |

VM-SET-1 leaves the reading of a set to the standard: a `SetOf` value, or a record of keys with empty values, which is the VM-PRES-2 example (`{("a", ()), ("b", ())}`). Both readings collide with something. A record of keys is equal to every record whose values are all empty, such as `{ a: {} }`. A `SetOf` value is equal to a map only when its members are non-empty sequences (TS-WIRE-1). This standard takes the `SetOf` reading because it collides less (VM-PRES-1).

**TS-TPL-1. Template literal types.** The grammar has no string operators (SI-14), so an alias of a template literal type is a described subdomain, `X sub String`. Only gaps of type `string` are covered. The `Extent:` note above it says: a `String` is a member exactly when it's the fixed pieces of text in order, with any `String`, the empty one too, in each gap. For `` `user-${string}` ``, that's every `String` that begins with `user-`.

**TS-FN-1. Function types.** A function isn't data (VM-PRES-3). A function type `(a: A) => R` is `Callable(P, R)`, where `P` is the parameter list read as the tuple type `[a: A]` (TS-TUP-1). `Callable` is a runtime domain: a tagged domain whose one member, `(callable, P, R)`, is a placeholder. It's never serialized or checked as data. The core `RUNTIME.md` wasn't available to this iteration, so `Callable` is this standard's own name for the concept.

## 10. Collisions, losses and what isn't covered

**TS-WIRE-1. One to one, or it says so (VM-WIRE-1).** §4 to §9 give each source value one `.sets` value. Two source values share a `.sets` value only in the cases below. Where a declared type admits both values, or two declared types differ only by one of these cases, the adapter reports it and writes an `Overlap:` note.

| Source values | `.sets` value | Why |
| ------------- | ------------- | --- |
| `1` and `1n` | `1` | The same number; a number has no representation (VM-AX-2). |
| `new Date(0)`, `0` and `0n`; an invalid `Date` and `NaN` | `0`; `NotANumber` | A `Date` wraps a number (VM-PRES-2). |
| `new Map([["a", 1]])` and `{ a: 1 }` | `{("a", 1)}` | A record is a map (VM-REC-1). |
| `[]`, `{}`, `new Set()` and `new Map()` | `{}` | Nothing inside (VM-AX-5, VM-REC-1). |
| A hole and `undefined` in an array | `unset` | Reading a hole gives `undefined`. |
| A class instance and a plain object with the same fields | the same record | The difference is behavior (VM-PRES-3). |
| Two functions of one type | the placeholder | A function isn't data (VM-PRES-3). |
| `new Set([[1, 2]])` and `new Map([[1, [2]]])` | `{(1, (2, ()))}` | **Not justified:** this is different data. A sequence is a pair (VM-SEQ-1), so a set of non-empty sequences with distinct first elements is also a map (VM-MAP-1). Open (L9). |

**TS-WIRE-2. What is lost.** Key order and insertion order (VM-MAP-1), and identity: two distinct objects with equal contents are one value. Where two keys of one `Map`, or two members of one `Set`, become equal `.sets` values, that's a diagnostic.

**TS-SCOPE-4. Not covered.** `symbol`; runtime types other than functions (`Promise<T>` is `Async(T)` by VM-PRES-3, but its declaration is in the unread `RUNTIME.md`); `WeakMap` and typed arrays; overloaded or generic function types; rest elements that aren't last; index signatures with other key types; template literal types with other gaps, or that aren't the whole of an alias; intersections of non-object types; and `#private` fields.

## 11. Known limits

| # | Limit |
| - | ----- |
| L1 | `number` gives one domain, so a target can't tell an integer from a fraction (UP-1). |
| L2 | A `Date` and a function are `object`s in TypeScript but aren't in `Composite`: a `Date` is a number, and a function isn't data. |
| L3 | `Anything`, `Something`, `Composite` and the unnamed fields of an `Object` are exact one level deep. What is nested inside them is `U`. |
| L4 | `Finite`, `NumericKey` and template literal types are prose only (SI-14). |
| L5 | A default type argument survives only as a note. |
| L6 | `Callable` is this standard's own runtime domain, because `RUNTIME.md` wasn't read (VM-PRES-3). |
| L7 | `base.sets` wasn't read. `base` names are used as the value model describes them, and a prelude name could clash with a `base` name the value model doesn't mention (VM-NAME-1). |
| L8 | The import form `X in prelude` is taken from VM-NAME-1's one mention of `X in file`. |
| L9 | The collision between a set of sequences and a map is unjustified and can't be fixed here: it follows from VM-SEQ-1, VM-SET-1 and VM-MAP-1 together. The record-of-keys reading of a set avoids it but collides more. |
