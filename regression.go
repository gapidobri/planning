package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"log"
)

type Goals mapset.Set[Predicate]

type GoalRegression struct {
	state   mapset.Set[Predicate]
	goals   Goals
	actions []Action
}

func NewGoalRegression(state mapset.Set[Predicate], goals Goals) *GoalRegression {
	return &GoalRegression{
		state: state,
		goals: goals,
	}
}

func (gr *GoalRegression) Run() {
	goals := gr.Goals(3)
	fmt.Println(goals)
}

func (gr *GoalRegression) Goals(i int) Goals {
	if i == 0 {
		return gr.goals
	}

	goals := gr.Goals(i - 1)
	if goals.IsSubset(gr.state) {
		return goals
	}

	for goal := range goals.Iter() {
		fmt.Printf("selected goal: %v\n", goal)

		actions := goal.GenerateActions()

		var bestAction Action
		for action := range actions.Iter() {
			if action.Deletes().Intersect(goals).IsEmpty() {
				bestAction = action
				break
			}
		}
		if bestAction == nil {
			log.Fatal("no actions left")
		}
		fmt.Printf("action: %v\n", bestAction)

		newGoals := goals.Union(bestAction.Can()).Difference(bestAction.Adds())

		if gr.ValidateGoals(newGoals) {
			return newGoals
		}
	}

	log.Fatal("no goals left")
	return nil
}

func (gr *GoalRegression) ValidateGoals(goals Goals) bool {
	for predicate := range goals.Iter() {
		if !goals.Intersect(predicate.InvalidPredicates()).IsEmpty() {
			return false
		}
	}
	return true
}
