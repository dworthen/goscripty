package generators

import (
	_ "embed"
	"fmt"
)

//go:embed templates/github/bash-install.go.tmpl
var ghTemplate string

func GenerateGithubInstallers() {
	fmt.Println(ghTemplate)
}
