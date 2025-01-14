package runtime

type StoryState struct {
	globalVars     map[string]any
	currentChoices []string
	tmpVars        map[string]any
}

func NewStoryState() *StoryState {
	s := &StoryState{
		globalVars:     make(map[string]any),
		currentChoices: []string{},
		tmpVars:        make(map[string]any),
	}
	return s
}

func (s *StoryState) SetVar(name string, val any) {
	s.tmpVars[name] = val
}
