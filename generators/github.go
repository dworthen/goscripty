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
var bashTemplate string

//go:embed templates/github/powershell-install.go.tmpl
var pwshTemplate string

type PromptVariables struct {
	RepoUrl             string
	ArchiveName         string
	LinuxPlatformDict   string
	LinuxArchDict       string
	WindowsPlatformDict string
	WindowsArchDict     string
	Destination         string
}

func NewPromptVariables() *PromptVariables {
	variables := new(PromptVariables)
	variables.RepoUrl = promptMessage("Provide the GitHub project in the format USERNAME/REPO", "", true)
	printSpacer()
	variables.ArchiveName = promptMessage(strings.TrimSpace(`
Provide the archive template name. The template may make use of the following variables:

  {PLATFORM}: The OS of the machine that results from running 'uname -s' on linux/darwin machines
    and is the literal value 'Windows' on windows.
  {ARCH}: The architecture of the machine that results from running 'uname -m' on linux/darwin
    machines and $ENV:PPROCESSOR_ACHITECTURE on windows.
  {TAG}: The specified GitHub release tag. Defaults to the latest tag if one was not specified
    when calling the install script.
  {ARCHIVE_EXT}: An extension specific to the PLATFORM. Defaults to '.tar.gz' on linux/darwin
    and '.zip' on windows.
  {BINARY_EXT}: The binary extension specific to the PLATFORM. Defaults to '' on linux/darwin
    and '.exe' on windows.

Example:
  myCoolApp_{PLATFORM}_{ARCH}_{TAG}{ARCHIVE_EXT}
Which might result in
  myCoolApp_Windows_x86_64_v1.0.5.zip when installing on windows with the produced install.ps1 script

Archive Name
  `), "", true)
	printSpacer()
	variables.LinuxPlatformDict = promptMessage(strings.TrimSpace(`
Provide a comma separated list of KEY=VALUE Platform mapping to use on Linux/Darwin machines.
This mapping helps bridge the gap from what is produced by 'uname -s | cut -d- -f1' and what
might be listed in GitHub releases. For example, the value produced on the machine might be the OS
capitalized, e.g., 'Linux' while the value listed in GitHub actions might be with a lowercase value,
'linux'. Or the project on GitHub might use the value 'mac' in their release names instead of
the OS value 'Darwin'. In this case, the mapping "Linux=linux,Darwin=mac" would be helpful.

Default: Linux=Linux,Darwin=Darwin

Linux/Darwin Platform Mapping
`), "Linux=Linux,Darwin=Darwin", false)
	printSpacer()
	variables.LinuxArchDict = promptMessage(strings.TrimSpace(`
Provide a comma separated list of KEY=VALUE architecture mapping for Linux/Darwin platforms.
Similar to the Platform mapping, this helps bridge the gap of what is produced on the machine
when the install script runs and the value that is used on GitHub releases. For example, the
machine might list the architecture 'x86_64' when running 'uname -m' but the GitHub project
might be using the value 'amd64' in their release names. In that case the mapping x86_64=amd64
would be helpful.

Default: x86_64=x86_64

Linux/Darwin Architecture Mapping
  `), "x86_64=x86_64", false)
	printSpacer()
	variables.WindowsPlatformDict = promptMessage(strings.TrimSpace(`
Provide a comma separated list of possible KEY=VALUE platform mappings to use on Windows.

Default: Windows=Windows

Windows Platform Mapping
  `), "Windows=Windows", false)
	printSpacer()
	variables.WindowsArchDict = promptMessage(strings.TrimSpace(`
Provide a comma separated list of possible KEY=VALUE architecture mappings to use on Windows.
The key is derived from the environment variable PROCESSOR_ARCHITECTURE. By default, we map
the AMD64 value to x86_64 to be more consistent with what is produced by 'uname -m' on
Linux/Darwin machines.

Default: AMD64=x86_64

Windows Architecture Mapping
  `), "AMD64=x86_64", false)
	printSpacer()
	variables.Destination = promptMessage("Location to create install scripts (defaults to '.')", ".", false)
	return variables
}

func (variables *PromptVariables) Validate() {
	repoUrlSlashCount := strings.Count(variables.RepoUrl, "/")
	if repoUrlSlashCount != 1 {
		fmt.Printf("Error: Expecting GitHub repo project to be in the form USERNAME/REPO but recieved %v", variables.RepoUrl)
		os.Exit(1)
	}
	variables.RepoUrl = fmt.Sprintf("https://github.com/%v", variables.RepoUrl)

	variables.LinuxPlatformDict = mapToBashDict(keyValueToMap(variables.LinuxPlatformDict))
	variables.LinuxArchDict = mapToBashDict(keyValueToMap(variables.LinuxArchDict))
	variables.WindowsPlatformDict = mapToPwshDict(keyValueToMap(variables.WindowsPlatformDict))
	variables.WindowsArchDict = mapToPwshDict(keyValueToMap(variables.WindowsArchDict))
	dest, err := filepath.Abs(utils.ToFilePath(variables.Destination))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	variables.Destination = dest
}

func GenerateGithubInstallers() {
	variables := NewPromptVariables()
	variables.Validate()

	os.MkdirAll(variables.Destination, 0744)

	writeTemplate(variables, bashTemplate, "install.sh")
	writeTemplate(variables, pwshTemplate, "install.ps1")
}

func writeTemplate(variables *PromptVariables, templateString string, scriptName string) {
	tmpl, err := template.New("template").Parse(templateString)
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

	dest := filepath.Join(variables.Destination, scriptName)
	if err := os.WriteFile(dest, buf.Bytes(), 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func promptMessage(message string, defaultValue string, required bool) string {
	var response string
	for {
		response = prompter.Prompt(message, defaultValue)
		if !required || response != "" {
			break
		}
	}
	return response
}

func printSpacer() {
	fmt.Println()
	fmt.Print("=====================================================")
	fmt.Println()
	fmt.Println()
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

func mapToBashDict(mapping map[string]string) string {
	output := ""
	for key, value := range mapping {
		output += fmt.Sprintf("\n\t[\"%s\"]=\"%s\"", key, value)
	}
	return output
}

func mapToPwshDict(mapping map[string]string) string {
	output := ""
	for key, value := range mapping {
		output += fmt.Sprintf("\n\t%s = \"%s\"", key, value)
	}
	return output
}
