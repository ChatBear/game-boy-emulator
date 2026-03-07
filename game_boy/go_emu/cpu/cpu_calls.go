package cpu

func (cpu *CPU) call(value, value2 uint8) {
	cpu.stackPointer--
	cpu.writeMemory(cpu.stackPointer, uint8(cpu.programCounter>>8))
	cpu.stackPointer--
	cpu.writeMemory(cpu.stackPointer, uint8(cpu.programCounter&0xFF))
	cpu.jump(value, value2)
}

func (cpu *CPU) initCallsOpCode() {
	cpu.opcodeTable[0xCD] = func(value, value2 uint8) {
		cpu.call(value, value2)
		cpu.cycle += 24
	}

	cpu.opcodeTable[0xC4] = func(value, value2 uint8) {
		if cpu.f&0x80 == 0 {
			cpu.call(value, value2)
		}
		cpu.cycle += 12
	}

	cpu.opcodeTable[0xCC] = func(value, value2 uint8) {
		if cpu.f&0x80 != 0 {
			cpu.call(value, value2)
		}
		cpu.cycle += 12
	}

	cpu.opcodeTable[0xD4] = func(value, value2 uint8) {
		if cpu.f&0x10 == 0 {
			cpu.call(value, value2)
		}
		cpu.cycle += 12
	}

	cpu.opcodeTable[0xDC] = func(value, value2 uint8) {
		if cpu.f&0x10 != 0 {
			cpu.call(value, value2)
		}
		cpu.cycle += 12
	}
}
