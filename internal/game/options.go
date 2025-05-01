package game

type GameOptions func(g *Game)

func OptSquareSize(size uint) GameOptions {
	return func(g *Game) {
		g.squareSize = size
	}
}

func OptSheetImageSize(size int) GameOptions {
	return func(g *Game) {
		g.sheetImageSize = size
	}
}

func OptFontHeight(height int) GameOptions {
	return func(g *Game) {
		g.fontHeight = uint(height)
	}
}
