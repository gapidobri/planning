package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"log"
	"strconv"
	"unicode"
)

type Predicate interface {
	Apply(w DisplayWorld) bool

	GenerateActions() mapset.Set[Action]
	InvalidPredicates() mapset.Set[Predicate]

	String() string
}

type On struct {
	top    rune
	bottom rune
}

func on(top, bottom rune) Predicate {
	return On{
		top:    top,
		bottom: bottom,
	}
}

func (o On) Apply(w DisplayWorld) bool {
	pos, err := strconv.Atoi(string(o.bottom))
	if err != nil {
		// boxes
		for p, boxes := range w {
			if len(boxes) == 0 {
				return false
			}
			if boxes[len(boxes)-1] == o.bottom {
				w[p] = append(w[p], o.top)
				return true
			}
		}
		return false
	}

	// numbers
	if len(w[pos]) > 0 {
		log.Fatalf("Block %c already on position %d\n", w[pos][0], pos)
	}
	w[pos] = append(w[pos], o.top)

	return true
}

func (o On) GenerateActions() mapset.Set[Action] {
	options := "1234abc"

	actions := mapset.NewSet[Action]()
	for _, from := range options {
		action := move(o.top, from, o.bottom)
		if action.CheckConstraints() {
			actions.Add(action)
		}
	}

	return actions
}

func (o On) InvalidPredicates() mapset.Set[Predicate] {
	options := "1234abc"
	predicates := mapset.NewSet[Predicate]()
	for _, top := range options {
		if top == o.top {
			continue
		}
		predicates.Add(on(top, o.bottom))
	}
	return predicates
}

func (o On) String() string {
	return fmt.Sprintf("on(%c, %c)", o.top, o.bottom)
}

type Clear struct {
	position rune
}

func clear(position rune) Predicate {
	return Clear{
		position: position,
	}
}

func (c Clear) Apply(w DisplayWorld) bool {
	if unicode.IsNumber(c.position) {
		col, err := strconv.Atoi(string(c.position))
		if err != nil {
			log.Fatalf("Error converting position to integer: %v", err)
		}
		w[col] = []rune{}
	}
	return true
}

func (c Clear) GenerateActions() mapset.Set[Action] {
	options := "1234abc"

	actions := mapset.NewSet[Action]()
	for _, x := range options {
		for _, to := range options {
			action := move(x, c.position, to)
			if action.CheckConstraints() {
				actions.Add(action)
			}
		}
	}

	return actions
}

func (c Clear) InvalidPredicates() mapset.Set[Predicate] {
	options := "1234abc"
	predicates := mapset.NewSet[Predicate]()
	for _, top := range options {
		predicates.Add(on(top, c.position))
	}
	return predicates
}

func (c Clear) String() string {
	return fmt.Sprintf("clear(%c)", c.position)
}
