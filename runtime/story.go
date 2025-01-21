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
	currentAddress  Address
	previousAddress stacks.Stack[Address]
	Finished        bool
}

func NewStory(ink types.Ink) Story {
	s := Story{
		ink:             ink,
		evaluationStack: arraystack.New[any](),
		outputBuffer:    arraystack.New[string](),
		mode:            None,
		stringMarker:    -1,
		state:           NewStoryState(),
		previousAddress: arraystack.New[Address](),
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
	str := mustPopStack[fmt.Stringer](s.evaluationStack)
	s.outputBuffer.Push(str.String())

}

func (s *Story) pushVisitCount() {
	count := s.state.visitCounts[s.currentAddress.C]
	s.evaluationStack.Push(types.IntVal(count))
}

func (s *Story) duplicateTopItem() {
	item, ok := s.evaluationStack.Peek()
	if ok {
		s.evaluationStack.Push(item)
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

func panicInvalidStackType[T any](expected T, actual any) {
	log.Panicf("non %s in stack: %s, %v", reflect.TypeOf(expected), reflect.TypeOf(actual), actual)
}

func (s *Story) Start() {
	// NEED TO PARSE ALL GLOBALS IN SPECIAL "global decal" Subcontainer of root
	s.enterContainer(Address{C: s.ink.Root.Contents[0].(*types.Container), I: 0})
}

func mustPopStack[T any](s stacks.Stack[any]) T {
	val, ok := s.Pop()
	if !ok {
		log.Panic("Popped empty stack")
	}
	log.Debugf("Popped %s %v", reflect.TypeOf(val), val)
	ret, ok := val.(T)
	if !ok {
		panicInvalidStackType[T](ret, val)
	}
	return ret
}

func (s *Story) ResolvePath(p types.Path) Address {
	c, i := types.ResolvePath(p, s.currentAddress.C)
	return Address{C: c, I: i}
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
		// End of the story?
		if pos, err := s.currentAddress.C.PositionInParent(); err != nil {
			// A choice is needed
			if len(s.state.CurrentChoices) > 0 {
				s.state.setDone(true)
				return
			}
			switch err.(type) {
			case types.EndOfSubContainer:
				// If we've reached the end of a sub-container, this is an implicit end of the story
				s.endStory()
			default:
				// NoParent, at the root container
				log.Debug("reached end of ink ", err)
				s.endStory()
			}
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
	log.Debugf("Ending Story. Located at %s %v", s.currentAddress.C.Name, s.state.CurrentChoices)
	s.state.text = ""
	s.writeToState()
	s.Finished = true
	s.state.done = true
}

func (s *Story) moveToPath(path types.Path) {
	a := s.ResolvePath(path)
	s.enterContainer(a)
}

func (s *Story) choose(c Choice) {
	s.enterContainer(c.Destination)
	s.state.TurnCount++
	s.state.CurrentChoices = s.state.CurrentChoices[:0]
	s.state.setDone(false)
}

func (s *Story) enterContainer(a Address) {
	s.currentAddress = a
	s.state.RecordContainer(a)
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
