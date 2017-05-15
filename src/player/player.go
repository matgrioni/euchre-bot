package player

import (
    "deck"
    "math/rand"
    "time"
)

// This method acts as a random player that follows all the appropriate rules of
// euchre. So it will follow suit, or in the case where it is the first player,
// it will randomly choose a card.
// hand   - The cards that are currently in the user's hand.
// played - The cards that have already been played on this trick.
// prior  - The cards that have been played in prior tricks.
// trump  - The trump suit.
// Returns the card to be played and the user's new hand.
func Random(hand []deck.Card, played []deck.Card, prior []deck.Card,
            trump deck.Suit) (deck.Card, []deck.Card) {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    var possible []int
    if len(played) > 0 {
        for i := range hand {
            if hand[i].AdjSuit(trump) == played[0].AdjSuit(trump) {
                possible = append(possible, i)
            }
        }
    }

    var chosen int
    if len(possible) > 0 {
        chosen = possible[r.Intn(len(possible))]
    } else {
        chosen = r.Intn(len(hand))
    }

    final := hand[chosen]
    hand[chosen] = hand[len(hand) - 1]
    hand = hand[:len(hand) - 1]

    return final, hand
}
