package euchre

import (
    "ai"
    "deck"
    "fmt"
    "testing"
)


/*
 * Tests the ai package in conjunction with the euchre package. Some of these
 * tests are not real tests, but are simple sanity checks whose output must be
 * manually checked.
 */


/*
 * Tests the output for a run playout for the first card played by the computer.
 */
func TestRunPlayout(t *testing.T) {
    setup := Setup {
        1,
        1,
        true,
        deck.Card { deck.D, deck.Nine },
        deck.D,
        deck.Card{ },
        -1,
    }

    hand := []deck.Card {
        deck.Card { deck.H, deck.Nine },
        deck.Card { deck.H, deck.Ten },
        deck.Card { deck.S, deck.A },
        deck.Card { deck.D, deck.Q },
        deck.Card { deck.C, deck.Q },
    }

    played := []deck.Card {
        deck.Card { deck.C, deck.J },
        deck.Card { deck.C, deck.A },
    }

    var prior []Trick

    s := NewUndeterminizedState(setup, 0, hand, played, prior)
    s.Determinize()
    n := ai.NewNode()
    m := ai.Move {
        nil,
        s,
    }
    n.Value(m)
    e := Engine{ }

    fmt.Println("Playout debug output")
    ai.RunPlayoutDebug(n, e)
}
