/*
▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄   ▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄
█▄▄▄▄▄▄▄█  █▄▄▄▄▄█   █▄▄▄▄▄█  █▄▄▄▄▄▄█  █▄▄▄▄▄▄█
   ▄▄▄     ▄▄▄  ▄▄▄  ▄▄▄▄▄▄▄  ▄▄▄       ▄▄▄  ▄▄▄
   ███     ███  ███  ██▄▄▄▄█  ███ █▀▀█  ███  ███
   ███     ███  ███  ███      ███ ▀███  ███  ███
   ███     ██▀▀▀▀██  ███      ██▀▀▀▀██  ██▀▀▀▀██
   ▀▀▀     ▀▀▀▀▀▀▀▀  ▀▀▀      ▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/digitallyserviced/tdfgo/tdf"
	"github.com/gookit/color"
	"github.com/gookit/goutil/errorx"
	"github.com/spf13/cobra"
)

// fontsCmd represents the fonts command
var fontsCmd = &cobra.Command{
	Use:   "fonts [-v] [-p] [pattern]",
	Short: "List fonts available",
	Long:  `List all available fonts found in the default directories`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 && len(args) != 1 {
			return errors.New("pattern provided can not contain spaces and must be a single word")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(errorx.Stacked(errorx.Newf("recovery: %v", err)))
			}
		}()
		var fonts []*tdf.FontInfo
		if len(args) == 1 {
			fonts = tdf.FindFonts(args[0])
		} else {
			fonts = tdf.GetFonts("")
		}

		for _, font := range fonts {
			tdff, err := tdf.LoadFont(font)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			verbose, err := cmd.InheritedFlags().GetBool("verbose")
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			checkChars, err := cmd.Flags().GetBool("checkChars")
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			preview, err := cmd.Flags().GetBool("preview")
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			txt, err := cmd.Flags().GetString("text")
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			if verbose {
				allFonts := tdff.GetAllFonts()
				fmt.Println(font.Path)
				if len(allFonts) > 0 {
					for _, f := range allFonts {
						inclFont := tdff.GetFont(f)
						SetFontOptions(inclFont)
						fmt.Print(inclFont.InfoString())
						if preview {
							fmt.Println("")
							inclFont.Print(txt)
						}
						if checkChars {
							characterStatus, missing := inclFont.HasChars(txt)
							if missing {
								fmt.Printf(color.Render(characterStatusFormat), characterStatus)
							}
						}
					}
					fmt.Println("")
				}
			} else {
				fmt.Print(tdff.String())
			}
			if preview && !verbose {
				allFonts := tdff.GetAllFonts()
				if len(allFonts) > 0 {
					for _, f := range allFonts {
						inclFont := tdff.GetFont(f)
						if checkChars {
							characterStatus, missing := inclFont.HasChars(txt)
							if missing {
								fmt.Printf(color.Render(characterStatusFormat), characterStatus)
							}
						}
						SetFontOptions(inclFont)
						if inclFont.Supported() {
							inclFont.Print(txt)
						} else {
							fmt.Printf("%s not supported because not Color type", f)
						}
					}
				}
			}
			tdff = nil
		}

	},
}

func init() {
	rootCmd.AddCommand(fontsCmd)

	fontsCmd.Example = color.Render(fontsExamples)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fontsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	fontsCmd.Flags().StringP("text", "t", "Preview", "Sample string to use for previewing fonts")
	fontsCmd.Flags().BoolP("preview", "p", false, "Output a preview for the fonts")
	fontsCmd.Flags().BoolP("checkChars", "X", false, "Check fonts to see if they are missing any characters in the defined preview text")
}
