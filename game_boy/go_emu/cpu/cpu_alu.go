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
	cpu.opcodeTable[0x87] = func(_, _ uint8) { cpu.addA(cpu.a); cpu.cycle += 4 }
	cpu.opcodeTable[0x80] = func(_, _ uint8) { cpu.addA(cpu.b); cpu.cycle += 4 }
	cpu.opcodeTable[0x81] = func(_, _ uint8) { cpu.addA(cpu.c); cpu.cycle += 4 }
	cpu.opcodeTable[0x82] = func(_, _ uint8) { cpu.addA(cpu.d); cpu.cycle += 4 }
	cpu.opcodeTable[0x83] = func(_, _ uint8) { cpu.addA(cpu.e); cpu.cycle += 4 }
	cpu.opcodeTable[0x84] = func(_, _ uint8) { cpu.addA(cpu.h); cpu.cycle += 4 }
	cpu.opcodeTable[0x85] = func(_, _ uint8) { cpu.addA(cpu.l); cpu.cycle += 4 }
	cpu.opcodeTable[0x86] = func(_, _ uint8) { cpu.addA(cpu.memory[cpu.getHL()]); cpu.cycle += 8 }
	cpu.opcodeTable[0xC6] = func(value, _ uint8) { cpu.addA(value); cpu.cycle += 8 }

	// ADC A, r
	cpu.opcodeTable[0x8F] = func(_, _ uint8) { cpu.adcA(cpu.a); cpu.cycle += 4 }
	cpu.opcodeTable[0x88] = func(_, _ uint8) { cpu.adcA(cpu.b); cpu.cycle += 4 }
	cpu.opcodeTable[0x89] = func(_, _ uint8) { cpu.adcA(cpu.c); cpu.cycle += 4 }
	cpu.opcodeTable[0x8A] = func(_, _ uint8) { cpu.adcA(cpu.d); cpu.cycle += 4 }
	cpu.opcodeTable[0x8B] = func(_, _ uint8) { cpu.adcA(cpu.e); cpu.cycle += 4 }
	cpu.opcodeTable[0x8C] = func(_, _ uint8) { cpu.adcA(cpu.h); cpu.cycle += 4 }
	cpu.opcodeTable[0x8D] = func(_, _ uint8) { cpu.adcA(cpu.l); cpu.cycle += 4 }
	cpu.opcodeTable[0x8E] = func(_, _ uint8) { cpu.adcA(cpu.memory[cpu.getHL()]); cpu.cycle += 8 }

	// SUB A, r
	cpu.opcodeTable[0x97] = func(_, _ uint8) { cpu.subA(cpu.a); cpu.cycle += 4 }
	cpu.opcodeTable[0x90] = func(_, _ uint8) { cpu.subA(cpu.b); cpu.cycle += 4 }
	cpu.opcodeTable[0x91] = func(_, _ uint8) { cpu.subA(cpu.c); cpu.cycle += 4 }
	cpu.opcodeTable[0x92] = func(_, _ uint8) { cpu.subA(cpu.d); cpu.cycle += 4 }
	cpu.opcodeTable[0x93] = func(_, _ uint8) { cpu.subA(cpu.e); cpu.cycle += 4 }
	cpu.opcodeTable[0x94] = func(_, _ uint8) { cpu.subA(cpu.h); cpu.cycle += 4 }
	cpu.opcodeTable[0x95] = func(_, _ uint8) { cpu.subA(cpu.l); cpu.cycle += 4 }
	cpu.opcodeTable[0x96] = func(_, _ uint8) { cpu.subA(cpu.memory[cpu.getHL()]); cpu.cycle += 8 }
	cpu.opcodeTable[0xD6] = func(_, _ uint8) {
		val := cpu.memory[cpu.programCounter]
		cpu.programCounter++
		cpu.subA(val)
		cpu.cycle += 8
	}

	// SBC A, r
	cpu.opcodeTable[0x9F] = func(_, _ uint8) { cpu.sbcA(cpu.a); cpu.cycle += 4 }
	cpu.opcodeTable[0x98] = func(_, _ uint8) { cpu.sbcA(cpu.b); cpu.cycle += 4 }
	cpu.opcodeTable[0x99] = func(_, _ uint8) { cpu.sbcA(cpu.c); cpu.cycle += 4 }
	cpu.opcodeTable[0x9A] = func(_, _ uint8) { cpu.sbcA(cpu.d); cpu.cycle += 4 }
	cpu.opcodeTable[0x9B] = func(_, _ uint8) { cpu.sbcA(cpu.e); cpu.cycle += 4 }
	cpu.opcodeTable[0x9C] = func(_, _ uint8) { cpu.sbcA(cpu.h); cpu.cycle += 4 }
	cpu.opcodeTable[0x9D] = func(_, _ uint8) { cpu.sbcA(cpu.l); cpu.cycle += 4 }
	cpu.opcodeTable[0x9E] = func(_, _ uint8) { cpu.sbcA(cpu.memory[cpu.getHL()]); cpu.cycle += 8 }
	cpu.opcodeTable[0xDE] = func(_, _ uint8) {
		val := cpu.memory[cpu.programCounter]
		cpu.programCounter++
		cpu.sbcA(val)
		cpu.cycle += 8
	}

	// AND A, r
	cpu.opcodeTable[0xA7] = func(_, _ uint8) { cpu.andA(cpu.a); cpu.cycle += 4 }
	cpu.opcodeTable[0xA0] = func(_, _ uint8) { cpu.andA(cpu.b); cpu.cycle += 4 }
	cpu.opcodeTable[0xA1] = func(_, _ uint8) { cpu.andA(cpu.c); cpu.cycle += 4 }
	cpu.opcodeTable[0xA2] = func(_, _ uint8) { cpu.andA(cpu.d); cpu.cycle += 4 }
	cpu.opcodeTable[0xA3] = func(_, _ uint8) { cpu.andA(cpu.e); cpu.cycle += 4 }
	cpu.opcodeTable[0xA4] = func(_, _ uint8) { cpu.andA(cpu.h); cpu.cycle += 4 }
	cpu.opcodeTable[0xA5] = func(_, _ uint8) { cpu.andA(cpu.l); cpu.cycle += 4 }
	cpu.opcodeTable[0xA6] = func(_, _ uint8) { cpu.andA(cpu.memory[cpu.getHL()]); cpu.cycle += 8 }
	cpu.opcodeTable[0xE6] = func(_, _ uint8) {
		val := cpu.memory[cpu.programCounter]
		cpu.programCounter++
		cpu.andA(val)
		cpu.cycle += 8
	}

	// OR A, r
	cpu.opcodeTable[0xB7] = func(_, _ uint8) { cpu.orA(cpu.a); cpu.cycle += 4 }
	cpu.opcodeTable[0xB0] = func(_, _ uint8) { cpu.orA(cpu.b); cpu.cycle += 4 }
	cpu.opcodeTable[0xB1] = func(_, _ uint8) { cpu.orA(cpu.c); cpu.cycle += 4 }
	cpu.opcodeTable[0xB2] = func(_, _ uint8) { cpu.orA(cpu.d); cpu.cycle += 4 }
	cpu.opcodeTable[0xB3] = func(_, _ uint8) { cpu.orA(cpu.e); cpu.cycle += 4 }
	cpu.opcodeTable[0xB4] = func(_, _ uint8) { cpu.orA(cpu.h); cpu.cycle += 4 }
	cpu.opcodeTable[0xB5] = func(_, _ uint8) { cpu.orA(cpu.l); cpu.cycle += 4 }
	cpu.opcodeTable[0xB6] = func(_, _ uint8) { cpu.orA(cpu.memory[cpu.getHL()]); cpu.cycle += 8 }
	cpu.opcodeTable[0xF6] = func(_, _ uint8) {
		val := cpu.memory[cpu.programCounter]
		cpu.programCounter++
		cpu.orA(val)
		cpu.cycle += 8
	}

	// XOR A, r
	cpu.opcodeTable[0xAF] = func(_, _ uint8) { cpu.xorA(cpu.a); cpu.cycle += 4 }
	cpu.opcodeTable[0xA8] = func(_, _ uint8) { cpu.xorA(cpu.b); cpu.cycle += 4 }
	cpu.opcodeTable[0xA9] = func(_, _ uint8) { cpu.xorA(cpu.c); cpu.cycle += 4 }
	cpu.opcodeTable[0xAA] = func(_, _ uint8) { cpu.xorA(cpu.d); cpu.cycle += 4 }
	cpu.opcodeTable[0xAB] = func(_, _ uint8) { cpu.xorA(cpu.e); cpu.cycle += 4 }
	cpu.opcodeTable[0xAC] = func(_, _ uint8) { cpu.xorA(cpu.h); cpu.cycle += 4 }
	cpu.opcodeTable[0xAD] = func(_, _ uint8) { cpu.xorA(cpu.l); cpu.cycle += 4 }
	cpu.opcodeTable[0xAE] = func(_, _ uint8) { cpu.xorA(cpu.memory[cpu.getHL()]); cpu.cycle += 8 }
	cpu.opcodeTable[0xEE] = func(_, _ uint8) {
		val := cpu.memory[cpu.programCounter]
		cpu.programCounter++
		cpu.xorA(val)
		cpu.cycle += 8
	}

	// CP A, r
	cpu.opcodeTable[0xBF] = func(_, _ uint8) { cpu.cpA(cpu.a); cpu.cycle += 4 }
	cpu.opcodeTable[0xB8] = func(_, _ uint8) { cpu.cpA(cpu.b); cpu.cycle += 4 }
	cpu.opcodeTable[0xB9] = func(_, _ uint8) { cpu.cpA(cpu.c); cpu.cycle += 4 }
	cpu.opcodeTable[0xBA] = func(_, _ uint8) { cpu.cpA(cpu.d); cpu.cycle += 4 }
	cpu.opcodeTable[0xBB] = func(_, _ uint8) { cpu.cpA(cpu.e); cpu.cycle += 4 }
	cpu.opcodeTable[0xBC] = func(_, _ uint8) { cpu.cpA(cpu.h); cpu.cycle += 4 }
	cpu.opcodeTable[0xBD] = func(_, _ uint8) { cpu.cpA(cpu.l); cpu.cycle += 4 }
	cpu.opcodeTable[0xBE] = func(_, _ uint8) { cpu.cpA(cpu.memory[cpu.getHL()]); cpu.cycle += 8 }
	cpu.opcodeTable[0xFE] = func(_, _ uint8) {
		val := cpu.memory[cpu.programCounter]
		cpu.programCounter++
		cpu.cpA(val)
		cpu.cycle += 8
	}

	// INC r
	cpu.opcodeTable[0x3C] = func(_, _ uint8) { cpu.a = cpu.inc(cpu.a); cpu.cycle += 4 }
	cpu.opcodeTable[0x04] = func(_, _ uint8) { cpu.b = cpu.inc(cpu.b); cpu.cycle += 4 }
	cpu.opcodeTable[0x0C] = func(_, _ uint8) { cpu.c = cpu.inc(cpu.c); cpu.cycle += 4 }
	cpu.opcodeTable[0x14] = func(_, _ uint8) { cpu.d = cpu.inc(cpu.d); cpu.cycle += 4 }
	cpu.opcodeTable[0x1C] = func(_, _ uint8) { cpu.e = cpu.inc(cpu.e); cpu.cycle += 4 }
	cpu.opcodeTable[0x24] = func(_, _ uint8) { cpu.h = cpu.inc(cpu.h); cpu.cycle += 4 }
	cpu.opcodeTable[0x2C] = func(_, _ uint8) { cpu.l = cpu.inc(cpu.l); cpu.cycle += 4 }
	cpu.opcodeTable[0x34] = func(_, _ uint8) {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.inc(cpu.memory[hl]))
		cpu.cycle += 12
	}
	// DEC r
	cpu.opcodeTable[0x3D] = func(_, _ uint8) { cpu.a = cpu.dec(cpu.a); cpu.cycle += 4 }
	cpu.opcodeTable[0x05] = func(_, _ uint8) { cpu.b = cpu.dec(cpu.b); cpu.cycle += 4 }
	cpu.opcodeTable[0x0D] = func(_, _ uint8) { cpu.c = cpu.dec(cpu.c); cpu.cycle += 4 }
	cpu.opcodeTable[0x15] = func(_, _ uint8) { cpu.d = cpu.dec(cpu.d); cpu.cycle += 4 }
	cpu.opcodeTable[0x1D] = func(_, _ uint8) { cpu.e = cpu.dec(cpu.e); cpu.cycle += 4 }
	cpu.opcodeTable[0x25] = func(_, _ uint8) { cpu.h = cpu.dec(cpu.h); cpu.cycle += 4 }
	cpu.opcodeTable[0x2D] = func(_, _ uint8) { cpu.l = cpu.dec(cpu.l); cpu.cycle += 4 }
	cpu.opcodeTable[0x35] = func(_, _ uint8) {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.dec(cpu.memory[hl]))
		cpu.cycle += 12
	}
}
