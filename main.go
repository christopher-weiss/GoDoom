package main

import (
	"fmt"
	"github.com/christopher-weiss/GoDoom/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
)

const (
	ScreenWidth  int = 640
	ScreenHeight int = 480
)

var things []engine.Things
var xOff int16 = 0
var yOff int16 = 0

func main() {
	engine.LoadWadFile("resources/doom1.wad")
	things = engine.ReadMapData("E1M1").Things
	for index, thing := range things {
		fmt.Println(fmt.Sprintf("%d x: %d y: %d, dir: %d, type: %d, flags: %d", index, thing.XPosition, thing.YPosition, thing.Direction, thing.ThingType, thing.Flags))
	}

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Go Doom")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}

type Game struct{}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		yOff = yOff + 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		yOff = yOff - 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		xOff = xOff - 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		xOff = xOff + 2
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Loading WAD file ...")
	for _, thing := range things {
		x := (thing.XPosition / 10) + xOff
		y := -((thing.YPosition / 10) + yOff)
		vector.DrawFilledCircle(screen, float32(x), float32(y), 2.0, color.RGBA{R: 128, G: 128, B: 128, A: 128}, false)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
