package cpu

func (cpu *CPU) add16HL(operand uint16) {
	hl := cpu.getHL()
	cpu.f &^= 0x70

	if (hl&0x0FFF)+(operand&0x0FFF) > 0x0FFF {
		cpu.f |= 0x20
	}
	if uint32(hl)+uint32(operand) > 0xFFFF {
		cpu.f |= 0x10
	}

	cpu.setHL(hl + operand)
}

func (cpu *CPU) add16SP(operand uint8) {
	n := uint16(operand)

	cpu.f = 0
	if (cpu.stackPointer&0x0F)+(n&0x0F) > 0x0F {
		cpu.f |= 0x20 // H flag
	}
	if (cpu.stackPointer&0xFF)+(n&0xFF) > 0xFF {
		cpu.f |= 0x10 // C flag
	}

	cpu.stackPointer = uint16(int32(int16(cpu.stackPointer)) + int32(int8(operand)))
}

func (cpu *CPU) init16BitArthmeticOpCode() {
	cpu.opcodeTable[0x09] = func(_, _ uint8) { cpu.add16HL(cpu.getBC()); cpu.cycle += 8 }
	cpu.opcodeTable[0x19] = func(_, _ uint8) { cpu.add16HL(cpu.getDE()); cpu.cycle += 8 }
	cpu.opcodeTable[0x29] = func(_, _ uint8) { cpu.add16HL(cpu.getHL()); cpu.cycle += 8 }
	cpu.opcodeTable[0x39] = func(_, _ uint8) { cpu.add16HL(cpu.stackPointer); cpu.cycle += 8 }
	cpu.opcodeTable[0xE8] = func(value, _ uint8) { cpu.add16SP(value); cpu.cycle += 16 }

	cpu.opcodeTable[0x03] = func(_, _ uint8) { cpu.setBC(cpu.getBC() + 1); cpu.cycle += 8 }
	cpu.opcodeTable[0x13] = func(_, _ uint8) { cpu.setDE(cpu.getDE() + 1); cpu.cycle += 8 }
	cpu.opcodeTable[0x23] = func(_, _ uint8) { cpu.setHL(cpu.getHL() + 1); cpu.cycle += 8 }

	cpu.opcodeTable[0x33] = func(_, _ uint8) { cpu.stackPointer = uint16(cpu.stackPointer + 1); cpu.cycle += 8 }
	cpu.opcodeTable[0x3B] = func(_, _ uint8) { cpu.stackPointer = uint16(cpu.stackPointer - 1); cpu.cycle += 8 }
	cpu.opcodeTable[0x0B] = func(_, _ uint8) { cpu.setBC(cpu.getBC() - 1); cpu.cycle += 8 }
	cpu.opcodeTable[0x1B] = func(_, _ uint8) { cpu.setDE(cpu.getDE() - 1); cpu.cycle += 8 }
	cpu.opcodeTable[0x2B] = func(_, _ uint8) { cpu.setHL(cpu.getHL() - 1); cpu.cycle += 8 }

}
