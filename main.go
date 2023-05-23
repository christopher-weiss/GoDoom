package main

import (
	"fmt"
	"github.com/christopher-weiss/GoDoom/engine"
)

func main() {
	fmt.Println("Loading Wad ...")
	engine.LoadWadFile("resources/doom1.wad")
	things := engine.ReadMapData("E1M1")
	for index, thing := range things.Things {
		fmt.Println(fmt.Sprintf("%d x: %d y: %d, dir: %d, type: %d, flags: %d", index, thing.XPosition, thing.YPosition, thing.Direction, thing.ThingType, thing.Flags))

	}
}
