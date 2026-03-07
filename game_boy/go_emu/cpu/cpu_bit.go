package cpu

func (cpu *CPU) bit(b uint8, val uint8) {
	cpu.f &= 0x10 // preserve C
	cpu.f |= 0x20 // set H
	if val&(1<<b) == 0 {
		cpu.f |= 0x80 // set Z
	}
}
func (cpu *CPU) set(b uint8, val uint8) uint8 {
	val |= 1 << b
	return val
}

func (cpu *CPU) res(b uint8, val uint8) uint8 {
	val &^= 1 << b
	return val
}
func (cpu *CPU) initBitOpCode() {
	for b := range uint8(8) {
		bit := b
		base := 0x40 + int(bit)*8

		cpu.opcodeTablePrefixed[base+0] = func(_, _ uint8) { cpu.bit(bit, cpu.b); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+1] = func(_, _ uint8) { cpu.bit(bit, cpu.c); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+2] = func(_, _ uint8) { cpu.bit(bit, cpu.d); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+3] = func(_, _ uint8) { cpu.bit(bit, cpu.e); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+4] = func(_, _ uint8) { cpu.bit(bit, cpu.h); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+5] = func(_, _ uint8) { cpu.bit(bit, cpu.l); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+6] = func(_, _ uint8) { cpu.bit(bit, cpu.memory[cpu.getHL()]); cpu.cycle += 16 }
		cpu.opcodeTablePrefixed[base+7] = func(_, _ uint8) { cpu.bit(bit, cpu.a); cpu.cycle += 8 }
	}

	for b := range uint8(8) {
		bit := b
		base := 0xC0 + int(bit)*8

		cpu.opcodeTablePrefixed[base+0] = func(_, _ uint8) { cpu.b = cpu.set(bit, cpu.b); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+1] = func(_, _ uint8) { cpu.c = cpu.set(bit, cpu.c); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+2] = func(_, _ uint8) { cpu.d = cpu.set(bit, cpu.d); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+3] = func(_, _ uint8) { cpu.e = cpu.set(bit, cpu.e); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+4] = func(_, _ uint8) { cpu.h = cpu.set(bit, cpu.h); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+5] = func(_, _ uint8) { cpu.l = cpu.set(bit, cpu.l); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+6] = func(_, _ uint8) {
			addr := cpu.getHL()
			cpu.memory[addr] = cpu.set(bit, cpu.memory[addr])
			cpu.cycle += 16
		}
		cpu.opcodeTablePrefixed[base+7] = func(_, _ uint8) { cpu.a = cpu.set(bit, cpu.a); cpu.cycle += 8 }
	}

	for b := range uint8(8) {
		bit := b
		base := 0x80 + int(bit)*8

		cpu.opcodeTablePrefixed[base+0] = func(_, _ uint8) { cpu.b = cpu.res(bit, cpu.b); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+1] = func(_, _ uint8) { cpu.c = cpu.res(bit, cpu.c); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+2] = func(_, _ uint8) { cpu.d = cpu.res(bit, cpu.d); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+3] = func(_, _ uint8) { cpu.e = cpu.res(bit, cpu.e); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+4] = func(_, _ uint8) { cpu.h = cpu.res(bit, cpu.h); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+5] = func(_, _ uint8) { cpu.l = cpu.res(bit, cpu.l); cpu.cycle += 8 }
		cpu.opcodeTablePrefixed[base+6] = func(_, _ uint8) {
			addr := cpu.getHL()
			cpu.memory[addr] = cpu.res(bit, cpu.memory[addr])
			cpu.cycle += 16
		}
		cpu.opcodeTablePrefixed[base+7] = func(_, _ uint8) { cpu.a = cpu.res(bit, cpu.a); cpu.cycle += 8 }
	}
}
