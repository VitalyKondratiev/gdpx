package main

import (
	"os"
	"fmt"
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
			CommandList()
		case "recreate":
			CreateMainConfig()
			CommandList()
		case "help":
			CommandHelp()
		default:
			CommandHelp()
	}
}

func CommandList() {
	config, _ := LoadMainConfig()
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

func CommandHelp()  {
	fmt.Println("List of avalaible commands:")
	fmt.Println("  ", "list\t\tlist of projects in active configuration")
	fmt.Println("  ", "recreate\trecreate configuration file")
}
