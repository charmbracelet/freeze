# Freeze

<p>
  <a href="https://github.com/charmbracelet/freeze/releases"><img src="https://img.shields.io/github/release/charmbracelet/freeze.svg" alt="Latest Release"></a>
  <a href="https://pkg.go.dev/github.com/charmbracelet/freeze?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="Go Docs"></a>
  <a href="https://github.com/charmbracelet/freeze/actions"><img src="https://github.com/charmbracelet/freeze/workflows/build/badge.svg" alt="Build Status"></a>
</p>

Capture and share your code on the command line.

<p align="center">
  <img alt="output of freeze command, haskell code block" src="https://github.com/charmbracelet/freeze/assets/42545625/6884a2ad-3d8e-4510-ad3d-e935e208d3b8" width="800" />
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

There are a bunch of different options:

* `Output`: where to output the SVG file.
* `Window`: whether to add window bar controls.
* `Border`: whether to add a pixel-wide border along the terminal.
* `Shadow`: whether to add a shadow under the terminal.
* `Radius`: the corner radius of the terminal.
* `Padding`: the terminal padding.
* `Margin`: the image margin.
* `FontFamily`: the terminal font family.
* `FontSize`: the terminal font size.
* `LineHeight`: the terminal line height.

<br />
<img alt="output of freeze command, haskell code block" src="https://github.com/charmbracelet/freeze/assets/42545625/f3fb212f-6629-4253-9c13-105055a4b6e8" width="600" />

## Output

To output different file formats: `.png`, `.jpg`, `.webp`, use the `--output`
flag with the desired extension.

```bash
freeze main.go -o out.png
freeze main.go -o out.jpg
freeze main.go -o out.webp
```

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
