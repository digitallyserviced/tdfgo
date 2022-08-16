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
	"embed"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ghostiam/binstruct"
	"github.com/gookit/color"
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

type TheDrawFontFile struct {
	FontInfo *FontInfo     `bin:"-"`
	FileSize int           `bin:"-"`
	Fonts    []TheDrawFont `bin:"ReadFonts"`
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
