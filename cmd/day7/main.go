package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type puzzle struct {
	orig   string
	answer int
	values []int
}

func (p puzzle) sum() int {
	total := 0
	for _, v := range p.values {
		total += v
	}
	return total
}

func solve(answer int, values []int, operator string) (string, error) {
	fmt.Printf("TRY: %d %+v %s\n", answer, values, operator)
	// start
	switch len(values) {
	case 0:
		// should never get here
		if answer == 0 {
			return operator, nil
		}
		return operator, fmt.Errorf("cannot solve, no values non-zero answer")
	case 1:
		// if the answer equals the last value, means we have a correct
		// solution
		fmt.Printf("case1: %d %+v %s\n", answer, values, operator)
		if answer == values[0] {
			return operator, nil
		}
		return operator, fmt.Errorf("cannot solve, no values non-zero answer")
	default:
		// do another calculation
	}
	// pop last value
	last, values := values[len(values)-1], values[:len(values)-1]

	// try multiply
	if answer%last == 0 {
		fmt.Printf("%d %d\n", answer, last)
		if operator, err := solve(answer/last, values, operator+fmt.Sprintf("%d*", last)); err == nil {
			// this worked
			return operator, err
		}
	}
	// now try addition
	if operator, err := solve(answer-last, values, operator+fmt.Sprintf("%d+", last)); err == nil {
		// this worked
		return operator, err
	}

	return operator, fmt.Errorf("eh?")
}
func (p *puzzle) solve() bool {
	// test for edge case on start
	if len(p.values) == 0 {
		// no values
		return false
	}
	if len(p.values) == 1 {
		// simples, just check the anwer matches the value
		return p.values[0] == p.answer
	}
	// start recursive
	// copy the values
	values := make([]int, len(p.values))
	copy(values, p.values)

	operations, err := solve(p.answer, values, "")
	if err != nil {
		fmt.Printf("error:%s\n", err)
		return false
	}
	fmt.Printf("X; res: %d %s %+v\n", p.answer, operations, p.values)
	return true
}

// func (p *puzzle) old() (int, bool) {
// 	if len(p.values) < 2 {
// 		panic("not enough values")
// 	}
// 	operators := ""
// 	for i := len(p.values) - 1; i > 0; i-- {
// 		// try multiply first
// 		// fmt.Printf("%d %d\n", p.values[i-1], p.values[i])
// 		if p.answer%p.values[i] == 0 {
// 			operators = operators + "*"
// 			p.answer = p.answer / p.values[i]
// 		} else {
// 			operators = operators + "+"
// 			p.answer = p.answer - p.values[i]
// 		}
// 		if p.answer < 0 {
// 			fmt.Printf("<0: %s:%d:%d: %s\n", p.orig, p.answer, i, operators)
// 			return 0, false
// 		}
// 	}
// 	if p.answer != p.values[0] {
// 		// not solveable
// 		return 0, false
// 	}
// 	total := 0

// 	for idx, oper := range operators {
// 		if idx == 0 {
// 			if oper == '*' {
// 				total = p.values[0] * p.values[1]
// 			} else {
// 				total = p.values[0] + p.values[1]
// 			}
// 		} else {
// 			if oper == '*' {
// 				total *= p.values[idx+1]
// 			} else {
// 				total += p.values[idx+1]
// 			}
// 		}
// 		// fmt.Printf("%c:%d:%d:%d\n", oper, a, b, total)

// 	}
// 	return total, true

// }

func getFileContents() (*os.File, error) {
	filename := os.Getenv("INPUT")
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file '%s', '%w'", filename, err)
	}
	return f, nil
}

func getPuzzle() ([]puzzle, error) {
	c, err := getFileContents()
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(bufio.NewReader(c))
	puzzles := make([]puzzle, 0)
	for s.Scan() {
		puzzle := puzzle{
			values: make([]int, 0),
			orig:   s.Text(),
		}
		split1 := strings.Split(s.Text(), ":")
		if len(split1) != 2 {
			return nil, fmt.Errorf("expected one colon in %s", s.Text())
		}

		answer, rest := split1[0], strings.Trim(split1[1], " ")
		answerInt, err := strconv.Atoi(answer)
		if err != nil {
			return nil, err
		}
		puzzle.answer = answerInt

		for _, val := range strings.Split(strings.Trim(rest, " "), " ") {
			valInt, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			puzzle.values = append(puzzle.values, valInt)
		}

		puzzles = append(puzzles, puzzle)
	}
	return puzzles, nil
}

func main() {
	puzzles, err := getPuzzle()
	if err != nil {
		panic(err)
	}
	total := 0
	for idx, puzzle := range puzzles {
		ok := puzzle.solve()
		if ok {
			total += puzzle.answer
			fmt.Printf("X;can solve %d: answer %d\n", idx, puzzle.answer)
		} else {
			fmt.Printf("X;cannot solve %s\n", puzzle.orig)
		}
	}
	fmt.Printf("%d\n", total)

}
