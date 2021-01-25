package helpers

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// SelectStringVariant : select string from array with UI
func SelectStringVariant(header string, variants []string) string {

	prompt := promptui.Select{
		Label:        header,
		Items:        variants,
		HideSelected: true,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return result
	}

	return result
}

func SuccessText(text string) string {
	return "\033[0;32m" + text + "\033[0m"
}

func FailText(text string) string {
	return "\033[0;31m" + text + "\033[0m"
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
