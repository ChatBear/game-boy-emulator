package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
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
	// Check if the attributes are all 8-bit unsigned
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
	// TODO : CODE A OPTIMISER PARCE QU ALLER RETOUR SUR STRING -> BINAIRE PAS BIEN
	// iterator := 0
	// for i := 0x0104; i <= 0x011b; i++ {
	// 	binaries := fmt.Sprintf("%02X", cpu.memory[i])
	// 	iterator++
	// 	hexaData[(i-0x0104)/2][(i-0x0104)/2] = string(binaries[0])
	// 	hexaData[(i-0x0104)/2+1][(i-0x0104)/2] = string(binaries[1])

	// }
	for i := 0; i <= 0x011b-0x0104; i++ {
		fmt.Print(i)
		fmt.Print("\n")
		fmt.Print("-----------------------------------------------------------------\n")
		binaries := fmt.Sprintf("%02X", cpu.memory[i+0x0104])
		// iterator++
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
		// iterator++
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
	for i := 0; i <= 7; i++ {
		for _, value := range hexaData[i] {
			val, _ := strconv.ParseUint(value, 16, 16)
			fmt.Printf("%04b", val)
			fmt.Print(" ")
		}

		if i == 3 {
			fmt.Print("\n")
			fmt.Print("-----------------------------------------------------------------\n")
			fmt.Print("\n")
		}
	}
	// var pixel []uint8
	var nibble string
	var nintendoScreen [][]uint8
	for i := 0; i < len(hexaData); i++ {
		for y := 0; y < len(hexaData[i]); y++ {
			val, _ := strconv.ParseUint(hexaData[i][y], 16, 16)
			nibble = fmt.Sprintf("%04b", val)

			for _, value := range nibble {
				if value == '0' {
					nintendoScreen[i+y] = []uint8{0, 0, 0, 0xff}
				}
				if value == '1' {
					nintendoScreen[i+y] = []uint8{0xff, 0xff, 0xff, 0xff}
				}
			}
		}
	}

	nintendoScreenFlat := make([]uint8, 8*12*4)
	fmt.Print(hexaData)
	fmt.Print("\n")

	for _, row := range nintendoScreen {
		nintendoScreenFlat = append(nintendoScreenFlat, row...)
	}

	for index, value := range nintendoScreenFlat {
		cpu.screen[index] = value
	}
	// for i := 0; i < len(hexaData); i++ {
	// 	for y := 0; y < len(hexaData[i]); y++ {
	// 		val, _ := strconv.ParseUint(hexaData[i][y], 16, 16)
	// 		nibble = fmt.Sprintf("%04b", val)
	// 		fmt.Printf("Nibble: %v", nibble)
	// 		fmt.Print("\n")
	// 		current_index := 4*i*ScreenW + y*4
	// 		for _, value := range nibble {
	// 			if value == '0' {
	// 				cpu.screen[current_index] = 0
	// 				cpu.screen[current_index+1] = 0
	// 				cpu.screen[current_index+2] = 0
	// 				cpu.screen[current_index+3] = 0xFF
	// 			}
	// 			if value == '1' {
	// 				cpu.screen[current_index] = 0xFF
	// 				cpu.screen[current_index+1] = 0xFF
	// 				cpu.screen[current_index+2] = 0xFF
	// 				cpu.screen[current_index+3] = 0xFF
	// 			}
	// 			current_index += 16
	// 		}

	// 	}
	// }

	fmt.Print("-----------------------------------------------------------------\n")
	fmt.Print("\n")
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
