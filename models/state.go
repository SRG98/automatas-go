package models

type State struct {
	Data      string
	IsInitial bool
	IsFinal   bool
	Adjacent  map[string][]string
}

func NewState(data string) *State {
	return &State{
		Data:      data,
		IsInitial: false,
		IsFinal:   false,
		Adjacent:  map[string][]string{},
	}
}

func (s *State) SetData(data string) {
	s.Data = data
}
func (s *State) GetData() string {
	return s.Data
}
func (s *State) GetIsInitial() bool {
	return s.IsInitial
}
func (s *State) SetIsInitial(initial bool) {
	s.IsInitial = initial
}
func (s *State) GetIsFinal() bool {
	return s.IsFinal
}
func (s *State) SetIsFinal(final bool) {
	s.IsFinal = final
}
func (s *State) GetAdjacent() map[string][]string {
	return s.Adjacent
}

func (s *State) SetAdjacent(adjacent map[string][]string) {
	s.Adjacent = adjacent
}

func (s *State) AddAdjacent(char []string, adjacent string) {

}

func (s *State) ToString() string {

	// brackets := "[ " + strings.Join(s.Adjacent, ", ") + " ]"
	data := s.Data /*  + " | " + brackets */
	if s.IsFinal {
		if s.IsFinal {
			return "((" + data + "))"
		}
	}
	return "(" + data + ")"

}

/*func (s *State) ToString() string {
	brackets := "[ " + strings.Join(s.Adjacent, ", ") + " ]"
	data := s.Data + " | " + brackets
	if s.IsFinal {
		return "((" + data + "))"
	}
	return "(" + data + ")"
}*/
