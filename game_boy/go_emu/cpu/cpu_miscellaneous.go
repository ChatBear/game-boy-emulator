package cpu

func (cpu *CPU) initMiscellaneousOpCodes() {
	cpu.opcodeTablePrefixed[0x37] = func(_, _ uint8) {
		cpu.a = (cpu.a >> 4) | (cpu.a << 4)
		cpu.f = 0
		if cpu.a == 0 {
			cpu.f |= 0x80
		}
		cpu.cycle += 8
	}

	cpu.opcodeTablePrefixed[0x30] = func(_, _ uint8) {
		cpu.b = (cpu.b >> 4) | (cpu.b << 4)
		cpu.f = 0
		if cpu.b == 0 {
			cpu.f |= 0x80
		}
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x31] = func(_, _ uint8) {
		cpu.c = (cpu.c >> 4) | (cpu.c << 4)
		cpu.f = 0
		if cpu.c == 0 {
			cpu.f |= 0x80
		}
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x32] = func(_, _ uint8) {
		cpu.d = (cpu.d >> 4) | (cpu.d << 4)
		cpu.f = 0
		if cpu.d == 0 {
			cpu.f |= 0x80
		}
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x33] = func(_, _ uint8) {
		cpu.e = (cpu.e >> 4) | (cpu.e << 4)
		cpu.f = 0
		if cpu.e == 0 {
			cpu.f |= 0x80
		}
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x34] = func(_, _ uint8) {
		cpu.h = (cpu.h >> 4) | (cpu.h << 4)
		cpu.f = 0
		if cpu.h == 0 {
			cpu.f |= 0x80
		}
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x35] = func(_, _ uint8) {
		cpu.l = (cpu.l >> 4) | (cpu.l << 4)
		cpu.f = 0
		if cpu.l == 0 {
			cpu.f |= 0x80
		}
		cpu.cycle += 8
	}
	cpu.opcodeTablePrefixed[0x36] = func(_, _ uint8) {
		hl := cpu.getHL()
		val := cpu.memory[hl]
		val = (val >> 4) | (val << 4)
		cpu.f = 0
		if val == 0 {
			cpu.f |= 0x80
		}
		cpu.writeMemory(hl, val)
		cpu.cycle += 16
	}

	cpu.opcodeTable[0x27] = func(_, _ uint8) {
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
		cpu.cycle += 4
	}

	cpu.opcodeTable[0x2F] = func(_, _ uint8) {
		cpu.a = ^cpu.a
		cpu.f |= 0x40
		cpu.f |= 0x20
		cpu.cycle += 4
	}

	cpu.opcodeTable[0x3F] = func(_, _ uint8) {
		cpu.f ^= 0x10
		cpu.f = cpu.f &^ 0x60
		cpu.cycle += 4
	}

	cpu.opcodeTable[0x37] = func(_, _ uint8) {
		cpu.f = cpu.f &^ 0x60
		cpu.f |= 0x10
		cpu.cycle += 4
	}

	cpu.opcodeTable[0x00] = func(_, _ uint8) {
		cpu.cycle += 4
	}

	cpu.opcodeTable[0x76] = func(_, _ uint8) {
		cpu.halt = true
		cpu.cycle += 4
	}

	cpu.opcodeTable[0x10] = func(_, _ uint8) {
		cpu.programCounter++
		cpu.stopped = true
		cpu.cycle += 4
	}

	cpu.opcodeTable[0xF3] = func(_, _ uint8) {
		// Main Loop -> disable the ime on the next instruction
		cpu.pendingDisableIME = true
		cpu.cycle += 4
	}

	cpu.opcodeTable[0xFB] = func(_, _ uint8) {
		// Main Loop -> disable the ime on the next instruction -> Same thing modify cme on the next instruction :
		// if cpu.pendingDisableIME {
		// 	cpu.ime = true
		// 	cpu.pendingDisableIME = true
		// }
		cpu.pendingDisableIME = false
		cpu.cycle += 4
	}
}
