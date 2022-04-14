package main

import (
	"encoding/json"
	"io/ioutil"
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
	byteValue, _ := ioutil.ReadAll(jsonFile)
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
	folderpath := filepath.Join(home, ".gossh")
	os.Mkdir(folderpath, 0700)
	// Crée le fichier ssh.json si il n'existe pas
	filepath := filepath.Join(folderpath, "ssh.json")
	file, _ := os.OpenFile(filepath, os.O_CREATE, 0600)
	file.Close()
	return filepath
}
