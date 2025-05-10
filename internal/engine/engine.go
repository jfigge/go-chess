package engine

type Engine struct {
	position *Position
}

func NewEngine() *Engine {
	e := &Engine{
		position: NewPosition(),
	}
	return e
}

func (e *Engine) Setup(fen string) {
	e.position.SetupBoard(fen)
}

func (e *Engine) GetBoards() []uint64 {
	boards := make([]uint64, 8, 8)
	for i := range 8 {
		boards[i] = e.position.bitboards[i]
	}
	return boards
}
