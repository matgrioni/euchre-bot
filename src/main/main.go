package main

import (
    "ai"
    "deck"
    "fmt"
)

func main() {
    fmt.Printf("Welcome to the Euchre AI!\n")

    var resp string
    fmt.Printf("Was it your deal?\n")
    fmt.Scanf("%s", &resp)
    friendly := resp == "y"

    fmt.Printf("Enter the top card.\n")

    var top deck.Card
    var value string
    var suit string
    fmt.Scanf("%s %s", &value, &suit)

    top.Value = deck.CreateValue(value);
    top.Suit = deck.CreateSuit(suit);

    fmt.Printf("Enter your hand to determine your call.\n")

    // Input the hand.
    var hand [5]deck.Card
    for i := 0; i < 5; i++ {
        fmt.Scanf("%s %s", &value, &suit)

        hand[i].Value = deck.CreateValue(value)
        hand[i].Suit = deck.CreateSuit(suit)
    }

    if ai.RPickUp(hand, top, friendly) {
        fmt.Printf("Pick it up!\n")
    } else {
        fmt.Printf("Pass...\n")
    }
}
