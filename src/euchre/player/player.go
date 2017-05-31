package player

import (
    "deck"
    "github.com/fighterlyt/permutation"
    "math/rand"
    "time"
    "fmt"
)

type situation struct {
    player1 []deck.Card
    player2 []deck.Card
    player3 []deck.Card
    kitty   []deck.Card
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
// top      -
// trump    -
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
        for situation := range situations(hand, played, prior, discard, dealer,
                                          pickedUp, top) {
            fmt.Println(situation.player1[0])
            break
        }
    }

    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    playable := possible(hand, played, trump)
    chosen := playable[r.Intn(len(playable))]
    final := hand[chosen]
    hand = append(hand[:chosen], hand[chosen + 1:]...)

    return final, hand
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
    // Create a set-like structure that has all the cards currently in play. Do
    // this by adding all cards and then removing those that are in your hand
    // have been played, and the card that was picked up if any.
    inPlay := make(map[deck.Card]bool)
    for _, card := range deck.CARDS {
        inPlay[card] = true
    }

    topPlayed := false
    for _, card := range append(played, prior...) {
        if top == card {
            topPlayed = true
            break
        }
    }

    // Remove all cards that we have already seen in some way.
    unavailable := append(hand, append(played, prior...)...)
    if !pickedUp || ((pickedUp && dealer != 0) && !topPlayed) {
        unavailable = append(unavailable, top)
    } else if pickedUp && dealer == 0 {
        unavailable = append(unavailable, discard)
    }

    for _, card := range unavailable {
        delete(inPlay, card)
    }

    // Set available to the keys of the inPlay map. These are the cards whose
    // distribution is unknown.
    available := make([]deck.Card, len(inPlay), len(inPlay))
    i := 0
    for card := range inPlay {
        available[i] = card
        i++
    }

    c := make(chan situation)
    gen, _ := permutation.NewPerm(available, permutation.Less(cardLess))

    // Figure out how many unknown cards are in each player's hand and in the
    // kitty.
    nums := [3]int{len(hand), len(hand), len(hand)}
    kit := 4
    for i := 2; i > 2 - len(played); i-- {
        nums[i]--
    }

    if pickedUp && dealer != 0 && !topPlayed {
        nums[dealer - 1]--
    } else if (pickedUp && dealer == 0) || (!pickedUp) {
        kit--
    }

    if nums[0] + nums[1] + nums[2] + kit != len(available) {
        panic("Number of freely fluxing cards is not correct.")
    }

    go func() {
        for t, err := gen.Next(); err == nil; t, err = gen.Next() {
            perm := t.([]deck.Card)
            next := situation {
                perm[:nums[0]],
                perm[nums[0]:nums[0] + nums[1]],
                perm[nums[0] + nums[1]:nums[0] + nums[1] + nums[2]],
                perm[nums[0] + nums[1] + nums[2]:],
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

func cardLess(value1, value2 interface{}) bool {
    return true
}
