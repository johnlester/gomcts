package gomcts

import (
	"fmt"
	"strconv"
)

// Game123ToTen implements a simple two-player game where players take turns
// adding 1, 2, or 3 to a shared total starting at 0. The first player to
// reach a total of 10 wins.
//
// Game theory: a player who brings the total to 2 or 6 can force a win,
// since they can always respond to keep the total on the 2-6-10 track.
type Game123ToTen struct {
	total  int
	player int // 0 or 1
}

// NewGame123ToTen returns the initial state for a 123ToTen game.
func NewGame123ToTen() *Game123ToTen {
	return &Game123ToTen{}
}

func (g *Game123ToTen) Actions() []string {
	remaining := 10 - g.total
	if remaining <= 0 {
		return nil
	}
	maxMove := min(3, remaining)
	actions := make([]string, maxMove)
	for i := range maxMove {
		actions[i] = strconv.Itoa(i + 1)
	}
	return actions
}

func (g *Game123ToTen) NextState(action string) GameState {
	n, _ := strconv.Atoi(action)
	return &Game123ToTen{
		total:  g.total + n,
		player: 1 - g.player,
	}
}

func (g *Game123ToTen) IsTerminal() bool {
	return g.total >= 10
}

func (g *Game123ToTen) Scores() [2]float64 {
	// The player whose turn it is did NOT make the winning move;
	// the other player did.
	winner := 1 - g.player
	var scores [2]float64
	scores[winner] = 1
	return scores
}

func (g *Game123ToTen) Player() int {
	return g.player
}

// Total returns the current game total.
func (g *Game123ToTen) Total() int {
	return g.total
}

// String returns a human-readable description of the game state.
func (g *Game123ToTen) String() string {
	return fmt.Sprintf("Player %d to move, total: %d", g.player, g.total)
}
