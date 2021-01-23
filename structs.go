package main

type EnvironmentConfig struct {
	NginxDomain string `json:"nginxDomain"`
	NginxPort   int `json:"nginxPort"`
}

type ProjectConfig struct {
	ProjectName   string `json:"projectName"`
	ProjectPath   string `json:"projectPath"`
	ConfigPath    string `json:"configPath"`
	DefaultConfig EnvironmentConfig `json:"defaultConfig"`
	isActive      bool
}

type GlobalConfig struct {
	WorkdirPath string `json:"workdirPath"`
	Projects []ProjectConfig `json:"projects"`
}
