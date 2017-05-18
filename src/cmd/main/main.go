package main

import (
    "bufio"
    "deck"
    "fmt"
    "os"
    "euchre/pickup"
    "euchre/player"
)

func main() {
    fmt.Println("Welcome to the Euchre AI!.")
    fmt.Println("Albert plays just like a human (sort of).")
    fmt.Println("This program plays a single hand (5 tricks) at a time.")
    fmt.Println("This includes picking the trump suit.")
    fmt.Println()

    fmt.Println("Enter the 5 cards in your hand...")
    var hand [5]deck.Card
    for i := 0; i< 5; i++ {
        var cardStr string
        fmt.Scanf("%s", &cardStr)
        hand[i] = deck.CreateCard(cardStr)
    }
    fmt.Println()

    fmt.Println("Enter the top card...")
    var topStr string
    fmt.Scanf("%s", &topStr)
    top := deck.CreateCard(topStr)
    fmt.Println()

    fmt.Println("You (2), Partner (1), Other (0)")
    fmt.Println("Enter whose deal it was...")
    var friend int
    fmt.Scanf("%d", friend)
    fmt.Println()

    var trump deck.Suit
    fn := os.Args[1]
    inputs, expected := pickup.LoadInputs(fn)
    if pickup.Perceptron(inputs, expected, hand, top, friend) {
        fmt.Println("Pick it up.")
        trump = top.Suit
    } else {
        fmt.Println("Pass.")
        fmt.Println("Enter the chosen trump suit...\n")
        fmt.Scanln()
        var trumpStr string
        fmt.Scanf("%s", &trumpStr)
        trump = deck.CreateSuit(trumpStr)
    }
    fmt.Println()

    var chosen deck.Card
    curHand := hand[:]
    scanner := bufio.NewScanner(os.Stdin)
    for i := 0; i < 5; i++ {
        fmt.Printf("Hand %d\n", i + 1)
        fmt.Println("Cards already played (blank line when done)...")
        fmt.Scanln()

        played := make([]deck.Card, 0, 0)
        for scanner.Scan() {
            line := scanner.Text()
            if line == "" {
                break
            }

            played = append(played, deck.CreateCard(line))
        }

        chosen, curHand = player.Random(curHand, played, trump)

        fmt.Printf("Play %s.\n", chosen)
        fmt.Println()
    }
}
