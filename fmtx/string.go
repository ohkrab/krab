package fmtx

import (
	"regexp"
	"strings"
)

var (
	ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

func StripANSI(str string) string {
	return ansiRegex.ReplaceAllString(str, "")
}

func Squish(str string) string {
	return strings.Join(strings.Fields(str), " ")
}
