# `.sets` Language Standards

Language standards are the contracts between programming languages and `.sets`.

```text
Source Language → [source standard] → .sets domain → [target standard] → Target Language
```

- A **source standard** (`<lang> → .sets`) says, for every type construct of a language, which values the construct admits and which canonical `.sets` text describes that domain.
- A **target standard** (`.sets → <lang>`) says, for every `.sets` domain form, which representation in the language holds exactly the members of that domain, and how membership is checked where the type system can't do it.

The two sides never refer to each other. A target standard doesn't assume any particular source language.

## Layout

```text
standards/
├── README.md                  this file
├── core/                      the core standard; every language standard extends it
│   ├── VALUE-MODEL.md         shared by both directions: values, lexical conventions, JSON projection
│   ├── base.sets              the `base` superdomain: operations, axioms and structural transformers
│   ├── source/STANDARD.md     core source standard  (any language → .sets)
│   ├── target/STANDARD.md     core target standard  (.sets → any language)
│   └── examples/*.sets        canonical domains with member / non-member test vectors
└── <language>/
    ├── source/STANDARD.md     extends core/source
    └── target/STANDARD.md     extends core/target
```

## How a language standard extends core

The core standard is the starting point for every language. A language standard only writes down what is specific to it.

1. **Declare the base.** The header names the core version it extends, for example `Extends: core/source 0.2.0`.
2. **Every core rule applies unless overridden.** Rules have stable IDs (`CS-REC-1`, `CT-UNI-2`, `VM-LEX-1`). A language standard doesn't restate rules it accepts.
3. **Rules have one of two classes:**
   - **Invariant**: can't be overridden. These rules keep `.sets` a shared boundary (disjoint axioms, no flattening, no unknown domains, deterministic output, no silent under-approximation, and so on). If a language seems to need an exception, that's a Spec Issue for core, not an override.
   - **Default**: can be overridden. The language standard lists each override in an **Overrides** table with the rule ID, the replacement rule, and the reason.
4. **Language axioms extend `base`; they don't replace it.** A source standard reuses `base` domains wherever their prose meaning fits (`String in base`, `Int32 in base`). It declares its own axioms or described subdomains in `<language>.sets` only for values that `base` doesn't cover, and it gives each new domain prose semantics so target standards can bind it by meaning.
5. **Rule IDs in language standards** use the language prefix (`TS-REC-3`, `GO-UNI-1`) so they never collide with core IDs.

### Overrides table (template)

| Core rule | Replacement | Reason |
| --------- | ----------- | ------ |
| `CS-NUM-3` | Numeric literals are emitted as the shortest round-trip decimal of the binary64 value. | Literal types in this language are binary64 values; exact expansions are unreadable. |

## Status

All standards are drafts against the v0 spec. Each standard records the `GUIDE.md` commit it targets. When the spec changes, re-check every rule that cites the changed section.
