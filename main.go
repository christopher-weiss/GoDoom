package main

import (
	"fmt"
	"github.com/christopher-weiss/GoDoom/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

type Game struct{}

var mapData = make(map[string]engine.Map)
var currentMap engine.Map

func main() {
	initializeGame()

	bTree := engine.NodeDemo{Value: 0}
	bTree.Insert(-15)
	bTree.Insert(-8)
	bTree.Insert(6)
	bTree.Insert(12)
	bTree.Insert(20)

	bTree.Traverse(4)

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func initializeGame() {
	ebiten.SetWindowSize(engine.ScreenWidth, engine.ScreenHeight)
	ebiten.SetWindowTitle("Go Doom")

	engine.LoadWadFile("resources/doom1.wad")
	startingMap := "E1M1"

	for level := 1; level <= 8; level++ {
		levelName := fmt.Sprintf("E1M%d", level)
		mapData[levelName] = engine.ReadMapData(levelName)
	}

	currentMap = mapData[startingMap]
}

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
	if ebiten.IsKeyPressed(ebiten.Key7) {
		currentMap = mapData["E1M7"]
	}
	if ebiten.IsKeyPressed(ebiten.Key8) {
		currentMap = mapData["E1M8"]
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		engine.YOffset = engine.YOffset - 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		engine.YOffset = engine.YOffset + 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		engine.XOffset = engine.XOffset + 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		engine.XOffset = engine.XOffset - 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyEqual) {
		if engine.ScalingFactor > 1 {
			engine.ScalingFactor = engine.ScalingFactor - 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyMinus) {
		engine.ScalingFactor = engine.ScalingFactor + 1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	engine.DrawMap(screen, &currentMap)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return engine.ScreenWidth, engine.ScreenHeight
}
