package clip

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	ScreenWidth        = 80
	ScreenHeight       = 40
	DefaultWrapOptions = WrapOptions{}
	ShowDebug          = false
)

func Debug(format string, args ...interface{}) {
	if ShowDebug {
		if format[len(format)-1] != '\n' {
			format += "\n"
		}
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

type WrapOptions struct {
	Width             int
	Indent            int // Indent for first line (and \r)
	Margin            int // Indent for 2+ lines (and \t)
	TrimNewlineSpaces bool
	MarginChar        byte
	IndentChar        byte
	Justify           bool
}

func init() {
	width, height, err := terminal.GetSize(1)
	if err != nil {
		ScreenWidth = width
		ScreenHeight = height
	}

	DefaultWrapOptions = *NewWrapOptions()
}

func NewWrapOptions() *WrapOptions {
	return &WrapOptions{
		Width:             ScreenWidth,
		Indent:            0,
		Margin:            0,
		TrimNewlineSpaces: false,
		MarginChar:        '\t',
		IndentChar:        '\r',
		Justify:           false,
	}
}

func Justify(text string, width int) string {
	addCR := ""

	if len(text) > 0 && text[len(text)-1] == '\n' {
		addCR = "\n"
	}

	text = strings.TrimRightFunc(text, unicode.IsSpace)
	size := len(text)

	if size >= width {
		return text + addCR
	}

	first := 0
	for ; first < size; first++ {
		if text[first] != ' ' {
			break
		}
	}
	if first+1 >= size {
		return text + addCR
	}

	need := width - size
	curPos := size - 1

	for need > 0 {
		initPos := curPos
		// Exit out of space sequence
		for text[curPos] == ' ' {
			curPos--
			if curPos == first {
				curPos = size - 1
			}
			if curPos == initPos {
				return text + addCR
			}
		}

		// Now look for space
		initPos = curPos
		for text[curPos] != ' ' {
			curPos--
			if curPos == first {
				curPos = size - 1
			}
			if curPos == initPos {
				return text + addCR
			}
		}

		// found a space
		text = text[:curPos] + " " + text[curPos:]
		need--
	}

	return text + addCR
}

func (wo *WrapOptions) Wrap(text string) string {
	// Stop weird cases
	if wo.MarginChar == 0 {
		wo.MarginChar = '\t'
	}

	if wo.IndentChar == 0 {
		wo.IndentChar = '\r'
	}

	if wo.Margin < 0 {
		wo.Margin = 0
	}

	if wo.Indent < 0 {
		wo.Indent = 0
	}

	if wo.Width == 0 {
		wo.Width = ScreenWidth
	}

	if wo.Width <= wo.Margin {
		wo.Width = wo.Margin + 1
	}

	if wo.Width <= wo.Indent {
		wo.Width = wo.Indent + 1
	}

	result := ""

	// Remove any trailing \n - but just one
	if len(text) > 0 && text[len(text)-1] == '\n' {
		text = text[:len(text)-1]
	}

	// Last line doesn't need trailing spaces, but add indent
	text = strings.Repeat(" ", wo.Indent) + strings.TrimRight(text, " ")

	chop, pos, start := wo.Indent, wo.Indent, wo.Indent

	for {
		// hitMarginChar := false
		hitIndentChar := false
		ch := byte(0)
		if pos < len(text) {
			ch = text[pos]
		}
		Debug("start:%d pos:%d len:%d ch: %q", start, pos, len(text), ch)

		if ch == wo.MarginChar {
			// hitMarginChar = true
			if pos < wo.Margin {
				replace := strings.Repeat(" ", wo.Margin-pos)
				text = text[:pos] + replace + text[pos+1:]
				if pos == start {
					start = wo.Margin
					pos, chop = start, start
				}
				continue
			}
			text = text[:pos] + "\n" + text[pos+1:]
			ch = '\n'
		} else if ch == wo.IndentChar {
			hitIndentChar = true
			if pos < wo.Indent {
				replace := strings.Repeat(" ", wo.Indent-pos)
				text = text[:pos] + replace + text[pos+1:]
				if pos == start {
					start = wo.Indent
					pos, chop = start, start
				}
				continue
			}
			text = text[:pos] + "\n" + text[pos+1:]
			ch = '\n'
		}

		if unicode.IsSpace(rune(ch)) {
			chop = pos
		}

		// These 3 cases force us to chop now
		if ch == 0 {
			chop = len(text)
			Debug("hit len of text - new chop: %d", chop)
		} else if ch == '\n' {
			chop = pos + 1
			Debug("Hit newline - new chop: %d", chop)
		} else if pos == wo.Width {
			if chop == start {
				chop = wo.Width
			}
			Debug("hit width - chop: %d", chop)
		} else {
			// Any other case we just move on to the next char
			pos++
			continue
		}

		add := text[:chop]
		rest := text[chop:]

		if chop > 0 && len(rest) != 0 {
			hadCR := add[chop-1] == '\n'

			if !hadCR {
				add = strings.TrimRight(add, " ")
				add += "\n"
			} else {
				add = strings.TrimRight(add[:len(add)-1], " ") + "\n"
			}

			if !hadCR || wo.TrimNewlineSpaces {
				rest = strings.TrimLeft(rest, " ")
			}

			if wo.Justify && !hadCR && len(add) > 2 {
				add = Justify(add, wo.Width)
			}
		} else {
			// rest == "" so we're at the end, trim spaces
			add = strings.TrimRight(add, " ")
		}
		Debug("Adding: %q  Rest: %q\n", add, rest)
		result += add

		text = rest
		if text == "" {
			return result
		}

		Debug("HitIndentChar: %v", hitIndentChar)
		if hitIndentChar {
			text = strings.Repeat(" ", wo.Indent) + text
			chop, pos, start = wo.Indent, wo.Indent, wo.Indent
		} else {
			text = strings.Repeat(" ", wo.Margin) + text
			chop, pos, start = wo.Margin, wo.Margin, wo.Margin
		}

	}
}

func WrapWithArgs(text string, width int, indent int, margin int) string {
	wo := DefaultWrapOptions

	wo.Width = width
	wo.Indent = indent
	wo.Margin = margin

	return wo.Wrap(text)
}

func Wrap(text string) string {
	return DefaultWrapOptions.Wrap(text)
}
