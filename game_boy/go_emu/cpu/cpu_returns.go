package cpu

func (cpu *CPU) ret() {
	low := cpu.memory[cpu.stackPointer]
	cpu.stackPointer++
	high := cpu.memory[cpu.stackPointer]
	cpu.stackPointer++
	cpu.jump(uint16(high)<<8 | uint16(low))
}

func (cpu *CPU) initReturnsOpCode() {
	cpu.opcodeTable[0xC9] = func() {
		cpu.ret()
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTable[0xC0] = func() {
		if cpu.f&0x80 == 0 {
			cpu.ret()
		}
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTable[0xC8] = func() {
		if cpu.f&0x80 != 0 {
			cpu.ret()
		}
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTable[0xD0] = func() {
		if cpu.f&0x10 == 0 {
			cpu.ret()
		}
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTable[0xD8] = func() {
		if cpu.f&0x10 != 0 {
			cpu.ret()
		}
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTable[0xD9] = func() {
		cpu.ret()
		cpu.ime = true
		cpu.UpdateTimer(8)
	}
}
