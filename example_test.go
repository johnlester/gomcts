package gomcts_test

import (
	"fmt"
	"math/rand/v2"

	"github.com/johnlester/gomcts"
)

func Example() {
	game := gomcts.NewGame123ToTen()
	m := gomcts.New(game, gomcts.WithRand(rand.New(rand.NewPCG(42, 0))))
	best := m.BestAction(10000)
	fmt.Println("Best opening move:", best)
	// Output: Best opening move: 2
}

func Example_fullGame() {
	rng := rand.New(rand.NewPCG(99, 0))
	var state gomcts.GameState = gomcts.NewGame123ToTen()

	for !state.IsTerminal() {
		m := gomcts.New(state, gomcts.WithRand(rng))
		action := m.BestAction(5000)
		fmt.Printf("Player %d plays %s\n", state.Player(), action)
		state = state.NextState(action)
	}

	scores := state.Scores()
	if scores[0] > scores[1] {
		fmt.Println("Player 0 wins!")
	} else {
		fmt.Println("Player 1 wins!")
	}
	// Output:
	// Player 0 plays 2
	// Player 1 plays 1
	// Player 0 plays 3
	// Player 1 plays 1
	// Player 0 plays 3
	// Player 0 wins!
}
