package gomcts

import (
	"fmt"
	"strconv"
)

// GameAddTillX implements a Nim-like two-player game where players take
// turns adding between MinAdd and MaxAdd to a shared total. The first
// player to reach the target total wins.
//
// Game theory: with target T, min M, and max X, positions where
// total % (M + X) == 0 are losing for the player to move.
type GameAddTillX struct {
	total  int
	player int // 0 or 1
	target int
	minAdd int
	maxAdd int
}

// NewGameAddTillX returns the initial state with the given parameters.
func NewGameAddTillX(target, minAdd, maxAdd int) *GameAddTillX {
	return &GameAddTillX{
		target: target,
		minAdd: minAdd,
		maxAdd: maxAdd,
	}
}

func (g *GameAddTillX) Actions() []string {
	remaining := g.target - g.total
	if remaining <= 0 {
		return nil
	}
	maxMove := min(g.maxAdd, remaining)
	actions := make([]string, 0, maxMove-g.minAdd+1)
	for i := g.minAdd; i <= maxMove; i++ {
		actions = append(actions, strconv.Itoa(i))
	}
	return actions
}

func (g *GameAddTillX) NextState(action string) GameState {
	n, _ := strconv.Atoi(action)
	return &GameAddTillX{
		total:  g.total + n,
		player: 1 - g.player,
		target: g.target,
		minAdd: g.minAdd,
		maxAdd: g.maxAdd,
	}
}

func (g *GameAddTillX) IsTerminal() bool {
	return g.total >= g.target
}

func (g *GameAddTillX) Scores() [2]float64 {
	winner := 1 - g.player
	var scores [2]float64
	scores[winner] = 1
	return scores
}

func (g *GameAddTillX) Player() int {
	return g.player
}

// Total returns the current game total.
func (g *GameAddTillX) Total() int {
	return g.total
}

// String returns a human-readable description of the game state.
func (g *GameAddTillX) String() string {
	return fmt.Sprintf("Player %d to move, total: %d (target: %d)", g.player, g.total, g.target)
}
