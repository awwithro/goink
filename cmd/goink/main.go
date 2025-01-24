package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/awwithro/goink/pkg/parser"
	"github.com/awwithro/goink/pkg/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debug bool
var defaultLogLevel = log.WarnLevel

var rootCmd = &cobra.Command{
	Use: "goink <ink_json>",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetLevel(defaultLogLevel)
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		if js, err := os.ReadFile(args[0]); err != nil {
			return err
		} else {
			ink := parser.Parse(js)
			s := runtime.NewStory(ink)
			runStory(s)
		}
		return nil
	},
}

func runStory(s runtime.Story) {
	log.Debug("Starting")
	reader := bufio.NewReader(os.Stdin)
	s.Start()

	for !s.IsFinished() {
		state, err := s.RunContinuous()
		if err != nil {
			log.Error(err)
		}
		fmt.Print(state.GetText())
		if len(state.CurrentChoices) > 0 {
			for x, choice := range state.CurrentChoices {
				fmt.Printf("%d: %s\n", x, choice.ChoiceText())
			}
			text, _ := reader.ReadString('\n')
			c, _ := strconv.Atoi(strings.TrimSpace(text))
			err = s.ChoseIndex(c)
			if err != nil {
				log.Error(err)
			}
		}
	}
	fmt.Println("THE END")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "set debug logging")
}
