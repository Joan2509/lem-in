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
