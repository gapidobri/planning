package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"log"
)

var state = mapset.NewSet(
	on('c', 'a'),
	on('a', '1'),
	on('b', '3'),
	clear('c'),
	clear('b'),
	clear('2'),
	clear('4'),
)

func main() {
	goals := mapset.NewSet(
		on('a', 'b'),
		on('b', 'c'),
	)

	//printState(state)

	//run(goals)
	NewGoalRegression(state, goals).Run()

	//printState(state)
}

func run(goals mapset.Set[Predicate]) {
	if goals.IsSubset(state) {
		return
	}

	// 1. select goal
	unsolvedGoals := goals.Difference(state)

	goal := unsolvedGoals.ToSlice()[0]

	// 2. select best action
	actions := goal.GenerateActions()

	maxCompleted := 0
	var bestAction Action
	for action := range actions.Iter() {
		completed := action.Can().Intersect(state).Cardinality()
		if bestAction == nil || completed > maxCompleted {
			maxCompleted = completed
			bestAction = action
		}
	}
	if bestAction == nil {
		log.Fatal("no actions left")
	}

	// 3. make action possible
	run(bestAction.Can())

	// 4. execute action
	state = state.Union(bestAction.Adds())
	state = state.Difference(bestAction.Deletes())

	// 5. return to step 1. if non completed goals exist
	run(goals)
}

func printState(state mapset.Set[Predicate]) {
	world := DisplayWorld{}
	world.Apply(state)
	fmt.Println()
	world.Print()
	fmt.Println()
}

func validateState(state mapset.Set[Predicate]) {
	for predicate := range state.Iter() {
		if !state.Intersect(predicate.InvalidPredicates()).IsEmpty() {
			log.Fatalf("invalid state: %v in not possible because of: %v\n", predicate, state.Intersect(predicate.InvalidPredicates()))
		}
	}
}
