package deck

import (
    "errors"
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

// Create a Suit from the input string. An error is provided if the input is not
// a valid Suit.
// s - The string value to convert to a Suit. Intuitive mapping.
// Returns the converted Suit or an error if something went wrong.
func CreateSuit(s string) (Suit, error){
    var res Suit
    switch s {
    case "H":
        res = H
    case "D":
        res = D
    case "S":
        res = S
    case "C":
        res = C
    default:
        return H, errors.New("Input is not a valid suit.")
    }

    return res, nil
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

func (s Suit) String() string {
    return string(s)
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
// s - The string to convert to a value. Intuitive mapping.
// Returns a Value type that represents the parameter and an error if anything
// went wrong.
func CreateValue(s string) (Value, error) {
    var res Value
    switch s {
    case "9":
        res = Nine
    case "10":
        res = Ten
    case "J":
        res = J
    case "Q":
        res = Q
    case "K":
        res = K
    case "A":
        res = A
    default:
        return Nine, errors.New("Input does not represent a valid value.")
    }

    return res, nil
}

func (v Value) String() string {
    switch v {
    case Nine:
        return "9"
    case Ten:
        return "10"
    case J:
        return "J"
    case Q:
        return "Q"
    case K:
        return "K"
    case A:
        return "A"
    }

    return ""
}

// A Card represents a playing card from a standard 52 card deck. It consists of
// a suit, such as Hearts (H), and a value such as J. The suit is represented by
// the Suit type, and the value is a simple int that should be in the range
// [9, 14], where 14 is A, 13 is K, and so on.
type Card struct {
    Suit Suit
    Value Value
}

// Create an array of the all the cards in the euchre deck.
var CARDS = createCards()

// Creates a card given the string in the format of VS, where V is the value, and
// S is the suit.
func CreateCard(s string) (Card, error) {
    var card Card
    var sErr, vErr error

    card.Suit, sErr = CreateSuit(s[len(s) - 1:])
    card.Value, vErr = CreateValue(s[:len(s) - 1])

    if sErr != nil || vErr != nil {
        return card, errors.New("There was an error in the input.")
    }

    return card, nil
}

func (c Card) String() string {
    return c.Value.String() + c.Suit.String()
}

// Adjusts suit of this card based on the trump suit. This is only really
// valuable when it matters if the card can be the left bower. In this case,
// this method returns that the suit of this card is the trump suit. For all
// other cards, the suit is simply outputted.
func (c Card) AdjSuit(t Suit) Suit {
    adjSuit := c.Suit
    if c.Value == J && c.Suit == t.Left() {
        adjSuit = t
    }

    return adjSuit
}

// Generate a random card.
func GenCard() Card {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    var card Card
    card.Suit = SUITS[r.Intn(4)]
    card.Value = VALUES[r.Intn(6)]

    return card
}

// TODO: Should this be in euchre package?
// Randomly generates a hand of cards, which is 5 cards in euchre.
func GenHand() [5]Card {
    var hand [5]Card
    present := make(map[Card]bool)
    for i := range hand {
        hand[i] = GenCard()

        // Ensure that any generated card is only included once in the result.
        for _, in := present[hand[i]]; in ; _, in = present[hand[i]] {
            hand[i] = GenCard()
        }

        present[hand[i]] = true
    }

    return hand
}

// TODO: This copies the array. Is that a problem / when should I use array vs
// slice in golang.
// A helper method that simply creates an array that has all the cards in a
// euchre deck.
func createCards() [24]Card {
    var cards [24]Card
    for i, value := range VALUES {
        for j, suit := range SUITS {
            cards[i * len(SUITS) + j] = Card { suit, value }
        }
    }

    return cards
}
