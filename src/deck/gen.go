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
func GenCard() Card {
    var card Card
    card.Suit = SUITS[r.Intn(4)]
    card.Value = VALUES[r.Intn(6)]

    return card
}

/*
 * Randomly generates n unique cards.
 *
 * Args:
 *  n, type(int): The number of unique cards to randomly generate.
 *
 * Returns:
 *  type([]deck.Card): A slice of n random cards.
 */
func GenHand(n int) []Card {
    hand := make([]deck.Card, n)
    present := make(map[Card]bool)

    for i := range hand {
        gen = GenCard()

        // Ensure that any generated card is only included once.
        for _, in := present[gen]; in ; _, in = present[gen] {
            gen = GenCard()
        }

        present[gen] = true
        hand[i] = gen
    }

    return hand
}
