package euchre

import (
    "deck"
    "testing"
)


/*
 * A test struct for going alone. This struct constitutes one test case for
 * the WinnerAlone method.
 */
type winnerTest struct {
    played []deck.Card
    trump deck.Suit
    led int
    alone int
    expected int
}


/*
 * Test euchre#Beat. A helper to determine which card wins in a head to head
 * faceoff where one card is chosen to be leading and there is a defined
 * trump suit.
 */


func TestBeatTrumpNonTrump(t *testing.T) {
    trump := deck.H
    card1 := deck.Card{ deck.H, deck.Ten }
    card2 := deck.Card{ deck.C, deck.A }

    if !Beat(card1, card2, trump) {
        t.Errorf("Expected %s to beat %s with trump %s", card1, card2, trump)
    }
}


func TestBeatNonTrumpTrump(t *testing.T) {
    trump := deck.C
    card1 := deck.Card{ deck.S, deck.Q }
    card2 := deck.Card{ deck.C, deck.A }

    if Beat(card1, card2, trump) {
        t.Errorf("Expected %s to beat %s with trump %s", card2, card1, trump)
    }
}


func TestBeatDifferentNonTrumpFirstLess(t *testing.T) {
    trump := deck.D
    card1 := deck.Card{ deck.S, deck.Nine }
    card2 := deck.Card{ deck.C, deck.A }

    if !Beat(card1, card2, trump) {
        t.Errorf("Expected %s to beat %s with trump %s", card1, card2, trump)
    }
}


func TestBeatDifferentNonTrumpFirstGreater(t *testing.T) {
    trump := deck.D
    card1 := deck.Card{ deck.S, deck.A }
    card2 := deck.Card{ deck.C, deck.Nine }

    if !Beat(card1, card2, trump) {
        t.Errorf("Expected %s to beat %s with trump %s", card1, card2, trump)
    }
}


func TestBeatSameSuit(t *testing.T) {
    trump := deck.D
    card1 := deck.Card{ deck.S, deck.A }
    card2 := deck.Card{ deck.S, deck.Nine }

    if !Beat(card1, card2, trump) {
        t.Errorf("Expected %s to beat %s with trump %s", card1, card2, trump)
    }
}


func TestBeatBowers(t *testing.T) {
    trump := deck.D
    card1 := deck.Card{ deck.H, deck.J }
    card2 := deck.Card{ deck.D, deck.J }

    if Beat(card1, card2, trump) {
        t.Errorf("Expected %s to beat %s with trump %s", card2, card1, trump)
    }
}


/*
 * Test Winner
 */

var winnerTests = []winnerTest {
    /*
     * The start of the non-going alone tests (i.e. where everybody plays).
     */

    // Only one trump
    winnerTest {
        []deck.Card {
            deck.Card { deck.D, deck.A },
            deck.Card { deck.H, deck.Q },
            deck.Card { deck.D, deck.Ten },
            deck.Card { deck.D, deck.Q },
        },
        deck.H,
        2,
        -1,
        3,
    },

    // Only one trump v2
    winnerTest {
        []deck.Card {
            deck.Card { deck.H, deck.A },
            deck.Card { deck.H, deck.Nine },
            deck.Card { deck.S, deck.A },
            deck.Card { deck.D, deck.Q },
        },
        deck.D,
        2,
        -1,
        1,
    },

    // Only one trump v3
    winnerTest {
        []deck.Card {
            deck.Card { deck.D, deck.Q },
            deck.Card { deck.D, deck.K },
            deck.Card { deck.C, deck.Ten },
            deck.Card { deck.D, deck.A },
        },
        deck.C,
        1,
        -1,
        3,
    },

    // Only one trump v4
    winnerTest {
        []deck.Card {
            deck.Card { deck.H, deck.A },
            deck.Card { deck.H, deck.J },
            deck.Card { deck.H, deck.Q },
            deck.Card { deck.C, deck.Ten },
        },
        deck.C,
        2,
        -1,
        1,
    },

    // Only one trump v5
    winnerTest {
        []deck.Card {
            deck.Card { deck.H, deck.A },
            deck.Card { deck.H, deck.J },
            deck.Card { deck.H, deck.Q },
            deck.Card { deck.C, deck.Ten },
        },
        deck.C,
        2,
        -1,
        1,
    },

    // All non-trump
    winnerTest {
        []deck.Card {
            deck.Card { deck.H, deck.Ten },
            deck.Card { deck.H, deck.A },
            deck.Card { deck.H, deck.Q },
            deck.Card { deck.H, deck.Ten },
        },
        deck.C,
        3,
        -1,
        0,
    },

    // Multiple trump in one trick.
    winnerTest {
        []deck.Card {
            deck.Card { deck.D, deck.J },
            deck.Card { deck.C, deck.Ten },
            deck.Card { deck.S, deck.J },
            deck.Card { deck.C, deck.A },
        },
        deck.C,
        1,
        -1,
        3,
    },

    /*
     * The start of going alone tests.
     */

    /*
     * A test where the discarded player doesn't matter since the winner is
     * before him.
     */
    winnerTest {
        []deck.Card {
            deck.Card { deck.H, deck.A },
            deck.Card { deck.C, deck.Ten },
            deck.Card { deck.C, deck.Q },
        },
        deck.H,
        0,
        0,
        0,
    },

    /*
     * A test where the winner is after the person who is cucked.
     */
    winnerTest {
        []deck.Card {
            deck.Card { deck.D, deck.J },
            deck.Card { deck.H, deck.J },
            deck.Card { deck.S, deck.Q },
        },
        deck.H,
        1,
        0,
        3,
    },

    /*
     * A test where the leader is the one going alone.
     */
    winnerTest {
        []deck.Card {
            deck.Card { deck.H, deck.Q },
            deck.Card { deck.S, deck.A },
            deck.Card { deck.H, deck.A },
        },
        deck.H,
        1,
        1,
        0,
    },

    /*
     * Just another normal test case for the winner when going alone.
     */
    winnerTest {
        []deck.Card {
            deck.Card { deck.C, deck.Nine },
            deck.Card { deck.H, deck.Q },
            deck.Card { deck.S, deck.J },
        },
        deck.S,
        2,
        1,
        1,
    },
}


/*
 *
 */
func TestWinnerAlone(t *testing.T) {
    for _, test := range winnerTests {
        res := Winner(test.played, test.trump, test.led, test.alone)

        if res != test.expected {
            t.Errorf("Expected error to be %d but got %d\n", test.expected, res)
        }
    }
}
