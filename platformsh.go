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

func platform() *container.TabItem {
	update()
	content := create()
	return container.NewTabItemWithIcon("Platform", theme.MenuIcon(), content)
}

func create() *fyne.Container {
	action := container.NewVBox(action())
	list := widget.NewList(projectsLength, createProjectsItem, updateProjectsItem)
	return container.NewVBox(action, container.NewMax(list))
}

func action() *fyne.Container {
	updateBtn := widget.NewButtonWithIcon("Mettre à jour", theme.ViewRefreshIcon(), update)
	loginBtn := widget.NewButtonWithIcon("Connexion", theme.LoginIcon(), login)
	row := container.NewGridWithRows(1, loginBtn, updateBtn)
	return container.NewVBox(row)
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

func update() {
	if !isLogged() {
		login()
	}
	setProjects()
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
		e := getEnvironments(s[0])
		p := Project{Name: s[1], Id: s[0], Environments: e}
		Projects = append(Projects, p)
	}
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
