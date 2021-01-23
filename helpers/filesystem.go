package helpers

import (
	"fmt"
	"io/ioutil"
	"strings"
	"path/filepath"
	"github.com/manifoldco/promptui"
)

func SelectDirectory(path string) string {
	var selectedDirectory string
	var names []string

	if (path != "/") {
		names = append(names, "..")
	}
	names = append(names, ".")

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if (f.IsDir() == true) {
			names = append(names, f.Name())
		}
	}

	prompt := promptui.Select{
		Label: "Current directory is '" + SuccessText(path) + "', select '" + SuccessText(".") + "' for select",
		Items: names,
		HideSelected: true,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return path
	}	
	path = strings.Trim(path, "/")
	if (path != "") {
		path = path + "/"
	}
	path = filepath.Clean("/" + path + result)
	if result != "." {
		return SelectDirectory(path)
	}

	selectedDirectory = path
	return selectedDirectory
}