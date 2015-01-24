package plugins

import "os"

var (
	IconDir   = os.Getenv("HOME")
	IconColor = "#78a4ff"
	Separator = "^fg(#2a2a2a) | ^fg()"
)

func icon(n string) string {
	return "^fg(" + IconColor + ")^i(" + IconDir + "/" + n + ")^fg() "
}
