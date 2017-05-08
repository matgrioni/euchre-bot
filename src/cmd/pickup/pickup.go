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

    var line string
    fmt.Scanf("%s", &line)
    top := deck.CreateCard(line)

    fmt.Printf("Enter your hand to determine your call.\n")

    // Input the hand.
    var hand [5]deck.Card
    for i := range hand {
        fmt.Scanf("%s", &line)
        hand[i] = deck.CreateCard(line)
    }

    if ai.RPickUp(hand, top, friendly) {
        fmt.Printf("Pick it up!\n")
    } else {
        fmt.Printf("Pass...\n")
    }
}
