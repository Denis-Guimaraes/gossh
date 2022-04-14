package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var name = widget.NewEntry()
var host = widget.NewEntry()
var port = widget.NewEntry()
var user = widget.NewEntry()
var password = widget.NewEntry()

func add() *container.TabItem {
	content := form()
	return container.NewTabItemWithIcon("Ajouter", theme.ContentAddIcon(), content)
}

func form() *fyne.Container {
	nameItem := widget.NewFormItem("Nom", name)
	hostItem := widget.NewFormItem("Hôte", host)
	portItem := widget.NewFormItem("Port", port)
	userItem := widget.NewFormItem("Utilisateur", user)
	passwordItem := widget.NewFormItem("Mot de passe", password)
	form := widget.NewForm(nameItem, hostItem, portItem, userItem, passwordItem)
	form.OnSubmit = submit
	return container.NewVBox(form)
}

func submit() {
	sshPort, err := strconv.Atoi(port.Text)
	if err != nil || sshPort == 0 {
		sshPort = 22
	}
	ssh := Ssh{
		Name:     name.Text,
		Host:     host.Text,
		Port:     sshPort,
		User:     user.Text,
		Password: password.Text,
	}
	MyList.Ssh = append(MyList.Ssh, ssh)
	writeSshData(MyList)
	name.SetText("")
	name.Refresh()
	host.SetText("")
	host.Refresh()
	port.SetText("")
	port.Refresh()
	user.SetText("")
	user.Refresh()
	password.SetText("")
	password.Refresh()
}
