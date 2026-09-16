# Sets Intermediate Representation Definition

This document elaborates `.sets` representation definition.

> This is NOT llm generated.

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

Seperate multiple parameters with comma.

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

**2. Predicates**

Predicates are functions that don't return a value, just answers a question.
Think of predicates as functions that return `true` or `false` like much of the other programming languages.

From the transformer example, every function call goes through one predicate naturally, we can name it as such

```sets
Belongs({x suchthat x in u}, y sub u): x in y
```

This predicate bridges the gap between the callsite and the declaration. `x` being the concrete argument, `y` being the declared domain.

