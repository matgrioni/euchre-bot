package main


import (
    "ai"
    "deck"
    "encoding/json"
    "euchre"
    "flag"
    "fmt"
    "math/rand"
    "time"
)


/*
 * This script randomly generates a data set to evaluate different strategies
 * on. This script will randomly generate a given number of samples along with
 * the utility for the user, player number 0, if euchre was a game of perfect
 * information and all players played optimally. The samples generated are used
 * to evaluate actual game play, not pre game play.
 *
 * Usage:
 *  ./gen_benchmark_play {numSamples} > data.txt
 *
 * numSamples are the number of situations you wish to compare.
 */


var r = rand.New(rand.NewSource(time.Now().UnixNano()))


/*
 * Create a random setup given the current situation. The random setup randomly
 * chooses who dealt and what trump is. Nobody is ever going alone though, and
 * nobody ever picks up. This allows us to focus on a given game state and not
 * worry about how other game logic components perform.
 *
 * Args:
 *  splits: The splits of the cards. Between players and the kitty.
 *
 * Returns:
 *  The randomized euchre setup.
 */
func randomNoAloneNoPickupSetup(splits [][]deck.Card) euchre.Setup {
    dealer := r.Intn(4)
    caller := r.Intn(4)
    pickedUp := false
    top := hands[4][3]
    trump := deck.SUITS[r.Intn(len(deck.SUITS))]
    var discard deck.Card

    return euchre.Setup {
        dealer,
        caller,
        pickedUp,
        top,
        trump,
        discard,
        -1,
    }
}


func main() {
    var samples int
    flag.IntVar(&samples, "samples", 0, "Number of sample games to simluate")
    flag.Parse()

    engine := euchre.Engine{ }
    for i := 0; i < samples; i++ {
        splits := euchre.GenSituation()
        setup := randomNoAloneNoPickupSetup(splits)

        played := make([]deck.Card, 0, 4)
        prior := make([]euchre.Trick, 0, 5)
        state := euchre.NewDeterminizedState(setup, (setup.Dealer + 1) % 4,
                                             hands, played, prior)

        score, _ := ai.Minimax(state, engine)

        stateStr, _ := json.Marshal(state)
        fmt.Printf("%s\t%f\n", stateStr, score)
    }
}
