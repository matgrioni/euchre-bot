package euchre


/*
 * Intersect two sets that represent player indices. Keys that are mapped to a
 * true value in both sets are kept in the result.
 *
 * Args:
 *  s1: The first set to consider.
 *  s2: The second set to consider in the intersection.
 *
 * Returns:
 *  The intersection of the two sets. Elements that were in both sets and were
 *  both mapped to a true value.
 */
func intersectPlayerSets(s1 map[int]bool, s2 map[int]bool) map[int]bool {
    result := make(map[int]bool)
    for k, _ := range s1 {
        if pres, ok := s2[k]; ok && pres {
            result[k] = true
        }
    }

    return result
}


/*
 * Return a random player index from a set representing some player indices.
 * Only players that correspond to a true value will be considered. Stuff like
 * this is why go is sometimes a huge PITA.
 *
 * Args:
 *  s: The set to randomly sample from.
 *
 * Returns:
 *  A random player index from the set from those that are mapped to true.
 */
func randomPlayerFromSet(s map[int]bool) int {
    keys := make([]int, 0, len(s))
    for k, _ := range s {
        keys = append(keys, k)
    }

    return keys[r.Intn(len(keys))]
}
