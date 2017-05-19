package discard

import (
    "deck"
    "math/rand"
    "time"
)

// Randomly choose a card to discard, including the card you are picking up.
// hand - Your current hand.
// top - The card the player is picking up.
// Returns two values. The first is the new player hand. The second is the card
// that was chosen to be discarded.
func Rand(hand [5]deck.Card, top deck.Card) ([5]deck.Card, deck.Card) {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    total := hand[:]
    total = append(total, top)

    i := r.Intn(len(total))
    chosen := total[i]
    total[i] = total[len(total) - 1]
    total = total[:len(total) - 1]

    copy(hand[:], total[:5])

    return hand, chosen
}
