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
			// fmt.Println("intial: ", state.ToString())
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
	// fmt.Println("Sólo un estado incial? ", initialStates == 1)
	return initialStates == 1
}

func (f *Function) hasFinalStates() bool {
	for _, state := range f.automata.GetStates() {
		if state.GetIsFinal() {
			return true
		}
	}
	// fmt.Println("El autómata debe tener al menos un estado final.")
	return false
}

func (f *Function) canTransition(char string) bool {
	transitions := f.automata.GetTransitions()
	counter := 0
	nextState := &State{}

	for _, transition := range transitions {
		// fmt.Println("actualState:", f.actualState.ToString(), "| trans:", transition.ToString())
		if f.actualState.GetData() == transition.GetStart() {
			for _, tChar := range transition.GetChars() {
				// fmt.Println("to: ", char, "| iterate: ", tChar)
				if tChar == char {
					nextState = f.automata.GetState(transition.GetEnd())
					counter++
					// fmt.Println("++ | newActualState: ", nextState.ToString())
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
		// fmt.Println("SENDED CHAR:", char)
		if !f.canTransition(char) {
			return false
		}
		idx++
		// fmt.Printf("travel index: %v | Char %v can transition.\n", idx, char)
	}
	return true
}

func (f *Function) Validate(inputString string) bool {
	if !f.hasOneInitialState() || !f.hasFinalStates() {
		return false
	}

	// fmt.Println("string to use: ", inputString)
	f.SetString(inputString)
	f.actualState = f.GetInitialState()

	if !f.travel() {
		return false
	}

	return f.actualState.GetIsFinal()
}
