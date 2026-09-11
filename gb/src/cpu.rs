use crate::memory;

pub mod registers;

const N_CYCLES: [u16; 4] = [1024, 16, 64, 256];

#[derive(Default)]
pub struct Cpu {
    a: u8,
    b: u8,
    c: u8,
    d: u8,
    e: u8,
    f: u8,
    h: u8,
    l: u8,
    cycle: u128,
    memory: memory::memory,
    halt: bool,
    stopeed: bool,
    pending_enable_ime: bool,
    ime: bool,
    serial_output: String,
}

impl Cpu {
    pub fn new() -> Self {
        Self::default()
    }
    pub fn get_af(&self) -> u16 {
        (self.a as u16) << 8 | (self.f as u16)
    }
    pub fn get_bc(&self) -> u16 {
        (self.b as u16) << 8 | (self.c as u16)
    }
    pub fn get_de(&self) -> u16 {
        (self.d as u16) << 8 | (self.e as u16)
    }
    pub fn get_hl(&self) -> u16 {
        (self.h as u16) << 8 | (self.l as u16)
    }
    pub fn set_af(&mut self, v: u16) {
        self.a = (v >> 8) as u8;
        self.f = (v & 0xF0) as u8;
    }
    pub fn set_bc(&mut self, v: u16) {
        self.b = (v >> 8) as u8;
        self.c = (v & 0xFF) as u8;
    }
    pub fn set_de(&mut self, v: u16) {
        self.d = (v >> 8) as u8;
        self.e = (v & 0xFF) as u8;
    }
    pub fn set_hl(&mut self, v: u16) {
        self.h = (v >> 8) as u8;
        self.l = (v & 0xFF) as u8;
    }

    pub fn display() {
        println!("Je suis un micro");
    }
}
