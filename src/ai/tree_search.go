package ai


type TSState interface { }


/*
 * A type to represent a move in a game. A move consists of an action, and a
 * successsor state. In the case where the game is in a terminal state, this is
 * represented via a nil action.
 */
type Move struct {
    Action interface{}
    State TSState
}


type TSEngine interface {
    Favorable(state TSState) bool
    IsTerminal(state TSState) bool
    Evaluation(state TSState) float64
    Successors(state TSState) []Move
}
