package cpu

func (cpu *CPU) initLoadOpcodes() {
	// LD r, n — load immediate into register
	cpu.opcodeTable[0x06] = func() { cpu.b = cpu.fetch8(); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x0E] = func() { cpu.c = cpu.fetch8(); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x16] = func() { cpu.d = cpu.fetch8(); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x1E] = func() { cpu.e = cpu.fetch8(); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x26] = func() { cpu.h = cpu.fetch8(); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x2E] = func() { cpu.l = cpu.fetch8(); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x3E] = func() { cpu.a = cpu.fetch8(); cpu.UpdateTimer(8) }

	// LD A, r
	cpu.opcodeTable[0x78] = func() { cpu.a = cpu.b; cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x79] = func() { cpu.a = cpu.c; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x7A] = func() { cpu.a = cpu.d; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x7B] = func() { cpu.a = cpu.e; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x7C] = func() { cpu.a = cpu.h; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x7D] = func() { cpu.a = cpu.l; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x7E] = func() { cpu.a = cpu.memory[cpu.getHL()]; cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x7F] = func() { cpu.UpdateTimer(4) }

	// LD B, r
	cpu.opcodeTable[0x40] = func() { cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x41] = func() { cpu.b = cpu.c; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x42] = func() { cpu.b = cpu.d; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x43] = func() { cpu.b = cpu.e; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x44] = func() { cpu.b = cpu.h; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x45] = func() { cpu.b = cpu.l; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x46] = func() { cpu.b = cpu.memory[cpu.getHL()]; cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x47] = func() { cpu.b = cpu.a; cpu.UpdateTimer(4) }

	// LD C, r
	cpu.opcodeTable[0x48] = func() { cpu.c = cpu.b; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x49] = func() { cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x4A] = func() { cpu.c = cpu.d; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x4B] = func() { cpu.c = cpu.e; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x4C] = func() { cpu.c = cpu.h; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x4D] = func() { cpu.c = cpu.l; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x4E] = func() { cpu.c = cpu.memory[cpu.getHL()]; cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x4F] = func() { cpu.c = cpu.a; cpu.UpdateTimer(4) }

	// LD D, r
	cpu.opcodeTable[0x50] = func() { cpu.d = cpu.b; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x51] = func() { cpu.d = cpu.c; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x52] = func() { cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x53] = func() { cpu.d = cpu.e; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x54] = func() { cpu.d = cpu.h; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x55] = func() { cpu.d = cpu.l; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x56] = func() { cpu.d = cpu.memory[cpu.getHL()]; cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x57] = func() { cpu.d = cpu.a; cpu.UpdateTimer(4) }

	// LD E, r
	cpu.opcodeTable[0x58] = func() { cpu.e = cpu.b; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x59] = func() { cpu.e = cpu.c; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x5A] = func() { cpu.e = cpu.d; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x5B] = func() { cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x5C] = func() { cpu.e = cpu.h; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x5D] = func() { cpu.e = cpu.l; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x5E] = func() { cpu.e = cpu.memory[cpu.getHL()]; cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x5F] = func() { cpu.e = cpu.a; cpu.UpdateTimer(4) }

	// LD H, r
	cpu.opcodeTable[0x60] = func() { cpu.h = cpu.b; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x61] = func() { cpu.h = cpu.c; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x62] = func() { cpu.h = cpu.d; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x63] = func() { cpu.h = cpu.e; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x64] = func() { cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x65] = func() { cpu.h = cpu.l; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x66] = func() { cpu.h = cpu.memory[cpu.getHL()]; cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x67] = func() { cpu.h = cpu.a; cpu.UpdateTimer(4) }

	// LD L, r
	cpu.opcodeTable[0x68] = func() { cpu.l = cpu.b; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x69] = func() { cpu.l = cpu.c; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x6A] = func() { cpu.l = cpu.d; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x6B] = func() { cpu.l = cpu.e; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x6C] = func() { cpu.l = cpu.h; cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x6D] = func() { cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x6E] = func() { cpu.l = cpu.memory[cpu.getHL()]; cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x6F] = func() { cpu.l = cpu.a; cpu.UpdateTimer(4) }

	// LD (HL), r
	cpu.opcodeTable[0x70] = func() { cpu.writeMemory(cpu.getHL(), cpu.b); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x71] = func() { cpu.writeMemory(cpu.getHL(), cpu.c); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x72] = func() { cpu.writeMemory(cpu.getHL(), cpu.d); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x73] = func() { cpu.writeMemory(cpu.getHL(), cpu.e); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x74] = func() { cpu.writeMemory(cpu.getHL(), cpu.h); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x75] = func() { cpu.writeMemory(cpu.getHL(), cpu.l); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x77] = func() { cpu.writeMemory(cpu.getHL(), cpu.a); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x36] = func() { cpu.writeMemory(cpu.getHL(), cpu.fetch8()); cpu.UpdateTimer(12) }

	// LD (nn), A
	cpu.opcodeTable[0xEA] = func() {
		cpu.writeMemory(cpu.fetch16(), cpu.a)
		cpu.UpdateTimer(16)
	}

	// LD A, (rr)
	cpu.opcodeTable[0x0A] = func() { cpu.a = cpu.memory[cpu.getBC()]; cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x1A] = func() { cpu.a = cpu.memory[cpu.getDE()]; cpu.UpdateTimer(8) }
	cpu.opcodeTable[0xFA] = func() {
		cpu.a = cpu.memory[cpu.fetch16()]
		cpu.UpdateTimer(16)
	}

	// LD (rr), A
	cpu.opcodeTable[0x02] = func() { cpu.writeMemory(cpu.getBC(), cpu.a); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0x12] = func() { cpu.writeMemory(cpu.getDE(), cpu.a); cpu.UpdateTimer(8) }

	// LDH — high memory
	cpu.opcodeTable[0xF2] = func() { cpu.a = cpu.memory[0xFF00+uint16(cpu.c)]; cpu.UpdateTimer(8) }
	cpu.opcodeTable[0xE2] = func() { cpu.writeMemory(0xFF00+uint16(cpu.c), cpu.a); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0xF0] = func() { cpu.a = cpu.memory[0xFF00+uint16(cpu.fetch8())]; cpu.UpdateTimer(12) }
	cpu.opcodeTable[0xE0] = func() { cpu.writeMemory(0xFF00+uint16(cpu.fetch8()), cpu.a); cpu.UpdateTimer(12) }

	// LDD / LDI — load with decrement/increment
	cpu.opcodeTable[0x3A] = func() {
		hl := cpu.getHL()
		cpu.a = cpu.memory[hl]
		cpu.setHL(hl - 1)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTable[0x32] = func() {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.a)
		cpu.setHL(hl - 1)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTable[0x2A] = func() {
		hl := cpu.getHL()
		cpu.a = cpu.memory[hl]
		cpu.setHL(hl + 1)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTable[0x22] = func() {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.a)
		cpu.setHL(hl + 1)
		cpu.UpdateTimer(8)
	}

	// LD rr, nn — 16-bit immediate loads
	cpu.opcodeTable[0x01] = func() { cpu.setBC(cpu.fetch16()); cpu.UpdateTimer(12) }
	cpu.opcodeTable[0x11] = func() { cpu.setDE(cpu.fetch16()); cpu.UpdateTimer(12) }
	cpu.opcodeTable[0x21] = func() { cpu.setHL(cpu.fetch16()); cpu.UpdateTimer(12) }
	cpu.opcodeTable[0x31] = func() { cpu.stackPointer = cpu.fetch16(); cpu.UpdateTimer(12) }

	// LD SP, HL
	cpu.opcodeTable[0xF9] = func() {
		cpu.stackPointer = cpu.getHL()
		cpu.UpdateTimer(8)
	}

	// LDHL SP,n — put SP + signed n into HL
	cpu.opcodeTable[0xF8] = func() {
		value := cpu.fetch8()
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
		cpu.UpdateTimer(12)
	}

	// LD (nn), SP
	cpu.opcodeTable[0x08] = func() {
		address := cpu.fetch16()
		cpu.writeMemory(address, uint8(cpu.stackPointer&0xFF))
		cpu.writeMemory(address+1, uint8(cpu.stackPointer>>8))
		cpu.UpdateTimer(20)
	}
}
