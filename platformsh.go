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
	return container.NewVBox(action())
}

func action() *fyne.Container {
	updateBtn := widget.NewButtonWithIcon("Mettre à jour", theme.ViewRefreshIcon(), update)
	loginBtn := widget.NewButtonWithIcon("Connexion", theme.LoginIcon(), login)
	row := container.NewGridWithRows(1, loginBtn, updateBtn)
	return container.NewVBox(row)
}

func update() {
	out, err := exec.Command("bash", "-c", "platform project:list --format csv").Output()
	if err != nil {
		login()
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	scanner.Scan()
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		p := Project{Name: s[1], Id: s[0]}
		Projects = append(Projects, p)
	}
	fmt.Println(Projects)
}

func login() {
	cmd := exec.Command("gnome-terminal", "--", "bash", "-c", "platform auth:browser-login")
	cmd.Start()
}
