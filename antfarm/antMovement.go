package antfarm
import (
	"fmt"
	"reflect"
	"sync"
)

var (
	sepMu sync.Mutex
)

// FindRoutes recursively finds all valid routes from the current room to the end room.
// It appends each valid route to the routes slice.
// curRoom: The current room being explored.
// curRoute: The current path being built.
// routes: A pointer to the slice of all valid routes.
// rooms: A pointer to the slice of all rooms in the farm.
func FindRoutes(curRoom Room, curRoute Route, routes *[]Route, rooms *[]Room) {
	// reached the end, add to routes
	if curRoom.Role == "end" {
		curRoute = append(curRoute, curRoom.Name)
		toSave := make(Route, len(curRoute))
		copy(toSave, curRoute) // copy values to a new route to avoid pointer problems
		*routes = append(*routes, toSave)
		return
	}

	// add new room to current route and proceed
	if !IsOnRoute(curRoute, curRoom) {
		curRoute = append(curRoute, curRoom.Name)
		for _, link := range curRoom.Links {
			nextRoom := (*rooms)[FindRoom(*rooms, link)]
			FindRoutes(nextRoom, curRoute, routes, rooms)
		}
	}
}
// SortRoutes sorts a slice of routes from shortest to longest.
// rts: A pointer to the slice of routes to be sorted.
// Returns an error if the slice is empty.
func SortRoutes(rts *[]Route) error {
	if len(*rts) < 1 {
		return fmt.Errorf("ERROR: invalid data format, no valid routes")
	}

	for i := 0; i < len(*rts)-1; i++ {
		for j := i + 1; j < len(*rts); j++ {
			if len((*rts)[i]) > len((*rts)[j]) {
				(*rts)[i], (*rts)[j] = (*rts)[j], (*rts)[i]
			}
		}
	}

	return nil
}
// IsOnRoute checks if a room is already in the current route.
// route: The current route being checked.
// room: The room to check for.
// Returns true if the room is in the route, otherwise false.
func IsOnRoute(route Route, room Room) bool {
	for _, r := range route {
		if room.Name == r {
			return true
		}
	}
	return false
}

