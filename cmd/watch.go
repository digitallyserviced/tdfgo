/*
   ┌───────┐  ┌──────┐   ┌───────┐  ┌───────┐  ┌───────┐
  ═╘═╕∙ ·╒═╛ ═│∙  ╒═╕└┐ ═│∙  ╒═╕∙│ ═│∙  ╒═══╛ ═│∙  ╒═╕∙│
  ░░▒│   │▒░ ░│   │█│ │  │   └┐└─┘  │   │┌──┐ ░│   │█│ │
  ░░▒│   │▒░ ▒│   │▓│ │ ░│   ┌┘░▒▓ ░│   │└┐·│ ▒│   │▓│ │
  ░░▒│   │▒░ ▓│   └─┘ │ ▒│   │░░▒▓ ▒│   │░│ │ ▓│   └─┘ │
  ═══│∙ ·│══ ═│∙     ┌┘ ═│∙  │════ ═│∙  ╘═╛∙│ ═│∙     ∙│
     ╘═══╛    ╘══════╛   ╘═══╛      ╘═══════╛  ╘═══════╛
*/
package cmd

import (
	// "fmt"

	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/go-cmd/cmd"
	"github.com/gookit/color"
	"github.com/gookit/goutil/dump"

	// sw "github.com/mattn/go-shellwords"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var Model = struct {
	app      *tview.Application
	flexV    *tview.Flex
	textView *tview.TextView
	scr      tcell.Screen
}{
	app:      nil,
	flexV:    nil,
	textView: nil,
	scr:      nil,
}

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch [-i interval] CMD args...",
	Short: "Repeat a command and use the output as text to print using defined font",
	Long:  `Run a command every # interval and use the output as the text rendered using a font`,
	Run: func(cmd *cobra.Command, args []string) {
		StartScreen()
		interval, err := cmd.Flags().GetDuration("interval")
		if err != nil {
			panic(err)
		}

		if useFont != nil {
			err = SetFontOptions(useFont)
			if err != nil {
				panic(err)
			}
			go watchLoop(interval, NewCmd(args))
			if err := Model.app.SetRoot(Model.flexV, true).Run(); err != nil {
				panic(err)
			}
		}
	},
}

func watchLoop(interval time.Duration, watchCmd *cmd.Cmd) {
	for {
		time.Sleep(interval)
		output := RunCmd(watchCmd.Clone())
		output = strings.TrimSpace(output)
		Model.app.QueueUpdateDraw(func() {
			w := Model.textView.BatchWriter()
			w.Clear()
			w.Write([]byte(tview.TranslateANSI(useFont.Render(output))))
			w.Close()
		})
	}
}

func init() {
	watchCmd.Example = color.Render(watchExamples)
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().DurationP("interval", "n", 1*time.Second, "Interval between executions of the specified command")
}

func StartScreen() {
	Model.app = tview.NewApplication()
	scr, err := tcell.NewTerminfoScreen()
	if err != nil {
		panic(err)
	}
	Model.scr = scr
	Model.scr.Init()
	Model.app.SetScreen(Model.scr)

	Model.flexV = tview.NewFlex()
	Model.flexV.SetDirection(tview.FlexRow)

	Model.textView = tview.NewTextView()
	Model.textView.SetDynamicColors(true)

	Model.flexV.AddItem(nil, 0, 1, false)
	Model.flexV.AddItem(Model.textView, 0, 4, false)
	Model.flexV.AddItem(nil, 0, 1, false)
}

func RunCmd(cmd *cmd.Cmd) string {
	status := cmd.Start()
	stat := <-status
	return strings.Join(stat.Stdout, " ")
}

func NewCmd(args []string) *cmd.Cmd {
	cm := "/bin/bash"
	arg := make([]string, 0)
	arg = append(arg, "-c")
	arg = append(arg, fmt.Sprintf(`%s`, strings.Join(args, " ")))
	cmdOptions := cmd.Options{
		Buffered:       true,
		LineBufferSize: 1,
		Streaming:      false,
	}

	dump.P(cm, arg)
	// return cmd.NewCmdOptions(cmdOptions, exe, args...)
	return cmd.NewCmdOptions(cmdOptions, cm, arg...)
}
