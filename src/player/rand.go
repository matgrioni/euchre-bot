package player

import (
    "deck"
    "euchre"
    "math/rand"
    "time"
)

type RandPlayer struct {
}


// TODO: This might be weird semantics. player.NewRand().
// Used to create a new RandPlayer struct that is properly constructed.
// Returns a RandPlayer pointer.
func NewRand() (*RandPlayer) {
    return &RandPlayer{ }
}


func (p *RandPlayer) Pickup(hand []deck.Card, top deck.Card, who int) bool {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    return r.Intn(2) == 1
}


func (p *RandPlayer) Discard(hand []deck.Card, top deck.Card) deck.Card {
    // TODO: Move r outside.
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    hand = append(hand, top)

    // Delete a random card not preserving order.
    i := r.Intn(len(hand))
    chosen := hand[i]
    hand[i] = hand[len(total) - 1]
    hand = hand[:len(total) - 1]

    return chosen
}


func (p *RandPlayer) Call(hand []deck.Card, top deck.Card,
                          who int) (deck.Suit, bool) {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    s := deck.SUITS[r.Intn(len(deck.SUITS))]
    for s == top.Suit {
        s = deck.SUITS[r.Intn(len(deck.SUITS))]
    }

    return s, r.Intn(2) == 1
}


func (p *RandPlayer) Play(setup euchre.Setup, hand, played []deck.Card,
                          prior []euchre.Trick) ([]deck.Card, deck.Card) {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    playable := euchre.Possible(hand, played, setup.Trump)

    chosen := playable[r.Intn(len(playable))]
    final := hand[chosen]
    hand[chosen] = hand[len(hand) - 1]
    hand = hand[:len(hand) - 1]

    return final
}
