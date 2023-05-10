package utils

import (
	"fmt"
)

func ToLink(text, address string) string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>", address, text)
}
