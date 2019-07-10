package source

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Source struct {
	name    string
	content string
	lines   []string
}

func sourceFromContent(name string, content string) (result *Source) {
	return &Source{
		name,
		content,
		strings.Split(content, "\n"),
	}
}

func SourceFromString(content string) (result *Source) {
	return sourceFromContent("[anonymous]", content)
}

func SourceFromFile(filename string) (result *Source, err error) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return sourceFromContent(
		filename,
		string(data),
	), nil
}

func GetSourceContent(source *Source) (result string) {
	return source.content
}

func FormatLine(source *Source, at int) (message string) {
	var lineNumber = 0
	var colNumber = 0
	var lineStart = 0

	for index, line := range source.lines {
		var nextLineStart = lineStart + len(line) + 1

		if at < nextLineStart {
			lineNumber = index
			colNumber = at - lineStart
			break
		}

		lineStart = nextLineStart
	}

	return fmt.Sprintf("%s: (%d, %d)\n%s\n%s", source.name, lineNumber+1, colNumber+1, source.lines[lineNumber], strings.Repeat(" ", colNumber)+"^")
}
