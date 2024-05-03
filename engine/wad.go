package engine

import (
	"bytes"
	_ "embed"
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
		readInt[int32](wad[4:8]),
		readInt[int32](wad[8:12]),
	}

	// Directories
	index := wadHeader.offFat
	i := 0
	for index+DirectoryBlockSize < int32(len(wad)) {
		name := string(bytes.Trim(wad[index+8:index+16], "\x00")) // trim null-terminated strings
		directories[name] = Directory{filepos: readInt[int32](wad[index : index+4]), size: readInt[int32](wad[index+4 : index+8]), name: name}
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
			XPosition: readInt[int16](thingsLumpData[0+entryOffset : 2+entryOffset]),
			YPosition: readInt[int16](thingsLumpData[2+entryOffset : 4+entryOffset]),
			Direction: readInt[int16](thingsLumpData[4+entryOffset : 6+entryOffset]),
			ThingType: readInt[int16](thingsLumpData[6+entryOffset : 8+entryOffset]),
			Flags:     readInt[int16](thingsLumpData[8+entryOffset : 10+entryOffset]),
		})
	}

	// Linedefs
	linedefsDirectory := ReadDirectoryForLumpIndex(lumpIndex + LineDefsOffset)
	linedefsLumpData := ReadLumpData(linedefsDirectory)
	var linedefs []Linedef
	for entryOffset := int32(0); entryOffset < linedefsDirectory.size; entryOffset += LinedefsBlockSize {
		linedefs = append(linedefs, Linedef{
			StartVertex:  readInt[int16](linedefsLumpData[0+entryOffset : 2+entryOffset]),
			EndVertex:    readInt[int16](linedefsLumpData[2+entryOffset : 4+entryOffset]),
			Flags:        readInt[int16](linedefsLumpData[4+entryOffset : 6+entryOffset]),
			SpecialType:  readInt[int16](linedefsLumpData[6+entryOffset : 8+entryOffset]),
			SectorTag:    readInt[int16](linedefsLumpData[8+entryOffset : 10+entryOffset]),
			FrontSideDef: readInt[int16](linedefsLumpData[10+entryOffset : 12+entryOffset]),
			BackSideDef:  readInt[int16](linedefsLumpData[12+entryOffset : 14+entryOffset]),
		})
	}

	// Vertexes
	vertexesDirectory := ReadDirectoryForLumpIndex(lumpIndex + VertexesOffset)
	vertexesLumpData := ReadLumpData(vertexesDirectory)
	var vertexes []Vertex
	for entryOffset := int32(0); entryOffset < vertexesDirectory.size; entryOffset += VertexesBlockSize {
		vertexes = append(vertexes, Vertex{
			XPosition: readInt[int16](vertexesLumpData[0+entryOffset : 2+entryOffset]),
			YPosition: readInt[int16](vertexesLumpData[2+entryOffset : 4+entryOffset]),
		})
	}

	// Segs
	var segs []Seg
	segsDirectory := ReadDirectoryForLumpIndex(lumpIndex + SegsOffset)
	segsLumpData := ReadLumpData(segsDirectory)
	for entryOffset := int32(0); entryOffset < segsDirectory.size; entryOffset += SegsBlockSize {
		segs = append(segs, Seg{
			startingVertexNumber: readInt[int16](segsLumpData[0+entryOffset : 2+entryOffset]),
			endingVertexNumber:   readInt[int16](segsLumpData[2+entryOffset : 4+entryOffset]),
			angle:                readInt[int16](segsLumpData[4+entryOffset : 6+entryOffset]),
			lineDefNumber:        readInt[int16](segsLumpData[6+entryOffset : 8+entryOffset]),
			direction:            readInt[int16](segsLumpData[8+entryOffset : 10+entryOffset]),
			offset:               readInt[int16](segsLumpData[10+entryOffset : 12+entryOffset]),
		})
	}

	// SubSectors
	var subSectors []SubSector
	subSectorDirectory := ReadDirectoryForLumpIndex(lumpIndex + SubSectorOffset)
	subSectorLumpData := ReadLumpData(subSectorDirectory)
	for entryOffset := int32(0); entryOffset < subSectorDirectory.size; entryOffset += SubSectorBlockSize {
		subSectors = append(subSectors, SubSector{
			segCount:       readInt[int16](subSectorLumpData[0+entryOffset : 2+entryOffset]),
			firstSegNumber: readInt[int16](subSectorLumpData[2+entryOffset : 4+entryOffset]),
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
			partitionLineX:   readInt[int16](nodeLumpData[0+entryOffset : 2+entryOffset]),
			partitionLineY:   readInt[int16](nodeLumpData[2+entryOffset : 4+entryOffset]),
			dxPartitionLineX: readInt[int16](nodeLumpData[4+entryOffset : 6+entryOffset]),
			dyPartitionLineY: readInt[int16](nodeLumpData[6+entryOffset : 8+entryOffset]),
			rightBoundingBox: readInt[int64](nodeLumpData[8+entryOffset : 16+entryOffset]),
			leftBoundingBox:  readInt[int64](nodeLumpData[16+entryOffset : 24+entryOffset]),
			rightChild:       readInt[int16](nodeLumpData[24+entryOffset : 26+entryOffset]),
			leftChild:        readInt[int16](nodeLumpData[26+entryOffset : 28+entryOffset]),
		})
		index += 1
	}

	// Sectors
	var sectors []Sector
	sectorDirectory := ReadDirectoryForLumpIndex(lumpIndex + SectorsOffset)
	sectorLumpData := ReadLumpData(sectorDirectory)
	for entryOffset := int32(0); entryOffset < sectorDirectory.size; entryOffset += SectorBlockSize {
		sectors = append(sectors, Sector{
			floorHeight:          readInt[int16](sectorLumpData[0+entryOffset : 2+entryOffset]),
			ceilingHeight:        readInt[int16](sectorLumpData[2+entryOffset : 4+entryOffset]),
			nameOfFloorTexture:   readString(sectorLumpData[4+entryOffset : 12+entryOffset]),
			nameOfCeilingTexture: readString(sectorLumpData[12+entryOffset : 20+entryOffset]),
			lightLevel:           readInt[int16](sectorLumpData[20+entryOffset : 22+entryOffset]),
			sectorType:           readInt[int16](sectorLumpData[22+entryOffset : 24+entryOffset]),
			tagNumber:            readInt[int16](sectorLumpData[24+entryOffset : 26+entryOffset]),
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

func readInt[T int16 | int32 | int64](data []byte) T {
	ret := T(0)
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.LittleEndian, &ret)
	if err != nil {
		return T(0)
	}

	return ret
}

func readString(data []byte) (ret string) {
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.LittleEndian, &ret)
	if err != nil {
		return ""
	}
	return
}

func isDebugModeEnabled() bool {
	return len(os.Args) > 1 && os.Args[1] == "debug"
}

// ConvertToBoundingBox converts data given as a 64-bit integer from the WAD file to the BoundingBox data structure.
func ConvertToBoundingBox(data int64) BoundingBox {
	boundingBoxRight := remapX(int16((data >> 48) & 0xffff))
	boundingBoxLeft := remapX(int16((data >> 32) & 0xffff))
	boundingBoxBottom := remapY(int16((data >> 16) & 0xffff))
	boundingBoxTop := remapY(int16(data & 0xffff))

	return BoundingBox{
		right:  boundingBoxRight,
		left:   boundingBoxLeft,
		bottom: boundingBoxBottom,
		top:    boundingBoxTop,
	}
}
