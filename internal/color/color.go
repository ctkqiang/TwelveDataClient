package color

import "fmt"

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	dim    = "\033[2m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
	white  = "\033[97m"
)

func Bold(s string) string   { return bold + s + reset }
func Dim(s string) string    { return dim + s + reset }
func Red(s string) string    { return red + s + reset }
func Green(s string) string  { return green + s + reset }
func Yellow(s string) string { return yellow + s + reset }
func Blue(s string) string   { return blue + s + reset }
func Cyan(s string) string   { return cyan + s + reset }
func White(s string) string  { return white + s + reset }

func Boldf(format string, a ...interface{}) string   { return bold + fmt.Sprintf(format, a...) + reset }
func Dimf(format string, a ...interface{}) string    { return dim + fmt.Sprintf(format, a...) + reset }
func Redf(format string, a ...interface{}) string    { return red + fmt.Sprintf(format, a...) + reset }
func Greenf(format string, a ...interface{}) string  { return green + fmt.Sprintf(format, a...) + reset }
func Yellowf(format string, a ...interface{}) string { return yellow + fmt.Sprintf(format, a...) + reset }
func Bluef(format string, a ...interface{}) string   { return blue + fmt.Sprintf(format, a...) + reset }
func Cyanf(format string, a ...interface{}) string   { return cyan + fmt.Sprintf(format, a...) + reset }
func Whitef(format string, a ...interface{}) string  { return white + fmt.Sprintf(format, a...) + reset }
