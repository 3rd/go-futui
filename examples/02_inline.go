package main

import (
	"fmt"
	"time"

	ui "github.com/3rd/go-futui"
	"github.com/gdamore/tcell/v2"
)

type InlineApp struct {
	ui.App
	counter int
}

func (app *InlineApp) Setup() {
	ticker := time.NewTicker(time.Millisecond * 250)
	go func() {
		for range ticker.C {
			app.counter++
			app.Update()
		}
	}()
}

func (app *InlineApp) Render() ui.Buffer {
	b := ui.Buffer{}
	b.Text(0, 0, "Inline mode (scrollback preserved)", ui.Style{Bold: true})
	b.Text(0, 1, fmt.Sprintf("Count: %d", app.counter), ui.Style{})
	return b
}

func (app *InlineApp) OnKeypress(ev tcell.EventKey) {
	if ev.Key() == tcell.KeyCtrlC {
		app.Quit()
	}
}

func main() {
	app := InlineApp{}
	app.RunWithOptions(&app, ui.RunOptions{
		Mode: ui.RunModeInline,
	})
}
