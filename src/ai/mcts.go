package ai

import (
    "container/heap"
    "deck"
    "fmt"
    "math"
    "math/rand"
    "time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type State interface {
    Hash() interface{}
}

// This is a Node that is used for the MCTS tree. It has the attributes necessary
// for this role such as, parent, children, wins, and simulations but also
// implements methods from PQItem, since the list of node children is a priority
// queue based on the node's UCB. Any data can be passed along with a node
// through the Value methods, which accept a blank interface type.
type Node struct {
    value State
    priority float64
    index int

    children PriorityQueue
    parent *Node

    wins int
    simulations int
}

// Return a new node that is properly initialized. Specifically, the priority
// queue for the children is properly made.
func NewNode() *Node {
    var n Node

    n.children = make(PriorityQueue, 0)
    heap.Init(&n.children)

    return &n
}

func (node *Node) GetValue() interface{} {
    return node.value
}

func (node *Node) Value(v interface{}) {
    node.value = v.(State)
}

// The priority of a node is based on how many times all of its siblings have
// been sampled, how many times it has been sampled, and how many times it and
// its siblings have been sampled.
func (node *Node) GetPriority() float64 {
    return node.priority
}

func (node *Node) Priority(priority float64) {
    node.priority = priority
}

func (node *Node) GetIndex() int {
    return node.index
}

func (node *Node) Index(i int) {
    node.index = i
}

// The UpperConfBound for a node's expected winnings are based on the node's
// current win average, how many times it has been played and how many times
// its siblings have been played. This upper confidence bound is a tradeoff
// between exploitation and explortion. As we exploit one node more and more
// that we think is a good choice, it's confidence bound narrows, and we become
// more curious of other nodes.
// node - The node which we are calculating the UCB of.
// Returns the UCB of the node using sqrt(2) as the bias parameter.
func UpperConfBound(node *Node) float64 {
    var ucb float64
    if node.simulations == 0 && node.parent != nil {
        ucb = math.Inf(1)
    } else if node.parent != nil {
        ucb = float64(node.wins) / float64(node.simulations) +
              math.Sqrt(2.0 * (math.Log(float64(node.parent.simulations)) + 1) / float64(node.simulations))
    }
    return ucb
}

type MCTSEngine interface {
    Favorable(state interface{}, eval int) bool
    IsTerminal(state interface{}) bool
    NextStates(state interface{}) []interface{}
    Evaluation(state interface{}) int
}

// TODO: I feel like node should just be abstracted away. The interface{} State
// object that you choose is the only thing that the caller should see.
func MCTS(s State, engine MCTSEngine, runs int) State {
    n := NewNode()
    n.Value(s)

    /*nextStates := engine.NextStates(n.GetValue())
    if len(nextStates) > 1 {
        for i := 0; i < runs; i++ {
            runPlayout(n, engine)
        }
    } else {
        return nextStates[0].(State)
    }*/
    for i := 0; i < runs; i++ {
        runPlayout(n, engine)
    }

    if n.children.Len() > 0 {
        for _, child := range n.children {
            fmt.Printf("%s\t%d\t%d\t%f\n", child.GetValue().(State).Hash().(deck.Card), child.(*Node).wins, child.(*Node).simulations, child.GetPriority())
        }

        return n.children.Poll().GetValue().(State)
    } else {
        return n.GetValue().(State)
    }
}


func runPlayout(node *Node, engine MCTSEngine) int {
    node.simulations++

    var eval int
    // We have been given a node state that is the last in the playout. Time to
    // return and backpropagate the results.
    if engine.IsTerminal(node.GetValue()) {
        eval = engine.Evaluation(node.GetValue())
    } else {
        nextStates := engine.NextStates(node.GetValue())
        var next *Node

        // If we don't have data on all the posssible next states, select one at
        // random. Otherwise, choose the one with the highest UCB.
        if len(nextStates) > node.children.Len() {
            takenMoves := make(map[interface{}]bool)
            for i := 0; i < node.children.Len(); i++ {
                takenMoves[node.children[i].GetValue().(State).Hash()] = true
            }

            availableStates := make([]interface{}, 0)
            for i := 0; i < len(nextStates); i++ {
                nextState := nextStates[i]
                if _, ok := takenMoves[nextState.(State).Hash()]; !ok {
                    availableStates = append(availableStates, nextState)
                }
            }

            nextState := availableStates[r.Intn(len(availableStates))]
            next = NewNode()
            next.Value(nextState)
            next.parent = node
            heap.Push(&node.children, next)
        } else {
            next = node.children.Poll().(*Node)
        }
        eval = runPlayout(next, engine)

        adjEval := eval
        if engine.Favorable(node.GetValue(), eval) {
            if eval < 0 {
                adjEval = -1 * eval
            }

            next.wins += adjEval
        } else {
            if eval < 0 {
                adjEval = -1 * eval
            }

            next.wins -= adjEval
        }

        next.Priority(UpperConfBound(next))
        node.children.Update(next, next.GetValue(), next.GetPriority())
    }

    return eval
}
