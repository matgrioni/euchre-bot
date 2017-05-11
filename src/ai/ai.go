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

func (p *Perceptron) Adjust(input Input, delta int, learningRate float32) {
    for i, input := range input.Features() {
        // TODO: Is this right?
        p.weights[i] += float32(input) * float32(delta) * learningRate
    }

    p.bias += float32(delta) * learningRate
}

// TODO: Maybe interface type for training data samples?
func (p *Perceptron) Train(inputs []Input, expected []int, rate float32) {
    for i, input := range inputs {
        actual := p.Process(input)
        del := expected[i] - actual

        p.Adjust(input, del, rate)
    }
}

func (p *Perceptron) Weights() []float32 {
    return p.weights
}

func (p *Perceptron) Bias() float32 {
    return p.bias
}
