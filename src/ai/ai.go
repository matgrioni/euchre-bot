package ai

import "deck"

// A rule based approach to deciding if one should call for the dealer to pick
// up the top card at the start of a deal. This approach takes into account the
// cards currently in the players hand, the top card, and whether one or one's
// partner picks up the card in question. The friendly parameter can take the
// following values. 2 for you are picking it up, 1 for your partner is picking
// it up, and 0 neither.
func RPickUp(hand [5]deck.Card, top deck.Card, friendly int, thres float32) float32 {
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

        // Adjust suit count for left bower.
        adjSuit := card.Suit
        if card.Value == deck.J && card.Suit == top.Suit.Left() {
            adjSuit = top.Suit
        }

        if _, ok := suitsPresent[adjSuit]; ok {
            suitsPresent[adjSuit] += 1
        } else {
            suitsPresent[adjSuit] = 1
        }
    }

    // Increase the confidence if you only have 2 suits as well, as that gives
    // more ability to use trump. This will only make a difference if you
    // already have a lot of trump cards that gets you close to 0.5 but not over
    // it. You can reach two suits either by already having them or by, picking
    // up the top card and removing one card if it is a friendly deal (to
    // yourself).
    suitCount := len(suitsPresent)
    if suitCount <= 2 && friendly != 2 {
        conf += 0.08
    } else if suitCount <= 3 && friendly == 2 {
        _, trumpPresent := suitsPresent[top.Suit]
        if suitCount <= 2 && trumpPresent {
            conf += 0.08
        } else if suitCount == 3 && trumpPresent {
            for _, card := range hand {
                // If this is the only card in the hand of a given suit, and it
                // is not a trump (left bower or of the same suit as the top
                // card), and it is not an A, then by removing it we get rid of
                // one suit. So if there were only 3 suits, then this means we
                // have 2 suits after discarding if we were to pick up. We also
                // have to make sure that this removal will result in at 2 suits
                // and that picking up the top card does not increase the suit
                // amount.
                if suitsPresent[card.Suit] == 1 && card.Suit != top.Suit &&
                   (card.Suit != top.Suit.Left() || card.Value != deck.J) &&
                   card.Value != deck.A {
                    conf += 0.08

                    // There could be more than one card that matches these
                    // requirements, we only care if one exists, so break on finding
                    // one.
                    break
                }
            }
        }
    }

    // If you are picking it up then consider the weight of the card on top as
    // well. The discarded card need not be considered since, it is most likely
    // non-trump, and will not affect this analysis which only considers trump
    // cards. If it is a trump card, then every other card is a trump as well
    // and the final decision should not be affected.
    if friendly != 0 {
        conf += weights[top.Value]
    }

    return conf >= thres
}
