package gomcts

import (
	"math/rand"
	"strconv"
	"fmt"
)

/////////////////////////////////////////////////////////////
// GameState implementation of my made-up game 123ToTen
/////////////////////////////////////////////////////////////

type GameState123ToTen struct {
	Total       int
	SecondPlayersTurn bool
	terminal	bool
	localRand *rand.Rand
}

func NewGameState123ToTen(seed int64) *GameState123ToTen {
	gs := new(GameState123ToTen)
	gs.terminal = false
	gs.SecondPlayersTurn = false
	rndSrc := rand.NewSource(seed)	// Static seed for testing
	gs.localRand = rand.New(rndSrc)
	return gs
}

func (gstate *GameState123ToTen) IsSecondPlayersTurn() bool {
	return gstate.SecondPlayersTurn
}


func (gstate *GameState123ToTen) PossibleMoves() []string {
	var moves []string
	switch {
	case gstate.Total < 8:
		moves = []string{"1", "2", "3"}
	case gstate.Total == 8:
		moves = []string{"1", "2"}
	case gstate.Total == 9:
		moves = []string{"1"}
	default:
		moves = []string{}
	}
	return moves
}


func (gstate *GameState123ToTen) NumberOfMoves() int {
	return len(gstate.PossibleMoves())
}

func (gstate *GameState123ToTen) PossibleMovesShuffled() []string {
	moves := gstate.PossibleMoves()
	for i, _ := range moves {
		j := gstate.localRand.Intn(i + 1)
		moves[i], moves[j] = moves[j], moves[i]
	}
	return moves
}


func (gstate *GameState123ToTen) IsTerminal() bool {
	return gstate.terminal
}

func (gstate *GameState123ToTen) LocalRand() *rand.Rand {
	return gstate.localRand
}

func (gstate *GameState123ToTen) TerminalReward() float64 {
	if !gstate.IsTerminal() {
		panic("reward called but gstate not terminal")
	}
	var reward float64
	if gstate.SecondPlayersTurn {
		reward = 0.0
	} else {
		reward = 1.0
	}
	return reward
}

func (gstate GameState123ToTen) NewGameStateFromMove(move string) GameState {
	cpy := gstate.Copy()
	cpy.DoMove(move)
	return &cpy
}

func (gstate *GameState123ToTen) RewardFromRandomPlayout() float64 {
	cpy := gstate.Copy()
	DoRandomPlayout(&cpy)
	return cpy.TerminalReward()
}

func (gstate *GameState123ToTen) DoMove(move string) {
	moveFromString, _ := strconv.Atoi(move)
	gstate.Total += moveFromString
	if gstate.Total >= 10 {
		gstate.terminal = true
	} else {
		gstate.SecondPlayersTurn = !(gstate.SecondPlayersTurn)
	}
}

func (gstate GameState123ToTen) Copy() GameState123ToTen {
	newCopy := gstate
	return newCopy
}

func (gstate *GameState123ToTen) Summary() string {
	result := fmt.Sprintf("Player %v, Total: %v", gstate.CurrentPlayer(), gstate.Total)
	return result
}

func (gstate *GameState123ToTen) CurrentPlayer() string {
	var result string
	if gstate.SecondPlayersTurn {
		result = "B"
	} else {
		result = "A"
	}
	return result
}