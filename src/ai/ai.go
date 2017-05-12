package ai

import (
    "math/rand"
    "time"
)

// Taken and inspired by appliedgo.net/perceptron

// A simple perceptron type that acts in the form w * x + b > 0. This means that
// this Perceptron acts as a binary classifier with its given weight and bias
// and that given some input vector, x, the output will follow the preceding
// inequality. Note that bias and weights are separate, sometimes they are
// included together.
type Perceptron struct {
    weights []float32
    bias float32
}

// The Input type to the perceptron breaks down any struct or type into a vector
// of features which function as the inputs to the perceptron. The features are
// either there or aren't, so they only take on a value of 1 or 0.
type Input interface {
    Features() []int
}

// Create a random perceptron that has n weights. The weight numbers are
// randomly given to start between the given numbers.
func CreatePerceptron(n int, low, high float32) *Perceptron {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    weights := make([]float32, n, n)
    for i := range weights {
        weights[i] = r.Float32() * (high - low) + low
    }

    return &Perceptron {
        weights,
        r.Float32() * (high - low) + low,
    }
}

// Process the given input and return a classification of either 1 or 0.
func (p *Perceptron) Process(input Input) int {
    s := p.bias
    for i, input := range input.Features() {
        s += float32(input) * p.weights[i]
    }

    if s > 0 {
        return 1
    }

    return 0
}

// Adjust a perceptron based on a given input, difference and learning rate. The
// input is what is given to the perceptron, delta is either 1, 0, or -1. It
// should be the result of expected - actual, and learningRate is how fast the
// perceptron should adjust. Therefore, this function moves a given perceptron
// in the direction of delta at a rate determined by learningRate. The use of
// input is that it provides the features that are activated by this input, and
// which are not. Therefore only weights that are from activated features are
// changed while weights related to features not in the input are unchanged.
func (p *Perceptron) Adjust(input Input, delta int, learningRate float32) {
    for i, input := range input.Features() {
        p.weights[i] += float32(input) * float32(delta) * learningRate
    }

    p.bias += float32(delta) * learningRate
}

// Based on an array of inputs and a parallel array of expected answers, train
// the perceptron at the given rate.
func (p *Perceptron) Train(inputs []Input, expected []int, rate float32) {
    for i, input := range inputs {
        actual := p.Process(input)
        del := expected[i] - actual

        p.Adjust(input, del, rate)
    }
}

// Train the perceptron until a certain level of convergence. Basically this
// method keeps training the perceptron on the same data until the percent of
// wrong classifications is less than percent or until it has trained the
// perceptron maxIter times.
func (p *Perceptron) Converge(inputs []Input, expected []int,
                              rate, percent float32, maxIter int) bool {
    thres := int(percent * float32(len(inputs)))

    iter := 0
    wrong := thres + 1
    for wrong > thres && iter < maxIter {
        wrong = 0

        for i, input := range inputs {
            actual := p.Process(input)
            del := expected[i] - actual

            p.Adjust(input, del, rate)

            if del != 0 {
                wrong++
            }
        }

        iter++
    }

    return wrong <= thres
}

func (p *Perceptron) Weights() []float32 {
    return p.weights
}

func (p *Perceptron) Bias() float32 {
    return p.bias
}
