package types

type Visitor interface {
	VisitString(StringVal)
	VisitControlCommand(ControlCommand)
	VisitContainer(Container)
	VisitDivertTarget(DivertTarget)
	VisitTmpVar(TempVar)
	VisitOperator(Operator)
}

type Acceptor interface {
	Accept(Visitor)
}
