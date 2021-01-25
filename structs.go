package main

// EnvironmentConfig : configuration for docker environment
type EnvironmentConfig struct {
	NginxDomain string `json:"nginxDomain"`
	NginxPort   int    `json:"nginxPort"`
}

// ProjectConfig : configuration for project
type ProjectConfig struct {
	ProjectName   string            `json:"projectName"`
	ProjectPath   string            `json:"projectPath"`
	ConfigPath    string            `json:"configPath"`
	DefaultConfig EnvironmentConfig `json:"defaultConfig"`
	isActive      bool
}

// GlobalConfig : Main application configuration
type GlobalConfig struct {
	WorkdirPath string          `json:"workdirPath"`
	Projects    []ProjectConfig `json:"projects"`
}
