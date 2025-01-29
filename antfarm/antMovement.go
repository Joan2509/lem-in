package antfarm

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
