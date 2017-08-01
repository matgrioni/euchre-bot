package util

import (
    "fmt"
    "testing"
)


/*
 * Test the Random multinomial for a normal situation for the sizes of the
 * returned slices.
 */
func TestRandMultinomialSizes(t *testing.T) {
    res := RandMultinomial(6, 3, 2, 1)
    sizes := []int { 3, 2, 1 }

    fmt.Printf("Example for (6, 3, 2, 1): %v\n", res)

    for i, slice := range res {
        if len(slice) != sizes[i] {
            t.Errorf("Expected size of slice %d to be %d but got %d.\n",
                     i + 1, sizes[i], len(slice))
        }
    }
}


/*
 * Tests to make sure all the numbers in the multinomial are not too big
 * or small and are not duplicates.
 */
func TestRandMultinomialValues(t *testing.T) {
    res := RandMultinomial(8, 3, 2, 3)
    fmt.Printf("Example for (8, 3, 2, 3): %v\n", res)

    seen := make(map[int]bool)

    for _, slice := range res {
        for _, item := range slice {
            if _, ok := seen[item]; ok {
                t.Errorf("Expected %d only once.\n", item)
            }
            seen[item] = true
        }
    }

    for i := 0; i < 8; i++ {
        if _, ok := seen[i]; !ok {
            t.Errorf("Expected %d at least once.\n", i)
        }
    }
}

/*
 * Test to make sure you can choose nothing.
 */
func TestRandMultinomialNothing(t *testing.T) {
    res := RandMultinomial(0, 0)
    fmt.Printf("Example for (0, 0): %v\n", res)

    if !(len(res) == 1 && len(res[0]) == 0) {
        t.Errorf("Expected an empty response.\n")
    }
}

/*
 * Test to make sure you can choose everything.
 */
func TestRandMultinomialEverything(t *testing.T) {
    res := RandMultinomial(10, 10)
    fmt.Printf("Example for (10, 10): %v\n", res)

    if !(len(res) == 1 && len(res[0]) == 10) {
        t.Errorf("Expected a full response.\n")
    }
}
