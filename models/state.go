package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type State struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Data      string
	IsInitial bool
	IsFinal   bool
	Adjacent  []string
}

func (s *State) SetData(data string) {
	s.Data = data
}

func (s *State) GetData() string {
	return s.Data
}

func (s *State) GetAdjacent() []string {
	return s.Adjacent
}

func (s *State) SetAdjacent(adjacent []string) {
	s.Adjacent = adjacent
}

func (s *State) SetIsInitial(initial bool) {
	s.IsInitial = initial
}

func (s *State) GetIsInitial() bool {
	return s.IsInitial
}

func (s *State) SetIsFinal(final bool) {
	s.IsFinal = final
}

func (s *State) GetIsFinal() bool {
	return s.IsFinal
}

func (s *State) ToString() string {
	brackets := ""
	if len(s.Adjacent) == 0 {
		brackets = "[ ]"
	} else {
		brackets = fmt.Sprintf("[ %v ]", s.Adjacent)
	}
	data := fmt.Sprintf("%s | %s", s.Data, brackets)
	if s.IsFinal {
		return fmt.Sprintf("((%s))", data)
	}
	return fmt.Sprintf("(%s)", data)
}
