package runtime

import (
	"strings"

	"github.com/awwithro/goink/pkg/parser/types"
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
	case types.Glue:
		s.glue()
	case types.Void:
		s.evaluationStack.Push(types.VoidVal{})
	case types.ReturnTunnel:
		s.returnTunnel()
	default:
		log.Panic("Unimplemented Command! ", cmd)
	}
	s.currentAddress.Increment()
}

func (s *Story) VisitTmpVar(v types.TempVar) {
	defer s.currentAddress.Increment()
	p := mustPopStack[any](s.evaluationStack)
	// If this is a reassignment to a variablePointer, dereference and set the target
	if v.ReAssign {
		val := s.state.GetVar(v.Name)
		if ref, ok := val.(types.VariablePointer); ok {
			s.setVariablePointerValue(ref, p)
			return
		}
	}
	switch val := p.(type) {
	case types.DivertTarget:
		s.state.SetVar(v.Name, types.Path(val))
	case types.VariablePointer:
		s.state.SetVar(v.Name, val)
	default:
		s.state.SetVar(v.Name, val)
	}
}

func (s *Story) VisitDivertTarget(divert types.DivertTarget) {
	s.evaluationStack.Push(divert)
	s.currentAddress.Increment()
}
func (s *Story) doDivert(divert types.Divert) {
	if divert.Conditional {
		visit := mustPopStack[types.Truthy](s.evaluationStack)
		if visit.AsBool() {
			s.moveToPath(divert.Path)
			return
		} else {
			log.Debug("Conditional divert failed")
			// If we don't divert, advance the index
			s.currentAddress.Increment()
			return
		}
	}
	log.Debug("Diverting to: ", divert.Path)
	s.moveToPath(divert.Path)
}
func (s *Story) VisitDivert(divert types.Divert) {
	s.doDivert(divert)
}

func (s *Story) VisitFunctionDivert(f types.FunctionDivert) {
	// Pushes the old address so we know where to return to
	// after the function runs
	oldAddr := s.currentAddress
	oldAddr.Increment()
	s.previousAddress.Push(oldAddr)
	s.doDivert(f.Divert)
}

func (s *Story) VisitTunnelDivert(t types.TunnelDivert) {
	// Pushes the old address so we know where to return to
	// after the function runs
	oldAddr := s.currentAddress
	oldAddr.Increment()
	s.previousAddress.Push(oldAddr)
	s.doDivert(t.Divert)
}

func (s *Story) VisitVariableDivert(divert types.VariableDivert) {
	// TODO: Can we have a variable divert with condition?
	// if so, do we check for the condition or divert path first?
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
		x := mustPopStack[types.Truthy](s.evaluationStack)
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

func (s *Story) VisitVoidVal(b types.VoidVal) {
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
	var finalVal any
	var intermediateValue any

	// first see if the var is set
	if val, ok := s.state.tmpVars[string(v)]; ok {
		intermediateValue = val
	} else if val, ok = s.state.globalVars[string(v)]; ok {
		intermediateValue = val
	} else {
		log.Debugf("temp %v\n", s.state.tmpVars)
		log.Debugf("global %v\n", s.state.tmpVars)
		log.Panic("ref to unset var ", string(v))
	}

	// if set, see if we need to dereference a variable pointer
	if p, ok := intermediateValue.(types.VariablePointer); ok {
		finalVal = s.getVariablePointerValue(p)
	} else {
		finalVal = intermediateValue
	}
	log.Debugf("Pushing val %v", finalVal)
	s.evaluationStack.Push(finalVal)
	s.currentAddress.Increment()
}

func (s *Story) VisitReadCount(r types.ReadCount) {
	addr := s.ResolvePath(types.Path(r))
	count := s.state.visitCounts[addr.C]
	s.evaluationStack.Push(types.IntVal(count))
	s.currentAddress.Increment()
}

func (s *Story) VisitVariablePointer(v types.VariablePointer) {
	s.evaluationStack.Push(v)
	s.currentAddress.Increment()

}

func (s *Story) getVariablePointerValue(p types.VariablePointer) any {
	// TODO: Use the ci of p to determine if global or local
	val, ok := s.state.globalVars[p.Name]
	if !ok {
		log.Panic("nil pointer, no var ", p.Name)
	}
	return val
}
func (s *Story) setVariablePointerValue(p types.VariablePointer, val any) {
	// TODO: Use the ci of p to determine if global or local
	s.state.globalVars[p.Name] = val
}
func (s *Story) glue() {
	val, _ := s.outputBuffer.Peek()
	for val == "\n" {
		s.outputBuffer.Pop()
		log.Debug("Popped newline")
		val, _ = s.outputBuffer.Peek()
	}
}

func (s *Story) returnTunnel() {
	s.currentAddress, _ = s.previousAddress.Pop()
}
