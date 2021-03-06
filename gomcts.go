package gomcts

import (
	"math"
	"fmt"
	"math/rand"
	"log"
)

var (
)

const (
	ExplorationFactor = 1.0
	dummyLowUCTValue = -10000.0 // Assumes all scores are positive
)


/////////////////////////////////////////////////////////////
// TreeNode structure, New function, and methods
/////////////////////////////////////////////////////////////

type TreeNode struct {
	State               GameState
	Parent              *TreeNode
	GeneratingMove      string
	ShuffledPossibleMoves []string
	NumberOfChildren    int
	Children            []*TreeNode
	NextMoveToTry       int
	VisitCount          float64
	CumulativeScore     [2]float64 //Scores are to be between 0 (loss) and 1 (win)
}

func NewNode(state GameState, parent *TreeNode, generatingMove string) *TreeNode {
	moves := state.PossibleMovesShuffled()
	numChilds := len(moves)
	childs := make([]*TreeNode, numChilds)
	return &TreeNode{State: state,
		Parent:              parent,
		GeneratingMove:      generatingMove,
		ShuffledPossibleMoves: moves,
		NumberOfChildren:    numChilds,
		Children:            childs}
	//rest of fields should start with "zeroed" values
}

func (node *TreeNode) averageScore() float64 {
	var result float64
	if node.State.IsSecondPlayersTurn() {
		result = node.CumulativeScore[0] / node.VisitCount		//Because parent is not second player's turn
	} else {
		result = node.CumulativeScore[1] / node.VisitCount
	}
	return result
}

func (node *TreeNode) BestMoveFromNIterations(iterationBudget int) string {
	node.DoNIterations(iterationBudget)
	return node.BestChild().GeneratingMove
}

func (node *TreeNode) DoNIterations(iterationBudget int) {
	for i := 0; i < iterationBudget; i++ {
		currentNode := node.doTreePolicy()
		reward := currentNode.State.RewardFromRandomPlayout()
		currentNode.backpropagateReward(reward)
	}
	log.Printf("Node's best move is %v", node.BestChild().GeneratingMove)
	log.Printf("%v", node.Summary())
}


func (node *TreeNode) doTreePolicy() *TreeNode {
	selectedNode := node
	for !(selectedNode.State.IsTerminal()) {
		if selectedNode.IsFullyExpanded() {
			selectedNode = selectedNode.ChildToExplore()
		} else {
			return selectedNode.NewChild()
		}
	}
	return selectedNode //returned node has terminal game state here, I think
}

//Two-player version using negamax 
func (node *TreeNode) backpropagateReward(scores [2]float64) {
	currentNode := node
	for currentNode.Parent != nil {
		currentNode.VisitCount += 1.0
		currentNode.CumulativeScore[0] += scores[0]
		currentNode.CumulativeScore[1] += scores[1]
		currentNode = currentNode.Parent
	}
	//Increment root node counter
	currentNode.VisitCount += 1.0
}

func (node *TreeNode) BestChild() *TreeNode {
	var bestAveScore float64 = dummyLowUCTValue		
	var selectedChild *TreeNode
	for i := 0; i < node.NextMoveToTry; i++ {
		child := node.Children[i]
		aveScore := child.averageScore()
		//Need random tie-breaker?
		if aveScore >= bestAveScore {
			selectedChild = child
			bestAveScore = aveScore
		}
	}
	return selectedChild
}

func (node *TreeNode) ChildToExplore() *TreeNode {
	var bestUctValue float64 = dummyLowUCTValue
	var selectedChild *TreeNode
	for _, child := range node.Children {
		uctValue := child.averageScore() + ExplorationFactor * math.Sqrt(2.0 * math.Log(node.VisitCount) / child.VisitCount)
		//Need random tie-breaker?
		if uctValue >= bestUctValue {
			selectedChild = child
			bestUctValue = uctValue
		}
	}
	return selectedChild
}

func (node *TreeNode) IsFullyExpanded() bool {
	return (node.NextMoveToTry == node.NumberOfChildren)
}

func (node *TreeNode) NewChild() *TreeNode {
	//NextMoveToTry should < NumberOfChildren
	nextMove := node.ShuffledPossibleMoves[node.NextMoveToTry]
	newState := node.State.NewGameStateFromMove(nextMove)
	newChild := NewNode(newState, node, nextMove)
	node.Children[node.NextMoveToTry] = newChild
	node.NextMoveToTry++
	return newChild
}

func (node *TreeNode) Summary() string {
	var summary string
	summary = fmt.Sprintf("Number of iterations done: %v\n", node.VisitCount)
	summary += fmt.Sprintf("Root has explored %v of %v children\n", node.NextMoveToTry, node.NumberOfChildren)
	summary += fmt.Sprintf("Root has  %v descendants (counting itself)\n", len(descendants(node)))
	for i := 0; i < node.NextMoveToTry; i++ {
		child := node.Children[i]
		summary += fmt.Sprintf(" %v) Move %v: %v visits, %v cum. score, %v average score, %v descendants\n", i, child.GeneratingMove, child.VisitCount, child.CumulativeScore, child.averageScore(), len(descendants(child)))
	}
	return summary
}

func descendants(node *TreeNode) []*TreeNode {
	result := []*TreeNode{node}
	for i := 0; i< node.NextMoveToTry; i++ {
		result = append(result, descendants(node.Children[i])...)
	}
	return result
}


/////////////////////////////////////////////////////////////
// GameState interface
/////////////////////////////////////////////////////////////

type GameState interface {
	NumberOfMoves() int
	PossibleMoves() []string
	PossibleMovesShuffled() []string
	IsTerminal() bool
	TerminalReward() [2]float64
	NewGameStateFromMove(move string) GameState
	RewardFromRandomPlayout() [2]float64
	IsSecondPlayersTurn() bool
	LocalRand() *rand.Rand
	DoMove(string)
	Summary() string
	CurrentPlayer() string
}

/////////////////////////////////////////////////////////////
// Playout 
/////////////////////////////////////////////////////////////


func DoRandomPlayout(gstate GameState) {
	// log.Printf("Playout\n%v", gstate.Summary())
	for !(gstate.IsTerminal()) {
		var i int
		if gstate.NumberOfMoves() == 0 {
			panic("random playout called on game state with zero possible moves")
		}
		if gstate.NumberOfMoves() == 1 {
			i = 0
		} else {
			i = gstate.LocalRand().Intn(gstate.NumberOfMoves())
		}
		// log.Printf("selected move %v of %v", gstate.PossibleMoves()[i], gstate.NumberOfMoves())
		gstate.DoMove(gstate.PossibleMoves()[i])
		// log.Printf(gstate.Summary())		
	}
}


