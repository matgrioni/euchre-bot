package euchre


import "deck"


/*
 * Create a random euchre situation. This means cards are randomly distributed
 * among the players, using our trusty player number assignment.
 *
 * Returns:
 *  A slice of card slices which corresponds to the player hands and the kitty
 *  in the last card slice.
 */
func GenSituation() [][]deck.Card {
    cards := deck.DrawN(24)

    hands := make([][]deck.Card, 5)
    for i := 0; i < 4; i++ {
        hands[i] = cards[i * 5: (i + 1) * 5]
    }

    hands[4] = cards[20:]

    return hands
}
