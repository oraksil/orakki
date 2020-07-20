package models

type Game struct {
	Id          int64
	Title       string
	Description string
	Maker       string
	MaxPlayers  int
}

type Player struct {
	Id         int64
	Name       string
	TotalCoins int
}

type RunningGame struct {
	Id        int64
	Game      *Game
	Players   []*Player
	CreatedAt int64
}
