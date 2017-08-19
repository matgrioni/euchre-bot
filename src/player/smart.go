package player

import (
    "ai"
    "deck"
    "euchre"
    "math"
)



type SmartPlayer struct {
    pickupConfidence float64
    callConfidence float64

    pickupRuns int
    pickupDeterminizations int

    callRuns int
    callDeterminizations int

    playRuns int
    playDeterminizations int
}


/*
 * TODO
 *
 * Args:
 *  pickupConfidence:
 *  callConfidence:
 *  pickupRuns:
 *  pickupDeterminizations:
 *  callRuns:
 *  callDeterminizations:
 *  playRuns:
 *  playDeterminizations:
 *
 * Returns:
 *  TODO
 */
func NewSmart(pickupConfidence float64, callConfidence float64,
              pickupRuns int, pickupDeterminizations int,
              callRuns int, callDeterminizations int,
              playRuns int, playDeterminizations int) (*SmartPlayer) {

    return &SmartPlayer{
        pickupConfidence,
        callConfidence,
        pickupRuns,
        pickupDeterminizations,
        callRuns,
        callDeterminizations,
        playRuns,
        playDeterminizations,
    }
}


func (p *SmartPlayer) Pickup(hand []deck.Card, top deck.Card, who int) bool {
    var setup euchre.Setup
    var discard deck.Card

    actualHand := make([]deck.Card, len(hand))

    played := make([]deck.Card, 0)
    prior := make([]euchre.Trick, 0)

    if who == 0 {
        copy(actualHand, hand)

        actualHand, discard = p.Discard(actualHand, top)
        setup = euchre.Setup {
            who,
            0,
            true,
            top,
            top.Suit,
            discard,
        }
    } else {
        copy(actualHand, hand)
        setup = euchre.Setup {
            who,
            0,
            true,
            top,
            top.Suit,
            deck.Card{ },
        }
    }

    nPlayer := (who + 1) % 4
    s := euchre.NewUndeterminizedState(setup, nPlayer, actualHand, played,
                                       prior, deck.Card{})
    e := euchre.Engine{ }
    _, expected := ai.MCTS(s, e, p.pickupRuns, p.pickupDeterminizations)

    return (nPlayer % 2 == 0 && expected > p.pickupConfidence) ||
           (nPlayer % 2 == 1 && expected < -1 * p.pickupConfidence)
}


func (p *SmartPlayer) Discard(hand []deck.Card,
                              top deck.Card) ([]deck.Card, deck.Card) {
    // First construct a map that holds counts for all of the suits and the
    // lowest card for each suit.
    suitsCount := make(map[deck.Suit]int)
    lowest := make(map[deck.Suit]int)

    hand = append(hand, top)
    for i, card := range hand {
        adjSuit := hand[i].AdjSuit(top.Suit)
        suitsCount[adjSuit]++

        if _, ok := lowest[adjSuit]; !ok {
            lowest[adjSuit] = i
        } else {
            // If the lowest card for a given suit can beat the current card,
            // the current card is smaller.
            curLowest := hand[lowest[adjSuit]]
            if euchre.Beat(curLowest, card, top.Suit) {
                lowest[adjSuit] = i
            }
        }
    }

    singleFound := false
    trumpExists := suitsCount[top.Suit] > 1
    minCard := deck.Card { top.Suit, deck.J }
    minIndex := -1

    for suit, i := range lowest {
        card := hand[i]

        // If there's only one of a card that is not trump and is not an A and
        // the current min card is of greater value (it is trump or its value is
        // less), or the current min card is not the only card of its suit then
        // update the trackers.
        if trumpExists && suitsCount[suit] == 1 && suit != top.Suit &&
           card.Value != deck.A && (minCard.IsTrump(top.Suit) ||
           deck.ValueCompare(card, minCard) < 0 ||
           suitsCount[minCard.Suit] > 1) {
            singleFound = true
            minCard = card
            minIndex = i
        } else if !singleFound {
        // If a single card that is non-trump and non-ace has not been found
        // then try to find the smallest card otherwise. In other words, any
        // single suit card is preferred to a multi-suit card as long as said
        // card is not trump or A.

            // If the two cards of the same suit we can compare them using Beat
            // since it wouldn't matter who led.
            if card.AdjSuit(top.Suit) == minCard.AdjSuit(top.Suit) {
                if euchre.Beat(minCard, card, top.Suit) {
                    minCard = card
                    minIndex = i
                }
            } else if minCard.IsTrump(top.Suit) || (!card.IsTrump(top.Suit) &&
                      deck.ValueCompare(card, minCard) < 0) {
            // If the cards are not of the same suit then if the current min is
            // trump, the new card must be less. Otherwise, as long as the
            // card in question is not trump and it's value is less then the
            // current min card, update the trackers.
                minCard = card
                minIndex = i
            }
        }
    }

    hand[minIndex] = top
    hand = hand[:len(hand) - 1]

    return hand, minCard
}


func (p *SmartPlayer) Call(hand []deck.Card, top deck.Card,
                           who int) (deck.Suit, bool) {
    played := make([]deck.Card, 0)
    prior := make([]euchre.Trick, 0)
    max := math.Inf(-1)
    var maxSuit deck.Suit

    for i := 0; i < len(deck.SUITS); i++ {
        suit := deck.SUITS[i]

        if suit != top.Suit {
            setup := euchre.Setup {
                who,
                0,
                false,
                top,
                suit,
                deck.Card{},
            }

            s := euchre.NewUndeterminizedState(setup, 0, hand, played, prior,
                                               deck.Card{})
            e := euchre.Engine{ }
            _, expected := ai.MCTS(s, e, p.callRuns, p.callDeterminizations)

            if expected > max {
                max = expected
                maxSuit = suit
            }
        }
    }

    return maxSuit, max > p.callConfidence
}


func (p *SmartPlayer) Play(setup euchre.Setup, hand, played []deck.Card,
                           prior []euchre.Trick) ([]deck.Card, deck.Card) {
    var card deck.Card

    s := euchre.NewUndeterminizedState(setup, 0, hand, played, prior,
                                       deck.Card{})
    e := euchre.Engine{ }
    chosenState, _ := ai.MCTS(s, e, p.playRuns, p.playDeterminizations)

    card = chosenState.(euchre.State).Move

    nHand := make([]deck.Card, 0)
    for i := 0; i < len(hand); i++ {
        if card != hand[i] {
            nHand = append(nHand, hand[i])
        }
    }

    return nHand, card
}
