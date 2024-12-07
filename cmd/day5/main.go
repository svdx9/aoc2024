package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getFileContents() (*os.File, error) {
	filename := os.Getenv("INPUT")
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file '%s', '%w'", filename, err)
	}
	return f, nil
}

type rule struct {
	p1 int
	p2 int
}
type update []int
type pageOrderMap map[int]map[int]struct{}

type ruleEvaluator struct {
	before pageOrderMap
	after  pageOrderMap
}

func newRuleEvaluator() *ruleEvaluator {
	r := &ruleEvaluator{
		before: make(map[int]map[int]struct{}),
		after:  make(map[int]map[int]struct{}),
	}
	return r
}
func (re *ruleEvaluator) addRule(rule rule) {
	update := func(m pageOrderMap, key int, val int) {
		if p, ok := m[key]; ok {
			p[val] = struct{}{}
		} else {
			m[key] = make(map[int]struct{})
			m[key][val] = struct{}{}
		}
	}
	// p2 comes after p1
	update(re.after, rule.p1, rule.p2)
	// p1 must be before p2
	update(re.before, rule.p2, rule.p1)
}

func (re *ruleEvaluator) allBefore(page int, pages []int) error {
	if pagesBefore, ok := re.before[page]; ok {
		// fmt.Printf("BB: check %+v are all in %+v\n", pages, pagesBefore)
		for _, p := range pages {
			if _, ok := pagesBefore[p]; !ok {
				return fmt.Errorf("page %d not before %d", p, page)
			}
		}
	}
	return nil
}

func (re *ruleEvaluator) allAfter(page int, pages []int) error {
	if pagesAfter, ok := re.after[page]; ok {
		// fmt.Printf("AA: check %+v are all in %+v\n", pages, pagesAfter)
		for _, p := range pages {
			if _, ok := pagesAfter[p]; !ok {
				return fmt.Errorf("page %d not after %d", p, page)
			}
		}
	}
	return nil
}

func (re *ruleEvaluator) checkOrder(u update) error {
	for idx, page := range u {
		// fmt.Printf("check %+v\n", u)
		pagesBefore := u[0:idx]
		// check that all the pagesBefore come before this page
		err := re.allBefore(page, pagesBefore)
		if err != nil {
			return err
		}
		pagesAfter := u[idx+1:]
		err = re.allAfter(page, pagesAfter)
		if err != nil {
			return err
		}
	}
	return nil
}

func (re *ruleEvaluator) isBefore(a, b int) bool {
	if m, ok := re.before[a]; ok {
		if _, ok := m[b]; ok {
			return true
		}
	}
	return false
}

func (re *ruleEvaluator) reOrder(u update) update {
	// new := make(update, len(u))
	ordered := false
	// keep going until they are all ordered:
	idx := 0
	for !ordered {
		// if all the numbers after this one are not present in the
		swap := false
		for i := 1; i < len(u); i++ {
			if !re.isBefore(u[i], u[i-1]) {
				u[i], u[i-1] = u[i-1], u[i]
				// fmt.Printf("swap %d %d %+v\n", i, i-1, u)
				swap = true
			}
		}
		ordered = !swap
		idx++
	}
	return u
}

func newRule(ruleString string) (rule, error) {
	r := strings.Split(ruleString, "|")
	if len(r) != 2 {
		return rule{}, fmt.Errorf("need 2 elements in rule")
	}
	p1, err := strconv.Atoi(r[0])
	if err != nil {
		return rule{}, err
	}
	p2, err := strconv.Atoi(r[1])
	if err != nil {
		return rule{}, err
	}
	return rule{p1: p1, p2: p2}, nil
}

func newUpdate(u string) (update, error) {
	rv := update{}
	for _, s := range strings.Split(u, ",") {
		i, err := strconv.Atoi(s)
		if err != nil {
			return update{}, err
		}
		rv = append(rv, i)
	}
	return rv, nil
}

func main() {
	f, err := getFileContents()
	if err != nil {
		panic(err)
	}
	defer f.Close()
	doRules := true
	// updates := make([]update, 0)
	part1Count := 0
	part2Count := 0
	re := newRuleEvaluator()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			doRules = false
			continue
		}
		if doRules {
			r, err := newRule(line)
			if err != nil {
				panic(err)
			}
			re.addRule(r)
		} else {
			fmt.Println("--")
			u, err := newUpdate(line)
			if err != nil {
				panic(err)
			}
			err = re.checkOrder(u)
			if err != nil {
				fmt.Printf("%s\n", err)
				r := re.reOrder(u)
				part2Count += r[len(r)/2]
			} else {
				fmt.Printf("%+v %d %d\n", u, len(u), u[len(u)/2])
				part1Count += u[len(u)/2]
			}
			// updates = append(updates, u)
		}
	}
	fmt.Println(part1Count, part2Count)
}
