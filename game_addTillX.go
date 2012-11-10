package gomcts

import (
	"math/rand"
	"strconv"
	"fmt"
)

/////////////////////////////////////////////////////////////
// GameState implementation of my add till X game
/////////////////////////////////////////////////////////////

const (
	TargetTotal = 31
	MinMove = 1
	MaxMove = 4
)

type GameAddTillX struct {
	Total       int
	SecondPlayersTurn bool
	terminal	bool
	localRand *rand.Rand
}

func NewGameStateAddTillX(seed int64) *GameAddTillX {
	gs := new(GameAddTillX)
	gs.terminal = false
	gs.SecondPlayersTurn = false
	rndSrc := rand.NewSource(seed)	// Static seed for testing
	gs.localRand = rand.New(rndSrc)
	return gs
}

func (gstate *GameAddTillX) IsSecondPlayersTurn() bool {
	return gstate.SecondPlayersTurn
}


func (gstate *GameAddTillX) PossibleMoves() []string {
	maxCurrentMove := TargetTotal - gstate.Total
	if maxCurrentMove > MaxMove {
		maxCurrentMove = MaxMove
	}
	moves := make([]string, maxCurrentMove - MinMove + 1)
	for i := MinMove; i <= maxCurrentMove; i++ {
		moves[i-MinMove] = strconv.Itoa(i)
	} 
	return moves
}


func (gstate *GameAddTillX) NumberOfMoves() int {
	return len(gstate.PossibleMoves())
}

func (gstate *GameAddTillX) PossibleMovesShuffled() []string {
	moves := gstate.PossibleMoves()
	for i, _ := range moves {
		j := gstate.localRand.Intn(i + 1)
		moves[i], moves[j] = moves[j], moves[i]
	}
	return moves
}


func (gstate *GameAddTillX) IsTerminal() bool {
	return gstate.terminal
}

func (gstate *GameAddTillX) LocalRand() *rand.Rand {
	return gstate.localRand
}

func (gstate *GameAddTillX) TerminalReward() [2]float64 {
	if !gstate.IsTerminal() {
		panic("reward called but gstate not terminal")
	}
	var reward [2]float64
	if gstate.SecondPlayersTurn {
		// Finished game state is second player's turn, which means first player just did winning move
		reward[0] = 1.0
		reward[1] = 0.0
	} else {
		reward[0] = 0.0
		reward[1] = 1.0
	}
	return reward
}

func (gstate GameAddTillX) NewGameStateFromMove(move string) GameState {
	cpy := gstate.Copy()
	cpy.DoMove(move)
	return &cpy
}

func (gstate *GameAddTillX) RewardFromRandomPlayout() [2]float64 {
	cpy := gstate.Copy()
	DoRandomPlayout(&cpy)
	return cpy.TerminalReward()
}

func (gstate *GameAddTillX) DoMove(move string) {
	moveFromString, _ := strconv.Atoi(move)
	gstate.Total += moveFromString
	if gstate.Total >= TargetTotal {
		gstate.terminal = true
	} 
	gstate.SecondPlayersTurn = !(gstate.SecondPlayersTurn)
}

func (gstate GameAddTillX) Copy() GameAddTillX {
	newCopy := gstate
	return newCopy
}

func (gstate *GameAddTillX) Summary() string {
	result := fmt.Sprintf("Player %v, Total: %v", gstate.CurrentPlayer(), gstate.Total)
	return result
}

func (gstate *GameAddTillX) CurrentPlayer() string {
	var result string
	if gstate.SecondPlayersTurn {
		result = "B"
	} else {
		result = "A"
	}
	return result
}