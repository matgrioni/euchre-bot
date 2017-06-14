package player

import (
    "deck"
    "euchre"
    "fmt"
    "github.com/klaidliadon/next"
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

    return r.Intn(2) == 1
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
    // There are two levels of reasoning in this method. When there are 5/4
    // cards in the player's hand, there are simply too many possibilities to
    // compute. So for this amount of cards, rules will be used. Otherwise, if
    // there are 3 or fewer cards, computational power can be used to figure out
    // the best move.

    if len(hand) <= 4 {
        winners := make(map[int]int)

        for situation := range situations(setup, hand, played, prior) {
            // TODO
            hands1 := [][]deck.Card{situation.player1, situation.player2, situation.player3}
            // This checks to make sure that this situation follows the rules of
            // the game. Namely, that you must follow the suit of the first card
            // played.
            if possibleSituation(hands1, setup.Trump, prior) {
                hands := [4][]deck.Card{hand, situation.player1, situation.player2, situation.player3}
                dec := minimax(hands, played, setup.Trump, 0)
                if dec.Value == 1 {
                    winners[dec.Move]++
                }
            }
        }

        chosen := -1
        winValue := -1
        for index, count := range winners {
            fmt.Printf("%d\t%d", index, count)
            if count > winValue {
                winValue = count
                chosen = index
            }
        }

        if chosen < 0 {
            poss := euchre.Possible(hand, played, setup.Trump)
            chosen = poss[len(poss) - 1]
        }

        // TODO: Remove repetition.
        final := hand[chosen]
        hand = append(hand[:chosen], hand[chosen + 1:]...)

        return hand, final
    } else {
        w := euchre.Winner(played, setup.Trump, (4 - len(played)) % 4)

        if w == 2 && len(played) == 3 {
        // The current winner is your partner and you are the last to play so
        // throw off your lowest card.
        // TODO:
        } else if w != 2 && len(played) == 3 {
        // The current winner is not your partner and you are the last to play
        // so try to win if you can reasonably, i.e. not throwing your highest
        // trump but a low one will do.
        // TODO.
        } else {
        // TODO: Not sure what to do here.
        }

        return hand[1:], hand[0]
    }
}

// A minimax implementation that will return what card to play and the resulting
// hand after playing that card.
// TODO: Move to AI package.
func minimax(hands [4][]deck.Card, played []deck.Card, trump deck.Suit,
             player int) Decision {
    hand := hands[player]

    // If there is only one card in the hand then simply return that.
    if len(hand) == 1 {
        hand = append(hand[:0], hand[1:]...)

        return Decision {
            0,
            1,
        }
    } else {
        var chosen int
        var bestValue int

        if player == 0 {
            bestValue = -1
        } else if player == 1 || player == 3 || player == 2 {
            bestValue = 2
        }

        poss := euchre.Possible(hand, played, trump)

        for _, i := range poss {
            card := hand[i]

            // Assume the current card is played on this players turn so add it
            // to the played list for the next player. The slice must be
            // duplicated since it is a pointer, and do not want duplicate
            // references in recursive call.
            newPlayed := make([]deck.Card, len(played))
            copy(newPlayed, played)
            newPlayed = append(newPlayed, card)

            // Then remove this card from the appropriate hand.
            newHand := make([]deck.Card, len(hand))
            copy(newHand, hand)
            newHand = append(newHand[:i], newHand[i + 1:]...)
            hands[player] = newHand

            var dec Decision
            if len(played) < 4 {
                dec = minimax(hands, newPlayed, trump, (player + 1) % 4)
            } else {
                w := euchre.Winner(played, trump, (player + 1) % 4)
                v := 0
                if w == 0 {
                    v = 1
                } else if w == 1 || w == 3 || w == 2 {
                    v = 0
                }

                dec = Decision {
                    3,
                    v,
                }
            }

            // If we are the current player, then try to maximize the results.
            if player == 0 {
                if dec.Value > bestValue {
                    bestValue = dec.Value
                    chosen = i
                }
            } else if player == 1 || player == 3 || player == 2 {
            // Otherwise we will try to minimize the results.
                if dec.Value < bestValue {
                    bestValue = dec.Value
                    chosen = i
                }
            }
        }

        return Decision {
            chosen,
            bestValue,
        }
    }
}

// A generator that iterates through all possible hands or situations given a
// player's current hand, the cards currently played, and the cards played in
// previous tricks. This is a generator for now for memory purposes. For example
// if everybody has 3 cards, there are about 369000 possilibilities since the
// kitty must be taken into account.
// hand     - The current hand of the player.
// played   - The cards that have already been played on this trick.
// prior    - The cards that have been played in previous tricks.
// discard  - The card that was discarded by the player if applicable.
// dealer   - The number designation for ther person who dealt the cards.
// pickedUp - Flag to designate if the top card was picked up by the dealer.
// top      - The card that was on top of the kitty.
// TODO: How permutation works?
func situations(setup euchre.Setup, hand []deck.Card, played []deck.Card,
                tricks []euchre.Trick) chan situation {
    var prior []deck.Card
    for i := 0; i < len(tricks); i++ {
        prior = append(prior, tricks[i].Cards[:]...)
    }

    // Figure out if the top card has been played, and is thus, a known card or
    // not.
    topPlayed := false
    for _, card := range append(played, prior...) {
        if setup.Top == card {
            topPlayed = true
            break
        }
    }

    // Figure out how many unknown cards are in each player's hand and in the
    // kitty.
    nums := [4]int{len(hand), len(hand), len(hand), 4}
    for i := 2; i > 2 - len(played); i-- {
        nums[i]--
    }
    if setup.PickedUp && setup.Dealer != 0 && !topPlayed {
    // If the top was picked up by somebody else and it has not been played yet
    // we know where it is.
        nums[setup.Dealer - 1]--
    } else if (setup.PickedUp && setup.Dealer == 0) || (!setup.PickedUp) {
    // If we picked up the top card or nobody picked it up, then we know about
    // one of the 4 cards not in play.
        nums[3]--
    }

    // Create a set-like structure that has all the cards currently in play. Do
    // this by adding all cards and then removing those that are in your hand
    // have been played, and the card that was picked up if any.
    unknowns := make(map[deck.Card]bool)
    for _, card := range deck.CARDS {
        unknowns[card] = true
    }

    // We know where the top is if wasn't picked up or if somebody else picked
    // it up and it hasn't been played yet.
    if !setup.PickedUp || (setup.PickedUp && setup.Dealer != 0 && !topPlayed) {
        delete(unknowns, setup.Top)
    } else if setup.PickedUp && setup.Dealer == 0 {
    // Similarly, we know one card is (the discarded card) if we picked it up.
        delete(unknowns, setup.Discard)
    }

    // Remove all other cards that we have already seen in some way.
    for _, card := range append(hand, append(played, prior...)...) {
        delete(unknowns, card)
    }

    // Set available to the keys of the unknowns map. These are the cards whose
    // distribution is unknown.
    available := make([]deck.Card, len(unknowns), len(unknowns))
    i := 0
    for card := range unknowns {
        available[i] = card
        i++
    }

    if nums[0] + nums[1] + nums[2] + nums[3] != len(available) {
        panic("Number of freely fluxing cards is not correct.")
    }

    c := make(chan situation)
    go func() {
        // TODO
        for multi := range multinomial(nums[0], nums[1], nums[2], nums[3]) {
            cards := make([][]deck.Card, 0)
            for i := 0; i < len(multi); i++ {
                cards = append(cards, make([]deck.Card, 0))
                for j := 0; j < len(multi[i]); j++ {
                    cards[i] = append(cards[i], available[multi[i][j].(int)])
                }
            }

            next := situation {
                cards[0],
                cards[1],
                cards[2],
                cards[3],
            }

            // If the top card has not been played yet and it was picked up by
            // somebody, else add it to their cards. Basically we know of a card
            // in somebody's hand, so it wasn't freely above but should be added
            // now. Similarly, the last two ifs add a card to the kitty if we
            // know about it, i.e. it isn't picked up or we discarded a card.
            if setup.PickedUp && setup.Dealer != 0 && !topPlayed {
                switch setup.Dealer {
                case 1:
                    next.player1 = append(next.player1, setup.Top)
                case 2:
                    next.player2 = append(next.player2, setup.Top)
                case 3:
                    next.player3 = append(next.player3, setup.Top)
                }
            } else if setup.PickedUp && setup.Dealer == 0 {
                next.kitty = append(next.kitty, setup.Discard)
            } else if !setup.PickedUp {
                next.kitty = append(next.kitty, setup.Top)
            }

            c <- next
        }
        close(c)
    }()

    return c
}

func initials(hand [5]deck.Card, top deck.Card) chan situation {
    // Create a dummy setup object. TODO: Is this a clear separation of
    // interfaces.
    var tmp deck.Card
    setup := euchre.Setup {
        0,
        false,
        top,
        top.Suit,
        tmp,
    }

    return situations(setup, hand[:], nil, nil)
}

// Given a list of integer sizes for multinomial choosing, return a channel that
// gives a slice of integer slices. The sizes of the integer slices correspond
// to the varidic arguments. This is the same as providing the
// (sum(ks); ks[0], ks[1], ..., ks[n]) ways to choose ks[0], ks[1], ..., ks[n]
// integers from sum(ks) integers.
// ks - The arguments to the multinomial function.
// Returns a channel you can range over that provides a slice of slices, where
// each entry has n slices for each selection.
// TODO: Channels or generator-consumer pattern?
// TODO: Improve this runtime, probably by not relying on underlying combination
//       logic.
func multinomial(ks ...int) chan [][]interface{} {
    c := make(chan [][]interface{})

    sum := 0
    for _, k := range ks {
        sum += k
    }

    idxs := make([]interface{}, sum)
    for i := 0; i < len(idxs); i++ {
        idxs[i] = i
    }

    if len(ks) > 1 {
        go func() {
            defer close(c)

            for comb := range next.Combination(idxs, ks[0], false) {
                // TODO: Check if this assumption is right. And it can change
                // any moment probably so be careful.
                // TODO: Type assertions.
                purged := make([]interface{}, sum)
                copy(purged, idxs)
                for i := len(comb) - 1; i >= 0; i-- {
                    chosen := comb[i].(int)
                    purged = append(purged[:chosen], purged[chosen + 1:]...)
                }

                for multi := range multinomial(ks[1:]...) {
                    for i := 0; i < len(multi); i++ {
                        // Must copy this choice when reassigning indexes since
                        // modifying the slice inside of multi will modify all of
                        // that element's appeareances in later combinations
                        // since it modifies the underlying array.
                        old := multi[i]
                        multi[i] = make([]interface{}, len(multi[i]))
                        for j := 0; j < len(multi[i]); j++ {
                            multi[i][j] = purged[old[j].(int)]
                        }
                    }

                    next := make([][]interface{}, 0)
                    next = append(next, comb)
                    next = append(next, multi...)

                    c <- next
                }
            }
        }()
    } else {
        go func() {
            defer close(c)

            c <- [][]interface{}{ idxs }
        }()
    }

    return c
}

func possibleSituation(hands [][]deck.Card, trump deck.Suit,
                       prior []euchre.Trick) bool {
    noSuits := make(map[int][]deck.Suit)

    for _, trick := range prior {
        top := trick.Cards[0]
        for i, card := range trick.Cards[1:] {
            if card.AdjSuit(trump) != top.AdjSuit(trump) {
                cur := noSuits[(trick.Led + i + 1) % 4]
                cur = append(cur, card.AdjSuit(trump))
                noSuits[(trick.Led + i + 1) % 4] = cur
            }
        }
    }

    for id, suits := range noSuits {
        if id != 0 {
            for _, card := range hands[id - 1] {
                for _, suit := range suits {
                    if card.AdjSuit(trump) == suit {
                        return false
                    }
                }
            }
        }
    }

    return true
}
