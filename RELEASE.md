# RELEASE.md — Detalhamento do Release Oficial

**Versão Atual:** 0.5.1  
**Data:** 04 de Outubro de 2026  
**Status do Release:** Phase 5 em Andamento (Subfase 5.1 — Automação "One-Click Run" Concluída)  
**Alvo:** Windows (x64) para o Toolkit & GUI / MSX 2 & MSX-DOS 2 para a Engine  

---

## 1. Resumo Executivo do Release v0.5.1

O release **v0.5.1** implementa a **Subfase 5.1 (Automação "One-Click Run")**, unificando o ciclo de desenvolvimento de cRPGs para MSX 2. Agora, com um único clique no editor ou um comando no terminal, qualquer projeto `.rpgproj` é exportado para binários nativos (`HEADER.BIN`, `GAME.DAT`), empacotado em um disquete virtual MSX-DOS 2 de 720 KB (`.dsk`) e inicializado automaticamente no emulador openMSX.

### Principais Inovações da Subfase 5.1:
1. **Pacote `pkg/runner` Autônomo:**
   - Recursos de sistema embutidos diretamente no binário via `//go:embed` (`COMMAND2.COM`, `MSXDOS2.SYS`, `autoexec.bat`, `zrealm.com`). O toolkit agora gera disquetes bootáveis em qualquer máquina sem dependências externas de arquivos DOS.
   - Integração com `msxtar` para criação ultrarrápida (<10ms) de imagens FAT12 de dupla face / 80 trilhas (737.280 bytes) compatíveis com MSX-DOS 2.
   - Detecção inteligente do executável do openMSX via `PATH`, Scoop shims (`%USERPROFILE%\scoop\shims\openmsx.exe`) e diretórios padrão de instalação do Windows (`Program Files`).
2. **Interface Gráfica (Fyne Desktop):**
   - Nova seção na aba **Exportar**: *🚀 Automação 'One-Click Run' (openMSX)* com atalho de teclado `[F5]` e botão proeminente.
   - Item no menu principal: *Ferramentas -> 🚀 Executar no openMSX (F5)*.
   - Seletor de perfil de máquina (`Philips_NMS_8250` [MSX2 clássico europeu] e `Panasonic_FS-A1GT` [MSXturboR]).
   - Campo para override de caminho customizado do emulador e console de logs em tempo real na própria GUI.
3. **Interface de Linha de Comando (CLI):**
   - Flag `-run <arquivo.rpgproj>` com suporte a `-machine <perfil>`, `-emu <caminho>` e `-script <script.tcl>` para automação e pipelines CI/CD de testes.
4. **Validação de Malha Fechada no openMSX (`sub51_test.tcl`):**
   - Teste automatizado de inicialização, navegação pelo grid e saída limpa via `ESC` gerando screenshots comprobatórias em `engine_msx/screenshots/`.

---

## 2. Histórico da Engine MSX 2 (C / MSXgl / SDCC) — Fase 4

A engine nativa do MSX 2 segue consolidada e validada:
* **Movimentação no Grid & Física (Subfase 4.1):** Controle discreto de 8x8 pixels na viewport 32x18 tiles (SCREEN 4), colisão contra tiles sólidos e transição cardeal contínua entre salas usando paginação dinâmica no Memory Mapper (Página 2: `0x8000-0xBFFF`).
* **Sistema de Entidades & IAs da Sala (Subfase 4.2):** Spawning de até 8 entidades simultâneas em Sprites Modo 2 (16x16 pixels com atributos de cor por scanline), NPCs estáticos, NPCs errantes com PRNG Z80, baús sólidos e gatilhos de piso.
* **Máquina Virtual de Eventos (Bytecode VM) (Subfase 4.3):** Interpretador de 11 opcodes canônicos (`OP_MSG`, `OP_GIVE_ITEM`, `OP_TAKE_ITEM`, `OP_SET_FLAG`, `OP_CHECK_FLAG`, etc.), inventário de 16 slots, 256 flags globais e atributos do herói.
* **Caixa de Diálogos & HUD no V9938 (Subfase 4.4):** HUD fixo nas linhas 18-19, caixa de diálogo emoldurada nas linhas 20-23 com *word-wrapping* automático de 30 colunas, paginação com prompt `[ESPACO] ▼` e debounce seguro.

---

## 3. Validação Automatizada de Malha Fechada

* **Subfase 5.1 (`sub51_test.tcl`):**
  - `sub51_01_oneclick_boot.png`: Inicialização do `.dsk` gerado via `msxtar`, transição de SCREEN 0 para SCREEN 4, carregamento de tileset, sala e HUD.
  - `sub51_02_oneclick_gameplay.png`: Movimentação do herói para oeste, colisão física e acionamento de diálogo com o Guardião.
  - `sub51_03_oneclick_exit.png`: Tecla `ESC` pressionada, finalização ordenada da engine e retorno seguro ao prompt `A:\>` com memória desalocada.

---

## 4. Próximos Passos (Fase 5: Pipeline Integrado de Build & Jogo de Referência)

1. **Subfase 5.1 — Automação "One-Click Run":** [CONCLUÍDA ✅]
2. **Subfase 5.2 — Backend MegaROM (Cartucho):** Suporte opcional à geração de arquivo `.ROM` unificado com chaveador de banco em cartucho (ASCII 16K / Konami).
3. **Subfase 5.3 — Jogo de Referência Completo:** Masmorra de demonstração com 20 salas, enigmas, NPCs, combate simples e trilha sonora PSG.
