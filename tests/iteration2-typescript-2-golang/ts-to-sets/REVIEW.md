# Self-review of the TS → `.sets` standard (pass 1 → pass 2)

Pass 1 is 0.1.0 and pass 2 is 0.2.0. Every finding was checked against the value model only. There is no third pass.

| # | Finding in 0.1.0 | Fix in 0.2.0 |
| - | ---------------- | ------------ |
| R1 | `TsNumber`, `TsDate`, `TsValue`, `TsObject` and `TsNonNullish` were used but not defined anywhere in the standard. A target adapter only sees `.sets`, so the meaning has to be in the prelude, and the standard must point to it. | §3 gives each prelude domain's definition in one line and says the prelude is normative. |
| R2 | TS-VAL-7 makes a `Date` a `number`, so a `Date` value isn't in `TsObject`, though `Date` is assignable to `object` in TS. `Page<Date>` would be empty. This follows from the VM-PRES-2 narrowing and can't be fixed without undoing it. | Recorded as a known limit in TS-PRIM-1. |
| R3 | `SetOf(T)` holds every subdomain of `T`, infinite ones included. JS sets are finite, and VM-SET-1 says the standard narrows to that, but the grammar can't state finiteness (VM SI-14). | Noted in TS-LIB-1 as a stated-but-unenforced narrowing. |
| R4 | String literal types didn't say they follow VM-LEX-1 escaping (`"a, b"` must be `"a\, b"`). | TS-LIT-1 says so. |
| R5 | Function types were only covered as class fields (TS-CLASS-1). An object type with a function-typed property had no rule. | `TsFunction` added to the TS-PRIM-1 table for any function type. |
| R6 | A leading or middle rest element (`[...A[], B]`) had no rule. | Out of scope (TS-TUP-1). |
| R7 | TS-WIRE-1 missed that a class instance and a plain object with the same fields are the same value. | Added to the list. |
