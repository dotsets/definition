# Core Runtime Concepts

- **Standard version:** 0.2.0
- **Targets `.sets` spec at commit:** `f5a56cc` (`GUIDE.md`)
- **Depends on:** [core source standard](source/STANDARD.md) (VM-PRES-3), [`base.sets`](base.sets) (`Tagged`)
- **Declares:** [`runtime.sets`](runtime.sets), the `runtime` superdomain
- **Status:** draft
- **Used by:** [core source standard](source/STANDARD.md), [core target standard](target/STANDARD.md), and every language standard

A runtime domain isn't data (VM-PRES-3). It's a tagged domain, `Tagged(tag, D)`, that only has to be unique and equal to itself. This document reserves the tags and gives each concept its meaning. Adapters bind a concept by that meaning, never by a source spelling.

Rule IDs are `RT-*`. **Invariant** rules can't be overridden; **Default** rules can be, through the language standard's Overrides table ([standards/README.md](../README.md)).

## Reserved tags

**RT-1 (Invariant). The tags in RT-2 are reserved.** The tag atoms listed in RT-2 (`RuntimeAsync`, `RuntimeMutex`, ...) mean only their runtime concept. No source standard emits one as data. If a language's data can hold an atom with the same spelling (a Ruby symbol `:RuntimeAsync`), the source standard reports a diagnostic instead of emitting it. Other atoms, such as `RuntimeError`, stay ordinary data.

Tags are written in canonical atom form (GUIDE, pending): `RuntimeAsync`, or `\RuntimeAsync` when a domain of that name is in scope. `\Runtime\Async` and `Runtime\Async` spell the same atom and may be used in prose; adapters emit the canonical form.

## Concepts

**RT-2 (Invariant). The registry.**

| Concept | Tag | Meaning |
| ------- | --- | ------- |
| `Async(T)` | `RuntimeAsync` | A computation whose result, of domain `T`, is available later. |
| `Function(A, R)` | `RuntimeFunction` | Behavior that, given a value of `A`, produces a value of `R`. A function of several parameters takes a tuple. |
| `Mutex` | `RuntimeMutex` | Mutual exclusion: at most one holder at a time. |
| `Channel(T)` | `RuntimeChannel` | A conduit that passes values of `T` between concurrent parts of a program. |
| `Handle` | `RuntimeHandle` | A reference to a resource outside the program's data, such as an open file or a socket. |

A concept without a parameter tags `Phi`. Adding a concept adds a row here and a declaration in `runtime.sets`.

## Binding

**RT-3 (Invariant). Language standards bind concepts, in both directions.** A source standard maps each runtime type of its language to a concept (TypeScript `Promise<T>` → `Async(T)`). A target standard represents each concept with its language's construct (`Async(T)` → a Rust `Future`, a Go channel). A target with no construct for a concept reports it as unsupported. A runtime type that matches no concept is the language standard's own wrapper (VM-PRES-2) and isn't emitted as a runtime domain.
