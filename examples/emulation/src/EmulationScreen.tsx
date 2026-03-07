import React from 'react';
import { Cpu, Monitor, Volume2, Gamepad2, HardDrive, Clock, Box, ArrowRight } from 'lucide-react';

export default function EmulatorArchitecture() {
  const components = [
    {
      icon: HardDrive,
      name: "ROM Loader",
      desc: "Read game/program files",
      color: "bg-blue-500",
      order: 1
    },
    {
      icon: Cpu,
      name: "CPU Core",
      desc: "Fetch-Decode-Execute cycle",
      color: "bg-purple-500",
      order: 2
    },
    {
      icon: Box,
      name: "Memory (RAM)",
      desc: "Address space mapping",
      color: "bg-green-500",
      order: 3
    },
    {
      icon: Monitor,
      name: "Graphics (PPU/GPU)",
      desc: "Video output simulation",
      color: "bg-red-500",
      order: 4
    },
    {
      icon: Volume2,
      name: "Audio (APU)",
      desc: "Sound generation",
      color: "bg-yellow-500",
      order: 5
    },
    {
      icon: Gamepad2,
      name: "Input Handler",
      desc: "Controller/keyboard mapping",
      color: "bg-indigo-500",
      order: 6
    },
    {
      icon: Clock,
      name: "Timing/Sync",
      desc: "Clock cycle management",
      color: "bg-pink-500",
      order: 7
    }
  ];

  return (
    <div className="w-full h-full bg-gradient-to-br from-slate-900 to-slate-800 p-8 overflow-auto">
      <div className="max-w-6xl mx-auto">
        <h1 className="text-3xl font-bold text-white mb-2">Emulator Implementation Guide</h1>
        <p className="text-slate-300 mb-8">Key components to build when creating an emulator</p>

        {/* Main Architecture Flow */}
        <div className="bg-slate-800 rounded-lg p-6 mb-8 border border-slate-700">
          <h2 className="text-xl font-semibold text-white mb-4">Implementation Order</h2>
          <div className="grid grid-cols-1 gap-4">
            {components.map((comp, idx) => {
              const Icon = comp.icon;
              return (
                <div key={idx} className="flex items-center gap-4">
                  <div className="flex items-center justify-center w-12 h-12 rounded-full bg-slate-700 text-white font-bold">
                    {comp.order}
                  </div>
                  <div className={`flex-1 ${comp.color} bg-opacity-20 border-l-4 ${comp.color.replace('bg-', 'border-')} rounded p-4`}>
                    <div className="flex items-center gap-3">
                      <Icon className="text-white" size={24} />
                      <div>
                        <h3 className="text-white font-semibold">{comp.name}</h3>
                        <p className="text-slate-300 text-sm">{comp.desc}</p>
                      </div>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Detailed Architecture Diagram */}
        <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h2 className="text-xl font-semibold text-white mb-6">System Architecture</h2>

          <div className="space-y-6">
            {/* Top Layer - Input */}
            <div className="flex justify-center">
              <div className="bg-blue-600 text-white px-6 py-3 rounded-lg font-semibold">
                ROM/Game File
              </div>
            </div>

            <div className="flex justify-center">
              <ArrowRight className="text-slate-500 rotate-90" size={32} />
            </div>

            {/* Core Layer */}
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-purple-600 bg-opacity-30 border-2 border-purple-500 rounded-lg p-4">
                <h3 className="text-white font-bold mb-2">CPU EMULATION</h3>
                <ul className="text-sm text-slate-300 space-y-1">
                  <li>• Instruction decoder</li>
                  <li>• Registers (A, X, Y, SP, PC)</li>
                  <li>• ALU operations</li>
                  <li>• Flags (N, Z, C, V)</li>
                </ul>
              </div>

              <div className="bg-green-600 bg-opacity-30 border-2 border-green-500 rounded-lg p-4">
                <h3 className="text-white font-bold mb-2">MEMORY</h3>
                <ul className="text-sm text-slate-300 space-y-1">
                  <li>• RAM array</li>
                  <li>• Memory mapping</li>
                  <li>• Read/Write methods</li>
                  <li>• Bank switching</li>
                </ul>
              </div>
            </div>

            <div className="flex justify-center">
              <ArrowRight className="text-slate-500 rotate-90" size={32} />
            </div>

            {/* Peripheral Layer */}
            <div className="grid grid-cols-3 gap-4">
              <div className="bg-red-600 bg-opacity-30 border-2 border-red-500 rounded-lg p-4">
                <h3 className="text-white font-bold mb-2">GRAPHICS</h3>
                <ul className="text-sm text-slate-300 space-y-1">
                  <li>• Framebuffer</li>
                  <li>• Sprite rendering</li>
                  <li>• Scanline timing</li>
                </ul>
              </div>

              <div className="bg-yellow-600 bg-opacity-30 border-2 border-yellow-500 rounded-lg p-4">
                <h3 className="text-white font-bold mb-2">AUDIO</h3>
                <ul className="text-sm text-slate-300 space-y-1">
                  <li>• Waveform gen</li>
                  <li>• Audio buffer</li>
                  <li>• Channel mixing</li>
                </ul>
              </div>

              <div className="bg-indigo-600 bg-opacity-30 border-2 border-indigo-500 rounded-lg p-4">
                <h3 className="text-white font-bold mb-2">INPUT</h3>
                <ul className="text-sm text-slate-300 space-y-1">
                  <li>• Button states</li>
                  <li>• Key mapping</li>
                  <li>• Controller poll</li>
                </ul>
              </div>
            </div>

            <div className="flex justify-center">
              <ArrowRight className="text-slate-500 rotate-90" size={32} />
            </div>

            {/* Timing Layer */}
            <div className="flex justify-center">
              <div className="bg-pink-600 bg-opacity-30 border-2 border-pink-500 rounded-lg p-4 w-full max-w-md">
                <h3 className="text-white font-bold mb-2 text-center">TIMING & SYNCHRONIZATION</h3>
                <ul className="text-sm text-slate-300 space-y-1">
                  <li>• Master clock simulation</li>
                  <li>• Frame rate control (60 FPS)</li>
                  <li>• Component cycle coordination</li>
                </ul>
              </div>
            </div>

            <div className="flex justify-center">
              <ArrowRight className="text-slate-500 rotate-90" size={32} />
            </div>

            {/* Output */}
            <div className="flex justify-center gap-4">
              <div className="bg-slate-700 text-white px-6 py-3 rounded-lg font-semibold">
                Display Output
              </div>
              <div className="bg-slate-700 text-white px-6 py-3 rounded-lg font-semibold">
                Audio Output
              </div>
            </div>
          </div>
        </div>

        {/* Getting Started Tips */}
        <div className="mt-8 bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h2 className="text-xl font-semibold text-white mb-4">Getting Started Tips</h2>
          <div className="space-y-3 text-slate-300">
            <p><strong className="text-white">1. Choose your target system:</strong> Start simple (Chip-8, Game Boy) before complex systems (PS2, N64)</p>
            <p><strong className="text-white">2. Study the specs:</strong> Get CPU instruction set, memory map, timing diagrams</p>
            <p><strong className="text-white">3. Begin with CPU:</strong> Implement the fetch-decode-execute cycle first</p>
            <p><strong className="text-white">4. Test incrementally:</strong> Use test ROMs to verify each component</p>
            <p><strong className="text-white">5. Add components gradually:</strong> CPU → Memory → Graphics → Audio → Input</p>
          </div>
        </div>
      </div>
    </div>
  );
}