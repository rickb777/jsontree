package jsontree

import (
	"encoding/json"
	"fmt"
)

type Option[T any] struct {
	V       T
	Present bool
}

func None[T any]() Option[T]    { return Option[T]{} }
func Some[T any](v T) Option[T] { return Option[T]{V: v, Present: true} }

func (o Option[T]) String() string {
	if !o.Present {
		return "None"
	}
	return fmt.Sprintf("Some(%v)", o.V)
}

func (o Option[T]) AsString() Option[string] {
	if o.Present {
		v := any(o.V)
		switch s := v.(type) {
		case string:
			return Some[string](s)
		}
	}
	return None[string]()
}

func (o Option[T]) AsInt() Option[int] {
	if o.Present {
		v := any(o.V)
		switch n := v.(type) {
		case int:
			return Some[int](n)
		case float64:
			return Some[int](int(n))
		case json.Number:
			i, err := n.Int64()
			if err == nil {
				return Some[int](int(i))
			}
		}
	}
	return None[int]()
}

func (o Option[T]) AsFloat64() Option[float64] {
	if o.Present {
		v := any(o.V)
		switch n := v.(type) {
		case int:
			return Some[float64](float64(n))
		case float64:
			return Some[float64](n)
		case json.Number:
			f, err := n.Float64()
			if err == nil {
				return Some[float64](f)
			}
		}
	}
	return None[float64]()
}
