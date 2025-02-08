package runtime

import (
	"math/rand"
	"slices"
	"strings"

	"github.com/awwithro/goink/pkg/parser/types"
	log "github.com/sirupsen/logrus"
)

func (s *Story) VisitString(str types.StringVal) {
	log.Debugf("Visiting String: \"%s\"", strings.Replace(str.String(), "\n", "\\n", -1))
	if s.mode == Eval {
		s.evaluationStack.Push(str)
	} else {
		s.outputBuffer.Push(str.String())
	}
	s.currentAddress.Increment()
}
func (s *Story) VisitControlCommand(cmd types.ControlCommand) {
	log.Debug("Visiting Control Command: ", cmd)
	switch cmd {
	case types.StartEvalMode:
		s.startEvalMode()
	case types.EndEvalMode:
		s.endEvalMode()
	case types.StartStrMode:
		s.startStrMode()
	case types.EndStrMode:
		s.endStrMode()
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
	case types.ReturnFunction:
		fallthrough
	case types.ReturnTunnel:
		s.returnTunnel()
	case types.StartTag:
		s.startTagMode()
	case types.EndTag:
		s.endTagMode()
	case types.Thread:
		log.Warn("bypassing thread")
		s.currentAddress.Increment()
	case types.Sequence:
		s.generateSequence()
	case types.PushTurnsSinceTarget:
		s.pushTurnsSinceTarget()
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

// pushes the previous address onto the stack as a return val
// increment controls if the previous address is incremented prior to pushing
func (s *Story) pushStackDivert(divert types.Divert, increment bool) {
	// Pushes the old address so we know where to return to
	// after the function runs
	oldAddr := s.currentAddress
	if increment {
		oldAddr.Increment()
	}
	s.previousState.Push(State{address: oldAddr, mode: s.mode})
	s.mode = None
	s.doDivert(divert)
}

func (s *Story) VisitFunctionDivert(f types.FunctionDivert) {
	s.pushStackDivert(f.Divert, true)
}

func (s *Story) VisitTunnelDivert(t types.TunnelDivert) {
	s.pushStackDivert(t.Divert, false)
}

func (s *Story) VisitVariableDivert(divert types.VariableDivert) {
	// TODO: Can we have a variable divert with condition?
	// if so, do we check for the condition or divert path first?
	log.Debug("Visit Variable Divert ", divert.Name)
	p := s.state.GetVar(divert.Name)
	if path, ok := p.(types.Path); ok {
		s.moveToPath(path)
	} else {
		panicInvalidStackType[types.Path](path, s)
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
		if _, ok := s.state.visitCounts[a.C]; ok {
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
	s.state.currentChoices = append(s.state.currentChoices, choice)
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
		log.Debugf("global %v\n", s.state.globalVars)
		s.Panicf("ref to unset var %s\nGlobal: %v", string(v), s.state.globalVars)
	}

	// if set, see if we need to dereference a variable pointer
	if p, ok := intermediateValue.(types.VariablePointer); ok {
		log.Debug("Deref a pointer")
		finalVal = s.getVariablePointerValue(p)
	} else {
		finalVal = intermediateValue
	}
	log.Debugf("Pushing val %v of type %T", finalVal, finalVal)
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

func (s *Story) VisitExternalFunctionDivert(e types.ExternalFunctionDivert) {
	if f, ok := s.extFuncs[string(e.Path)]; !ok {
		log.Warnf("External func %s not registered, using fallback", string(e.Path))
		s.pushStackDivert(e.Divert, true)
	} else {
		defer s.currentAddress.Increment()
		args := []any{}
		if e.Args > 0 {
			for x := 0; x < e.Args; x++ {
				a := mustPopStack[any](s.evaluationStack)
				switch arg := a.(type) {
				case types.BoolVal:
					args = append(args, arg.AsBool())
				case types.NumericVal:
					if arg.IsFloat() {
						args = append(args, arg.AsFloat())
					} else {
						args = append(args, arg.AsInt())
					}
				case types.StringVal:
					args = append(args, arg.String())
				default:
					s.Panicf("Unrecognized type for external func %T", arg)
				}
			}
		}
		slices.Reverse(args)
		res := f(args)
		if res != nil {
			switch val := res.(type) {
			case int:
				s.evaluationStack.Push(types.IntVal(val))
			case bool:
				s.evaluationStack.Push(types.BoolVal(val))
			case float64:
				s.evaluationStack.Push(types.FloatVal(val))
			case string:
				s.evaluationStack.Push(types.StringVal(val))
			default:
				s.Panicf("unrecognized return value for external func %T", val)
			}
		} else {
			log.Debug("No return val from external func, pushing void")
			s.evaluationStack.Push(types.VoidVal{})
		}
	}
}

func (s *Story) getVariablePointerValue(p types.VariablePointer) any {
	// TODO: Use the ci of p to determine if global or local
	if val, ok := s.state.globalVars[p.Name]; ok {
		log.Debugf("Pointer to %T: %v", val, val)
		return val
	} else if val, ok = s.state.tmpVars[p.Name];ok{
		log.Debugf("Pointer to %T: %v", val, val)
		return val
	}
	s.Panicf("nil pointer, no var %s", p.Name)
	return nil
}
func (s *Story) setVariablePointerValue(p types.VariablePointer, val any) {
	// TODO: Use the ci of p to determine if global or local
	log.Debugf("Setting Pointer Var, %T: %v named %s", val, val, p.Name)
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
	oldState, _ := s.previousState.Pop()
	s.currentAddress = oldState.address
	s.mode = oldState.mode
	log.Debugf("Tunnel Returned to Name: %s Idx: %d", s.currentAddress.C.Name, s.currentAddress.I)
}

func (s *Story) VisitListInit(l types.ListInit) {
	log.Debugf("Visiting ListInit %v", l)
	list := s.initializeList(l)
	s.evaluationStack.Push(list)
	s.currentAddress.Increment()
}

func (s *Story) VisitListValItem(l *types.ListValItem) {
	log.Debugf("Visiting ListValItem %v", l)
	s.evaluationStack.Push(l)
	s.currentAddress.Increment()
}

func (s *Story) VisitListVal(l types.ListVal) {
	log.Debugf("Visiting ListVal %v", l)
	s.evaluationStack.Push(l)
	s.currentAddress.Increment()
}

func (s *Story) initializeList(init types.ListInit) (list types.ListVal) {
	list = types.ListVal{}
	log.Debug("init list: ", init)
	for _, name := range init.Origins {
		if def, ok := s.computedLists[name]; !ok {
			s.Panicf("origin referenced an undefined list %s", name)
		} else {
			for _, v := range def {
				log.Debugf("asspending %v to list", v)
				list = append(list, v)
			}
		}
	}
	// the list setup is odd, it references the global list name with the item
	// name as well as the item val
	for name := range init.List {
		segs := strings.Split(name, ".")
		list = append(list, s.computedLists[segs[0]].Get(segs[1]))
	}
	return list
}

func (s *Story) generateSequence() {
	seq := mustPopStack[types.IntVal](s.evaluationStack)
	val := rand.Intn(seq.AsInt()-1) + 1
	res := types.IntVal(val)
	log.Debugf("Generated Sequence number: %d", res.AsInt())
	s.evaluationStack.Push(res)
}

func (s *Story) pushTurnsSinceTarget() {
	divert := mustPopStack[types.DivertTarget](s.evaluationStack)
	target := s.ResolvePath(types.Path(divert))
	turns := s.state.LastTurnVisited(target.C)
	delta := s.state.TurnCount - turns
	s.evaluationStack.Push(types.IntVal(delta))
}
