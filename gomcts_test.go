package gomcts

import (
	"testing"
)

/////////////////////////////////////////////////////////////
// 123ToTen Game tests
/////////////////////////////////////////////////////////////

func TestGame123ToTen(t *testing.T) {
	gs := new(GameState123ToTen)
	moves := gs.PossibleMoves()
	if len(moves) != 3 {
		t.Errorf("123ToTen GameState should have 3 possible moves, but instead has %v", len(moves))
	}

}
