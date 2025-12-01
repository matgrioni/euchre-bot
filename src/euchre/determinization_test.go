package euchre

import (
	"euchre/deck"
	"fmt"
	"testing"
)

/*
 * Test the determinization of the euchre State when there have been no cards
 * played in the current trick.
 */
func TestDeterminization(t *testing.T) {
	setup := Setup{
		1,
		2,
		false,
		deck.Card{Suit: deck.C, Value: deck.A},
		deck.S,
		deck.Card{},
		-1,
	}

	hand := []deck.Card{
		deck.Card{Suit: deck.H, Value: deck.A},
		deck.Card{Suit: deck.D, Value: deck.Nine},
		deck.Card{Suit: deck.D, Value: deck.A},
		deck.Card{Suit: deck.C, Value: deck.Ten},
	}

	cards1 := []deck.Card{
		deck.Card{Suit: deck.H, Value: deck.K},
		deck.Card{Suit: deck.C, Value: deck.Q},
		deck.Card{Suit: deck.H, Value: deck.Q},
		deck.Card{Suit: deck.H, Value: deck.Nine},
	}

	var played []deck.Card

	trick1 := Trick{
		cards1,
		2,
		deck.S,
		-1,
	}

	prior := []Trick{trick1}

	state := NewUndeterminizedState(setup, 2, hand, played, prior)
	state.Determinize()

	fmt.Printf("Sanity check: %v\n", state)
}

/*
 * Test that a determinization in the middle of a trick will result in the
 * proper amount of cards in each hand.
 */

func TestDeterminizationMiddleOfTrick(t *testing.T) {
	setup := Setup{
		1,
		2,
		false,
		deck.Card{Suit: deck.C, Value: deck.A},
		deck.S,
		deck.Card{},
		-1,
	}

	hand := []deck.Card{
		{Suit: deck.H, Value: deck.A},
		{Suit: deck.D, Value: deck.Nine},
		{Suit: deck.D, Value: deck.A},
		{Suit: deck.C, Value: deck.Ten},
	}

	cards1 := []deck.Card{
		{Suit: deck.H, Value: deck.K},
		{Suit: deck.C, Value: deck.Q},
		{Suit: deck.H, Value: deck.Q},
		{Suit: deck.H, Value: deck.Nine},
	}

	played := []deck.Card{
		{Suit: deck.C, Value: deck.A},
		{Suit: deck.C, Value: deck.Ten},
	}

	trick1 := Trick{
		cards1,
		2,
		deck.S,
		-1,
	}

	prior := []Trick{trick1}

	state := NewUndeterminizedState(setup, 0, hand, played, prior)
	state.Determinize()

	if len(state.Hands[0]) != 4 || len(state.Hands[1]) != 4 ||
		len(state.Hands[2]) != 3 || len(state.Hands[3]) != 3 {
		t.Errorf("Incorrect number of cards in some players hand.\n")
	}
}

/*
 * Test a problematic instance that is giving problems in generating next valid
 * states.
 */
func TestDeterminizationProblematic(t *testing.T) {
	setup := Setup{
		1,
		0,
		true,
		deck.Card{Suit: deck.C, Value: deck.K},
		deck.C,
		deck.Card{},
		-1,
	}

	cards1 := []deck.Card{
		{Suit: deck.C, Value: deck.J},
		{Suit: deck.C, Value: deck.Ten},
		{Suit: deck.C, Value: deck.Q},
		{Suit: deck.C, Value: deck.K},
	}

	cards2 := []deck.Card{
		{Suit: deck.D, Value: deck.A},
		{Suit: deck.D, Value: deck.K},
		{Suit: deck.D, Value: deck.Q},
		{Suit: deck.D, Value: deck.J},
	}

	cards3 := []deck.Card{
		{Suit: deck.D, Value: deck.Ten},
		{Suit: deck.C, Value: deck.Nine},
		{Suit: deck.C, Value: deck.A},
		{Suit: deck.S, Value: deck.J},
	}

	played := []deck.Card{
		{Suit: deck.S, Value: deck.Q},
		{Suit: deck.S, Value: deck.K},
		{Suit: deck.H, Value: deck.A},
	}

	hand := []deck.Card{
		{Suit: deck.S, Value: deck.Nine},
		{Suit: deck.S, Value: deck.A},
	}

	trick1 := Trick{
		cards1,
		2,
		deck.C,
		-1,
	}

	trick2 := Trick{
		cards2,
		2,
		deck.C,
		-1,
	}

	trick3 := Trick{
		cards3,
		2,
		deck.C,
		-1,
	}

	prior := []Trick{trick1, trick2, trick3}

	state := NewUndeterminizedState(setup, 0, hand, played, prior)
	state.Determinize()

	fmt.Printf("Sanity check: %v\n", state)
	if len(state.Hands[0]) != 2 || len(state.Hands[1]) != 1 ||
		len(state.Hands[2]) != 1 || len(state.Hands[3]) != 1 {
		t.Errorf("Incorrect number of cards in some players hand.\n")
	}
}
