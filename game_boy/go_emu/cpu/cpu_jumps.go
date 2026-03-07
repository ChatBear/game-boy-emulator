package cpu

func (cpu *CPU) jump(value, value2 uint8) {
	cpu.programCounter = uint16(value2)<<8 | uint16(value)
}

func (cpu *CPU) jumpRelative(value uint8) {
	cpu.programCounter = uint16(int32(cpu.programCounter) + int32(int8(value)))
}

func (cpu *CPU) initJumpsOpCodes() {
	cpu.opcodeTable[0xC3] = func(value, value2 uint8) {
		cpu.jump(value, value2)
		cpu.cycle += 12
	}

	cpu.opcodeTable[0xC2] = func(value, value2 uint8) {
		if cpu.f&0x80 == 0 {
			cpu.jump(value, value2)
		}
		cpu.cycle += 12
	}

	cpu.opcodeTable[0xCA] = func(value, value2 uint8) {
		if cpu.f&0x80 != 0 {
			cpu.jump(value, value2)
		}
		cpu.cycle += 12
	}

	cpu.opcodeTable[0xD2] = func(value, value2 uint8) {
		if cpu.f&0x10 == 0 {
			cpu.jump(value, value2)
		}
		cpu.cycle += 12
	}

	cpu.opcodeTable[0xDA] = func(value, value2 uint8) {
		if cpu.f&0x10 != 0 {
			cpu.jump(value, value2)
		}
		cpu.cycle += 12
	}

	cpu.opcodeTable[0xE9] = func(_, _ uint8) {
		cpu.programCounter = cpu.getHL()
		cpu.cycle += 4
	}

	cpu.opcodeTable[0x18] = func(value uint8, _ uint8) {
		cpu.jumpRelative(value)
		cpu.cycle += 8
	}

	cpu.opcodeTable[0x20] = func(value uint8, _ uint8) {
		if cpu.f&0x80 == 0 {
			cpu.jumpRelative(value)
		}
		cpu.cycle += 8
	}

	cpu.opcodeTable[0x28] = func(value uint8, _ uint8) {
		if cpu.f&0x80 != 0 {
			cpu.jumpRelative(value)
		}
		cpu.cycle += 8
	}

	cpu.opcodeTable[0x30] = func(value uint8, _ uint8) {
		if cpu.f&0x10 == 0 {
			cpu.jumpRelative(value)
		}
		cpu.cycle += 8
	}

	cpu.opcodeTable[0x38] = func(value uint8, _ uint8) {
		if cpu.f&0x10 != 0 {
			cpu.jumpRelative(value)
		}
		cpu.cycle += 8
	}

}
