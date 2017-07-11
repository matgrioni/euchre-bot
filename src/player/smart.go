package player

import (
    "ai"
    "deck"
    "euchre"
    "fmt"
    "math"
)

type Decision struct {
    Move int
    Value int
}

type SmartPlayer struct { }

func NewSmart() (*SmartPlayer) {
    return &SmartPlayer{ }
}

func (p *SmartPlayer) Pickup(hand [5]deck.Card, top deck.Card, who int) bool {
    var copyHand [5]deck.Card
    copy(copyHand[:], hand[:])
    played := make([]deck.Card, 0)
    prior := make([]euchre.Trick, 0)

    if who == 0 {
        newHand, discard := p.Discard(copyHand, top)
        setup := euchre.Setup {
            who,
            0,
            true,
            top,
            top.Suit,
            discard,
        }

        s := euchre.NewState(setup, (who + 1) % 4, newHand[:], played, prior, deck.Card{}, 0)
        e := euchre.Engine{ }
        _, expected := ai.MCTS(s, e, 50000)

        fmt.Printf("%f\n", expected)
        return expected > 0.6
    } else {
        setup := euchre.Setup {
            who,
            0,
            true,
            top,
            top.Suit,
            deck.Card{},
        }

        s := euchre.NewState(setup, 0, hand[:], played, prior, deck.Card{}, 0)
        e := euchre.Engine{ }
        _, expected := ai.MCTS(s, e, 50000)

        fmt.Printf("%f\n", expected)
        return expected > 0.6
    }
}

// TODO: Change contract to just give back number.
func (p *SmartPlayer) Discard(hand [5]deck.Card,
                              top deck.Card) ([5]deck.Card, deck.Card) {
    // First construct a map that holds counts for all of the suits.
    suitsCount := make(map[deck.Suit]int)
    lowest := make(map[deck.Suit]int)
    for i , card := range hand {
        adjSuit := hand[i].AdjSuit(top.Suit)
        suitsCount[adjSuit]++
        if _, ok := lowest[adjSuit]; !ok {
            lowest[adjSuit] = i
        } else if int(card.Value) < int(hand[lowest[adjSuit]].Value) {
            // TODO: Isn't this an error?
            lowest[adjSuit] = i
        }
    }

    minCard := deck.Card { top.Suit, deck.J }
    minIndex := -1
    singleFound := false
    for suit, i := range lowest {
        card := hand[i]

        // If there's only one of a card that is not trump and is not an A and
        // the current min card is of greater value (it is trump or its value is
        // less), or the current min card is not the only card of its suit then
        // update the trackers.
        if suitsCount[suit] == 1 && suit != top.Suit && card.Value != deck.A &&
           (minCard.AdjSuit(top.Suit) == top.Suit ||
            int(card.Value) < int(minCard.Value) ||
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
            } else {
            // If the cards are not of the same suit then if the current min is
            // trump, the new card must be less. Otherwise, as long as the
            // card in question is not trump and it's value is less then the
            // current min card, update the trackers.
                if minCard.AdjSuit(top.Suit) == top.Suit || (card.AdjSuit(top.Suit) != top.Suit && int(card.Value) < int(minCard.Value)) {
                    minCard = card
                    minIndex = i
                }
            }
        }
    }

    hand[minIndex] = top
    return hand, minCard
}

func (p *SmartPlayer) Call(hand [5]deck.Card, top deck.Card, who int) (deck.Suit, bool) {
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

            s := euchre.NewState(setup, 0, hand[:], played, prior, deck.Card{}, 0)
            e := euchre.Engine{ }
            _, expected := ai.MCTS(s, e, 5000)

            if expected > max {
                max = expected
                maxSuit = suit
            }
        }
    }

    return maxSuit, max > 0.6
}

// TODO: Maybe don't edit the hand. Return a completely new version?
// Not sure yet.
func (p *SmartPlayer) Play(setup euchre.Setup, hand, played []deck.Card,
                           prior []euchre.Trick) ([]deck.Card, deck.Card) {
    var card deck.Card

    s := euchre.NewState(setup, 0, hand, played, prior, deck.Card{}, 0)
    e := euchre.Engine{ }
    chosenState, _ := ai.MCTS(s, e, 75000)

    card = chosenState.(euchre.State).Move

    nHand := make([]deck.Card, 0)
    for i := 0; i < len(hand); i++ {
        if card != hand[i] {
            nHand = append(nHand, hand[i])
        }
    }

    return nHand, card
}
