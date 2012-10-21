package gomcts

import (
	"bufio"
	"log"
	"math/rand"
	"os"
)

var ()

const ()

// type Tree struct {
// 	root *treeNode
// }

/////////////////////////////////////////////////////////////
// TreeNode structure, New function, and methods

type TreeNode struct {
	State               *GameState
	Parent              *TreeNode
	GeneratingMove      string
	SortedPossibleMoves []string
	NumberOfChildren    int
	Children            []*TreeNode
	orderToTryMoves     []int
	NextMoveToTry       int
	VisitCount          []uint64
	CumulativeScore     []float64
	BestChildMove       int
	BestChildScore      float64
}

func NewNode(state *GameState, parent *TreeNode, generatingMove string) *TreeNode {
	sortedMoves := state.PossibleMoves
	numChilds := len(sortedMoves)
	return TreeNode{State: state,
		Parent:              parent,
		GeneratingMove:      generatingMove,
		SortedPossibleMoves: sortedMoves,
		NumberOfChildren:    numChilds,
		Children:            make([]*TreeNode, numChilds),
		orderToTryMoves:     Perm(numChilds),
		VisitCount:          make([]unit64, numChilds),
		CumulativeScore:     make([]float64, numChilds),
		BestChildMove:       -1,
		BestChildScore:      -1E10} //rest of fields should start with "zeroed" values
}

func (node *TreeNode) bestMoveFromNIterations(iterationBudget uint64) string {

}

func (node *TreeNode) doTreePolicy() *TreeNode {
	selectedNode := node
	for selectedNode.State.IsNotTerminal() {
		if selectedNode.IsFullyExpanded() {
			selectedNode = selectedNode.BestChild(1.0)
		}
	}
}

//Two-player version using negamax 
func (node *TreeNode) backpropagate(reward float64) {

}

func (node *TreeNode) BestChild(factor float64) *TreeNode {
	var bestUCTValue float64 = -1E10
	var selectedChild *TreeNode
}

/////////////////////////////////////////////////////////////
// GameState structure and interface

type GameState interface {
	PossibleMoves() []string
	IsNotTerminal() bool
	TerminalReward() float64
	//copy() GameState
	NewGameStateFromMove(move string) GameState
	RewardFromRandomPlayout() float64
}
