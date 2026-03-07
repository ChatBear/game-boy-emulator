package cpu

func (cpu *CPU) restart(value uint8) {
	cpu.stackPointer--
	cpu.writeMemory(cpu.stackPointer, uint8(cpu.programCounter>>8))
	cpu.stackPointer--
	cpu.writeMemory(cpu.stackPointer, uint8(cpu.programCounter&0xFF))
	cpu.programCounter = uint16(value)
}

func (cpu *CPU) initRestartOpCode() {
	cpu.opcodeTable[0xC7] = func(_, _ uint8) { cpu.restart(0x00); cpu.cycle += 32 }
	cpu.opcodeTable[0xCF] = func(_, _ uint8) { cpu.restart(0x08); cpu.cycle += 32 }
	cpu.opcodeTable[0xD7] = func(_, _ uint8) { cpu.restart(0x10); cpu.cycle += 32 }
	cpu.opcodeTable[0xDF] = func(_, _ uint8) { cpu.restart(0x18); cpu.cycle += 32 }
	cpu.opcodeTable[0xE7] = func(_, _ uint8) { cpu.restart(0x20); cpu.cycle += 32 }
	cpu.opcodeTable[0xEF] = func(_, _ uint8) { cpu.restart(0x28); cpu.cycle += 32 }
	cpu.opcodeTable[0xF7] = func(_, _ uint8) { cpu.restart(0x30); cpu.cycle += 32 }
	cpu.opcodeTable[0xFF] = func(_, _ uint8) { cpu.restart(0x38); cpu.cycle += 32 }
}
