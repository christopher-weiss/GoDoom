package main

import (
	"fmt"
	"github.com/christopher-weiss/GoDoom/engine"
)

func main() {
	engine.LoadWadFile("resources/doom1.wad")
	fmt.Println(engine.ReadLumpData("E1M1"))
}
