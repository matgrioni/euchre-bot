package euchre

import (
    "deck"
    "testing"
)


/*
 * A test struct for going alone. This struct constitutes one test case for
 * the WinnerAlone method.
 */
type winnerAloneTest struct {
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
 * euchre#Winner. A helper to determine which player is the winner given the
 * player who led, the list of cards, and the trump suit.
 */


func TestWinnerOneTrump1(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.D, deck.A },
        deck.Card{ deck.H, deck.Q },
        deck.Card{ deck.D, deck.Ten },
        deck.Card{ deck.D, deck.Q },
    }

    trump := deck.H
    led := 2
    answer := 3

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}


func TestWinnerOneTrump2(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.H, deck.Nine },
        deck.Card{ deck.S, deck.A },
        deck.Card{ deck.D, deck.Q },
    }

    trump := deck.D
    led := 2
    answer := 1

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}


func TestWinnerOneTrump3(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.D, deck.Q },
        deck.Card{ deck.D, deck.K },
        deck.Card{ deck.C, deck.Ten },
        deck.Card{ deck.D, deck.A },
    }

    led := 1
    trump := deck.C
    answer := 3

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}


func TestWinnerOneTrump4(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.H, deck.J },
        deck.Card{ deck.H, deck.Q },
        deck.Card{ deck.C, deck.Ten },
    }

    led := 2
    trump := deck.C
    answer := 1

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}


func TestWinnerOneTrump5(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.H, deck.J },
        deck.Card{ deck.H, deck.Q },
        deck.Card{ deck.C, deck.Ten },
    }

    led := 2
    trump := deck.C
    answer := 1

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}


func TestWinnerNonTrump(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.H, deck.Ten },
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.H, deck.Q },
        deck.Card{ deck.H, deck.Ten },
    }

    led := 3
    trump := deck.C
    answer := 0

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}


func TestWinnerManyTrump(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.D, deck.J },
        deck.Card{ deck.C, deck.Ten },
        deck.Card{ deck.S, deck.J },
        deck.Card{ deck.C, deck.A },
    }

    led := 1
    trump := deck.C
    answer := 3

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}


/*
 * Test WinnerAlone.
 */

var winnerAloneTests = []winnerAloneTest {
    /*
     * A test where the discarded player doesn't matter since the winner is
     * before him.
     */
    winnerAloneTest {
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
    winnerAloneTest {
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
    winnerAloneTest {
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
     * Just another normal test case.
     */
    winnerAloneTest {
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
 * Test we are able to properly determine the winner when there is a player
 * who is going alone.
 */
func TestWinnerAlone(t *testing.T) {
    for _, test := range winnerAloneTests {
        res := WinnerAlone(test.played, test.trump, test.led, test.alone)

        if res != test.expected {
            t.Errorf("Expected error to be %d but got %d\n", test.expected, res)
        }
    }
}
