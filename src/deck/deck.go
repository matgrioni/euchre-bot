package deck

import (
    "math/rand"
    "time"
)

// Define a Suit type off of the string type. This Suit type should only take on
// 4 values. These are the constants defined below for the 4 suits of a standard
// 52 card deck, H, D, S, C.
type Suit string
const (
    H Suit = "H"
    D Suit = "D"
    S Suit = "S"
    C Suit = "C"
)

// An array of all the suits. In order they are hearts, diamonds, spades, clubs.
var SUITS = [4]Suit { H, D, S, C, }

// Create a Suit from the input string. In this case since, Suit is already an
// alias for the string type, it is just a matter of casting.
func CreateSuit(s string) Suit {
    return Suit(s)
}

// Return the suit of the left bower given what the current suit is.
func (s Suit) Left() Suit {
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

// Define a Value type off the int type. Each Value corresponds to the different
// cards used in euchre. A is high at value 14, and Nine is low at value 9.
type Value int
const (
    Nine Value = iota + 9
    Ten
    J
    Q
    K
    A
)

// An array of all the values in ascending order of value.
var VALUES = [6]Value { Nine, Ten, J, Q, K, A }

// Returns a Value type from the input string. The mapping is evident from the
// standard 52 card deck.
func CreateValue(s string) Value {
    switch s {
    case "9":
        return Nine
    case "10":
        return Ten
    case "J":
        return J
    case "Q":
        return Q
    case "K":
        return K
    case "A":
        return A
    }

    return Nine
}

// A Card represents a playing card from a standard 52 card deck. It consists of
// a suit, such as Hearts (H), and a value such as J. The suit is represented by
// the Suit type, and the value is a simple int that should be in the range
// [9, 14], where 14 is A, 13 is K, and so on.
type Card struct {
    Suit Suit
    Value Value
}

// Creates a card given the string in the format of VS, where V is the value, and
// S is the suit.
func CreateCard(s string) Card {
    var card Card
    card.Suit = CreateSuit(s[len(s) - 1:])
    card.Value = CreateValue(s[:len(s) - 1])

    return card
}

// Generate a random card.
func GenCard() Card {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    var card Card
    card.Suit = SUITS[r.Intn(4)]
    card.Value = VALUES[r.Intn(6)]

    return card
}

// Randomly generates a hand of cards, which is 5 cards.
func GenHand() [5]Card {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    var hand [5]Card
    for i := range hand {
        hand[i].Suit = SUITS[r.Intn(4)]
        hand[i].Value = VALUES[r.Intn(6)]
    }

    return hand
}
