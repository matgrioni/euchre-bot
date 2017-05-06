package main

import (
    "euchre"
    "fmt"
)

func main() {
    fmt.Printf("Welcome to the Euchre AI!\n")

    fmt.Printf("Was it your deal?\n")
    var resp string
    fmt.Scanf("%s", &resp)
    friendly := resp == "y"

    fmt.Printf("Enter the top card.\n")

    var top euchre.Card
    var value int
    var suitStr string
    fmt.Scanf("%d %s", &value, &suitStr)

    top.Value = value;
    top.Suit = euchre.NewSuit(suitStr);

    fmt.Printf("Enter your hand to determine your call.\n")

    // Input the hand.
    var hand [5]euchre.Card
    for i := 0; i < 5; i++ {
        fmt.Scanf("%d %s", &value, &suitStr)

        hand[i].Value = value
        hand[i].Suit = euchre.NewSuit(suitStr)
    }

    if euchre.RPickUp(hand, top, friendly) {
        fmt.Printf("Pick it up!\n")
    } else {
        fmt.Printf("Pass...\n")
    }
}
