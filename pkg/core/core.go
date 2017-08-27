package core

import (
	"fmt"
	"log"
	"time"
)

const DEFAULT_TIMEOUT = 1 * time.Minute

type Engine struct {
	CurrentGame *Game
	GameOn      bool //use mutex
	timer       *time.Timer
}

type State struct {
	GameBoard string `json:"board"`
	Status    string `json:"status"`
}

//TODO: implement timeout
func NewEngine() *Engine {
	retval := &Engine{}
	return retval
}

func (e *Engine) Start(p1, p2 *Player) (*State, error) {
	if e.GameOn {
		return nil, fmt.Errorf("Game is still going on")
	}
	e.CurrentGame = NewGame(p1, p2)
	e.GameOn = true
	e.timer = time.AfterFunc(DEFAULT_TIMEOUT, e.TearDown)

	retval := &State{
		GameBoard: e.CurrentGame.GameBoard.PrettyPrint(),
		Status:    fmt.Sprintf("Player: %s 's turn to make a move", p1.ID),
	}
	return retval, nil
}

func (e *Engine) Move(p *Player, pos rune) (*State, error) {
	var err error
	if p.ID == e.CurrentGame.Alice.ID {
		err = e.CurrentGame.AliceMove(pos)
	} else if p.ID == e.CurrentGame.Bob.ID {
		err = e.CurrentGame.BobMove(pos)
	} else {
		err = fmt.Errorf("Invalid Player")
	}
	if err != nil {
		return nil, err
	}
	if e.CurrentGame.IsTie() {
		e.GameOn = false
		log.Printf("Game is a tie.")
		retval := &State{
			GameBoard: e.CurrentGame.GameBoard.PrettyPrint(),
			Status:    "Game is a tie.",
		}
		return retval, nil
	}
	if player := e.CurrentGame.GetWinner(); player != nil {
		e.GameOn = false
		log.Printf("Player: %s has won!", player.ID)
		retval := &State{
			GameBoard: e.CurrentGame.GameBoard.PrettyPrint(),
			Status:    fmt.Sprintf("Player: %s has won!", player.ID),
		}
		return retval, nil
	}
	retval := &State{
		GameBoard: e.CurrentGame.GameBoard.PrettyPrint(),
		Status:    fmt.Sprintf("Player: %s 's turn to make a move'", e.CurrentGame.Move.ID),
	}
	return retval, nil
}

/*
func (e *Engine) Forfiet(p *Player) {
	e.GameOn = false
}
*/

func (e *Engine) GameState() *State {
	if e.GameOn {
		return &State{
			GameBoard: e.CurrentGame.GameBoard.PrettyPrint(),
			Status:    fmt.Sprintf("Player: %s 's turn to make a move'", e.CurrentGame.Move.ID),
		}
	}
	return &State{
		GameBoard: NewBoard().PrettyPrint(),
		Status:    "No Game Running",
	}
}

func (e *Engine) TearDown() {
	log.Printf("Game did not complete in a minute. Game is a tie.")
	e.GameOn = false
	e.timer.Stop()
}
