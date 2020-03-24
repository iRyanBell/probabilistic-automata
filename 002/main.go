package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

const (
	WIDTH  = 640
	HEIGHT = 360
	SCALE  = 2
	RULE   = 54
)

func ProbabilisticBit(truth bool, thresh float32) bool {
	if truth {
		return rand.Float32() < thresh
	}
	return rand.Float32() < (1 - thresh)
}

var grid [HEIGHT][WIDTH]bool

func ruleNumToBits(r uint8) [7]bool {
	return [7]bool{
		r&64 == 64,
		r&32 == 32,
		r&16 == 16,
		r&8 == 8,
		r&4 == 4,
		r&2 == 2,
		r&1 == 1,
	}
}

func applyRule(row [WIDTH]bool) [WIDTH]bool {
	rule := ruleNumToBits(RULE)

	rowNext := [WIDTH]bool{}

	for i := 0; i < len(row); i++ {
		var l, r bool
		// Wrap
		if i == 0 {
			l = row[WIDTH-1]
		} else {
			l = row[i-1]
		}

		if i == WIDTH-1 {
			r = row[0]
		} else {
			r = row[i+1]
		}

		c := row[i]

		switch [3]bool{l, c, r} {
		case [3]bool{true, true, true}:
			rowNext[i] = ProbabilisticBit(true, 0.5)
		case [3]bool{true, true, false}:
			rowNext[i] = ProbabilisticBit(rule[0], 1.0)
		case [3]bool{true, false, true}:
			rowNext[i] = ProbabilisticBit(rule[1], 1.0)
		case [3]bool{true, false, false}:
			rowNext[i] = ProbabilisticBit(rule[2], 1.0)
		case [3]bool{false, true, true}:
			rowNext[i] = ProbabilisticBit(rule[3], 1.0)
		case [3]bool{false, true, false}:
			rowNext[i] = ProbabilisticBit(rule[4], 1.0)
		case [3]bool{false, false, true}:
			rowNext[i] = ProbabilisticBit(rule[5], 1.0)
		case [3]bool{false, false, false}:
			rowNext[i] = ProbabilisticBit(rule[6], 1.0)
		}
	}

	return rowNext
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			if grid[y][x] {
				screen.Set(x, y, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
			}
		}

		if y < HEIGHT-1 {
			// Scroll
			grid[y] = grid[y+1]
		}
	}
	grid[HEIGHT-1] = applyRule(grid[HEIGHT-1])

	return nil
}

func initialize() {
	// Flip center bit on top row
	grid[HEIGHT-1][int(WIDTH/2)] = true
}

func main() {
	initialize()

	err := ebiten.Run(update, WIDTH, HEIGHT, SCALE, "PA002")
	if err != nil {
		fmt.Println(err)
	}
}
