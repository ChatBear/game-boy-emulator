package main

import (
	"fmt"
	"go_emu/config"
	"go_emu/cpu"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	cpu *cpu.CPU
}

func main() {
	start := time.Now()
	cpu, err := cpu.NewCPU(0, 0, 0, 0, 0, 0, 0, 0)
	if err != nil {
		fmt.Printf("Error creating CPU: %v\n", err)
		return
	}

	romPath := "rom.gb"
	bytes, err := os.ReadFile(romPath)
	if err != nil {
		fmt.Printf("Error reading ROM: %v\n", err)
		return
	}

	hexas := make([]int, len(bytes))
	for i, b := range bytes {
		hexas[i] = int(b) & 0xFF
	}

	cpu.UploadROM(hexas)
	cpu.InitializeRegisterValues()
	fmt.Print("-----------------------------------------------------------------\n")
	cpu.Boot()
	end := time.Since(start)
	fmt.Printf("----------------------------- %v ----------------------------------------------\n", end)

	ebiten.SetWindowSize(config.ScreenW, config.ScreenH)
	ebiten.SetWindowTitle("Game Boy")
	ebiten.SetWindowResizingMode(1)
	if err := ebiten.RunGame(&Game{cpu}); err != nil {
		log.Fatal(err)
	}

}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	screen.WritePixels(g.cpu.Screen)
}
func (g *Game) Layout(w, h int) (int, int) {
	return config.ScreenW, config.ScreenH
}
func (g *Game) Update() error {
	return nil
}
