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

func (p *Perceptron) Process(inputs []int) int {
    s := p.bias
    for i, input := range inputs {
        s += float32(input) * p.weights[i]
    }

    if s > 0 {
        return 1
    }

    return 0
}

func (p *Perceptron) Adjust(inputs []int, delta int, learningRate float32) {
    // TODO: Is this right?
    for i, input := range inputs {
        p.weights[i] += float32(input) * float32(delta) * learningRate
    }

    p.bias += float32(delta) * learningRate
}

// TODO: Maybe interface type for training data samples.
func (p *Perceptron) Train(inputs [][]int, expected []int, rate float32) {
    for i, input := range inputs {
        actual := p.Process(input)
        del := expected[i] - actual

        p.Adjust(input, del, rate)
    }
}
