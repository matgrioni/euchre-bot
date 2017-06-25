package player

import (
    "ai"
    "deck"
    "euchre"
    "math/rand"
    "time"
)

type situation struct {
    player1 []deck.Card
    player2 []deck.Card
    player3 []deck.Card
    kitty   []deck.Card
}

type Decision struct {
    Move int
    Value int
}

type SmartPlayer struct {
}

func NewSmart() (*SmartPlayer) {
    return &SmartPlayer{ }
}

func (p *SmartPlayer) Pickup(hand [5]deck.Card, top deck.Card, who int) bool {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    return r.Intn(2) == 1
    /*limit := 100000
    i := 0

    for situation := range initials(hand, top) {
        if i >= limit {
            break
        }

        if r.Intn(2) == 0 {
            i++

            //hands := [4][]deck.Card{hand[:], situation.player1, situation.player2, situation.player3}
            //dec := minimax(hands, nil, top.Suit, 0)
        }
    }*/
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
        }

        if int(card.Value) < int(hand[lowest[adjSuit]].Value) {
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
        // less), then update the trackers.
        if suitsCount[suit] == 1 && suit != top.Suit && card.Value != deck.A &&
           (minCard.AdjSuit(top.Suit) == top.Suit || int(card.Value) < int(minCard.Value)) {
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

func (p *SmartPlayer) Call(hand [5]deck.Card, top deck.Card) (deck.Suit, bool) {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    s := deck.SUITS[r.Intn(len(deck.SUITS))]
    for s == top.Suit {
        s = deck.SUITS[r.Intn(len(deck.SUITS))]
    }

    return s, r.Intn(2) == 1
}

// TODO: Maybe don't edit the hand. Return a completely new version?
// Not sure yet.
func (p *SmartPlayer) Play(setup euchre.Setup, hand, played []deck.Card,
                           prior []euchre.Trick) ([]deck.Card, deck.Card) {
    var card deck.Card

    s := euchre.State {
        setup,
        0,
        hand,
        played,
        prior,
        deck.Card{},
    }

    n := ai.NewNode()
    e := euchre.Engine{}

    n.Value(s)
    chosenState := ai.MCTS(n, e, 10000)
    card = chosenState.Move

    nHand := make([]deck.Card, 0)
    for i := 0; i < len(hand); i++ {
        if card != hand[i] {
            nHand = append(nHand, hand[i])
        }
    }

    return nHand, card
}