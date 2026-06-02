package jsontree

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/rickb777/expect"
)

func ExampleTreeNode() {
	// some JSON
	s := `{
		"meta":{"code":200, "status":"OK"},
		"response":{
			"csrf_token":"x.y.z",
			"user":{"auth_token":"a.b.c"}
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
	fmt.Println("meta,code               :", TreeNode(tree, "meta", "code").AsInt().V)
	fmt.Println("meta,status             :", TreeNode(tree, "meta", "status").AsString().V)
	fmt.Println("meta,status,absent      :", TreeNode(tree, "meta", "status", "absent")) // value is absent
	fmt.Println("response,user,auth_token:", TreeNode(tree, "response", "user", "auth_token").AsString().V)
	fmt.Println("props,0,a               :", TreeNode(tree, "props", 0, "a").AsInt().V)
	fmt.Println("props,1,b               :", TreeNode(tree, "props", 1, "b").AsFloat64().V)
	fmt.Println("nest,0,1                :", TreeNode(tree, "nest", 0, 1).AsInt().V)
	fmt.Println("nest,0,2                :", TreeNode(tree, "nest", 0, 2).AsFloat64().V)
	fmt.Println("nest,0                  :", TreeNode(tree, "nest", 0).AsInts().V)
	// Output:
	// meta,code               : 200
	// meta,status             : OK
	// meta,status,absent      : None(meta,status,absent not found)
	// response,user,auth_token: a.b.c
	// props,0,a               : 1
	// props,1,b               : 2.1
	// nest,0,1                : 5
	// nest,0,2                : 8
	// nest,0                  : [3 5 8]
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
	if err != nil {
		panic(err)
	}

	expect.Value(TreeNode(tree, "meta", "code").AsInt()).ToBe(t, Some[int](200))
	expect.Value(TreeNode(tree, "meta", "status").AsString()).ToBe(t, Some[string]("OK"))
	expect.Value(TreeNode(tree, "meta", "status", "absent")).ToBe(t,
		Option[any]{Err: errors.New("meta,status,absent not found")})
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

	d := json.NewDecoder(strings.NewReader(s))
	d.UseNumber()

	tree := map[string]any{}
	err := d.Decode(&tree)
	if err != nil {
		panic(err)
	}

	expect.Value(TreeNode(tree, "meta", "code").AsInt()).ToBe(t, Some[int](200))
	expect.Value(TreeNode(tree, "meta", "status").AsString()).ToBe(t, Some[string]("OK"))
	expect.Value(TreeNode(tree, "meta", "status", "absent")).ToBe(t,
		Option[any]{Err: errors.New("meta,status,absent not found")})
	expect.Value(TreeNode(tree, "response", "user", "authentication_token").AsString()).ToBe(t,
		Some[string]("a.b.c"))
	expect.Value(TreeNode(tree, "props", 0, "a").AsInt()).ToBe(t, Some[int](1))
	expect.Value(TreeNode(tree, "props", 1, "b").AsFloat64()).ToBe(t, Some[float64](2.1))
	expect.Value(TreeNode(tree, "nest", 0, 2).AsInt()).ToBe(t, Some[int](8))
	expect.Value(TreeNode(tree, "nest", 0, 2).AsFloat64()).ToBe(t, Some[float64](8))
}

func TestSlices(t *testing.T) {
	s := `{
		"full": {
			"0": [],
			"1": [],
			"2": [],
			"3": [],
			"4": [],
			"5": [],
			"6": [],
			"7": [],
			"8": [],
			"9": []
		},
		"line": {
			"0": ["1","2","3","4","5","6","7","8","9"],
			"1": ["2","3","4","5","6","7","8","9"],
			"2": ["3","4","5","6","7","8","9"],
			"3": ["4","5","6","7","8","9"],
			"4": [5,6,7,8,9],
			"5": [6,7,8,9],
			"6": [7,8,9],
			"7": [8,9],
			"8": [9],
			"9": []
		}
	}`

	d := json.NewDecoder(strings.NewReader(s))
	d.UseNumber()

	tree := map[string]any{}
	err := d.Decode(&tree)
	if err != nil {
		panic(err)
	}

	expect.Value(TreeNode(tree, "line", "5", "2").AsInt()).ToBe(t, Some[int](8))
	expect.Value(TreeNode(tree, "line", "1", "2").AsString()).ToBe(t, Some[string]("4"))
	expect.Value(TreeNode(tree, "line", "1", "2").AsInt()).ToBe(t, Option[int]{Err: errors.New("string is not an int")})
	expect.Value(TreeNode(tree, "line", "1", "2").CoerceInt()).ToBe(t, Some[int](4))
	expect.Value(TreeNode(tree, "line", "1", "2").CoerceFloat64()).ToBe(t, Some[float64](4))
	expect.Value(TreeNode(tree, "line", "7").AsInts()).ToBe(t, Some[[]int]([]int{8, 9}))
	expect.Value(TreeNode(tree, "line", "3").CoerceInts()).ToBe(t, Some[[]int]([]int{4, 5, 6, 7, 8, 9}))
	expect.Value(TreeNode(tree, "line", "3").AsInts()).ToBe(t, Option[[]int]{Err: errors.New("[]interface {} is not []int")})
	expect.Value(TreeNode(tree, "line", "7").AsFloat64s()).ToBe(t, Some[[]float64]([]float64{8, 9}))
	expect.Value(TreeNode(tree, "line", "3").CoerceFloat64s()).ToBe(t, Some[[]float64]([]float64{4, 5, 6, 7, 8, 9}))
	expect.Value(TreeNode(tree, "line", "3").AsFloat64s()).ToBe(t, Option[[]float64]{Err: errors.New("[]interface {} is not []float64")})
}
