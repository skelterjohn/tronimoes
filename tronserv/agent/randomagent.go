package main

import "github.com/skelterjohn/tronimoes/tronserv/game"

type RandomAgent struct {
}

func (RandomAgent) Ready() {

}
func (RandomAgent) Update(g *game.Game) {

}
func (RandomAgent) GetMove(g *game.Game, p *game.Player) Move {
	if p.JustDrew {
		return Move{
			Pass: true,
			Selected: Selected{
				X: g.BoardWidth / 2,
				Y: (g.BoardHeight / 2) - 1,
			},
		}
	}
	return Move{
		Draw: true,
	}
}
