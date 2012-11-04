package gomcts

import (
	"testing"
	"time"
)

/////////////////////////////////////////////////////////////
// 123ToTen Game tests
/////////////////////////////////////////////////////////////

func TestGame123ToTen_Moves(t *testing.T) {
	gs := NewGameState123ToTen(42)
	if gs.NumberOfMoves() != 3 {
		t.Errorf("123ToTen GameState should have 3 possible moves, but instead has %v", gs.NumberOfMoves())
	}
	moves := gs.PossibleMoves()
	if len(moves) != 3 {
		t.Errorf("123ToTen GameState should have 3 possible moves, but instead has %v", len(moves))
	}
}

func TestGame123ToTen_Copying(t *testing.T) {
	gs := NewGameState123ToTen(42)
	cpy := gs.Copy()
	cpy.Total++
	if cpy.Total != 1 {
		t.Errorf("cpy.Total should be 1 but is %v", cpy.Total)
	}
	if gs.Total != 0 {
		t.Errorf("gs.Total should still be 0 but is %v", gs.Total)
	}
}

func TestGame123ToTen_DoMove(t *testing.T) {
	gs := NewGameState123ToTen(42)
	gs.DoMove("1")
	if gs.Total != 1 {
		t.Errorf("Total should be 1 but is %v", gs.Total)
	}
	if gs.SecondPlayersTurn != true {
		t.Errorf("Should be Player B's move now")
	}
}

func TestGame123ToTen_RandomPlayout(t *testing.T) {
	gs := NewGameState123ToTen(4)
	if gs.Total != 0 {
		t.Errorf("Starting total should be 0 but is %v", gs.Total)
	}
	if gs.NumberOfMoves() != 3 {
		t.Errorf("123ToTen GameState should have 3 possible moves, but instead has %v", gs.NumberOfMoves())
	}
	DoRandomPlayout(gs)
	if gs.Total < 10 {
		t.Errorf("Total should be 10 but is %v", gs.Total)
	}
	if !(gs.IsTerminal()) {
		t.Errorf("Should be terminal")
	}
}


/////////////////////////////////////////////////////////////
// 123ToTen MCTS Tests
/////////////////////////////////////////////////////////////

func TestGame123ToTenGoMCTS_RootNode(t *testing.T) {
	gs := NewGameState123ToTen(42)
	rn := NewNode(gs, nil, "")
	if rn.VisitCount != 0.0 {
		t.Errorf("Root node at creation should have VisitCount of 0.0")
	}
	if rn.NextMoveToTry != 0 {
		t.Errorf("Root node at creation should have NextMoveToTry of 0")
	}
	if len(rn.Children) != 3 {
		t.Errorf("Root node at creation should have 3 Children")
	}}

func TestGame123ToTenGoMCTS_OneManualIteration(t *testing.T) {
	gs := NewGameState123ToTen(42)
	rn := NewNode(gs, nil, "")
	child := rn.NewChild()
	move := child.GeneratingMove
	if rn.NextMoveToTry != 1 {
		t.Errorf("Root node at 1 iteration should have NextMoveToTry of 1")
	}
	if (move != "1") && (move != "2") && (move != "3") {
		t.Errorf("Move is something wrong")
	}
	if rn.Children[0] != child {
		t.Errorf("Something is wrong")
	}
	if child.Parent != rn {
		t.Errorf("Something is wrong")
	}
	if rn.Parent != nil {
		t.Errorf("Something is wrong")
	}
	if child.State.IsTerminal() {
		t.Errorf("Something is wrong")
	}
	if child.State.IsSecondPlayersTurn() != true {
		t.Errorf("Something is wrong")
	}
	// Now do random playout
	reward := child.State.RewardFromRandomPlayout()
	if reward[0] + reward[1] != 1.0 {
		t.Errorf("Something is wrong")
	}
	// Now propagate reward from child to parent (root node in this case)
	child.backpropagateReward(reward)
	if rn.VisitCount != 1.0 {
		t.Errorf("Root node should have VisitCount of 1.0, not %v", rn.VisitCount)
	}
	if child.CumulativeScore != reward {
		t.Errorf("Child node should have CumulativeScore of %v, not %v", reward, child.CumulativeScore)
	}
	if child.VisitCount != 1.0 {
		t.Errorf("Child node should have VisitCount of 1.0, not %v", child.VisitCount)
	}
}

func TestGame123ToTenGoMCTS_DoOneIteration(t *testing.T) {
	gs := NewGameState123ToTen(42)
	rn := NewNode(gs, nil, "")
	rn.DoNIterations(1)
	if rn.VisitCount != 1.0 {
		t.Errorf("Root node at 1 iteration should have VisitCount of 1.0")
	}
	if rn.NextMoveToTry != 1 {
		t.Errorf("Root node at 1 iteration should have NextMoveToTry of 1")
	}
	// if len(rn.Children) != 1 {
	// 	t.Errorf("Root node at 1 iteration should have 1 child")
	// }
	if rn.Children[0].State.IsSecondPlayersTurn() != true {
		t.Errorf("State at child node after 1 iteration should have be second players turn")
	}
}

func TestGame123ToTenGoMCTS_DoManyIterationsEndgame(t *testing.T) {
	iters := 10000
	gs := NewGameState123ToTen(int64(time.Now().Nanosecond()))
	gs.Total = 0
	rn := NewNode(gs, nil, "")
	rn.DoNIterations(iters)
	if rn.State.IsSecondPlayersTurn() != false {
		t.Errorf("Something is wrong")
	}
	if rn.VisitCount != float64(iters) {
		t.Errorf("Root node at %v iterations should have VisitCount of %v", iters, iters)
	}
	if rn.NextMoveToTry != rn.NumberOfChildren {
		t.Errorf("Root node at %v iterations should have NextMoveToTry of %v, not %v", iters, rn.NumberOfChildren, rn.NextMoveToTry)
	}
}


// func TestGame123ToTenGoMCTS_DoManyIterations(t *testing.T) {
// 	iters := 10000
// 	gs := NewGameState123ToTen(int64(time.Now().Nanosecond()))
// 	rn := NewNode(gs, nil, "")
// 	rn.DoNIterations(iters)
// 	if rn.VisitCount != float64(iters) {
// 		t.Errorf("Root node at %v iterations should have VisitCount of %v", iters, iters)
// 	}
// 	if rn.NextMoveToTry != 3 {
// 		t.Errorf("Root node at %v iterations should have NextMoveToTry of 3", iters)
// 	}
// 	move := rn.BestChild().GeneratingMove
// 	t.Logf("Root node's best move is %v", move)
// 	t.Logf("%v", rn.Summary())
// }
