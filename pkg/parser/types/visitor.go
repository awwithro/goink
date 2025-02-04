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
	VisitExternalFunctionDivert(ExternalFunctionDivert)
	VisitTunnelDivert(TunnelDivert)
	VisitListInit(ListInit)
	VisitListVal(ListVal)
	VisitListValItem(*ListValItem)
}

type Acceptor interface {
	Accept(Visitor)
}
