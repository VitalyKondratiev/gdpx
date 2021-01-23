package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"./helpers"
	"encoding/json"
)

const userConfigDir = "/.config/gdpx"

func main() {
	config, _ := LoadMainConfig()
	projects, err := LoadDockerComposeConfigs(config.WorkdirPath)
	if err != nil {
		fmt.Println(err)
	}
	if (len(projects) > 0 ){
		for _, project := range projects {
			fmt.Println(
				helpers.SuccessText(project.DefaultConfig.NginxDomain + ":" + strconv.Itoa(project.DefaultConfig.NginxPort)),
			)
		}
		SaveMainConfig(config)
	} else {
		fmt.Println(
			helpers.FailText("Location '" + config.WorkdirPath + "' don't contain any supported project"),
		)
	}
}

func LoadMainConfig() (GlobalConfig, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	userConfigDir := filepath.Join(dir, userConfigDir)
  var globalConfig GlobalConfig;
  configFile, err := ioutil.ReadFile(userConfigDir + "/gdpx.json")
  if err != nil {
		globalConfig.WorkdirPath = helpers.SelectDirectory("/")
		return globalConfig, err
  }
	_ = json.Unmarshal([]byte(configFile), &globalConfig)
  return globalConfig, nil
}

func SaveMainConfig(config GlobalConfig)  {
	usr, _ := user.Current()
	dir := usr.HomeDir
	userConfigDir := filepath.Join(dir, userConfigDir)
	_, err := os.Open(userConfigDir + "/gdpx.json")
	if err != nil {
		os.MkdirAll(userConfigDir, os.ModePerm)
	}
	file, _ := json.MarshalIndent(config, "", "\t")
	_ = ioutil.WriteFile(userConfigDir + "/gdpx.json", file, 0644)
}

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
				projectPath := root+"/"+dir.Name()+"/.env"
				configPath := root+"/"+dir.Name()
				config, _ := GetProjectEnvironment(projectPath)
				project := ProjectConfig{projectPath, configPath, config}
				projects = append(projects, project)
			}
		}
	}
	return projects, nil
}

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
		if strings.HasPrefix(scanner.Text(), "SITE_DOMAIN=") {
			siteDomain = strings.Trim(
				strings.ReplaceAll(scanner.Text(), "SITE_DOMAIN=", ""), " ",
			)
		}
		if strings.HasPrefix(scanner.Text(), "NGINX_PORT=") {
			nginxPort = strings.Trim(
				strings.ReplaceAll(scanner.Text(), "NGINX_PORT=", ""), " ",
			)
		}
	}
	nginxPortInt, err := strconv.Atoi(nginxPort)
	config = EnvironmentConfig{siteDomain, nginxPortInt}
	return config, nil
}
