package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
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
		DrawSubSector(screen, subSector, currentMap.Segs, currentMap.Vertexes, depthColor)
		depthColor++
		time.Sleep(500 * time.Microsecond)
		return
	}

	node := currentMap.Nodes[nodeId]

	if !IsPlayerLeftOfSplitter(x, y, node, offsetX, offsetY) {
		Traverse(node.leftChild, currentMap, x, y, screen)
		Traverse(node.rightChild, currentMap, x, y, screen)
	} else {
		Traverse(node.rightChild, currentMap, x, y, screen)
		Traverse(node.leftChild, currentMap, x, y, screen)
	}
}
