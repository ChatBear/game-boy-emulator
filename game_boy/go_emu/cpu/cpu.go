package cpu

import (
	"fmt"
	"go_emu/apu"
	"go_emu/config"
	"strings"
)

type CPU struct {
	a, b, c, d, e, f, h, l uint8
	cycle                  int
	timerAcc               int
	programCounter         uint16
	stackPointer           uint16
	scx, scy               int
	memory                 []uint8
	Screen                 []uint8
	apu                    *apu.APU
	opcodeTable            [256]func()
	opcodeTablePrefixed    [256]func()
	halt                   bool
	stopped                bool
	pendingEnableIME       bool
	ime                    bool
	serialOutput           strings.Builder
}

var nCycles = map[uint8]int{
	0b00: 1024,
	0b01: 16,
	0b10: 64,
	0b11: 256,
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
		memory: make([]uint8, 0x10000),
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

func (cpu *CPU) ie() uint8  { return cpu.memory[0xFFFF] }
func (cpu *CPU) if_() uint8 { return cpu.memory[0xFF0F] }

func (cpu *CPU) tick(tCycle int) {
	cpu.UpdateTimer(tCycle)
}

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
	for i := 0; i < 0x8000 && i < len(rom); i++ {
		cpu.memory[i] = uint8(rom[i])
	}
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
	// fmt.Print("-----------------------------------------------------------------\n")
	cpu.stackPointer = 0xFFFE
	cpu.programCounter = 0x100
	cpu.cycle = 0
	cpu.timerAcc = 0
	// fmt.Print("  \nEnd of initialization\n")
}

// fetch8 reads the byte at PC and advances PC by 1.
func (cpu *CPU) fetch8() uint8 {
	v := cpu.memory[cpu.programCounter]
	cpu.programCounter++
	return v
}

// fetch16 reads a little-endian 16-bit value at PC and advances PC by 2.
func (cpu *CPU) fetch16() uint16 {
	lo := cpu.fetch8()
	hi := cpu.fetch8()
	return uint16(hi)<<8 | uint16(lo)
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
			cpu.serialOutput.WriteByte(cpu.memory[0xFF01])
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

// Step fetches a single opcode at PC, advances PC past it, and dispatches.
// Each handler is responsible for consuming its own immediate operands
// (via fetch8/fetch16) and advancing PC accordingly.
func (cpu *CPU) serviceInterrupt() {
	ie := cpu.memory[0xFFFF]
	iflag := cpu.memory[0xFF0F]
	pending := ie & iflag & 0x1F
	var bit uint8
	for bit = 0; bit < 5; bit++ {
		if pending&(1<<bit) != 0 {
			break
		}
	}
	cpu.ime = false
	cpu.memory[0xFF0F] = iflag &^ (1 << bit)
	cpu.tick(8)
	cpu.stackPointer--
	_ = cpu.writeMemory(cpu.stackPointer, uint8(cpu.programCounter>>8))
	cpu.tick(4)
	cpu.stackPointer--
	_ = cpu.writeMemory(cpu.stackPointer, uint8(cpu.programCounter&0xFF))
	cpu.tick(4)
	cpu.programCounter = 0x40 + uint16(bit)*8
	cpu.tick(4)
}

func (cpu *CPU) UpdateTimer(cycle_increment int) {
	cpu.cycle += cycle_increment
	tac := cpu.memory[0xFF07]
	if tac&0b100 == 0 {
		return
	}
	period := nCycles[tac&0b11]
	cpu.timerAcc += cycle_increment
	for cpu.timerAcc >= period {
		cpu.timerAcc -= period
		tima := cpu.memory[0xFF05]
		if tima == 0xFF {
			_ = cpu.writeMemory(0xFF05, cpu.memory[0xFF06])
			_ = cpu.writeMemory(0xFF0F, cpu.memory[0xFF0F]|0x04)
		} else {
			_ = cpu.writeMemory(0xFF05, tima+1)
		}
	}
}

func (cpu *CPU) Step() error {
	if int(cpu.programCounter) >= len(cpu.memory) {
		return fmt.Errorf("program counter out of bounds: 0x%04X", cpu.programCounter)
	}
	// 0x1F is for mask the bit 5 to 7 since they are not used in IE and IF.
	pending := cpu.ie() & cpu.if_() & 0x1F
	if pending != 0 {
		cpu.halt = false
		if cpu.ime {
			cpu.serviceInterrupt()
		}
	}
	if cpu.halt {
		cpu.UpdateTimer(4)
		return nil
	}
	opcode := cpu.fetch8()

	if opcode == 0xCB {
		prefixed := cpu.fetch8()
		handler := cpu.opcodeTablePrefixed[prefixed]
		if handler == nil {
			return fmt.Errorf("Unimplemented prefixed opcode: 0xCB%02X", prefixed)
		}
		handler()
		return nil
	}
	handler := cpu.opcodeTable[opcode]
	if handler == nil {
		return fmt.Errorf("Unimplemented opcode: 0x%02X", opcode)
	}

	handler()
	if cpu.pendingEnableIME {
		cpu.ime = true
		cpu.pendingEnableIME = false
	}
	return nil
}

func (cpu *CPU) Run(maxCycles int) {
	for cpu.cycle < maxCycles && !cpu.stopped {
		if err := cpu.Step(); err != nil {
			fmt.Println(err)
			return
		}

	}
}
