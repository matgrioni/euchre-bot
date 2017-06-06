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

func (p *RandPlayer) Pickup(hand [5]deck.Card, top deck.Card, who int) bool {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    return r.Intn(2) == 1
}

func (p *RandPlayer) Discard(hand [5]deck.Card,
                             top deck.Card) ([5]deck.Card, deck.Card) {
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

func (p *RandPlayer) Call(hand [5]deck.Card, top deck.Card) (deck.Suit, bool) {
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

    return hand, final
}
