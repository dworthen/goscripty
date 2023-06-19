package generators

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Songmu/prompter"
	"github.com/dworthen/goscripty/utils"
)

//go:embed templates/github/bash-install.go.tmpl
var ghTemplate string

func GenerateGithubInstallers() {
	variables := prompt()
	shapeVariables(variables)

	tmpl, err := template.New("bash").Parse(ghTemplate)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, variables)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dest := filepath.Join(variables["dest"], "install.sh")
	os.WriteFile(dest, []byte(buf.String()), 0755)
}

func prompt() map[string]string {
	variables := map[string]string{}
	promptMessage(variables, "repoUrl", "repoUrl", true)
	promptMessage(variables, "archiveName", "Archive template", true)
	promptMessage(variables, "platformDict", "Platform map (optional, comma separated key=value. Example: win=Windows)", false)
	promptMessage(variables, "archDict", "Architecture map (Optional, comma separated key=value. Example: x86_64=amd64)", false)
	promptMessage(variables, "dest", "Destination directory (Optional, defaults to '.')", false)
	return variables
}

func promptMessage(variables map[string]string, key string, message string, required bool) {
	for {
		variables[key] = prompter.Prompt(message, "")
		if val, ok := variables[key]; !required || (ok && val != "") {
			break
		}
	}
}

func shapeVariables(variables map[string]string) {
	var platformMapping map[string]string
	var archMapping map[string]string

	if val, ok := variables["platformDict"]; ok && val != "" {
		platformMapping = keyValueToMap(val)
	} else {
		platformMapping = getDefaultPlatformMapping()
	}

	if val, ok := variables["archDict"]; ok && val != "" {
		archMapping = keyValueToMap(val)
	} else {
		archMapping = getDefaultArchMapping()
	}

	if val, ok := variables["dest"]; !ok || val == "" {
		variables["dest"] = "."
	}

	val, err := filepath.Abs(utils.ToFilePath(variables["dest"]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	variables["dest"] = val

	variables["platformDict"] = mapToBashDict(platformMapping)
	variables["archDict"] = mapToBashDict(archMapping)
}

func keyValueToMap(keyVals string) map[string]string {
	listKeyVals := strings.Split(keyVals, ",")
	mapping := map[string]string{}

	for _, val := range listKeyVals {
		if val == "" {
			break
		}
		equalCount := strings.Count(val, "=")
		if equalCount != 1 {
			fmt.Printf("Error: malformed list of key=value pairs. Expecting comma separated list of key=value pairs.")
			os.Exit(1)
		}
		kv := strings.Split(val, "=")
		mapping[kv[0]] = kv[1]
	}

	return mapping
}

func getDefaultPlatformMapping() map[string]string {
	return map[string]string{
		"Linux":  "Linux",
		"Darwin": "Darwin",
	}
}

func getDefaultArchMapping() map[string]string {
	return map[string]string{
		"x86_64": "x86_64",
	}
}

func mapToBashDict(mapping map[string]string) string {
	output := ""
	for key, value := range mapping {
		output += fmt.Sprintf("\n\t[\"%s\"]=\"%s\"", key, value)
	}
	return output
}
