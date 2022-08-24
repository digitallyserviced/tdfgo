/*
 ▄▄▄▄▄▄▄▄▄▄▄    ▄▄▄▄▄▄▄       ▄▄▄▄▄▄▄▄       ▄▄▄▄▄▄         ▄▄▄▄▄
▐▓▓▓▓▓▓▓▓▓▓▓▌  ▐▓▓▓▓▓▓▓█▄    ▐▓▓▓▓▓▓▓▓▌    ▄█▓▓▓▓▓▓█▄     ▄█▓▓▓▓▓█▄
 ▀▀▀▐▒▒▒▌▀▀▀   ▐▒▒▒▌ ▀▒▒▒█   ▐▒▒▒▌▀▀▀▀    █▒▒▒▀  ▐▒▒▒▌   █▒▒▒▀▀▀▒▒▒█
    ▐░░░▌      ▐░░░▌  ▐░░░▌  ▐░░░░░░▌    ▐░░░▌    ▀▀▀   ▐░░░▌   ▐░░░▌
    ▐▓▓▓▌      ▐▓▓▓▌  ▐▓▓▓▌  ▐▓▓▓▓▓▓▌    ▐▓▓▓▌ ▐▓▓▓▓▓▌  ▐▓▓▓▌   ▐▓▓▓▌
    ▐▒▒▒▌      ▐▒▒▒▌ ▄▒▒▒█   ▐▒▒▒▌        █▒▒▒▄ ▀▀▒▒▒▌   █▒▒▒▄▄▄▒▒▒█
    ▐░░░▌      ▐░░░░░░░█▀    ▐░░░▌         ▀█░░░░░░█▀     ▀█░░░░░█▀
     ▀▀▀        ▀▀▀▀▀▀▀       ▀▀▀            ▀▀▀▀▀▀         ▀▀▀▀▀
*/
package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/digitallyserviced/tdfgo/tdf"
	cc "github.com/crinklywrappr/coloredcobra"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var useFont *tdf.TheDrawFont

var rootCmd = &cobra.Command{
	Use:   "tdfgo [-f font] [] TEXT",
	Short: "TheDraw font console printer",
	Long:  `For more information use tdfgo help`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()
		font, err := f.GetString("font")
		if err != nil {
			return err
		}
		fontIndex, err := f.GetInt("fontIndex")
		if err != nil {
			return err
		}
		random, err := f.GetBool("random")
		if err != nil {
			return err
		}
		var fontInfo *tdf.FontInfo
		if random {
			rand.Seed(time.Now().UnixNano())
			fontInfos := tdf.GetFonts("")
			fontInfo = lo.Sample(fontInfos)
			if fontInfo == nil {
				return fmt.Errorf(`Could not find any fonts to choose from randomly`)
			}
		} else {
			fontInfo = tdf.FindFont(font)
			if fontInfo == nil {
				return fmt.Errorf(`Could not find any font with the pattern "%s"`, font)
			}
		}
		tdff, err := tdf.LoadFont(fontInfo)
		if err != nil {
			panic(err)
		}
		allFonts := tdff.GetAllFonts()
		if len(allFonts) == 0 {
			return fmt.Errorf(`Could not find any supported fonts in "%s"`, fontInfo.Path)
		}
		if random {
			fontIndex = rand.Int() % len(allFonts)
		}
		if len(allFonts)-1 >= fontIndex {
			useFont = tdff.GetFont(allFonts[fontIndex])
		} else {
			return fmt.Errorf(`Could not find font supported fonts in "%s" #%d`, fontInfo.Path, fontIndex)
		}
		return nil
	},
}

func SetFontOptions(theFont *tdf.TheDrawFont) error {
	f := rootCmd.Flags()
	spacing, err := f.GetInt("spacing")
	if err != nil {
		return err
	}
	monochrome, err := f.GetBool("monochrome")
	if err != nil {
		return err
	}
	justify, err := f.GetString("justify")
	if err != nil {
		return err
	}
	spaceWidth, err := f.GetInt("spaceWidth")
	if err != nil {
		return err
	}
	columns, err := f.GetInt("columns")
	if err != nil {
		return err
	}
	if !f.Changed("columns") {
		w, h, err := term.GetSize(0)
		if err == nil {
			columns = w
			f.Set("columns", fmt.Sprintf("%d", columns))
		}
		_, _, _ = w, h, err
	}
	theFont.SetOptions(func(tdfso *tdf.TheDrawFontStringOptions) {
		if tdfso == nil {
			return
		}
		tdfso.SpaceWidth = uint8(spaceWidth)
		tdfso.Width = uint16(columns)

		switch justify {
		case "left":
			tdfso.Justify = tdf.AlignLeft
		case "center":
			tdfso.Justify = tdf.AlignCenter
		case "right":
			tdfso.Justify = tdf.AlignRight
		}
		if spacing != -1 {
			tdfso.Spacing = uint8(spacing)
		}
		tdfso.NoColor = monochrome
	})
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	tpl := usageTemplate
	rootCmd.SetUsageTemplate(tpl)
	cc.Init(&cc.Config{
		RootCmd:       rootCmd,
		Headings:      cc.HiCyan + cc.Bold + cc.Underline,
		Commands:      cc.HiYellow + cc.Bold,
		CmdShortDescr: cc.HiGreen,
		// Example:       ,
		ExecName:        cc.Bold,
		Flags:           cc.Yellow,
		FlagsDataType:   cc.Bold + cc.Magenta,
		FlagsDescr:      cc.HiBlue,
		Aliases:         cc.Bold + cc.Green,
		NoExtraNewlines: true,
	})

  rootCmd.Flags().SortFlags = false
	tdff, err := tdf.LoadBuiltinFont("yazoox", &tdf.BuiltinFontsFiles)
	if err != nil {
		panic(err)
	}
	allFonts := tdff.GetAllFonts()
	if len(allFonts) > 0 {
		// for _, f := range allFonts {
		n := time.Now()
		inclFont := tdff.GetFont(allFonts[n.Second()%len(allFonts)])
		hdr := inclFont.Render("tdfgo")
		tpl = rootCmd.UsageTemplate()
		// s := regexp.MustCompile(`FlagsUsage`)
		tpl = fmt.Sprintf("%s%s", hdr, tpl)
		rootCmd.SetUsageTemplate(tpl)
	}
	// tdfStr := tdf.NewTheDrawFontStringFontInfo(fontInfo)
	// hdr := tdfStr.RenderString("tdfgo")

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	setFlags()

	setCompletions()

	// rootCmd.Flags().SortFlags = true
}
func setFlags() {
	// rootCmd.Flags().FlagUsagesWrapped
	rootCmd.PersistentFlags().StringP("font", "f", "mindstax", "Specify font to use for TEXT if multiple uses first")
	rootCmd.MarkFlagRequired("font")
	rootCmd.PersistentFlags().StringP("justify", "j", "left", "Specify the justification for rendered content, computes padding if necessary given current terminal width. {left, center, justify}")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Print more information about fonts")
	rootCmd.PersistentFlags().BoolP("random", "r", false, "Use a random font as the selected font when rendering text content")
	rootCmd.PersistentFlags().BoolP("monochrome", "m", false, "Render text in monochrome by stripping color escape sequences. Will then be using the default foreground color")

	rootCmd.PersistentFlags().IntP("columns", "w", 80, "Specify the amount of columns or width that the text will be rendered into")
	rootCmd.PersistentFlags().IntP("fontIndex", "i", 0, "If multiple fonts per file specify the index of the font to use (default 0)")
	rootCmd.PersistentFlags().IntP("spacing", "s", -1, "Override the fonts own specified spacing used between characters.")
	rootCmd.PersistentFlags().IntP("spaceWidth", "W", 3, "Set the spacing for a space character in the text provided")

	// viper.SetEnvPrefix("")
	// viper.AutomaticEnv()
	//
	// viper.BindPFlag("columns", rootCmd.PersistentFlags().Lookup("columns"))
}

func setCompletions() {
	rootCmd.RegisterFlagCompletionFunc("fontIndex", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		font, _ := rootCmd.Flags().GetString("font")
		fonts := tdf.FindFont(font)
		if fonts == nil {
			return []string{}, cobra.ShellCompDirectiveError
		}
		tdff, _ := tdf.LoadFont(fonts)
		fontNames := make([]string, 0)
		for i, v := range tdff.GetAllFonts() {
			fontNames = append(fontNames, fmt.Sprintf("%d\t%s", i, v))
		}
		return fontNames, cobra.ShellCompDirectiveDefault
	})

	rootCmd.RegisterFlagCompletionFunc("justify", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"left", "center", "right"}, cobra.ShellCompDirectiveDefault
	})

	rootCmd.RegisterFlagCompletionFunc("font", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		fonts := tdf.FindFonts(toComplete)
		var fontNames []string
		fontNames = lo.Map(fonts, func(fi *tdf.FontInfo, i int) string {
			return fi.File[:len(fi.File)-4]
		})

		return fontNames, cobra.ShellCompDirectiveDefault
	})
}
