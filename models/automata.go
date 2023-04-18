package models

import (
	"fmt"
	"strings"
)

type Automata struct {
	Name        string
	States      []*State
	Transitions []*Transition
	Alphabet    []string
}

func NewAutomaton() *Automata {
	return &Automata{
		Name:        "",
		States:      []*State{},
		Transitions: []*Transition{},
		Alphabet:    []string{},
	}
}

func (a *Automata) SetName(name string) {
	a.Name = name
}

func (a *Automata) GetName() string {
	return a.Name
}

func (a *Automata) SetStates(states []*State) {
	a.States = states
}

func (a *Automata) GetStates() []*State {
	return a.States
}

func (a *Automata) GetState(data string) *State {
	for _, state := range a.States {
		if state.GetData() == data {
			return state
		}
	}
	return nil
}

func (a *Automata) SetTransitions(transitions []*Transition) {
	a.Transitions = transitions
}

func (a *Automata) GetTransitions() []*Transition {
	return a.Transitions
}

func (a *Automata) SetAlphabet(symbols []string) {
	a.Alphabet = symbols
}

func (a *Automata) GetAlphabet() []string {
	return a.Alphabet
}

func (a *Automata) ExistObj(obj interface{}, list interface{}) bool {
	switch list.(type) {
	case []*State:
		states := list.([]*State)
		for _, state := range states {
			if state.GetData() == obj.(*State).GetData() {
				return true
			}
		}
	case []*Transition:
		transitions := list.([]*Transition)
		for _, transition := range transitions {
			if transition.GetStart() == obj.(*Transition).GetStart() &&
				transition.GetEnd() == obj.(*Transition).GetEnd() {
				return true
			}
		}
	}
	return false
}

func (a *Automata) ExistState(data string) bool {
	for _, state := range a.States {
		if state.GetData() == data {
			return true
		}
	}
	return false
}

func (a *Automata) NewState(data string, final bool, initial bool) {
	newState := NewState(data)
	newState.SetIsFinal(final)
	newState.SetIsInitial(initial)

	if a.ExistObj(newState, a.States) {
		fmt.Println("State already exists.")
		return
	}
	a.States = append(a.States, newState)
	fmt.Println("New state added.")
}

func (a *Automata) NewTransition(start, end string, data []string) {
	newTransition := NewTransition(start, end, data)

	if a.ExistObj(newTransition, a.Transitions) {
		fmt.Println("Transition already exists.")
		return
	}

	if a.ExistState(start) && a.ExistState(end) {
		validChars := true
		for _, char := range data {
			if !contains(a.Alphabet, char) {
				validChars = false
				break
			}
		}
		if validChars {
			state := a.GetState(start)
			state.SetAdjacent(append(state.GetAdjacent(), end))
			a.Transitions = append(a.Transitions, newTransition)

			fmt.Println("New transition added.")
			return
		}
	}
	fmt.Println("No transition added.")
}

func (a *Automata) SeeStates() string {
	var result strings.Builder
	for _, state := range a.States {
		result.WriteString(state.ToString() + " ")
	}
	return result.String()
}

func (a *Automata) SeeTransitions() string {
	var result strings.Builder
	for _, transition := range a.Transitions {
		result.WriteString(transition.ToString() + " ")
	}
	return result.String()
}

func (a *Automata) ToString() string {
	return fmt.Sprintf(
		"\nName: %s\nStates: %s\nTransitions: %s\nAlphabet: %v",
		a.GetName(), a.SeeStates(), a.SeeTransitions(), a.Alphabet,
	)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
