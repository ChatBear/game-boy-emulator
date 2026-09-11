pub mod cpu;
pub mod gb;
pub mod memory;

fn main() {
    let gb = gb::GameBoy::new();
    gb.display();
}
