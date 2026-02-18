package cpu

import (
	"fmt"
	"go_emu/apu"
	"go_emu/config"
)

type CPU struct {
	a, b, c, d, e, f, h, l uint8
	af, bc, de, hl         uint16
	cycle                  int
	programCounter         int
	stackPointer           int
	scx, scy               int
	memory                 []uint8
	Screen                 []uint8
	apu                    *apu.APU
}

var Palette = [4][3]uint8{
	{0xE0, 0xF8, 0xD0}, // 0 - lightest
	{0x88, 0xC0, 0x70}, // 1
	{0x34, 0x68, 0x56}, // 2
	{0x08, 0x18, 0x20}, // 3 - darkest
}

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
		memory: make([]uint8, 0xFFFF),
		Screen: make([]uint8, 4*config.ScreenW*config.ScreenH),
	}, nil
}

// TODO: Need to add a banking transition system on the memory not done yet
// Look for MBC1 and MBC2 in the page 13
func (cpu *CPU) UploadROM(rom []int) {
	fmt.Println("Writing the first 32Kb on the Memory")
	for i := 0; i < 0x8000 && i < len(rom); i++ {
		cpu.memory[i] = rom[i]
	}
	fmt.Println("Done")
}

func (cpu *CPU) Boot() {
	for i := 0x0104; i <= 0x011B; i++ {
		fmt.Printf("%02X", cpu.memory[i])
		fmt.Printf(" : %v", cpu.memory[i])
		fmt.Print("\n")
	}
	var hexaData [8][12]string

	for i := 0; i <= 0x011b-0x0104; i++ {
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
	nintendoScreenData := []byte{}
	for _, row := range hexaData {
		for _, hexChar := range row {
			// Convert hex character to number
			var value int
			fmt.Sscanf(hexChar, "%x", &value)
			for bit := 3; bit >= 0; bit-- {
				if (value & (1 << bit)) != 0 {
					nintendoScreenData = append(nintendoScreenData, 0xFF, 0xFF, 0xFF, 0xFF)
				} else {
					nintendoScreenData = append(nintendoScreenData, 0x00, 0x00, 0x00, 0xFF)
				}
			}
		}
	}
	var offset = 4*config.ScreenW*72 + 280
	var multi_48 = 0
	for index, value := range nintendoScreenData {
		if index != 0 && index%192 == 0 {
			multi_48 += 1
		}
		cpu.Screen[4*multi_48*config.ScreenW+index%192+offset] = value
	}
}

func (cpu *CPU) InitializeRegisterValues() {
	fmt.Print("-----------------------------------------------------------------\n")
	cpu.stackPointer = 0xFFFE
	cpu.programCounter = 0
	cpu.cycle = 0
	fmt.Print("  \nEnd of initialization\n")
}

func (cpu *CPU) opCodes(code uint16, value uint8) {
	switch code {
	case 0x06:
		cpu.b = value
		cpu.cycle += 8
	case 0x0E:
		cpu.c = value
		cpu.cycle += 8
	case 0x16:
		cpu.d = value
		cpu.cycle += 8
	case 0x1E:
		cpu.e = value
		cpu.cycle += 8
	case 0x26:
		cpu.h = value
		cpu.cycle += 8
	case 0x2E:
		cpu.l = value
		cpu.cycle += 8
	case 0x7E:
		cpu.a = cpu.memory[cpu.hl]
		cpu.cycle += 8
	case 0x7F:
		cpu.cycle += 8
	case 0x78:
		cpu.a = cpu.b
	case 0x79:
		cpu.a = cpu.c
		cpu.cycle += 4
	case 0x7A:
		cpu.a = cpu.d
		cpu.cycle += 4
	case 0x7B:
		cpu.a = cpu.e
		cpu.cycle += 4
	case 0x7C:
		cpu.a = cpu.h
		cpu.cycle += 4
	case 0x7D:
		cpu.a = cpu.l
		cpu.cycle += 4
	case 0x40:
		cpu.cycle += 4
	case 0x41:
		cpu.b = cpu.c
		cpu.cycle += 4
	case 0x42:
		cpu.b = cpu.d
		cpu.cycle += 4
	case 0x43:
		cpu.b = cpu.e
		cpu.cycle += 4
	case 0x44:
		cpu.b = cpu.h
		cpu.cycle += 4
	case 0x45:
		cpu.b = cpu.l
		cpu.cycle += 4
	case 0x46:
		cpu.b = cpu.memory[cpu.hl]
		cpu.cycle += 8
	case 0x48:
		cpu.c = cpu.b
		cpu.cycle += 4
	case 0x49:
		cpu.cycle += 4
	case 0x4A:
		cpu.c = cpu.d
		cpu.cycle += 4
	case 0x4B:
		cpu.c = cpu.e
		cpu.cycle += 4
	case 0x4C:
		cpu.c = cpu.h
		cpu.cycle += 4
	case 0x4D:
		cpu.c = cpu.l
		cpu.cycle += 4
	case 0x4E:
		cpu.c = cpu.memory[cpu.hl]
		cpu.cycle += 8
	case 0x50:
		cpu.d = cpu.b
		cpu.cycle += 4
	case 0x51:
		cpu.d = cpu.c
		cpu.cycle += 4
	case 0x52:
		cpu.cycle += 4
	case 0x53:
		cpu.d = cpu.e
		cpu.cycle += 4
	case 0x54:
		cpu.d = cpu.h
		cpu.cycle += 4
	case 0x55:
		cpu.d = cpu.l
		cpu.cycle += 4
	case 0x56:
		cpu.d = cpu.memory[cpu.hl]
		cpu.cycle += 8
	case 0x58:
		cpu.e = cpu.b
		cpu.cycle += 4
	case 0x59:
		cpu.e = cpu.c
		cpu.cycle += 4
	case 0x5A:
		cpu.e = cpu.d
		cpu.cycle += 4
	case 0x5B:
		cpu.cycle += 4
	case 0x5C:
		cpu.e = cpu.h
		cpu.cycle += 4
	case 0x5D:
		cpu.e = cpu.l
		cpu.cycle += 4
	case 0x5E:
		cpu.e = cpu.memory[cpu.hl]
		cpu.cycle += 8
	case 0x60:
		cpu.h = cpu.b
		cpu.cycle += 4
	case 0x61:
		cpu.h = cpu.c
		cpu.cycle += 4
	case 0x62:
		cpu.h = cpu.d
		cpu.cycle += 4
	case 0x63:
		cpu.h = cpu.e
		cpu.cycle += 4
	case 0x64:
		cpu.cycle += 4
	case 0x65:
		cpu.h = cpu.l
		cpu.cycle += 4
	case 0x66:
		cpu.h = cpu.memory[cpu.hl]
		cpu.cycle += 8
	case 0x68:
		cpu.l = cpu.b
		cpu.cycle += 4
	case 0x69:
		cpu.l = cpu.c
		cpu.cycle += 4
	case 0x6A:
		cpu.l = cpu.d
		cpu.cycle += 4
	case 0x6B:
		cpu.l = cpu.e
		cpu.cycle += 4
	case 0x6C:
		cpu.l = cpu.h
		cpu.cycle += 4
	case 0x6D:
		cpu.cycle += 4
	case 0x6E:
		cpu.l = cpu.memory[cpu.hl]
		cpu.cycle += 8
	case 0x70:
		cpu.memory[cpu.hl] = cpu.b
		cpu.cycle += 8
	case 0x71:
		cpu.memory[cpu.hl] = cpu.c
		cpu.cycle += 8
	case 0x72:
		cpu.memory[cpu.hl] = cpu.d
		cpu.cycle += 8
	case 0x73:
		cpu.memory[cpu.hl] = cpu.e
		cpu.cycle += 8
	case 0x74:
		cpu.memory[cpu.hl] = cpu.h
		cpu.cycle += 8
	case 0x75:
		cpu.memory[cpu.hl] = cpu.l
		cpu.cycle += 8
	case 0x36:
		cpu.memory[cpu.hl] = value
		cpu.cycle += 12
	}

}

func (cpu *CPU) writeMemory(adress uint16, value uint8) error {
	switch {
	case adress < 0x8000:
		// ROM — read only, ignore writes (or handle MBC later)
		return fmt.Errorf("Multiple Bank Cartridge has not been implemented yet: 0x%04X", adress)
	case adress >= 0x8000 && adress < 0xA000:
		cpu.memory[adress] = value
	case adress >= 0xA000 && adress < 0xC000:
		cpu.memory[adress] = value
	case adress >= 0xC000 && adress < 0xE000:
		cpu.memory[adress] = value
	case adress >= 0xE000 && adress < 0xFE00:
		// Echo RAM — mirrors C000-DDFF
		cpu.memory[adress] = value
		cpu.memory[adress-0x2000] = value
	case adress >= 0xFE00 && adress < 0xFEA0:
		cpu.memory[adress] = value
	case adress >= 0xFEA0 && adress < 0xFF00:
		return fmt.Errorf("Write to unusable memory: 0x%04X\n", adress)
	case adress >= 0xFF00 && adress < 0xFF4C:
		// I/O mapping — not implemented yet
		return fmt.Errorf("I/O mapping has not been implemented yet: 0x%04X", adress)
	case adress >= 0xFF4C && adress < 0xFF80:
		fmt.Errorf("Write to unusable memory: 0x%04X\n", adress)
	case adress >= 0xFF80 && adress < 0xFFFF:
		cpu.memory[adress] = value
	case adress == 0xFFFF:
		cpu.memory[adress] = value
	default:
		return fmt.Errorf("Unknown adress cannot write into the memory: 0x%04X", adress)
	}
	return nil
}
