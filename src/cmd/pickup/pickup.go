package main

import (
    "ai"
    "deck"
    "fmt"
)

func main() {
    fmt.Printf("Welcome to the Euchre AI!\n")

    var friendly int
    fmt.Printf("Did you(2) or your partner(1) or neither(0) deal (2/1/0)?\n")
    fmt.Scanf("%d", &friendly)

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

    conf := ai.RPickUp(hand, top, friendly)
    fmt.Printf("Confidence of %f\n", conf)
    if conf >= 0.5 {
        fmt.Printf("Pick it up!\n")
    } else {
        fmt.Printf("Pass...\n")
    }
}
