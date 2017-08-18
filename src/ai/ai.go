package ai

import (
    "math/rand"
    "time"
)

/*
 * Taken and inspired by appliedgo.net/perceptron
 */


/*
 * The Input type to the perceptron breaks down any struct or type into a vector
 * of features which function as the inputs to the perceptron. The features are
 * either there or aren't, so they only take on a value of 1 or 0.
 */
type Input interface {
    Features() []int
}



/*
 * A simple perceptron type that acts in the form w * x + b > 0. This means that
 * this Perceptron acts as a binary classifier with its given weight and bias
 * and that given some input vector, x, the output will follow the preceding
 * inequality. Note that bias and weights are separate, sometimes they are
 * included together.
 */
type Perceptron struct {
    weights []float32
    bias float32
}


/*
 * Create a random perceptron that has n weights. The weight numbers are
 * randomly given to start between the given numbers.
 *
 * Args:
 *  n: The number of weights.
 *  low: The lower bound for the random weights.
 *  high: The upper bound for the random weights.
 *
 * Returns:
 *  A pointer to a Perceptron with n random weights and a random bias.
 */
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


/*
 * Process the given input and return a classification of either 1 or 0.
 *
 * Args:
 *  input: The input object that is broken down into binary features.
 *
 * Returns:
 *  1 if the input was classified as yes (passed the threshold), and 0
 *  otherwise.
 */
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


/*
 * Based on an array of inputs and a parallel array of expected answers, train
 * the perceptron at the given rate.
 *
 * Args:
 *  inputs: The slice of inputs to train this Perceptron on.
 *  expected: A parallel slice to inputs for the expected values.
 *  rate: The rate at which the Perceptron learns.
 */
func (p *Perceptron) Train(inputs []Input, expected []int, rate float32) {
    for i, input := range inputs {
        actual := p.Process(input)
        del := expected[i] - actual

        p.adjust(input, del, rate)
    }
}


/*
 * Train the perceptron until a certain level of convergence. Basically this
 * method keeps training the perceptron on the same data until the percent of
 * wrong classifications is less than percent or until it has trained the
 * perceptron maxIter times.
 *
 * Args:
 *  inputs: A slice of inputs to train a Perceptron on.
 *  expected: The parallel slice to inputs that has the expected values for the
 *            inputs.
 *  rate: The rate at which the Perceptron should learn. This is parallel to
 *        rate used in Perceptron#Adjust.
 *  percent: The inaccuracy allowed for a Perceptron on the training set.
 *  maxIter: The maximum number of times to go through the entire training set
 *           before giving up on convergence.
 *
 * Returns:
 *  True if the Perceptron converged within the desired limit and false
 *  otherwise.
 */
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

            p.adjust(input, del, rate)

            if del != 0 {
                wrong++
            }
        }

        iter++
    }

    return wrong <= thres
}


/*
 * Get the weights for a Perceptron.
 *
 * Returns:
 *  The weights of a Perceptron.
 */
func (p *Perceptron) Weights() []float32 {
    return p.weights
}


/*
 * Get the bias of a Perceptron.
 *
 * Returns:
 *  The bias field of a Perceptron.
 */
func (p *Perceptron) Bias() float32 {
    return p.bias
}


/*
 * Adjust a perceptron based on a given input, difference and learning rate. The
 * input is what is given to the perceptron, delta is either 1, 0, or -1. It
 * should be the result of expected - actual, and learningRate is how fast the
 * perceptron should adjust. Therefore, this function moves a given perceptron
 * in the direction of delta at a rate determined by learningRate. The use of
 * input is that it provides the features that are activated by this input, and
 * which are not. Therefore only weights that are from activated features are
 * changed while weights related to features not in the input are unchanged.
 *
 * Args:
 *  input: The input to adjust to.
 *  delta: The direction in which to adjust.
 *  learningRate: The rate at which to adjust the perceptron weights.
 */
func (p *Perceptron) adjust(input Input, delta int, learningRate float32) {
    for i, input := range input.Features() {
        p.weights[i] += float32(input) * float32(delta) * learningRate
    }

    p.bias += float32(delta) * learningRate
}
