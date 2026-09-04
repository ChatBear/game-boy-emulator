package cpu

func (cpu *CPU) restart(value uint8) {
	cpu.stackPointer--
	cpu.writeMemory(cpu.stackPointer, uint8(cpu.programCounter>>8))
	cpu.stackPointer--
	cpu.writeMemory(cpu.stackPointer, uint8(cpu.programCounter&0xFF))
	cpu.programCounter = uint16(value)
}

func (cpu *CPU) initRestartOpCode() {
	cpu.opcodeTable[0xC7] = func() { cpu.restart(0x00); cpu.UpdateTimer(32) }
	cpu.opcodeTable[0xCF] = func() { cpu.restart(0x08); cpu.UpdateTimer(32) }
	cpu.opcodeTable[0xD7] = func() { cpu.restart(0x10); cpu.UpdateTimer(32) }
	cpu.opcodeTable[0xDF] = func() { cpu.restart(0x18); cpu.UpdateTimer(32) }
	cpu.opcodeTable[0xE7] = func() { cpu.restart(0x20); cpu.UpdateTimer(32) }
	cpu.opcodeTable[0xEF] = func() { cpu.restart(0x28); cpu.UpdateTimer(32) }
	cpu.opcodeTable[0xF7] = func() { cpu.restart(0x30); cpu.UpdateTimer(32) }
	cpu.opcodeTable[0xFF] = func() { cpu.restart(0x38); cpu.UpdateTimer(32) }
}
