package jsontree

// TreeNode is for traversing JSON data without needing to type-convert each level. Example
//
//	// tree is JSON `{"meta":{"status":"OK"}}`
//	_ = AsString(TreeNode(tree, "meta", "status"))
func TreeNode(tree map[string]any, key ...any) Option[any] {
	return treeNode(tree, key...)
}

func treeNode(node any, key ...any) Option[any] {
	if len(key) > 0 {
		switch t1 := node.(type) {
		case map[string]any:
			switch kk := key[0].(type) {
			case string:
				v, ok := t1[kk]
				if !ok {
					return None[any]()
				}

				if len(key) == 1 {
					return Some[any](v)
				}

				switch t2 := v.(type) {
				case map[string]any:
					return treeNode(t2, key[1:]...)
				case []any:
					return treeNode(t2, key[1:]...)
				default:
					return None[any]()
				}
			}

		case []any:
			switch kk := key[0].(type) {
			case int:
				return arrayElement(t1[kk], key[1:]...)
			case float64:
				return arrayElement(t1[int(kk)], key[1:]...)
			}
		}
	}

	return None[any]()
}

func arrayElement(node any, key ...any) Option[any] {
	switch v := node.(type) {
	case map[string]any:
		return treeNode(v, key...)
	case []any:
		return treeNode(v, key...)
	default:
		return Some[any](v)
	}
}
