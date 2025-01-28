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
		choiceCounts    []int // expected choice counts at each point we're offered choices
		expectedText    string
		externalFuncs   map[string]func([]any) any
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
		{
			desc:            "Loop With Default Only Choices",
			inkJsonFilePath: "../../examples/fallback.json",
			choices:         []int{0, 0},
			choiceCounts:    []int{2, 1},
			expectedText:    "four test\n\ntwo test\n\nthree ",
		},
		{
			desc:            "External Function Fallback",
			inkJsonFilePath: "../../examples/externalfunc.json",
			expectedText:    "Calling Func Internal Hello world\n\n",
		},
		{
			desc:            "External Function",
			inkJsonFilePath: "../../examples/externalfunc.json",
			expectedText:    "Calling Func External Hello world\n",
			externalFuncs: map[string]func([]any) any{
				"Hello": hello,
			},
		},
	}
	parsed := map[string]types.Ink{}
	for _, tC := range testCases {
		assert := assert.New(t)
		t.Run(tC.desc, func(t *testing.T) {
			var assesChoiceCounts bool
			if len(tC.choiceCounts) > 0 {
				assesChoiceCounts = true
			}
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

			// Register functions if we have them
			if tC.externalFuncs != nil {
				for k, v := range tC.externalFuncs {
					s.RegisterExternalFunction(k, v)
				}
			}
			s.Start()
			choiceIdx := 0
			var state StoryState
			var err error
			for !s.IsFinished() {
				state, err = s.RunContinuous()
				assert.NoError(err)
				numChoices := len(state.GetChoices())
				if numChoices > 0 {
					if assesChoiceCounts {
						assert.Equal(tC.choiceCounts[choiceIdx], numChoices, "Wrong number of choices encountered")
					}
					err = s.ChoseIndex(tC.choices[choiceIdx])
					assert.NoError(err)
					choiceIdx++
				}
			}
			assert.Equal(tC.expectedText, state.GetText())
		})
	}
}

func hello(x []any) any {
	return "External Hello " + x[0].(string)
}
