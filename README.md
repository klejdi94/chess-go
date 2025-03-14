# Chess Game in Go

A beautiful and feature-rich command-line chess game implemented in Go. Play chess with a friend in your terminal with support for time controls, game saving, and more!

## ♟️ Features

- Interactive command-line interface
- Unicode chess pieces (with ASCII fallback)
- Time control with increment
- Game save/load functionality
- Player names support
- Comprehensive test coverage

## 🚀 Quick Start

```bash
# Install the game
go install github.com/klejdi94/chess-go/cmd/chess@latest

# Start a new game
chess

# Start with custom settings
chess -player1 "Alice" -player2 "Bob"    # Set player names
chess -save "game.json"                  # Save game to file
chess -load "game.json"                  # Load game from file
chess -time "5,3"                        # 5 minutes + 3 seconds increment
chess -no-timer                          # Disable time control
```

## ⚙️ Time Control

The game supports chess clocks with increment:
- Each player starts with a main time bank (default: 10 minutes)
- After each move, the player gets additional time (default: 5 seconds)
- The clock shows remaining time for the current player
- A player loses if their time runs out

Example time controls:
```bash
chess -time "3,2"     # 3 minutes + 2 seconds per move
chess -time "5,0"     # 5 minutes with no increment
chess -time "15,10"   # 15 minutes + 10 seconds per move
chess -no-timer       # Disable time control
```

## 📋 Game Commands

During the game, you can use these commands:
- Move pieces using algebraic notation (e.g., `e2e4`, `Nf3`)
- Type `quit` to exit the game
- Type `save` to save the current game
- Type `help` to see all commands

## 🎯 Roadmap

Completed:
- [x] Basic chess rules and movement
- [x] Game save/load functionality
- [x] Player names support
- [x] Time control with increment
- [x] Clear and intuitive interface
- [x] Comprehensive test coverage

Planned:
- [ ] PGN notation support
- [ ] AI opponent
- [ ] Network play
- [ ] Undo/redo functionality
- [ ] Game analysis tools
- [ ] Tournament mode

## 🤝 Contributing

Contributions are welcome! Feel free to:
1. Fork the repository
2. Create a feature branch
3. Submit a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 