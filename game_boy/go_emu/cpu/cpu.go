package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type CPU struct {
	a, b, c, d, e, f, h, l uint8
	af, bc, de, hl         uint16
	cycle                  int
	programCounter         int
	stackPointer           int
	scx, scy               int
	memory                 []int
	screen                 []uint8
}

var Palette = [4][3]uint8{
	{0xE0, 0xF8, 0xD0}, // 0 - lightest
	{0x88, 0xC0, 0x70}, // 1
	{0x34, 0x68, 0x56}, // 2
	{0x08, 0x18, 0x20}, // 3 - darkest
}

const (
	ScreenW = 160
	ScreenH = 144
	Scale   = 1
)

func NewCPU(a, b, c, d, e, f, h, l uint8) (*CPU, error) {
	return &CPU{
		a:      a,
		b:      b,
		c:      c,
		d:      d,
		e:      e,
		f:      f,
		h:      h,
		l:      l,
		memory: make([]int, 0xFFFF),
		screen: make([]uint8, 4*ScreenW*ScreenH),
	}, nil
}

// TODO: Need to add a banking transition system on the memory not done yet
// Look for MBC1 and MBC2 in the page 13
func (cpu *CPU) uploadROM(rom []int) {
	fmt.Println("Writing the first 32Kb on the Memory")
	for i := 0; i < 0x8000 && i < len(rom); i++ {
		cpu.memory[i] = rom[i]
	}
	fmt.Println("Done")
}

func (cpu *CPU) OpCodes(code int) {
	if code == 0x0031 {
		// Implementation here
	}
}

func (cpu *CPU) boot() {
	for i := 0x0104; i <= 0x011B; i++ {
		fmt.Printf("%02X", cpu.memory[i])
		fmt.Printf(" : %v", cpu.memory[i])
		fmt.Print("\n")
	}
	fmt.Print("\n")
	fmt.Print("-----------------------------------------------------------------\n")
	var hexaData [8][12]string

	for i := 0; i <= 0x011b-0x0104; i++ {
		fmt.Print(i)
		fmt.Print("\n")
		fmt.Print("-----------------------------------------------------------------\n")
		binaries := fmt.Sprintf("%02X", cpu.memory[i+0x0104])
		if i%2 == 0 {
			hexaData[0][i/2] = string(binaries[0])
			hexaData[1][i/2] = string(binaries[1])
		} else {
			hexaData[2][i/2] = string(binaries[0])
			hexaData[3][i/2] = string(binaries[1])
		}
	}

	for i := 0; i <= 0x011b-0x0104; i++ {
		fmt.Print(i)
		fmt.Print("\n")
		binaries := fmt.Sprintf("%02X", cpu.memory[i+0x011b+1])
		if i%2 == 0 {
			hexaData[4][i/2] = string(binaries[0])
			hexaData[5][i/2] = string(binaries[1])
		} else {
			hexaData[6][i/2] = string(binaries[0])
			hexaData[7][i/2] = string(binaries[1])
		}
	}
	fmt.Print("\n")
	fmt.Print("-----------------------------------------------------------------\n")

	nintendoScreenData := []byte{}
	for _, row := range hexaData {
		for _, hexChar := range row {
			// Convert hex character to number
			var value int
			fmt.Sscanf(hexChar, "%x", &value)
			// Extract each bit (from MSB to LSB)
			for bit := 3; bit >= 0; bit-- {
				if (value & (1 << bit)) != 0 {
					nintendoScreenData = append(nintendoScreenData, 0xFF, 0xFF, 0xFF, 0xFF)
				} else {
					nintendoScreenData = append(nintendoScreenData, 0x00, 0x00, 0x00, 0xFF)
				}
			}
		}
	}
	var offset = 10000
	var multi_48 = 0
	for index, value := range nintendoScreenData {
		if index != 0 && index%188 == 0 {
			multi_48 += 1
		}
		cpu.screen[multi_48*ScreenW+index+offset] = value
	}
}

// System of bank switching: Two types of Cartridge: MBC1 and MBC2 (3, 4, 5)
// depending on the size of the game
// It is also named in the header of the card -> in the rom binary

func (cpu *CPU) initialize() {
	fmt.Print("-----------------------------------------------------------------\n")
	cpu.stackPointer = 0xFFFE
	cpu.programCounter = 0
	cpu.cycle = 0
	fmt.Print("  \nEnd of initialization\n")
}

type Game struct {
	cpu *CPU
}

func (g *Game) Update() error {
	return nil
}

func main() {
	start := time.Now()
	cpu, err := NewCPU(0, 0, 0, 0, 0, 0, 0, 0)
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

	cpu.uploadROM(hexas)
	cpu.initialize()
	fmt.Print("-----------------------------------------------------------------\n")
	cpu.boot()
	end := time.Since(start)
	fmt.Printf("----------------------------- %v ----------------------------------------------\n", end)

	ebiten.SetWindowSize(ScreenW*Scale, ScreenH*Scale)
	ebiten.SetWindowTitle("Game Boy")
	ebiten.SetWindowResizingMode(1)
	if err := ebiten.RunGame(&Game{cpu}); err != nil {
		log.Fatal(err)
	}

}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	screen.WritePixels(g.cpu.screen)
}
func (g *Game) Layout(w, h int) (int, int) {
	return ScreenW, ScreenH
}
