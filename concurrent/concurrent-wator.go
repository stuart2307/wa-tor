// Wator Concurrent Implementation
// Created: 23/11/25
//	Copyright (C) 2025 Stuart Rossiter
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
)

// Occupied: No = 0, Yes = 1
// Occupant: Water = 0, Fish = 1, Shark = 2
// Energy: Sharks' energy meter
// Breed: Turns since last breed
type Square struct {
	occupied int
	occupant int
	energy   int
	breed    int
}

const scale int = 1

var NumShark = 20000
var NumFish = 300000
var FishBreed = 5
var SharkBreed = 6
var Starve = 4

const width = 960
const height = 540

var startTime time.Time
var chronon = 0

var Threads = 2

var blue color.Color = color.RGBA{69, 145, 196, 255}
var yellow color.Color = color.RGBA{255, 230, 120, 255}
var red color.Color = color.RGBA{255, 0, 0, 255}
var grid [width][height]Square = [width][height]Square{}
var buffer [width][height]Square = [width][height]Square{}
var count int = 0

func wrap(num int, edge int) int {
	if num == -1 {
		return edge - 1
	} else if num == edge {
		return 0
	} else {
		return num
	}
}

func moveFish(x, y int) {
	if grid[x][y].occupied == 1 && grid[x][y].occupant == 1 {
		newBreed := grid[x][y].breed + 1
		n := (grid[x][wrap(y+1, height)].occupied << 3) + (grid[wrap(x+1, width)][y].occupied << 2) + (grid[x][wrap(y-1, height)].occupied << 1) + grid[wrap(x-1, width)][y].occupied
		if buffer[x][y].occupant == 2 {

		} else if n == 15 {
			buffer[x][y] = grid[x][y]
		} else {
			if grid[x][y].breed == FishBreed {
				newBreed = 0
				buffer[x][y] = grid[x][y]
				buffer[x][y].breed = 0
			}
			directions := []int{1, 2, 4, 8}
			rand.Shuffle(4, func(i int, j int) {
				directions[i], directions[j] = directions[j], directions[i]
			})
			moved := false
			for choice := 0; choice < 4; choice++ {
				if n&directions[choice] == 0 {
					switch directions[choice] {
					case 1:
						if buffer[wrap(x-1, width)][y].occupied == 0 {
							buffer[wrap(x-1, width)][y] = grid[x][y]
							buffer[wrap(x-1, width)][y].breed = newBreed
							choice = 5
							moved = true
						}

					case 2:
						if buffer[x][wrap(y-1, height)].occupied == 0 {
							buffer[x][wrap(y-1, height)] = grid[x][y]
							buffer[x][wrap(y-1, height)].breed = newBreed
							choice = 5
							moved = true
						}

					case 4:
						if buffer[wrap(x+1, width)][y].occupied == 0 {
							buffer[wrap(x+1, width)][y] = grid[x][y]
							buffer[wrap(x+1, width)][y].breed = newBreed
							choice = 5
							moved = true
						}

					case 8:
						if buffer[x][wrap(y+1, height)].occupied == 0 {
							buffer[x][wrap(y+1, height)] = grid[x][y]
							buffer[x][wrap(y+1, height)].breed = newBreed
							choice = 5
							moved = true
						}

					}
				}

			}
			if !moved {
				buffer[x][y] = grid[x][y]
				buffer[x][y].breed = newBreed
			}
		}
	}
}

func moveSharks(x, y int) {
	if grid[x][y].occupied == 1 && grid[x][y].occupant == 2 {
		newEnergy := grid[x][y].energy - 1
		newBreed := grid[x][y].breed + 1
		if newEnergy > 0 {
			n := (grid[x][wrap(y+1, height)].occupant << 3) + (grid[wrap(x+1, width)][y].occupant << 2) + (grid[x][wrap(y-1, height)].occupant << 1) + grid[wrap(x-1, width)][y].occupant
			if n == 30 {
				buffer[x][y] = grid[x][y]
			} else {
				if grid[x][y].breed == SharkBreed {
					newBreed = 0
					buffer[x][y] = grid[x][y]
					buffer[x][y].breed = 0
					buffer[x][y].energy = Starve
				}
				//fish are 0, else is 1
				var n = (^((grid[x][wrap(y+1, height)].occupant % 2 << 3) + (grid[wrap(x+1, width)][y].occupant % 2 << 2) + (grid[x][wrap(y-1, height)].occupant % 2 << 1) + grid[wrap(x-1, width)][y].occupant%2)) & 15
				if n == 15 {
					n = (grid[x][wrap(y+1, height)].occupied << 3) + (grid[wrap(x+1, width)][y].occupied << 2) + (grid[x][wrap(y-1, height)].occupied << 1) + grid[wrap(x-1, width)][y].occupied
				}
				directions := []int{1, 2, 4, 8}
				rand.Shuffle(4, func(i int, j int) {
					directions[i], directions[j] = directions[j], directions[i]
				})
				moved := false
				for choice := 0; choice < 4; choice++ {
					if n&directions[choice] == 0 {
						switch directions[choice] {
						case 1:
							if buffer[wrap(x-1, width)][y].occupant != 2 {
								newEnergy += Starve * (grid[wrap(x-1, width)][y].occupant % 2)
								buffer[wrap(x-1, width)][y] = grid[x][y]
								buffer[wrap(x-1, width)][y].energy = newEnergy
								buffer[wrap(x-1, width)][y].breed = newBreed
								if buffer[wrap(x-1, width)][y].energy > Starve {
									buffer[wrap(x-1, width)][y].energy = Starve
								}
								choice = 5
								moved = true
							}

						case 2:
							if buffer[x][wrap(y-1, height)].occupant != 2 {
								newEnergy += Starve * (grid[x][wrap(y-1, height)].occupant % 2)
								buffer[x][wrap(y-1, height)] = grid[x][y]
								buffer[x][wrap(y-1, height)].energy = newEnergy
								buffer[x][wrap(y-1, height)].breed = newBreed
								if buffer[x][wrap(y-1, height)].energy > Starve {
									buffer[x][wrap(y-1, height)].energy = Starve
								}
								choice = 5
								moved = true
							}

						case 4:
							if buffer[wrap(x+1, width)][y].occupant != 2 {
								newEnergy += Starve * (grid[wrap(x+1, width)][y].occupant % 2)
								buffer[wrap(x+1, width)][y] = grid[x][y]
								buffer[wrap(x+1, width)][y].energy = newEnergy
								buffer[wrap(x+1, width)][y].breed = newBreed
								if buffer[wrap(x+1, width)][y].energy > Starve {
									buffer[wrap(x+1, width)][y].energy = Starve
								}
								choice = 5
								moved = true
							}

						case 8:
							if buffer[x][wrap(y+1, height)].occupant != 2 {
								newEnergy += Starve * (grid[x][wrap(y+1, height)].occupant % 2)
								buffer[x][wrap(y+1, height)] = grid[x][y]
								buffer[x][wrap(y+1, height)].energy = newEnergy
								buffer[x][wrap(y+1, height)].breed = newBreed
								if buffer[x][wrap(y+1, height)].energy > Starve {
									buffer[x][wrap(y+1, height)].energy = Starve
								}
								choice = 5
								moved = true
							}

						}
					}
				}
				if !moved {
					buffer[x][y].breed = newBreed
					buffer[x][y] = grid[x][y]
					buffer[x][y].energy = newEnergy
				}
			}
		}
	}
}

func update() error {
	buffer = [width][height]Square{}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if grid[x][y].occupant == 1 {
				moveFish(x, y)
			} else if grid[x][y].occupant == 2 {
				moveSharks(x, y)
			}
		}
	}
	grid = buffer
	return nil
}

func processData(ystart, ylen int) {
	for x := 0; x < width; x++ {
		for y := ystart; y < ystart+ylen; y++ {
			if grid[x][y].occupant == 1 {
				moveFish(x, y)
			} else if grid[x][y].occupant == 2 {
				moveSharks(x, y)
			}
		}
	}
}

func concUpdate() error {
	buffer = [width][height]Square{}
	var wg = sync.WaitGroup{}
	jobs := make(chan int, Threads)
	for i := 0; i < Threads; i++ {
		wg.Add(1)
		go func(id int) {
			for row := range jobs {
				processData(row, 10)
			}
			wg.Done()
		}(i)
	}
	for row := 0; row < height; row += 10 {
		jobs <- row
	}
	close(jobs)

	wg.Wait()
	grid = buffer
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
	chronon++
	var err error = nil
	if count == 1 {
		err = concUpdate()
		count = 0
	}
	if !ebiten.IsDrawingSkipped() {
		display(window)
	}
	if chronon == 1000 {
		var now = time.Since(startTime)
		fmt.Println(now)
		os.Exit(0)
	}
	return err
}

func main() {
	flatGrid := make([]Square, width*height)
	for i := 0; i < NumFish; i++ {
		flatGrid[i].occupant = 1
		flatGrid[i].occupied = 1
		flatGrid[i].breed = 0
		flatGrid[i].energy = 0
	}
	for i := NumFish; i < NumFish+NumShark; i++ {
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
	startTime = time.Now()
	if err := ebiten.Run(frame, width, height, 2, "Concurrent Wa-Tor"); err != nil {
		log.Fatal(err)
	}
}
