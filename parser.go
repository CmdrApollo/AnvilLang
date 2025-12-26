package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type ParseNode struct {
	Name string
	Data string

	Parent   *ParseNode
	Children []*ParseNode
}

type Parser struct {
	Tokens     []Token
	Pos        int
	TreeNode   ParseNode
	BasicTypes []string
	ToPath     string
}

func (p *Parser) Peek() Token {
	if p.Pos < len(p.Tokens) {
		return p.Tokens[p.Pos]
	}
	return Token{}
}

func (p *Parser) PeekNext() Token {
	if p.Pos+1 < len(p.Tokens) {
		return p.Tokens[p.Pos+1]
	}
	return Token{}
}

func (p *Parser) Next() Token {
	if p.Pos+1 < len(p.Tokens) {
		p.Pos++
		return p.Tokens[p.Pos]
	}
	return Token{}
}

func (p *Parser) Expect(tType string) Token {
	t := p.Next()
	if t.TokenType != tType {
		fmt.Println(fmt.Sprintf("ERROR: Expected token of type '%s', got type '%s'.", tType, t.TokenType))
		return Token{}
	}
	return t
}

func (p *Parser) ExpectNoAdvance(tType string) Token {
	t := p.PeekNext()
	if t.TokenType != tType {
		fmt.Println(fmt.Sprintf("ERROR: Expected token of type '%s', got type '%s'.", tType, t.TokenType))
		return Token{}
	}
	return t
}

func (p *Parser) ParseAtom(currentNode *ParseNode) *ParseNode {
	currTok := p.Next()

	switch currTok.TokenType {
	case "left_square_token":
		arrNode := &ParseNode{Name: "array", Parent: currentNode, Children: make([]*ParseNode, 0)}
		if p.PeekNext().TokenType != "right_square_token" {
			for p.Pos < len(p.Tokens) {
				expr := p.ParseExpr(arrNode)
				arrNode.Children = append(arrNode.Children, expr)
				if p.PeekNext().TokenType == "right_square_token" {
					break
				}
				if p.Expect("comma_token").TokenType != "comma_token" {
					break
				}
			}
		}
		p.Expect("right_square_token")
		return arrNode
	case "left_curly_token":
		mapNode := &ParseNode{Name: "map", Parent: currentNode, Children: make([]*ParseNode, 0)}
		if p.PeekNext().TokenType != "right_curly_token" {
			for p.Pos < len(p.Tokens) {
				name := p.ParseExpr(mapNode)
				p.Expect("equal_token")
				expr := p.ParseExpr(mapNode)
				mapNode.Children = append(mapNode.Children, &ParseNode{Name: "map_element", Parent: mapNode, Children: []*ParseNode{name, expr}})
				if p.PeekNext().TokenType == "right_curly_token" {
					break
				}
				if p.Expect("comma_token").TokenType != "comma_token" {
					break
				}
			}
		}
		p.Expect("right_curly_token")
		return mapNode
	case "int_literal", "float_literal", "bool_literal", "string_literal", "null_literal":
		return &ParseNode{Name: currTok.TokenType, Data: currTok.TokenValue, Parent: currentNode}
	case "muldiv_operator", "addsub_operator", "comparison_operator", "equality_operator":
		return &ParseNode{Name: currTok.TokenType, Data: currTok.TokenValue, Parent: currentNode}
	case "and_operator", "or_operator", "not_operator":
		return &ParseNode{Name: currTok.TokenType, Data: currTok.TokenType, Parent: currentNode}
	case "new_keyword":
		structNode := &ParseNode{Name: "struct_instance", Parent: currentNode, Children: make([]*ParseNode, 0)}
		nameNode := p.ParseAtom(structNode)
		if nameNode.Name != "name" {
			fmt.Println("ERROR: Expected struct name after new keyword.")
		}

		structNode.Children = append(structNode.Children, nameNode)

		p.Expect("left_curly_token")

		if p.PeekNext().TokenType != "right_curly_token" {
			for p.Pos < len(p.Tokens) {
				expr := p.ParseExpr(structNode)
				structNode.Children = append(structNode.Children, expr)
				if p.PeekNext().TokenType == "right_curly_token" {
					break
				}
				if p.Expect("comma_token").TokenType != "comma_token" {
					break
				}
			}
		}

		p.Expect("right_curly_token")

		return structNode
	case "identifier":
		if p.PeekNext().TokenType == "left_paren_token" {
			// func call
			callNode := &ParseNode{Name: "func_call", Parent: currentNode}

			funcName := currTok.TokenValue

			p.Next()

			expressions := make([]*ParseNode, 0)
			nextToken := p.PeekNext()
			if nextToken.TokenType != "right_paren_token" {
				for p.Pos < len(p.Tokens) {
					expr := p.ParseExpr(callNode)
					expressions = append(expressions, expr)
					nextToken := p.PeekNext()
					if nextToken.TokenType == "right_paren_token" {
						break
					}
					if p.Expect("comma_token").TokenType != "comma_token" {
						break
					}
				}
			}
			p.Expect("right_paren_token")

			nameNode := ParseNode{Name: "name", Data: funcName, Parent: callNode}

			callNode.Children = []*ParseNode{&nameNode}

			for _, expr := range expressions {
				callNode.Children = append(callNode.Children, expr)
			}

			return callNode
		} else if p.PeekNext().TokenType == "period_token" {
			// accessfield
			fieldNode := ParseNode{Name: "access_field", Parent: currentNode}

			structName := currTok.TokenValue
			nameNode := ParseNode{Name: "name", Data: structName, Parent: &fieldNode}

			fieldNode.Children = []*ParseNode{&nameNode}

			for p.PeekNext().TokenType == "period_token" {
				p.Expect("period_token")
				nextNode := ParseNode{Name: "name", Data: p.Next().TokenValue, Parent: &fieldNode}
				fieldNode.Children = append(fieldNode.Children, &nextNode)
			}

			return &fieldNode
		} else if slices.Contains(p.BasicTypes, currTok.TokenValue) {
			typeNode := &ParseNode{Name: "type", Data: currTok.TokenValue}
			children := []*ParseNode{{Name: "bool_literal", Data: "false", Parent: typeNode}, {Name: "name", Data: currTok.TokenValue, Parent: typeNode}}

			if p.PeekNext().TokenType == "comparison_operator" && p.PeekNext().TokenValue == "lesser" {
				p.Next()

				children = append(children, p.ParseAtom(currentNode))
				
				if p.PeekNext().TokenType == "comma_token" {
					p.Next()
					children = append(children, p.ParseAtom(currentNode))
				}

				if !(p.PeekNext().TokenType == "comparison_operator" && p.PeekNext().TokenValue == "greater") {
					fmt.Println("ERROR: Missing '>' in type.")
				}
				p.Next()
			}
			if p.PeekNext().TokenType == "question_token" {
				p.Next()
				children[0] = &ParseNode{Name: "bool_literal", Data: "true", Parent: typeNode}
			}
			typeNode.Children = children
			return typeNode
		}
		return &ParseNode{Name: "name", Data: currTok.TokenValue, Parent: currentNode}
	case "left_paren_token":
		expr := p.ParseExpr(currentNode)
		p.Expect("right_paren_token")
		return expr
	}

	return &ParseNode{}
}

func (p *Parser) ParseGetItem(currentNode *ParseNode) *ParseNode {
	left := p.ParseAtom(currentNode)

	if p.PeekNext().TokenType == "left_square_token" {
		// getitem
		itemNode := ParseNode{Name: "get_item", Parent: currentNode}

		itemNode.Children = append(itemNode.Children, left)

		for p.Pos < len(p.Tokens) {
			p.Next()

			itemNode.Children = append(itemNode.Children, p.ParseExpr(&itemNode))
			p.Expect("right_square_token")
			if p.PeekNext().TokenType != "left_square_token" {
				break
			}
		}

		return &itemNode
	}

	return left
}

func (p *Parser) ParseTerm(currentNode *ParseNode) *ParseNode {
	left := p.ParseGetItem(currentNode)

	for p.PeekNext().TokenType == "muldiv_operator" {
		operator := p.ParseGetItem(currentNode)
		right := p.ParseGetItem(currentNode)

		children := []*ParseNode{left, operator, right}

		left = &ParseNode{Name: "muldiv_operation", Parent: currentNode, Children: children}
	}

	return left
}

func (p *Parser) ParseAddSub(currentNode *ParseNode) *ParseNode {
	left := p.ParseTerm(currentNode)

	if left.Name == "addsub_operator" && left.Data == "minus" {
		expr := p.ParseEquality(currentNode)
		children := []*ParseNode{expr}
		left = &ParseNode{Name: "unary_negate_operation", Parent: currentNode, Children: children}
		return left;
	}

	for p.PeekNext().TokenType == "addsub_operator" {
		operator := p.ParseTerm(currentNode)
		right := p.ParseTerm(currentNode)

		children := []*ParseNode{left, operator, right}

		left =& ParseNode{Name: "addsub_operation", Parent: currentNode, Children: children}
	}

	return left
}

func (p *Parser) ParseComp(currentNode *ParseNode) *ParseNode {
	left := p.ParseAddSub(currentNode)

	for p.PeekNext().TokenType == "comparison_operator" {
		operator := p.ParseAddSub(currentNode)
		right := p.ParseAddSub(currentNode)

		children := []*ParseNode{left, operator, right}

		left = &ParseNode{Name: "comparison_operation", Parent: currentNode, Children: children}
	}

	return left
}

func (p *Parser) ParseEquality(currentNode *ParseNode) *ParseNode {
	left := p.ParseComp(currentNode)

	for p.PeekNext().TokenType == "equality_operator" {
		operator := p.ParseComp(currentNode)
		right := p.ParseComp(currentNode)

		children := []*ParseNode{left, operator, right}

		left = &ParseNode{Name: "equality_operation", Parent: currentNode, Children: children}
	}

	return left
}

func (p *Parser) ParseNotFactor(currentNode *ParseNode) *ParseNode {
	left := p.ParseEquality(currentNode)

	if left.Name == "not_operator" {
		expr := p.ParseEquality(currentNode)
		children := []*ParseNode{expr}
		left = &ParseNode{Name: "not_operation", Parent: currentNode, Children: children}
	}

	return left
}

func (p *Parser) ParseAndFactor(currentNode *ParseNode) *ParseNode {
	left := p.ParseNotFactor(currentNode)

	for p.PeekNext().TokenType == "and_operator" {
		p.ParseNotFactor(currentNode)
		right := p.ParseNotFactor(currentNode)

		children := []*ParseNode{left, right}

		left = &ParseNode{Name: "and_operation", Parent: currentNode, Children: children}
	}

	return left
}

func (p *Parser) ParseExpr(currentNode *ParseNode) *ParseNode {
	left := p.ParseAndFactor(currentNode)

	for p.PeekNext().TokenType == "or_operator" {
		p.ParseAndFactor(currentNode)
		right := p.ParseAndFactor(currentNode)

		children := []*ParseNode{left, right}

		left = &ParseNode{Name: "or_operation", Parent: currentNode, Children: children}
	}

	return left
}

func (p *Parser) Parse(currentNode *ParseNode) *ParseNode {
	startPos := p.Pos
	tok := p.Next()
	
	if tok.TokenType == "right_curly_token" {
		// let the enclosing block consume it
		p.Pos--
		return &ParseNode{}
	}

	if tok.TokenType == "pass_keyword" {
		stmt := &ParseNode{Name: "pass_stmt", Parent: currentNode}
		p.Expect("semicolon_token")
		return stmt
	} else if tok.TokenType == "continue_keyword" {
		stmt := &ParseNode{Name: "continue_stmt", Parent: currentNode}
		p.Expect("semicolon_token")
		return stmt
	} else if tok.TokenType == "break_keyword" {
		stmt := &ParseNode{Name: "break_stmt", Parent: currentNode}
		p.Expect("semicolon_token")
		return stmt
	} else if tok.TokenType == "return_keyword" {
		stmt := &ParseNode{Name: "return_stmt", Parent: currentNode}
		stmt.Children = append(stmt.Children, p.ParseExpr(stmt))
		p.Expect("semicolon_token")
		return stmt
	} else if tok.TokenType == "import_keyword" {
		fileName := p.Expect("string_literal").TokenValue
		p.ExpectNoAdvance("semicolon_token")

		// FIXME imports are super duper ultra simple right now

		newPath := p.ToPath + "/" + fileName

		fileText, err := os.ReadFile(newPath)

		if err != nil {
			fmt.Println(fmt.Sprintf("ERROR: Failed to import file '%s'.", newPath))
			os.Exit(1)
		}

		newLexer := Lexer{Stream: string(fileText)}
		newLexer.Lex()

		newParser := Parser{Tokens: newLexer.Tokens, ToPath: newPath}
		newParser.ParseAll()

		for _, child := range newParser.TreeNode.Children {
			p.TreeNode.Children = append(p.TreeNode.Children, child)
		}

		return &ParseNode{Name: "pass_stmt", Parent: currentNode}
	} else if tok.TokenType == "fn_keyword" {
		declNode := &ParseNode{Name: "func_decl", Parent: currentNode}

		funcName := p.Expect("identifier").TokenValue
		p.Expect("left_paren_token")
		// arguments
		args := &ParseNode{Name: "args", Parent: declNode}

		if p.PeekNext().TokenType != "right_paren_token" {
			for p.Pos < len(p.Tokens) {
				arg := &ParseNode{Name: "argument", Parent: args}

				p.ExpectNoAdvance("identifier")
				arg.Children = append(arg.Children, p.ParseAtom(declNode))
				p.Expect("colon_token")
				p.ExpectNoAdvance("identifier")
				arg.Children = append(arg.Children, p.ParseAtom(declNode))

				if p.PeekNext().TokenType == "equal_token" {
					p.Next()
					arg.Children = append(arg.Children, p.ParseExpr(declNode))
				}
				args.Children = append(args.Children, arg)
				if p.PeekNext().TokenType == "right_paren_token" {
					break
				}
				if p.Expect("comma_token").TokenType != "comma_token" {
					break
				}
			}
		}

		p.Expect("right_paren_token")
		p.Expect("arrow_token")
		p.ExpectNoAdvance("identifier")
		returnType := p.ParseAtom(declNode)

		p.Expect("left_curly_token")
		// suite

		suite := &ParseNode{Name: "suite", Parent: declNode}

		for p.Pos < len(p.Tokens) {
			suite.Children = append(suite.Children, p.Parse(suite))
			if p.PeekNext().TokenType == "right_curly_token" {
				break
			}
		}
		// FIXME
		p.Expect("right_curly_token")

		nameNode := &ParseNode{Name: "name", Data: funcName, Parent: declNode}

		declNode.Children = []*ParseNode{nameNode, args, returnType, suite}

		return declNode
	} else if tok.TokenType == "st_keyword" {
		declNode := &ParseNode{Name: "struct_decl", Parent: currentNode}

		structName := p.Expect("identifier").TokenValue
		p.Expect("left_curly_token")
		// fields
		fields := &ParseNode{Name: "fields", Parent: declNode}

		if p.PeekNext().TokenType != "right_curly_token" {
			for p.Pos < len(p.Tokens) {
				arg := &ParseNode{Name: "field", Parent: fields}

				p.ExpectNoAdvance("identifier")
				arg.Children = append(arg.Children, p.ParseAtom(declNode))
				p.Expect("colon_token")
				p.ExpectNoAdvance("identifier")
				arg.Children = append(arg.Children, p.ParseAtom(declNode))

				if p.PeekNext().TokenType == "equal_token" {
					p.Next()
					arg.Children = append(arg.Children, p.ParseExpr(declNode))
				}
				fields.Children = append(fields.Children, arg)
				if p.PeekNext().TokenType == "right_curly_token" {
					break
				}
				if p.Expect("comma_token").TokenType != "comma_token" {
					break
				}
			}
		}

		p.Expect("right_curly_token")

		nameNode := &ParseNode{Name: "name", Data: structName, Parent: declNode}

		declNode.Children = []*ParseNode{nameNode, fields}

		return declNode
	} else if tok.TokenType == "while_keyword" {
		whileNode :=& ParseNode{Name: "while_stmt", Parent: currentNode}

		p.Expect("left_paren_token")
		// arguments
		condition := p.ParseExpr(whileNode)

		p.Expect("right_paren_token")

		p.Expect("left_curly_token")
		// suite

		suite := &ParseNode{Name: "suite", Parent: whileNode}

		for p.Pos < len(p.Tokens) {
			suite.Children = append(suite.Children, p.Parse(suite))
			if p.PeekNext().TokenType == "right_curly_token" {
				break
			}
		}
		p.Expect("right_curly_token")

		whileNode.Children = []*ParseNode{condition, suite}

		return whileNode
	} else if tok.TokenType == "if_keyword" {

		ifNode := &ParseNode{Name: "if_stmt", Parent: currentNode}

		p.Expect("left_paren_token")
		// arguments
		condition := p.ParseExpr(ifNode)

		p.Expect("right_paren_token")

		p.Expect("left_curly_token")
		// suite

		suite := &ParseNode{Name: "suite", Parent: ifNode}

		for p.Pos < len(p.Tokens) {
			suite.Children = append(suite.Children, p.Parse(suite))
			if p.PeekNext().TokenType == "right_curly_token" {
				break
			}
		}
		p.Expect("right_curly_token")

		elseSuite := &ParseNode{Name: "suite", Parent: ifNode}
		if p.PeekNext().TokenType == "else_keyword" {
			p.Next()
			if p.PeekNext().TokenType == "if_keyword" {
				elseSuite.Children = append(elseSuite.Children, p.Parse(suite))
			} else {
				p.Expect("left_curly_token")
				for p.Pos < len(p.Tokens) {
					elseSuite.Children = append(elseSuite.Children, p.Parse(suite))
					if p.PeekNext().TokenType == "right_curly_token" {
						break
					}
				}
				p.Expect("right_curly_token")
			}
		}

		ifNode.Children = []*ParseNode{condition, suite, elseSuite}

		return ifNode
	} else if tok.TokenType == "loop_keyword" {
		loopNode := &ParseNode{Name: "loop_stmt", Parent: currentNode}

		p.Expect("left_paren_token")
		values := p.ParseExpr(loopNode)
		p.Expect("right_paren_token")

		p.Expect("arrow_token")

		p.Expect("left_paren_token")
		// arguments
		varName := p.ParseExpr(loopNode)
		p.Expect("colon_token")
		varType := p.ParseAtom(loopNode)

		p.Expect("right_paren_token")

		p.Expect("left_curly_token")
		// suite

		suite := &ParseNode{Name: "suite", Parent: loopNode}

		for p.Pos < len(p.Tokens) {
			suite.Children = append(suite.Children, p.Parse(suite))
			if p.PeekNext().TokenType == "right_curly_token" {			
				break
			}
		}
		p.Expect("right_curly_token")

		loopNode.Children = []*ParseNode{values, varName, varType, suite}

		return loopNode
	} else if tok.TokenType == "var_keyword" {
		declNode := &ParseNode{Name: "variable_decl", Parent: currentNode}

		varName := p.Expect("identifier").TokenValue // expect identifier and get its value
		p.Expect("colon_token")                      // expect colon
		p.ExpectNoAdvance("identifier")
		varType := p.ParseAtom(declNode)
		p.Expect("equal_token") // expect equals
		expr := p.ParseExpr(declNode)
		p.ExpectNoAdvance("semicolon_token") // expect semicolon after the expression

		nameNode := &ParseNode{Name: "name", Data: varName, Parent: declNode}

		declNode.Children = []*ParseNode{nameNode, varType, expr}

		return declNode
	} else if tok.TokenType == "identifier" {
		// this can be either mut stmt or func call
		nextToken := p.PeekNext()

		if nextToken.TokenType == "left_paren_token" {
			p.Next()

			// func call
			callNode := &ParseNode{Name: "func_call", Parent: currentNode}

			funcName := tok.TokenValue // expect identifier and get its value
			expressions := make([]*ParseNode, 0)
			nextToken := p.PeekNext()
			if nextToken.TokenType != "right_paren_token" {
				for p.Pos < len(p.Tokens) {
					expr := p.ParseExpr(callNode)
					expressions = append(expressions, expr)
					nextToken := p.PeekNext()
					if nextToken.TokenType == "comma_token" {
						p.Next()
					} else if nextToken.TokenType == "right_paren_token" {
						break
					}
				}
			}
			p.Expect("right_paren_token")
			p.Expect("semicolon_token") // expect semicolon after the expression

			nameNode := &ParseNode{Name: "name", Data: funcName, Parent: currentNode}

			callNode.Children = []*ParseNode{nameNode}

			for _, expr := range expressions {
				callNode.Children = append(callNode.Children, expr)
			}

			return callNode
		} else if nextToken.TokenType == "equal_token" || nextToken.TokenType == "left_square_token" || nextToken.TokenType == "period_token" {
			mutNode := &ParseNode{Name: "mut_stmt", Parent: currentNode}

			// mut stmt
			p.Pos--
			varName := p.ParseGetItem(mutNode) // expect identifier and get its value
			p.Expect("equal_token")
			expr := p.ParseExpr(mutNode)
			p.ExpectNoAdvance("semicolon_token") // expect semicolon after the expression

			mutNode.Children = []*ParseNode{varName, expr}

			return mutNode
		}
	}

	if p.Pos == startPos {
		// emergencies only
		p.Pos++
	}

	return &ParseNode{}
}

func (p *Parser) LookUntil(tokenType string, accepted []string) *Token {
	for p.Pos < len(p.Tokens) {
		tok := p.Tokens[p.Pos]
		if tok.TokenType == tokenType {
			return &tok
		} else if !slices.Contains(accepted, tok.TokenType) {
			return nil
		}
	}
	return nil
}

func (p *Parser) ParseAll() {
	p.BasicTypes = []string{
		"any",
		"int",
		"float",
		"bool",
		"string",
		"arr",
		"map",
	}

	p.TreeNode.Name = "start"

	p.TreeNode.Parent = nil
	p.TreeNode.Children = make([]*ParseNode, 0)

	p.Pos = -1

	for p.Pos < len(p.Tokens) {
		p.TreeNode.Children = append(p.TreeNode.Children, p.Parse(&p.TreeNode))
	}
}

func (n *ParseNode) PrintNode(depth int) {
	if len(n.Name) == 0 {
		return
	}

	repr := strings.Repeat("| ", depth) + n.Name
	if len(n.Data) > 0 {
		repr += "  " + n.Data
	}

	fmt.Println(repr)
}

func PrintParseNode(n *ParseNode, depth int) {
	n.PrintNode(depth)

	for _, c := range n.Children {
		PrintParseNode(c, depth+1)
	}
}

func (p *Parser) PrettyPrint() {
	PrintParseNode(&p.TreeNode, 0)
}
