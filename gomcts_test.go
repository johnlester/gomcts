package gomcts

import (
	"testing"
)

/////////////////////////////////////////////////////////////
// 123ToTen Game tests
/////////////////////////////////////////////////////////////

func TestGame123ToTenMoves(t *testing.T) {
	gs := new(GameState123ToTen)
	moves := gs.PossibleMoves()
	if len(moves) != 3 {
		t.Errorf("123ToTen GameState should have 3 possible moves, but instead has %v", len(moves))
	}
}

func TestGame123ToTenCopying(t *testing.T) {
	gs := new(GameState123ToTen)
	cpy := gs.Copy()
	cpy.Total++
	if cpy.Total != 1 {
		t.Errorf("cpy.Total should be 1 but is %v", cpy.Total)
	}
	if gs.Total != 0 {
		t.Errorf("gs.Total should still be 0 but is %v", gs.Total)
	}
}

func TestGame123ToTenDoMove(t *testing.T) {
	gs := new(GameState123ToTen)
	gs.DoMove("1")
	if gs.Total != 1 {
		t.Errorf("Total should be 1 but is %v", gs.Total)
	}
	if gs.PlayerBMove != true {
		t.Errorf("Should be Player B's move now")
	}
}
