package interfaces

type ISurface interface {
	GetSize() (int, int)

	MoveCursor(rows int, column int)

	Draw(rows int, column int, c rune)

	DrawColor(rows int, column int, c rune, fg ColorDef, bg ColorDef, mode ColorMode)

	DrawText(rows int, column int, c string)

	DrawTextColor(rows int, column int, c string, fg ColorDef, bg ColorDef, mode ColorMode)

	DrawSeries(data []float64, w int, h int, min float64, max float64)

	Begin()

	End()
}
