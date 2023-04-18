package models

type Function struct {
	automata    *Automata
	inputString string
	index       int
	actualState *State
}

// func NewFunction(automata *Automata) *Function {
// 	return &Function{
// 		automata:    automata,
// 		inputString: "",
// 		index:       0,
// 		actualState: automata.GetInitialState(),
// 	}
// }

// func (f *Function) SetString(inputString string) {
// 	f.inputString = inputString
// }

// func (f *Function) GetInitialState() *State {
// 	for _, state := range f.automata.states {
// 		if state.GetIsStart() {
// 			return state
// 		}
// 	}
// 	return nil
// }

// func (f *Function) Validate(inputString string) {
// 	f.actualState = f.automata.states[0]
// 	acceptable := false
// 	for i := 0; i < len(inputString); i++ {
// 		if f.actualState.GetIsEnd() {
// 			acceptable = f.GetNextState()
// 		}
// 	}
// }

// func (f *Function) GetNextState() bool {
// 	if f.index == len(f.inputString) {
// 		fmt.Println("All string has been read.")
// 		return true
// 	}
// 	flag := true

// 	transitions := f.automata.GetTransitions()
// 	for _, transition := range transitions {
// 		if f.actualState.GetData() == transition.GetStart() && flag {
// 			for _, char := range transition.GetChars() {
// 				if char == string(f.inputString[f.index]) && flag {
// 					if contains(f.automata.GetAlphabet(), string(f.inputString[f.index])) {
// 						f.actualState = f.automata.GetState(transition.GetEnd())
// 						fmt.Printf("Valid char, state to be: %s\n", f.actualState)
// 						f.index++
// 						flag = false
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return false
// }

// func (f *Function) SetActualState(actualState *State) {
// 	f.actualState = actualState
// }

// func (f *Function) GetActualState() *State {
// 	return f.actualState
// }
