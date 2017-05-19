package call

import (
    "deck"
    "math/rand"
    "time"
)

// Randomly returns whether one should call trump and what that call should be.
// If the flag that indicates whether trump should be called is false, then the
// suit that is returned has no meaning and can be ignored.
// passed - The trump suit that was passed on the original go around and
//          therefore cannot be used as a trump suit.
// Returns two values. The first is whether trump should be called. The second
// is the actual suit. Both of these values are random.
func Rand(passed deck.Suit) (bool, deck.Suit) {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    bools := [2]bool{ true, false, }

    s := deck.SUITS[r.Intn(len(deck.SUITS))]
    for s == passed {
        s = deck.SUITS[r.Intn(len(deck.SUITS))]
    }

    return bools[r.Intn(len(bools))], s
}

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
