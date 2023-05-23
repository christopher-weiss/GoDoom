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
	ThingsOffset   int = 1
	LineDefsOffset int = 2
	SideDefsOffset int = 3
	VertexesOffset int = 4
	SegsOffset     int = 5
	SSectorsOffset int = 6
	NodesOffset    int = 7
	SectorsOffset  int = 8
	RejectOffset   int = 9
	BlockmapOffset int = 10
)

// Block sizes
const (
	ThingsBlockSize    int32 = 10
	DirectoryBlockSize int32 = 16
)

type WadHeader struct {
	// 4 character identification, either 'IWAD' or 'PWAD'
	identification string

	// Integer specifying the number of lumps (files) in the WAD
	numLumps int32

	// Integer holding a pointer to the location of the directory.
	offFat int32
}

type Directory struct {
	// An integer holding a pointer to the start of the lump's data in the file
	filepos int32

	// An integer representing the size of the lump in bytes
	size int32

	// A string defining the lump's name
	name string
}

type Map struct {
	// Map name
	name string
	// Things (see: https://doomwiki.org/wiki/Thing)
	Things []Things
}

type Things struct {
	XPosition int16
	YPosition int16
	Direction int16
	ThingType int16
	Flags     int16
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

func ReadHeader() WadHeader {
	return wadHeader
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

	var things []Things
	for entryOffset := int32(0); entryOffset < thingsDirectory.size; entryOffset += ThingsBlockSize {
		things = append(things, Things{
			XPosition: readInt16(thingsLumpData[0+entryOffset : 2+entryOffset]),
			YPosition: readInt16(thingsLumpData[2+entryOffset : 4+entryOffset]),
			Direction: readInt16(thingsLumpData[4+entryOffset : 6+entryOffset]),
			ThingType: readInt16(thingsLumpData[6+entryOffset : 8+entryOffset]),
			Flags:     readInt16(thingsLumpData[8+entryOffset : 10+entryOffset]),
		})
	}

	return Map{
		name:   mapName,
		Things: things,
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

func isDebugModeEnabled() bool {
	return len(os.Args) > 1 && os.Args[1] == "debug"
}
