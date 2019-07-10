package parser

import (
	"testing"
	"zen/source"
	"zen/tokenizer"
)

func checkToken(t *testing.T, token *tokenizer.Token, value string, tokenType tokenizer.TokenType) {
	if token.TokenType != tokenType || token.Value != value {
		t.Errorf("Expected token type to be %d with value '%s' got %d with value '%s'", token.TokenType, token.Value, tokenType, value)
	}
}

func checkIdentifierAsStatement(t *testing.T, typeExpr Statement, value string) {
	var asIdentifier, ok = typeExpr.(*Identifier)

	if asIdentifier == nil || !ok {
		t.Errorf("Type is not an identifer")
	} else if asIdentifier.Token.Value != value {
		t.Errorf("Expected %s got %s", value, asIdentifier.Token.Value)
	}
}

func checkIdentifier(t *testing.T, typeExpr Expression, value string) {
	var asIdentifier, ok = typeExpr.(*Identifier)

	if asIdentifier == nil || !ok {
		t.Errorf("Type is not an identifer")
	} else if asIdentifier.Token.Value != value {
		t.Errorf("Expected %s got %s", value, asIdentifier.Token.Value)
	}
}

func checkTypeIdentifier(t *testing.T, typeExpr TypeExpression, value string) {
	var asNamedType, ok = typeExpr.(*NamedType)

	if asNamedType == nil || !ok {
		t.Errorf("Type is not a named type")
	} else if asNamedType.Token.Value != value {
		t.Errorf("Expected %s got %s", value, asNamedType.Token.Value)
	}
}

func TestIdentifier(t *testing.T) {
	var tokens = tokenizer.Tokenize(source.SourceFromString("a a0 _a 0"))
	var state = createState(&tokens)
	var result = createParseResult()

	_, okA := parseIdentifier(&result, &state)
	_, okB := parseIdentifier(&result, &state)
	_, okC := parseIdentifier(&result, &state)
	_, okD := parseIdentifier(&result, &state)
	if !okA {
		t.Errorf("Identifier A should not be nil")
	}

	if !okB {
		t.Errorf("Identifier B should not be nil")
	}

	if !okC {
		t.Errorf("Identifier C should not be nil")
	}

	if okD {
		t.Errorf("Identifier D should be nil")
	}

	if len(result.errors) != 1 {
		t.Errorf("Expected 1 parse error")
	}
}

func TestNumber(t *testing.T) {
	var tokens = tokenizer.Tokenize(source.SourceFromString("0 0.0"))
	var state = createState(&tokens)
	var result = createParseResult()

	var tokenA, _ = parseNumber(&result, &state)
	var tokenB, _ = parseNumber(&result, &state)

	if tokenA == nil {
		t.Errorf("Number A should not be nil")
	}

	if tokenB == nil {
		t.Errorf("Number B should not be nil")
	}
}

func TestStructure(t *testing.T) {
	var tokens = tokenizer.Tokenize(source.SourceFromString("[i32, u8, u8]"))
	var state = createState(&tokens)
	var result = createParseResult()

	var typeExp, ok = parseType(&result, &state)
	structureType, ok := typeExp.(*StructureType)

	if !ok || structureType == nil {
		t.Errorf("Returned type is not a structure")
	}

	if len(structureType.Entries) != 3 {
		t.Error("Wrong number of parameters")
	} else {
		checkTypeIdentifier(t, structureType.Entries[0].TypeExp, "i32")
		checkTypeIdentifier(t, structureType.Entries[1].TypeExp, "u8")
		checkTypeIdentifier(t, structureType.Entries[2].TypeExp, "u8")
	}

	tokens = tokenizer.Tokenize(source.SourceFromString("[a: i32, b: u8, c: u8]"))
	state = createState(&tokens)
	result = createParseResult()

	typeExp, ok = parseType(&result, &state)
	structureType, ok = typeExp.(*StructureType)

	if len(structureType.Entries) != 3 {
		t.Error("Wrong number of parameters")
	} else {
		checkToken(t, structureType.Entries[0].Name, "a", tokenizer.IDToken)
		checkTypeIdentifier(t, structureType.Entries[0].TypeExp, "i32")
		checkToken(t, structureType.Entries[1].Name, "b", tokenizer.IDToken)
		checkTypeIdentifier(t, structureType.Entries[1].TypeExp, "u8")
		checkToken(t, structureType.Entries[2].Name, "c", tokenizer.IDToken)
		checkTypeIdentifier(t, structureType.Entries[2].TypeExp, "u8")
	}

	tokens = tokenizer.Tokenize(source.SourceFromString("[i32, b: u8, c: u8]"))
	state = createState(&tokens)
	result = createParseResult()

	typeExp, ok = parseType(&result, &state)
	structureType, ok = typeExp.(*StructureType)

	if len(structureType.Entries) != 3 {
		t.Error("Wrong number of parameters")
	} else {
		checkTypeIdentifier(t, structureType.Entries[0].TypeExp, "i32")
		checkToken(t, structureType.Entries[1].Name, "b", tokenizer.IDToken)
		checkTypeIdentifier(t, structureType.Entries[1].TypeExp, "u8")
		checkToken(t, structureType.Entries[2].Name, "c", tokenizer.IDToken)
		checkTypeIdentifier(t, structureType.Entries[2].TypeExp, "u8")
	}
}

func TestFunctionType(t *testing.T) {
	var tokens = tokenizer.Tokenize(source.SourceFromString("[i32, u8, u8] => [u8]"))
	var state = createState(&tokens)
	var result = createParseResult()

	var typeExp, ok = parseType(&result, &state)
	functionType, ok := typeExp.(*FunctionType)

	if functionType == nil || !ok {
		t.Error("Did not parse as a function")
	}

	tokens = tokenizer.Tokenize(source.SourceFromString("[a:bool, b:bool] => [c:bool] => [d:u8]"))
	state = createState(&tokens)
	result = createParseResult()

	typeExp, ok = parseType(&result, &state)
	fnType, ok := typeExp.(*FunctionType)

	if fnType == nil {
		t.Errorf("Function type first")
		return
	}

	fnType, _ = fnType.Output.(*FunctionType)

	if fnType == nil {
		t.Error("Function type expected")
	}
}

func TestBody(t *testing.T) {
	var tokens = tokenizer.Tokenize(source.SourceFromString("{a; b c}"))
	var state = createState(&tokens)
	var result = createParseResult()

	var body = parseBody(&result, &state)

	if body == nil {
		t.Error("Expected body got nil")
		return
	}

	if len(body.Statements) != 3 {
		t.Errorf("Expected 3 statements got %d", len(body.Statements))
		return
	}

	checkIdentifierAsStatement(t, body.Statements[0], "a")
	checkIdentifierAsStatement(t, body.Statements[1], "b")
	checkIdentifierAsStatement(t, body.Statements[2], "c")
}

func TestWhere(t *testing.T) {
	var tokens = tokenizer.Tokenize(source.SourceFromString("[a:bool, b:bool] where a || b"))
	var state = createState(&tokens)
	var result = createParseResult()

	var typeExp, _ = parseType(&result, &state)
	whereExp, ok := typeExp.(*WhereType)

	if whereExp == nil || !ok {
		t.Errorf("Where expression did not parse")
		return
	}

	tokens = tokenizer.Tokenize(source.SourceFromString("[a:bool, b:bool] => [c:bool] where c"))
	state = createState(&tokens)
	result = createParseResult()

	typeExp, ok = parseType(&result, &state)
	whereExp, ok = typeExp.(*WhereType)

	if whereExp == nil {
		t.Errorf("Where expression did not parse")
		return
	}

	fnType, _ := whereExp.TypeExp.(*FunctionType)

	if fnType == nil {
		t.Error("Where should apply to the function type")
	}

	tokens = tokenizer.Tokenize(source.SourceFromString("[a:bool, b:bool] => ([c:bool] where c)"))
	state = createState(&tokens)
	result = createParseResult()

	typeExp, ok = parseType(&result, &state)
	fnType, ok = typeExp.(*FunctionType)

	if fnType == nil {
		t.Errorf("Function type first")
		return
	}

	whereType, _ := fnType.Output.(*WhereType)

	if whereType == nil {
		t.Error("Where expression expected")
	}
}

func TestFile(t *testing.T) {
	source, err := source.SourceFromFile("../../../test/FileTest.zen")

	if err != nil {
		t.Errorf("Source is error %s", err)
		return
	}
	var tokens = tokenizer.Tokenize(source)
	var state = createState(&tokens)
	var result = createParseResult()

	var fileDef = parseFileDefinition(&result, &state)

	if fileDef == nil {
		t.Error("File def is nil")
		return
	}

	if len(fileDef.Definitions) != 3 {
		t.Error("Not right number of definitions got ", len(fileDef.Definitions))
		return
	}
}

func TestOperators(t *testing.T) {
	var tokens = tokenizer.Tokenize(source.SourceFromString("a - b - c"))
	var state = createState(&tokens)
	var result = createParseResult()

	expType, _ := parseExpression(&result, &state)
	binaryExp, ok := expType.(*BinaryExpression)

	if binaryExp == nil && !ok {
		t.Error("binary expression is nil")
		return
	}

	checkIdentifier(t, binaryExp.Right, "c")
	checkToken(t, binaryExp.Operator, "-", tokenizer.MinusToken)

	binaryExp, _ = binaryExp.Left.(*BinaryExpression)

	if binaryExp == nil {
		t.Error("left binary expression is nil")
		return
	}

	checkIdentifier(t, binaryExp.Left, "a")
	checkIdentifier(t, binaryExp.Right, "b")
	checkToken(t, binaryExp.Operator, "-", tokenizer.MinusToken)

	tokens = tokenizer.Tokenize(source.SourceFromString("a - b * c"))
	state = createState(&tokens)
	result = createParseResult()

	expType, ok = parseExpression(&result, &state)
	binaryExp, ok = expType.(*BinaryExpression)

	if binaryExp == nil {
		t.Error("binary expression is nil")
		return
	}

	checkIdentifier(t, binaryExp.Left, "a")
	checkToken(t, binaryExp.Operator, "-", tokenizer.MinusToken)

	binaryExp, _ = binaryExp.Right.(*BinaryExpression)

	if binaryExp == nil {
		t.Error("right binary expression is nil")
		return
	}

	checkIdentifier(t, binaryExp.Left, "b")
	checkIdentifier(t, binaryExp.Right, "c")
	checkToken(t, binaryExp.Operator, "*", tokenizer.MultiplyToken)
}
