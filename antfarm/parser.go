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
				Point:      determineRoomRole(prev),
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

func ValidateRooms(rooms []Room) error {
	nameCount := make(map[string]int)
	startCount, endCount := 0, 0

	for _, room := range rooms {
		// Count role occurrences
		switch room.Point {
		case "start":
			startCount++
		case "end":
			endCount++
		}

		// Check for duplicate names
		nameCount[room.Name]++
		if nameCount[room.Name] > 1 {
			return errors.New("ERROR: duplicate room name: " + room.Name)
		}

		// Validate links
		for _, neighbour := range room.Neighbours {
			if !roomExists(rooms, neighbour) {
				return errors.New("ERROR: invalid link: " + room.Name + " > " + neighbour)
			}
		}
	}

	// Validate room roles
	switch {
	case startCount == 0:
		return errors.New("ERROR: no start room")
	case startCount > 1:
		return errors.New("ERROR: too many start rooms")
	case endCount == 0:
		return errors.New("ERROR: no end room")
	case endCount > 1:
		return errors.New("ERROR: too many end rooms")
	}

	return nil
}

func roomExists(rooms []Room, name string) bool {
	for _, room := range rooms {
		if room.Name == name {
			return true
		}
	}
	return false
}

func StartIndex(rooms []Room) int {
	for i, room := range rooms {
		if room.Point == "start" {
			return i
		}
	}
	return -1
}
