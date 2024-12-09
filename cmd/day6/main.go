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

const (
	up direction = 1 << iota
	down
	left
	right
)

type direction int

func (d direction) vec() vec {
	switch d {
	case up:
		return vec{0, -1}
	case right:
		return vec{1, 0}
	case down:
		return vec{0, 1}
	case left:
		return vec{-1, 0}
	}
	return vec{0, 0}
}

func (d direction) rotateRight() direction {
	switch d {
	case up:
		return right
	case right:
		return down
	case down:
		return left
	case left:
		return up
	}
	return up
}

func (d direction) char() rune {
	switch d {
	case up:
		return '|'
	case down:
		return '|'
	case left:
		return '-'
	case right:
		return '-'
	default:
		return '+'
	}
}

func getFileContents() (*os.File, error) {
	filename := os.Getenv("INPUT")
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file '%s', '%w'", filename, err)
	}
	return f, nil
}

type vec struct {
	x int
	y int
}

func (v vec) addVec(a vec) vec {
	return vec{
		x: v.x + a.x,
		y: v.y + a.y,
	}
}

type grid struct {
	rows  [][]rune
	sz    vec
	guard *guard
}

type guard struct {
	dir direction
	pos vec
	// path map[vec]direction
}

func (g *guard) move() {
	g.pos = g.pos.addVec(g.dir.vec())
}

func (g *guard) position() vec {
	return g.pos
}

func (g *guard) direction() direction {
	return g.dir
}

func (g *guard) rotateRight() {
	g.dir = g.dir.rotateRight()
}

func newGuard(dir direction, pos vec) *guard {
	return &guard{
		dir: dir,
		pos: pos,
	}
}

type moveable interface {
	move()
	direction() direction
	position() vec
	rotateRight()
}

//	func (g *guard) marker() rune {
//		if d, ok := g.path[g.pos]; ok {
//			return d.char()
//		}
//		return '0'
//	}
// func (g *guard) wouldLoop() bool {
// 	loopDir := g.dir.rotateRight()
// 	// keep walking until either we hit an obstacle (no loop)
// 	// or we get onto a previous path, and we are going inthe
// 	// same direction
// 	loopPos := g.pos.addVec(loopDir.vec())

// 	// 	if cur, ok := g.path[loopPos]; ok {
// 	// 		fmt.Printf("check curPos: %+v %d l: %+v:%d: %d\n", g.pos, cur, loopPos, loopDir, cur&loopDir)
// 	// 		return cur&loopDir == loopDir
// 	// 	}
// 	// 	return false
// 	// }
// }

type path struct {
	p map[vec]direction
}

func (p path) update(pos vec, dir direction) {
	if cur, ok := p.p[pos]; ok {
		p.p[pos] = dir & cur
	} else {
		p.p[pos] = dir
	}
}

func (p path) len() int {
	return len(p.p)
}

func (p path) getPos(pos vec) (direction, bool) {
	d, ok := p.p[pos]
	return d, ok
}

func newPath() *path {
	return &path{p: make(map[vec]direction)}
}

func NewGrid(r io.Reader) (*grid, error) {
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
			if char == guardChar {
				g.guard = newGuard(up, vec{x, g.sz.y})
			} else if char == guardChar || char == obstacleChar || char == freeChar {
				g.rows[g.sz.y][x] = char
			} else {
				return nil, fmt.Errorf("unexpected char found at {%d,%d}: '%c'", x, g.sz.y, char)
			}
		}
		g.sz.y += 1
	}
	g.dump(nil)
	return g, nil
}

func (g *grid) dump(p *path) {
	for j, r := range g.rows {
		for i, c := range r {
			if p != nil {
				if d, ok := p.getPos(vec{i, j}); ok {
					fmt.Printf("%c", d.char())
				} else {
					fmt.Printf("%c", c)
				}
			} else {
				fmt.Printf("%c", c)
			}
		}
		fmt.Println()
	}
}

func (g *grid) offMap(pos vec) bool {
	return pos.x < 0 || pos.x >= g.sz.x || pos.y < 0 || pos.y >= g.sz.y
}

func (g *grid) posFree(pos vec) bool {
	return g.rows[pos.y][pos.x] != obstacleChar
}

func (g *grid) move(m moveable) (vec, error) {
	d := m.direction()
	pos := m.position().addVec(d.vec())
	if g.offMap(pos) {
		return vec{}, fmt.Errorf("off map")
	}
	if g.posFree(pos) {
		// we can move
		m.move()
	} else {
		// rotate
		m.rotateRight()
	}
	return m.position(), nil
}

func main() {
	f, err := getFileContents()
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	g, err := NewGrid(r)

	g.dump(nil)

	if err != nil {
		panic(err)
	}

	path := newPath()
	fmt.Printf("guard is at {%d,%d}\n", g.guard.pos.x, g.guard.pos.y)
	for {
		// check if off map
		path.update(g.guard.position(), g.guard.direction())
		if _, err := g.move(g.guard); err == nil {

		} else {
			fmt.Println("guard off map")
			break
		}
	}
	g.dump(path)
	fmt.Printf("moved %d times", path.len())
}
