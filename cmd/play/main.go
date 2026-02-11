// Command play lets a human play 123ToTen against an MCTS-powered AI.
//
// Usage:
//
//	go run ./cmd/play
package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/johnlester/gomcts"
)

const aiIterations = 50000

func main() {
	fmt.Println("=== 123ToTen ===")
	fmt.Println("Take turns adding 1, 2, or 3 to reach 10. First to 10 wins!")
	fmt.Println()

	rng := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	scanner := bufio.NewScanner(os.Stdin)
	var state gomcts.GameState = gomcts.NewGame123ToTen()

	for !state.IsTerminal() {
		g := state.(*gomcts.Game123ToTen)
		fmt.Printf("Total: %d | ", g.Total())

		if state.Player() == 0 {
			fmt.Printf("Your turn. Choose %v: ", state.Actions())
			if !scanner.Scan() {
				fmt.Println("\nGoodbye!")
				return
			}
			move := strings.TrimSpace(scanner.Text())
			valid := false
			for _, a := range state.Actions() {
				if a == move {
					valid = true
					break
				}
			}
			if !valid {
				fmt.Println("Invalid move, try again.")
				continue
			}
			state = state.NextState(move)
		} else {
			m := gomcts.New(state, gomcts.WithRand(rng))
			move := m.BestAction(aiIterations)
			fmt.Printf("AI plays: %s\n", move)
			state = state.NextState(move)
		}
	}

	g := state.(*gomcts.Game123ToTen)
	fmt.Printf("Total: %d | ", g.Total())
	scores := state.Scores()
	if scores[0] > scores[1] {
		fmt.Println("You win!")
	} else {
		fmt.Println("AI wins!")
	}
}
