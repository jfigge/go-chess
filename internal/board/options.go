package board

type BoardOptions func(b *Board)

func OptSetup(fen string) BoardOptions {
	return func(b *Board) {
		b.fen = fen
	}
}
