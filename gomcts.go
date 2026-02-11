// Package gomcts implements Monte Carlo Tree Search for two-player
// perfect-information games.
//
// To use this package, implement the [GameState] interface for your game,
// then call [New] to create a search tree and [MCTS.BestAction] to find
// the best move.
package gomcts

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
)

// DefaultExplorationFactor is the default UCB1 exploration constant.
// Higher values favor exploration of less-visited nodes; lower values
// favor exploitation of high-scoring nodes.
const DefaultExplorationFactor = 1.0

// GameState represents the state of a two-player perfect-information game.
type GameState interface {
	// Actions returns the available actions from this state.
	// Returns nil or an empty slice if the game is over.
	Actions() []string

	// NextState returns a new GameState resulting from taking the given action.
	// The receiver must not be modified.
	NextState(action string) GameState

	// IsTerminal reports whether the game is over.
	IsTerminal() bool

	// Scores returns the terminal rewards as [player0, player1].
	// Each score should be in the range [0, 1].
	// This method is only called when IsTerminal returns true.
	Scores() [2]float64

	// Player returns the index (0 or 1) of the player whose turn it is.
	Player() int
}

// RolloutPolicy is an optional interface that a [GameState] can implement
// to provide a custom simulation strategy. If a GameState does not
// implement RolloutPolicy, uniform random rollouts are used.
type RolloutPolicy interface {
	// Rollout simulates a random game from this state to completion
	// and returns the terminal scores as [player0, player1].
	Rollout() [2]float64
}

// MCTS holds the search tree and configuration for Monte Carlo Tree Search.
type MCTS struct {
	root              *node
	explorationFactor float64
	rng               *rand.Rand
}

// Option configures an [MCTS] instance.
type Option func(*MCTS)

// WithExplorationFactor sets the UCB1 exploration constant.
// Higher values encourage exploration; lower values favor exploitation.
// The default is [DefaultExplorationFactor].
func WithExplorationFactor(c float64) Option {
	return func(m *MCTS) { m.explorationFactor = c }
}

// WithRand sets the random number generator used for move shuffling,
// tie-breaking, and default rollouts.
func WithRand(r *rand.Rand) Option {
	return func(m *MCTS) { m.rng = r }
}

// New creates a new MCTS search tree rooted at the given game state.
func New(state GameState, opts ...Option) *MCTS {
	m := &MCTS{
		explorationFactor: DefaultExplorationFactor,
		rng:               rand.New(rand.NewPCG(0, 0)),
	}
	for _, opt := range opts {
		opt(m)
	}
	m.root = m.newNode(state, nil, "")
	return m
}

// BestAction runs the specified number of MCTS iterations and returns
// the action with the highest average score. Returns an empty string
// if no actions are available.
func (m *MCTS) BestAction(iterations int) string {
	m.Search(iterations)
	best := m.root.bestChild(m.rng)
	if best == nil {
		return ""
	}
	return best.action
}

// Search runs the specified number of MCTS iterations on the tree.
// Multiple calls to Search are cumulative.
func (m *MCTS) Search(iterations int) {
	for range iterations {
		leaf := m.selectAndExpand()
		scores := m.rollout(leaf.state)
		leaf.backpropagate(scores)
	}
}

// Stats returns a human-readable summary of the current search tree.
func (m *MCTS) Stats() string {
	r := m.root
	var b strings.Builder
	fmt.Fprintf(&b, "Iterations: %.0f\n", r.visits)
	totalActions := len(r.children) + len(r.unexpanded)
	fmt.Fprintf(&b, "Children explored: %d/%d\n", len(r.children), totalActions)
	fmt.Fprintf(&b, "Tree size: %d nodes\n", r.treeSize())
	for i, child := range r.children {
		fmt.Fprintf(&b, "  %d) Action %s: %.0f visits, avg %.4f, %d descendants\n",
			i, child.action, child.visits, child.averageScore(), child.treeSize())
	}
	return b.String()
}

// node is an internal MCTS tree node.
type node struct {
	state      GameState
	parent     *node
	action     string   // action from parent that led to this node
	children   []*node
	unexpanded []string // actions not yet expanded (consumed from end)
	visits     float64
	scores     [2]float64
}

// newNode creates a tree node with shuffled unexpanded actions.
func (m *MCTS) newNode(state GameState, parent *node, action string) *node {
	actions := append([]string(nil), state.Actions()...)
	m.rng.Shuffle(len(actions), func(i, j int) {
		actions[i], actions[j] = actions[j], actions[i]
	})
	return &node{
		state:      state,
		parent:     parent,
		action:     action,
		unexpanded: actions,
	}
}

// selectAndExpand traverses the tree using UCB1 until it finds a node
// that can be expanded, then expands it. Returns the new leaf node,
// or a terminal node if the selection reaches a game-over state.
func (m *MCTS) selectAndExpand() *node {
	n := m.root
	for !n.state.IsTerminal() {
		if len(n.unexpanded) > 0 {
			return m.expand(n)
		}
		n = n.childToExplore(m.explorationFactor, m.rng)
	}
	return n
}

// expand creates a new child node from the next unexpanded action.
func (m *MCTS) expand(n *node) *node {
	// Pop the last unexpanded action.
	idx := len(n.unexpanded) - 1
	action := n.unexpanded[idx]
	n.unexpanded = n.unexpanded[:idx]

	child := m.newNode(n.state.NextState(action), n, action)
	n.children = append(n.children, child)
	return child
}

// rollout simulates a random game from the given state and returns the
// terminal scores. Uses the state's [RolloutPolicy] if available,
// otherwise performs a uniform random rollout.
func (m *MCTS) rollout(state GameState) [2]float64 {
	if rp, ok := state.(RolloutPolicy); ok {
		return rp.Rollout()
	}
	s := state
	for !s.IsTerminal() {
		actions := s.Actions()
		s = s.NextState(actions[m.rng.IntN(len(actions))])
	}
	return s.Scores()
}

// backpropagate updates visit counts and scores from this node up to the root.
func (n *node) backpropagate(scores [2]float64) {
	for cur := n; cur != nil; cur = cur.parent {
		cur.visits++
		cur.scores[0] += scores[0]
		cur.scores[1] += scores[1]
	}
}

// averageScore returns the average score from the perspective of the parent
// player (the one who chose to move to this node).
func (n *node) averageScore() float64 {
	if n.visits == 0 {
		return 0
	}
	parentPlayer := 1 - n.state.Player()
	return n.scores[parentPlayer] / n.visits
}

// ucb1 computes the Upper Confidence Bound for Trees value for this node.
func (n *node) ucb1(explorationFactor float64) float64 {
	if n.visits == 0 {
		return math.Inf(1)
	}
	return n.averageScore() + explorationFactor*math.Sqrt(2*math.Log(n.parent.visits)/n.visits)
}

// bestChild returns the child with the highest average score.
// Ties are broken randomly using the provided RNG.
func (n *node) bestChild(rng *rand.Rand) *node {
	if len(n.children) == 0 {
		return nil
	}
	best := math.Inf(-1)
	var candidates []*node
	for _, child := range n.children {
		avg := child.averageScore()
		if avg > best {
			best = avg
			candidates = candidates[:0]
			candidates = append(candidates, child)
		} else if avg == best {
			candidates = append(candidates, child)
		}
	}
	return candidates[rng.IntN(len(candidates))]
}

// childToExplore returns the child with the highest UCB1 value.
// Ties are broken randomly using the provided RNG.
func (n *node) childToExplore(explorationFactor float64, rng *rand.Rand) *node {
	best := math.Inf(-1)
	var candidates []*node
	for _, child := range n.children {
		ucb := child.ucb1(explorationFactor)
		if ucb > best {
			best = ucb
			candidates = candidates[:0]
			candidates = append(candidates, child)
		} else if ucb == best {
			candidates = append(candidates, child)
		}
	}
	return candidates[rng.IntN(len(candidates))]
}

// treeSize returns the number of nodes in the subtree rooted at this node.
func (n *node) treeSize() int {
	size := 1
	for _, child := range n.children {
		size += child.treeSize()
	}
	return size
}
