package cpu

func (cpu *CPU) initStackOpcodes() {
	// PUSH
	cpu.opcodeTable[0xF5] = func(_, _ uint8) {
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.a)
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.f)
		cpu.cycle += 16
	}
	cpu.opcodeTable[0xC5] = func(_, _ uint8) {
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.b)
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.c)
		cpu.cycle += 16
	}
	cpu.opcodeTable[0xD5] = func(_, _ uint8) {
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.d)
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.e)
		cpu.cycle += 16
	}
	cpu.opcodeTable[0xE5] = func(_, _ uint8) {
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.h)
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.l)
		cpu.cycle += 16
	}

	// POP
	cpu.opcodeTable[0xF1] = func(_, _ uint8) {
		cpu.f = cpu.memory[cpu.stackPointer] & 0xF0
		cpu.stackPointer++
		cpu.a = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.cycle += 12
	}
	cpu.opcodeTable[0xC1] = func(_, _ uint8) {
		cpu.c = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.b = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.cycle += 12
	}
	cpu.opcodeTable[0xD1] = func(_, _ uint8) {
		cpu.e = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.d = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.cycle += 12
	}
	cpu.opcodeTable[0xE1] = func(_, _ uint8) {
		cpu.l = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.h = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.cycle += 12
	}
}
