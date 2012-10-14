package gomcts

import (
	"bufio"
	"log"
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
	State          *gameState
	generatingMove string
	Parent         *TreeNode
	Children       [string]TreeNode
	BestChildMove  string
	bestChildScore float64
	VisitCount     uint64
	shuffledMoves  []string
	nextMoveToTry  uint16
}

func NewNode(state *gameState, parent *TreeNode, generatingMove string) *TreeNode {
	newNode := TreeNode{State: state, Parent: parent, generatingMove: generatingMove} //rest of fields should start with "zeroed" values
	newNode.shuffledMoves = shuffle(newNode.PossibleMoves)
}

func (node *TreeNode) bestMoveFromNIterations(iterationBudget uint64) string {

}

func (node *TreeNode) doTreePolicy() *TreeNode {

}

//Two-player version using negamax 
func (node *TreeNode) backpropagate(reward float64) {

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
