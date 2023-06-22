package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

// Node see: https://doom.fandom.com/wiki/Node
type Node struct {
	id               int16
	partitionLineX   int16
	partitionLineY   int16
	dxPartitionLineX int16
	dyPartitionLineY int16
	rightBoundingBox int64
	leftBoundingBox  int64
	rightChild       int16
	leftChild        int16
}

// Sector see: https://doom.fandom.com/wiki/Sector
type Sector struct {
	floorHeight          int16
	ceilingHeight        int16
	nameOfFloorTexture   string
	nameOfCeilingTexture string
	lightLevel           int16
	sectorType           int16
	tagNumber            int16
}

// SubSector see: https://doom.fandom.com/wiki/Subsector
type SubSector struct {
	segCount       int16
	firstSegNumber int16
}

// Seg see: https://doom.fandom.com/wiki/Seg
type Seg struct {
	startingVertexNumber int16
	endingVertexNumber   int16
	angle                int16
	lineDefNumber        int16
	direction            int16
	offset               int16
}

var depthColor uint8 = 0

func Traverse(nodeId int16, currentMap *Map, x float32, y float32, screen *ebiten.Image) {
	if nodeId < 0 {
		subSectorId := uint16(nodeId) - uint16(0x8000)
		subSector := currentMap.SubSectors[subSectorId]
		//if !DrawBoundingBoxesInMap {
		DrawSubSector(screen, subSector, currentMap.Segs, currentMap.Vertexes, depthColor)
		//	depthColor++
		//time.Sleep(500 * time.Microsecond)
		//}
		return
	}

	node := currentMap.Nodes[nodeId]

	if !IsPlayerLeftOfSplitter(x, y, node, offsetX, offsetY) {
		Traverse(node.leftChild, currentMap, x, y, screen)
		if collidesWithBoundingBox(node.rightBoundingBox) {
			Traverse(node.rightChild, currentMap, x, y, screen)
		}
	} else {
		Traverse(node.rightChild, currentMap, x, y, screen)
		if collidesWithBoundingBox(node.leftBoundingBox) {
			Traverse(node.leftChild, currentMap, x, y, screen)
		}
	}
}

func collidesWithBoundingBox(data int64) bool {
	boundingBox := ConvertToBoundingBox(data)

	// bounding box vertices
	a := Vec2{boundingBox.left, boundingBox.bottom}
	b := Vec2{boundingBox.left, boundingBox.top}
	c := Vec2{boundingBox.right, boundingBox.top}
	d := Vec2{boundingBox.right, boundingBox.bottom}

	var boundingBoxSides [4]Vec2

	if PlayerOffsetX < boundingBox.left {
		if PlayerOffsetY > boundingBox.top {
			boundingBoxSides = [4]Vec2{b, a, c, b}
		} else if PlayerOffsetY < boundingBox.bottom {
			boundingBoxSides = [4]Vec2{b, a, a, d}
		} else {
			boundingBoxSides = [4]Vec2{b, a, {}, {}}
		}
	} else if PlayerOffsetX > boundingBox.right {
		if PlayerOffsetY > boundingBox.top {
			boundingBoxSides = [4]Vec2{c, b, d, c}
		} else if PlayerOffsetY > boundingBox.bottom {
			boundingBoxSides = [4]Vec2{a, d, d, c}
		} else {
			boundingBoxSides = [4]Vec2{d, c, {}, {}}
		}
	} else {
		if PlayerOffsetY > boundingBox.top {
			boundingBoxSides = [4]Vec2{c, b, {}, {}}
		} else if PlayerOffsetY < boundingBox.bottom {
			boundingBoxSides = [4]Vec2{a, d, {}, {}}
		} else {
			panic("Could not determine bounding box collision")
		}
	}

	for i := 0; i < 4; i += 2 {

		if boundingBoxSides[i] == (Vec2{}) {
			continue
		}

		angle1 := angleFromPlayer(boundingBoxSides[i])
		angle2 := angleFromPlayer(boundingBoxSides[i+1])

		span := normalizeAngle(int32(angle1 - angle2))

		angle1 -= PlayerAngle
		span1 := normalizeAngle(int32(angle1) + HalfFieldOfView)

		if span1 > FieldOfView {
			if span1 >= span+FieldOfView {
				continue
			}
		}
		return true
	}
	return false
}

// angleFromPlayer determines the angle between the players position and the given vertex
func angleFromPlayer(vertex Vec2) float64 {
	deltaX := float64(vertex.x - PlayerOffsetX)
	deltaY := float64(vertex.y - PlayerOffsetY)
	return RadToDeg(math.Atan2(deltaY, deltaX))
}

func normalizeAngle(angle int32) int32 {
	angle %= 360
	if angle < 0 {
		return angle + 360
	} else {
		return angle
	}
}

type BoundingBox struct {
	right  float32
	left   float32
	bottom float32
	top    float32
}

type Vec2 struct {
	x float32
	y float32
}
