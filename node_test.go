package jsontree

import (
	"encoding/json"
	"errors"
	"strconv"
	"testing"

	"github.com/rickb777/expect"
)

func TestOption(t *testing.T) {
	t.Run("AsString good", func(t *testing.T) {
		expect.Value(Some[any]("foo").AsString()).ToBe(t, Some[string]("foo"))
	})

	t.Run("AsString bad", func(t *testing.T) {
		expect.Value(Some[any](true).AsString()).ToBe(t,
			Option[string]{Err: errors.New("bool is not a string")})
		expect.Value(None[any]().AsString()).ToBe(t,
			Option[string]{Err: ErrAbsent})
	})

	t.Run("AsInt good", func(t *testing.T) {
		expect.Value(Some[any](44).AsInt()).ToBe(t, Some[int](44))
		expect.Value(Some[any](json.Number("45")).AsInt()).ToBe(t, Some[int](45))
	})

	t.Run("AsInt bad", func(t *testing.T) {
		expect.Value(Some[any](json.Number("xyz")).AsInt()).ToBe(t,
			Option[int]{Err: &strconv.NumError{Func: "ParseInt", Num: "xyz", Err: errors.New("invalid syntax")}})
		expect.Value(None[any]().AsInt()).ToBe(t, None[int]())
	})

	t.Run("AsFloat64 good", func(t *testing.T) {
		expect.Value(Some[any](44).AsFloat64()).ToBe(t, Some[float64](44))
		expect.Value(Some[any](float64(44)).AsFloat64()).ToBe(t, Some[float64](44))
		expect.Value(Some[any](json.Number("45")).AsFloat64()).ToBe(t, Some[float64](45))
	})

	t.Run("AsFloat64 bad", func(t *testing.T) {
		expect.Value(Some[any](json.Number("xyz")).AsFloat64()).ToBe(t,
			Option[float64]{Err: &strconv.NumError{Func: "ParseFloat", Num: "xyz", Err: errors.New("invalid syntax")}})
		expect.Value(None[any]().AsFloat64()).ToBe(t, None[float64]())
	})

	t.Run("AsBool good", func(t *testing.T) {
		expect.Value(Some[any](false).AsBool()).ToBe(t, Some[bool](false))
		expect.Value(Some[any](true).AsBool()).ToBe(t, Some[bool](true))
	})

	t.Run("AsBool bad", func(t *testing.T) {
		expect.Value(Some[any]("true").AsBool()).ToBe(t,
			Option[bool]{Err: errors.New("string is not a bool")})
		expect.Value(Some[any](json.Number("xyz")).AsBool()).ToBe(t,
			Option[bool]{Err: errors.New("json.Number is not a bool")})
		expect.Value(None[any]().AsBool()).ToBe(t, None[bool]())
	})

	t.Run("AsStrings good", func(t *testing.T) {
		expect.Value(Some[any]([]any{"1", "2", "3"}).AsStrings()).ToBe(t, Some[[]string]([]string{"1", "2", "3"}))
	})

	t.Run("AsStrings bad", func(t *testing.T) {
		expect.Value(Some[any]([]bool{true, false}).AsStrings()).ToBe(t,
			Option[[]string]{Err: errors.New("[]bool is not []string")})
	})

	t.Run("AsInts good", func(t *testing.T) {
		expect.Value(Some[any]([]any{1, 2, 3}).AsInts()).ToBe(t, Some[[]int]([]int{1, 2, 3}))
		expect.Value(Some[any]([]any{1.0, 2.0, 3.0}).AsInts()).ToBe(t, Some[[]int]([]int{1, 2, 3}))
	})

	t.Run("AsInts bad", func(t *testing.T) {
		expect.Value(Some[any]([]any{"1", "2", "3"}).AsInts()).ToBe(t,
			Option[[]int]{Err: errors.New("[]interface {} is not []int")})
		expect.Value(Some[any]([]bool{true, false}).AsInts()).ToBe(t,
			Option[[]int]{Err: errors.New("[]bool is not []int")})
	})

	t.Run("AsFloat64s good", func(t *testing.T) {
		expect.Value(Some[any]([]any{1, 2, 3}).AsFloat64s()).ToBe(t, Some[[]float64]([]float64{1, 2, 3}))
		expect.Value(Some[any]([]any{1.0, 2.0, 3.0}).AsFloat64s()).ToBe(t, Some[[]float64]([]float64{1, 2, 3}))
	})

	t.Run("AsFloat64s bad", func(t *testing.T) {
		expect.Value(Some[any]([]any{"1", "2", "3"}).AsFloat64s()).ToBe(t,
			Option[[]float64]{Err: errors.New("[]interface {} is not []float64")})
		expect.Value(Some[any]([]bool{true, false}).AsFloat64s()).ToBe(t,
			Option[[]float64]{Err: errors.New("[]bool is not []float64")})
	})

	t.Run("CoerceString good", func(t *testing.T) {
		expect.Value(Some[any]("foo").CoerceString()).ToBe(t, Some[string]("foo"))
		expect.Value(Some[any](123).CoerceString()).ToBe(t, Some[string]("123"))
		expect.Value(Some[any](float64(123.4)).CoerceString()).ToBe(t, Some[string]("123.4"))
		expect.Value(Some[any](json.Number("23.45")).CoerceString()).ToBe(t, Some[string]("23.45"))
	})

	t.Run("CoerceString bad", func(t *testing.T) {
		expect.Value(Some[any](true).CoerceString()).ToBe(t,
			Option[string]{Err: errors.New("cannot coerce bool to string")})
		expect.Value(None[any]().CoerceString()).ToBe(t, None[string]())
	})

	t.Run("CoerceInt good", func(t *testing.T) {
		expect.Value(Some[any](34).CoerceInt()).ToBe(t, Some[int](34))
		expect.Value(Some[any]("35").CoerceInt()).ToBe(t, Some[int](35))
		expect.Value(Some[any](float64(36)).CoerceInt()).ToBe(t, Some[int](36))
		expect.Value(Some[any](json.Number("37")).CoerceInt()).ToBe(t, Some[int](37))
	})

	t.Run("CoerceInt bad", func(t *testing.T) {
		expect.Value(Some[any](json.Number("xyz")).CoerceInt()).ToBe(t,
			Option[int]{Err: &strconv.NumError{Func: "ParseInt", Num: "xyz", Err: errors.New("invalid syntax")}})
		expect.Value(Some[any](true).CoerceInt()).ToBe(t,
			Option[int]{Err: errors.New("cannot coerce bool to int")})
		expect.Value(None[any]().CoerceInt()).ToBe(t, None[int]())
	})

	t.Run("CoerceFloat64 good", func(t *testing.T) {
		expect.Value(Some[any](44).CoerceFloat64()).ToBe(t, Some[float64](44))
		expect.Value(Some[any]("45").CoerceFloat64()).ToBe(t, Some[float64](45))
		expect.Value(Some[any](float64(46)).CoerceFloat64()).ToBe(t, Some[float64](46))
		expect.Value(Some[any](json.Number("47")).CoerceFloat64()).ToBe(t, Some[float64](47))
	})

	t.Run("CoerceFloat64 bad", func(t *testing.T) {
		expect.Value(Some[any](json.Number("xyz")).CoerceFloat64()).ToBe(t,
			Option[float64]{Err: &strconv.NumError{Func: "ParseFloat", Num: "xyz", Err: errors.New("invalid syntax")}})
		expect.Value(Some[any](true).CoerceFloat64()).ToBe(t,
			Option[float64]{Err: errors.New("cannot coerce bool to float64")})
		expect.Value(None[any]().CoerceFloat64()).ToBe(t, None[float64]())
	})

	t.Run("CoerceBool good", func(t *testing.T) {
		expect.Value(Some[any](true).CoerceBool()).ToBe(t, Some[bool](true))
		expect.Value(Some[any]("true").CoerceBool()).ToBe(t, Some[bool](true))
	})

	t.Run("CoerceBool bad", func(t *testing.T) {
		expect.Value(Some[any](json.Number("xyz")).CoerceBool()).ToBe(t,
			Option[bool]{Err: errors.New("cannot coerce json.Number to bool")})
		expect.Value(Some[any](101).CoerceBool()).ToBe(t,
			Option[bool]{Err: errors.New("cannot coerce int to bool")})
		expect.Value(None[any]().CoerceBool()).ToBe(t, None[bool]())
	})

	t.Run("CoerceStrings good", func(t *testing.T) {
		expect.Value(Some[any]([]any{"1", "2", "3"}).CoerceStrings()).ToBe(t, Some[[]string]([]string{"1", "2", "3"}))
		expect.Value(Some[any]([]any{1, 2, 3}).CoerceStrings()).ToBe(t, Some[[]string]([]string{"1", "2", "3"}))
		expect.Value(Some[any]([]any{1.0, 2.0, 3.0}).CoerceStrings()).ToBe(t, Some[[]string]([]string{"1", "2", "3"}))
	})

	t.Run("CoerceStrings bad", func(t *testing.T) {
		expect.Value(Some[any]([]bool{true, false}).CoerceStrings()).ToBe(t,
			Option[[]string]{Err: errors.New("[]bool is not []string")})
		expect.Value(None[any]().CoerceStrings()).ToBe(t, None[[]string]())
	})

	t.Run("CoerceInts good", func(t *testing.T) {
		expect.Value(Some[any]([]any{1, 2, 3}).CoerceInts()).ToBe(t, Some[[]int]([]int{1, 2, 3}))
		expect.Value(Some[any]([]int{1, 2, 3}).CoerceInts()).ToBe(t, Some[[]int]([]int{1, 2, 3}))
		expect.Value(Some[any]([]any{1.0, 2.0, 3.0}).CoerceInts()).ToBe(t, Some[[]int]([]int{1, 2, 3}))
		expect.Value(Some[any]([]any{"1", "2", "3"}).CoerceInts()).ToBe(t, Some[[]int]([]int{1, 2, 3}))
	})

	t.Run("CoerceInts bad", func(t *testing.T) {
		expect.Value(Some[any]([]bool{true, false}).CoerceInts()).ToBe(t,
			Option[[]int]{Err: errors.New("[]bool is not []int")})
		expect.Value(None[any]().CoerceInts()).ToBe(t, None[[]int]())
	})

	t.Run("CoerceFloat64s good", func(t *testing.T) {
		expect.Value(Some[any]([]any{1, 2, 3}).CoerceFloat64s()).ToBe(t, Some[[]float64]([]float64{1, 2, 3}))
		expect.Value(Some[any]([]any{1.0, 2.0, 3.0}).CoerceFloat64s()).ToBe(t, Some[[]float64]([]float64{1, 2, 3}))
		expect.Value(Some[any]([]float64{1.0, 2.0, 3.0}).CoerceFloat64s()).ToBe(t, Some[[]float64]([]float64{1, 2, 3}))
		expect.Value(Some[any]([]any{"1", "2", "3"}).CoerceFloat64s()).ToBe(t, Some[[]float64]([]float64{1, 2, 3}))
	})

	t.Run("CoerceFloat64s bad", func(t *testing.T) {
		expect.Value(Some[any]([]bool{true, false}).CoerceFloat64s()).ToBe(t,
			Option[[]float64]{Err: errors.New("[]bool is not []float64")})
		expect.Value(None[any]().CoerceFloat64s()).ToBe(t, None[[]float64]())
	})
}
