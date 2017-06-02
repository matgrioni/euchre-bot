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
