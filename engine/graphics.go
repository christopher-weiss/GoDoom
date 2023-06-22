package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

const (
	NativeResX      = 320
	NativeResY      = 200
	ScaleFactor     = 5
	ScreenResX      = NativeResX * ScaleFactor
	ScreenRexY      = NativeResY * ScaleFactor
	ScreenCenterX   = ScreenResX / 2
	ScreenCenterY   = ScreenRexY / 2
	FieldOfView     = 90
	HalfFieldOfView = FieldOfView / 2
)

var DrawBoundingBoxesInMap bool = false
var offsetX float32 = 0
var offsetY float32 = 0

func DrawMap(screen *ebiten.Image, currentMap *Map) {
	calculateMapOffset(&currentMap.Things[0])
	drawPlayer(screen)
	drawThings(screen, &currentMap.Things)
	drawLineDefs(screen, &currentMap.Linedefs, &currentMap.Vertexes)
	drawNodeBoundingBoxes(screen, &currentMap.Nodes) //TODO remove once debug no longer necessary
	drawBspTraversal(screen, currentMap)             //TODO remove once debug no longer necessary
	drawFov(screen)
}

func drawFov(screen *ebiten.Image) {
	fovLen := float64(100)
	sinAlpha := math.Sin(DegToRad(PlayerAngle - float64(HalfFieldOfView)))
	cosAlpha := math.Cos(DegToRad(PlayerAngle - float64(HalfFieldOfView)))
	sinBeta := math.Sin(DegToRad(PlayerAngle + float64(HalfFieldOfView)))
	cosBeta := math.Cos(DegToRad(PlayerAngle + float64(HalfFieldOfView)))

	playerX := float64(800)
	playerY := float64(500)
	x1 := float32(playerX + fovLen*sinAlpha)
	y1 := float32(playerY + fovLen*cosAlpha)
	x2 := float32(playerX + fovLen*sinBeta)
	y2 := float32(playerY + fovLen*cosBeta)
	vector.StrokeLine(screen, float32(playerX), float32(playerY), x1, y1, 1, color.RGBA{R: 128, G: 128, A: 128}, true)
	vector.StrokeLine(screen, float32(playerX), float32(playerY), x2, y2, 1, color.RGBA{R: 128, G: 128, A: 128}, true)
}

func DegToRad(angle float64) float64 {
	return angle * (math.Pi / 180)
}

func RadToDeg(angle float64) float64 {
	return angle * (180 / math.Pi)
}

func drawBspTraversal(screen *ebiten.Image, currentMap *Map) {
	Traverse(int16(len(currentMap.Nodes)-1), currentMap, 800, 500, screen)
}

func drawNodeBoundingBoxes(screen *ebiten.Image, nodes *[]Node) {
	if DrawBoundingBoxesInMap {
		for _, node := range *nodes {
			drawBoundingBoxes(screen, node)
		}
	}
}

func drawLineDefs(screen *ebiten.Image, linedefs *[]Linedef, vertexes *[]Vertex) {
	for _, linedef := range *linedefs {
		x1 := remapX((*vertexes)[linedef.StartVertex].XPosition)
		y1 := remapY((*vertexes)[linedef.StartVertex].YPosition)
		x2 := remapX((*vertexes)[linedef.EndVertex].XPosition)
		y2 := remapY((*vertexes)[linedef.EndVertex].YPosition)
		vector.StrokeLine(screen, x1, y1, x2, y2, 2.0, color.RGBA{R: 128, G: 128, B: 128, A: 128}, true)
	}
}

func drawThings(screen *ebiten.Image, things *[]Thing) {
	for _, thing := range *things {
		x := remapX(thing.XPosition)
		y := remapY(thing.YPosition)
		vector.DrawFilledCircle(screen, x, y, 2.0, color.RGBA{R: 128, G: 128, B: 128, A: 128}, true)
	}
}

func drawPlayer(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(ScreenCenterX), float32(ScreenCenterY), 4.0, color.RGBA{R: 128, A: 128}, true)
}

// DrawBoundingBoxes draws the bounding boxes (left/right) of a given node
func drawBoundingBoxes(screen *ebiten.Image, node Node) {
	drawBoundingBox(screen, node.leftBoundingBox, color.RGBA{R: 128, A: 128})
	drawBoundingBox(screen, node.rightBoundingBox, color.RGBA{R: 128, A: 128})
}

func drawBoundingBox(screen *ebiten.Image, data int64, color color.RGBA) {
	boundingBox := ConvertToBoundingBox(data)
	vector.StrokeRect(screen, boundingBox.left, boundingBox.bottom, boundingBox.right-boundingBox.left, boundingBox.top-boundingBox.bottom, 1.0, color, true)
}

func DrawSubSector(screen *ebiten.Image, subSector SubSector, segs []Seg, vertexes []Vertex, depthColor uint8) {
	for i := 0; i < int(subSector.segCount); i++ {
		seg := segs[int(subSector.firstSegNumber)+i]
		v1 := vertexes[seg.startingVertexNumber]
		v2 := vertexes[seg.endingVertexNumber]
		vector.StrokeLine(screen, remapX(v1.XPosition), remapY(v1.YPosition), remapX(v2.XPosition), remapY(v2.YPosition), 3, color.RGBA{R: depthColor, A: 128}, true)
	}
}

func IsPlayerLeftOfSplitter(playerX float32, playerY float32, node Node, offsetX float32, offsetY float32) bool {
	partitionLineXRemapped := remapX(node.partitionLineX)
	partitionLineYRemapped := remapY(node.partitionLineY)
	partitionLineDxRemapped := float32(node.dxPartitionLineX)
	partitionLineDyRemapped := float32(-node.dyPartitionLineY)

	dx := playerX - partitionLineXRemapped
	dy := playerY - partitionLineYRemapped

	return (dx*partitionLineDyRemapped)-(dy*partitionLineDxRemapped) <= 0
}

// Remap WAD X-coordinate match resolution and make more of the map visible
func remapX(x int16) float32 {
	return float32(x*int16(ScaleFactor)/20) + offsetX
}

// Remap WAD Y-coordinate match resolution, make more of the map visible and invert (in WAD: positive-y values mean up,
// not down).
func remapY(y int16) float32 {
	return float32(-y*int16(ScaleFactor)/20) - offsetY
}

func calculateMapOffset(player *Thing) {
	offsetX = float32(ScreenCenterX) - float32(player.XPosition*int16(ScaleFactor)/20) + PlayerOffsetX
	offsetY = -(float32(ScreenCenterY) - float32(-player.YPosition*int16(ScaleFactor)/20)) + PlayerOffsetY
}
