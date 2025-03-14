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
	game      *game.Game
	scanner   *bufio.Scanner
	useAscii  bool
	whiteName string
	blackName string
}

// NewUI creates a new UI
func NewUI(game *game.Game) *UI {
	return &UI{
		game:      game,
		scanner:   bufio.NewScanner(os.Stdin),
		useAscii:  false,
		whiteName: "White",
		blackName: "Black",
	}
}

// SetAsciiMode sets whether to use ASCII characters instead of Unicode
func (ui *UI) SetAsciiMode(ascii bool) {
	ui.useAscii = ascii
}

// SetPlayerNames sets the names of the players
func (ui *UI) SetPlayerNames(white, black string) {
	ui.whiteName = white
	ui.blackName = black
}

// Start starts the UI
func (ui *UI) Start() {
	fmt.Println("Welcome to Chess in Go!")
	fmt.Printf("Players: %s (White) vs %s (Black)\n", ui.whiteName, ui.blackName)
	fmt.Println("Enter moves in algebraic notation (e.g., 'e2 e4' to move from e2 to e4)")
	fmt.Println("Type 'quit' to exit")

	for {
		if ui.useAscii {
			ui.game.Board.PrintASCII()
		} else {
			ui.game.Board.Print()
		}
		fmt.Println(ui.getGameStatus())

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

// getGameStatus returns a string representation of the game status with player names and time
func (ui *UI) getGameStatus() string {
	state := ui.game.GetGameStatus()
	timeLeft := ui.game.GetTimeLeft()

	if timeLeft != "" {
		if ui.game.CurrentPlayer == board.White {
			return fmt.Sprintf("%s's turn [%s] (%s)", ui.whiteName, timeLeft, state)
		}
		return fmt.Sprintf("%s's turn [%s] (%s)", ui.blackName, timeLeft, state)
	}

	if ui.game.CurrentPlayer == board.White {
		return fmt.Sprintf("%s's turn (%s)", ui.whiteName, state)
	}
	return fmt.Sprintf("%s's turn (%s)", ui.blackName, state)
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
