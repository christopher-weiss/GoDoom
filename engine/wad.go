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

type WadHeader struct {
	// 4 character identification, either 'IWAD' or 'PWAD'
	identification string

	// integer specifying the number of lumps (files) in the WAD
	numLumps int32

	// integer holding a pointer to the location of the directory.
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
	for index+16 < int32(len(wad)) {
		name := string(wad[index+8 : index+16])
		directories[name] = Directory{filepos: readInt32(wad[index : index+4]), size: readInt32(wad[index+4 : index+8]), name: name}
		index += 16
	}

	lumpData = wad[12:wadHeader.offFat]

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

func ReadLumpData(name string) []byte {
	directory := directories[name]
	return lumpData[directory.filepos : directory.filepos+directory.size]
}

func readInt32(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func isDebugModeEnabled() bool {
	return len(os.Args) > 1 && os.Args[1] == "debug"
}
