package jsontree

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Option wraps some value, which may be present or absent, of some type.
type Option[T any] struct {
	// V is the value contained, if present.
	V T
	// Err is nil for values that are present; when non-nil, the
	// optional value is absent and this error provides a reason.
	Err error
}

var ErrAbsent = fmt.Errorf("absent")

// None creates an empty optional value.
func None[T any]() Option[T] { return Option[T]{Err: ErrAbsent} }

// Some wraps some value to create an option that is present.
func Some[T any](v T) Option[T] { return Option[T]{V: v} }

func wrongType[T any](v any, exp string) Option[T] {
	return Option[T]{Err: fmt.Errorf("%T is not %s", v, exp)}
}

func cannotCoerce[T any](v any, exp string) Option[T] {
	return Option[T]{Err: fmt.Errorf("cannot coerce %T to %s", v, exp)}
}

func (o Option[T]) String() string {
	if o.Err != nil {
		return fmt.Sprintf("None(%v)", o.Err)
	}
	return fmt.Sprintf("Some(%v)", o.V)
}

// Present returns true iff the option contains a value. Otherwise, [Option.Err] will contain
// more information about why not.
func (o Option[T]) Present() bool {
	return o.Err == nil
}

//-------------------------------------------------------------------------------------------------

// SubTree traverses deeper from some intermediate node. For example if
//
//	response := TreeNode(tree, "response")
//
// does not return a leaf node but instead returns an intermediate node, then
//
//	response.SubTree("user", "auth_token")
//
// can traverse to nodes below response. It does not matter whether the
// intermediate node is a JSON object or a JSON array.
//
// For objects, the keys should be strings; for arrays, the keys should be integers.
func (o Option[T]) SubTree(keys ...any) Option[any] {
	if o.Err != nil {
		return Option[any]{Err: o.Err}
	}

	switch t1 := any(o.V).(type) {
	case map[string]any:
		return treeNode(t1, 0, keys)
	case []any:
		return traverseArray(t1, 0, keys)
	}

	return notFound(keys)
}

//-------------------------------------------------------------------------------------------------

// AsString obtains an optional string, provided that o contains a string value.
func (o Option[T]) AsString() Option[string] {
	if o.Err != nil {
		return Option[string]{Err: o.Err}
	}
	return asString(o.V)
}

func asString(v any) Option[string] {
	switch s := v.(type) {
	case string:
		return Some[string](s)
	}
	return wrongType[string](v, "a string")
}

// CoerceString obtains an optional string, converting to string if required.
func (o Option[T]) CoerceString() Option[string] {
	if o.Err != nil {
		return Option[string]{Err: o.Err}
	}
	return coerceString(o.V)
}

func coerceString(v any) Option[string] {
	switch s := v.(type) {
	case string:
		return Some[string](s)
	// coercive types
	case int:
		return Some[string](strconv.Itoa(s))
	case float64:
		return Some[string](strconv.FormatFloat(s, 'g', -1, 64))
	case json.Number:
		return Some[string](s.String())
	}
	return cannotCoerce[string](v, "string")
}

//-------------------------------------------------------------------------------------------------

// AsInt obtains an optional int, provided that o contains a numeric value.
func (o Option[T]) AsInt() Option[int] {
	if o.Err != nil {
		return Option[int]{Err: o.Err}
	}
	return asInt(o.V)
}

func asInt(v any) Option[int] {
	switch n := v.(type) {
	case int:
		return Some[int](n)
	case float64:
		return Some[int](int(n))
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			return Option[int]{Err: err}
		}
		return Some[int](int(i))
	}
	return wrongType[int](v, "an int")
}

// CoerceInt obtains an optional int, converting to int if required.
func (o Option[T]) CoerceInt() Option[int] {
	if o.Err != nil {
		return Option[int]{Err: o.Err}
	}
	return coerceInt(o.V)
}

func coerceInt(v any) Option[int] {
	switch n := v.(type) {
	case int:
		return Some[int](n)
	case float64:
		return Some[int](int(n))
	case json.Number:
		return intOrError(n.Int64())
	// coercive types
	case string:
		return intOrError(strconv.ParseInt(n, 10, 64))
	}
	return cannotCoerce[int](v, "int")
}

func intOrError(i int64, err error) Option[int] {
	if err != nil {
		return Option[int]{Err: err}
	}
	return Some[int](int(i))
}

//-------------------------------------------------------------------------------------------------

// AsFloat64 obtains an optional float64, provided that o contains a numeric value.
func (o Option[T]) AsFloat64() Option[float64] {
	if o.Err != nil {
		return Option[float64]{Err: o.Err}
	}
	return asFloat64(o.V)
}

func asFloat64(v any) Option[float64] {
	switch n := v.(type) {
	case int:
		return Some[float64](float64(n))
	case float64:
		return Some[float64](n)
	case json.Number:
		return valueOrError(n.Float64())
	}
	return wrongType[float64](v, "a float64")
}

// CoerceFloat64 obtains an optional float64, converting to float64 if required.
func (o Option[T]) CoerceFloat64() Option[float64] {
	if o.Err != nil {
		return Option[float64]{Err: o.Err}
	}
	return coerceFloat64(o.V)
}

func coerceFloat64(v any) Option[float64] {
	switch n := v.(type) {
	case int:
		return Some[float64](float64(n))
	case float64:
		return Some[float64](n)
	case json.Number:
		return valueOrError(n.Float64())
	// coercive types
	case string:
		return valueOrError(strconv.ParseFloat(n, 64))
	}
	return cannotCoerce[float64](v, "float64")
}

//-------------------------------------------------------------------------------------------------

// AsBool obtains an optional bool, provided that o contains a boolean value.
func (o Option[T]) AsBool() Option[bool] {
	if o.Err != nil {
		return Option[bool]{Err: o.Err}
	}
	return asBool(o.V)
}

func asBool(v any) Option[bool] {
	switch n := v.(type) {
	case bool:
		return Some[bool](n)
	}
	return wrongType[bool](v, "a bool")
}

// CoerceBool obtains an optional bool, converting to bool if required.
func (o Option[T]) CoerceBool() Option[bool] {
	if o.Err != nil {
		return Option[bool]{Err: o.Err}
	}
	return coerceBool(o.V)
}

func coerceBool(v any) Option[bool] {
	switch n := v.(type) {
	case bool:
		return Some[bool](n)
	// coercive types
	case string:
		return valueOrError(strconv.ParseBool(n))
	}
	return cannotCoerce[bool](v, "bool")
}

func valueOrError[T string | int | float64 | bool](v T, err error) Option[T] {
	if err != nil {
		return Option[T]{Err: err}
	}
	return Some[T](v)
}

//=================================================================================================

// AsStrings obtains an optional string slice, provided that o contains a slice of strings.
// This only handles []string or []any; see [Option.CoerceStrings] for value formatting capability.
func (o Option[T]) AsStrings() Option[[]string] {
	if o.Err != nil {
		return Option[[]string]{Err: o.Err}
	}
	return asStrings(o.V)
}

func asStrings(v any) Option[[]string] {
	switch n := v.(type) {
	case []string:
		return Some[[]string](n)
	case []any:
		ii := make([]string, 0, len(n))
		for _, i := range n {
			oi := asString(i)
			if oi.Err != nil {
				return wrongType[[]string](n, "[]string")
			}
			ii = append(ii, oi.V)
		}
		return Some[[]string](ii)
	}
	return wrongType[[]string](v, "[]string")
}

// CoerceStrings obtains an optional string slice, provided that o contains a slice of strings, ints or float64s.
func (o Option[T]) CoerceStrings() Option[[]string] {
	if o.Err != nil {
		return Option[[]string]{Err: o.Err}
	}
	return coerceStrings(o.V)
}

func coerceStrings(v any) Option[[]string] {
	switch n := v.(type) {
	case []string:
		return Some[[]string](n)
	case []any:
		ii := make([]string, 0, len(n))
		for _, i := range n {
			oi := coerceString(i)
			if oi.Err != nil {
				return wrongType[[]string](n, "[]string")
			}
			ii = append(ii, oi.V)
		}
		return Some[[]string](ii)
	}
	return wrongType[[]string](v, "[]string")
}

//-------------------------------------------------------------------------------------------------

// AsInts obtains an optional int slice, provided that o contains a slice of numbers.
// This only handles []int or []any; see [Option.CoerceInts] for value parsing capability.
func (o Option[T]) AsInts() Option[[]int] {
	if o.Err != nil {
		return Option[[]int]{Err: o.Err}
	}
	return asInts(o.V)
}

func asInts(v any) Option[[]int] {
	switch n := v.(type) {
	case []int:
		return Some[[]int](n)
	case []any:
		ii := make([]int, 0, len(n))
		for _, i := range n {
			oi := asInt(i)
			if oi.Err != nil {
				return wrongType[[]int](n, "[]int")
			}
			ii = append(ii, oi.V)
		}
		return Some[[]int](ii)
	}
	return wrongType[[]int](v, "[]int")
}

// CoerceInts obtains an optional int slice, provided that o contains a slice of ints.
func (o Option[T]) CoerceInts() Option[[]int] {
	if o.Err != nil {
		return Option[[]int]{Err: o.Err}
	}
	return coerceInts(o.V)
}

func coerceInts(v any) Option[[]int] {
	switch n := v.(type) {
	case []int:
		return Some[[]int](n)
	case []any:
		ii := make([]int, 0, len(n))
		for _, i := range n {
			oi := coerceInt(i)
			if oi.Err != nil {
				return wrongType[[]int](n, "[]int")
			}
			ii = append(ii, oi.V)
		}
		return Some[[]int](ii)
	}
	return wrongType[[]int](v, "[]int")
}

//-------------------------------------------------------------------------------------------------

// AsFloat64s obtains an optional float64 slice, provided that o contains a slice of numbers.
// This only handles []float64 or []any; see [Option.CoerceFloat64s] for value parsing capability.
func (o Option[T]) AsFloat64s() Option[[]float64] {
	if o.Err != nil {
		return Option[[]float64]{Err: o.Err}
	}
	return asFloat64s(o.V)
}

func asFloat64s(v any) Option[[]float64] {
	switch n := v.(type) {
	case []float64:
		return Some[[]float64](n)
	case []any:
		ii := make([]float64, 0, len(n))
		for _, i := range n {
			oi := asFloat64(i)
			if oi.Err != nil {
				return wrongType[[]float64](n, "[]float64")
			}
			ii = append(ii, oi.V)
		}
		return Some[[]float64](ii)
	}
	return wrongType[[]float64](v, "[]float64")
}

// CoerceFloat64s obtains an optional float64 slice, provided that o contains a slice of float64s.
func (o Option[T]) CoerceFloat64s() Option[[]float64] {
	if o.Err != nil {
		return Option[[]float64]{Err: o.Err}
	}
	return coerceFloat64s(o.V)
}

func coerceFloat64s(v any) Option[[]float64] {
	switch n := v.(type) {
	case []float64:
		return Some[[]float64](n)
	case []any:
		ii := make([]float64, 0, len(n))
		for _, i := range n {
			oi := coerceFloat64(i)
			if oi.Err != nil {
				return wrongType[[]float64](n, "[]float64")
			}
			ii = append(ii, oi.V)
		}
		return Some[[]float64](ii)
	}
	return wrongType[[]float64](v, "[]float64")
}

//-------------------------------------------------------------------------------------------------

// AsBools obtains an optional bool slice, provided that o contains a slice of bools.
// This only handles []bool or []any; see [Option.CoerceBools] for value parsing capability.
func (o Option[T]) AsBools() Option[[]bool] {
	if o.Err != nil {
		return Option[[]bool]{Err: o.Err}
	}
	return asBools(o.V)
}

func asBools(v any) Option[[]bool] {
	switch n := v.(type) {
	case []bool:
		return Some[[]bool](n)
	case []any:
		ii := make([]bool, 0, len(n))
		for _, i := range n {
			oi := asBool(i)
			if oi.Err != nil {
				return wrongType[[]bool](n, "[]bool")
			}
			ii = append(ii, oi.V)
		}
		return Some[[]bool](ii)
	}
	return wrongType[[]bool](v, "[]bool")
}

// CoerceBools obtains an optional bool slice, provided that o contains a slice of bools.
func (o Option[T]) CoerceBools() Option[[]bool] {
	if o.Err != nil {
		return Option[[]bool]{Err: o.Err}
	}
	return coerceBools(o.V)
}

func coerceBools(v any) Option[[]bool] {
	switch n := v.(type) {
	case []bool:
		return Some[[]bool](n)
	case []any:
		ii := make([]bool, 0, len(n))
		for _, i := range n {
			oi := coerceBool(i)
			if oi.Err != nil {
				return wrongType[[]bool](n, "[]bool")
			}
			ii = append(ii, oi.V)
		}
		return Some[[]bool](ii)
	}
	return wrongType[[]bool](v, "[]bool")
}
