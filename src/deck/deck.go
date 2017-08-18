package deck

import "errors"

/*
 * Define a Suit type off of the string type. This Suit type should only take on
 * 4 values. These are the constants defined below for the 4 suits of a standard
 * 52 card deck, H, D, S, C.
 */
type Suit string
const (
    H Suit = "H"
    D Suit = "D"
    S Suit = "S"
    C Suit = "C"
)


/*
 * An array of all the suits. In order they are hearts, diamonds, spades, clubs.
 */
var SUITS = [4]Suit { H, D, S, C, }


/*
 * Create a Suit from the input string. An error is provided if the input is not
 * a valid Suit.
 *
 * Args:
 *  s: The string value to convert to a Suit. Intuitive mapping.
 *
 * Returns:
 *  Returns the converted Suit or an error if something went wrong.
 */
func CreateSuit(s string) (Suit, error) {
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


/*
 * The left bower suit given the current suit.
 *
 * Returns:
 *  The suit of the left bower if the right bower is this.
 */
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


/*
 * Convert the suit to a string representation.
 *
 * Returns:
 *  A string representation of the suit.
 */
func (s Suit) String() string {
    return string(s)
}



/*
 * Define a Value type off the int type. Each Value corresponds to the different
 * cards used in euchre. A is high at value 14, and Nine is low at value 9.
 */
type Value int
const (
    Nine Value = iota + 9
    Ten
    J
    Q
    K
    A
)


/*
 * An array of all the values in ascending order of value.
 */
var VALUES = [6]Value { Nine, Ten, J, Q, K, A }


/*
 * Compares two card values. The order of cards is: 9, 10, J, Q, K, A. If this
 * value (v1) is greater then v2, then a positive number is returned. If v1 is
 * less than v2 then negative number is returned, and if they are equal 0 is
 * returned.
 *
 * Args:
 *  v2: The value to compare this to.
 *
 * Returns:
 *  A positive, 0, or negative number if this is greater than, equal to, or less
 *  than v2, respectively.
 */
func (v1 Value) Compare(v2 Value) int {
    return int(v1) - int(v2)
}


/*
 * Returns a Value type from the input string. The mapping is evident from the
 * standard 52 card deck.
 *
 * Args:
 *  s: The string to convert to a value. Intuitive mapping.
 *
 * Returns:
 *  A Value type that represents the parameter and an error if anything went
 *  wrong.
 */
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


/*
 * Converts a value type to a string. Nine goes to "9", Ten to "10", Q to "Q",
 * and so on.
 *
 * Returns:
 *  A string representation of the Value. Intuitive mapping.
 */
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



/*
 * A Card represents a playing card from a standard 52 card deck. It consists of
 * a suit, such as Hearts (H), and a value such as J. The suit is represented by
 * the Suit type, and the value is a simple int that should be in the range
 * [9, 14], where 14 is A, 13 is K, and so on.
 */
type Card struct {
    Suit Suit
    Value Value
}


/*
 * Create an array of the all the cards in the euchre deck.
 */
var CARDS = createCards()
var CARDS_SET = createCardsSet()


/*
 * Creates a card given the string in the format of {V}{S}, where V is the value
 * and S is the suit.
 *
 * Args:
 *  s: The string to convert to a card. This string is in the format {V}{S}.
 *
 * Returns:
 *  A Card whose value and suit match those in the string provied. If there is
 *  any error in creating the value of string from the substrings, an error is
 *  bubbled up.
 */
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


/*
 * Converts a card to a string representation.
 *
 * Returns:
 *  A string representation of a card, which is essentially {V}{S}, where V is
 *  the value of the card and {S} is the Suit of the card.
 */
func (c Card) String() string {
    return c.Value.String() + c.Suit.String()
}


/*
 * Checks if a card is a trump card. This method accounts for the left bower
 * oddity in suits.
 *
 * Args:
 *  t: The trump suit.
 *
 * Returns:
 *  True if the card has a trump suit and false otherwise. A card is a trump if
 *  it's suit matches that of the trump suit, or it is the other bower suit and
 *  and the value of the card is J.
 */
func (c Card) IsTrump(t Suit) bool {
    return c.AdjSuit(t) == t
}


/*
 * Adjusts suit of this card based on the trump suit. This is only really
 * valuable when it matters if the card can be the left bower. In this case,
 * this method returns that the suit of this card is the trump suit. For all
 * other cards, the suit is simply outputted.
 *
 * Args:
 *  t: The trump suit.
 *
 * Returns:
 *  The effective suit of the card accounting for the left bower oddity.
 */
func (c Card) AdjSuit(t Suit) Suit {
    adjSuit := c.Suit
    if c.Value == J && c.Suit == t.Left() {
        adjSuit = t
    }

    return adjSuit
}



/*
 * Creates a set of the cards, initialized to true. This method does not give
 * back a new set each time, there is only one global set behind it all. This
 * means this method is not thread safe.
 *
 * Returns:
 *  A set (map[Card]bool) with each value initialized to true. All cards in
 *  a euchre game are in this set.
 */
func NewCardsSet() map[Card]bool {
    for k, _ := range CARDS_SET {
        CARDS_SET[k] = true
    }

    return CARDS_SET
}


/*
 * Creates a new card set with each value initialized to true.
 *
 * Returns:
 *  A new set (map[Card]bool) with each value initialized to true. All cards in
 *  a euchre game are in this set.
 */
func createCardsSet() map[Card]bool {
    set := make(map[Card]bool)
    for i := 0; i < len(CARDS); i++ {
        set[CARDS[i]] = true
    }

    return set
}

/*
 * A helper method that simply creates an array that has all the cards in a
 * euchre deck.
 *
 * Returns:
 *  A new array with the 24 cards used in euchre.
 */
func createCards() [24]Card {
    var cards [24]Card
    for i, value := range VALUES {
        for j, suit := range SUITS {
            cards[i * len(SUITS) + j] = Card { suit, value }
        }
    }

    return cards
}
