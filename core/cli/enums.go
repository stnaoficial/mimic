package cli

type LogSeverity int

const (
	LogSeverityInfo  LogSeverity = 0
	LogSeverityWarn  LogSeverity = 1
	LogSeverityError LogSeverity = 2
)

type ANSIColorCode string

const (
	ANSIColorCodeReset  ANSIColorCode = "\033[0;0m"
	ANSIColorCodeBlack  ANSIColorCode = "\033[0;30m"
	ANSIColorCodeRed    ANSIColorCode = "\033[0;31m"
	ANSIColorCodeGreen  ANSIColorCode = "\033[0;32m"
	ANSIColorCodeYellow ANSIColorCode = "\033[0;33m"
	ANSIColorCodeBlue   ANSIColorCode = "\033[0;34m"
	ANSIColorCodePurple ANSIColorCode = "\033[0;35m"
	ANSIColorCodeCyan   ANSIColorCode = "\033[0;36m"
	ANSIColorCodeWhite  ANSIColorCode = "\033[0;37m"
)
