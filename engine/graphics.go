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

	// debug
	for _, node := range currentMap.Nodes {
		//if i == len(currentMap.Nodes)-1 {
		DrawBoundingBoxes(screen, node, offsetX, offsetY)
		//}
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

// DrawNode Only for debugging
func DrawNode(node Node) {

}

func DrawBoundingBoxes(screen *ebiten.Image, node Node, offsetX float32, offsetY float32) {
	leftBoundingBoxRight := remapX(int16((node.leftBoundingBox>>48)&0xffff), offsetX)
	leftBoundingBoxLeft := remapX(int16((node.leftBoundingBox>>32)&0xffff), offsetX)
	leftBoundingBoxBottom := remapY(int16((node.leftBoundingBox>>16)&0xffff), offsetY)
	leftBoundingBoxTop := remapY(int16(node.leftBoundingBox&0xffff), offsetY)
	vector.StrokeRect(screen, leftBoundingBoxLeft, leftBoundingBoxBottom, leftBoundingBoxRight-leftBoundingBoxLeft, leftBoundingBoxTop-leftBoundingBoxBottom, 1.0, color.RGBA{R: 128, A: 128}, true)

	rightBoundingBoxRight := remapX(int16((node.rightBoundingBox>>48)&0xffff), offsetX)
	rightBoundingBoxLeft := remapX(int16((node.rightBoundingBox>>32)&0xffff), offsetX)
	rightBoundingBoxBottom := remapY(int16((node.rightBoundingBox>>16)&0xffff), offsetY)
	rightBoundingBoxTop := remapY(int16(node.rightBoundingBox&0xffff), offsetY)
	vector.StrokeRect(screen, rightBoundingBoxLeft, rightBoundingBoxBottom, rightBoundingBoxRight-rightBoundingBoxLeft, rightBoundingBoxTop-rightBoundingBoxBottom, 1.0, color.RGBA{G: 128, A: 128}, true)
}

func highByte(rVal int64) uint8 {
	return uint8(rVal >> 8)
}

func lowByte(rVal uint16) uint8 {
	return uint8(rVal & 0xff)
}
