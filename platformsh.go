package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Project struct {
	Name         string
	Id           string
	Environments []Environment
}

type Environment struct {
	Name string
	Id   string
}

var Projects []Project
var wListPlatform *widget.List
var projectLoader *widget.ProgressBarInfinite

func platform() *container.TabItem {
	return container.NewTabItemWithIcon("Platform", theme.MenuIcon(), create())
}

func create() *fyne.Container {
	go updateProjects()
	projectLoader = widget.NewProgressBarInfinite()
	action := action()
	wListPlatform = widget.NewList(projectsLength, createProjectsItem, updateProjectsItem)
	return container.NewBorder(action, nil, nil, nil, wListPlatform, projectLoader)
}

func action() *fyne.Container {
	updateBtn := widget.NewButtonWithIcon("Mettre à jour", theme.ViewRefreshIcon(), updateProjects)
	loginBtn := widget.NewButtonWithIcon("Connexion", theme.LoginIcon(), login)
	row := container.NewGridWithRows(1, loginBtn, updateBtn)
	return row
}

func projectsLength() int {
	return len(Projects)
}

func createProjectsItem() fyne.CanvasObject {
	style := fyne.TextStyle{Bold: true}
	title := widget.NewLabelWithStyle("", fyne.TextAlignLeading, style)
	return container.NewGridWithRows(1, title)
}

func updateProjectsItem(i int, o fyne.CanvasObject) {
	item := Projects[i]
	o.(*fyne.Container).Objects = nil

	style := fyne.TextStyle{Bold: true}
	title := widget.NewLabelWithStyle(item.Name, fyne.TextAlignLeading, style)
	o.(*fyne.Container).Add(title)
}

func updateProjects() {
	projectLoader.Start()
	if !isLogged() {
		login()
	}
	setProjects()
	projectLoader.Stop()
}

func login() {
	cmd := exec.Command("gnome-terminal", "--wait", "--", "bash", "-c", "platform auth:browser-login")
	cmd.Run()
}

func isLogged() bool {
	out, err := exec.Command("bash", "-c", "platform auth:info --no-auto-login --format csv").Output()
	if err != nil || len(out) <= 0 {
		return false
	}
	return true
}

func setProjects() {
	out, _ := exec.Command("bash", "-c", "platform project:list --format csv").Output()
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	scanner.Scan()
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		p := Project{Name: s[1], Id: s[0], Environments: getEnvironments(s[0])}
		Projects = append(Projects, p)
	}
	fmt.Println(Projects)
}

func getEnvironments(id string) []Environment {
	var environments []Environment
	c := fmt.Sprintf("platform environment:list -p %s -I --format csv", id)
	out, _ := exec.Command("bash", "-c", c).Output()
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	scanner.Scan()
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		e := Environment{Name: s[1], Id: s[0]}
		environments = append(environments, e)
	}
	return environments
}
