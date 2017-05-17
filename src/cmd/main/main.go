package main

import (
    "bufio"
    "deck"
    "fmt"
    "os"
    "pickup"
    "player"
)

func main() {
    fmt.Print("Welcome to the Euchre AI!.\n")
    fmt.Print("Albert plays just like a human (sort of).\n")
    fmt.Print("This program plays a single hand (5 tricks) at a time.\n")
    fmt.Print("This includes picking the trump suit.\n")
    fmt.Println()

    fmt.Print("Enter the 5 cards in your hand...\n")
    var hand [5]deck.Card
    for i := 0; i< 5; i++ {
        var cardStr string
        fmt.Scanf("%s", &cardStr)
        hand[i] = deck.CreateCard(cardStr)
    }
    fmt.Println()

    fmt.Print("Enter the top card...\n")
    var topStr string
    fmt.Scanf("%s", &topStr)
    top := deck.CreateCard(topStr)
    fmt.Println()

    fmt.Print("You (2), Partner (1), Other (0)\n")
    fmt.Print("Enter whose deal it was...\n")
    var friend int
    fmt.Scanf("%d", friend)
    fmt.Println()

    var trump deck.Suit
    fn := os.Args[1]
    inputs, expected := pickup.LoadInputs(fn)
    if pickup.Perceptron(inputs, expected, hand, top, friend) {
        fmt.Print("Pick it up.\n")
        trump = top.Suit
    } else {
        fmt.Print("Pass.\n")
        fmt.Print("Enter the chosen trump suit...\n")
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
        fmt.Print("Cards already played (blank line when done)...\n")
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
