package consumer

import (
	"fmt"
	"github.com/romangurevitch/go-training/internal/basics/interface/producerside/producer"
)

type GameServer struct {
	game producer.Strategy
}

func NewGameServer(s producer.Strategy) *GameServer {
	return &GameServer{
		game: s,
	}
}

func Start() {
	gb := &producer.GameBoard{}
	gs := NewGameServer(gb)
	gs.StartGame()
}

func (g *GameServer) StartGame() {
	c := g.game.Play()
	fmt.Println(c)
}
