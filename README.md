# ♔ Chess in Go

A beautiful command-line chess game implemented in Go, featuring:
- Full chess rules implementation
- Unicode chess pieces (♔♕♖♗♘♙ ♚♛♜♝♞♟)
- Algebraic notation for moves
- Game state tracking (check, checkmate, stalemate)
- Clean and modular code structure

## 📦 Installation

1. Make sure you have Go installed (version 1.16 or later)
2. Clone the repository:
```bash
git clone https://github.com/yourusername/chess.git
cd chess
```

3. Build and run the game:
```bash
go run cmd/chess/main.go
```

## 🎮 How to Play

1. The game uses standard algebraic notation for moves
2. Enter moves in the format: `source destination` (e.g., `e2 e4`)
3. White pieces are uppercase (♔♕♖♗♘♙)
4. Black pieces are lowercase (♚♛♜♝♞♟)
5. Type `quit` to exit the game

### Example Moves
- `e2 e4` - Move pawn from e2 to e4
- `g1 f3` - Move knight from g1 to f3
- `e7 e5` - Move pawn from e7 to e5

## 🏗️ Project Structure

```
chess/
├── cmd/
│   └── chess/
│       └── main.go       # Entry point
├── pkg/
│   ├── board/
│   │   └── board.go      # Board representation
│   ├── game/
│   │   └── game.go       # Game logic
│   ├── piece/
│   │   └── piece.go      # Piece movement
│   └── ui/
│       └── ui.go         # User interface
├── go.mod
└── README.md
```

## ✨ Features

- [x] Complete chess rules implementation
- [x] Valid move checking
- [x] Check detection
- [x] Checkmate detection
- [x] Stalemate detection
- [x] Beautiful Unicode chess pieces
- [x] Clear and intuitive interface

## 🎯 Future Enhancements

- [ ] PGN notation support
- [ ] Game save/load functionality
- [ ] Time controls
- [ ] AI opponent
- [ ] Network play
- [ ] Opening book
- [ ] Move history display

## 🤝 Contributing

Contributions are welcome! Feel free to:
1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Inspired by traditional chess implementations
- Built with Go's standard library
- Unicode chess symbols for beautiful display 