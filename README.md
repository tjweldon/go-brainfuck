# Readme

Go-brainfuck is a suuuper simple implementation of the [brainfuck](https://en.wikipedia.org/wiki/Brainfuck) interpreter in golang.

## Usage

Build using

```shell
./build.sh
```

To execute some brainfuck:

```shell
./go-brainfuck '++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.'
Hello World!
```

Brainfuck can read from stdin, this is implemened using stdin redirect (`<`)
```shell
program='-,+[-[>>++++[>++++++++<-]<+<-[>+>+>-[>>>]<[[>+<-]>>+>]<<<<<-]]>>>[-]+>--[-[<->+++[-]]]<[++++++++++++<[>-[>+>>]>[+[<+>-]>+>>]<<<<<-]>>[<+>-]>[-[-<<[-]>>]<<[<<->>-]>>]<<[<<+>>-]]<[-]<.[-]<-,+]'

echo hello > hello.txt
./go-brainfuck $program < hello.txt
uryyb

echo uryyb > uryyb.txt
./go-brainfuck $program < uryyb.txt
hello
```
