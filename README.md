# Sets Intermediate Representation Definition

This document elaborates `.sets` representation definition.

> This is NOT llm generated.

<!-- mtoc-start -->

* [Purpose](#purpose)
  * [Proposed solution](#proposed-solution)
* [Specification](#specification)
  * [1. Grammar](#1-grammar)
  * [2. Domains](#2-domains)
    * [Reserved Keywords](#reserved-keywords)
    * [Using domain builder notation](#using-domain-builder-notation)
    * [Special sets](#special-sets)
  * [3. Functions ](#3-functions-)
    * [Understanding parameters](#understanding-parameters)
      * [Calling a function](#calling-a-function)
  * [4. Expressions](#4-expressions)
    * [Membership](#membership)
    * [Subdomain](#subdomain)
    * [Logical operations](#logical-operations)
  * [5. Domain Operations](#5-domain-operations)
    * [Union](#union)
    * [Difference](#difference)
    * [Intersection](#intersection)
    * [Product](#product)
  * [6. Tuples and Vectors](#6-tuples-and-vectors)
  * [7. Domain Builder Notation](#7-domain-builder-notation)
  * [8. Singleton Domains](#8-singleton-domains)
  * [9. Special Domains](#9-special-domains)
    * [Empty domain](#empty-domain)
    * [Universal domain](#universal-domain)
  * [10. Domain Axioms](#10-domain-axioms)
  * [11. Domain Validation](#11-domain-validation)
  * [12. Source and Target Independence](#12-source-and-target-independence)

<!-- mtoc-end -->

## Purpose

`.sets`  is a domain representation of the type information collected from anyetsource language.

I wanted to convert `typescript` types to `Golang` for a project of mine. But soon enough, crossed a fundamental issue of language barriers. Deterministic transfomation was becoming more complicated.

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
// disjoint sets, don't require defining
String
RealNumber

// defining here as example but are part of standard operation also defined in sets lang.
union(A, B, ...): union({x suchthat x in A or x in B}, ...)
difference(A, B): {x suchthat x in A and x not in B}
intersect(A, B, ...): intersect(
    difference(
        union(A, B),
        union(
            difference(A, B),
            difference(B, A)
        )
    ),
    ...
)

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

`union(Person, Aged)` expands to `{x suchthat x in Person or x in Aged}`, the kind of `x` here and in `("abc", 123)` in fight. Our input is a tuple or a vector, but the union works on scalars.

`Person` using vector elements makes sense. `{name: string}` is not the same as a `string` scalar. However, when we are considering an intersected domain, language specifics come into play.

We are trying to remove language specific flattening because they have no place in `.sets`. When we receive a `("abc", 123)`, value comes first, we check if `("abc")` is in `Person`? or `(123)` is in `Aged`? or check the cartesian product itself.

And the transpiler of the specific language-to-sets is responsible for constructing the right domain off of its language knowledge, instead of constructing the visually similar representation.

Another way to represt this would be

```sets
Human = {(x, y) suchthat (x) in Person or (y) in Aged or ((x), (y)) in product2(Person, Aged)} // nested tuples because that's what Person and Aged are

// while product being
product2(A, B): {(x, y) suchthat x in A and y in B}
```


Or,

```sets
Human = {(x, y) suchthat (x) in Person or (y) in Aged or (x, y) in product2(String, RealNumber)} // flattening Person and RealNumber domains

// while product being
product2(A, B): {(x, y) suchthat x in A and y in B}
```

This is where `.sets` becomes meaningful.

## Specification

### 1. Grammar

`.sets` uses very similar grammar to LaTeX. That is intentional.

**1. Every statement is a declaration**

There is nothing in `.sets` that is not a declaration. Every statement itself is a declaration.

**2. A declaration is either a domain or a function**

**3. A top level declaration without an assignment is a unique domain**

Example
```sets
String
Integer
Real
Character
```
The above each represent a domain each.


**4. Every top level unique domain are disjoint**

From \#3, `intersect(String, Integer) = phi`.

However, the following is not inherently disjoint in relative to the other sets.

```sets
Number = union(Integer, Real)
```

**5. Building sets**

The fundamental object in `.sets` is a domain. An anonymous domain can be constructed with `{}` (curly braces).

`.sets` support both rooster notation and builder notation. 

Rooster notation: `{1, 2, 3, 4}`.
Builder notation is explained in a different section.

### 2. Domains

Declaring domains is one of the central purposes of `.sets`. We start with reserved contextually relevant keywords:
> Contextually relevant because you don't need to shy away from using the keywords, can be used just fine, as their meaning is context dependent.

#### Reserved Keywords

- `not`: negation of any conditional.
- `in`: lvalue to be a member of rvalue. rvalue must be a domain.
- `suchthat`: used to elaborate on a named element's properties. used in domain builder notation.
- `sub`: used to check if lhs is a sub domain of rhs.
- `and`: logical and operation.
- `or`: logical or operation.
- `Phi` or `phi`: an empty domain.
- `u` or `U`: the universal domain.

#### Using domain builder notation

Set builder notation follows a couple of simple rules
1. Describe the form of the elements, scalar, vector or tuple.
2. Followed by `suchthat` to start describing the right properties of the element itself

Example:

```sets
EventNumbers = {
    x suchthat x in union(Integer, Real) and x%2 = 0
}
```

#### Special sets

An empty domain is effectively `{}`. There is a reserved keyword for this `Phi` or `phi`. Which is effectively just an alias to `{}`.
```sets
Phi = {}
phi = Phi
```

Universal domain is a domain of all known and unknown elements. Represented primarily as `{ ... }`, whitespace doesn't matter. Similarly `.sets` understands the alias `U` or `u`.
```sets
U = {...}
u = U
```

One concept worth pointing out now, there is a reason `{...}` is chosen over `{x}` for a universal domain. `{x}` is different from `{...}`. 
Generally `{x suchthat` the `suchthat` that follows spreads the elements' domain. If there is no description of the domain, `{x}` becomes a domain with just one element, therefore can not contain any element from the domain of discourse. `...` is the variadic symbol in `.sets` used in different situations.

### 3. Functions 

Functions are transformations or predicates over explicitly defined domains.

Write an identifier, or called name of the function normally. Followed by parenthesis to accept the input. Functions don't follow a ` = ` syntax, rather it uses `:` to separate the lhs from the rhs.

Example
```sets
f(): {}
```

Here `f` is a function that always evaluates to `phi`.

##### Understanding parameters

The language has an inherently quirky parameter passing grammar.

A function can work over a scalar, vector (essentially single element) or a domain.

In general, all parameters to a function is considered to be a domain. 

```sets
f(x): x * 2
```

In the above example, `.sets` will look for a domain aliased `x`, if can't find one, the code is invalid.

This is because every function works over a domain. The passed value at callsite will be validated through a simple `in` predicate.

So when `f` is called like `f(2)`, first check that happens is `2 in x`, if `x` is not a valid domain, the function becomes a transformation over an unknown domain, which is invalid.

The right way to write the function would be
```sets
f({x suchthat x in Number}): x*2
```

Multiple parameters are separated by commas.

```sets
f(
    {x suchthat x in Number},
    {y suchthat y in Number}
): x + y
```

If you want to allow any element, use the universal set.

```sets
f({x suchthat x in u}): x
```

##### Calling a function

Calling a function is as simple as the same, followed by parenthesis. Inside the parenthesis, the arguments.


There can be two types of functions.

**1. Transformers**

These kinds of functions produce either a scalar, a vector, or another domain, from `n` number of inputs where `n >=0`.

The key is that transformers produce some output. The rhs becomes what a transformer returns. So wrapping in `()` makes it return a vector, wrapping in `{}` returns a domain, otherwise a scalar.

```sets
Double({x suchthat x in Number}): x * 2
```

A scalar expression on the right hand side produces a scalar.

Parentheses can be used when the result needs to be a vector or tuple.

```sets
Pair({x suchthat x in Number}): (x, x)
```

Curly braces produce a domain.

```sets
EvenNumbers({x suchthat x in Number}): {
    y suchthat y in x and y % 2 = 0
}
```

The result of a transformer is therefore determined by the expression used on the right hand side.

**2. Predicates**

Predicates are functions that don't return a value, just answers a question.
Think of predicates as functions that return `true` or `false` like much of the other programming languages.

From the transformer example, every function call goes through one predicate naturally, we can name it as such

```sets
Belongs({x suchthat x in u}, y sub u): x in y
```

```sets
IsEven({x suchthat x in Number}): x % 2 = 0
```

The distinction between a transformer and predicate is primarily semantic. A transformer produces a value or domain, while a predicate establishes whether a condition holds.
This predicate bridges the gap between the callsite and the declaration. `x` being the concrete argument, `y` being the declared domain.

### 4. Expressions

Expressions are used to describe domains, function parameters, function results, and predicates.

An expression can represent a value, a domain, or a condition depending on its context.

#### Membership

The `in` keyword checks whether the left hand side belongs to the domain on the right hand side.

```sets
x in Number
```

The right hand side must resolve to a domain.

Membership is a fundamental operation in `.sets`.

For a concrete value:

```sets
2 in Integer
```

is a membership predicate.

Membership can also be used against a constructed domain:

```sets
2 in union(Integer, Real)
```

#### Subdomain

The `sub` keyword checks whether the left hand side is a subdomain of the right hand side.

```sets
Integer sub Number
```

A domain can be a subdomain of another domain without being equal to it.

For example:

```sets
Integer
Real

Number = union(Integer, Real)
```

therefore establishes that:

```sets
Integer sub Number
Real sub Number
```

while `Number` and `Integer` are not disjoint.

`sub` operates on domains, while `in` operates on elements and domains.

#### Logical operations

`.sets` provides logical operations for predicates.

`and` requires both conditions to hold.

```sets
x in Number and x % 2 = 0
```

`or` requires at least one condition to hold.

```sets
x in Integer or x in Real
```

`not` negates a condition.

```sets
not x in Integer
```

Logical operations operate on conditions and do not themselves construct domains.

### 5. Domain Operations

Domains can be combined using the standard domain operations.

These operations are functions and therefore use function-call syntax.

#### Union

`union` constructs a domain containing elements belonging to any of its input domains.

```sets
union(A, B)
```

is equivalent to:

```sets
{x suchthat x in A or x in B}
```

`union` accepts a variable number of domains.

```sets
union(A, B, C, ...)
```

The domains supplied to `union` do not need to be disjoint.

For example:

```sets
A = {1, 2}
B = {2, 3}

union(A, B)
```

produces:

```sets
{1, 2, 3}
```

#### Difference

`difference` constructs a domain containing elements belonging to the first domain but not the second.

```sets
difference(A, B)
```

is equivalent to:

```sets
{x suchthat x in A and x not in B}
```

#### Intersection

`intersect` constructs a domain containing elements common to all supplied domains.

```sets
intersect(A, B)
```

is equivalent to:

```sets
difference(
    union(A, B),
    union(
        difference(A, B),
        difference(B, A)
    )
)
```

The operation can be extended to multiple domains.

```sets
intersect(A, B, C, ...)
```

The result is a domain containing only elements that belong to every supplied domain.

#### Product

`product` constructs a domain of tuples from its input domains.

For two domains:

```sets
product2(A, B)
```

can be represented as:

```sets
product2(A, B): {
    (x, y) suchthat x in A and y in B
}
```

For example:

```sets
Integer
String

Pairs = product2(Integer, String)
```

contains values such as:

```sets
(1, "hello")
(2, "world")
```

but does not contain:

```sets
1
"hello"
```

The tuple itself is the element of the resulting domain.

### 6. Tuples and Vectors

Parentheses are semantically significant in `.sets`.

They are not merely used for grouping expressions.

A scalar and a one-element tuple are different values:

```sets
"abc"
("abc")
```

Therefore:

```sets
"abc" != ("abc")
```

A domain can explicitly describe the shape of its elements.

```sets
Person = {
    (name) suchthat name in String
}
```

Here `Person` contains one-element tuples whose element belongs to `String`.

It does not contain the scalar `"abc"`.

```sets
("abc") in Person
```

is valid, while:

```sets
"abc" in Person
```

is not.

Likewise, a two-element tuple:

```sets
(x, y)
```

is a different value from either of its individual elements.

This distinction is important when representing structural types. `.sets` does not implicitly flatten tuples or domains.

Whether a target language represents a tuple as fields, an embedded structure, an array, or another representation is a concern of the target-language transformation and is not implied by the tuple itself.

### 7. Domain Builder Notation

Domain builder notation describes the possible elements of a domain.

The general form is:

```sets
{
    element suchthat condition
}
```

The expression before `suchthat` describes the form of the elements.

The condition after `suchthat` describes the properties those elements must satisfy.

For example:

```sets
EvenNumbers = {
    x suchthat x in Integer and x % 2 = 0
}
```

The elements of `EvenNumbers` are integers satisfying the given condition.

The element expression may describe a scalar:

```sets
{x suchthat x in Integer}
```

or a tuple:

```sets
{(x, y) suchthat x in Integer and y in String}
```

The latter describes a domain whose members are tuples.

The distinction between the element expression and its conditions is important. `suchthat` does not merely introduce a boolean filter; it establishes the domain over which the described element is considered.

### 8. Singleton Domains

A singleton domain contains exactly one element.

A singleton can be constructed using ordinary roster notation:

```sets
{1}
```

or expressed through a helper function:

```sets
singleton(x): {x}
```

This also allows equality to be expressed through membership.

```sets
x in singleton(y)
```

is true exactly when `x` is equal to `y`.

For example:

```sets
IsOne(x): x in singleton(1)
```

can then be used to construct a domain:

```sets
One = {
    x suchthat IsOne(x)
}
```

This demonstrates that equality does not need to be a separate domain construction primitive. It can be expressed through membership in a singleton domain.

### 9. Special Domains

`.sets` defines two special domains.

#### Empty domain

The empty domain contains no elements.

It is represented by:

```sets
{}
```

`Phi` and `phi` are aliases for the empty domain.

```sets
Phi = {}
phi = Phi
```

Therefore:

```sets
intersect(A, Phi) = Phi
```

for any domain `A`.

The empty domain is a subdomain of every domain.

```sets
Phi sub A
```

is true for every domain `A`.

#### Universal domain

The universal domain contains all known and unknown elements in the domain of discourse.

It is represented by:

```sets
{...}
```

`U` and `u` are aliases for the universal domain.

```sets
U = {...}
u = U
```

The distinction between `{...}` and `{x}` is intentional.

```sets
{x}
```

is a singleton domain containing `x`.

It does not mean an arbitrary domain containing `x`.

The `...` notation indicates that the domain is not restricted to the explicitly represented element.

Therefore:

```sets
{x} sub U
```

is true, while `U` itself cannot be reduced to a singleton.

`...` is also the variadic symbol used elsewhere in `.sets`. Its exact meaning depends on the construct in which it appears.

### 10. Domain Axioms

Some domains are introduced without an assignment.

```sets
String
Integer
Real
Character
```

Each such declaration introduces a unique domain.

Unique domains declared at the top level are disjoint.

Therefore:

```sets
intersect(String, Integer) = Phi
```

and:

```sets
intersect(Integer, Real) = Phi
```

These domains may subsequently be used to construct other domains.

For example:

```sets
Number = union(Integer, Real)
```

`Number` is a derived domain and is not subject to the disjointness rule that applies to unique domains.

Therefore:

```sets
intersect(Number, Integer) = Integer
```

while:

```sets
intersect(Integer, Real) = Phi
```

### 11. Domain Validation

`.sets` does not permit unknown domains.

Every domain referenced by a declaration must either:

* be introduced by a top-level domain declaration;
* be constructed by a domain expression;
* be provided by the standard domain environment; or
* explicitly resolve to the universal domain.

For example:

```sets
f(x): x * 2
```

is invalid if `x` does not resolve to a domain.

The language does not infer that `x` must be numeric because it is multiplied by `2`.

Instead, the domain must be explicitly specified:

```sets
f({x suchthat x in Number}): x * 2
```

If arbitrary values are intended, the universal domain must be explicitly used:

```sets
f({x suchthat x in u}): ...
```

This is intentional. `.sets` is domain-first and does not introduce dynamic or unknown domains to defer semantic decisions until a later stage.

### 12. Source and Target Independence

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

