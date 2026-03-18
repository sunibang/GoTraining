package producer

type GameBoard struct {
	Players []Player
}

type Player struct {
	X int
	Y int
}

type Command string

const (
	Forward Command = "F"
	Left    Command = "L"
	Right   Command = "R"
	Shoot   Command = "S"
)

func (b *GameBoard) Play() Command {
	// Implement some magic strategy
	return Forward
}

func (b *GameBoard) SomeOtherUnrelatedFunction() {}
