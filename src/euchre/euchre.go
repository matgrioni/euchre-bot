package euchre

import (
    "math/rand"
    "time"
)

type Suit int
const (
    H Suit = iota
    D Suit = iota
    S Suit = iota
    C Suit = iota
)
func (s Suit) String() string {
    if s == H {
        return "H"
    } else if s == D {
        return "D"
    } else if s == S {
        return "S"
    } else if s == C {
        return "C"
    }

    return " "
}
var SUITS = [4]Suit { H, D, S, C, }

// A Card represents a playing card from a standard 52 card deck. It consists of
// a suit, such as Hearts (H), and a value such as J. The suit is represented by
// the Suit type, and the value is a simple int that should be in the range
// [9, 14], where 14 is A, 13 is K, and so on.
type Card struct {
    Suit Suit
    Value int
}

func GenHand() [5]Card {
    // TODO: Look at differences in seeding and using random library.
    rand.Seed(time.Now().UnixNano())

    var hand [5]Card
    // TODO: How to make a better loop.
    for i, _ := range hand {
        hand[i].Suit = SUITS[rand.Intn(4)]
        hand[i].Value = rand.Intn(6) + 9
    }

    return hand
}

// A rule based approach to deciding if one should call for the dealer to pick
// up the top card at the start of a deal. This approach takes into account the
// cards currently in the players hand, the top card, and whether one or one's
// partner picks up the card in question.
func RPickUp(hand [5]Card, top Card, friendly bool) bool {
    return true
}
