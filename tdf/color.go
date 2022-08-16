package tdf

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
)

var (
	fgcolors = []uint8{30, 34, 32, 36, 31, 35, 33, 37, 90, 94, 92, 96, 91, 95, 93, 97}
	bgcolors = []uint8{40, 44, 42, 46, 41, 45, 43, 47, 40, 44, 42, 46, 41, 45, 43, 47}
)

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

func (cpg *ColorPairGroup) AddPair(cp *ColorPair) {
	// cpg.Pairs = append(cpg.Pairs, cp)
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

