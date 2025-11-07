package main

import (
    "image/color"
    "log"
    "math/rand"

    "github.com/hajimehoshi/ebiten"
)
//Water = 0, Fish = 1, Shark = 2
type Square struct {
	occupied int
	occupant int
	energy int
    breed int
}

const scale int = 2
var NumShark = 10
var NumFish = 10
var FishBreed = 5
var SharkBreed = 10
var Starve = 10
var SharkEnergyRestore = 5
const width = 300
const height = 300
var Threads = 1

var blue color.Color = color.RGBA{69, 145, 196, 255}
var yellow color.Color = color.RGBA{255, 230, 120, 255}
var grid [width][height]Square = [width][height]Square{}
var buffer [width][height]Square = [width][height]Square{}
var count int = 0

func moveFish() {
    for x := 1; x < width-1; x++ {
        for y := 1; y < height-1; y++ {
            if (grid[x][y].occupied == 1 && grid[x][y].occupant == 1) {
                n := (grid[x][y+1].occupied << 3) + (grid[x+1][y].occupied << 2) + (grid[x][y-1].occupied << 1) + grid[x-1][y].occupied
                if (n == 15) {
                    buffer[x][y] = grid[x][y]
                } else {
                    if (grid[x][y].breed == FishBreed) {
                        grid[x][y].breed = -1
                        buffer[x][y] = grid[x][y]
                    } else {
                        buffer[x][y].occupant = 0
                        buffer[x][y].occupied = 0
                    }
                    var choice = 1 << rand.Intn(4)
                    for (choice != -1) {
                        if (n & choice == 0) {
                            switch choice {
                            case 1:
                                buffer[x-1][y] = grid[x][y]
                                buffer[x-1][y].breed++;
                                choice = -1
                                break
                            case 2:
                                buffer[x][y-1] = grid[x][y]
                                buffer[x][y-1].breed++;
                                choice = -1
                                break
                            case 4:
                                buffer[x+1][y] = grid[x][y]
                                buffer[x+1][y].breed++;
                                choice = -1
                                break
                            case 8: 
                                buffer[x][y+1] = grid[x][y]
                                buffer[x][y+1].breed++;
                                choice = -1
                                break
                            }
                        }
                        if (choice != -1) {
                            choice = choice << 1
                            if (choice > 8) {
                                choice = 1
                            }
                        }
                    }
                }
            }
        }
    }
}

/* func moveSharks() {
    for x := 1; x < width-1; x++ {
        for y := 1; y < height-1; y++ {
            if (grid[x][y].occupied == 1 && grid[x][y].occupant == 2) {
                n := (grid[x][y+1].occupant << 3) + (grid[x+1][y].occupant << 2) + (grid[x][y-1].occupant << 1) + grid[x-1][y].occupant
                if (n == 30) {
                    buffer[x][y] = grid[x][y]
                } else {
                    if (grid[x][y].breed == SharkBreed) {
                        grid[x][y].breed = -1
                        buffer[x][y] = grid[x][y]
                    } else {
                        buffer[x][y].occupant = 0
                        buffer[x][y].occupied = 0
                    }
                    n := (grid[x][y+1].occupant%2 << 3) + (grid[x+1][y].occupant%2 << 2) + (grid[x][y-1].occupant%2 << 1) + grid[x-1][y].occupant%2
                    if (n > 0)
                }
            }
        }
    }
} */

func update() error {
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            buffer[x][y] = Square{}
        }
    }
    moveFish()
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
                    if grid[x][y].occupied == 1 {
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
    if count == 1 {
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
            if rand.Float32() < 0.1 {
                grid[x][y].occupant = 1
                grid[x][y].occupied = 1
            }
        }
    }

    if err := ebiten.Run(frame, width, height, 2, "Game of Life"); err != nil {
        log.Fatal(err)
    }
}