package main

import (
	"io"
	"os"
)

type Interpreter struct {
	iPointer int
	dPointer int
}

var (
	instructions []byte
	memory       [1024]byte
	err          error
	interpreter  Interpreter
)

func (interpreter *Interpreter) next() {
	readBuf := []byte{0}
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
		_, err = os.Stdout.Write([]byte{memory[interpreter.dPointer]})
	case ',':
		_, err = os.Stdin.Read(readBuf)
		if err == io.EOF {
			os.Exit(0)
		}
		memory[interpreter.dPointer] = readBuf[0]
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

func (interpreter *Interpreter) interpret() {
	for interpreter.iPointer < len(instructions) {
		interpreter.next()
	}
}

func main() {
	instructions = []byte(os.Args[1])
	interpreter.interpret()
}