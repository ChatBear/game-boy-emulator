pub mod cpu;
pub mod gb;
pub mod memory;
pub mod opcodes;

fn main() {
    let gb = gb::GameBoy::new();
    gb.display();
}
