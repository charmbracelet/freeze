# Freeze

<p>
  <a href="https://stuff.charm.sh/freeze/freeze-4k.png"><img src="https://github.com/charmbracelet/freeze/assets/25087/de76b799-fa67-4b5b-8da2-d990ca5b4e06" width="500" /></a><br>
  <a href="https://github.com/charmbracelet/freeze/releases"><img src="https://img.shields.io/github/release/charmbracelet/freeze.svg" alt="Latest Release"></a>
  <a href="https://pkg.go.dev/github.com/charmbracelet/freeze?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="Go Docs"></a>
</p>


Generate images of code and terminal output.

<p align="left">
  <a><img width="600" src="https://vhs.charm.sh/vhs-1C6z5SUKlTdqdj4KL1ADlH.gif" alt="Freeze code screenshot"></a>
</p>

## Examples

Freeze generates PNGs, SVGs, and WebPs of code and terminal output alike.

### Generate an image of code

```sh
freeze artichoke.hs -o artichoke.png
```

<p align="center">
  <a href="https://github.com/charmbracelet/freeze/assets/42545625/f15efdda-8e9b-4cb1-9e87-3d32b692eb7c">
    <img alt="output of freeze command, haskell code block" src="./test/golden/shadow.svg" width="800" />
  </a>
</p>

### Generate an image of terminal output

You can use `freeze` to capture ANSI output of a terminal command with the
`--execute` flag.

```bash
freeze --execute "eza -lah"
```
<p align="center">
  <img alt="output of freeze command, ANSI" src="./test/golden/eza.svg" width="800" />
</p>


Freeze is also [super customizable](#customization) and ships with an [interactive TUI](#interactive-mode).

## Installation

```sh
# macOS or Linux
brew install charmbracelet/tap/freeze

# Arch Linux (btw)
pacman -S freeze

# Nix
nix-env -iA nixpkgs.charm-freeze
```

Or, download it:

- [Packages][releases] are available in Debian and RPM formats
- [Binaries][releases] are available for Linux, macOS, and Windows

Or, just install it with `go`:

```sh
go install github.com/charmbracelet/freeze@latest
```

[releases]: https://github.com/charmbracelet/freeze/releases

## Customization

### Interactive mode

Freeze features a fully interactive mode for easy customization.

```bash
freeze --interactive
```

<img alt="freeze interactive mode" src="https://vhs.charm.sh/vhs-1AGhIlc2Mtn9Ltc8vPtaAP.gif" width="400" />

Settings are written to `$XDG_CONFIG/freeze/user.json` and can be accessed with
`freeze --config user`.

### Flags

Screenshots can be customized with `--flags` or [Configuration](#configuration) files.

> [!NOTE]
> You can view all freeze customization with `freeze --help`.

- [`-o`](#output), [`--output`](#output): Output location for .svg, .png, .jpg.
- [`-c`](#configuration), [`--config`](#configuration): Base configuration file or template.
- [`-t`](#theme), [`--theme`](#theme): Theme to use for syntax highlighting.
- [`-l`](#language), [`--language`](#language): Language to apply to code
- [`-w`](#window), [`--window`](#window): Display window controls.
- [`-m`](#margin), [`--margin`](#margin): Apply margin to the window.
- [`-p`](#padding), [`--padding`](#padding): Apply padding to the code.
- [`-r`](#border-radius), [`--border.radius`](#border-radius): Corner radius of window.
- [`--border.width`](#border-width): Border width thickness.
- [`--border.color`](#border-width): Border color.
- [`--shadow.blur`](#shadow): Shadow Gaussian Blur.
- [`--shadow.x`](#shadow): Shadow offset x coordinate.
- [`--shadow.y`](#shadow): Shadow offset y coordinate.
- [`--font.family`](#font): Font family to use for code.
- [`--font.size`](#font): Font size to use for code.
- [`--font.file`](#font): File path to the font to use (embedded in the SVG).
- [`--line-height`](#font): Line height relative to font size.

### Language

If possible, `freeze` auto-detects the language from the file name or analyzing
the file contents. Override this inference with the `--language` flag.

```bash
cat artichoke.hs | freeze --language haskell
```

<br />
<img alt="output of freeze command, haskell code block" src="./test/golden/haskell.svg" width="600" />

### Theme

Change the color theme.

```bash
freeze artichoke.hs --theme dracula
```

<br /><img alt="output of freeze command, haskell code block" src="./test/golden/dracula.svg" width="600" />

### Output

Change the output file location, defaults to `out.svg` or stdout if piped. This
value supports `.svg`, `.png`, `.webp`.

```bash
freeze main.go --output out.svg
freeze main.go --output out.png
freeze main.go --output out.webp

# or all of the above
freeze main.go --output out.{svg,png,webp}
```

### Font

Specify the font family, font size, and font line height of the output image.
Defaults to `JetBrains Mono`, `14`(px), `1.2`(em).

```bash
freeze artichoke.hs \
  --font.family "SF Mono" \
  --font.size 16 \
  --line-height 1.4
```

You can also embed a font file (in TTF, WOFF, or WOFF2 format) using the
`--font.file` flag.

### Border Radius

Add rounded corners to the terminal.

```bash
freeze artichoke.hs --border.radius 8
```

<br />
<img alt="code screenshot with corner radius of 8px" src="./test/golden/border-radius.svg" width="600" />

### Window

Add window controls to the terminal, macOS-style.

```bash
freeze artichoke.hs --window
```

<br />
<img alt="output of freeze command, haskell code block with window controls applied" src="./test/golden/window.svg" width="600" />

### Border Width

Add a border outline to the terminal window.

```bash
freeze artichoke.hs --border.width 1 --border.color "#515151" --border.radius 8
```

<br />
<img alt="output of freeze command, haskell code block with border applied" src="./test/golden/border-width.svg" width="600" />

### Padding

Add padding to the terminal window. You can provide 1, 2, or 4 values.

```bash
freeze main.go --padding 20          # all sides
freeze main.go --padding 20,40       # vertical, horizontal
freeze main.go --padding 20,60,20,40 # top, right, bottom, left
```

<br />
<img alt="output of freeze command, haskell code block with padding applied" src="./test/golden/padding.svg" width="600" />

### Margin

Add margin to the terminal window. You can provide 1, 2, or 4 values.

```bash
freeze main.go --margin 20          # all sides
freeze main.go --margin 20,40       # vertical, horizontal
freeze main.go --margin 20,60,20,40 # top, right, bottom, left
```

<br />
<img alt="output of freeze command, haskell code block with margin applied" src="./test/golden/margin.svg" width="720" />

### Shadow

Add a shadow under the terminal window.

```bash
freeze artichoke.hs --shadow.blur 20 --shadow.x 0 --shadow.y 10
```

<br />
<img alt="output of freeze command, haskell code block with a shadow" src="./test/golden/shadow.svg" width="720" />

## Configuration

Freeze also supports configuration via a JSON file which can be passed with the
`--config` / `-c` flag. In general, all `--flag` options map directly to keys
and values in the config file

There are also some default configurations built into `freeze` which can be passed by name.

- `base`: Simple screenshot of code.
- `full`: macOS-like screenshot.
- `user`: Uses `~/.config/freeze/user.json`.

If you use `--interactive` mode, a configuration file will be created for you at
`~/.config/freeze/user.json`. This will be the default configuration file used
in your screenshots.

```bash
freeze -c base main.go
freeze -c full main.go
freeze -c user main.go # alias for ~/.config/freeze/user.json
freeze -c ./custom.json main.go
```

Here's what an example configuration looks like:

```json
{
  "window": false,
  "border": {
    "radius": 0,
    "width": 0,
    "color": "#515151"
  },
  "shadow": false,
  "padding": [20, 40, 20, 20],
  "margin": "0",
  "font": {
    "family": "JetBrains Mono",
    "size": 14
  },
  "line_height": 1.2
}
```

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

- [Twitter](https://twitter.com/charmcli)
- [The Fediverse](https://mastodon.social/@charmcli)
- [Discord](https://charm.sh/chat)

## License

[MIT](https://github.com/charmbracelet/freeze/raw/main/LICENSE)

---

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/">
  <img
    alt="The Charm logo"
    width="400"
    src="https://stuff.charm.sh/charm-badge.jpg"
  />
</a>

Charm热爱开源 • Charm loves open source
