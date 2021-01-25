package main

import (
	"fmt"
	"os"

	"./helpers"
)

const siteDomainVariable = "SITE_DOMAIN"
const nginxPortVariable = "NGINX_PORT"

func main() {

	if len(os.Args) < 2 {
		CommandHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "list":
		config, _ := LoadMainConfig()
		CommandList(config)
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
		fmt.Printf(
			"  %v\t%v:%v\tactive: %v\n",
			helpers.SuccessText(project.ProjectName),
			project.DefaultConfig.NginxDomain,
			project.DefaultConfig.NginxPort,
			IsActiveProject(project),
		)
	}
}

// CommandHelp : shows all available commands
func CommandHelp() {
	fmt.Println("List of avalaible commands:")
	fmt.Println("  ", "list\t\tlist of projects in active configuration")
	fmt.Println("  ", "recreate\trecreate configuration file")
}
