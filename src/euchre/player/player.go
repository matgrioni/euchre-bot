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
// trump  - The trump suit.
// Returns the card to be played and the user's new hand.
func Random(hand []deck.Card, played []deck.Card, trump deck.Suit) (deck.Card,
            []deck.Card) {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    playable := possible(hand, played, trump)
    chosen := playable[r.Intn(len(playable))]
    final := hand[chosen]
    hand = append(hand[:chosen], hand[chosen + 1:]...)

    return final, hand
}

// Given a player's current hand and the cards that have been played, the
// possible cards for a player to play are returned. In other words, all cards
// in the player's hand that match the suit of the led card are returned or all
// cards otherwise. Also, the actual cards are not returned, rather their
// position in the hand is returned. This is to make deletion easier.
// hand   - The player's current cards.
// played - The cards that have already been played.
// trump  - The suit that is currently trump.
// Returns the index of cards that can be played according to euchre rules.
func possible(hand []deck.Card, played []deck.Card, trump deck.Suit) []int {
    possible := make([]int, 0, len(hand))
    if len(played) > 0 {
        for i := range hand {
            if hand[i].AdjSuit(trump) == played[0].AdjSuit(trump) {
                possible = append(possible, i)
            }
        }
    }

    if len(possible) == 0 {
        for i := range hand {
            possible = append(possible, i)
        }
    }

    return possible
}
