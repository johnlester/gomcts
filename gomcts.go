package gomcts

import (
	// "bufio"
	// "log"
	"math"
	"math/rand"
	// "os"
	"strconv"
)

var ()

const (
	startingBestChildScore = -1E9
	oneDividedBySqrtOfTwo  = 0.70710678118
)

/////////////////////////////////////////////////////////////
// TreeNode structure, New function, and methods
/////////////////////////////////////////////////////////////

type TreeNode struct {
	State               GameState
	Parent              *TreeNode
	GeneratingMove      string
	SortedPossibleMoves []string
	NumberOfChildren    int
	Children            []*TreeNode
	orderToTryMoves     []int
	NextMoveToTry       int
	VisitCount          float64
	CumulativeScore     float64 //Scores are to be between 0 (loss) and 1 (win) 
}

func NewNode(state GameState, parent *TreeNode, generatingMove string) *TreeNode {
	sortedMoves := state.PossibleMoves()
	numChilds := len(sortedMoves)
	return &TreeNode{State: state,
		Parent:              parent,
		GeneratingMove:      generatingMove,
		SortedPossibleMoves: sortedMoves,
		NumberOfChildren:    numChilds,
		Children:            make([]*TreeNode, numChilds),
		orderToTryMoves:     rand.Perm(numChilds)}
	//rest of fields should start with "zeroed" values
}

func (node *TreeNode) bestMoveFromNIterations(iterationBudget int) string {
	for i := 0; i < iterationBudget; i++ {
		currentNode := node.doTreePolicy()
		reward := currentNode.State.RewardFromRandomPlayout()
		currentNode.backpropagateReward(reward)
	}
	return node.BestChild(0.0).GeneratingMove
}

func (node *TreeNode) doTreePolicy() *TreeNode {
	selectedNode := node
	for selectedNode.State.IsNotTerminal() {
		if selectedNode.IsFullyExpanded() {
			selectedNode = selectedNode.BestChild(1.0)
		} else {
			return selectedNode.NewChild()
		}
	}
	return selectedNode //must be terminal, I think
}

//Two-player version using negamax 
func (node *TreeNode) backpropagateReward(reward float64) {
	currentNode := node
	currentReward := reward
	for currentNode.Parent != nil {
		currentNode.VisitCount += 1.0
		currentNode.CumulativeScore += reward
		currentReward = currentReward * -1.0
		currentNode = currentNode.Parent
	}
}

func (node *TreeNode) BestChild(factor float64) *TreeNode {
	var bestUctValue float64 = -1E10
	var selectedChild *TreeNode
	for _, child := range node.Children {
		uctValue := (child.CumulativeScore / child.VisitCount) + factor*math.Sqrt(2.0*math.Log(node.VisitCount)/child.VisitCount)
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
	nextMove := node.SortedPossibleMoves[node.orderToTryMoves[node.NextMoveToTry]]
	newState := node.State.NewGameStateFromMove(nextMove)
	newChild := NewNode(newState, node, nextMove)
	node.Children[node.NextMoveToTry] = newChild
	node.NextMoveToTry++
	return newChild
}

/////////////////////////////////////////////////////////////
// GameState interface
/////////////////////////////////////////////////////////////

type GameState interface {
	PossibleMoves() []string
	IsNotTerminal() bool
	TerminalReward() float64
	//copy() GameState
	NewGameStateFromMove(move string) GameState
	RewardFromRandomPlayout() float64
}

/////////////////////////////////////////////////////////////
// GameState implementation of my made-up game 123ToTen
/////////////////////////////////////////////////////////////

type GameState123ToTen struct {
	Total       int
	PlayerBMove bool
}

func (gstate GameState123ToTen) PossibleMoves() []string {
	moves := []string{"1", "2", "3"}
	return moves
}

func (gstate GameState123ToTen) IsNotTerminal() bool {
	return gstate.Total < 10
}

func (gstate GameState123ToTen) TerminalReward() float64 {
	// error if not terminal?
	var reward float64
	if gstate.PlayerBMove {
		reward = 0.0
	} else {
		reward = 1.0
	}
	return reward
}

func (gstate GameState123ToTen) NewGameStateFromMove(move string) GameState123ToTen {
	cpy := gstate.Copy()
	cpy.DoMove(move)
	return cpy
}

func (gstate *GameState123ToTen) DoMove(move string) {
	moveFromString, _ := strconv.Atoi(move)
	gstate.Total += moveFromString
	gstate.PlayerBMove = !gstate.PlayerBMove
}


func (gstate GameState123ToTen) RewardFromRandomPlayout() float64 {
	cpy := gstate.Copy()
	for cpy.IsNotTerminal() {
		cpy.DoMove(cpy.RandomMove())
	}
	return cpy.TerminalReward()
}

func (gstate GameState123ToTen) RandomMove() string {
	rndIdx := rand.Intn(len(gstate.PossibleMoves()))
	return gstate.PossibleMoves()[rndIdx]
}


func (gstate GameState123ToTen) Copy() GameState123ToTen {
	newCopy := gstate
	return newCopy
}


