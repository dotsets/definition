// Runtime declarations required by sets-to-go/STANDARD.md GO-RT-1.

package model

// Opt is an optional record field: Present is false when the field's tuple is missing.
type Opt[T any] struct {
	Present bool
	Value   T
}

// Either holds a member of a two-branch union (GO-UNION-6).
type Either[A, B any] struct {
	IsB bool
	A   A
	B   B
}

// Pair is one entry of a map whose key type isn't value-comparable (GO-CALL-1).
type Pair[K, V any] struct {
	Key   K
	Value V
}
