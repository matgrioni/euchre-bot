package ai

import (
    "container/heap"
    "fmt"
    "math"
    "math/rand"
    "time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type State interface {
    Hash() interface{}
    Determinize()
    // TODO: Is there a more efficient way than copying to get new
    //       determinizations?
    Copy() State
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
    depth int
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


/*
 * Performs a Monte Carlo Tree search on the given the state and the game engine.
 * This MCTS is on a non-deterministic game so specify the amount of random
 * determinizations to test it on, and the number of runs to do per
 * determinization.
 *
 * Args:
 *  s: The current state from which to start simulation.
 *  engine: The game engine with which to step through game logic.
 *  runs: The number of times run a given determinization.
 *  deters: The number of determinizations to run through.
 *
 * Returns:
 *  The state with the highest value and the expected value associated with it.
 */
func MCTS(s State, engine MCTSEngine, runs int, deters int) (State, float64) {
    // TODO: Is there a better way than this dual map way. This probably isn't
    // a bottleneck however.
    weights := make(map[interface{}]float64)
    conv := make(map[interface{}]State)
    counts := make(map[interface{}]int)

    for i := 0; i < deters; i++ {
        copyState := s.Copy()
        copyState.Determinize()
        n := NewNode()
        n.Value(copyState)

        for j := 0; j < runs; j++ {
            RunPlayout(n, engine)

            if n.children.Len() > 0 {
                topNode := n.children.Poll().(*Node)
                topState := topNode.GetState()

                conv[topState.Hash()] = topState
                weights[topState.Hash()] += topNode.GetPriority()
                counts[topState.Hash()] += 1
            } else {
                fmt.Println("HEre")
                nState := n.GetState()
                conv[nState.Hash()] = nState
                weights[nState.Hash()] += n.GetPriority()
                counts[nState.Hash()] += 1
            }
        }
    }

    maxWeight := math.Inf(-1)
    var maxState State
    for hash, weight := range weights {
        if weight > maxWeight {
            maxState = conv[hash]
            maxWeight = weight
        }
    }

    return maxState, maxWeight / float64(counts[maxState.Hash()])
}


/*
 * This method is for internal testing of an MCTS playout. It prints out debug
 * info so that the MCTS process can be verified.
 *
 * Args:
 *  node   - A node in the MCTS tree to start from.
 *  engine - The engine for traversing through the MCTS tree.
 *
 * Returns:
 *  An integer that represents the final terminal state of the playout per the
 *  engine's computation.
 */
func RunPlayoutDebug(node *Node, engine MCTSEngine) int {
    return runPlayout(node, engine, true)
}


/*
 * The normal MCTS playout method. This method does not provide any logging.
 *
 * Args:
 *  node   - A node in the MCTS tree to start from.
 *  engine - The engine for traversing the MCTS tree.
 *
 * Returns:
 *  An integer that represents the final terminal state of the playout per the
 *  engine's computation.
 */
func RunPlayout(node *Node, engine MCTSEngine) int {
    return runPlayout(node, engine, false)
}


/*
 * The internal logic for the MCTS tree logic. Provides a logging flag for
 * debugging purposes.
 *
 * Args:
 *  node   - A node in the MCTS tree to start from.
 *  engine - The engine for traversing the MCTS tree.
 *  log    - A flag to indicate whether the function should log.
 *
 * Returns:
 *  An integer that represents the final terminal state of the playout per the
 *  engine's computation.
 */
func runPlayout(node *Node, engine MCTSEngine, log bool) int {
    if log {
        fmt.Println(node.GetState())
    }

    node.simulations++

    var eval int
    // We have been given a node state that is the last in the playout. Time to
    // return and backpropagate the results.
    if engine.IsTerminal(node.GetState()) {
        eval = engine.Evaluation(node.GetState())
    } else {
        var nextStates []State
        if node.depth <= 2 {
            if node.memoize == nil {
                node.memoize = engine.NextStates(node.GetState())
            }

            nextStates = node.memoize
        } else {
            nextStates = engine.NextStates(node.GetState())
        }

        var next *Node

        // If we don't have data on all the posssible next states, select one at
        // random. Otherwise, choose the one with the highest UCB.
        if len(nextStates) > node.children.Len() {
            takenMoves := make(map[interface{}]int)

            for i := 0; i < node.children.Len(); i++ {
                takenMoves[node.children[i].(*Node).GetState().Hash()] = i
            }

            nextState := nextStates[r.Intn(len(nextStates))]

            if _, ok := takenMoves[nextState.Hash()]; ok {
                next = node.children[takenMoves[nextState.Hash()]].(*Node)
            } else {
                next = NewNode()
                next.Value(nextState)
                next.parent = node
                next.depth = node.depth + 1
                heap.Push(&node.children, next)
            }
        } else {
            next = node.children.Poll().(*Node)
        }
        eval = runPlayout(next, engine, log)

        adjEval := eval
        if adjEval < 0 {
            adjEval *= -1
        }
        if engine.Favorable(node.GetValue().(State), eval) {
            next.wins += adjEval
        } else {
            next.wins -= adjEval
        }

        next.Priority(UpperConfBound(next))
        node.children.Update(next, next.GetValue(), next.GetPriority())
    }

    return eval
}
