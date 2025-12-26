package main

import (
	"fmt"
	"strconv"
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
	}
	// case "arr":
	// 	output := "[ "
	// 	for _, element := range v.Value.([]any) {
	// 		output += StringValue(element.(*IValue)) + " "
	// 	}
	// 	output += "]"
	// 	return output
	// }
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

func (r *IRInterpreter) Run() {
	r.Scopes = make([]map[string]CValue, 0)
	r.CurrentScope = make(map[string]CValue)
	r.Scopes = append(r.Scopes, r.CurrentScope)

	r.Stack = make([]CValue, 0)
	r.CallStack = make([]int, 0)

	r.Labels = make(map[string]int)

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
		case OP_PUSH:
			value := r.LookupVar(instruction.A.StringValue)
			r.Push(value)
		case OP_STORE:
			value := r.Pop()
			r.CurrentScope[instruction.A.StringValue] = value
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
			value := r.Pop()
			PrintValue(value, "\n")
		}
		r.ProgramCounter++
	}
}