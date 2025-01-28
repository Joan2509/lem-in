package antfarm

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetStartValues(file *os.File) (int, []Room, error) {
	scanner := bufio.NewScanner(file)
	rooms := []Room{}

	// First line must be number of ants
	if !scanner.Scan() {
		return 0, rooms, errors.New("ERROR: empty file")
	}

	ants, err := strconv.Atoi(scanner.Text())
	if err != nil || ants < 1 {
		return 0, rooms, errors.New("ERROR: invalid number of ants")
	}

	var prev string
	roomRegex := regexp.MustCompile(`^(\w+)\s([-+]?\d+)\s([-+]?\d+)$`)
	linkRegex := regexp.MustCompile(`^(\w+)-(\w+)$`)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			prev = line
			continue
		}

		if matches := roomRegex.FindStringSubmatch(line); matches != nil {
			room := Room{
				Name:      matches[1],
				Occupants: make(map[int]bool),
				Point:     determineRoomRole(prev),
			}

			room.Cordinates[0], _ = strconv.Atoi(matches[2])
			room.Cordinates[1], _ = strconv.Atoi(matches[3])

			rooms = append(rooms, room)
		}

		if matches := linkRegex.FindStringSubmatch(line); matches != nil {
			LinkRoom(rooms, matches[1], matches[2])
		}

		prev = line
	}

	return ants, rooms, nil
}

func determineRoomRole(prev string) string {
	switch {
	case strings.HasPrefix(prev, "##start"):
		return "start"
	case strings.HasPrefix(prev, "##end"):
		return "end"
	default:
		return "normal"
	}
}

func LinkRoom(rooms []Room, room1, room2 string) {
	for i := range rooms {
		if rooms[i].Name == room1 || rooms[i].Name == room2 {
			if rooms[i].Name == room1 {
				rooms[i].Neighbours = append(rooms[i].Neighbours, room2)
			} else {
				rooms[i].Neighbours = append(rooms[i].Neighbours, room1)
			}
		}
	}
}

// ValidateRooms validates the rooms to ensure there is exactly one start room, one end room,
// no duplicate room names, and no invalid Neighbours.
func ValidateRooms(rooms []Room) error {
	starts := 0
	ends := 0

	for i := 0; i < len(rooms); i++ {
		if rooms[i].Point == "start" {
			starts++
		}
		if rooms[i].Point == "end" {
			ends++
		}

		// Check for duplicate room names
		for j := i + 1; j < len(rooms); j++ {
			if rooms[i].Name == rooms[j].Name {
				return errors.New("ERROR: invalid data format, duplicate room name: " + rooms[i].Name)
			}
		}

		// Check if all linked rooms exist in the list of rooms
		for _, ln := range rooms[i].Neighbours {
			found := false
			for _, rm := range rooms {
				if ln == rm.Name {
					found = true
					break
				}
			}
			if !found {
				return errors.New("ERROR: invalid data format, bad link: " + rooms[i].Name + " > " + ln)
			}
		}
	}

	// Validate that there is exactly one start/end room
	if starts != 1 {
		if starts == 0 {
			return errors.New("ERROR: invalid data format, no start room")
		} else {
			return errors.New("ERROR: invalid data format, too many start rooms")
		}
	}

	if ends != 1 {
		if ends == 0 {
			return errors.New("ERROR: invalid data format, no end room")
		} else {
			return errors.New("ERROR: invalid data format, too many end rooms")
		}
	}

	return nil
}

// ReturnStartIndex returns the index of the "start" room in the slice of rooms.
func ReturnStartIndex(rooms []Room) int {
	for i, r := range rooms {
		if r.Point == "start" {
			return i
		}
	}
	return -1
}
