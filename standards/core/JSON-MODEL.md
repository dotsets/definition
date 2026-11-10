# Core JSON Model

- **Status:** draft, moved out of [VALUE-MODEL.md](VALUE-MODEL.md) §8 unchanged; not yet reviewed
- **Used by:** [core source standard](source/STANDARD.md), [core target standard](target/STANDARD.md)

The canonical JSON projection of core values (VALUE-MODEL.md, VM-WIRE-1). Rule IDs are kept from the value model.

## Canonical JSON projection

JSON is the most common wire. This projection is the default wire for both directions. A language standard that uses a different wire format defines its own projection and keeps to the same rules: it's one-to-one on each domain, and it rejects non-members.

**VM-JSON-1 (Default).**

| Value | JSON |
| ----- | ---- |
| String | string |
| Number | number, written per VM-LEX-2 |
| `true`, `false`, `null` | `true`, `false`, `null` |
| Other atom `A` | the string `"A"` |
| `NotANumber`, `PositiveInfinity`, `NegativeInfinity` | not representable; see VM-JSON-2 |
| `NegativeZero` | `-0` |
| Record | object with one member per field, including fields the domain doesn't mention |
| `(Tag, payload)` variant | `{"Tag": payload}` |
| Sequence | array |
| Set | array, sorted by the canonical JSON text of its elements |
| Map with string keys | object |
| Map with other keys | array of `[key, value]` arrays, sorted by key text |
| Empty tuple `()` | `[]` |
| Other n-tuple | array of length n |

**VM-JSON-2 (Invariant). Ambiguity is a diagnostic, not a guess.** If a domain admits two different values with the same projection (for example the atom `Circle` and the string `"Circle"`, a 2-tuple and a 2-element sequence, or a variant `(Tag, p)` and a record with the single key `"Tag"`), the projection isn't one-to-one for that domain. The adapter then emits a diagnostic and either rejects the domain or uses a projection override from its language standard. The same applies to non-finite floats on plain JSON.
