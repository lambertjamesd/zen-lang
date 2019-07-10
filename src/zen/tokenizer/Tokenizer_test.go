package tokenizer

import (
	"testing"
	"zen/source"
)

func checkToken(t *testing.T, token Token, value string, tokenType TokenType) {
	if token.TokenType != tokenType || token.Value != value {
		t.Errorf("Expected token type to be %d with value '%s' got %d with value '%s'", tokenType, value, token.TokenType, token.Value)
	}
}

func TestIdentifier(t *testing.T) {
	var source = source.SourceFromString("a b c")

	tokenizeResult := Tokenize(source)

	if len(tokenizeResult.Tokens) != 4 {
		t.Errorf("Expected token length to be 4 but was %d", len(tokenizeResult.Tokens))
	}

	checkToken(t, tokenizeResult.Tokens[0], "a", IDToken)
	checkToken(t, tokenizeResult.Tokens[1], "b", IDToken)
	checkToken(t, tokenizeResult.Tokens[2], "c", IDToken)
	checkToken(t, tokenizeResult.Tokens[3], "", EOFToken)
}

func TestSquareBracket(t *testing.T) {
	var source = source.SourceFromString("[a, b, c]")
	tokenizeResult := Tokenize(source)
	if len(tokenizeResult.Tokens) != 8 {
		t.Errorf("Expected token length to be 8 but was %d", len(tokenizeResult.Tokens))
	}
	checkToken(t, tokenizeResult.Tokens[0], "[", OpenSqaureToken)
	checkToken(t, tokenizeResult.Tokens[1], "a", IDToken)
	checkToken(t, tokenizeResult.Tokens[2], ",", CommaToken)
	checkToken(t, tokenizeResult.Tokens[3], "b", IDToken)
	checkToken(t, tokenizeResult.Tokens[4], ",", CommaToken)
	checkToken(t, tokenizeResult.Tokens[5], "c", IDToken)
	checkToken(t, tokenizeResult.Tokens[6], "]", CloseSquareToken)
}

func TestOperators(t *testing.T) {
	var source = source.SourceFromString("=>")
	tokenizeResult := Tokenize(source)
	if len(tokenizeResult.Tokens) != 2 {
		t.Errorf("Expected token length to be 1 but was %d", len(tokenizeResult.Tokens))
	}
	checkToken(t, tokenizeResult.Tokens[0], "=>", FatArrowToken)
}
