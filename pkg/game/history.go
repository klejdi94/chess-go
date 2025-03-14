package game

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/chess/pkg/board"
)

// GameHistory represents a complete game with its moves
type GameHistory struct {
	Date        time.Time `json:"date"`
	Moves       []string  `json:"moves"`
	Result      string    `json:"result"`
	WhitePlayer string    `json:"white_player"`
	BlackPlayer string    `json:"black_player"`
}

// SaveGame saves the current game to a JSON file
func (g *Game) SaveGame(filename string) error {
	// Convert []Move to []string
	moveStrings := make([]string, len(g.moveHistory))
	for i, move := range g.moveHistory {
		moveStrings[i] = fmt.Sprintf("%s %s", move.From.String(), move.To.String())
	}

	history := GameHistory{
		Date:        time.Now(),
		Moves:       moveStrings,
		Result:      g.GetGameStatus(),
		WhitePlayer: "Player 1",
		BlackPlayer: "Player 2",
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling game history: %v", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error saving game history: %v", err)
	}

	return nil
}

// LoadGame loads a game from a JSON file
func (g *Game) LoadGame(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading game file: %v", err)
	}

	var history GameHistory
	err = json.Unmarshal(data, &history)
	if err != nil {
		return fmt.Errorf("error unmarshaling game history: %v", err)
	}

	// Reset the game
	g.Board = board.NewBoard()
	g.CurrentPlayer = board.White
	g.State = InProgress
	g.moveHistory = make([]Move, 0)

	// Replay all moves
	for _, moveStr := range history.Moves {
		// Parse and apply each move
		var from, to string
		fmt.Sscanf(moveStr, "%s %s", &from, &to)

		fromPos, err := board.NewPosition(from)
		if err != nil {
			return fmt.Errorf("invalid move in history %s: %v", moveStr, err)
		}

		toPos, err := board.NewPosition(to)
		if err != nil {
			return fmt.Errorf("invalid move in history %s: %v", moveStr, err)
		}

		err = g.MakeMove(fromPos, toPos)
		if err != nil {
			return fmt.Errorf("error replaying move %s: %v", moveStr, err)
		}
	}

	return nil
}
