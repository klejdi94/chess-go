package board

import (
	"testing"
)

func TestNewPosition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantRow int
		wantCol int
		wantErr bool
	}{
		{"valid e2", "e2", 6, 4, false},
		{"valid a1", "a1", 7, 0, false},
		{"valid h8", "h8", 0, 7, false},
		{"invalid i9", "i9", 0, 0, true},
		{"invalid a0", "a0", 0, 0, true},
		{"invalid format", "aa", 0, 0, true},
		{"too short", "a", 0, 0, true},
		{"too long", "a1b", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPosition(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPosition(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Row != tt.wantRow || got.Col != tt.wantCol {
					t.Errorf("NewPosition(%q) = {%d,%d}, want {%d,%d}",
						tt.input, got.Row, got.Col, tt.wantRow, tt.wantCol)
				}
			}
		})
	}
}

func TestNewBoard(t *testing.T) {
	board := NewBoard()

	// Test pawns
	for col := 0; col < 8; col++ {
		if board.Squares[1][col].Type != Pawn || board.Squares[1][col].Color != Black {
			t.Errorf("Expected black pawn at position {1,%d}, got %v", col, board.Squares[1][col])
		}
		if board.Squares[6][col].Type != Pawn || board.Squares[6][col].Color != White {
			t.Errorf("Expected white pawn at position {6,%d}, got %v", col, board.Squares[6][col])
		}
	}

	// Test back rank pieces
	backRank := []PieceType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for col := 0; col < 8; col++ {
		if board.Squares[0][col].Type != backRank[col] || board.Squares[0][col].Color != Black {
			t.Errorf("Expected black %v at position {0,%d}, got %v",
				backRank[col], col, board.Squares[0][col])
		}
		if board.Squares[7][col].Type != backRank[col] || board.Squares[7][col].Color != White {
			t.Errorf("Expected white %v at position {7,%d}, got %v",
				backRank[col], col, board.Squares[7][col])
		}
	}
}

func TestMovePiece(t *testing.T) {
	board := NewBoard()
	from := Position{6, 4} // e2
	to := Position{4, 4}   // e4

	// Test moving a pawn
	piece := board.GetPiece(from)
	board.MovePiece(from, to)

	// Check if piece moved correctly
	if board.GetPiece(to) != piece {
		t.Errorf("Piece not moved correctly to destination")
	}

	// Check if original position is empty
	if !board.IsEmpty(from) {
		t.Errorf("Original position not empty after move")
	}
}
