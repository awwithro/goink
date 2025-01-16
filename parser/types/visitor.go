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
}

type Acceptor interface {
	Accept(Visitor)
}
