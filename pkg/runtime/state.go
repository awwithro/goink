package runtime

import (
	"github.com/awwithro/goink/pkg/parser/types"
	"github.com/sirupsen/logrus"
)

type StoryState struct {
	globalVars     map[string]any
	currentChoices []Choice
	tmpVars        map[string]any
	done           bool
	Finished       bool
	visitCounts    map[*types.Container]int
	lastTurn       map[*types.Container]int
	TurnCount      int
	text           string
}

func NewStoryState() *StoryState {
	s := &StoryState{
		globalVars:     make(map[string]any),
		currentChoices: []Choice{},
		tmpVars:        make(map[string]any),
		visitCounts:    make(map[*types.Container]int),
		lastTurn:       make(map[*types.Container]int),
		TurnCount:      1,
	}
	return s
}

func (s *StoryState) GetChoices() []Choice {
	fallback := []Choice{}
	choices := []Choice{}
	includeFallBack := true
	for _, choice := range s.currentChoices {
		if choice.OnlyDefault {
			fallback = append(fallback, choice)
		} else {
			choices = append(choices, choice)
			includeFallBack = false
		}
	}
	if includeFallBack {
		choices = append(choices, fallback...)
	}
	return choices
}

func (s *StoryState) SetVar(name string, val any) {
	s.tmpVars[name] = val
}

func (s *StoryState) GetVar(name string) any {
	if v, ok := s.tmpVars[name]; ok {
		return v
	}
	logrus.Panicf("no var named %s", name)
	return nil
}

type Choice struct {
	text           string
	choiceOnlyText string
	Destination    Address
	OnlyDefault    bool
}

func (c Choice) ChoiceText() string {
	return c.text + c.choiceOnlyText
}
func (c Choice) storyText() string {
	return c.text
}

func (s *StoryState) RecordContainer(a Address) {
	c := a.C
	idx := a.I
	if c.RecordVisits() {
		if !c.CountStartOnly() || idx == 0 {
			s.visitCounts[c] += 1
		}
	}
	if c.RecordTurns() {
		s.lastTurn[c] = s.TurnCount
	}
}

func (s *StoryState) setDone(x bool) {
	s.done = x
}

func (s *StoryState) CanContinue() bool {
	if len(s.currentChoices) > 0 && s.done || s.Finished {
		return false
	}
	return true
}

func (s *StoryState) GetText() string {
	text := s.text
	s.text = ""
	return text
}
