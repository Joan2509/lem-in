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
