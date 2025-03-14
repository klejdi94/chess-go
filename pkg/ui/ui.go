package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/user/chess/pkg/board"
	"github.com/user/chess/pkg/game"
)

// UI represents the user interface for the chess game
type UI struct {
	game    *game.Game
	scanner *bufio.Scanner
}

// NewUI creates a new UI
func NewUI(game *game.Game) *UI {
	return &UI{
		game:    game,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// Start starts the UI
func (ui *UI) Start() {
	fmt.Println("Welcome to Chess in Go!")
	fmt.Println("Enter moves in algebraic notation (e.g., 'e2 e4' to move from e2 to e4)")
	fmt.Println("Type 'quit' to exit")

	for {
		ui.game.Board.Print()
		fmt.Println(ui.game.GetGameStatus())

		if ui.game.State == game.Checkmate || ui.game.State == game.Stalemate || ui.game.State == game.Draw {
			break
		}

		move := ui.getMove()
		if move == "quit" {
			break
		}
	}

	fmt.Println("Game over")
}

// getMove gets a move from the user
func (ui *UI) getMove() string {
	for {
		fmt.Print("Enter move: ")
		ui.scanner.Scan()
		input := ui.scanner.Text()
		input = strings.TrimSpace(input)

		if input == "quit" {
			return "quit"
		}

		// Parse move
		parts := strings.Fields(input)
		if len(parts) != 2 {
			fmt.Println("Invalid input. Please enter source and destination (e.g., 'e2 e4')")
			continue
		}

		from, err := board.NewPosition(parts[0])
		if err != nil {
			fmt.Println("Invalid source position:", err)
			continue
		}

		to, err := board.NewPosition(parts[1])
		if err != nil {
			fmt.Println("Invalid destination position:", err)
			continue
		}

		// Try to make the move
		err = ui.game.MakeMove(from, to)
		if err != nil {
			fmt.Println("Invalid move:", err)
			continue
		}

		return input
	}
}
