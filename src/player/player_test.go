package player

import (
    "deck"
    "euchre"
    "testing"
)

// TODO: Format test cases, especially arrays, properly.

func TestMinimaxNoWinningChoice(t *testing.T) {
    setup := euchre.Setup {
        3,
        true,
        deck.Card {
            deck.S,
            deck.Q,
        },
        deck.S,
        deck.Card {
            deck.H,
            deck.Nine,
        },
    }

    hand := [...]deck.Card{ deck.Card{ deck.H, deck.Nine },
                            deck.Card{ deck.H, deck.Q },
                            deck.Card{ deck.D, deck.Q } }

    played := [...]deck.Card{ deck.Card{ deck.H, deck.A } }


    cards1 := [...]deck.Card{ deck.Card{ deck.S, deck.K },
                              deck.Card{ deck.S, deck.Ten },
                              deck.Card{ deck.C, deck.J },
                              deck.Card{ deck.S, deck.Nine } }
    cards2 := [...]deck.Card{ deck.Card{ deck.C, deck.A },
                              deck.Card{ deck.S, deck.Q },
                              deck.Card{ deck.C, deck.K },
                              deck.Card{ deck.C, deck.Ten } }

    prior := [...]euchre.Trick{
                    euchre.Trick {
                        cards1,
                        0,
                        deck.S,
                    },
                    euchre.Trick {
                        cards2,
                        2,
                        deck.S,
                    },
                }

    desired := hand[1]
    p := NewSmart()
    _, play := p.Play(setup, hand[:], played[:], prior[:])

    if play != desired {
        t.Errorf("Wanted %s but got %s", hand[1], play)
    }
}

func TestMinimaxMultipleWinOptions(t *testing.T) {
    setup := euchre.Setup {
        2,
        true,
        deck.Card {
            deck.S,
            deck.K,
        },
        deck.S,
        deck.Card {
            deck.H,
            deck.J,
        },
    }

    hand := [...]deck.Card{ deck.Card{ deck.S, deck.Nine },
                            deck.Card{ deck.S, deck.A } }

    played := [...]deck.Card{ deck.Card{ deck.H, deck.Q }, }


    cards1 := [...]deck.Card{ deck.Card{ deck.H, deck.A },
                              deck.Card{ deck.H, deck.J },
                              deck.Card{ deck.H, deck.Nine },
                              deck.Card{ deck.S, deck.K } }
    cards2 := [...]deck.Card{ deck.Card{ deck.C, deck.A },
                              deck.Card{ deck.C, deck.Ten },
                              deck.Card{ deck.D, deck.Nine },
                              deck.Card{ deck.C, deck.K } }
    cards3 := [...]deck.Card{ deck.Card{ deck.D, deck.A },
                              deck.Card{ deck.S, deck.Ten },
                              deck.Card{ deck.D, deck.Ten },
                              deck.Card{ deck.D, deck.J } }

    prior := [...]euchre.Trick{
                    euchre.Trick {
                        cards1,
                        3,
                        deck.S,
                    },
                    euchre.Trick {
                        cards2,
                        2,
                        deck.S,
                    },
                    euchre.Trick {
                        cards3,
                        2,
                        deck.S,
                    },
                }

    desired := hand[1]
    p := NewSmart()
    _, play := p.Play(setup, hand[:], played[:], prior[:])

    if play != desired {
        t.Errorf("Wanted %s but got %s", hand[1], play)
    }
}
