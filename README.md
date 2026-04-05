# go-futui: Futuristic TUI library for Go.

> [!IMPORTANT]
> Very much WIP, you probably don't want to build stuff with this yet.

## Installation

```sh
go get github.com/3rd/go-futui
```

## Usage

```go
package main

import (
	"strconv"
	"time"

	ui "github.com/3rd/go-futui"
	"github.com/gdamore/tcell/v2"
)

// basic component
type Counter struct {
	ui.Component
	value int
}

func (c *Counter) Render() ui.Buffer {
	b := ui.Buffer{}
	b.Text(0, 0, "Count: "+strconv.Itoa(c.value), ui.Style{
		Background: "#ff0000",
		Foreground: "#000000",
		Bold:       true,
		Italic:     true,
		Underline:  true,
	})
	return b
}

// app with state
type State struct {
	counter int
}

type MyApp struct {
	ui.App
	state State
}

func (app *MyApp) Setup() {
	ticker := time.NewTicker(time.Millisecond * 16)
	go func() {
		for range ticker.C {
			app.state.counter++
			app.Update()
		}
	}()
}

func (app *MyApp) Render() ui.Buffer {
	b := ui.Buffer{}
	b.Resize(app.Width(), app.Height())
	b.FillStyle(ui.Style{
		Background: "#00ff00",
		Foreground: "#000000",
	})
	b.Text(0, 0, "Hello World!", ui.Style{Foreground: "#220ee0"})
	b.DrawComponent(0, 1, &Counter{value: app.state.counter})
	return b
}

func (app *MyApp) OnKeypress(ev tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyCtrlC:
		app.Quit()
	}
	app.state.counter++
	app.Update()
}

func (app *MyApp) OnResize() {
	app.Update()
}

func main() {
	app := MyApp{}
	app.Run(&app)
}
```

## Inline (scrollback) mode

By default, go-futui renders in a full-screen TUI. To preserve terminal history and only
rewrite the last lines your app rendered, run in inline mode:

```go
func main() {
	app := MyApp{}
	app.RunWithOptions(&app, ui.RunOptions{
		Mode: ui.RunModeInline,
	})
}
```

If `TERM` isn't recognized, inline mode falls back to `xterm-256color`/`xterm` terminfo
for input handling. Set `TERM` correctly for best results.

Inline mode trims trailing empty lines (no content or styling) so apps that resize
their buffer to full height don’t accidentally rewrite the whole screen.

Inline mode tips:
- Avoid `Resize(app.Width(), app.Height())` unless you really want full-screen updates.
- If you apply `FillStyle` across the whole buffer, those lines are treated as content
  and will be rewritten (by design).
- Inline mode hides the cursor while running and restores it on exit.

To cap inline output to the last N lines, use `InlineMaxHeight`:

```go
app.RunWithOptions(&app, ui.RunOptions{
	Mode:            ui.RunModeInline,
	InlineMaxHeight: 5,
})
```
