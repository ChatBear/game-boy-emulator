package main.java.com;

import java.io.DataInputStream;
import java.io.FileInputStream;

public class CPU {
    private int a, b, c, d, e, f, h, l;
    private int af, bc, de, hl;
    private int cycle;
    private int program_counter;
    private int stack_pointer; 
    private int scx, scy; 

    private int[] memory;

    public CPU(int a, int b, int c, int d, int e, int f, int h, int l) {
        // // Check if the attributes are all 8-Bytes non signed
        if ((a < 0 || a > 255) || (b < 0 || b > 255) || (c < 0 || c > 255) || (d < 0 || d > 255) || (e < 0 || e > 255)
                || (f < 0 || f > 255) || (h < 0 || h > 255) || (l < 0 || l > 255) || (a < 0 || a > 255)) {
            throw new IllegalArgumentException("Register values must be between 0 and 255 so 8-bit unsigned");
        }
        this.a = a;
        this.b = b;
        this.c = c;
        this.d = d;
        this.e = e;
        this.f = f;
        this.h = h;
        this.l = l;
        this.memory = new int[0xFFFF];
    }

    // TODO : Need to add a banking transition system on the memory not done yet
    // Look for MBC1 and MBC2 in the page 13
    private void upload_rom(int[] rom) {
        System.out.println("Writing the first 32Kb on the Memory");
        for (int i = 0; i < 0x8000; i++) {
            this.memory[i] = rom[i];
        }
        System.out.println("Done");
    }

    public void OpCodes(int code) {
    };

    // System of bank switching : Two types of Cartridge : MBC1 and MBC2 (3, 4, 5)
    // depending on the size of the game
    // It is also named in the header of the card -> in the rom binary
    //
    public void execute_command(int opCodes) {

    };
    public void initialize() {
        System.out.print("----------------------------------------------------------------- \n");
        this.stack_pointer = 0xFFFE;
        this.program_counter = 0; 
        this.cycle = 0; 
        System.out.print("  \n End of initialization \n");
    } 

    public void boot() {
        for (int i=0x0000; i<=0x00FF; i++) {
            int instruction = this.memory[i]; 
        }
    }
    public static void main(String[] args) {
        DataInputStream reader = null;
        final CPU cpu = new CPU(0, 0, 0, 0, 0, 0, 0, 0);
        try {
            final String rom_path = "rom.gb";

            reader = new DataInputStream(new FileInputStream((rom_path)));
            int nBytesToRead = reader.available();

            if (nBytesToRead > 0) {
                byte[] bytes = new byte[nBytesToRead];
                reader.read(bytes);
                int[] hexas = new int[nBytesToRead];

                for (int i = 0; i < nBytesToRead; i++) {
                    hexas[i] = bytes[i] & 0xFF;
                }
                cpu.upload_rom(hexas);
            }
            reader.close();
            cpu.initialize();
        } catch (Exception e) {
            System.out.println("Error: " + e.getMessage());

        }

    }
}