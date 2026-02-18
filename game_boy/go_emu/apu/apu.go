package apu

type APU struct {
	NR10 uint8 // 0xFF10 — Ch1 sweep
	NR11 uint8 // 0xFF11 — Ch1 duty/length
	NR12 uint8 // 0xFF12 — Ch1 envelope
	NR13 uint8 // 0xFF13 — Ch1 freq low
	NR14 uint8 // 0xFF14 — Ch1 freq high + trigger

	NR21 uint8 // 0xFF16 — Ch2 duty/length
	NR22 uint8 // 0xFF17 — Ch2 envelope
	NR23 uint8 // 0xFF18 — Ch2 freq low
	NR24 uint8 // 0xFF19 — Ch2 freq high + trigger

	NR30 uint8 // 0xFF1A — Ch3 DAC on/off
	NR31 uint8 // 0xFF1B — Ch3 length
	NR32 uint8 // 0xFF1C — Ch3 output level
	NR33 uint8 // 0xFF1D — Ch3 freq low
	NR34 uint8 // 0xFF1E — Ch3 freq high + trigger

	NR41 uint8 // 0xFF20 — Ch4 length
	NR42 uint8 // 0xFF21 — Ch4 envelope
	NR43 uint8 // 0xFF22 — Ch4 polynomial counter
	NR44 uint8 // 0xFF23 — Ch4 trigger

	NR50 uint8 // 0xFF24 — Master volume
	NR51 uint8 // 0xFF25 — Panning
	NR52 uint8 // 0xFF26 — Sound on/off

	WaveRAM [16]uint8 // 0xFF30–0xFF3F
}

func (apu *APU) WriteRegister(addr uint16, value uint8) {
	switch addr {
	case 0xFF10:
		apu.NR10 = value
	case 0xFF11:
		apu.NR11 = value
	case 0xFF12:
		apu.NR12 = value
	case 0xFF13:
		apu.NR13 = value
	case 0xFF14:
		apu.NR14 = value
	case 0xFF16:
		apu.NR21 = value
	case 0xFF17:
		apu.NR22 = value
	case 0xFF18:
		apu.NR23 = value
	case 0xFF19:
		apu.NR24 = value
	case 0xFF1A:
		apu.NR30 = value
	case 0xFF1B:
		apu.NR31 = value
	case 0xFF1C:
		apu.NR32 = value
	case 0xFF1D:
		apu.NR33 = value
	case 0xFF1E:
		apu.NR34 = value
	case 0xFF20:
		apu.NR41 = value
	case 0xFF21:
		apu.NR42 = value
	case 0xFF22:
		apu.NR43 = value
	case 0xFF23:
		apu.NR44 = value
	case 0xFF24:
		apu.NR50 = value
	case 0xFF25:
		apu.NR51 = value
	case 0xFF26:
		apu.NR52 = value
	default:
		if addr >= 0xFF30 && addr <= 0xFF3F {
			apu.WaveRAM[addr-0xFF30] = value
		}
	}
}
