package game

import (
	"fmt"
	"time"
)

// TimeControl represents the time control settings for a game
type TimeControl struct {
	InitialTime      time.Duration // Initial time per player
	IncrementPerMove time.Duration // Time added after each move
	WhiteTimeLeft    time.Duration
	BlackTimeLeft    time.Duration
	lastMoveTime     time.Time
	isRunning        bool
}

// NewTimeControl creates a new time control with the specified initial time and increment
func NewTimeControl(initialMinutes int, incrementSeconds int) *TimeControl {
	initialTime := time.Duration(initialMinutes) * time.Minute
	increment := time.Duration(incrementSeconds) * time.Second

	return &TimeControl{
		InitialTime:      initialTime,
		IncrementPerMove: increment,
		WhiteTimeLeft:    initialTime,
		BlackTimeLeft:    initialTime,
		isRunning:        false,
	}
}

// Start starts the timer for the current player
func (tc *TimeControl) Start() {
	tc.lastMoveTime = time.Now()
	tc.isRunning = true
}

// Stop stops the timer and returns the elapsed time
func (tc *TimeControl) Stop() time.Duration {
	if !tc.isRunning {
		return 0
	}

	elapsed := time.Since(tc.lastMoveTime)
	tc.isRunning = false
	return elapsed
}

// SwitchPlayer switches the active timer and adds the increment
func (tc *TimeControl) SwitchPlayer(isWhite bool) {
	elapsed := tc.Stop()

	// Subtract elapsed time from the previous player
	if isWhite {
		tc.WhiteTimeLeft -= elapsed
	} else {
		tc.BlackTimeLeft -= elapsed
	}

	// Add increment to the previous player
	if tc.IncrementPerMove > 0 {
		if isWhite {
			tc.WhiteTimeLeft += tc.IncrementPerMove
		} else {
			tc.BlackTimeLeft += tc.IncrementPerMove
		}
	}

	tc.Start()
}

// IsTimeUp checks if a player has run out of time
func (tc *TimeControl) IsTimeUp(isWhite bool) bool {
	if isWhite {
		return tc.WhiteTimeLeft <= 0
	}
	return tc.BlackTimeLeft <= 0
}

// FormatTime formats the remaining time in a human-readable format
func (tc *TimeControl) FormatTime(isWhite bool) string {
	var timeLeft time.Duration
	if isWhite {
		timeLeft = tc.WhiteTimeLeft
	} else {
		timeLeft = tc.BlackTimeLeft
	}

	minutes := int(timeLeft.Minutes())
	seconds := int(timeLeft.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
