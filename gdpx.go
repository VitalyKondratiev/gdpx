package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"text/tabwriter"

	"./helpers"
	rice "github.com/GeertJohan/go.rice"
)

const siteDomainVariable = "SITE_DOMAIN"
const nginxPortVariable = "NGINX_PORT"
const nginxSslPortVariable = "NGINX_PORT_SSL"

func main() {

	if len(os.Args) < 2 {
		CommandHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "list":
		config, _ := LoadMainConfig()
		CommandList(config)
	case "start":
		config, _ := LoadMainConfig()
		CommandStart(config)
	case "stop":
		config, _ := LoadMainConfig()
		CommandStop(config)
	case "recreate":
		config, _ := CreateMainConfig()
		CommandList(config)
	case "install-autocomplete":
		CommandInstallAutocomplete()
	case "help":
		CommandHelp()
	default:
		CommandHelp()
	}
}

// CommandList : shows projects
func CommandList(config GlobalConfig) {
	var showStarted = true
	var showStopped = true
	if len(os.Args) == 3 && os.Args[2] == "--started" {
		showStopped = false
	} else if len(os.Args) == 3 && os.Args[2] == "--stopped" {
		showStarted = false
	}

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)
	for _, project := range config.Projects {
		isActive := IsActiveProject(project)
		activePort, _ := GetActiveProjectPort(project)
		var colorProjectName string
		if isActive {
			colorProjectName = helpers.SuccessText(project.ProjectName)
			if !showStarted {
				continue
			}
		} else {
			colorProjectName = helpers.FailText(project.ProjectName)
			if !showStopped {
				continue
			}
		}
		line := fmt.Sprintf(
			"\t%v\thttp://%v:%v\t",
			colorProjectName,
			project.DefaultConfig.NginxDomain,
			strconv.Itoa(project.DefaultConfig.NginxPort)+" ("+strconv.Itoa(activePort)+")",
		)
		fmt.Fprintln(w, line)
	}

	w.Flush()
}

// CommandStart : starts project
func CommandStart(config GlobalConfig) {
	var unactiveProjectNames []string
	for _, project := range config.Projects {
		if !IsActiveProject(project) {
			unactiveProjectNames = append(unactiveProjectNames, project.ProjectName)
		}
	}
	var selectedProjectName string
	if len(os.Args) == 3 && helpers.InArray(os.Args[2], unactiveProjectNames) {
		selectedProjectName = os.Args[2]
	} else if len(unactiveProjectNames) > 0 {
		selectedProjectName = helpers.SelectStringVariant("Select project to activate", unactiveProjectNames)
	} else {
		fmt.Println(
			helpers.FailText("Hasn't projects to activate"),
		)
		return
	}
	for _, project := range config.Projects {
		if project.ProjectName == selectedProjectName {
			StartProject(project)
			ReloadEnvironment()
			break
		}
	}
}

// CommandStop : stops project
func CommandStop(config GlobalConfig) {
	var activeProjectNames []string
	for _, project := range config.Projects {
		if IsActiveProject(project) {
			activeProjectNames = append(activeProjectNames, project.ProjectName)
		}
	}
	var selectedProjectName string
	if len(os.Args) == 3 && helpers.InArray(os.Args[2], activeProjectNames) {
		selectedProjectName = os.Args[2]
	} else if len(activeProjectNames) > 0 {
		selectedProjectName = helpers.SelectStringVariant("Select project to deactivate", activeProjectNames)
	} else {
		fmt.Println(
			helpers.FailText("Hasn't projects to deactivate"),
		)
		return
	}
	for _, project := range config.Projects {
		if project.ProjectName == selectedProjectName {
			StopProject(project)
			ReloadEnvironment()
			break
		}
	}
}

// CommandHelp : shows all available commands
func CommandHelp() {

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)
	fmt.Println("List of avalaible commands:")

	fmt.Fprintln(w, fmt.Sprintf(
		"\t%v\t%v\t",
		"list",
		"list of projects in active configuration",
	))
	fmt.Fprintln(w, fmt.Sprintf(
		"\t%v\t%v\t",
		"start",
		"start unactive project (you can select project interactively)",
	))
	fmt.Fprintln(w, fmt.Sprintf(
		"\t%v\t%v\t",
		"stop",
		"stop active project (you can select project interactively)",
	))
	fmt.Fprintln(w, fmt.Sprintf(
		"\t%v\t%v\t",
		"recreate",
		"recreate configuration file",
	))
	fmt.Fprintln(w, fmt.Sprintf(
		"\t%v\t%v\t",
		"install-autocomplete",
		"install bash autocomlpetion(need privelegies)",
	))
	w.Flush()
}

// CommandInstallAutocomplete : install autocomplete.sh to /etc/bash_completion.d/gdpx
func CommandInstallAutocomplete() {
	cmd := exec.Command("id", "-u")
	output, _ := cmd.Output()
	idU, _ := strconv.Atoi(string(output[:len(output)-1]))
	if idU == 0 {
		conf := rice.Config{
			LocateOrder: []rice.LocateMethod{rice.LocateEmbedded, rice.LocateAppended, rice.LocateFS},
		}
		box := conf.MustFindBox("./shell")
		autocompleteString, _ := box.String("autocomplete.sh")
		autocompleteFile, _ := os.Create("/etc/bash_completion.d/gdpx")
		defer autocompleteFile.Close()
		autocompleteFile.WriteString(autocompleteString)
		fmt.Println(helpers.SuccessText("Awesome! You can relaunch your terminal session and get autocomplete superpower!"))
	} else {
		fmt.Println(helpers.FailText("This command must be run as root! (sudo)"))
	}
}
