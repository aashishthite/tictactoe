package core

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
)

var (
	defaultBoard = []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9'}
)

type Board []rune

func NewBoard() Board {
	newBoard := make([]rune, 9)
	copy(newBoard, defaultBoard)
	return newBoard
}

func (b Board) PrettyPrint() string {

	buffer := bytes.NewBufferString("")
	buffer.WriteString(fmt.Sprintf("%s | %s | %s \n", string(b[0]), string(b[1]), string(b[2])))
	buffer.WriteString("__  __  __ \n")
	buffer.WriteString(fmt.Sprintf("%s | %s | %s \n", string(b[3]), string(b[4]), string(b[5])))
	buffer.WriteString("__  __  __ \n")
	buffer.WriteString(fmt.Sprintf("%s | %s | %s \n", string(b[6]), string(b[7]), string(b[8])))

	return buffer.String()
}

func (b Board) IsWinner(symb rune) bool {

	if symb != 'X' && symb != 'O' {
		log.Printf("Ivalid Symbol: %s", string(symb))
		return false
	}

	res := (b[0] == symb && b[1] == symb && b[2] == symb) ||
		(b[3] == symb && b[4] == symb && b[5] == symb) ||
		(b[6] == symb && b[7] == symb && b[8] == symb) ||
		(b[0] == symb && b[3] == symb && b[6] == symb) ||
		(b[1] == symb && b[4] == symb && b[7] == symb) ||
		(b[2] == symb && b[5] == symb && b[8] == symb) ||
		(b[0] == symb && b[4] == symb && b[8] == symb) ||
		(b[2] == symb && b[4] == symb && b[6] == symb)
	return res
}

func (b Board) IsFull() bool {
	for k, v := range b {
		if v == defaultBoard[k] {
			return false
		}
	}
	return true
}

func (b Board) SetPosition(symb, pos rune) error {

	if symb != 'X' && symb != 'O' {
		return fmt.Errorf("Ivalid Symbol: %s", string(symb))
	}

	n, _ := strconv.Atoi(string(pos))
	if n < 1 || n > 9 {
		return fmt.Errorf("Invalid position: %s", string(pos))
	}

	if b[n-1] != defaultBoard[n-1] {
		return fmt.Errorf("Position already used: %s", string(pos))
	}

	b[n-1] = symb

	return nil
}
