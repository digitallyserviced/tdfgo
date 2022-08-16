/*
                ▄▀▄   ▄▄▄▄   ▄       ▄
    ▄        ▄▀█ █   █    ▀▀▀█   ▄▄▄▀█▀▄
  ▄▀ █        █  █   ▀█ █▀▄▄▀  ▄▀ █▄██ █
▄▀   ▀▀▀▄    ▄█  █    █ █▄ ▄    █  ██  █     ▄▀▀▄
 ▀█  █▀▀   ▄▀ ▀  █    █ ▄▄▀▄▀   █  ██  █   ▄▀█▄▄ ▀▄
  █  █    █  ▄█  █    █ █ ▀     █  ▀█  █ ▄▀ ▄▀  ▀▄█▀▄
  █  █    █  ██  █    █ █        ▀▄ █  █ ██ █   ▄▀█▄▀
  █  ▀▄▀▄ █  ▀█  ▀▄   ███          ▀█  █  ▀▄█▀▄▀█▄▀
  ▀▀▄ ▄▀  ▀▀▄ ▄▀▄ ▄▀ ▄▀█           ▄▀  █    ▀▄█▄▀
     ▀       ▀   ▀   ▀▀          ▄▀  ▄▀       ▀
                                ▀▄ ▄▀
                                  ▀
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// clockCmd represents the clock command
var clockCmd = &cobra.Command{
	Use:   "clock",
  Short: "Run the watch command with the time format H:m:S",
	Run: func(cmd *cobra.Command, args []string) {
    cmd.Parent().PersistentFlags().Set("justify", "center")
    if !cmd.Parent().PersistentFlags().Changed("font") {
      cmd.Parent().PersistentFlags().Set("font", "yazoox")
    }
    cmd.Parent().PersistentPreRunE(cmd.Parent(), []string{"date", "+%H:%m:%S"})
		watchCmd.Run(watchCmd, []string{"date", "+%H:%m:%S"})
	},
}

func init() {
	rootCmd.AddCommand(clockCmd)
}
