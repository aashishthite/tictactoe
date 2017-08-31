package core

import (
	"fmt"
	"log"
	"time"
)

const DEFAULT_TIMEOUT = 1 * time.Minute

type Engine struct {
	OngoingGames map[string]*Game
}

type State struct {
	GameBoard string `json:"board"`
	Status    string `json:"status"`
}

//TODO: implement timeout
func NewEngine() *Engine {
	retval := &Engine{
		OngoingGames: make(map[string]*Game),
	}
	return retval
}

func (e *Engine) Start(p1, p2 *Player, channelID string) (*State, error) {
	if _, ok := e.OngoingGames[channelID]; ok {
		return nil, fmt.Errorf("Game is still going on")
	}
	e.OngoingGames[channelID] = NewGame(p1, p2)

	retval := &State{
		GameBoard: e.OngoingGames[channelID].GameBoard.PrettyPrint(),
		Status:    fmt.Sprintf("Player: %s 's turn to make a move", p1.ID),
	}
	return retval, nil
}

func (e *Engine) Move(p *Player, pos rune, channelID string) (*State, error) {
	currentGame, ok := e.OngoingGames[channelID]
	if !ok {
		return &State{}, fmt.Errorf("Game does not exist")
	}
	var err error
	if p.ID == currentGame.Alice.ID {
		err = currentGame.AliceMove(pos)
	} else if p.ID == currentGame.Bob.ID {
		err = currentGame.BobMove(pos)
	} else {
		err = fmt.Errorf("Invalid Player")
	}
	if err != nil {
		return nil, err
	}
	if currentGame.IsTie() {
		log.Printf("Game is a tie.")
		retval := &State{
			GameBoard: currentGame.GameBoard.PrettyPrint(),
			Status:    "Game is a tie.",
		}
		delete(e.OngoingGames, channelID)
		return retval, nil
	}
	if player := currentGame.GetWinner(); player != nil {

		log.Printf("Player: %s has won!", "@"+player.ID)
		retval := &State{
			GameBoard: currentGame.GameBoard.PrettyPrint(),
			Status:    fmt.Sprintf("Player: %s has won!", "@"+player.ID),
		}
		delete(e.OngoingGames, channelID)
		return retval, nil
	}
	retval := &State{
		GameBoard: currentGame.GameBoard.PrettyPrint(),
		Status:    fmt.Sprintf("Player: @%s 's turn to make a move", currentGame.Move.ID),
	}
	return retval, nil
}

func (e *Engine) GameState(channelID string) *State {
	if v, ok := e.OngoingGames[channelID]; ok {
		return &State{
			GameBoard: v.GameBoard.PrettyPrint(),
			Status:    fmt.Sprintf("Player: @%s 's turn to make a move", v.Move.ID),
		}
	}
	return &State{
		GameBoard: NewBoard().PrettyPrint(),
		Status:    "No Game Running",
	}
}

func (e *Engine) Forfeit(player string, channelID string) string {
	if v, ok := e.OngoingGames[channelID]; ok {
		if v.Alice.ID == player || v.Bob.ID == player {
			delete(e.OngoingGames, channelID)
			return fmt.Sprintf("Player: @%s has forfeited the game", player)
		}
		return "Only participating players can forfeit"

	}
	return "Game not found"

}
