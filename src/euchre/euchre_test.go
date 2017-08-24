package euchre

import (
    "deck"
    "testing"
)


// Test euchre#Possible. Returns the possible cards in a hand given the
// currently played cards.
// -----------------------------------------------------------------------------
func TestPossibleCantFollow(t *testing.T) {
    hand := []deck.Card {
        deck.Card{ deck.H, deck.Ten },
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.C, deck.J },
        deck.Card{ deck.C, deck.K },
        deck.Card{ deck.D, deck.J },
    }
    played := []deck.Card {
        deck.Card{ deck.S, deck.Ten },
    }
    trump := deck.H

    res := Possible(hand, played, trump)

    if len(res) != 5 {
        t.Errorf("Expected all cards to be possible.")
    }
}

func TestPossibleCanFollow(t *testing.T) {
    hand := []deck.Card {
        deck.Card{ deck.H, deck.Ten },
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.C, deck.J },
        deck.Card{ deck.C, deck.K },
        deck.Card{ deck.S, deck.J },
    }
    played := []deck.Card {
        deck.Card{ deck.S, deck.Ten },
    }
    trump := deck.H

    res := Possible(hand, played, trump)

    if len(res) != 1 || res[0] != 4 {
        t.Errorf("Expected only #4 not %d option(s).", len(res))
    }
}
