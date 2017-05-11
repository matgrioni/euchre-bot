package main

import (
    "ai"
    "bufio"
    "deck"
    "pickup"
    "fmt"
    "os"
)

func check(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    // The first (and only) argument should be the location of the training
    // samples of the perceptron.
    filename := os.Args[1]
    file, err := os.Open(filename)
    check(err)

    // Scan all the training data from the file into the samples slice.
    var samples []ai.Input
    var expected []int
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()

        var nextInput pickup.Input
        var tmpTop string
        var tmpHand [5]string
        var up int
        // TODO: Are parenthesis needed?
        // Read in a line from the file and parse it for the different needed
        // fields for a pickup problem instance.
        fmt.Sscanf(line, "%s %s %s %s %s %s %d %d", &tmpTop, &tmpHand[0],
                                                    &tmpHand[1], &tmpHand[2],
                                                    &tmpHand[3], &tmpHand[4],
                                                    &(nextInput.Friend), &up)

        // Initialize the card from the values read in and add it to the samples
        // slice.
        nextInput.Top = deck.CreateCard(tmpTop)
        for i, tmpCard := range tmpHand {
            nextInput.Hand[i] = deck.CreateCard(tmpCard)
        }

        samples = append(samples, nextInput)
        expected = append(expected, up)
    }


    fmt.Printf("Welcome to the Euchre AI!\n")
    fmt.Printf("This is the perceptron based approach to picking up or not\n")
    fmt.Printf("Control-C to quit.\n")

    p := pickup.InitialPerceptron()
    fmt.Print("These are the initial weights of the perceptron.\n")
    for _, weight := range p.Weights() {
        fmt.Printf("%.3f ", weight)
    }
    fmt.Println()

    for {
        var friendly int
        fmt.Printf("Did you(2), your partner(1), or neither(0) deal (2/1/0)?\n")
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

        if pickup.PPickUp(p, samples, expected, hand, top, friendly, 10) {
            fmt.Printf("Pick it up!\n")
        } else {
            fmt.Printf("Pass...\n")
        }

        fmt.Print("These are the new weights.\n")
        for _, weight := range p.Weights() {
            fmt.Printf("%.3f ", weight)
        }
        fmt.Println()
    }
}
