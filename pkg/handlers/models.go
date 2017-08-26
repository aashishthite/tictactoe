package handlers

type StartGameRequest struct {
	Player1ID string `json:"player_1"`
	Player2ID string `json:"player_2"`
}

type MoveRequest struct {
	PlayerID string `json:"player"`
	Position string `json:"position"`
}
