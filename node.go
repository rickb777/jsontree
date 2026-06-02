package jsontree

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Option[T any] struct {
	V   T
	Err error
}

var ErrAbsent = fmt.Errorf("absent")

func None[T any]() Option[T]    { return Option[T]{Err: ErrAbsent} }
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
		return float64OrError(n.Float64())
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
		return float64OrError(n.Float64())
	// coercive types
	case string:
		return float64OrError(strconv.ParseFloat(n, 64))
	}
	return cannotCoerce[float64](v, "float64")
}

func float64OrError(f float64, err error) Option[float64] {
	if err != nil {
		return Option[float64]{Err: err}
	}
	return Some[float64](f)
}

//-------------------------------------------------------------------------------------------------

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
