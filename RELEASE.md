# RELEASE.md — Detalhamento do Release Oficial

**Versão Atual:** 0.4.4  
**Data:** 03 de Outubro de 2026  
**Status do Release:** Phase 4 Complete (Gameplay Engine & Máquina de Eventos no MSX 100% Concluída)  
**Alvo:** Windows (x64) para o Toolkit & GUI / MSX 2 & MSX-DOS 2 para a Engine  

---

## 1. Resumo Executivo do Release

A **Fase 4 (Gameplay Engine & Máquina de Eventos no MSX)** está **100% concluída**. Este release entrega a engine completa em C compilada com SDCC e MSXgl para **MSX-DOS 2** e chips gráficos **Yamaha V9938**, trazendo:

1. **Movimentação no Grid & Física (Subfase 4.1):** Controle discreto de 8x8 pixels na viewport 32x18 tiles (SCREEN 4), colisão contra tiles sólidos e transição cardeal contínua entre salas usando paginação dinâmica no Memory Mapper (Página 2: `0x8000-0xBFFF`).
2. **Sistema de Entidades & IAs da Sala (Subfase 4.2):** Spawning de até 8 entidades simultâneas em Sprites Modo 2 (16x16 pixels com atributos de cor por scanline), NPCs estáticos, NPCs errantes com PRNG Z80, baús sólidos e gatilhos de piso acionados por aproximação ou tecla de ação (`ESPAÇO` / Botão A).
3. **Máquina Virtual de Eventos (Bytecode VM) (Subfase 4.3):** Interpretador de 11 opcodes canônicos (`OP_MSG`, `OP_GIVE_ITEM`, `OP_TAKE_ITEM`, `OP_SET_FLAG`, `OP_CHECK_FLAG`, `OP_TELEPORT`, `OP_HEAL`, `OP_DAMAGE`, etc.), inventário de 16 slots, 256 flags de evento globais e atributos do herói em buffer imune a chaveamentos de memória.
4. **Caixa de Diálogos & HUD no V9938 (Subfase 4.4):** HUD fixo em tempo real nas linhas 18-19 (`♥ HP: 075/100`, `★ MP: 030/030`, `LV: 01`, `🗝: 0..9`), caixa de diálogo emoldurada nas linhas 20-23 com *word-wrapping* automático de 30 colunas por linha, paginação com prompt `[ESPACO] ▼`, debounce seguro e painel de repouso ("Standby").

---

## 2. Componentes da Engine MSX 2 (C / MSXgl / SDCC)

* **`hero.h` & `hero.c`:** Controle do herói no grid 32x18, cooldown suave de passos, orientação direcional (Norte/Sul/Leste/Oeste) e integração de colisão física contra tiles e atores.
* **`world.h` & `world.c`:** Gerenciador do mundo ativo e paginação transparente de salas via Memory Mapper do MSX-DOS 2 (`MAPPER_SetPage2`), com suporte a limites cardeais e transições dinâmicas.
* **`entity.h` & `entity.c`:** Alocação de hardware para sprites 1 a 8 no V9938, rotinas de IA para NPCs errantes e interações contextuais com baús e NPCs.
* **`vm.h` & `vm.c`:** Interpretador da Bytecode VM com buffer isolado na Página 1 (`s_VMScriptBuffer[256]`), tabela de 256 flags globais e inventário de 16 posições.
* **`ui.h` & `ui.c`:** Gestão visual do Banco 2 da VRAM (linhas 16 a 23 da SCREEN 4), carga de fonte ASCII TMS9900 8x8 e glifos especiais de UI, renderizador de HUD em tempo real e caixa de diálogo emoldurada com paginação e word-wrapping.
* **`mapper.h` & `mapper.c`:** Gerenciamento do Memory Mapper via `EXTBIOS` do DOS 2, com alocação dinâmica e liberação integral garantida.
* **`vdp_screen4.h` & `vdp_screen4.c`:** Inicialização do modo Graphic 3 no V9938, bancos de padrões e cores, e viewport de 576 bytes contíguos.

---

## 3. Validação Automatizada no Emulador openMSX

Todas as 4 subfases foram submetidas a testes de malha fechada via scripts de automação TCL no openMSX (`Panasonic_FS-A1GT`, `ram512k`, `msxdos2`), capturando screenshots sequenciais e inspecionando a memória RAM:

* **Subfase 4.1 (`sub41_test.tcl`):** 5 capturas comprovando movimentação no grid, colisão física, transição entre Sala 1 e Sala 2 e encerramento limpo via `ESC`.
* **Subfase 4.2 (`sub42_test.tcl`):** 6 capturas comprovando spawn de entidades, colisão contra o Guardião, interação com baú, IA errante e persistência.
* **Subfase 4.3 (`sub43_test.tcl`):** 8 capturas comprovando quest do Guardião, obtenção de itens, branch condicional por Flag 1, cura ao abrir o baú (+25 HP) e detecção de baú esvaziado.
* **Subfase 4.4 (`sub44_test.tcl`):** 9 capturas comprovando HUD em tempo real, abertura de caixa de diálogo com moldura ciano e word-wrapping, atualização dinâmica do contador de chaves (`🗝:1`), diálogo ramificado, cura no HUD (`♥ 100/100`), diálogo de baú vazio e saída limpa ao prompt `A:\>`.

---

## 4. Próximos Passos (Fase 5: Pipeline Integrado de Build & Jogo de Referência)

Com a Gameplay Engine 100% concluída, o projeto ruma à sua fase final:
1. **Subfase 5.1 — Automação "One-Click Run":** Botão de compilação e teste automático disparando exportação de `.rpgproj`, montagem de `.DSK` e boot no openMSX.
2. **Subfase 5.2 — Backend MegaROM (Cartucho):** Suporte opcional à geração de arquivo `.ROM` unificado com chaveador de banco em cartucho (ASCII 16K / Konami).
3. **Subfase 5.3 — Jogo de Referência Completo:** Masmorra de demonstração com 20 salas, enigmas, NPCs, combate simples e trilha sonora PSG.
