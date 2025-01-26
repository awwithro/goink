package runtime

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/awwithro/goink/pkg/parser/types"
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
	if s.mode == Eval {
		log.Warn("starting eval mode while already in eval mode")
		return
	}
	if s.mode != None {
		panicInvalidModeTransition(s.mode, Eval)
	}
	s.mode = Eval
}

func (s *Story) endEvalMode() {
	// FIXME: I think i see what's happening
	// when we switch to a function, we should have a fresh state of
	// tmp vars and modes (probably more). When we return, we should be
	// back in whatever mode we were in before instead of the mode we were
	// in while in the function. This must be the call stack referenced in the runtime
	if s.mode == None {
		log.Warn("ending eval mode while not in eval mode")
		return
	}
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

func panicInvalidStackType[T any](actual any) {
	log.Panicf("non %T in stack: %T, %v", new(T), actual, actual)
}

func (s *Story) Start() {
	s.setupGlobalVars()
	s.enterContainer(Address{C: s.ink.Root.Contents[0].(*types.Container), I: 0})
}

func mustPopStack[T any](s stacks.Stack[any]) T {
	val, ok := s.Pop()
	if !ok {
		log.Panic("Popped empty stack")
	}
	log.Debugf("Popped %T %v", val, val)
	ret, ok := val.(T)
	if !ok {
		panicInvalidStackType[T](val)
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
		if !s.state.CanContinue() && len(s.state.GetChoices()) > 0 {
			// check if only default choices remain
			onlyDefaults := true
			choices := s.state.GetChoices()
			for _, choice := range choices {
				if !choice.OnlyDefault {
					onlyDefaults = false
					break
				}
			}
			// if we only have a default choices, chose it
			if onlyDefaults {
				s.choose(choices[0])
				s.state.setDone(true)
			} else {
				// we have a choice to be made, write the story so far
				s.writeToState()
			}

		}
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
			if len(s.state.GetChoices()) > 0 {
				s.state.setDone(true)
				return
			}
			switch err.(type) {
			case types.EndOfSubContainer:
				// If we've reached the end of a sub-container, this is an implicit end of the story
				// unless there is a previous address on the stack (we're at the end of a function call)
				if s.previousAddress.Size() > 0 {
					s.currentAddress, _ = s.previousAddress.Pop()
					// Assuming we're always returning from a function,
					// this assumption likely doesn't hold up
					s.evaluationStack.Push(types.Void)
					s.reEnterStory()
				} else {
					s.endStory()
				}
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
	log.Debugf("Ending Story. Located at %s %v", s.currentAddress.C.Name, s.state.currentChoices)
	s.state.text = ""
	s.writeToState()
	s.state.Finished = true
	s.state.done = true
}

func (s *Story) moveToPath(path types.Path) {
	a := s.ResolvePath(path)
	s.enterContainer(a)
}

func (s *Story) choose(c Choice) {
	s.enterContainer(c.Destination)
	s.state.TurnCount++
	s.state.currentChoices = s.state.currentChoices[:0]
	s.state.setDone(false)
}

func (s *Story) enterContainer(a Address) {
	s.currentAddress = a
	s.state.RecordContainer(a)
}

func (s *Story) ChoseIndex(idx int) error {
	if idx < 0 || idx >= len(s.state.GetChoices()) {
		return fmt.Errorf("%d is out of range of choices: %d", idx, len(s.state.GetChoices()))
	}
	s.choose(s.state.GetChoices()[idx])
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

func (s *Story) IsFinished() bool {
	return s.state.Finished
}

func (s *Story) setupGlobalVars() {
	c, err := s.ink.Root.GetNamedContainer(types.GlobalVarKey)
	// no global vars to work parse
	if err != nil {
		return
	}
	s.currentAddress = Address{C: c, I: 0}
	for s.state.CanContinue() {
		if _, err := s.Step(); err != nil {
			log.Panic("failed while parsing globals ", err)
		}
	}
	// Globals run until a "end" statement
	s.state.Finished = false
}

func (s *Story) RunContinuous() (state StoryState, err error) {
	run := true
	for run {
		if state, err = s.Step(); err != nil {
			return state, err
		} else if state.CanContinue() && !state.Finished {
			continue
		}
		run = false
	}
	return state, err
}
