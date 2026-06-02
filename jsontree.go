package jsontree

import (
	jsonpkg "encoding/json"
	"fmt"
	io "io"
	"strconv"
	"strings"
)

// TreeNode is for traversing JSON objects and arrays without needing to type-convert each sub-level.
//
//	// tree is JSON `{"meta":{"status":"OK"}}`
//	_ = TreeNode(tree, "meta", "status").AsString() // returns Some("OK")
//
// or a JSON array might be used:
//
//	// array is JSON `[[3,5,8], [1,9,0,0], {"ok":true}]`
//	_ = TreeNode(array, 0, 2).AsInt() // returns Some(8)
//
// The first key should be a string. For each deeper node, the key should be a
// string if the node is a JSON object, or an integer if the node is a JSON array.
//
// The input can be a literal JSON string or an [io.Reader] that provides JSON text. In this case,
// the JSON will be parsed using the standard JSON decoder. Note that the input must be
// well-formed JSON; this function will panic it there is a parse error.
func TreeNode(json any, key ...any) Option[any] {
	switch s := json.(type) {
	case string:
		return TreeNode(strings.NewReader(s), key...)

	case io.Reader:
		d := jsonpkg.NewDecoder(s)
		d.UseNumber()

		var decoded any
		err := d.Decode(&decoded)
		if err != nil {
			panic("TreeNode cannot parse JSON: " + err.Error())
		}
		json = decoded
	}

	return treeNode(json, 0, key)
}

func treeNode(node any, ki int, keys []any) Option[any] {
	if len(keys) > ki {
		switch t1 := node.(type) {
		case map[string]any:
			return traverseMap(t1, ki, keys)
		case []any:
			return traverseArray(t1, ki, keys)
		}
	}

	return notFound(keys)
}

func traverseMap(t1 map[string]any, ki int, keys []any) Option[any] {
	switch kk := keys[ki].(type) {
	case string:
		v, ok := t1[kk]
		if !ok {
			return notFound(keys)
		}

		if len(keys) == ki+1 {
			return Some[any](v)
		}

		switch t2 := v.(type) {
		case map[string]any:
			return treeNode(t2, ki+1, keys)
		case []any:
			return treeNode(t2, ki+1, keys)
		default:
			return notFound(keys)
		}
	}

	return notFound(keys)
}

func traverseArray(t1 []any, ki int, keys []any) Option[any] {
	switch kk := keys[ki].(type) {
	case int:
		return arrayElement(t1[kk], ki+1, keys)
	case float64:
		return arrayElement(t1[int(kk)], ki+1, keys)
	case string:
		i, err := strconv.Atoi(kk)
		if err == nil {
			return arrayElement(t1[i], ki+1, keys)
		}
	}
	return notFound(keys)
}

func arrayElement(node any, ki int, keys []any) Option[any] {
	if ki == len(keys) {
		return Some[any](node)
	}

	switch v := node.(type) {
	case map[string]any:
		return treeNode(v, ki, keys)
	case []any:
		return treeNode(v, ki, keys)
	default:
		return Some[any](v)
	}
}

func notFound(keys []any) Option[any] {
	return Option[any]{Err: fmt.Errorf("%s not found", strings.Join(coerceStringSlice(keys), ","))}
}

func coerceStringSlice(vv []any) []string {
	ss := make([]string, len(vv))
	for i, v := range vv {
		ss[i] = fmt.Sprintf("%v", v)
	}
	return ss
}
