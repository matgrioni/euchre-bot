# Albert: The Euchre AI

An attempt to create a euchre playing AI. This is a learning experience and as such, I've also chosen to use golang to implement it. So it should be fun, let's see what happens! The software will try to compare two several different approaches, several different AI methods, a rule based method, and a random method.

The AI algorithms in use now or that have already been implemented but are not currently in use are the perceptron for picking up cards and calling suit. Given some training set of the game state on reaching the AI, the perceptron will mold to match this input using 11 predefined features. It is not known whether these samples are linearly separable, so it remains to be seen whether some SVM or more complex model must be created to perfectly separate picking up from not picking up. For actual play,  a MCTS algorithm has been used with some success, but it remains to be seen how to elevate it above human ability. It seems like a combination of expected values and non-uniform random sampling may help.

Another approach that should be attempted is minimax across randomly sampled hands. In other words, iterate over all the possible hands of opponents and run minimax assuming these hands. Weight the chosen card to be played based on the likelihood of the hands and find the average expected win value for each card. This was actually the first thing attempted but I kinda sucked and couldn't get it to work because of non-iterative development. Another approach could combine this minimax based approach and MCTS based on the hand, since minimax for 5 cards has a much larger search space than for 3 cards.

## TODO

- Improve MCTS
    * Improve efficiency for more possible iterations.
    * Improve ability to win.
    * Add tests and sanity checks to make sure current code works.
- Add new minimax approach
- Figure out way to evaluate each approach automatically.

Enjoy!

-- Matias Grioni
