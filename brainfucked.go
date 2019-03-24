package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

// Commands in the brainfuck language
const (
	OpIncPointer = 43 // +	Increment the memory cell under the pointer
	OpDecPinter  = 45 // -	Decrement the memory cell under the pointer
	OpMoveLeft   = 60 // <	Move the pointer to the left
	OpMoveRight  = 62 // >	Move the pointer to the right
	OpOutput     = 46 // .	Output the character signified by the cell at the pointer
	OpInput      = 44 // ,	Input a character and store it in the cell at the pointer
	OpJumpFwd    = 91 // [	Jump past the matching ] if the cell under the pointer is 0
	OpJumpBack   = 93 // ]	Jump back to the matching [ if the cell under the pointer is nonzero
	OpComment    = 35 // #  Ignores the rest of the line until newline
)

// Tape is a representation of the virtual machine.
type Tape struct {
	data  []byte
	cells []int16
	ptr   int16
}

// NewInterpreter creates a new instance of a tape
func NewInterpreter(data []byte) (*Tape, error) {
	if len(data) <= 0 {
		return &Tape{}, errors.New("No input")
	}
	return &Tape{
		data:  data,
		cells: make([]int16, 65535),
	}, nil
}

// Interpret brainfuck code
func (tape *Tape) Interpret() {
	// Tape is currently hardcoded to a size of 65535.
	// the pointer to
	ptr := &tape.ptr

	for i := 0; i < len(tape.data); i++ {
		switch tape.data[i] {
		case OpIncPointer: // +
			tape.cells[*ptr]++

		case OpDecPinter: // -
			tape.cells[*ptr]--

		case OpMoveLeft: // <
			*ptr--

		case OpMoveRight: // >
			*ptr++

		case OpOutput: // .
			fmt.Print(string(tape.cells[*ptr]))

		case OpInput: // ,
			var character int16
			_, err := fmt.Scanf("%c", &character)
			if err != nil {
				log.Fatal(err.Error())
			}
			tape.cells[*ptr] = character

		case OpJumpFwd: // [
			for j := 1; j > 0 && tape.cells[*ptr] == 0; {
				i++
				if tape.data[i] == 91 {
					j++
				} else if tape.data[i] == 93 {
					j--
				}
			}

		case OpJumpBack: // ]
			for j := 1; j > 0; {
				i--
				if tape.data[i] == 91 {
					j--
				} else if tape.data[i] == 93 {
					j++
				}
			}
			i--

		case OpComment: // #
			// if the buffer encounters a hashtag, increment i until a newline is found.
			for i < len(tape.data) && tape.data[i] != 10 {
				i++
			}
		}
	}
}

// RunInterpreter creates a new interpreter and runs the code
func RunInterpreter(buf []byte) {
	tape, err := NewInterpreter(buf)
	if err != nil {
		log.Fatal(err.Error())
	}

	tape.Interpret()
}

func main() {
	var file, instructions string

	// Command line flags, use the interpreter in two ways either by providing a file, or an instruction set.
	flag.StringVar(&file, "file", "", "path to brainfuck file")
	flag.StringVar(&instructions, "instructions", "", "string of instructions to interpret")
	flag.Parse()

	if file != "" {
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err.Error())
		}

		RunInterpreter(buf)
	} else if instructions != "" {
		RunInterpreter([]byte(instructions))
	}
}
