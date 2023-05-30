package models

import (
	"strings"
)

type Transition struct {
	Start string
	End   string
	Chars []string
}

func NewTransition(start string, end string, chars []string) *Transition {
	return &Transition{
		Start: start,
		End:   end,
		Chars: chars,
	}
}

func (t *Transition) SetStart(start string) {
	t.Start = start
}

func (t *Transition) GetStart() string {
	return t.Start
}

func (t *Transition) SetEnd(end string) {
	t.End = end
}

func (t *Transition) GetEnd() string {
	return t.End
}

func (t *Transition) NewChars(data string) {
	t.Chars = append(t.Chars, data)
}

func (t *Transition) SetChars(data []string) {
	t.Chars = data
}

func (t *Transition) GetChars() []string {
	return t.Chars
}

func (t *Transition) ToString() string {
	return "(" + t.Start + " → " + t.End + " ⟨" + strings.Join(t.Chars, ", ") + "⟩)"
}
