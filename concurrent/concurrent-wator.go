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

package concurrent_wator

import (
	"image/color"
	"log"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten"
)

// Occupied: No = 0, Yes = 1
// Occupant: Water = 0, Fish = 1, Shark = 2
// Energy: Sharks' energy meter
// Breed: Turns since last breed
type Square struct {
	Occupied int
	Occupant int
	Energy   int
	Breed    int
}

// How many pixels a square's dimensions are/
const scale = 1

// Number of sharks to begin the simulation with.
var NumShark = 20000

// Number of fish to begin the simulation with.
var NumFish = 200000

// Number of turns until a fish can breed.
var FishBreed = 5

// Number of turns until a shark can breed.
var SharkBreed = 6

// Number of turns a shark can go without eating until it dies of starvation.
var Starve = 4

const Width = 1600
const Height = 900

// The amount of threads to be used for grid processing
var Threads = 4

var blue color.Color = color.RGBA{69, 145, 196, 255}
var yellow color.Color = color.RGBA{255, 230, 120, 255}
var red color.Color = color.RGBA{255, 0, 0, 255}
var grid [Width][Height]Square = [Width][Height]Square{}
var buffer [Width][Height]Square = [Width][Height]Square{}
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
	if grid[x][y].Occupied == 1 && grid[x][y].Occupant == 1 {
		newBreed := grid[x][y].Breed + 1
		n := (grid[x][Wrap(y+1, Height)].Occupied << 3) + (grid[Wrap(x+1, Width)][y].Occupied << 2) + (grid[x][Wrap(y-1, Height)].Occupied << 1) + grid[Wrap(x-1, Width)][y].Occupied
		if buffer[x][y].Occupant == 2 {

		} else if n == 15 {
			buffer[x][y] = grid[x][y]
		} else {
			if grid[x][y].Breed == FishBreed {
				newBreed = 0
				buffer[x][y] = grid[x][y]
				buffer[x][y].Breed = 0
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
						if buffer[Wrap(x-1, Width)][y].Occupied == 0 {
							buffer[Wrap(x-1, Width)][y] = grid[x][y]
							buffer[Wrap(x-1, Width)][y].Breed = newBreed
							choice = 5
							moved = true
						}

					case 2:
						if buffer[x][Wrap(y-1, Height)].Occupied == 0 {
							buffer[x][Wrap(y-1, Height)] = grid[x][y]
							buffer[x][Wrap(y-1, Height)].Breed = newBreed
							choice = 5
							moved = true
						}

					case 4:
						if buffer[Wrap(x+1, Width)][y].Occupied == 0 {
							buffer[Wrap(x+1, Width)][y] = grid[x][y]
							buffer[Wrap(x+1, Width)][y].Breed = newBreed
							choice = 5
							moved = true
						}

					case 8:
						if buffer[x][Wrap(y+1, Height)].Occupied == 0 {
							buffer[x][Wrap(y+1, Height)] = grid[x][y]
							buffer[x][Wrap(y+1, Height)].Breed = newBreed
							choice = 5
							moved = true
						}

					}
				}

			}
			if !moved {
				buffer[x][y] = grid[x][y]
				buffer[x][y].Breed = newBreed
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
	if grid[x][y].Occupied == 1 && grid[x][y].Occupant == 2 {
		newEnergy := grid[x][y].Energy - 1
		newBreed := grid[x][y].Breed + 1
		if newEnergy > 0 {
			n := (grid[x][Wrap(y+1, Height)].Occupant << 3) + (grid[Wrap(x+1, Width)][y].Occupant << 2) + (grid[x][Wrap(y-1, Height)].Occupant << 1) + grid[Wrap(x-1, Width)][y].Occupant
			if n == 30 {
				buffer[x][y] = grid[x][y]
			} else {
				if grid[x][y].Breed == SharkBreed {
					newBreed = 0
					buffer[x][y] = grid[x][y]
					buffer[x][y].Breed = 0
					buffer[x][y].Energy = Starve
				}
				//fish are 0, else is 1
				var n = (^((grid[x][Wrap(y+1, Height)].Occupant % 2 << 3) + (grid[Wrap(x+1, Width)][y].Occupant % 2 << 2) + (grid[x][Wrap(y-1, Height)].Occupant % 2 << 1) + grid[Wrap(x-1, Width)][y].Occupant%2)) & 15
				if n == 15 {
					n = (grid[x][Wrap(y+1, Height)].Occupied << 3) + (grid[Wrap(x+1, Width)][y].Occupied << 2) + (grid[x][Wrap(y-1, Height)].Occupied << 1) + grid[Wrap(x-1, Width)][y].Occupied
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
							if buffer[Wrap(x-1, Width)][y].Occupant != 2 {
								newEnergy += Starve * (grid[Wrap(x-1, Width)][y].Occupant % 2)
								buffer[Wrap(x-1, Width)][y] = grid[x][y]
								buffer[Wrap(x-1, Width)][y].Energy = newEnergy
								buffer[Wrap(x-1, Width)][y].Breed = newBreed
								if buffer[Wrap(x-1, Width)][y].Energy > Starve {
									buffer[Wrap(x-1, Width)][y].Energy = Starve
								}
								choice = 5
								moved = true
							}

						case 2:
							if buffer[x][Wrap(y-1, Height)].Occupant != 2 {
								newEnergy += Starve * (grid[x][Wrap(y-1, Height)].Occupant % 2)
								buffer[x][Wrap(y-1, Height)] = grid[x][y]
								buffer[x][Wrap(y-1, Height)].Energy = newEnergy
								buffer[x][Wrap(y-1, Height)].Breed = newBreed
								if buffer[x][Wrap(y-1, Height)].Energy > Starve {
									buffer[x][Wrap(y-1, Height)].Energy = Starve
								}
								choice = 5
								moved = true
							}

						case 4:
							if buffer[Wrap(x+1, Width)][y].Occupant != 2 {
								newEnergy += Starve * (grid[Wrap(x+1, Width)][y].Occupant % 2)
								buffer[Wrap(x+1, Width)][y] = grid[x][y]
								buffer[Wrap(x+1, Width)][y].Energy = newEnergy
								buffer[Wrap(x+1, Width)][y].Breed = newBreed
								if buffer[Wrap(x+1, Width)][y].Energy > Starve {
									buffer[Wrap(x+1, Width)][y].Energy = Starve
								}
								choice = 5
								moved = true
							}

						case 8:
							if buffer[x][Wrap(y+1, Height)].Occupant != 2 {
								newEnergy += Starve * (grid[x][Wrap(y+1, Height)].Occupant % 2)
								buffer[x][Wrap(y+1, Height)] = grid[x][y]
								buffer[x][Wrap(y+1, Height)].Energy = newEnergy
								buffer[x][Wrap(y+1, Height)].Breed = newBreed
								if buffer[x][Wrap(y+1, Height)].Energy > Starve {
									buffer[x][Wrap(y+1, Height)].Energy = Starve
								}
								choice = 5
								moved = true
							}

						}
					}
				}
				if !moved {
					buffer[x][y].Breed = newBreed
					buffer[x][y] = grid[x][y]
					buffer[x][y].Energy = newEnergy
				}
			}
		}
	}
}

// ProcessData:
// Loops through the grid, processing movement for any fish or sharks it comes across.
func ProcessData(ystart, ylen int) {
	for x := 0; x < Width; x++ {
		for y := ystart; y < ystart+ylen; y++ {
			if grid[x][y].Occupant == 1 {
				MoveFish(x, y)
			} else if grid[x][y].Occupant == 2 {
				MoveSharks(x, y)
			}
		}
	}
}

// ConcUpdate:
// Creates workers and assigns them rows of the screen to process.
func ConcUpdate() error {
	buffer = [Width][Height]Square{}
	var wg = sync.WaitGroup{}
	jobs := make(chan int, Threads)
	for i := 0; i < Threads; i++ {
		wg.Add(1)
		go func(id int) {
			for row := range jobs {
				ProcessData(row, 10)
			}
			wg.Done()
		}(i)
	}
	for row := 0; row < Height; row += 10 {
		jobs <- row
	}
	close(jobs)

	wg.Wait()
	grid = buffer
	return nil
}

// Display:
// Displays the grid using the ebiten engine
func Display(window *ebiten.Image) {
	window.Fill(blue)

	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			for i := 0; i < scale; i++ {
				for j := 0; j < scale; j++ {
					if grid[x][y].Occupant == 1 {
						window.Set(x*scale+i, y*scale+j, yellow)
					} else if grid[x][y].Occupant == 2 {
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
	var err error = nil
	if count == 1 {
		err = ConcUpdate()
		count = 0
	}
	if !ebiten.IsDrawingSkipped() {
		Display(window)
	}
	return err
}

// Randomises the grid, filling it with the appropriate amount of fish and sharks.
// Runs the ebiten engine once the grid has been shuffled.
func MainFunc() {
	flatGrid := make([]Square, Width*Height)
	for i := 0; i < NumFish; i++ {
		flatGrid[i].Occupant = 1
		flatGrid[i].Occupied = 1
		flatGrid[i].Breed = 0
		flatGrid[i].Energy = 0
	}
	for i := NumFish; i < NumFish+NumShark; i++ {
		flatGrid[i].Occupant = 2
		flatGrid[i].Occupied = 1
		flatGrid[i].Breed = 0
		flatGrid[i].Energy = Starve
	}
	rand.Shuffle(len(flatGrid), func(i int, j int) {
		flatGrid[i], flatGrid[j] = flatGrid[j], flatGrid[i]
	})
	for x := 0; x < Width*Height; x++ {
		grid[x%Width][x/Width].Occupied = flatGrid[x].Occupied
		grid[x%Width][x/Width].Occupant = flatGrid[x].Occupant
		grid[x%Width][x/Width].Breed = flatGrid[x].Breed
		grid[x%Width][x/Width].Energy = flatGrid[x].Energy
	}
	if err := ebiten.Run(Frame, Width, Height, scale, "Concurrent Wa-Tor"); err != nil {
		log.Fatal(err)
	}
}
