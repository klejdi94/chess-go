package board

import (
	"fmt"
	"strings"
)

// PieceType represents the type of chess piece
type PieceType int

const (
	Empty PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

// Color represents the color of a chess piece
type Color int

const (
	NoColor Color = iota
	White
	Black
)

// Piece represents a chess piece
type Piece struct {
	Type  PieceType
	Color Color
}

func (p Piece) String() string {
	symbols := map[PieceType]map[Color]string{
		Pawn:   {White: "♙", Black: "♟"},
		Knight: {White: "♘", Black: "♞"},
		Bishop: {White: "♗", Black: "♝"},
		Rook:   {White: "♖", Black: "♜"},
		Queen:  {White: "♕", Black: "♛"},
		King:   {White: "♔", Black: "♚"},
		Empty:  {NoColor: " ", White: " ", Black: " "},
	}
	return symbols[p.Type][p.Color]
}

// ASCIIString returns a plain ASCII representation of the piece
func (p Piece) ASCIIString() string {
	if p.Type == Empty {
		return " "
	}

	symbols := map[PieceType]string{
		Pawn:   "P",
		Knight: "N",
		Bishop: "B",
		Rook:   "R",
		Queen:  "Q",
		King:   "K",
	}

	symbol := symbols[p.Type]
	if p.Color == Black {
		return strings.ToLower(symbol)
	}
	return symbol
}

// Position represents a position on the chess board
type Position struct {
	Row int // 0-7
	Col int // 0-7
}

func NewPosition(algebraic string) (Position, error) {
	if len(algebraic) != 2 {
		return Position{}, fmt.Errorf("invalid algebraic notation: %s", algebraic)
	}

	col := int(algebraic[0] - 'a')
	row := int('8' - algebraic[1])

	if col < 0 || col > 7 || row < 0 || row > 7 {
		return Position{}, fmt.Errorf("position out of bounds: %s", algebraic)
	}

	return Position{Row: row, Col: col}, nil
}

func (p Position) String() string {
	return fmt.Sprintf("%c%c", 'a'+rune(p.Col), '8'-rune(p.Row))
}

// Board represents a chess board
type Board struct {
	Squares [8][8]Piece
}

// NewBoard creates a new chess board with pieces in the initial position
func NewBoard() *Board {
	board := &Board{}

	// Set up pawns
	for col := 0; col < 8; col++ {
		board.Squares[1][col] = Piece{Type: Pawn, Color: Black}
		board.Squares[6][col] = Piece{Type: Pawn, Color: White}
	}

	// Set up other pieces
	backRank := []PieceType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for col := 0; col < 8; col++ {
		board.Squares[0][col] = Piece{Type: backRank[col], Color: Black}
		board.Squares[7][col] = Piece{Type: backRank[col], Color: White}
	}

	return board
}

// GetPiece returns the piece at the given position
func (b *Board) GetPiece(pos Position) Piece {
	return b.Squares[pos.Row][pos.Col]
}

// SetPiece sets a piece at the given position
func (b *Board) SetPiece(pos Position, piece Piece) {
	b.Squares[pos.Row][pos.Col] = piece
}

// MovePiece moves a piece from one position to another
func (b *Board) MovePiece(from, to Position) {
	piece := b.GetPiece(from)
	b.SetPiece(to, piece)
	b.SetPiece(from, Piece{Type: Empty, Color: NoColor})
}

// IsEmpty checks if a position is empty
func (b *Board) IsEmpty(pos Position) bool {
	return b.Squares[pos.Row][pos.Col].Type == Empty
}

// Print prints the current state of the board
func (b *Board) Print() {
	fmt.Println("  a b c d e f g h")
	fmt.Println(" +-----------------+")
	for row := 0; row < 8; row++ {
		// Print the row number
		fmt.Printf("%d|", 8-row)

		// Print each piece in the row
		for col := 0; col < 8; col++ {
			fmt.Printf(" %s", b.Squares[row][col])
		}

		// Print the right border and row number
		fmt.Printf(" |%d\n", 8-row)
	}
	fmt.Println(" +-----------------+")
	fmt.Println("  a b c d e f g h")
}

// PrintASCII prints the board using ASCII characters for better console compatibility
func (b *Board) PrintASCII() {
	fmt.Println("  a b c d e f g h")
	fmt.Println(" +---------------+")
	for row := 0; row < 8; row++ {
		fmt.Printf("%d|", 8-row)
		for col := 0; col < 8; col++ {
			piece := b.Squares[row][col]
			fmt.Printf(" %s", piece.ASCIIString())
		}
		fmt.Printf(" |%d\n", 8-row)
	}
	fmt.Println(" +---------------+")
	fmt.Println("  a b c d e f g h")
}
