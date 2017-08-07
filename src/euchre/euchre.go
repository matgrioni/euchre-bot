package euchre

import (
    "ai"
    "deck"
    "math/rand"
    "time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

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
    Hands [][]deck.Card
    Kitty []deck.Card
    Played []deck.Card
    Prior []Trick
    Move deck.Card
}


func (s State) Hash() interface{} {
    return s.Move
}

/*
 * Create a determinized euchre state off of an incomplete (non-terminal) euchre
 * state. Information can range from the first move to the last move, and the
 * incomplete information will be filled in through a determinization process.
 */
func (s State) Determinize() {
    cardsSet := deck.NewCardsSet()

    // Remove all prior cards from contention.
    for _, trick := range s.Prior {
        for _, card := range trick.Cards {
            cardsSet[card] = false
        }
    }

    // Remove all played cards from contention.
    for _, card := range s.Played {
        cardsSet[card] = false
    }

    // Remove the top card from contention if it was flipped over, or remove
    // the discarded card if you were the one who put it down.
    if s.Setup.Dealer == 0 && s.Setup.PickedUp {
        cardsSet[s.Setup.Discard] = false
    } else if !s.Setup.PickedUp {
        cardsSet[s.Setup.Top] = false
    }

    idxs := r.Perm(len(cardsSet))
    cards := extractAvailableCards(cardsSet)

    noSuits := noSuits(s.Prior, s.Setup.Trump)
    // Go through each opponent player.
    for player := 1; player < 4; player++ {

        // Go through each card in the random permutation. Go backwards so that
        // the idxs can be deleted as cards are chosen as valid.
        for i := len(idxs) - 1; i >= 0; i-- {
            cardIndex := idxs[i]
            curCard := cards[cardIndex]

            // If the current card is possible for the current player, then add
            // it to the current player's hand.
            possible := true
            for _, suit := range noSuits[player] {
                if curCard.AdjSuit(s.Setup.Trump) == suit {
                    possible = false
                }
            }

            // If the current card is possible by the rules of following suit,
            // then add it. If we have all the cards for this player, then move
            // on to the next as well. Also remove, the card index so that it
            // isn't included in further players.
            if possible {
                s.Hands[player] = append(s.Hands[player], curCard)
                idxs = append(idxs[:i], idxs[:i+1]...)
                if len(s.Hands[player]) == 5 {
                    break
                }
            }
        }
    }
}


/*
 * Create a new state that only has the known information for any given instance
 * and is therefore undeterminized.
 *
 * Args:
 *  setup: The setup for the game. Information such as top card, dealer, etc.
 *  player: The current player number.
 *  hand: The current cards in your hand.
 *  played: The cards played in the current trick.
 *  prior: The prior tricks.
 *  move: The prior card played to get to this state.
 *
 * Returns:
 *  A state that has undeterminized information. This means that not all fields
 *  of the state have valid values.
 */
func NewUndeterminizedState(setup Setup, player int, hand, played []deck.Card,
                            prior []Trick, move deck.Card) State {
    // Create blank values for the kitty and other hands other than your own.
    kitty := make([]deck.Card, 0)
    hands := make([][]deck.Card, 4)
    hands[0] = hand
    hands[1] = make([]deck.Card, 0)
    hands[2] = make([]deck.Card, 0)
    hands[3] = make([]deck.Card, 0)

    return State {
        setup,
        player,
        hands,
        kitty,
        played,
        prior,
        move,
    }
}


/*
 * Create a new determinized state. This state may have missing information,
 * such as in the middle of euchre game. This is fine as long as the structure
 * of the incomplete information is consistent with actual game play. For
 * example, players can not have two different card counts, except when it is
 * the middle of a turn where they differ by one card.
 *
 * Args:
 *  setup: The setup for the game. Information such as top card, dealer, etc.
 *  player: The current player number.
 *  hands: A slice of each players cards.
 *  kitty: The slice of cards in the kitty. TODO: This might not be needed.
 *  played: The slice of cards played in the current trick.
 *  prior: The prior tricks.
 *  move: The move that was played prior to get the current state.
 *
 * Returns:
 *  A new state that is determinized, ie has created some known universe where
 *  information that is usually not known by an opponent is now assumed.
 */
func NewDeterminizedState(setup Setup, player int, hands [][]deck.Card, kitty,
                          played []deck.Card, prior []Trick,
                          move deck.Card) State {
    return State {
        setup,
        player,
        hands,
        kitty,
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
    nextStates := make([]ai.State, 0)

    curHand := cState.Hands[cState.Player]

    var nPlayed []deck.Card
    var nPrior []Trick
    var nPlayer int
    nmPlayer := (cState.Player + 1) % 4

    for i, card := range curHand {
        if len(cState.Played) < 3 {
            // Copy the old played cards into memory and add the new card.
            nPlayed = make([]deck.Card, len(cState.Played))
            copy(nPlayed, cState.Played)
            nPlayed = append(nPlayed, card)

            // The prior tricks stay the same (no new tricks) and the next
            // player is just the next modulo player.
            nPrior = cState.Prior
            nPlayer = nmPlayer
        } else if len(cState.Played) == 3 {
            // If this next card ends the trick then copy the old tricks over
            // and make a new trick out of the current cards.
            nPrior = make([]Trick, len(cState.Prior))
            copy(nPrior, cState.Prior)

            var arrPlayed [4]deck.Card
            copy(arrPlayed[:], cState.Played)
            arrPlayed[3] = card

            nPlayed = make([]deck.Card, 0, 4)
            nPlayer = Winner(arrPlayed[:], cState.Setup.Trump, nmPlayer)

            // If this is the last trick to be played, then don't move the
            // player to the next one. This is because when we evaluate the
            // final state, we want it to be in reference to this last player,
            // not the next one. TODO: Actually I don't really know if this
            // makes a difference anymore.
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

        nHand := make([]deck.Card, len(curHand))
        nHand = append(nHand[:i], nHand[i+1:]...)
        nHands := cState.Hands
        nHands[cState.Player] = nHand

        nextState := NewDeterminizedState(cState.Setup, nPlayer, nHands,
                                          cState.Kitty, nPlayed, nPrior, card)

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


/*
 * A private helper method that returns what suits a given player cannot have.
 *
 * Args:
 *  prior: The list of prior tricks.
 *  trump: The current trump suit.
 *
 * Returns:
 *  A list of the suits that a player cannot have indexed by player numbers
 *  in a map.
 */
func noSuits(prior []Trick, trump deck.Suit) map[int][]deck.Suit {
    noSuits := make(map[int][]deck.Suit)

    for i := 0; i < len(prior); i++ {
        // For each trick, find out if a user did not follow suit and therefore
        // does not have this suit.
        trick := prior[i]
        first := trick.Cards[0]

        for player := 0; player < 4; player++ {
            // Add 4 to player, so that it is guaranteed to be after trick.Led,
            // but it does not change the final result mod 4.
            playedCard := trick.Cards[(player + 4 - trick.Led) % 4]
            if first.AdjSuit(trump) != playedCard.AdjSuit(trump) {
                noSuits[player] = append(noSuits[player], first.AdjSuit(trump))
            }
        }
    }

    return noSuits
}


/*
 * Extracts all cards that exist in the given set and have a value of true, into
 * a list of cards.
 *
 * Args:
 *  cardsSet: A set of cards. A card exists if it is in the set and its value is
 *            true.
 *
 * Returns:
 *  A slice of the existing cards in the given set.
 */
func extractAvailableCards(cardsSet map[deck.Card]bool) []deck.Card {
    cards := make([]deck.Card, 0, len(cardsSet))
    for card, exists := range cardsSet {
        if exists {
            cards = append(cards, card)
        }
    }

    return cards
}
