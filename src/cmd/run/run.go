package main

import (
    "deck"
    "fmt"
    "math/rand"
    "os"
    "player"
    "strconv"
    "time"
)


/*
 * A runner that runs a provided type of player's pickup (inital pickup or pass)
 * on a random sampling of cards so as to better understand its behavior. By
 * providing a player type and the percentage by which to choose each hand. This
 * percentage is the chance a given hand will be evaluated.
 *
 * Usage:
 *  ./run_pickup {playerType} {samples}
 *
 * playerType corresponds to the approach that will be ran. For a random player
 * use 1, for a rule based player use 2, for the supposedly smart approach use
 * 3. samples corresponds to how many different hands will be ran.
 */


// TODO: This should be moved into some configuration file.
const (
    PICKUP_CONF = 0.6
    CALL_CONF = 0.6
    ALONE_CONF = 1.2
    PICKUP_RUNS = 50
    PICKUP_DETERMINIZATIONS = 50
    CALL_RUNS = 50
    CALL_DETERMINIZATIONS = 50
    PLAY_RUNS = 50
    PLAY_DETERMINIZATIONS = 50
    ALONE_RUNS = 50
    ALONE_DETERMINIZATIONS = 50
)


func main() {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    players := make(map[int]player.Player)
    players[0] = player.NewRand(0.5, 0.5, 0)
    // TODO: Make non-hardcoded.
    players[1] = player.NewRule("data/train.dat")
    players[2] = player.NewSmart(PICKUP_CONF, CALL_CONF, ALONE_CONF,
                                 PICKUP_RUNS, PICKUP_DETERMINIZATIONS,
                                 CALL_RUNS, CALL_DETERMINIZATIONS,
                                 PLAY_RUNS, PLAY_DETERMINIZATIONS,
                                 ALONE_RUNS, ALONE_DETERMINIZATIONS)

    // TODO: Use more robust library rather than command line arguments.
    playerType, _ := strconv.Atoi(os.Args[1])
    samples, _ := strconv.Atoi(os.Args[2])
    player := players[playerType - 1]

    for i := 0; i < samples; i++ {
        situation := deck.DrawN(6)

        hand := situation[:5]
        top := situation[5]

        copyHand := make([]deck.Card, len(hand))
        copy(copyHand, hand)

        dealer := r.Intn(4)
        pickup := player.Pickup(copyHand, top, dealer)

        fmt.Printf("%v\t%s\t%d\t%t\n", hand, top, dealer, pickup)
    }
}
