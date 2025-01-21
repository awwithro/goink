package types

type ControlCommand int

var controlCommandMap = map[string]ControlCommand{
	"ev":        StartEvalMode,
	"/ev":       EndEvalMode,
	"str":       StartStrMode,
	"/str":      EndStrMode,
	"out":       PopOutput,
	"pop":       Pop,
	"->->":      ReturnTunnel,
	"~ret":      ReturnFunction,
	"du":        Duplicate,
	"nop":       NoOp,
	"choiceCnt": ChoiceCount,
	"turn":      PushTurn,
	"turns":     PushTurnsSinceTarget,
	"visit":     VisitCount,
	"seq":       Sequence,
	"thread":    Thread,
	"done":      Done,
	"end":       End,
	"void":      Void,
}

const (
	StartEvalMode ControlCommand = iota
	EndEvalMode
	StartStrMode
	EndStrMode
	PopOutput
	Pop
	ReturnTunnel
	ReturnFunction
	Duplicate
	NoOp
	ChoiceCount
	PushTurn
	PushTurnsSinceTarget
	VisitCount
	Sequence
	Thread
	Done
	End
	Void
)

func IsControlCommand(str string) (ControlCommand, bool) {
	c, ok := controlCommandMap[str]
	return c, ok
}

func (c ControlCommand) Accept(v Visitor) {
	v.VisitControlCommand(c)
}
