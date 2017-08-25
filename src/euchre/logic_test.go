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


type noSuitsTest struct {
    prior []Trick
    trump deck.Suit
    expected map[int][]deck.Suit
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
 * Tests all the above tests to make sure the expected winner is actually seen.
 */
func TestWinnerAlone(t *testing.T) {
    for _, test := range winnerTests {
        res := Winner(test.played, test.trump, test.led, test.alone)

        if res != test.expected {
            t.Errorf("Expected error to be %d but got %d\n", test.expected, res)
        }
    }
}


/*
 * Test noSuits.
 *
 * Checks whether we correctly learn what suits a player has based on their
 * inability to follow.
 */

var noSuitsTests = []noSuitsTest {
    // No tricks played yet returns no results.
    noSuitsTest {
        nil,
        deck.H,
        map[int][]deck.Suit { },
    },

    // Over two tricks, only one player does not follow suit.
    noSuitsTest {
        []Trick {
            Trick {
                []deck.Card {
                    deck.Card { deck.H, deck.A },
                    deck.Card { deck.C, deck.K },
                    deck.Card { deck.H, deck.K },
                    deck.Card { deck.H, deck.Q },
                },
                0,
                deck.C,
                -1,
            },
            Trick {
                []deck.Card {
                    deck.Card { deck.S, deck.K },
                    deck.Card { deck.S, deck.Ten },
                    deck.Card { deck.S, deck.Q },
                    deck.Card { deck.S, deck.A },
                },
                1,
                deck.C,
                -1,
            },
        },
        deck.C,
        map[int][]deck.Suit {
            1: []deck.Suit {
                deck.H,
            },
        },
    },

    // Player 3 (last in modulo order) and 3 players in total are missing suits.
    noSuitsTest {
        []Trick {
            Trick {
                []deck.Card {
                    deck.Card { deck.H, deck.A },
                    deck.Card { deck.C, deck.K },
                    deck.Card { deck.H, deck.K },
                    deck.Card { deck.H, deck.Q },
                },
                0,
                deck.C,
                -1,
            },
            Trick {
                []deck.Card {
                    deck.Card { deck.S, deck.K },
                    deck.Card { deck.S, deck.Ten },
                    deck.Card { deck.C, deck.Ten },
                    deck.Card { deck.C, deck.A },
                },
                1,
                deck.C,
                -1,
            },
        },
        deck.C,
        map[int][]deck.Suit {
            1: []deck.Suit {
                    deck.H,
            },
            3: []deck.Suit {
                    deck.S,
            },
            0: []deck.Suit {
                    deck.S,
            },
        },
    },

    // Check that modulo player numbers work in keeping track.
    noSuitsTest {
        []Trick {
            Trick {
                []deck.Card {
                    deck.Card { deck.H, deck.A },
                    deck.Card { deck.C, deck.K },
                    deck.Card { deck.H, deck.K },
                    deck.Card { deck.H, deck.Q },
                },
                3,
                deck.C,
                -1,
            },
            Trick {
                []deck.Card {
                    deck.Card { deck.S, deck.K },
                    deck.Card { deck.S, deck.Ten },
                    deck.Card { deck.C, deck.Ten },
                    deck.Card { deck.C, deck.A },
                },
                0,
                deck.C,
                -1,
            },
        },
        deck.C,
        map[int][]deck.Suit {
            0: []deck.Suit {
                deck.H,
            },
            2: []deck.Suit {
                deck.S,
            },
            3: []deck.Suit {
                deck.S,
            },
        },
    },

    // One player is missing more than one suit.
    noSuitsTest {
        []Trick {
            Trick {
                []deck.Card {
                    deck.Card { deck.D, deck.A },
                    deck.Card { deck.S, deck.J },
                    deck.Card { deck.D, deck.Ten },
                    deck.Card { deck.D, deck.J },
                },
                0,
                deck.S,
                -1,
            },
            Trick {
                []deck.Card {
                    deck.Card { deck.C, deck.A },
                    deck.Card { deck.C, deck.Q },
                    deck.Card { deck.C, deck.K },
                    deck.Card { deck.S, deck.Ten },
                },
                1,
                deck.S,
                -1,
            },
            Trick {
                []deck.Card {
                    deck.Card { deck.H, deck.K },
                    deck.Card { deck.H, deck.Ten },
                    deck.Card { deck.H, deck.Nine },
                    deck.Card { deck.H, deck.A },
                },
                0,
                deck.S,
                -1,
            },
            Trick {
                []deck.Card {
                    deck.Card { deck.H, deck.Q },
                    deck.Card { deck.S, deck.A },
                    deck.Card { deck.C, deck.Ten },
                    deck.Card { deck.S, deck.K },
                },
                3,
                deck.S,
                -1,
            },
        },
        deck.S,
        map[int][]deck.Suit {
            0: []deck.Suit {
                deck.C,
                deck.H,
            },
            1: []deck.Suit {
                deck.D,
                deck.H,
            },
            2: []deck.Suit {
                deck.H,
            },
        },
    },

    // Going alone v1
    noSuitsTest {
        []Trick {
            Trick {
                []deck.Card {
                    deck.Card { deck.C, deck.A },
                    deck.Card { deck.C, deck.K },
                    deck.Card { deck.H, deck.Nine },
                },
                0,
                deck.H,
                3,
            },
        },
        deck.H,
        map[int][]deck.Suit {
            3: []deck.Suit {
                deck.C,
            },
        },
    },

    // Going alone v2, where we must pass the modulo wrap
    noSuitsTest {
        []Trick {
            Trick {
                []deck.Card {
                    deck.Card { deck.S, deck.J },
                    deck.Card { deck.C, deck.J },
                    deck.Card { deck.H, deck.Nine },
                },
                3,
                deck.S,
                2,
            },
        },
        deck.S,
        map[int][]deck.Suit {
            2: []deck.Suit {
                deck.S,
            },
        },
    },

    // Going alone v3, where two players do not follow suit across two tricks.
    noSuitsTest {
        []Trick {
            Trick {
                []deck.Card {
                    deck.Card { deck.H, deck.K },
                    deck.Card { deck.S, deck.Nine },
                    deck.Card { deck.H, deck.A },
                },
                1,
                deck.S,
                2,
            },
            Trick {
                []deck.Card {
                    deck.Card { deck.S, deck.K },
                    deck.Card { deck.S, deck.Nine },
                    deck.Card { deck.D, deck.Q },
                },
                2,
                deck.S,
                2,
            },
        },
        deck.S,
        map[int][]deck.Suit {
            2: []deck.Suit {
                deck.H,
            },
            1: []deck.Suit {
                deck.S,
            },
        },
    },

    noSuitsTest {
        []Trick {
            Trick {
                []deck.Card {
                    deck.Card { deck.C, deck.A },
                    deck.Card { deck.C, deck.Q },
                    deck.Card { deck.H, deck.Nine, },
                    deck.Card { deck.C, deck.Ten, },
                },
                0,
                deck.C,
                -1,
            },
        },
        deck.C,
        map[int][]deck.Suit {
            2: []deck.Suit {
                deck.C,
            },
        },
    },

    noSuitsTest {
        []Trick {
            Trick {
                []deck.Card {
                    deck.Card { deck.D, deck.J },
                    deck.Card { deck.H, deck.J },
                    deck.Card { deck.S, deck.A },
                    deck.Card { deck.C, deck.Nine },
                },
                2,
                deck.D,
                -1,
            },
        },
        deck.D,
        map[int][]deck.Suit {
            0: []deck.Suit {
                deck.D,
            },
            1: []deck.Suit {
                deck.D,
            },
        },
    },

    noSuitsTest {
        []Trick {
            Trick {
                []deck.Card {
                    deck.Card { deck.H, deck.Ten },
                    deck.Card { deck.D, deck.K },
                    deck.Card { deck.H, deck.Nine },
                    deck.Card { deck.D, deck.A },
                },
                3,
                deck.C,
                -1,
            },
            Trick {
                []deck.Card {
                    deck.Card { deck.S, deck.K },
                    deck.Card { deck.S, deck.A },
                    deck.Card { deck.S, deck.Nine },
                    deck.Card { deck.C, deck.A },
                },
                3,
                deck.C,
                -1,
            },
        },
        deck.C,
        map[int][]deck.Suit {
            0: []deck.Suit {
                deck.H,
            },
            2: []deck.Suit {
                deck.H,
                deck.S,
            },
        },
    },

}


/*
 * Test the noSuits method. The expected result matches the actual result if the
 * maps contain the same content but not necessairily in the same order.
 */
func TestNoSuits(t *testing.T) {
    for _, test := range noSuitsTests {
        res := noSuits(test.prior, test.trump)

        // Check that the two maps are the same size.
        if len(res) != len(test.expected) {
            errorOut(t, test.expected, res, "no suits")
        }

        // For each player check that the suits are the same.
        for player, suits := range test.expected {
            // Check if this player does exists in the actual results.
            actualTmp, ok := res[player]
            if !ok {
                errorOut(t, test.expected, res, "no suits")
            }
            actual := suitsSliceToSet(actualTmp)
            expected := suitsSliceToSet(suits)

            // Check that the two sets of suits are the same length, and then
            // have the same contents.
            if len(actual) != len(expected) {
                errorOut(t, test.expected, res, "no suits")
            }

            for suit, _ := range expected {
                _, ok := actual[suit]
                if !ok {
                    errorOut(t, test.expected, res, "no suits")
                }
            }
        }
    }
}


/*
 * Converts a list of suits into a set of those suits.
 *
 * Args:
 *  suits: The list of suits.
 *
 * Returns:
 *  The set of suits. This is returned in the form of a map, where the key is
 *  the suit and the value is true. If a suit was not in the list it is not in
 *  the map.
 */
func suitsSliceToSet(suits []deck.Suit) map[deck.Suit]bool {
    set := make(map[deck.Suit]bool)

    for _, suit := range suits {
        set[suit] = true
    }

    return set
}


/*
 * Errors out the given test.
 *
 * Args:
 *  t: The current testing context.
 *  expected: The expected value from the test case.
 *  actual: The value actually calculated.
 *  test: The test name to output for clarification purposes.
 */
func errorOut(t *testing.T, expected interface{}, actual interface{},
              test string) {
    t.Errorf("Expected %v but got %v for %s\n", expected, actual, test)
}
