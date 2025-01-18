package runtime

import (
	"github.com/awwithro/goink/parser/types"
	"github.com/sirupsen/logrus"
)

type StoryState struct {
	globalVars     map[string]any
	CurrentChoices []Choice
	tmpVars        map[string]any
	done           bool
	visitCounts    map[*types.Container]int
	lastTurn       map[*types.Container]int
	TurnCount      int
	text           string
}

func NewStoryState() *StoryState {
	s := &StoryState{
		globalVars:     make(map[string]any),
		CurrentChoices: []Choice{},
		tmpVars:        make(map[string]any),
		visitCounts:    make(map[*types.Container]int),
		lastTurn:       make(map[*types.Container]int),
		TurnCount:      1,
	}
	return s
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
	Destination    types.Path
	OnlyDefault    bool
}

func (c Choice) ChoiceText() string {
	return c.text + c.choiceOnlyText
}
func (c Choice) storyText() string {
	return c.text
}

func (s *StoryState) RecordContainer(c *types.Container, idx int) {
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
	if len(s.CurrentChoices) > 0 && s.done {
		return false
	// } else if s.done && len(s.CurrentChoices) == 0 {
	// 	return false
	}
	return true
}

func (s *StoryState) GetText() string {
	text := s.text
	s.text = ""
	return text
}
