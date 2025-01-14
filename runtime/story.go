package runtime

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/awwithro/goink/parser/types"
	"github.com/emirpasic/gods/v2/stacks"
	"github.com/emirpasic/gods/v2/stacks/arraystack"
	"github.com/sirupsen/logrus"
)

type Story struct {
	ink              types.Ink
	evaluationStack  stacks.Stack[any]
	outputBuffer     stacks.Stack[string]
	mode             Mode
	stringMarker     int //Used to track index of stack to concatenate into a string
	state            *StoryState
	currentContainer *types.Container
}

func NewStory(ink types.Ink) Story {
	s := Story{
		ink:             ink,
		evaluationStack: arraystack.New[any](),
		outputBuffer:    arraystack.New[string](),
		mode:            Str,
		stringMarker:    -1,
		state:           NewStoryState(),
	}
	return s
}

func (s *Story) VisitString(str types.StringVal) {
	if s.mode == Eval {
		logrus.Panicf("String encountered while in Eval mode: %s", str)
	}
	s.outputBuffer.Push(str.String())
}
func (s *Story) VisitControlCommand(cmd types.ControlCommand) {
	switch cmd {
	case types.StartEvalMode:
		s.startEvalMode()
	case types.StartStrMode:
		s.startStrMode()
	case types.EndStrMode:
		s.endStrMode()
	case types.EndEvalMode:

	case types.PopOutput:
		s.popOutput()
	case types.Pop:
		_ = mustPopStack(s.evaluationStack)
	case types.NoOp:
	default:
		logrus.Panic("Unimplemented Command! ", cmd)
	}
}

func (s *Story) VisitTmpVar(v types.TempVar) {
	pth := mustPopStack(s.evaluationStack)
	path := pth.(types.Path)
	s.state.SetVar(v.Name, path)

}

func (s *Story) VisitDivertTarget(divert types.DivertTarget) {
	s.evaluationStack.Push(divert)
}

func (s *Story) startEvalMode() {
	if s.mode != None {
		panicInvalidModeTransition(s.mode, Eval)
	}
	s.mode = Eval
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
		logrus.Panic("Couldn't push non Stringer type to output ", reflect.TypeOf(val))
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
		logrus.Panic("No elements could be popped from the output buffer")
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

func (s *Story) VisitContainer(cmd types.Container) {}

func panicInvalidModeTransition(current, attempted Mode) {
	logrus.Panicf("Invalid Mode transition. Can't go from %s to %s", current, attempted)
}

func (s *Story) Start() {
	_ = s.ink.Root.Contents[0]
}

func mustPopStack(s stacks.Stack[any]) any {
	val, ok := s.Pop()
	if !ok {
		logrus.Panic("Popped empty stack")
	}
	return val
}

func mustPopNumeric(s stacks.Stack[any]) types.NumericVal {
	x := mustPopStack(s)
	num, ok := x.(types.NumericVal)
	if !ok {
		logrus.Panic("Non number popper ", x)
	}
	return num
}

func (s *Story) ResolvePath(p types.Path) *types.Container {
	return types.ResolvePath(p, &s.ink.Root, s.currentContainer)
}
