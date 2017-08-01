package util

import (
    "math/rand"
    "time"
)

/*
 * A utility package for combinatorial functions. Includes items such as
 * combinations and permutations.
 */


var r = rand.New(rand.NewSource(time.Now().UnixNano()))


/*
 * Generates a random multinomial given the number of items to choose from and
 * the way to choose the items. Note that sum(ks) should equal n. Also note that,
 * this algorithm has n! responses. In other words, duplicates with different
 * orderings are possible. But since, they are used one at a time the
 * probability of getting any specific one stays the same.
 *
 * Args:
 *  n, type(int): The number of items to choose from.
 *  ks, type(...int): The way to choose to the n items.
 *
 * Returns:
 *  A slice of integer slices. There are len(ks) integer slices, each of which
 *  corresponds to a random choice from the n items.
 */
func RandMultinomial(n int, ks ...int) [][]int {
    perm := r.Perm(n)

    res := make([][]int, len(ks))

    last := 0
    for i := 0; i < len(ks); i++ {
        next := make([]int, 0)
        k := ks[i]

        for j := 0; j < k; j++ {
            next = append(next, perm[last + j])
        }

        last += k
        res[i] = next
    }

    return res
}
