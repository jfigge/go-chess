package board

type Options func(b *Board)

func OptSetup(fen string) Options {
	return func(b *Board) {
		b.fen = fen
	}
}
