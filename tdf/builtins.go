package tdf

import (
	"embed"
	_ "embed"
	"log"
	"path"
	"path/filepath"
	"strings"
)

//go:embed fonts/mindstax.tdf
//go:embed fonts/yazoox.tdf
var BuiltinFontsFiles embed.FS

func SearchBuiltinFonts(s string) []*FontInfo {
	dirs, err := BuiltinFontsFiles.ReadDir("fonts")
	if err != nil {
		log.Fatal(err)
	}
	BuiltInFonts = make([]*FontInfo, 0)
	pat := "*"
	if s != "" {
		pat = s
	}
	for _, d := range dirs {
		match, err := filepath.Match(pat, strings.ToLower(d.Name()))
		// fmt.Println(match, pat, dirs)
		if err != nil {
			log.Fatalf("Error matching %s %v", d.Name(), err)
		}
		if !match {
			continue
		}

		fi := NewFontInfo(d.Name(), path.Join("fonts", d.Name()))
		fi.File = d.Name()
		fi.FontDir = "BUILTIN"
		fi.BuiltIn = true
		// dump.P(fi)
		BuiltInFonts = append(BuiltInFonts, fi)
	}
	return BuiltInFonts
}

func init() {
	fontsDirs = append(fontsDirs, "BUILTIN")
}
