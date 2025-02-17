package runtime

import (
	"strings"

	"github.com/awwithro/goink/pkg/parser/types"
	log "github.com/sirupsen/logrus"
)

type StoryState struct {
	globalVars     map[string]any
	currentChoices []Choice
	currentTags    []types.Tag
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
	log.Debugf("setting %s to %v", name, val)
	s.tmpVars[name] = val
}

func (s *StoryState) GetVar(name string) any {
	if v, ok := s.tmpVars[name]; ok {
		return v
	}
	log.Panicf("no var named %s", name)
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
			// A little odd, counts get recorded on entry but the visit operator
			// treats the count as prior visits. Therefore, a !ok map means we've never visited
			// a 0 means this is our first visit ...etc.
			if _, ok := s.visitCounts[c]; ok {
				s.visitCounts[c] += 1
			} else {
				s.visitCounts[c] = 0
			}
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

func (s *StoryState) GetTextAndTags() (string, []types.Tag) {
	text := CleanOutput(s.text)
	tags := s.currentTags
	s.text = ""
	s.currentTags = []types.Tag{}
	return text, tags
}

func (s *StoryState) LastTurnVisited(c *types.Container) int {
	return s.lastTurn[c]
}

func CleanOutput(str string) string {
	sb := strings.Builder{}
	currentWhitespaceStart := -1
	startOfLine := 0

	for i := 0; i < len(str); i++ {
		c := string(str[i])
		isInlineWhitespace := (c == " " || c == "\t")

		if isInlineWhitespace && currentWhitespaceStart == -1 {
			currentWhitespaceStart = i
		}
		if !isInlineWhitespace {
			if c == "\n" && i != len(str)-1 && string(str[i+1])=="\n"{
				continue
			}
			if c != "\n" && currentWhitespaceStart > 0 && currentWhitespaceStart != startOfLine {
				sb.WriteString(" ")
			}
			currentWhitespaceStart = -1
		}

		if c == "\n" {
			startOfLine = i + 1
		}
		if !isInlineWhitespace {
			sb.WriteString(c)
		}

	}

	return sb.String()
}
