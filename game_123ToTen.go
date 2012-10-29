package gomcts

import (
	"math/rand"
	"strconv"
)

/////////////////////////////////////////////////////////////
// GameState implementation of my made-up game 123ToTen
/////////////////////////////////////////////////////////////

type GameState123ToTen struct {
	Total       int
	SecondPlayersTurn bool
	Finished	bool
	LocalRand *rand.Rand
}

func NewGameState123ToTen() *GameState123ToTen {
	gs := new(GameState123ToTen)
	rndSrc := rand.NewSource(43)	// Static seed for testing
	gs.LocalRand = rand.New(rndSrc)
	return gs
}

func (gstate GameState123ToTen) IsSecondPlayersTurn() bool {
	return gstate.SecondPlayersTurn
}


func (gstate GameState123ToTen) PossibleMoves() []string {
	moves := []string{"1", "2", "3"}
	return moves
}

func (gstate GameState123ToTen) PossibleMovesShuffled() []string {
	moves := gstate.PossibleMoves()
	for i := range moves {
		j := gstate.LocalRand.Intn(i + 1)
		moves[i], moves[j] = moves[j], moves[i]
	}
	return moves
}


func (gstate GameState123ToTen) IsNotTerminal() bool {
	return !gstate.Finished
}

func (gstate GameState123ToTen) TerminalReward() float64 {
	// error if not terminal?
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
	return cpy
}

func (gstate GameState123ToTen) RewardFromRandomPlayout() float64 {
	cpy := gstate.Copy()
	DoRandomPlayout(&cpy)
	return cpy.TerminalReward()
}

func (gstate *GameState123ToTen) DoMove(move string) {
	moveFromString, _ := strconv.Atoi(move)
	gstate.Total += moveFromString
	if gstate.Total >= 10 {
		gstate.Finished = true
	} else {
		gstate.SecondPlayersTurn = !gstate.SecondPlayersTurn
	}
}

func (gstate GameState123ToTen) Copy() GameState123ToTen {
	newCopy := gstate
	return newCopy
}

func DoRandomPlayout(gstate *GameState123ToTen) {	
	for gstate.IsNotTerminal() {
		moveCount := len(gstate.PossibleMoves())
		rndIdx := gstate.LocalRand.Intn(moveCount)
		gstate.DoMove(gstate.PossibleMoves()[rndIdx])
	}
}
