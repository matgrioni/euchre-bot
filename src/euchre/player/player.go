package player

import (
    "deck"
    "github.com/klaidliadon/next"
    "math/rand"
    "fmt"
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

// This method acts as a random player that follows all the appropriate rules of
// euchre. So it will follow suit, or in the case where it is the first player,
// it will randomly choose a card.
// hand   - The cards that are currently in the user's hand.
// played - The cards that have already been played on this trick.
// trump  - The trump suit.
// Returns the card to be played and the user's new hand.
func Random(hand []deck.Card, played []deck.Card, trump deck.Suit) (deck.Card,
            []deck.Card) {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    playable := possible(hand, played, trump)
    chosen := playable[r.Intn(len(playable))]
    final := hand[chosen]
    hand = append(hand[:chosen], hand[chosen + 1:]...)

    return final, hand
}

// This method acts as a combintorially smart player (uses the power of machine
// power).
// hand     - The cards that are currently in the user's hand.
// played   - The cards that have already been played on this trick.
// prior    - The cards that have been played in prior tricks.
// dealer   - Who dealt. 0 is for self, then moves in clockwise directions so 2
//            is partner, 1 is opponent to the left, and 3 is opponent to the
//            right.
// discard  - The card that was discarded if you were the dealer and were told to
//            pick it up. This value can simply be nil if there is no suitable
//            value.
// pickedUp - Whether the top card was picked up or not by the who above.
// top      - The card on top of the kitty after dealing.
// trump    - The eventual suit that is chosen as trump.
// Returns the card to be played and the user's new hand.
func AI(hand []deck.Card, played []deck.Card, prior []deck.Card,
        discard deck.Card, dealer int, pickedUp bool, top deck.Card,
        trump deck.Suit) (deck.Card, []deck.Card) {
    // There are two levels of reasoning in this method. When there are 5/4
    // cards in the player's hand, there are simply too many possibilities to
    // compute. So for this amount of cards, rules will be used. Otherwise, if
    // there are 3 or fewer cards, computational power can be used to figure out
    // the best move.

    if len(hand) <= 3 {
        winners := make(map[int]int)

        for situation := range situations(hand, played, prior, discard, dealer, pickedUp, top) {
            hands := [4][]deck.Card{hand, situation.player1, situation.player2, situation.player3}
            dec := minimax(hands, played, trump, 0)
            if dec.Value == 1 {
                winners[dec.Move]++
            }
        }

        chosen := 0
        winValue := -1
        for index, count := range winners {
            fmt.Printf("%d\t%d\n", index, count)
            if winValue > count {
                winValue = count
                chosen = index
            }
        }

        for _, card := range hand {
            fmt.Print(card)
            fmt.Print(" ")
        }
        fmt.Println()
        fmt.Println(chosen)

        // TODO: Remove repetition.
        final := hand[chosen]
        hand = append(hand[:chosen], hand[chosen + 1:]...)

        return final, hand
    } else {
        r := rand.New(rand.NewSource(time.Now().UnixNano()))

        playable := possible(hand, played, trump)
        chosen := playable[r.Intn(len(playable))]
        final := hand[chosen]
        hand = append(hand[:chosen], hand[chosen + 1:]...)

        return final, hand
    }
}

// A minimax implementation that will return what card to play and the resulting
// hand after playing that card.
func minimax(hands [4][]deck.Card, played []deck.Card, trump deck.Suit,
             player int) Decision {
    hand := hands[player]

    // If there is only one card in the hand then simply return that.
    if len(hand) == 1 {
        hand = append(hand[:0], hand[1:]...)

        // TODO: what is played here?
        w := winner(played, trump, (player + 1) % 4)
        v := 0
        if w == 0 || w == 2 {
            v = 1
        } else if w == 1 || w == 3 {
            v = 0
        }

        return Decision {
            0,
            v,
        }
    } else {
        var chosen int
        var bestValue int

        if player == 0 {
            bestValue = -1
        } else if player == 1 || player == 3 || player == 2 {
            bestValue = 2
        }

        poss := possible(hand, played, trump)
        if player == 0 {
            for _, p := range poss {
                fmt.Print(p)
                fmt.Print(" ")
            }
        }
        fmt.Println()

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
                w := winner(played, trump, (player + 1) % 4)
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

// A function that returns the winning player (using the same number designation
// as before) based on the trump suit, the cards that have been played, and
// what the player number is for the first player.
func winner(played []deck.Card, trump deck.Suit, led int) int {
    highest := played[0]
    highPlayer := led
    for i, card := range played[1:] {
        if Beat(highest, card, trump) {
            highest = card
            highPlayer = (led + i + 1) % 4
        }
    }

    return highPlayer
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
    if a.AdjSuit(trump) == trump && b.AdjSuit(trump) != trump {
        res = true
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
func situations(hand []deck.Card, played []deck.Card, prior []deck.Card,
                discard deck.Card, dealer int, pickedUp bool,
                top deck.Card) chan situation {
    // Figure out if the top card has been played, and is thus, a known card or
    // not.
    topPlayed := false
    for _, card := range append(played, prior...) {
        if top == card {
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
    if pickedUp && dealer != 0 && !topPlayed {
    // If the top was picked up by somebody else and it has not been played yet
    // we know where it is.
        nums[dealer - 1]--
    } else if (pickedUp && dealer == 0) || (!pickedUp) {
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
    if !pickedUp || (pickedUp && dealer != 0 && !topPlayed) {
        delete(unknowns, top)
    } else if pickedUp && dealer == 0 {
    // Similarly, we know one card is (the discarded card) if we picked it up.
        delete(unknowns, discard)
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
            if pickedUp && dealer != 0 && !topPlayed {
                switch dealer {
                case 1:
                    next.player1 = append(next.player1, top)
                case 2:
                    next.player2 = append(next.player2, top)
                case 3:
                    next.player3 = append(next.player3, top)
                }
            } else if pickedUp && dealer == 0 {
                next.kitty = append(next.kitty, discard)
            } else if !pickedUp {
                next.kitty = append(next.kitty, top)
            }

            c <- next
        }
        close(c)
    }()

    return c
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
func possible(hand, played []deck.Card, trump deck.Suit) []int {
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
