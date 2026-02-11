package gomcts

import (
	"math/rand/v2"
	"testing"
)

// --- Game123ToTen unit tests ---

func TestGame123ToTen_Actions(t *testing.T) {
	g := NewGame123ToTen()
	if got := len(g.Actions()); got != 3 {
		t.Errorf("expected 3 actions at start, got %d", got)
	}
}

func TestGame123ToTen_ActionsNearEnd(t *testing.T) {
	tests := []struct {
		total    int
		expected int
	}{
		{0, 3}, {5, 3}, {7, 3},
		{8, 2}, {9, 1}, {10, 0},
	}
	for _, tt := range tests {
		g := &Game123ToTen{total: tt.total}
		if got := len(g.Actions()); got != tt.expected {
			t.Errorf("total %d: expected %d actions, got %d", tt.total, tt.expected, got)
		}
	}
}

func TestGame123ToTen_NextState(t *testing.T) {
	g := NewGame123ToTen()
	next := g.NextState("2").(*Game123ToTen)
	if next.Total() != 2 {
		t.Errorf("expected total 2, got %d", next.Total())
	}
	if next.Player() != 1 {
		t.Errorf("expected player 1, got %d", next.Player())
	}
	// Original must be unchanged.
	if g.Total() != 0 {
		t.Errorf("original state mutated: total %d", g.Total())
	}
	if g.Player() != 0 {
		t.Errorf("original state mutated: player %d", g.Player())
	}
}

func TestGame123ToTen_Terminal(t *testing.T) {
	// Player 0 moved to total 10, now it's player 1's turn at a terminal state.
	g := &Game123ToTen{total: 10, player: 1}
	if !g.IsTerminal() {
		t.Error("expected terminal at total 10")
	}
	scores := g.Scores()
	if scores != [2]float64{1, 0} {
		t.Errorf("expected player 0 wins [1 0], got %v", scores)
	}
}

func TestGame123ToTen_TerminalPlayer1Wins(t *testing.T) {
	// Player 1 moved to total 10, now it's player 0's turn at terminal.
	g := &Game123ToTen{total: 10, player: 0}
	scores := g.Scores()
	if scores != [2]float64{0, 1} {
		t.Errorf("expected player 1 wins [0 1], got %v", scores)
	}
}

// --- GameAddTillX unit tests ---

func TestGameAddTillX_Actions(t *testing.T) {
	g := NewGameAddTillX(31, 1, 4)
	if got := len(g.Actions()); got != 4 {
		t.Errorf("expected 4 actions at start, got %d", got)
	}
}

func TestGameAddTillX_ActionsNearEnd(t *testing.T) {
	g := &GameAddTillX{total: 29, target: 31, minAdd: 1, maxAdd: 4}
	if got := len(g.Actions()); got != 2 {
		t.Errorf("expected 2 actions at total 29, got %d", got)
	}
}

func TestGameAddTillX_NextState(t *testing.T) {
	g := NewGameAddTillX(31, 1, 4)
	next := g.NextState("3").(*GameAddTillX)
	if next.Total() != 3 {
		t.Errorf("expected total 3, got %d", next.Total())
	}
	if next.Player() != 1 {
		t.Errorf("expected player 1, got %d", next.Player())
	}
	if g.Total() != 0 {
		t.Errorf("original mutated")
	}
}

// --- MCTS core tests ---

func TestMCTS_SingleIteration(t *testing.T) {
	g := NewGame123ToTen()
	m := New(g, WithRand(rand.New(rand.NewPCG(42, 0))))
	m.Search(1)
	if m.root.visits != 1 {
		t.Errorf("expected 1 root visit, got %.0f", m.root.visits)
	}
	if len(m.root.children) != 1 {
		t.Errorf("expected 1 child after 1 iteration, got %d", len(m.root.children))
	}
}

func TestMCTS_FullExpansion(t *testing.T) {
	g := NewGame123ToTen()
	m := New(g, WithRand(rand.New(rand.NewPCG(42, 0))))
	m.Search(100)
	if len(m.root.children) != 3 {
		t.Errorf("expected all 3 children expanded after 100 iterations, got %d", len(m.root.children))
	}
}

func TestMCTS_Backpropagation(t *testing.T) {
	g := NewGame123ToTen()
	m := New(g, WithRand(rand.New(rand.NewPCG(42, 0))))
	m.Search(50)
	if m.root.visits != 50 {
		t.Errorf("expected 50 root visits, got %.0f", m.root.visits)
	}
	var childVisits float64
	for _, c := range m.root.children {
		childVisits += c.visits
	}
	if childVisits != 50 {
		t.Errorf("child visits sum %.0f != root visits 50", childVisits)
	}
}

func TestMCTS_CumulativeSearch(t *testing.T) {
	g := NewGame123ToTen()
	m := New(g, WithRand(rand.New(rand.NewPCG(42, 0))))
	m.Search(100)
	m.Search(100)
	if m.root.visits != 200 {
		t.Errorf("expected 200 cumulative visits, got %.0f", m.root.visits)
	}
}

func TestMCTS_Stats(t *testing.T) {
	g := NewGame123ToTen()
	m := New(g, WithRand(rand.New(rand.NewPCG(42, 0))))
	m.Search(100)
	stats := m.Stats()
	if stats == "" {
		t.Error("expected non-empty stats")
	}
}

func TestMCTS_BestActionEmpty(t *testing.T) {
	// Terminal state has no actions.
	g := &Game123ToTen{total: 10, player: 1}
	m := New(g, WithRand(rand.New(rand.NewPCG(42, 0))))
	if got := m.BestAction(10); got != "" {
		t.Errorf("expected empty action for terminal state, got %q", got)
	}
}

// --- MCTS optimality tests ---

func TestMCTS_123ToTen_FindsOptimalOpening(t *testing.T) {
	// From total 0, player 0 should play "2" to reach total 2
	// (a losing position for the opponent).
	g := NewGame123ToTen()
	m := New(g, WithRand(rand.New(rand.NewPCG(42, 0))))
	action := m.BestAction(50000)
	if action != "2" {
		t.Errorf("expected optimal opening move '2', got '%s'\n%s", action, m.Stats())
	}
}

func TestMCTS_123ToTen_FindsWinningMoveAt7(t *testing.T) {
	// From total 7, current player can win immediately with "3" → 10.
	g := &Game123ToTen{total: 7}
	m := New(g, WithRand(rand.New(rand.NewPCG(42, 0))))
	action := m.BestAction(10000)
	if action != "3" {
		t.Errorf("from total 7, expected '3', got '%s'", action)
	}
}

func TestMCTS_123ToTen_FindsWinningMoveAt8(t *testing.T) {
	g := &Game123ToTen{total: 8}
	m := New(g, WithRand(rand.New(rand.NewPCG(42, 0))))
	action := m.BestAction(10000)
	if action != "2" {
		t.Errorf("from total 8, expected '2', got '%s'", action)
	}
}

func TestMCTS_123ToTen_FindsWinningMoveAt9(t *testing.T) {
	g := &Game123ToTen{total: 9}
	m := New(g, WithRand(rand.New(rand.NewPCG(42, 0))))
	action := m.BestAction(5000)
	if action != "1" {
		t.Errorf("from total 9, expected '1', got '%s'", action)
	}
}

func TestMCTS_AddTillX_Completion(t *testing.T) {
	g := NewGameAddTillX(31, 1, 4)
	m := New(g, WithRand(rand.New(rand.NewPCG(42, 0))))
	m.Search(10000)
	if m.root.visits != 10000 {
		t.Errorf("expected 10000 visits, got %.0f", m.root.visits)
	}
}

func TestMCTS_AddTillX_FindsWinningMoveAt27(t *testing.T) {
	// From total 27 (target 31), current player wins with "4" → 31.
	g := &GameAddTillX{total: 27, target: 31, minAdd: 1, maxAdd: 4}
	m := New(g, WithRand(rand.New(rand.NewPCG(42, 0))))
	action := m.BestAction(5000)
	if action != "4" {
		t.Errorf("from total 27, expected '4', got '%s'", action)
	}
}

func TestMCTS_AddTillX_FindsWinningMoveAt22(t *testing.T) {
	// From total 22 (target 31, min 1, max 4), current player should
	// play "4" to reach 26 (a losing position for the opponent).
	// Losing positions are where (target - total) % (min + max) == 0,
	// i.e., totals 1, 6, 11, 16, 21, 26, 31. The opponent at 26
	// can only reach 27-30, and from any of those the current player
	// reaches 31 and wins.
	g := &GameAddTillX{total: 22, target: 31, minAdd: 1, maxAdd: 4}
	m := New(g, WithRand(rand.New(rand.NewPCG(42, 0))))
	action := m.BestAction(20000)
	if action != "4" {
		t.Errorf("from total 22, expected '4' (→26), got '%s'\n%s", action, m.Stats())
	}
}

// --- Interface satisfaction ---

func TestGame123ToTen_ImplementsGameState(t *testing.T) {
	var _ GameState = (*Game123ToTen)(nil)
}

func TestGameAddTillX_ImplementsGameState(t *testing.T) {
	var _ GameState = (*GameAddTillX)(nil)
}

// --- WithExplorationFactor option ---

func TestMCTS_WithExplorationFactor(t *testing.T) {
	g := NewGame123ToTen()
	m := New(g, WithExplorationFactor(2.0), WithRand(rand.New(rand.NewPCG(42, 0))))
	if m.explorationFactor != 2.0 {
		t.Errorf("expected exploration factor 2.0, got %f", m.explorationFactor)
	}
	// Should still produce a valid result.
	action := m.BestAction(1000)
	if action == "" {
		t.Error("expected a non-empty action")
	}
}
