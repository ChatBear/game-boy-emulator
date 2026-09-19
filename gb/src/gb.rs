use std::any::Any;
use std::collections::HashMap;

use crate::{cpu::Cpu, memory::Memory};

pub struct GameBoy {
    cpu: Cpu,
    memory: Memory,
}

impl GameBoy {
    pub fn new() -> Self {
        GameBoy {
            cpu: Cpu::default(),
            memory: Memory::new(),
        }
    }

    pub fn display(self) {
        self.cpu.display();
        println!("Je suis une game boy");
    }

    fn decode(code: u8) {
        
    }
}
