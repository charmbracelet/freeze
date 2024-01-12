# Freeze

<p>
  <a href="https://github.com/charmbracelet/freeze/releases"><img src="https://img.shields.io/github/release/charmbracelet/freeze.svg" alt="Latest Release"></a>
  <a href="https://pkg.go.dev/github.com/charmbracelet/freeze?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="Go Docs"></a>
</p>

Capture and share your code on the command line.

<p align="center">
  <img alt="output of freeze command, haskell code block" src="./examples/shadow.svg" width="800" />
</p>

## Tutorial

To generate an image of your code, run:

<img alt="a terminal running: freeze -wsm 20 -r 8 artichoke.hs" src="https://vhs.charm.sh/vhs-3yL9rCVdQNUjYwVBqzW7Kc.gif" width="600" />

```sh
freeze main.go -o out.svg
```

Your image file will live in `out.svg`.

## Installation

```sh
# macOS or Linux
brew install freeze

# Arch Linux (btw)
pacman -S freeze

# Nix
nix-env -iA nixpkgs.freeze
```

Or, download it:

* [Packages][releases] are available in Debian and RPM formats
* [Binaries][releases] are available for Linux, macOS, and Windows

Or, just install it with `go`:

```sh
go install github.com/charmbracelet/freeze@latest
```

[releases]: https://github.com/charmbracelet/freeze/releases

## Customization

Screenshots can be customized with `--flags`.

> [!NOTE]
> You can view all freeze customization with `freeze --help`.

There are many different configuration options:

* [`-l`, `--language`](#language): code language.
* [`-t`, `--theme`](#theme): theme to use.
* [`-o`, `--output`](#output): output file.
* [`-r`, `--radius`](#corner-radius): add corner radius.
* [`-w`, `--window`](#window): show window controls.
* [`-b`, `--border`](#border): add a border to the window.
* [`-p`, `--padding`](#padding): terminal padding.
* [`-m`, `--margin`](#margin): window margin.
* [`-s`, `--shadow`](#shadow): add a shadow to the window.
* [`--font.family`](#font): font family.
* [`--font.size`](#font): font size.
* [`--line-height`](#font): line height.


### Language

If possible, `freeze` auto-detects the language from the file name or analyzing
the file contents. Override this inference with the `--language` flag.

```bash
cat artichoke.hs | freeze --language haskell
```

<br />
<img alt="output of freeze command, haskell code block" src="https://github.com/charmbracelet/freeze/assets/42545625/f3fb212f-6629-4253-9c13-105055a4b6e8" width="600" />

### Theme

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

### Corner Radius

Add rounded corners to the terminal.

```bash
freeze artichoke.hs --radius 8
```

<br />
<img alt="code screenshot with corner radius of 8px" src="./examples/corner-radius.svg" width="600" />

### Window

Add window controls to the terminal, macOS-style.

```bash
freeze artichoke.hs --window
```

<br />
<img alt="output of freeze command, haskell code block with window controls applied" src="./examples/window.svg" width="600" />

### Border

Add a border to the terminal window.

```bash
freeze artichoke.hs --border --radius 8
```

<br />
<img alt="output of freeze command, haskell code block with border applied" src="./examples/border.svg" width="600" />

### Padding

Add padding to the terminal window.

Given 1 value, padding applies to all sides. Given 2 values, the first value
applies to vertical padding and the second to horizontal padding. Given 4
values, the padding applies to the top, right, bottom, left sides respectively.

```bash
freeze main.go --padding 20          # all sides
freeze main.go --padding 20,40       # vertical, horizontal
freeze main.go --padding 20,60,20,40 # top, right, bottom, left
```

<br />
<img alt="output of freeze command, haskell code block with padding applied" src="./examples/padding.svg" width="600" />


### Margin

Add margin to the terminal window.

Given 1 value, margin applies to all sides. Given 2 values, the first value
applies to vertical margin and the second to horizontal margin. Given 4
values, the margin applies to the top, right, bottom, left sides respectively.

```bash
freeze main.go --margin 20          # all sides
freeze main.go --margin 20,40       # vertical, horizontal
freeze main.go --margin 20,60,20,40 # top, right, bottom, left
```

<br />
<img alt="output of freeze command, haskell code block with margin applied" src="./examples/margin.svg" width="600" />

### Shadow

Add a shadow under the terminal window.

```bash
freeze artichoke.hs --shadow --margin 20
```

<br />
<img alt="output of freeze command, haskell code block with a shadow" src="./examples/shadow.svg" width="600" />

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

* [Twitter](https://twitter.com/charmcli)
* [The Fediverse](https://mastodon.social/@charmcli)
* [Discord](https://charm.sh/chat)

## License

[MIT](https://github.com/charmbracelet/freeze/raw/main/LICENSE)

***

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/">
  <img
    alt="The Charm logo"
    width="400"
    src="https://stuff.charm.sh/charm-badge.jpg"
  />
</a>

Charm热爱开源 • Charm loves open source
