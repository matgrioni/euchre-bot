package main

import (
    "bufio"
    "deck"
    "fmt"
    "math/rand"
    "os"
    "pickup"
    "time"
)

// TODO: Isolate randomness to it's own local package.
var r *rand.Rand

func GenPickupInput() pickup.Input {
    r = rand.New(rand.NewSource(time.Now().UnixNano()))

    hand := deck.GenHand()

    var top deck.Card
    inHand := true
    for inHand {
        top = deck.GenCard()

        for _, card := range hand {
            inHand = false
            if card == top {
                inHand = true
                break
            }
        }
    }

    return pickup.Input {
        top,
        hand,
        r.Intn(3),
    }
}

func check(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    filename := os.Args[1]
    file, err := os.OpenFile(filename, os.O_RDWR, 0600)
    check(err)

    // Read in the current existing samples into a map that easily tracks which
    // problem instances have already been determined.
    var samples map[pickup.Input]bool = make(map[pickup.Input]bool)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()

        var nextSample pickup.Input
        var tmpTop string
        var tmpHand [5]string
        var tmpFriend int
        fmt.Sscanf(line, "%s %s %s %s %s %s %d", &tmpTop, &tmpHand[0],
                                                 &tmpHand[1], &tmpHand[2],
                                                 &tmpHand[3], &tmpHand[4],
                                                 &tmpFriend)

        nextSample.Top = deck.CreateCard(tmpTop)
        for i, tmpCard := range tmpHand {
            nextSample.Hand[i] = deck.CreateCard(tmpCard)
        }
        nextSample.Friend = tmpFriend

        samples[nextSample] = true
    }

    // Move the file pointer to the end of the file.
    file.Seek(0, 2)

    fmt.Print("Each line that is generated is a new test sample.\n")
    fmt.Print("Enter 1 for it is picked up, 0 to pass, and -1 to quit.\n")
    fmt.Printf("%-10s\t%-20s\t%-10s\n", "Top", "Hand", "Friend")
    for {
        ps := GenPickupInput()

        // If this generated sample already exists generate a new one until it
        // it is a new one.
        for _, in := samples[ps]; in ; _, in = samples[ps] {
            ps = GenPickupInput()
        }

        // Output the new sample and get the user's decision.
        var handStr string
        for _, card := range ps.Hand {
            handStr += card.String() + " "
        }
        handStr = handStr[:len(handStr) - 1]
        fmt.Printf("%-10s\t%-20s\t%-10d ", ps.Top, handStr, ps.Friend)

        var orderedUp int
        fmt.Scanf("%d", &orderedUp)
        if orderedUp < 0 {
            break
        }

        fmt.Fprintf(file, "%s %s %d %d\n", ps.Top, handStr, ps.Friend, orderedUp)
    }

    file.Close()
}
