package main

import (
	"strconv"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"./helpers"
	"encoding/json"
)

const userConfigPath = "/.config/gdpx"

func LoadMainConfig() (GlobalConfig, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	configPath := filepath.Join(dir, userConfigPath)
  var config GlobalConfig;
  configFile, err := ioutil.ReadFile(configPath + "/gdpx.json")
  if err != nil {
		config, _ := CreateMainConfig()
		return config, err
  }
	_ = json.Unmarshal([]byte(configFile), &config)
  return config, nil
}

func SaveMainConfig(config GlobalConfig)  {
	usr, _ := user.Current()
	dir := usr.HomeDir
	configPath := filepath.Join(dir, userConfigPath)
	_, err := os.Open(configPath + "/gdpx.json")
	if err != nil {
		os.MkdirAll(configPath, os.ModePerm)
	}
	file, _ := json.MarshalIndent(config, "", "\t")
	_ = ioutil.WriteFile(configPath + "/gdpx.json", file, 0644)
}

func CreateMainConfig() (GlobalConfig, error) {
	var config GlobalConfig;	
	usr, _ := user.Current()
	dir := usr.HomeDir
	config.WorkdirPath = helpers.SelectDirectory(dir)
	projects, _ := LoadDockerComposeConfigs(config.WorkdirPath)
	if (len(projects) > 0 ) {
		config.Projects = projects
		SaveMainConfig(config)
		fmt.Println(
			helpers.SuccessText("Succesfully added " + strconv.Itoa(len(projects)) + " projects\n"),
		)
	} else {
		fmt.Println(
			helpers.FailText("Location '" + config.WorkdirPath + "' don't contain any supported project"),
		)
	}
	return config, nil
}