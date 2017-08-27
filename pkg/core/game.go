package core

import (
	"fmt"
	"time"
)

type Player struct {
	ID string
}

type Game struct {
	Alice *Player //always 'X'
	Bob   *Player //always 'O'

	Move *Player

	Timeout time.Duration

	GameBoard Board
}

func NewGame(p1, p2 *Player) *Game {
	return &Game{
		Alice:     p1,
		Bob:       p2,
		Move:      p1,
		GameBoard: NewBoard(),
	}
}

func (g *Game) AliceMove(pos rune) error {
	if g.Move != g.Alice {
		return fmt.Errorf("Player 2's turn")
	}
	err := g.GameBoard.SetPosition('X', pos)
	if err != nil {
		return err
	}
	g.Move = g.Bob
	return nil
}

func (g *Game) BobMove(pos rune) error {
	if g.Move != g.Bob {
		return fmt.Errorf("Player 1's turn")
	}
	err := g.GameBoard.SetPosition('O', pos)
	if err != nil {
		return err
	}
	g.Move = g.Alice
	return nil
}

//TODO: Detect unwinnable games
func (g *Game) IsTie() bool {
	return !g.GameBoard.IsWinner('X') && !g.GameBoard.IsWinner('O') && g.GameBoard.IsFull()
}

func (g *Game) GetWinner() *Player {
	if g.GameBoard.IsWinner('X') {
		return g.Alice
	} else if g.GameBoard.IsWinner('O') {
		return g.Bob
	}
	return nil
}
