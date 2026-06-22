package cli

import (
	"fmt"
	"time"
)

type LogSeverity int

const (
	LogSeverityInfo LogSeverity = iota
	LogSeveritySuccess
	LogSeverityWarn
	LogSeverityError
)

func prefix(severity LogSeverity) string {
	now := time.Now().UTC().Format(time.RFC3339)

	switch severity {
	case LogSeverityInfo:
		return fmt.Sprintf("%s [INFO]", now)
	case LogSeveritySuccess:
		return fmt.Sprintf("%s [SUCCESS]", now)
	case LogSeverityWarn:
		return fmt.Sprintf("%s [WARN]", now)
	case LogSeverityError:
		return fmt.Sprintf("%s [ERROR]", now)
	default:
		return fmt.Sprintf("%s [UNKNOWN]", now)
	}
}

func Log(severity LogSeverity, cause string) {
	switch severity {
	case LogSeverityInfo:
		Printf(Normal, Cyan, "%s %s", prefix(severity), cause)
	case LogSeveritySuccess:
		Printf(Normal, Green, "%s %s", prefix(severity), cause)
	case LogSeverityWarn:
		Printf(Normal, Yellow, "%s %s", prefix(severity), cause)
	case LogSeverityError:
		Printf(Normal, Red, "%s %s", prefix(severity), cause)
	}
}

func Logf(severity LogSeverity, cause string, args ...any) {
	switch severity {
	case LogSeverityInfo:
		Printf(Normal, Cyan, fmt.Sprintf("%s %s", prefix(severity), cause), args...)
	case LogSeveritySuccess:
		Printf(Normal, Green, fmt.Sprintf("%s %s", prefix(severity), cause), args...)
	case LogSeverityWarn:
		Printf(Normal, Yellow, fmt.Sprintf("%s %s", prefix(severity), cause), args...)
	case LogSeverityError:
		Printf(Normal, Red, fmt.Sprintf("%s %s", prefix(severity), cause), args...)
	}
}

func Logln(severity LogSeverity, cause string) {
	switch severity {
	case LogSeverityInfo:
		Println(Normal, Cyan, fmt.Sprintf("%s %s", prefix(severity), cause))
	case LogSeveritySuccess:
		Println(Normal, Green, fmt.Sprintf("%s %s", prefix(severity), cause))
	case LogSeverityWarn:
		Println(Normal, Yellow, fmt.Sprintf("%s %s", prefix(severity), cause))
	case LogSeverityError:
		Println(Normal, Red, fmt.Sprintf("%s %s", prefix(severity), cause))
	}
}
