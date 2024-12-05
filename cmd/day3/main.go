package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

var re = regexp.MustCompile(`(don't\(\)|do\(\)|mul\((\d+),(\d+)\))`)

func getFileContents() (*os.File, error) {
	filename := os.Getenv("INPUT")
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file '%s', '%w'", filename, err)
	}
	return f, nil
}

func mul(x, y string) (int, error) {
	xInt, err := strconv.Atoi(x)
	if err != nil {
		return 0, err
	}
	yInt, err := strconv.Atoi(y)
	if err != nil {
		return 0, err
	}
	return xInt * yInt, nil

}

func getMul(s string) (int, error) {
	total := 0
	b := []byte(s)
	do := true
	for _, loc := range re.FindAllSubmatch(b, -1) {
		cmd := string(loc[1])
		// fmt.Printf("%+v!%+v!%+v\n", string(loc[1]), string(loc[2]), string(loc[3]))
		switch cmd {
		case "do()":
			do = true
		case "don't()":
			do = false
		default:
			// must be mul
			if cmd[:3] != "mul" {
				return 0, fmt.Errorf("bad command: %s", cmd)
			}
			if do {
				fmt.Printf("%s:%t\n", cmd, do)
				m, err := mul(string(loc[2]), string(loc[3]))
				if err != nil {
					return 0, err
				}
				total += m
			}
		}
	}
	return total, nil
}

func main() {
	f, err := getFileContents()
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	b, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	t, err := getMul(string(b))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", t)
}
