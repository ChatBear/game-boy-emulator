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
	cpu.opcodeTable[0x07] = func() {
		flagC := cpu.a&0x80 != 0
		cpu.a = cpu.a << 1
		if flagC {
			cpu.a |= 0x01
		}
		cpu.f = 0
		if flagC {
			cpu.f |= 0x10
		}
		cpu.UpdateTimer(4)
	}
	cpu.opcodeTable[0x17] = func() {
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
		cpu.UpdateTimer(4)
	}

	cpu.opcodeTable[0x0F] = func() {
		flagC := cpu.a&0x01 != 0
		cpu.a = cpu.a >> 1
		if flagC {
			cpu.a |= 0x80
		}

		cpu.f = 0
		if flagC {
			cpu.f |= 0x10
		}
		cpu.UpdateTimer(4)
	}

	cpu.opcodeTable[0x1F] = func() {
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
		cpu.UpdateTimer(4)
	}

	cpu.opcodeTablePrefixed[0x07] = func() {
		cpu.a = cpu.rlc(cpu.a)
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTablePrefixed[0x00] = func() {
		cpu.b = cpu.rlc(cpu.b)
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTablePrefixed[0x01] = func() {
		cpu.c = cpu.rlc(cpu.c)
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTablePrefixed[0x02] = func() {
		cpu.d = cpu.rlc(cpu.d)
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTablePrefixed[0x03] = func() {
		cpu.e = cpu.rlc(cpu.e)
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTablePrefixed[0x04] = func() {
		cpu.h = cpu.rlc(cpu.h)
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTablePrefixed[0x05] = func() {
		cpu.l = cpu.rlc(cpu.l)
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTablePrefixed[0x06] = func() {
		addr := uint16(cpu.h)<<8 | uint16(cpu.l)
		cpu.memory[addr] = cpu.rlc(cpu.memory[addr])
		cpu.UpdateTimer(16)
	}

	cpu.opcodeTablePrefixed[0x17] = func() {
		cpu.a = cpu.rl(cpu.a)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x10] = func() {
		cpu.b = cpu.rl(cpu.b)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x11] = func() {
		cpu.c = cpu.rl(cpu.c)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x12] = func() {
		cpu.d = cpu.rl(cpu.d)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x13] = func() {
		cpu.e = cpu.rl(cpu.e)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x14] = func() {
		cpu.h = cpu.rl(cpu.h)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x15] = func() {
		cpu.l = cpu.rl(cpu.l)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x16] = func() {
		addr := uint16(cpu.h)<<8 | uint16(cpu.l)
		cpu.memory[addr] = cpu.rl(cpu.memory[addr])
		cpu.UpdateTimer(16)
	}

	cpu.opcodeTablePrefixed[0x0F] = func() {
		cpu.a = cpu.rrc(cpu.a)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x08] = func() {
		cpu.b = cpu.rrc(cpu.b)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x09] = func() {
		cpu.c = cpu.rrc(cpu.c)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x0A] = func() {
		cpu.d = cpu.rrc(cpu.d)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x0B] = func() {
		cpu.e = cpu.rrc(cpu.e)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x0C] = func() {
		cpu.h = cpu.rrc(cpu.h)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x0D] = func() {
		cpu.l = cpu.rrc(cpu.l)
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x0E] = func() {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.rrc(cpu.memory[hl]))
		cpu.UpdateTimer(16)
	}
	cpu.opcodeTablePrefixed[0x1F] = func() { cpu.a = cpu.rr(cpu.a); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x18] = func() { cpu.b = cpu.rr(cpu.b); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x19] = func() { cpu.c = cpu.rr(cpu.c); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x1A] = func() { cpu.d = cpu.rr(cpu.d); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x1B] = func() { cpu.e = cpu.rr(cpu.e); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x1C] = func() { cpu.h = cpu.rr(cpu.h); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x1D] = func() { cpu.l = cpu.rr(cpu.l); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x1E] = func() {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.rr(cpu.memory[hl]))
		cpu.UpdateTimer(16)
	}
	cpu.opcodeTablePrefixed[0x27] = func() { cpu.a = cpu.sla(cpu.a); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x20] = func() { cpu.b = cpu.sla(cpu.b); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x21] = func() { cpu.c = cpu.sla(cpu.c); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x22] = func() { cpu.d = cpu.sla(cpu.d); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x23] = func() { cpu.e = cpu.sla(cpu.e); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x24] = func() { cpu.h = cpu.sla(cpu.h); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x25] = func() { cpu.l = cpu.sla(cpu.l); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x26] = func() {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.sla(cpu.memory[hl]))
		cpu.UpdateTimer(16)
	}
	cpu.opcodeTablePrefixed[0x2F] = func() { cpu.a = cpu.sra(cpu.a); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x28] = func() { cpu.b = cpu.sra(cpu.b); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x29] = func() { cpu.c = cpu.sra(cpu.c); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x2A] = func() { cpu.d = cpu.sra(cpu.d); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x2B] = func() { cpu.e = cpu.sra(cpu.e); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x2C] = func() { cpu.h = cpu.sra(cpu.h); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x2D] = func() { cpu.l = cpu.sra(cpu.l); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x2E] = func() {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.sra(cpu.memory[hl]))
		cpu.UpdateTimer(16)
	}
	cpu.opcodeTablePrefixed[0x3F] = func() { cpu.a = cpu.srl(cpu.a); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x38] = func() { cpu.b = cpu.srl(cpu.b); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x39] = func() { cpu.c = cpu.srl(cpu.c); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x3A] = func() { cpu.d = cpu.srl(cpu.d); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x3B] = func() { cpu.e = cpu.srl(cpu.e); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x3C] = func() { cpu.h = cpu.srl(cpu.h); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x3D] = func() { cpu.l = cpu.srl(cpu.l); cpu.UpdateTimer(8) }
	cpu.opcodeTablePrefixed[0x3E] = func() {
		hl := cpu.getHL()
		cpu.writeMemory(hl, cpu.srl(cpu.memory[hl]))
		cpu.UpdateTimer(16)
	}
}
