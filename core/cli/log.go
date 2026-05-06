package cli

import (
	"fmt"
	"os"
)

type LogSeverity int

const (
	LogSeverityInfo LogSeverity = iota
	LogSeveritySuccess
	LogSeverityWarn
	LogSeverityError
)

func Log(cause string, severity LogSeverity) {
	switch severity {
	case LogSeverityInfo:
		fmt.Printf("%s%s", ANSIColorCodeCyan, cause)
	case LogSeveritySuccess:
		fmt.Printf("%s%s", ANSIColorCodeGreen, cause)
	case LogSeverityWarn:
		fmt.Printf("%s%s", ANSIColorCodeYellow, cause)
	case LogSeverityError:
		fmt.Printf("%s%s", ANSIColorCodeRed, cause)
	}

	fmt.Printf("%s\n", ANSIColorCodeReset)
}

func LogWithPrefix(cause string, severity LogSeverity) {
	switch severity {
	case LogSeverityInfo:
		Log(fmt.Sprintf("[INFO] %s", cause), severity)
	case LogSeveritySuccess:
		Log(fmt.Sprintf("[SUCCESS] %s", cause), severity)
	case LogSeverityWarn:
		Log(fmt.Sprintf("[WARN] %s", cause), severity)
	case LogSeverityError:
		Log(fmt.Sprintf("[ERROR] %s", cause), severity)
	}
}

func LogAndExit(cause string, severity LogSeverity) {
	LogWithPrefix(cause, severity)
	os.Exit(1)
}
