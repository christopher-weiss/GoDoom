package main

import (
	"fmt"
	"github.com/christopher-weiss/GoDoom/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math"
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
	sinA := math.Sin(engine.DegToRad(engine.PlayerAngle))
	cosA := math.Cos(engine.DegToRad(engine.PlayerAngle))
	speedSin := engine.PlayerMovementSpeed * sinA
	speedCos := engine.PlayerMovementSpeed * cosA
	dx := 0.0
	dy := 0.0

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
		dx += -speedSin
		dy += speedCos

	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dx += speedSin
		dy += -speedCos
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx += -speedCos
		dy += -speedSin
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx += speedCos
		dy += speedSin
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		engine.PlayerAngle += engine.PlayerRotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		engine.PlayerAngle -= engine.PlayerRotationSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyB) {
		engine.DrawBoundingBoxesInMap = !engine.DrawBoundingBoxesInMap
	}

	// speed correction for both x and y movement
	if dx != 0 && dy != 0 {
		dx *= 1 / math.Sqrt(2)
		dy *= 1 / math.Sqrt(2)
	}
	engine.PlayerOffsetX += float32(dx)
	engine.PlayerOffsetY += float32(dy)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	engine.DrawMap(screen, &currentMap)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return engine.ScreenResX, engine.ScreenRexY
}
