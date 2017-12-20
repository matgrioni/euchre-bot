package deck

import "testing"


/*
 * Test the card drawing functionality.
 * TODO: Add test for more cards than in deck, and test for 0 cards.
 */


/*
 * Test that all cards in the randomly drawn cards are unique.
 */
func TestDrawNUnique(t *testing.T) {
    n := 5
    hand := DrawN(n)
    present := make(map[Card]bool)

    for _, card := range hand {
        if _, ok := present[card]; ok {
            t.Errorf("%s has been seen twice.\n", card)
        }
    }
}


/*
 * Test that the randomly drawn cards are of the expected length.
 */
func TestDrawNLength(t *testing.T) {
    ns := []int { 1, 2, 3, 4, 5 }

    for _, n := range ns {
        hand := DrawN(n)

        if len(hand) != n {
            t.Errorf("Expected hand length to be %s but it is %s.\n", n, len(hand))
        }
    }
}
