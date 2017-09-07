package euchre

import (
    "ai"
    "deck"
)

/*
 * The minimax logic for a euchre playout. This functionality will not be used
 * for entire games as it requires a lot of time even for determinized games on
 * the first trick. This is more for evaluation of different methods other than
 * perfect play on a perfect information game.
 *
 * Later this code might be used however, after 2 tricks and randomly
 * determinized games since performance will not be as much of an issue.
 */


type MinimaxEngine struct { }


func (e MinimaxEngine) Favorable(state ai.MinimaxState) bool {
    cState := state.(State)
    return cState.Player % 2 == 0
}


func (e MinimaxEngine) IsTerminal(state ai.MinimaxState) bool {
    cState := state.(State)
    return len(cState.Played) == 0 && len(cState.Prior) == 5
}


func (e MinimaxEngine) Evaluation(state ai.MinimaxState) float64 {
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
        return 4.0
    } else if winCounts1 == 5 && cState.Setup.AlonePlayer % 2 == 1 {
        return 4.0
    }

    // If nobody who went alone won, but somebody won 5 hands or got euched then
    // that's two points.
    if winCounts0 == 5 || (winCounts0 >= 3 && cState.Setup.Caller % 2 == 1) {
        return 2.0
    } else if winCounts0 == 0 || (winCounts0 < 3 && cState.Setup.Caller % 2 == 0) {
        return -2.0
    }

    // For a normal win, that's worth one point.
    if winCounts0 > winCounts1 {
        return 1.0
    } else {
        return -1.0
    }
}


func (e MinimaxEngine) Successors(state ai.MinimaxState) []ai.Move {
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
                                          nPlayed, nPrior, card)

        nextMove := ai.Move {
            card,
            nextState,
        }
        nextMoves= append(nextMoves, nextMove)
    }

    return nextMoves
}
