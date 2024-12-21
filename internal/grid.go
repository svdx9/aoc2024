package internal

import (
	"bufio"
	"fmt"
	"io"
)

type Grid struct {
	rows [][]rune
	Size Vec
}

func NewGrid(r io.Reader) (*Grid, error) {
	scanner := bufio.NewScanner(r)
	g := &Grid{}
	g.rows = make([][]rune, 0)
	for scanner.Scan() {
		t := scanner.Text()
		tLen := len(t)
		if tLen == 0 {
			continue
		}
		g.rows = append(g.rows, []rune(t))
		g.Size.Y += 1
		if g.Size.X > 0 && g.Size.X != tLen {
			return nil, fmt.Errorf("mismatch line length on line %d", g.Size.Y)
		}
		g.Size.X = tLen
	}
	return g, nil
}

func (g Grid) RuneAt(pos Vec) (rune, error) {

	if g.OutOfBounds(pos) {
		return 0, fmt.Errorf("%v out of bounds %v", pos, g.Size)
	}
	fmt.Printf("+%v +%v\n", pos, g.Size)
	return g.rows[pos.Y][pos.X], nil
}

func (g Grid) OutOfBounds(pos Vec) bool {
	return pos.X < 0 || pos.X >= g.Size.X || pos.Y < 0 || pos.Y >= g.Size.Y
}
