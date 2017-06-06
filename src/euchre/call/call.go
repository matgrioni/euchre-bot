package call

import (
    "deck"
    "math/rand"
    "time"
)

// Uses a rule based approach to determine whether one should call trump and
// thereafter what the call should be.
// top - The card that was on the top, until everybody passed.
// hand - The player's current hand.
// Returns two values. The first is whether trump should be called. The second
// is the actual suit.
func Rule(top deck.Card, hand [5]deck.Card) (bool, deck.Suit) {
    chosen := false
    maxT := top.Suit
    maxConf := float32(0)

    for _, trump := range deck.SUITS {
        if trump == top.Suit {
            continue
        }

        conf := float32(0.08)

        // The following maps create a relation between a certain feature and its
        // position within the features vector.
        weights := map[deck.Value]float32 {
            deck.Nine: 0.05,
            deck.Ten: 0.07,
            deck.J: 0.3,
            deck.Q: 0.12,
            deck.K: 0.15,
            deck.A: 0.2,
        }

        // Used to keep track of how many suits are in the hand.
        suitsPresent := make(map[deck.Suit]int)

        for _, card := range hand {
            if card.Suit == trump {
                conf += weights[card.Value]
            } else if card.AdjSuit(trump) == trump {
                conf += 0.25
            } else if card.Value == deck.A {
                conf += 0.04
            }

            // Adjust suit count for left bower and increment the count for the
            // current card's suit.
            adjSuit := card.AdjSuit(trump)
            if _, ok := suitsPresent[adjSuit]; ok {
                suitsPresent[adjSuit] += 1
            } else {
                suitsPresent[adjSuit] = 1
            }
        }

        // If the hand has 2 suits or less we can have more confidence.
        if len(suitsPresent) <= 2 {
            conf += 0.08
        }

        if conf > 0.5 && conf > maxConf {
            maxConf = conf
            maxT = trump
            chosen = true
        }
    }

    return chosen, maxT
}
