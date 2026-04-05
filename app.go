package futui

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/terminfo"
)

type IApp interface {
	Render() Buffer
}
type IAppWithSetup interface {
	Setup()
}
type IAppWithKeypressHandler interface {
	OnKeypress(event tcell.EventKey)
}
type IAppWithResizeHandler interface {
	OnResize(w, h int)
}

type RunMode int

const (
	RunModeScreen RunMode = iota
	RunModeInline
)

type RunOptions struct {
	Mode            RunMode
	InlineMaxHeight int
}

type Handlers struct {
	render   func() Buffer
	setup    *func()
	keypress *func(tcell.EventKey)
	resize   *func(int, int)
}

type App struct {
	handlers    Handlers
	Screen      tcell.Screen
	QuitChannel chan struct{}
	mode        RunMode
	inline      *inlineRenderer
}

func (app *App) render() {
	// user render
	if app.handlers.render == nil {
		panic("app.handlers.render is nil")
	}
	buff := app.handlers.render()

	if app.mode == RunModeInline && app.inline != nil {
		if err := app.inline.Render(buff); err != nil {
			panic(err)
		}
		return
	}

	// proxy
	width, height := buff.Size()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if cell := buff.GetCell(x, y); cell != nil {
				app.Screen.SetContent(x, y, cell.r, nil, cell.style.build())
			}
		}
	}

	app.Screen.Show()
}

func (app *App) loop() {
	for {
		event := app.Screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventKey:
			if app.handlers.keypress != nil {
				(*app.handlers.keypress)(*event)
			} else {
				// <c-c> to quit by default
				if event.Key() == tcell.KeyCtrlC {
					app.Quit()
				}
			}
		case *tcell.EventResize:
			if app.handlers.resize != nil {
				width, height := app.Screen.Size()
				(*app.handlers.resize)(width, height)
				if app.mode == RunModeScreen {
					app.Screen.Sync()
				}
			}
		}
	}
}

func (app *App) Update() {
	app.render()
}

func (app *App) Quit() {
	app.Screen.Fini()
	if app.mode == RunModeInline {
		writeANSI(os.Stdout, ansiShowCursor)
	}
	close(app.QuitChannel)
}

func (app *App) WaitForExit() {
	for range app.QuitChannel {
		break
	}
}

func (app *App) Size() (width int, height int) {
	return app.Screen.Size()
}

func (app *App) Width() int {
	width, _ := app.Screen.Size()
	return width
}

func (app *App) Height() int {
	_, height := app.Screen.Size()
	return height
}

func (app *App) Clear() {
	if app.mode == RunModeInline && app.inline != nil {
		if err := app.inline.Clear(); err != nil {
			panic(err)
		}
		return
	}
	app.Screen.Clear()
	app.Screen.Show()
}

func (app *App) Sync() {
	if app.mode == RunModeInline {
		app.render()
		return
	}
	app.Screen.Sync()
}

func (app *App) Beep() {
	app.Screen.Beep()
}

func (app *App) Run(userApp IApp) {
	app.RunWithOptions(userApp, RunOptions{})
}

func (app *App) RunWithOptions(userApp IApp, opts RunOptions) {
	// setup handlers
	app.handlers.render = userApp.Render

	if userAppWithSetup, ok := userApp.(IAppWithSetup); ok {
		setupHandler := func() {
			userAppWithSetup.Setup()
		}
		app.handlers.setup = &setupHandler
	}

	if userAppWithKeypressHandler, ok := userApp.(IAppWithKeypressHandler); ok {
		keypressHandler := func(ev tcell.EventKey) {
			userAppWithKeypressHandler.OnKeypress(ev)
		}
		app.handlers.keypress = &keypressHandler
	}
	if userAppWithResizeHandler, ok := userApp.(IAppWithResizeHandler); ok {
		resizeHandler := func(w, h int) {
			userAppWithResizeHandler.OnResize(w, h)
		}
		app.handlers.resize = &resizeHandler
	}

	mode := RunModeScreen
	if opts.Mode == RunModeInline {
		mode = RunModeInline
	}

	// setup screen
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	var (
		screen tcell.Screen
		err    error
	)
	if mode == RunModeInline {
		writeANSI(os.Stdout, ansiSaveCursor)
		screen, err = newInlineScreen()
	} else {
		screen, err = tcell.NewScreen()
	}
	if err != nil {
		panic("Cannot create screen.")
	}
	if err = screen.Init(); err != nil {
		panic("Cannot initialize screen.")
	}
	if mode == RunModeInline {
		writeANSI(os.Stdout, ansiRestore)
		writeANSI(os.Stdout, ansiHideCursor)
	}
	screen.SetStyle(tcell.StyleDefault)
	if mode == RunModeScreen {
		screen.Clear()
	}

	// internal state
	app.Screen = screen
	app.QuitChannel = make(chan struct{})
	app.mode = mode
	if mode == RunModeInline {
		app.inline = newInlineRenderer(os.Stdout, opts.InlineMaxHeight)
	} else {
		app.inline = nil
	}

	if app.handlers.setup != nil {
		(*app.handlers.setup)()
	}

	// boot
	app.render()
	go app.loop()
	app.WaitForExit()
}

func newInlineScreen() (tcell.Screen, error) {
	tty, err := tcell.NewStdIoTty()
	if err != nil {
		return nil, err
	}

	term := os.Getenv("TERM")
	candidates := []string{}
	if term != "" {
		candidates = append(candidates, term)
	}
	candidates = append(candidates, "xterm-256color", "xterm")

	for _, name := range candidates {
		ti, err := terminfo.LookupTerminfo(name)
		if err != nil || ti == nil {
			continue
		}
		clone := *ti
		clone.EnterCA = ""
		clone.ExitCA = ""
		clone.Clear = ""
		return tcell.NewTerminfoScreenFromTtyTerminfo(tty, &clone)
	}

	if err == nil {
		err = terminfo.ErrTermNotFound
	}
	return nil, err
}
