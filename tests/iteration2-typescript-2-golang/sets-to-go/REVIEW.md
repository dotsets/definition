# Self-review of the `.sets` → Go standard (pass 1 → pass 2)

Pass 1 is 0.1.0 and pass 2 is 0.2.0. Every finding was checked against the value model only, with no knowledge of any source language. There is no third pass.

| # | Finding in 0.1.0 | Fix in 0.2.0 |
| - | ---------------- | ------------ |
| R1 | GO-UNION-1 didn't say how the rest of the union is classified once `{Nothing}` is removed. | The rest goes through GO-UNION-0…7 again. |
| R2 | GO-UNION-1 left out pointer, slice and map types. A `nil` slice already means `()` (GO-CALL-1), so it can't also mean `Nothing`. | Every type except `any` and interfaces is wrapped in `*G`, slices, maps and pointers included. |
| R3 | GO-UNION-3 said "a float domain", but the branch is usually a named domain whose Go type is a named float type (a refinement of `Float64`). | It matches any branch whose Go type has a float underlying type, and keeps that type. |
| R4 | GO-DECL-1 made `X = F(A)` a defined type. `type P Page[User]` isn't identical to `Page[User]`, so values don't pass between them. | A call to an input function is an alias too. |
| R5 | A field or element binder with no `in` clause (`(k, v) in r` alone) had no type. | That's `any` (the binder ranges over `U`). |
| R6 | GO-VAL-1 didn't say how an `Integer` member is represented in a float or big type. | One line per numeric Go type. |
| R7 | GO-ROSTER-1 put the `Integer` and `Binary64Finite` checks next to each other without saying which wins for `{1}`. | The rules apply in order; the first match wins. |
