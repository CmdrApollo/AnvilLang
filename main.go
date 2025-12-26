package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	arguments := os.Args[1:]

	if len(arguments) > 0 {
		fileName := arguments[0]

		splitUp := strings.Split(fileName, "/")
		pathForParser := strings.Join(splitUp[:len(splitUp)-1], "/")

		data, err := os.ReadFile(fileName)
		if err == nil {
			stringData := string(data)

			lexer := Lexer{Stream: stringData}
			lexer.Lex()

			parser := Parser{Tokens: lexer.Tokens, ToPath: "./" + pathForParser}
			parser.ParseAll()

			// compiler := Compiler{Tree: parser.TreeNode}
			// compiler.MakeBuiltins()
			// compiler.Start()
			// compiler.PrintHiLevelOutput()

			// irInterpreter := IRInterpreter{Stream: compiler.FinalOutput}
			// irInterpreter.Run()

			interpreter := Interpreter{Tree: parser.TreeNode}
			interpreter.MakeBuiltins()
			interpreter.Start()
		}
	} else {
		interpreter := Interpreter{}
		interpreter.MakeBuiltins()
		interpreter.Start()

		for {
			fmt.Print(">>> ")
			data, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				return
			}

			lexer := Lexer{Stream: data}
			lexer.Lex()

			parser := Parser{Tokens: lexer.Tokens, ToPath: "."}
			parser.ParseAll()

			interpreter.Tree = parser.TreeNode
			interpreter.RunTree()
		}
	}
}
