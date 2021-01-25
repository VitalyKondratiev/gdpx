package helpers

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// SelectDirectory : select directory with console UI
func SelectDirectory(path string) string {
	var selectedDirectory string
	var names []string

	if path != "/" {
		names = append(names, "..")
	}
	names = append(names, ".")

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() == true {
			names = append(names, f.Name())
		}
	}

	result := SelectStringVariant("Current directory is '"+SuccessText(path)+"', select '"+SuccessText(".")+"' for select", names)

	path = strings.Trim(path, "/")
	if path != "" {
		path = path + "/"
	}
	path = filepath.Clean("/" + path + result)
	if result != "." {
		return SelectDirectory(path)
	}

	selectedDirectory = path
	return selectedDirectory
}
