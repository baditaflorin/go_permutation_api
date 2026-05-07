//go:build js && wasm

package main

import (
	"sort"
	"syscall/js"
)

func main() {
	js.Global().Set("goPermutations", js.FuncOf(generate))
	// Keep the Go runtime alive
	select {}
}

// generate is called from JavaScript as:
//
//	goPermutations(["a","b","c"])  → [["a","b","c"], ...]
func generate(_ js.Value, args []js.Value) any {
	if len(args) == 0 {
		return errorResult("no arguments provided")
	}

	jsArr := args[0]
	if jsArr.Type() != js.TypeObject || jsArr.Get("length").IsUndefined() {
		return errorResult("argument must be an array")
	}

	n := jsArr.Get("length").Int()
	if n == 0 {
		return errorResult("array must not be empty")
	}
	if n > 12 {
		return errorResult("too many elements (max 12)")
	}

	elements := make([]string, n)
	for i := range n {
		elements[i] = jsArr.Index(i).String()
	}

	perms := generateAll(elements)

	// Convert [][]string → JS array of arrays
	outer := js.Global().Get("Array").New(len(perms))
	for i, perm := range perms {
		inner := js.Global().Get("Array").New(len(perm))
		for j, s := range perm {
			inner.SetIndex(j, s)
		}
		outer.SetIndex(i, inner)
	}
	return outer
}

func errorResult(msg string) js.Value {
	err := js.Global().Get("Error").New(msg)
	panic(err) // surfaces as a thrown JS Error
}

func generateAll(elements []string) [][]string {
	copied := make([]string, len(elements))
	copy(copied, elements)
	sort.Strings(copied)

	var results [][]string
	results = append(results, clone(copied))

	for nextPermutation(copied) {
		results = append(results, clone(copied))
	}
	return results
}

func clone(s []string) []string {
	c := make([]string, len(s))
	copy(c, s)
	return c
}

func nextPermutation(arr []string) bool {
	n := len(arr)
	i := n - 2
	for i >= 0 && arr[i] >= arr[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := n - 1
	for arr[j] <= arr[i] {
		j--
	}
	arr[i], arr[j] = arr[j], arr[i]
	for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
		arr[l], arr[r] = arr[r], arr[l]
	}
	return true
}
