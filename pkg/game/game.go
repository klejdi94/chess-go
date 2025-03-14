package game

import (
	"errors"
	"fmt"

	"github.com/user/chess/pkg/board"
)

// GameState represents the state of a chess game
type GameState int

const (
	InProgress GameState = iota
	Check
	Checkmate
	Stalemate
	Draw
	TimeOut
)

// Game represents a chess game
type Game struct {
	Board           *board.Board
	CurrentPlayer   board.Color
	moveHistory     []Move
	castlingRights  map[board.Color]CastlingRights
	enPassantTarget *board.Position
	halfMoveClock   int // For 50-move rule
	fullMoveNumber  int
	State           GameState
	TimeControl     *TimeControl
}

// Move represents a chess move
type Move struct {
	From          board.Position
	To            board.Position
	PromotionType board.PieceType
	Notation      string
}

// CastlingRights represents the castling rights for a player
type CastlingRights struct {
	KingSide  bool
	QueenSide bool
}

// NewGame creates a new chess game
func NewGame() *Game {
	return &Game{
		Board:         board.NewBoard(),
		CurrentPlayer: board.White,
		castlingRights: map[board.Color]CastlingRights{
			board.White: {KingSide: true, QueenSide: true},
			board.Black: {KingSide: true, QueenSide: true},
		},
		halfMoveClock:  0,
		fullMoveNumber: 1,
		State:          InProgress,
		TimeControl:    NewTimeControl(10, 5), // 10 minutes + 5 seconds increment
	}
}

// IsValidMove checks if a move is valid
func (g *Game) IsValidMove(from, to board.Position) bool {
	// Get the piece at the source position
	piece := g.Board.GetPiece(from)

	// Check if there's a piece at the source position and it belongs to the current player
	if piece.Type == board.Empty || piece.Color != g.CurrentPlayer {
		return false
	}

	// Check if the destination is occupied by a piece of the same color
	destPiece := g.Board.GetPiece(to)
	if destPiece.Type != board.Empty && destPiece.Color == g.CurrentPlayer {
		return false
	}

	// Check if the move is valid for the specific piece type
	return g.isValidPieceMove(from, to, piece)
}

// isValidPieceMove checks if a move is valid for a specific piece
func (g *Game) isValidPieceMove(from, to board.Position, piece board.Piece) bool {
	rowDiff := to.Row - from.Row
	colDiff := to.Col - from.Col
	absRowDiff := abs(rowDiff)
	absColDiff := abs(colDiff)

	switch piece.Type {
	case board.Pawn:
		return g.isValidPawnMove(from, to, rowDiff, colDiff, absRowDiff, absColDiff)
	case board.Knight:
		// Knights move in an L-shape: 2 squares in one direction and 1 square perpendicular
		return (absRowDiff == 2 && absColDiff == 1) || (absRowDiff == 1 && absColDiff == 2)
	case board.Bishop:
		// Bishops move diagonally
		if absRowDiff != absColDiff {
			return false
		}
		return g.isClearPath(from, to)
	case board.Rook:
		// Rooks move horizontally or vertically
		if from.Row != to.Row && from.Col != to.Col {
			return false
		}
		return g.isClearPath(from, to)
	case board.Queen:
		// Queens move like bishops or rooks
		if from.Row != to.Row && from.Col != to.Col && absRowDiff != absColDiff {
			return false
		}
		return g.isClearPath(from, to)
	case board.King:
		// Check for castling
		if absRowDiff == 0 && absColDiff == 2 {
			return g.isValidCastling(from, to)
		}
		// Kings move one square in any direction
		return absRowDiff <= 1 && absColDiff <= 1
	}

	return false
}

// isValidPawnMove checks if a pawn move is valid
func (g *Game) isValidPawnMove(from, to board.Position, rowDiff, colDiff, absRowDiff, absColDiff int) bool {
	// Determine the direction pawns move based on their color
	forwardDir := -1
	startRow := 6
	if g.CurrentPlayer == board.Black {
		forwardDir = 1
		startRow = 1
	}

	// Check for normal pawn move (1 square forward)
	if colDiff == 0 && rowDiff == forwardDir && g.Board.IsEmpty(to) {
		return true
	}

	// Check for pawn's first move (2 squares forward)
	if colDiff == 0 && from.Row == startRow && rowDiff == 2*forwardDir {
		midPos := board.Position{Row: from.Row + forwardDir, Col: from.Col}
		return g.Board.IsEmpty(midPos) && g.Board.IsEmpty(to)
	}

	// Check for capture (1 square diagonally)
	if absColDiff == 1 && rowDiff == forwardDir {
		// Normal capture
		if !g.Board.IsEmpty(to) && g.Board.GetPiece(to).Color != g.CurrentPlayer {
			return true
		}

		// En passant capture
		if g.enPassantTarget != nil && to.Row == g.enPassantTarget.Row && to.Col == g.enPassantTarget.Col {
			return true
		}
	}

	return false
}

// isValidCastling checks if a castling move is valid
func (g *Game) isValidCastling(from, to board.Position) bool {
	// Check if castling rights are available
	rights := g.castlingRights[g.CurrentPlayer]

	// Determine if it's kingside or queenside castling
	var rookCol int
	if to.Col > from.Col { // Kingside
		if !rights.KingSide {
			return false
		}
		rookCol = 7
	} else { // Queenside
		if !rights.QueenSide {
			return false
		}
		rookCol = 0
	}

	rookRow := from.Row
	rookPos := board.Position{Row: rookRow, Col: rookCol}

	// Check if there are pieces between the king and the rook
	return g.isClearPath(from, rookPos) && !g.isInCheck(g.CurrentPlayer)
}

// isClearPath checks if there are no pieces between the start and end positions
func (g *Game) isClearPath(from, to board.Position) bool {
	// Determine the direction
	rowDir := 0
	if to.Row > from.Row {
		rowDir = 1
	} else if to.Row < from.Row {
		rowDir = -1
	}

	colDir := 0
	if to.Col > from.Col {
		colDir = 1
	} else if to.Col < from.Col {
		colDir = -1
	}

	row, col := from.Row+rowDir, from.Col+colDir
	for {
		if row == to.Row && col == to.Col {
			break
		}

		if !g.Board.IsEmpty(board.Position{Row: row, Col: col}) {
			return false
		}

		row += rowDir
		col += colDir
	}

	return true
}

// MakeMove makes a move on the board and updates the game state
func (g *Game) MakeMove(from, to board.Position) error {
	if g.State != InProgress {
		return fmt.Errorf("game is already finished")
	}

	if g.TimeControl != nil && g.TimeControl.IsTimeUp(g.CurrentPlayer == board.White) {
		g.State = TimeOut
		return fmt.Errorf("time is up for %s", g.GetCurrentPlayerName())
	}

	if !g.IsValidMove(from, to) {
		return errors.New("invalid move")
	}

	piece := g.Board.GetPiece(from)
	capturedPiece := g.Board.GetPiece(to)

	// Handle en passant capture
	if piece.Type == board.Pawn && g.enPassantTarget != nil && to.Row == g.enPassantTarget.Row && to.Col == g.enPassantTarget.Col {
		capturePos := board.Position{Row: from.Row, Col: to.Col}
		g.Board.SetPiece(capturePos, board.Piece{Type: board.Empty, Color: board.NoColor})
	}

	// Update en passant target
	g.enPassantTarget = nil
	if piece.Type == board.Pawn && abs(to.Row-from.Row) == 2 {
		enPassantRow := (from.Row + to.Row) / 2
		g.enPassantTarget = &board.Position{Row: enPassantRow, Col: from.Col}
	}

	// Handle castling
	if piece.Type == board.King && abs(to.Col-from.Col) == 2 {
		rookRow := from.Row
		var oldRookCol, newRookCol int
		if to.Col > from.Col { // Kingside
			oldRookCol = 7
			newRookCol = 5
		} else { // Queenside
			oldRookCol = 0
			newRookCol = 3
		}
		rookFrom := board.Position{Row: rookRow, Col: oldRookCol}
		rookTo := board.Position{Row: rookRow, Col: newRookCol}
		g.Board.MovePiece(rookFrom, rookTo)
	}

	// Update castling rights
	if piece.Type == board.King {
		g.castlingRights[g.CurrentPlayer] = CastlingRights{false, false}
	} else if piece.Type == board.Rook {
		rights := g.castlingRights[g.CurrentPlayer]
		if from.Col == 0 {
			rights.QueenSide = false
		} else if from.Col == 7 {
			rights.KingSide = false
		}
		g.castlingRights[g.CurrentPlayer] = rights
	}

	// Make the move
	g.Board.MovePiece(from, to)

	// Handle pawn promotion (default to Queen)
	if piece.Type == board.Pawn && (to.Row == 0 || to.Row == 7) {
		g.Board.SetPiece(to, board.Piece{Type: board.Queen, Color: g.CurrentPlayer})
	}

	// Update half-move clock
	if piece.Type == board.Pawn || capturedPiece.Type != board.Empty {
		g.halfMoveClock = 0
	} else {
		g.halfMoveClock++
	}

	// Update time control
	if g.TimeControl != nil {
		g.TimeControl.SwitchPlayer(g.CurrentPlayer == board.White)
	}

	// Switch player and update game state
	if g.CurrentPlayer == board.White {
		g.CurrentPlayer = board.Black
	} else {
		g.CurrentPlayer = board.White
		g.fullMoveNumber++
	}

	g.updateGameState()
	return nil
}

// updateGameState updates the state of the game (check, checkmate, etc.)
func (g *Game) updateGameState() {
	// Check if the current player is in check
	inCheck := g.isInCheck(g.CurrentPlayer)

	// Check if the current player has any valid moves
	hasValidMoves := g.hasValidMoves()

	if inCheck && !hasValidMoves {
		g.State = Checkmate
	} else if inCheck {
		g.State = Check
	} else if !hasValidMoves {
		g.State = Stalemate
	} else if g.halfMoveClock >= 50 {
		g.State = Draw // 50-move rule
	} else {
		g.State = InProgress
	}
}

// isInCheck checks if a player is in check
func (g *Game) isInCheck(color board.Color) bool {
	// Find the king
	var kingPos board.Position
	found := false
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			pos := board.Position{Row: row, Col: col}
			piece := g.Board.GetPiece(pos)
			if piece.Type == board.King && piece.Color == color {
				kingPos = pos
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	// Check if any opponent piece can capture the king
	opponentColor := board.White
	if color == board.White {
		opponentColor = board.Black
	}

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			pos := board.Position{Row: row, Col: col}
			piece := g.Board.GetPiece(pos)
			if piece.Type != board.Empty && piece.Color == opponentColor {
				// Temporarily switch player to check move validity
				savedPlayer := g.CurrentPlayer
				g.CurrentPlayer = opponentColor
				isValid := g.isValidPieceMove(pos, kingPos, piece)
				g.CurrentPlayer = savedPlayer
				if isValid {
					return true
				}
			}
		}
	}

	return false
}

// hasValidMoves checks if the current player has any valid moves
func (g *Game) hasValidMoves() bool {
	for fromRow := 0; fromRow < 8; fromRow++ {
		for fromCol := 0; fromCol < 8; fromCol++ {
			fromPos := board.Position{Row: fromRow, Col: fromCol}
			piece := g.Board.GetPiece(fromPos)
			if piece.Type != board.Empty && piece.Color == g.CurrentPlayer {
				for toRow := 0; toRow < 8; toRow++ {
					for toCol := 0; toCol < 8; toCol++ {
						toPos := board.Position{Row: toRow, Col: toCol}
						if g.IsValidMove(fromPos, toPos) {
							// Check if the move would leave the player in check
							if !g.wouldBeInCheck(fromPos, toPos) {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

// wouldBeInCheck checks if a move would leave the player in check
func (g *Game) wouldBeInCheck(from, to board.Position) bool {
	// Save the current state
	origPiece := g.Board.GetPiece(from)
	destPiece := g.Board.GetPiece(to)

	// Make the move temporarily
	g.Board.MovePiece(from, to)

	// Check if the player is in check
	inCheck := g.isInCheck(g.CurrentPlayer)

	// Restore the board
	g.Board.SetPiece(from, origPiece)
	g.Board.SetPiece(to, destPiece)

	return inCheck
}

// GetGameStatus returns a string representation of the game status
func (g *Game) GetGameStatus() string {
	switch g.State {
	case InProgress:
		if g.CurrentPlayer == board.White {
			return "White to move"
		}
		return "Black to move"
	case Check:
		if g.CurrentPlayer == board.White {
			return "White is in check"
		}
		return "Black is in check"
	case Checkmate:
		if g.CurrentPlayer == board.White {
			return "Black wins by checkmate"
		}
		return "White wins by checkmate"
	case Stalemate:
		return "Draw by stalemate"
	case Draw:
		return "Draw"
	case TimeOut:
		return "Time out"
	default:
		return "Unknown game state"
	}
}

// GetTimeLeft returns the formatted time left for the current player
func (g *Game) GetTimeLeft() string {
	if g.TimeControl == nil {
		return ""
	}
	return g.TimeControl.FormatTime(g.CurrentPlayer == board.White)
}

// GetCurrentPlayerName returns the name of the current player
func (g *Game) GetCurrentPlayerName() string {
	if g.CurrentPlayer == board.White {
		return "White"
	}
	return "Black"
}

// Utility function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
