package internal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type InputLineHandler interface {
	HandleInput([]int) error
}

func parseInputString(text string) ([]int, error) {
	fields := strings.Split(text, " ")
	rv := make([]int, 0)
	for _, field := range fields {
		if len(field) == 0 {
			continue
		}
		i, err := strconv.Atoi(field)
		if err != nil {
			return nil, err
		}
		rv = append(rv, i)
	}
	return rv, nil
}

func ParseInputFile(filename string, handler InputLineHandler) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open file '%s', '%w'", filename, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for {
		ok := scanner.Scan()
		if ok {
			// fmt.Println(scanner.Text())
			values, err := parseInputString(scanner.Text())
			if err != nil {
				return err
			}
			err = handler.HandleInput(values)
			if err != nil {
				return err
			}
		} else {
			break
		}
	}
	return nil

}

func GetFileContents() (*os.File, error) {
	filename := os.Getenv("INPUT")
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file '%s', '%w'", filename, err)
	}
	return f, nil
}
