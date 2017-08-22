package euchre

import (
    "deck"
)


/*
 * Returns whether a beats b given the current trump suit. a and b are assumed
 * to be different cards. Also it is assumed a leads before b, such that if a
 * and b are two different non-trump suits, a wins automatically.
 *
 * Args:
 *  a: The card that we are asking if it is greater.
 *  b: The card that we are asking if it beats a if it is led.
 *  trump: The current trump suit.
 *
 * Returns:
 *  True if a beats b, if a is led and we are given the trump suit.
 */
func Beat(a deck.Card, b deck.Card, trump deck.Suit) bool {
    var res bool
    // If a is a trump card but b is not, then a wins.
    if (a.AdjSuit(trump) == trump && b.AdjSuit(trump) != trump) ||
       (a.AdjSuit(trump) != trump && b.AdjSuit(trump) == trump) {
        res = a.AdjSuit(trump) == trump
    } else if a.AdjSuit(trump) == trump && b.AdjSuit(trump) == trump {
    // If a is a trump and so is b, then we must compare their values knowing
    // that right and left bower are a rule.
        if a.Value == deck.J || b.Value == deck.J {
            // If a is right bower, then it must win.
            if a.Value == deck.J && a.Suit == trump {
                res = true
            } else if a.Value == deck.J && a.Suit == trump.Left() {
            // If a is left bower, then it wins as long as b is not the right
            // bower.
                res = b.Value != deck.J
            } else {
            // Otherwise, a is not a J, so it is b so b must win.
                res = false
            }
        } else {
        // If neither are one of the bowers, then the values of the cards are
        // compared as normal.
            res = a.Value.Compare(b.Value) > 0
        }
    } else if a.Suit == b.Suit {
    // Otherwise, if they are both the same and they are not both trump, then
    // whoever has the higher value will win.
        res = a.Value.Compare(b.Value) > 0
    } else {
    // And lastly if they have different suits, then a wins automatically since
    // b did not lead.
        res = true
    }

    return res
}


/*
 * A function that returns the winning player (using the same number designation
 * as before) based on the trump suit, the cards that have been played, and
 * what the player number is for the first player.
 *
 * Args:
 *  played: The cards that were played.
 *  trump: The trump suit.
 *  led: The number designation of the person who played the first card.
 *
 * Returns:
 *  The number designation of the person who won the trick.
 */
func Winner(played []deck.Card, trump deck.Suit, led int) int {
    return winner(played, trump, led, -1)
}


/*
 * A function that returns the winning player (using the same number designation
 * as before) based on the trump suit, the cards that have been played, and
 * what the player number is for the first player. This method assumes that one
 * player has gone alone. Also provide what player is going alone.
 *
 * Args:
 *  played: The cards that were played.
 *  trump: The trump suit.
 *  led: The number designation of the person who played the first card.
 *  alone: The number designation of the player going alone.
 *
 * Returns:
 *  The number designation of the person who won the trick.
 */
func WinnerAlone(played []deck.Card, trump deck.Suit, led int, alone int) int {
    return winner(played, trump, led, alone)
}


/*
 * A function that returns the winning player (using the same number designation
 * as before) based on the trump suit, the cards that have been played, what the
 * player number is for the first player, and if anybody went alone. This method
 * is for internal purposes. The caller chooses between either Winner or
 * WinnerAlone.
 *
 * Args:
 *  played: The cards that were played.
 *  trump: The trump suit.
 *  led: The number designation of the person who played the first card.
 *  alone: The alone player if there is any. If there is not then put in an
 *         invalid player number.
 *
 * Returns:
 *  The number designation of the person who won the trick.
 */
func winner(played []deck.Card, trump deck.Suit, led int, alone int) int {
    highPlayer := led

    if len(played) >= 2 {
        highest := played[0]
        for i, card := range played[1:] {
            if !Beat(highest, card, trump) {
                highest = card
                highPlayer = (led + i + 1) % 4
            }
        }
    }

    nextAfterHigh := (highPlayer + 1) % 4
    if alone >= 0 && alone < 4 {
        for player := led; player != nextAfterHigh; player = (player + 1) % 4 {
            // If somewhere between the leading player, and the highPlayer the
            // player who is cucked by their partner calling going alone is
            // found, then the highPlayer should be moved up one, since one
            // player is not actually playing.
            if player == (alone + 2) % 4 {
                highPlayer = (highPlayer + 1) % 4
                break
            }
        }
    }

    return highPlayer
}


/*
 * Given a player's current hand and the cards that have been played, the
 * possible cards for a player to play are returned. In other words, all cards
 * in the player's hand that match the suit of the led card are returned or all
 * cards otherwise. Also, the actual cards are not returned, rather their
 * position in the hand is returned. This is to make deletion easier.
 *
 * Args:
 *  hand: The player's current cards.
 *  played: The cards that have already been played.
 *  trump: The suit that is currently trump.
 *
 * Returns:
 *  The index of cards in hand that can be played according to euchre rules.
 */
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
