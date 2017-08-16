package player

import (
    "deck"
    "testing"
)


/*
 * Tests the different types of players functionality.
 */


/*
 * A test fixture that defines the inputs and the expected output for a given
 * test case.
 */
type fixture struct {
    hand []deck.Card
    top deck.Card
    expected deck.Card
}


var fixtures = []fixture {
    /*
     * Test that the lowest trump card is discarded if all cards are trump.
     */
    fixture {
        []deck.Card {
            deck.Card { deck.H, deck.J  },
            deck.Card { deck.H, deck.Ten },
            deck.Card { deck.H, deck.Nine },
            deck.Card { deck.H, deck.A },
            deck.Card { deck.D, deck.J },
        },
        deck.Card { deck.H, deck.Q },
        deck.Card { deck.H, deck.Nine },
    },

    /*
     * Whitebox testing where trumps are in ascending order. The test is to see
     * if the lowest of a trump suit will be chosen independent of order.
     */
    fixture {
        []deck.Card {
            deck.Card { deck.C, deck.Nine },
            deck.Card { deck.C, deck.Ten },
            deck.Card { deck.C, deck.Q },
            deck.Card { deck.C, deck.K },
            deck.Card { deck.C, deck.J },
        },
        deck.Card { deck.C, deck.A },
        deck.Card { deck.C, deck.Nine },
    },

    /*
     * Another whitebox test to assure that discarding properly handles bowers.
     */
    fixture {
        []deck.Card {
            deck.Card { deck.C, deck.J },
            deck.Card { deck.S, deck.J },
            deck.Card { deck.C, deck.Q },
            deck.Card { deck.C, deck.K },
            deck.Card { deck.C, deck.A },
        },
        deck.Card { deck.C, deck.Ten },
        deck.Card { deck.C, deck.Ten },
    },

    /*
     * When there is a suit with only one card in it, but it is an A, do not discard
     * it since it is valuable. Discard the lowest card of a non-trump suit.
     */
    fixture {
        []deck.Card {
            deck.Card { deck.C, deck.A },
            deck.Card { deck.H, deck.Nine },
            deck.Card { deck.H, deck.J },
            deck.Card { deck.S, deck.K },
            deck.Card { deck.S, deck.Q },
        },
        deck.Card { deck.H, deck.Q },
        deck.Card { deck.S, deck.Q },
    },

    /*
     * Test that a suit with only one card is discarded. This makes sense if you
     * have trumps. TODO: Otherwise, maybe something else should happen.
     */
    fixture {
        []deck.Card {
            deck.Card { deck.C, deck.Q },
            deck.Card { deck.H, deck.Nine },
            deck.Card { deck.H, deck.J },
            deck.Card { deck.S, deck.K },
            deck.Card { deck.S, deck.Q },
        },
        deck.Card { deck.H, deck.Q },
        deck.Card { deck.C, deck.Q },
    },
}

// TODO: Improve test modularity. Logic can be condensed together.

func TestDiscard(t *testing.T) {
    players := getTestablePlayers()

    for i, fixture := range fixtures {
        for j, player := range players {
            copyHand := make([]deck.Card, len(fixture.hand))
            copy(copyHand, fixture.hand)

            _, discarded := player.Discard(copyHand, fixture.top)
            if discarded != fixture.expected {
                t.Logf("Fixture %d, implementation %d failed.\n", i + 1, j + 1)
                t.Errorf("Gave %s instead of %s.\n", discarded, fixture.expected)
            }
        }
    }
}


/*
 * Returns a list of all the different Player implementations to test.
 *
 * Returns:
 *  A list of the different player implementations to test in this file. The
 *  order of the players is [rule, smart].
 */
func getTestablePlayers() []Player {
    rule := NewRule("")
    smart := NewSmart()

    players := []Player { rule, smart }

    return players
}
