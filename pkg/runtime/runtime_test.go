package runtime

import (
	"os"
	"testing"

	"github.com/awwithro/goink/pkg/parser"
	"github.com/awwithro/goink/pkg/parser/types"
	"github.com/stretchr/testify/assert"
)

// runs the given ink story. Uses the choices in order until the story is finished
func TestFullInkStory(t *testing.T) {
	testCases := []struct {
		desc            string
		inkJsonFilePath string
		choices         []int
		expectedText    string
	}{
		{
			desc:            "Hello World",
			inkJsonFilePath: "../../examples/hello.json",
			choices:         []int{0},
			expectedText:    "Hello, world.\n",
		},
		{
			desc:            "Easy First Path",
			inkJsonFilePath: "../../examples/easy.json",
			choices:         []int{0},
			expectedText:    "There were two choices.\nThey lived happily ever after.\n",
		},
		{
			desc:            "Easy Second Path",
			inkJsonFilePath: "../../examples/easy.json",
			choices:         []int{1},
			expectedText:    "There were four lines of content.\nThey lived happily ever after.\n",
		},
	}
	parsed := map[string]types.Ink{}
	for _, tC := range testCases {
		assert := assert.New(t)
		t.Run(tC.desc, func(t *testing.T) {
			var ink types.Ink
			var ok bool
			// see if we've already parsed this in another test
			if ink, ok = parsed[tC.inkJsonFilePath]; !ok {
				js, err := os.ReadFile(tC.inkJsonFilePath)
				assert.NoError(err)
				ink = parser.Parse(js)
				parsed[tC.inkJsonFilePath] = ink
			}
			s := NewStory(ink)
			s.Start()
			choiceIdx := 0
			var state StoryState
			var err error
			for !s.IsFinished() {
				state, err = s.RunContinuous()
				assert.NoError(err)
				if len(state.CurrentChoices) > 0 {
					err = s.ChoseIndex(tC.choices[choiceIdx])
					assert.NoError(err)
					choiceIdx++
				}
			}
			assert.Equal(tC.expectedText, state.GetText())
		})
	}
}
