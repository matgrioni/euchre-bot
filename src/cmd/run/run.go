package main

import (
    "deck"
    "fmt"
    "math/rand"
    "os"
    "player"
    "strconv"
    "time"
    "util"
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


/*
 * Transform a list of indices into the corresponding cards.
 *
 * Args:
 *  indices, type(int): A list of indices that correspond to cards in deck.CARDS
 *
 * Returns:
 *  type([]deck.Card): A slice that is the mapping from the indices to cards.
 */
func indicesToCards(indices []int) []deck.Card {
    cards := make([]deck.Card, len(indices))
    for i, idx := range indices {
        cards[i] = deck.CARDS[idx]
    }

    return cards
}


func main() {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    players := make(map[int]player.Player)
    players[0] = player.NewRand()
    // TODO: Make non-hardcoded.
    players[1] = player.NewRule("data/train.dat")
    players[2] = player.NewSmart()

    // TODO: Use more robust library rather than command line arguments.
    playerType, _ := strconv.Atoi(os.Args[1])
    samples, _ := strconv.Atoi(os.Args[2])
    player := players[playerType - 1]

    for i := 0; i < samples; i++ {
        situation := util.RandMultinomial(24, 5, 1)

        hands0 := indicesToCards(situation[0])
        top := indicesToCards(situation[1])[0]

        var copyHand [5]deck.Card
        copy(copyHand[:], hands0)

        dealer := r.Intn(4)
        player.Pickup(copyHand, top, dealer)

        fmt.Printf("%v\t%s\t%d\t%t\n", hands0, top, dealer, pickup)
    }
}
