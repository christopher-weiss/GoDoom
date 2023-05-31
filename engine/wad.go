package engine

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

var wad []byte
var wadHeader WadHeader
var lumpData []byte
var directories = make(map[string]Directory)
var directory = make(map[int]Directory)
var directoryIndex = make(map[string]int)

// Index offsets between Map-marker and Map-objects
const (
	ThingsOffset    int = 1
	LineDefsOffset  int = 2
	SideDefsOffset  int = 3
	VertexesOffset  int = 4
	SegsOffset      int = 5
	SubSectorOffset int = 6
	NodesOffset     int = 7
	SectorsOffset   int = 8
	RejectOffset    int = 9
	BlockmapOffset  int = 10
)

// Block sizes (in bytes)
const (
	ThingsBlockSize    int32 = 10
	DirectoryBlockSize int32 = 16
	VertexesBlockSize  int32 = 4
	LinedefsBlockSize  int32 = 14
	SegsBlockSize      int32 = 12
	SubSectorBlockSize int32 = 4
	NodeBlockSize      int32 = 28
	SectorBlockSize    int32 = 24
)

type WadHeader struct {
	identification string
	numLumps       int32
	offFat         int32
}

type Directory struct {
	filepos int32
	size    int32
	name    string
}

type Map struct {
	// Map Name
	Name       string
	Things     []Thing
	Linedefs   []Linedef
	Vertexes   []Vertex
	Segs       []Seg
	SubSectors []SubSector
	Nodes      []Node
	Sectors    []Sector
}

// Thing (see: https://doomwiki.org/wiki/Thing)
type Thing struct {
	XPosition int16
	YPosition int16
	Direction int16
	ThingType int16
	Flags     int16
}

type Vertex struct {
	XPosition int16
	YPosition int16
}

type Linedef struct {
	StartVertex  int16
	EndVertex    int16
	Flags        int16
	SpecialType  int16
	SectorTag    int16
	FrontSideDef int16
	BackSideDef  int16
}

func LoadWadFile(path string) {
	wad, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Could not read WAD file from path `%s`", path))
	}

	// Header
	wadHeader = WadHeader{
		string(wad[0:4]),
		readInt32(wad[4:8]),
		readInt32(wad[8:12]),
	}

	// Directories
	index := wadHeader.offFat
	i := 0
	for index+DirectoryBlockSize < int32(len(wad)) {
		name := string(bytes.Trim(wad[index+8:index+16], "\x00")) // trim null-terminated strings
		directories[name] = Directory{filepos: readInt32(wad[index : index+4]), size: readInt32(wad[index+4 : index+8]), name: name}
		directoryIndex[name] = i
		directory[i] = directories[name]
		i++
		index += DirectoryBlockSize
	}

	lumpData = wad[0:wadHeader.offFat]

	if isDebugModeEnabled() {
		fmt.Println("--- HEADER ---")
		fmt.Println(fmt.Sprintf("file type: %s", wadHeader.identification))
		fmt.Println(fmt.Sprintf("num lumps: %d", wadHeader.numLumps))
		fmt.Println(fmt.Sprintf("off FAT:   %d", wadHeader.offFat))

		fmt.Println("--- DIRECTORIES ---")
		for _, directory := range directories {
			fmt.Println(directory)
		}
	}
}

func readLumpIndexForName(name string) int {
	return directoryIndex[name]
}

func ReadDirectoryForLumpIndex(index int) Directory {
	return directory[index]
}

func ReadMapData(mapName string) Map {
	lumpIndex := readLumpIndexForName(mapName)
	thingsDirectory := ReadDirectoryForLumpIndex(lumpIndex + ThingsOffset)
	thingsLumpData := ReadLumpData(thingsDirectory)

	// Thing
	var things []Thing
	for entryOffset := int32(0); entryOffset < thingsDirectory.size; entryOffset += ThingsBlockSize {
		things = append(things, Thing{
			XPosition: readInt16(thingsLumpData[0+entryOffset : 2+entryOffset]),
			YPosition: readInt16(thingsLumpData[2+entryOffset : 4+entryOffset]),
			Direction: readInt16(thingsLumpData[4+entryOffset : 6+entryOffset]),
			ThingType: readInt16(thingsLumpData[6+entryOffset : 8+entryOffset]),
			Flags:     readInt16(thingsLumpData[8+entryOffset : 10+entryOffset]),
		})
	}

	// Linedefs
	linedefsDirectory := ReadDirectoryForLumpIndex(lumpIndex + LineDefsOffset)
	linedefsLumpData := ReadLumpData(linedefsDirectory)
	var linedefs []Linedef
	for entryOffset := int32(0); entryOffset < linedefsDirectory.size; entryOffset += LinedefsBlockSize {
		linedefs = append(linedefs, Linedef{
			StartVertex:  readInt16(linedefsLumpData[0+entryOffset : 2+entryOffset]),
			EndVertex:    readInt16(linedefsLumpData[2+entryOffset : 4+entryOffset]),
			Flags:        readInt16(linedefsLumpData[4+entryOffset : 6+entryOffset]),
			SpecialType:  readInt16(linedefsLumpData[6+entryOffset : 8+entryOffset]),
			SectorTag:    readInt16(linedefsLumpData[8+entryOffset : 10+entryOffset]),
			FrontSideDef: readInt16(linedefsLumpData[10+entryOffset : 12+entryOffset]),
			BackSideDef:  readInt16(linedefsLumpData[12+entryOffset : 14+entryOffset]),
		})
	}

	// Vertexes
	vertexesDirectory := ReadDirectoryForLumpIndex(lumpIndex + VertexesOffset)
	vertexesLumpData := ReadLumpData(vertexesDirectory)
	var vertexes []Vertex
	for entryOffset := int32(0); entryOffset < vertexesDirectory.size; entryOffset += VertexesBlockSize {
		vertexes = append(vertexes, Vertex{
			XPosition: readInt16(vertexesLumpData[0+entryOffset : 2+entryOffset]),
			YPosition: readInt16(vertexesLumpData[2+entryOffset : 4+entryOffset]),
		})
	}

	// Segs
	var segs []Seg
	segsDirectory := ReadDirectoryForLumpIndex(lumpIndex + SegsOffset)
	segsLumpData := ReadLumpData(segsDirectory)
	for entryOffset := int32(0); entryOffset < segsDirectory.size; entryOffset += SegsBlockSize {
		segs = append(segs, Seg{
			startingVertexNumber: readInt16(segsLumpData[0+entryOffset : 2+entryOffset]),
			endingVertexNumber:   readInt16(segsLumpData[2+entryOffset : 4+entryOffset]),
			angle:                readInt16(segsLumpData[4+entryOffset : 6+entryOffset]),
			lineDefNumber:        readInt16(segsLumpData[6+entryOffset : 8+entryOffset]),
			direction:            readInt16(segsLumpData[8+entryOffset : 10+entryOffset]),
			offset:               readInt16(segsLumpData[10+entryOffset : 12+entryOffset]),
		})
	}

	// SubSectors
	var subSectors []SubSector
	subSectorDirectory := ReadDirectoryForLumpIndex(lumpIndex + SubSectorOffset)
	subSectorLumpData := ReadLumpData(subSectorDirectory)
	for entryOffset := int32(0); entryOffset < subSectorDirectory.size; entryOffset += SubSectorBlockSize {
		subSectors = append(subSectors, SubSector{
			segCount:       readInt16(subSectorLumpData[0+entryOffset : 2+entryOffset]),
			firstSegNumber: readInt16(subSectorLumpData[2+entryOffset : 4+entryOffset]),
		})
	}

	// Nodes
	var nodes []Node
	nodeDirectory := ReadDirectoryForLumpIndex(lumpIndex + NodesOffset)
	nodeLumpData := ReadLumpData(nodeDirectory)
	index := int16(0)
	for entryOffset := int32(0); entryOffset < nodeDirectory.size; entryOffset += NodeBlockSize {
		nodes = append(nodes, Node{
			id:               index,
			partitionLineX:   readInt16(nodeLumpData[0+entryOffset : 2+entryOffset]),
			partitionLineY:   readInt16(nodeLumpData[2+entryOffset : 4+entryOffset]),
			dxPartitionLineX: readInt16(nodeLumpData[4+entryOffset : 6+entryOffset]),
			dyPartitionLineY: readInt16(nodeLumpData[6+entryOffset : 8+entryOffset]),
			rightBoundingBox: readInt64(nodeLumpData[8+entryOffset : 16+entryOffset]),
			leftBoundingBox:  readInt64(nodeLumpData[16+entryOffset : 24+entryOffset]),
			rightChild:       readInt16(nodeLumpData[24+entryOffset : 26+entryOffset]),
			leftChild:        readInt16(nodeLumpData[26+entryOffset : 28+entryOffset]),
		})
		index += 1
	}

	// Sectors
	var sectors []Sector
	sectorDirectory := ReadDirectoryForLumpIndex(lumpIndex + SectorsOffset)
	sectorLumpData := ReadLumpData(sectorDirectory)
	for entryOffset := int32(0); entryOffset < sectorDirectory.size; entryOffset += SectorBlockSize {
		sectors = append(sectors, Sector{
			floorHeight:          readInt16(sectorLumpData[0+entryOffset : 2+entryOffset]),
			ceilingHeight:        readInt16(sectorLumpData[2+entryOffset : 4+entryOffset]),
			nameOfFloorTexture:   readString(sectorLumpData[4+entryOffset : 12+entryOffset]),
			nameOfCeilingTexture: readString(sectorLumpData[12+entryOffset : 20+entryOffset]),
			lightLevel:           readInt16(sectorLumpData[20+entryOffset : 22+entryOffset]),
			sectorType:           readInt16(sectorLumpData[22+entryOffset : 24+entryOffset]),
			tagNumber:            readInt16(sectorLumpData[24+entryOffset : 26+entryOffset]),
		})
	}

	return Map{
		Name:       mapName,
		Things:     things,
		Vertexes:   vertexes,
		Linedefs:   linedefs,
		Segs:       segs,
		SubSectors: subSectors,
		Nodes:      nodes,
	}
}

func ReadLumpData(directory Directory) []byte {
	return lumpData[directory.filepos : directory.filepos+directory.size]
}

func readInt16(data []byte) (ret int16) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func readInt32(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func readInt64(data []byte) (ret int64) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func readString(data []byte) (ret string) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func isDebugModeEnabled() bool {
	return len(os.Args) > 1 && os.Args[1] == "debug"
}
