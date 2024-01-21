package project

import (
	"bufio"
	"encoding/json"
	"gossh/internal/platformsh"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	Name         string
	Id           string
	Environments []Environment
}

type Projects struct {
	Items []Project
}

func (p *Projects) Count() int {
	return len(p.Items)
}

func GetProjects() Projects {
	if hasCache() {
		go refreshCache()
		return getCache()
	}
	return RefreshProjects()
}

func RefreshProjects() Projects {
	return refreshCache()
}

func hasCache() bool {
	fs, _ := os.Stat(getCachePath())
	return fs.Size() > 0
}

func getCache() Projects {
	var projects Projects
	jsonFile, _ := os.Open(getCachePath())
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &projects)
	return projects
}

func setCache(projects Projects) {
	file, _ := os.OpenFile(getCachePath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	jsonData, _ := json.Marshal(projects)
	file.Write(jsonData)
	file.Close()
}

func refreshCache() Projects {
	var projects []Project
	scanner := bufio.NewScanner(strings.NewReader(string(platformsh.ProjectList())))
	scanner.Scan()

	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		p := Project{Name: s[1], Id: s[0], Environments: getEnvironments(s[0])}
		projects = append(projects, p)
	}

	p := Projects{Items: projects}
	setCache(p)
	return p
}

func getCachePath() string {
	// Crée le dossier .gossh si il n'existe pas
	home, _ := os.UserHomeDir()
	folderPath := filepath.Join(home, ".gossh")
	os.Mkdir(folderPath, 0700)
	// Crée le fichier ssh.json si il n'existe pas
	filepath := filepath.Join(folderPath, "projects.json")
	file, _ := os.OpenFile(filepath, os.O_CREATE, 0600)
	file.Close()
	return filepath
}
