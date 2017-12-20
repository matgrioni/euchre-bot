package deck

import (
    "math/rand"
    "time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

/*
 * This package interacts with the deck definitions in order to generate cards,
 * hands and so on using the random package in go.
 */


/*
 * Generates a random card. The suit is generated randomly and so is the value.
 *
 * Returns:
 *  A random deck.Card.
 */
func Draw() Card {
    var card Card
    card.Suit = SUITS[r.Intn(4)]
    card.Value = VALUES[r.Intn(6)]

    return card
}


/*
 * Randomly generates n unique cards.
 *
 * Args:
 *  n: The number of unique cards to randomly generate.
 *
 * Returns:
 *  A slice of n random, unique cards.
 */
func DrawN(n int) []Card {
    hand := make([]Card, n)
    perm := r.Perm(len(CARDS))

    for i := 0; i < n; i++ {
        hand[i] = CARDS[perm[i]]
    }

    return hand
}
