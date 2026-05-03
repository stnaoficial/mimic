package util

func IsLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func IsDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func IsWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\n' || ch == '\t'
}
