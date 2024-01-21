package platformsh

import (
	"fmt"
	"os"
	"os/exec"
)

func Login() {
	command := fmt.Sprintf("%s auth:browser-login", platformPath())
	fmt.Println(command)
	cmd := exec.Command("gnome-terminal", "--wait", "--", "bash", "-c", command)
	e := cmd.Run()
	fmt.Println(e)
}

func Ssh(projectId string, environmentId string) func() {
	return func() {
		if !isLogged() {
			Login()
		}
		command := fmt.Sprintf("%s ssh  -p %s -e %s", platformPath(), projectId, environmentId)
		cmd := exec.Command("gnome-terminal", "--", "bash", "-c", command)
		cmd.Start()
	}
}

func ProjectList() []byte {
	if !isLogged() {
		Login()
	}
	command := fmt.Sprintf("%s project:list --format csv --count 0", platformPath())
	out, _ := exec.Command("bash", "-c", command).Output()
	return out
}

func EnvironmentList(projectId string) []byte {
	if !isLogged() {
		Login()
	}
	c := fmt.Sprintf("%s environment:list -p %s -I --format csv", platformPath(), projectId)
	out, _ := exec.Command("bash", "-c", c).Output()
	return out
}

func isLogged() bool {
	command := fmt.Sprintf("%s auth:info --no-auto-login --format csv", platformPath())
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil || len(out) <= 0 {
		return false
	}
	return true
}

func platformPath() string {
	home, _ := os.UserHomeDir()
	platformPath := fmt.Sprintf("%s/.platformsh/bin/platform", home)
	_, err := os.Stat(platformPath)

	if err != nil {
		platformPath = fmt.Sprintf("%s/.local/bin/platform", home)
		_, err = os.Stat(platformPath)
	}

	if err != nil {
		platformPath = "/usr/bin/platform"
	}
	return platformPath
}
