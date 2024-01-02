package futui

import "github.com/gdamore/tcell/v2"

type Style struct {
	Background    Color
	Foreground    Color
	Bold          bool
	Blink         bool
	Reverse       bool
	Underline     bool
	Dim           bool
	Italic        bool
	StrikeThrough bool
}

func (s *Style) build() tcell.Style {
	style := tcell.StyleDefault

	// colors
	background := s.Background.Build()
	if background != nil {
		style = style.Background(*background)
	}
	foreground := s.Foreground.Build()
	if foreground != nil {
		style = style.Foreground(*foreground)
	}

	// styles
	if s.Bold {
		style = style.Bold(true)
	}
	if s.Blink {
		style = style.Bold(true)
	}
	if s.Reverse {
		style = style.Reverse(true)
	}
	if s.Underline {
		style = style.Underline(true)
	}
	if s.Dim {
		style = style.Dim(true)
	}
	if s.Italic {
		style = style.Italic(true)
	}
	if s.StrikeThrough {
		style = style.StrikeThrough(true)
	}

	return style
}
