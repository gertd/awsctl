package shared

// Copyright Jon Bodner (https://github.com/jonbodner)
//
// Adopted from:
// Closures are the Generics for Go
// Medium, Jun 7, 2017
// https://medium.com/capital-one-developers/closures-are-the-generics-for-go-cb32021fb5b5

import "sort"

type sorter struct {
	len  int
	swap func(i, j int)
	less func(i, j int) bool
}

func (x sorter) Len() int           { return x.len }
func (x sorter) Swap(i, j int)      { x.swap(i, j) }
func (x sorter) Less(i, j int) bool { return x.less(i, j) }

// Sort --
func Sort(n int, swap func(i, j int), less func(i, j int) bool) {
	sort.Sort(sorter{len: n, swap: swap, less: less})
}
