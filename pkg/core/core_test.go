package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoardPrint(t *testing.T) {
	fmt.Println(NewBoard().PrettyPrint())
}

func TestBoardFull(t *testing.T) {
	assert := assert.New(t)

	board := NewBoard()
	res := board.IsFull()
	assert.False(res)

	board.SetPosition('X', '1')
	board.SetPosition('X', '2')
	board.SetPosition('X', '3')
	board.SetPosition('X', '4')
	board.SetPosition('X', '5')

	res = board.IsFull()
	assert.False(res)

	board.SetPosition('X', '6')
	board.SetPosition('X', '7')
	board.SetPosition('X', '8')
	board.SetPosition('X', '9')

	res = board.IsFull()
	assert.True(res)

	fmt.Println(board.PrettyPrint())
}

func TestBoardWinner(t *testing.T) {
	assert := assert.New(t)
	board := NewBoard()
	res := board.IsWinner('X')
	assert.False(res)

	board.SetPosition('X', '1')
	board.SetPosition('X', '2')
	board.SetPosition('X', '3')

	res = board.IsWinner('X')
	assert.True(res)
}

func TestGameMove(t *testing.T) {
	assert := assert.New(t)
	p1 := Player{ID: "p1"}
	p2 := Player{ID: "p2"}

	game := NewGame(&p1, &p2)
	err := game.AliceMove('1')
	assert.Nil(err)

	err = game.AliceMove('9')
	assert.NotNil(err)

	err = game.BobMove('1')
	assert.NotNil(err)

	err = game.BobMove('9')
	assert.Nil(err)

}

func TestGameWinner(t *testing.T) {
	assert := assert.New(t)
	p1 := Player{ID: "p1"}
	p2 := Player{ID: "p2"}

	game := NewGame(&p1, &p2)

	game.GameBoard.SetPosition('X', '1')
	game.GameBoard.SetPosition('X', '2')
	game.GameBoard.SetPosition('X', '3')

	winner := game.GetWinner()
	assert.Equal(&p1, winner)
}

func TestCoreMove(t *testing.T) {
	t.Skip()
}
