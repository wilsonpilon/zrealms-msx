# RELEASE.md — Detalhamento do Release Oficial

**Versão Atual:** 0.5.3  
**Data:** 04 de Outubro de 2026  
**Status do Release:** Phase 5 Concluída com Sucesso (Subfase 5.3 — Jogo de Referência Completo Concluída)  
**Alvo:** Windows (x64) para o Toolkit & GUI Desktop / MSX 2, MSX-DOS 2 (.COM / .DSK) e Cartucho MegaROM ASCII-16 (.ROM) para a Engine  

---

## 1. Resumo Executivo do Release v0.5.3

O release **v0.5.3** conclui com excelência a **Subfase 5.3 (Jogo de Referência Completo)** e coroa a entrega da **Fase 5 (Pipeline Integrado de Build & Jogo de Referência)** do Z-Realm.

Apresentamos **"As Catacumbas de Cristal: O Desafio do Rei Esquecido"**, um cRPG completo e funcional para o padrão MSX 2 demonstrando todas as capacidades tecnológicas desenvolvidas ao longo do projeto:
- Masmorra profunda com **20 salas interligadas em grid 4x5** com transições cardeais contínuas (Norte, Sul, Leste, Oeste).
- Driver nativo de som **PSG AY-3-8910** com efeitos sonoros em tempo real nos registradores `0xA0`/`0xA1`.
- Mecânicas de **combate em tempo real com IA hostil** (`BEHAVIOR_HOSTILE`), perseguição inteligente, dano de combate (-8 HP com SFX 2 e atualização do HUD) e golpe de espada do herói.
- Dano de perigo em tiles de espinho (`COLLISION_DAMAGE`).
- Enigmas em cadeia, NPCs com inteligência artificial errante e de patrulha, resgate de prisioneiro, decifração de profecia, restauração em fonte sagrada e conquista do **Cristal Primordial** de vitória.

Ambos os alvos de distribuição — **Disquete MSX-DOS 2 (`DOS2_zrealm.dsk`)** e **Cartucho MegaROM ASCII-16 (`zrealm_demo.rom`)** — estão gerados, validados e inclusos no pacote oficial de release.

---

## 2. Destaques Técnicos da Subfase 5.3

### 2.1 Masmorra Completa de 20 Salas em Grid 4x5 (`pkg/project/demo.go`)
- **Topologia Espacial 4x5:**
  - Linha 0: Armaria dos Antigos (4), Corredor das Sombras (3), Entrada das Catacumbas (1), Câmara dos Pilares (2), Salão dos Reis (5).
  - Linha 1: Galeria Subterrânea (8), Cripta dos Heróis (7), Vale das Almas (6), Labirinto de Pedra (9), Fosso de Espinhos (10).
  - Linha 2: Fonte Sagrada (12), Refúgio do Eremita (11), Celas Subterrâneas (13), Cárcere do Prisioneiro (14), Câmara de Tortura (15).
  - Linha 3: Esgotos da Cidadela (16), Pórtico Antigo (17), Caverna de Cristais (18), Antecâmara Real (19), Santuário do Rei Esquecido (20).
- **Corredores Desobstruídos:** Ajustes arquiteturais nos portais (Eixos X=15..16 e Y=8..9) permitindo navegação cardeal fluida sem colisões acidentais em paredes de salas vizinhas.

### 2.2 Catálogo Expandido de Recursos
- **13 Tiles Customizados:** Chão Limpo (Passable), Alvenaria (Solid), Portais de Arco (Trigger), Água e Canais (Water), Espinhos (Damage), Tochas Acesa (Solid), Estátuas Antigas (Solid), Altar Místico (Solid), Grades de Ferro (Solid), Lajotas Antigas (Passable), Cristais Azuis (Solid), Fonte de Água Benta (Solid) e Alavancas de Bronze (Solid).
- **8 Sprites Modo 2 (16x16 pixels):** Herói Guerreiro, Guardião Sentinela, Baú com Ferragens, Eremita Sábio, Esqueleto Guerreiro, Prisioneiro Trancafiado, Cristal Primordial (Ciano e Azul Místico) e Goblin Ladino.
- **7 Itens de RPG e Progressão:** Chave de Bronze, Poção de Vida (+25 HP), Chave de Ferro, Chave Real Dourada, Amuleto de Cristal, Elixir Mágico (+30 MP) e Cristal Primordial do Rei.
- **15 Strings de Diálogo:** Textos contextuais com moldura estética, quebra automática inteligente de linha e paginação dinâmica.
- **8 Bytecode Event Scripts:** Cadeia completa de quests e eventos gerenciados pela Bytecode VM.

### 2.3 Driver de Áudio PSG Nativo AY-3-8910 (`engine_msx/audio.h` e `audio.c`)
- Acesso de baixo nível aos registradores do PSG através das portas de I/O SDCC:
  - `__sfr __at(0xA0) g_PSG_RegPort;`
  - `__sfr __at(0xA1) g_PSG_DataPort;`
- **Banco de 4 Efeitos Sonoros:**
  - `SFX 1` (Item / Chime / Vitória): Arpejo brilhante ascendente com envelope suave (baús, itens recebidos, cura e triunfo).
  - `SFX 2` (Dano / Combate): Ruído percussivo de impacto com decaimento rápido (golpe recebido de monstro ou dano de espinho).
  - `SFX 3` (Porta / Mecanismo): Ruído mecânico de baixa frequência para destrancamento de selos e portas de pedra.
  - `SFX 4` (Diálogo / Blip): Pulso curto e suave emitido na abertura de janelas de texto e transição de páginas.
- Totalmente integrado na VM (`VM_PlaySFX`), na interface (`ui.c`) e nos loops principais da engine (`AUDIO_Update()`).

### 2.4 Combate em Tempo Real e Dano de Piso
- **IA Hostil (`BEHAVIOR_HOSTILE`, tipo 6):**
  - Detecção do herói dentro de raio de visão de 7 tiles.
  - Perseguição com aproximação Manhattan inteligente.
  - Ataque corpo a corpo automático desferindo -8 HP ao herói com SFX 2 e atualização instantânea do HUD.
- **Contra-Ataque do Herói:**
  - Pressionar a tecla de Ação (Espaço ou Gatilho do joystick) desfere golpe de espada no inimigo adjacente.
  - Abate imediato do monstro, ocultação do sprite de hardware no VDP, emissão do SFX 1 e execução do script de drop de quest (Chave Real Dourada).
- **Dano de Piso em Espinhos (`COLLISION_DAMAGE`):**
  - Penalidade de -5 HP ao pisar sobre armadilhas, com alerta sonoro e sincronização no HUD.

---

## 3. Validação Automatizada em Malha Fechada no openMSX

A integridade do jogo de referência e do pipeline de compilação foi verificada através da suite automatizada de testes `engine_msx/sub53_test.tcl`, executada no emulador openMSX sobre o cartucho MegaROM compilado (`zrealm_demo.rom`).

### Capturas de Tela Comprovatórias (SCREEN 4 - V9938):
1. **`sub53_01_spawn.png`**: Spawn inicial na Sala 1 (Entrada das Catacumbas) com tochas iluminando a masmorra, sprites do Herói e do Guardião ativos, portas abertas e HUD funcional (`♥ 075/100`, `★ 030/030`, `LV:01`, `🗝:0`).
2. **`sub53_02_dialogue.png`**: Abordagem ao Guardião Real em `(14, 7)` e interação via ESPAÇO; caixa de diálogo emoldurada nas linhas 19-23 com a fala do Guardião e concessão da Chave de Bronze (HUD atualizado para `🗝:1` e SFX 1 emitido).
3. **`sub53_03_armory.png`**: Travessia pelas Salas 1 e 3 até a Armaria dos Antigos (Sala 4); abordagem ao Baú de Tesouro em `(8, 5)` com estatua decorativa ao norte e sentinela patrulhando; abertura do baú recebendo Chave de Ferro e Poção de Vida com cura de +25 HP.
4. **`sub53_04_fountain.png`**: Descida pelas galerias até a Fonte Sagrada (Sala 12); tanque de águas cristalinas com fonte sagrada; ingestão das águas curativas restaurando a vida do herói para `♥ 100/100`.
5. **`sub53_05_combat.png`**: Jornada pelas celas até a Câmara de Tortura (Sala 15); confronto e combate em tempo real com o Esqueleto Guerreiro hostil; golpe de espada do herói destruindo a criatura e conquistando a Chave Real Dourada entre os ossos.
6. **`sub53_06_victory.png`**: Acesso ao Santuário do Rei Esquecido (Sala 20); santuário iluminado com cristais azuis cintilantes e altar central contendo o Cristal Primordial; heroica conquista do cristal com mensagem triunfal de vitória: `"VITORIA! VOCE ERGUEU O CRISTAL PRIMORDIAL E SALVOU O REINO DE Z-REALM!"`.

---

## 4. Pacote de Distribuição Oficial

O pacote oficial para download é gerado pelo script `build.ps1` e inclui:
- `dist/bin/zrealm.exe`: Utilitário executável de linha de comando para Windows x64 assinado digitalmente.
- `dist/msx/zrealm_demo.rom`: Cartucho MegaROM ASCII-16 completo (128 KB) pronto para gravação em FlashROM ou execução direta em emuladores.
- `dist/msx/DOS2_zrealm.dsk`: Imagem de disquete de 720 KB para MSX-DOS 2 com boot automático via `autoexec.bat`.
- `dist/msx/zrealm.com`, `HEADER.BIN`, `GAME.DAT`: Binários avulsos do jogo.
- `dist/tools/msxtar.exe`: Ferramenta para manipulação de imagens DSK.
- `dist/zrealm-msx-v0.5.3-windows-amd64.zip`: Arquivo compactado com todos os binários e documentação completa.
