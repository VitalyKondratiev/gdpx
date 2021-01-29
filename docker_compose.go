package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"./helpers"
	rice "github.com/GeertJohan/go.rice"
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
	var nginxSslPort string
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
		if strings.HasPrefix(scanner.Text(), nginxSslPortVariable+"=") {
			nginxSslPort = strings.Trim(
				strings.ReplaceAll(scanner.Text(), nginxSslPortVariable+"=", ""), " ",
			)
		}
	}
	nginxPortInt, err := strconv.Atoi(nginxPort)
	nginxSslPortInt, err := strconv.Atoi(nginxSslPort)
	config = EnvironmentConfig{siteDomain, nginxPortInt, nginxSslPortInt}
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

	if len(strings.Join(projectContainerIds, "")) == 0 {
		return false
	}

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

// GetActiveProjectPort : get port for active project
func GetActiveProjectPort(project ProjectConfig) (int, error) {
	envFile, _ := os.Open(project.ConfigPath)
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), nginxPortVariable+"=") {
			return strconv.Atoi(strings.Trim(
				strings.ReplaceAll(scanner.Text(), nginxPortVariable+"=", ""), " ",
			))
		}
	}
	return project.DefaultConfig.NginxPort, nil
}

// StartProject : prepare and launch project docker-compose
func StartProject(project ProjectConfig) {

	envFile, _ := os.Open(project.ConfigPath)
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)

	var newEnvContent []string

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), nginxPortVariable+"=") {
			newNginxPort := helpers.GetOpenedPort()
			oldNginxPort := strings.Trim(
				strings.ReplaceAll(scanner.Text(), nginxPortVariable+"=", ""), " ",
			)
			nginxPortString := strings.Replace(scanner.Text(), oldNginxPort, strconv.Itoa(newNginxPort), -1)
			newEnvContent = append(newEnvContent, nginxPortString)
		} else if strings.HasPrefix(scanner.Text(), nginxSslPortVariable+"=") {
			newNginxSslPort := helpers.GetOpenedPort()
			oldNginxSslPort := strings.Trim(
				strings.ReplaceAll(scanner.Text(), nginxSslPortVariable+"=", ""), " ",
			)
			nginxSslPortString := strings.Replace(scanner.Text(), oldNginxSslPort, strconv.Itoa(newNginxSslPort), -1)
			newEnvContent = append(newEnvContent, nginxSslPortString)
		} else {
			newEnvContent = append(newEnvContent, scanner.Text())
		}
	}

	newEnvFile, _ := os.Create(project.ConfigPath)
	newEnvFile.Truncate(0)
	newEnvFile.WriteString(strings.Join(newEnvContent, "\n"))
	args := []string{"up", "--force-recreate", "--build", "-d"}
	fmt.Println(
		helpers.SuccessText("("+project.ProjectName+"):"),
		"docker-compose", strings.Join(args, " "),
	)
	composeCmd := exec.Command("docker-compose", args...)
	composeCmd.Dir = project.ProjectPath
	composeCmd.Output()
}

// StopProject : stop docker-compose, and return defaults env settings
func StopProject(project ProjectConfig) {
	envFile, _ := os.Open(project.ConfigPath)
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)

	var newEnvContent []string

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), nginxPortVariable+"=") {
			nginxPortString := scanner.Text()
			oldNginxPort := strings.Trim(
				strings.ReplaceAll(scanner.Text(), nginxPortVariable+"=", ""), " ",
			)
			nginxPortString = strings.Replace(nginxPortString, oldNginxPort, strconv.Itoa(project.DefaultConfig.NginxPort), -1)
			newEnvContent = append(newEnvContent, nginxPortString)
		} else if strings.HasPrefix(scanner.Text(), nginxSslPortVariable+"=") {
			nginxSslPortString := scanner.Text()
			oldSslNginxPort := strings.Trim(
				strings.ReplaceAll(scanner.Text(), nginxSslPortVariable+"=", ""), " ",
			)
			nginxSslPortString = strings.Replace(nginxSslPortString, oldSslNginxPort, strconv.Itoa(project.DefaultConfig.NginxSslPort), -1)
			newEnvContent = append(newEnvContent, nginxSslPortString)
		} else {
			newEnvContent = append(newEnvContent, scanner.Text())
		}
	}

	newEnvFile, _ := os.Create(project.ConfigPath)
	newEnvFile.Truncate(0)
	newEnvFile.WriteString(strings.Join(newEnvContent, "\n"))

	args := []string{"stop"}
	fmt.Println(
		helpers.SuccessText("("+project.ProjectName+"):"),
		"docker-compose", strings.Join(args, " "),
	)
	composeCmd := exec.Command("docker-compose", args...)
	composeCmd.Dir = project.ProjectPath
	composeCmd.Output()
}

// UnpackEnvironment : unpack docker-environment files
func UnpackEnvironment(path string) {
	conf := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded, rice.LocateAppended, rice.LocateFS},
	}

	box := conf.MustFindBox("./docker-environment")
	dockerComposeString, _ := box.String("docker-compose.yml")
	dockerComposeFile, _ := os.Create(path + "/docker-compose.yml")
	defer dockerComposeFile.Close()
	dockerComposeFile.WriteString(dockerComposeString)

	nginxConfString, _ := box.String("nginx.conf")
	nginxConfFile, _ := os.Create(path + "/nginx.conf")
	nginxConfFile.WriteString(nginxConfString)
	defer nginxConfFile.Close()
}

// UnpackProxyFiles : unpack nginx proxy files
func UnpackProxyFiles(path string, projects []ProjectConfig) {
	conf := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded, rice.LocateAppended, rice.LocateFS},
	}
	box := conf.MustFindBox("./docker-environment")
	upstreamString, _ := box.String("upstreams.conf")
	var upstreamContent []string

	proxyString, _ := box.String("proxy.conf")
	var proxyContent []string

	for _, project := range projects {
		port, _ := GetActiveProjectPort(project)
		if IsActiveProject(project) {
			upstreamContent = append(
				upstreamContent,
				fmt.Sprintf(
					upstreamString,
					project.ProjectName,
					port,
				),
			)
			proxyContent = append(
				proxyContent,
				fmt.Sprintf(
					proxyString,
					project.DefaultConfig.NginxDomain,
					project.ProjectName,
				),
			)
		}
	}

	upstreamConfFile, _ := os.Create(path + "/upstreams.conf")
	upstreamConfFile.WriteString(strings.Join(upstreamContent, "\n"))
	defer upstreamConfFile.Close()

	proxyConfFile, _ := os.Create(path + "/proxy.conf")
	proxyConfFile.WriteString(strings.Join(proxyContent, "\n"))
	defer proxyConfFile.Close()
}

// ReloadEnvironment : reload main environment
func ReloadEnvironment() bool {
	config, _ := LoadMainConfig()
	UnpackProxyFiles(GetConfigPath(), config.Projects)

	checkEnvCmd := exec.Command("docker-compose", "ps", "-q")
	checkEnvCmd.Dir = GetConfigPath() + "/"
	checkEnvCmdOut, _ := checkEnvCmd.Output()
	var envRunning bool
	if len(checkEnvCmdOut) == 0 {
		envRunning = false
	} else {
		envRunning = true
	}
	status := false
	command := "down"
	args := []string{}
	for _, project := range config.Projects {
		if IsActiveProject(project) {
			if envRunning {
				command = "restart"
			} else {
				command = "up"
				args = []string{"--force-recreate", "--build", "-d"}
			}
			status = true
			break
		}
	}

	args = append(args, command)
	copy(args[1:], args)
	args[0] = command
	fmt.Println(
		helpers.SuccessText("(reverse-proxy):"),
		"docker-compose", strings.Join(args, " "),
	)
	mainEnvCmd := exec.Command("docker-compose", args...)
	mainEnvCmd.Dir = GetConfigPath() + "/"
	mainEnvCmd.Output()

	return status
}
