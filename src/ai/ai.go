package ai

import "deck"

// A rule based approach to deciding if one should call for the dealer to pick
// up the top card at the start of a deal. This approach takes into account the
// cards currently in the players hand, the top card, and whether one or one's
// partner picks up the card in question.
// TODO: Special edge case where you can drop the only card of a certain suit.
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
    // action. Start off with a little confidence that your partner will win at
    // least one trick.
    var conf float32 = 0.08

    // Used to keep track of how many suits are in the hand.
    suitsPresent := make(map[deck.Suit]int)

    for _, card := range hand {
        if card.Suit == top.Suit {
            // Use the normal trump weights to increase confidence of winning a
            // trick for each trump you have.
            conf += weights[card.Value]
        } else if card.Suit.Left() == top.Suit && card.Value == deck.J {
            // Don't forget the left bower if you have it.
            conf += 0.25
        } else if card.Value == deck.A {
            // Then for every A, one has, it slightly increases the chances of
            // winning said trick.
            conf += 0.04
        }

        suitsPresent[card.Suit] = 1
    }

    // Count the actual number of suits.
    suitCount := 0
    for _, present := range suitsPresent {
        suitCount += present
    }

    // Increase the confidence if you only have 2 suits as well, as that gives
    // more ability to use trump. This will only make a difference if you
    // already have a lot of trump cards that gets you close to 0.5 but not over
    // it.
    if suitCount <= 2 {
        conf += 0.08
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
