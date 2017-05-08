package euchre

import (
    "math/rand"
    "time"
)

// TODO: Where should these definitions go?
type Suit int
const (
    H Suit = iota
    D Suit = iota
    S Suit = iota
    C Suit = iota
)

func NewSuit(s string) Suit {
    switch s {
    case "H":
        return H
    case "D":
        return D
    case "S":
        return S
    case "C":
        return C
    }

    return H
}

func (s Suit) Right() Suit {
    switch s {
    case H:
        return D
    case D:
        return H
    case S:
        return C
    case C:
        return S
    }

    return H
}

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

    return "H"
}

var SUITS = [4]Suit { H, D, S, C, }

// A Card represents a playing card from a standard 52 card deck. It consists of
// a suit, such as Hearts (H), and a value such as J. The suit is represented by
// the Suit type, and the value is a simple int that should be in the range
// [9, 14], where 14 is A, 13 is K, and so on.
// TODO: Make value type?
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
// cards currently in the players hand, the top card, and whether one picks up
// the card in question.
func RPickUp(hand [5]Card, top Card, friendly bool) bool {
    // Arbitrary weights that I felt like using based on experience.
    weights := map[int]float32 {
        9: 0.05,
        10: 0.07,
        11: 0.3,
        12: 0.12,
        13: 0.15,
        14: 0.2,
    }

    // One wants to tell the dealer to pick up if you believe that you can win 3
    // tricks if the dealer picks it up. So if your confidence / probability is
    // greater than 50% (in a non-technical sense) then you should take such an
    // action.
    var conf float32 = 0

    // TODO: How to make a better loop.
    for i, _ := range hand {
        card := hand[i]
        if card.Suit == top.Suit {
            conf += weights[card.Value]
        } else if card.Suit.Right() == top.Suit && card.Value == 11 {
            conf += 0.25
        }
    }

    if friendly {
        conf += weights[top.Value]
    }

    return conf >= 0.5
}
