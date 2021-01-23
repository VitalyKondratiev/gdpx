package helpers

func SuccessText(text string) (string) {
	return "\033[0;32m" + text + "\033[0m"
}

func FailText(text string) (string) {
	return "\033[0;31m" + text + "\033[0m"
}