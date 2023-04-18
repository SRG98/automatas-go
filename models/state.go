package models

import (
	"strings"
)

type State struct {
	Data      string
	IsInitial bool
	IsFinal   bool
	Adjacent  []string
}

func NewState(data string) *State {
	return &State{
		Data:      data,
		IsInitial: false,
		IsFinal:   false,
		Adjacent:  []string{},
	}
}

func (s *State) SetData(data string) {
	s.Data = data
}
func (s *State) GetData() string {
	return s.Data
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
func (s *State) GetAdjacent() []string {
	return s.Adjacent
}

func (s *State) SetAdjacent(adjacent []string) {
	s.Adjacent = adjacent
}

func (s *State) ToString() string {
	brackets := "[ " + strings.Join(s.Adjacent, ", ") + " ]"
	data := s.Data + " | " + brackets
	if s.IsFinal {
		return "((" + data + "))"
	}
	return "(" + data + ")"
}
