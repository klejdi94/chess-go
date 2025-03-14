package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/chess/pkg/game"
	"github.com/user/chess/pkg/ui"
)

func main() {
	// Command line flags
	playerNames := flag.String("names", "Player1,Player2", "Names of the two players (comma-separated)")
	saveFile := flag.String("save", "", "Save game to specified file")
	loadFile := flag.String("load", "", "Load game from specified file")
	ascii := flag.Bool("ascii", false, "Use ASCII characters instead of Unicode")
	help := flag.Bool("help", false, "Show help message")
	timeControl := flag.String("time", "10,5", "Time control in minutes,increment_seconds (e.g., '10,5' for 10 minutes + 5 seconds increment)")
	noTimer := flag.Bool("no-timer", false, "Disable time control")

	flag.Parse()

	if *help {
		fmt.Println("Chess Game in Go")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Parse player names
	names := strings.Split(*playerNames, ",")
	whiteName := "White"
	blackName := "Black"
	if len(names) >= 2 {
		whiteName = strings.TrimSpace(names[0])
		blackName = strings.TrimSpace(names[1])
	}

	// Create new game
	g := game.NewGame()

	// Configure time control
	if !*noTimer {
		times := strings.Split(*timeControl, ",")
		if len(times) >= 2 {
			initialMinutes, err1 := strconv.Atoi(strings.TrimSpace(times[0]))
			incrementSeconds, err2 := strconv.Atoi(strings.TrimSpace(times[1]))

			if err1 == nil && err2 == nil && initialMinutes > 0 {
				g.TimeControl = game.NewTimeControl(initialMinutes, incrementSeconds)
			} else {
				fmt.Println("Invalid time control format, using default (10 minutes + 5 seconds)")
			}
		}
	} else {
		g.TimeControl = nil
	}

	// Load game if specified
	if *loadFile != "" {
		err := g.LoadGame(*loadFile)
		if err != nil {
			fmt.Printf("Error loading game: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Game loaded from %s\n", *loadFile)
	}

	// Create UI
	ui := ui.NewUI(g)
	ui.SetAsciiMode(*ascii)
	ui.SetPlayerNames(whiteName, blackName)

	// Start the game
	ui.Start()

	// Save game if specified
	if *saveFile != "" {
		err := g.SaveGame(*saveFile)
		if err != nil {
			fmt.Printf("Error saving game: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Game saved to %s\n", *saveFile)
	}
}
