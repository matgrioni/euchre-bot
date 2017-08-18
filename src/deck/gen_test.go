package deck

import "testing"


/*
 * Test the generate hand functionality.
 */


/*
 * Test that all cards in the generated hand are unique.
 */
func TestGenHandUnique(t *testing.T) {
    n := 5
    hand := GenHand(n)
    present := make(map[Card]bool)

    for _, card := range hand {
        if _, ok := present[card]; ok {
            t.Errorf("%s has been seen twice.\n", card)
        }
    }
}


/*
 * Test that the generated hand is of the expected length.
 */
func TestGenHandLength(t *testing.T) {
    ns := []int { 1, 2, 3, 4, 5 }

    for _, n := range ns {
        hand := GenHand(n)

        if len(hand) != n {
            t.Errorf("Expected hand length to be %s but it is %s.\n", n, len(hand))
        }
    }
}
