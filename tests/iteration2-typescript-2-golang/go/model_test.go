// Go-side test vectors (VECTORS.md §4, G-*). Each case builds the Go value
// that GO-VAL-1 gives for a .sets member and checks it compiles and holds.

package model

import (
	"math"
	"math/big"
	"testing"
)

func TestG01ExactDecimalConstantsRoundToTheSameDouble(t *testing.T) {
	if float64(PiValue) != 3.14 {
		t.Fatalf("PiValue = %v", PiValue)
	}
	if float64(TenthValue) != 0.1 {
		t.Fatalf("TenthValue = %v", TenthValue)
	}
}

func TestG02LargeWholeNumbersAreExact(t *testing.T) {
	// 1000000000000000000000 and 8640000000000000 are Integer-form members of TsNumber.
	if TsNumber(1e21) != TsNumber(1000000000000000000000) {
		t.Fatal("1e21 not exact")
	}
	var d TsDate = 8640000000000000
	if float64(d) != 8.64e15 {
		t.Fatal("max date not exact")
	}
}

func TestG03IEEESpecials(t *testing.T) {
	var nan TsDate = TsDate(math.NaN())
	if !math.IsNaN(float64(nan)) {
		t.Fatal("NotANumber")
	}
	negZero := TsNumber(math.Copysign(0, -1))
	if !math.Signbit(float64(negZero)) {
		t.Fatal("NegativeZero")
	}
}

func TestG04UserThreeAbsences(t *testing.T) {
	missing := User{Id: "u1", Name: "A", Tags: []string{}}
	undef := missing
	undef.Age = Opt[Either[TsNumber, TsUndefined]]{Present: true, Value: Either[TsNumber, TsUndefined]{IsB: true, B: TsUndefinedValue}}
	num := missing
	num.Age = Opt[Either[TsNumber, TsUndefined]]{Present: true, Value: Either[TsNumber, TsUndefined]{A: 30}}
	if missing.Age.Present || !undef.Age.Value.IsB || num.Age.Value.A != 30 {
		t.Fatal("absence states collapsed")
	}
	if missing.Email != nil {
		t.Fatal("email Nothing should be nil")
	}
}

func TestG05OpenRecordKeepsExtraKeys(t *testing.T) {
	u := User{Id: "u1", Extra: map[TsKey]any{"nick": "a"}}
	if u.Extra["nick"] != "a" {
		t.Fatal("extra key lost")
	}
}

func TestG06TaggedUnions(t *testing.T) {
	shapes := []Shape{ShapeCircle{Radius: 1}, ShapeRect{W: 1, H: 2}}
	if _, ok := shapes[0].(ShapeCircle); !ok {
		t.Fatal("circle")
	}
	var r Result[TsNumber, string] = ResultFalse[TsNumber, string]{Error: "e"}
	if _, ok := r.(ResultFalse[TsNumber, string]); !ok {
		t.Fatal("result")
	}
}

func TestG07SequencesAndTuples(t *testing.T) {
	p := Point{1, 2}
	a := Ambiguous{1, 2} // one Go value for the JS value [1, 2]
	if p[1] != a[1] {
		t.Fatal("point vs ambiguous")
	}
	short := Span{F0: 1}
	if short.F1.Present {
		t.Fatal("span length 1")
	}
	args := Args{F0: "a", Rest: []TsNumber{1, 2}}
	if len(args.Rest) != 2 {
		t.Fatal("args")
	}
}

func TestG08BigAndSets(t *testing.T) {
	var b Big = big.NewInt(1)
	if b.Int64() != 1 {
		t.Fatal("big")
	}
	s := TagSet{"a": {}}
	if _, ok := s["a"]; !ok {
		t.Fatal("set")
	}
	tree := Tree[string]{Value: "a", Children: []Tree[string]{{Value: "b"}}}
	if tree.Children[0].Value != "b" {
		t.Fatal("tree")
	}
	var j Json = map[string]Json{"a": []Json{1.0, nil}}
	if _, ok := j.(map[string]Json); !ok {
		t.Fatal("json")
	}
}
