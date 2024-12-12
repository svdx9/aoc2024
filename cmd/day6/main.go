package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	obstacleChar = '#'
	freeChar     = '.'
	visitedChar  = 'X'
	guardChar    = '^'
)

type direction int

const (
	up direction = 1 << iota
	left
	down
	right
)

type vec struct {
	x int
	y int
}

func (v vec) turnRight() vec {
	// {0,-1}, {1,0}, {0, 1}, {-1, 0}
	return vec{v.y * -1, v.x}
}

func (v vec) add(o vec) vec {
	return vec{v.x + o.x, v.y + o.y}
}

func (v vec) direction() direction {
	switch v {
	case vec{0, -1}:
		return up
	case vec{0, 1}:
		return down
	case vec{1, 0}:
		return right
	case vec{-1, 0}:
		return left
	}
	return up
}

func getFileContentsFromEnv(envKey string) (io.Reader, error) {
	fn := os.Getenv(envKey)
	if fn == "" {
		return nil, fmt.Errorf("expected value for env key %s", envKey)
	}
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)
	return r, nil
}

func newGrid(r io.Reader) (*grid, error) {
	scanner := bufio.NewScanner(r)
	g := &grid{}
	g.rows = make([][]rune, 0)
	for scanner.Scan() {
		t := scanner.Text()
		if g.sz.x == 0 {
			g.sz.x = len(t)
		}
		row := make([]rune, g.sz.x)
		g.rows = append(g.rows, row)
		for x, char := range t {
			switch char {
			case guardChar:
				g.startPos = vec{x: x, y: g.sz.y}
				g.rows[g.sz.y][x] = freeChar
			case freeChar:
				g.rows[g.sz.y][x] = char
			case obstacleChar:
				g.rows[g.sz.y][x] = char
			default:
				return nil, fmt.Errorf("unexpected char found at {%d,%d}: '%c'", x, g.sz.y, char)
			}
		}
		g.sz.y += 1
	}
	return g, nil
}

type grid struct {
	rows     [][]rune
	sz       vec
	startPos vec
}

func (g *grid) onGrid(pos vec) bool {
	return pos.x >= 0 && pos.x < g.sz.x && pos.y >= 0 && pos.y < g.sz.y
}

func (g *grid) isObs(pos vec) bool {
	return g.rows[pos.y][pos.x] == obstacleChar
}

type path struct {
	p map[vec]direction
}

func (p *path) add(pos vec, dir vec) {
	if d, ok := p.p[pos]; ok {
		p.p[pos] = d | dir.direction()
	} else {
		p.p[pos] = dir.direction()
	}
}

func (p *path) visited(pos vec, dir vec) bool {
	if d, ok := p.p[pos]; ok {
		return d&dir.direction() == dir.direction()
	}
	return false
}

func (p *path) len() int {
	return len(p.p)
}

func newPath() *path {
	return &path{
		p: make(map[vec]direction),
	}
}

func runSimulation(g *grid, p *path) bool {
	dir := vec{0, -1}
	pos := g.startPos
	// return false if looped
	// run a simulation starting at pos with direction dir
	for g.onGrid(pos) {
		// add to path
		p.add(pos, dir)
		// now check if we're blocked:
		next := pos.add(dir)
		if g.onGrid(next) && g.isObs(next) {
			dir = dir.turnRight()
		} else {
			pos = next
		}
		if p.visited(pos, dir) {
			// been here already, in a loop
			return false
		}
	}
	return true
}

func main() {
	r, err := getFileContentsFromEnv("INPUT")
	if err != nil {
		panic(err)
	}
	g, err := newGrid(r)
	if err != nil {
		panic(err)
	}
	p := newPath()
	if runSimulation(g, p) {
		fmt.Printf("stage1 moves: %d\n", p.len())
	}
	loop := 0
	for obstaclePos := range p.p {
		if obstaclePos.x == g.startPos.x && obstaclePos.y == g.startPos.y {
			// ignore start
			continue
		}
		g.rows[obstaclePos.y][obstaclePos.x] = obstacleChar
		if !runSimulation(g, newPath()) {
			// looped
			loop++
		}
		g.rows[obstaclePos.y][obstaclePos.x] = freeChar
	}
	fmt.Println(loop)
}
