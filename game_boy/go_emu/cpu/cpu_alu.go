package cpu

func (cpu *CPU) addA(operand uint8) {
	result := uint16(cpu.a) + uint16(operand)
	cpu.f = 0
	if uint8(result) == 0 {
		cpu.f |= 0x80
	}
	if (cpu.a&0x0F)+(operand&0x0F) > 0x0F {
		cpu.f |= 0x20
	}
	if result > 0xFF {
		cpu.f |= 0x10
	}
	cpu.a = uint8(result)
}

func (cpu *CPU) adcA(operand uint8) {
	carry := uint8(0)
	if (cpu.f & 0x10) != 0 {
		carry = 1
	}
	result := uint16(cpu.a) + uint16(operand) + uint16(carry)
	cpu.f = 0
	if uint8(result) == 0 {
		cpu.f |= 0x80
	}
	if (cpu.a&0x0F)+(operand&0x0F)+carry > 0x0F {
		cpu.f |= 0x20
	}
	if result > 0xFF {
		cpu.f |= 0x10
	}
	cpu.a = uint8(result)
}

func (cpu *CPU) subA(operand uint8) {
	result := uint16(cpu.a) - uint16(operand)
	cpu.f = 0x40
	if uint8(result) == 0 {
		cpu.f |= 0x80
	}
	if (cpu.a & 0x0F) < (operand & 0x0F) {
		cpu.f |= 0x20
	}
	if result > 0xFF {
		cpu.f |= 0x10
	}
	cpu.a = uint8(result)
}

func (cpu *CPU) sbcA(operand uint8) {
	carry := uint8(0)
	if (cpu.f & 0x10) != 0 {
		carry = 1
	}
	result := uint16(cpu.a) - uint16(operand) - uint16(carry)
	cpu.f = 0x40
	if uint8(result) == 0 {
		cpu.f |= 0x80
	}
	if (cpu.a & 0x0F) < ((operand & 0x0F) + carry) {
		cpu.f |= 0x20
	}
	if result > 0xFF {
		cpu.f |= 0x10
	}
	cpu.a = uint8(result)
}

func (cpu *CPU) andA(operand uint8) {
	cpu.a = cpu.a & operand
	cpu.f = 0x20
	if cpu.a == 0 {
		cpu.f |= 0x80
	}
}

func (cpu *CPU) orA(operand uint8) {
	cpu.a |= operand
	cpu.f = 0
	if cpu.a == 0 {
		cpu.f = 0x80
	}
}

func (cpu *CPU) xorA(operand uint8) {
	cpu.a = operand ^ cpu.a
	cpu.f = 0x00
	if cpu.a == 0 {
		cpu.f = 0x80
	}
}

func (cpu *CPU) cpA(operand uint8) {
	result := uint16(cpu.a) - uint16(operand)
	cpu.f = 0x40
	if uint8(result) == 0 {
		cpu.f |= 0x80
	}
	if (cpu.a & 0x0F) < (operand & 0x0F) {
		cpu.f |= 0x20
	}
	if result > 0xFF {
		cpu.f |= 0x10
	}
}

func (cpu *CPU) inc(operand uint8) uint8 {
	result := operand + 1
	cpu.f &= 0x10
	if result == 0 {
		cpu.f |= 0x80
	}
	if (operand & 0x0F) == 0x0F {
		cpu.f |= 0x20
	}
	return result
}

func (cpu *CPU) dec(operand uint8) uint8 {
	result := operand - 1
	cpu.f &= 0x10
	cpu.f |= 0x40
	if result == 0 {
		cpu.f |= 0x80
	}
	if (operand & 0x0F) == 0x00 {
		cpu.f |= 0x20
	}
	return result
}

func (cpu *CPU) initALUOpcodes() {
	// ADD A, r
	cpu.opcodeTable[0x87] = func() { cpu.addA(cpu.a); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x80] = func() { cpu.addA(cpu.b); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x81] = func() { cpu.addA(cpu.c); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x82] = func() { cpu.addA(cpu.d); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x83] = func() { cpu.addA(cpu.e); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x84] = func() { cpu.addA(cpu.h); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x85] = func() { cpu.addA(cpu.l); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x86] = func() { cpu.addA(cpu.memory[cpu.getHL()]); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0xC6] = func() { cpu.addA(cpu.fetch8()); cpu.UpdateTimer(8) }

	// ADC A, r
	cpu.opcodeTable[0x8F] = func() { cpu.adcA(cpu.a); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x88] = func() { cpu.adcA(cpu.b); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x89] = func() { cpu.adcA(cpu.c); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x8A] = func() { cpu.adcA(cpu.d); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x8B] = func() { cpu.adcA(cpu.e); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x8C] = func() { cpu.adcA(cpu.h); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x8D] = func() { cpu.adcA(cpu.l); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x8E] = func() { cpu.adcA(cpu.memory[cpu.getHL()]); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0xCE] = func() { cpu.adcA(cpu.fetch8()); cpu.UpdateTimer(8) }

	// SUB A, r
	cpu.opcodeTable[0x97] = func() { cpu.subA(cpu.a); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x90] = func() { cpu.subA(cpu.b); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x91] = func() { cpu.subA(cpu.c); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x92] = func() { cpu.subA(cpu.d); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x93] = func() { cpu.subA(cpu.e); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x94] = func() { cpu.subA(cpu.h); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x95] = func() { cpu.subA(cpu.l); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x96] = func() { cpu.subA(cpu.memory[cpu.getHL()]); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0xD6] = func() { cpu.subA(cpu.fetch8()); cpu.UpdateTimer(8) }

	// SBC A, r
	cpu.opcodeTable[0x9F] = func() { cpu.sbcA(cpu.a); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x98] = func() { cpu.sbcA(cpu.b); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x99] = func() { cpu.sbcA(cpu.c); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x9A] = func() { cpu.sbcA(cpu.d); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x9B] = func() { cpu.sbcA(cpu.e); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x9C] = func() { cpu.sbcA(cpu.h); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x9D] = func() { cpu.sbcA(cpu.l); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x9E] = func() { cpu.sbcA(cpu.memory[cpu.getHL()]); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0xDE] = func() { cpu.sbcA(cpu.fetch8()); cpu.UpdateTimer(8) }

	// AND A, r
	cpu.opcodeTable[0xA7] = func() { cpu.andA(cpu.a); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xA0] = func() { cpu.andA(cpu.b); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xA1] = func() { cpu.andA(cpu.c); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xA2] = func() { cpu.andA(cpu.d); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xA3] = func() { cpu.andA(cpu.e); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xA4] = func() { cpu.andA(cpu.h); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xA5] = func() { cpu.andA(cpu.l); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xA6] = func() { cpu.andA(cpu.memory[cpu.getHL()]); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0xE6] = func() { cpu.andA(cpu.fetch8()); cpu.UpdateTimer(8) }

	// OR A, r
	cpu.opcodeTable[0xB7] = func() { cpu.orA(cpu.a); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xB0] = func() { cpu.orA(cpu.b); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xB1] = func() { cpu.orA(cpu.c); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xB2] = func() { cpu.orA(cpu.d); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xB3] = func() { cpu.orA(cpu.e); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xB4] = func() { cpu.orA(cpu.h); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xB5] = func() { cpu.orA(cpu.l); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xB6] = func() { cpu.orA(cpu.memory[cpu.getHL()]); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0xF6] = func() { cpu.orA(cpu.fetch8()); cpu.UpdateTimer(8) }

	// XOR A, r
	cpu.opcodeTable[0xAF] = func() { cpu.xorA(cpu.a); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xA8] = func() { cpu.xorA(cpu.b); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xA9] = func() { cpu.xorA(cpu.c); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xAA] = func() { cpu.xorA(cpu.d); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xAB] = func() { cpu.xorA(cpu.e); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xAC] = func() { cpu.xorA(cpu.h); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xAD] = func() { cpu.xorA(cpu.l); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xAE] = func() { cpu.xorA(cpu.memory[cpu.getHL()]); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0xEE] = func() { cpu.xorA(cpu.fetch8()); cpu.UpdateTimer(8) }

	// CP A, r
	cpu.opcodeTable[0xBF] = func() { cpu.cpA(cpu.a); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xB8] = func() { cpu.cpA(cpu.b); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xB9] = func() { cpu.cpA(cpu.c); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xBA] = func() { cpu.cpA(cpu.d); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xBB] = func() { cpu.cpA(cpu.e); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xBC] = func() { cpu.cpA(cpu.h); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xBD] = func() { cpu.cpA(cpu.l); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0xBE] = func() { cpu.cpA(cpu.memory[cpu.getHL()]); cpu.UpdateTimer(8) }
	cpu.opcodeTable[0xFE] = func() { cpu.cpA(cpu.fetch8()); cpu.UpdateTimer(8) }

	// INC r
	cpu.opcodeTable[0x3C] = func() { cpu.a = cpu.inc(cpu.a); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x04] = func() { cpu.b = cpu.inc(cpu.b); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x0C] = func() { cpu.c = cpu.inc(cpu.c); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x14] = func() { cpu.d = cpu.inc(cpu.d); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x1C] = func() { cpu.e = cpu.inc(cpu.e); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x24] = func() { cpu.h = cpu.inc(cpu.h); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x2C] = func() { cpu.l = cpu.inc(cpu.l); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x34] = func() {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.inc(cpu.memory[hl]))
		cpu.UpdateTimer(12)
	}
	// DEC r
	cpu.opcodeTable[0x3D] = func() { cpu.a = cpu.dec(cpu.a); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x05] = func() { cpu.b = cpu.dec(cpu.b); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x0D] = func() { cpu.c = cpu.dec(cpu.c); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x15] = func() { cpu.d = cpu.dec(cpu.d); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x1D] = func() { cpu.e = cpu.dec(cpu.e); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x25] = func() { cpu.h = cpu.dec(cpu.h); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x2D] = func() { cpu.l = cpu.dec(cpu.l); cpu.UpdateTimer(4) }
	cpu.opcodeTable[0x35] = func() {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.dec(cpu.memory[hl]))
		cpu.UpdateTimer(12)
	}
}
