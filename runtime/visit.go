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
	s.currentIdx++
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
		_ = mustPopStack(s.evaluationStack)
	case types.Done:
		s.state.setDone(true)
	case types.End:
		s.endStory()
	case types.NoOp:
	default:
		log.Panic("Unimplemented Command! ", cmd)
	}
	s.currentIdx++
}

func (s *Story) VisitTmpVar(v types.TempVar) {
	p := mustPopStack(s.evaluationStack)
	switch val := p.(type) {
	case types.Path:
		s.state.SetVar(v.Name, val)
	case types.DivertTarget:
		s.state.SetVar(v.Name, types.Path(val))
	default:
		log.Panic("don't know how to set tmp var to ", reflect.TypeOf(val))
	}
	s.currentIdx++
}

func (s *Story) VisitDivertTarget(divert types.DivertTarget) {
	s.evaluationStack.Push(divert)
	s.currentIdx++
}

func (s *Story) VisitDivert(divert types.Divert) {
	if divert.Conditional {
		item := mustPopStack(s.evaluationStack)
		switch ok := item.(type) {
		case bool:
			if ok {
				s.moveToPath(divert.Path)
				return
			} else {
				log.Debug("Conditional divert failed")
				// If we don't divert, advance the index
				s.currentIdx++
				return
			}
		default:
			panicInvalidStackType(true, ok)
		}
	}
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
	if p.HasCondition() {
		x := mustPopNumeric(s.evaluationStack)
		if !x.AsBool() {
			return
		}
	}
	choice := Choice{Destination: p.Path}
	if p.OnceOnly() {
		cnt, _ := s.ResolvePath(p.Path)
		if s.state.visitCounts[cnt] > 0 {
			return
		}
	}
	if p.HasChoiceOnly() {
		x := mustPopStack(s.evaluationStack)
		if txt, ok := x.(types.StringVal); ok {
			choice.choiceOnlyText = txt.String()
		} else {
			panicInvalidStackType(txt, x)
		}
	}
	if p.HasStartContent() {
		x := mustPopStack(s.evaluationStack)
		if txt, ok := x.(types.StringVal); ok {
			choice.text = txt.String()
		} else {
			panicInvalidStackType(txt, x)
		}
	}
	if p.IsInvisibleDefault() {
		choice.OnlyDefault = true
	}
	s.state.CurrentChoices = append(s.state.CurrentChoices, choice)
	s.currentIdx++
}

func (s *Story) VisitContainer(c *types.Container) {
	log.Debug("Visiting Container: ", c.Name)
	s.enterContainer(c, 0)
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
	s.currentIdx++
}

func (s *Story) VisitFloatBoolVal(f types.FloatVal) {
	s.visitNumber(f)
}

func (s *Story) VisitBoolVal(b types.BoolVal) {
	s.evaluationStack.Push(b)
	s.currentIdx++
}

func (s *Story) VisitGlobalVar(v types.GlobalVar) {
	log.Debug("Visiting Global Var ", v.Name)
	val := mustPopStack(s.evaluationStack)
	s.state.globalVars[v.Name] = val
	s.currentIdx++
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
	s.currentIdx++
}
