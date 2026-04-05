package futui

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	ansiReset      = "\x1b[0m"
	ansiClear      = "\x1b[2K"
	ansiSaveCursor = "\x1b7"
	ansiRestore    = "\x1b8"
	ansiHideCursor = "\x1b[?25l"
	ansiShowCursor = "\x1b[?25h"
)

type inlineRenderer struct {
	out        *bufio.Writer
	prevHeight int
	maxHeight  int
}

func newInlineRenderer(w io.Writer, maxHeight int) *inlineRenderer {
	return &inlineRenderer{
		out:       bufio.NewWriter(w),
		maxHeight: maxHeight,
	}
}

func (r *inlineRenderer) Render(buff Buffer) error {
	lines := bufferToANSILines(buff)
	if r.maxHeight > 0 && len(lines) > r.maxHeight {
		lines = lines[len(lines)-r.maxHeight:]
	}
	return r.renderLines(lines)
}

func (r *inlineRenderer) Clear() error {
	return r.renderLines(nil)
}

func (r *inlineRenderer) renderLines(lines []string) error {
	currentHeight := len(lines)

	if r.prevHeight > 0 {
		if _, err := r.out.WriteString("\r"); err != nil {
			return err
		}
		if r.prevHeight > 1 {
			if _, err := r.out.WriteString(cursorUp(r.prevHeight - 1)); err != nil {
				return err
			}
		}
	}

	maxHeight := currentHeight
	if r.prevHeight > maxHeight {
		maxHeight = r.prevHeight
	}

	for i := 0; i < maxHeight; i++ {
		if _, err := r.out.WriteString("\r"); err != nil {
			return err
		}
		if _, err := r.out.WriteString(ansiClear); err != nil {
			return err
		}
		if i < currentHeight {
			if _, err := r.out.WriteString(lines[i]); err != nil {
				return err
			}
		}
		if i < maxHeight-1 {
			if _, err := r.out.WriteString("\n"); err != nil {
				return err
			}
		}
	}

	r.prevHeight = currentHeight
	return r.out.Flush()
}

func cursorUp(lines int) string {
	if lines <= 0 {
		return ""
	}
	return fmt.Sprintf("\x1b[%dA", lines)
}

func bufferToANSILines(buff Buffer) []string {
	width, height := buff.Size()
	if width <= 0 || height <= 0 {
		return nil
	}

	lastContent := lastContentLine(buff, width, height)
	if lastContent < 0 {
		return nil
	}

	lines := make([]string, lastContent+1)
	for y := 0; y <= lastContent; y++ {
		var sb strings.Builder
		sb.WriteString(ansiReset)

		lastStyle := Style{}
		for x := 0; x < width; x++ {
			cell := buff.GetCell(x, y)
			r := ' '
			style := Style{}
			if cell != nil {
				if cell.r != 0 {
					r = cell.r
				}
				style = cell.style
			}

			if !stylesEqual(style, lastStyle) {
				sb.WriteString(ansiReset)
				sb.WriteString(styleToANSI(style))
				lastStyle = style
			}
			sb.WriteRune(r)
		}

		if !isZeroStyle(lastStyle) {
			sb.WriteString(ansiReset)
		}
		lines[y] = sb.String()
	}

	return lines
}

func lastContentLine(buff Buffer, width int, height int) int {
	for y := height - 1; y >= 0; y-- {
		if lineHasContent(buff, y, width) {
			return y
		}
	}
	return -1
}

func lineHasContent(buff Buffer, y int, width int) bool {
	for x := 0; x < width; x++ {
		cell := buff.GetCell(x, y)
		if cell == nil {
			continue
		}
		if !isZeroStyle(cell.style) {
			return true
		}
		if cell.r != 0 && cell.r != ' ' {
			return true
		}
	}
	return false
}

func styleToANSI(style Style) string {
	if isZeroStyle(style) {
		return ""
	}

	var sb strings.Builder

	if style.Foreground != "" {
		r, g, b := style.Foreground.Noire().RGB()
		sb.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm", int(r), int(g), int(b)))
	}
	if style.Background != "" {
		r, g, b := style.Background.Noire().RGB()
		sb.WriteString(fmt.Sprintf("\x1b[48;2;%d;%d;%dm", int(r), int(g), int(b)))
	}
	if style.Bold {
		sb.WriteString("\x1b[1m")
	}
	if style.Dim {
		sb.WriteString("\x1b[2m")
	}
	if style.Italic {
		sb.WriteString("\x1b[3m")
	}
	if style.Underline {
		sb.WriteString("\x1b[4m")
	}
	if style.Blink {
		sb.WriteString("\x1b[5m")
	}
	if style.Reverse {
		sb.WriteString("\x1b[7m")
	}
	if style.StrikeThrough {
		sb.WriteString("\x1b[9m")
	}

	return sb.String()
}

func stylesEqual(a Style, b Style) bool {
	return a == b
}

func isZeroStyle(s Style) bool {
	return s == (Style{})
}

func writeANSI(w io.Writer, seq string) {
	if w == nil || seq == "" {
		return
	}
	_, _ = w.Write([]byte(seq))
}
