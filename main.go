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

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func initializeGame() {
	ebiten.SetWindowSize(engine.ScreenResX, engine.ScreenRexY)
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
		engine.PlayerOffsetX = 0
		engine.PlayerOffsetY = 0
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		currentMap = mapData["E1M2"]
		engine.PlayerOffsetX = 0
		engine.PlayerOffsetY = 0
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		currentMap = mapData["E1M3"]
		engine.PlayerOffsetX = 0
		engine.PlayerOffsetY = 0
	}
	if ebiten.IsKeyPressed(ebiten.Key4) {
		currentMap = mapData["E1M4"]
		engine.PlayerOffsetX = 0
		engine.PlayerOffsetY = 0
	}
	if ebiten.IsKeyPressed(ebiten.Key5) {
		currentMap = mapData["E1M5"]
		engine.PlayerOffsetX = 0
		engine.PlayerOffsetY = 0
	}
	if ebiten.IsKeyPressed(ebiten.Key6) {
		currentMap = mapData["E1M6"]
		engine.PlayerOffsetX = 0
		engine.PlayerOffsetY = 0
	}
	if ebiten.IsKeyPressed(ebiten.Key7) {
		currentMap = mapData["E1M7"]
		engine.PlayerOffsetX = 0
		engine.PlayerOffsetY = 0
	}
	if ebiten.IsKeyPressed(ebiten.Key8) {
		currentMap = mapData["E1M8"]
		engine.PlayerOffsetX = 0
		engine.PlayerOffsetY = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		engine.PlayerOffsetY = engine.PlayerOffsetY - float32(1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		engine.PlayerOffsetY = engine.PlayerOffsetY + float32(1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		engine.PlayerOffsetX = engine.PlayerOffsetX + float32(1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		engine.PlayerOffsetX = engine.PlayerOffsetX - float32(1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		engine.PlayerAngle = engine.PlayerAngle + float64(5)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		engine.PlayerAngle = engine.PlayerAngle - float64(5)
	}

	if ebiten.IsKeyPressed(ebiten.KeyB) {
		engine.DrawBoundingBoxesInMap = !engine.DrawBoundingBoxesInMap
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	engine.DrawMap(screen, &currentMap)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return engine.ScreenResX, engine.ScreenRexY
}
