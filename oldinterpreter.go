package main

// import (
// 	"bufio"
// 	"fmt"
// 	"math"
// 	"math/rand"
// 	"os"
// 	"slices"
// 	"strconv"
// 	"time"

// 	"github.com/gdamore/tcell/v2/termbox"
// )

// type IType struct {
// 	Type     string
// 	Subtype  any
// 	Nullable bool
// }

// func (t IType) DefaultValue() *IValue {
// 	switch t.Type {
// 	case "int":
// 		return &IValue{Type: t, Value: 0}
// 	case "float":
// 		return &IValue{Type: t, Value: 0.0}
// 	case "string":
// 		return &IValue{Type: t, Value: ""}
// 	case "bool":
// 		return &IValue{Type: t, Value: false}
// 	case "arr":
// 		return &IValue{Type: t, Value: make([]*IValue, 0)}
// 	case "map":
// 		return &IValue{Type: t, Value: make(map[any]any, 0)}
// 	}
// 	return &IValue{Type: IType{Type: "any", Nullable: true}, Value: nil}
// }

// func TypeMismatch(a IType, b IType) bool {
// 	if a.Subtype != nil && b.Subtype != nil {
// 		if a.Type == "any" || b.Type == "any" {
// 			return false
// 		}
// 		if a.Subtype.(IType).Type == "any" || b.Subtype.(IType).Type == "any" {
// 			return false
// 		}
// 		if a.Nullable != b.Nullable || a.Subtype.(IType).Nullable != b.Subtype.(IType).Nullable {
// 			return true
// 		}
// 		return !(a.Type == b.Type && a.Subtype.(IType).Type == b.Subtype.(IType).Type)
// 	}
// 	if a.Type == "any" || b.Type == "any" {
// 		return false
// 	}
// 	if a.Nullable != b.Nullable {
// 		return true
// 	}
// 	return !(a.Type == b.Type)
// }

// type IValue struct {
// 	Type  IType
// 	Value any
// }

// func StringValue(v *IValue) string {
// 	if v.Value == nil {
// 		return "null"
// 	}

// 	switch v.Type.Type {
// 	case "int", "float", "bool", "string":
// 		return fmt.Sprint(v.Value)
// 	case "arr":
// 		output := "[ "
// 		for _, element := range v.Value.([]any) {
// 			output += StringValue(element.(*IValue)) + " "
// 		}
// 		output += "]"
// 		return output
// 	}
// 	sInstance, ok := v.Value.(IStructInstance)
// 	if ok {
// 		output := v.Type.Type + " { "
// 		for _, element := range sInstance.Fields {
// 			output += StringValue(element) + " "
// 		}
// 		output += "}"
// 		return output
// 	}
// 	return fmt.Sprint(v.Value)
// }

// func PrintValue(v *IValue, end string) {
// 	fmt.Print(StringValue(v), end)
// }

// func toFloat(v *IValue) float64 {
// 	if v.Type.Type == "float" {
// 		return v.Value.(float64)
// 	} else if v.Type.Type == "int" {
// 		intVal, ok := v.Value.(int)
// 		if !ok {
// 			return 0
// 		}
// 		return float64(intVal)
// 	}
// 	return 0
// }

// func MultiplyValues(a IValue, b IValue) (IType, any) {
// 	if a.Type.Type == "float" || b.Type.Type == "float" {
// 		fa := toFloat(&a)
// 		fb := toFloat(&b)
// 		return IType{Type: "float"}, fa * fb
// 	} else if a.Type.Type == "int" && b.Type.Type == "int" {
// 		return IType{Type: "int"}, a.Value.(int) * b.Value.(int)
// 	}
// 	return IType{Type: "any", Nullable: true}, nil
// }

// func DivideValues(a IValue, b IValue) (IType, any) {
// 	if a.Type.Type == "float" || b.Type.Type == "float" {
// 		fa := toFloat(&a)
// 		fb := toFloat(&b)
// 		return IType{Type: "float"}, fa / fb
// 	} else if a.Type.Type == "int" && b.Type.Type == "int" {
// 		return IType{Type: "int"}, a.Value.(int) / b.Value.(int)
// 	}
// 	return IType{Type: "any", Nullable: true}, nil
// }

// func ModuloValues(a IValue, b IValue) (IType, any) {
// 	if a.Type.Type == "int" && b.Type.Type == "int" {
// 		return IType{Type: "int"}, a.Value.(int) % b.Value.(int)
// 	}
// 	return IType{Type: "any", Nullable: true}, nil
// }

// func AddValues(a IValue, b IValue) (IType, any) {
// 	if a.Type.Type == "string" && b.Type.Type == "string" {
// 		if a.Value == nil && b.Value == nil {
// 			return IType{Type: "string"}, ""
// 		} else if a.Value == nil {
// 			return IType{Type: "string"}, b.Value.(string)
// 		} else if b.Value == nil {
// 			return IType{Type: "string"}, a.Value.(string)
// 		}
// 		return IType{Type: "string"}, a.Value.(string) + b.Value.(string)
// 	} else if a.Type.Type == "float" || b.Type.Type == "float" || a.Type.Type == "int" || b.Type.Type == "int" {
// 		fa := toFloat(&a)
// 		fb := toFloat(&b)
// 		if a.Type.Type == "int" && b.Type.Type == "int" {
// 			return IType{Type: "int"}, a.Value.(int) + b.Value.(int)
// 		} else {
// 			return IType{Type: "float"}, fa + fb
// 		}
// 	}
// 	return IType{Type: "any", Nullable: true}, nil
// }

// func SubtractValues(a IValue, b IValue) (IType, any) {
// 	if a.Type.Type == "float" || b.Type.Type == "float" {
// 		fa := toFloat(&a)
// 		fb := toFloat(&b)
// 		return IType{Type: "float"}, fa - fb
// 	} else if a.Type.Type == "int" && b.Type.Type == "int" {
// 		return IType{Type: "int"}, a.Value.(int) - b.Value.(int)
// 	}
// 	return IType{Type: "any", Nullable: true}, nil
// }

// func CompareEquality(a *IValue, op string, b *IValue) bool {
// 	switch op {
// 	case "equal_equal":
// 		if a.Type.Type == b.Type.Type || a.Type.Type == "any" || b.Type.Type == "any" {
// 			return a.Value == b.Value
// 		} else if a.Type.Type == "int" && b.Type.Type == "float" {
// 			return toFloat(a) == toFloat(b)
// 		} else if a.Type.Type == "float" && b.Type.Type == "int" {
// 			return toFloat(a) == toFloat(b)
// 		}
// 		return false
// 	case "not_equal":
// 		if a.Type.Type == b.Type.Type {
// 			return a.Value != b.Value
// 		} else if a.Type.Type == "int" && b.Type.Type == "float" {
// 			return toFloat(a) != toFloat(b)
// 		} else if a.Type.Type == "float" && b.Type.Type == "int" {
// 			return toFloat(a) != toFloat(b)
// 		}
// 		return true
// 	}
// 	return false
// }

// func CompareValues(a *IValue, op string, b *IValue) bool {
// 	var fa, fb float64
// 	if !(a.Type.Type == "int" || a.Type.Type == "float") {
// 		return false
// 	}
// 	if !(b.Type.Type == "int" || b.Type.Type == "float") {
// 		return false
// 	}
// 	fa, fb = toFloat(a), toFloat(b)
// 	switch op {
// 	case "lesser":
// 		return fa < fb
// 	case "greater":
// 		return fa > fb
// 	case "lesser_equal":
// 		return fa <= fb
// 	case "greater_equal":
// 		return fa >= fb
// 	}
// 	return false
// }

// func IsTruthy(value any) bool {
// 	return value == true
// }

// type IParam struct {
// 	Name    string
// 	Type    IType
// 	Default *IValue
// }

// type IFunction struct {
// 	Params     []IParam
// 	ReturnType IType
// 	Suite      ParseNode
// }

// type IStruct struct {
// 	Fields []IParam
// }

// type IStructInstance struct {
// 	Fields map[string]*IValue
// }

// type Interpreter struct {
// 	Tree              ParseNode
// 	Scopes            []map[string]IValue
// 	CurrentScope      map[string]IValue
// 	Functions         map[string]IFunction
// 	Structs           map[string]IStruct
// 	Builtins          map[string]func(...any) *IValue
// 	ReturnFromFunc    bool
// 	ReturnValue       *IValue
// 	BreakFromLoop     bool
// 	ContinueLoop      bool
// 	BasicTypes        []string
// 	HasOpenedTerminal bool
// }

// func (i *Interpreter) ThrowError(msg string) {
// 	panic(msg)
// }

// func (i *Interpreter) MakeBuiltins() {
// 	i.Builtins = make(map[string]func(...any) *IValue, 0)
// 	i.Builtins["unpack"] = func(a ...any) *IValue {
// 		toUnpack := a[0].(*IValue)
// 		if toUnpack.Type.Type == "int" {
// 			value := toUnpack.Value.(int)
// 			elements := make([]any, 0)
// 			for j := 0; j < value; j++ {
// 				elements = append(elements, &IValue{Type: IType{Type: "int", Nullable: false}, Value: j})
// 			}
// 			return &IValue{Type: IType{Type: "arr", Subtype: IType{Type: "int"}, Nullable: false}, Value: elements}
// 		} else if toUnpack.Type.Type == "string" {
// 			stringValue := toUnpack.Value.(string)
// 			elements := make([]any, 0)
// 			for j := 0; j < len(stringValue); j++ {
// 				elements = append(elements, &IValue{Type: IType{Type: "string", Nullable: false}, Value: string(stringValue[j])})
// 			}
// 			return &IValue{Type: IType{Type: "arr", Subtype: IType{Type: "string"}, Nullable: false}, Value: elements}
// 		}
// 		return nil
// 	}
// 	i.Builtins["pairs"] = func(a ...any) *IValue {
// 		toPairify := a[0].(*IValue)
// 		if toPairify.Type.Type == "map" {
// 			value := toPairify.Value.(map[any]any)
// 			elements := make([]any, 0)
// 			for k, v := range value {
// 				elements = append(
// 					elements,
// 					&IValue{
// 						Type: IType{Type: "arr", Subtype: IType{Type: "any", Nullable: true}, Nullable: false},
// 						Value: []any{
// 							&IValue{Type: toPairify.Type.Subtype.(IType), Value: k},
// 							&IValue{Type: toPairify.Type.Subtype.(IType).Subtype.(IType), Value: v.(*IValue).Value},
// 						},
// 					},
// 				)
// 			}
// 			return &IValue{Type: IType{Type: "arr", Subtype: IType{Type: "arr", Subtype: IType{Type: "any", Nullable: true}}, Nullable: false}, Value: elements}
// 		}
// 		return nil
// 	}
// 	i.Builtins["input"] = func(a ...any) *IValue {
// 		// FIXME this is a hack
// 		fmt.Print(a[0].(*IValue).Value)
// 		name, err := bufio.NewReader(os.Stdin).ReadString('\n')
// 		if err != nil {
// 			return nil
// 		}
// 		return &IValue{Type: IType{Type: "string", Nullable: false}, Value: name[:len(name)-2]}
// 	}
// 	i.Builtins["startTerminal"] = func(a ...any) *IValue {
// 		err := termbox.Init()
// 		if err != nil {
// 			i.ThrowError("ERROR: Couldn't start terminal.")
// 		}
// 		i.HasOpenedTerminal = true
// 		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
// 		return nil
// 	}
// 	i.Builtins["terminalSize"] = func(a ...any) *IValue {
// 		width, height := termbox.Size()
// 		widthVal := &IValue{Type: IType{Type: "int", Nullable: false}, Value: width}
// 		heightVal := &IValue{Type: IType{Type: "int", Nullable: false}, Value: height}
// 		return &IValue{Type: IType{Type: "arr", Subtype: IType{Type: "int", Nullable: false}, Nullable: false}, Value: []any{widthVal, heightVal}}
// 	}
// 	i.Builtins["stopTerminal"] = func(a ...any) *IValue {
// 		termbox.Close()
// 		return nil
// 	}
// 	i.Builtins["clearTerminal"] = func(a ...any) *IValue {
// 		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
// 		return nil
// 	}
// 	i.Builtins["flushTerminal"] = func(a ...any) *IValue {
// 		err := termbox.Flush()
// 		if err != nil {
// 			i.ThrowError("ERROR: Failed to flush terminal.")
// 		}
// 		return nil
// 	}
// 	i.Builtins["readFromTerminal"] = func(a ...any) *IValue {
// 		event := termbox.PollEvent()
// 		key := string(event.Ch)
// 		specialKeys := map[termbox.Key]string{
// 			termbox.KeyArrowUp:    "UP",
// 			termbox.KeyArrowDown:  "DOWN",
// 			termbox.KeyArrowLeft:  "LEFT",
// 			termbox.KeyArrowRight: "RIGHT",
// 			termbox.KeyCtrlC:      "CTRL+C",
// 		}
// 		special, ok := specialKeys[event.Key]
// 		if ok {
// 			key = special
// 		}
// 		return &IValue{Type: IType{Type: "string", Nullable: false}, Value: key}
// 	}
// 	i.Builtins["writeToTerminal"] = func(a ...any) *IValue {
// 		x := a[0].(*IValue).Value.(int)
// 		y := a[1].(*IValue).Value.(int)
// 		char := a[2].(*IValue).Value.(string)[0]
// 		colors := map[string]termbox.Attribute{
// 			"white":   termbox.ColorWhite,
// 			"black":   termbox.ColorBlack,
// 			"red":     termbox.ColorRed,
// 			"green":   termbox.ColorGreen,
// 			"blue":    termbox.ColorBlue,
// 			"yellow":  termbox.ColorYellow,
// 			"cyan":    termbox.ColorCyan,
// 			"magenta": termbox.ColorMagenta,
// 			"gray":    termbox.ColorDefault,
// 		}
// 		fg := termbox.ColorWhite
// 		bg := termbox.ColorBlack
// 		var ok bool
// 		if len(a) > 3 {
// 			colorFg := a[3].(*IValue).Value.(string)
// 			fg, ok = colors[colorFg]
// 			if !ok {
// 				i.ThrowError("ERROR: Invalid color passed to 'writeToTerminal'")
// 			}
// 		}
// 		if len(a) > 4 {
// 			colorBg := a[4].(*IValue).Value.(string)
// 			bg, ok = colors[colorBg]
// 			if !ok {
// 				i.ThrowError("ERROR: Invalid color passed to 'writeToTerminal'")
// 			}
// 		}
// 		r := rune(char)
// 		termbox.SetCell(x, y, r, fg, bg)
// 		return nil
// 	}
// 	i.Builtins["print"] = func(a ...any) *IValue {
// 		PrintValue(a[0].(*IValue), "")
// 		return nil
// 	}
// 	i.Builtins["println"] = func(a ...any) *IValue {
// 		PrintValue(a[0].(*IValue), "\n")
// 		return nil
// 	}
// 	i.Builtins["parseInt"] = func(a ...any) *IValue {
// 		value, err := strconv.Atoi(a[0].(*IValue).Value.(string))
// 		if err != nil {
// 			return &IValue{Type: IType{Type: "int", Nullable: true}, Value: nil}
// 		}
// 		return &IValue{Type: IType{Type: "int", Nullable: true}, Value: value}
// 	}
// 	i.Builtins["parseFloat"] = func(a ...any) *IValue {
// 		value, err := strconv.ParseFloat(a[0].(*IValue).Value.(string), 64)
// 		if err != nil {
// 			return &IValue{Type: IType{Type: "float", Nullable: true}, Value: nil}
// 		}
// 		return &IValue{Type: IType{Type: "float", Nullable: true}, Value: value}
// 	}
// 	i.Builtins["toString"] = func(a ...any) *IValue {
// 		value := StringValue(a[0].(*IValue))
// 		return &IValue{Type: IType{Type: "string", Nullable: false}, Value: value}
// 	}
// 	i.Builtins["getTime"] = func(a ...any) *IValue {
// 		now := time.Now()
// 		value := float64(now.UnixNano()) / 1e9
// 		return &IValue{Type: IType{Type: "float", Nullable: false}, Value: value}
// 	}
// 	i.Builtins["length"] = func(a ...any) *IValue {
// 		if a[0].(*IValue).Value != nil {
// 			if a[0].(*IValue).Type.Type == "arr" {
// 				return &IValue{Type: IType{Type: "int", Nullable: true}, Value: len(a[0].(*IValue).Value.([]any))}
// 			} else if a[0].(*IValue).Type.Type == "string" {
// 				return &IValue{Type: IType{Type: "int", Nullable: true}, Value: len(a[0].(*IValue).Value.(string))}
// 			}
// 		}
// 		return nil
// 	}
// 	i.Builtins["append"] = func(a ...any) *IValue {
// 		// todo work in subtype
// 		if len(a) < 2 {
// 			return &IValue{Type: IType{Type: "arr", Subtype: IType{Type: "any"}, Nullable: false}, Value: []any{}}
// 		}

// 		v := a[0].(*IValue)
// 		t := v.Type

// 		slice, ok := v.Value.([]any)
// 		if !ok {
// 			return &IValue{Type: t, Value: slice}
// 		}

// 		return &IValue{Type: t, Value: append(slice, a[1])}
// 	}
// 	// TODO pop and popAt
// 	i.Builtins["unicode"] = func(a ...any) *IValue {
// 		value, ok := a[0].(*IValue).Value.(int)
// 		if !ok {
// 			return &IValue{Type: IType{Type: "string", Nullable: true}, Value: nil}
// 		}
// 		return &IValue{Type: IType{Type: "string", Nullable: true}, Value: string(value)}
// 	}
// 	i.Builtins["ordinal"] = func(a ...any) *IValue {
// 		value, ok := a[0].(*IValue).Value.(string)
// 		if !ok {
// 			return &IValue{Type: IType{Type: "int", Nullable: true}, Value: nil}
// 		}
// 		return &IValue{Type: IType{Type: "int", Nullable: true}, Value: int(value[0])}
// 	}
// 	i.Builtins["randomValue"] = func(a ...any) *IValue {
// 		return &IValue{Type: IType{Type: "float", Nullable: false}, Value: rand.Float64()}
// 	}
// 	i.Builtins["mathSin"] = func(a ...any) *IValue {
// 		floatVal := a[0].(*IValue).Value.(float64)
// 		return &IValue{Type: IType{Type: "float", Nullable: false}, Value: math.Sin(floatVal)}
// 	}
// 	i.Builtins["mathCos"] = func(a ...any) *IValue {
// 		floatVal := a[0].(*IValue).Value.(float64)
// 		return &IValue{Type: IType{Type: "float", Nullable: false}, Value: math.Cos(floatVal)}
// 	}
// 	i.Builtins["mathFloor"] = func(a ...any) *IValue {
// 		floatVal1 := toFloat(a[0].(*IValue))
// 		return &IValue{Type: IType{Type: "int", Nullable: false}, Value: int(floatVal1)}
// 	}
// 	i.Builtins["mathCeil"] = func(a ...any) *IValue {
// 		floatVal1 := toFloat(a[0].(*IValue))
// 		return &IValue{Type: IType{Type: "int", Nullable: false}, Value: int(floatVal1 + 1)}
// 	}
// 	i.Builtins["mathPow"] = func(a ...any) *IValue {
// 		floatVal1 := toFloat(a[0].(*IValue))
// 		floatVal2 := toFloat(a[1].(*IValue))
// 		return &IValue{Type: IType{Type: "float", Nullable: false}, Value: math.Pow(floatVal1, floatVal2)}
// 	}
// 	i.Builtins["min"] = func(a ...any) *IValue {
// 		floatVal1 := toFloat(a[0].(*IValue))
// 		floatVal2 := toFloat(a[1].(*IValue))
// 		return &IValue{Type: IType{Type: "float", Nullable: false}, Value: math.Min(floatVal1, floatVal2)}
// 	}
// 	i.Builtins["max"] = func(a ...any) *IValue {
// 		floatVal1 := toFloat(a[0].(*IValue))
// 		floatVal2 := toFloat(a[1].(*IValue))
// 		return &IValue{Type: IType{Type: "float", Nullable: false}, Value: math.Max(floatVal1, floatVal2)}
// 	}
// 	i.Builtins["abs"] = func(a ...any) *IValue {
// 		val := a[0].(*IValue).Value.(int)
// 		if val < 0 {
// 			val *= -1
// 		}
// 		return &IValue{Type: IType{Type: "int", Nullable: false}, Value: val}
// 	}
// 	i.Builtins["floatAbs"] = func(a ...any) *IValue {
// 		val := toFloat(a[0].(*IValue))
// 		if val < 0 {
// 			val *= -1
// 		}
// 		return &IValue{Type: IType{Type: "float", Nullable: false}, Value: val}
// 	}
// }

// func (i *Interpreter) LookupVar(name string) *IValue {
// 	value, ok := i.CurrentScope[name]
// 	if ok {
// 		return &value
// 	}
// 	for j := len(i.Scopes) - 1; j > -1; j-- {
// 		scope := i.Scopes[j]
// 		value, ok := scope[name]
// 		if ok {
// 			return &value
// 		}
// 	}
// 	i.ThrowError(fmt.Sprintf("ERROR: No such variable '%s'.", name))
// 	return nil
// }

// func (i *Interpreter) MutVar(name string, index int, value *IValue) {
// 	t, ok := i.CurrentScope[name]
// 	if ok {
// 		if index != -1 {
// 			i.CurrentScope[name].Value.([]any)[index] = value
// 		} else {
// 			if TypeMismatch(t.Type, value.Type) {
// 				i.ThrowError("ERROR: Type mismatched in mutation statement.")
// 			}
// 			i.CurrentScope[name] = IValue{Type: t.Type, Value: value.Value}
// 		}
// 		return
// 	}
// 	for j := len(i.Scopes) - 1; j > -1; j-- {
// 		scope := i.Scopes[j]
// 		t, ok := scope[name]
// 		if ok {
// 			if index != -1 {
// 				if scope[name].Type.Type == "string" {
// 					original := scope[name].Value.(string)
// 					modified := original[:index] + value.Value.(string) + original[index+1:]
// 					scope[name] = IValue{Type: IType{Type: "string", Nullable: false}, Value: modified}
// 				} else {
// 					scope[name].Value.([]any)[index] = value
// 				}
// 			} else {
// 				if TypeMismatch(t.Type, value.Type) {
// 					i.ThrowError("ERROR: Type mismatched in mutation statement.")
// 				}
// 				scope[name] = IValue{Type: t.Type, Value: value.Value}
// 			}
// 			return
// 		}
// 	}
// 	i.ThrowError(fmt.Sprintf("ERROR: No such variable '%s'.", name))
// }

// func (i *Interpreter) RunInstruction(node ParseNode) *IValue {
// 	switch node.Name {
// 	case "func_decl":
// 		funcName := node.Children[0].(ParseNode).Data
// 		funcArgs := node.Children[1].(ParseNode).Children
// 		returnType := i.RunInstruction(node.Children[2].(ParseNode)).Type
// 		suite := node.Children[3].(ParseNode)

// 		params := make([]IParam, 0)

// 		for _, a := range funcArgs {
// 			arg := a.(ParseNode)
// 			if len(arg.Children) > 2 {
// 				params = append(params, IParam{Name: arg.Children[0].(ParseNode).Data, Type: i.RunInstruction(arg.Children[1].(ParseNode)).Type, Default: i.RunInstruction(arg.Children[2].(ParseNode))})
// 			} else {
// 				params = append(params, IParam{Name: arg.Children[0].(ParseNode).Data, Type: i.RunInstruction(arg.Children[1].(ParseNode)).Type})
// 			}
// 		}

// 		i.Functions[funcName] = IFunction{Params: params, ReturnType: returnType, Suite: suite}
// 		break
// 	case "struct_decl":
// 		structName := node.Children[0].(ParseNode).Data
// 		structFields := node.Children[1].(ParseNode).Children

// 		fields := make([]IParam, 0)

// 		for _, a := range structFields {
// 			field := a.(ParseNode)
// 			if len(field.Children) > 2 {
// 				fields = append(fields, IParam{Name: field.Children[0].(ParseNode).Data, Type: i.RunInstruction(field.Children[1].(ParseNode)).Type, Default: i.RunInstruction(field.Children[2].(ParseNode))})
// 			} else {
// 				fields = append(fields, IParam{Name: field.Children[0].(ParseNode).Data, Type: i.RunInstruction(field.Children[1].(ParseNode)).Type})
// 			}
// 		}

// 		i.Structs[structName] = IStruct{Fields: fields}
// 		break
// 	case "struct_instance":
// 		structName := node.Children[0].(ParseNode).Data

// 		toInstance, ok := i.Structs[structName]

// 		if !ok {
// 			i.ThrowError(fmt.Sprintf("ERROR: No such struct '%s'.", structName))
// 		}

// 		fieldValues := make(map[string]*IValue, 0)

// 		for _, field := range toInstance.Fields {
// 			if field.Default != nil {
// 				fieldValues[field.Name] = field.Default
// 			} else {
// 				fieldValues[field.Name] = field.Type.DefaultValue()
// 			}
// 		}

// 		for j, a := range node.Children {
// 			if j == 0 {
// 				continue
// 			}
// 			field := a.(ParseNode)
// 			fieldName := toInstance.Fields[j-1].Name
// 			fieldValues[fieldName] = i.RunInstruction(field)
// 		}

// 		return &IValue{Type: IType{Type: "struct"}, Value: IStructInstance{Fields: fieldValues}}
// 	case "access_field":
// 		varName := node.Children[0].(ParseNode).Data
// 		variable := i.LookupVar(varName).Value.(IStructInstance)
// 		var value *IValue
// 		var ok bool
// 		for j, child := range node.Children {
// 			if j == 0 {
// 				continue
// 			}
// 			fieldName := child.(ParseNode).Data
// 			value, ok = variable.Fields[fieldName]
// 			if !ok {
// 				i.ThrowError(fmt.Sprintf("ERROR: No such struct field '%s'.", fieldName))
// 			}
// 			variable, ok = value.Value.(IStructInstance)
// 			if ok {
// 				continue
// 			} else {
// 				break
// 			}
// 		}
// 		return value
// 	case "return_stmt": // done
// 		i.ReturnFromFunc = true
// 		i.ReturnValue = i.RunInstruction(node.Children[0].(ParseNode))
// 		break
// 	case "continue_stmt": // done
// 		i.ContinueLoop = true
// 		break
// 	case "break_stmt": // done
// 		i.BreakFromLoop = true
// 		break
// 	case "while_stmt": // done
// 		condition := node.Children[0].(ParseNode)
// 		suite := node.Children[1].(ParseNode)
// 		for IsTruthy(i.RunInstruction(condition).Value) {
// 			i.Scopes = append(i.Scopes, make(map[string]IValue))
// 			i.CurrentScope = i.Scopes[len(i.Scopes)-1]

// 			for _, child := range suite.Children {
// 				i.RunInstruction(child.(ParseNode))

// 				if i.ReturnFromFunc || i.BreakFromLoop {
// 					break
// 				}
// 				if i.ContinueLoop {
// 					i.ContinueLoop = false
// 					break
// 				}
// 			}

// 			i.Scopes = i.Scopes[:len(i.Scopes)-1]
// 			i.CurrentScope = i.Scopes[len(i.Scopes)-1]

// 			if i.ReturnFromFunc || i.BreakFromLoop {
// 				i.BreakFromLoop = false
// 				break
// 			}
// 		}
// 		break
// 	case "loop_stmt":
// 		values := i.RunInstruction(node.Children[0].(ParseNode))
// 		varName := node.Children[1].(ParseNode).Data
// 		varType := i.RunInstruction(node.Children[2].(ParseNode)).Type
// 		suite := node.Children[3].(ParseNode)

// 		if values.Type.Type == "arr" {
// 			arr := values.Value.([]any)
// 			for _, item := range arr {
// 				i.Scopes = append(i.Scopes, make(map[string]IValue))
// 				i.CurrentScope = i.Scopes[len(i.Scopes)-1]

// 				i.CurrentScope[varName] = IValue{Type: varType, Value: item.(*IValue).Value}

// 				for _, child := range suite.Children {
// 					i.RunInstruction(child.(ParseNode))

// 					if i.ReturnFromFunc || i.BreakFromLoop {
// 						break
// 					}
// 					if i.ContinueLoop {
// 						i.ContinueLoop = false
// 						break
// 					}
// 				}

// 				i.Scopes = i.Scopes[:len(i.Scopes)-1]
// 				i.CurrentScope = i.Scopes[len(i.Scopes)-1]

// 				if i.ReturnFromFunc || i.BreakFromLoop {
// 					i.BreakFromLoop = false
// 					break
// 				}
// 			}
// 		} else {
// 			i.ThrowError("ERROR: Attempt to loop over non-array.")
// 		}
// 		break
// 	case "if_stmt": // done
// 		suite := node.Children[1].(ParseNode)
// 		elseSuite := node.Children[2].(ParseNode)
// 		if IsTruthy(i.RunInstruction(node.Children[0].(ParseNode)).Value) {
// 			i.Scopes = append(i.Scopes, make(map[string]IValue))
// 			i.CurrentScope = i.Scopes[len(i.Scopes)-1]

// 			for _, child := range suite.Children {
// 				i.RunInstruction(child.(ParseNode))

// 				if i.ReturnFromFunc || i.BreakFromLoop || i.ContinueLoop {
// 					break
// 				}
// 			}

// 			i.Scopes = i.Scopes[:len(i.Scopes)-1]
// 			i.CurrentScope = i.Scopes[len(i.Scopes)-1]
// 		} else {
// 			i.Scopes = append(i.Scopes, make(map[string]IValue))
// 			i.CurrentScope = i.Scopes[len(i.Scopes)-1]

// 			for _, child := range elseSuite.Children {
// 				i.RunInstruction(child.(ParseNode))

// 				if i.ReturnFromFunc || i.BreakFromLoop || i.ContinueLoop {
// 					break
// 				}
// 			}

// 			i.Scopes = i.Scopes[:len(i.Scopes)-1]
// 			i.CurrentScope = i.Scopes[len(i.Scopes)-1]
// 		}
// 		break
// 	case "variable_decl": // done
// 		varName := node.Children[0].(ParseNode).Data
// 		varType := i.RunInstruction(node.Children[1].(ParseNode))
// 		varValue := i.RunInstruction(node.Children[2].(ParseNode))
// 		_, ok := i.CurrentScope[varName]
// 		if ok {
// 			i.ThrowError(fmt.Sprintf("ERROR: Attempt to multiple-initialize variable '%s'.", varName))
// 		}
// 		if !varType.Type.Nullable && varValue.Value == nil {
// 			i.ThrowError("ERROR: Attempt to use null-value in non-nullable type.")
// 		}
// 		varValue.Type.Nullable = varType.Type.Nullable
// 		if TypeMismatch(varType.Type, varValue.Type) {
// 			fmt.Println(varType.Type)
// 			fmt.Println(varValue.Type)
// 			i.ThrowError("ERROR: Type mismatch in variable initialization.")
// 		}
// 		i.CurrentScope[varName] = IValue{Type: varType.Type, Value: varValue.Value}
// 		break
// 	case "mut_stmt": // done
// 		varName := node.Children[0].(ParseNode)
// 		var varIndex int
// 		var name string
// 		if varName.Name == "get_item" {
// 			value := i.RunInstruction(varName.Children[0].(ParseNode))
// 			index := i.RunInstruction(varName.Children[1].(ParseNode))
// 			t := value.Type
// 			if t.Type == "arr" {
// 				varIndex = index.Value.(int)
// 				varValue := i.RunInstruction(node.Children[1].(ParseNode))
// 				value.Value.([]any)[varIndex] = varValue
// 			} else if t.Type == "map" {
// 				varValue := i.RunInstruction(node.Children[1].(ParseNode))
// 				if index.Type.Type == t.Subtype.(IType).Type && varValue.Type.Type == t.Subtype.(IType).Subtype.(IType).Type {
// 					value.Value.(map[any]any)[index.Value] = varValue
// 				} else {
// 					i.ThrowError("ERROR: Wrong type for map mutation in either key or value.")
// 				}
// 			} else {
// 				i.ThrowError("ERROR: Couldn't mutate item.")
// 			}
// 		} else if varName.Name == "access_field" {
// 			structName := varName.Children[0].(ParseNode).Data

// 			variable := i.LookupVar(structName).Value.(IStructInstance)
// 			var value *IValue
// 			var fieldName string
// 			var ok bool
// 			for j, child := range varName.Children {
// 				if j == 0 {
// 					continue
// 				}
// 				fieldName = child.(ParseNode).Data
// 				value, ok = variable.Fields[fieldName]
// 				if !ok {
// 					i.ThrowError(fmt.Sprintf("ERROR: No such struct field '%s'.", fieldName))
// 				}
// 				oldVar := variable
// 				variable, ok = value.Value.(IStructInstance)
// 				if ok {
// 					continue
// 				} else {
// 					varValue := i.RunInstruction(node.Children[1].(ParseNode))
// 					oldVar.Fields[fieldName] = varValue
// 					break
// 				}
// 			}
// 		} else if varName.Name == "name" {
// 			name = varName.Data
// 			varIndex = -1
// 			varValue := i.RunInstruction(node.Children[1].(ParseNode))
// 			i.MutVar(name, varIndex, varValue)
// 		}
// 		break
// 	case "get_item":
// 		variable := i.RunInstruction(node.Children[0].(ParseNode))
// 		for j, child := range node.Children {
// 			if j == 0 {
// 				continue
// 			}
// 			expr := i.RunInstruction(child.(ParseNode))
// 			var array []any
// 			if variable.Type.Type == "arr" {
// 				array = variable.Value.([]any)
// 				if expr.Type.Type != "int" {
// 					i.ThrowError("ERROR: Usage of non-int type to index into array.")
// 				}
// 			} else if variable.Type.Type == "string" {
// 				value, ok := variable.Value.(string)
// 				if ok {
// 					array = make([]any, len(value))
// 					for k, element := range value {
// 						array[k] = &IValue{Type: IType{Type: "string", Nullable: false}, Value: string(element)}
// 					}
// 				}
// 				if expr.Type.Type != "int" {
// 					i.ThrowError("ERROR: Usage of non-int type to index into array.")
// 				}
// 			} else if variable.Type.Type == "map" {
// 				if expr.Type.Type != variable.Type.Subtype.(IType).Type {
// 					i.ThrowError("ERROR: Incorrect type for map lookup.")
// 				}
// 				variable, ok := variable.Value.(map[any]any)[expr.Value].(*IValue)
// 				if ok {
// 					return variable
// 				}
// 				return &IValue{Type: IType{Type: "any", Nullable: true}, Value: nil}
// 			} else {
// 				if variable.Type.Type == "any" {
// 					value, ok := variable.Value.([]any)
// 					if ok {
// 						array = value
// 						if expr.Type.Type != "int" {
// 							i.ThrowError("ERROR: Usage of non-int type to index into array.")
// 						}
// 					}
// 				} else {
// 					fmt.Println(variable.Type.Type)
// 					i.ThrowError("ERROR: Attempt to get item from non-iterable.")
// 				}
// 			}
// 			intValue := expr.Value.(int)
// 			if intValue < 0 || intValue >= len(array) {
// 				fmt.Println(node.Children[0].(ParseNode).Data)
// 				i.ThrowError(fmt.Sprintf("ERROR: Index %d is out of bounds.", intValue))
// 			}
// 			variable = array[intValue].(*IValue)
// 		}
// 		return variable
// 	case "string_literal": // done
// 		return &IValue{Type: IType{Type: "string", Nullable: false}, Value: node.Data}
// 	case "int_literal": // done
// 		value, err := strconv.Atoi(node.Data)
// 		if err != nil {
// 			return nil
// 		}
// 		return &IValue{Type: IType{Type: "int", Nullable: false}, Value: value}
// 	case "float_literal": // done
// 		value, err := strconv.ParseFloat(node.Data, 64)
// 		if err != nil {
// 			return nil
// 		}
// 		return &IValue{Type: IType{Type: "float", Nullable: false}, Value: value}
// 	case "bool_literal", "nullable": // done
// 		value, err := strconv.ParseBool(node.Data)
// 		if err != nil {
// 			return nil
// 		}
// 		return &IValue{Type: IType{Type: "bool", Nullable: false}, Value: value}
// 	case "type": // done
// 		typeNode := &IValue{}

// 		typeNode.Type = i.RunInstruction(node.Children[1].(ParseNode)).Type

// 		if len(node.Children) > 2 {
// 			typeNode.Type.Subtype = i.RunInstruction(node.Children[2].(ParseNode)).Type
// 		}

// 		if len(node.Children) > 3 {
// 			typeNode.Type.Subtype = IType{Type: typeNode.Type.Subtype.(IType).Type, Subtype: i.RunInstruction(node.Children[3].(ParseNode)).Type}
// 		}

// 		typeNode.Type.Nullable = i.RunInstruction(node.Children[0].(ParseNode)).Value.(bool)

// 		return typeNode
// 	case "array":
// 		elements := make([]any, 0)
// 		for _, child := range node.Children {
// 			elements = append(elements, i.RunInstruction(child.(ParseNode)))
// 		}
// 		// TODO make sure all elements are the same type
// 		if len(elements) != 0 {
// 			return &IValue{Type: IType{Type: "arr", Subtype: elements[0].(*IValue).Type, Nullable: false}, Value: elements}
// 		} else {
// 			return &IValue{Type: IType{Type: "arr", Subtype: IType{Type: "any"}, Nullable: false}, Value: elements}
// 		}
// 	case "map":
// 		elements := make(map[any]any, 0)
// 		for _, child := range node.Children {
// 			childKey := i.RunInstruction(child.(ParseNode).Children[0].(ParseNode))
// 			childValue := i.RunInstruction(child.(ParseNode).Children[1].(ParseNode))

// 			if !(childKey.Type.Type == "int" || childKey.Type.Type == "string" || childKey.Type.Type == "bool") {
// 				i.ThrowError("ERROR: Unaccepted key type for map.")
// 			}

// 			elements[childKey] = childValue
// 		}
// 		// TODO make sure all elements are the same type
// 		if len(elements) != 0 {
// 			var kType IType
// 			var vType IType
// 			newElements := make(map[any]any, 0)
// 			for k, v := range elements {
// 				kType = k.(*IValue).Type
// 				vType = v.(*IValue).Type
// 				newElements[k.(*IValue).Value] = v
// 			}
// 			return &IValue{Type: IType{Type: "map", Subtype: IType{Type: kType.Type, Subtype: vType}, Nullable: false}, Value: newElements}
// 		} else {
// 			return &IValue{Type: IType{Type: "map", Subtype: IType{Type: "any", Subtype: IType{Type: "any"}}, Nullable: false}, Value: elements}
// 		}
// 	case "func_call": // done
// 		funcName := node.Children[0].(ParseNode).Data
// 		arguments := make([]*IValue, 0)
// 		for j := 1; j < len(node.Children); j++ {
// 			arguments = append(arguments, i.RunInstruction(node.Children[j].(ParseNode)))
// 		}
// 		function, ok := i.Functions[funcName]
// 		if ok {
// 			i.Scopes = append(i.Scopes, make(map[string]IValue))
// 			i.CurrentScope = i.Scopes[len(i.Scopes)-1]

// 			i.ReturnFromFunc = false
// 			i.ReturnValue = nil

// 			for j, param := range function.Params {
// 				varName := param.Name
// 				varType := param.Type
// 				var varValue *IValue
// 				if j < len(arguments) {
// 					varValue = arguments[j]
// 				} else {
// 					if param.Default != nil {
// 						varValue = param.Default
// 					} else {
// 						i.ThrowError("ERROR: Wrong number of arguments passed to function.")
// 					}
// 				}
// 				if !varType.Nullable && varValue.Value == nil {
// 					i.ThrowError("ERROR: Null passed for non-nullable argument.")
// 				}
// 				varValue.Type.Nullable = varType.Nullable
// 				if TypeMismatch(varType, varValue.Type) {
// 					i.ThrowError("ERROR: Type mismatch in function call.")
// 				}
// 				i.CurrentScope[varName] = IValue{Type: varType, Value: varValue.Value}
// 			}

// 			for _, child := range function.Suite.Children {
// 				i.RunInstruction(child.(ParseNode))

// 				if i.ReturnFromFunc {
// 					break
// 				}
// 			}
// 			i.ReturnFromFunc = false
// 			i.Scopes = i.Scopes[:len(i.Scopes)-1]
// 			i.CurrentScope = i.Scopes[len(i.Scopes)-1]

// 			if i.ReturnValue == nil {
// 				i.ReturnValue = &IValue{Type: IType{Type: "any", Nullable: true}, Value: nil}
// 			}

// 			if function.ReturnType.Nullable && !i.ReturnValue.Type.Nullable {
// 				i.ReturnValue.Type.Nullable = true
// 			}

// 			if TypeMismatch(i.ReturnValue.Type, function.ReturnType) {
// 				i.ThrowError(fmt.Sprintf("ERROR: Type mismatch between return type and value returned. Got %s, expected %s.", i.ReturnValue.Type.Type, function.ReturnType.Type))
// 			}

// 			return i.ReturnValue
// 		} else {
// 			function, ok := i.Builtins[funcName]

// 			if ok {
// 				argumentValues := make([]any, 0)
// 				for _, arg := range arguments {
// 					argumentValues = append(argumentValues, arg)
// 				}
// 				i.ReturnValue = function(argumentValues...)

// 				if i.ReturnValue == nil {
// 					i.ReturnValue = &IValue{Type: IType{Type: "any", Nullable: true}, Value: nil}
// 				}

// 				return i.ReturnValue
// 			} else {
// 				i.ThrowError(fmt.Sprintf("ERROR: No such function: '%s'.", funcName))
// 			}
// 		}
// 		break
// 	case "name": // done
// 		if slices.Contains(i.BasicTypes, node.Data) {
// 			return &IValue{Type: IType{Type: node.Data}}
// 		} else if node.Data == "null" {
// 			return &IValue{Type: IType{Type: "any", Nullable: true}, Value: nil}
// 		}
// 		return i.LookupVar(node.Data)
// 	case "muldiv_operation": // done
// 		a := i.RunInstruction(node.Children[0].(ParseNode))
// 		op := node.Children[1].(ParseNode).Data
// 		b := i.RunInstruction(node.Children[2].(ParseNode))
// 		if op == "star" {
// 			t, v := MultiplyValues(*a, *b)
// 			return &IValue{Type: t, Value: v}
// 		}
// 		if op == "slash" {
// 			t, v := DivideValues(*a, *b)
// 			return &IValue{Type: t, Value: v}
// 		}
// 		t, v := ModuloValues(*a, *b)
// 		return &IValue{Type: t, Value: v}
// 	case "addsub_operation": // done
// 		a := i.RunInstruction(node.Children[0].(ParseNode))
// 		op := node.Children[1].(ParseNode).Data
// 		b := i.RunInstruction(node.Children[2].(ParseNode))
// 		if op == "plus" {
// 			t, v := AddValues(*a, *b)
// 			return &IValue{Type: t, Value: v}
// 		}
// 		t, v := SubtractValues(*a, *b)
// 		return &IValue{Type: t, Value: v}
// 	case "unary_negate_operation": // done
// 		a := i.RunInstruction(node.Children[0].(ParseNode))
// 		t, v := SubtractValues(IValue{Type: IType{Type: "int", Nullable: false}, Value: 0}, *a)
// 		return &IValue{Type: t, Value: v}
// 	case "or_operation": // done
// 		a := i.RunInstruction(node.Children[0].(ParseNode)).Value.(bool)
// 		b := i.RunInstruction(node.Children[1].(ParseNode)).Value.(bool)
// 		return &IValue{Type: IType{Type: "bool"}, Value: a || b}
// 	case "and_operation": // done
// 		a := i.RunInstruction(node.Children[0].(ParseNode)).Value.(bool)
// 		b := i.RunInstruction(node.Children[1].(ParseNode)).Value.(bool)
// 		return &IValue{Type: IType{Type: "bool"}, Value: a && b}
// 	case "not_operation": // done
// 		a := i.RunInstruction(node.Children[0].(ParseNode)).Value.(bool)
// 		return &IValue{Type: IType{Type: "bool"}, Value: !a}
// 	case "comparison_operation": // done
// 		a := i.RunInstruction(node.Children[0].(ParseNode))
// 		op := node.Children[1].(ParseNode).Data
// 		b := i.RunInstruction(node.Children[2].(ParseNode))
// 		v := CompareValues(a, op, b)
// 		return &IValue{Type: IType{Type: "bool"}, Value: v}
// 	case "equality_operation": // done
// 		a := i.RunInstruction(node.Children[0].(ParseNode))
// 		op := node.Children[1].(ParseNode).Data
// 		b := i.RunInstruction(node.Children[2].(ParseNode))
// 		v := CompareEquality(a, op, b)
// 		return &IValue{Type: IType{Type: "bool"}, Value: v}
// 	}
// 	return &IValue{Type: IType{Type: "any", Nullable: true}, Value: nil}
// }

// func (i *Interpreter) Start() {
// 	i.Functions = make(map[string]IFunction)
// 	i.Structs = make(map[string]IStruct)

// 	i.Scopes = make([]map[string]IValue, 0)
// 	i.CurrentScope = make(map[string]IValue)
// 	i.Scopes = append(i.Scopes, i.CurrentScope)

// 	i.BasicTypes = []string{
// 		"any",
// 		"int",
// 		"float",
// 		"bool",
// 		"string",
// 		"arr",
// 		"map",
// 	}

// 	i.HasOpenedTerminal = false

// 	i.RunTree()
// }

// func (i *Interpreter) RunTree() {
// 	for _, node := range i.Tree.Children {
// 		i.RunInstruction(node.(ParseNode))
// 	}
// 	if i.HasOpenedTerminal {
// 		termbox.Close()
// 	}
// }