package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const (
	ScreenWidth  int = 640
	ScreenHeight int = 480
)

func DrawMap(screen *ebiten.Image, currentMap *Map) {
	// player
	vector.DrawFilledCircle(screen, float32(ScreenWidth)/2.0, float32(ScreenHeight)/2.0, 2.0, color.RGBA{R: 128, A: 128}, true)

	var scalingFactor int16 = 3
	for _, thing := range currentMap.Things {
		x := (thing.XPosition / scalingFactor) + XOffset
		y := -((thing.YPosition / scalingFactor) + YOffset)
		vector.DrawFilledCircle(screen, float32(x), float32(y), 2.0, color.RGBA{R: 128, G: 128, B: 128, A: 128}, true)
	}

	for _, linedef := range currentMap.Linedefs {
		x1 := (currentMap.Vertexes[linedef.StartVertex].XPosition / scalingFactor) + XOffset
		y1 := -((currentMap.Vertexes[linedef.StartVertex].YPosition / scalingFactor) + YOffset)
		x2 := (currentMap.Vertexes[linedef.EndVertex].XPosition / scalingFactor) + XOffset
		y2 := -((currentMap.Vertexes[linedef.EndVertex].YPosition / scalingFactor) + YOffset)
		vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 2.0, color.RGBA{R: 128, G: 128, B: 128, A: 128}, true)
	}
}
