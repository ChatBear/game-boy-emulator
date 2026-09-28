use crate::cpu::Cpu;
use crate::gb::GameBoy;
use crate::memory::Memory;

pub type OpCodeExecution = fn(&mut GameBoy, u8) -> u128;

pub struct OpCodeTable {
    pub opcodes: [OpCodeExecution; 256],
}

const FLAG_Z: u8 = 0x80;
const FLAG_N: u8 = 0x40;
const FLAG_H: u8 = 0x20;
const FLAG_C: u8 = 0x10;

// const REG8: [&str; 8] = ["B", "C", "D", "E", "H", "L", "(HL)", "A"];

pub struct Decoded {
    pub operation: &'static str,
    pub register: &'static str,
}

// pub fn decode(code: u8) -> Option<Decoded> {
//     let register = REG8[((code >> 3) & 7) as usize];
//     Some(Decoded {
//         operation,
//         register,
//     })
// }

fn fetch8(cpu: &mut Cpu, mem: &Memory) -> u8 {
    let value = mem.read(cpu.pc);
    cpu.pc = cpu.pc.wrapping_add(1);
    value
}

fn fetch16(cpu: &mut Cpu, mem: &Memory) -> u16 {
    let lo = fetch8(cpu, mem) as u16;
    let hi = fetch8(cpu, mem) as u16;
    lo | (hi << 8)
}

pub(crate) fn push16(cpu: &mut Cpu, mem: &mut Memory, value: u16) {
    cpu.sp = cpu.sp.wrapping_sub(1);
    mem.write(cpu.sp, (value >> 8) as u8);
    cpu.sp = cpu.sp.wrapping_sub(1);
    mem.write(cpu.sp, value as u8);
}

fn pop16(cpu: &mut Cpu, mem: &Memory) -> u16 {
    let lo = mem.read(cpu.sp) as u16;
    cpu.sp = cpu.sp.wrapping_add(1);
    let hi = mem.read(cpu.sp) as u16;
    cpu.sp = cpu.sp.wrapping_add(1);
    lo | (hi << 8)
}

fn flag(cpu: &Cpu, mask: u8) -> bool {
    cpu.f & mask != 0
}

fn set_flags(cpu: &mut Cpu, z: bool, n: bool, h: bool, c: bool) {
    cpu.f = 0;
    if z {
        cpu.f |= FLAG_Z;
    }
    if n {
        cpu.f |= FLAG_N;
    }
    if h {
        cpu.f |= FLAG_H;
    }
    if c {
        cpu.f |= FLAG_C;
    }
}

fn r8(cpu: &Cpu, mem: &Memory, r: u8) -> u8 {
    match r & 7 {
        0 => cpu.b,
        1 => cpu.c,
        2 => cpu.d,
        3 => cpu.e,
        4 => cpu.h,
        5 => cpu.l,
        6 => mem.read(cpu.get_hl()),
        _ => cpu.a,
    }
}

fn set_r8(cpu: &mut Cpu, mem: &mut Memory, r: u8, value: u8) {
    match r & 7 {
        0 => cpu.b = value,
        1 => cpu.c = value,
        2 => cpu.d = value,
        3 => cpu.e = value,
        4 => cpu.h = value,
        5 => cpu.l = value,
        6 => mem.write(cpu.get_hl(), value),
        _ => cpu.a = value,
    }
}

fn r16(cpu: &Cpu, rp: u8) -> u16 {
    match rp & 3 {
        0 => cpu.get_bc(),
        1 => cpu.get_de(),
        2 => cpu.get_hl(),
        _ => cpu.sp,
    }
}

fn set_r16(cpu: &mut Cpu, rp: u8, value: u16) {
    match rp & 3 {
        0 => cpu.set_bc(value),
        1 => cpu.set_de(value),
        2 => cpu.set_hl(value),
        _ => cpu.sp = value,
    }
}

fn illegal(_gb: &mut GameBoy, op: u8) -> u128 {
    panic!("illegal opcode: 0x{op:02X}");
}

fn nop(_gb: &mut GameBoy, _op: u8) -> u128 {
    4
}

fn stop(gb: &mut GameBoy, _op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let _ = fetch8(cpu, mem);
    cpu.stopped = true;
    4
}

fn halt(gb: &mut GameBoy, _op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    cpu.halt = true;
    4
}

fn di(gb: &mut GameBoy, _op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    cpu.ime = false;
    cpu.ei_pending = false;
    4
}

fn ei(gb: &mut GameBoy, _op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    cpu.ei_pending = true;
    4
}

fn ld_r_r(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let dst = (op >> 3) & 7;
    let src = op & 7;
    let value = r8(cpu, mem, src);
    set_r8(cpu, mem, dst, value);
    if dst == 6 || src == 6 { 8 } else { 4 }
}

fn ld_r_n(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let dst = (op >> 3) & 7;
    let value = fetch8(cpu, mem);
    set_r8(cpu, mem, dst, value);
    if dst == 6 { 12 } else { 8 }
}

fn ld_rr_nn(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let value = fetch16(cpu, mem);
    set_r16(cpu, (op >> 4) & 3, value);
    12
}

fn ld_mem_a(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let addr = match op {
        0x02 => cpu.get_bc(),
        0x12 => cpu.get_de(),
        0x22 => {
            let hl = cpu.get_hl();
            cpu.set_hl(hl.wrapping_add(1));
            hl
        }
        0x32 => {
            let hl = cpu.get_hl();
            cpu.set_hl(hl.wrapping_sub(1));
            hl
        }
        0xEA => fetch16(cpu, mem),
        _ => unreachable!(),
    };
    mem.write(addr, cpu.a);
    if op == 0xEA { 16 } else { 8 }
}

fn ld_a_mem(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let addr = match op {
        0x0A => cpu.get_bc(),
        0x1A => cpu.get_de(),
        0x2A => {
            let hl = cpu.get_hl();
            cpu.set_hl(hl.wrapping_add(1));
            hl
        }
        0x3A => {
            let hl = cpu.get_hl();
            cpu.set_hl(hl.wrapping_sub(1));
            hl
        }
        0xFA => fetch16(cpu, mem),
        _ => unreachable!(),
    };
    cpu.a = mem.read(addr);
    if op == 0xFA { 16 } else { 8 }
}

fn ldh_n_a(gb: &mut GameBoy, _op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let n = fetch8(cpu, mem) as u16;
    mem.write(0xFF00 + n, cpu.a);
    12
}

fn ldh_a_n(gb: &mut GameBoy, _op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let n = fetch8(cpu, mem) as u16;
    cpu.a = mem.read(0xFF00 + n);
    12
}

fn ldh_c_a(gb: &mut GameBoy, _op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    mem.write(0xFF00 + cpu.c as u16, cpu.a);
    8
}

fn ldh_a_c(gb: &mut GameBoy, _op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    cpu.a = mem.read(0xFF00 + cpu.c as u16);
    8
}

fn ld_nn_sp(gb: &mut GameBoy, _op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let addr = fetch16(cpu, mem);
    mem.write(addr, cpu.sp as u8);
    mem.write(addr.wrapping_add(1), (cpu.sp >> 8) as u8);
    20
}

fn ld_sp_hl(gb: &mut GameBoy, _op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    cpu.sp = cpu.get_hl();
    8
}

fn add_sp_offset(cpu: &Cpu, offset: i8) -> (u16, bool, bool) {
    let n = offset as u16;
    let h = (cpu.sp & 0x0F) + (n & 0x0F) > 0x0F;
    let c = (cpu.sp & 0xFF) + (n & 0xFF) > 0xFF;
    (cpu.sp.wrapping_add(offset as u16), h, c)
}

fn ld_hl_sp_e(gb: &mut GameBoy, _op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let e = fetch8(cpu, mem) as i8;
    let (value, h, c) = add_sp_offset(cpu, e);
    cpu.set_hl(value);
    set_flags(cpu, false, false, h, c);
    12
}

fn add_sp_e(gb: &mut GameBoy, _op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let e = fetch8(cpu, mem) as i8;
    let (value, h, c) = add_sp_offset(cpu, e);
    cpu.sp = value;
    set_flags(cpu, false, false, h, c);
    16
}

fn inc8(cpu: &mut Cpu, value: u8) -> u8 {
    let result = value.wrapping_add(1);
    set_flags(
        cpu,
        result == 0,
        false,
        value & 0x0F == 0x0F,
        flag(cpu, FLAG_C),
    );
    result
}

fn dec8(cpu: &mut Cpu, value: u8) -> u8 {
    let result = value.wrapping_sub(1);
    set_flags(cpu, result == 0, true, value & 0x0F == 0, flag(cpu, FLAG_C));
    result
}

fn inc_r(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let r = (op >> 3) & 7;
    let value = inc8(cpu, r8(cpu, mem, r));
    set_r8(cpu, mem, r, value);
    if r == 6 { 12 } else { 4 }
}

fn dec_r(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let r = (op >> 3) & 7;
    let value = dec8(cpu, r8(cpu, mem, r));
    set_r8(cpu, mem, r, value);
    if r == 6 { 12 } else { 4 }
}

fn inc_rr(gb: &mut GameBoy, op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    let rp = (op >> 4) & 3;
    set_r16(cpu, rp, r16(cpu, rp).wrapping_add(1));
    8
}

fn dec_rr(gb: &mut GameBoy, op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    let rp = (op >> 4) & 3;
    set_r16(cpu, rp, r16(cpu, rp).wrapping_sub(1));
    8
}

fn add_hl_rr(gb: &mut GameBoy, op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    let hl = cpu.get_hl();
    let operand = r16(cpu, (op >> 4) & 3);
    let result = hl as u32 + operand as u32;
    cpu.set_hl(result as u16);
    set_flags(
        cpu,
        flag(cpu, FLAG_Z),
        false,
        (hl & 0x0FFF) + (operand & 0x0FFF) > 0x0FFF,
        result > 0xFFFF,
    );
    8
}

fn add_a(cpu: &mut Cpu, operand: u8) {
    let result = cpu.a as u16 + operand as u16;
    let h = (cpu.a & 0x0F) + (operand & 0x0F) > 0x0F;
    cpu.a = result as u8;
    set_flags(cpu, cpu.a == 0, false, h, result > 0xFF);
}

fn adc_a(cpu: &mut Cpu, operand: u8) {
    let carry = flag(cpu, FLAG_C) as u16;
    let result = cpu.a as u16 + operand as u16 + carry;
    let h = (cpu.a & 0x0F) as u16 + (operand & 0x0F) as u16 + carry > 0x0F;
    cpu.a = result as u8;
    set_flags(cpu, cpu.a == 0, false, h, result > 0xFF);
}

fn sub_a(cpu: &mut Cpu, operand: u8) {
    let (res, c) = cpu.a.overflowing_sub(operand);
    let h = (cpu.a & 0x0F) < (operand & 0x0F);
    cpu.a = res;
    set_flags(cpu, res == 0, true, h, c);
}

fn sbc_a(cpu: &mut Cpu, operand: u8) {
    let carry = flag(cpu, FLAG_C) as u8;
    let result = cpu.a as i16 - operand as i16 - carry as i16;
    let h = (cpu.a & 0x0F) < (operand & 0x0F) + carry;
    cpu.a = result as u8;
    set_flags(cpu, cpu.a == 0, true, h, result < 0);
}

fn and_a(cpu: &mut Cpu, operand: u8) {
    cpu.a &= operand;
    set_flags(cpu, cpu.a == 0, false, true, false);
}

fn xor_a(cpu: &mut Cpu, operand: u8) {
    cpu.a ^= operand;
    set_flags(cpu, cpu.a == 0, false, false, false);
}

fn or_a(cpu: &mut Cpu, operand: u8) {
    cpu.a |= operand;
    set_flags(cpu, cpu.a == 0, false, false, false);
}

fn cp_a(cpu: &mut Cpu, operand: u8) {
    let (res, c) = cpu.a.overflowing_sub(operand);
    let h = (cpu.a & 0x0F) < (operand & 0x0F);
    set_flags(cpu, res == 0, true, h, c);
}

fn alu_r(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let operand = r8(cpu, mem, op & 7);
    match (op >> 3) & 7 {
        0 => add_a(cpu, operand),
        1 => adc_a(cpu, operand),
        2 => sub_a(cpu, operand),
        3 => sbc_a(cpu, operand),
        4 => and_a(cpu, operand),
        5 => xor_a(cpu, operand),
        6 => or_a(cpu, operand),
        _ => cp_a(cpu, operand),
    }
    if op & 7 == 6 { 8 } else { 4 }
}

fn alu_n(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let operand = fetch8(cpu, mem);
    match op {
        0xC6 => add_a(cpu, operand),
        0xCE => adc_a(cpu, operand),
        0xD6 => sub_a(cpu, operand),
        0xDE => sbc_a(cpu, operand),
        0xE6 => and_a(cpu, operand),
        0xEE => xor_a(cpu, operand),
        0xF6 => or_a(cpu, operand),
        0xFE => cp_a(cpu, operand),
        _ => unreachable!(),
    }
    8
}

fn rlc(cpu: &mut Cpu, value: u8) -> u8 {
    let carry = value >> 7;
    let result = (value << 1) | carry;
    set_flags(cpu, result == 0, false, false, carry == 1);
    result
}

fn rrc(cpu: &mut Cpu, value: u8) -> u8 {
    let carry = value & 1;
    let result = (value >> 1) | (carry << 7);
    set_flags(cpu, result == 0, false, false, carry == 1);
    result
}

fn rl(cpu: &mut Cpu, value: u8) -> u8 {
    let carry = value >> 7;
    let result = (value << 1) | flag(cpu, FLAG_C) as u8;
    set_flags(cpu, result == 0, false, false, carry == 1);
    result
}

fn rr(cpu: &mut Cpu, value: u8) -> u8 {
    let carry = value & 1;
    let result = (value >> 1) | ((flag(cpu, FLAG_C) as u8) << 7);
    set_flags(cpu, result == 0, false, false, carry == 1);
    result
}

fn sla(cpu: &mut Cpu, value: u8) -> u8 {
    let carry = value >> 7;
    let result = value << 1;
    set_flags(cpu, result == 0, false, false, carry == 1);
    result
}

fn sra(cpu: &mut Cpu, value: u8) -> u8 {
    let carry = value & 1;
    let result = (value >> 1) | (value & 0x80);
    set_flags(cpu, result == 0, false, false, carry == 1);
    result
}

fn srl(cpu: &mut Cpu, value: u8) -> u8 {
    let carry = value & 1;
    let result = value >> 1;
    set_flags(cpu, result == 0, false, false, carry == 1);
    result
}

fn swap(cpu: &mut Cpu, value: u8) -> u8 {
    let result = (value >> 4) | (value << 4);
    set_flags(cpu, result == 0, false, false, false);
    result
}

fn rlca(gb: &mut GameBoy, _op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    cpu.a = rlc(cpu, cpu.a);
    cpu.f &= !FLAG_Z;
    4
}

fn rrca(gb: &mut GameBoy, _op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    cpu.a = rrc(cpu, cpu.a);
    cpu.f &= !FLAG_Z;
    4
}

fn rla(gb: &mut GameBoy, _op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    cpu.a = rl(cpu, cpu.a);
    cpu.f &= !FLAG_Z;
    4
}

fn rra(gb: &mut GameBoy, _op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    cpu.a = rr(cpu, cpu.a);
    cpu.f &= !FLAG_Z;
    4
}

fn daa(gb: &mut GameBoy, _op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    let mut adjust = 0u8;
    let n = flag(cpu, FLAG_N);
    let mut carry = flag(cpu, FLAG_C);
    if flag(cpu, FLAG_H) || (!n && cpu.a & 0x0F > 0x09) {
        adjust |= 0x06;
    }
    if carry || (!n && cpu.a > 0x99) {
        adjust |= 0x60;
        carry = true;
    }
    if n {
        cpu.a = cpu.a.wrapping_sub(adjust);
    } else {
        cpu.a = cpu.a.wrapping_add(adjust);
    }
    set_flags(cpu, cpu.a == 0, n, false, carry);
    4
}

fn cpl(gb: &mut GameBoy, _op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    cpu.a = !cpu.a;
    set_flags(cpu, flag(cpu, FLAG_Z), true, true, flag(cpu, FLAG_C));
    4
}

fn scf(gb: &mut GameBoy, _op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    set_flags(cpu, flag(cpu, FLAG_Z), false, false, true);
    4
}

fn ccf(gb: &mut GameBoy, _op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    set_flags(cpu, flag(cpu, FLAG_Z), false, false, !flag(cpu, FLAG_C));
    4
}

fn jr(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let e = fetch8(cpu, mem) as i8;
    let taken = match op {
        0x18 => true,
        0x20 => !flag(cpu, FLAG_Z),
        0x28 => flag(cpu, FLAG_Z),
        0x30 => !flag(cpu, FLAG_C),
        0x38 => flag(cpu, FLAG_C),
        _ => unreachable!(),
    };
    if taken {
        cpu.pc = cpu.pc.wrapping_add(e as u16);
        12
    } else {
        8
    }
}

fn jp_nn(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let addr = fetch16(cpu, mem);
    let taken = match op {
        0xC3 => true,
        0xC2 => !flag(cpu, FLAG_Z),
        0xCA => flag(cpu, FLAG_Z),
        0xD2 => !flag(cpu, FLAG_C),
        0xDA => flag(cpu, FLAG_C),
        _ => unreachable!(),
    };
    if taken {
        cpu.pc = addr;
        16
    } else {
        12
    }
}

fn jp_hl(gb: &mut GameBoy, _op: u8) -> u128 {
    let cpu = &mut gb.cpu;
    cpu.pc = cpu.get_hl();
    4
}

fn call(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let addr = fetch16(cpu, mem);
    let taken = match op {
        0xCD => true,
        0xC4 => !flag(cpu, FLAG_Z),
        0xCC => flag(cpu, FLAG_Z),
        0xD4 => !flag(cpu, FLAG_C),
        0xDC => flag(cpu, FLAG_C),
        _ => unreachable!(),
    };
    if taken {
        push16(cpu, mem, cpu.pc);
        cpu.pc = addr;
        24
    } else {
        12
    }
}

fn ret(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let (check, taken_cycles) = match op {
        0xC9 | 0xD9 => (true, 16),
        0xC0 => (!flag(cpu, FLAG_Z), 20),
        0xC8 => (flag(cpu, FLAG_Z), 20),
        0xD0 => (!flag(cpu, FLAG_C), 20),
        0xD8 => (flag(cpu, FLAG_C), 20),
        _ => unreachable!(),
    };
    if check {
        cpu.pc = pop16(cpu, mem);
        if op == 0xD9 {
            cpu.ime = true;
        }
        taken_cycles
    } else {
        8
    }
}

fn rst(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    push16(cpu, mem, cpu.pc);
    cpu.pc = (op & 0x38) as u16;
    16
}

fn push(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let value = match op {
        0xC5 => cpu.get_bc(),
        0xD5 => cpu.get_de(),
        0xE5 => cpu.get_hl(),
        0xF5 => cpu.get_af(),
        _ => unreachable!(),
    };
    push16(cpu, mem, value);
    16
}

fn pop(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let value = pop16(cpu, mem);
    match op {
        0xC1 => cpu.set_bc(value),
        0xD1 => cpu.set_de(value),
        0xE1 => cpu.set_hl(value),
        0xF1 => cpu.set_af(value),
        _ => unreachable!(),
    }
    12
}

fn prefix_cb(gb: &mut GameBoy, _op: u8) -> u128 {
    let cb = fetch8(&mut gb.cpu, &gb.memory);
    CB_OPCODES[cb as usize](gb, cb)
}

fn cb_shift(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let r = op & 7;
    let value = r8(cpu, mem, r);
    let result = match op >> 3 {
        0 => rlc(cpu, value),
        1 => rrc(cpu, value),
        2 => rl(cpu, value),
        3 => rr(cpu, value),
        4 => sla(cpu, value),
        5 => sra(cpu, value),
        6 => swap(cpu, value),
        _ => srl(cpu, value),
    };
    set_r8(cpu, mem, r, result);
    if r == 6 { 16 } else { 8 }
}

fn cb_bit(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let r = op & 7;
    let bit = (op >> 3) & 7;
    let z = r8(cpu, mem, r) & (1 << bit) == 0;
    set_flags(cpu, z, false, true, flag(cpu, FLAG_C));
    if r == 6 { 12 } else { 8 }
}

fn cb_res(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let r = op & 7;
    let value = r8(cpu, mem, r) & !(1 << ((op >> 3) & 7));
    set_r8(cpu, mem, r, value);
    if r == 6 { 16 } else { 8 }
}

fn cb_set(gb: &mut GameBoy, op: u8) -> u128 {
    let (cpu, mem) = (&mut gb.cpu, &mut gb.memory);
    let r = op & 7;
    let value = r8(cpu, mem, r) | (1 << ((op >> 3) & 7));
    set_r8(cpu, mem, r, value);
    if r == 6 { 16 } else { 8 }
}

fn build_cb_table() -> [OpCodeExecution; 256] {
    let mut table = [illegal as OpCodeExecution; 256];
    for op in 0u8..=0x3F {
        table[op as usize] = cb_shift;
    }
    for op in 0x40u8..=0x7F {
        table[op as usize] = cb_bit;
    }
    for op in 0x80u8..=0xBF {
        table[op as usize] = cb_res;
    }
    for op in 0xC0u8..=0xFF {
        table[op as usize] = cb_set;
    }
    table
}

static CB_OPCODES: std::sync::LazyLock<[OpCodeExecution; 256]> =
    std::sync::LazyLock::new(build_cb_table);

impl OpCodeTable {
    pub fn get_code(&self, code: u8) -> OpCodeExecution {
        self.opcodes[code as usize]
    }

    pub fn new() -> Self {
        let mut opcodes = [illegal as OpCodeExecution; 256];

        opcodes[0x00] = nop;
        opcodes[0x10] = stop;
        opcodes[0x76] = halt;
        opcodes[0xF3] = di;
        opcodes[0xFB] = ei;
        opcodes[0x07] = rlca;
        opcodes[0x0F] = rrca;
        opcodes[0x17] = rla;
        opcodes[0x1F] = rra;
        opcodes[0x27] = daa;
        opcodes[0x2F] = cpl;
        opcodes[0x37] = scf;
        opcodes[0x3F] = ccf;
        opcodes[0xCB] = prefix_cb;

        opcodes[0x04] = inc_r; // INC B
        opcodes[0x0C] = inc_r; // INC C
        opcodes[0x14] = inc_r; // INC D
        opcodes[0x1C] = inc_r; // INC E
        opcodes[0x24] = inc_r; // INC H
        opcodes[0x2C] = inc_r; // INC L
        opcodes[0x34] = inc_r; // INC (HL)
        opcodes[0x3C] = inc_r; // INC A

        opcodes[0x05] = dec_r; // DEC B
        opcodes[0x0D] = dec_r; // DEC C
        opcodes[0x15] = dec_r; // DEC D
        opcodes[0x1D] = dec_r; // DEC E
        opcodes[0x25] = dec_r; // DEC H
        opcodes[0x2D] = dec_r; // DEC L
        opcodes[0x35] = dec_r; // DEC (HL)
        opcodes[0x3D] = dec_r; // DEC A

        opcodes[0x06] = ld_r_n; // LD B, n
        opcodes[0x0E] = ld_r_n; // LD C, n
        opcodes[0x16] = ld_r_n; // LD D, n
        opcodes[0x1E] = ld_r_n; // LD E, n
        opcodes[0x26] = ld_r_n; // LD H, n
        opcodes[0x2E] = ld_r_n; // LD L, n
        opcodes[0x36] = ld_r_n; // LD (HL), n
        opcodes[0x3E] = ld_r_n; // LD A, n

        opcodes[0x01] = ld_rr_nn; // LD BC, nn
        opcodes[0x11] = ld_rr_nn; // LD DE, nn
        opcodes[0x21] = ld_rr_nn; // LD HL, nn
        opcodes[0x31] = ld_rr_nn; // LD SP, nn

        opcodes[0x03] = inc_rr; // INC BC
        opcodes[0x13] = inc_rr; // INC DE
        opcodes[0x23] = inc_rr; // INC HL
        opcodes[0x33] = inc_rr; // INC SP

        opcodes[0x0B] = dec_rr; // DEC BC
        opcodes[0x1B] = dec_rr; // DEC DE
        opcodes[0x2B] = dec_rr; // DEC HL
        opcodes[0x3B] = dec_rr; // DEC SP

        opcodes[0x09] = add_hl_rr; // ADD HL, BC
        opcodes[0x19] = add_hl_rr; // ADD HL, DE
        opcodes[0x29] = add_hl_rr; // ADD HL, HL
        opcodes[0x39] = add_hl_rr; // ADD HL, SP

        opcodes[0x40] = ld_r_r; // LD B, B
        opcodes[0x41] = ld_r_r; // LD B, C
        opcodes[0x42] = ld_r_r; // LD B, D
        opcodes[0x43] = ld_r_r; // LD B, E
        opcodes[0x44] = ld_r_r; // LD B, H
        opcodes[0x45] = ld_r_r; // LD B, L
        opcodes[0x46] = ld_r_r; // LD B, (HL)
        opcodes[0x47] = ld_r_r; // LD B, A
        opcodes[0x48] = ld_r_r; // LD C, B
        opcodes[0x49] = ld_r_r; // LD C, C
        opcodes[0x4A] = ld_r_r; // LD C, D
        opcodes[0x4B] = ld_r_r; // LD C, E
        opcodes[0x4C] = ld_r_r; // LD C, H
        opcodes[0x4D] = ld_r_r; // LD C, L
        opcodes[0x4E] = ld_r_r; // LD C, (HL)
        opcodes[0x4F] = ld_r_r; // LD C, A
        opcodes[0x50] = ld_r_r; // LD D, B
        opcodes[0x51] = ld_r_r; // LD D, C
        opcodes[0x52] = ld_r_r; // LD D, D
        opcodes[0x53] = ld_r_r; // LD D, E
        opcodes[0x54] = ld_r_r; // LD D, H
        opcodes[0x55] = ld_r_r; // LD D, L
        opcodes[0x56] = ld_r_r; // LD D, (HL)
        opcodes[0x57] = ld_r_r; // LD D, A
        opcodes[0x58] = ld_r_r; // LD E, B
        opcodes[0x59] = ld_r_r; // LD E, C
        opcodes[0x5A] = ld_r_r; // LD E, D
        opcodes[0x5B] = ld_r_r; // LD E, E
        opcodes[0x5C] = ld_r_r; // LD E, H
        opcodes[0x5D] = ld_r_r; // LD E, L
        opcodes[0x5E] = ld_r_r; // LD E, (HL)
        opcodes[0x5F] = ld_r_r; // LD E, A
        opcodes[0x60] = ld_r_r; // LD H, B
        opcodes[0x61] = ld_r_r; // LD H, C
        opcodes[0x62] = ld_r_r; // LD H, D
        opcodes[0x63] = ld_r_r; // LD H, E
        opcodes[0x64] = ld_r_r; // LD H, H
        opcodes[0x65] = ld_r_r; // LD H, L
        opcodes[0x66] = ld_r_r; // LD H, (HL)
        opcodes[0x67] = ld_r_r; // LD H, A
        opcodes[0x68] = ld_r_r; // LD L, B
        opcodes[0x69] = ld_r_r; // LD L, C
        opcodes[0x6A] = ld_r_r; // LD L, D
        opcodes[0x6B] = ld_r_r; // LD L, E
        opcodes[0x6C] = ld_r_r; // LD L, H
        opcodes[0x6D] = ld_r_r; // LD L, L
        opcodes[0x6E] = ld_r_r; // LD L, (HL)
        opcodes[0x6F] = ld_r_r; // LD L, A
        opcodes[0x70] = ld_r_r; // LD (HL), B
        opcodes[0x71] = ld_r_r; // LD (HL), C
        opcodes[0x72] = ld_r_r; // LD (HL), D
        opcodes[0x73] = ld_r_r; // LD (HL), E
        opcodes[0x74] = ld_r_r; // LD (HL), H
        opcodes[0x75] = ld_r_r; // LD (HL), L
        opcodes[0x77] = ld_r_r; // LD (HL), A
        opcodes[0x78] = ld_r_r; // LD A, B
        opcodes[0x79] = ld_r_r; // LD A, C
        opcodes[0x7A] = ld_r_r; // LD A, D
        opcodes[0x7B] = ld_r_r; // LD A, E
        opcodes[0x7C] = ld_r_r; // LD A, H
        opcodes[0x7D] = ld_r_r; // LD A, L
        opcodes[0x7E] = ld_r_r; // LD A, (HL)
        opcodes[0x7F] = ld_r_r; // LD A, A

        opcodes[0x80] = alu_r; // ADD A, B
        opcodes[0x81] = alu_r; // ADD A, C
        opcodes[0x82] = alu_r; // ADD A, D
        opcodes[0x83] = alu_r; // ADD A, E
        opcodes[0x84] = alu_r; // ADD A, H
        opcodes[0x85] = alu_r; // ADD A, L
        opcodes[0x86] = alu_r; // ADD A, (HL)
        opcodes[0x87] = alu_r; // ADD A, A
        opcodes[0x88] = alu_r; // ADC A, B
        opcodes[0x89] = alu_r; // ADC A, C
        opcodes[0x8A] = alu_r; // ADC A, D
        opcodes[0x8B] = alu_r; // ADC A, E
        opcodes[0x8C] = alu_r; // ADC A, H
        opcodes[0x8D] = alu_r; // ADC A, L
        opcodes[0x8E] = alu_r; // ADC A, (HL)
        opcodes[0x8F] = alu_r; // ADC A, A
        opcodes[0x90] = alu_r; // SUB B
        opcodes[0x91] = alu_r; // SUB C
        opcodes[0x92] = alu_r; // SUB D
        opcodes[0x93] = alu_r; // SUB E
        opcodes[0x94] = alu_r; // SUB H
        opcodes[0x95] = alu_r; // SUB L
        opcodes[0x96] = alu_r; // SUB (HL)
        opcodes[0x97] = alu_r; // SUB A
        opcodes[0x98] = alu_r; // SBC A, B
        opcodes[0x99] = alu_r; // SBC A, C
        opcodes[0x9A] = alu_r; // SBC A, D
        opcodes[0x9B] = alu_r; // SBC A, E
        opcodes[0x9C] = alu_r; // SBC A, H
        opcodes[0x9D] = alu_r; // SBC A, L
        opcodes[0x9E] = alu_r; // SBC A, (HL)
        opcodes[0x9F] = alu_r; // SBC A, A
        opcodes[0xA0] = alu_r; // AND B
        opcodes[0xA1] = alu_r; // AND C
        opcodes[0xA2] = alu_r; // AND D
        opcodes[0xA3] = alu_r; // AND E
        opcodes[0xA4] = alu_r; // AND H
        opcodes[0xA5] = alu_r; // AND L
        opcodes[0xA6] = alu_r; // AND (HL)
        opcodes[0xA7] = alu_r; // AND A
        opcodes[0xA8] = alu_r; // XOR B
        opcodes[0xA9] = alu_r; // XOR C
        opcodes[0xAA] = alu_r; // XOR D
        opcodes[0xAB] = alu_r; // XOR E
        opcodes[0xAC] = alu_r; // XOR H
        opcodes[0xAD] = alu_r; // XOR L
        opcodes[0xAE] = alu_r; // XOR (HL)
        opcodes[0xAF] = alu_r; // XOR A
        opcodes[0xB0] = alu_r; // OR B
        opcodes[0xB1] = alu_r; // OR C
        opcodes[0xB2] = alu_r; // OR D
        opcodes[0xB3] = alu_r; // OR E
        opcodes[0xB4] = alu_r; // OR H
        opcodes[0xB5] = alu_r; // OR L
        opcodes[0xB6] = alu_r; // OR (HL)
        opcodes[0xB7] = alu_r; // OR A
        opcodes[0xB8] = alu_r; // CP B
        opcodes[0xB9] = alu_r; // CP C
        opcodes[0xBA] = alu_r; // CP D
        opcodes[0xBB] = alu_r; // CP E
        opcodes[0xBC] = alu_r; // CP H
        opcodes[0xBD] = alu_r; // CP L
        opcodes[0xBE] = alu_r; // CP (HL)
        opcodes[0xBF] = alu_r; // CP A

        opcodes[0xC6] = alu_n; // ADD A, n
        opcodes[0xCE] = alu_n; // ADC A, n
        opcodes[0xD6] = alu_n; // SUB n
        opcodes[0xDE] = alu_n; // SBC A, n
        opcodes[0xE6] = alu_n; // AND n
        opcodes[0xEE] = alu_n; // XOR n
        opcodes[0xF6] = alu_n; // OR n
        opcodes[0xFE] = alu_n; // CP n

        opcodes[0x02] = ld_mem_a;
        opcodes[0x12] = ld_mem_a;
        opcodes[0x22] = ld_mem_a;
        opcodes[0x32] = ld_mem_a;
        opcodes[0xEA] = ld_mem_a;
        opcodes[0x0A] = ld_a_mem;
        opcodes[0x1A] = ld_a_mem;
        opcodes[0x2A] = ld_a_mem;
        opcodes[0x3A] = ld_a_mem;
        opcodes[0xFA] = ld_a_mem;
        opcodes[0xE0] = ldh_n_a;
        opcodes[0xF0] = ldh_a_n;
        opcodes[0xE2] = ldh_c_a;
        opcodes[0xF2] = ldh_a_c;
        opcodes[0x08] = ld_nn_sp;
        opcodes[0xF9] = ld_sp_hl;
        opcodes[0xF8] = ld_hl_sp_e;
        opcodes[0xE8] = add_sp_e;

        opcodes[0x18] = jr;
        opcodes[0x20] = jr;
        opcodes[0x28] = jr;
        opcodes[0x30] = jr;
        opcodes[0x38] = jr;
        opcodes[0xC3] = jp_nn;
        opcodes[0xC2] = jp_nn;
        opcodes[0xCA] = jp_nn;
        opcodes[0xD2] = jp_nn;
        opcodes[0xDA] = jp_nn;
        opcodes[0xE9] = jp_hl;
        opcodes[0xCD] = call;
        opcodes[0xC4] = call;
        opcodes[0xCC] = call;
        opcodes[0xD4] = call;
        opcodes[0xDC] = call;
        opcodes[0xC9] = ret;
        opcodes[0xD9] = ret;
        opcodes[0xC0] = ret;
        opcodes[0xC8] = ret;
        opcodes[0xD0] = ret;
        opcodes[0xD8] = ret;

        for rst_op in [0xC7, 0xCF, 0xD7, 0xDF, 0xE7, 0xEF, 0xF7, 0xFF] {
            opcodes[rst_op] = rst;
        }
        for op in [0xC5, 0xD5, 0xE5, 0xF5] {
            opcodes[op] = push;
        }
        for op in [0xC1, 0xD1, 0xE1, 0xF1] {
            opcodes[op] = pop;
        }

        Self { opcodes }
    }
}
