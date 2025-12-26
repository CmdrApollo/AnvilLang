package main

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
)

type CValueType uint8
type OpCode uint8

const (
	VAL_ANY CValueType = iota
	VAL_INT
	VAL_FLOAT
	VAL_BOOL
	VAL_STRING
	VAL_ARR
	VAL_MAP
)

const (
	OP_NOP OpCode = iota
	OP_HALT
	OP_CONST
	OP_LOAD
	OP_STORE
	OP_DUP
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_MOD
	OP_JUMP
	OP_JUMP_IF
	OP_JUMP_IF_NOT
	OP_CALL
	OP_CALL_BUILTIN
	OP_RETURN
	OP_AND
	OP_OR
	OP_NOT
	OP_PUSH
	OP_POP
	OP_CONCAT
	OP_EQUAL
	OP_NOT_EQUAL
	OP_LESS
	OP_LESS_EQUAL
	OP_MORE
	OP_MORE_EQUAL
	OP_GETITEM
	OP_PACK
	OP_LABEL
)

func (o OpCode) ToString() string {
	return []string{
		"OP_NOP",
		"OP_HALT",
		"OP_CONST",
		"OP_LOAD",
		"OP_STORE",
		"OP_DUP",
		"OP_ADD",
		"OP_SUB",
		"OP_MUL",
		"OP_DIV",
		"OP_MOD",
		"OP_JUMP",
		"OP_JUMP_IF",
		"OP_JUMP_IF_NOT",
		"OP_CALL",
		"OP_CALL_BUILTIN",
		"OP_RETURN",
		"OP_AND",
		"OP_OR",
		"OP_NOT",
		"OP_PUSH",
		"OP_POP",
		"OP_CONCAT",
		"OP_EQUAL",
		"OP_NOT_EQUAL",
		"OP_LESS",
		"OP_LESS_EQUAL",
		"OP_MORE",
		"OP_MORE_EQUAL",
		"OP_GETITEM",
		"OP_PACK",
		"OP_LABEL",
	}[int(o)]
}

type CType struct {
	Type       CValueType
	Subtype    *CType
	KeyType    *CType
	ValueType  *CType
	IsNullable bool
}

type CValue struct {
	Type CType

	IsNull bool

	IntValue    int
	FloatValue  float64
	BoolValue   bool
	StringValue string
	ArrValue    []CValue
	MapValue    map[string]CValue
}

func (v CValue) ToString() string {
	switch v.Type.Type {
	case VAL_INT:
		return strconv.Itoa(v.IntValue)
	case VAL_FLOAT:
		return strconv.FormatFloat(v.FloatValue, 'f', -1, 64)
	case VAL_BOOL:
		return strconv.FormatBool(v.BoolValue)
	case VAL_STRING:
		return v.StringValue
	}
	return "< >"
}

type CInstruction struct {
	Op OpCode
	A  CValue
	B  CValue
}

type CParam struct {
	Name       string
	Type       CValueType
	IsNullable bool
	Default    *CValue
}

type CFunction struct {
	Params     []CParam
	ReturnType CValueType
	Suite      []CInstruction
}

type CStruct struct {
	Fields []CParam
}

type Compiler struct {
	Tree           ParseNode
	Variables      map[string]bool
	Functions      map[string]CFunction
	Structs        map[string]CStruct
	Builtins       []string
	ReturnFromFunc bool
	ReturnValue    CValue
	BreakFromLoop  bool
	ContinueLoop   bool
	BasicTypes 	   []string
	HiLevelOutput  map[string][]CInstruction
	CurrentOutput  []CInstruction
	Labels         []string
	FinalOutput    []CInstruction
}

func (c *Compiler) ThrowError(msg string) {
	panic(msg)
}

func (c *Compiler) MakeBuiltins() {
	c.Builtins = []string{
		"unpack",
		"pairs",
		"input",
		"startTerminal",
		"terminalSize",
		"stopTerminal",
		"clearTerminal",
		"flushTerminal",
		"readFromTerminal",
		"writeToTerminal",
		"print",
		"println",
		"parseInt",
		"parseFloat",
		"toString",
		"getTime",
		"length",
		"append",
		"unicode",
		"ordinal",
		"randomValue",
		"mathSin",
		"mathCos",
		"mathTan",
		"mathFloor",
		"mathCeil",
		"mathPow",
		"min",
		"max",
		"abs",
		"floatAbs",
	}
}

func (c *Compiler) NewLabel() string {
	label := "_label" + strconv.FormatInt(int64(len(c.Labels)), 16)
	c.Labels = append(c.Labels, label)
	return label
}

func (c *Compiler) DummyVar() string {
	dummy := ""
	for i := 0; i < 10; i++ {
		dummy += string(byte(rand.Intn(256)))
	}
	return dummy
}

func (c * Compiler) MakeLabel(label string) CValue {
	return CValue{
		Type: CType{
			Type: VAL_STRING,
		},
		StringValue: label,
	}
}

func (c * Compiler) MakeInt(x int) CValue {
	return CValue{
		Type: CType{
			Type: VAL_INT,
		},
		IntValue: x,
	}
}

func (c *Compiler) Emit(instruction CInstruction) {
	c.CurrentOutput = append(c.CurrentOutput, instruction)
}

func (c *Compiler) EmitLabel(label string) {
	c.Emit(
		CInstruction{
			Op: OP_LABEL,
			A: c.MakeLabel(label),
		},
	)
}

func (c *Compiler) RunInstruction(node *ParseNode, startLabel string, endLabel string) {
	switch node.Name {
	// literals
	case "int_literal":
		value, err := strconv.Atoi(node.Data)
		if err != nil {
			return
		}
		c.Emit(
			CInstruction{
				Op: OP_CONST,
				A: CValue{
					Type: CType{
						Type: VAL_INT,
						IsNullable: false,
					},
					IntValue: value,
				},
			},
		)
		break
	case "float_literal":
		value, err := strconv.ParseFloat(node.Data, 64)
		if err != nil {
			return
		}
		c.Emit(
			CInstruction{
				Op: OP_CONST,
				A: CValue{
					Type: CType{
						Type: VAL_FLOAT,
						IsNullable: false,
					},
					FloatValue: value,
				},
			},
		)
		break
	case "bool_literal", "nullable":
		value, err := strconv.ParseBool(node.Data)
		if err != nil {
			return
		}
		c.Emit(
			CInstruction{
				Op: OP_CONST,
				A: CValue{
					Type: CType{
						Type: VAL_BOOL,
						IsNullable: false,
					},
					BoolValue: value,
				},
			},
		)
		break
	case "string_literal":
		c.Emit(
			CInstruction{
				Op: OP_CONST,
				A: CValue{
					Type: CType{
						Type: VAL_STRING,
						IsNullable: false,
					},
					StringValue: node.Data,
				},
			},
		)
		break
	case "array":
		for i := len(node.Children) - 1; i >= 0; i--{
			c.RunInstruction(node.Children[i], startLabel, endLabel)
		}
		c.Emit(
			CInstruction{
				Op: OP_PACK,
				A: CValue{
					Type: CType{
						Type: VAL_INT,
						IsNullable: false,
					},
					IntValue: len(node.Children),
				},
			},
		)
	// variable stuff
	case "variable_decl":
		varName := node.Children[0].Data

		_, ok := c.Variables[varName]

		if ok {
			c.ThrowError(fmt.Sprintf("ERROR: Attempt to multiply-initialize variable '%s'.", varName));
		}
		
		// type expection
		c.RunInstruction(node.Children[1], startLabel, endLabel)

		// value
		c.RunInstruction(node.Children[2], startLabel, endLabel)
		
		c.Emit(
			CInstruction{
				Op: OP_STORE,
				A: c.MakeLabel(varName),
			},
		)

		c.Variables[varName] = true
	case "mut_stmt":
		varName := node.Children[0].Data
		c.RunInstruction(node.Children[1], startLabel, endLabel)
		
		c.Emit(
			CInstruction{
				Op: OP_STORE,
				A: c.MakeLabel(varName),
			},
		)
	case "name":
		if node.Data == "null" {
			c.Emit(
				CInstruction{
					Op: OP_CONST,
					A: CValue{
						Type: CType{
							Type: VAL_ANY,
							IsNullable: true,
						},
						IsNull: true,
					},
				},
			)
		} else {
			c.Emit(
				CInstruction{
					Op: OP_PUSH,
					A: c.MakeLabel(node.Data),
				},
			)
		}
		break
	// math stuff
	case "addsub_operation":
		c.RunInstruction(node.Children[0], startLabel, endLabel)
		c.RunInstruction(node.Children[2], startLabel, endLabel)
		op := node.Children[1].Data
		if op == "plus" {
			c.Emit(
				CInstruction{
					Op: OP_ADD,
				},
			)
		} else {
			c.Emit(
				CInstruction{
					Op: OP_SUB,
				},
			)
		}
		break
	case "unary_negate_operation":
		// 0
		c.Emit(
			CInstruction{
				Op: OP_PUSH,
				A: CValue{
					Type: CType{
						Type: VAL_INT,
						IsNullable: false,
					},
					IntValue: 0,
				},
			},
		)
		// whatever to negate
		c.RunInstruction(node.Children[0], startLabel, endLabel)
		// emit 0 - X
		c.Emit(
			CInstruction{
				Op: OP_SUB,
			},
		)
		break
	case "muldiv_operation":
		c.RunInstruction(node.Children[0], startLabel, endLabel)
		c.RunInstruction(node.Children[2], startLabel, endLabel)
		op := node.Children[1].Data
		if op == "star" {
			c.Emit(
				CInstruction{
					Op: OP_MUL,
				},
			)
		} else if op == "slash" {
			c.Emit(
				CInstruction{
					Op: OP_DIV,
				},
			)
		} else {
			c.Emit(
				CInstruction{
					Op: OP_MOD,
				},
			)
		}
		break
	case "comparison_operation":
		c.RunInstruction(node.Children[0], startLabel, endLabel)
		c.RunInstruction(node.Children[2], startLabel, endLabel)
		op := node.Children[1].Data

		var opcode OpCode

		switch op {
		case "greater":
			opcode = OP_MORE
		case "lesser":
			opcode = OP_LESS
		case "greater_equal":
			opcode = OP_MORE_EQUAL
		case "lesser_equal":
			opcode = OP_LESS_EQUAL
		}

		c.Emit(
			CInstruction{
				Op: opcode,
			},
		)
	case "or_operation":
		c.RunInstruction(node.Children[0], startLabel, endLabel)
		c.RunInstruction(node.Children[2], startLabel, endLabel)	

		c.Emit(
			CInstruction{
				Op: OP_OR,
			},
		)
	case "and_operation":
		c.RunInstruction(node.Children[0], startLabel, endLabel)
		c.RunInstruction(node.Children[2], startLabel, endLabel)	

		c.Emit(
			CInstruction{
				Op: OP_AND,
			},
		)
	case "not_operation":
		c.RunInstruction(node.Children[0], startLabel, endLabel)

		c.Emit(
			CInstruction{
				Op: OP_NOT,
			},
		)
	case "equality_operation":
		c.RunInstruction(node.Children[0], startLabel, endLabel)
		c.RunInstruction(node.Children[2], startLabel, endLabel)
		op := node.Children[1].Data
		if op == "equal_equal" {
			c.Emit(
				CInstruction{
					Op: OP_EQUAL,
				},
			)
		} else {
			c.Emit(
				CInstruction{
					Op: OP_NOT_EQUAL,
				},
			)
		}
	// get_item
	case "get_item":
		c.RunInstruction(node.Children[0], startLabel, endLabel)
		c.RunInstruction(node.Children[1], startLabel, endLabel)
		c.Emit(
			CInstruction{
				Op: OP_GETITEM,
			},
		)
	// control flow
	case "while_stmt":
		/*
		label_start:
			EVAL condition
			JUMP_IF_NOT label_done
			... body ...
			JUMP label_start
		label_done:
			...
		*/
		labelStart := c.NewLabel()
		labelDone := c.NewLabel()

		suite := node.Children[1]

		c.EmitLabel(labelStart)

		c.RunInstruction(node.Children[0], startLabel, endLabel)

		c.Emit(
			CInstruction{
				Op: OP_JUMP_IF_NOT,
				A: c.MakeLabel(labelDone),
			},
		)

		for _, child := range suite.Children {
			c.RunInstruction(child, labelStart, labelDone)
		}

		c.Emit(
			CInstruction{
				Op: OP_JUMP,
				A: c.MakeLabel(labelStart),
			},
		)

		c.EmitLabel(labelDone)
		break
	case "loop_stmt":
		/*
		CONST 0
		STORE newLabel
		CONST [...]
		STORE arr
		label_start
			PUSH arr
			PUSH newLabel
			GET_ITEM
			STORE varName
			... body ...
			PUSH newLabel
			CONST 1
			ADD
			STORE newLabel
			PUSH newLabel
			PUSH arr
			CALL_BUILTIN length 1
			LESS
			JUMP_IF label_start
		...
		*/
		labelStart := c.NewLabel()
		labelDone := c.NewLabel()
		arrayVariable := c.DummyVar()
		indexVariable := c.DummyVar()

		varName := node.Children[1].Data

		// TODO give a shit about type, again

		suite := node.Children[3]

		c.Emit(
			CInstruction{
				Op: OP_CONST,
				A: c.MakeInt(0),
			},
		)
		
		c.Emit(
			CInstruction{
				Op: OP_STORE,
				A: c.MakeLabel(indexVariable),
			},
		)

		c.RunInstruction(node.Children[0], startLabel, endLabel)
		
		c.Emit(
			CInstruction{
				Op: OP_STORE,
				A: c.MakeLabel(arrayVariable),
			},
		)

		c.EmitLabel(labelStart)	
		
		c.Emit(
			CInstruction{
				Op: OP_PUSH,
				A: c.MakeLabel(arrayVariable),
			},
		)

		c.Emit(
			CInstruction{
				Op: OP_PUSH,
				A: c.MakeLabel(indexVariable),
			},
		)

		c.Emit(
			CInstruction{
				Op: OP_GETITEM,
			},
		)
		
		c.Emit(
			CInstruction{
				Op: OP_STORE,
				A: c.MakeLabel(varName),
			},
		)

		// suite

		for _, child := range suite.Children {
			c.RunInstruction(child, labelStart, labelDone)
		}

		c.Emit(
			CInstruction{
				Op: OP_PUSH,
				A: c.MakeLabel(indexVariable),
			},
		)

		c.Emit(
			CInstruction{
				Op: OP_CONST,
				A: c.MakeInt(1),
			},
		)

		c.Emit(
			CInstruction{
				Op: OP_ADD,
			},
		)

		c.Emit(
			CInstruction{
				Op: OP_STORE,
				A: c.MakeLabel(indexVariable),
			},
		)

		c.Emit(
			CInstruction{
				Op: OP_PUSH,
				A: c.MakeLabel(indexVariable),
			},
		)
		

		c.Emit(
			CInstruction{
				Op: OP_PUSH,
				A: c.MakeLabel(arrayVariable),
			},
		)

		c.Emit(
			CInstruction{
				Op: OP_CALL_BUILTIN,
				A: c.MakeLabel("length"),
				B: c.MakeInt(1),
			},
		)

		c.Emit(
			CInstruction{
				Op: OP_LESS,
			},
		)

		c.Emit(
			CInstruction{
				Op: OP_JUMP_IF,
				A: c.MakeLabel(labelStart),
			},
		)

		c.EmitLabel(labelDone)

		break
	case "continue_stmt":
		c.Emit(
			CInstruction{
				Op: OP_JUMP,
				A: c.MakeLabel(startLabel),
			},
		)
	case "break_stmt":
		c.Emit(
			CInstruction{
				Op: OP_JUMP,
				A: c.MakeLabel(endLabel),
			},
		)
	case "return_stmt":
		c.RunInstruction(node.Children[0], startLabel, endLabel)
		c.Emit(
			CInstruction{
				Op: OP_RETURN,
			},
		)
	case "if_stmt":
		/*
		EVAL condition
		JUMP_IF_NOT label_false
		... body ...
		JUMP label_done
		label_false:
			...
		label_done:
			...
		*/
		suite := node.Children[1]
		elseSuite := node.Children[2]

		labelFalse := c.NewLabel()
		labelDone := c.NewLabel()

		// EVAL condition
		c.RunInstruction(node.Children[0], startLabel, endLabel)

		// JUMP_IF_NOT label_false
		c.Emit(
			CInstruction{
				Op: OP_JUMP_IF_NOT,
				A: c.MakeLabel(labelFalse),
			},
		)

		// truthy body

		for _, child := range suite.Children {
			c.RunInstruction(child, startLabel, endLabel)
		}

		c.Emit(
			CInstruction{
				Op: OP_JUMP,
				A: c.MakeLabel(labelDone),
			},
		)

		c.EmitLabel(labelFalse)

		// falsey body

		for _, child := range elseSuite.Children {
			c.RunInstruction(child, startLabel, endLabel)
		}

		c.EmitLabel(labelDone)

		break
	// function bullshit
	case "func_call":
		funcName := node.Children[0].Data

		// pushing arguments onto the stack
		for j := 1; j < len(node.Children); j++ {
			c.RunInstruction(node.Children[j], startLabel, endLabel)
		}

		_, ok := c.Functions[funcName]

		if ok {
			// emit user func
			c.Emit(
				CInstruction{
					Op: OP_CALL,
					A: CValue{
						Type: CType{
							Type: VAL_STRING,
						}, // name of func
						StringValue: funcName,
					},
					B: CValue{
						Type: CType{
							Type: VAL_INT,
						},  // number of args
						IntValue: len(node.Children) - 1,
					},
				},
			)
		} else {
			ok = slices.Contains(c.Builtins, funcName)

			if !ok {
				c.ThrowError(fmt.Sprintf("ERROR: No such function: '%s'.", funcName))
			} else {
				// emit native func
				c.Emit(
					CInstruction{
						Op: OP_CALL_BUILTIN,
						A: CValue{
							Type: CType{
								Type: VAL_STRING,
							}, // name of func
							StringValue: funcName,
						},
						B: CValue{
							Type: CType{
								Type: VAL_INT,
							}, // number of args
							IntValue: len(node.Children) - 1,
						},
					},
				)
			}
		}
		break
	}
}

func (c *Compiler) Start() {
	c.Functions = make(map[string]CFunction)
	c.Structs = make(map[string]CStruct)

	c.Variables = make(map[string]bool)

	c.BasicTypes = []string{
		"any",
		"int",
		"float",
		"bool",
		"string",
		"arr",
		"map",
	}

	c.RunTree()
}

func (c *Compiler) RunTree() {
	for _, node := range c.Tree.Children {
		c.RunInstruction(node, "", "")
	}
	c.CurrentOutput = append(c.CurrentOutput, CInstruction{Op: OP_NOP})
	c.FinalOutput = c.CurrentOutput
}

func (c *Compiler) PrintHiLevelOutput() {
	for _, instruction := range c.CurrentOutput {
		fmt.Println(
			fmt.Sprintf("%s %s %s",
				instruction.Op.ToString(),
				instruction.A.ToString(),
				instruction.B.ToString(),
			),
		)
	}
}