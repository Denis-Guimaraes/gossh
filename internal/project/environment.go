package project

import (
	"bufio"
	"gossh/internal/platformsh"
	"strings"
)

type Environment struct {
	Name string
	Id   string
}

func getEnvironments(projectId string) []Environment {
	var environments []Environment
	scanner := bufio.NewScanner(strings.NewReader(string(platformsh.EnvironmentList(projectId))))
	scanner.Scan()

	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		e := Environment{Name: s[1], Id: s[0]}
		environments = append(environments, e)
	}
	return environments
}
