/*
┌─────┐ ┌───┐  ┌────┐ ┌───.  ┌────┐
└─┐░┌─┘ │░├┐└┐ │░┌─┬┘ │░┌_|_ │░┌┐░│
  │▒│   │▒├┘░│ │▒┌─┘  │▒└┘ │ │▒└┘▒│
  └─┘   └────┘ └─┘    └────┘ └────┘
*/

package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	// "github.com/digitallyserviced/tdfgo/tdf"
	// "github.com/gookit/goutil/dump"

	"github.com/spf13/cobra"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:     "print TEXT args...",
	Short:   "Print text arguments using the selected 'TheDrawFont' to stdout",
	Long:    `Render the TEXT arguments to stdout with the selected font and options`,
	Args:    cobra.ArbitraryArgs, //cobra.MinimumNArgs(1),
	Aliases: []string{"render"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// var stdinData []byte
		argTxt := strings.Join(args, " ")
		stdin := cmd.InOrStdin()
		if len(argTxt) == 0 && stdin != nil {
			stdinData, err := ioutil.ReadAll(stdin)
			if err == nil {
				if len(argTxt) == 0 {
					argTxt = fmt.Sprintf("%s", stdinData)
				}
			}
		}

		if useFont == nil {
			return fmt.Errorf("Could not find or load font to use.")
		}
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		if verbose {
			fmt.Println("")
			fmt.Print(useFont.InfoString())
			fmt.Println("")
		}
		useFont.Print(argTxt)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
	// Here you will define your flags and configuration settings.
}
