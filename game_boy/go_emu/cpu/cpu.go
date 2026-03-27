package cpu

import (
	"fmt"
	"go_emu/apu"
	"go_emu/config"
)

type CPU struct {
	a, b, c, d, e, f, h, l uint8
	cycle                  int
	programCounter         uint16
	stackPointer           uint16
	scx, scy               int
	memory                 []uint8
	Screen                 []uint8
	apu                    *apu.APU
	opcodeTable            [256]func(value uint8, value2 uint8)
	opcodeTablePrefixed    [256]func(value uint8, value2 uint8)
	halt                   bool
	stopped                bool
	pendingDisableIME      bool
	ime                    bool
}

var Palette = [4][3]uint8{
	{0xE0, 0xF8, 0xD0},
	{0x88, 0xC0, 0x70},
	{0x34, 0x68, 0x56},
	{0x08, 0x18, 0x20},
}

func NewCPU(a, b, c, d, e, f, h, l uint8) (*CPU, error) {
	cpu := &CPU{
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
	}
	cpu.initOpcodes()
	cpu.InitializeRegisterValues()
	return cpu, nil
}

func (cpu *CPU) getAF() uint16 { return uint16(cpu.a)<<8 | uint16(cpu.f) }
func (cpu *CPU) getBC() uint16 { return uint16(cpu.b)<<8 | uint16(cpu.c) }
func (cpu *CPU) getDE() uint16 { return uint16(cpu.d)<<8 | uint16(cpu.e) }
func (cpu *CPU) getHL() uint16 { return uint16(cpu.h)<<8 | uint16(cpu.l) }

func (cpu *CPU) setAF(v uint16) { cpu.a = uint8(v >> 8); cpu.f = uint8(v & 0xF0) }
func (cpu *CPU) setBC(v uint16) { cpu.b = uint8(v >> 8); cpu.c = uint8(v & 0xFF) }
func (cpu *CPU) setDE(v uint16) { cpu.d = uint8(v >> 8); cpu.e = uint8(v & 0xFF) }
func (cpu *CPU) setHL(v uint16) { cpu.h = uint8(v >> 8); cpu.l = uint8(v & 0xFF) }

func (cpu *CPU) initOpcodes() {
	cpu.initLoadOpcodes()
	cpu.initALUOpcodes()
	cpu.initStackOpcodes()
	cpu.init16BitArthmeticOpCode()
	cpu.initMiscellaneousOpCodes()
	cpu.initRotateShiftOpCode()
	cpu.initBitOpCode()
	cpu.initJumpsOpCodes()
	cpu.initCallsOpCode()
	cpu.initRestartOpCode()
	cpu.initReturnsOpCode()
}

// TODO: Need to add a banking transition system on the memory not done yet
// Look for MBC1 and MBC2 in the page 13
func (cpu *CPU) UploadROM(rom []byte) {
	fmt.Println("Writing the first 32Kb on the Memory")
	for i := 0; i < 0x8000 && i < len(rom); i++ {
		cpu.memory[i] = uint8(rom[i])
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

func (cpu *CPU) opCodes(code uint16, value uint8, value2 uint8) error {
	if code == 0xCB {
		cpu.programCounter++
		cpu.opcodeTablePrefixed[cpu.programCounter](value, value2)
		return nil
	}
	if handler := cpu.opcodeTable[code&0xFF]; handler != nil {
		handler(value, value2)
		return nil
	} else {
		return fmt.Errorf("Unimplemented opcode: 0x%02X\n", code)
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
	// case adress >= 0xFF00 && adress < 0xFF4C:
	// 	// I/O mapping — not implemented yet
	// 	return fmt.Errorf("I/O mapping has not been implemented yet: 0x%04X", adress)
	case adress >= 0xFF00 && adress < 0xFF4C:
		cpu.memory[adress] = value
		if adress == 0xFF02 && value == 0x81 {
			fmt.Printf("%c", cpu.memory[0xFF01]) // print serial output
		}
	case adress >= 0xFF4C && adress < 0xFF80:
		return fmt.Errorf("Write to unusable memory: 0x%04X\n", adress)
	case adress >= 0xFF80 && adress < 0xFFFF:
		cpu.memory[adress] = value
	case adress == 0xFFFF:
		cpu.memory[adress] = value
	default:
		return fmt.Errorf("Unknown adress cannot write into the memory: 0x%04X", adress)
	}
	return nil
}

func (cpu *CPU) Step() error {
	fmt.Print("----------------------------")
	fmt.Print("\n")
	fmt.Print(cpu.programCounter)
	fmt.Print("\n")
	fmt.Print(cpu.stackPointer)
	fmt.Print("\n")
	fmt.Print("----------------------------")
	if cpu.programCounter == 515 {
		fmt.Print("DEBUGGER")
	}
	opcode := cpu.memory[cpu.programCounter]
	cpu.programCounter++
	v1 := cpu.memory[int(cpu.programCounter)]
	v2 := cpu.memory[int(cpu.programCounter+1)]
	return cpu.opCodes(uint16(opcode), v1, v2)
}

func (cpu *CPU) Run(maxCycles int) {
	for cpu.cycle < maxCycles && !cpu.stopped {
		if err := cpu.Step(); err != nil {
			fmt.Println(err)
			return
		}
	}
}
