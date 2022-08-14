/*
Copyright © 2020 Steve Francia <spf@spf13.com>
This file is part of CLI application foo.
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// junkCmd represents the junk command
var junkCmd = &cobra.Command{
	Use:    "junk",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// fonts := tdf.FindFonts(strings.Join(args, ""))
		// var fontNames []string
		// fontNames = lo.Map[*tdf.FontInfo, string](fonts, func(fi *tdf.FontInfo, i int) string {
		// 	return fi.File[:len(fi.File)-4]
		// })
		// fmt.Println(strings.Join(fontNames, " "))

		err := doc.GenMarkdownTree(rootCmd, "./")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(junkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// junkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// junkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
