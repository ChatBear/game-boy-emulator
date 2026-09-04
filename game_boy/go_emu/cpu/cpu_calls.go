package cpu

func (cpu *CPU) call(addr uint16) {
	cpu.stackPointer--
	cpu.writeMemory(cpu.stackPointer, uint8(cpu.programCounter>>8))
	cpu.stackPointer--
	cpu.writeMemory(cpu.stackPointer, uint8(cpu.programCounter&0xFF))
	cpu.jump(addr)
}

func (cpu *CPU) initCallsOpCode() {
	cpu.opcodeTable[0xCD] = func() {
		addr := cpu.fetch16()
		cpu.call(addr)
		cpu.UpdateTimer(24)
	}

	cpu.opcodeTable[0xC4] = func() {
		addr := cpu.fetch16()
		if cpu.f&0x80 == 0 {
			cpu.call(addr)
		}
		cpu.UpdateTimer(12)
	}

	cpu.opcodeTable[0xCC] = func() {
		addr := cpu.fetch16()
		if cpu.f&0x80 != 0 {
			cpu.call(addr)
		}
		cpu.UpdateTimer(12)
	}

	cpu.opcodeTable[0xD4] = func() {
		addr := cpu.fetch16()
		if cpu.f&0x10 == 0 {
			cpu.call(addr)
		}
		cpu.UpdateTimer(12)
	}

	cpu.opcodeTable[0xDC] = func() {
		addr := cpu.fetch16()
		if cpu.f&0x10 != 0 {
			cpu.call(addr)
		}
		cpu.UpdateTimer(12)
	}
}
