package main

import (
	"fmt"

	"github.com/svdx9/aoc2024/internal"
)

type antinodeScanner interface {
	scan(a, b internal.Vec) []internal.Vec
}

type antennaGrid struct {
	grid     *internal.Grid
	antennas map[rune][]internal.Vec
}

func (a *antennaGrid) scan() []internal.Vec {
	rv := make(map[internal.Vec]struct{}, 0)
	for freq, positions := range a.antennas {
		for i := 0; i < len(positions); i++ {
			for j := i + 1; j < len(positions); j++ {
				fmt.Printf("freq: %s compare: %d %d\n", string(freq), i, j)
				a1 := positions[i]
				a2 := positions[j]
				delta := a2.SubVec(a1)
				for _, antinode := range []internal.Vec{a2.AddVec(delta), a1.SubVec(delta)} {
					fmt.Printf("freq: %s compare: %+v %+v %+v %+v\n", string(freq), a1, a2, delta, antinode)
					if a.grid.OutOfBounds(antinode) {
						continue
					}
					rv[antinode] = struct{}{}
				}
			}
		}
	}
	r := make([]internal.Vec, 0)
	for antinode, _ := range rv {
		r = append(r, antinode)
	}
	return r
}

func newAntennaGrid(g *internal.Grid) (*antennaGrid, error) {
	a := antennaGrid{
		grid:     g,
		antennas: make(map[rune][]internal.Vec),
	}
	for j := 0; j < g.Size.Y; j++ {
		for i := 0; i < g.Size.X; i++ {
			pos := internal.Vec{X: i, Y: j}
			c, err := g.RuneAt(pos)
			if c == '.' {
				continue
			}
			if err != nil {
				return nil, err
			}
			if _, ok := a.antennas[c]; !ok {
				a.antennas[c] = make([]internal.Vec, 0)
			}
			a.antennas[c] = append(a.antennas[c], pos)
		}
	}
	return &a, nil
}

func main() {
	r, err := internal.GetFileContents()
	if err != nil {
		panic(err)
	}
	g, err := internal.NewGrid(r)
	if err != nil {
		panic(err)
	}
	a, err := newAntennaGrid(g)
	if err != nil {
		panic(err)
	}
	antinodes := a.scan()
	fmt.Printf("%+v %d\n", antinodes, len(antinodes))

}
