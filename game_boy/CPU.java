package game_boy;

import java.io.DataInputStream;
import java.io.FileInputStream;
import java.util.Arrays;

public class CPU {
    private int a, b, c, d, e, f, h, l;
    private int af, bc, de, hl;
    private int cycle;
    private int program_counter;

    public CPU(int a, int b, int c, int d, int e, int f, int h, int l, int cycle, int program_counter) {
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
        this.cycle = cycle;
        this.program_counter = program_counter;
    }

    public int get_bc() {
        String b_string = Integer.toString(b);
        String c_string = Integer.toString(c);

        String bc_string = b_string + c_string;
        return Integer.parseInt(bc_string);
    }

    public void OpCodes(int code) {
    };

    public void execute_command(int opCodes) {

    };

    public static void main(String[] args) {
        DataInputStream reader = null;
        final CPU cpu = new CPU(0, 0, 0, 0, 0, 0, 0, 0, 0, 0);
        try {
            final String rom_path = args[0];

            reader = new DataInputStream(new FileInputStream((rom_path)));
            int nBytesToRead = reader.available();

            if (nBytesToRead > 0) {
                byte[] bytes = new byte[nBytesToRead];
                reader.read(bytes);
                // result = new String(bytes);
                // System.out.println(bytes);
                int[] hexas = new int[nBytesToRead];

                // for (byte b : bytes) {
                // // System.out.printf("%02X ", b & 0xFF);
                // hexas[]
                // }

                for (int i = 0; i < nBytesToRead; i++) {
                    hexas[i] = bytes[i] & 0xFF;
                }
                System.out.println(Arrays.toString(hexas));
            }

            reader.close();
        } catch (Exception e) {
            System.out.println("Error: " + e.getMessage());

        }

    }
}