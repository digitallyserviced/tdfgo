package tdf

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/ghostiam/binstruct"
)

type TheDrawFontString struct {
	Text    string
	Font    *TheDrawFont
	Chars   []*TheDrawFontCharacter
	Output  strings.Builder
	Options *TheDrawFontStringOptions
}

type TheDrawFontStringOptions struct {
	Justify                              TheDrawFontAlignment
	Width, LineWidth                     uint16
	SpaceWidth, Padding, Height, Spacing uint8
	NoColor                              bool
}

func (tdfs *TheDrawFontString) RenderString(s string) string {
	chars := make([]*TheDrawFontCharacter, len(s))
	tdfs.Text = s
	tdfs.Options.Height = tdfs.Font.MetaData.FontMaxHeight
	tdfs.Output = strings.Builder{}
	tdfs.Output.WriteRune('\n')
	for _, char := range s {
		var glyph *TheDrawFontCharacter
		if g, ok := tdfs.Font.HasChar(char); !ok {
			continue
		} else {
			glyph = g
			chars = append(chars, g)
		}
		tdfs.Options.LineWidth = tdfs.Options.LineWidth + uint16(glyph.Width+tdfs.Options.Spacing)
	}
	tdfs.Chars = chars
	tdfs.Options.Padding = tdfs.GetPadding()

	for i := 0; i < int(tdfs.Font.MetaData.FontMaxHeight); i++ {
		tdfs.PrintRow(i)
	}
	tdfs.Output.WriteRune('\n')
	return tdfs.Output.String()
}

func (tdfs *TheDrawFontString) PrintString(s string) {
	fmt.Print(tdfs.RenderString(s))
}

func (tdfs *TheDrawFontString) LoadFont(s string) {
	fontname := fmt.Sprintf("./%s.tdf", s)
	file, err := os.Open(fontname)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// var fonthdr TheDrawFontMetaData
	var font TheDrawFont
	decoder := binstruct.NewDecoder(file, binary.LittleEndian)
	// decoder.SetDebug(true)
	err = decoder.Decode(&font)
	if err != nil {
		panic(err)
	}

	tdfs.Font = &font
}
func NewTheDrawFontStringOptions(font *TheDrawFont) *TheDrawFontStringOptions {
	tdfso := &TheDrawFontStringOptions{
		Justify:    0,
		Width:      80,
		Padding:    0,
		LineWidth:  0,
		Height:     font.MetaData.FontMaxHeight,
		Spacing:    font.MetaData.FontSpacing,
		NoColor:    false,
		SpaceWidth: 3,
	}
	return tdfso
}

func NewTheDrawFontStringFontInfo(fontInfo *FontInfo) *TheDrawFontString {

	tdfs := &TheDrawFontString{
		Output: strings.Builder{},
	}

	tdfs.LoadFontFontInfo(fontInfo)
	tdfs.Options = tdfs.Font.Options

	return tdfs
}

func NewTheDrawFontStringFont(font *TheDrawFont) *TheDrawFontString {
	tdfs := &TheDrawFontString{
		Font:   font,
		Chars:  make([]*TheDrawFontCharacter, 0),
		Output: strings.Builder{},
	}

	tdfs.Font = font
	tdfs.Options = tdfs.Font.Options

	return tdfs
}

func NewTheDrawFontString(font string) *TheDrawFontString {

	tdfs := &TheDrawFontString{
		Output: strings.Builder{},
	}

	tdfs.LoadFont(font)
	tdfs.Options = tdfs.Font.Options

	return tdfs
}

func (tdfs *TheDrawFontString) GetPadding() uint8 {
	if tdfs.Options.Justify == AlignCenter {
		center := math.Floor(float64(tdfs.Options.Width) / 2)
		halfWidth := math.Floor(float64(tdfs.Options.LineWidth) / 2)
		return uint8(center - halfWidth)
	}
	if tdfs.Options.Justify == AlignRight {
		pad := tdfs.Options.Width - tdfs.Options.LineWidth - 1
		return uint8(pad)
	}
	// if tdfs.Options.Justify == 0 {
	// 	return 0
	// }
	//
	// return 0
	return 0
}
func (tdfs *TheDrawFontString) PrintRow(row int) {
	for i := 0; i < int(tdfs.Options.Padding); i++ {
		tdfs.Output.WriteString(" ")
	}

	for ci, v := range tdfs.Chars {
		if v == nil {
			// for s := 0; s < int(tdfs.Options.SpaceWidth); s++ {
			// 	tdfs.Output.WriteString("W")
			// }
			// tdfs.Output.WriteString("!")
			continue
		}
		lastcolor := ""
		for col := 0; col < int(v.Width); col++ {
			coloring, char := v.GetFor(uint8(row), uint8(col))
			if (lastcolor == coloring && col != 0) || tdfs.Options.NoColor {
				coloring = ""
			}
			tdfs.Output.WriteString(fmt.Sprintf("%s%s", coloring, char))
			lastcolor = coloring
		}
		tdfs.Output.WriteString("\x1b[0m")
		if ci == len(tdfs.Chars)-1 {
			continue
		}
		for s := 0; s < int(tdfs.Options.Spacing); s++ {
			tdfs.Output.WriteString(" ")
		}
	}
	tdfs.Output.WriteRune('\n')
}

func (tdfs *TheDrawFontString) LoadFontFontInfo(fi *FontInfo) {
	file, err := os.Open(fi.Path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// var fonthdr TheDrawFontMetaData
	var font TheDrawFont
	decoder := binstruct.NewDecoder(file, binary.LittleEndian)
	// decoder.SetDebug(true)
	err = decoder.Decode(&font)
	if err != nil {
		panic(err)
	}

	tdfs.Font = &font
}
