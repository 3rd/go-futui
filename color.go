package futui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/teacat/noire"
)

type Color string

func (c Color) Build() *tcell.Color {
	if len(c) == 0 {
		return nil
	}
	color := noire.NewHex(string(c))
	r, g, b := color.RGB()
	tcellColor := tcell.NewRGBColor(int32(r), int32(g), int32(b))
	return &tcellColor
}

func (c Color) Noire() noire.Color {
	return noire.NewHex(string(c))
}

func (c Color) Lighten(val float64) Color {
	return Color(c.Noire().Lighten(val).Hex())
}

func (c Color) Brighten(val float64) Color {
	return Color(c.Noire().Brighten(val).Hex())
}

func (c Color) Tint(val float64) Color {
	return Color(c.Noire().Tint(val).Hex())
}

func (c Color) Darken(val float64) Color {
	return Color(c.Noire().Darken(val).Hex())
}

func (c Color) Shade(val float64) Color {
	return Color(c.Noire().Shade(val).Hex())
}

func (c Color) Saturate(val float64) Color {
	return Color(c.Noire().Saturate(val).Hex())
}

func (c Color) Desaturate(val float64) Color {
	return Color(c.Noire().Desaturate(val).Hex())
}

func (c Color) AdjustHue(val float64) Color {
	return Color(c.Noire().AdjustHue(val).Hex())
}

func (c Color) Invert() Color {
	return Color(c.Noire().Invert().Hex())
}

func (c Color) Complement() Color {
	return Color(c.Noire().Complement().Hex())
}

func (c Color) Grayscale() Color {
	return Color(c.Noire().Grayscale().Hex())
}

func (c Color) OptimalForeground() Color {
	return Color(c.Noire().Foreground().Hex())
}
