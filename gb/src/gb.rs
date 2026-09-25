use std::path::Path;

use crate::{
    cpu::Cpu,
    memory::Memory,
    opcodes::{Decoded, OpCodeTable, decode},
};

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
        if let Some(decoded) = decode(*code){
            
        }
        self.cpu.cycle += opcode(self, *code) as u128;
    }

    pub fn execute_game(&self, rom_path: &Path) {
        // println!("Opening the rom in this path {rom_path}")
    }

    pub fn run(&self) {
        let mut current_adress: u16 = self.cpu.pc;
        let mut code: u8 = self.memory.read(current_adress);
        while !self.cpu.halt {
            self.handler(&code);
        }
    }
}
