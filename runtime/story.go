package runtime

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/awwithro/goink/parser/types"
	"github.com/emirpasic/gods/v2/stacks"
	"github.com/emirpasic/gods/v2/stacks/arraystack"
	log "github.com/sirupsen/logrus"
)

// Index of a container
type Address struct {
	C *types.Container
	I int
}

func (a Address) AtEnd() bool {
	return a.I >= len(a.C.Contents)
}

func (a *Address) Set(c *types.Container, i int) {
	a.C = c
	a.I = i
}
func (a *Address) Increment() {
	a.I++
}

type Story struct {
	ink             types.Ink
	evaluationStack stacks.Stack[any]
	outputBuffer    stacks.Stack[string]
	mode            Mode
	stringMarker    int //Used to track index of stack to concatenate into a string
	state           *StoryState
	// Where in the ink we're located
	currentAddress Address
	Finished       bool
}

func NewStory(ink types.Ink) Story {
	s := Story{
		ink:             ink,
		evaluationStack: arraystack.New[any](),
		outputBuffer:    arraystack.New[string](),
		mode:            None,
		stringMarker:    -1,
		state:           NewStoryState(),
	}
	return s
}

func (s *Story) startEvalMode() {
	if s.mode != None {
		panicInvalidModeTransition(s.mode, Eval)
	}
	s.mode = Eval
}

func (s *Story) endEvalMode() {
	if s.mode != Eval {
		panicInvalidModeTransition(s.mode, None)
	}
	s.mode = None
}

func (s *Story) startStrMode() {
	if s.mode != Eval {
		panicInvalidModeTransition(s.mode, Str)
	}
	s.mode = Str
	s.stringMarker = s.outputBuffer.Size()
}
func (s *Story) popOutput() {
	val := mustPopStack(s.evaluationStack)
	if str, ok := val.(fmt.Stringer); !ok {
		panicInvalidStackType(str, val)
	} else {
		s.outputBuffer.Push(str.String())
	}
}

func (s *Story) endStrMode() {
	if s.mode != Str {
		panicInvalidModeTransition(s.mode, Str)
	}
	s.mode = Eval
	if s.outputBuffer.Size() == s.stringMarker {
		log.Panic("No elements could be popped from the output buffer")
	}
	result := strings.Builder{}
	items := s.outputBuffer.Size() - s.stringMarker
	for items > 0 {
		val, _ := s.outputBuffer.Pop()
		result.WriteString(val)
		items--
	}
	s.evaluationStack.Push(types.StringVal(result.String()))
	s.stringMarker = -1
}

func panicInvalidModeTransition(current, attempted Mode) {
	log.Panicf("Invalid Mode transition. Can't go from %s to %s", current, attempted)
}

func panicInvalidStackType(expected any, actual any) {
	log.Panicf("non %s in stack: %s, %v", reflect.TypeOf(expected), reflect.TypeOf(actual), actual)
}

func (s *Story) Start() {
	s.enterContainer(s.ink.Root.Contents[0].(*types.Container), 0)
}

func mustPopStack(s stacks.Stack[any]) any {
	val, ok := s.Pop()
	if !ok {
		log.Panic("Popped empty stack")
	}
	log.Debug("Popped ", val)
	return val
}

func mustPopNumeric(s stacks.Stack[any]) types.NumericVal {
	x := mustPopStack(s)
	num, ok := x.(types.NumericVal)
	if !ok {
		panicInvalidStackType(num, x)
	}
	return num
}

func (s *Story) ResolvePath(p types.Path) (*types.Container, int) {
	return types.ResolvePath(p, s.currentAddress.C)
}

func (s *Story) Step() (StoryState, error) {
	if s.state.CanContinue() {
		// flush any already presented text
		s.state.text = ""
		s.reEnterStory()
		// if a choice is needed after taking a step, send text to the state
		if !s.state.CanContinue() {
			s.writeToState()
		}
		// } else if s.state.done && len(s.state.CurrentChoices) == 0 {
		// 	s.endStory()
	} else {
		return *s.state, fmt.Errorf("can't continue")
	}
	return *s.state, nil
}

func (s *Story) reEnterStory() {
	s.state.setDone(false)
	if s.currentAddress.AtEnd() {
		log.Debug("Reached end of Container: ", s.currentAddress.C.Name)
		// End of the story
		if pos, err := s.currentAddress.C.PositionInParent(); err != nil {
			log.Debug("reached end of ink ", err)
			s.endStory()
		} else {
			// pick up at the position just after the container we left
			s.currentAddress.Set(s.currentAddress.C.ParentContainer, pos+1)
			s.reEnterStory()
		}
		return
	} else {
		log.Debugf("Entering idx %d of Container: %v", s.currentAddress.I, s.currentAddress.C.Name)
		log.Debugf("Item is %q, %s", s.currentAddress.C.Contents[s.currentAddress.I], reflect.TypeOf(s.currentAddress.C.Contents[s.currentAddress.I]))
		s.currentAddress.C.Contents[s.currentAddress.I].Accept(s)
	}
}

func (s *Story) endStory() {
	log.Debug("Ending Story")
	s.state.text = ""
	s.writeToState()
	s.Finished = true
	s.state.done = true
}

func (s *Story) moveToPath(path types.Path) {
	target, idx := s.ResolvePath(path)
	s.enterContainer(target, idx)
}

func (s *Story) choose(c Choice) {
	cnt, idx := s.ResolvePath(c.Destination)
	s.enterContainer(cnt, idx)
	s.state.TurnCount++
	s.state.CurrentChoices = s.state.CurrentChoices[:0]
	s.state.setDone(false)
}

func (s *Story) enterContainer(c *types.Container, idx int) {
	s.currentAddress.Set(c, idx)
	s.state.RecordContainer(c, idx)
}

func (s *Story) ChoseIndex(idx int) error {
	if idx < 0 || idx >= len(s.state.CurrentChoices) {
		return fmt.Errorf("%d is out of range of choices: %d", idx, len(s.state.CurrentChoices))
	}
	s.choose(s.state.CurrentChoices[idx])
	return nil
}

func (s *Story) writeToState() {
	str := ""
	if !s.outputBuffer.Empty() && s.mode == None {

		for !s.outputBuffer.Empty() {
			text, _ := s.outputBuffer.Pop()
			str = text + str
		}
		s.state.text = str
		log.Debugf("Wrote: \"%s\"", strings.Replace(str, "\n", "\\n", -1))
		s.outputBuffer.Clear()
	}
}
