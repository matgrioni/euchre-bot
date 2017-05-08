package ai

import "deck"

// A rule based approach to deciding if one should call for the dealer to pick
// up the top card at the start of a deal. This approach takes into account the
// cards currently in the players hand, the top card, and whether one picks up
// the card in question.
func RPickUp(hand [5]deck.Card, top deck.Card, friendly bool) bool {
    // Arbitrary weights that I felt like using based on experience.
    weights := map[deck.Value]float32 {
        deck.Nine: 0.05,
        deck.Ten: 0.07,
        deck.J: 0.3,
        deck.Q: 0.12,
        deck.K: 0.15,
        deck.A: 0.2,
    }

    // One wants to tell the dealer to pick up if you believe that you can win 3
    // tricks if the dealer picks it up. So if your confidence / probability is
    // greater than 50% (in a non-technical sense) then you should take such an
    // action.
    var conf float32 = 0

    for i := range hand {
        card := hand[i]
        if card.Suit == top.Suit {
            conf += weights[card.Value]
        } else if card.Suit.Left() == top.Suit && card.Value == deck.J {
            conf += 0.25
        }
    }

    // If you are picking it up then consider the weight of the card onn top as
    // well. The discarded card need not be considered since, it is most likely
    // non-trump, and will not affect this analysis which only considers trump
    // cards. If it is a trump card, then every other card is a trump as well
    // and the final decision should not be affected.
    if friendly {
        conf += weights[top.Value]
    }

    return conf >= 0.5
}
