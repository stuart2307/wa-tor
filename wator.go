package main

import (
    "image/color"
    "log"
    "math"

    "github.com/hajimehoshi/ebiten"
)
//Water = 0, Fish = 1, Shark = 2
type Square struct {
	occupied int
	occupant int
	energy int
}

const scale int = 2
var NumShark = 10
var NumFish = 10
var FishBreed = 10
var SharkBreed = 10
var Starve = 10
const width = 300
const height = 300
var Threads = 1

var blue color.Color = color.RGBA{69, 145, 196, 255}
var yellow color.Color = color.RGBA{255, 230, 120, 255}
var grid [width][height]Square = [width][height]Square{}
var buffer [width][height]Square = [width][height]Square{}
var count int = 0

var FishBreedTimer = FishBreed
var SharkBreedTimer = SharkBreed

func update() error {
    for x := 1; x < width-1; x++ {
        for y := 1; y < height-1; y++ {
            buffer[x][y].occupant = 0
            n := (grid[x][y+1].occupied << 3) + (grid[x+1][y].occupied << 2) + (grid[x][y-1].occupied << 1) + grid[x-1][y].occupied

			//12 is the highest 4 bit number with at least 2 bits set to 0
			if (n == 16) {
				buffer[x-1][y].occupant = 1
				buffer[x-1][y].occupied = 1
			}
			else {
				choice = math.pow(2, math.rand.Intn(4))
				if (n & choice == 0)
				{
					switch choice {
					case 1:
						buffer[x-1][y].occupant = 1
						buffer[x-1][y].occupied = 1
						break
					case 2:
						buffer[x][y-1].occupant = 1
						buffer[x][y-1].occupied = 1
						break
					case 4:
						buffer[x+1][y].occupant = 1
						buffer[x+1][y].occupied = 1
						break
					case 8: 
						buffer[x][y+1].occupant = 1
						buffer[x][y+1].occupied = 1
						break
					}
				}
			}
        }
    }

    temp := buffer
    buffer = grid
    grid = temp
    return nil
}

func display(window *ebiten.Image) {
    window.Fill(blue)

    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            for i := 0; i < scale; i++ {
                for j := 0; j < scale; j++ {
                    if grid[x][y].occupant == 1 {
                        window.Set(x*scale+i, y*scale+j, yellow)
                    }
                }
            }
        }
    }
}

func frame(window *ebiten.Image) error {
    count++
    var err error = nil
    if count == 20 {
        err = update()
        count = 0
    }
    if !ebiten.IsDrawingSkipped() {
        display(window)
    }

    return err
}

func main() {
    for x := 1; x < width-1; x++ {
        for y := 1; y < height-1; y++ {
            if math.rand.Float32() < 0.5 {
                grid[x][y].occupant = 1
            }
        }
    }

    if err := ebiten.Run(frame, width, height, 2, "Game of Life"); err != nil {
        log.Fatal(err)
    }
}