package main

type EnvironmentConfig struct {
	NginxDomain string `json:"nginxDomain"`
	NginxPort   int `json:"nginxPort"`
}

type ProjectConfig struct {
	ProjectPath   string `json:"projectPath"`
	ConfigPath    string `json:"configPath"`
	DefaultConfig EnvironmentConfig `json:"defaultConfig"`
}

type GlobalConfig struct {
	WorkdirPath string `json:"workdirPath"`
}
