package main

import (
	"flag"
	"fmt"
	"slices"

	"github.com/svdx9/aoc2024/internal"
)

type DistanceCalculator struct {
	a []int
	b []int
}

func (d *DistanceCalculator) HandleInput(input []int) error {
	if len(input) != 2 {
		return fmt.Errorf("unexpected number of fields: %d", len(input))
	}
	d.a = append(d.a, input[0])
	d.b = append(d.b, input[1])
	return nil
}

func (d *DistanceCalculator) similarity() (int, error) {
	// build a map of the counts
	bCounts := make(map[int]int)
	for _, bVal := range d.b {
		if cnt, ok := bCounts[bVal]; ok {
			bCounts[bVal] = cnt + 1
		} else {
			bCounts[bVal] = 1
		}
	}
	for k, v := range bCounts {
		fmt.Printf("%d counted %d times\n", k, v)
	}
	similarity := 0
	for _, aVal := range d.a {
		if bCnt, ok := bCounts[aVal]; ok {
			similarity += aVal * bCnt
		}
	}
	return similarity, nil
}

func (d *DistanceCalculator) calculate() (int, error) {
	if len(d.a) != len(d.b) {
		return 0, fmt.Errorf("mismatched slice length for a and b")
	}

	a := make([]int, len(d.a))
	copy(a, d.a)
	slices.Sort(a)

	b := make([]int, len(d.b))
	copy(b, d.b)
	slices.Sort(b)

	distance := 0
	for i, aVal := range a {
		delta := b[i] - aVal
		if delta < 0 {
			delta = delta * -1
		}
		distance += delta
	}
	return distance, nil
}

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "input", "", "specify file to read as input")
	flag.Parse()

	d := &DistanceCalculator{}
	err := internal.ParseInputFile(inputFile, d)
	if err != nil {
		panic(err)
	}
	distance, err := d.calculate()
	if err != nil {
		panic(err)
	}
	fmt.Printf("distance is %d\n", distance)

	similarity, err := d.similarity()
	if err != nil {
		panic(err)
	}
	fmt.Printf("similarity is %d\n", similarity)

}
