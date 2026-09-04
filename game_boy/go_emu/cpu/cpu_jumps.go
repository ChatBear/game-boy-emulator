package cpu

func (cpu *CPU) jump(addr uint16) {
	cpu.programCounter = addr
}

func (cpu *CPU) jumpRelative(offset uint8) {
	cpu.programCounter = uint16(int32(cpu.programCounter) + int32(int8(offset)))
}

func (cpu *CPU) initJumpsOpCodes() {
	cpu.opcodeTable[0xC3] = func() {
		addr := cpu.fetch16()
		cpu.jump(addr)
		cpu.UpdateTimer(12)
	}

	cpu.opcodeTable[0xC2] = func() {
		addr := cpu.fetch16()
		if cpu.f&0x80 == 0 {
			cpu.jump(addr)
		}
		cpu.UpdateTimer(12)
	}

	cpu.opcodeTable[0xCA] = func() {
		addr := cpu.fetch16()
		if cpu.f&0x80 != 0 {
			cpu.jump(addr)
		}
		cpu.UpdateTimer(12)
	}

	cpu.opcodeTable[0xD2] = func() {
		addr := cpu.fetch16()
		if cpu.f&0x10 == 0 {
			cpu.jump(addr)
		}
		cpu.UpdateTimer(12)
	}

	cpu.opcodeTable[0xDA] = func() {
		addr := cpu.fetch16()
		if cpu.f&0x10 != 0 {
			cpu.jump(addr)
		}
		cpu.UpdateTimer(12)
	}

	cpu.opcodeTable[0xE9] = func() {
		cpu.programCounter = cpu.getHL()
		cpu.UpdateTimer(4)
	}

	cpu.opcodeTable[0x18] = func() {
		offset := cpu.fetch8()
		cpu.jumpRelative(offset)
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTable[0x20] = func() {
		offset := cpu.fetch8()
		if cpu.f&0x80 == 0 {
			cpu.jumpRelative(offset)
		}
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTable[0x28] = func() {
		offset := cpu.fetch8()
		if cpu.f&0x80 != 0 {
			cpu.jumpRelative(offset)
		}
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTable[0x30] = func() {
		offset := cpu.fetch8()
		if cpu.f&0x10 == 0 {
			cpu.jumpRelative(offset)
		}
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTable[0x38] = func() {
		offset := cpu.fetch8()
		if cpu.f&0x10 != 0 {
			cpu.jumpRelative(offset)
		}
		cpu.UpdateTimer(8)
	}
}
