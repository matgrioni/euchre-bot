package euchre

import (
    "ai"
    "deck"
)

// Contains all the relevant information the setup portion of a euchre game.
// This includes who was dealer, who called it up, what the top card was, if it
// was picked up, what the trump suit is and if anything was discarded. Not all
// of these values will be valid. For example, discard only makes since if the
// top was picked up and you are the dealer, and in this case trump is not
// necessary. However, together these 6 fields cover all possible starting
// scenarios of interest.
type Setup struct {
    Dealer int
    Caller int
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

type State struct {
    Setup Setup
    Player int
    Hand []deck.Card
    Played []deck.Card
    Prior []Trick
    Move deck.Card
}

func (s State) Hash() interface{} {
    return s.Move
}

func NewState(setup Setup, player int, hand, played []deck.Card,
              prior []Trick, move deck.Card) State {
    return State {
        setup,
        player,
        hand,
        played,
        prior,
        move,
    }
}

// Returns whether a beats b given the current trump suit. a and b are assumed
// to be different cards. Also it is assumed a leads before b, such that if a
// and b are two different non-trump suits, a wins automatically.
// a     - The card that we are asking if it is greater.
// b     - The card that we are asking if it beats a if it is led.
// trump - The current trump suit.
// Returns if a beats b, if a is led and we are given the trump suit.
// TODO: int casting?
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
            res = int(a.Value) > int(b.Value)
        }
    } else if a.Suit == b.Suit {
    // Otherwise, if they are both the same and they are not both trump, then
    // whoever has the higher value will win.
        res = int(a.Value) > int(b.Value)
    } else {
    // And lastly if they have different suits, then a wins automatically since
    // b did not lead.
        res = true
    }

    return res
}

// Given a player's current hand and the cards that have been played, the
// possible cards for a player to play are returned. In other words, all cards
// in the player's hand that match the suit of the led card are returned or all
// cards otherwise. Also, the actual cards are not returned, rather their
// position in the hand is returned. This is to make deletion easier.
// hand   - The player's current cards.
// played - The cards that have already been played.
// trump  - The suit that is currently trump.
// Returns the index of cards that can be played according to euchre rules.
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

// A function that returns the winning player (using the same number designation
// as before) based on the trump suit, the cards that have been played, and
// what the player number is for the first player.
func Winner(played []deck.Card, trump deck.Suit, led int) int {
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

    return highPlayer
}

type Engine struct { }

func (engine Engine) Favorable(state ai.State, eval int) bool {
    cState := state.(State)
    return (cState.Player % 2 == 0 && eval > 0) ||
           (cState.Player % 2 == 1 && eval < 0)
}

func (engine Engine) IsTerminal(state ai.State) bool {
    cState := state.(State)
    return len(cState.Played) == 0 && len(cState.Prior) == 5
}

func (engine Engine) NextStates(state ai.State) []ai.State {
    cState := state.(State)
    nextStates := make([]ai.State, 0, len(deck.CARDS))
    var pCards []deck.Card

    if cState.Player == 0 {
        pIdxs := Possible(cState.Hand, cState.Played, cState.Setup.Trump)
        pCards = make([]deck.Card, len(pIdxs))

        for i, idx := range pIdxs {
            pCards[i] = cState.Hand[idx]
        }
    } else {
        noSuits := make([]deck.Suit, 0, 4)
        all := deck.NewCardsSet()
        pCards = make([]deck.Card, 0, len(all))

        for i := 0; i < len(cState.Hand); i++ {
            all[cState.Hand[i]] = false
        }

        if cState.Setup.Dealer != cState.Player || !cState.Setup.PickedUp {
            all[cState.Setup.Top] = false
        }

        if cState.Setup.Dealer == 0 && cState.Setup.PickedUp {
            all[cState.Setup.Discard] = false
        }

        for i := 0; i < len(cState.Prior); i++ {
            // For each trick, find out if a user did not follow suit and
            // therefore does not have this suit.
            trick := cState.Prior[i]
            first := trick.Cards[0]

            playedCard := trick.Cards[(4 - (cState.Player - trick.Led)) % 4]
            if first.AdjSuit(cState.Setup.Trump) != playedCard.AdjSuit(cState.Setup.Trump) {
                noSuits = append(noSuits, first.AdjSuit(cState.Setup.Trump))
            }

            for j := 0; j < len(trick.Cards); j++ {
                all[trick.Cards[j]] = false
            }
        }

        for i := 0; i < len(cState.Played); i++ {
            all[cState.Played[i]] = false
        }

        for i := 0; i < len(noSuits); i++ {
            for card, _ := range all {
                if card.AdjSuit(cState.Setup.Trump) == noSuits[i] {
                    all[card] = false
                    break
                }
            }
        }

        for card, present := range all {
            if present {
                pCards = append(pCards, card)
            }
        }
    }

    for i := 0; i < len(pCards); i++ {
        nCard := pCards[i]

        var nHand []deck.Card
        if cState.Player == 0 {
            nHand = make([]deck.Card, 0, len(cState.Hand))
            for j := 0; j < len(cState.Hand); j++ {
                if i != j {
                    nHand = append(nHand, cState.Hand[j])
                }
            }
        } else {
            nHand = cState.Hand
        }

        var nPrior []Trick
        var nPlayed []deck.Card
        var nPlayer int
        nmPlayer := (cState.Player + 1) % 4

        if len(cState.Played) < 3 {
            nPlayed = make([]deck.Card, len(cState.Played))
            copy(nPlayed, cState.Played)
            nPrior = cState.Prior
            nPlayed = append(nPlayed, nCard)
            nPlayer = nmPlayer
        } else if len(cState.Played) == 3 {
            nPrior = make([]Trick, len(cState.Prior))
            copy(nPrior, cState.Prior)

            var arrPlayed [4]deck.Card
            copy(arrPlayed[:], cState.Played)
            arrPlayed[3] = nCard

            nPlayed = make([]deck.Card, 0, 4)
            nPlayer = Winner(arrPlayed[:], cState.Setup.Trump, nmPlayer)

            if len(cState.Prior) == 4 {
                nPlayer = cState.Player
            }

            nextPrior := Trick {
                arrPlayed,
                nmPlayer,
                cState.Setup.Trump,
            }
            nPrior = append(nPrior, nextPrior)
        }

        nextState := NewState(cState.Setup, nPlayer, nHand, nPlayed,
                              nPrior, nCard)

        nextStates = append(nextStates, nextState)
    }

    return nextStates
}

func (engine Engine) Evaluation(state ai.State) int {
    winCounts0 := 0
    winCounts1 := 0

    cState := state.(State)
    for i := 0; i < len(cState.Prior); i++ {
        trick := cState.Prior[i]

        w := Winner(trick.Cards[:], cState.Setup.Trump, trick.Led)
        if w % 2 == 0 {
            winCounts0++
        } else {
            winCounts1++
        }
    }

    if winCounts0 == 5 || (winCounts0 >= 3 && cState.Setup.Caller % 2 == 1) {
        return 2
    } else if winCounts0 == 0 || (winCounts0 < 3 && cState.Setup.Caller % 2 == 0) {
        return -2
    }

    if winCounts0 > winCounts1 {
        return 1
    } else {
        return -1
    }
}
