/*
▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄     ·  ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄     ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄     ▄▄▄▄▄▄▄▄▄▄▄▄
▄             ░ ░▒▓░  ▄              ▀▄    ▄▀         ░ ░▒▓░   ▄▀         ░ ░▒▓░   ▄▀            ▀▄
▒                ░▒▒  █               ▐▌  ▐▌             ░▒▒  ▐▌             ░▒▒  ▐▌              ▐▌
░                 ░▓  ▓      ▒▀▀░     ░█  ▓               ░▓  ▓               ░▓  ▓      █▀▀▒     ░█
█▄▄▄▄▄▄       ▄▄▄▄▄█  ▒      ░  ▒      █  ▒      ▄▄▄▄▄▄▄▄▄▄█  ▒      ▄▄▄▄▄▄▄▄▄▄█  ▒      █  ░      █
      █┼      █       ░      █  █▄     █  ░ ┼    ▒▄▄▄▄▄▄▄     ░┼     █   ▄▄▄▄▄▄▄  ░      █  █▄     █
  ░   █┼┼     █ · ·   █┼┼    █   ░    ┼░  █┼┼┼       ┼┼┼▒ ·   █┼     █   ▓┼┼  ┼░  █┼     █   ░    ┼░
 ░▒░· █┼┼    ┼█   ░   █┼┼┼┼ ┼█   ▒  ┼┼┼▒  █┼┼   ┼┼┼┼┼┼┼┼▒     █┼┼   ┼█▄▄▄▒┼┼┼ ┼░  █┼┼┼┼ ┼█   ▒┼  ┼┼▒
░▒▓▒░ █┼┼┼ ┼┼┼░ ·░░░  █┼┼┼┼┼┼█   ▓┼┼┼┼┼▓  █┼┼┼ ┼┼░▀▀▀▀▀▀▀ ·   █┼┼┼ ┼┼┼┼┼┼┼┼┼┼┼┼░  █┼┼┼┼┼┼█   ▓┼┼ ┼┼▓
 ░▒░  █┼┼┼┼┼┼┼▒ ░▒▒░  █┼┼┼┼┼┼█▄▄▄▀┼┼┼┼▐▌  █┼┼┼┼┼┼▒   ░▒▓▒░    ▐▌┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼▒  ▐▌┼┼┼┼┼█▄▄▄▀┼┼┼┼▐▌
  ░ · █┼┼┼┼┼┼┼▓ ░░░   █┼┼┼┼┼┼┼┼┼┼┼┼┼┼▄▀   █┼┼┼┼┼┼▓  ░▒▓█▓▒░    ▀▄┼┼┼┼┼┼┼┼┼┼┼┼┼┼▓   ▀▄┼┼┼┼┼┼┼┼┼┼┼┼▄▀
 ·    ▀▀▀▀▀▀▀▀▀  ░ ·  ▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀     ▀▀▀▀▀▀▀▀  ·░▒▓▒░       ▀▀▀▀▀▀▀▀▀▀▀▀▀▀      ▀▀▀▀▀▀▀▀▀▀▀▀
*/

package tdf

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gookit/goutil/fsutil"
)

var (
	fontsDirs    []string = []string{".", "./fonts", "/usr/share/tdfgo/fonts", "/usr/local/share/tdfgo/fonts", "~/.config/tdfgo/fonts"}
	BuiltInFonts []*FontInfo
)

type FontInfo struct {
	File, Path string
	FontDir    string
	BuiltIn    bool
}

func NewFontInfo(file, path string) *FontInfo {
	fi := &FontInfo{
		File: file,
		Path: path,
	}
	return fi
}
func FindFonts(font string) []*FontInfo {
	pat := fmt.Sprintf("*%s*.tdf", font)
	fonts := GetFonts(pat)

	return fonts
}
func FindFont(font string) *FontInfo {
	fonts := FindFonts(font)
	if len(fonts) == 0 {
		return nil
	}
	return fonts[0]
}

func GetFonts(pat string) []*FontInfo {
	fonts := make([]*FontInfo, 0)
	if pat == "" {
		pat = "*.tdf"
	}
	for _, fontPath := range fontsDirs {
		if fontPath == "BUILTIN" {
			fis := SearchBuiltinFonts(pat)
			// fmt.Println(fis)
			fonts = append(fonts, fis...)
		}

		fsutil.FindInDir(fsutil.Expand(fontPath), func(fPath string, fi os.FileInfo) error {
			if fPath == "BUILTIN" {
				return fmt.Errorf("Cannot search embedded FS")
			}
			fontP := path.Base(fPath)
			fInfo := NewFontInfo(fontP, fPath)
			fInfo.FontDir = fontPath
			fonts = append(fonts, fInfo)
			return nil
		}, func(fPath string, fi os.FileInfo) bool {
			if fPath == "BUILTIN" {
				return false
			}
			match, err := filepath.Match(pat, strings.ToLower(fi.Name()))
			if err != nil {
				return false
			}
			return match
		})
	}
	return fonts
}
