// A Probabilistic Take On Conway's Game Of Life

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
)

// ProbabilisticBit returns a boolean value from a random float threshold.
func ProbabilisticBit(thresh float32) bool {
	return rand.Float32() < thresh
}

// Grid: current
var grid [WIDTH][HEIGHT]bool

// Grid: previous
var gridP [WIDTH][HEIGHT]bool

func drawScreen(screen *ebiten.Image) {
	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			if grid[x][y] {
				screen.Set(x, y, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
			}
		}
	}
}

func step() {
	var nextGrid [WIDTH][HEIGHT]bool

	for x := 1; x < WIDTH-1; x++ {
		for y := 1; y < HEIGHT-1; y++ {
			// Check center value
			c := grid[x][y]

			// Check prev center value
			cp := gridP[x][y]

			// Check neighboring cells
			t := grid[x][y-1]
			tr := grid[x+1][y-1]
			r := grid[x+1][y]
			br := grid[x+1][y+1]
			b := grid[x][y+1]
			bl := grid[x-1][y+1]
			l := grid[x-1][y]
			tl := grid[x-1][y-1]

			// Count the sum
			neighbors := [9]bool{cp, t, r, b, l, tr, br, bl, tl}
			var sum int
			for idx, neighbor := range neighbors {
				if neighbor {
					if idx < 5 {
						sum += 3
					} else {
						sum += 2
					}
				}
			}

			// Apply rules
			if c {
				nextGrid[x][y] = ProbabilisticBit(0.025)
			} else if sum == 4 || sum == 5 {
				nextGrid[x][y] = ProbabilisticBit(0.975)
			}
		}
	}

	// Apply step
	gridP = grid
	grid = nextGrid
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	step()
	drawScreen(screen)

	return nil
}

func Init() {
	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			grid[x][y] = rand.Intn(2) == 0
		}
	}
	gridP = grid
}

func main() {
	Init()

	err := ebiten.Run(update, WIDTH, HEIGHT, SCALE, "PA001")
	if err != nil {
		fmt.Println(err)
	}
}
