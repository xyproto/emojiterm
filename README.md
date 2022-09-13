# Emojiterm [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/emojiterm)](https://goreportcard.com/report/github.com/xyproto/emojiterm) [![License](https://img.shields.io/badge/License-BSD3-brightgreen)](https://raw.githubusercontent.com/xyproto/emojiterm/main/LICENSE)

* Do you want to find a suitable Emoji for use on GitHub, using only a terminal that supports 256 colors?
* Then this application is for you.

`emojiterm` can list all available emoji names, or search for a keyword and display the emoji directly on the terminal.

![recording](img/recording.gif)

## Requirements

* Go >= 1.10.
* A terminal emulator that supports 256 colors.

## Installation

Install with your favorite package manager, if possible.

### Manual installation of the development version

    go get -u github.com/xyproto/emojiterm

## Supported terminal emulators

These terminal emulators are known to work:

* `konsole`
* `xfce4-terminal`

This one does not work:

* `urxvt`

This one works, but does not look quite right:

* `st`

## Request limit

If you reach the request limit for using the GitHub API, placing a valid token in the `GITHUB_TOKEN` environment variable should solve the issue.

For generating a token, just visit [github.com/settings/tokens](https://github.com/settings/tokens) and click "Generate new token". None of the boxes with extra access needs to be checked, since `emojiterm` only fetches emoji-related information.

## Slideshow

You can display a slideshow of all available GitHub emojis with the following command:

```bash
GITHUB_TOKEN="asdf" emojiterm -a
```

This can also take an optional search term:

```bash
GITHUB_TOKEN="asdf" emojiterm -a cat
```

## Setting up an alias

One way of using `emojiterm` is through an alias, that also sets the `GITHUB_TOKEN`, like this:

```bash
alias et='GITHUB_TOKEN="asdf" emojiterm'
```

Replace `asdf` with your own token, that can be generated here: [github.com/settings/tokens](https://github.com/settings/tokens).

## General Info

* Developed on Arch Linux, using Go 1.18.
* Uses [pixterm](https://github.com/eliukblau/pixterm), [imaging](https://github.com/disintegration/imaging), [go-colorful](https://github.com/lucasb-eyer/go-colorful) and [go-github](https://github.com/google/go-github).
* The `display` function in `main.go` is based on code from [pixterm](https://github.com/eliukblau/pixterm) (which is also licensed under `Mozilla Public License 2.0`).
* License: BSD-3
* Version: 0.3.2
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
