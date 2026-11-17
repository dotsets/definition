// Input for iteration 2 (copied unchanged from iteration 1).
// Each block is numbered; sets/types.sets and go/model.go use the same numbers.

// 1. Primitives and number literals
type Id = string;
type Count = number;
type Big = bigint;
type Pi = 3.14;
type Tenth = 0.1;
type One = 1;
type OneN = 1n;

// 2. Literal unions and enums
type Status = "active" | "suspended" | "deleted";
enum Color { Red = "RED", Green = "GREEN" }
enum Level { Low, Mid, High }

// 3. Interface with optional, nullable, readonly, array and Date fields
interface User {
  id: Id;
  name: string;
  age?: number;
  email: string | null;
  tags: string[];
  readonly createdAt: Date;
}

// 4. Discriminated union
type Shape =
  | { kind: "circle"; radius: number }
  | { kind: "rect"; w: number; h: number };

// 5. Tuples: fixed, optional element, rest element
type Point = [number, number];
type Span = [start: number, end?: number];
type Args = [string, ...number[]];

// 6. Index signatures
type Dictionary = { [k: string]: number };
type ByIndex = { [k: number]: string };
type Headers = { contentType: string; [k: string]: string };
type Flags = Record<"read" | "write", boolean>;

// 7. Generics
type Box<T> = { value: T };
type Result<T, E = string> = { ok: true; value: T } | { ok: false; error: E };
type Page<T extends object> = { items: T[]; next: string | null };
type UserPage = Page<User>;

// 8. Recursion
interface Tree<T> { value: T; children: Tree<T>[] }
type Json = string | number | boolean | null | Json[] | { [k: string]: Json };

// 9. Library wrappers
type TagSet = Set<string>;
type Scores = Map<string, number>;

// 10. Intersection
type Timestamped = { updatedAt: Date };
type AuditedUser = User & Timestamped;

// 11. Template literal
type UserKey = `user-${string}`;

// 12. Class with private field, method and function-typed field
class Account {
  id: string = "";
  private balance: number = 0;
  onChange: (n: number) => void = () => {};
  deposit(n: number): void { this.balance += n; }
}

// 13. Collision probes
type Ambiguous = [number, number] | number[];
type EmptyThings = {} | null | undefined;
type Nullish = null | undefined | void;
type SetOrDict = Set<string> | Record<string, {}>;
