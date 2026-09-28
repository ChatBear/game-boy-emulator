use std::path::Path;

use crate::{
    cpu::Cpu,
    memory::Memory,
    opcodes::{self, OpCodeTable},
};

// const REG8: [&str; 8] = ["B", "C", "D", "E", "H", "L", "(HL)", "A"];
const N_CYCLES: [u16; 4] = [1024, 16, 64, 256];

const TIMA: u16 = 0xFF05;
const TMA: u16 = 0xFF06;
const TAC: u16 = 0xFF07;
const IF: u16 = 0xFF0F;
const IE: u16 = 0xFFFF;

pub struct GameBoy {
    pub cpu: Cpu,
    pub memory: Memory,
    op_code_table: OpCodeTable,
}

impl GameBoy {
    pub fn new() -> Self {
        GameBoy {
            cpu: Cpu::new(),
            memory: Memory::new(),
            op_code_table: OpCodeTable::new(),
        }
    }

    pub fn display(self) {
        self.cpu.display();
    }

    pub fn handler(&mut self, code: &u8) {
        let opcode = self.op_code_table.get_code(*code);
        self.cpu.cycle += opcode(self, *code);
    }

    pub fn execute_game(&self, _: &Path) {
        // println!("Opening the rom in this path {rom_path}")
    }

    pub fn update_timer(&mut self, cycles: u128) {
        let tac = self.memory.read(TAC);
        if tac & 0x04 == 0 {
            return;
        }

        let period = N_CYCLES[(tac & 0x03) as usize] as u64;
        self.cpu.timer_acc += cycles as u64;

        while self.cpu.timer_acc >= period {
            self.cpu.timer_acc -= period;

            let tima = self.memory.read(TIMA);
            if tima == 0xFF {
                let tma = self.memory.read(TMA);
                self.memory.write(TIMA, tma);
                let if_ = self.memory.read(IF);
                self.memory.write(IF, if_ | 0x04);
            } else {
                self.memory.write(TIMA, tima + 1);
            }
        }
    }

    /// Returns the cycles spent if an interrupt was serviced.
    pub fn service_interrupt(&mut self) -> Option<u128> {
        let if_ = self.memory.read(IF);
        let pending = self.memory.read(IE) & if_ & 0x1F;
        if pending == 0 {
            return None;
        }

        // HALT wakes up on any pending interrupt, even with IME off.
        self.cpu.halt = false;
        if !self.cpu.ime {
            return None;
        }

        let bit = pending.trailing_zeros();
        self.cpu.ime = false;
        self.memory.write(IF, if_ & !(1u8 << bit));
        let pc = self.cpu.pc;
        opcodes::push16(&mut self.cpu, &mut self.memory, pc);
        self.cpu.pc = 0x0040 + 8 * bit as u16;
        Some(20)
    }

    fn execute_next(&mut self) -> u128 {
        let enable_ime = self.cpu.ei_pending;

        let code = self.memory.read(self.cpu.pc);
        self.cpu.pc = self.cpu.pc.wrapping_add(1);
        let opcode = self.op_code_table.get_code(code);
        let cycles = opcode(self, code);

        // EI takes effect after the instruction that follows it (unless a DI cancelled it).
        if enable_ime && self.cpu.ei_pending {
            self.cpu.ime = true;
            self.cpu.ei_pending = false;
        }
        cycles
    }

    pub fn step(&mut self) {
        let cycles = match self.service_interrupt() {
            Some(cycles) => cycles,
            None if self.cpu.halt => 4,
            None => self.execute_next(),
        };
        self.cpu.cycle += cycles;
        self.update_timer(cycles);
    }

    pub fn run(&mut self) {
        loop {
            self.step();
        }
    }
}

#[cfg(test)]
mod test {
    use std::{fs, panic, path::PathBuf};

    use super::*;

    const MAX_STEP: u64 = 35_000_000;

    #[test]
    fn run_blargg_roms() {
        let roms_dir = PathBuf::from(env!("CARGO_MANIFEST_DIR"))
            .join("test_roms/gb-test-roms/cpu_instrs/individual");

        let mut roms: Vec<PathBuf> = fs::read_dir(&roms_dir)
            .expect("cannot read ROM directory")
            .filter_map(|entry| entry.ok().map(|e| e.path()))
            .filter(|p| p.extension().is_some_and(|ext| ext == "gb"))
            .collect();
        roms.sort();

        let mut failures = Vec::new();

        for path in &roms {
            let name: std::borrow::Cow<'_, str> = path.file_stem().unwrap().to_string_lossy();
            let rom_bytes = fs::read(path).expect("cannot read ROM");

            match panic::catch_unwind(|| run_rom(&rom_bytes)) {
                Ok(Ok(())) => println!("[PASS] {name}"),
                Ok(Err(report)) => failures.push(format!("[FAIL] {name}\n{report}")),
                Err(_) => failures.push(format!("[PANIC] {name}")),
            }
        }

        assert!(failures.is_empty(), "\n{}", failures.join("\n\n"));
    }

    fn run_rom(rom_bytes: &[u8]) -> Result<(), String> {
        let mut gb: GameBoy = GameBoy::new();
        gb.memory.load_rom(rom_bytes);

        let mut last_pc: u16 = 0xFFFE;
        let mut same_count = 0;

        for _ in 0..MAX_STEP {
            gb.step();

            if gb.cpu.halt {
                same_count = 0;
            } else if gb.cpu.pc == last_pc {
                same_count += 1;
                if same_count > 10 {
                    break;
                }
            } else {
                same_count = 0;
            }
            last_pc = gb.cpu.pc;

            let out = &gb.memory.serial_output;
            if out.contains("Passed") || out.contains("Failed") {
                break;
            }
        }

        let out = &gb.memory.serial_output;
        if out.contains("Passed") {
            return Ok(());
        }
        Err(format!(
            "cycle={} TIMA={:02X} TAC={:02X} TMA={:02X} IF={:02X} IE={:02X} PC={:04X}\nserial output:\n{}",
            gb.cpu.cycle,
            gb.memory.read(0xFF05),
            gb.memory.read(0xFF07),
            gb.memory.read(0xFF06),
            gb.memory.read(0xFF0F),
            gb.memory.read(0xFFFF),
            gb.cpu.pc,
            out
        ))
    }
}
