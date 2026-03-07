package cpu

func (cpu *CPU) initLoadOpcodes() {
	// LD r, n — load immediate into register
	cpu.opcodeTable[0x06] = func(value, value2 uint8) { cpu.b = value; cpu.cycle += 8 }
	cpu.opcodeTable[0x0E] = func(value, value2 uint8) { cpu.c = value; cpu.cycle += 8 }
	cpu.opcodeTable[0x16] = func(value, value2 uint8) { cpu.d = value; cpu.cycle += 8 }
	cpu.opcodeTable[0x1E] = func(value, value2 uint8) { cpu.e = value; cpu.cycle += 8 }
	cpu.opcodeTable[0x26] = func(value, value2 uint8) { cpu.h = value; cpu.cycle += 8 }
	cpu.opcodeTable[0x2E] = func(value, value2 uint8) { cpu.l = value; cpu.cycle += 8 }
	cpu.opcodeTable[0x3E] = func(value, value2 uint8) { cpu.a = value; cpu.cycle += 8 }

	// LD A, r
	cpu.opcodeTable[0x78] = func(_, _ uint8) { cpu.a = cpu.b; cpu.cycle += 8 }
	cpu.opcodeTable[0x79] = func(_, _ uint8) { cpu.a = cpu.c; cpu.cycle += 4 }
	cpu.opcodeTable[0x7A] = func(_, _ uint8) { cpu.a = cpu.d; cpu.cycle += 4 }
	cpu.opcodeTable[0x7B] = func(_, _ uint8) { cpu.a = cpu.e; cpu.cycle += 4 }
	cpu.opcodeTable[0x7C] = func(_, _ uint8) { cpu.a = cpu.h; cpu.cycle += 4 }
	cpu.opcodeTable[0x7D] = func(_, _ uint8) { cpu.a = cpu.l; cpu.cycle += 4 }
	cpu.opcodeTable[0x7E] = func(_, _ uint8) { cpu.a = cpu.memory[cpu.getHL()]; cpu.cycle += 8 }
	cpu.opcodeTable[0x7F] = func(_, _ uint8) { cpu.cycle += 4 }

	// LD B, r
	cpu.opcodeTable[0x40] = func(_, _ uint8) { cpu.cycle += 4 }
	cpu.opcodeTable[0x41] = func(_, _ uint8) { cpu.b = cpu.c; cpu.cycle += 4 }
	cpu.opcodeTable[0x42] = func(_, _ uint8) { cpu.b = cpu.d; cpu.cycle += 4 }
	cpu.opcodeTable[0x43] = func(_, _ uint8) { cpu.b = cpu.e; cpu.cycle += 4 }
	cpu.opcodeTable[0x44] = func(_, _ uint8) { cpu.b = cpu.h; cpu.cycle += 4 }
	cpu.opcodeTable[0x45] = func(_, _ uint8) { cpu.b = cpu.l; cpu.cycle += 4 }
	cpu.opcodeTable[0x46] = func(_, _ uint8) { cpu.b = cpu.memory[cpu.getHL()]; cpu.cycle += 8 }
	cpu.opcodeTable[0x47] = func(_, _ uint8) { cpu.b = cpu.a; cpu.cycle += 4 }

	// LD C, r
	cpu.opcodeTable[0x48] = func(_, _ uint8) { cpu.c = cpu.b; cpu.cycle += 4 }
	cpu.opcodeTable[0x49] = func(_, _ uint8) { cpu.cycle += 4 }
	cpu.opcodeTable[0x4A] = func(_, _ uint8) { cpu.c = cpu.d; cpu.cycle += 4 }
	cpu.opcodeTable[0x4B] = func(_, _ uint8) { cpu.c = cpu.e; cpu.cycle += 4 }
	cpu.opcodeTable[0x4C] = func(_, _ uint8) { cpu.c = cpu.h; cpu.cycle += 4 }
	cpu.opcodeTable[0x4D] = func(_, _ uint8) { cpu.c = cpu.l; cpu.cycle += 4 }
	cpu.opcodeTable[0x4E] = func(_, _ uint8) { cpu.c = cpu.memory[cpu.getHL()]; cpu.cycle += 8 }
	cpu.opcodeTable[0x4F] = func(_, _ uint8) { cpu.c = cpu.a; cpu.cycle += 4 }

	// LD D, r
	cpu.opcodeTable[0x50] = func(_, _ uint8) { cpu.d = cpu.b; cpu.cycle += 4 }
	cpu.opcodeTable[0x51] = func(_, _ uint8) { cpu.d = cpu.c; cpu.cycle += 4 }
	cpu.opcodeTable[0x52] = func(_, _ uint8) { cpu.cycle += 4 }
	cpu.opcodeTable[0x53] = func(_, _ uint8) { cpu.d = cpu.e; cpu.cycle += 4 }
	cpu.opcodeTable[0x54] = func(_, _ uint8) { cpu.d = cpu.h; cpu.cycle += 4 }
	cpu.opcodeTable[0x55] = func(_, _ uint8) { cpu.d = cpu.l; cpu.cycle += 4 }
	cpu.opcodeTable[0x56] = func(_, _ uint8) { cpu.d = cpu.memory[cpu.getHL()]; cpu.cycle += 8 }
	cpu.opcodeTable[0x57] = func(_, _ uint8) { cpu.d = cpu.a; cpu.cycle += 4 }

	// LD E, r
	cpu.opcodeTable[0x58] = func(_, _ uint8) { cpu.e = cpu.b; cpu.cycle += 4 }
	cpu.opcodeTable[0x59] = func(_, _ uint8) { cpu.e = cpu.c; cpu.cycle += 4 }
	cpu.opcodeTable[0x5A] = func(_, _ uint8) { cpu.e = cpu.d; cpu.cycle += 4 }
	cpu.opcodeTable[0x5B] = func(_, _ uint8) { cpu.cycle += 4 }
	cpu.opcodeTable[0x5C] = func(_, _ uint8) { cpu.e = cpu.h; cpu.cycle += 4 }
	cpu.opcodeTable[0x5D] = func(_, _ uint8) { cpu.e = cpu.l; cpu.cycle += 4 }
	cpu.opcodeTable[0x5E] = func(_, _ uint8) { cpu.e = cpu.memory[cpu.getHL()]; cpu.cycle += 8 }
	cpu.opcodeTable[0x5F] = func(_, _ uint8) { cpu.e = cpu.a; cpu.cycle += 4 }

	// LD H, r
	cpu.opcodeTable[0x60] = func(_, _ uint8) { cpu.h = cpu.b; cpu.cycle += 4 }
	cpu.opcodeTable[0x61] = func(_, _ uint8) { cpu.h = cpu.c; cpu.cycle += 4 }
	cpu.opcodeTable[0x62] = func(_, _ uint8) { cpu.h = cpu.d; cpu.cycle += 4 }
	cpu.opcodeTable[0x63] = func(_, _ uint8) { cpu.h = cpu.e; cpu.cycle += 4 }
	cpu.opcodeTable[0x64] = func(_, _ uint8) { cpu.cycle += 4 }
	cpu.opcodeTable[0x65] = func(_, _ uint8) { cpu.h = cpu.l; cpu.cycle += 4 }
	cpu.opcodeTable[0x66] = func(_, _ uint8) { cpu.h = cpu.memory[cpu.getHL()]; cpu.cycle += 8 }
	cpu.opcodeTable[0x67] = func(_, _ uint8) { cpu.h = cpu.a; cpu.cycle += 4 }

	// LD L, r
	cpu.opcodeTable[0x68] = func(_, _ uint8) { cpu.l = cpu.b; cpu.cycle += 4 }
	cpu.opcodeTable[0x69] = func(_, _ uint8) { cpu.l = cpu.c; cpu.cycle += 4 }
	cpu.opcodeTable[0x6A] = func(_, _ uint8) { cpu.l = cpu.d; cpu.cycle += 4 }
	cpu.opcodeTable[0x6B] = func(_, _ uint8) { cpu.l = cpu.e; cpu.cycle += 4 }
	cpu.opcodeTable[0x6C] = func(_, _ uint8) { cpu.l = cpu.h; cpu.cycle += 4 }
	cpu.opcodeTable[0x6D] = func(_, _ uint8) { cpu.cycle += 4 }
	cpu.opcodeTable[0x6E] = func(_, _ uint8) { cpu.l = cpu.memory[cpu.getHL()]; cpu.cycle += 8 }
	cpu.opcodeTable[0x6F] = func(_, _ uint8) { cpu.l = cpu.a; cpu.cycle += 4 }

	// LD (HL), r
	cpu.opcodeTable[0x70] = func(_, _ uint8) { cpu.writeMemory(cpu.getHL(), cpu.b); cpu.cycle += 8 }
	cpu.opcodeTable[0x71] = func(_, _ uint8) { cpu.writeMemory(cpu.getHL(), cpu.c); cpu.cycle += 8 }
	cpu.opcodeTable[0x72] = func(_, _ uint8) { cpu.writeMemory(cpu.getHL(), cpu.d); cpu.cycle += 8 }
	cpu.opcodeTable[0x73] = func(_, _ uint8) { cpu.writeMemory(cpu.getHL(), cpu.e); cpu.cycle += 8 }
	cpu.opcodeTable[0x74] = func(_, _ uint8) { cpu.writeMemory(cpu.getHL(), cpu.h); cpu.cycle += 8 }
	cpu.opcodeTable[0x75] = func(_, _ uint8) { cpu.writeMemory(cpu.getHL(), cpu.l); cpu.cycle += 8 }
	cpu.opcodeTable[0x77] = func(_, _ uint8) { cpu.writeMemory(cpu.getHL(), cpu.a); cpu.cycle += 8 }
	cpu.opcodeTable[0x36] = func(value, _ uint8) { cpu.writeMemory(cpu.getHL(), value); cpu.cycle += 12 }

	// LD (nn), A
	cpu.opcodeTable[0xEA] = func(value, value2 uint8) {
		addr := uint16(value2)<<8 | uint16(value)
		cpu.writeMemory(addr, cpu.a)
	}

	// LD A, (rr)
	cpu.opcodeTable[0x0A] = func(_, _ uint8) { cpu.a = cpu.memory[cpu.getBC()]; cpu.cycle += 8 }
	cpu.opcodeTable[0x1A] = func(_, _ uint8) { cpu.a = cpu.memory[cpu.getDE()]; cpu.cycle += 8 }
	cpu.opcodeTable[0xFA] = func(value, value2 uint8) {
		addr := uint16(value2)<<8 | uint16(value)
		cpu.a = cpu.memory[addr]
		cpu.cycle += 16
	}

	// LD (rr), A
	cpu.opcodeTable[0x02] = func(_, _ uint8) { cpu.writeMemory(cpu.getBC(), cpu.a); cpu.cycle += 8 }
	cpu.opcodeTable[0x12] = func(_, _ uint8) { cpu.writeMemory(cpu.getDE(), cpu.a); cpu.cycle += 8 }

	// LDH — high memory
	cpu.opcodeTable[0xF2] = func(_, _ uint8) { cpu.a = cpu.memory[0xFF00+uint16(cpu.c)]; cpu.cycle += 8 }
	cpu.opcodeTable[0xE2] = func(_, _ uint8) { cpu.writeMemory(0xFF00+uint16(cpu.c), cpu.a); cpu.cycle += 8 }
	cpu.opcodeTable[0xE0] = func(value, _ uint8) { cpu.a = cpu.memory[0xFF00+uint16(value)]; cpu.cycle += 16 }
	cpu.opcodeTable[0xF0] = func(value, _ uint8) { cpu.writeMemory(0xFF00+uint16(value), cpu.a); cpu.cycle += 16 }

	// LDD / LDI — load with decrement/increment
	cpu.opcodeTable[0x3A] = func(_, _ uint8) {
		hl := cpu.getHL()
		cpu.a = cpu.memory[hl]
		cpu.setHL(hl - 1)
		cpu.cycle += 8
	}
	cpu.opcodeTable[0x32] = func(_, _ uint8) {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.a)
		cpu.setHL(hl - 1)
		cpu.cycle += 8
	}
	cpu.opcodeTable[0x2A] = func(_, _ uint8) {
		hl := cpu.getHL()
		cpu.a = cpu.memory[hl]
		cpu.setHL(hl + 1)
		cpu.cycle += 8
	}
	cpu.opcodeTable[0x22] = func(_, _ uint8) {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.a)
		cpu.setHL(hl + 1)
		cpu.cycle += 8
	}

	// LD rr, nn — 16-bit immediate loads
	cpu.opcodeTable[0x01] = func(value, value2 uint8) {
		cpu.b = value2
		cpu.c = value
		cpu.cycle += 12
	}
	cpu.opcodeTable[0x11] = func(value, value2 uint8) {
		cpu.d = value2
		cpu.e = value
		cpu.cycle += 12
	}
	cpu.opcodeTable[0x21] = func(value, value2 uint8) {
		cpu.h = value2
		cpu.l = value
		cpu.cycle += 12
	}
	cpu.opcodeTable[0x31] = func(value, value2 uint8) {
		cpu.stackPointer = uint16(value2)<<8 | uint16(value)
		cpu.cycle += 12
	}

	// LD SP, HL
	cpu.opcodeTable[0xF9] = func(_, _ uint8) {
		cpu.stackPointer = cpu.getHL()
		cpu.cycle += 8
	}

	// LDHL SP,n — put SP + signed n into HL
	cpu.opcodeTable[0xF8] = func(value, _ uint8) {
		sp := uint16(cpu.stackPointer)
		n := uint16(value)

		cpu.f = 0
		if (sp&0x0F)+(n&0x0F) > 0x0F {
			cpu.f |= 0x20 // H flag
		}
		if (sp&0xFF)+(n&0xFF) > 0xFF {
			cpu.f |= 0x10 // C flag
		}

		cpu.setHL(uint16(int32(int16(sp)) + int32(int8(value))))
		cpu.cycle += 12
	}

	// LD (nn), SP
	cpu.opcodeTable[0x08] = func(value, value2 uint8) {
		address := uint16(value2)<<8 | uint16(value)
		cpu.writeMemory(address, uint8(cpu.stackPointer&0xFF))
		cpu.writeMemory(address+1, uint8(cpu.stackPointer>>8))
		cpu.cycle += 20
	}
}
