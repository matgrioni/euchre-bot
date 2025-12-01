package euchre

import (
	"euchre/ai"
	"euchre/deck"
	"fmt"
	"testing"
)

/*
 * Tests the ai package in conjunction with the euchre package. Some of these
 * tests are not real tests, but are simple sanity checks whose output must be
 * manually checked.
 */

/*
 * Tests the output for a run playout for the first card played by the computer.
 */
func TestRunPlayout(t *testing.T) {
	setup := Setup{
		1,
		1,
		true,
		deck.Card{Suit: deck.D, Value: deck.Nine},
		deck.D,
		deck.Card{},
		-1,
	}

	hand := []deck.Card{
		{Suit: deck.H, Value: deck.Nine},
		{Suit: deck.H, Value: deck.Ten},
		{Suit: deck.S, Value: deck.A},
		{Suit: deck.D, Value: deck.Q},
		{Suit: deck.C, Value: deck.Q},
	}

	played := []deck.Card{
		{Suit: deck.C, Value: deck.J},
		{Suit: deck.C, Value: deck.A},
	}

	var prior []Trick

	s := NewUndeterminizedState(setup, 0, hand, played, prior)
	s.Determinize()
	n := ai.NewNode()
	m := ai.Move{
		Action: nil,
		State:  s,
	}
	n.Value(m)
	e := Engine{}

	fmt.Println("Playout debug output")
	ai.RunPlayoutDebug(n, e)
}
