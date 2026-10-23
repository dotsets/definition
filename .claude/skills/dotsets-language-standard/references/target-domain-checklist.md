# `.sets` domain-form checklist (target standards)

Walk every form. For each, record: exact / over + check / under (documented) / unsupported / gap. Forms come from `GUIDE.md`; re-check it for new ones.

## Declarations
- Top-level axiomatic domain (`String`)
- Alias (`Phi = phi`, `U = u`)
- Superdomain import (`Integer in axioms`)
- A file as a superdomain (a whole `.sets` file → package/module)

## Roster
- Roster of bare symbols (`{Succeeded, Failed}`)
- Roster of quoted symbols (`{"Number"}`), escaped separators (`\,`, `\\`), UTF-8 symbols
- Roster of numbers/mixed symbols
- Singleton (`{1}`, `singleton(x)`)
- Empty roster / `phi`
- Roster whose member is a domain symbol (`{Number}` = a domain of domains)

## Builder
- Scalar member (`{ x suchthat x in String }`)
- Tuple member, arity 1 (`{ (name) suchthat name in String }`)
- Tuple member, arity n
- Nested tuples (`((x), (y))`)
- Domain member (`{ {name, age} suchthat ... }`): a domain of domains / unordered member
- Rules combined with `and`
- Rules combined with `or` (optional components, the README `Human` shape)
- `not in` rules
- Membership in a singleton (equality)
- Rules that call predicates

## Operations (from `base`)
- `union` of disjoint domains
- `union` of overlapping domains
- n-ary unions (`union5`, nested `union`)
- `intersect`
- `difference`
- `product2` / cartesian products
- `u` (universal domain)
- `phi` (empty domain)

## Functions
- Transformer with domain parameters returning a domain (`Box({T sub u}): {...}`): generics
- Bounded domain parameter (`{T sub User}`)
- Value parameter captured via builder (`{x suchthat x in Number}`)
- Transformer returning a vector `( )` or scalar
- Predicate
- Call site applying a transformer to produce a named domain

## Not yet expressible (expect gap)
- Variable-length sequences / arrays
- Maps / dictionaries
- Field names distinct from positional variables
- Recursive / mutually recursive domains
- Cardinality constraints
