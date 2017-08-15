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
    }

    hand := []deck.Card {
        deck.Card { deck.H, deck.A },
        deck.Card { deck.D, deck.Nine },
        deck.Card { deck.D, deck.A },
        deck.Card { deck.C, deck.Ten },
    }

    cards1 := [4]deck.Card {
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
    }

    prior := []Trick { trick1 }
    move := deck.Card { deck.H, deck.Nine }

    state := NewUndeterminizedState(setup, 2, hand, played, prior, move)
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
    }

    hand := []deck.Card {
        deck.Card { deck.H, deck.A },
        deck.Card { deck.D, deck.Nine },
        deck.Card { deck.D, deck.A },
        deck.Card { deck.C, deck.Ten },
    }

    cards1 := [4]deck.Card {
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
    }

    prior := []Trick { trick1 }
    move := deck.Card { deck.C, deck.Ten }

    state := NewUndeterminizedState(setup, 0, hand, played, prior, move)
    state.Determinize()

    if len(state.Hands[0]) != 4 || len(state.Hands[1]) != 4 ||
       len(state.Hands[2]) != 3 || len(state.Hands[3]) != 3 {
        t.Errorf("Incorrect number of cards in some players hand.\n")
    }
}
