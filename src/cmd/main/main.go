package main

import (
    "bufio"
    "deck"
    "flag"
    "fmt"
    "log"
    "euchre"
    "player"
    "os"
    "runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func inputValidCard() deck.Card {
    var cardStr string
    fmt.Scanf("%s", &cardStr)
    card, err := deck.CreateCard(cardStr)

    for err != nil {
        fmt.Println("Invalid input.")
        fmt.Scanf("%s", &cardStr)
        card, err = deck.CreateCard(cardStr)
    }

    return card
}

func main() {
    player := player.NewSmart()

    fmt.Println("Welcome to the Euchre AI!.")
    fmt.Println("Albert is basically the best euchre player ever.")
    fmt.Println("This program plays a single hand (5 tricks) at a time and")
    fmt.Println("includes picking the trump suit.")
    fmt.Println()

    fmt.Println("Enter the 5 cards in your hand...")
    var hand [5]deck.Card
    for i := 0; i < 5; i++ {
        hand[i] = inputValidCard()
    }
    fmt.Println()

    fmt.Println("Enter the top card...")
    top := inputValidCard()
    fmt.Println()

    fmt.Println("You (0), Left (1), Partner (2), Right (3)")
    fmt.Println("Enter whose deal it was...")
    var dealer int
    var caller int
    fmt.Scanf("%d", &dealer)
    fmt.Println()

    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    var (
        trump deck.Suit
        d deck.Card
    )
    pickedUp := player.Pickup(hand, top, dealer)
    if pickedUp {
        fmt.Println("Order it up.")
        caller = 0
    } else {
        fmt.Println("Pass.")
        fmt.Println()
        var pickedUpIn int
        fmt.Println("But did somebody else order it up?")
        fmt.Scanf("%d", &pickedUpIn)
        pickedUp = pickedUpIn == 1

        if pickedUp {
            fmt.Println("Who was it?")
            fmt.Scanf("%d", &caller)
        }

        if !pickedUp {
            if chosenSuit, call := player.Call(hand, top, dealer); call {
                fmt.Printf("If possible call %s on second go around.\n", chosenSuit)
                caller = 0
            } else {
                fmt.Println("Pass if it makes its way back to you.")
            }

            fmt.Println("Enter the eventual chosen trump suit...")
            var trumpStr string
            var err error
            fmt.Scanf("%s", &trumpStr)
            trump, err = deck.CreateSuit(trumpStr)

            for err != nil {
                fmt.Println("Invalid input.")
                fmt.Scanf("%s", &trumpStr)
                trump, err = deck.CreateSuit(trumpStr)
            }

            fmt.Println("Who called this suit?")
            fmt.Scanf("%d", &caller)
        }
    }

    if pickedUp {
        trump = top.Suit

        if dealer == 0 {
            hand, d = player.Discard(hand, top)
            fmt.Printf("Discard %s.\n", d)
        }
    }

    setup := euchre.Setup {
        dealer,
        caller,
        pickedUp,
        top,
        trump,
        d,
    }

    led := (dealer + 1) % 4
    var prior []euchre.Trick
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

            card, err := deck.CreateCard(line)
            for err != nil {
                fmt.Println("Invalid input.")
                scanner.Scan()
                line := scanner.Text()
                card, err = deck.CreateCard(line)
            }

            played = append(played, card)
        }

        curHand, chosen = player.Play(setup, curHand, played, prior)
        played = append(played, chosen)

        fmt.Printf("Play %s.\n", chosen)
        fmt.Println()

        fmt.Println("Enter the rest of the cards played.")
        left := 4 - len(played)
        for i := 0; i < left; i++ {
            scanner.Scan()
            line := scanner.Text()

            card, err := deck.CreateCard(line)
            for err != nil {
                fmt.Println("Invalid input.")
                scanner.Scan()
                line := scanner.Text()
                card, err = deck.CreateCard(line)
            }

            played = append(played, card)
        }

        var playedArr [4]deck.Card
        copy(playedArr[:], played[:])
        trick := euchre.Trick {
            playedArr,
            led,
            trump,
        }
        prior = append(prior, trick)
        led = euchre.Winner(played, trump, led)
    }
}
