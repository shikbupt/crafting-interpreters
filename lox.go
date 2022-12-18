package main

import (
	"bufio"
	"craftinginterpreters/lox/interpreter"
	"craftinginterpreters/lox/parser"
	"craftinginterpreters/lox/scanner"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	promptFlag := flag.Bool("p", false, "process model")
	flag.Parse()

	if len(flag.Args()) == 0 {
		panic("need input lox file")
	}
	lox := newLox()

	if *promptFlag {
		lox.RunPrompt()
	}
	lox.RunFile(flag.Arg(0))
}

type Lox struct {
	interpreter     *interpreter.Interpreter
	hadError        bool
	hadRuntimeError bool
}

func newLox() *Lox {
	return &Lox{
		interpreter: interpreter.NewInterpreter(),
	}
}

func (l *Lox) RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf(">> ")
		if !scanner.Scan() {
			break
		}
		l.run(scanner.Text())
		l.hadError = false
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func (l *Lox) RunFile(name string) {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	loxContext, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	l.run(string(loxContext))
	if l.hadError {
		panic("lox has error")
	}
}

func (l *Lox) run(loxContext string) {
	scanner := scanner.New()
	tokens := scanner.ScanAll(loxContext)

	parsers := parser.NewParser(tokens)
	expression, err := parsers.Parse()
	if err != nil {
		l.Error(err)
		return
	}

	fmt.Println(parser.AstPrinter{}.Print(expression))
	err = l.interpreter.Interpret(expression)
	if err != nil {
		l.runtimeError(err)
		return
	}

}

func (l *Lox) Error(err error) {
	l.hadError = true
	fmt.Println(err)
}

func (l *Lox) runtimeError(err error) {
	l.hadRuntimeError = true
	fmt.Println(err)
}

func (l *Lox) report(line int, where, message string) {
	fmt.Printf("[line %d ] Error %s: %s", line, where, message)
	l.hadError = true
}
