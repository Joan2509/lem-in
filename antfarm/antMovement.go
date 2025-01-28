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
