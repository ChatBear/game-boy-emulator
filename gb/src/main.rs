pub mod cpu;
pub mod memory;

fn main() {
    let cpu = cpu::Cpu::new();
    cpu.display();
}
