package cpu

func (cpu *CPU) ret() {
	low := cpu.memory[cpu.stackPointer]
	cpu.stackPointer++
	high := cpu.memory[cpu.stackPointer]
	cpu.stackPointer++
	cpu.jump(low, high)
}

func (cpu *CPU) initReturnsOpCode() {
	cpu.opcodeTable[0xC9] = func(_, _ uint8) {
		cpu.ret()
		cpu.cycle += 8
	}
	cpu.opcodeTable[0xC0] = func(_, _ uint8) {
		if cpu.f&0x80 == 0 {
			cpu.ret()
		}
		cpu.cycle += 8
	}

	cpu.opcodeTable[0xC8] = func(_, _ uint8) {
		if cpu.f&0x80 != 0 {
			cpu.ret()
		}
		cpu.cycle += 8
	}

	cpu.opcodeTable[0xD0] = func(_, _ uint8) {
		if cpu.f&0x10 == 0 {
			cpu.ret()
		}
		cpu.cycle += 8
	}

	cpu.opcodeTable[0xD8] = func(_, _ uint8) {
		if cpu.f&0x10 != 0 {
			cpu.ret()
		}
		cpu.cycle += 8
	}
	cpu.opcodeTable[0xD9] = func(_, _ uint8) {
		cpu.ret()
		cpu.ime = true
		cpu.cycle += 8
	}
}
