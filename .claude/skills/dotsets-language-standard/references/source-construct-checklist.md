# Source type-system construct checklist

Walk every category. For each construct the language has, record mapped / out of scope / gap. Record N/A for categories the language lacks.

## Scalars and atoms
- Strings, characters, runes, bytes
- Integers (each width/signedness), big integers
- Floating point (each width; NaN, ±Inf, −0 handling)
- Decimals / fixed point
- Booleans (and whether they subtype integers)
- Null / nil / None / undefined / unit / void — which are distinct values
- Date/time, UUID, and other "well-known" scalar types the ecosystem treats as atoms

## Literals and enumerations
- String / number / boolean literal types
- Enums (plain, string-valued, numeric, with explicit values)
- Template-literal / pattern types (e.g. TS template literals, regex-constrained strings)

## Products
- Records / structs / objects / classes (data shape only)
- Positional tuples, named tuples
- Optional fields, default values, readonly/immutability markers
- Open vs closed shapes (excess properties, `additionalProperties`, `#[serde(deny_unknown_fields)]`)
- Embedding / inheritance / spread (`struct { Base }`, `extends`, `...`)
- Index signatures / dynamic keys

## Sums and combinations
- Untagged unions (`A | B`, `Union[A, B]`) — check overlap semantics
- Tagged / discriminated unions, Rust enums with payloads, sealed classes, `oneof`
- Intersections (`A & B`) and mixins
- Nullable / optional wrappers (`T?`, `Option<T>`, `*T`, `Optional[T]`)
- Top and bottom types (`any`, `unknown`, `interface{}`, `object`, `never`, `!`)

## Collections
- Arrays / lists / slices (variable length)
- Fixed-length arrays
- Sets
- Maps / dictionaries / records-of-K-to-V
- Strings-as-sequences (if the language treats them so)

## Abstraction
- Type aliases vs new/nominal types (`type A = B` vs `type A B`, newtypes, branded types)
- Generics / type parameters, constraints/bounds, defaults, variance
- Higher-kinded / conditional / mapped types (TS `Partial<T>`, `keyof`, conditional types)
- Recursive and mutually recursive types
- Interfaces / traits / protocols (structural or nominal method sets)

## Behavior (usually out of scope for wire domains — decide explicitly)
- Function / method types
- Channels, promises/futures, iterators/generators
- Pointers / references (identity vs value)

## Refinements and validation
- Range/length/pattern constraints (from the type system or common schema annotations)
- Non-empty, integer-valued floats, etc.

## Modules
- Packages / modules / namespaces → superdomain files
- Imports / re-exports → `X in file`
- Visibility (exported vs private) — usually not domain information
