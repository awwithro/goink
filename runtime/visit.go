package runtime

import (
	"reflect"
	"strings"

	"github.com/awwithro/goink/parser/types"
	log "github.com/sirupsen/logrus"
)

func (s *Story) VisitString(str types.StringVal) {
	log.Debugf("Visiting String: \"%s\"", strings.Replace(str.String(), "\n", "\\n", -1))
	if s.mode == Eval {
		log.Panicf("String encountered while in Eval mode: %s", str)
	}
	s.outputBuffer.Push(str.String())
	s.currentAddress.Increment()
}
func (s *Story) VisitControlCommand(cmd types.ControlCommand) {
	log.Debug("Visiting Control Command: ", cmd)
	switch cmd {
	case types.StartEvalMode:
		s.startEvalMode()
	case types.StartStrMode:
		s.startStrMode()
	case types.EndStrMode:
		s.endStrMode()
	case types.EndEvalMode:
		s.endEvalMode()
	case types.PopOutput:
		s.popOutput()
	case types.Pop:
		_ = mustPopStack[any](s.evaluationStack)
	case types.Done:
		s.state.setDone(true)
	case types.End:
		s.endStory()
	case types.NoOp:
	case types.VisitCount:
		s.pushVisitCount()
	case types.Duplicate:
		s.duplicateTopItem()
	default:
		log.Panic("Unimplemented Command! ", cmd)
	}
	s.currentAddress.Increment()
}

func (s *Story) VisitTmpVar(v types.TempVar) {
	p := mustPopStack[any](s.evaluationStack)
	switch val := p.(type) {
	case types.Path:
		s.state.SetVar(v.Name, val)
	case types.DivertTarget:
		s.state.SetVar(v.Name, types.Path(val))
	case types.VariablePointer:
		// ref := s.state.GetVar(val.Name)
		// x := ref.(types.Acceptor)
		// s.state.SetVar(v.Name, *x)
		// Need to set a var that is a pointer to what this points to
	default:
		log.Panic("don't know how to set tmp var to ", reflect.TypeOf(val))
	}
	s.currentAddress.Increment()
}

func (s *Story) VisitDivertTarget(divert types.DivertTarget) {
	s.evaluationStack.Push(divert)
	s.currentAddress.Increment()
}

func (s *Story) VisitDivert(divert types.Divert) {
	if divert.Conditional {
		var visit bool
		item := mustPopStack[any](s.evaluationStack)
		switch ok := item.(type) {
		case bool:
			visit = ok
		case types.NumericVal:
			visit = ok.AsBool()
		default:
			panicInvalidStackType(true, ok)
		}
		if visit {
			s.moveToPath(divert.Path)
			return
		} else {
			log.Debug("Conditional divert failed")
			// If we don't divert, advance the index
			s.currentAddress.Increment()
			return
		}
	}
	if s.mode == None {
		oldAddr := s.currentAddress
		oldAddr.Increment()
		s.previousAddress.Push(oldAddr)
	}
	log.Debug("Diverting to: ", divert.Path)
	s.moveToPath(divert.Path)

}

func (s *Story) VisitVariableDivert(divert types.VariableDivert) {
	log.Debug("Visit Variable Divert ", divert.Name)
	p := s.state.GetVar(divert.Name)
	if path, ok := p.(types.Path); ok {
		s.moveToPath(path)
	} else {
		panicInvalidStackType(path, p)
	}
}

func (s *Story) VisitChoicePoint(p types.ChoicePoint) {
	log.Debug("Visit Choice Point ", p.Path)
	a := s.ResolvePath(p.Path)
	defer s.currentAddress.Increment()
	if p.HasCondition() {
		x := mustPopStack[types.NumericVal](s.evaluationStack)
		if !x.AsBool() {
			return
		}
	}
	choice := Choice{Destination: a}
	if p.OnceOnly() {
		if s.state.visitCounts[a.C] > 0 {
			return
		}
	}
	if p.HasChoiceOnly() {
		txt := mustPopStack[types.StringVal](s.evaluationStack)
		choice.choiceOnlyText = txt.String()

	}
	if p.HasStartContent() {
		txt := mustPopStack[types.StringVal](s.evaluationStack)
		choice.text = txt.String()

	}
	if p.IsInvisibleDefault() {
		choice.OnlyDefault = true
	}
	s.state.CurrentChoices = append(s.state.CurrentChoices, choice)
}

func (s *Story) VisitContainer(c *types.Container) {
	log.Debug("Visiting Container: ", c.Name)
	s.enterContainer(Address{C: c, I: 0})
	// no advance needed
}

func (s *Story) VisitIntVal(i types.IntVal) {
	s.visitNumber(i)
}
func (s *Story) VisitFloatVal(f types.FloatVal) {
	s.visitNumber(f)
}

func (s *Story) visitNumber(i types.NumericVal) {
	s.evaluationStack.Push(i)
	s.currentAddress.Increment()
}

func (s *Story) VisitFloatBoolVal(f types.FloatVal) {
	s.visitNumber(f)
}

func (s *Story) VisitBoolVal(b types.BoolVal) {
	s.evaluationStack.Push(b)
	s.currentAddress.Increment()
}

func (s *Story) VisitGlobalVar(v types.GlobalVar) {
	log.Debug("Visiting Global Var ", v.Name)
	val := mustPopStack[any](s.evaluationStack)
	s.state.globalVars[v.Name] = val
	s.currentAddress.Increment()
}

func (s *Story) VisitVarRef(v types.VarRef) {
	log.Debug("Visiting Var Ref ", v)
	var val any
	var ok bool
	if val, ok = s.state.tmpVars[string(v)]; ok {
	} else if val, ok = s.state.globalVars[string(v)]; ok {
	} else {
		val = false
	}
	log.Debugf("Pushing val %v", val)
	s.evaluationStack.Push(val)
	s.currentAddress.Increment()
}

func (s *Story) VisitReadCount(r types.ReadCount) {
	addr := s.ResolvePath(types.Path(r))
	count := s.state.visitCounts[addr.C]
	s.evaluationStack.Push(types.IntVal(count))
	s.currentAddress.Increment()
}

func (s *Story) VisitVariablePointer(v types.VariablePointer) {

}
