// Generated from sets/types.sets by sets-to-go/STANDARD.md 0.2.0.
// Block numbers match the .sets file. Each declaration names its rule.

package model

import "math/big"

// 1.

type Id = string      // GO-DECL-1 alias
type Count = TsNumber // GO-DECL-1 alias
type Big = *big.Int   // GO-BASE-1, alias (pointer)

// Pi: a one-member roster in Binary64Finite (GO-ROSTER-1).
type Pi float64

const PiValue Pi = 3.140000000000000124344978758017532527446746826171875

// Tenth: a one-member roster in Binary64Finite (GO-ROSTER-1).
type Tenth float64

const TenthValue Tenth = 0.1000000000000000055511151231257827021181583404541015625

// One: an Integer roster (GO-ROSTER-1).
type One int8

const OneValue One = 1

// OneN: an Integer roster (GO-ROSTER-1). Its definition is the same as One's.
type OneN int8

const OneNValue OneN = 1

// 2.

// Status: a String roster (GO-ROSTER-1).
type Status string

const (
	StatusActive    Status = "active"
	StatusSuspended Status = "suspended"
	StatusDeleted   Status = "deleted"
)

// Color: a String roster (GO-ROSTER-1).
type Color string

const (
	ColorRED   Color = "RED"
	ColorGREEN Color = "GREEN"
)

// Level: an Integer roster (GO-ROSTER-1).
type Level int8

const (
	Level0 Level = 0
	Level1 Level = 1
	Level2 Level = 2
)

// 3. User: a record builder (GO-REC-1).
type User struct {
	Id        Id                                 `sets:"id"`
	Name      string                             `sets:"name"`
	Age       Opt[Either[TsNumber, TsUndefined]] `sets:"age"`   // OptionalField + GO-UNION-6
	Email     *string                            `sets:"email"` // GO-UNION-1
	Tags      []string                           `sets:"tags"`
	CreatedAt TsDate                             `sets:"createdAt"`
	Extra     map[TsKey]any                      `sets:",extra"`
}

// 4. Shape: tagged records on key "kind" (GO-UNION-5).
type Shape interface{ isShape() }

type ShapeCircle struct { // kind = "circle"
	Radius TsNumber      `sets:"radius"`
	Extra  map[TsKey]any `sets:",extra"`
}

type ShapeRect struct { // kind = "rect"
	W     TsNumber      `sets:"w"`
	H     TsNumber      `sets:"h"`
	Extra map[TsKey]any `sets:",extra"`
}

func (ShapeCircle) isShape() {}
func (ShapeRect) isShape()   {}

// 5.

// Point: a fixed-length sequence pattern, all elements the same type (GO-SEQ-1).
type Point [2]TsNumber

// Span: prefix sequence patterns (GO-UNION-4).
type Span struct {
	F0 TsNumber
	F1 Opt[Either[TsNumber, TsUndefined]]
}

// Args: a prefix and a rest (GO-SEQ-2).
type Args struct {
	F0   string
	Rest []TsNumber
}

// 6.

type Dictionary map[string]TsNumber  // GO-CALL-1
type ByIndex map[TsNumericKey]string // GO-CALL-1

// Headers: a record with Rest (GO-REC-1).
type Headers struct {
	ContentType string            `sets:"contentType"`
	Rest        map[string]string `sets:",rest"`
}

// Flags: a record builder (GO-REC-1).
type Flags struct {
	Read  bool          `sets:"read"`
	Write bool          `sets:"write"`
	Extra map[TsKey]any `sets:",extra"`
}

// 7.

// Box: a generic record (GO-DECL-1, GO-REC-1).
type Box[T any] struct {
	Value T             `sets:"value"`
	Extra map[TsKey]any `sets:",extra"`
}

// Result: generic tagged records on key "ok" (GO-UNION-5).
type Result[T, E any] interface{ isResult() }

type ResultTrue[T, E any] struct { // ok = true
	Value T             `sets:"value"`
	Extra map[TsKey]any `sets:",extra"`
}

type ResultFalse[T, E any] struct { // ok = false
	Error E             `sets:"error"`
	Extra map[TsKey]any `sets:",extra"`
}

func (ResultTrue[T, E]) isResult()  {}
func (ResultFalse[T, E]) isResult() {}

// Page: a generic record. {T sub TsObject} becomes any (GO-DECL-1, wider).
type Page[T any] struct {
	Items []T           `sets:"items"`
	Next  *string       `sets:"next"`
	Extra map[TsKey]any `sets:",extra"`
}

type UserPage = Page[User] // GO-DECL-1 alias (instantiation)

// 8.

// Tree: a generic recursive record through a slice (GO-RECUR-1).
type Tree[T any] struct {
	Value    T             `sets:"value"`
	Children []Tree[T]     `sets:"children"`
	Extra    map[TsKey]any `sets:",extra"`
}

// Json: a recursive 6-branch union with Nothing (GO-UNION-1, GO-UNION-7).
// Branches: String, TsNumber, Boolean, Seq(Json), Map(String, Json); Nothing is nil.
type Json = any

// 9.

type TagSet map[string]struct{} // GO-CALL-1 SetOf
type Scores map[string]TsNumber // GO-CALL-1

// 10.

// Timestamped: a record builder (GO-REC-1).
type Timestamped struct {
	UpdatedAt TsDate        `sets:"updatedAt"`
	Extra     map[TsKey]any `sets:",extra"`
}

// AuditedUser = intersect(User, Timestamped): fields merged (GO-INTER-1).
type AuditedUser struct {
	Id        Id                                 `sets:"id"`
	Name      string                             `sets:"name"`
	Age       Opt[Either[TsNumber, TsUndefined]] `sets:"age"`
	Email     *string                            `sets:"email"`
	Tags      []string                           `sets:"tags"`
	CreatedAt TsDate                             `sets:"createdAt"`
	UpdatedAt TsDate                             `sets:"updatedAt"`
	Extra     map[TsKey]any                      `sets:",extra"`
}

// 11.

// UserKey sub String (GO-DECL-1).
// Exact members: every String whose first five characters are `user-`.
type UserKey string

// 12. Account: a record builder (GO-REC-1).
type Account struct {
	Id       string        `sets:"id"`
	Balance  TsNumber      `sets:"balance"`
	OnChange TsFunction    `sets:"onChange"`
	Extra    map[TsKey]any `sets:",extra"`
}

// 13.

// Ambiguous: Point is visibly inside Seq(TsNumber), so it's dropped (GO-UNION-0).
type Ambiguous []TsNumber

// EmptyThings: Nothing, then two branches (GO-UNION-1, GO-UNION-6). Alias (pointer).
type EmptyThings = *Either[TsNonNullish, TsUndefined]

// Nullish: the duplicate is dropped, then Nothing (GO-UNION-0, GO-UNION-1).
type Nullish = *TsUndefined

// SetOrDict: two branches (GO-UNION-6). {} is in both and is held as A.
type SetOrDict Either[map[string]struct{}, map[string]TsNonNullish]
