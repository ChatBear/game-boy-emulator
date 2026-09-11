pub struct Memory {
    memory: [u8; 65536],
}

impl Memory {
    pub fn new() -> Self {
        Memory { memory: [0; 65536] }
    }

    fn write(adress: u16, value: u8) {
        if adress < 0x8000 {}
    }
    fn read(adress: u8) {}
}
