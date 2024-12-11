package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

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

func (v vec) reflect() vec {
	return vec{
		x: v.x * -1,
		y: v.y * -1,
	}
}

var vecs = []vec{
	{1, 0},
	{1, 1},
	{0, 1},
	{-1, 1},
	{-1, 0},
	{-1, -1},
	{0, -1},
	{1, -1},
}

var cornerVecs = []vec{
	{1, 1},
	{-1, 1},
}

type runePeeker interface {
	runeAt(vec) (rune, error)
}

type wordSeeker interface {
	seek(r rune, pos vec)
}

func getFileContents() (*os.File, error) {
	filename := os.Getenv("INPUT")
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file '%s', '%w'", filename, err)
	}
	return f, nil
}

type Grid struct {
	rows [][]rune
	x    int
	y    int
}

func NewGrid(r io.Reader) (*Grid, error) {
	scanner := bufio.NewScanner(r)
	g := &Grid{}
	g.rows = make([][]rune, 0)
	for scanner.Scan() {
		t := scanner.Text()
		tLen := len(t)
		g.rows = append(g.rows, []rune(t))
		g.y += 1
		if g.x > 0 && g.x != tLen {
			return nil, fmt.Errorf("mismatch line length")
		}
		g.x = tLen
	}
	return g, nil
}

func (g *Grid) runeAt(pos vec) (rune, error) {
	if pos.x < 0 || pos.x >= g.x {
		return 0, fmt.Errorf("x %d out of bounds [%d]", pos.x, g.x)
	}
	if pos.y < 0 || pos.y >= g.y {
		return 0, fmt.Errorf("y %d out of bounds [%d]", pos.y, g.y)
	}
	return g.rows[pos.x][pos.y], nil
}

func (g *Grid) scan(seeker wordSeeker) error {
	for i := 0; i < g.x; i++ {
		for j := 0; j < g.y; j++ {
			r, err := g.runeAt(vec{i, j})
			if err != nil {
				return err
			}
			seeker.seek(r, vec{i, j})
		}
	}
	return nil
}

type wordFinderStrategy struct {
	word  string
	rp    runePeeker
	found int
	debug bool
}

func (w *wordFinderStrategy) seek(char rune, pos vec) {
	w.log(" found char: %c (%d,%d)\n", char, pos.x, pos.y)
	// get the first rune:
	wanted, width0 := utf8.DecodeRuneInString(w.word[0:])
	if char != wanted {
		return
	}
	for _, vec := range vecs {
		p := pos
		found := true
		for i, wi := width0, width0; i < len(w.word); i += wi {
			wanted, wi = utf8.DecodeRuneInString(w.word[i:])
			p = p.addVec(vec)
			if c, err := w.rp.runeAt(p); err != nil || wanted != c {
				w.log(">      char: %c[%c] (%d,%d)\n", wanted, c, p.x, p.y)
				found = false
				break
			} else {
				w.log("!      char: %c[%c] (%d,%d)\n", wanted, c, p.x, p.y)
			}
		}
		if found {
			w.found++
		}
	}
}

type xmasFinderStrategy struct {
	found int
	rp    runePeeker
}

func (w *xmasFinderStrategy) seek(char rune, pos vec) {
	// look for the 'A'
	if char != 'A' {
		return
	}
	found := true
	for _, vec := range cornerVecs {
		c1, c1err := w.rp.runeAt(pos.addVec(vec))
		c2, c2err := w.rp.runeAt(pos.addVec(vec.reflect()))
		if errors.Join(c1err, c2err) != nil {
			found = false
			break
		}
		if !((c1 == 'M' && c2 == 'S') || (c1 == 'S' && c2 == 'M')) {
			found = false
			break
		}
	}
	if found {
		w.found++
	}

}

func (w *wordFinderStrategy) log(format string, a ...any) {
	if w.debug {
		fmt.Printf(format, a...)
	}
}

func main() {
	f, err := getFileContents()
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	g, err := NewGrid(r)
	if err != nil {
		panic(err)
	}
	seeker := &wordFinderStrategy{
		word:  "XMAS",
		rp:    g,
		debug: true,
	}
	g.scan(seeker)
	fmt.Printf("found: %d\n", seeker.found)

	xmasSeeker := &xmasFinderStrategy{
		rp: g,
	}
	g.scan(xmasSeeker)
	fmt.Printf("found: %d\n", xmasSeeker.found)

}
