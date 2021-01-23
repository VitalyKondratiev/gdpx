package main

type EnvironmentConfig struct {
	nginxDomain string
	nginxPort   int
}

type ProjectConfig struct {
	projectPath   string
	configPath    string
	defaultConfig EnvironmentConfig
}

type GlobalConfig struct {
	workdirPath string
	projects []ProjectConfig
}
