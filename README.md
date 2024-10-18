# Freeze

<p>
  <a href="https://stuff.charm.sh/freeze/freeze-4k.png"><img src="https://github.com/charmbracelet/freeze/assets/25087/de76b799-fa67-4b5b-8da2-d990ca5b4e06" width="500" /></a><br>
  <a href="https://github.com/charmbracelet/freeze/releases"><img src="https://img.shields.io/github/release/charmbracelet/freeze.svg" alt="Latest Release"></a>
  <a href="https://github.com/charmbracelet/freeze/actions"><img src="https://github.com/charmbracelet/freeze/workflows/build/badge.svg" alt="Build Status"></a>
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
    <img alt="output of freeze command, Haskell code block" src="./test/golden/svg/shadow.svg" width="800" />
  </a>
</p>

### Generate an image of terminal output

You can use `freeze` to capture ANSI output of a terminal command with the
`--execute` flag.

```bash
freeze --execute "eza -lah"
```

<p align="center">
  <a href="https://github.com/charmbracelet/freeze/assets/42545625/aa5447ed-999a-4809-909d-67093d758f5a">
    <img alt="output of freeze command, ANSI" src="./test/golden/svg/eza.svg" width="800" />
  </a>
</p>

Freeze is also [super customizable](#customization) and ships with an [interactive TUI](#interactive-mode).

## Installation

```sh
# macOS or Linux
brew install charmbracelet/tap/freeze

# Arch Linux (btw)
yay -S freeze

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

- [`-b`](#background), [`--background`](#background): Apply a background fill.
- [`-c`](#configuration), [`--config`](#configuration): Base configuration file or template.
- [`-l`](#language), [`--language`](#language): Language to apply to code
- [`-m`](#margin), [`--margin`](#margin): Apply margin to the window.
- [`-o`](#output), [`--output`](#output): Output location for .svg, .png, .jpg.
- [`-p`](#padding), [`--padding`](#padding): Apply padding to the code.
- [`-r`](#border-radius), [`--border.radius`](#border-radius): Corner radius of window.
- [`-t`](#theme), [`--theme`](#theme): Theme to use for syntax highlighting.
- [`-w`](#window), [`--window`](#window): Display window controls.
- [`-H`](#height), [`--height`](#height): Height of terminal window.
- [`--border.width`](#border-width): Border width thickness.
- [`--border.color`](#border-width): Border color.
- [`--shadow.blur`](#shadow): Shadow Gaussian Blur.
- [`--shadow.x`](#shadow): Shadow offset x coordinate.
- [`--shadow.y`](#shadow): Shadow offset y coordinate.
- [`--font.family`](#font): Font family to use for code.
- [`--font.ligatures`](#font): Use ligatures in the font.
- [`--font.size`](#font): Font size to use for code.
- [`--font.file`](#font): File path to the font to use (embedded in the SVG).
- [`--line-height`](#font): Line height relative to font size.
- [`--show-line-numbers`](#line-numbers): Show line numbers.
- [`--lines`](#line-numbers): Lines to capture (start,end).

### Language

If possible, `freeze` auto-detects the language from the file name or analyzing
the file contents. Override this inference with the `--language` flag.

```bash
cat artichoke.hs | freeze --language haskell
```

<br />

<a href="https://github.com/charmbracelet/freeze/assets/42545625/d746f028-3d51-4bfd-ba81-94bbc47b3587">
  <img alt="output of freeze command, Haskell code block" src="./test/golden/svg/haskell.svg" width="600" />
</a>

### Theme

Change the color theme.

```bash
freeze artichoke.hs --theme dracula
```

<br />

<a href="https://github.com/charmbracelet/freeze/assets/42545625/c693dc80-89b8-43c6-a34e-2d49a332d0c7">
  <img alt="output of freeze command, Haskell code block with dracula theme" src="./test/golden/svg/dracula.svg" width="600" />
</a>

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

To use ligatures in the font, you can apply the `--font.ligatures` flag.

### Line Numbers

Show line numbers in the terminal window with the `--show-line-numbers` flag.

```bash
freeze artichoke.hs --show-line-numbers
```

To capture only a specific range of line numbers you can use the `--lines` flag.

```bash
freeze artichoke.hs --show-line-numbers --lines 2,3
```

### Border Radius

Add rounded corners to the terminal.

```bash
freeze artichoke.hs --border.radius 8
```

<br />

<a href="https://github.com/charmbracelet/freeze/assets/42545625/50932b10-56fd-4312-adbd-f64b36a15441">
  <img alt="code screenshot with corner radius of 8px" src="./test/golden/svg/border-radius.svg" width="600" />
</a>

### Window

Add window controls to the terminal, macOS-style.

```bash
freeze artichoke.hs --window
```

<a href="https://github.com/charmbracelet/freeze/assets/42545625/ba46bc4a-fb36-4718-88d4-f63b93343615">
  <img alt="output of freeze command, Haskell code block with window controls applied" src="./test/golden/svg/window.svg" width="600" />
</a>

### Background

Set the background color of the terminal window.

```bash
freeze artichoke.hs --background "#08163f"
```

### Height

Set the height of the terminal window.

```bash
freeze artichoke.hs --height 400
```

### Border Width

Add a border outline to the terminal window.

```bash
freeze artichoke.hs --border.width 1 --border.color "#515151" --border.radius 8
```

<br />

<a href="https://github.com/charmbracelet/freeze/assets/42545625/b9207976-50fe-479b-974b-e1bf2ad1684c">
  <img alt="output of freeze command, Haskell code block with border applied" src="./test/golden/svg/border-width.svg" width="600" />
</a>

### Padding

Add padding to the terminal window. You can provide 1, 2, or 4 values.

```bash
freeze main.go --padding 20          # all sides
freeze main.go --padding 20,40       # vertical, horizontal
freeze main.go --padding 20,60,20,40 # top, right, bottom, left
```

<br />

<a href="https://github.com/charmbracelet/freeze/assets/42545625/41da80bb-a6d8-402b-9f14-c73050720c0f">
  <img alt="output of freeze command, Haskell code block with padding applied" src="./test/golden/svg/padding.svg" width="600" />
</a>

### Margin

Add margin to the terminal window. You can provide 1, 2, or 4 values.

```bash
freeze main.go --margin 20          # all sides
freeze main.go --margin 20,40       # vertical, horizontal
freeze main.go --margin 20,60,20,40 # top, right, bottom, left
```

<br />

<a href="https://github.com/charmbracelet/freeze/assets/42545625/eca78a8e-2dbc-4cfa-81da-6c2194925238">
  <img alt="output of freeze command, Haskell code block with margin applied" src="./test/golden/svg/margin.svg" width="720" />
</a>

### Shadow

Add a shadow under the terminal window.

```bash
freeze artichoke.hs --shadow.blur 20 --shadow.x 0 --shadow.y 10
```

<br />

<a href="https://github.com/charmbracelet/freeze/assets/42545625/b9a8cc17-5ec6-4004-941b-5ae1b65d2b18">
  <img alt="output of freeze command, Haskell code block with a shadow" src="./test/golden/svg/shadow.svg" width="720" />
</a>

## Screenshot TUIs

Use `tmux capture-pane` to generate screenshots of TUIs.

Run your TUI in `tmux` and get it to the state you want to capture.
Next, use `capture-pane` to capture the pane and pipe that to freeze.

```bash
hx # in a separate pane
tmux capture-pane -pet 1 | freeze -c full
```

<img width="650px" src="./test/golden/svg/helix.svg" alt="helix captured with freeze">

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

## Contributing

See [contributing][contribute].

[contribute]: https://github.com/charmbracelet/freeze/contribute

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
