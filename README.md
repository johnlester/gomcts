# gomcts

A Monte Carlo Tree Search (MCTS) library for two-player perfect-information
games in Go.

## Usage

Implement the `GameState` interface for your game:

```go
type GameState interface {
    Actions() []string
    NextState(action string) GameState
    IsTerminal() bool
    Scores() [2]float64
    Player() int
}
```

Then run MCTS to find the best move:

```go
game := NewMyGame()
m := gomcts.New(game)
bestMove := m.BestAction(10000)
```

### Options

```go
// Tune the exploration/exploitation tradeoff (default: 1.0)
gomcts.New(game, gomcts.WithExplorationFactor(1.4))

// Provide a seeded RNG for reproducibility
gomcts.New(game, gomcts.WithRand(rand.New(rand.NewPCG(42, 0))))
```

### Custom rollout policy

If your game benefits from a domain-specific simulation strategy, implement
`RolloutPolicy` on your `GameState`:

```go
func (g *MyGame) Rollout() [2]float64 {
    // your custom simulation logic
}
```

## Included games

Two example games are included as reference implementations:

- **Game123ToTen**: Players add 1, 2, or 3 to a total starting at 0.
  First to 10 wins. Positions 2 and 6 are winning for the player who
  just moved.

- **GameAddTillX**: A parameterizable Nim variant. Players add between
  `minAdd` and `maxAdd` to reach a target total.

## Interactive play

Play 123ToTen against the MCTS AI:

```
go run ./cmd/play
```

## Testing

```
go test ./...
```
