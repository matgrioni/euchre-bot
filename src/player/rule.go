package player

import (
    "ai"
    "bufio"
    "deck"
    "euchre"
    "fmt"
    "math/rand"
    "os"
    "time"
)

// The input data type is used to represent an input into the perceptron
// based approach to calling for the dealer to pick it up or not. It implements
// the ai.Input interface by providing a method to export its data to a slice of
// features.
type Input struct {
    Top deck.Card
    Hand [5]deck.Card
    Dealer int
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
    if suitCount <= 2 && i.Dealer != 0 {
        features[11] = 1
    } else if suitCount <= 3 && i.Dealer == 0 {
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
    if i.Dealer == 0 && i.Dealer == 2 {
        features[indexes[i.Top.Value]] = 1
    }

    return features
}

type RulePlayer struct {
    pickupFn string
}

// Used to create a new RulePlayer struct that is properly constructed.
// pickupFn - The location of the file with the pickup / answer data samples.
// Returns a RulePlayer pointer.
func NewRule(pickupFn string) (*RulePlayer) {
    return &RulePlayer{ pickupFn, }
}

func (p *RulePlayer) Pickup(hand [5]deck.Card, top deck.Card, who int) bool {
    inputs, expected := loadInputs(p.pickupFn)
    prcp := ai.CreatePerceptron(12, 0, 1)

    // Move the perceptron toward linear separability if possible and then
    // return the result for the given input.
    prcp.Converge(inputs, expected, 0.005, 0.07, 10000)
    nextInput := Input {
        top,
        hand,
        who,
    }
    res := prcp.Process(nextInput)

    return res == 1
}

func (p *RulePlayer) Discard(hand [5]deck.Card,
                             top deck.Card) ([5]deck.Card, deck.Card) {
    // TODO: For now just use the random approach.
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    total := hand[:]
    total = append(total, top)

    i := r.Intn(len(total))
    chosen := total[i]
    total[i] = total[len(total) - 1]
    total = total[:len(total) - 1]

    copy(hand[:], total[:5])

    return hand, chosen
}

func (p *RulePlayer) Call(hand [5]deck.Card, top deck.Card, who int) (deck.Suit, bool) {
    chosen := false
    maxT := top.Suit
    maxConf := float32(0)

    for _, trump := range deck.SUITS {
        if trump == top.Suit {
            continue
        }

        conf := float32(0.08)

        // The following maps create a relation between a certain feature and its
        // position within the features vector.
        weights := map[deck.Value]float32 {
            deck.Nine: 0.05,
            deck.Ten: 0.07,
            deck.J: 0.3,
            deck.Q: 0.12,
            deck.K: 0.15,
            deck.A: 0.2,
        }

        // Used to keep track of how many suits are in the hand.
        suitsPresent := make(map[deck.Suit]int)

        for _, card := range hand {
            if card.Suit == trump {
                conf += weights[card.Value]
            } else if card.AdjSuit(trump) == trump {
                conf += 0.25
            } else if card.Value == deck.A {
                conf += 0.04
            }

            // Adjust suit count for left bower and increment the count for the
            // current card's suit.
            adjSuit := card.AdjSuit(trump)
            if _, ok := suitsPresent[adjSuit]; ok {
                suitsPresent[adjSuit] += 1
            } else {
                suitsPresent[adjSuit] = 1
            }
        }

        // If the hand has 2 suits or less we can have more confidence.
        if len(suitsPresent) <= 2 {
            conf += 0.08
        }

        if conf > 0.5 && conf > maxConf {
            maxConf = conf
            maxT = trump
            chosen = true
        }
    }

    return maxT, chosen
}

func (p *RulePlayer) Play(setup euchre.Setup, hand, played []deck.Card,
                          prior []euchre.Trick) ([]deck.Card, deck.Card) {
    // TODO: For now just use the random approach.
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    playable := euchre.Possible(hand, played, setup.Trump)

    chosen := playable[r.Intn(len(playable))]
    final := hand[chosen]
    hand[chosen] = hand[len(hand) - 1]
    hand = hand[:len(hand) - 1]

    return hand, final
}

// Loads the inputs in a file and returns a slice of the inputs and their
// expected values.
// fn - The filename where the inputs are located.
// Returns a slice of ai.Input values and a parallel slice of the expected
// values for each of these inputs.
func loadInputs(fn string) ([]ai.Input, []int) {
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
                                                    &nextInput.Dealer, &up)

        // Initialize the card from the values read in and add it to the samples
        // slice.
        nextInput.Top, _ = deck.CreateCard(tmpTop)
        for i, tmpCard := range tmpHand {
            nextInput.Hand[i], _ = deck.CreateCard(tmpCard)
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
