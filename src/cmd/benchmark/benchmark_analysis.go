package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
)


/*
 * Analyze the results of the benchmarking. The benchmarking results file has
 * the initial setup of a Euchre hand and then the difference between the
 * specified player type and the optimal, all knowing player. This script
 * calculates the average, and the distribution of the differences amongst the
 * values 0 to 6 inclusive.
 *
 * Usage:
 *  ./benchmark_analysis -dataLoc={dataLoc}
 */


func main() {
    var dataLoc string
    flag.StringVar(&dataLoc, "dataLoc", "", "Location of the benchmark results.")
    flag.Parse()

    dataFile, err := os.Open(dataLoc)
    if err != nil {
        log.Fatal(err)
    }
    defer dataFile.Close()

    count := 0
    sum := 0.0
    bins := make([]float64, 7)
    scanner := bufio.NewScanner(dataFile)
    for scanner.Scan() {
        line := scanner.Text()
        tabIndex := strings.IndexRune(line, '\t')

        diff, _ := strconv.ParseFloat(line[tabIndex + 1:], 64)
        bin := int(diff)

        count++
        sum += diff
        bins[bin]++
    }

    for i := 0; i < len(bins); i++ {
        bins[i] = bins[i] / float64(count)
    }

    avg := sum / float64(count)
    fmt.Printf("%f, %f, %f, %f, %f, %f, %f, %f\n", avg, bins[0], bins[1], bins[2],
                                                   bins[3], bins[4], bins[5],
                                                   bins[6])
}
