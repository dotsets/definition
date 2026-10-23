<!-- vi: set wrap: -->
<!-- mtoc-start -->

* [How to write `.sets`](#how-to-write-sets)
  * [The dot sets fundamental](#the-dot-sets-fundamental)
  * [Ways of strongly describing domains](#ways-of-strongly-describing-domains)
    * [Roster method](#roster-method)
    * [Builder method](#builder-method)
  * [Functions](#functions)
    * [Transformers](#transformers)
    * [Predicates](#predicates)
  * [Domain operations built on axioms](#domain-operations-built-on-axioms)
    * [unions of domains](#unions-of-domains)
    * [difference of domains](#difference-of-domains)
    * [intersection of domains](#intersection-of-domains)
  * [phi](#phi)
  * [`.sets` files are super domains](#sets-files-are-super-domains)
  * [Singletons](#singletons)
  * [Domain Validation](#domain-validation)
  * [Tuples and Vectors](#tuples-and-vectors)

<!-- mtoc-end -->

# How to write `.sets`

`.sets` describes domains, domains of "possible" values.

## The dot sets fundamental

The first thing to understand about `.sets` is that this language has no values until described. Therefore, recursively, there is no non-domain values or identifiers.

This brings us to describing through axioms.

**1. Top domain axiom**

Before describing a domain, i.e. a collection of values, one needs to visualize it and name it.

The name can be anything, this top level domain does not need a descriptor.

As an example, to describe a collection of numbers, declare a top level domain as an axiom

```.sets
Number
```

Similarly, for integers

```.sets
Integer
```

**Every top level domain is disjoint to one another.**

You can declare two domains such as

```.sets
Number
Integer
```

However, conceptually it is incorrect, unless your concepts for these two domains are different from the math.

Here any member of `Number` **can not be a member of** `Integer`. Instead, it should be rather written as

```.sets
Integer
Decimal // or FloatingPoint
Number = union(Integer, Decimal)
```

I will talk about the `union` syntax later.

The idea of union is the same as in math, it declares a new domain, such that the new space occupies the total of the member spaces.

The union does not change the fact that `Integer` and `Decimal` are disjoint. They still are. However, describing `Number` as an union allows for `Number` to not be disjoint with respect to the available spaces `Integer` and `Decimal`.

**There is no concept of primitive in `.sets`**

In this language, you always describe your domain. `Number` here is one representation of it.

You don't need to describe a number, I would argue, describing it will put an unnecessary language scope understanding constraint on it.

> Why?
>
> Fair question.
>
> Because when we break a language barrier, symbol sets change. While yes, numbers are very much universal, they also put a constraint
> where this language understands numbers.
>
> Numbers themselves are nothing special, even in math, decimal system is essentially `{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}`
> binary as `{0, 1}`, I can declare a number system named `BladiBla` with symbol domain `{a, *, 8, j}` which would be totally fair and valid.
>
> The point is, numbers are a bunch of symbols, whom we gave a meaning through counting. And that can be preserved through simple
> agreement of an axiomatic domain called `Number`. Value-only symbols cause different problems. The zero to one collection does NOT collect full collection of every possible permutation of the symbols. Which becomes necessary. Otherwise it becomes very difficult to differentiate between a domain of `{0, 1}` containing `01` vs a domain of `{a, b}` NOT containing `ab`.
>
> This is v0 of the language. To support this, a new grammar willl be necessary and may come in the future.

## Ways of strongly describing domains

There are two ways to describe domains,

1. Roster method
2. Builder method

### Roster method

This is when you write out the exact members by hand.

Example

```.sets
OrderStatus = {Succeeded, Failed, Processing}
```

Inside curly braces, type out the value-members.

**Strings are not top level "kinds"**

The first intuition must have been why not `"Succeeded"` why `Succeeded`. No particular reason.

Again, languages chose to use quotation marks to denote a collection of characters as strings.

Think about `C`, the language.

```c
char * a = "a";
// is effectively
char aa[1] = {'a'};
// expanding
char ac = 'a';
// is the same as
int ai = 97;
```

`"a"` is not special, just another way of writing `97`, which is another way of writing `9` and `7` together.

Writing `Succeeded` as-is already uniquely identifies the symbol, a transpiler may treat it differently.

```.sets
OrderStatus = {"Succeeded", "Failed", "Processing"}
```

Is also valid. Transpiler to transpiler, using quotes may be a valid way to communicate strings.

One caveat, using quotes helps distinguish the member symbol from a global symbol. For example, `Number` may be top level axiomatic domain, while a member may require `Number` as a string. Writing `{Number}` will cause the member to be treated as the `Number` domain, instead, adding any character to the symbol changes it. So `{"Number"}` becomes a valid escape.

I am emphasizing this because `"Number"` is not special. The quotes here just distinguish `Number` as a domain from it as a value. It could very well be written as `"Number`, or `@Number`, if the transpiler ecosystem agreed on those prefixes to distinguish domain vs values. Using opening and closing quotations both, helps with human intuition and carries existing comprehension memory. That is all.

**What are symbols?**

Yes I waited until now to define this.

_A symbol is a collection of utf8 characters, put forth without an internal domain of reserved symbols._

The internal domain of reserved symbols identify `, { }` _for now_.

UTF8, because many languages like Python3.15 now supports it, it is important to support it.

To include the natural separators, use `\` escape. Like `{SomeValue\,,Somevalue2}` makes first element `SomeValue,`. Escape itself through `\\`.

In roster method, `SomeValue Somevalue2` are ONE symbol, comma is necessary to separate members. I do not want to put many responsibility to whitespaces in `.sets`.

> `{S1, S2}` means `S1` one value, ` S2` another. I don't want that to necessarily be the case, such that one always must write `{S1,S2}` as that ends up giving whitespaces a value. Therefore, whitespaces `{ ,\b,\r,\n}` one byte around a reserved symbol is IGNORED.

### Builder method

This one is where the power lies.

In builder method, you describe your domains through pure mathematical rules. Grammar is roughly like

```
Domain = { member suchthat rule }
```

Generally with the constraints of `.sets` so far, it is incredibly hard to show practical examples of builder method. So we look at a transpilation example.

The following typescript code

```ts
type User = {
  name: string;
  age: number;
};
```

The domain name in `.sets` will be the same, so

```.sets
User = {
```

Now comes the members. The members can simply be a tuple. Or even a domain.

```.sets
User = { {name, age} }
```

Or

```.sets
User = { (name, age) }
```

Because a "struct" may have a meaning for ordering depending on source language, a tuple might make more sense here, however for typescript, it is irrelevant, so I will keep to my promise and give the source language the right importance, use domain.

When you describe members, use the keyword `suchthat`.

```.sets
User = { {name, age} suchthat ... }
```

After `suchthat`, `name` and `age` become scoped symbols to work with to define the rules.

**Membership can be denoted by using keyword `in`.**

Let me start by just the name.

```.sets
String // top level

User = { {name} suchthat name in String }
```

Here, single `{name}` is different than `name suchthat`. `{name}` conveys a domain of domains, while `name suchthat` makes `User` the same as a domain involving `String` members.

**Rules can be compounded through `and` and `or`.**

Now back to using both `name` and `age`:

```.sets
String
Number
User = { {name, age} suchthat name in String and age in Number }
```

## Functions

`.sets` supports functions. There are two kinds of functions

1. **predicates** tell a truth, whether something is true or false.
2. **transformers** transforms input value to an output.

### Transformers

Syntax in general for functions are as follows

```.sets
FunctionName(Parameters): ReturnValue
// callsite
FunctionName(Arguments)
```

When declaring function parameters, the parameter **must be a domain**. When calling said function
arguments are checked against the declared domain parameter for membeship. That membership decides
whether the argument is allowed or not.

Transformers can return anything, a particular value/scalar, a domain, a vector etc.

The key is that transformers produce some output. The rhs becomes what a transformer returns. So wrapping in `()` makes it return a vector, wrapping in `{}` returns a domain, otherwise a scalar.

When a return value is typed, same rule applies, unless the symbol is predefined to be a domain, it acts as a value.

The result of a transformer is therefore determined by the expression used on the right hand side.

Examples

```.sets
ToUser(x, y): { {x, y} }
```

Here `ToUser` is the only semantic elaboration of the return domain being `User` as described before. Luckily there is a better way.

```.sets
ToUser(x, y): { {x, y} suchthat {x, y} in User }
```

A shortcut to the above is `{ {x, y} sub User }`. `sub` declares the lvalue is a subdomain of rvalue.

We take this `sub` to force passing a domain as an argument to a function.

**`sub` explained**

`sub` marks lhs the subdomain of rhs. This is different from using `in` itself.

`x sub U` does NOT mean `x in U`, nor `{x} in U`. However,
`{x} in U` *can* be written as `{x} sub U`. 

```.sets
TakeUser({D sub User})
```

`D sub User` marks `D` a subdomain of `User`, however by itself, at callsite `TakeUser` could take a value that is a member of `D`. By wrapping it in `{}` we force the callsite to pass a domain that is a subdomain of `User`, e.g. `{{"v1", "v2"}} sub {D sub User}`.

### Predicates

Predicates are functions that don't return a value, just answers a question.
Think of predicates as functions that return `true` or `false` like much of the other programming languages.

The distinction between a transformer and predicate is primarily semantic. A transformer produces a value or domain, while a predicate establishes whether a condition holds.
This predicate bridges the gap between the callsite and the declaration. `x` being the concrete argument, `y` being the declared domain.

## Domain operations built on axioms

From the existing rules, we can build some important operations on the `.sets` language.

E.g.

### unions of domains

When we take union of two domains, the result is a domain, such that, it describes the combined space each of the participant domain occupies.

For example, image domains

```.sets
N1 = {1, 2, 3}
N2 = {4, 5, 6}
```

Now a domain `NN` a union of `N1` and `N2`, all of the following are true predicates

```.sets
In({x suchthat x in Number}, {D sub U}): x in D
In(1, NN) // true
In(2, NN) // true
In(3, NN) // true
In(4, NN) // true
In(5, NN) // true
In(6, NN) // true
```

**Note on the `In` declaration**
There are two parameters, `{x suchthat x in Number}`, since there is no `Number` global domain that includes number specifically, this is a semantic choice, it is just a global domain.

Using the builder syntax, during parameterization, you can capture the argument passed. The variable[s], that is/are part of the builder notation, gets captured into the function body to use, like we used `x` later in the function body. Same goes for `D`.

`U` is the same as top level domains. It is a semantic choice to declare an universal domain.

**Back to unions**

Now, how do we define `NN` in `.sets`?

Through a function. A union of two domains as follows

```.sets
union({A sub U}, {B sub U}): { x suchthat x in A or x in B }
```

The resulting domain declaration checks both domains for membership. Any one of them satisfies the predicate.

for more, this can easily be combined recursively. a five domain union for example

```.sets
union5({A sub U}, {B sub U}, {C sub U}, {D sub U}, {E sub U}): {
                x suchthat
                x in A or
                x in union(B, union(C, union(D, E)))
            }
```

### difference of domains

```.sets
difference({A sub u}, {B sub u}): {x suchthat x in A and x not in B}
```

This is self explanatory. This helps us get to the more important definition.

### intersection of domains

When one element is a member of multiple domains simultaneously, the resulting domain is called the intersection of the domains in context.

Defined
```.sets
intersect({A sub U}, {B sub U}): { x suchthat x in A and x in B }
// or can also be defined through previous functions
intersect({A sub U}, {B sub U}): difference(
        union(A, B),
        union(
            difference(A, B),
            difference(B, A)
        )
    )
```

Now with n-ary function, an intersection of 3 domains
```.sets
intersect3({A sub U},{B sub U}, {C sub U}):
    intersect(A, intersect(B, C))
```

## phi

`phi` is an empty domain, not a special domain. It can be defined as
```.sets
phi = {}
Phi = phi // alias
```

## `.sets` files are super domains

Each `.set` file is a superdomain. The name of the file becomes the name of the domain.

Therefore one file named `axioms.sets`, containing
```sets
String
Integer
Real
Number = union(Integer, Real) // implementation defines union
```

Another file `any.sets` can perform
```sets
Integer in axioms
```

This top level `in` declares, not defines, declares that `Integer` is a member of `axioms`. 

```sets
Complex in axioms
```

The above, is invalid. As `axioms` have no member domain named `Complex`.

> This is similar to "importing" from other modules in programming languages.

Now an implementation can expose a superset, instead of magically defining a function or domain. For instance

```sets
// base.sets
u = {...}
U = u
phi = {}
Phi = phi
union({A sub u}, {B sub u}, ...): // ...
intersect({A sub u}, {B sub u}, ...): //...
///... others
```

And the implementation for instance (dummy code) could

```go
_ = ctx.Scope.Domains.Add(constants.BaseDomainPath)
```

This allows any code that uses the implementation for transpiling or diagnostics (like an lsp server) to perform

```sets
phi in base
u in base
union in base
intersect in base
```

Instead of relying the implementation to make the symbols special and handle differently.

**Difference between `sub` and `in` in this context**

`sub` means members of lhs is also a member of rhs. Therefore, `phi in base` and `phi sub base` are not the same. When saying `phi in base` you are declaring `phi` is a member of `base` which an implementation can check. But saying `phi sub base` while valid `.sets` code, just declares a new domain named `phi` which at most can be the same as whole of `base`, essentially becoming an alias of `base` itself.

## Singletons

A singleton domain contains exactly one element.

A singleton can be constructed using ordinary roster notation:

```sets
{1}
```

or expressed through a helper function:

```sets
singleton({x suchthat x in u}): {x}
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

> It is critical to understand, `.sets` has no capability to define cardinality of a domain yet. So a call like `singleton({1, 2})` is totally valid based on the transformation parameter domain itself.

> An implementation may be cardinality aware for diagnostic purposes, but the language grammar does not recognize it.

> For example a domain using roster notation can have its cardinality tracked under the diagnostic system to further limit callsite usage.

> `Foo = { 1, 2, 3 }`, a function `f({x sub or = Foo}): x%2 = 0`, at callsite `f({1, 2, 3, 4})` would be invalid without evaluating

> the expression, through cardinality of parameter domain and argument domain.

## Domain Validation

`.sets` does not permit unknown domains.

Every domain referenced by a declaration must either:

* be introduced by a top-level domain declaration;
* be constructed by a domain expression;
* be provided by the standard domain environment; or
* explicitly resolve to the universal domain.

## Tuples and Vectors

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
