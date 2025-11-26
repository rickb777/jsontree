package jsontree

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/rickb777/expect"
)

func ExampleTreeNode_builtin_numbers() {
	// some JSON
	s := `{
		"meta":{"code":200, "status":"OK"},
		"response":{
			"csrf_token":"x.y.z",
			"user":{"authentication_token":"a.b.c"}
		},
		"props":[
			{"a":1},
			{"b":2.1}
		],
		"nest":[
			[3,5,8]
		]
	}`

	// decode the JSON as a tree
	tree := map[string]any{}
	err := json.NewDecoder(strings.NewReader(s)).Decode(&tree)
	if err != nil {
		panic(err)
	}

	// traverse the tree for a selection of nodes
	fmt.Println(TreeNode(tree, "meta", "code").AsInt().V)
	fmt.Println(TreeNode(tree, "meta", "status").AsString().V)
	fmt.Println(TreeNode(tree, "meta", "status", "absent"))
	fmt.Println(TreeNode(tree, "response", "user", "authentication_token").AsString().V)
	fmt.Println(TreeNode(tree, "props", 0, "a").AsInt().V)
	fmt.Println(TreeNode(tree, "props", 1, "b").AsFloat64().V)
	fmt.Println(TreeNode(tree, "nest", 0, 2).AsInt().V)
	fmt.Println(TreeNode(tree, "nest", 0, 2).AsFloat64().V)
	// Output:
	// 200
	// OK
	// None
	// a.b.c
	// 1
	// 2.1
	// 8
	// 8
}

func TestTree_builtin_numbers(t *testing.T) {
	s := `{
		"meta":{"code":200, "status":"OK"},
		"response":{
			"csrf_token":"x.y.z",
			"user":{"authentication_token":"a.b.c"}
		},
		"props":[
			{"a":1},
			{"b":2.1}
		],
		"nest":[
			[3,5,8]
		]
	}`
	tree := map[string]any{}
	err := json.NewDecoder(strings.NewReader(s)).Decode(&tree)
	expect.Error(err).ToBeNil(t)
	expect.Value(TreeNode(tree, "meta", "code").AsInt()).ToBe(t, Some[int](200))
	expect.Value(TreeNode(tree, "meta", "status").AsString()).ToBe(t, Some[string]("OK"))
	expect.Value(TreeNode(tree, "meta", "status", "absent")).ToBe(t, None[any]())
	expect.Value(TreeNode(tree, "response", "user", "authentication_token").AsString()).ToBe(t, Some[string]("a.b.c"))
	expect.Value(TreeNode(tree, "props", 0, "a").AsInt()).ToBe(t, Some[int](1))
	expect.Value(TreeNode(tree, "props", 1, "b").AsFloat64()).ToBe(t, Some[float64](2.1))
	expect.Value(TreeNode(tree, "nest", 0, 2).AsInt()).ToBe(t, Some[int](8))
	expect.Value(TreeNode(tree, "nest", 0, 2).AsFloat64()).ToBe(t, Some[float64](8))
}

func TestTree_json_numbers(t *testing.T) {
	s := `{
		"meta":{"code":200, "status":"OK"},
		"response":{
			"csrf_token":"x.y.z",
			"user":{"authentication_token":"a.b.c"}
		},
		"props":[
			{"a":1},
			{"b":2.1}
		],
		"nest":[
			[3,5,8]
		]
	}`
	tree := map[string]any{}
	d := json.NewDecoder(strings.NewReader(s))
	d.UseNumber()
	err := d.Decode(&tree)
	expect.Error(err).ToBeNil(t)
	expect.Value(TreeNode(tree, "meta", "code").AsInt()).ToBe(t, Some[int](200))
	expect.Value(TreeNode(tree, "meta", "status").AsString()).ToBe(t, Some[string]("OK"))
	expect.Value(TreeNode(tree, "meta", "status", "absent")).ToBe(t, None[any]())
	expect.Value(TreeNode(tree, "response", "user", "authentication_token").AsString()).ToBe(t, Some[string]("a.b.c"))
	expect.Value(TreeNode(tree, "props", 0, "a").AsInt()).ToBe(t, Some[int](1))
	expect.Value(TreeNode(tree, "props", 1, "b").AsFloat64()).ToBe(t, Some[float64](2.1))
	expect.Value(TreeNode(tree, "nest", 0, 2).AsInt()).ToBe(t, Some[int](8))
	expect.Value(TreeNode(tree, "nest", 0, 2).AsFloat64()).ToBe(t, Some[float64](8))
}

func TestAsXxx(t *testing.T) {
	expect.Value(Some[any]("foo").AsString()).ToBe(t, Some[string]("foo"))
	expect.Value(None[any]().AsString()).ToBe(t, None[string]())

	expect.Value(Some[any](44).AsInt()).ToBe(t, Some[int](44))
	expect.Value(Some[any](json.Number("45")).AsInt()).ToBe(t, Some[int](45))
	expect.Value(Some[any](json.Number("xyz")).AsInt()).ToBe(t, None[int]())
	expect.Value(None[any]().AsInt()).ToBe(t, None[int]())

	expect.Value(Some[any](44).AsFloat64()).ToBe(t, Some[float64](44))
	expect.Value(Some[any](float64(44)).AsFloat64()).ToBe(t, Some[float64](44))
	expect.Value(Some[any](json.Number("45")).AsFloat64()).ToBe(t, Some[float64](45))
	expect.Value(Some[any](json.Number("xyz")).AsFloat64()).ToBe(t, None[float64]())
	expect.Value(None[any]().AsFloat64()).ToBe(t, None[float64]())
}
