package main

import "io/ioutil"
import "bufio"
import "os"
import "fmt"
import "strings"

func main() {
	files, err := GetDockerComposeConfigs("/home/user/Documents/GitProjects");
	if err != nil {
		fmt.Println(err)
	}
	for _, filePath := range files {
		fmt.Println("\033[0;31m" + filePath + "\033[0m")
		ReadEnvFile(filePath)
	}
}

func GetDockerComposeConfigs(root string) ([]string, error) {
	var files []string
	rootDirInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}
	for _, dir := range rootDirInfo {
		if (dir.IsDir()) {
			dirInfo, err := ioutil.ReadDir(root + "/" + dir.Name())
			if err != nil {
				return files, err
			}
			hasComposeFile := false
			hasEnvironmentFile := false
			for _, file := range dirInfo {
				if (file.Name() == "docker-compose.yml") {
					hasComposeFile = true
				}
				if (file.Name() == ".env") {
					hasEnvironmentFile = true
				}
			}
			if (hasComposeFile && hasEnvironmentFile) {
				files = append(files, root + "/" + dir.Name() + "/.env")
			}
		}
	}
	return files, nil
}

func ReadEnvFile(filePath string) ([]string, error) {
	var lines []string
	file, err := os.Open(filePath)
    if err != nil {
        return lines, err
    }
    defer file.Close()

	scanner := bufio.NewScanner(file)
	var siteDomain string
	var nginxPort string
    for scanner.Scan() {
		if (strings.HasPrefix(scanner.Text(), "SITE_DOMAIN=")) {
			siteDomain = strings.Trim(strings.ReplaceAll(scanner.Text(), "SITE_DOMAIN=", ""), " ")
		}
		if (strings.HasPrefix(scanner.Text(), "NGINX_PORT=")) {
			nginxPort = strings.Trim(strings.ReplaceAll(scanner.Text(), "NGINX_PORT=", ""), " ")
		}
	}
	fmt.Println(siteDomain + ":" + nginxPort)
	return lines, nil
}