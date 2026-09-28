pub struct Memory {
    memory: [u8; 65536],
    pub serial_output: String,
}

impl Memory {
    pub fn new() -> Self {
        Memory {
            memory: [0; 65536],
            serial_output: String::new(),
        }
    }

    pub fn write(&mut self, address: u16, value: u8) {
        match address {
            0x0000..=0x7FFF => panic!(
                "Multiple Bank Cartridge has not been implemented yet: 0x{:04X}",
                address
            ),
            0x8000..=0x9FFF | 0xA000..=0xBFFF | 0xC000..=0xDFFF => {
                self.memory[address as usize] = value;
            }
            0xE000..=0xFDFF => {
                self.memory[address as usize] = value;
                self.memory[(address - 0x2000) as usize] = value;
            }
            0xFE00..=0xFE9F => {
                self.memory[address as usize] = value;
            }
            0xFEA0..=0xFEFF | 0xFF4C..=0xFF7F => {
                panic!("Write to unusable memory: 0x{:04X}", address)
            }
            0xFF00..=0xFF4B => {
                self.memory[address as usize] = value;
                if address == 0xFF02 && value == 0x81 {
                    let c = self.memory[0xFF01] as char;
                    print!("{}", c);
                    self.serial_output.push(c);
                }
            }
            0xFF80..=0xFFFF => {
                self.memory[address as usize] = value;
            }
        }
    }
    pub fn read(&self, address: u16) -> u8 {
        self.memory[address as usize]
    }

    pub fn load_rom(&mut self, rom: &[u8]) {
        let n = rom.len().min(0x8000);
        self.memory[..n].copy_from_slice(&rom[..n]);
    }
}
