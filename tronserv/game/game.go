package game

type Game struct {
	Code        string   `json:"code"`
	Players     []Player `json:"players"`
	Rounds      []Round  `json:"rounds"`
	BoardWidth  int      `json:"board_width"`
	BoardHeight int      `json:"board_height"`
}

type Player struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type Round struct {
	Turn  int    `json:"turn"`
	Tiles []Tile `json:"tiles"`
}

type Tile struct {
	PipsA       int    `json:"pips_a"`
	PipsB       int    `json:"pips_b"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Orientation string `json:"orientation"`
	PlayerName  string `json:"player_name"`
}
