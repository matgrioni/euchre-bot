# Albert: The Euchre AI

An attempt to create a euchre playing AI. This is a learning experience and as such, I've also chosen to use golang to implement it. So it should be fun, let's see what happens! The software will try to compare two several different approaches, several different AI methods, a rule based method, and a random method.

The AI algorithms in use now or that have already been implemented but are not currently in use are the perceptron for picking up cards and calling suit. Given some training set of the game state on reaching the AI, the perceptron will mold to match this input using 11 predefined features. It is not known whether these samples are linearly separable, so it remains to be seen whether some SVM or more complex model must be created to perfectly separate picking up from not picking up. For actual play,  a MCTS algorithm has been used with some success, but it remains to be seen how to elevate it above human ability. It seems like a combination of expected values and non-uniform random sampling may help.

Another approach that should be attempted is minimax across randomly sampled hands. In other words, iterate over all the possible hands of opponents and run minimax assuming these hands. Weight the chosen card to be played based on the likelihood of the hands and find the average expected win value for each card. This was actually the first thing attempted but I kinda sucked and couldn't get it to work because of non-iterative development. Another approach could combine this minimax based approach and MCTS based on the hand, since minimax for 5 cards has a much larger search space than for 3 cards.

## Structure

The game agnostic AI libraries live in `src/ai`. These libraries are used by the different player logic found in `src/player`. Each file in this module represents a different playing strategy embodied by the player interface which each file implements through a structure like `SmartPlayer` or `RandomPlayer`. The `src/deck` module has definitions and helper functions for manipulating a card deck, and the `src/euchre` module has helper methods on top of the `deck` module that deal with euchre play.


## Results

Initial results are now available. Results were obtained as follows:

First a test dataset had to be created. I randomly generated various euchre situations and randomly assigned the dealer, caller, and trump suit (this is not realistic but it equalizes the negative effects amongst all players). Then I ran through this situation with an open hand (no hidden information) and Minimax players. There are 1000 such annotated situations. Evaluation boiled down to replacing one of the optimal players with a given implementation and playing through the game again with all information known to all players except the player being tested. This is meant to show how close a given implementation comes to being optimal.

Below is a bar chart that shows the average difference between the optimal player and different player implementations.

INSERT IMAGE

Below is a distribution of the differences. The differences are the integers between 0 (no difference) and 6 (minimax gets a loner, but player implementation is euched).

INSERT IMAGE

So MCTS is definitely the best out of the three, but not much better than the rule implementation. The only difference between the two is 19 situations which have a difference of 3 under the rule player and a difference of 0 under the MCTS player. Across all possible euchre hands this might be more noticeable. It should be noted however that not all of these hands are equally likely especially with regards to the caller and suit picked, and so it is possible that these differences are not very important. So while MCTS seems to be better, it still has a way to go before it is worth as an alternative to rule based systems for euchre, especially if complex rules can be used.

Note that this system is using a very, very vanilla MCTS approach. The choosing of successor nodes is done purely by UCB, whereas top notch Google implementations use a deep learning approach to focus on successor nodes that seem promising. Improving on the current MCTS can be done in a few ways that I'm interested in. One is determinizations. Currently determinization are polled uniformly, given the constraints of prior tricks. However, we can leverage prior play and calling information to weight the results of more likely determinizations. Also, while I do not have data for this yet, it seemed like increasing determinizations did not help as much as increasing number of playouts. The average case in Euchre is a shitty hand, so polling more of these will not help, but investigating them more is helpful. This leads me to believe that focusing on the distribution of high value cards rather than all cards left is key (i.e. where are the trumps, not the non-trumps). Another area for improvement as already mentioned would be the choosing of successors with more than vanilla UCB. A key requirement for this to be useful, however, would be a general approach and non domain specific approach as is currently implemented.


## TODO

- Improve MCTS
    * Improve efficiency for more possible iterations.
    * Improve ability to win.
    * Add tests and sanity checks to make sure current code works.
- ~~Add new minimax approach~~
- ~~Figure out way to evaluate each approach automatically.~~
- Evaluate picking up, calling, and going alone.
- Add a MLPlayer along with the Smart, combinatorially good player.

Enjoy!

-- Matias Grioni
