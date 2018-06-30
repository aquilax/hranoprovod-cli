package parser

import (
	"strings"
)

func mytrim(s string) string {
	return strings.Trim(s, "\t \n:")
}
