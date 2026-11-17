// Generated from sets/ts.sets by sets-to-go/STANDARD.md 0.2.0.
// Each declaration names the rule that produced it.

package model

import "math/big"

// TsSymbol is an opaque top-level axiom (GO-DECL-1).
type TsSymbol struct{ Handle any }

// TsFunction is an opaque top-level axiom (GO-DECL-1).
type TsFunction struct{ Handle any }

// TsUndefined = {undefined}: an atom roster (GO-ROSTER-1).
type TsUndefined string

const TsUndefinedValue TsUndefined = "undefined"

// TsWholeFractional: a refinement of Fractional (GO-REF-1), an alias because it's a pointer.
type TsWholeFractional = *big.Rat

// TsNumber: a refinement of Float64 (GO-REF-1). Wider.
type TsNumber float64

// TsKey = String (GO-DECL-1 alias).
type TsKey = string

// TsRecord = Map(TsKey, U) (GO-CALL-1).
type TsRecord map[TsKey]any

// TsNumericKey sub String (GO-DECL-1). Exact members: the strings that
// JavaScript's String(n) returns for some number n, for example "0", "-1",
// "1.5", "1e+21", "NaN", "Infinity", "-Infinity". String(-0) is "0", so "-0"
// isn't a member.
type TsNumericKey string

// TsDate = union(refinement of TsNumber, {NotANumber}) (GO-UNION-3).
type TsDate TsNumber

// TsPrimitive: a 7-branch union with Nothing (GO-UNION-1, GO-UNION-7).
// Branches: String, TsNumber, Integer, Boolean, TsSymbol, TsUndefined; Nothing is nil.
type TsPrimitive = any

// TsValue: a recursive 5-branch union (GO-UNION-7).
// Branches: TsPrimitive, Map(TsValue, TsValue), Seq(TsValue), SetOf(TsValue), TsFunction.
type TsValue = any

// TsObject: a refinement of TsValue (GO-REF-1).
type TsObject = any

// TsNonNullish: a refinement of TsValue (GO-REF-1).
type TsNonNullish = any
