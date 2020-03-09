# Emojiterm [![Build Status](https://travis-ci.com/xyproto/emojiterm.svg?branch=master)](https://travis-ci.com/xyproto/emojiterm) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/emojiterm)](https://goreportcard.com/report/github.com/xyproto/emojiterm) [![License](https://img.shields.io/badge/License-MPL2-brightgreen)](https://raw.githubusercontent.com/xyproto/emojiterm/master/LICENSE)

Want to find a suitable Emoji for use on GitHub, using only a terminal that supports 256 colors? Then this application is for you.

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

### General Info

* Developed on Arch Linux, using Go 1.14.
* Uses [pixterm](https://github.com/eliukblau/pixterm), [imaging](https://github.com/disintegration/imaging), [go-colorful](https://github.com/lucasb-eyer/go-colorful) and [go-github](https://github.com/google/go-github).
* The `display` function in `main.go` is based on code from [pixterm](https://github.com/eliukblau/pixterm) (which is also licensed under `Mozilla Public License 2.0`).

* License: Mozilla Public License 2.0
* Version: 0.1.0
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
