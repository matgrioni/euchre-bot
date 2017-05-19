package main

import (
    "bufio"
    "deck"
    "euchre/call"
    "euchre/discard"
    "euchre/pickup"
    "euchre/player"
    "fmt"
    "os"
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
    fmt.Scanf("%d", &friend)
    fmt.Println()

    var trump deck.Suit
    fn := os.Args[1]
    inputs, expected := pickup.LoadInputs(fn)
    orderedUp := pickup.Perceptron(inputs, expected, hand, top, friend)
    if orderedUp {
        fmt.Println("Order it up.")
        trump = top.Suit

        if friend == 2 {
            var d deck.Card
            hand, d = discard.Rand(hand, top)

            fmt.Printf("Discard %s.\n", d)
        }
    } else {
        fmt.Println("Pass.")
        fmt.Println()

        if call, chosenSuit := call.Rule(top, hand); call {
            fmt.Printf("If possible call on second go around %s.\n", chosenSuit)
        } else {
            fmt.Println("Pass if it makes its way back to you.")
        }

        fmt.Println("Enter the eventual chosen trump suit...")
        var trumpStr string
        fmt.Scanf("%s", &trumpStr)
        trump = deck.CreateSuit(trumpStr)
        fmt.Printf("This is trump: %s.\n", trumpStr)
    }

    if !orderedUp && friend == 2 {
        var pickedUp int
        fmt.Println("Did you pick it up (1/0)?")
        fmt.Scanf("%d", &pickedUp)

        if pickedUp == 1 {
            var d deck.Card
            hand, d = discard.Rand(hand, top)

            fmt.Printf("Discard %s.\n", d)
        }
    }

    fmt.Println()

    var chosen deck.Card
    curHand := hand[:]
    scanner := bufio.NewScanner(os.Stdin)
    for i := 0; i < 5; i++ {
        fmt.Printf("Trick %d\n", i + 1)
        fmt.Println("Cards already played (blank line when done)...")

        played := make([]deck.Card, 0, 4)
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
