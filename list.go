package main

import (
	"fmt"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var wList *widget.List

func list(w fyne.Window) *container.TabItem {
	wList = widget.NewList(length, createItem, updateItem)
	return container.NewTabItemWithIcon("Connexion", theme.ListIcon(), wList)
}

func length() int {
	return len(MyList.Ssh)
}

func createItem() fyne.CanvasObject {
	style := fyne.TextStyle{Bold: true}
	title := widget.NewLabelWithStyle("", fyne.TextAlignLeading, style)
	loginBtn := widget.NewButtonWithIcon("Connexion", theme.LoginIcon(), func() { fmt.Println("Connexion") })
	deleteBtn := widget.NewButtonWithIcon("Supprimer", theme.DeleteIcon(), func() { fmt.Println("Supprimer") })
	return container.NewGridWithRows(1, title, loginBtn, deleteBtn)
}

func updateItem(i int, o fyne.CanvasObject) {
	item := MyList.Ssh[i]
	o.(*fyne.Container).Objects = nil

	style := fyne.TextStyle{Bold: true}
	title := widget.NewLabelWithStyle(item.Name, fyne.TextAlignLeading, style)
	loginBtn := widget.NewButtonWithIcon("Connexion", theme.LoginIcon(), func() { loginSsh(i) })
	deleteBtn := widget.NewButtonWithIcon("Supprimer", theme.ContentCopyIcon(), func() { removeSsh(i) })
	o.(*fyne.Container).Add(title)
	o.(*fyne.Container).Add(loginBtn)
	o.(*fyne.Container).Add(deleteBtn)
}

func loginSsh(i int) {
	item := MyList.Ssh[i]
	command := fmt.Sprintf(
		"sshpass -p '%s' ssh -o StrictHostKeyChecking=no %s@%s -p %d",
		item.Password,
		item.User,
		item.Host,
		item.Port,
	)

	cmd := exec.Command("gnome-terminal", "--", "bash", "-c", command)
	cmd.Start()
}

func removeSsh(index int) {
	newSsh := []Ssh{}
	for i := 0; i < len(MyList.Ssh); i++ {
		if i != index {
			newSsh = append(newSsh, MyList.Ssh[i])
		}
	}
	MyList.Ssh = newSsh
	writeSshData(MyList)
	wList.Refresh()
}
