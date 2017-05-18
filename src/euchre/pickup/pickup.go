package pickup

import (
    "ai"
    "bufio"
    "deck"
    "fmt"
    "os"
)

// The input data type is used to represent an input into the perceptron
// based approach to calling for the dealer to pick it up or not. It implements
// the ai.Input interface by providing a method to export its data to a slice of
// features.
type Input struct {
    Top deck.Card
    Hand [5]deck.Card
    Friend int
}

// Converts an input to a Perceptron used to tell a player to tell the dealer to
// pick up a card or not to a vector of binary features. The slice is of size 12
// and has the following features in order:
// - Nine of trump
// - Ten of trump
// - J of trump
// - Q of trump
// - K of trump
// - A of trump
// - J of same color as trump
// - AH
// - AD
// - AS
// - AC
// - Hand only has 2 suits.
func (i Input) Features() []int {
    features := make([]int, 12, 12)

    // The following maps create a relation between a certain feature and its
    // position within the features vector.
    indexes := map[deck.Value]int {
        deck.Nine: 0,
        deck.Ten: 1,
        deck.J: 2,
        deck.Q: 3,
        deck.K: 4,
        deck.A: 5,
    }

    aces := map[deck.Suit]int {
        deck.H: 7,
        deck.D: 8,
        deck.S: 9,
        deck.C: 10,
    }

    // Used to keep track of how many suits are in the hand.
    suitsPresent := make(map[deck.Suit]int)

    for _, card := range i.Hand {
        if card.Suit == i.Top.Suit {
            features[indexes[card.Value]] = 1
        } else if card.Suit == i.Top.Suit.Left() && card.Value == deck.J {
            features[6] = 1
        } else if card.Value == deck.A {
            features[aces[card.Suit]] = 1
        }

        // Adjust suit count for left bower and increment the count for the
        // current card's suit.
        adjSuit := card.AdjSuit(i.Top.Suit)
        if _, ok := suitsPresent[adjSuit]; ok {
            suitsPresent[adjSuit] += 1
        } else {
            suitsPresent[adjSuit] = 1
        }
    }

    suitCount := len(suitsPresent)
    // If the hand has less then 2 suits and no card is being picked up then one
    // can be sure we have 2 suits.
    if suitCount <= 2 && i.Friend != 2 {
        features[11] = 1
    } else if suitCount <= 3 && i.Friend == 2 {
    // Else, if there are less than 3 suits and we are picking up, one suit
    // might be gotten rid of yet. Even if we had 2 suits we would to check that
    // we already had the trump suit if we are picking up.
        _, trumpPresent := suitsPresent[i.Top.Suit]
        if suitCount <= 2 && trumpPresent {
            features[11] = 1
        } else if suitCount == 3 && trumpPresent {
            for _, card := range i.Hand {
                // If this is the only card in the hand of a given suit, and it
                // is not a trump (left bower or of the same suit as the top
                // card), and it is not an A, then by removing it we get rid of
                // one suit. So if there were only 3 suits, then this means we
                // have 2 suits after discarding if we were to pick up. We also
                // have to make sure that this removal will result in at 2 suits
                // and that picking up the top card does not increase the suit
                // amount.
                if suitsPresent[card.Suit] == 1 && card.Value != deck.A &&
                   card.AdjSuit(i.Top.Suit) != i.Top.Suit {
                    features[11] = 1

                    // There could be more than one card that matches these
                    // requirements, we only care if one exists, so break on finding
                    // one.
                    break
                }
            }
        }
    }

    // If you are picking it up then consider the weight of the card on top as
    // well. The discarded card need not be considered since, it is most likely
    // non-trump, and will not affect this analysis which only considers trump
    // cards. If it is a trump card, then every other card is a trump as well
    // and the final decision should not be affected.
    if i.Friend != 0 {
        features[indexes[i.Top.Value]] = 1
    }

    return features
}

// A rule based approach to deciding if one should call for the dealer to pick
// up the top card at the start of a deal. This approach takes into account the
// cards currently in the players hand, the top card, and whether one or one's
// partner picks up the card in question. The friend parameter can take the
// following values. 2 for you are picking it up, 1 for your partner is picking
// it up, and 0 neither.
func Rule(hand [5]deck.Card, top deck.Card, friend int) bool {
    i := Input {
        top,
        hand,
        friend,
    }

    // One wants to tell the dealer to pick up if you believe that you can win 3
    // tricks if the dealer picks it up. So if your confidence / probability is
    // greater than 50% (in a non-technical sense) then you should take such an
    // action. Start off with a little confidence that your partner will win at
    // least one trick.
    bias := float32(0.08)

    // Arbitrary weights for the features that I felt like using based on
    // experience.
    weights := map[int]float32 {
        0: 0.05,
        1: 0.07,
        2: 0.3,
        3: 0.12,
        4: 0.15,
        5: 0.2,
        6: 0.25,
        7: 0.04,
        8: 0.04,
        9: 0.04,
        10: 0.04,
        11: 0.08,
    }

    // Do the dot product between the binary feature vector and the weights for
    // each respective feature.
    conf := bias
    for i, feat := range i.Features() {
        conf += weights[i] * float32(feat)
    }

    return conf >= 0.5
}

// Determine if the given hand should be picked up or not using Perceptron logic.
// Provide the inputs to the perceptron and a parallel array of the expected
// answers, along with the current problem instance and true is returned if it
// should be picked up and false otherwise.
func Perceptron(inputs []ai.Input, expected []int, hand [5]deck.Card,
                top deck.Card, friend int) bool {
    p := ai.CreatePerceptron(12, 0, 1)

    // Output debugging info.
    fmt.Print("These are the initial weights of the perceptron.\n")
    for _, weight := range p.Weights() {
        fmt.Printf("%.3f ", weight)
    }
    fmt.Printf("%.3f\n", p.Bias())

    // Check if the perceptron converged.
    conv := p.Converge(inputs, expected, 0.005, 0.07, 10000)
    if conv {
        fmt.Print("Converged\n")
    } else {
        fmt.Print("Did not converge\n")
    }

    fmt.Print("These are the final weights.\n")
    for _, weight := range p.Weights() {
        fmt.Printf("%.3f ", weight)
    }
    fmt.Printf("%.3f\n", p.Bias())

    nextInput := Input {
        top,
        hand,
        friend,
    }
    res := p.Process(nextInput)

    return res == 1
}

// Loads the inputs in a file and returns a slice of the inputs and their
// expected values.
// fn - The filename where the inputs are located.
// Returns a slice of ai.Input values and a parallel slice of the expected
// values for each of these inputs.
func LoadInputs(fn string) ([]ai.Input, []int) {
    file, err := os.Open(fn)
    check(err)
    scanner := bufio.NewScanner(file)

    // Scan all the training data from the file into the samples slice.
    var samples []ai.Input
    var expected []int
    for scanner.Scan() {
        line := scanner.Text()

        // Declare all variables needed to input a sample instance with the
        // input (top and hand) and the answer as well (up).
        var nextInput Input
        var tmpTop string
        var tmpHand [5]string
        var up int
        // Read in a line from the file and parse it for the different needed
        // fields for a pickup problem instance.
        fmt.Sscanf(line, "%s %s %s %s %s %s %d %d", &tmpTop, &tmpHand[0],
                                                    &tmpHand[1], &tmpHand[2],
                                                    &tmpHand[3], &tmpHand[4],
                                                    &nextInput.Friend, &up)

        // Initialize the card from the values read in and add it to the samples
        // slice.
        nextInput.Top = deck.CreateCard(tmpTop)
        for i, tmpCard := range tmpHand {
            nextInput.Hand[i] = deck.CreateCard(tmpCard)
        }

        samples = append(samples, nextInput)
        expected = append(expected, up)
    }

    return samples, expected
}

// Simple utility method to check and abort if there was an error.
func check(err error) {
    if err != nil {
        panic(err)
    }
}
