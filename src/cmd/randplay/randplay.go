package main

import (
    "deck"
    "fmt"
    "player"
)

func main() {
    fmt.Print("Welcome to the Euchre AI!\n")
    fmt.Print("This is the random approach to playing a card.\n")

    var numCards int
    fmt.Print("Please enter your hand.\n")
    fmt.Print("Cards in hand > ")
    fmt.Scanf("%d", &numCards)

    fmt.Printf("Now please enter %d cards.\n", numCards)
    hand := make([]deck.Card, numCards, numCards)
    for i := 0; i < numCards; i++ {
        var cardStr string
        fmt.Scanf("%s", &cardStr)
        hand[i] = deck.CreateCard(cardStr)
    }

    var numPlayed int
    fmt.Print("Please enter the cards already played in this trick.\n")
    fmt.Print("Cards played > ")
    fmt.Scanf("%d", &numPlayed)

    fmt.Printf("Now please enter %d cards.\n", numPlayed)
    played := make([]deck.Card, numPlayed, numPlayed)
    for i := 0; i < numPlayed; i++ {
        var cardStr string
        fmt.Scanf("%s", &cardStr)
        played[i] = deck.CreateCard(cardStr)
    }

    var trumpStr string
    fmt.Print("Please enter the trump suit.\n")
    fmt.Print("Trump suit > ")
    fmt.Scanf("%s", &trumpStr)
    trump := deck.CreateSuit(trumpStr)

    chosen, _ := player.Random(hand, played, nil, trump)
    fmt.Printf("Play %s.\n", chosen)
}
