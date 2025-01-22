package types

type Visitor interface {
	VisitString(StringVal)
	VisitControlCommand(ControlCommand)
	VisitContainer(*Container)
	VisitDivertTarget(DivertTarget)
	VisitVariableDivert(VariableDivert)
	VisitDivert(Divert)
	VisitTmpVar(TempVar)
	VisitOperator(Operator)
	VisitChoicePoint(ChoicePoint)
	VisitIntVal(IntVal)
	VisitFloatVal(FloatVal)
	VisitBoolVal(BoolVal)
	VisitGlobalVar(GlobalVar)
	VisitVarRef(VarRef)
	VisitReadCount(ReadCount)
	VisitVariablePointer(VariablePointer)
	VisitFunctionDivert(FunctionDivert)
}

type Acceptor interface {
	Accept(Visitor)
}
