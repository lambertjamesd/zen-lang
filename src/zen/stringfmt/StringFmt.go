package stringfmt

import "strings"

func Indent(indent string, input string) string {
	return indent + strings.ReplaceAll(input, "\n", indent+"\n")
}
