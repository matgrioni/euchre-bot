package util

/*
 * A utility package for combinatorial functions. Includes items such as
 * combinations and permutations.
 */


/*
 * Compute the multinomial choices for a given n, and ks. This function returns
 * a chan which is used as a generator to iterate over all multinomial choices
 * of n numbers given ks. Note that the sum(ks) = n.
 *
 * Args:
 *  n,  type(int): The number of total choices there are.
 *  ks, type(int): A varidic argument for the numbers on the bottom of the
 *                 coefficient.
 *
 * Returns:
 *  type(chan []int): The integer slices are the multinomial choices. So every
 *                    len(ks) items from the channel make up one partition.
 */
func Multinomial(n int, ks ...int) chan []int {
    idxs := make([]int, n)
    for i := 0; i < len(idxs); i++ {
        idxs[i] = i
    }

    return multinomialHelper(idxs, ks...)
}


/*
 * The main logic for multinomials. This helper needed to be made since it was
 * much easier to make a function that took the numbers list to be used as the
 * sample space.
 *
 * Args:
 *  idxs,  type([]int): The numbers on which to do a multinomial choice.
 *  ks,   type(...int): The k's for selection on the n idxs.
 *
 * Returns:
 *  type(chan []int): A channel that provides the partitions of the idxs. This
 *                    one whole partition occurs every len(ks) chan items.
 */
func multinomialHelper(idxs []int, ks ...int) chan []int {
    ch := make(chan []int)
    n := len(idxs)

    go func() {
        if n > 0 {
            bch := Binomial(n, ks[0])

            for b := range bch {
                trans := make([]int, len(b))
                for i := 0; i < len(b); i++ {
                    trans[i] = idxs[b[i]]
                }

                left := make([]int, n)
                copy(left, idxs)
                l := len(b) - 1
                for i := l; i >= 0; i-- {
                    left = append(left[:b[i]], left[b[i] + 1:]...)
                }

                if len(ks) == 1 {
                    ch <- trans
                } else {
                    count := 0

                    for rc := range multinomialHelper(left, ks[1:]...) {
                        if (count % (len(ks) - 1) == 0) {
                            ch <- trans
                        }

                        c := make([]int, len(rc))
                        copy(c, rc)
                        ch <- c
                        count++
                    }
                }
            }
        }

        close(ch);
    }()

    return ch
}


/*
 * Creates a channel that lists all of the binomial combinations given some n
 * and some k. They are provided two at at time, so get them two at time.
 *
 * Args:
 *  n, type(int): The number of items to choose from.
 *  k, type(int): The number of items to choose.
 *
 * Returns:
 *  type(chan[]int): A channel that provides a slice that represent a choice
 *                   of k and n-k integers from n integers. These combinations
 *                   are provided one after the other. So get them two at a
 *                   time.
 */
func Binomial(n, k int) chan []int {
    ch := make(chan []int)

    go func() {
        comb := make([]int, k)
        if k > 0 {
            last := k - 1

            var rc func(int, int)
            rc = func(i, next int) {
                for j := next; j < n; j++ {
                    comb[i] = j

                    if i == last {
                        c := make([]int, len(comb))
                        copy(c, comb)
                        ch <- c
                    } else {
                        rc(i+1, j+1)
                    }
                }
            }
            rc(0, 0)
        } else {
            ch <- []int { }
        }

        close(ch)
    }()

    return ch
}
