package cpu

func (cpu *CPU) rlc(val uint8) uint8 {
	flagC := val&0x80 != 0
	val = val << 1
	if flagC {
		val |= 0x01
	}

	cpu.f = 0
	if val == 0 {
		cpu.f |= 0x80
	}
	if flagC {
		cpu.f |= 0x10
	}
	return val
}

func (cpu *CPU) rl(val uint8) uint8 {
	flagC := val&0x80 != 0
	oldCarry := cpu.f&0x10 != 0
	val = val << 1
	cpu.f = 0
	if flagC {
		cpu.f |= 0x10
	}
	if oldCarry {
		val |= 0x01
	}
	if val == 0 {
		cpu.f |= 0x80
	}
	return val
}

func (cpu *CPU) rrc(val uint8) uint8 {
	flagC := val&0x01 != 0
	val = val >> 1
	if flagC {
		val |= 0x80
	}
	cpu.f = 0
	if val == 0 {
		cpu.f |= 0x80
	}
	if flagC {
		cpu.f |= 0x10
	}
	return val
}

func (cpu *CPU) rr(val uint8) uint8 {
	flagC := val&0x01 != 0
	oldCarry := cpu.f&0x10 != 0
	val = val >> 1
	cpu.f = 0
	if flagC {
		cpu.f |= 0x10
	}
	if oldCarry {
		val |= 0x80
	}
	if val == 0 {
		cpu.f |= 0x80
	}
	return val
}

func (cpu *CPU) sla(val uint8) uint8 {
	flagC := val&0x80 != 0
	val = val << 1
	cpu.f = 0
	if val == 0 {
		cpu.f |= 0x80
	}
	if flagC {
		cpu.f |= 0x10
	}
	return val
}

func (cpu *CPU) sra(val uint8) uint8 {
	msb := val&0x80 != 0
	flagC := val&0x01 != 0
	val = val >> 1
	cpu.f = 0
	if msb {
		val |= 0x80
	}
	if val == 0 {
		cpu.f |= 0x80
	}
	if flagC {
		cpu.f |= 0x10
	}

	return val
}

func (cpu *CPU) srl(val uint8) uint8 {
	flagC := val&0x01 != 0
	val = val >> 1
	cpu.f = 0
	if val == 0 {
		cpu.f |= 0x80
	}
	if flagC {
		cpu.f |= 0x10
	}

	return val
}

func (cpu *CPU) initRotateShiftOpCode() {
	cpu.opcodeTable[0x07] = func(_, _ uint8) {
		flagC := cpu.a&0x80 != 0
		cpu.a = cpu.a << 1
		if flagC {
			cpu.a |= 0x01
		}
		cpu.f = 0
		if cpu.a == 0 {
			cpu.f |= 0x80
		}
		if flagC {
			cpu.f |= 0x10
		}
		cpu.cycle += 4
	}
	cpu.opcodeTable[0x17] = func(_, _ uint8) {
		flagC := cpu.a&0x80 != 0
		oldCarry := cpu.f&0x10 != 0
		cpu.a = cpu.a << 1
		cpu.f = 0
		if flagC {
			cpu.f |= 0x10
		}

		if oldCarry {
			cpu.a |= 0x01
		}
		if cpu.a == 0 {
			cpu.f |= 0x80
		}
		cpu.cycle += 4
	}

	cpu.opcodeTable[0x0F] = func(_, _ uint8) {
		flagC := cpu.a&0x01 != 0
		cpu.a = cpu.a >> 1
		if flagC {
			cpu.a |= 0x80
		}

		cpu.f = 0
		if cpu.a == 0 {
			cpu.f |= 0x80
		}
		if flagC {
			cpu.f |= 0x10
		}
		cpu.cycle += 4
	}

	cpu.opcodeTable[0x1F] = func(_, _ uint8) {
		flagC := cpu.a&0x01 != 0
		oldCarry := cpu.f&0x10 != 0
		cpu.a = cpu.a >> 1
		cpu.f = 0
		if flagC {
			cpu.f |= 0x10
		}

		if oldCarry {
			cpu.a |= 0x80
		}
		if cpu.a == 0 {
			cpu.f |= 0x80
		}
		cpu.cycle += 4
	}

	cpu.opcodeTablePrefixed[0x07] = func(_, _ uint8) {
		cpu.a = cpu.rlc(cpu.a)
		cpu.cycle += 8
	}

	cpu.opcodeTablePrefixed[0x00] = func(_, _ uint8) {
		cpu.b = cpu.rlc(cpu.b)
		cpu.cycle += 8
	}

	cpu.opcodeTablePrefixed[0x01] = func(_, _ uint8) {
		cpu.c = cpu.rlc(cpu.c)
		cpu.cycle += 8
	}

	cpu.opcodeTablePrefixed[0x02] = func(_, _ uint8) {
		cpu.d = cpu.rlc(cpu.d)
		cpu.cycle += 8
	}

	cpu.opcodeTablePrefixed[0x03] = func(_, _ uint8) {
		cpu.e = cpu.rlc(cpu.e)
		cpu.cycle += 8
	}

	cpu.opcodeTablePrefixed[0x04] = func(_, _ uint8) {
		cpu.h = cpu.rlc(cpu.h)
		cpu.cycle += 8
	}

	cpu.opcodeTablePrefixed[0x05] = func(_, _ uint8) {
		cpu.l = cpu.rlc(cpu.l)
		cpu.cycle += 8
	}

	cpu.opcodeTablePrefixed[0x06] = func(_, _ uint8) {
		addr := uint16(cpu.h)<<8 | uint16(cpu.l)
		cpu.memory[addr] = cpu.rlc(cpu.memory[addr])
		cpu.cycle += 16
	}

	cpu.opcodeTablePrefixed[0x17] = func(_, _ uint8) {
		cpu.a = cpu.rl(cpu.a)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x10] = func(_, _ uint8) {
		cpu.b = cpu.rl(cpu.b)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x11] = func(_, _ uint8) {
		cpu.c = cpu.rl(cpu.c)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x12] = func(_, _ uint8) {
		cpu.d = cpu.rl(cpu.d)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x13] = func(_, _ uint8) {
		cpu.e = cpu.rl(cpu.e)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x14] = func(_, _ uint8) {
		cpu.h = cpu.rl(cpu.h)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x15] = func(_, _ uint8) {
		cpu.l = cpu.rl(cpu.l)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x16] = func(_, _ uint8) {
		addr := uint16(cpu.h)<<8 | uint16(cpu.l)
		cpu.memory[addr] = cpu.rl(cpu.memory[addr])
		cpu.cycle += 16
	}

	cpu.opcodeTablePrefixed[0x0F] = func(_, _ uint8) {
		cpu.a = cpu.rrc(cpu.a)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x08] = func(_, _ uint8) {
		cpu.b = cpu.rrc(cpu.b)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x09] = func(_, _ uint8) {
		cpu.c = cpu.rrc(cpu.c)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x0A] = func(_, _ uint8) {
		cpu.d = cpu.rrc(cpu.d)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x0B] = func(_, _ uint8) {
		cpu.e = cpu.rrc(cpu.e)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x0C] = func(_, _ uint8) {
		cpu.h = cpu.rrc(cpu.h)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x0D] = func(_, _ uint8) {
		cpu.l = cpu.rrc(cpu.l)
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x0E] = func(_, _ uint8) {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.rrc(cpu.memory[hl]))
		cpu.cycle += 16
	}
	cpu.opcodeTablePrefixed[0x1F] = func(_, _ uint8) { cpu.a = cpu.rr(cpu.a); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x18] = func(_, _ uint8) { cpu.b = cpu.rr(cpu.b); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x19] = func(_, _ uint8) { cpu.c = cpu.rr(cpu.c); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x1A] = func(_, _ uint8) { cpu.d = cpu.rr(cpu.d); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x1B] = func(_, _ uint8) { cpu.e = cpu.rr(cpu.e); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x1C] = func(_, _ uint8) { cpu.h = cpu.rr(cpu.h); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x1D] = func(_, _ uint8) { cpu.l = cpu.rr(cpu.l); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x1E] = func(_, _ uint8) {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.rr(cpu.memory[hl]))
		cpu.cycle += 16
	}
	cpu.opcodeTablePrefixed[0x27] = func(_, _ uint8) { cpu.a = cpu.sla(cpu.a); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x20] = func(_, _ uint8) { cpu.b = cpu.sla(cpu.b); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x21] = func(_, _ uint8) { cpu.c = cpu.sla(cpu.c); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x22] = func(_, _ uint8) { cpu.d = cpu.sla(cpu.d); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x23] = func(_, _ uint8) { cpu.e = cpu.sla(cpu.e); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x24] = func(_, _ uint8) { cpu.h = cpu.sla(cpu.h); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x25] = func(_, _ uint8) { cpu.l = cpu.sla(cpu.l); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x26] = func(_, _ uint8) {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.sla(cpu.memory[hl]))
		cpu.cycle += 16
	}
	cpu.opcodeTablePrefixed[0x2F] = func(_, _ uint8) { cpu.a = cpu.sra(cpu.a); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x28] = func(_, _ uint8) { cpu.b = cpu.sra(cpu.b); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x29] = func(_, _ uint8) { cpu.c = cpu.sra(cpu.c); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x2A] = func(_, _ uint8) { cpu.d = cpu.sra(cpu.d); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x2B] = func(_, _ uint8) { cpu.e = cpu.sra(cpu.e); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x2C] = func(_, _ uint8) { cpu.h = cpu.sra(cpu.h); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x2D] = func(_, _ uint8) { cpu.l = cpu.sra(cpu.l); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x2E] = func(_, _ uint8) {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.sra(cpu.memory[hl]))
		cpu.cycle += 16
	}
	cpu.opcodeTablePrefixed[0x3F] = func(_, _ uint8) { cpu.a = cpu.srl(cpu.a); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x38] = func(_, _ uint8) { cpu.b = cpu.srl(cpu.b); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x39] = func(_, _ uint8) { cpu.c = cpu.srl(cpu.c); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x3A] = func(_, _ uint8) { cpu.d = cpu.srl(cpu.d); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x3B] = func(_, _ uint8) { cpu.e = cpu.srl(cpu.e); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x3C] = func(_, _ uint8) { cpu.h = cpu.srl(cpu.h); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x3D] = func(_, _ uint8) { cpu.l = cpu.srl(cpu.l); cpu.cycle += 8 }
	cpu.opcodeTablePrefixed[0x3E] = func(_, _ uint8) {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.srl(cpu.memory[hl]))
		cpu.cycle += 16
	}
}
