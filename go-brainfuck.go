package main

import (
	"bytes"
	"fmt"
	"github.com/alexflint/go-arg"
	tm "github.com/buger/goterm"
	"log"
	"os"
	"strings"
	"time"
)

type Interpreter struct {
	iPointer int
	dPointer int
	output   []byte
	input    []byte
}

var (
	instructions []byte
	memory       [1024]byte
	err          error
	w            int
)

var args struct {
	File         string `arg:"-f, --file" help:"The file containing the code"`
	Instructions string `arg:"-i, --instructions" help:"The brainfuck code"`
	Tick         int    `arg:"-t, --tick" help:"Execute in slow motion, time per instruction in milliseconds" default:"0"`
}

type InterpreterHook func(i Interpreter)

var interpreter = Interpreter{}

func main() {
	arg.MustParse(&args)
	instructions = []byte(args.Instructions)

	if len(instructions) == 0 {
		instructions, err = os.ReadFile(args.File)
	}
	if bytes.Contains(instructions, []byte(",")) {
		buff := make([]byte, 32)
		_, err := os.Stdin.Read(buff[:32])
		if err != nil {
			log.Fatalln(err)
		}
		interpreter.input = buff
	}

	interpreter.interpret(
		[]InterpreterHook{
			heartbeat(args.Tick),
			dump,
		},
	)
}

func (interpreter *Interpreter) next() {
	var inByte byte
	instruction := instructions[interpreter.iPointer]
	switch instruction {
	case '>':
		interpreter.dPointer++
	case '<':
		interpreter.dPointer--
	case '+':
		memory[interpreter.dPointer]++
	case '-':
		memory[interpreter.dPointer]--
	case '.':
		interpreter.output = append(interpreter.output, memory[interpreter.dPointer])
	case ',':
		if len(interpreter.input[:]) == 0 {
			os.Exit(0)
		}
		inByte, interpreter.input = interpreter.input[0], interpreter.input[1:]
		memory[interpreter.dPointer] = inByte
	case '[':
		if memory[interpreter.dPointer] == 0 {
			depth := 0
			for offset, instruction := range instructions[interpreter.iPointer:] {
				switch instruction {
				case '[':
					depth++
				case ']':
					depth--
				}

				if depth == 0 {
					interpreter.iPointer += offset
					return
				}
			}
		}
	case ']':
		if memory[interpreter.dPointer] != 0 {
			depth := 0
			for offset := 0; offset > -interpreter.iPointer; offset-- {
				switch instructions[interpreter.iPointer+offset] {
				case '[':
					depth++
				case ']':
					depth--
				}

				if depth == 0 {
					interpreter.iPointer += offset
					return
				}
			}
		}
	}

	if err != nil {
		panic(err)
	}

	interpreter.iPointer++
}

func (interpreter *Interpreter) interpret(hooks []InterpreterHook) {
	for interpreter.iPointer < len(instructions) {
		for _, hook := range hooks {
			hook(*interpreter)
		}
		interpreter.next()
	}
}

func (interpreter *Interpreter) String() string {
	w = tm.Width()

	padding := []byte(strings.Repeat(" ", max(len(instructions)-interpreter.iPointer-1, 0)))
	ticker := append(instructions, padding...)
	ticker2 := make([]byte, w)
	copy(ticker2, ticker[interpreter.iPointer:len(instructions)-1])
	nothings := strings.Repeat(" ", w)
	s2 := fmt.Sprintf(
		"%s\r%s\n%s\r%s\n% x\n%s\n%s\n",
		strings.Repeat(" ", w),
		string(interpreter.input),
		nothings,
		ticker2,
		memory[:32],
		renderPointer([]byte("___"), []byte("_⇡_"), 32),
		interpreter.output,
	)
	return s2
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Sgr string

const (
	PtrStyle Sgr = "\x1b[38;5;12m"
	BgStyle      = "\x1b[1m\x1b[38;5;8m"
	Default      = "\x1b[0m"
)

func (s Sgr) String() string {
	return string(s)
}

func renderPointer(bgEncoded []byte, ptrEncoded []byte, pRange int) string {
	increment := len(bgEncoded)

	styledEncodedPtr := fmt.Sprintf(
		"%s%s%s",
		PtrStyle,
		string(ptrEncoded),
		Default+BgStyle,
	)

	renderedBg := bytes.Repeat(bgEncoded, pRange-1)
	ptrStrOffset := increment * interpreter.dPointer

	return fmt.Sprintf(
		"%s%s%s%s%s",
		BgStyle,
		renderedBg[:ptrStrOffset],
		styledEncodedPtr,
		renderedBg[ptrStrOffset:],
		Default,
	)
}

func nothing(_ Interpreter) {
	return
}

func heartbeat(tick int) InterpreterHook {
	if tick == 0 {
		return nothing
	}

	return func(_ Interpreter) {
		time.Sleep(time.Duration(tick) * time.Millisecond)
	}
}

func dump(i Interpreter) {
	fmt.Print(tm.MoveTo(i.String(), 0, 10))
}
