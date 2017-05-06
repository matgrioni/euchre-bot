package main

import (
    "euchre"
    "fmt"
)

func main() {
    hand := euchre.GenHand()

    fmt.Printf(hand[4].Suit.String())
}
