# jsontree - Simple traversal of JSON data

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](http://pkg.go.dev/github.com/rickb777/jsontree)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/jsontree)](https://goreportcard.com/report/github.com/rickb777/jsontree)
[![Build](https://github.com/rickb777/jsontree/actions/workflows/go.yml/badge.svg)](https://github.com/rickb777/jsontree/actions)
[![Coverage](https://coveralls.io/repos/github/rickb777/jsontree/badge.svg?branch=main)](https://coveralls.io/github/rickb777/jsontree?branch=main)
[![Issues](https://img.shields.io/github/issues/rickb777/jsontree.svg)](https://github.com/rickb777/jsontree/issues)

Sometimes, you don't want to unmarshal a whole JSON document to find just a few values. This API provides easy
access to items of interest within JSON documents.

JSON documents are trees containing leaves and intermediate nodes. Leaves can be

 * string - UTF-8 characters surrounded by `"` marks
 * number - a signed decimal number
 * boolean - `true` or `false`
 * null - a special case

Intermediate nodes can be

 * object containing key:value pairs; the keys are strings and the values are other nodes
 * array of other nodes that are indexed by zero-based integer

## `jsontree.TreeNode`

This function explores a JSON tree to find a given path within it. Its input is a JSON tree and a list of path items of arbitrary length.
The tree is obtained from a source document using the standard `encoding/json` API.
If there was no decoding error, then `tree` can be passed to `jsontree.TreeNode`, for example:

```Go
s := `{"a": ["ignored", {"b": 101}]}`
var tree any
err := json.NewDecoder(strings.NewReader(s)).Decode(&tree)
...
value := jsontree.TreeNode(tree, "a", 1, "b") // finds "a" then element 1 then "b"
```

The result of this function is an optional value that may contain the node wanted, which will typically be a leaf node.
Please see the examples below for more usage hints.

### Arrays are handled too

The `TreeNode` function also explores a JSON array to find a given path within it. Its input is a JSON array and
a list of path items of arbitrary length.
If there was no decoding error, then `array` can be passed to `jsontree.ArrayNode`, for example:

```Go
s := `[{"ignore": "element 0"}, {"a": {"b": 101}}]`
var array any
err := json.NewDecoder(strings.NewReader(s)).Decode(&array)
...
value := jsontree.TreeNode(array, 1, "a", "b") // finds element 1 then "a" then "b"
```

The result of this function is an optional value that may contain the node wanted, which will typically be a leaf node.
Please see the examples below for more usage hints.

## `jsontree.Option`

An `Option` may hold a leaf value or an array value that can be used as Go data. It has a generic type: the value
returned from `TreeNode` is `Option[any]` but can be converted to

 * `Option[string]` using `AsString` or `CoerceString`
 * `Option[int]` using `AsInt` or `CoerceInt`
 * `Option[float64]` using `AsFloat64` or `CoerceFloat64`
 * `Option[bool]` using `AsBool` or `CoerceBool`

There are also plural converters, e.g. `AsStrings` returns `Option[[]string]` in which the value is a slice of strings,
provided that the input was an array of strings.

The `AsXxx` methods *simply alter the type* of the result based on inspection of the value; the value is absent if it is not
of the correct type or does not have a direct conversion (e.g. number → `int`).

The `CoerceXxx` methods *convert the type* of the result to the required type (e.g. by using `strconv.ParseInt`); the resulting
value is absent if conversion failed.

The `Option.Present()` method returns true iff the value is present. If not, `Option.Err` will provide information about why not.

It is OK to call methods on absent option values; the earlier error will be propagated to the result.
