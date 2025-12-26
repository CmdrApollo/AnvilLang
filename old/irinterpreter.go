package main

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func StringValue(v CValue) string {
	if v.IsNull {
		return "null"
	}

	switch v.Type.Type {
	case VAL_INT:
		return strconv.Itoa(v.IntValue)
	case VAL_FLOAT:
		return strconv.FormatFloat(v.FloatValue, 'f', -1, 64)
	case VAL_BOOL:
		return strconv.FormatBool(v.BoolValue)
	case VAL_STRING:
		return v.StringValue
	case VAL_ARR:
		output := "[ "
		for _, element := range v.ArrValue {
			output += StringValue(element) + " "
		}
		output += "]"
		return output
	}
	// sInstance, ok := v.Value.(IStructInstance)
	// if ok {
	// 	output := v.Type.Type + " { "
	// 	for _, element := range sInstance.Fields {
	// 		output += StringValue(element) + " "
	// 	}
	// 	output += "}"
	// 	return output
	// }
	// return fmt.Sprint(v.Value)
	return ""
}

func PrintValue(v CValue, end string) {
	fmt.Print(StringValue(v), end)
}

type IRInterpreter struct {
	Stream         []CInstruction
	Stack          []CValue
	CallStack      []int
	ProgramCounter int
	Labels         map[string]int
	Scopes         []map[string]CValue
	CurrentScope   map[string]CValue
	Builtins       map[string]func([]CValue) *CValue
}

func (r *IRInterpreter) MakeBuiltins() {
	r.Builtins = make(map[string]func([]CValue) *CValue, 0)
	r.Builtins["unpack"] = func(a []CValue) *CValue {
		toUnpack := a[0]
		if toUnpack.Type.Type == VAL_INT {
			value := toUnpack.IntValue
			elements := make([]CValue, 0)
			for j := 0; j < value; j++ {
				elements = append(elements, CValue{Type: CType{Type: VAL_INT, IsNullable: false}, IntValue: j})
			}
			return &CValue{Type: CType{Type: VAL_ARR, Subtype: &CType{Type: VAL_INT, IsNullable: false}},  ArrValue: elements}
		} else if toUnpack.Type.Type == VAL_STRING {
			stringValue := toUnpack.StringValue
			elements := make([]CValue, 0)
			for j := 0; j < len(stringValue); j++ {
				elements = append(elements, CValue{Type: CType{Type: VAL_STRING, IsNullable: false}, StringValue: string(stringValue[j])})
			}
			return &CValue{Type: CType{Type: VAL_ARR, Subtype: &CType{Type: VAL_STRING, IsNullable: false}},  ArrValue: elements}
		}
		return nil
	}
	// r.Builtins["pairs"] = func(a ...any) *IValue {
	// 	toPairify := a[0].(*IValue)
	// 	if toPairify.Type.Type == "map" {
	// 		value := toPairify.Value.(map[any]any)
	// 		elements := make([]any, 0)
	// 		for k, v := range value {
	// 			elements = append(
	// 				elements,
	// 				&IValue{
	// 					Type: IType{Type: "arr", Subtype: IType{Type: "any", Nullable: true}, Nullable: false},
	// 					Value: []any{
	// 						&IValue{Type: toPairify.Type.Subtype.(IType), Value: k},
	// 						&IValue{Type: toPairify.Type.Subtype.(IType).Subtype.(IType), Value: v.(*IValue).Value},
	// 					},
	// 				},
	// 			)
	// 		}
	// 		return &IValue{Type: IType{Type: "arr", Subtype: IType{Type: "arr", Subtype: IType{Type: "any", Nullable: true}}, Nullable: false}, Value: elements}
	// 	}
	// 	return nil
	// }
	// r.Builtins["input"] = func(a ...any) *IValue {
	// 	// FIXME this is a hack
	// 	fmt.Print(a[0].(*IValue).Value)
	// 	name, err := bufio.NewReader(os.Stdin).ReadString('\n')
	// 	if err != nil {
	// 		return nil
	// 	}
	// 	return &IValue{Type: IType{Type: "string", Nullable: false}, Value: name[:len(name)-2]}
	// }
	// r.Builtins["startTerminal"] = func(a ...any) *IValue {
	// 	err := termbox.Init()
	// 	if err != nil {
	// 		r.ThrowError("ERROR: Couldn't start terminal.")
	// 	}
	// 	r.HasOpenedTerminal = true
	// 	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	// 	return nil
	// }
	// r.Builtins["terminalSize"] = func(a ...any) *IValue {
	// 	width, height := termbox.Size()
	// 	widthVal := &IValue{Type: IType{Type: "int", Nullable: false}, Value: width}
	// 	heightVal := &IValue{Type: IType{Type: "int", Nullable: false}, Value: height}
	// 	return &IValue{Type: IType{Type: "arr", Subtype: IType{Type: "int", Nullable: false}, Nullable: false}, Value: []any{widthVal, heightVal}}
	// }
	// r.Builtins["stopTerminal"] = func(a ...any) *IValue {
	// 	termbox.Close()
	// 	return nil
	// }
	// r.Builtins["clearTerminal"] = func(a ...any) *IValue {
	// 	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	// 	return nil
	// }
	// r.Builtins["flushTerminal"] = func(a ...any) *IValue {
	// 	err := termbox.Flush()
	// 	if err != nil {
	// 		r.ThrowError("ERROR: Failed to flush terminal.")
	// 	}
	// 	return nil
	// }
	// r.Builtins["readFromTerminal"] = func(a ...any) *IValue {
	// 	event := termbox.PollEvent()
	// 	key := string(event.Ch)
	// 	specialKeys := map[termbox.Key]string{
	// 		termbox.KeyArrowUp:    "UP",
	// 		termbox.KeyArrowDown:  "DOWN",
	// 		termbox.KeyArrowLeft:  "LEFT",
	// 		termbox.KeyArrowRight: "RIGHT",
	// 		termbox.KeyCtrlC:      "CTRL+C",
	// 	}
	// 	special, ok := specialKeys[event.Key]
	// 	if ok {
	// 		key = special
	// 	}
	// 	return &IValue{Type: IType{Type: "string", Nullable: false}, Value: key}
	// }
	// r.Builtins["writeToTerminal"] = func(a ...any) *IValue {
	// 	x := a[0].(*IValue).Value.(int)
	// 	y := a[1].(*IValue).Value.(int)
	// 	char := a[2].(*IValue).Value.(string)[0]
	// 	colors := map[string]termbox.Attribute{
	// 		"white":   termbox.ColorWhite,
	// 		"black":   termbox.ColorBlack,
	// 		"red":     termbox.ColorRed,
	// 		"green":   termbox.ColorGreen,
	// 		"blue":    termbox.ColorBlue,
	// 		"yellow":  termbox.ColorYellow,
	// 		"cyan":    termbox.ColorCyan,
	// 		"magenta": termbox.ColorMagenta,
	// 		"gray":    termbox.ColorDefault,
	// 	}
	// 	fg := termbox.ColorWhite
	// 	bg := termbox.ColorBlack
	// 	var ok bool
	// 	if len(a) > 3 {
	// 		colorFg := a[3].(*IValue).Value.(string)
	// 		fg, ok = colors[colorFg]
	// 		if !ok {
	// 			r.ThrowError("ERROR: Invalid color passed to 'writeToTerminal'")
	// 		}
	// 	}
	// 	if len(a) > 4 {
	// 		colorBg := a[4].(*IValue).Value.(string)
	// 		bg, ok = colors[colorBg]
	// 		if !ok {
	// 			r.ThrowError("ERROR: Invalid color passed to 'writeToTerminal'")
	// 		}
	// 	}
	// 	r := rune(char)
	// 	termbox.SetCell(x, y, r, fg, bg)
	// 	return nil
	// }
	r.Builtins["print"] = func(a []CValue) *CValue {
		PrintValue(a[0], "")
		return nil
	}
	r.Builtins["println"] = func(a []CValue) *CValue {
		PrintValue(a[0], "\n")
		return nil
	}
	// r.Builtins["parseInt"] = func(a ...any) *IValue {
	// 	value, err := strconv.Atoi(a[0].(*IValue).Value.(string))
	// 	if err != nil {
	// 		return &IValue{Type: IType{Type: "int", Nullable: true}, Value: nil}
	// 	}
	// 	return &IValue{Type: IType{Type: "int", Nullable: true}, Value: value}
	// }
	// r.Builtins["parseFloat"] = func(a ...any) *IValue {
	// 	value, err := strconv.ParseFloat(a[0].(*IValue).Value.(string), 64)
	// 	if err != nil {
	// 		return &IValue{Type: IType{Type: "float", Nullable: true}, Value: nil}
	// 	}
	// 	return &IValue{Type: IType{Type: "float", Nullable: true}, Value: value}
	// }
	r.Builtins["toString"] = func(a []CValue) *CValue {
		value := StringValue(a[0])
		return &CValue{Type: CType{Type: VAL_STRING, IsNullable: false}, StringValue: value}
	}
	r.Builtins["getTime"] = func(a []CValue) *CValue {
		now := time.Now()
		value := float64(now.UnixNano()) / 1e9
		return &CValue{Type: CType{Type: VAL_FLOAT, IsNullable: false}, FloatValue: value}
	}
	r.Builtins["length"] = func(a []CValue) *CValue {
		if !a[0].IsNull {
			if a[0].Type.Type == VAL_ARR {
				return &CValue{Type: CType{Type: VAL_INT, IsNullable: true}, IntValue: len(a[0].ArrValue)}
			} else if a[0].Type.Type == VAL_STRING {
				return &CValue{Type: CType{Type: VAL_INT, IsNullable: true}, IntValue: len(a[0].StringValue)}
			}
		}
		return nil
	}
	r.Builtins["append"] = func(a []CValue) *CValue {
		// todo work in subtype
		if len(a) < 2 {
			return &CValue{Type: CType{Type: VAL_ARR, Subtype: &CType{Type: VAL_ANY}, IsNullable: false}, ArrValue: []CValue{}}
		}

		v := a[0]
		t := v.Type

		return &CValue{Type: t, ArrValue: append(v.ArrValue, a[1])}
	}
	// // TODO pop and popAt
	// r.Builtins["unicode"] = func(a ...any) *IValue {
	// 	value, ok := a[0].(*IValue).Value.(int)
	// 	if !ok {
	// 		return &IValue{Type: IType{Type: "string", Nullable: true}, Value: nil}
	// 	}
	// 	return &IValue{Type: IType{Type: "string", Nullable: true}, Value: string(value)}
	// }
	// r.Builtins["ordinal"] = func(a ...any) *IValue {
	// 	value, ok := a[0].(*IValue).Value.(string)
	// 	if !ok {
	// 		return &IValue{Type: IType{Type: "int", Nullable: true}, Value: nil}
	// 	}
	// 	return &IValue{Type: IType{Type: "int", Nullable: true}, Value: int(value[0])}
	// }
	// r.Builtins["randomValue"] = func(a ...any) *IValue {
	// 	return &IValue{Type: IType{Type: "float", Nullable: false}, Value: rand.Float64()}
	// }
	r.Builtins["mathSin"] = func(a []CValue) *CValue {
		return &CValue{Type: CType{Type: VAL_FLOAT, IsNullable: false}, FloatValue: math.Sin(a[0].FloatValue)}
	}
	r.Builtins["mathCos"] = func(a []CValue) *CValue {
		return &CValue{Type: CType{Type: VAL_FLOAT, IsNullable: false}, FloatValue: math.Cos(a[0].FloatValue)}
	}
	r.Builtins["mathFloor"] = func(a []CValue) *CValue {
		return &CValue{Type: CType{Type: VAL_INT, IsNullable: false}, IntValue: int(a[0].FloatValue)}
	}
	r.Builtins["min"] = func(a []CValue) *CValue {
		return &CValue{Type: CType{Type: VAL_FLOAT, IsNullable: false}, FloatValue: math.Min(a[0].FloatValue, a[1].FloatValue)}
	}
	r.Builtins["max"] = func(a []CValue) *CValue {
		return &CValue{Type: CType{Type: VAL_FLOAT, IsNullable: false}, FloatValue: math.Max(a[0].FloatValue, a[1].FloatValue)}
	}
	// r.Builtins["mathCeil"] = func(a ...any) *IValue {
	// 	floatVal1 := toFloat(a[0].(*IValue))
	// 	return &IValue{Type: IType{Type: "int", Nullable: false}, Value: int(floatVal1 + 1)}
	// }
	// r.Builtins["mathPow"] = func(a ...any) *IValue {
	// 	floatVal1 := toFloat(a[0].(*IValue))
	// 	floatVal2 := toFloat(a[1].(*IValue))
	// 	return &IValue{Type: IType{Type: "float", Nullable: false}, Value: math.Pow(floatVal1, floatVal2)}
	// }
	// r.Builtins["min"] = func(a ...any) *IValue {
	// 	floatVal1 := toFloat(a[0].(*IValue))
	// 	floatVal2 := toFloat(a[1].(*IValue))
	// 	return &IValue{Type: IType{Type: "float", Nullable: false}, Value: math.Min(floatVal1, floatVal2)}
	// }
	// r.Builtins["max"] = func(a ...any) *IValue {
	// 	floatVal1 := toFloat(a[0].(*IValue))
	// 	floatVal2 := toFloat(a[1].(*IValue))
	// 	return &IValue{Type: IType{Type: "float", Nullable: false}, Value: math.Max(floatVal1, floatVal2)}
	// }
	// r.Builtins["abs"] = func(a ...any) *IValue {
	// 	val := a[0].(*IValue).Value.(int)
	// 	if val < 0 {
	// 		val *= -1
	// 	}
	// 	return &IValue{Type: IType{Type: "int", Nullable: false}, Value: val}
	// }
	// r.Builtins["floatAbs"] = func(a ...any) *IValue {
	// 	val := toFloat(a[0].(*IValue))
	// 	if val < 0 {
	// 		val *= -1
	// 	}
	// 	return &IValue{Type: IType{Type: "float", Nullable: false}, Value: val}
	// }
}

func (r *IRInterpreter) Push(x CValue) {
	r.Stack = append(r.Stack, x)
}

func (r *IRInterpreter) Pop() CValue {
	x := r.Stack[len(r.Stack)-1]
	r.Stack = r.Stack[:len(r.Stack)-1]
	return x
}

func (r *IRInterpreter) ThrowError(msg string) {
	panic(msg)
}

func (r *IRInterpreter) LookupVar(name string) CValue {
	value, ok := r.CurrentScope[name]
	if ok {
		return value
	}
	for j := len(r.Scopes) - 1; j > -1; j-- {
		scope := r.Scopes[j]
		value, ok := scope[name]
		if ok {
			return value
		}
	}
	r.ThrowError(fmt.Sprintf("ERROR: No such variable '%s'.", name))
	return CValue{}
}

/*
	OP_CALL
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
	OP_LABEL
*/

func (r *IRInterpreter) DoGetItem() {
	item := r.Pop()
	iterable := r.Pop()

	if item.Type.Type != VAL_INT {
		r.ThrowError("ERROR: IR Interpreter currently doesn't support maps or non-integer indexing.")
	}

	if iterable.Type.Type == VAL_ARR {
		r.Push(iterable.ArrValue[item.IntValue])
	} else if iterable.Type.Type == VAL_STRING {
		r.Push(
			CValue{
				Type: iterable.Type,
				StringValue: string(iterable.StringValue[item.IntValue]),
			},
		)
	}
}

func (r *IRInterpreter) Run() {
	r.Scopes = make([]map[string]CValue, 0)
	r.CurrentScope = make(map[string]CValue)
	r.Scopes = append(r.Scopes, r.CurrentScope)

	r.Stack = make([]CValue, 0)
	r.CallStack = make([]int, 0)

	r.Labels = make(map[string]int)

	// FIXME this should be handled in the compiler, not the IR interpreter
	finalInstructions := make([]CInstruction, 0)

	r.ProgramCounter = 0

	for r.ProgramCounter < len(r.Stream) {
		instr := r.Stream[r.ProgramCounter]
		if instr.Op == OP_LABEL {
			r.Labels[instr.A.StringValue] = len(finalInstructions) - 1
		} else {
			finalInstructions = append(finalInstructions, instr)
		}
		r.ProgramCounter++
	}

	r.ProgramCounter = 0
	for r.ProgramCounter < len(finalInstructions) {
		instruction := finalInstructions[r.ProgramCounter]

		switch instruction.Op {
		case OP_HALT:
			r.ProgramCounter = len(finalInstructions)
		case OP_DUP:
			// duplicate the top of the stack
			r.Push(r.Stack[len(r.Stack)-1])
		case OP_CONST:
			r.Push(instruction.A)
		case OP_PACK:
			array := make([]CValue, 0)
			for i := 0; i < instruction.A.IntValue; i++ {
				array = append(array, r.Pop())
			}
			r.Push(CValue{Type: CType{Type: VAL_ARR}, ArrValue: array})
		case OP_PUSH:
			value := r.LookupVar(instruction.A.StringValue)
			r.Push(value)
		case OP_STORE:
			value := r.Pop()
			r.CurrentScope[instruction.A.StringValue] = value
		case OP_GETITEM:
			r.DoGetItem()
		case OP_ADD:
			b := r.Pop()
			a := r.Pop()
			if a.Type.Type == VAL_INT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: a.Type, IntValue: a.IntValue + b.IntValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: a.Type, FloatValue: a.FloatValue + float64(b.IntValue)})
			} else if a.Type.Type == VAL_INT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: b.Type, FloatValue: float64(a.IntValue) + b.FloatValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: b.Type, FloatValue: a.FloatValue + b.FloatValue})
			} else if a.Type.Type == VAL_STRING && b.Type.Type == VAL_STRING{
				r.Push(CValue{Type: a.Type, StringValue: a.StringValue + b.StringValue})	
			} else {
				r.ThrowError("ERROR: Attempt to add values with types unsupported by the add operator.")
			}
		case OP_SUB:
			b := r.Pop()
			a := r.Pop()
			if a.Type.Type == VAL_INT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: a.Type, IntValue: a.IntValue - b.IntValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: a.Type, FloatValue: a.FloatValue - float64(b.IntValue)})
			} else if a.Type.Type == VAL_INT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: b.Type, FloatValue: float64(a.IntValue) - b.FloatValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: b.Type, FloatValue: a.FloatValue - b.FloatValue})
			} else {
				r.ThrowError("ERROR: Attempt to subtract values with types unsupported by the subtract operator.")
			}
		case OP_MUL:
			b := r.Pop()
			a := r.Pop()
			if a.Type.Type == VAL_INT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: a.Type, IntValue: a.IntValue * b.IntValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: a.Type, FloatValue: a.FloatValue * float64(b.IntValue)})
			} else if a.Type.Type == VAL_INT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: b.Type, FloatValue: float64(a.IntValue) * b.FloatValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: b.Type, FloatValue: a.FloatValue * b.FloatValue})
			} else {
				r.ThrowError("ERROR: Attempt to multiply values with types unsupported by the multiply operator.")
			}
		case OP_DIV:
			b := r.Pop()
			a := r.Pop()
			if a.Type.Type == VAL_INT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: a.Type, IntValue: a.IntValue / b.IntValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: a.Type, FloatValue: a.FloatValue / float64(b.IntValue)})
			} else if a.Type.Type == VAL_INT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: b.Type, FloatValue: float64(a.IntValue) / b.FloatValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: b.Type, FloatValue: a.FloatValue / b.FloatValue})
			} else {
				r.ThrowError("ERROR: Attempt to divide values with types unsupported by the divide operator.")
			}
		case OP_MOD:
			b := r.Pop()
			a := r.Pop()
			if a.Type.Type == VAL_INT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: a.Type, IntValue: a.IntValue % b.IntValue})
			} else {
				r.ThrowError("ERROR: Attempt to modulo values with types unsupported by the modulo operator.")
			}
		case OP_MORE:
			b := r.Pop()
			a := r.Pop()
			if a.Type.Type == VAL_INT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.IntValue > b.IntValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.FloatValue > float64(b.IntValue)})
			} else if a.Type.Type == VAL_INT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: float64(a.IntValue) > b.FloatValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.FloatValue > b.FloatValue})
			} else {
				r.ThrowError("ERROR: Attempt to compare values with types unsupported by comparison.")
			}
		case OP_LESS:
			b := r.Pop()
			a := r.Pop()
			if a.Type.Type == VAL_INT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.IntValue < b.IntValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.FloatValue < float64(b.IntValue)})
			} else if a.Type.Type == VAL_INT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: float64(a.IntValue) < b.FloatValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.FloatValue < b.FloatValue})
			} else {
				r.ThrowError("ERROR: Attempt to compare values with types unsupported by comparison.")
			}
		case OP_MORE_EQUAL:
			b := r.Pop()
			a := r.Pop()
			if a.Type.Type == VAL_INT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.IntValue >= b.IntValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.FloatValue >= float64(b.IntValue)})
			} else if a.Type.Type == VAL_INT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: float64(a.IntValue) >= b.FloatValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.FloatValue >= b.FloatValue})
			} else {
				r.ThrowError("ERROR: Attempt to compare values with types unsupported by comparison.")
			}
		case OP_LESS_EQUAL:
			b := r.Pop()
			a := r.Pop()
			if a.Type.Type == VAL_INT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.IntValue <= b.IntValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.FloatValue <= float64(b.IntValue)})
			} else if a.Type.Type == VAL_INT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: float64(a.IntValue) <= b.FloatValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.FloatValue <= b.FloatValue})
			} else {
				r.ThrowError("ERROR: Attempt to compare values with types unsupported by comparison.")
			}
		case OP_EQUAL:
			b := r.Pop()
			a := r.Pop()
			if a.Type.Type == VAL_INT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.IntValue == b.IntValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.FloatValue == float64(b.IntValue)})
			} else if a.Type.Type == VAL_INT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: float64(a.IntValue) == b.FloatValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.FloatValue == b.FloatValue})
			} else {
				r.ThrowError("ERROR: Attempt to compare values with types unsupported by comparison.")
			}
		case OP_NOT_EQUAL:
			b := r.Pop()
			a := r.Pop()
			if a.Type.Type == VAL_INT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.IntValue != b.IntValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_INT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.FloatValue != float64(b.IntValue)})
			} else if a.Type.Type == VAL_INT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: float64(a.IntValue) != b.FloatValue})
			} else if a.Type.Type == VAL_FLOAT && b.Type.Type == VAL_FLOAT {
				r.Push(CValue{Type: CType{Type: VAL_BOOL}, BoolValue: a.FloatValue != b.FloatValue})
			} else {
				r.ThrowError("ERROR: Attempt to compare values with types unsupported by comparison.")
			}
		case OP_JUMP:
			r.ProgramCounter = r.Labels[instruction.A.StringValue]
		case OP_JUMP_IF:
			cond := r.Pop()
			if cond.BoolValue {
				r.ProgramCounter = r.Labels[instruction.A.StringValue]
			}
		case OP_JUMP_IF_NOT:
			cond := r.Pop()
			if !cond.BoolValue {
				r.ProgramCounter = r.Labels[instruction.A.StringValue]
			}
		case OP_CALL_BUILTIN:
			// right now just assume the builtin being called is println
			funcName := instruction.A.StringValue
			arguments := make([]CValue, 0)
			for i := 0; i < instruction.B.IntValue; i++ {
				arguments = append(arguments, r.Pop())
			}
			returned := r.Builtins[funcName](arguments)
			if returned != nil {
				r.Push(*returned)
			}
		}
		r.ProgramCounter++
	}
}