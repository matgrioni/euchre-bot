package player

import (
    "deck"
    "euchre"
)

// The center of the euchre bot is the different types of players. When running
// the program one can choose between several different strategies that are used
// throughout all decision making processes. These options will be for example,
// random, AI, rule based and possibly more. These options are chosen through a
// struct that implements this Player interface. The underlying logic will then
// follow different models.
type Player interface {
    // Returns whether or not a player should tell the top card to be ordered up
    // based on their current cards and who is picking it up.
    // hand - The 5 cards currently in the player's hand. TODO: Change to slice.
    // top  - The card on top of the kitty and currenty in question to be
    //        picked up.
    // who  - Who is picking up the card (the dealer). The number designation
    //        for each player is as follows. Yourself(0), partner(2), opp to
    //        left (1), opp to right (3). So a clockwise order.
    // Returns true if the card should be ordered up when the player gets the
    // chance and false otherwise.
    Pickup(hand [5]deck.Card, top deck.Card, who int) bool

    // Determines what to discard if a player has just picked up the top card
    // after their deal.
    // hand - The 5 cards currently in the dealer's hand.
    // top  - The card that was on top and is now to be picked up.
    // Returns the new hand after discarding and the card to be discarded.
    Discard(hand [5]deck.Card, top deck.Card) ([5]deck.Card, deck.Card)

    // Determines whether a player should call a certain suit, such as when all
    // players have passed on picking it up.
    // hand - The 5 cards currently in the player's hand.
    // top  - The card that was passed by all players and was on the kitty.
    // who  - The player who dealt the cards.
    // Returns the suit that should be called if given the chance. This result
    // valid iff true is returned as well. Otherwise, the returned suit is
    // meaningless.
    Call(hand [5]deck.Card, top deck.Card, who int) (deck.Suit, bool)

    // Determines which card to play given the current euchre situation.
    // setup  - The setup of the euchre game before any tricks which consists
    //          of who was dealer, what the top card was, etc.
    // hand   - The cards currently in the user's hand.
    // played - The cards that have already been played in this trick.
    // prior  - The cards that have been played in previous tricks.
    // Returns the user's new hand and the card that was chosen from the user's
    // hand.
    Play(setup euchre.Setup, hand []deck.Card, played []deck.Card,
         prior []euchre.Trick) ([]deck.Card, deck.Card)
}
