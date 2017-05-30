package main

import (
    "deck"
    "euchre/pickup"
    "fmt"
    "os"
)

func main() {
    // The first (and only) argument should be the location of the training
    // samples of the perceptron.
    fn := os.Args[1]
    inputs, expected := pickup.LoadInputs(fn)

    fmt.Printf("Welcome to the Euchre AI!\n")
    fmt.Printf("This is the perceptron based approach to picking up or not\n")

    var dealer int
    fmt.Printf("Did you(0), your partner(2), or neither(1/3) deal?\n")
    fmt.Scanf("%d", &dealer)

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

    if pickup.Perceptron(inputs, expected, hand, top, dealer) {
        fmt.Printf("Pick it up!\n")
    } else {
        fmt.Printf("Pass...\n")
    }
}
