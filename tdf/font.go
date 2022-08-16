/*
┌───────┐┐  ┌─────┐┐      ┌──────┐┐    ┌────┐┐    ┌───┐┐
│▓▓▓▓▓▓▓││  │▓▓▓▓▓└─┐┐  ┌─┘▓▓▓▓▓▓││  ┌─┘▓▓▓▓││  ┌─┘▓▓▓└─┐┐
└─┐░░───┘└  │░┌───┐░││  │░░──┬───┘└  │░░┌───┘└  │░░┌─┐░░││
  │▓▓││     │▓││ ┌┘▓││  │▓▓▓▓││      │▓▓││┬──┐  │▓▓│││▓▓││
  │░░││     │░│┌─┘░░││  │░░┌─┘└      │░░└─┘░┌│  │░░└─┘░░││
  │▒▒││     │▒└┘▒▒┌─┘└  │▒▒││        └─┐▒▒▒▒││  └─┐▒▒▒┌─┘
  └──┘└     └─────┘└    └──┘└          └────┘└    └───┘└
*/

package tdf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/ghostiam/binstruct"
	"github.com/gookit/color"

	// "github.com/gookit/goutil/arrutil"

	"golang.org/x/text/encoding/charmap"
)

type TheDrawFontAlignment int

const (
	AlignLeft TheDrawFontAlignment = iota + 1
	AlignCenter
	AlignRight
)

func (a TheDrawFontAlignment) String() string {
	return [...]string{"left", "center", "right"}[a-1]
}

var (
	Alignments       []TheDrawFontAlignment = []TheDrawFontAlignment{AlignLeft, AlignCenter, AlignRight}
	AlignmentStrings []string               = []string{"left", "center", "right"}
)
type TheDrawFontHeader struct {
}

type TheDrawFontMetaData struct {
	Namelen       uint8    `bin:"len:1"`
	FontName      string   `bin:"len:12"`
	_             byte     `bin:"len:4"`
	FontType      uint8    `bin:"len:1"` //outl block color
	FontSpacing   uint8    `bin:"len:1"`
	FontBlockSize uint16   `bin:"len:2"`
	CharList      []uint16 `bin:"len:94,[len:2]"`
	FontMaxHeight uint8    `bin:"-"`
	FontMaxWidth  uint8    `bin:"-"`
}

type TheDrawFontCell struct {
	Char      byte                `bin:"len:1"`
	Character string              `bin:"-"`
	Coloring  TheDrawFontColoring `bin:"-"`
}

type TheDrawFontCharacter struct {
	Offset  uint32            `bin:"-"`
	Width   uint8             `bin:"len:1"`
	Height  uint8             `bin:"len:1"`
	Content []TheDrawFontCell `bin:"-"`
	Colors  []string          `bin:"-"`
}

// type TheDrawFontGlyphs []TheDrawFontCharacter
type TheDrawFont struct {
	Offset      int                       `bin:"-"`
	CharOffsets int                       `bin:"-"`
	MetaData    TheDrawFontMetaData       // `bin:"len:233"`
	Glyphs      []TheDrawFontCharacter    `bin:"-"` // `bin:"len:fontBlockSize"` //
	Space       TheDrawFontCharacter      `bin:"-"` // `bin:"len:fontBlockSize"` //
	Colors      *ColorPairGroup           `bin:"-"`
	Options     *TheDrawFontStringOptions `bin:"-"`
}

func (tdf *TheDrawFont) FillCharContent(d *TheDrawFontCharacter) error {
	d.Content = make([]TheDrawFontCell, d.Width*d.Height)

	for i := 0; i < int(d.Height)*int(d.Width); i++ {
		cell := TheDrawFontCell{
			Char:      byte(' '),
			Character: fmt.Sprintf("%s", " "),
			Coloring: TheDrawFontColoring{
				BgFg: 0,
				Bg:   0,
				Fg:   0,
			},
		}
		if len(d.Content)-1 < i {
			continue
		}
		d.Content[i] = cell
	}
	return nil
}

func (tdf *TheDrawFont) NullTerminatedString(r binstruct.Reader, d *TheDrawFontCharacter) error {
	err := tdf.FillCharContent(d)
	if err != nil {
		return err
	}
	row, col := 0, 0
	for {
		readByte, err := r.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// panic(err)
				return nil
			}
			panic(err)
		}
		pos := row*int(d.Width) + col
		if len(d.Content)-1 < pos {
			return nil
		}
		switch {
		case bytes.Equal([]byte{readByte}, []byte{0x00}):
			return nil
		case bytes.Equal([]byte{readByte}, []byte{0x0d}):
			row = row + 1
			col = 0
			continue
		default:
			if readByte < 0x20 {
				readByte = uint8(' ')
			}
			d.Content[pos].Char = readByte
			d.Content[pos].Character = fmt.Sprintf("%c", d.Content[pos].FixEncoding(readByte))
			// d.Content[pos].Character = string(bytes.Runes([]byte{readByte}))
			color, err := r.ReadByte()
			if err != nil {
				if errors.Is(err, io.EOF) {
					// panic(err)
					return nil
				}
				panic(err)
			}
			// d.Content[pos].Coloring.BgFg = color
			d.Content[pos].Coloring.ParseColor(color, tdf)
		}
		col = col + 1
	}
}

func (tdfc *TheDrawFontCell) GetCharacter(readByte byte, second bool) {

}

func (tdfc *TheDrawFontCharacter) GetFor(row, col uint8) (string, string) {
	pos := (tdfc.Width*row + col)
	if len(tdfc.Content)-1 < int(pos) {
		return "", " "
	}
	fgcol := tdfc.Content[pos].Coloring.Fg
	bgcol := tdfc.Content[pos].Coloring.Bg
	// coloring := fmt.Sprintf("\x1b[%d;%dm", fgcolors[(fgcol+1)%16], bgcolors[(bgcol+1)%16])
	// colors := lo.Shuffle[uint8](fgcolors)
	coloring := fmt.Sprintf("\x1b[%d;%dm", fgcolors[fgcol], bgcolors[bgcol])
	// coloring := fmt.Sprintf("%d;%dm", fgcolors[fgcol], bgcolors[bgcol])
	if fgcol == 0 && bgcol == 0 {
		// coloring = ""
	}
	// coloring := fmt.Sprintf("%d;%dm", fgcol, bgcol)
	return coloring, tdfc.Content[pos].Character
}
func (tdf *TheDrawFont) ReadCharacters(r binstruct.Reader) error {
	if !tdf.Supported() {
		return fmt.Errorf("Only supporting colored fonts currently... %s not supported", tdf.MetaData.FontName)
	}
	tdf.Glyphs = make([]TheDrawFontCharacter, 94)
	for i := 0; i < 94; i++ {
		var offset uint16 = 0
		if num := tdf.MetaData.CharList[i]; num == 65535 {
			continue
		} else {
			offset = num
		}
		charOffset := tdf.CharOffsets + int(offset)
		var tfc TheDrawFontCharacter
		cpos, err := r.Seek(int64(charOffset), io.SeekStart)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println(err)
				panic(err)
			}
			fmt.Println(err)
			panic(err)
		}
		tfc.Offset = uint32(cpos)
		err = r.Unmarshal(&tfc)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println(err)
				panic(err)
			}
			fmt.Println(err)
			panic(err)
		}

		err = tdf.NullTerminatedString(r, &tfc)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println(err)
				panic(err)
			}
			fmt.Println(err)
			panic(err)
		}

		if tdf.MetaData.FontMaxWidth < tfc.Width {
			tdf.MetaData.FontMaxWidth = tfc.Width
		}
		if tdf.MetaData.FontMaxHeight < tfc.Height {
			tdf.MetaData.FontMaxHeight = tfc.Height
		}

		tdf.Glyphs[i] = tfc
	}
	tdf.Space = TheDrawFontCharacter{
		Offset: 0,
		Width:  tdf.Options.SpaceWidth,
		Height: tdf.MetaData.FontMaxHeight,
	}
	tdf.FillCharContent(&tdf.Space)
	return nil
}

var (
	lock       = sync.Mutex{}
	lastChange = time.Now()
)

// func ColorPairsString() string {
//
// }
// func ColorPairsString() string {
// 	previews := make([]string, 0)
// 	colors := make([]string, 0)
// 	count := 1
// 	for _, v := range ColorPairs {
// 		prev := v.Style.Sprint("▄▄")
// 		pair := fmt.Sprintf("pair %d (fg: %s bg: %s)", count, v.Fg.Sprint(v.Fg.Name()), v.Bg.ToFg().Sprint(v.Bg.ToFg().Name()))
// 		previews = append(previews, prev)
// 		colors = append(colors, pair)
// 		count += 1
// 	}
// 	return fmt.Sprintf("pairs:\n%s \npreview: %s", strings.Join(colors, "\n"), strings.Join(previews, ""))
// }

func (tdf *TheDrawFont) InfoString() string {
	info := fmt.Sprintf(`
<yellow>font:</> <cyan>%s</>
<yellow>char list:</>
<cyan>%s</>
<yellow>width:</> <cyan>%d</> <yellow>height:</> <cyan>%d</> <yellow>spacing:</> <cyan>%d</>
<yellow>preview:</> %s
	`, tdf.MetaData.FontName, tdf.Chars(true), tdf.MetaData.FontMaxWidth, tdf.MetaData.FontMaxHeight, tdf.MetaData.FontSpacing, tdf.Colors.String())
	return color.Render(info)
}

var TheDrawFontNotSupported = errors.New("Font type not supported. Only support Color types.")

func (tdf *TheDrawFont) Supported() bool {
	return tdf.MetaData.FontType == 2
}
func (tdf *TheDrawFont) Render(s string) string {
	tdfs := NewTheDrawFontStringFont(tdf)
	tdfs.Output = strings.Builder{}
	tdfs.Options.LineWidth = 0
	return tdfs.RenderString(s)
}

func (tdf *TheDrawFont) SetOptions(o func(*TheDrawFontStringOptions)) {
	o(tdf.Options)
}

func (tdf *TheDrawFont) Print(s string) {
	tdfs := NewTheDrawFontStringFont(tdf)
	tdfs.Output = strings.Builder{}
	tdfs.PrintString(s)
}

func (tdf *TheDrawFont) HasChars(chars string) (string, bool) {
	charlist := make([]string, len(chars))
	missing := false
	for i, char := range chars {
		if _, ok := tdf.HasChar(char); !ok {
			missing = true
			charlist[i] = color.Red.Render(fmt.Sprintf("%c", char))
		} else {
			charlist[i] = color.Green.Render(fmt.Sprintf("%c", char))
		}
	}
	return strings.Join(charlist, ""), missing
}

func (tdf *TheDrawFont) HasChar(char rune) (*TheDrawFontCharacter, bool) {
	charOrd := int(char) - 33
	if charOrd < 0 || charOrd > 93 {
		if char == ' ' {
			return &tdf.Space, true
		}
		return nil, false
	}
	charOffset := tdf.MetaData.CharList[charOrd]
	// fmt.Printf("%c, %d %d", char, charOrd, charOffset)
	return &tdf.Glyphs[charOrd], charOffset != 65535
}
func (tdfc *TheDrawFontCell) FixEncoding(b byte) rune {
	return charmap.CodePage437.DecodeByte(b)
}
func (tdf *TheDrawFont) Chars(inclMissing bool) string {
	chars := make([]string, 0)
	for ord, v := range tdf.MetaData.CharList {
		ordChar := rune(int(ord) + 33)
		char := fmt.Sprintf("%c", ordChar)
		if v == 65535 {
			if inclMissing {
				chars = append(chars, color.Red.Render(char))
			}
			continue
		}
		chars = append(chars, color.Green.Render(char))
	}
	return strings.Join(chars, "")
}

