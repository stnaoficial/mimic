package cli

import "fmt"

type Color string
type Weight string

const (
	Normal Weight = "0"
	Bold   Weight = "1"
)

const (
	Black  Color = "30"
	Red    Color = "31"
	Green  Color = "32"
	Yellow Color = "33"
	Blue   Color = "34"
	Purple Color = "35"
	Cyan   Color = "36"
	White  Color = "37"
	Gray   Color = "90"
)

func Sprint(weight Weight, color Color, text string) string {
	return fmt.Sprintf("\033[%s;%sm%s\033[0m", weight, color, text)
}

func Sprintf(weight Weight, color Color, text string, args ...any) string {
	return Sprint(weight, color, fmt.Sprintf(text, args...))
}

func Sprintln(weight Weight, color Color, text string) string {
	return Sprint(weight, color, fmt.Sprintln(text))
}

func Print(weight Weight, color Color, text string) {
	fmt.Print(Sprint(weight, color, text))
}

func Printf(weight Weight, color Color, text string, args ...any) {
	fmt.Print(Sprintf(weight, color, text, args...))
}

func Println(weight Weight, color Color, text string) {
	fmt.Print(Sprintln(weight, color, text))
}
