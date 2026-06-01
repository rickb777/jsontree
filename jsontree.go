package jsontree

import (
	"fmt"
	"strconv"
	"strings"
)

// TreeNode is for traversing JSON data without needing to type-convert each level. Example
//
//	// tree is JSON `{"meta":{"status":"OK"}}`
//	_ = TreeNode(tree, "meta", "status").AsString()
func TreeNode(tree map[string]any, key ...any) Option[any] {
	return treeNode(tree, 0, key)
}

func treeNode(node any, ki int, keys []any) Option[any] {
	if len(keys) > ki {
		switch t1 := node.(type) {
		case map[string]any:
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

		case []any:
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
	return Option[any]{Err: fmt.Errorf("%s not found", strings.Join(coerceStrings(keys), ","))}
}

func coerceStrings(vv []any) []string {
	ss := make([]string, len(vv))
	for i, v := range vv {
		ss[i] = fmt.Sprintf("%v", v)
	}
	return ss
}
