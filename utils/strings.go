package utils

func StringStartsWith(str string, comp string) bool {
	return len(str) >= len(comp) && str[:len(comp)] == comp
}

func StringEndsWith(str string, comp string) bool {
	return len(str) >= len(comp) && str[len(str)-len(comp):] == comp
}
