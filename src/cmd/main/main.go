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
    fmt.Println("Albert is basically the best euchre player ever.")
    fmt.Println("This program plays a single hand (5 tricks) at a time and")
    fmt.Println("includes picking the trump suit.")
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

    fmt.Println("You (0), Left (1), Partner (2), Right (3)")
    fmt.Println("Enter whose deal it was...")
    var dealer int
    fmt.Scanf("%d", &dealer)
    fmt.Println()

    fn := os.Args[1]
    inputs, expected := pickup.LoadInputs(fn)

    var (
        trump deck.Suit
        d deck.Card
    )
    pickedUp := pickup.Perceptron(inputs, expected, hand, top, dealer)
    if pickedUp {
        fmt.Println("Order it up.")
    } else {
        fmt.Println("Pass.")
        fmt.Println()

        var pickedUpIn int
        fmt.Println("But did somebody else order it up?")
        fmt.Scanf("%d", &pickedUpIn)
        pickedUp = pickedUpIn == 1

        if !pickedUp {
            if call, chosenSuit := call.Rule(top, hand); call {
                fmt.Printf("If possible call %s on second go around.\n", chosenSuit)
            } else {
                fmt.Println("Pass if it makes its way back to you.")
            }

            fmt.Println("Enter the eventual chosen trump suit...")
            var trumpStr string
            fmt.Scanf("%s", &trumpStr)
            trump = deck.CreateSuit(trumpStr)
        }
    }

    if pickedUp {
        trump = top.Suit

        if dealer == 0 {
            hand, d = discard.Rand(hand, top)
            fmt.Printf("Discard %s.\n", d)
        }
    }

    var prior []deck.Card
    var chosen deck.Card
    curHand := hand[:]
    scanner := bufio.NewScanner(os.Stdin)
    for i := 0; i < 5; i++ {
        fmt.Println()
        fmt.Printf("Trick %d\n", i + 1)
        fmt.Println("Cards already played (blank line when done)...")

        played := make([]deck.Card, 0, 3)
        for i := 0; i < 3; i++ {
            scanner.Scan()
            line := scanner.Text()
            if line == "" {
                break
            }

            card := deck.CreateCard(line)
            played = append(played, card)
        }

        chosen, curHand = player.AI(curHand, played, prior, d, dealer, pickedUp, top, trump)

        fmt.Printf("Play %s.\n", chosen)
        fmt.Println()

        fmt.Println("Enter the rest of the cards played.")
        left := 3 - len(played)
        for i := 0; i < left; i++ {
            scanner.Scan()
            line := scanner.Text()

            card := deck.CreateCard(line)
            played = append(played, card)
        }

        prior = append(prior, played...)
        prior = append(prior, chosen)
    }
}
