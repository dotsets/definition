# Sets Intermediate Representation Definition

This document elaborates `.sets` representation definition.

> This is NOT llm generated.

<!-- mtoc-start -->

* [Purpose](#purpose)
  * [Proposed solution](#proposed-solution)
* [Guide](#guide)
* [Source and Target Independence](#source-and-target-independence)

<!-- mtoc-end -->

## Purpose

`.sets`  is a domain representation of the type information collected from any source language.

I wanted to convert `typescript` types to `Golang` for a project of mine. But soon enough, crossed a fundamental issue of language barrier. Deterministic transfomation was becoming more complicated.

Eventually I crossed a specific problem. "Intersections".

In `typescript` we generally use intersections as an alternative to `or`. Not intersections themselves.

For example the following is one specific use case

```ts
type Foo = string | number;
```

This is easier to understand, `const x: Foo = 1` or `const x: Foo = ''`. 

But consider the following type

```ts
type Person = { name: string };
type Aged = { age: number };
type Human = Person | Aged;

function isPerson(p: unknown): p is Person {
    return 'name' in p;
}

function isAged(p: unknown): p is Aged {
    return 'age' in p;
}

function foo(x: Human) {
    if (isPerson(x)) {
        (void) x.name;
        //        ^?
    } else if (isAged(x)) {
        (void) x.age;
        //        ^?
    }
}
```

The bug in the above code is assuming `Person | Aged` is meant to be an "OR". While `{name: "abc", age: 123}` is a perfectly valid `Human` value. It happens to satisfy both `Person` and `Age` structure.

Intersections in typescript happen to be mathematical intersections, not logical ORs. But this, specifcally, changes how we would generate the respective Go code itself.

For a logical or, a constraint makes sense

```go
type stringOrIntConstraint interface {
    ~string | ~int
}
```

But when transpiling typescript's value-first type schema to Go, a constraint-like system applied over structural schema breaks how data gets deseralized into the runtime. We could start off with a broken type system; one that fails to work correctly over wire; or propose limitations that require significantly more complexity.

As an example, instead of making the Go code look like typescript's `Human` example, if we wanted to preseve the domain itself, one choice would be

```go
type Person struct {
    Name string `json:"name,omitempty"`
}

type Aged struct {
    Age int `json:"age,omitempty"`
}

type Human struct {
    *Person
    *Aged
}

func (h *Human) IsPerson() bool {
    return h.Person != nil
}

func (h *Human) IsAged() bool {
    return h.Age != nil
}

func Foo(x Human) {
    if x.IsPerson() {
        // handle x.Person.Name
    } else if x.IsAged() {
        // handle x.Aged.Age
    }
}
```

This version uses pointer embeddings. Visually significantly different, domain wise this will be simpler and can actually contain the wire data.

The bug is intentional, `x.IsPerson() == true` does not mean `x.IsAged() == false`. Kept it to keep the example the same as typescript version.

That brings me to the proposed solution.

### Proposed solution

`.sets` is a mathematical representation of Type information (schema) from any language. Each type is a domain, how we derive that domain is mathematical information, therefore is encoded as such in `.sets` files.

**Example**  
The previous `Human` example could be written as

```sets
// example code, not meant to be correct. The guide is what actually establishes the rules
// disjoint domains, don't require defining
String in base
RealNumber in base
union in base
intersect in base

Person = {
    (name) suchthat
    name in String
}

Aged = {
    (age) suchthat
    age in RealNumber
}

Human = {(x, y) suchthat {(x), (y)} subsetof intersect({(x)}, {(y)})}
```

From this, any other language alternative can be generated. Source language is no longer relevant as long as the domain description is correct.

But the very first question you will ask is that why is this not `Human = intersect(Person, Aged)`?

Let's get into that discussion.

Here `("abc")` is under valid domain `Person`. While `(123)` is under `Aged`. 

Wheras `("abc", 123)` is different. let's untangle the math here.

`union(Person, Aged)` expands to `{x suchthat x in Person or x in Aged}`, the kind of `x` here and in `("abc", 123)`, are in fight. Our input is a tuple or a vector, but the union works on scalars.

`Person` using vector elements makes sense. `{name: string}` is not the same as a `string` scalar. However, when we are considering an intersected domain, language specifics come into play.

We are trying to remove language specific flattening because they have no place in `.sets`. When we receive a `("abc", 123)`, value comes first, we check if `("abc")` is in `Person`? or `(123)` is in `Aged`? or check the cartesian product itself.

And the transpiler of the specific language-to-sets is responsible for constructing the right domain off of its language knowledge, instead of constructing the visually similar representation.

Another way to represent this would be

```sets
Human = {(x, y) suchthat 
    (x) in Person or
    (y) in Aged or 
    ((x), (y)) in product2(Person, Aged)} // nested tuples because that's what Person and Aged are

// while product being
product2({A sub u}, {B sub u}): {(x, y) suchthat x in A and y in B}
```

Or,

```sets
Human = {(x, y) suchthat 
    (x) in Person or 
    (y) in Aged or 
    (x, y) in product2(String, RealNumber)} // flattening Person and RealNumber domains

// while product being
product2({A sub u}, {B sub u}): {(x, y) suchthat x in A and y in B}
```

This is where `.sets` becomes meaningful.

## Guide

Read the fill specification and language guide in [Guide](./GUIDE.md).

## Source and Target Independence

A `.sets` representation describes a domain independently of the language from which the domain originated.

A source-language adapter is responsible for interpreting the source language's type system and producing the corresponding `.sets` domain.

A target-language adapter is responsible for interpreting the `.sets` domain and producing a representation in the target language.

The target adapter does not need to know the source language.

```text
Source Language
      |
      v
Source → .sets
      |
      v
   Domain
      |
      v
.sets → Target
      |
      v
Target Language
```

The `.sets` representation is therefore a semantic boundary between the two languages.

Source-language-specific representation decisions do not belong in `.sets` unless they are themselves part of the mathematical domain being represented.

Similarly, target-language-specific representation decisions are not implied by a `.sets` domain.

The purpose of `.sets` is to preserve the domain semantics while allowing the representation of that domain to differ between languages.

