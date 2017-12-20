package main


import (
    "ai"
    "bufio"
    "deck"
    "euchre"
    "flag"
    "fmt"
    "encoding/json"
    "log"
    "os"
    "player"
    "strconv"
    "strings"
)


/*
 * Benchmark an implementation against optimal players. Provide the location to
 * data that has a given state and the eventual score assuming perfect
 * information and perfect players. Given a player type mapped to a number this
 * script outputs the overall difference between the output of the given player
 * type and the optimal player and also outputs the results for each individual
 * data sample.
 *
 * Usage:
 *  ./benchmark_play {dataFile} {playerType}
 *
 * dataFile is the location of the minimax evaluated hands. playerType is the
 * type of player to run on these situations. The mapping from playerType to
 * player is as follows:
 *  0: MCTS
 *  1: RULE
 *  2: RANDOM
 */


const (
    PICKUP_CONF = 0.6
    CALL_CONF = 0.6
    ALONE_CONF = 1.2
    PICKUP_RUNS = 5000
    PICKUP_DETERMINIZATIONS = 50
    CALL_RUNS = 5000
    CALL_DETERMINIZATIONS = 50
    PLAY_RUNS = 5000
    PLAY_DETERMINIZATIONS = 50
    ALONE_RUNS = 5000
    ALONE_DETERMINIZATIONS = 50
)


func main() {
    var dataLoc string
    var playerType int
    flag.StringVar(&dataLoc, "dataLoc", "", "Location of minimax evaluated games.")
    flag.IntVar(&playerType, "playerType", 0, "Tye type of player to evaluate.")
    flag.Parse()

    // Create the mapping of playerType to player object and get the desired
    // player to evaluate.
    players := make(map[int]player.Player)
    players[0] = player.NewSmart(PICKUP_CONF, CALL_CONF, ALONE_CONF,
                                  PICKUP_RUNS, PICKUP_DETERMINIZATIONS,
                                  CALL_RUNS, CALL_DETERMINIZATIONS,
                                  PLAY_RUNS, PLAY_DETERMINIZATIONS,
                                  ALONE_RUNS, ALONE_DETERMINIZATIONS)
    players[1] = player.NewRule("data/train.dat")
    players[2] = player.NewRand(0.5, 0.5, 0)
    chosenPlayer := players[playerType]

    dataFile, err := os.Open(dataLoc)
    if err != nil {
        log.Fatal(err)
    }
    defer dataFile.Close()

    scanner := bufio.NewScanner(dataFile)
    for scanner.Scan() {
        // Parse each line to get the minimax evaluation, and the actual game
        // state.
        line := scanner.Text()
        tabIndex := strings.IndexRune(line, '\t')

        var state euchre.State
        stateStr := line[:tabIndex]
        json.Unmarshal([]byte(stateStr), &state)
        minimaxEval, _ := strconv.ParseFloat(line[tabIndex + 1:], 64)


        // Now that we have the game state, we can simulate the game. Using
        // Minimax players for the opponents and the desired player strategy
        // for user 0.
        engine := euchre.Engine{ }

        for i := 0; i < 5; i++ {

            var last int
            for j := 0; j < 4; j++ {
                last = state.Player
                // If it is the AI's turn, use the chosen player logic to choose
                // what card to use next. Then keep the state updated, so that the
                // Minimax agents know what is going on.
                if state.Player == 0 {
                    curHand, chosen := chosenPlayer.Play(state.Setup, state.Hands[0], state.Played, state.Prior)
                    state.Played = append(state.Played, chosen)
                    state.Hands[0] = curHand
                    state.Player = (state.Player + 1) % 4
                } else {
                // If it is the Minimax agents' turn, use their logic. This agent
                // provides the successor state as well so just use that.
                    _, chosenMove := ai.Minimax(state, engine)
                    state = chosenMove.State.(euchre.State)
                }
            }

            // If the last turn was the AI's then the last trick was not added
            // automatically to the state. Further we know, who led, since we
            // were last. Namely, player 1.
            // TODO: This is internal logic that is being handled very closely
            // by the outside program. This should be encapsulated.
            if last == 0 {
                trick := euchre.Trick {
                    state.Played,
                    1,
                    state.Setup.Trump,
                    state.Setup.AlonePlayer,
                }
                state.Played = make([]deck.Card, 0, 4)
                state.Prior = append(state.Prior, trick)

                state.Player = euchre.Winner(state.Played, state.Setup.Trump, 1,
                                    state.Setup.AlonePlayer)
            }
        }

        playerScore := engine.Evaluation(state)
        diff := minimaxEval - playerScore

        fmt.Printf("%s\t%f\n", stateStr, diff)
    }


    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}
