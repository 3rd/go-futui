package futui

import (
	"github.com/gdamore/tcell/v2"
)

type IApp interface {
	Render() Buffer
}
type IAppWithSetup interface {
	Setup()
}
type IAppWithKeypressHandler interface {
	OnKeypress(key tcell.Key, r rune)
}
type IAppWithResizeHandler interface {
	OnResize(w, h int)
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
}

func (app *App) render() {
	// user render
	if app.handlers.render == nil {
		panic("app.handlers.render is nil")
	}
	buff := app.handlers.render()

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
				app.Screen.Sync()
			}
		}
	}
}

func (app *App) Update() {
	app.render()
}

func (app *App) Quit() {
	app.Screen.Fini()
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
	app.Screen.Clear()
	app.Screen.Show()
}

func (app *App) Sync() {
	app.Screen.Sync()
}

func (app *App) Beep() {
	app.Screen.Beep()
}

func (app *App) Run(userApp IApp) {
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
			userAppWithKeypressHandler.OnKeypress(ev.Key(), ev.Rune())
		}
		app.handlers.keypress = &keypressHandler
	}
	if userAppWithResizeHandler, ok := userApp.(IAppWithResizeHandler); ok {
		resizeHandler := func(w, h int) {
			userAppWithResizeHandler.OnResize(w, h)
		}
		app.handlers.resize = &resizeHandler
	}

	// setup screen
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, err := tcell.NewScreen()
	if err != nil {
		panic("Cannot create screen.")
	}
	if err = screen.Init(); err != nil {
		panic("Cannot initialize screen.")
	}
	screen.SetStyle(tcell.StyleDefault)
	screen.Clear()

	// internal state
	app.Screen = screen
	app.QuitChannel = make(chan struct{})

	if app.handlers.setup != nil {
		(*app.handlers.setup)()
	}

	// boot
	app.render()
	go app.loop()
	app.WaitForExit()
}
