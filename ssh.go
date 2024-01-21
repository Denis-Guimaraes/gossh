package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type Ssh struct {
	Name     string
	Host     string
	Port     int
	User     string
	Password string
}

type SshList struct {
	Ssh []Ssh
}

var MyList SshList

func load() {
	fs, _ := os.Stat(getFilepath())
	if fs.Size() == 0 {
		newList := SshList{
			Ssh: []Ssh{},
		}
		writeSshData(newList)
	}
	jsonFile, _ := os.Open(getFilepath())
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &MyList)
}

func writeSshData(sshList SshList) {
	file, _ := os.OpenFile(getFilepath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	jsonData, _ := json.Marshal(sshList)
	file.Write(jsonData)
	file.Close()
}

func getFilepath() string {
	// Crée le dossier .gossh si il n'existe pas
	home, _ := os.UserHomeDir()
	folderPath := filepath.Join(home, ".gossh")
	os.Mkdir(folderPath, 0700)
	// Crée le fichier ssh.json si il n'existe pas
	filePath := filepath.Join(folderPath, "ssh.json")
	file, _ := os.OpenFile(filePath, os.O_CREATE, 0600)
	file.Close()
	return filePath
}
