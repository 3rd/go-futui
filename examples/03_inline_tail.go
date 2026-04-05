package main

import (
	"fmt"
	"time"

	ui "github.com/3rd/go-futui"
	"github.com/gdamore/tcell/v2"
)

type TailApp struct {
	ui.App
	lines []string
}

func (app *TailApp) Setup() {
	ticker := time.NewTicker(time.Millisecond * 300)
	go func() {
		for range ticker.C {
			app.lines = append(app.lines, fmt.Sprintf("Line %d", len(app.lines)))
			app.Update()
		}
	}()
}

func (app *TailApp) Render() ui.Buffer {
	b := ui.Buffer{}
	for i, line := range app.lines {
		b.Text(0, i, line, ui.Style{})
	}
	return b
}

func (app *TailApp) OnKeypress(ev tcell.EventKey) {
	if ev.Key() == tcell.KeyCtrlC {
		app.Quit()
	}
}

func main() {
	app := TailApp{}
	app.RunWithOptions(&app, ui.RunOptions{
		Mode:            ui.RunModeInline,
		InlineMaxHeight: 5,
	})
}
