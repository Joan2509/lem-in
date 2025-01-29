package antfarm

import (
	"reflect"
)

func OptimalsToRooms(optimals [][]Route, rooms *[]Room) [][][]*Room {
	roomCombos := [][][]*Room{} // multiple combinations of routes

	for i, combo := range optimals {
		roomCombos = append(roomCombos, [][]*Room{}) // combination of routes

		for j, route := range combo {
			roomCombos[i] = append(roomCombos[i], []*Room{}) // one route

			for _, roomName := range route {
				thisRoom := &(*rooms)[FindRoom(*rooms, roomName)]
				roomCombos[i][j] = append(roomCombos[i][j], thisRoom)
			}
		}
	}

	return roomCombos
}

func MakeAnts(optimals [][]Route, n int) [][]Ant {
	setsOfAnts := [][]Ant{}

	for i := range optimals {
		setsOfAnts = append(setsOfAnts, []Ant{})
		for j := 0; j < n; j++ {
			setsOfAnts[i] = append(setsOfAnts[i], Ant{Name: j + 1})
		}
	}
	return setsOfAnts
}

func AssignRoutes(optimals [][]Route, optiRooms [][][]*Room, setsOfAnts *[][]Ant) {
	for i, routeCombo := range optimals {
		// how many ants on each route in this combo
		onRoutes := make([]int, len(routeCombo))

		// loop over the set of ants pertaining to the combo of routes
		for j := 0; j < len((*setsOfAnts)[i]); j++ {
			// find the shortest route for this ant (length = route length + ants already taking it)
			shortest := 0
			shortD := len(routeCombo[0]) + onRoutes[0]
			for k := 0; k < len(routeCombo); k++ {
				if len(routeCombo[k])+onRoutes[k] < shortD {
					shortest = k
					shortD = len(routeCombo[k]) + onRoutes[k]
				}
			}
			(*setsOfAnts)[i][j].Route = optiRooms[i][shortest]

			onRoutes[shortest]++
		}
	}
}

func BestSolution(optimals [][][]*Room, setsOfAnts [][]Ant) int {
	if len(optimals) == 1 {
		return 0
	}

	longestRoutes := make([]int, len(optimals))
	for i, combo := range optimals {
		longest := 0
		for _, route := range combo {
			// count ants on this route
			ants := 0
			for _, ant := range setsOfAnts[i] {
				if reflect.DeepEqual(ant.Route, route) {
					ants++
				}
			}

			// turns to complete this route (only compare to longest if active)
			turns := len(route) - 1 + ants
			if ants > 0 && turns > longest {
				longest = turns
			}
		}
		longestRoutes[i] = longest
	}

	// find which optimal route is the quickest for these ants
	quickI := 0
	shortestLong := longestRoutes[0]
	for i, n := range longestRoutes {
		if n < shortestLong {
			shortestLong = n
			quickI = i
		}
	}

	return quickI
}

func PopulateStart(rooms *[]Room, ants []Ant) {
	start := &(*rooms)[StartIndex(*rooms)]
	for _, a := range ants {
		start.Occupants[a.Name] = true
	}
}

// func MoveAnts(ants []Ant) []string {
// 	var turns []string
// 	antsAtEnd := 0

// 	for antsAtEnd < len(ants) {
// 		var moves []string
// 		linksUsed := make(map[string]bool)

// 		for i := range ants {
// 			if ants[i].AtEnd {
// 				continue
// 			}

// 			if canMove, currentRoom, nextRoom := NextIsOk(ants[i], linksUsed); canMove {
// 				// Remove ant from current room
// 				delete(currentRoom.Occupants, ants[i].Name)

// 				// Mark link as used
// 				linksUsed[currentRoom.Name+"-"+nextRoom.Name] = true

// 				// Move ant to next room
// 				nextRoom.Occupants[ants[i].Name] = true
// 				ants[i].RoomIndex++

// 				// Check if ant reached end
// 				if nextRoom.Point == "end" {
// 					ants[i].AtEnd = true
// 					antsAtEnd++
// 				}

// 				// Record move
// 				moves = append(moves, fmt.Sprintf("L%d-%s", ants[i].Name, nextRoom.Name))
// 			}
// 		}

// 		if len(moves) > 0 {
// 			turns = append(turns, joinMoves(moves))
// 		}
// 	}

// 	return turns
// }

// func joinMoves(moves []string) string {
// 	return strings.Join(moves, " ")
// }

// func NextIsOk(a Ant, linksUsed map[string]bool) (bool, *Room, *Room) {
// 	currentRoom := a.Route[a.RoomIndex]

// 	// No move possible if at end
// 	if currentRoom.Point == "end" {
// 		return false, currentRoom, nil
// 	}

// 	nextRoom := a.Route[a.RoomIndex+1]

// 	// Check if link is already used
// 	if linksUsed[currentRoom.Name+"-"+nextRoom.Name] {
// 		return false, currentRoom, nextRoom
// 	}

// 	// Check if next room is empty or is end room
// 	return len(nextRoom.Occupants) < 1 || nextRoom.Point == "end", currentRoom, nextRoom
// }
