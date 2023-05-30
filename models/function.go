package models

import (
	"fmt"
)

type Function struct {
	automata    *Automata
	inputString string
	actualState *State
}

func NewFunction(automata *Automata) *Function {
	return &Function{
		automata:    automata,
		inputString: "",
		actualState: nil,
	}
}

func (f *Function) SetAutomata(automata *Automata) {
	f.automata = automata
}

func (f *Function) SetString(inputString string) {
	f.inputString = inputString
}

func (f *Function) GetInitialState() *State {
	for _, state := range f.automata.GetStates() {
		if state.GetIsInitial() {
			return state
		}
	}
	return nil
}

func (f *Function) hasOneInitialState() bool {
	initialStates := 0
	for _, state := range f.automata.GetStates() {
		if state.GetIsInitial() {
			initialStates++
		}
	}
	return initialStates == 1
}

func (f *Function) hasFinalStates() bool {
	for _, state := range f.automata.GetStates() {
		if state.GetIsFinal() {
			return true
		}
	}
	return false
}

func (f *Function) canTransition(char string) bool {
	transitions := f.automata.GetTransitions()
	counter := 0
	nextState := &State{}

	for _, transition := range transitions {
		if f.actualState.GetData() == transition.GetStart() {
			for _, tChar := range transition.GetChars() {
				if tChar == char {
					nextState = f.automata.GetState(transition.GetEnd())
					counter++
				}
			}
		}
	}

	if counter == 1 {
		f.actualState = nextState
		return true
	}

	return false
}

func (f *Function) travel() bool {
	if f.inputString == "" {
		fmt.Println("La cadena de entrada está vacía.")
		return f.actualState.GetIsFinal()
	}

	idx := 0
	for idx < len(f.inputString) {
		char := string(f.inputString[idx])
		if !f.canTransition(char) {
			return false
		}
		idx++
	}
	return true
}

func (f *Function) Validate(inputString string) bool {
	if !f.hasOneInitialState() || !f.hasFinalStates() {
		return false
	}

	f.SetString(inputString)
	f.actualState = f.GetInitialState()

	if !f.travel() {
		return false
	}

	return f.actualState.GetIsFinal()
}
