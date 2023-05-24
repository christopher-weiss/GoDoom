package main

import (
	"fmt"
	"github.com/christopher-weiss/GoDoom/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
)

const (
	ScreenWidth  int = 640
	ScreenHeight int = 480
)

var xOff int16 = 0
var yOff int16 = 0
var mapData = make(map[string]engine.Map)
var currentMap engine.Map

func main() {
	engine.LoadWadFile("resources/doom1.wad")
	mapName := "E1M1"
	for level := 1; level <= 6; level++ {
		levelName := fmt.Sprintf("E1M%d", level)
		mapData[levelName] = engine.ReadMapData(levelName)
	}
	currentMap = mapData["E1M1"]

	for index, thing := range mapData[mapName].Things {
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
	if ebiten.IsKeyPressed(ebiten.Key1) {
		currentMap = mapData["E1M1"]
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		currentMap = mapData["E1M2"]
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		currentMap = mapData["E1M3"]
	}
	if ebiten.IsKeyPressed(ebiten.Key4) {
		currentMap = mapData["E1M4"]
	}
	if ebiten.IsKeyPressed(ebiten.Key5) {
		currentMap = mapData["E1M5"]
	}
	if ebiten.IsKeyPressed(ebiten.Key6) {
		currentMap = mapData["E1M6"]
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		yOff = yOff + 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		yOff = yOff - 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		xOff = xOff - 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		xOff = xOff + 4
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, thing := range currentMap.Things {
		x := (thing.XPosition / 10) + xOff
		y := -((thing.YPosition / 10) + yOff)
		vector.DrawFilledCircle(screen, float32(x), float32(y), 2.0, color.RGBA{R: 128, G: 128, B: 128, A: 128}, true)
	}

	for _, linedef := range currentMap.Linedefs {
		x1 := (currentMap.Vertexes[linedef.StartVertex].XPosition / 10) + xOff
		y1 := -((currentMap.Vertexes[linedef.StartVertex].YPosition / 10) + yOff)
		x2 := (currentMap.Vertexes[linedef.EndVertex].XPosition / 10) + xOff
		y2 := -((currentMap.Vertexes[linedef.EndVertex].YPosition / 10) + yOff)
		vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 2.0, color.RGBA{R: 128, G: 128, B: 128, A: 128}, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
