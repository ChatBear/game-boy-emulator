package main

import (
	"fmt"
	"os"
	"time"
)

type CPU struct {
	a, b, c, d, e, f, h, l uint8
	af, bc, de, hl         uint16
	cycle                  int
	programCounter         int
	stackPointer           int
	scx, scy               int
	memory                 []int
}

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

	var hexaData [4][12]string
	// TODO : CODE A OPTIMISER PARCE QU ALLER RETOUR SUR STRING -> BINAIRE PAS BIEN
	for i := 0x0104; i <= 0x011B; i++ {
		binaries := fmt.Sprintf("%02X", cpu.memory[i])
		hexaData[i%2][(i-0x0104)/2] = string(binaries[0])
		hexaData[i%2+1][(i-0x0104)/2] = string(binaries[1])
	}
	fmt.Print("\n")
	for i := 0; i <= 3; i++ {
		fmt.Print(hexaData[i])
		fmt.Print("\n")
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
}
