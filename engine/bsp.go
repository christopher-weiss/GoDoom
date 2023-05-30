package engine

import "fmt"

// Node see: https://doom.fandom.com/wiki/Node
type Node struct {
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

type NodeDemo struct {
	Value int
	left  *NodeDemo
	right *NodeDemo
}

func (n *NodeDemo) Insert(value int) {
	if value < n.Value {
		if n.left != nil {
			n.left.Insert(value)
		} else {
			n.left = &NodeDemo{Value: value}
		}
	} else if value > n.Value {
		if n.right != nil {
			n.right.Insert(value)
		} else {
			n.right = &NodeDemo{Value: value}
		}
	}
}

func (n *NodeDemo) Traverse(playerPosition int) {
	if n != nil {
		if playerPosition <= n.Value {
			n.left.Traverse(playerPosition)
			fmt.Println(n.Value)
			n.right.Traverse(playerPosition)
		} else {
			n.right.Traverse(playerPosition)
			fmt.Println(n.Value)
			n.left.Traverse(playerPosition)
		}
	}
}

func GetRootNode(nodes []Node) Node {
	return nodes[len(nodes)-1]
}
