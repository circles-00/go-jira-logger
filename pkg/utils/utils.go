package utils

import "fmt"

func ConstructOsc8Hyperlink(url, displayText string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, displayText)
}
