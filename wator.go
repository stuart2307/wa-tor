package main

import (
    "image/color"
    "log"
    "math/rand"
    "time"

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
var NumShark = 30000
var NumFish = 10000
var FishBreed = 5
var SharkBreed = 10
var Starve = 10
var SharkEnergyRestore = 2
const width = 300
const height = 300
var Threads = 1

var blue color.Color = color.RGBA{69, 145, 196, 255}
var yellow color.Color = color.RGBA{255, 230, 120, 255}
var red color.Color = color.RGBA{255, 0, 0, 255}
var grid [width][height]Square = [width][height]Square{}
var buffer [width][height]Square = [width][height]Square{}
var count int = 0

func wrap(num int, edge int) int {
    if (num == -1) {
        return edge - 1
    } else if (num == edge) {
        return 0
    } else {
        return num
    }
}

func moveFish() {
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            if (grid[x][y].occupied == 1 && grid[x][y].occupant == 1) {
                n := (grid[x][wrap(y+1, height)].occupied << 3) + (grid[wrap(x+1, width)][y].occupied << 2) + (grid[x][wrap(y-1, height)].occupied << 1) + grid[wrap(x-1, width)][y].occupied
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
                                buffer[wrap(x-1, width)][y] = grid[x][y]
                                buffer[wrap(x-1, width)][y].breed++;
                                choice = -1
                                break
                            case 2:
                                buffer[x][wrap(y-1, height)] = grid[x][y]
                                buffer[x][wrap(y-1, height)].breed++;
                                choice = -1
                                break
                            case 4:
                                buffer[wrap(x+1, width)][y] = grid[x][y]
                                buffer[wrap(x+1, width)][y].breed++;
                                choice = -1
                                break
                            case 8: 
                                buffer[x][wrap(y+1, height)] = grid[x][y]
                                buffer[x][wrap(y+1, height)].breed++;
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

func moveSharks() {
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            if (grid[x][y].occupied == 1 && grid[x][y].occupant == 2) {
                grid[x][y].energy--;
                if (grid[x][y].energy == 0) {
                    buffer[x][y] = Square{}
                } else {
                    n := (grid[x][wrap(y+1, height)].occupant << 3) + (grid[wrap(x+1, width)][y].occupant << 2) + (grid[x][wrap(y-1, height)].occupant << 1) + grid[wrap(x-1, width)][y].occupant
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
                        //fish are 0, else is 1
                        var n = ^((grid[x][wrap(y+1, height)].occupant%2 << 3) + (grid[wrap(x+1, width)][y].occupant%2 << 2) + (grid[x][wrap(y-1, height)].occupant%2 << 1) + grid[wrap(x-1, width)][y].occupant%2) & 15
                        if (n == 15) {
                            n = (grid[x][wrap(y+1, height)].occupied << 3) + (grid[wrap(x+1, width)][y].occupied << 2) + (grid[x][wrap(y-1, height)].occupied << 1) + grid[wrap(x-1, width)][y].occupied
                        }
                        var choice = 1 << rand.Intn(4)
                        for (choice != -1) {
                            if (n & choice == 0) {
                                switch choice {
                                case 1:
                                    buffer[wrap(x-1, width)][y] = grid[x][y]
                                    buffer[wrap(x-1, width)][y].breed++;
                                    buffer[wrap(x-1, width)][y].energy += SharkEnergyRestore * grid[wrap(x-1, width)][y].occupant
                                    if (buffer[wrap(x-1, width)][y].energy > Starve) {buffer[wrap(x-1, width)][y].energy = Starve}
                                    choice = -1
                                    break
                                case 2:
                                    buffer[x][wrap(y-1, height)] = grid[x][y]
                                    buffer[x][wrap(y-1, height)].breed++;
                                    buffer[x][wrap(y-1, height)].energy += SharkEnergyRestore * grid[x][wrap(y-1, height)].occupant
                                    if (buffer[x][wrap(y-1, height)].energy > Starve) {buffer[x][wrap(y-1, height)].energy = Starve}
                                    choice = -1
                                    break
                                case 4:
                                    buffer[wrap(x+1, width)][y] = grid[x][y]
                                    buffer[wrap(x+1, width)][y].breed++;
                                    buffer[wrap(x+1, width)][y].energy += SharkEnergyRestore * grid[wrap(x+1, width)][y].occupant
                                    if (buffer[wrap(x+1, width)][y].energy > Starve) {buffer[wrap(x+1, width)][y].energy = Starve}
                                    choice = -1
                                    break
                                case 8:  
                                    buffer[x][wrap(y+1, height)] = grid[x][y]
                                    buffer[x][wrap(y+1, height)].breed++;
                                    buffer[x][wrap(y+1, height)].energy += SharkEnergyRestore * grid[x][wrap(y+1, height)].occupant
                                    if (buffer[x][wrap(y+1, height)].energy > Starve) {buffer[x][wrap(y+1, height)].energy = Starve}
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
}


func update() error {
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            buffer[x][y] = Square{}
        }
    }
    moveFish()
    moveSharks()
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
                    } else if grid[x][y].occupant == 2 {
                        window.Set(x*scale+i, y*scale+j, red)
                    }
                }
            }
        }
    }
}

func frame(window *ebiten.Image) error {
    count++
    var err error = nil
    if count == 10 {
        err = update()
        count = 0
    }
    if !ebiten.IsDrawingSkipped() {
        display(window)
    }
    return err
}

func main() {
    rand.Seed(time.Now().UnixMicro())
    flatGrid := make([]Square, width*height)
    for i := 0; i < NumFish; i++ {
        flatGrid[i].occupant = 1
        flatGrid[i].occupied = 1
        flatGrid[i].breed = 0
    }
    for i := NumFish; i < NumFish + NumShark; i++ {
        flatGrid[i].occupant = 2
        flatGrid[i].occupied = 1
        flatGrid[i].breed = 0
        flatGrid[i].energy = Starve
    }
    rand.Shuffle(len(flatGrid), func(i int, j int) {
        flatGrid[i], flatGrid[j] = flatGrid[j], flatGrid[i]
    })
    for x := 0; x < width*height; x++ {
        grid[x%width][x/width].occupied = flatGrid[x].occupied
        grid[x%width][x/width].occupant = flatGrid[x].occupant
        grid[x%width][x/width].breed = flatGrid[x].breed
        grid[x%width][x/width].energy = flatGrid[x].energy
    }
    if err := ebiten.Run(frame, width, height, 2, "Wa-Tor"); err != nil {
        log.Fatal(err)
    }
}