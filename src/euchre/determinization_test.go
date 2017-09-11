package euchre

import (
    "deck"
    "fmt"
    "testing"
)


/*
 * Test the determinization of the euchre State when there have been no cards
 * played in the current trick.
 */
func TestDeterminization(t *testing.T) {
    setup := Setup {
        1,
        2,
        false,
        deck.Card { deck.C, deck.A },
        deck.S,
        deck.Card { },
        -1,
    }

    hand := []deck.Card {
        deck.Card { deck.H, deck.A },
        deck.Card { deck.D, deck.Nine },
        deck.Card { deck.D, deck.A },
        deck.Card { deck.C, deck.Ten },
    }

    cards1 := []deck.Card {
        deck.Card { deck.H, deck.K },
        deck.Card { deck.C, deck.Q },
        deck.Card { deck.H, deck.Q },
        deck.Card { deck.H, deck.Nine },
    }

    var played []deck.Card

    trick1 := Trick {
        cards1,
        2,
        deck.S,
        -1,
    }

    prior := []Trick { trick1 }

    state := NewUndeterminizedState(setup, 2, hand, played, prior)
    state.Determinize()

    fmt.Printf("Sanity check: %v\n", state)
}


/*
 * Test that a determinization in the middle of a trick will result in the
 * proper amount of cards in each hand.
 */

func TestDeterminizationMiddleOfTrick(t *testing.T) {
    setup := Setup {
        1,
        2,
        false,
        deck.Card { deck.C, deck.A },
        deck.S,
        deck.Card { },
        -1,
    }

    hand := []deck.Card {
        deck.Card { deck.H, deck.A },
        deck.Card { deck.D, deck.Nine },
        deck.Card { deck.D, deck.A },
        deck.Card { deck.C, deck.Ten },
    }

    cards1 := []deck.Card {
        deck.Card { deck.H, deck.K },
        deck.Card { deck.C, deck.Q },
        deck.Card { deck.H, deck.Q },
        deck.Card { deck.H, deck.Nine },
    }

    played := []deck.Card {
        deck.Card { deck.C, deck.A },
        deck.Card { deck.C, deck.Ten },
    }

    trick1 := Trick {
        cards1,
        2,
        deck.S,
        -1,
    }

    prior := []Trick { trick1 }

    state := NewUndeterminizedState(setup, 0, hand, played, prior)
    state.Determinize()

    if len(state.Hands[0]) != 4 || len(state.Hands[1]) != 4 ||
       len(state.Hands[2]) != 3 || len(state.Hands[3]) != 3 {
        t.Errorf("Incorrect number of cards in some players hand.\n")
    }
}


/*
 * Test a problematic instance that is giving problems in generating next valid
 * states.
 */
func TestDeterminizationProblematic(t *testing.T) {
    setup := Setup {
        1,
        0,
        true,
        deck.Card { deck.C, deck.K },
        deck.C,
        deck.Card { },
        -1,
    }

    cards1 := []deck.Card {
        deck.Card { deck.C, deck.J },
        deck.Card { deck.C, deck.Ten },
        deck.Card { deck.C, deck.Q },
        deck.Card { deck.C, deck.K },
    }

    cards2 := []deck.Card {
        deck.Card { deck.D, deck.A },
        deck.Card { deck.D, deck.K },
        deck.Card { deck.D, deck.Q },
        deck.Card { deck.D, deck.J },
    }

    cards3 := []deck.Card {
        deck.Card { deck.D, deck.Ten },
        deck.Card { deck.C, deck.Nine },
        deck.Card { deck.C, deck.A },
        deck.Card { deck.S, deck.J },
    }

    played := []deck.Card {
        deck.Card { deck.S, deck.Q },
        deck.Card { deck.S, deck.K },
        deck.Card { deck.H, deck.A },
    }

    hand := []deck.Card {
        deck.Card { deck.S, deck.Nine },
        deck.Card { deck.S, deck.A },
    }


    trick1 := Trick {
        cards1,
        2,
        deck.C,
        -1,
    }

    trick2 := Trick {
        cards2,
        2,
        deck.C,
        -1,
    }

    trick3 := Trick {
        cards3,
        2,
        deck.C,
        -1,
    }

    prior := []Trick { trick1, trick2, trick3 }

    state := NewUndeterminizedState(setup, 0, hand, played, prior)
    state.Determinize()

    fmt.Printf("Sanity check: %v\n", state)
    if len(state.Hands[0]) != 2 || len(state.Hands[1]) != 1 ||
       len(state.Hands[2]) != 1 || len(state.Hands[3]) != 1 {
        t.Errorf("Incorrect number of cards in some players hand.\n")
    }
}
