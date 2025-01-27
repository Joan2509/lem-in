package antfarm

type Room struct {
	Name      string
	Cordinates    [2]int
	Occupants map[int]bool
	Neighbours     []string // neighbouring room names
	Point      string   // "start", "normal" or "end"
}

type Ant struct {
	Name      int
	Route     []*Room
	RoomIndex int
	AtEnd     bool
}

type Route []string