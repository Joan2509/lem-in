package main

import (
	"fmt"
	"os"
	"sync"

	"antfarm/antfarm"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: provide the input file in one argument")
		return
	}

	// Open the input file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer file.Close()

	// Read and save the number of ants and information about the rooms
	nAnts, rooms, err := antfarm.GetStartValues(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Verify that the rooms are valid (one start, one end, no duplicates, etc.)
	err = antfarm.ValidateRooms(rooms)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Find all routes connecting "start" to "end"
	var routes []antfarm.Route
	startRoom := rooms[antfarm.StartIndex(rooms)]
	antfarm.FindRoutes(startRoom, antfarm.Route{}, &routes, &rooms)
	err = antfarm.SortRoutes(&routes)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Find all combinations of non-overlapping routes
	var combosOfSeparates [][]antfarm.Route
	wg := sync.WaitGroup{}
	for i := range routes {
		wg.Add(1)
		go antfarm.FindSeparates(routes, []antfarm.Route{}, &combosOfSeparates, i, &wg)
	}
	wg.Wait()

	// Reduce the amount of solutions to test
	optimals := antfarm.ShortCombos(combosOfSeparates, routes)
	optimals = append(optimals, antfarm.LowAverages(combosOfSeparates)...) // lowAverages() also finds the longest combinations
	optimals = antfarm.RemoveRedundant(optimals)

	// Convert optimal routes to slices of rooms
	optiRooms := antfarm.OptimalsToRooms(optimals, &rooms)

	// Create ants and assign routes
	setsOfAnts := antfarm.MakeAnts(optimals, nAnts)
	antfarm.AssignRoutes(optimals, optiRooms, &setsOfAnts)

	// Find the best solution for the given number of ants
	optI := antfarm.BestSolution(optiRooms, setsOfAnts)

	// Populate the start room with ants
	antfarm.PopulateStart(&rooms, setsOfAnts[optI])

	// Move ants and save the moves
	turns := antfarm.MoveAnts(setsOfAnts[optI])

	// Print out the file contents and the moves
	antfarm.PrintSolution(turns)
}
