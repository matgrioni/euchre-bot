package euchre

import (
    "ai"
    "deck"
    "math/rand"
    "time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))


/*
 * Contains all the relevant information the setup portion of a euchre game.
 * This includes who was dealer, who called it up, what the top card was, if it
 * was picked up, what the trump suit is and if anything was discarded. Not all
 * of these values will be valid. For example, discard only makes since if the
 * top was picked up and you are the dealer, and in this case trump is not
 * necessary. However, together these 6 fields cover all possible starting
 * scenarios of interest.
 */
type Setup struct {
    Dealer int
    Caller int
    PickedUp bool
    Top deck.Card
    Trump deck.Suit
    Discard deck.Card
    AlonePlayer int
}



/*
 * A Trick in Euchre consists of the cards that were played and some context.
 * Namely, who led in the trick (using our famililar number designation), what
 * the trump suit was, and if anybody was going alone.
 */
type Trick struct {
    Cards []deck.Card
    Led int
    Trump deck.Suit
    Alone int
}



/*
 * All informative state in the euchre state tree. A state contains all
 * information prior to this moment including the move that created this current
 * state.
 */
type State struct {
    Setup Setup
    Player int
    Hands [][]deck.Card
    Played []deck.Card
    Prior []Trick
}


/*
 * Create a determinized euchre state off of an incomplete (non-terminal) euchre
 * state. Information can range from the first move to the last move, and the
 * incomplete information will be filled in through a determinization process.
 */
func (s State) Determinize() {
    cardsSet := deck.NewCardsSet()
    left := len(cardsSet)

    // Remove all prior cards from contention.
    for _, trick := range s.Prior {
        for _, card := range trick.Cards {
            cardsSet[card] = false
            left--
        }
    }

    // Remove all played cards from contention.
    for _, card := range s.Played {
        cardsSet[card] = false
        left--
    }

    // Remove all known cards of a player's hand.
    for _, card := range s.Hands[0] {
        cardsSet[card] = false
        left--
    }

    // Remove the top card from contention if it was flipped over, or remove
    // the discarded card if you were the one who put it down.
    if s.Setup.Dealer == 0 && s.Setup.PickedUp {
        cardsSet[s.Setup.Discard] = false
        left--
    } else if !s.Setup.PickedUp {
        cardsSet[s.Setup.Top] = false
        left--
    }

    // If the top card was picked up, and the top card has not already been
    // excluded due to it already being played in the current or previous
    // tricks, then remove it from contention. It can only be with the person
    // who picked it up at this moment.
    if s.Setup.PickedUp && cardsSet[s.Setup.Top] {
        cardsSet[s.Setup.Top] = false
        left--
    }

    idxs := r.Perm(left)
    cards := extractAvailableCards(cardsSet)

    noSuits := noSuits(s.Prior, s.Setup.Trump)
    // Go through each opponent player giving them cards until they have a
    // full hand given the current trick and player.
    for player := 1; player < 4; player++ {
        if s.Setup.PickedUp && s.Setup.Dealer == player {
            s.Hands[player] = append(s.Hands[player], s.Setup.Top)
        }

        playerHandSize := 5 - len(s.Prior)

        // TODO: This expression seems like it can be simplified. LM-A0
        start := ((s.Player + 4) - len(s.Played)) % 4
        dist := ((s.Player + 4) - start) % 4
        for i := start; i < start + dist; i++ {
            if i % 4 == player {
                playerHandSize--
                break
            }
        }

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
                idxs = idxs[:len(idxs) - 1]
                if len(s.Hands[player]) == playerHandSize {
                    break
                }
            }
        }
    }
}


/*
 * Creates a copy of this state. This copy is deep so the value returned is a
 * whole new state in memory with the same value as the caller.
 *
 * Returns:
 *  A new state that is identical to the current one in value, but not in
 *  memory. In other words, a deep copy!
 */
func (s State) Copy() ai.State {
    copyHands := copyAllHands(s)

    copyPlayed := make([]deck.Card, len(s.Played))
    copy(copyPlayed, s.Played)

    copyPrior := make([]Trick, len(s.Prior))
    copy(copyPrior, s.Prior)

    return State {
        s.Setup,
        s.Player,
        copyHands,
        copyPlayed,
        copyPrior,
    }
}


/*
 * Create a new state that only has the known information for any given instance
 * and is therefore undeterminized.
 *
 * Args:
 *  setup: The setup for the game. Information such as top card, dealer, etc.
 *  player: The current player number. Note that this player number must be
 *          valid as the current player number given other constraints in the
 *          parameters.
 *  hand: The current cards in your hand.
 *  played: The cards played in the current trick.
 *  prior: The prior tricks.
 *
 * Returns:
 *  A state that has undeterminized information. This means that not all fields
 *  of the state have valid values.
 */
func NewUndeterminizedState(setup Setup, player int, hand, played []deck.Card,
                            prior []Trick) State {
    // Create blank values for hands other than your own.
    hands := make([][]deck.Card, 4)
    hands[0] = hand
    hands[1] = make([]deck.Card, 0)
    hands[2] = make([]deck.Card, 0)
    hands[3] = make([]deck.Card, 0)

    return State {
        setup,
        player,
        hands,
        played,
        prior,
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
 *  played: The slice of cards played in the current trick.
 *  prior: The prior tricks.
 *
 * Returns:
 *  A new state that is determinized, ie has created some known universe where
 *  information that is usually not known by an opponent is now assumed.
 */
func NewDeterminizedState(setup Setup, player int, hands [][]deck.Card,
                          played []deck.Card, prior []Trick) State {
    return State {
        setup,
        player,
        hands,
        played,
        prior,
    }
}


/*
 * A TreeSearchEngine. This engine encapsulates all the game logic needed for
 * decision making in euchre in order to traverse the state tree.
 */
type Engine struct { }


func (engine Engine) Favorable(state ai.TSState) bool {
    cState := state.(State)
    return cState.Player % 2 == 0
}


func (engine Engine) IsTerminal(state ai.TSState) bool {
    cState := state.(State)
    return len(cState.Played) == 0 && len(cState.Prior) == 5
}


func (engine Engine) Successors(state ai.TSState) []ai.Move {
    cState := state.(State)
    nextMoves := make([]ai.Move, 0)

    curHand := cState.Hands[cState.Player]
    possibleIdxs := Possible(curHand, cState.Played, cState.Setup.Trump)

    var nPlayed []deck.Card
    var nPrior []Trick
    var nPlayer int
    nmPlayer := (cState.Player + 1) % 4

    for _, idx := range possibleIdxs {
        card := curHand[idx]
        nHands := copyAllHands(cState)

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

            cState.Played = append(cState.Played, card)

            nPlayed = make([]deck.Card, 0, 4)
            nPlayer = Winner(cState.Played, cState.Setup.Trump, nmPlayer,
                             cState.Setup.AlonePlayer)

            nextPrior := Trick {
                cState.Played,
                nmPlayer,
                cState.Setup.Trump,
                cState.Setup.AlonePlayer,
            }
            nPrior = append(nPrior, nextPrior)
        }

        nHand := make([]deck.Card, len(curHand))
        copy(nHand, curHand)
        nHand[idx] = nHand[len(nHand) - 1]
        nHand = nHand[:len(nHand) - 1]
        nHands[cState.Player] = nHand

        nextState := NewDeterminizedState(cState.Setup, nPlayer, nHands,
                                          nPlayed, nPrior)

        nextMove := ai.Move {
            card,
            nextState,
        }
        nextMoves = append(nextMoves, nextMove)
    }

    return nextMoves
}


func (engine Engine) Evaluation(state ai.TSState) float64 {
    cState := state.(State)

    winCounts0 := 0
    winCounts1 := 0

    for i := 0; i < len(cState.Prior); i++ {
        trick := cState.Prior[i]

        w := Winner(trick.Cards, cState.Setup.Trump, trick.Led,
                    cState.Setup.AlonePlayer)
        if w % 2 == 0 {
            winCounts0++
        } else {
            winCounts1++
        }
    }

    // If a player calls going alone and wins all 5 tricks then they get 4
    // points.
    if winCounts0 == 5 && cState.Setup.AlonePlayer % 2 == 0 {
        return 4
    } else if winCounts1 == 5 && cState.Setup.AlonePlayer % 2 == 1 {
        return 4
    }

    // If nobody who went alone won, but somebody won 5 hands or got euched then
    // that's two points.
    if winCounts0 == 5 || (winCounts0 >= 3 && cState.Setup.Caller % 2 == 1) {
        return 2
    } else if winCounts0 == 0 || (winCounts0 < 3 && cState.Setup.Caller % 2 == 0) {
        return -2
    }

    // For a normal win, that's worth one point.
    if winCounts0 > winCounts1 {
        return 1
    } else {
        return -1
    }
}
