package cpu

func (cpu *CPU) initStackOpcodes() {
	// PUSH
	cpu.opcodeTable[0xF5] = func() {
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.a)
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.f)
		cpu.UpdateTimer(16)
	}
	cpu.opcodeTable[0xC5] = func() {
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.b)
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.c)
		cpu.UpdateTimer(16)
	}
	cpu.opcodeTable[0xD5] = func() {
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.d)
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.e)
		cpu.UpdateTimer(16)
	}
	cpu.opcodeTable[0xE5] = func() {
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.h)
		cpu.stackPointer--
		cpu.writeMemory(uint16(cpu.stackPointer), cpu.l)
		cpu.UpdateTimer(16)
	}

	// POP
	cpu.opcodeTable[0xF1] = func() {
		cpu.f = cpu.memory[cpu.stackPointer] & 0xF0
		cpu.stackPointer++
		cpu.a = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.UpdateTimer(12)
	}
	cpu.opcodeTable[0xC1] = func() {
		cpu.c = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.b = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.UpdateTimer(12)
	}
	cpu.opcodeTable[0xD1] = func() {
		cpu.e = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.d = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.UpdateTimer(12)
	}
	cpu.opcodeTable[0xE1] = func() {
		cpu.l = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.h = cpu.memory[cpu.stackPointer]
		cpu.stackPointer++
		cpu.UpdateTimer(12)
	}
}
