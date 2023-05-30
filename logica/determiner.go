package logic

import (
	"strconv"

	"github.com/SRG98/automatas-go/models"
)

type Determiner struct {
	Automata  *models.Automata
	adjacent  []string
	autoTable []map[string][]string
}

func NewDeterminer() *Determiner {
	return &Determiner{
		Automata:  nil,
		adjacent:  []string{},
		autoTable: []map[string][]string{},
	}
}

func (d *Determiner) SetAutomata(automata *models.Automata) {
	d.Automata = automata
}

func (d *Determiner) Determine() *models.Automata {
	states := d.Automata.GetStates()
	d.setRealAdjacent(states)
	d.autoTable = d.createTable()

	auto := d.formatAutomata(d.autoTable)
	return auto
}

/* ----------------------------------------------------------------------------------------- */

func (d *Determiner) setRealAdjacent(states []*models.State) {
	language := d.Automata.GetLanguage()
	for _, state := range states {
		d.stateTuples(state, language)
	}
}

func (d *Determiner) stateTuples(state *models.State, language []string) {
	tCharSet := make(map[string][]string)
	for _, symbol := range language {
		d.adjacent = []string{}
		d.adjacentByChar(state, symbol, true)
		tCharSet[symbol] = d.adjacent
	}
	d.Automata.AddAdjacent(state.GetData(), tCharSet)
}

func (d *Determiner) adjacentByChar(state *models.State, tChar string, hasTChar bool) {
	lamWay := false
	transitions := d.selectTrans(state)

	if len(transitions) == 0 {
		d.adjacent = append(d.adjacent, state.GetData())
		return
	}

	for _, tran := range transitions {
		if contains(tran.GetChars(), "_") {
			lamWay = true
		}
	}

	// hadOneTChar := false
	// for _, tran := range transitions {
	// 	if contains(tran.GetChars(), tChar) {
	// 		// hadOneTChar = true
	// 	}
	// }

	for _, tran := range transitions {
		if contains(tran.GetChars(), "_") {
			nextState := d.Automata.GetState(tran.GetEnd())
			d.adjacentByChar(nextState, tChar, hasTChar)
		} else if contains(tran.GetChars(), tChar) {
			if hasTChar {
				nextState := d.Automata.GetState(tran.GetEnd())
				d.adjacentByChar(nextState, tChar, false)
			} else {
				if !contains(d.adjacent, state.GetData()) {
					d.adjacent = append(d.adjacent, state.GetData())
				}
				return
			}
		} else {
			if lamWay {
				return
			}

			if !contains(d.adjacent, state.GetData()) {
				d.adjacent = append(d.adjacent, state.GetData())
			}
			return
		}
	}
	// return
}

func (d *Determiner) selectTrans(state *models.State) []*models.Transition {
	transitions := []*models.Transition{}
	for _, transition := range d.Automata.GetTransitions() {
		if transition.GetStart() == state.GetData() {
			transitions = append(transitions, transition)
		}
	}
	return transitions
}

// Util function
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

/* ----------------------------------------------------------------------------------------- */

func (d *Determiner) createTable() []map[string][]string {
	table := []map[string][]string{}
	states := d.Automata.GetStates()
	for _, state := range states {
		newRow := d.newRowTable(state)
		table = append(table, newRow)
	}
	return d.mergeSubsets(table)
}

func (d *Determiner) newRowTable(state *models.State) map[string][]string {
	row := make(map[string][]string)
	row["origin"] = []string{state.GetData()}
	transitions := d.selectTrans(state)
	for _, transition := range transitions {
		for _, char := range transition.GetChars() {
			row[char] = append(row[char], transition.GetEnd())
		}
	}
	return row
}
func (d *Determiner) mergeSubsets(table []map[string][]string) []map[string][]string {
	newTable := []map[string][]string{}
	for _, row := range table {
		found := false
		for _, nRow := range newTable {
			if d.sameElements(row["origin"], nRow["origin"]) {
				found = true
				break
			}
		}
		if !found {
			newTable = append(newTable, row)
		}
	}
	return newTable
}

func (d *Determiner) colToRow(table []map[string][]string) map[string][]string {
	row := make(map[string][]string)
	for _, symbol := range d.Automata.GetLanguage() {
		for _, subset := range table {
			row[symbol] = append(row[symbol], subset[symbol]...)
		}
		row[symbol] = d.uniqueElements(row[symbol])
	}
	return row
}

func (d *Determiner) sameElements(list1, list2 []string) bool {
	if len(list1) != len(list2) {
		return false
	}
	for _, val := range list1 {
		if !contains(list2, val) {
			return false
		}
	}
	return true
}

func (d *Determiner) uniqueElements(list []string) []string {
	u := make([]string, 0, len(list))
	m := make(map[string]bool)

	for _, val := range list {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

/* ----------------------------------------------------------------------------------------- */

func (d *Determiner) formatAutomata(table []map[string][]string) *models.Automata {
	newAutomata := models.NewAutomaton()

	// Crear registro de datos de estado
	stateDataRegistry := make(map[string][]string)
	for i, row := range table {
		index := "Q" + strconv.Itoa(i)
		stateDataRegistry[index] = row["origin"]
	}

	// Crear estados
	for stateIndex, statesData := range stateDataRegistry {
		isInitial, isFinal := d.getStatesAtt(statesData)
		newAutomata.NewState(stateIndex, isInitial, isFinal)
	}

	// Crear transiciones
	for i, row := range table {
		stateIndex := "Q" + strconv.Itoa(i)
		for symbol, destinations := range row {
			if symbol != "origin" {
				for _, destination := range destinations {
					destStateIndex := d.stateIndexByData(stateDataRegistry, destination)
					newAutomata.NewTransition(stateIndex, destStateIndex, []string{symbol})
				}
			}
		}
	}
	return newAutomata
}

func (d *Determiner) getStatesAtt(statesData []string) (bool, bool) {
	isInitial := true
	isFinal := false
	for _, state := range d.Automata.GetStates() {
		if contains(statesData, state.GetData()) {
			isFinal = isFinal || state.GetIsFinal()
			isInitial = isInitial && state.GetIsInitial()
		}
	}
	return isInitial, isFinal
}

func (d *Determiner) stateIndexByData(stateDataRegistry map[string][]string, destination string) string {
	for stateIndex, statesData := range stateDataRegistry {
		if contains(statesData, destination) {
			return stateIndex
		}
	}
	return ""
}

/* ----------------------------------------------------------------------------------------- */

// func (d *Determiner) contains(slice []string, item string) bool {
// 	for _, a := range slice {
// 		if a == item {
// 			return true
// 		}
// 	}
// 	return false
// }
