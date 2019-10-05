package tokenizer

import "zen/source"

type SourceLocation struct {
	Source *source.Source
	At     int
}

func (at SourceLocation) WithOffset(offset int) SourceLocation {
	result := at.At + offset

	if result < 0 {
		result = 0
	} else if result > at.Source.Length() {
		result = at.Source.Length()
	}

	return SourceLocation{at.Source, result}
}
