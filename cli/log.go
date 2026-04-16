package cli

import (
	"fmt"
	"os"
	"strings"
)

func Log(cause string, severity LogSeverity) {
	switch severity {
	case LogSeverityError:
		fmt.Print(ANSIColorCodeRed)
		fmt.Printf("%s[ERROR] %s%s\n", ANSIColorCodeRed, cause, ANSIColorCodeReset)
	case LogSeverityWarn:
		fmt.Print(ANSIColorCodeYellow)
		fmt.Printf("%s[WARN] %s%s\n", ANSIColorCodeYellow, cause, ANSIColorCodeReset)
	case LogSeverityInfo:
		fmt.Printf("%s[INFO] %s%s\n", ANSIColorCodeYellow, cause, ANSIColorCodeReset)
	}

	fmt.Print(ANSIColorCodeReset)
}

func LogAndExit(cause string, severity LogSeverity) {
	Log(cause, severity)
	os.Exit(1)
}

func LogFileNameAt(name string) {
	fmt.Print(ANSIColorCodeCyan)
	fmt.Printf("@ %s\n", name)
	fmt.Print(ANSIColorCodeReset)
}

func LogFileNameAdded(name string) {
	fmt.Print(ANSIColorCodeGreen)
	fmt.Printf("+ %s\n", name)
	fmt.Print(ANSIColorCodeReset)
}

func LogFileDataAdded(data string) {
	fmt.Print(ANSIColorCodeGreen)

	// must be updated later for other OS compatibility
	for _, line := range strings.Split(data, "\n") {
		fmt.Printf("+ %s\n", line)
	}

	fmt.Printf("%s\n", ANSIColorCodeReset)
}
