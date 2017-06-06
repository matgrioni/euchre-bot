package euchre

import (
    "deck"
)

// Contains all the relevant information the setup portion of a euchre game.
// This includes who was dealer, what the top card was, if it was picked up,
// what the trump suit is and if anything was discarded. Not all of these values
// will be valid. For example, discard only makes since if the top was picked up
// and you are the dealer, and in this case trump is not necessary. However,
// together these 5 fields cover all possible starting scenarios of interest.
type Setup struct {
    Dealer int
    PickedUp bool
    Top deck.Card
    Trump deck.Suit
    Discard deck.Card
}

// A Trick in Euchre consists of the cards that were played and some context.
// Namely, who led in the trick (using our famililar number designation) and
// what the trump suit was.
type Trick struct {
    Cards [4]deck.Card
    Led int
    Trump deck.Suit
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
func Possible(hand, played []deck.Card, trump deck.Suit) []int {
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
