package main

import (
	"fmt"
	"os"
	"strconv"

	"./helpers"
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
	case "help":
		CommandHelp()
	default:
		CommandHelp()
	}
}

// CommandList : shows projects
func CommandList(config GlobalConfig) {
	for _, project := range config.Projects {
		isActive := IsActiveProject(project)
		activePort, _ := GetActiveProjectPort(project)
		var activePortString string
		if isActive {
			activePortString = helpers.SuccessText(strconv.Itoa(activePort))
		} else {
			activePortString = helpers.FailText(strconv.Itoa(activePort))
		}
		fmt.Printf(
			"  %v\t%v:%v\tactive: %v\n",
			helpers.SuccessText(project.ProjectName),
			project.DefaultConfig.NginxDomain,
			strconv.Itoa(project.DefaultConfig.NginxPort)+" ("+activePortString+")",
			isActive,
		)
	}
}

// CommandStart : starts project
func CommandStart(config GlobalConfig) {
	var unactiveProjectNames []string
	for _, project := range config.Projects {
		if !IsActiveProject(project) {
			unactiveProjectNames = append(unactiveProjectNames, project.ProjectName)
		}
	}
	if len(unactiveProjectNames) > 0 {
		selectedProjectName := helpers.SelectStringVariant("Select project to activate", unactiveProjectNames)
		for _, project := range config.Projects {
			if project.ProjectName == selectedProjectName {
				StartProject(project)
				ReloadEnvironment()
				break
			}
		}
	} else {
		fmt.Println(
			helpers.FailText("Hasn't projects to activate"),
		)
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
	if len(activeProjectNames) > 0 {
		selectedProjectName := helpers.SelectStringVariant("Select project to deactivate", activeProjectNames)
		for _, project := range config.Projects {
			if project.ProjectName == selectedProjectName {
				StopProject(project)
				ReloadEnvironment()
				break
			}
		}
	} else {
		fmt.Println(
			helpers.FailText("Hasn't projects to deactivate"),
		)
	}
}

// CommandHelp : shows all available commands
func CommandHelp() {
	fmt.Println("List of avalaible commands:")
	fmt.Println("  ", "list\t\tlist of projects in active configuration")
	fmt.Println("  ", "start\tstart unactive project (you can select project interactively)")
	fmt.Println("  ", "stop\tstop active project (you can select project interactively)")
	fmt.Println("  ", "recreate\trecreate configuration file")
}
