package game

type Options func(g *Game)

func OptSquareSize(size int) Options {
	return func(g *Game) {
		g.squareSize = size
	}
}

func OptSheetImageSize(size int) Options {
	return func(g *Game) {
		g.sheetImageSize = size
	}
}

func OptFontHeight(height int) Options {
	return func(g *Game) {
		g.fontHeight = height
	}
}

func OptEnableDebug(enabled bool) Options {
	return func(g *Game) {
		g.debugEnabled = enabled
	}
}

func OptHighlightAttacks(highlightAttacks bool) Options {
	return func(g *Game) {
		g.highlightAttacks = highlightAttacks
	}
}

func OptShowStrength(enabled bool) Options {
	return func(g *Game) {
		g.showStrength = enabled
	}
}

func OptShowFPS(showFPS bool) Options {
	return func(g *Game) {
		g.showFPS = showFPS
	}
}

func OptShowLabels(showLabels bool) Options {
	return func(g *Game) {
		g.showLabels = showLabels
	}
}
