package futui

type IComponent interface {
	Render() Buffer
}

type Component struct {
	IComponent
	Style Style
}

func (*Component) Render(*Buffer) {
	panic("Component does not implement the render function.")
}
