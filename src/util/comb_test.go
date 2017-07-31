package util

import (
    "reflect"
    "testing"
)

/*
 * Test that the Binomial output for small numbers is in lexicographic order and
 * includes all the expected choices.
 */
func TestBinomialChoices(t *testing.T) {
    res := [][]int {
        []int { 0, 1 },
        []int { 0, 2 },
        []int { 0, 3 },
        []int { 0, 4 },
        []int { 1, 2 },
        []int { 1, 3 },
        []int { 1, 4 },
        []int { 2, 3 },
        []int { 2, 4 },
        []int { 3, 4 },
    }

    count := 0
    for choice := range Binomial(5, 2) {
        if !reflect.DeepEqual(choice, res[count]) {
            t.Errorf("Expected %s but got %s\n", res[count], choice)
        }

        count++
    }
}


/*
 * Test that there is only one empty combination when choosing 0 from n.
 */
func TestBinomialChoicesNone(t *testing.T) {
    count := 0
    for c := range Binomial(7, 0) {
        if len(c) > 0 {
            t.Errorf("Expected an empty combination but got %s\n", c)
        }

        count++
    }

    if count > 1 {
        t.Errorf("Expected one choice but got %d\n", count)
    }
}


/*
 * Test that there is only one way to choose all n elements from n.
 */
func TestBinomialChoicesAll(t *testing.T) {
    count := 0
    for c := range Binomial(9, 9) {
        if len(c) < 9 {
            t.Errorf("Expected all inclusive choice but got %s\n", c)
        }

        count++
    }

    if count > 1 {
        t.Errorf("Expected only one choice but got %d\n", count)
    }
}


/*
 * Test that for large binomial n and k, the amount of binomials generated is
 * correct.
 */
func TestBinomialSize(t *testing.T) {
    count := 0
    for _ = range Binomial(15, 5) {
        count++
    }

    res := 3003
    if count != res {
        t.Errorf("Expected %d combinations but got %d", res, count)
    }
}
