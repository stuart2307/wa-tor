// Wator Sequential Implementation
// Created: 3/11/25
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

// How many pixels a square's dimensions are.
const scale int = 1

// NumShark = Number of starting sharks.
// NumFish = Number of starting fish.
// FishBreed = Number of turns it takes for a fish to breed.
// SharkBreed = Number of turns it takes for a shark to breed.
// Starve = How many turns it takes for a shark to starve (Max Energy)
var NumShark = 10000
var NumFish = 200000
var FishBreed = 5
var SharkBreed = 6
var Starve = 4
var startTime time.Time

const width = 1000
const height = 800

var Threads = 1
var chronon = 0

var blue color.Color = color.RGBA{69, 145, 196, 255}
var yellow color.Color = color.RGBA{255, 230, 120, 255}
var red color.Color = color.RGBA{255, 0, 0, 255}
var grid [width][height]Square = [width][height]Square{}
var buffer [width][height]Square = [width][height]Square{}
var count int = 0

// Wrap:
// Wraps around a co-ordinate. Makes it so that if a co-ordinate goes off the screen, it wraps back around to the other side.
func Wrap(num int, edge int) int {
	if num == -1 {
		return edge - 1
	} else if num == edge {
		return 0
	} else {
		return num
	}
}

// MoveFish:
// Processes the movement of fish.
// If the fish has already been eaten, it does nothing.
// If the fish is surrounded, it stays put.
// If the fish can move, it does.
// If the fish moves and is able to breed, it leaves behind a new fish as it moves.
func MoveFish(x, y int) {
	if grid[x][y].occupied == 1 && grid[x][y].occupant == 1 {
		newBreed := grid[x][y].breed + 1
		n := (grid[x][Wrap(y+1, height)].occupied << 3) + (grid[Wrap(x+1, width)][y].occupied << 2) + (grid[x][Wrap(y-1, height)].occupied << 1) + grid[Wrap(x-1, width)][y].occupied
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
						if buffer[Wrap(x-1, width)][y].occupied == 0 {
							buffer[Wrap(x-1, width)][y] = grid[x][y]
							buffer[Wrap(x-1, width)][y].breed = newBreed
							choice = 5
							moved = true
						}
					case 2:
						if buffer[x][Wrap(y-1, height)].occupied == 0 {
							buffer[x][Wrap(y-1, height)] = grid[x][y]
							buffer[x][Wrap(y-1, height)].breed = newBreed
							choice = 5
							moved = true
						}
					case 4:
						if buffer[Wrap(x+1, width)][y].occupied == 0 {
							buffer[Wrap(x+1, width)][y] = grid[x][y]
							buffer[Wrap(x+1, width)][y].breed = newBreed
							choice = 5
							moved = true
						}
					case 8:
						if buffer[x][Wrap(y+1, height)].occupied == 0 {
							buffer[x][Wrap(y+1, height)] = grid[x][y]
							buffer[x][Wrap(y+1, height)].breed = newBreed
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

// MoveSharks:
// Processes movement for the sharks.
// If the shark is completely surrounded by other sharks, stay put.
// If the shark does move and can breed, leave a new shark behind as it moves.
// The shark prioritises fish squares.
// If it sees no fish squares around it, it looks for empty squares.
// When it moves, if it lands on a fish, it eats it and replenishes energy.
func MoveSharks(x, y int) {
	if grid[x][y].occupied == 1 && grid[x][y].occupant == 2 {
		newEnergy := grid[x][y].energy - 1
		newBreed := grid[x][y].breed + 1
		if newEnergy > 0 {
			n := (grid[x][Wrap(y+1, height)].occupant << 3) + (grid[Wrap(x+1, width)][y].occupant << 2) + (grid[x][Wrap(y-1, height)].occupant << 1) + grid[Wrap(x-1, width)][y].occupant
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
				var n = (^((grid[x][Wrap(y+1, height)].occupant % 2 << 3) + (grid[Wrap(x+1, width)][y].occupant % 2 << 2) + (grid[x][Wrap(y-1, height)].occupant % 2 << 1) + grid[Wrap(x-1, width)][y].occupant%2)) & 15
				if n == 15 {
					n = (grid[x][Wrap(y+1, height)].occupied << 3) + (grid[Wrap(x+1, width)][y].occupied << 2) + (grid[x][Wrap(y-1, height)].occupied << 1) + grid[Wrap(x-1, width)][y].occupied
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
							if buffer[Wrap(x-1, width)][y].occupant != 2 {
								newEnergy += Starve * (grid[Wrap(x-1, width)][y].occupant % 2)
								buffer[Wrap(x-1, width)][y] = grid[x][y]
								buffer[Wrap(x-1, width)][y].energy = newEnergy
								buffer[Wrap(x-1, width)][y].breed = newBreed
								if buffer[Wrap(x-1, width)][y].energy > Starve {
									buffer[Wrap(x-1, width)][y].energy = Starve
								}
								choice = 5
								moved = true
							}
						case 2:
							if buffer[x][Wrap(y-1, height)].occupant != 2 {
								newEnergy += Starve * (grid[x][Wrap(y-1, height)].occupant % 2)
								buffer[x][Wrap(y-1, height)] = grid[x][y]
								buffer[x][Wrap(y-1, height)].energy = newEnergy
								buffer[x][Wrap(y-1, height)].breed = newBreed
								if buffer[x][Wrap(y-1, height)].energy > Starve {
									buffer[x][Wrap(y-1, height)].energy = Starve
								}
								choice = 5
								moved = true
							}
						case 4:
							if buffer[Wrap(x+1, width)][y].occupant != 2 {
								newEnergy += Starve * (grid[Wrap(x+1, width)][y].occupant % 2)
								buffer[Wrap(x+1, width)][y] = grid[x][y]
								buffer[Wrap(x+1, width)][y].energy = newEnergy
								buffer[Wrap(x+1, width)][y].breed = newBreed
								if buffer[Wrap(x+1, width)][y].energy > Starve {
									buffer[Wrap(x+1, width)][y].energy = Starve
								}
								choice = 5
								moved = true
							}
						case 8:
							if buffer[x][Wrap(y+1, height)].occupant != 2 {
								newEnergy += Starve * (grid[x][Wrap(y+1, height)].occupant % 2)
								buffer[x][Wrap(y+1, height)] = grid[x][y]
								buffer[x][Wrap(y+1, height)].energy = newEnergy
								buffer[x][Wrap(y+1, height)].breed = newBreed
								if buffer[x][Wrap(y+1, height)].energy > Starve {
									buffer[x][Wrap(y+1, height)].energy = Starve
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

// Update:
// Loops through the grid, processing movement for any fish or sharks it comes across.
func Update() error {
	buffer = [width][height]Square{}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if grid[x][y].occupant == 1 {
				MoveFish(x, y)
			} else if grid[x][y].occupant == 2 {
				MoveSharks(x, y)
			}
		}
	}
	grid = buffer
	return nil
}

// Display:
// Displays the grid using the ebiten engine
func Display(window *ebiten.Image) {
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

// Frame:
// Runs every Frame as part of the ebiten window
func Frame(window *ebiten.Image) error {
	count++
	chronon++
	var err error = nil
	if count == 1 {
		err = Update()
		count = 0
	}
	if !ebiten.IsDrawingSkipped() {
		Display(window)
	}
	if chronon == 1000 {
		var now = time.Since(startTime)
		fmt.Println(now)
	}
	return err
}

// Randomises the grid, filling it with the appropriate amount of fish and sharks.
// Runs the ebiten engine once the grid has been shuffled.
func main() {
	flatGrid := make([]Square, width*height)
	for i := 0; i < NumFish; i++ {
		flatGrid[i].occupant = 1
		flatGrid[i].occupied = 1
		flatGrid[i].breed = 0
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
	buffer = grid
	startTime = time.Now()
	if err := ebiten.Run(Frame, width, height, 2, "Wa-Tor"); err != nil {
		log.Fatal(err)
	}
}
