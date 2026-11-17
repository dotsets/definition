# `.sets` → Go Target Standard (iteration 2)

- **Standard version:** 0.2.0 (second and final pass; see [REVIEW.md](REVIEW.md))
- **Based on:** `standards/core/VALUE-MODEL.md` 0.2.0 only.
- **Input:** one or more `.sets` files, with `base` in scope. Nothing is known about where the files came from. Names are only names; a domain's meaning is read from its `.sets` text.
- **Output:** one Go package (Go 1.21 or later) with one declaration per named domain, plus the runtime file of §9.
- **Status:** test draft

## 0. Principles

**GO-0. Read the text, not the name.** Base names (§2) have fixed meanings. Every other domain is classified by its definition (§3). Two domains with the same definition map to the same Go type, whatever they're called.

**GO-EXACT.** Each Go type is one of:
- **exact**: Go values and domain members correspond one to one;
- **wider**: some Go values aren't members (Go doesn't check membership in this iteration);
- **narrower**: some members have no Go value.

Wider is allowed. Narrower has to be stated in the rule. Projection follows VM-WIRE-1: where two members map to one Go value, the rule says so.

## 1. Declarations

**GO-DECL-1.** Every named domain in the input, except `base`'s, becomes a Go declaration with the same name:
- `X = Y`, where `Y` is a name → `type X = Y` (alias).
- `X sub Y` → `type X G(Y)`, with the prose comment above it as the doc comment (wider unless the prose is checked).
- A bare top-level axiom `X` → `type X struct{ Handle any }`, an opaque handle. Its members have no notation, so there are no constants.
- `F({T sub D}, …): body` → a generic type `F[T any, …]`, where body is mapped with `T` as a type parameter. `{T sub D}` becomes `any`, because Go constraints can't state domain membership (wider).
- Any other `X = expr` → `type X G(expr)`. It's an alias instead when `G(expr)` is a pointer type, `any`, or an instantiation of an input function (`X = F(A)`), so values pass between them unchanged.

`G(e)` is the Go type of expression `e`, by §2–§7.

**GO-NAME-1.** Field and constant names are PascalCase: the first letter is upper-cased, and characters not allowed in Go identifiers become `_`. If a name collides, it gets a numeric suffix.

## 2. Base domains

**GO-BASE-1.**

| `.sets` | Go | Exactness |
| ------- | -- | --------- |
| `String` | `string` | wider (Go strings may hold invalid UTF-8) |
| `Char` | `rune` | wider |
| `Integer`, `Natural` | `*big.Int` | wider (`nil`, negatives for `Natural`) |
| `Int8`…`Int64`, `Uint8`…`Uint64` | `int8`…`int64`, `uint8`…`uint64` | exact |
| `Fractional` | `*big.Rat` | wider |
| `Real` | `*big.Rat` | merges `1` and `1.0` (diagnostic, VM-WIRE-1) |
| `Binary64Finite`, `Float64` | `float64` | `Binary64Finite` is wider (NaN, ±Inf); both merge `2` and `2.0` |
| `Binary32Finite`, `Float32` | `float32` | as above |
| `Boolean` | `bool` | exact |
| `U` | `any` | wider |
| `Phi` | `struct{}` | wider (the domain has no members) |
| `Record` | `map[any]any` | wider |

The `IEEESpecial` atoms are the Go float values `math.NaN()`, `math.Inf(1)`, `math.Inf(-1)` and `math.Copysign(0, -1)` (VM-AX-3).

## 3. Rosters

**GO-ROSTER-1.** A roster is classified by its members. The rules apply in order, and the first match wins:
- Every member in `Boolean` → `bool`.
- All `String`s → `string`.
- All atoms → `string`, where each atom's value is its decoded spelling (VM-LEX-6).
- All `Integer`s → the smallest of `int8`…`int64` that holds every member, else `*big.Int`.
- All numbers, and every member in `Binary64Finite` (VM-SUB-0) → `float64`; otherwise `*big.Rat`.
- Anything else → `any`.

A named roster also gets one constant per member. With one member it's `XValue`. Otherwise it's `X` + PascalCase(the spelling), where for numbers `-` → `Neg` and `.` → `_` (`Level0`, `Size2_5`). Each constant's literal is the exact member value. Go constants are exact, so a decimal literal rounds correctly to the float type.

**GO-ROSTER-2.** `{Nothing}` on its own → `struct{}`. Inside a union, see GO-UNION-1.

## 4. Builders

A builder `{ p suchthat rule }` is classified by its member pattern `p` and its rule, in this order.

**GO-REC-1. Record builders → structs.** The member is a binder `r`, the rule contains `r in R` where `R` resolves to `Record` or `Map(K, U)`, and every field pattern `(key, v) in r` has a literal `String` key. The struct has, in clause order:
- `(k, v) in r and v in D` → field `K G(D)`, tag `` `sets:"k"` ``. With no `v in D` clause, the field is `any`. The same holds for element binders in §4's patterns.
- `(k, lit) in r` with no binder → not stored. It's a constant of the domain, written as a comment, or the variant tag in GO-UNION-1.
- `r in OptionalField(k, D)` → field `K Opt[G(D)]`. `OptionalFields({k1, …}, D)` gives one field per key.
- `r in Map(K2, V)` with `V` not `U` → field `Rest map[G(K2)]G(V)`, tag `` `sets:",rest"` ``. It holds every key that isn't a named field.
- Without `Rest`, there's a field `Extra map[G(K)]any`, tag `` `sets:",extra"` ``, for keys the rule doesn't name. Records are open (VM-REC-3), so this keeps the projection one to one (VM-WIRE-1).
- `r in OnlyKeys(…)` drops `Extra`.

Exactness: wider. Go doesn't stop `Extra` from holding a named key.

**GO-SEQ-1. Sequence patterns.** The member is a nested pair that ends in `()`, such as `(a, (b, ()))`. That's a fixed-length sequence (VM-SEQ-1):
- every element's Go type the same `G` → `[N]G`;
- otherwise → `struct { F0 G0; F1 G1; … }`.

**GO-SEQ-2. Prefix and rest.** The pattern is a nested pair whose last component is a binder `rest` with `rest in Seq(T)` → `struct { F0 G0; …; Rest []G(T) }`.

**GO-TUP-1. Tuple patterns.** Any other tuple pattern `(a, b, …)` → `struct { F0 G0; F1 G1; … }`. A 1-tuple is `struct { F0 G0 }`.

**GO-REF-1. Refinements.** The member is a binder `x`, and the rule starts with `x in D` → `G(D)`, wider. If the rest of the rule bounds an `Integer` `D` with literal `>=`/`<=` comparisons, the smallest fitting `intN`/`uintN` is used instead.

**GO-BUILD-X.** Any other builder → `any`, with a comment saying it couldn't be classified.

## 5. Structural calls

**GO-CALL-1.**

| `.sets` | Go |
| ------- | -- |
| `Seq(T)` | `[]G(T)`. `nil` and empty are both `()`. |
| `SetOf(T)` | `map[G(T)]struct{}` if `G(T)` is value-comparable, else `[]G(T)`. |
| `Map(K, V)` | `map[G(K)]G(V)` if `G(K)` is value-comparable, else `[]Pair[G(K), G(V)]` |
| `F(A, …)` for a function in the input | `F[G(A), …]` |

Value-comparable types are strings, booleans, fixed-size integers and floats, named types over these, and arrays and structs of them. Pointers (`*big.Int`), interfaces, slices and maps aren't, because Go compares pointers by identity. With float keys, NaN can't be looked up.

## 6. Unions

**GO-UNION-0. Normalize.** Flatten nested `union`s and drop duplicate branches. Then drop a branch that's visibly inside another:
- the same name;
- a sequence pattern whose elements are all in `T`, next to `Seq(T)`;
- a roster whose members are all in another branch's base domain.

If one branch is left, the union is that branch.

**GO-UNION-1. Nullable.** If a branch is `{Nothing}`, remove it and classify the rest again with GO-UNION-0…7, giving `G`. If `G` is `any` or an interface, `Nothing` is `nil`. Otherwise, slices, maps and pointers included, the type is `*G`, with `nil` for `Nothing`. A `nil` slice already means `()`.

**GO-UNION-2. Rosters.** All branches rosters → one merged roster (GO-ROSTER-1).

**GO-UNION-3. Floats.** One branch whose Go type has a float underlying type, with the other branches all rosters of `IEEESpecial` atoms → that branch's Go type.

**GO-UNION-4. Prefix sequences.** All branches are sequence patterns, and each is the previous one with more elements → a struct with fields `F0…Fn`, where the fields after the shortest pattern are `Opt[…]`.

**GO-UNION-5. Tagged records.** All branches are GO-REC-1 records, and there's one key `k` where each branch has a literal `(k, litᵢ)` with distinct literals → a sealed interface:

```go
type X interface{ isX() }
type XLit1 struct { … }   // one struct per branch, named X + PascalCase(litᵢ)
func (XLit1) isX() {}
```

The tag isn't stored: the variant type carries it. For a generic `X[T, …]`, every variant takes all of `X`'s type parameters.

**GO-UNION-6. Two branches** → `Either[G(A), G(B)]`. A value in both is held as `A` (diagnostic, VM-WIRE-1).

**GO-UNION-7. Otherwise** → `any`, holding the Go type of the first branch that contains the value. The comment lists the branches. Wider.

## 7. Intersections

**GO-INTER-1.** `intersect(A, B)` of two GO-REC-1 records → one struct with the fields of both, `A`'s first. A key in both takes `G(intersect(Da, Db))`, or the shared type when the domains are the same. There's one `Extra` or `Rest`. Anything else → `G(A)`, wider.

## 8. Recursion

**GO-RECUR-1.** Recursive domains map as written. Recursion through a slice, map, pointer or interface needs nothing more. A struct that contains itself directly gets a `*X` field. Values are finitely built (VM-RECUR-1), so the pointers are never cyclic.

## 9. Runtime

**GO-RT-1.** The package contains `sets_runtime.go`:

```go
type Opt[T any] struct { Present bool; Value T }   // a missing field vs. a present one
type Either[A, B any] struct { IsB bool; A A; B B }
type Pair[K, V any] struct { Key K; Value V }
```

## 10. Value projection

**GO-VAL-1.** A member maps to a Go value of its domain's type:
- A string → Go string (UTF-8).
- An `Integer`: in an `intN`/`uintN` type, that integer; in a float type, its float value (exact inside `Binary64Finite`); in `*big.Int`, a `big.Int` holding it.
- A `Fractional`: in a float type, its exact binary value; in `*big.Rat`, the exact rational.
- `true`/`false` → `bool`. An atom → its constant.
- `Nothing` → `nil`.
- A record → its struct: named fields set, other keys in `Extra`/`Rest`, an optional field missing → `Opt{}`.
- A sequence → a slice or array, a set → a map with `struct{}` values, a map → a Go map or pairs.

The reverse direction rejects Go values that aren't members, where the type is wider. Validation code isn't generated in this iteration, so the check is the reader's job.
