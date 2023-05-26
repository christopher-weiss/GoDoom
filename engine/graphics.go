package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const (
	NativeResX             = 320
	NativeResY             = 200
	ScaleFactor            = 5
	mapScaleFactor float32 = ScaleFactor / 20
	ScreenResX             = NativeResX * ScaleFactor
	ScreenRexY             = NativeResY * ScaleFactor
	ScreenCenterX          = ScreenResX / 2
	ScreenCenterY          = ScreenRexY / 2
)

func DrawMap(screen *ebiten.Image, currentMap *Map) {
	// player (always centered on screen)
	vector.DrawFilledCircle(screen, float32(ScreenCenterX), float32(ScreenCenterY), 4.0, color.RGBA{R: 128, A: 128}, true)

	var offsetX float32 = 0.0
	var offsetY float32 = 0.0

	for i, thing := range currentMap.Things {

		if i == 0 {
			offsetX, offsetY = CalculateOffset(thing.XPosition, thing.YPosition)
		}
		x := remapX(thing.XPosition, offsetX)
		y := remapY(thing.YPosition, offsetY)
		vector.DrawFilledCircle(screen, x, y, 2.0, color.RGBA{R: 128, G: 128, B: 128, A: 128}, true)
	}

	for _, linedef := range currentMap.Linedefs {
		x1 := remapX(currentMap.Vertexes[linedef.StartVertex].XPosition, offsetX)
		y1 := remapY(currentMap.Vertexes[linedef.StartVertex].YPosition, offsetY)
		x2 := remapX(currentMap.Vertexes[linedef.EndVertex].XPosition, offsetX)
		y2 := remapY(currentMap.Vertexes[linedef.EndVertex].YPosition, offsetY)
		vector.StrokeLine(screen, x1, y1, x2, y2, 2.0, color.RGBA{R: 128, G: 128, B: 128, A: 128}, true)
	}
}

// Remap WAD X-coordinate match resolution and make more of the map visible
func remapX(x int16, offset float32) float32 {
	return float32(x*int16(ScaleFactor)/20) + offset
}

// Remap WAD Y-coordinate match resolution, make more of the map visible and invert (in WAD: positive-y values mean up,
// not down).
func remapY(y int16, offset float32) float32 {
	return float32(-y*int16(ScaleFactor)/20) - offset
}

func CalculateOffset(x int16, y int16) (float32, float32) {
	return float32(ScreenCenterX) - float32(x*int16(ScaleFactor)/20) + PlayerOffsetX, -(float32(ScreenCenterY) - float32(-y*int16(ScaleFactor)/20)) + PlayerOffsetY
}
