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
	"embed"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"math"
	"os"
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

var (
	fgcolors = []uint8{30, 34, 32, 36, 31, 35, 33, 37, 90, 94, 92, 96, 91, 95, 93, 97}
	bgcolors = []uint8{40, 44, 42, 46, 41, 45, 43, 47, 40, 44, 42, 46, 41, 45, 43, 47}
)

func init() {
	// arrutil.Uniong
	// rand.Seed(time.Now().Unix())
	// fgcolors = lo.Shuffle(fgcolors)
}

func (cpg *ColorPairGroup) Detail() string {
	// previews := make([]string, 0)
	colors := make([]string, 0)
	count := 1
	for _, v := range cpg.Pairs {
		// prev := v.Style.Sprint("▄▄")
		// previews = append(previews, prev)
		pair := fmt.Sprintf("pair %d (fg: %s bg: %s)", count, v.Fg.Sprint(v.Fg.Name()), v.Bg.ToFg().Sprint(v.Bg.ToFg().Name()))
		colors = append(colors, pair)
		count += 1
	}
	return fmt.Sprintf("%s", strings.Join(colors, "\n"))
	// return fmt.Sprintf("pairs:\n%s \npreview: %s", strings.Join(colors, "\n"), strings.Join(previews, ""))
}
func (cpg *ColorPairGroup) String() string {
	previews := make([]string, 0)
	// colors := make([]string, 0)
	count := 1
	for _, v := range cpg.Pairs {
		prev := v.Style.Sprint("▄▄")
		previews = append(previews, prev)
		// pair := fmt.Sprintf("pair %d (fg: %s bg: %s)", count, v.Fg.Sprint(v.Fg.Name()), v.Bg.ToFg().Sprint(v.Bg.ToFg().Name()))
		// colors = append(colors, pair)
		count += 1
	}
	return fmt.Sprintf("%s", strings.Join(previews, ""))
	// return fmt.Sprintf("pairs:\n%s \npreview: %s", strings.Join(colors, "\n"), strings.Join(previews, ""))
}

func (cpg *ColorPairGroup) AddPair(cp *ColorPair) {
	// cpg.Pairs = append(cpg.Pairs, cp)
}

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

type ColorPairGroup struct {
	Pairs map[string]ColorPair
}

type ColorPair struct {
	color.Style
	ColorPair4
}
type ColorPair4 struct {
	Fg color.Color
	Bg color.Color
}

type ColorPairRGB struct {
	fg color.RGBColor
	bg color.RGBColor
}

func NewColorPair8(fg, bg uint8) *ColorPair {
	cp4 := &ColorPair4{
		Fg: color.Bit4(fg),
		Bg: color.Bit4(bg),
	}
	sty := color.New(cp4.Fg, cp4.Bg)
	cp := &ColorPair{
		Style:      sty,
		ColorPair4: *cp4,
	}
	return cp
}

// func NewColorPairRGB(fg, bg string) *ColorPair {
// 	cp := &ColorPair{
// 		fg: color.HEX(fg, false),
// 		bg: color.HEX(bg, true),
// 	}
// 	cp.style = color.NewRGBStyle(cp.fg, cp.bg)
// 	return cp
// }

type TheDrawFontColoring struct {
	BgFg uint8      `bin:"len:1,ParseColor"`
	Bg   uint8      `bin:"-"`
	Fg   uint8      `bin:"-"`
	Pair *ColorPair `bin:"-"`
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
type TheDrawFontFile struct {
	FontInfo *FontInfo     `bin:"-"`
	FileSize int           `bin:"-"`
	Fonts    []TheDrawFont `bin:"ReadFonts"`
}

type TheDrawFont struct {
	Offset      int                       `bin:"-"`
	CharOffsets int                       `bin:"-"`
	MetaData    TheDrawFontMetaData       // `bin:"len:233"`
	Glyphs      []TheDrawFontCharacter    `bin:"-"` // `bin:"len:fontBlockSize"` //
	Space       TheDrawFontCharacter      `bin:"-"` // `bin:"len:fontBlockSize"` //
	Colors      *ColorPairGroup           `bin:"-"`
	Options     *TheDrawFontStringOptions `bin:"-"`
}

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

func (tdfc *TheDrawFontColoring) MakeColorPair(tdf *TheDrawFont) error {
	key := fmt.Sprintf("%d%d", fgcolors[tdfc.Fg], bgcolors[tdfc.Bg])
	val, ok := tdf.Colors.Pairs[key]
	if !ok {
		cp := NewColorPair8(fgcolors[tdfc.Fg], bgcolors[tdfc.Bg])
		val = *cp
		tdf.Colors.Pairs[key] = val
	}
	tdfc.Pair = &val
	return nil
}

func (tdfc *TheDrawFontColoring) ParseColor(col uint8, tdf *TheDrawFont) error {
	tdfc.BgFg = col
	fgcol := col & 0x0f
	bgcol := (col & 0xf0) >> 4
	tdfc.Bg = bgcol
	tdfc.Fg = fgcol
	tdfc.MakeColorPair(tdf)
	return nil
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

func NewTheDrawFontFile(fi *FontInfo) *TheDrawFontFile {
	tdff := &TheDrawFontFile{
		FontInfo: fi,
		Fonts:    make([]TheDrawFont, 0),
	}

	return tdff
}

func (tdff *TheDrawFontFile) ParseFont(r binstruct.Reader, tdf *TheDrawFont) error {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("Recovered ParseFont error: %v", err))
		}
	}()

	// tdf.Colors.Pairs = make(map[string]ColorPair)
	tdf.Colors = &ColorPairGroup{
		Pairs: make(map[string]ColorPair),
	}
	err := r.Unmarshal(tdf)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		fmt.Println(err)
		panic(err)
	}
	tdf.Options = NewTheDrawFontStringOptions(tdf)
	loc, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		fmt.Println(err)
		panic(err)
	}
	tdf.CharOffsets = int(loc)
	if !tdf.Supported() {
		return TheDrawFontNotSupported
	}

	err = tdf.ReadCharacters(r)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		fmt.Println(err)
		panic(err)
	}
	// dump.P(TDFColoring)
	// fmt.Println(TDFColoring.Preview(tdf.MetaData.FontName))
	return nil
}
func (tdff *TheDrawFontFile) ReadFonts(r binstruct.Reader) error {
	offset := 24
	_, err := r.Seek(int64(offset), io.SeekStart)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		fmt.Println(err)
		panic(err)
	}
	for {
		cpos, err := r.Seek(0, io.SeekCurrent)
		var tdf TheDrawFont
		tdf.Offset = int(cpos)
		err = tdff.ParseFont(r, &tdf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println(err)
				return nil
			}
			if !errors.Is(err, TheDrawFontNotSupported) {
				return nil
			}
		}
		if tdf.CharOffsets == 0 || tdf.MetaData.FontBlockSize == 0 {
			return nil
		}
		tdff.Fonts = append(tdff.Fonts, tdf)
		// dump.P(tdf.Offset, tdf.CharOffsets, tdf.MetaData.FontBlockSize)
		_, err = r.Seek(int64(tdf.CharOffsets)+int64(tdf.MetaData.FontBlockSize)+4, io.SeekStart)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println(err)
				return nil
			}
			fmt.Println(err)
			panic(err)
		}
	}
}

func ParseFont(file fs.File) (*TheDrawFontFile, error) {
	// var fonthdr TheDrawFontMetaData
	var fontFile TheDrawFontFile
	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fontFile.FileSize = int(fileStat.Size())
	bits, _ := ioutil.ReadAll(file)
	reader := binstruct.NewReaderFromBytes(bits, binary.LittleEndian, false)
	err = reader.Unmarshal(&fontFile)
	if err != nil {
		return nil, err
	}
	// tdff.FileSize = int(fileStat.Size())
	// decoder := binstruct.NewDecoder(file.(io.ReadSeeker), binary.LittleEndian)
	// decoder.SetDebug(true)
	// err = decoder.Decode(&fontFile)
	// if err != nil {
	// 	return nil, err
	// }
	return &fontFile, nil
}

func (tdff *TheDrawFontFile) String() string {
	// dump.P(tdff.FontInfo)
	str := fmt.Sprintf(`
<yellow>file:</> <green>%s</>
<yellow>embedded fonts:</> <green>%d</>
<yellow>font names:</> <cyan>%s</>
`, tdff.FontInfo.Path, len(tdff.Fonts), strings.Join(tdff.GetAllFonts(), " | "))
	return color.Render(str)
}

func (tdff *TheDrawFontFile) GetFont(f string) *TheDrawFont {
	for _, v := range tdff.Fonts {
		if f == v.MetaData.FontName {
			return &v
		}
	}
	return nil
}

func (tdff *TheDrawFontFile) GetAllFonts() []string {
	fonts := make([]string, 0)
	for _, v := range tdff.Fonts {
		if v.MetaData.FontType != 2 {
			continue
		}
		fonts = append(fonts, v.MetaData.FontName)
	}
	return fonts
}

func LoadBuiltinFont(name string, efs *embed.FS) (*TheDrawFontFile, error) {
	if efs != nil {
		file, err := efs.Open(fmt.Sprintf("fonts/%s.tdf", name))
		if err != nil {
			return nil, err
		}
		defer file.Close()
		fontFile, err := ParseFont(file)
		if err != nil {
			return nil, err
		}
		return fontFile, nil
	}
	return nil, errors.New("No proper embed FS provided")
}
func LoadFont(fi *FontInfo) (*TheDrawFontFile, error) {
	if fi != nil {
		var file fs.File
		var err error
		if fi.BuiltIn {
			file, err = BuiltinFontsFiles.Open(fmt.Sprintf("%s", fi.Path))
		} else {
			file, err = os.Open(fi.Path)
		}
		if err != nil {
			return nil, err
		}
		defer file.Close()
		fontFile, err := ParseFont(file)
		if err != nil {
			return nil, err
		}
		fontFile.FontInfo = fi
		return fontFile, nil
	}
	return nil, errors.New("No font info for file")
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
