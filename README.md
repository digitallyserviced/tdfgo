## tdfgo

[TheDraw](https://en.wikipedia.org/wiki/TheDraw) font parser and console text renderer. This utility was inspired by the [`tdfiglet`](https://github.com/tat3r/tdfiglet) CLI utility. I wanted to also use the fonts in another utility I am working on. 

<img align="right" style="margin:6px;" width="500" src="./assets/tdfgo-preview.gif">

`tdfiglet` was missing some features and the ability to use more than a single font per file. I realized there was multiple fonts with different color schemes in some of them as well as font packs containing files with lots of different fonts within one font file.

> The CLI utility and library currently only support **COLOR** font types. **BLOCK** and **OUTLINE** will fail with a note.

*Features*

- Render __TheDraw__ fonts in the terminal
- Supports **COLOR** font types only as they are more plentiful and cooler looking
- Configure spacing between characters, the size a ` `(space) character takes up in text 
- Use a random font for content
- Supported characters and a missing character view available when listing fonts to find a font that contains the necessary characters.
- `print` - use arguments or stdin as the text content to render
- `watch` - specify a command to repeat and render the output (clocks, status, live banner)
- `clock` - alias to print `date +%H:%m:%S` to the console (tmux lock, console clock) `./tdfgo clock`
- `fonts` - dump the list of fonts available with information about the font and coloring previews `./tdfgo fonts -v`
- Preview fonts when dumping the list with customizable content `./tdfgo fonts -vp -t text`
- Shell completion of all flags and commands, completion of font file names found, completion of fonts found within a single font file 
<!-- ![tdfgo gif preview](assets/tdfgo-preview.gif) -->


### Rendering Text Options 

![help preview](assets/tdfgo-help.png)

#### Options

```bash
  -w, --columns int      Specify the amount of columns or width that the text will be rendered into (default 80)
  -f, --font string      Specify font to use for TEXT if multiple uses first (default "mindstax")
  -i, --fontIndex int    If multiple fonts per file specify the index of the font to use (default 0)
  -h, --help             help for tdfgo
  -j, --justify string   Specify the justification for rendered content, computes padding if necessary given current terminal width. {left, center, justify} (default "left")
  -m, --monochrome       Render text in monochrome by stripping color escape sequences. Will then be using the default foreground color
  -r, --random           Use a random font as the selected font when rendering text content
  -W, --spaceWidth int   Set the spacing for a space character in the text provided (default 3)
  -s, --spacing int      Override the fonts own specified spacing used between characters. (default -1)
  -v, --verbose          Print more information about fonts
```

### `Clock`

![clock preview](assets/tdfgo-clock.png)

Using the selected font and specified options output a clock to the terminal

```bash
# Center the clock using the yazoox font and use font #4 (0 based index)
tdfgo -j center -f yazoox -i 3 clock
```

### `Watch`

Run a command every # interval and use the output as the text rendered using a font

```
tdfgo watch [-i interval] CMD args... [flags]
```

#### Examples

```
# console clock
tdfgo watch -n 1s 'date "+%H:%m:%S"'
tdfgo clock

# print ONLINE or OFFLINE based off ping success
tdfgo watch 'ping -c1 -w 1 google.com>/dev/null 2>&1 && echo -n ONLINE || echo -n OFFLINE'

# print 1min load avg every 0.5 seconds
tdfgo watch 'cat /proc/loadavg | cut -d" " -f 1'

```

#### Options

```
  -h, --help                help for watch
  -n, --interval duration   Interval between executions of the specified command (default 1s)
```

### `Fonts`

![fonts preview](assets/tdfgo-fonts.png)

List all available fonts found in the default directories

```bash
tdfgo fonts [-v] [-p] [-t text] [-X] [pattern] 
```

#### Examples

```bash
# list fonts with impact in name
tdfgo fonts impact

# list fonts with much more information
tdfgo fonts -v

# list fonts with a preview output using default string "Preview"
tdfgo fonts -p

# list fonts with information about required characters specified by -t
# useful to find fonts having all the characters you need
tdfgo fonts -p -X -t "!@#$JDKALFK@{}"
```

#### Options

```
  -X, --checkChars    Check fonts to see if they are missing any characters in the defined preview text
  -h, --help          help for fonts
  -p, --preview       Output a preview for the fonts
  -t, --text string   Sample string to use for previewing fonts (default "Preview")
```

##### The font directories searched:

- `.` 
- `./fonts`
- `/usr/share/tdfgo/fonts`
- `/usr/local/share/tdfgo/fonts`
- `~/.config/tdfgo/fonts`
- Built-in (There are two fonts `mindstax.tdf` and `yazoox.tdf` embedded within the binary for use anywhere)

### Options

```
  -w, --columns int      Specify the amount of columns or width that the text will be rendered into (default 80)
  -f, --font string      Specify font to use for TEXT if multiple uses first (default "mindstax")
  -i, --fontIndex int    If multiple fonts per file specify the index of the font to use (default 0)
  -h, --help             help for tdfgo
  -j, --justify string   Specify the justification for rendered content, computes padding if necessary given current terminal width. {left, center, justify} (default "left")
  -m, --monochrome       Render text in monochrome by stripping color escape sequences. Will then be using the default foreground color
  -r, --random           Use a random font as the selected font when rendering text content
  -W, --spaceWidth int   Set the spacing for a space character in the text provided (default 3)
  -s, --spacing int      Override the fonts own specified spacing used between characters. (default -1)
  -v, --verbose          Print more information about fonts
```


