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
		expectedTags    []types.Tag
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
		{
			desc:            "Basic lists",
			inkJsonFilePath: "../../examples/list1.json",
			expectedText:    "The Kettle is cold\n",
		},
		{
			desc:            "Complex lists",
			inkJsonFilePath: "../../examples/list2.json",
			expectedText: `three,six
true
get the representation of a list object: one
two
get the value of a list element: 3
compare two list objects: false
one
Pre-Increment three
Post increment four
four
`,
		},
		{
			desc:            "List All Func",
			inkJsonFilePath: "../../examples/all.json",
			expectedText:    "a,f\n\na,e,b,f,c,g,d,h,e,i,j\n",
		},
		{
			desc:            "Range Func",
			inkJsonFilePath: "../../examples/range.json",
			expectedText:    "b,c,d\n",
		},
		{
			desc:            "List Invert Func",
			inkJsonFilePath: "../../examples/invert.json",
			expectedText:    "Pre: Smith,Jones\n\n\nPost: Carter,Braithwaite\n",
		},
		{
			desc: "Sequence",
			inkJsonFilePath: "../../examples/seq.json",
			expectedText: "\"Three!\"\n\"Two!\"\n\"One!\"\nThere was the white noise racket of an explosion.\nBut it was just static.\n",
		},
		{
			desc: "Tags",
			inkJsonFilePath: "../../examples/tag.json",
			expectedText: "Hello \n",
			expectedTags: []types.Tag{types.Tag("world "), types.Tag("another")},
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
			actualText, actualTags := state.GetTextAndTags()
			assert.Equal(tC.expectedText, actualText)
			assert.Equal(tC.expectedTags, actualTags)
		})
	}
}

func hello(x []any) any {
	return "External Hello " + x[0].(string)
}
