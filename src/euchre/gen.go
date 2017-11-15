package euchre


import (
    "deck"
)

/*
 * Create a random euchre situation. This means cards are randomly distributed
 * among the players, using our trusty player number assignment.
 *
 * Returns:
 *  A slice of card slices which corresponds to the player hands and the kitty
 *  in the last card slice.
 */
func GenSituation() [][]deck.Card {
    situation := r.Perm(len(deck.CARDS))
    cards := indicesToCards(situation)

    hands := make([][]deck.Card, 5)
    for i := 0; i < 4; i++ {
        hands[i] = cards[i * 5: (i + 1) * 5]
    }

    hands[4] = cards[20:]

    return hands
}


/*
 * Transform a list of indices into the corresponding cards.
 *
 * Args:
 *  indices, type(int): A list of indices that correspond to cards in deck.CARDS
 *
 * Returns:
 *  type([]deck.Card): A slice that is the mapping from the indices to cards.
 */
func indicesToCards(indices []int) []deck.Card {
    cards := make([]deck.Card, len(indices))
    for i, idx := range indices {
        cards[i] = deck.CARDS[idx]
    }

    return cards
}
