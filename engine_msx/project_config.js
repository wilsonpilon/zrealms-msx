// ____________________________
// Z-Realm (zrealm-msx) - MSX 2 Engine Build Configuration
// Baseado em MSXgl (CC BY-SA) - Guillaume 'Aoineko' Blanchard
//─────────────────────────────────────────────────────────────────────────────

//*****************************************************************************
// BUILD STEPS
//*****************************************************************************
DoClean   = false;
DoCompile = true;
DoMake    = true;
DoPackage = true;
DoDeploy  = true;
DoRun     = false;

//*****************************************************************************
// TOOLS SETTINGS
//*****************************************************************************
SDCCPath  = "C:/Users/barney/scoop/apps/sdcc/current/";
Compiler  = `${SDCCPath}bin/sdcc.exe`;
Assembler = `${SDCCPath}bin/sdasz80.exe`;
Linker    = `${SDCCPath}bin/sdcc.exe`;
MakeLib   = `${SDCCPath}bin/sdar.exe`;
Hex2Bin   = `${RootDir}tools/MSXtk/bin/MSXhex.exe`;
MSXDOS    = `${RootDir}tools/build/DOS/`;
DskTool   = `${RootDir}tools/build/msxtar/msxtar.exe`;
Emulator  = "openmsx";

//*****************************************************************************
// PROJECT SETTINGS
//*****************************************************************************
ProjName = "zrealm";
ProjModules = [ ProjName, "mapper", "vdp_screen4", "loader", "world", "hero", "entity" ];
ProjSegments = ProjName;

// Módulos da biblioteca MSXgl necessários para a Engine
LibModules = [ "dos", "dos_mapper", "vdp", "system", "bios", "keyboard", "joystick" ];

//*****************************************************************************
// TARGET & HARDWARE SETTINGS
//*****************************************************************************
Machine = "2";               // MSX 2 (V9938)
Target  = "DOS2";            // MSX-DOS 2 (.COM com suporte a Memory Mapper nativo)
ROMSize = 256;               // Perfil padrão: 256 KB Memory Mapper

DiskFiles = [ "emul/dos2/HEADER.BIN", "emul/dos2/GAME.DAT" ];
DiskSize = "720K";

// Argumentos de linha de comando no DOS 2
DOSParseArg = true;
AppSignature = false;
