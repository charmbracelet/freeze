# Freeze

<p>
  <img src="https://user-images.githubusercontent.com/42545625/198402537-12ca2f6c-0779-4eb8-a67c-8db9cb3df13c.png#gh-dark-mode-only" width="500" />
  <img src="https://user-images.githubusercontent.com/42545625/198402542-a305f669-a05a-4d91-b18b-ca76e72b655a.png#gh-light-mode-only" width="500" />
  <br>
  <a href="https://github.com/charmbracelet/freeze/releases"><img src="https://img.shields.io/github/release/charmbracelet/freeze.svg" alt="Latest Release"></a>
  <a href="https://pkg.go.dev/github.com/charmbracelet/freeze?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="Go Docs"></a>
  <a href="https://github.com/charmbracelet/freeze/actions"><img src="https://github.com/charmbracelet/freeze/workflows/build/badge.svg" alt="Build Status"></a>
</p>

Capture screenshots and share images of your code on the command line.

<img alt="Welcome to VHS" src="https://vhs.charm.sh/vhs-3yL9rCVdQNUjYwVBqzW7Kc.gif" width="600" />

## Tutorial

To snap an image of your code, run the following:

```sh
freeze main.go > out.svg
```

The `out.svg` file will contain your code screenshot.

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

[MIT](https://github.com/charmbracelet/vhs/raw/main/LICENSE)

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
