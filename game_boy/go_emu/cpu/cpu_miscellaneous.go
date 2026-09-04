package cpu

func (cpu *CPU) initMiscellaneousOpCodes() {
	cpu.opcodeTablePrefixed[0x37] = func() {
		cpu.a = (cpu.a >> 4) | (cpu.a << 4)
		cpu.f = 0
		if cpu.a == 0 {
			cpu.f |= 0x80
		}
		cpu.UpdateTimer(8)
	}

	cpu.opcodeTablePrefixed[0x30] = func() {
		cpu.b = (cpu.b >> 4) | (cpu.b << 4)
		cpu.f = 0
		if cpu.b == 0 {
			cpu.f |= 0x80
		}
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x31] = func() {
		cpu.c = (cpu.c >> 4) | (cpu.c << 4)
		cpu.f = 0
		if cpu.c == 0 {
			cpu.f |= 0x80
		}
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x32] = func() {
		cpu.d = (cpu.d >> 4) | (cpu.d << 4)
		cpu.f = 0
		if cpu.d == 0 {
			cpu.f |= 0x80
		}
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x33] = func() {
		cpu.e = (cpu.e >> 4) | (cpu.e << 4)
		cpu.f = 0
		if cpu.e == 0 {
			cpu.f |= 0x80
		}
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x34] = func() {
		cpu.h = (cpu.h >> 4) | (cpu.h << 4)
		cpu.f = 0
		if cpu.h == 0 {
			cpu.f |= 0x80
		}
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x35] = func() {
		cpu.l = (cpu.l >> 4) | (cpu.l << 4)
		cpu.f = 0
		if cpu.l == 0 {
			cpu.f |= 0x80
		}
		cpu.UpdateTimer(8)
	}
	cpu.opcodeTablePrefixed[0x36] = func() {
		hl := cpu.getHL()
		val := cpu.memory[hl]
		val = (val >> 4) | (val << 4)
		cpu.f = 0
		if val == 0 {
			cpu.f |= 0x80
		}
		cpu.writeMemory(hl, val)
		cpu.UpdateTimer(16)
	}

	cpu.opcodeTable[0x27] = func() {
		correction := uint8(0)
		nFlag := cpu.f&0x40 != 0
		carry := cpu.f&0x10 != 0
		halfCarry := cpu.f&0x20 != 0

		if halfCarry || (!nFlag && cpu.a&0x0F > 0x09) {
			correction |= 0x06
		}
		if carry || (!nFlag && cpu.a > 0x99) {
			correction |= 0x60
			carry = true
		}

		if nFlag {
			cpu.a -= correction
		} else {
			cpu.a += correction
		}

		cpu.f &= 0x40 // preserve N only
		if cpu.a == 0 {
			cpu.f |= 0x80
		}
		if carry {
			cpu.f |= 0x10
		}
		cpu.UpdateTimer(4)
	}

	cpu.opcodeTable[0x2F] = func() {
		cpu.a = ^cpu.a
		cpu.f |= 0x40
		cpu.f |= 0x20
		cpu.UpdateTimer(4)
	}

	cpu.opcodeTable[0x3F] = func() {
		cpu.f ^= 0x10
		cpu.f = cpu.f &^ 0x60
		cpu.UpdateTimer(4)
	}

	cpu.opcodeTable[0x37] = func() {
		cpu.f = cpu.f &^ 0x60
		cpu.f |= 0x10
		cpu.UpdateTimer(4)
	}

	cpu.opcodeTable[0x00] = func() {
		cpu.UpdateTimer(4)
	}

	cpu.opcodeTable[0x76] = func() {
		cpu.halt = true
		cpu.UpdateTimer(4)
	}

	cpu.opcodeTable[0x10] = func() {
		cpu.stopped = true
		cpu.UpdateTimer(4)
	}

	cpu.opcodeTable[0xF3] = func() {
		// Main Loop -> disable the ime on the next instruction
		cpu.pendingEnableIME = false
		cpu.ime = false
		cpu.UpdateTimer(4)
	}

	cpu.opcodeTable[0xFB] = func() {
		// Main Loop -> disable the ime on the next instruction -> Same thing modify cme on the next instruction :
		cpu.pendingEnableIME = true
		cpu.UpdateTimer(4)
	}
}
