package types

type Visitor interface {
	VisitString(StringVal)
	VisitVoidVal(VoidVal)
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
	VisitTunnelDivert(TunnelDivert)
}

type Acceptor interface {
	Accept(Visitor)
}
