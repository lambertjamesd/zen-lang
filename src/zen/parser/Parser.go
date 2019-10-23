package parser

import (
	"zen/source"
	"zen/tokenizer"
)

type parseState struct {
	tokens   *tokenizer.TokenizeResult
	location uint
}

type parseResult struct {
	errors []ParseError
}

func createState(tokens *tokenizer.TokenizeResult) (result parseState) {
	return parseState{
		tokens,
		0,
	}
}

func createParseResult() (result parseResult) {
	return parseResult{
		nil,
	}
}

func advance(state *parseState) {
	state.location = state.location + 1
}

func peek(state *parseState, offset uint) (token *tokenizer.Token) {
	if state.location+offset >= uint(len(state.tokens.Tokens)) {
		return &state.tokens.Tokens[len(state.tokens.Tokens)-1]
	} else {
		return &state.tokens.Tokens[state.location+offset]
	}
}

func optional(state *parseState, tokenType tokenizer.TokenType) (result *tokenizer.Token) {
	var token = peek(state, 0)
	if token.TokenType == tokenType {
		advance(state)
		return token
	} else {
		return nil
	}
}

func optionalIdentifier(state *parseState, value string) (result *tokenizer.Token) {
	var token = peek(state, 0)
	if token.Value == value {
		advance(state)
		return token
	} else {
		return nil
	}
}

func expect(parseResult *parseResult, state *parseState, tokenType tokenizer.TokenType) (result *tokenizer.Token) {
	var maybeToken = optional(state, tokenType)
	if maybeToken == nil {
		parseResult.errors = append(parseResult.errors, CreateError(peek(state, 0).At, "Unexpected token '"+peek(state, 0).Value+"'"))
		advance(state)
	}

	return maybeToken
}

func expectIdentifier(parseResult *parseResult, state *parseState, value string) (result *tokenizer.Token) {
	var maybeToken = optional(state, tokenizer.IDToken)
	if maybeToken == nil || maybeToken.Value != value {
		parseResult.errors = append(parseResult.errors, CreateError(peek(state, 0).At, "Unexpected token '"+peek(state, 0).Value+"' expected '"+value+"'"))
		if maybeToken == nil {
			advance(state)
		}
	}

	return maybeToken
}

func parseIdentifier(parseResult *parseResult, state *parseState) (result *Identifier, okResult bool) {
	var maybeToken = expect(parseResult, state, tokenizer.IDToken)

	if maybeToken != nil {
		return &Identifier{
			maybeToken,
			&UndefinedType{},
		}, true
	} else {
		return nil, false
	}
}

func parseNumber(parseResult *parseResult, state *parseState) (result *Number, okResult bool) {
	var maybeToken = expect(parseResult, state, tokenizer.NumberToken)

	if maybeToken != nil {
		return &Number{
			maybeToken,
			&UndefinedType{},
		}, true
	} else {
		return nil, false
	}
}

func parseStructureExpression(parseResult *parseResult, state *parseState) (result *StructureExpression, okResult bool) {
	var openBracket = expect(parseResult, state, tokenizer.OpenSqaureToken)

	if openBracket == nil {
		return nil, false
	}

	var entries []StructureExpressionEntry = nil

	var hasNext = peek(state, 0).TokenType != tokenizer.CloseSquareToken

	for hasNext {
		var nameToken *tokenizer.Token = nil
		if peek(state, 0).TokenType == tokenizer.IDToken && peek(state, 1).TokenType == tokenizer.ColonToken {
			nameToken = expect(parseResult, state, tokenizer.IDToken)
			expect(parseResult, state, tokenizer.ColonToken)
		}

		expression, ok := parseExpression(parseResult, state)

		if !ok {
			return nil, false
		}

		entries = append(entries, StructureExpressionEntry{
			nameToken,
			expression,
		})

		hasNext = optional(state, tokenizer.CommaToken) != nil && peek(state, 0).TokenType != tokenizer.CloseSquareToken
	}

	var closeBracket = expect(parseResult, state, tokenizer.CloseSquareToken)

	return &StructureExpression{
		openBracket,
		entries,
		closeBracket,
		nil,
	}, closeBracket != nil
}

type typeOperatorPrecedence uint

const (
	minTypePrecedence  typeOperatorPrecedence = 0
	wherePrecedence    typeOperatorPrecedence = 1
	fatArrowPrecedence typeOperatorPrecedence = 2
	noTypePrecedence   typeOperatorPrecedence = 100
)

func parseSingleType(parseResult *parseResult, state *parseState) (result TypeExpression, okResult bool) {
	var next = peek(state, 0)

	if next.TokenType == tokenizer.OpenSqaureToken {
		result = parseStructureType(parseResult, state)
	} else if next.TokenType == tokenizer.IDToken {
		advance(state)
		result = &NamedType{
			next,
			&UndefinedType{},
		}
	} else if next.TokenType == tokenizer.OpenParenToken {
		advance(state)
		result, ok := parseType(parseResult, state)

		if !ok {
			return nil, false
		}

		var closeParen = expect(parseResult, state, tokenizer.CloseParenToken)
		if closeParen == nil {
			return nil, false
		}

		return result, true
	} else {
		parseResult.errors = append(parseResult.errors, CreateError(next.At, "Expcted type got '"+next.Value+"'"))
		return nil, false
	}

	return result, true
}

func parseUnaryType(parseResult *parseResult, state *parseState) (result TypeExpression, okResult bool) {
	return parseSingleType(parseResult, state)
}

func getTypeOperatorPrecedence(token *tokenizer.Token) (result typeOperatorPrecedence) {
	switch token.TokenType {
	case tokenizer.FatArrowToken:
		return fatArrowPrecedence
	case tokenizer.IDToken:
		switch token.Value {
		case "where":
			return wherePrecedence
		}
	}

	return noTypePrecedence
}

func parseBinaryType(parseResult *parseResult, state *parseState, precedence typeOperatorPrecedence) (result TypeExpression, okResult bool) {
	result, ok := parseUnaryType(parseResult, state)

	if !ok {
		return nil, false
	}

	var maybeOperator = peek(state, 0)
	var opPrecedence = getTypeOperatorPrecedence(maybeOperator)

	for opPrecedence != noTypePrecedence && opPrecedence >= precedence {
		advance(state)

		switch maybeOperator.Value {
		case "=>":
			right, ok := parseBinaryType(parseResult, state, opPrecedence)

			if !ok {
				return nil, false
			}

			result = &FunctionType{
				result,
				right,
				nil,
			}
		case "where":
			right, ok := parseExpression(parseResult, state)

			if !ok {
				return nil, false
			}

			result = &WhereType{
				maybeOperator,
				result,
				right,
			}
		}

		maybeOperator = peek(state, 0)
		opPrecedence = getTypeOperatorPrecedence(maybeOperator)
	}

	return result, true
}

func parseType(parseResult *parseResult, state *parseState) (result TypeExpression, okResult bool) {
	return parseBinaryType(parseResult, state, minTypePrecedence)
}

func checkHasNext(state *parseState, closingToken tokenizer.TokenType) (result bool) {
	if optional(state, tokenizer.CommaToken) == nil {
		return false
	} else {
		return peek(state, 0).TokenType != closingToken
	}
}

func parseStructureType(parseResult *parseResult, state *parseState) (result *StructureType) {
	openToken := expect(parseResult, state, tokenizer.OpenSqaureToken)
	if openToken == nil {
		return nil
	}

	var entries []*StructureNamedEntry = nil

	var hasNext = peek(state, 0).TokenType != tokenizer.CloseSquareToken

	for hasNext {
		if peek(state, 0).TokenType == tokenizer.IDToken && peek(state, 1).TokenType == tokenizer.ColonToken {
			hasNext = false
		} else {
			typeExp, ok := parseType(parseResult, state)

			if !ok {
				return nil
			}

			entries = append(entries, &StructureNamedEntry{nil, typeExp, nil})

			hasNext = checkHasNext(state, tokenizer.CloseSquareToken)
		}
	}

	hasNext = peek(state, 0).TokenType != tokenizer.CloseSquareToken

	for hasNext {
		var name = expect(parseResult, state, tokenizer.IDToken)

		if name == nil {
			return nil
		}

		if expect(parseResult, state, tokenizer.ColonToken) == nil {
			return nil
		}

		typeExp, ok := parseType(parseResult, state)

		if !ok {
			return nil
		}

		entries = append(entries, &StructureNamedEntry{name, typeExp, nil})

		hasNext = checkHasNext(state, tokenizer.CloseSquareToken)
	}

	closingToken := expect(parseResult, state, tokenizer.CloseSquareToken)

	if closingToken == nil {
		return nil
	}

	return &StructureType{
		openToken,
		entries,
		nil,
		closingToken,
	}
}

func parseTypeDefinition(parseResult *parseResult, state *parseState) (result *TypeDefinition) {
	openToken := expectIdentifier(parseResult, state, "type")
	if openToken == nil {
		return nil
	}
	var name = expect(parseResult, state, tokenizer.IDToken)
	if name == nil {
		return nil
	}
	typeExp, ok := parseType(parseResult, state)

	if !ok {
		return nil
	}

	return &TypeDefinition{
		openToken,
		name,
		typeExp,
		CreateScope(),
		&UndefinedType{},
	}
}

type expressionOperatorPrecedence uint

const (
	minExpressionPrecedence expressionOperatorPrecedence = 0
	boolOrPrecedence        expressionOperatorPrecedence = 1
	boolAndPrecedence       expressionOperatorPrecedence = 2
	equalityPrecedence      expressionOperatorPrecedence = 3
	comparePrecedence       expressionOperatorPrecedence = 4
	addExpPrecedence        expressionOperatorPrecedence = 5
	multiplyExpPrecedence   expressionOperatorPrecedence = 6
	noExpressionPrecedence  expressionOperatorPrecedence = 100
)

func getExpressionOperatorPrecedence(token *tokenizer.Token) (result expressionOperatorPrecedence) {
	switch token.TokenType {
	case tokenizer.EqualToken:
		fallthrough
	case tokenizer.NotEqualToken:
		return equalityPrecedence
	case tokenizer.GTEqToken:
		fallthrough
	case tokenizer.GTToken:
		fallthrough
	case tokenizer.LTEqToken:
		fallthrough
	case tokenizer.LTToken:
		return comparePrecedence
	case tokenizer.AddToken:
		fallthrough
	case tokenizer.MinusToken:
		return addExpPrecedence
	case tokenizer.MultiplyToken:
		fallthrough
	case tokenizer.DivideToken:
		return multiplyExpPrecedence
	case tokenizer.BooleanOrToken:
		return boolOrPrecedence
	case tokenizer.BooleanAndToken:
		return boolAndPrecedence
	}

	return noExpressionPrecedence
}

func isPostfixOperator(token *tokenizer.Token) bool {
	return token.TokenType == tokenizer.DotToken
}

func parsePostfixExpression(parseResult *parseResult, state *parseState) (result Expression, okResult bool) {
	result, ok := parseSingleExpression(parseResult, state)

	if !ok {
		return nil, false
	}

	var next = peek(state, 0)

	for isPostfixOperator(next) {
		if next.TokenType == tokenizer.DotToken {
			advance(state)
			var propertyName = expect(parseResult, state, tokenizer.IDToken)
			result = &PropertyExpression{
				result,
				propertyName,
				&UndefinedType{},
			}
		} else {
			parseResult.errors = append(parseResult.errors, CreateError(next.At, "Unkown postfix operator '"+next.Value+"'"))
			return result, false
		}
		next = peek(state, 0)
	}

	return result, true
}

func parseSingleExpression(parseResult *parseResult, state *parseState) (result Expression, okResult bool) {
	var next = peek(state, 0)

	if next.Value == "if" {
		return parseIfStatement(parseResult, state)
	} else if next.TokenType == tokenizer.IDToken {
		return parseIdentifier(parseResult, state)
	} else if next.TokenType == tokenizer.NumberToken {
		return parseNumber(parseResult, state)
	} else if next.TokenType == tokenizer.OpenParenToken {
		advance(state)
		result, ok := parseExpression(parseResult, state)
		ok = ok && expect(parseResult, state, tokenizer.CloseParenToken) != nil
		return result, ok
	} else if next.TokenType == tokenizer.OpenSqaureToken {
		return parseStructureExpression(parseResult, state)
	} else {
		advance(state)
		parseResult.errors = append(parseResult.errors, CreateError(next.At, "Expected expression got '"+next.Value+"'"))
		return nil, false
	}
}

func parseUnaryExpression(parseResult *parseResult, state *parseState) (result Expression, okResult bool) {
	var maybeOperator = peek(state, 0)

	if maybeOperator.TokenType == tokenizer.MinusToken {
		advance(state)
		expr, ok := parseUnaryExpression(parseResult, state)

		if !ok {
			return nil, false
		} else {
			return &UnaryExpression{
				expr,
				maybeOperator,
				&UndefinedType{},
			}, true
		}
	} else {
		return parsePostfixExpression(parseResult, state)
	}
}

func parseBinaryExpression(parseResult *parseResult, state *parseState, precedence expressionOperatorPrecedence) (result Expression, okResult bool) {
	result, ok := parseUnaryExpression(parseResult, state)

	if !ok {
		return nil, false
	}

	var maybeOperator = peek(state, 0)
	var opPrecedence = getExpressionOperatorPrecedence(maybeOperator)

	for opPrecedence != noExpressionPrecedence && opPrecedence > precedence {
		advance(state)

		right, ok := parseBinaryExpression(parseResult, state, opPrecedence)

		if !ok {
			return nil, false
		}

		result = &BinaryExpression{
			result,
			maybeOperator,
			right,
			&UndefinedType{},
		}

		maybeOperator = peek(state, 0)
		opPrecedence = getExpressionOperatorPrecedence(maybeOperator)
	}

	return result, true
}

func parseExpression(parseResult *parseResult, state *parseState) (result Expression, okResult bool) {
	return parseBinaryExpression(parseResult, state, minExpressionPrecedence)
}

func parseStatement(parseResult *parseResult, state *parseState) (result Statement, okResult bool) {
	var next = peek(state, 0)

	if next.Value == "return" {
		advance(state)

		var after = peek(state, 0)

		if after.TokenType == tokenizer.CloseCurlyToken || after.TokenType == tokenizer.SemicolonToken {
			return &ReturnStatement{
				next,
				nil,
				nil,
			}, true
		}

		var expressions []Expression = nil
		var hasNext = true

		for hasNext {
			var returnValue, ok = parseExpression(parseResult, state)
			expressions = append(expressions, returnValue)
			hasNext = ok && optional(state, tokenizer.CommaToken) != nil
		}

		return &ReturnStatement{
			next,
			expressions,
			nil,
		}, true
	}

	return parseExpression(parseResult, state)
}

func parseBody(parseResult *parseResult, state *parseState) (result *Body) {
	openToken := expect(parseResult, state, tokenizer.OpenCurlyToken)
	var statements []Statement = nil

	if openToken == nil {
		return nil
	}

	closeToken := optional(state, tokenizer.CloseCurlyToken)

	for closeToken == nil && peek(state, 0).TokenType != tokenizer.EOFToken {
		statement, ok := parseStatement(parseResult, state)
		if !ok {
			return nil
		}
		optional(state, tokenizer.SemicolonToken)
		statements = append(statements, statement)
		closeToken = optional(state, tokenizer.CloseCurlyToken)
	}

	return &Body{
		openToken,
		statements,
		&UndefinedType{},
		CreateScope(),
		closeToken,
	}
}

func parseFunction(parseResult *parseResult, state *parseState) (result *Function, okResult bool) {
	var typeExp, ok = parseType(parseResult, state)

	if !ok {
		return nil, false
	}

	var body = parseBody(parseResult, state)

	if body == nil {
		return nil, false
	}

	return &Function{
		typeExp,
		body,
		CreateScope(),
		nil,
	}, true
}

func parseFunctionDefinition(parseResult *parseResult, state *parseState) (result *FunctionDefinition, okResult bool) {
	if expectIdentifier(parseResult, state, "func") == nil {
		return nil, false
	}
	var name = expect(parseResult, state, tokenizer.IDToken)
	if name == nil {
		return nil, false
	}
	var function, ok = parseFunction(parseResult, state)

	if !ok {
		return nil, false
	}

	return &FunctionDefinition{
		name,
		function,
	}, true
}

func parseIfStatement(parseResult *parseResult, state *parseState) (result *IfStatement, okResult bool) {
	ifKeyword := expectIdentifier(parseResult, state, "if")

	if ifKeyword == nil {
		return nil, false
	}

	if expect(parseResult, state, tokenizer.OpenParenToken) == nil {
		return nil, false
	}

	var exp, ok = parseExpression(parseResult, state)

	if !ok {
		return nil, false
	}

	if expect(parseResult, state, tokenizer.CloseParenToken) == nil {
		return nil, false
	}

	var body = parseBody(parseResult, state)

	if body == nil {
		return nil, false
	}

	if optionalIdentifier(state, "else") != nil {
		var next = peek(state, 0)

		if next.Value == "if" {
			elseIf, ok := parseIfStatement(parseResult, state)

			if !ok {
				return nil, false
			}

			return &IfStatement{
				ifKeyword,
				exp,
				body,
				elseIf,
				&UndefinedType{},
			}, true
		} else {
			var elseBody = parseBody(parseResult, state)

			if elseBody == nil {
				return nil, false
			}

			return &IfStatement{
				ifKeyword,
				exp,
				body,
				elseBody,
				&UndefinedType{},
			}, true
		}
	}

	return &IfStatement{
		ifKeyword,
		exp,
		body,
		&Body{
			body.end,
			nil,
			&UndefinedType{},
			CreateScope(),
			body.end,
		},
		&UndefinedType{},
	}, true
}

func parseFileDefinition(parseResult *parseResult, state *parseState) (result *FileDefinition) {
	fileStart := peek(state, 0).At

	var definitions []Definition = nil

	var inError = false

	for peek(state, 0).TokenType != tokenizer.EOFToken {
		var next = peek(state, 0)

		if next.Value == "func" {
			var funcDef, ok = parseFunctionDefinition(parseResult, state)

			if !ok {
				inError = true
			} else {
				inError = false
				definitions = append(definitions, funcDef)
			}
		} else if next.Value == "type" {
			var typeDef = parseTypeDefinition(parseResult, state)

			if typeDef == nil {
				inError = true
			} else {
				inError = false
				definitions = append(definitions, typeDef)
			}
		} else {
			if !inError {
				parseResult.errors = append(parseResult.errors, CreateError(next.At, "Unexpected token '"+next.Value+"'"))
			}
			advance(state)

			inError = true
		}
	}

	return &FileDefinition{
		fileStart,
		definitions,
		CreateScope(),
		peek(state, 0).At,
	}
}

func Parse(src *source.Source) (result *FileDefinition, errors []ParseError) {
	var tokens = tokenizer.Tokenize(src)
	var state = createState(&tokens)
	var parseResult = createParseResult()

	return parseFileDefinition(&parseResult, &state), parseResult.errors
}

func ParseExpression(src *source.Source) (result Expression, errors []ParseError) {
	var tokens = tokenizer.Tokenize(src)
	var state = createState(&tokens)
	var parseResult = createParseResult()
	result, _ = parseExpression(&parseResult, &state)
	return result, parseResult.errors
}

func ParseTest(sourceString string) (result Expression, okResult bool) {
	var tokens = tokenizer.Tokenize(source.SourceFromString(sourceString))
	var state = createState(&tokens)
	var parseResult = createParseResult()

	return parseExpression(&parseResult, &state)
}
