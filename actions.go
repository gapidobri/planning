package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"unicode"
)

type Action interface {
	Can() mapset.Set[Predicate]
	Adds() mapset.Set[Predicate]
	Deletes() mapset.Set[Predicate]
	CheckConstraints() bool

	String() string
}

type Move struct {
	x    rune
	from rune
	to   rune
}

func move(x, from, to rune) Action {
	return Move{
		x:    x,
		from: from,
		to:   to,
	}
}

func (m Move) Can() mapset.Set[Predicate] {
	return mapset.NewSet(
		on(m.x, m.from),
		clear(m.to),
		clear(m.x),
	)
}

func (m Move) Adds() mapset.Set[Predicate] {
	return mapset.NewSet(
		on(m.x, m.to),
		clear(m.from),
	)
}

func (m Move) Deletes() mapset.Set[Predicate] {
	return mapset.NewSet(
		clear(m.to),
		on(m.x, m.from),
	)
}

func (m Move) CheckConstraints() bool {
	return unicode.IsLetter(m.x) && m.from != m.to && m.x != m.from && m.x != m.to
}

func (m Move) String() string {
	return fmt.Sprintf("move(%c, %c, %c)", m.x, m.from, m.to)
}
