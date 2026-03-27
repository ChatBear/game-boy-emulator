package cpu

import (
	"fmt"
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	cpu, err := NewCPU(0, 0, 0, 0, 0, 0, 0, 0)

	if err != nil {
		t.Fail()
	}
	romPath := "/Users/shiraz/Desktop/Bureau - MacBook Pro de Shiraz (3)/myProject/tuto_java/game_boy/go_emu/blargg-test-roms/cpu_instrs/individual/06-ld r,r.gb"
	bytes, err := os.ReadFile(romPath)
	if err != nil {
		fmt.Printf("Error reading ROM: %v\n", err)
		return
	}

	cpu.UploadROM(bytes)
	cpu.Run(30000)
}
