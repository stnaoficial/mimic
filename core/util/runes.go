package util

func IsLetter(ch rune) bool {
	if ch >= 'a' && ch <= 'z' {
		return true
	}

	if ch >= 'A' && ch <= 'Z' {
		return true
	}

	if ch == '_' {
		return true
	}

	return false
}

func IsWhitespace(ch rune) bool {
	if ch == ' ' {
		return true
	}

	if ch == '\n' {
		return true
	}

	if ch == '\t' {
		return true
	}

	return false
}

func IsDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
