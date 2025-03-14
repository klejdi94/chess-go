package piece

import (
	"github.com/user/chess/pkg/board"
)

// MoveValidator defines an interface for validating piece moves
type MoveValidator interface {
	ValidateMoves(pos board.Position, b *board.Board) []board.Position
}

// GetValidMoves returns all valid moves for a piece at the given position
func GetValidMoves(pos board.Position, b *board.Board) []board.Position {
	piece := b.GetPiece(pos)
	var moves []board.Position

	switch piece.Type {
	case board.Pawn:
		moves = getValidPawnMoves(pos, piece.Color, b)
	case board.Knight:
		moves = getValidKnightMoves(pos, piece.Color, b)
	case board.Bishop:
		moves = getValidBishopMoves(pos, piece.Color, b)
	case board.Rook:
		moves = getValidRookMoves(pos, piece.Color, b)
	case board.Queen:
		bishopMoves := getValidBishopMoves(pos, piece.Color, b)
		rookMoves := getValidRookMoves(pos, piece.Color, b)
		moves = append(bishopMoves, rookMoves...)
	case board.King:
		moves = getValidKingMoves(pos, piece.Color, b)
	}

	return moves
}

// Helper functions for each piece type
func getValidPawnMoves(pos board.Position, color board.Color, b *board.Board) []board.Position {
	var moves []board.Position

	// Determine the direction pawns move based on their color
	forwardDir := -1
	startRow := 6
	if color == board.Black {
		forwardDir = 1
		startRow = 1
	}

	// Check one square forward
	forwardPos := board.Position{Row: pos.Row + forwardDir, Col: pos.Col}
	if isValidPosition(forwardPos) && b.IsEmpty(forwardPos) {
		moves = append(moves, forwardPos)

		// Check two squares forward from starting position
		if pos.Row == startRow {
			twoForwardPos := board.Position{Row: pos.Row + 2*forwardDir, Col: pos.Col}
			if b.IsEmpty(twoForwardPos) {
				moves = append(moves, twoForwardPos)
			}
		}
	}

	// Check diagonal captures
	for _, colOffset := range []int{-1, 1} {
		capturePos := board.Position{Row: pos.Row + forwardDir, Col: pos.Col + colOffset}
		if isValidPosition(capturePos) {
			capturePiece := b.GetPiece(capturePos)
			if capturePiece.Type != board.Empty && capturePiece.Color != color {
				moves = append(moves, capturePos)
			}
		}
	}

	return moves
}

func getValidKnightMoves(pos board.Position, color board.Color, b *board.Board) []board.Position {
	var moves []board.Position

	// Knight moves in L-shape
	knightOffsets := []struct{ row, col int }{
		{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2},
		{1, -2}, {1, 2}, {2, -1}, {2, 1},
	}

	for _, offset := range knightOffsets {
		newPos := board.Position{Row: pos.Row + offset.row, Col: pos.Col + offset.col}
		if isValidPosition(newPos) {
			piece := b.GetPiece(newPos)
			if piece.Type == board.Empty || piece.Color != color {
				moves = append(moves, newPos)
			}
		}
	}

	return moves
}

func getValidBishopMoves(pos board.Position, color board.Color, b *board.Board) []board.Position {
	// Bishop moves diagonally
	directions := []struct{ rowDir, colDir int }{
		{-1, -1}, {-1, 1}, {1, -1}, {1, 1},
	}

	return getSlidingMoves(pos, color, b, directions)
}

func getValidRookMoves(pos board.Position, color board.Color, b *board.Board) []board.Position {
	// Rook moves horizontally and vertically
	directions := []struct{ rowDir, colDir int }{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	}

	return getSlidingMoves(pos, color, b, directions)
}

func getValidKingMoves(pos board.Position, color board.Color, b *board.Board) []board.Position {
	var moves []board.Position

	// King moves one square in any direction
	for rowOffset := -1; rowOffset <= 1; rowOffset++ {
		for colOffset := -1; colOffset <= 1; colOffset++ {
			if rowOffset == 0 && colOffset == 0 {
				continue
			}

			newPos := board.Position{Row: pos.Row + rowOffset, Col: pos.Col + colOffset}
			if isValidPosition(newPos) {
				piece := b.GetPiece(newPos)
				if piece.Type == board.Empty || piece.Color != color {
					moves = append(moves, newPos)
				}
			}
		}
	}

	return moves
}

// Helper function for sliding pieces (bishop, rook, queen)
func getSlidingMoves(pos board.Position, color board.Color, b *board.Board, directions []struct{ rowDir, colDir int }) []board.Position {
	var moves []board.Position

	for _, dir := range directions {
		for i := 1; i < 8; i++ {
			newPos := board.Position{Row: pos.Row + i*dir.rowDir, Col: pos.Col + i*dir.colDir}
			if !isValidPosition(newPos) {
				break
			}

			piece := b.GetPiece(newPos)
			if piece.Type == board.Empty {
				moves = append(moves, newPos)
			} else if piece.Color != color {
				moves = append(moves, newPos)
				break
			} else {
				break
			}
		}
	}

	return moves
}

// Helper function to check if a position is on the board
func isValidPosition(pos board.Position) bool {
	return pos.Row >= 0 && pos.Row < 8 && pos.Col >= 0 && pos.Col < 8
}
