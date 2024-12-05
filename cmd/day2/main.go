package main

import (
	"flag"
	"fmt"

	"github.com/svdx9/aoc2024/internal"
)

type Report struct {
	values         []int
	errorTolerance int // maximum bad level count
	levelTolerance int // maximum gap between levels allowed
}

type minMaxReportTester struct {
	minLevelDeltaTolerance int
	maxLevelDeltaTolerance int
}

func (m minMaxReportTester) test(delta int, increasing bool) bool {
	if ((delta < 0) && increasing) || ((delta > 0) && !increasing) {
		return false
	}
	if delta < 0 {
		delta = delta * -1
	}
	if delta < m.minLevelDeltaTolerance || delta > m.maxLevelDeltaTolerance {
		// must always diff by at least the min level
		return false
	}
	return true

}

type reportTester interface {
	test(int, bool) bool
}

func isSafe(reportValues []int, rt reportTester) bool {
	// check the report and count unsafe values in the series
	if len(reportValues) < 2 {
		panic("not enough values")
	}
	// test to see if this is an increasing or decreasing range
	increasing := reportValues[0] < reportValues[1]
	for curIdx := range reportValues {
		if curIdx == 0 {
			continue
		}
		delta := reportValues[curIdx] - reportValues[curIdx-1]
		if !rt.test(delta, increasing) {
			return false
		}
	}
	return true
}

type ReportCalculator struct {
	countSafe         int
	badLevelTolerance int
	reportTester      reportTester
}

func newSlice(s []int, without int) []int {
	rv := make([]int, len(s)-1)
	newIdx := 0
	for idx, v := range s {
		if idx == without {
			continue
		}
		rv[newIdx] = v
		newIdx++
	}
	return rv
}

func (rc *ReportCalculator) HandleInput(values []int) error {
	for i := 0; i < rc.badLevelTolerance; i++ {
		if isSafe(values, rc.reportTester) {
			rc.countSafe++
			return nil
		} else {
			// try again with a level removed
			for idx := 0; idx < len(values); idx++ {
				smaller := newSlice(values, idx)
				if isSafe(smaller, rc.reportTester) {
					rc.countSafe++
					return nil
				}
			}
		}
	}
	fmt.Printf("%+v %d\n", values, rc.countSafe)
	return nil
}

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "input", "", "specify file to read as input")
	flag.Parse()

	r := &ReportCalculator{
		badLevelTolerance: 1,
		reportTester: &minMaxReportTester{
			minLevelDeltaTolerance: 1,
			maxLevelDeltaTolerance: 3,
		},
	}
	err := internal.ParseInputFile(inputFile, r)
	if err != nil {
		panic(err)
	}
	fmt.Println(r.countSafe)
}
