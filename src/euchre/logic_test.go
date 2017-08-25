package euchre

import (
    "deck"
    "fmt"
    "testing"
)


type beatTest struct {
    a deck.Card
    b deck.Card
    trump deck.Suit
    expected bool
}


type noSuitsTest struct {
    prior []Trick
    trump deck.Suit
    expected map[int][]deck.Suit
}


type possibleTest struct {
    hand []deck.Card
    played []deck.Card
    trump deck.Suit
    expected []int
}


type winnerTest struct {
    played []deck.Card
    trump deck.Suit
    led int
    alone int
    expected int
}


/*
 * Test Beat. A helper to determine which card wins in a head to head faceoff
 * where one card is chosen to be leading and there is a defined trump suit.
 */

var beatTests = []beatTest {
    // Trump beats non-trump
    beatTest {
        deck.Card { deck.H, deck.Ten },
        deck.Card { deck.C, deck.A },
        deck.H,
        true,
    },

    // Trump beats non-trump v2, with orders flipped
    beatTest {
        deck.Card { deck.S, deck.Q },
        deck.Card { deck.C, deck.A },
        deck.C,
        false,
    },

    // If suits don't match and no trumps, then leader wins
    beatTest {
        deck.Card { deck.S, deck.Nine },
        deck.Card { deck.C, deck.A },
        deck.D,
        true,
    },

    // The higher card of the same suit wins
    beatTest {
        deck.Card { deck.S, deck.A },
        deck.Card { deck.S, deck.Nine },
        deck.D,
        true,
    },

    // The right bower beats the left bower
    beatTest {
        deck.Card { deck.H, deck.J },
        deck.Card { deck.D, deck.J },
        deck.D,
        false,
    },
}


/*
 * Run all the beat tests defined above and output any errors in reference to
 * the expected and actual values and the index of the test.
 */
func TestBeat(t *testing.T) {
    for i, test := range beatTests {
        res := Beat(test.a, test.b, test.trump)

        if res != test.expected {
            errorOut(t, test.expected, res, "beat", i)
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
    for i, test := range noSuitsTests {
        res := noSuits(test.prior, test.trump)

        // Check that the two maps are the same size.
        if len(res) != len(test.expected) {
            errorOut(t, test.expected, res, "no suits", i)
        }

        // For each player check that the suits are the same.
        for player, suits := range test.expected {
            // Check if this player does exists in the actual results.
            actualTmp, ok := res[player]
            if !ok {
                errorOut(t, test.expected, res, "no suits", i)
            }
            actual := suitsSliceToSet(actualTmp)
            expected := suitsSliceToSet(suits)

            // Check that the two sets of suits are the same length, and then
            // have the same contents.
            if len(actual) != len(expected) {
                errorOut(t, test.expected, res, "no suits", i)
            }

            for suit, _ := range expected {
                _, ok := actual[suit]
                if !ok {
                    errorOut(t, test.expected, res, "no suits", i)
                }
            }
        }
    }
}


/*
 * Test Possible. This returns the cards that can be played given a hand and the
 * cards currently played. The actual cards are not returned, the indices of the
 * possible cards are provided.
 */

var possibleTests = []possibleTest {
    // Can't follow lead card, so all cards are possible.
    possibleTest {
        []deck.Card {
            deck.Card { deck.H, deck.Ten },
            deck.Card { deck.H, deck.A },
            deck.Card { deck.C, deck.J },
            deck.Card { deck.C, deck.K },
            deck.Card { deck.D, deck.J },
        },
        []deck.Card {
            deck.Card { deck.S, deck.Ten },
        },
        deck.H,
        []int {
            0, 1, 2, 3, 4,
        },
    },

    // Can follow with one card.
    possibleTest {
        []deck.Card {
            deck.Card { deck.H, deck.Ten },
            deck.Card { deck.H, deck.A },
            deck.Card { deck.C, deck.J },
            deck.Card { deck.C, deck.K },
            deck.Card { deck.S, deck.J },
        },
        []deck.Card {
            deck.Card { deck.S, deck.Ten },
        },
        deck.H,
        []int {
            4,
        },
    },
}


func TestPossible(t *testing.T) {
    for i, test := range possibleTests {
        res := Possible(test.hand, test.played, test.trump)

        if len(res) != len(test.expected) {
            errorOut(t, test.expected, res, "possible", i)
        }

        for j := 0; j < len(test.expected); j++ {
            if test.expected[j] != res[j] {
                errorOut(t, test.expected, res, "possible", i)
            }
        }
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
 * Run all the above tests to make sure the expected winner is actually seen.
 */
func TestWinner(t *testing.T) {
    for _, test := range winnerTests {
        res := Winner(test.played, test.trump, test.led, test.alone)

        if res != test.expected {
            t.Errorf("Expected error to be %d but got %d\n", test.expected, res)
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
 *  index: The index of the test so that the exact test is known.
 */
func errorOut(t *testing.T, expected interface{}, actual interface{},
              test string, index int) {
    testName := fmt.Sprintf("%s[%d]", test, index)
    t.Errorf("Expected %v but got %v for %s\n", expected, actual, testName)
}
