package main

import (
	"bufio"
	"fmt"
	"os"
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
	top := container.NewVBox(action, projectLoader)
	wListPlatform = widget.NewList(projectsLength, createProjectsItem, updateProjectsItem)
	return container.NewBorder(top, nil, nil, nil, wListPlatform)
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
	login := container.NewAdaptiveGrid(3)
	return container.NewGridWithRows(2, title, login)
}

func updateProjectsItem(i int, o fyne.CanvasObject) {
	item := Projects[i]
	o.(*fyne.Container).Objects = nil

	style := fyne.TextStyle{Bold: true}
	title := widget.NewLabelWithStyle(item.Name, fyne.TextAlignLeading, style)
	login := container.NewAdaptiveGrid(3)
	for _, e := range item.Environments {
		login.Add(widget.NewButtonWithIcon(e.Name, theme.LoginIcon(), startPlatformSsh(item.Id, e.Id)))
	}
	o.(*fyne.Container).Add(title)
	o.(*fyne.Container).Add(login)
}

func updateProjects() {
	projectLoader.Show()
	projectLoader.Start()
	if !isLogged() {
		login()
	}
	setProjects()
	projectLoader.Stop()
	projectLoader.Hide()
}

func platformPath() string {
	home, _ := os.UserHomeDir()
	platformPath := fmt.Sprintf("%s/.platformsh/bin/platform", home)
	_, err := os.Stat(platformPath)

	if err != nil {
		platformPath = fmt.Sprintf("%s/.local/bin/platform", home)
	}
	return platformPath
}

func login() {
	command := fmt.Sprintf("%s auth:browser-login", platformPath())
	fmt.Println(command)
	cmd := exec.Command("gnome-terminal", "--wait", "--", "bash", "-c", command)
	e := cmd.Run()
	fmt.Println(e)
}

func isLogged() bool {
	command := fmt.Sprintf("%s auth:info --no-auto-login --format csv", platformPath())
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil || len(out) <= 0 {
		return false
	}
	return true
}

func setProjects() {
	Projects = nil
	command := fmt.Sprintf("%s project:list --format csv --count 0", platformPath())
	out, _ := exec.Command("bash", "-c", command).Output()
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	scanner.Scan()
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		p := Project{Name: s[1], Id: s[0], Environments: getEnvironments(s[0])}
		Projects = append(Projects, p)
		wListPlatform.Refresh()
	}
}

func getEnvironments(id string) []Environment {
	var environments []Environment
	c := fmt.Sprintf("%s environment:list -p %s -I --format csv", platformPath(), id)
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

func startPlatformSsh(projectId string, environmentId string) func() {
	return func() {
		command := fmt.Sprintf("%s ssh  -p %s -e %s", platformPath(), projectId, environmentId)
		cmd := exec.Command("gnome-terminal", "--", "bash", "-c", command)
		cmd.Start()
	}
}
