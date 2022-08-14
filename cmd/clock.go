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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		watchCmd.Run(watchCmd, []string{"date", "+%H:%m:%S"})
		// watchCmd.Parent().Traverse([]string{"print", "watch", "date", "+%H:%m:%S"})
	},
}

func init() {
	// printCmd.AddCommand(clockCmd)
	rootCmd.AddCommand(clockCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
