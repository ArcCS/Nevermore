// Copyright 2016 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package text

import (
	"bytes"
)

// Mixed ANSI and 256 Color codes for colorizing the game
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Normal    = "\033[22m"
	Black     = "\033[30m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Cyan      = "\033[36m"
	White     = "\033[37m"
	BGBlack   = "\033[40m"
	BGRed     = "\033[41m"
	BGGreen   = "\033[42m"
	BGYellow  = "\033[43m"
	BGBlue    = "\033[44m"
	BGMagenta = "\033[45m"
	BGCyan    = "\033[46m"
	BGWhite   = "\033[47m"
	Gray 	  = "\u001b[38;5;241m"
	LightGreen = "\u001b[38;5;36m"
	Turquoise = "\u001b[38;5;31m"
	Pink 	  = "\u001b[38;5;168m"

	// Setup brown as an alias for yellow
	Brown   = Yellow
	BGBrown = BGYellow

	// WolfMUD specific meta colors
	Good   = Green
	Info   = Yellow
	Bad    = Red
	Prompt = Magenta
)

// colorTable maps color place holders to color escape sequences. Colorize uses
// this map to to substitute color placeholders of the form [COLOR] with the
// matching ANSI escape sequence.
var colorTable = map[string]string{
	"[RESET]":     Reset,
	"[BOLD]":      Bold,
	"[NORMAL]":    Normal,
	"[BLACK]":     Black,
	"[RED]":       Red,
	"[GREEN]":     Green,
	"[BROWN]":     Brown,
	"[YELLOW]":    Yellow,
	"[BLUE]":      Blue,
	"[MAGENTA]":   Magenta,
	"[CYAN]":      Cyan,
	"[WHITE]":     White,
	"[BGBLACK]":   BGBlack,
	"[BGRED]":     BGRed,
	"[BGGREEN]":   BGGreen,
	"[BGBROWN]":   BGBrown,
	"[BGYELLOW]":  BGYellow,
	"[BGBLUE]":    BGBlue,
	"[BGMAGENTA]": BGMagenta,
	"[BGCYAN]":    BGCyan,
	"[BGWHITE]":   BGWhite,
}

// Colorize returns a []byte with color place holders replaced with their ANSI
// escape sequence equivalent. Color place holders have the format [COLOR]
// where COLOR represents the name of the color (uppercased) to be used. For
// example:
//
//	Colorize([]byte("[RED]Hello [GREEN]World![DEFAULT]"))
//
// Would return a []byte with [RED] and [GREEN] replaced with the ANSI escape
// sequences \033[31m and \033[32m respectively causing Hello to be displayed
// in red and World! to be displayed in green.
//
// The returned slice is always a copy even it the original contains no colors.
//
// Use of this function is discouraged due to relatively poor performance. It's
// main use is to render text from files, such as those loaded when the server
// is initially started. In code it is better to use the ANSI escape sequence
// constants directly.
func Colorize(in []byte) []byte {
	out := make([]byte, len(in))
	copy(out, in)

	if bytes.IndexByte(out, ']') == -1 {
		return out
	}

	p := 0
	for color, code := range colorTable {
		// Quick exit? Check for ']' as '[' also in replacement text
		if bytes.IndexByte(out, ']') == -1 {
			break
		}
		// Shortcut id we can't find an instance of the current color?
		if p = bytes.Index(out, []byte(color)); p == -1 {
			continue
		}
		// If no shortcut available we can still use the position of the check to
		// shorten the length of the slice we are doing replacements on
		out = append(out[:p], bytes.Replace(out[p:], []byte(color), []byte(code), -1)...)
	}
	return out
}
