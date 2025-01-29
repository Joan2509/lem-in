package antfarm

import "reflect"

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
