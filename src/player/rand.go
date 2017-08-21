package player

import (
    "deck"
    "euchre"
    "math/rand"
    "time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))


/*
 * A RandomPlayer. This player is non determinisitic, so for the same input, you
 * can get different outputs. The distribution of picking up, or discarding can
 * are configurable. Discarding or playing is not since they would naturally
 * have to depend on card order, which does not necessairily have any meaning.
 */
type RandPlayer struct {
    pickupProb float64
    callProb float64
}


/*
 * Used to create a new RandPlayer struct that is properly constructed. Returns
 * a RandPlayer pointer.
 *
 * Args:
 *  pickupProb: The probability that the player says the top card should be
 *              picked up.
 *  callProb: The probability that the player will call the suit after everybody
 *            skips the first round.
 */
func NewRand(pickupProb float64, callProb float64) (*RandPlayer) {
    return &RandPlayer{ }
}


/*
 * Player decides to pickup with probability pickupProb.
 *
 * Args:
 *  hand: The current hand of the player.
 *  top: The card on top of the kitty.
 *  who: The player number designation for the dealer.
 *
 * Returns:
 *  True if the player should tell the dealer to pickup the card. False if pass.
 *  The player decides the dealer should pickup it up with probability
 *  pickupProb.
 */
func (p *RandPlayer) Pickup(hand []deck.Card, top deck.Card, who int) bool {
    return r.Float64() < p.pickupProb
}


func (p *RandPlayer) Discard(hand []deck.Card, top deck.Card) ([]deck.Card, deck.Card) {
    hand = append(hand, top)

    // Delete a random card not preserving order.
    i := r.Intn(len(hand))
    chosen := hand[i]
    hand[i] = hand[len(hand) - 1]
    hand = hand[:len(hand) - 1]

    return hand, chosen
}


/*
 * Player decision method to call suit once all other players have passed on
 * telling the dealer to pickup.
 *
 * Args:
 *  hand: The hand of the player.
 *  top: The card on the top of the kitty.
 *  who: The player number designation for the dealer.
 *
 * Returns:
 *  The suit to call and True if the player should call this suit. False if the
 *  player should not call it. The player decides to call the suit with
 *  probability callProb.
 */
func (p *RandPlayer) Call(hand []deck.Card, top deck.Card,
                          who int) (deck.Suit, bool) {
    s := deck.SUITS[r.Intn(len(deck.SUITS))]
    for s == top.Suit {
        s = deck.SUITS[r.Intn(len(deck.SUITS))]
    }

    return s, r.Float64() < p.callProb
}


func (p *RandPlayer) Play(setup euchre.Setup, hand, played []deck.Card,
                          prior []euchre.Trick) ([]deck.Card, deck.Card) {
    playable := euchre.Possible(hand, played, setup.Trump)

    chosen := playable[r.Intn(len(playable))]
    final := hand[chosen]
    hand[chosen] = hand[len(hand) - 1]
    hand = hand[:len(hand) - 1]

    return hand, final
}
