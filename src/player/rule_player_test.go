package player

import (
    "deck"
    "euchre"
    "testing"
)

// TODO: still add the tests for AlonePlayers. Much ado about that.


/*
 * Tests the expectations specific to the RulePlayer. This includes the normal
 * game play. Since it is rule based it can be exactly predicted what a rule
 * based player will do.
 */
type playTest struct {
    player int
    setup euchre.Setup
    hand []deck.Card
    played []deck.Card
    prior []euchre.Trick
    expected deck.Card
}

var playTests = []playTest {
    /*
     * Routine test with only one card played.
     */
    playTest {
        0,
        euchre.Setup {
            2,
            3,
            false,
            deck.Card { deck.C, deck.Ten },
            deck.H,
            deck.Card { },
            -1,
        },
        []deck.Card {
            deck.Card { deck.D, deck.J },
            deck.Card { deck.H, deck.J },
            deck.Card { deck.S, deck.Q },
            deck.Card { deck.D, deck.K },
            deck.Card { deck.C, deck.Nine },
        },
        []deck.Card {
            deck.Card { deck.C, deck.Q },
        },
        []euchre.Trick { },
        deck.Card { deck.C, deck.Nine },
    },

    /*
     * Test for playing first. The most valuable card in the hand should be
     * played.
     */
    playTest {
        0,
        euchre.Setup {
            3,
            0,
            true,
            deck.Card { deck.H, deck.Nine },
            deck.H,
            deck.Card { },
            -1,
        },
        []deck.Card {
            deck.Card { deck.H, deck.Ten },
            deck.Card { deck.H, deck.J },
            deck.Card { deck.S, deck.Ten },
            deck.Card { deck.S, deck.K },
            deck.Card { deck.C, deck.Ten },
        },
        []deck.Card { },
        []euchre.Trick { },
        deck.Card { deck.H, deck.J },
    },

    /*
     * Test for playing when partner is already winning. The lowest value card
     * should be thrown.
     */
    playTest {
        0,
        euchre.Setup {
            3,
            0,
            true,
            deck.Card { deck.H, deck.Nine },
            deck.H,
            deck.Card { },
            -1,
        },
        []deck.Card {
            deck.Card { deck.H, deck.Ten },
            deck.Card { deck.H, deck.J },
            deck.Card { deck.S, deck.Ten },
            deck.Card { deck.C, deck.Ten },
        },
        []deck.Card {
            deck.Card { deck.H, deck.Q },
            deck.Card { deck.H, deck.Nine },
        },
        []euchre.Trick {
            euchre.Trick {
                []deck.Card {
                    deck.Card { deck.S, deck.K },
                    deck.Card { deck.S, deck.J },
                    deck.Card { deck.S, deck.A },
                    deck.Card { deck.S, deck.Nine },
                },
                0,
                deck.H,
                -1,
            },
        },
        deck.Card { deck.H, deck.Ten },
    },

    /*
     * Test for playing when partner is already winning on non trump hand. The
     * lowest value card should be thrown, which is not trump if available.
     */
    playTest {
        0,
        euchre.Setup {
            3,
            0,
            true,
            deck.Card { deck.H, deck.Nine },
            deck.H,
            deck.Card { },
            -1,
        },
        []deck.Card {
            deck.Card { deck.H, deck.Ten },
            deck.Card { deck.H, deck.J },
            deck.Card { deck.S, deck.Ten },
            deck.Card { deck.C, deck.Nine },
        },
        []deck.Card {
            deck.Card { deck.C, deck.A },
            deck.Card { deck.C, deck.J },
        },
        []euchre.Trick {
            euchre.Trick {
                []deck.Card {
                    deck.Card { deck.S, deck.K },
                    deck.Card { deck.S, deck.J },
                    deck.Card { deck.S, deck.A },
                    deck.Card { deck.S, deck.Nine },
                },
                0,
                deck.H,
                -1,
            },
        },
        deck.Card { deck.C, deck.Nine },
    },

    /*
     * Test for playing lowest trump when that is all you have and partner is
     * winning.
     */
    playTest {
        0,
        euchre.Setup {
            3,
            0,
            true,
            deck.Card { deck.H, deck.Nine },
            deck.H,
            deck.Card { },
            -1,
        },
        []deck.Card {
            deck.Card { deck.H, deck.Ten },
            deck.Card { deck.H, deck.J },
            deck.Card { deck.H, deck.Q },
            deck.Card { deck.H, deck.K },
        },
        []deck.Card {
            deck.Card { deck.C, deck.A },
            deck.Card { deck.C, deck.J },
        },
        []euchre.Trick {
            euchre.Trick {
                []deck.Card {
                    deck.Card { deck.S, deck.K },
                    deck.Card { deck.S, deck.J },
                    deck.Card { deck.S, deck.A },
                    deck.Card { deck.S, deck.Nine },
                },
                0,
                deck.H,
                -1,
            },
        },
        deck.Card { deck.H, deck.Ten },
    },

    /*
     * Test for playing lowest card when there is no way to win.
     */
    playTest {
        0,
        euchre.Setup {
            0,
            0,
            false,
            deck.Card { deck.D, deck.Q },
            deck.S,
            deck.Card { },
            -1,
        },
        []deck.Card {
            deck.Card { deck.S, deck.Ten },
            deck.Card { deck.C, deck.J },
            deck.Card { deck.D, deck.Nine },
        },
        []deck.Card {
            deck.Card { deck.S, deck.J },
            deck.Card { deck.S, deck.K },
        },
        []euchre.Trick {
            euchre.Trick {
                []deck.Card {
                    deck.Card { deck.H, deck.Ten },
                    deck.Card { deck.H, deck.Q },
                    deck.Card { deck.H, deck.A },
                    deck.Card { deck.H, deck.K },
                },
                1,
                deck.S,
                -1,
            },
            euchre.Trick {
                []deck.Card {
                    deck.Card { deck.C, deck.K },
                    deck.Card { deck.C, deck.Q },
                    deck.Card { deck.C, deck.Nine },
                    deck.Card { deck.S, deck.A },
                },
                3,
                deck.S,
                -1,
            },
        },
        deck.Card { deck.S, deck.Ten },
    },

    /*
     * Test for playing the highest possible non trump possible if both can beat
     * the current winner.
     */
    playTest {
        0,
        euchre.Setup {
            1,
            3,
            false,
            deck.Card { deck.S, deck.A },
            deck.C,
            deck.Card { },
            -1,
        },
        []deck.Card {
            deck.Card { deck.H, deck.A },
            deck.Card { deck.S, deck.Q },
            deck.Card { deck.S, deck.A },
        },
        []deck.Card {
            deck.Card { deck.S, deck.Ten },
            deck.Card { deck.S, deck.Nine },
        },
        []euchre.Trick {
            euchre.Trick {
                []deck.Card {
                    deck.Card { deck.D, deck.K },
                    deck.Card { deck.C, deck.Q },
                    deck.Card { deck.D, deck.A },
                    deck.Card { deck.D, deck.Q },
                },
                2,
                deck.C,
                -1,
            },
            euchre.Trick {
                []deck.Card {
                    deck.Card { deck.S, deck.J },
                    deck.Card { deck.C, deck.Nine },
                    deck.Card { deck.C, deck.A },
                    deck.Card { deck.C, deck.J },
                },
                1,
                deck.C,
                -1,
            },
        },
        deck.Card { deck.S, deck.Q },
    },
}


/*
 * The main driver to test the rule players general playing logic.
 */
func TestPlay(t *testing.T) {
    player := NewRule("")

    for i, fixture := range playTests {
        _, chosen := player.Play(fixture.player, fixture.setup, fixture.hand,
                                 fixture.played, fixture.prior)
        if chosen != fixture.expected {
            t.Logf("Fixture %d failed.\n", i + 1)
            t.Errorf("Gave %s instead of %s", chosen, fixture.expected)
        }
    }
}
