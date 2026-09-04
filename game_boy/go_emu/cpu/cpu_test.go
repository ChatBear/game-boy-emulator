package cpu

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestCpuInstrs runs every individual Blargg cpu_instrs ROM as a subtest.
// Each subtest fails independently if its ROM does not report "Passed".
func TestCpuInstrs(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file path")
	}
	romsDir := filepath.Join(filepath.Dir(thisFile), "..", "gb-test-roms", "cpu_instrs", "individual")

	entries, err := os.ReadDir(romsDir)
	if err != nil {
		t.Fatalf("cannot read %s: %v", romsDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".gb") {
			continue
		}
		romPath := filepath.Join(romsDir, entry.Name())
		name := strings.TrimSuffix(entry.Name(), ".gb")

		t.Run(name, func(t *testing.T) {
			runBlarggROM(t, romPath)
		})
	}
}

// runBlarggROM executes a Blargg test ROM and fails the test
// if its serial output does not contain "Passed".
func runBlarggROM(t *testing.T, romPath string) {
	t.Helper()

	cpu, err := NewCPU(0, 0, 0, 0, 0, 0, 0, 0)
	if err != nil {
		t.Fatalf("NewCPU error: %v", err)
	}

	bytes, err := os.ReadFile(romPath)
	if err != nil {
		t.Fatalf("cannot read ROM %s: %v", romPath, err)
	}
	cpu.UploadROM(bytes)

	const maxSteps = 35_000_000
	lastPC := uint16(0xFFFE)
	sameCount := 0

	for i := 0; i < maxSteps; i++ {
		pc := cpu.programCounter
		if err := cpu.Step(); err != nil {
			t.Fatalf("step %d (PC=0x%04X): %v\nserial=%q",
				i, pc, err, cpu.serialOutput.String())
		}

		if cpu.halt {
			sameCount = 0
		} else if cpu.programCounter == lastPC {
			sameCount++
			if sameCount > 10 {
				break
			}
		} else {
			sameCount = 0
		}
		lastPC = cpu.programCounter

		out := cpu.serialOutput.String()
		if strings.Contains(out, "Passed") || strings.Contains(out, "Failed") {
			break
		}
	}

	out := cpu.serialOutput.String()
	if !strings.Contains(out, "Passed") {
		t.Logf("cycle=%d TIMA=%02X TAC=%02X TMA=%02X IF=%02X IE=%02X PC=%04X",
			cpu.cycle,
			cpu.memory[0xFF05],
			cpu.memory[0xFF07],
			cpu.memory[0xFF06],
			cpu.memory[0xFF0F],
			cpu.memory[0xFFFF],
			cpu.programCounter)
		t.Fatalf("ROM did not pass\nserial output:\n%s", out)
	}
}
