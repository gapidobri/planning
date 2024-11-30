package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
)

type DisplayWorld map[int][]rune

func (w DisplayWorld) Apply(p mapset.Set[Predicate]) {
	applied := mapset.NewSet[Predicate]()

	for {
		for predicate := range p.Iter() {
			if applied.Contains(predicate) {
				continue
			}
			if predicate.Apply(w) {
				applied.Add(predicate)
			}
		}
		if applied.Equal(p) {
			break
		}
	}
}

func (w DisplayWorld) Print() {
	height := 0
	for _, col := range w {
		if len(col) > height {
			height = len(col)
		}
	}

	for i := height - 1; i >= 0; i-- {
		for j := range len(w) {
			col, ok := w[j+1]
			if !ok {
				continue
			}
			if len(col)-1 < i {
				fmt.Print("  ")
			} else {
				fmt.Printf(" %c", col[i])
			}
		}
		fmt.Println()
	}
	for i := range len(w) {
		fmt.Printf(" %d", i+1)
	}
	fmt.Println()
}
