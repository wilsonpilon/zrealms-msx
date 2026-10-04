# RELEASE.md — Detalhamento do Release Oficial

**Versão Atual:** 0.5.2  
**Data:** 04 de Outubro de 2026  
**Status do Release:** Phase 5 em Andamento (Subfase 5.2 — Backend MegaROM .ROM Concluída com Sucesso)  
**Alvo:** Windows (x64) para o Toolkit & GUI / MSX 2, MSX-DOS 2 e Cartucho MegaROM (ASCII-16) para a Engine  

---

## 1. Resumo Executivo do Release v0.5.2

O release **v0.5.2** conclui a **Subfase 5.2 (Backend MegaROM - Cartucho .ROM)**, expandindo o Z-Realm para suportar nativamente tanto o ambiente de disco **MSX-DOS 2 (`.COM`)** quanto cartuchos **MegaROM (`.ROM`)** no formato padrão da indústria **ASCII-16**.

Com esta entrega, os criadores de jogos no Z-Realm podem exportar e testar seus cRPGs em formato de cartucho com inicialização instantânea no MSX, sem depender de leitor de disquetes ou sistema operacional em disco.

### Principais Inovações da Subfase 5.2:

1. **Engine C Multi-Target (`engine_msx`):**
   - Suporte nativo e transparente para compilação multi-alvo:
     - `target=DOS2 make package`: gera `zrealm.com` e imagem `DOS2_zrealm.dsk` para MSX-DOS 2.
     - `target=ROM_ASCII16 make package`: gera `zrealm.rom` para cartuchos MegaROM ASCII-16.
   - Chaveamento de bancos MegaROM ASCII-16 implementado em `mapper.c`: escrita no endereço `0x77FF` para chavear o banco de 16 KB correspondente na Página 2 (`0x8000..0xBFFF`), mapeando `logicalSegment + 1`.
   - Carregador de ROM direto em `loader.c` (`LOADER_LoadROM`): mapeamento imediato da Tabela Mestra (`MasterHeader`) e Diretório de Recursos diretamente a partir do endereço de CPU `0x7C00` (offset `0x3C00` do arquivo ROM), sem necessidade de cópias intermediárias em RAM.
   - Ponto de entrada dedicado para cartucho em `zrealm.c` inicializando VDP em SCREEN 4, entidades, VM de eventos, herói e loop de 60 fps sincronizado via `halt` e V-Blank.
   - Ajustes de compatibilidade em `MSXgl` (`vdp.c`, `system.c`, `bios.c`, `compiler.js`) expandindo macros com quebras de linha e adicionando flag `-g` ao assembler `sdasz80` para tratamento automático de símbolos de runtime como globais.

2. **Exportador Go MegaROM (`pkg/exporter/megarom.go`):**
   - Embutimento da imagem base de cartucho de 16 KB compilada (`zrealm_base.rom`).
   - Injeção atômica da Tabela Mestra compilada (`HEADER.BIN`) na janela reservada de 1 KB em `0x3C00` (CPU `0x7C00`).
   - Concatenação sequencial dos segmentos de 16 KB empacotados (`GAME.DAT`) a partir do Banco 1 (`0x4000` em diante).
   - Preenchimento padronizado (*padding*) com `0xFF` através da função `CalculateStandardROMSize` para as capacidades comerciais clássicas de cartucho: 128 KB (8 bancos), 256 KB (16 bancos), 512 KB (32 bancos), 1024 KB (64 bancos), 2048 KB (128 bancos) e 4096 KB (256 bancos).
   - Suíte abrangente de testes unitários e de integração em `megarom_test.go`.

3. **Automação One-Click Run para Cartuchos (`pkg/runner`):**
   - Adição dos métodos `BuildOpenMSXROMArgs`, `LaunchOpenMSXCart` e `OneClickRunROM`.
   - Disparo do emulador openMSX com argumento `-cart <caminho.rom>`.

4. **Interface Gráfica Desktop Fyne (`pkg/gui`):**
   - Novo botão de destaque **"🕹️ Testar Cartucho MegaROM (.ROM)"** na aba Exportar.
   - Suporte à exportação manual de cartuchos `.ROM` com relatório visual detalhado de tamanho total, quantidade de bancos de 16 KB e contagem de recursos.

5. **Interface de Linha de Comando (CLI - `cmd/zrealm/main.go`):**
   - `-export-rom <arquivo.rpgproj>`: exporta o projeto diretamente para cartucho `.ROM`.
   - `-run-rom <arquivo.rpgproj>`: pipeline completo de exportação e inicialização instantânea no openMSX.
   - `-pad-size <KB>`: opção para forçar o tamanho exato do cartucho (ex.: 128, 256, 512).

---

## 2. Validação Automatizada de Malha Fechada no openMSX

A conformidade do cartucho MegaROM gerado (`demo.rom`) foi validada em malha fechada via script de controle TCL (`engine_msx/sub52_test.tcl`) no emulador openMSX (modelo `Philips_NMS_8250`), gerando 5 screenshots sequenciais em SCREEN 4:

1. `sub52_01_spawn_cart.png`: Inicialização do cartucho MegaROM na máquina virtual, ativação do modo SCREEN 4 (Graphic 3 do V9938), herói no ponto de partida `(16, 9)` da Sala 1 (Catacumbas), baú de tesouro em `(8, 4)`, NPC Guardião em `(12, 9)` e HUD funcional (`♥ 075/100`, `★ 030/030`, `LV:01`, `🗝:0`).
2. `sub52_02_facing_guardian_cart.png`: Navegação contínua no grid para oeste até a posição `(13, 9)`, ficando em frente ao Guardião.
3. `sub52_03_dialogue_cart.png`: Interação via barra de espaço acionando a Bytecode VM; caixa de diálogo emoldurada nas linhas 19-23 com o texto do Guardião: `"GUARDIAO: AS PROFUNDEZAS SAO PERIGOSAS. PEGUE A CHAVE DE FERRO! [ESPACO]"`.
4. `sub52_04_item_received_cart.png`: Avanço do diálogo, execução dos opcodes `GIVE_ITEM 1`, `SET_FLAG 1 1` e `PLAY_SFX 1`; restauração do painel de repouso e atualização dinâmica do HUD para `🗝:1`.
5. `sub52_05_room2_cart.png`: Travessia do portal Leste da Sala 1; acionamento do chaveamento de banco MegaROM via `0x77FF`; carga instantânea e renderização da Sala 2 (Câmara dos Pilares), com estrutura central em cruz, 4 pilares de pedra, baú místico, NPC errante e herói na entrada Oeste `(0, 9)` mantendo os dados de inventário.

---

## 3. Próximos Passos (Fase 5: Pipeline Integrado de Build & Jogo de Referência)

1. **Subfase 5.1 — Automação "One-Click Run" (Disquete DOS2):** [CONCLUÍDA ✅]
2. **Subfase 5.2 — Backend MegaROM (Cartucho .ROM ASCII-16):** [CONCLUÍDA ✅]
3. **Subfase 5.3 — Jogo de Referência Completo:** Masmorra de demonstração com 20 salas interligadas, quebra-cabeças de chaves e alavancas, múltiplos tipos de NPCs, combate simples em tempo real e trilha sonora em PSG.
