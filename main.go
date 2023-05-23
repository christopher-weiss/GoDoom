package main

import (
	"fmt"
	"github.com/christopher-weiss/GoDoom/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

const (
	ScreenWidth  int = 640
	ScreenHeight int = 480
)

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Go Doom")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

	engine.LoadWadFile("resources/doom1.wad")
	things := engine.ReadMapData("E1M1")
	for index, thing := range things.Things {
		fmt.Println(fmt.Sprintf("%d x: %d y: %d, dir: %d, type: %d, flags: %d", index, thing.XPosition, thing.YPosition, thing.Direction, thing.ThingType, thing.Flags))
	}
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Loading WAD file ...")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
