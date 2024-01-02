package futui

import (
	"strconv"
)

type Cell struct {
	r     rune
	style Style
}

type Buffer struct {
	w, h  int
	cells [][]*Cell
}

func (b *Buffer) ensureSize(w, h int) {
	colCount := len(b.cells)
	for col := 0; col < w-colCount; col++ {
		b.cells = append(b.cells, []*Cell{})
	}

	for col := 0; col < len(b.cells); col++ {
		if len(b.cells[col]) < h {
			rowCount := len(b.cells[col])
			for row := 0; row < h-rowCount; row++ {
				cell := Cell{}
				b.cells[col] = append(b.cells[col], &cell)
			}
		}
	}

	if h > b.h {
		b.h = h
	}
	if w > b.w {
		b.w = w
	}
}

func (b *Buffer) Resize(w, h int) {
	if w <= 0 || h <= 0 || (w == b.w && h == b.h) {
		return
	}

	b.ensureSize(w, h)
	b.w, b.h = w, h
}

func (b *Buffer) GetCell(x, y int) *Cell {
	if x >= len(b.cells) {
		return nil
	}
	if y >= len(b.cells[x]) {
		return nil
	}
	return b.cells[x][y]
}

func (b *Buffer) SetCell(x, y int, c *Cell) {
	b.ensureSize(x+1, y+1)

	if prev := b.GetCell(x, y); prev != nil {
		if c.style.Background == "" {
			c.style.Background = prev.style.Background
		}
	}

	b.cells[x][y] = c
}

func (b *Buffer) DrawCell(x, y int, r rune, s Style) *Cell {
	b.ensureSize(x+1, y+1)

	if prev := b.GetCell(x, y); prev != nil {
		if s.Background == "" {
			s.Background = prev.style.Background
		}
		if s.Foreground == "" {
			s.Foreground = prev.style.Background
		}
	}

	cell := Cell{r, s}
	b.cells[x][y] = &cell
	return &cell
}

func (b *Buffer) Size() (int, int) {
	return b.w, b.h
}

func (b *Buffer) Width() int {
	return b.w
}

func (b *Buffer) Height() int {
	return b.h
}

func (b *Buffer) Fill(r rune, s Style) {
	b.ensureSize(b.w, b.h)

	cell := &Cell{r, s}
	for x := range b.cells {
		for y := range b.cells[x] {
			b.cells[x][y] = cell
		}
	}
}

func (b *Buffer) FillStyle(s Style) {
	b.Fill(' ', s)
}

func (b *Buffer) ApplyStyle(s Style) {
	b.ensureSize(b.w, b.h)
	for x := range b.cells {
		for y := range b.cells[x] {
			b.cells[x][y].style = s
		}
	}
}

func (b *Buffer) DrawBuffer(x, y int, buff Buffer) {
	if x < 0 {
		panic("DrawBuffer: x < 0 = " + strconv.Itoa(x))
	}
	if y < 0 {
		panic("DrawBuffer: y < 0 = " + strconv.Itoa(y))
	}
	bw, bh := buff.Size()
	b.ensureSize(x+bw, y+bh)

	for bx := 0; bx < bw; bx++ {
		for by := 0; by < bh; by++ {
			c := buff.GetCell(bx, by)
			if c != nil {
				b.SetCell(x+bx, y+by, c)
			}
		}
	}
}

func (b *Buffer) DrawComponent(x int, y int, c IComponent) Buffer {
	buffer := c.Render()

	cw, ch := buffer.Size()
	b.ensureSize(x+cw, y+ch)

	b.DrawBuffer(x, y, buffer)

	return buffer
}

func (b *Buffer) Rect(x, y, w, h int, s Style) {
	if x < 0 || y < 0 || w <= 0 || h <= 0 {
		return
	}
	buff := Buffer{w: w, h: h}
	buff.Fill(' ', s)
	b.DrawBuffer(x, y, buff)
}

func (b *Buffer) Text(x int, y int, text string, s Style) {
	buff := Buffer{}
	for i, r := range text {
		buff.DrawCell(i, 0, r, s)
	}
	b.DrawBuffer(x, y, buff)
}
