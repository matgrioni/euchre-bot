package ai

import (
    "container/heap"
    "math"
    "math/rand"
    "time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type State interface {
    Hash() interface{}
    Weight() float64
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

    memoize []State
}

// Return a new node that is properly initialized. Specifically, the priority
// queue for the children is properly made.
func NewNode() *Node {
    var n Node

    n.children = make(PriorityQueue, 0)
    heap.Init(&n.children)

    return &n
}

func (node *Node) GetState() State {
    return node.value
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
    Favorable(state State, eval int) bool
    IsTerminal(state State) bool
    NextStates(state State) []State
    Evaluation(state State) int
}

func MCTS(s State, engine MCTSEngine, runs int) (State, float64) {
    n := NewNode()
    n.Value(s)

    for i := 0; i < runs; i++ {
        runPlayout(n, engine)
    }

    if n.children.Len() > 0 {
        topNode := n.children.Poll()
        for i := 0; i < n.children.Len(); i++ {
            child := n.children[i]
            fmt.Printf("%d\t%d\t%f\n", child.(*Node).wins, child.(*Node).simulations, child.GetPriority())
        }

        return topNode.GetValue().(State), topNode.GetPriority()
    } else {
        return n.GetValue().(State), n.GetPriority()
    }
}


func runPlayout(node *Node, engine MCTSEngine) int {
    node.simulations++

    var eval int
    // We have been given a node state that is the last in the playout. Time to
    // return and backpropagate the results.
    if engine.IsTerminal(node.GetState()) {
        eval = engine.Evaluation(node.GetState())
    } else {
        if node.memoize == nil {
            node.memoize = engine.NextStates(node.GetState())
        }
        nextStates := node.memoize

        var next *Node

        // If we don't have data on all the posssible next states, select one at
        // random. Otherwise, choose the one with the highest UCB.
        if len(nextStates) > node.children.Len() {
            nextState := nextStates[r.Intn(len(nextStates))]

            next = NewNode()
            next.Value(nextState)
            next.parent = node
            heap.Push(&node.children, next)
        } else {
            next = node.children.Poll().(*Node)
        }
        eval = runPlayout(next, engine)

        adjEval := eval
        if engine.Favorable(node.GetValue().(State), eval) {
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
