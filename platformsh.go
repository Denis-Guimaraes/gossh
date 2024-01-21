package main

import (
	"gossh/internal/platformsh"
	"gossh/internal/project"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var Projects project.Projects
var wListPlatform *widget.List
var projectLoader *widget.ProgressBarInfinite

func platform() *container.TabItem {
	return container.NewTabItemWithIcon("Platform", theme.MenuIcon(), create())
}

func create() *fyne.Container {
	go loadProjects()
	projectLoader = widget.NewProgressBarInfinite()
	action := action()
	top := container.NewVBox(action, projectLoader)
	wListPlatform = widget.NewList(Projects.Count, createProjectsItem, updateProjectsItem)
	return container.NewBorder(top, nil, nil, nil, wListPlatform)
}

func action() *fyne.Container {
	updateBtn := widget.NewButtonWithIcon("Mettre à jour", theme.ViewRefreshIcon(), refreshProjects)
	loginBtn := widget.NewButtonWithIcon("Connexion", theme.LoginIcon(), platformsh.Login)
	row := container.NewGridWithRows(1, loginBtn, updateBtn)
	return row
}

func createProjectsItem() fyne.CanvasObject {
	style := fyne.TextStyle{Bold: true}
	title := widget.NewLabelWithStyle("", fyne.TextAlignLeading, style)
	login := container.NewAdaptiveGrid(3)
	return container.NewGridWithRows(2, title, login)
}

func updateProjectsItem(i int, o fyne.CanvasObject) {
	item := Projects.Items[i]
	o.(*fyne.Container).Objects = nil

	style := fyne.TextStyle{Bold: true}
	title := widget.NewLabelWithStyle(item.Name, fyne.TextAlignLeading, style)
	login := container.NewAdaptiveGrid(3)
	for _, e := range item.Environments {
		login.Add(widget.NewButtonWithIcon(e.Name, theme.LoginIcon(), platformsh.Ssh(item.Id, e.Id)))
	}
	o.(*fyne.Container).Add(title)
	o.(*fyne.Container).Add(login)
}

func loadProjects() {
	startLoader()
	Projects = project.GetProjects()
	stoptLoader()
}

func refreshProjects() {
	startLoader()
	Projects = project.Projects{}
	Projects = project.RefreshProjects()
	stoptLoader()
}

func startLoader() {
	projectLoader.Start()
	projectLoader.Show()
}

func stoptLoader() {
	projectLoader.Stop()
	projectLoader.Hide()
}
