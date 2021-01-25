package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"./helpers"
)

// LoadDockerComposeConfigs : gets configuration for all projects
func LoadDockerComposeConfigs(root string) ([]ProjectConfig, error) {
	var projects []ProjectConfig
	rootDirInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return projects, err
	}
	for _, dir := range rootDirInfo {
		if dir.IsDir() {
			dirInfo, err := ioutil.ReadDir(root + "/" + dir.Name())
			if err != nil {
				return projects, err
			}
			hasComposeFile := false
			hasEnvironmentFile := false
			for _, file := range dirInfo {
				if file.Name() == "docker-compose.yml" {
					hasComposeFile = true
				}
				if file.Name() == ".env" {
					hasEnvironmentFile = true
				}
			}
			if hasComposeFile && hasEnvironmentFile {
				projectName := dir.Name()
				projectPath := root + "/" + dir.Name()
				configPath := root + "/" + dir.Name() + "/.env"
				config, _ := GetProjectEnvironment(configPath)
				project := ProjectConfig{projectName, projectPath, configPath, config, false}
				projects = append(projects, project)
			}
		}
	}
	return projects, nil
}

// GetProjectEnvironment : get environment from .env file in project directory
func GetProjectEnvironment(filePath string) (EnvironmentConfig, error) {
	var config EnvironmentConfig
	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var siteDomain string
	var nginxPort string
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), siteDomainVariable+"=") {
			siteDomain = strings.Trim(
				strings.ReplaceAll(scanner.Text(), siteDomainVariable+"=", ""), " ",
			)
		}
		if strings.HasPrefix(scanner.Text(), nginxPortVariable+"=") {
			nginxPort = strings.Trim(
				strings.ReplaceAll(scanner.Text(), nginxPortVariable+"=", ""), " ",
			)
		}
	}
	nginxPortInt, err := strconv.Atoi(nginxPort)
	config = EnvironmentConfig{siteDomain, nginxPortInt}
	return config, nil
}

// IsActiveProject : Check project for active on host machine
func IsActiveProject(project ProjectConfig) bool {
	isActiveProject := false
	composeCmd := exec.Command("docker-compose", "ps", "-q")
	composeCmd.Dir = project.ProjectPath
	composeCmdOut, _ := composeCmd.Output()
	projectContainerIds := strings.Split(
		strings.Trim(string(composeCmdOut), "\n"),
		"\n",
	)

	dockerCmd := exec.Command("docker", "ps", "--filter=status=running", "--no-trunc", "-q")
	dockerCmdOut, _ := dockerCmd.Output()
	runningIds := strings.Split(
		strings.Trim(string(dockerCmdOut), "\n"),
		"\n",
	)

	for _, projectContainerID := range projectContainerIds {
		if helpers.Contains(runningIds, projectContainerID) {
			isActiveProject = true
		}
	}
	return isActiveProject
}
