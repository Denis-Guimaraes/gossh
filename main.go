package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	load()
	a := app.New()
	Window := a.NewWindow("GOSSH")
	Window.Resize(fyne.NewSize(600, 600))
	Window.SetPadded(true)
	Window.CenterOnScreen()

	tabs := container.NewAppTabs(list(Window), add(), platform())
	tabs.SetTabLocation(container.TabLocationLeading)
	Window.SetContent(tabs)
	Window.ShowAndRun()
}
