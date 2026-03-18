package consumer

import (
	"fmt"

	"github.com/romangurevitch/go-training/internal/basics/interface/consumerside/producer"
)

type Strategy interface {
	Play() producer.Command
}

type GameServer struct {
	game Strategy
}

func NewGameServer(s Strategy) *GameServer {
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
